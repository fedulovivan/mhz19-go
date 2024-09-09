package arguments

import "strings"

func isTemplate(in string) bool {
	return strings.Contains(in, "{{") && strings.Contains(in, "}}")
}
