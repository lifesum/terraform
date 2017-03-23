package google

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBigQueryDataset_importBasic(t *testing.T) {
	resourceName := "google_bigquery_dataset.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryDatasetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccBigQueryDataset,
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"retain_on_delete"},
			},
		},
	})
}
