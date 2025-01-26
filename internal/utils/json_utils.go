package utils

// TraverseDynamicJSON processes a JSON object of type `interface{}` to extract field paths, treating arrays as a single field.
func TraverseDynamicJSON(data interface{}) []string {
	var fields []string

	var traverse func(interface{}, string)
	traverse = func(data interface{}, prefix string) {
		switch v := data.(type) {
		case map[string]interface{}:
			// If the value is a map, iterate through its keys
			for key, value := range v {
				// Build the field path
				fullPath := key
				if prefix != "" {
					fullPath = prefix + "." + key
				}
				// Recurse into nested structures
				traverse(value, fullPath)
			}
		case map[string]string:
			// Handle map[string]string
			for key := range v {
				fullPath := key
				if prefix != "" {
					fullPath = prefix + "." + key
				}
				// Since map[string]string has leaf nodes, append directly
				fields = append(fields, fullPath)
			}
		case []interface{}:
			// If the value is a slice, process each element
			for _, value := range v {
				traverse(value, prefix) // Keep the same prefix for arrays
			}
		default:
			// Base case: Add the full path to the fields list
			fields = append(fields, prefix)
		}
	}

	traverse(data, "")
	return fields
}
