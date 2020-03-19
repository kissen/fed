package main

import "strings"

func IsEmpty(s *string) bool {
	if s == nil {
		return true
	}

	trimmed := strings.TrimSpace(*s)
	return len(trimmed) == 0
}
