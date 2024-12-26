package assets

import (
	"regexp"
	"strings"
)

// Check if word is present in asset body
func (asset *Asset) BodyMatchWord(word string) bool {
	return strings.Contains(asset.Body, word)
}

// Check if asset body match given regex pattern
func (asset *Asset) BodyMatchRegex(pattern string) bool {
	match, err := regexp.MatchString(pattern, asset.Body)
	if err != nil {
		return false
	}

	return match
}
