package main

import (
	"gitlab.cs.fau.de/kissen/fed/template"
	"gitlab.cs.fau.de/kissen/fed/util"
	"net/http"
)

// HTTP handler that handles a not found error.
func NotFound(w http.ResponseWriter, r *http.Request) {
	DoError(w, r, nil, http.StatusNotFound)
}

// HTTP handler that handles a method not allowed error.
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	DoError(w, r, nil, http.StatusMethodNotAllowed)
}

// Return an HTTP error that indicates that the type was wrong.
func WrongContentType(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		DoError(w, r, nil, http.StatusNotAcceptable)
	case "POST":
		fallthrough
	case "PUT":
		DoError(w, r, nil, http.StatusUnsupportedMediaType)
	}
}

func DoError(w http.ResponseWriter, r *http.Request, cause interface{}, status int) {
	ct := util.ContentType(r)

	switch ct {
	case util.HTML_TYPE:
		template.Error(w, r, status, nil, nil)
	case util.AP_TYPE:
		fallthrough
	default:
		ApiError(w, r, cause, status)
	}
}
