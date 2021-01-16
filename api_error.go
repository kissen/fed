package main

import (
	"encoding/json"
	"fmt"
	"github.com/kissen/fed/errors"
	"log"
	"net/http"
)

// Answer the HTTP request to some JSON with a nice JSON-encoded error.
// This is similar to http.Error except the client has less trouble
// parsing the error.
//
// This function tries its best to infer a cause string from argument
// cause.
func ApiError(w http.ResponseWriter, r *http.Request, cause interface{}, status int) {
	// find out how to represent the cause

	var causestr string

	if s, ok := cause.(string); ok {
		causestr = s
	}

	if st, ok := cause.(fmt.Stringer); ok {
		causestr = st.String()
	}

	if ce, ok := cause.(error); ok {
		causestr = ce.Error()

		// if cause contains a status code, use that instead;
		// this is kind of hacky I admit
		if es, ok := errors.Status(ce); ok {
			status = es
		}
	}

	// build up reply

	reply := map[string]interface{}{
		"description": http.StatusText(status),
		"status":      status,
	}

	if causestr != "" {
		reply["cause"] = causestr
	}

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
