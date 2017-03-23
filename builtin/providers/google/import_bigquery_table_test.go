package google

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBigQueryTable_importBasic(t *testing.T) {
	resourceName := "google_bigquery_table.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccBigQueryTable,
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
