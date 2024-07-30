package pkg

import "unicode/utf8"

func isASCIIPrintable(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < 32 || s[i] > 126 {
			return false
		}
	}
	return true
}

func ValidateContent(content string) bool {

	if !utf8.ValidString(content) {
		return false
	}

	return isASCIIPrintable(content)
}
