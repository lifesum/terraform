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
	pipelineStages = expandPipelineStages(d)
	pipelineArtificatStore = expandPipelineArtifactStore(d)

	pipeline := &codepipeline.PipelineDeclaration{
		Name:    aws.String(d.Get("name").(string)),
		RoleArn: aws.String(d.Get("role_arn").(string)),
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

func expandPipelineArtifactStore(d *schema.ResourceData) []codepipeline.ArtifactStore {
	pipelineArtificatStore := []codepipeline.ArtifactStore{}
	return pipelineArtificatStore
}

func expandPipelineStages(d *schema.ResourceData) []codepipeline.StageDeclaration {
	pipelineArtificatStore := []codepipeline.StageDeclaration{{}}
	return pipelineArtificatStore
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
