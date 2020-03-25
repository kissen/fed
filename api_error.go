package main

import (
	"encoding/json"
	"github.com/kissen/httpstatus"
	"log"
	"net/http"
)

// Answer the HTTP request to some JSON with a nice JSON-encoded error.
// This is similar to http.Error except the client has less trouble
// parsing the error.
func ApiError(w http.ResponseWriter, r *http.Request, description string, err error, status int) {
	reply := map[string]interface{}{
		"status": status,
	}

	if err != nil {
		reply["error"] = err.Error()
	}

	if len(description) > 0 {
		reply["description"] = description
	}

	bs, err := json.Marshal(&reply)
	if err != nil {
		log.Fatal("unexpected: marshal of error message failed:", err)
	}

	http.Error(w, string(bs), status)
}

func ApiNotFound(w http.ResponseWriter, r *http.Request) {
	ApiError(w, r, httpstatus.Describe(http.StatusNotFound), nil, http.StatusNotFound)
}

func ApiNotAllowed(w http.ResponseWriter, r *http.Request) {
	ApiError(w, r, httpstatus.Describe(http.StatusMethodNotAllowed), nil, http.StatusMethodNotAllowed)
}
