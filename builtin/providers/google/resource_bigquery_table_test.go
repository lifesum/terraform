package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBigQueryTable_Basic(t *testing.T) {
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccBigQueryTable(datasetID, tableID),
				Check: resource.ComposeTestCheckFunc(
					testAccBigQueryTableExists(
						"google_bigquery_table.test"),
				),
			},

			resource.TestStep{
				Config: testAccBigQueryTableUpdated(datasetID, tableID),
				Check: resource.ComposeTestCheckFunc(
					testAccBigQueryTableExists(
						"google_bigquery_table.test"),
				),
			},
		},
	})
}

func testAccCheckBigQueryTableDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_bigquery_table" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		_, err := config.clientBigQuery.Tables.Get(
			config.Project,
			rs.Primary.Attributes["dataset_id"],
			rs.Primary.Attributes["name"],
		).Do()
		if err == nil {
			return fmt.Errorf("Table still present")
		}
	}

	return nil
}

func testAccBigQueryTableExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		_, err := config.clientBigQuery.Tables.Get(
			config.Project,
			rs.Primary.Attributes["dataset_id"],
			rs.Primary.Attributes["name"],
		).Do()
		if err != nil {
			return fmt.Errorf("BigQuery Table not present")
		}

		return nil
	}
}

func testAccBigQueryTable(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
	dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
	table_id = "%s"
	dataset_id = "${google_bigquery_dataset.test.dataset_id}"

	time_partitioning {
	  type = "DAY"
	}

	schema = <<EOH
[
  {
    "name": "event",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "date",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "timestamp",
    "type": "TIMESTAMP",
    "mode": "NULLABLE"
  },
  {
    "name": "data",
    "type": "RECORD",
    "mode": "NULLABLE",
    "fields": [
      {
        "name": "action",
        "type": "RECORD",
        "mode": "NULLABLE",
        "fields": [
          {
            "name": "name",
            "type": "STRING",
            "mode": "NULLABLE"
          },
          {
            "name": "result",
            "type": "STRING",
            "mode": "NULLABLE"
          }
        ]
      },
      {
        "name": "test_name",
        "type": "STRING",
        "mode": "NULLABLE"
      }
    ]
  }
]
EOH
}`, datasetID, tableID)
}

func testAccBigQueryTableUpdated(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
	dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
	table_id = "%s"
	dataset_id = "${google_bigquery_dataset.test.dataset_id}"

	time_partitioning {
	  type = "DAY"
	}

	schema = <<EOH
[
  {
    "name": "event",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "date",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "timestamp",
    "type": "TIMESTAMP",
    "mode": "NULLABLE"
  },
  {
    "name": "data",
    "type": "RECORD",
    "mode": "NULLABLE",
    "fields": [
      {
        "name": "action",
        "type": "RECORD",
        "mode": "NULLABLE",
        "fields": [
          {
            "name": "name",
            "type": "STRING",
            "mode": "NULLABLE"
          },
          {
            "name": "result",
            "type": "STRING",
            "mode": "NULLABLE"
          }
        ]
      },
      {
        "name": "test_name",
        "type": "STRING",
        "mode": "NULLABLE"
      }
    ]
  },
  {
    "name": "screen",
    "type": "RECORD",
    "mode": "NULLABLE",
    "fields": [
      {
        "name": "height",
        "type": "INTEGER",
        "mode": "NULLABLE"
      },
      {
        "name": "width",
        "type": "INTEGER",
        "mode": "NULLABLE"
      }
    ]
  }
]
EOH
}`, datasetID, tableID)
}
