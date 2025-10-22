package utils

import (
	"strings"

	"github.com/javiorfo/nilo"
)

// GetValueAsCommaSeparated attempts to convert a value to a slice of strings
// by splitting it on commas if the value is a comma-separated string.
// Returns an Option containing the slice if successful, or None otherwise.
func GetValueAsCommaSeparated(value any) nilo.Option[[]string] {
	str, ok := value.(string)
	if !ok || !strings.Contains(str, ",") {
		return nilo.None[[]string]()
	}
	return nilo.Some(strings.Split(str, ","))
}
