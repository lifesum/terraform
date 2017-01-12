package aws

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsCodeBuildProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsCodeBuildProjectCreate,
		Read:   resourceAwsCodeBuildProjectRead,
		Update: resourceAwsCodeBuildProjectUpdate,
		Delete: resourceAwsCodeBuildProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"artifacts": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"location": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"namespace_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAwsCodeBuildArifactsNamespaceType,
						},
						"packaging": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAwsCodeBuildArifactsType,
						},
					},
				},
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAwsCodeBuildProjectDescription,
			},
			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"environment": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compute_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAwsCodeBuildEnvironmentComputeType,
						},
						"environment_variable": &schema.Schema{
							Type: schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
							Optional: true,
						},
						"image": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAwsCodeBuildEnvironmentType,
						},
					},
				},
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAwsCodeBuildProjectName,
			},
			"service_role": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth": &schema.Schema{
							Type: schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateAwsCodeBuildSourceAuthType,
									},
								},
							},
							Optional: true,
						},
						"buildspec": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"location": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAwsCodeBuildSourceType,
						},
					},
				},
				Required: true,
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAwsCodeBuildTimeout,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAwsCodeBuildProjectCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).codebuildconn

	log.Printf("[DEBUG] CodeBuild Create Project: %s", d.Id())

	params := &codebuild.CreateProjectInput{
		Environment: expandProjectEnvironment(d.Get("environment").(*schema.Set).List()[0].(map[string]interface{})),
		Name:        aws.String(d.Get("name").(string)),
		Source:      expandProjectSource(d.Get("source").(*schema.Set).List()[0].(map[string]interface{})),
		Artifacts:   expandProjectArtifacts(d.Get("artifacts").(*schema.Set).List()[0].(map[string]interface{})),
	}

	if v, ok := d.GetOk("description"); ok {
		params.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("encryption_key"); ok {
		params.EncryptionKey = aws.String(v.(string))
	}

	if v, ok := d.GetOk("service_role"); ok {
		params.ServiceRole = aws.String(v.(string))
	}

	if v, ok := d.GetOk("timeout"); ok {
		params.TimeoutInMinutes = aws.Int64(int64(v.(int)))
	}

	var resp *codebuild.CreateProjectOutput
	err := resource.Retry(30*time.Second, func() *resource.RetryError {
		var err error

		resp, err = conn.CreateProject(params)

		if err != nil {
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("[ERROR] Error creating CodeBuild project: %s", err)
	}

	if resp.Project.Name == nil {
		return fmt.Errorf("[ERROR] Project name was nil")
	}

	d.SetId(*resp.Project.Name)

	return resourceAwsCodeBuildProjectUpdate(d, meta)
}

func expandProjectArtifacts(m map[string]interface{}) *codebuild.ProjectArtifacts {

	projectArtifacts := &codebuild.ProjectArtifacts{
		Type: aws.String(m["type"].(string)),
	}

	if len(m["location"].(string)) > 0 {
		projectArtifacts.Location = aws.String(m["location"].(string))
	}

	if len(m["name"].(string)) > 0 {
		projectArtifacts.Name = aws.String(m["name"].(string))
	}

	if len(m["namespace_type"].(string)) > 0 {
		projectArtifacts.NamespaceType = aws.String(m["namespace_type"].(string))
	}

	if len(m["packaging"].(string)) > 0 {
		projectArtifacts.Packaging = aws.String(m["packaging"].(string))
	}

	if len(m["path"].(string)) > 0 {
		projectArtifacts.Path = aws.String(m["path"].(string))
	}

	return projectArtifacts
}

func expandProjectEnvironment(m map[string]interface{}) *codebuild.ProjectEnvironment {
	projectEnv := &codebuild.ProjectEnvironment{
		ComputeType: aws.String(m["compute_type"].(string)),
		Image:       aws.String(m["image"].(string)),
		Type:        aws.String(m["type"].(string)),
	}

	envVariables := m["environment_variable"].(*schema.Set).List()
	projectEnv.EnvironmentVariables = make([]*codebuild.EnvironmentVariable, len(envVariables))

	for i := 0; i < len(envVariables); i++ {
		v := envVariables[i].(map[string]interface{})
		projectEnv.EnvironmentVariables[i] = &codebuild.EnvironmentVariable{
			Name:  aws.String(v["name"].(string)),
			Value: aws.String(v["value"].(string)),
		}
	}

	return projectEnv
}

func expandProjectSource(m map[string]interface{}) *codebuild.ProjectSource {

	projectSource := &codebuild.ProjectSource{
		Type:      aws.String(m["type"].(string)),
		Location:  aws.String(m["location"].(string)),
		Buildspec: aws.String(m["buildspec"].(string)),
	}

	if v, ok := m["auth"]; ok {
		if len(v.(*schema.Set).List()) > 0 {
			auth := v.(*schema.Set).List()[0].(map[string]interface{})

			projectSource.Auth = &codebuild.SourceAuth{
				Type:     aws.String(auth["type"].(string)),
				Resource: aws.String(auth["resource"].(string)),
			}
		}
	}

	return projectSource
}

func resourceAwsCodeBuildProjectRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).codebuildconn

	log.Printf("[DEBUG] CodeBuild Read Project: %s", d.Id())

	resp, err := conn.BatchGetProjects(&codebuild.BatchGetProjectsInput{
		Names: []*string{
			aws.String(d.Get("name").(string)),
		},
	})

	if err != nil {
		return fmt.Errorf("[ERROR] Error retreiving Projects: %q", err)
	}

	// if nothing was found, then return no state
	if len(resp.Projects) == 0 {
		log.Printf("[INFO]: No projects were found, removing from state")
		d.SetId("")
		return nil
	}

	project := resp.Projects[0]

	artifacts := flatternAwsCodebuildProjectArtifacts(project.Artifacts)
	if artifacts != nil {
		d.Set("artifacts", artifacts)
	}

	environment := flattenAwsCodebuildProjectEnvironment(project.Environment)
	if environment != nil {
		d.Set("environment", environment)
	}

	if err := d.Set("source", flattenAwsCodebuildProjectSource(project.Source)); err != nil {
		return err
	}

	d.Set("description", project.Description)
	d.Set("encryption_key", project.EncryptionKey)
	d.Set("name", project.Name)
	d.Set("service_role", project.ServiceRole)
	d.Set("timeout", project.TimeoutInMinutes)

	if err := d.Set("tags", tagsToMapCodeBuild(project.Tags)); err != nil {
		return err
	}

	return nil
}

func resourceAwsCodeBuildProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).codebuildconn

	log.Printf("[DEBUG] CodeBuild Update Project: %s", d.Id())

	params := &codebuild.UpdateProjectInput{
		Environment: expandProjectEnvironment(d.Get("environment").(*schema.Set).List()[0].(map[string]interface{})),
		Name:        aws.String(d.Get("name").(string)),
		Source:      expandProjectSource(d.Get("source").(*schema.Set).List()[0].(map[string]interface{})),
		Artifacts:   expandProjectArtifacts(d.Get("artifacts").(*schema.Set).List()[0].(map[string]interface{})),
	}

	if v, ok := d.GetOk("description"); ok {
		params.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("encryption_key"); ok {
		params.EncryptionKey = aws.String(v.(string))
	}

	if v, ok := d.GetOk("service_role"); ok {
		params.ServiceRole = aws.String(v.(string))
	}

	if v, ok := d.GetOk("timeout"); ok {
		params.TimeoutInMinutes = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("tags"); ok {
		params.Tags = tagsFromMapCodeBuild(v.(map[string]interface{}))
	}

	_, err := conn.UpdateProject(params)

	if err != nil {
		return fmt.Errorf(
			"[ERROR] Error updating CodeBuild project (%s): %s",
			d.Id(), err)
	}

	return resourceAwsCodeBuildProjectRead(d, meta)
}

func resourceAwsCodeBuildProjectDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).codebuildconn

	_, err := conn.DeleteProject(&codebuild.DeleteProjectInput{
		Name: aws.String(d.Id()),
	})

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func flatternAwsCodebuildProjectArtifacts(artifacts *codebuild.ProjectArtifacts) *schema.Set {

	artifactSet := schema.Set{
		F: resourceAwsCodeBuildProjectArtifactsHash,
	}

	values := map[string]interface{}{}

	values["type"] = *artifacts.Type

	if artifacts.Location != nil {
		values["location"] = *artifacts.Location
	}

	if artifacts.Name != nil {
		values["name"] = *artifacts.Name
	}

	if artifacts.NamespaceType != nil {
		values["namespace_type"] = *artifacts.NamespaceType
	}

	if artifacts.Packaging != nil {
		values["packaging"] = *artifacts.Packaging
	}

	if artifacts.Path != nil {
		values["path"] = *artifacts.Path
	}

	artifactSet.Add(values)

	return &artifactSet
}

func flattenAwsCodebuildProjectEnvironment(environment *codebuild.ProjectEnvironment) *schema.Set {

	environmentSet := schema.Set{
		F: resourceAwsCodeBuildProjectEnvironmentHash,
	}

	envConfig := map[string]interface{}{}

	envConfig["type"] = *environment.Type
	envConfig["compute_type"] = *environment.ComputeType
	envConfig["image"] = *environment.Image
	envConfig["environment_variable"] = environmentVariablesToMap(environment.EnvironmentVariables)

	environmentSet.Add(envConfig)

	return &environmentSet

}

func flattenAwsCodebuildProjectSource(source *codebuild.ProjectSource) *schema.Set {

	sourceSet := schema.Set{
		F: resourceAwsCodeBuildProjectSourceHash,
	}

	sourceConfig := map[string]interface{}{}

	sourceConfig["type"] = *source.Type

	if source.Auth != nil {
		sourceConfig["auth"] = sourceAuthToMap(source.Auth)
	}

	if source.Buildspec != nil {
		sourceConfig["buildspec"] = *source.Buildspec
	}

	if source.Location != nil {
		sourceConfig["location"] = *source.Location
	}

	sourceSet.Add(sourceConfig)

	return &sourceSet

}

func resourceAwsCodeBuildProjectArtifactsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	artifactType := m["type"].(string)

	buf.WriteString(fmt.Sprintf("%s-", artifactType))

	return hashcode.String(buf.String())
}

func resourceAwsCodeBuildProjectEnvironmentHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	environmentType := m["type"].(string)
	computeType := m["compute_type"].(string)
	image := m["image"].(string)

	buf.WriteString(fmt.Sprintf("%s-", environmentType))
	buf.WriteString(fmt.Sprintf("%s-", computeType))
	buf.WriteString(fmt.Sprintf("%s-", image))

	return hashcode.String(buf.String())
}

func resourceAwsCodeBuildProjectSourceHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	sourceType := m["type"].(string)
	buildspec := m["buildspec"].(string)
	location := m["location"].(string)

	buf.WriteString(fmt.Sprintf("%s-", sourceType))
	buf.WriteString(fmt.Sprintf("%s-", buildspec))
	buf.WriteString(fmt.Sprintf("%s-", location))

	return hashcode.String(buf.String())
}

func environmentVariablesToMap(environmentVariables []*codebuild.EnvironmentVariable) []map[string]interface{} {

	envVariables := make([]map[string]interface{}, len(environmentVariables))

	if len(environmentVariables) > 0 {
		for i := 0; i < len(environmentVariables); i++ {
			env := environmentVariables[i]
			item := make(map[string]interface{})
			item["name"] = *env.Name
			item["value"] = *env.Value
			envVariables = append(envVariables, item)
		}
	}

	return envVariables
}

func sourceAuthToMap(sourceAuth *codebuild.SourceAuth) map[string]interface{} {

	auth := map[string]interface{}{}
	auth["type"] = *sourceAuth.Type

	if sourceAuth.Type != nil {
		auth["resource"] = *sourceAuth.Resource
	}

	return auth
}

func validateAwsCodeBuildArifactsType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	types := map[string]bool{
		"CODEPIPELINE": true,
		"NO_ARTIFACTS": true,
		"S3":           true,
	}

	if !types[value] {
		errors = append(errors, fmt.Errorf("CodeBuild: Arifacts Type can only be CODEPIPELINE / NO_ARTIFACTS / S3"))
	}
	return
}

func validateAwsCodeBuildArifactsNamespaceType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	types := map[string]bool{
		"NONE":     true,
		"BUILD_ID": true,
	}

	if !types[value] {
		errors = append(errors, fmt.Errorf("CodeBuild: Arifacts Namespace Type can only be NONE / BUILD_ID"))
	}
	return
}

func validateAwsCodeBuildProjectName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[A-Za-z0-9]`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"first character of %q must be a letter or number", value))
	}

	if !regexp.MustCompile(`^[A-Za-z0-9\-_]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only alphanumeric characters, hyphens and underscores allowed in %q", value))
	}

	if len(value) > 255 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be greater than 255 characters", value))
	}

	return
}

func validateAwsCodeBuildProjectDescription(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 255 {
		errors = append(errors, fmt.Errorf("%q cannot be greater than 255 characters", value))
	}
	return
}

func validateAwsCodeBuildEnvironmentComputeType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	types := map[string]bool{
		"BUILD_GENERAL1_SMALL":  true,
		"BUILD_GENERAL1_MEDIUM": true,
		"BUILD_GENERAL1_LARGE":  true,
	}

	if !types[value] {
		errors = append(errors, fmt.Errorf("CodeBuild: Environment Compute Type can only be BUILD_GENERAL1_SMALL / BUILD_GENERAL1_MEDIUM / BUILD_GENERAL1_LARGE"))
	}
	return
}

func validateAwsCodeBuildEnvironmentType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	types := map[string]bool{
		"LINUX_CONTAINER": true,
	}

	if !types[value] {
		errors = append(errors, fmt.Errorf("CodeBuild: Environment Type can only be LINUX_CONTAINER"))
	}
	return
}

func validateAwsCodeBuildSourceType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	types := map[string]bool{
		"CODECOMMIT":   true,
		"CODEPIPELINE": true,
		"GITHUB":       true,
		"S3":           true,
	}

	if !types[value] {
		errors = append(errors, fmt.Errorf("CodeBuild: Source Type can only be CODECOMMIT / CODEPIPELINE / GITHUB / S3"))
	}
	return
}

func validateAwsCodeBuildSourceAuthType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	types := map[string]bool{
		"OAUTH": true,
	}

	if !types[value] {
		errors = append(errors, fmt.Errorf("CodeBuild: Source Auth Type can only be OAUTH"))
	}
	return
}

func validateAwsCodeBuildTimeout(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)

	if value < 5 || value > 480 {
		errors = append(errors, fmt.Errorf("%q must be greater than 5 minutes and less than 480 minutes (8 hours)", value))
	}
	return
}
