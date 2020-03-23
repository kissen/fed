package oauth

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/fedd/db"
	"log"
	"net/http"
	"net/url"
	"sync"
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

	// Map from code to associated username.
	Codes     map[string]string
	CodesLock sync.Mutex

	// Map from token to associated username.
	Tokens     map[string]string
	TokensLock sync.Mutex
}

func New(storage db.FedStorage) FedOAuther {
	return &fedoauther{
		Storage: storage,
		Codes:   make(map[string]string),
		Tokens:  make(map[string]string),
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

	// check wheter credentials are valid

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

	// all good; set set and send out reply

	oa.CodesLock.Lock()
	defer oa.CodesLock.Unlock()

	oa.Codes[code] = user.Name
	log.Printf("recording code=%v for user=%v", code, user.Name)

	http.Redirect(w, r, redirect.String(), http.StatusFound)
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

	oa.CodesLock.Lock()
	defer oa.CodesLock.Unlock()

	username, ok := oa.Codes[code]
	if !ok {
		http.Error(w, "unknown code", http.StatusUnauthorized)
		return
	}

	// create token

	token, err := oa.random()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create reply message

	reply := map[string]interface{}{
		"access_token": token,
		"token_type":   "Bearer",
		"scope":        "all",
		"created_at":   time.Now().Unix(),
	}

	replybytes, err := json.Marshal(&reply)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// record token

	oa.TokensLock.Lock()
	defer oa.TokensLock.Unlock()

	oa.Tokens[token] = username
	log.Printf("recording token=%v for user=%v", token, username)

	// write out response

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(replybytes); err != nil {
		log.Printf("repy with token failed: %v", err)
	}
}

func (oa *fedoauther) UserForCode(code string) (username string, ok bool) {
	log.Printf("UserForCode(%v)", code)
	return oa.userFor(code, oa.Codes, &oa.CodesLock)
}

func (oa *fedoauther) UserForToken(token string) (username string, ok bool) {
	log.Printf("UserForToken(%v)", token)
	return oa.userFor(token, oa.Tokens, &oa.TokensLock)
}

func (oa *fedoauther) userFor(key string, m map[string]string, l *sync.Mutex) (string, bool) {
	l.Lock()
	defer l.Unlock()

	username, ok := m[key]
	return username, ok
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
