package aws

import (
	"fmt"
	"time"

	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsCodePipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsCodePipelineCreate,
		Read:   resourceAwsCodePipelineRead,
		Update: resourceAwsCodePipelineUpdate,
		Delete: resourceAwsCodePipelineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"role_arn": {
				Type:     schema.TypeString,
				Required: true,
			},

			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"artifact_store": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Type:     schema.TypeString,
							Required: true,
						},

						"type": {
							Type:     schema.TypeString,
							Required: true,
						},

						"encryption_key": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Required: true,
									},

									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"stage": {
				Type:     schema.TypeList,
				MinItems: 2,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"action": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"configuration": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"category": {
										Type:     schema.TypeString,
										Required: true,
									},
									"owner": {
										Type:     schema.TypeString,
										Required: true,
									},
									"provider": {
										Type:     schema.TypeString,
										Required: true,
									},
									"version": {
										Type:     schema.TypeString,
										Required: true,
									},
									"input_artifacts": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"output_artifacts": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAwsCodePipelineCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).codepipelineconn
	pipelineStages := expandPipelineStages(d)
	pipelineArtifactStore := expandPipelineArtifactStore(d)

	pipeline := &codepipeline.PipelineDeclaration{
		Name:          aws.String(d.Get("name").(string)),
		RoleArn:       aws.String(d.Get("role_arn").(string)),
		ArtifactStore: pipelineArtifactStore,
		Stages:        pipelineStages,
	}
	params := &codepipeline.CreatePipelineInput{
		Pipeline: pipeline,
	}

	var resp *codepipeline.CreatePipelineOutput
	err := resource.Retry(30*time.Second, func() *resource.RetryError {
		var err error

		resp, err = conn.CreatePipeline(params)

		if err != nil {
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating CodePipeline: %s", err)
	}
	fmt.Printf("%#v", resp)
	d.SetId(*resp.Pipeline.Name)

	return resourceAwsCodePipelineUpdate(d, meta)
}

func expandPipelineArtifactStore(d *schema.ResourceData) *codepipeline.ArtifactStore {
	configs := d.Get("artifact_store").([]interface{})
	data := configs[0].(map[string]interface{})
	pipelineArtifactStore := codepipeline.ArtifactStore{
		Location: aws.String(data["location"].(string)),
		Type:     aws.String(data["type"].(string)),
	}
	return &pipelineArtifactStore
}

func expandPipelineStages(d *schema.ResourceData) []*codepipeline.StageDeclaration {
	configs := d.Get("stage").([]interface{})
	pipelineStages := []*codepipeline.StageDeclaration{}

	for _, stage := range configs {
		data := stage.(map[string]interface{})
		a := data["action"].([]interface{})
		actions := expandPipelineActions(a)
		pipelineStages = append(pipelineStages, &codepipeline.StageDeclaration{
			Name:    aws.String(data["name"].(string)),
			Actions: actions,
		})
	}
	fmt.Printf("%#v", pipelineStages)
	return pipelineStages
}

func expandPipelineActions(s []interface{}) []*codepipeline.ActionDeclaration {
	actions := []*codepipeline.ActionDeclaration{}
	for _, taction := range s {
		action := taction.(map[string]interface{})
		conf := map[string]*string{}
		if action["configuration"].(string) != "" {
			json.Unmarshal([]byte(action["configuration"].(string)), &conf)
		}
		oa := action["output_artifacts"].([]interface{})
		outputArtifacts := expandPipelineActionsOutputArtifacts(oa)

		ia := action["input_artifacts"].([]interface{})
		inputArtifacts := expandPipelineActionsInputArtifacts(ia)

		actions = append(actions, &codepipeline.ActionDeclaration{
			ActionTypeId: &codepipeline.ActionTypeId{
				Category: aws.String(action["category"].(string)),
				Owner:    aws.String(action["owner"].(string)),

				Provider: aws.String(action["provider"].(string)),
				Version:  aws.String(action["version"].(string)),
			},
			Name:            aws.String(action["name"].(string)),
			Configuration:   conf,
			OutputArtifacts: outputArtifacts,
			InputArtifacts:  inputArtifacts,
		})
	}
	return actions
}
func expandPipelineActionsOutputArtifacts(s []interface{}) []*codepipeline.OutputArtifact {
	outputArtifacts := []*codepipeline.OutputArtifact{}
	for _, artifact := range s {
		outputArtifacts = append(outputArtifacts, &codepipeline.OutputArtifact{
			Name: aws.String(artifact.(string)),
		})
	}
	return outputArtifacts
}
func expandPipelineActionsInputArtifacts(s []interface{}) []*codepipeline.InputArtifact {
	outputArtifacts := []*codepipeline.InputArtifact{}
	for _, artifact := range s {
		outputArtifacts = append(outputArtifacts, &codepipeline.InputArtifact{
			Name: aws.String(artifact.(string)),
		})
	}
	return outputArtifacts
}

func resourceAwsCodePipelineRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsCodePipelineUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsCodePipelineDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
