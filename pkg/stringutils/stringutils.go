package stringutils

import (
	"bytes"
	"math/rand"
	"strings"
)

// GenerateRandomAlphaOnlyString generates an alphabetical random string with length n.
func GenerateRandomAlphaOnlyString(n int) string {
	// make a really long string
	letters := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// GenerateRandomASCIIString generates an ASCII random string with length n.
func GenerateRandomASCIIString(n int) string {
	chars := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:` "
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = chars[rand.Intn(len(chars))]
	}

	return string(res)
}

// Ellipsis truncates a string to fit within maxlen, and appends ellipsis (...).
// For maxlen of 3 and lower, no ellipsis is appended.
func Ellipsis(s string, maxlen int) string {
	r := []rune(s)
	if len(r) <= maxlen {
		return s
	}

	if maxlen <= 3 {
		return string(r[:maxlen])
	}

	return string(r[:maxlen-3]) + "..."
}

// Truncate truncates a string to maxlen.
func Truncate(s string, maxlen int) string {
	r := []rune(s)
	if len(r) <= maxlen {
		return s
	}

	return string(r[:maxlen])
}

// InSlice tests whether a string is contained in a slice of strings or not.
// Comparison is case insensitive
func InSlice(slice []string, s string) bool {
	for _, ss := range slice {
		if strings.EqualFold(s, ss) {
			return true
		}
	}

	return false
}

// RemoveFromSlice removes a string from a slice.  The string can be present
// multiple times.  The entire slice is iterated.
func RemoveFromSlice(slice []string, s string) (ret []string) {
	for _, ss := range slice {
		if !strings.EqualFold(s, ss) {
			ret = append(ret, ss)
		}
	}

	return ret
}
