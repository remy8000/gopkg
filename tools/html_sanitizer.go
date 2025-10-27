package tools

import (
	"regexp"
	"strings"
)

// SanitizedHtmlWithRegexp removes all HTML tags from the given string `doc`
// It performs better than the `sanitizer.StrictPolicy` function by:
// - Removing all HTML tags (e.g., `<h3>`, `<p>`) and replacing them with spaces
// - Removing redundant whitespace created from removed tags or existing spaces
// - Trimming leading and trailing spaces for a clean result
func SanitizedHtmlWithRegexp(doc string) string {
    // Define a regular expression pattern to match any HTML tag.
    // `<[^>]+>`:
    // - `<` matches the opening tag character.
    // - `[^>]+` matches one or more characters that are NOT `>`, ensuring it captures the entire tag until the closing `>`.
    // - `>` matches the closing tag character.
    re := regexp.MustCompile(`<[^>]+>`)

    // Replace all HTML tags with a single space.
    // This helps in cases where adjacent tags (e.g., `</h3><p>`) need to be separated to avoid concatenated text.
    doc = re.ReplaceAllString(doc, " ")

    // Call RemoveRedundantWhitespaces to clean up any extra whitespace created
    // by removing HTML tags, as well as any redundant spaces already in the text.
    doc = RemoveRedundantWhitespaces(doc)

    // Trim any remaining leading or trailing whitespace for a polished result.
    doc = strings.TrimSpace(doc)

    return doc
}