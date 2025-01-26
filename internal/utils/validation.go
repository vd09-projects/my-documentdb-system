package utils

// IsValidRecord checks if the data meets minimal JSON criteria using generics
func IsValidRecord[T any](rec map[string]T) bool {
	// Example: Must have a "userId" field
	_, ok := rec["userId"]
	return ok
}
