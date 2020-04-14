package fedcontext

import "net/http"

// Redirect to addr preserving the persisted part of
// the current FedContext.
//
// That means writing the cookie.
func Redirect(w http.ResponseWriter, r *http.Request, addr string) {
	Context(r).WriteToCookie(w)
	http.Redirect(w, r, addr, http.StatusFound)
}
