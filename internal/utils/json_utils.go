package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

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

// extractFloatValues converts a slice of interfaces into a slice of float64 values,
// ignoring non-numeric elements.
func extractFloatValues(data []interface{}) []float64 {
	var floats []float64
	for _, val := range data {
		switch v := val.(type) {
		case float64:
			floats = append(floats, v)
		case int:
			// Handle integers as floats
			floats = append(floats, float64(v))
		case string:
			fl, err := strconv.ParseFloat(v, 64)
			if err == nil {
				floats = append(floats, fl)
			}
		}
	}
	return floats
}

// GetNestedValue marshals bson.M to JSON, unmarshals it into a Go map, and then performs traversal.
func GetNestedValue(data bson.M, path string) ([]float64, bool) {
	// Step 1: Marshal bson.M to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling BSON to JSON:", err)
		return nil, false
	}

	// Step 2: Unmarshal JSON into a generic Go map
	var genericData map[string]interface{}
	err = json.Unmarshal(jsonData, &genericData)
	if err != nil {
		fmt.Println("Error unmarshaling JSON to map:", err)
		return nil, false
	}

	// Step 3: Traverse the unmarshaled map
	return traverseMap(genericData, path)
}

// traverseMap performs recursive traversal on a generic map based on a dot-separated path.
func traverseMap(data interface{}, path string) ([]float64, bool) {
	keys := strings.Split(path, ".")

	var traverse func(current interface{}, depth int) ([]interface{}, bool)
	traverse = func(current interface{}, depth int) ([]interface{}, bool) {
		// Base case: if we've reached the end of the path
		if depth > len(keys) {
			return nil, false
		}

		switch v := current.(type) {
		case map[string]interface{}:
			// Traverse into the map using the key
			nextVal, ok := v[keys[depth]]
			if !ok {
				return nil, false
			}
			return traverse(nextVal, depth+1)

		case []interface{}:
			// Traverse into each element of the array
			var collected []interface{}
			for _, elem := range v {
				result, ok := traverse(elem, depth)
				if ok {
					collected = append(collected, result...)
				}
			}
			if len(collected) == 0 {
				return nil, false
			}
			return collected, true

		default:
			// Unsupported type for traversal
			return []interface{}{current}, true
		}
	}

	// Start traversal
	collected, ok := traverse(data, 0)
	if !ok {
		return nil, false
	}

	// Extract float64 values from the collected data
	return extractFloatValues(collected), true
}
