package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAppEngine_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppEngineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAppEngine,
				Check: resource.ComposeTestCheckFunc(
					testAccAppEngineExists(
						"google_appengine_version.foobar"),
				),
			},
		},
	})
}

func testAccCheckAppEngineDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_appengine_version" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		ver, _ := config.clientAppEngine.Apps.Services.Versions.Get(rs.Primary.Attributes["apps_id"], rs.Primary.Attributes["services_id"], rs.Primary.Attributes["versions_id"]).Do()
		if ver != nil {
			return fmt.Errorf("Version still present")
		}
	}

	return nil
}

func testAccAppEngineExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)

		_, err := config.clientAppEngine.Apps.Services.Versions.Get(rs.Primary.Attributes["apps_id"], rs.Primary.Attributes["services_id"], rs.Primary.Attributes["versions_id"]).Do()
		if err != nil {
			return fmt.Errorf("BigQuery Dataset not present")
		}

		return nil
	}
}

var testAccAppEngine = `
resource "google_appengine_version" "foobar" {
  apps_id = "foo"
}
`
