package main

import (
	"encoding/json"
	"fmt"
	"gitlab.cs.fau.de/kissen/fed/db"
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
	"gitlab.cs.fau.de/kissen/fed/template"
	"gitlab.cs.fau.de/kissen/fed/util"
	"log"
	"net/http"
	"net/url"
)

func GetOAuthAuthorize(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetOAuthAuthorize()")

	if done := validateOAuthAuthorize(w, r); done {
		return
	}

	template.Render(w, r, "res/authorize.page.tmpl", nil)
}

func PostOAuthAuthorize(w http.ResponseWriter, r *http.Request) {
	log.Printf("PostOAuthAuthorize()")

	if done := validateOAuthAuthorize(w, r); done {
		return
	}

	// get login credentials

	username, ok := util.FormValue(r, "username")
	if !ok {
		fedcontext.FlashWarning(r, "missing username")
		fedcontext.Status(r, http.StatusBadRequest)

		GetOAuthAuthorize(w, r)
		return
	}

	password, ok := util.FormValue(r, "password")
	if !ok {
		fedcontext.FlashWarning(r, "missing password")
		fedcontext.Status(r, http.StatusBadRequest)

		GetOAuthAuthorize(w, r)
		return
	}

	// generate a code

	storage := fedcontext.Context(r).Storage
	code, err := db.NewFedOAuthCode(username, password, storage)
	if err != nil {
		ApiError(w, r, err, http.StatusInternalServerError)
		return
	}

	// prepare the redirect addr

	redirectUris, ok := r.URL.Query()["redirect_uri"]
	if !ok {
		ApiError(w, r, "missing redirect_uri", http.StatusBadRequest)
		return
	}

	redirect, err := url.Parse(redirectUris[0])
	if err != nil {
		ApiError(w, r, err, http.StatusBadRequest)
		return
	}

	redirect = util.WithParam(redirect, "code", code.Code)

	// send out reply

	http.Redirect(w, r, redirect.String(), http.StatusFound)
	log.Printf("recorded code=%v for username=%v", code.Code, username)
}

func PostOAuthToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("PostOAuthToken()")

	// parse out args

	args := []string{
		"client_id", "client_secret", "redirect_uri",
		"code", "grant_type",
	}

	qs, err := query(args, r)
	if err != nil {
		ApiError(w, r, err, http.StatusBadRequest)
		return
	}

	// validate args

	code := qs["code"]
	grantType := qs["grant_type"]

	if grantType != "authorization_code" {
		ApiError(w, r, "unsupported grant_type", http.StatusBadRequest)
		return
	}

	// look up user

	cm, err := fedcontext.Context(r).Storage.RetrieveCode(code)
	if err != nil {
		ApiError(w, r, err, http.StatusUnauthorized)
		return
	}

	// create token

	storage := fedcontext.Context(r).Storage
	tokenmeta, err := db.NewFedOAuthTokenFor(cm.Username, storage)
	if err != nil {
		ApiError(w, r, err, http.StatusInternalServerError)
		return
	}

	// create reply message

	reply := map[string]interface{}{
		"access_token": tokenmeta.Token,
		"token_type":   "Bearer",
		"scope":        "all",
		"created_at":   tokenmeta.IssuedOn.Unix(),
	}

	replybytes, err := json.Marshal(&reply)
	if err != nil {
		ApiError(w, r, err, http.StatusInternalServerError)
		return
	}

	// write out response

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(replybytes); err != nil {
		log.Printf("repy with token failed: %v", err)
	}

	log.Printf("recorded token=%v for user=%v", tokenmeta.Token, cm.Username)
}

// Validate a request to /oauth/authorize, that is look at whether all
// mandatory query parameters are present.
func validateOAuthAuthorize(w http.ResponseWriter, r *http.Request) (ok bool) {
	ks := []string{
		"response_type", "client_id", "redirect_uri",
	}

	ps, err := query(ks, r)
	if err != nil {
		ApiError(w, r, err, http.StatusInternalServerError)
		return true
	}

	if ps["response_type"] != "code" {
		ApiError(w, r, "unsupported response_type", http.StatusBadRequest)
		return true
	}

	return false
}

// Given a set of keys, return a map that maps each key to the query
// parameter value in request r. Returns an r if at least one key
// in keys is not present in request r.
func query(keys []string, r *http.Request) (params map[string]string, err error) {
	params = make(map[string]string)
	q := r.URL.Query()

	for _, k := range keys {
		vs, ok := q[k]

		if !ok {
			return nil, fmt.Errorf("missing required param=%v", k)
		}

		if len(vs) != 1 {
			return nil, fmt.Errorf("param=%v needs to appear exactly once", k)
		}

		params[k] = vs[0]
	}

	return params, nil
}
