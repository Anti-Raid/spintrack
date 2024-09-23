package strutils

import "strings"

// Reverse replaces the first n occurrences of old with new in s from the back of the string
func ReplaceFromBack(s string, old string, new string, n int) string {
	// Find the last n occurrences of old
	// Replace them with new
	// Return the string
	for i := 0; i < n; i++ {
		lastIdx := strings.LastIndex(s, old)

		if lastIdx == -1 {
			break
		}

		s = s[:lastIdx] + new + s[lastIdx+len(old):]
	}

	return s
}
