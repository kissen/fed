package ap

import "strings"

func isEmpty(s string) bool {
	trimmed := strings.TrimSpace(s)
	return len(trimmed) == 0
}
