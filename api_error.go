package main

import (
	"encoding/json"
	"gitlab.cs.fau.de/kissen/fed/errors"
	"log"
	"net/http"
)

// Answer the HTTP request to some JSON with a nice JSON-encoded error.
// This is similar to http.Error except the client has less trouble
// parsing the error.
func ApiError(w http.ResponseWriter, r *http.Request, description string, err error, status int) {
	// if err contains a status code, use that instead

	if estatus, ok := errors.Status(err); ok {
		status = estatus
	}

	// build up json map

	reply := map[string]interface{}{
		"status": status,
	}

	if err != nil {
		reply["error"] = err.Error()
	}

	if len(description) > 0 {
		reply["description"] = description
	}

	// create json bytes; this really should never fail

	bs, err := json.Marshal(&reply)
	if err != nil {
		log.Fatal("unexpected: marshal of error message failed:", err)
	}

	// write out response

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(bs); err != nil {
		log.Printf("writing err json to client failed: %v", err)
	}
}
