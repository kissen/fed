package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"gitlab.cs.fau.de/kissen/fed/db"
	"gitlab.cs.fau.de/kissen/fed/errors"
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
	"gitlab.cs.fau.de/kissen/fed/template"
	"log"
	"net/http"
	"net/url"
	"time"
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

	username, ok := FormValue(r, "username")
	if !ok {
		fedcontext.FlashWarning(r, "missing username")
		fedcontext.Status(r, http.StatusBadRequest)

		GetOAuthAuthorize(w, r)
		return
	}

	password, ok := FormValue(r, "password")
	if !ok {
		fedcontext.FlashWarning(r, "missing password")
		fedcontext.Status(r, http.StatusBadRequest)

		GetOAuthAuthorize(w, r)
		return
	}

	// check whether credentials are valid

	user, err := fedcontext.Context(r).Storage.RetrieveUser(username)
	if err != nil {
		fedcontext.FlashError(r, "no such user")
		fedcontext.Status(r, http.StatusUnauthorized)

		GetOAuthAuthorize(w, r)
		return
	}

	if !user.PasswordOK(password) {
		fedcontext.FlashError(r, "bad password")
		fedcontext.Status(r, http.StatusUnauthorized)

		GetOAuthAuthorize(w, r)
		return
	}

	// generate a code

	code, err := random()
	if err != nil {
		ApiError(w, r, err, http.StatusBadRequest)
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

	qu := redirect.Query()
	qu.Add("code", code)
	redirect.RawQuery = qu.Encode()

	// all good; write code to db

	c := db.FedOAuthCode{
		Code:     code,
		Username: user.Name,
		IssuedOn: time.Now().UTC(),
	}

	if err := fedcontext.Context(r).Storage.StoreCode(&c); err != nil {
		ApiError(w, r, err, http.StatusInternalServerError)
		return
	}

	// send out reply

	http.Redirect(w, r, redirect.String(), http.StatusFound)
	log.Printf("recorded code=%v for user=%v", code, user.Name)
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

	codemeta, err := fedcontext.Context(r).Storage.RetrieveCode(code)
	if err != nil {
		ApiError(w, r, err, http.StatusUnauthorized)
		return
	}

	// create token

	token, err := random()
	if err != nil {
		ApiError(w, r, err, http.StatusInternalServerError)
		return
	}

	tokenmeta := db.FedOAuthToken{
		Token:    token,
		Username: codemeta.Username,
		IssuedOn: time.Now().UTC(),
	}

	if err := fedcontext.Context(r).Storage.StoreToken(&tokenmeta); err != nil {
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

	log.Printf("recorded token=%v for user=%v", token, codemeta.Username)
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

// Return a random string from a secure source that is sufficently long
// for use as OAuth codes or tokens.
func random() (string, error) {
	nbytes := 16
	b := make([]byte, nbytes)

	if _, err := rand.Read(b); err != nil {
		return "", errors.Wrap(err, "could not generate random string")
	}

	return fmt.Sprintf("%x", b), nil
}
