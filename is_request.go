package main

import "net/http"

// Return whether r is a request serving or requesting ActivityPub
// data.
func IsActivityPubRequest(r *http.Request) bool {
	return isRequest(r, AP_TYPE)
}

// Return whether r is a rquest serving or requesting HTML.
func IsHTMLRequest(r *http.Request) bool {
	return isRequest(r, HTML_TYPE)
}

func isRequest(r *http.Request, contentType string) bool {
	switch r.Method {
	case "GET":
		return r.Header.Get("Accept") == contentType
	case "POST":
		fallthrough
	case "PUT":
		return r.Header.Get("Content-Type") == contentType
	default:
		return false
	}
}
