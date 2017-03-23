package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBigQueryDataset_Basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryDatasetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccBigQueryDataset,
				Check: resource.ComposeTestCheckFunc(
					testAccBigQueryDatasetExists(
						"google_bigquery_dataset.foobar"),
				),
			},

			resource.TestStep{
				Config: testAccBigQueryDatasetUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccBigQueryDatasetExists(
						"google_bigquery_dataset.foobar"),
				),
			},
		},
	})
}

func testAccCheckBigQueryDatasetDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_bigquery_dataset" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		ds, _ := config.clientBigQuery.Datasets.Get(config.Project, rs.Primary.Attributes["dataset_id"]).Do()
		if ds != nil {
			return fmt.Errorf("Dataset still present")
		}
	}

	return nil
}

func testAccCheckBigQueryDatasetExistsThenDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_bigquery_dataset" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		ds, _ := config.clientBigQuery.Datasets.Get(config.Project, rs.Primary.Attributes["dataset_id"]).Do()
		if ds == nil {
			return fmt.Errorf("Dataset was deleted when it shouldn't have been!")
		}

		err := config.clientBigQuery.Datasets.Delete(config.Project, rs.Primary.Attributes["dataset_id"]).Do()
		if err != nil {
			return fmt.Errorf("Failed to hard delete soft delete target after check.")
		}
	}
	return nil
}

func testAccBigQueryDatasetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		_, err := config.clientBigQuery.Datasets.Get(config.Project, rs.Primary.Attributes["dataset_id"]).Do()
		if err != nil {
			return fmt.Errorf("BigQuery Dataset not present")
		}

		return nil
	}
}

var testAccBigQueryDataset = fmt.Sprintf(`
resource "google_bigquery_dataset" "foobar" {
	dataset_id = "foo_%[1]v"
	friendly_name = "foo"
	description = "This is a foo description"
	location = "EU"
	default_table_expiration_ms = 3600000

	labels {
		env = "foo"
		default_table_expiration_ms = 3600000
	}
}`, acctest.RandString(10))

var testAccBigQueryDatasetUpdate = fmt.Sprintf(`
resource "google_bigquery_dataset" "foobar" {
	dataset_id = "bar__%[1]v"
	friendly_name = "bar"
	description = "This is a bar description"
	location = "EU"
	default_table_expiration_ms = 7200000

	labels {
		env = "bar"
		default_table_expiration_ms = 3600000
	}
}`, acctest.RandString(10))
