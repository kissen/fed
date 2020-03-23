package main

import (
	"gitlab.cs.fau.de/kissen/fed/fedd/oauth"
	"log"
	"net/http"
)

func newAuthorizeHandler(oa oauth.FedOAuther) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AuthorizeHandler(%v)", r.URL)
		oa.HandleAuthorizeRequest(w, r)
	}
}

func newTokenHandler(oa oauth.FedOAuther) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("TokenHandler(%v)", r.URL)
		oa.HandleTokenRequest(w, r)
	}
}
