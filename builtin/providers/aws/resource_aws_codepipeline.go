package aws

import (
	"fmt"
	"time"

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
	var pipelineStages []*codepipeline.StageDeclaration

	for _, stage := range configs {
		data := stage.(map[string]interface{})
		actionData := data["action"].([]interface{})
		var actions []*codepipeline.ActionDeclaration
		for _, taction := range actionData {
			action := taction.(map[string]interface{})
			actions = append(actions, &codepipeline.ActionDeclaration{
				ActionTypeId: &codepipeline.ActionTypeId{
					Category: aws.String(action["category"].(string)),
					Owner:    aws.String(action["owner"].(string)),
					Provider: aws.String(action["provider"].(string)),
					Version:  aws.String(action["version"].(string)),
				},
				Name: aws.String(action["name"].(string)),
			})
		}
		pipelineStages = append(pipelineStages, &codepipeline.StageDeclaration{
			Name:    aws.String(data["name"].(string)),
			Actions: actions,
		})
	}
	return pipelineStages
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
