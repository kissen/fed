package template

import "strings"

// Return whether s is either nil or ponting to
// an, if trimmed, empty string.
func IsEmpty(s *string) bool {
	if s == nil {
		return true
	}

	trimmed := strings.TrimSpace(*s)
	return len(trimmed) == 0
}
