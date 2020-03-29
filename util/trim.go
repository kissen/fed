package util

import "strings"

// Return a trimmed version of the contents of sp. ok is true
// only when pointer sp is not nil and pointing to a string that
// consists of at least none-whitespace character.
func Trim(sp *string) (trimmed string, ok bool) {
	if sp == nil {
		return "", false
	}

	ts := strings.TrimSpace(*sp)
	if len(ts) == 0 {
		return ts, false
	}

	return ts, true
}
