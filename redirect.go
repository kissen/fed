package main

import (
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
	"net/http"
)

// Redirect to addr. Makes sure a browser writes our cookie before
// leaving.
func Redirect(w http.ResponseWriter, r *http.Request, addr string) {
	// try to write out the cookies
	context := fedcontext.Context(r)
	context.WriteToCookie(w)

	// actual redirect
	http.Redirect(w, r, addr, http.StatusFound)
}
