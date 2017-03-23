package google

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/hashicorp/terraform/helper/validation"
	bigquery "google.golang.org/api/bigquery/v2"
	"google.golang.org/api/googleapi"
)

func resourceBigQueryTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigQueryTableCreate,
		Read:   resourceBigQueryTableRead,
		Delete: resourceBigQueryTableDelete,
		Update: resourceBigQueryTableUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// https://godoc.org/google.golang.org/api/bigquery/v2#TableReference
			"table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"dataset_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// Description: [Optional] The field description. The maximum length is
			// 16K characters.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// ExpirationTime: [Optional] The time when this table expires, in
			// milliseconds since the epoch. If not present, the table will persist
			// indefinitely. Expired tables will be deleted and their storage
			// reclaimed.
			"expiration_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			// FriendlyName: [Optional] A descriptive name for this table.
			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Labels: [Experimental] The labels associated with this table. You can
			// use these to organize and group your tables. Label keys and values
			// can be no longer than 63 characters, can only contain lowercase
			// letters, numeric characters, underscores and dashes. International
			// characters are allowed. Label values are optional. Label keys must
			// start with a letter and each label in the list must have a different
			// key.
			"labels": &schema.Schema{
				Elem:     schema.TypeString,
				Type:     schema.TypeMap,
				Optional: true,
			},

			// Schema: [Optional] Describes the schema of this table.
			"schema": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.ValidateJsonString,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
			},

			// TimePartitioning: [Experimental] If specified, configures time-based
			// partitioning for this table.
			"time_partitioning": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expiration_ms": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"DAY"}, false),
						},
					},
				},
			},

			// Don't delete the table
			"retain_on_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			// CreationTime: [Output-only] The time when this table was created, in
			// milliseconds since the epoch.
			"creation_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// Etag: [Output-only] A hash of this resource.
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Kind: [Output-only] The type of the resource.
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// LastModifiedTime: [Output-only] The time when this table was last
			// modified, in milliseconds since the epoch.
			"last_modified_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// Location: [Output-only] The geographic location where the table
			// resides. This value is inherited from the dataset.
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// NumBytes: [Output-only] The size of this table in bytes, excluding
			// any data in the streaming buffer.
			"num_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// NumLongTermBytes: [Output-only] The number of bytes in the table that
			// are considered "long-term storage".
			"num_long_term_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// NumRows: [Output-only] The number of rows of data in this table,
			// excluding any data in the streaming buffer.
			"num_rows": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// SelfLink: [Output-only] A URL that can be used to access this
			// resource again.
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Type: [Output-only] Describes the table type. The following values
			// are supported: TABLE: A normal BigQuery table. VIEW: A virtual table
			// defined by a SQL query. EXTERNAL: A table that references data stored
			// in an external storage system, such as Google Cloud Storage. The
			// default value is TABLE.
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTable(d *schema.ResourceData, meta interface{}) (*bigquery.Table, error) {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	table := &bigquery.Table{TableReference: &bigquery.TableReference{
		DatasetId: d.Get("dataset_id").(string),
		TableId:   d.Get("table_id").(string),
		ProjectId: project,
	}}

	if v, ok := d.GetOk("description"); ok {
		table.Description = v.(string)
	}

	if v, ok := d.GetOk("expiration_time"); ok {
		table.ExpirationTime = v.(int64)
	}

	if v, ok := d.GetOk("friendly_name"); ok {
		table.FriendlyName = v.(string)
	}

	if v, ok := d.GetOk("labels"); ok {
		table.Labels = expandBigQueryLabels(v)
	}

	if v, ok := d.GetOk("schema"); ok {
		schema, err := expandSchema(v.(string))
		if err != nil {
			return nil, err
		}

		table.Schema = schema
	}

	if v, ok := d.GetOk("time_partitioning"); ok {
		table.TimePartitioning = expandTimePartitioning(v)
	}

	return table, nil
}

func resourceBigQueryTableCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	table, err := resourceTable(d, meta)
	if err != nil {
		return err
	}

	dID := d.Get("dataset_id").(string)

	res, err := config.clientBigQuery.Tables.Insert(project, dID, table).Do()
	if err != nil {
		return err
	}

	d.SetId(res.Id)

	return resourceBigQueryTableRead(d, meta)
}

func resourceBigQueryTableParseID(id string) (string, string, string) {
	// projectID, datasetID, tableID separately
	parts := strings.FieldsFunc(id, func(r rune) bool { return r == ':' || r == '.' })
	return parts[0], parts[1], parts[2]
}

func resourceBigQueryTableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	pID, dID, tID := resourceBigQueryTableParseID(d.Id())

	res, err := config.clientBigQuery.Tables.Get(pID, dID, tID).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			// The resource doesn't exist anymore
			d.SetId("")

			return nil
		}

		return err
	}

	d.Set("description", res.Description)
	d.Set("expiration_time", res.ExpirationTime)
	d.Set("friendly_name", res.FriendlyName)
	d.Set("labels", res.Labels)
	d.Set("creation_time", res.CreationTime)
	d.Set("etag", res.Etag)
	d.Set("kind", res.Kind)
	d.Set("last_modified_time", res.LastModifiedTime)
	d.Set("location", res.Location)
	d.Set("num_bytes", res.NumBytes)
	d.Set("table_id", res.TableReference.TableId)
	d.Set("dataset_id", res.TableReference.DatasetId)
	d.Set("num_long_term_bytes", res.NumLongTermBytes)
	d.Set("num_rows", res.NumRows)
	d.Set("self_link", res.SelfLink)
	d.Set("type", res.Type)

	if res.TimePartitioning != nil {
		if err := d.Set("time_partitioning", flattenTimePartitioning(res.TimePartitioning)); err != nil {
			return err
		}
	}

	if res.Schema != nil {
		schema, err := flattenSchema(res.Schema)
		if err != nil {
			return err
		}

		d.Set("schema", schema)
	}

	return nil
}

func resourceBigQueryTableUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	table, err := resourceTable(d, meta)
	if err != nil {
		return err
	}

	pID, dID, tID := resourceBigQueryTableParseID(d.Id())

	_, err = config.clientBigQuery.Tables.Update(pID, dID, tID, table).Do()
	if err != nil {
		return err
	}

	return resourceBigQueryTableRead(d, meta)
}

func resourceBigQueryTableDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if d.Get("retain_on_delete").(bool) {
		log.Printf("[WARN] Removing BigQuery table ID %q with `retain_on_delete` set. Please delete this table manually.", d.Id())
		d.SetId("")
		return nil
	}

	pID, dID, tID := resourceBigQueryTableParseID(d.Id())

	err := config.clientBigQuery.Tables.Delete(pID, dID, tID).Do()
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func expandField(field map[string]interface{}) (*bigquery.TableFieldSchema, error) {
	fieldSchema := &bigquery.TableFieldSchema{}

	if v, ok := field["description"]; ok {
		fieldSchema.Description = v.(string)
	}

	if v, ok := field["mode"]; ok {
		fieldSchema.Mode = v.(string)
	}

	if v, ok := field["name"]; ok {
		fieldSchema.Name = v.(string)
	}

	if v, ok := field["type"]; ok {
		fieldSchema.Type = v.(string)
	}

	if v, ok := field["fields"]; ok {
		fieldList, err := expandFieldList(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		fieldSchema.Fields = fieldList
	}

	return fieldSchema, nil
}

func expandFieldList(schema []interface{}) ([]*bigquery.TableFieldSchema, error) {
	var tableFieldList []*bigquery.TableFieldSchema
	for _, fieldInterface := range schema {
		fieldParsed, err := expandField(fieldInterface.(map[string]interface{}))
		if err != nil {
			return nil, err
		}

		tableFieldList = append(tableFieldList, fieldParsed)
	}

	return tableFieldList, nil
}

func flattenTimePartitioning(tp *bigquery.TimePartitioning) []map[string]interface{} {
	result := []map[string]interface{}{}

	result = append(result, map[string]interface{}{"type": tp.Type})

	if tp.ExpirationMs != 0 {
		result[0]["expiration_ms"] = tp.ExpirationMs
	}

	return result
}

func expandTimePartitioning(configured interface{}) *bigquery.TimePartitioning {
	tpRaw := configured.([]interface{})[0].(map[string]interface{})
	tp := &bigquery.TimePartitioning{
		Type: tpRaw["type"].(string),
	}

	if eMs, ok := tpRaw["expiration_ms"]; ok {
		tp.ExpirationMs = int64(eMs.(int))
	}

	return tp
}

func flattenSchema(tableSchema *bigquery.TableSchema) (string, error) {
	schema, err := json.Marshal(tableSchema.Fields)
	if err != nil {
		return "", err
	}

	return string(schema), nil
}

func expandSchema(rawSchema string) (*bigquery.TableSchema, error) {
	var schemaJSONInterface []interface{}
	if err := json.Unmarshal([]byte(rawSchema), &schemaJSONInterface); err != nil {
		return nil, err
	}

	fieldList, err := expandFieldList(schemaJSONInterface)
	if err != nil {
		return nil, err
	}
	return &bigquery.TableSchema{Fields: fieldList}, nil
}
