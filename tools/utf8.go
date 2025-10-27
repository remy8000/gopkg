package tools
import (
    "unicode/utf8"
    "html"
)


// Utf8Validate validates a UTF-8 encoded string and decodes HTML entities.
// If the string contains invalid UTF-8 characters, they are removed.
// After validation, any HTML entities (e.g., &#8230;) are decoded to their UTF-8 equivalents.
func Utf8Validate(content string) string {
	// Check if the content string is valid UTF-8
	if !utf8.ValidString(content) {
		// Create a slice of runes with an initial capacity of the string length
		// (Using runes to properly handle multi-byte characters)
		v := make([]rune, 0, len(content))
		
		// Iterate through each rune in the string
		for i, r := range content {
			// Check if the current rune is a UTF-8 error marker
			if r == utf8.RuneError {
				// Decode the next rune in the string to determine its size
				_, size := utf8.DecodeRuneInString(content[i:])
				
				// If the size is 1, it's an invalid rune, so skip it
				// (Size of 1 indicates an invalid rune in the input)
				if size == 1 {
					continue
				}
			}
			// Append valid runes to the new slice
			v = append(v, r)
		}
		
		// Convert the slice of valid runes back to a string
		content = string(v)
	}

	// Decode any HTML entities (e.g., &#8230;) into their UTF-8 characters
	content = html.UnescapeString(content)

	// Return the validated and decoded string
	return content
}
