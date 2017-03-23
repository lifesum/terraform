package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAppEngineApplication_Basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppEngineApplicationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAppEngineApplication,
				Check: resource.ComposeTestCheckFunc(
					testAccAppEngineApplicationExists(
						"google_appengine_application.foobar"),
				),
			},
		},
	})
}

func testAccCheckAppEngineApplicationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_appengine_application" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		app, _ := config.clientAppEngine.Apps.Get(rs.Primary.ID).Do()
		if app != nil {
			return fmt.Errorf("Application still present")
		}
	}

	return nil
}

func testAccAppEngineApplicationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		_, err := config.clientAppEngine.Apps.Get(rs.Primary.ID).Do()
		if err != nil {
			return fmt.Errorf("Application does not exist")
		}

		return nil
	}
}

var testAccAppEngineApplication = `
resource "google_appengine_application" "foobar" {}
`
