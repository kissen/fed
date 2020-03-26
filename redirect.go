package main

import (
	"gitlab.cs.fau.de/kissen/fed/config"
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
	"net/http"
	"path"
)

// Redirect to addr.
func Redirect(w http.ResponseWriter, r *http.Request, addr string) {
	// try to write out the cookies
	context := fedcontext.Context(r)
	context.WriteToCookie(w)

	// actual redirect
	http.Redirect(w, r, addr, http.StatusFound)
}

// Redirect to addr on our instance. So if addr=/following and our instance
// is runing on example.com/fed/, this function issues a redirect to
// example.com/fed/following.
func RedirectLocal(w http.ResponseWriter, r *http.Request, target string) {
	addr := *config.Get().Base

	addr.Path = path.Join(addr.Path, target)
	Redirect(w, r, addr.String())
}
