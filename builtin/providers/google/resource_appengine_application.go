package google

import (
	"github.com/hashicorp/terraform/helper/schema"
	appengine "google.golang.org/api/appengine/v1"
)

func resourceAppEngineApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppEngineCreate,
		Read:   resourceAppEngineRead,
		Delete: resourceAppEngineDelete,

		Schema: map[string]*schema.Schema{
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"location_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAppEngineCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	call := config.clientAppEngine.Apps.Create(&appengine.Application{
		Id:         project,
		LocationId: d.Get("location_id").(string),
	})

	res, err := call.Do()
	if err != nil {
		return err
	}

	d.SetId(res.Name)

	return resourceAppEngineRead(d, meta)
}

func resourceAppEngineRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	call := config.clientAppEngine.Apps.Get(d.Id())
	res, err := call.Do()
	if err != nil {
		return err
	}

	d.Set("name", res.Name)
	d.Set("auth_domain", res.AuthDomain)
	d.Set("code_bucket", res.CodeBucket)

	return nil
}

func resourceAppEngineDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
