package utils

import (
	"strings"

	"github.com/javiorfo/nilo"
)

func GetValueAsCommaSeparated(value any) nilo.Option[[]string] {
	str, ok := value.(string)
	if !ok || !strings.Contains(str, ",") {
		return nilo.None[[]string]()
	}
	return nilo.Some(strings.Split(str, ","))
}
