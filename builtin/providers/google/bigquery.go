package google

func expandBigQueryLabels(configured interface{}) map[string]string {
	labels := map[string]string{}

	for k, v := range configured.(map[string]interface{}) {
		labels[k] = v.(string)
	}

	return labels
}
