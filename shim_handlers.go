package main

import (
	"github.com/kissen/fed/util"
	"log"
	"net/http"
	"path"
)

// GET /shim/ostatus_subscribe?acct={uri}
//
// A certain elephant is doing weird things; see
// https://git.pleroma.social/pleroma/pleroma/issues/286
func GetOStatusSubscribe(w http.ResponseWriter, r *http.Request) {
	log.Println("GetOStatusSubscribe()")

	acct, ok := util.FormValue(r, "acct")
	if !ok {
		ApiError(w, r, "missing acct", http.StatusBadRequest)
		return
	}

	location := path.Join("/remote", acct)
	http.Redirect(w, r, location, http.StatusFound)
}
