package main

import (
	"net/http"
	"strings"
)

// Return the form value with given key from request r. The value
// may not be empty or only contain whitespace characters.
func FormValue(r *http.Request, key string) (value string, ok bool) {
	value = r.FormValue(key)
	tv := strings.TrimSpace(value)

	if len(tv) == 0 {
		return "", false
	} else {
		return tv, true
	}
}
