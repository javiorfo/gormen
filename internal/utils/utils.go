package utils

import "strings"

// GetValueAsCommaSeparated attempts to convert a value to a slice of strings
// by splitting it on commas if the value is a comma-separated string.
// Returns a slice if successful, or an empty slice otherwise.
func GetValueAsCommaSeparated(value any) []string {
	str, ok := value.(string)
	if !ok || !strings.Contains(str, ",") {
		return nil
	}
	return strings.Split(str, ",")
}
