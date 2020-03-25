package oauth

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/db"
	"log"
	"net/http"
	"net/url"
	"time"
)

// https://docs.joinmastodon.org/methods/apps/oauth/
// https://www.oauth.com/oauth2-servers/server-side-apps/authorization-code/
// https://github.com/go-oauth2/oauth2/blob/master/example/server/server.go
// https://github.com/openshift/osin
// https://auth0.com/docs/protocols/oauth2#how-response-type-works

type FedOAuther interface {
	// HTTP handler that takes care of GET requests to the /authorize
	// endpoint.
	GetAuthorize(w http.ResponseWriter, r *http.Request)

	// HTTP handler that takes care of the POST requests to the /authorize
	// endpoint.
	PostAuthorize(w http.ResponseWriter, r *http.Request)

	// HTTP handler that takes care of PoSt requests to the /token
	// endpoint.
	PostToken(w http.ResponseWriter, r *http.Request)

	// Return the registerted user for the given code. If the code
	// is invalid or expired, ok will be false.
	UserForCode(code string) (username string, ok bool)

	// Return the registerted user for the given token. If the token
	// is invalid or expired, ok will be false.
	UserForToken(token string) (username string, ok bool)
}

type fedoauther struct {
	// Storage for accessing user database.
	Storage db.FedStorage
}

func New(storage db.FedStorage) FedOAuther {
	return &fedoauther{
		Storage: storage,
	}
}

func (oa *fedoauther) GetAuthorize(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetAuthorize()")

	if done := oa.validateAuthorize(w, r); done {
		return
	}

	// HACK AROUND THE FAC THAT BASIC AUTH DOES NOT TRIGGER POST
	//
	// ANYWAYS; THIS WORKS; THINK ABOUT MERGING FED AND FEDWEB?!
	// OR WRITE UI/LOGIN PROMPT
	if _, _, ok := r.BasicAuth(); ok {
		oa.PostAuthorize(w, r)
		return
	}

	w.Header().Add("WWW-Authenticate", `Basic realm="Gondor"`)
	w.WriteHeader(http.StatusUnauthorized)

	w.Write([]byte("Please Authenticate. Thank you!\n"))
}

func (oa *fedoauther) PostAuthorize(w http.ResponseWriter, r *http.Request) {
	log.Printf("PostAuthorize()")

	if done := oa.validateAuthorize(w, r); done {
		return
	}

	// get login credentials

	username, password, ok := r.BasicAuth()
	if !ok {
		oa.GetAuthorize(w, r) // try again
		return
	}

	// check whether credentials are valid

	user, err := oa.Storage.RetrieveUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !user.PasswordOK(password) {
		http.Error(w, "bad password", http.StatusUnauthorized)
		return
	}

	// generate a code

	code, err := oa.random()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// prepare the redirect addr

	redirectUris, ok := r.URL.Query()["redirect_uri"]
	if !ok {
		http.Error(w, "missing redirect_uri", http.StatusBadRequest)
		return
	}

	redirect, err := url.Parse(redirectUris[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	if err := oa.Storage.StoreCode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send out reply

	http.Redirect(w, r, redirect.String(), http.StatusFound)
	log.Printf("recorded code=%v for user=%v", code, user.Name)
}

func (oa *fedoauther) PostToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("PostToken()")

	// parse out args

	args := []string{
		"client_id", "client_secret", "redirect_uri",
		"code", "grant_type",
	}

	qs, err := oa.query(args, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate args

	code := qs["code"]
	grantType := qs["grant_type"]

	if grantType != "authorization_code" {
		http.Error(w, "unsupported grant_type", http.StatusBadRequest)
		return
	}

	// look up user

	codemeta, err := oa.Storage.RetrieveCode(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// create token

	token, err := oa.random()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenmeta := db.FedOAuthToken{
		Token:    token,
		Username: codemeta.Username,
		IssuedOn: time.Now().UTC(),
	}

	if err := oa.Storage.StoreToken(&tokenmeta); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func (oa *fedoauther) UserForCode(code string) (username string, ok bool) {
	log.Printf("UserForCode(%v)", code)

	if meta, err := oa.Storage.RetrieveCode(code); err != nil {
		return "", false
	} else {
		return meta.Username, true
	}
}

func (oa *fedoauther) UserForToken(token string) (username string, ok bool) {
	log.Printf("UserForToken(%v)", token)

	if meta, err := oa.Storage.RetrieveToken(token); err != nil {
		return "", false
	} else {
		return meta.Username, true
	}
}

func (oa *fedoauther) validateAuthorize(w http.ResponseWriter, r *http.Request) (handled bool) {
	ks := []string{
		"response_type", "client_id", "redirect_uri",
	}

	ps, err := oa.query(ks, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return true
	}

	if ps["response_type"] != "code" {
		http.Error(w, "unsupported response_type", http.StatusBadRequest)
		return true
	}

	return false
}

func (oa *fedoauther) query(keys []string, r *http.Request) (params map[string]string, err error) {
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

// Return a random string from a secure source.
func (oa *fedoauther) random() (string, error) {
	nbytes := 16
	b := make([]byte, nbytes)

	if _, err := rand.Read(b); err != nil {
		return "", errors.Wrap(err, "could not generate random string")
	}

	return fmt.Sprintf("%x", b), nil
}
