package aws

import (
	"fmt"
	"testing"

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

func TestAccAWSCodePipeline_artifactsTypeValidation(t *testing.T) {
	// cases := []struct {
	// 	Value    string
	// 	ErrCount int
	// }{
	// 	{Value: "CODEPIPELINE", ErrCount: 0},
	// 	{Value: "NO_ARTIFACTS", ErrCount: 0},
	// 	{Value: "S3", ErrCount: 0},
	// 	{Value: "XYZ", ErrCount: 1},
	// }

	// for _, tc := range cases {
	// 	_, errors := validateAwsCodeBuildArifactsType(tc.Value, "aws_codebuild_project")

	// 	if len(errors) != tc.ErrCount {
	// 		t.Fatalf("Expected the AWS CodeBuild project artifacts type to trigger a validation error")
	// 	}
	// }
}

func TestAccAWSCodePipeline_artifactsNamespaceTypeValidation(t *testing.T) {
	// cases := []struct {
	// 	Value    string
	// 	ErrCount int
	// }{
	// 	{Value: "NONE", ErrCount: 0},
	// 	{Value: "BUILD_ID", ErrCount: 0},
	// 	{Value: "XYZ", ErrCount: 1},
	// }

	// for _, tc := range cases {
	// 	_, errors := validateAwsCodeBuildArifactsNamespaceType(tc.Value, "aws_codebuild_project")

	// 	if len(errors) != tc.ErrCount {
	// 		t.Fatalf("Expected the AWS CodeBuild project artifacts namepsace_type to trigger a validation error")
	// 	}
	// }
}

// func longTestData() string {
// 	data := `
// 	test-test-test-test-test-test-test-test-test-test-
// 	test-test-test-test-test-test-test-test-test-test-
// 	test-test-test-test-test-test-test-test-test-test-
// 	test-test-test-test-test-test-test-test-test-test-
// 	test-test-test-test-test-test-test-test-test-test-
// 	test-test-test-test-test-test-test-test-test-test-
// 	`

// 	return strings.Map(func(r rune) rune {
// 		if unicode.IsSpace(r) {
// 			return -1
// 		}
// 		return r
// 	}, data)
// }

func TestAccAWSCodePipeline_nameValidation(t *testing.T) {
	// cases := []struct {
	// 	Value    string
	// 	ErrCount int
	// }{
	// 	{Value: "_test", ErrCount: 1},
	// 	{Value: "test", ErrCount: 0},
	// 	{Value: "1_test", ErrCount: 0},
	// 	{Value: "test**1", ErrCount: 1},
	// 	{Value: longTestData(), ErrCount: 1},
	// }

	// for _, tc := range cases {
	// 	_, errors := validateAwsCodePipelineName(tc.Value, "aws_codebuild_project")

	// 	if len(errors) != tc.ErrCount {
	// 		t.Fatalf("Expected the AWS CodeBuild project name to trigger a validation error - %s", errors)
	// 	}
	// }
}

func TestAccAWSCodePipeline_descriptionValidation(t *testing.T) {
	// cases := []struct {
	// 	Value    string
	// 	ErrCount int
	// }{
	// 	{Value: "test", ErrCount: 0},
	// 	{Value: longTestData(), ErrCount: 1},
	// }

	// for _, tc := range cases {
	// 	_, errors := validateAwsCodePipelineDescription(tc.Value, "aws_codebuild_project")

	// 	if len(errors) != tc.ErrCount {
	// 		t.Fatalf("Expected the AWS CodeBuild project description to trigger a validation error")
	// 	}
	// }
}

func TestAccAWSCodePipeline_environmentComputeTypeValidation(t *testing.T) {
	// cases := []struct {
	// 	Value    string
	// 	ErrCount int
	// }{
	// 	{Value: "BUILD_GENERAL1_SMALL", ErrCount: 0},
	// 	{Value: "BUILD_GENERAL1_MEDIUM", ErrCount: 0},
	// 	{Value: "BUILD_GENERAL1_LARGE", ErrCount: 0},
	// 	{Value: "BUILD_GENERAL1_VERYLARGE", ErrCount: 1},
	// }

	// for _, tc := range cases {
	// 	_, errors := validateAwsCodeBuildEnvironmentComputeType(tc.Value, "aws_codebuild_project")

	// 	if len(errors) != tc.ErrCount {
	// 		t.Fatalf("Expected the AWS CodeBuild project environment compute_type to trigger a validation error")
	// 	}
	// }
}

func TestAccAWSCodePipeline_environmentTypeValidation(t *testing.T) {
	// cases := []struct {
	// 	Value    string
	// 	ErrCount int
	// }{
	// 	{Value: "LINUX_CONTAINER", ErrCount: 0},
	// 	{Value: "WINDOWS_CONTAINER", ErrCount: 1},
	// }

	// for _, tc := range cases {
	// 	_, errors := validateAwsCodeBuildEnvironmentType(tc.Value, "aws_codebuild_project")

	// 	if len(errors) != tc.ErrCount {
	// 		t.Fatalf("Expected the AWS CodeBuild project environment type to trigger a validation error")
	// 	}
	// }
}

func TestAccAWSCodePipeline_sourceTypeValidation(t *testing.T) {
	// cases := []struct {
	// 	Value    string
	// 	ErrCount int
	// }{
	// 	{Value: "CODECOMMIT", ErrCount: 0},
	// 	{Value: "CODEPIPELINE", ErrCount: 0},
	// 	{Value: "GITHUB", ErrCount: 0},
	// 	{Value: "S3", ErrCount: 0},
	// 	{Value: "GITLAB", ErrCount: 1},
	// }

	// for _, tc := range cases {
	// 	_, errors := validateAwsCodeBuildSourceType(tc.Value, "aws_codebuild_project")

	// 	if len(errors) != tc.ErrCount {
	// 		t.Fatalf("Expected the AWS CodeBuild project source type to trigger a validation error")
	// 	}
	// }
}

func TestAccAWSCodePipeline_sourceAuthTypeValidation(t *testing.T) {
	// cases := []struct {
	// 	Value    string
	// 	ErrCount int
	// }{
	// 	{Value: "OAUTH", ErrCount: 0},
	// 	{Value: "PASSWORD", ErrCount: 1},
	// }

	// for _, tc := range cases {
	// 	_, errors := validateAwsCodeBuildSourceAuthType(tc.Value, "aws_codebuild_project")

	// 	if len(errors) != tc.ErrCount {
	// 		t.Fatalf("Expected the AWS CodeBuild project source auth to trigger a validation error")
	// 	}
	// }
}

func TestAccAWSCodePipeline_timeoutValidation(t *testing.T) {
	// cases := []struct {
	// 	Value    int
	// 	ErrCount int
	// }{
	// 	{Value: 10, ErrCount: 0},
	// 	{Value: 200, ErrCount: 0},
	// 	{Value: 1, ErrCount: 1},
	// 	{Value: 500, ErrCount: 1},
	// }

	// for _, tc := range cases {
	// 	_, errors := validateAwsCodeBuildTimeout(tc.Value, "aws_codebuild_project")

	// 	if len(errors) != tc.ErrCount {
	// 		t.Fatalf("Expected the AWS CodeBuild project timeout to trigger a validation error")
	// 	}
	// }
}

func testAccCheckAWSCodePipelineExists(n string) resource.TestCheckFunc {
	// 	return func(s *terraform.State) error {
	// 		rs, ok := s.RootModule().Resources[n]
	// 		if !ok {
	// 			return fmt.Errorf("Not found: %s", n)
	// 		}

	// 		if rs.Primary.ID == "" {
	// 			return fmt.Errorf("No CodeBuild Project ID is set")
	// 		}

	// 		conn := testAccProvider.Meta().(*AWSClient).codebuildconn

	// 		out, err := conn.BatchGetProjects(&codebuild.BatchGetProjectsInput{
	// 			Names: []*string{
	// 				aws.String(rs.Primary.ID),
	// 			},
	// 		})

	// 		if err != nil {
	// 			return err
	// 		}

	// 		if len(out.Projects) < 1 {
	// 			return fmt.Errorf("No project found")
	// 		}

	return nil
	// 	}
}

func testAccCheckAWSCodePipelineDestroy(s *terraform.State) error {
	// conn := testAccProvider.Meta().(*AWSClient).codepipelineconn

	// for _, rs := range s.RootModule().Resources {
	// 	if rs.Type != "aws_codepipeline" {
	// 		continue
	// 	}

	// 	out, err := conn.BatchGetProjects(&codebuild.BatchGetProjectsInput{
	// 		Names: []*string{
	// 			aws.String(rs.Primary.ID),
	// 		},
	// 	})

	// 	if err != nil {
	// 		return err
	// 	}

	// 	if out != nil && len(out.Projects) > 0 {
	// 		return fmt.Errorf("Expected AWS CodeBuild Project to be gone, but was still found")
	// 	}

	return nil
	// }

	// return fmt.Errorf("Default error in CodeBuild Test")
}

func testAccAWSCodePipelineConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "codepipeline_role" {
  name = "codebuild-role-%s"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "codebuild.amazonaws.com"
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
  		location = "test"
  		type = "S3"
  }

  stage {
	name = "first_stage"
	action {
			name = "notsure"
			category = "Deploy"
			owner = "AWS"
			provider = "CodeDeploy"
			version = "foo"
		}
  }
}
`, rName, rName)
}
