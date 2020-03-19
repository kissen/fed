package main

import (
	"encoding/base64"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/fedweb/fedclient"
	"log"
	"net/http"
)

const _FEDWEB_CONTEXT_KEY = "FedWebContext"
const _COOKIE_CONTEXT_KEY = "CookieContext"

type FedWebContext struct {
	VolatileContext
	CookieContext
}

type VolatileContext struct {
	// The name of the tab that should be highlighted in the
	// navigation bar. If empty, not tab will be highlighted.
	Selected string

	// The title that should be used.
	Title string

	// The HTTP status code to reply with. After being set the
	// first time with function Status(), it will not change anymore.
	// This means a handler can (1) set the status and then (2)
	// just call another handler to take care of the request w/o
	// changing the HTTP status code.
	//
	// Initialized to 200.
	Status int

	// Currently logged in user for this session. Might be nil.
	Client fedclient.FedClient
}

type CookieContext struct {
	// Username from login, if there is one.
	Username *string

	// ActorIRI from login, if there is one.
	ActorIRI *string

	// Flashes to display on top of the page. Might be nil.
	Flashs []string

	// Warning (yellow) flashes to display on the top of the page.
	// Might be nil.
	Warnings []string

	// Error (red) flashes to display on the top of the page.
	// Might be nil.
	Errors []string
}

// Load persisted fields from cookie into context.
func (cc *CookieContext) LoadFromCookie(r *http.Request) error {
	// check if a cookie is even set

	cookie, err := r.Cookie(_COOKIE_CONTEXT_KEY)
	if err == http.ErrNoCookie {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "error retrieving cookie")
	}

	// convert the base64 value of the cookie to json binary data

	text, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return errors.Wrap(err, "cookie value malformed base64")
	}

	// try to interpret the contents of the cookie

	var buf CookieContext

	if err := json.Unmarshal(text, &buf); err != nil {
		return errors.Wrap(err, "cookie unmarshal failed")
	}

	// set fields; we either want both username and actor iri set
	// or none of them

	cc.Flashs = buf.Flashs
	cc.Warnings = buf.Warnings
	cc.Errors = buf.Errors

	cc.Username = nil
	cc.ActorIRI = nil

	if !IsEmpty(cc.Username) && !IsEmpty(cc.ActorIRI) {
		cc.Username = buf.Username
		cc.ActorIRI = buf.ActorIRI
	}

	return nil
}

// Write persisted fields from context into cookie.
func (cc *CookieContext) WriteToCookie(w http.ResponseWriter) error {
	// prepare a sanitized copy of cc in buf

	var buf CookieContext = *cc

	if IsEmpty(buf.Username) || IsEmpty(buf.ActorIRI) {
		buf.Username = nil
		buf.ActorIRI = nil
	}

	// conver buf to json

	text, err := json.Marshal(&buf)
	if err != nil {
		return errors.Wrap(err, "cookie marshal failed")
	}

	// convert json to base64

	encoded := base64.StdEncoding.EncodeToString(text)

	// build up cookie

	cookie := http.Cookie{
		Name:  _COOKIE_CONTEXT_KEY,
		Value: encoded,
	}

	// send out cookie wit the response

	http.SetCookie(w, &cookie)
	return nil
}

func (cc *CookieContext) NewClient() (fedclient.FedClient, error) {
	if IsEmpty(cc.Username) || IsEmpty(cc.ActorIRI) {
		return nil, nil
	}

	client, err := fedclient.New(*cc.ActorIRI)
	if err != nil {
		return nil, errors.Wrap(err, "creating client with credentials failed")
	}

	return client, nil
}

// Set all flash slices to nil.
func (cc *CookieContext) ClearFlashes() {
	cc.Flashs = nil
	cc.Warnings = nil
	cc.Errors = nil
}

func AddContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(_FEDWEB_CONTEXT_KEY) == nil {
			var err error

			// create context w/ default values

			fc := &FedWebContext{}
			fc.Status = http.StatusOK

			// try to read cookies if there are any any

			if err = fc.LoadFromCookie(r); err != nil {
				log.Println(err)
			}

			// try creating a client

			if fc.Client, err = fc.NewClient(); err != nil {
				log.Println(err)
			}

			// install context into request

			c := context.WithValue(r.Context(), _FEDWEB_CONTEXT_KEY, fc)
			r = r.WithContext(c)
		}

		next.ServeHTTP(w, r)
	})
}

func Context(r *http.Request) *FedWebContext {
	return r.Context().Value(_FEDWEB_CONTEXT_KEY).(*FedWebContext)
}

func Selected(r *http.Request, tab string) {
	if fc := Context(r); len(fc.Selected) == 0 {
		fc.Selected = tab
	}
}

func Title(r *http.Request, title string) {
	if fc := Context(r); len(fc.Title) == 0 {
		fc.Title = title
	}
}

func Status(r *http.Request, status int) {
	if fc := Context(r); fc.Status == 0 {
		fc.Status = status
	}
}

func Username(r *http.Request, username string) {
	Context(r).Username = &username
}

func ActorIRI(r *http.Request, actorIRI string) {
	Context(r).ActorIRI = &actorIRI
}

func Flash(r *http.Request, s string) {
	fc := Context(r)
	fc.Flashs = append(fc.Flashs, s)
}

func FlashWarning(r *http.Request, s string) {
	fc := Context(r)
	fc.Warnings = append(fc.Warnings, s)
}

func FlashError(r *http.Request, s string) {
	fc := Context(r)
	fc.Errors = append(fc.Errors, s)
}
