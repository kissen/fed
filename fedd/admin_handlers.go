package main

import (
	"gitlab.cs.fau.de/kissen/fed/fedd/ap"
	"gitlab.cs.fau.de/kissen/fed/fedd/db"
	"log"
	"net/http"
)

func newAdminHandler(admin *ap.FedAdminProtocol, store db.FedStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AdminHandler(%v)", r.URL)

		// AdminProtocol takes care of everything :)
		admin.Handle(r.Context(), w, r)
	}
}
