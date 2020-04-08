package util

import "net/http"

// Return the content type requested (GET) or povided (POST)
// by request r.
//
// Returns an empty string if content type could not be
// determined.
func ContentType(r *http.Request) string {
	switch r.Method {
	case "GET":
		return r.Header.Get("Accept")
	case "POST":
		fallthrough
	case "PUT":
		return r.Header.Get("Content-Type")
	default:
		return ""
	}
}
