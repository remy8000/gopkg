package tools

import (
	"regexp"
	"strings"
)

// RemoveRedundantWhitespaces takes a string `s` and removes any redundant whitespace.
// It trims leading and trailing spaces, and consolidates multiple spaces within the string to a single space.
func RemoveRedundantWhitespaces(s string) string {
	// First, remove leading and trailing whitespace using TrimSpace.
	// This removes spaces, tabs, and other whitespace from both ends of the string.
	s = strings.TrimSpace(s)

	// Define a regular expression to identify sequences of two or more whitespace characters
	// (`[\s\p{Zs}]{2,}`). Here:
	// - `\s` matches any whitespace character (spaces, tabs, newlines, etc.)
	// - `\p{Zs}` matches Unicode whitespace (e.g., non-breaking spaces)
	// - `{2,}` indicates two or more occurrences.
	reInsideWhitespace := regexp.MustCompile(`[\s\p{Zs}]{2,}`)

	// Use the regex to replace all sequences of two or more spaces with a single space.
	s = reInsideWhitespace.ReplaceAllString(s, " ")

	// Return the cleaned-up string, with leading/trailing whitespace removed
	// and multiple spaces inside reduced to a single space.
	return s
}

// Excerpt builds an excerpt from any string by specifying count of words
// Only adds "..." if the text was actually truncated
func Excerpt(s string, wordsCount int) string {
	var excerpt string
	var truncated bool
	if len(s) > 0 {
		split := strings.Split(s, " ")
		if len(split) > wordsCount {
			excerpt = strings.Join(split[:wordsCount], " ")
			truncated = true
		} else {
			excerpt = s
		}
	}
	if truncated {
		return excerpt + "..."
	}
	return excerpt
}
