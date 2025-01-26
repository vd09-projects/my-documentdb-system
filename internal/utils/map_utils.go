package utils

func GetKeysFromMap(inputMap map[string][]interface{}) []string {
	keys := make([]string, 0, len(inputMap)) // Preallocate slice with map size
	for key := range inputMap {
		keys = append(keys, key)
	}
	return keys
}
