package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSCodePipeline_basic(t *testing.T) {
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSCodePipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSCodePipelineConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSCodePipelineExists("aws_codepipeline.foo"),
				),
			},
		},
	})
}

func testAccCheckAWSCodePipelineExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CodePipeline ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).codepipelineconn

		_, err := conn.GetPipeline(&codepipeline.GetPipelineInput{
			Name: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckAWSCodePipelineDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).codepipelineconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_codepipeline" {
			continue
		}

		_, err := conn.GetPipeline(&codepipeline.GetPipelineInput{
			Name: aws.String(rs.Primary.ID),
		})

		if err == nil {
			return fmt.Errorf("Expected AWS CodePipeline to be gone, but was still found")
		}
		return nil
	}

	return fmt.Errorf("Default error in CodePipeline Test")
}

// JSON from AWS Console created pipeline
// {
//     "pipeline": {
//         "roleArn": "arn:aws:iam::xxx:role/AWS-CodePipeline-Service",
//         "stages": [
//             {
//                 "name": "Source",
//                 "actions": [
//                     {
//                         "inputArtifacts": [],
//                         "name": "Source",
//                         "actionTypeId": {
//                             "category": "Source",
//                             "owner": "ThirdParty",
//                             "version": "1",
//                             "provider": "GitHub"
//                         },
//                         "outputArtifacts": [
//                             {
//                                 "name": "MyApp"
//                             }
//                         ],
//                         "configuration": {
//                             "Owner": "lifesum-terraform",
//                             "Repo": "test",
//                             "Branch": "master",
//                             "OAuthToken": "****"
//                         },
//                         "runOrder": 1
//                     }
//                 ]
//             },
//             {
//                 "name": "Build",
//                 "actions": [
//                     {
//                         "inputArtifacts": [
//                             {
//                                 "name": "MyApp"
//                             }
//                         ],
//                         "name": "CodeBuild",
//                         "actionTypeId": {
//                             "category": "Build",
//                             "owner": "AWS",
//                             "version": "1",
//                             "provider": "CodeBuild"
//                         },
//                         "outputArtifacts": [
//                             {
//                                 "name": "MyAppBuild"
//                             }
//                         ],
//                         "configuration": {
//                             "ProjectName": "test"
//                         },
//                         "runOrder": 1
//                     }
//                 ]
//             }
//         ],
//         "artifactStore": {
//             "type": "S3",
//             "location": "codepipeline-us-west-2-679037673204"
//         },
//         "name": "test",
//         "version": 1
//     }
// }

func testAccAWSCodePipelineConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "main" {
    bucket = "tf-test-pipeline-%s"
    acl = "private"
}


resource "aws_iam_role" "codepipeline_role" {
  name = "codepipeline-role-%s"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "codepipeline.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}


resource "aws_codepipeline" "foo" {
  name         = "test-pipeline-%s"
  role_arn = "${aws_iam_role.codepipeline_role.arn}"

  artifact_store {
  		location = "${aws_s3_bucket.main.bucket}"
  		type = "S3"
  }

  stage {
	name = "Source"
	action {
			name = "Source"
			category = "Source"
			owner = "ThirdParty"
			provider = "GitHub"
			version = "1"
			output_artifacts = ["test"]
			configuration = <<EOF
{
    "Owner": "lifesum-terraform",
    "Repo": "test",
    "Branch": "master",
    "OAuthToken": "0000000000000000000000000000000000000000"
}
EOF
		}
  } 

  stage {
	name = "Build"
	action {
			name = "Build"
			category = "Build"
			owner = "AWS"
			provider = "CodeBuild"
			input_artifacts = ["test"]
			version = "1"
			configuration = <<EOF
{
    "ProjectName": "test"
}
EOF
		}
  }
}
`, rName, rName, rName)
}
