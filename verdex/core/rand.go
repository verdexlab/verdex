package core

import "math/rand"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

// Generate a random string (a-z) of given length
func RandomAlphaString(length int) string {
	b := make([]rune, length)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
