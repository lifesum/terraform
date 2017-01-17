package aws

import "github.com/hashicorp/terraform/helper/schema"

func resourceAwsCodePipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsCodepipelineCreate,
		Read:   resourceAwsCodepipelineRead,
		Update: resourceAwsCodepipelineUpdate,
		Delete: resourceAwsCodepipelineDelete,
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

						"stages": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"actions": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action_type_id": {
													Type: schema.TypeString,
												},
											},
										},
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

func resourceAwsCodepipelineCreate(d *schema.ResourceData, meta interface{}) error {
	codepipelineconn := meta.(*AWSClient).codepipelineconn

	return resourceAwsCodepipelineUpdate(d, meta)
}

func resourceAwsCodepipelineRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsCodepipelineUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsCodepipelineDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
