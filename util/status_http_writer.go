package util

import "net/http"

// An http.ResponseWriter that makes it easy to retrieve the status
// code.
type StatusHTTPWriter struct {
	// The response writer to which all calls are forwarded.
	Target http.ResponseWriter

	// The status set by WriteHeader.
	status int
}

func (sw *StatusHTTPWriter) Header() http.Header {
	return sw.Target.Header()
}

func (sw *StatusHTTPWriter) Write(bs []byte) (int, error) {
	return sw.Target.Write(bs)
}

func (sw *StatusHTTPWriter) WriteHeader(status int) {
	sw.status = status
	sw.Target.WriteHeader(status)
}

// Return the status last supplied to WriteHeader or http.StatusOK
// if WriteHeader was not called before.
func (sw *StatusHTTPWriter) Status() int {
	if sw.status == 0 {
		return http.StatusOK
	} else {
		return sw.status
	}
}
