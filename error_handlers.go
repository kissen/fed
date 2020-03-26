package main

import (
	"gitlab.cs.fau.de/kissen/fed/template"
	"net/http"
)

// HTTP handler that handles a not found error.
func NotFound(w http.ResponseWriter, r *http.Request) {
	Error(w, r, http.StatusNotFound)
}

// HTTP handler that handles a method not allowed error.
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	Error(w, r, http.StatusMethodNotAllowed)
}

// Return an HTTP error that indicates that the type was wrong.
func WrongContentType(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		Error(w, r, http.StatusNotAcceptable)
	case "POST":
		fallthrough
	case "PUT":
		Error(w, r, http.StatusUnsupportedMediaType)
	}
}

// HTTP handler that handles error with code status.
func Error(w http.ResponseWriter, r *http.Request, status int) {
	if IsHTMLRequest(r) {
		template.Error(w, r, status, nil, nil)
	} else {
		ApiError(w, r, "routing error", nil, status)
	}
}
