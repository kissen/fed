package main

import (
	"github.com/kissen/fed/ap"
	"log"
	"net/http"
)

func PutUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("PutUser(%v)", r.URL)

	// AdminProtocol takes care of everything :)
	adm := &ap.FedAdminProtocol{}
	adm.Handle(r.Context(), w, r)
}
