package aws

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAWSCodeBuildProject_importBasic(t *testing.T) {
	resourceName := "aws_codebuild_project.foo"
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSCodeBuildProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSCodeBuildProjectConfig_basic(name),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
