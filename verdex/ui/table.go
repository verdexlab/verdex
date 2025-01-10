package ui

import (
	"regexp"
	"unicode"
)

func tableWidthFunc(input string) int {
	// Remove non-ascii chars
	input = regexp.
		MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`).
		ReplaceAllString(input, "")

	// Count printable chars
	count := 0
	for _, r := range input {
		if unicode.IsGraphic(r) {
			count++
		}
	}

	return count
}
