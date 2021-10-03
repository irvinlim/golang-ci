package stringutils

import (
	"strings"
)

// MakeLines joins multiple strings with a newline character.
func MakeLines(lines ...string) string {
	return strings.Join(lines, "\n")
}
