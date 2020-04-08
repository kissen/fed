package main

import (
	"gitlab.cs.fau.de/kissen/fed/util"
	"log"
	"net/http"
	"path"
)

// GET /ostatus_subscribe?acct={uri}
//
// A certain elephant is doing weird things; see
// https://git.pleroma.social/pleroma/pleroma/issues/286
func GetOStatusSubscribe(w http.ResponseWriter, r *http.Request) {
	log.Println("GetOstatusSubscribe()")

	acct, ok := util.FormValue(r, "acct")
	if !ok {
		ApiError(w, r, "missing acct", http.StatusBadRequest)
		return
	}

	location := path.Join("/remote", acct)
	http.Redirect(w, r, location, http.StatusFound)
}
