package main

import (
	"gitlab.cs.fau.de/kissen/fed/util"
	"log"
	"net/http"
)

func SignResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pw := util.NewSigningWriter()
		next.ServeHTTP(pw, r)

		if err := pw.ApplyTo(w); err != nil {
			log.Println(err)
		}
	})
}
