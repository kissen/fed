package main

import (
	"fmt"
	"github.com/go-fed/activity/streams"
	"github.com/gorilla/mux"
	"gitlab.cs.fau.de/kissen/fed/db"
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
	"gitlab.cs.fau.de/kissen/fed/template"
	"gitlab.cs.fau.de/kissen/fed/util"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GET /
func WebGetIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebGetIndex(%v)", r.URL)
	WebGetStream(w, r)
}

// GET /stream
func WebGetStream(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebGetStream(%v)", r.URL)

	fedcontext.Title(r, "Your Stream")
	fedcontext.Selected(r, "Stream")

	// get client; if we are not logged in /stream does not make any sense

	client := fedcontext.Context(r).Client
	if client == nil {
		fedcontext.FlashWarning(r, "authorization required")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// render out the collection

	stream, err := client.Stream()
	if err != nil {
		template.Error(w, r, http.StatusBadGateway, err, nil)
		return
	}

	template.Iter(w, r, stream)
}

// GET /liked
func WebGetLiked(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebGetLiked(%v)", r.URL)

	fedcontext.Title(r, "You Liked")
	fedcontext.Selected(r, "Liked")

	client := fedcontext.Context(r).Client
	if client == nil {
		fedcontext.FlashWarning(r, "authorization requried")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	template.Remote(w, r, client.LikedIRI())
}

// GET /following
func WebGetFollowing(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebGetFollowing(%v)", r.URL)

	fedcontext.Title(r, "Following")
	fedcontext.Selected(r, "Following")

	template.Error(w, r, http.StatusNotImplemented, nil, nil)
}

// GET /followers
func WebGetFollowers(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebGetFollowers(%v)", r.URL)

	fedcontext.Title(r, "Followers")
	fedcontext.Selected(r, "Followers")

	template.Error(w, r, http.StatusNotImplemented, nil, nil)
}

// GET /remote/{remote_path}
func WebGetRemote(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebGetRemote(%v)", r.URL)

	// get and sanitize iri

	query := mux.Vars(r)["remote_path"]

	iri, err := url.Parse(query)
	if err != nil {
		template.Error(w, r, http.StatusBadRequest, err, nil)
		return
	}

	iri.Path = strings.TrimLeft(iri.Path, "/")

	// re-add query params for the remote if there were any

	s := iri.String()

	for key, value := range r.URL.Query() {
		s += fmt.Sprintf("?%v=%v", key, value[0])
	}

	iri, err = url.Parse(s)
	if err != nil {
		template.Error(w, r, http.StatusInternalServerError, err, nil)
		return
	}

	// let our friend Remote take care of it

	template.Remote(w, r, iri)
}

// GET /login
func WebGetLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebGetLogin(%v)", r.URL)

	// if we are logged in, forward to stream

	if fedcontext.Context(r).Client != nil {
		fedcontext.Flash(r, "already logged in")
		http.Redirect(w, r, "/stream", http.StatusFound)
		return
	}

	// we are not logged in; show the login form

	fedcontext.Title(r, "Login")
	template.Render(w, r, "res/login.page.tmpl", nil)
}

// POST /login
func WebPostLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebPostLogin(%v)", r.URL)
	context := fedcontext.Context(r)

	// check whether we have valid input
	username, ok := util.FormValue(r, "username")
	if !ok {
		fedcontext.FlashWarning(r, "missing username")
		fedcontext.Status(r, http.StatusBadRequest)
		WebGetLogin(w, r)
		return
	}
	password, ok := util.FormValue(r, "password")
	if !ok {
		fedcontext.FlashWarning(r, "missing password")
		fedcontext.Status(r, http.StatusBadRequest)
		WebGetLogin(w, r)
		return
	}

	// try to create a session with the supplied credentials
	tm, err := db.NewFedOAuthToken(username, password, context.Storage)
	if err != nil {
		fedcontext.FlashWarning(r, err.Error())
		fedcontext.Status(r, http.StatusUnauthorized)
		WebGetLogin(w, r)
		return
	}

	// success; set context and write cookie for later
	context.Token = &tm.Token

	// we are just logged on; forward to stream page for now
	fedcontext.Flash(r, "successfully logged in")
	http.Redirect(w, r, "/stream", http.StatusFound)
}

// POST /logout
func WebPostLogout(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebPostLogout()")

	// remove credentials and client from context
	context := fedcontext.Context(r)
	context.Token = nil

	// redirect to login page
	fedcontext.Flash(r, "logged out")
	http.Redirect(w, r, "/login", http.StatusFound)
}

// POST /submit
func WebPostSubmit(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebPostSubmit()")

	// check whether we have valid input

	payload, ok := util.FormValue(r, "postinput")

	if !ok {
		fedcontext.FlashWarning(r, "cowardly refusing to create an empty note")
		fedcontext.Status(r, http.StatusBadRequest)

		WebGetStream(w, r)
		return
	}

	if len(payload) > 1024 {
		StatusPayloadTooLarge := 413
		template.Error(w, r, StatusPayloadTooLarge, nil, nil)
	}

	// retreive the client session

	client := fedcontext.Context(r).Client
	if client == nil {
		fedcontext.FlashWarning(r, "authorization requried")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// build up the note

	note := streams.NewActivityStreamsNote()

	attrib := streams.NewActivityStreamsAttributedToProperty()
	attrib.AppendIRI(client.IRI())
	note.SetActivityStreamsAttributedTo(attrib)

	content := streams.NewActivityStreamsContentProperty()
	content.AppendXMLSchemaString(payload)
	note.SetActivityStreamsContent(content)

	published := streams.NewActivityStreamsPublishedProperty()
	published.Set(time.Now())
	note.SetActivityStreamsPublished(published)

	// post it to the server

	if err := client.Create(note); err != nil {
		template.Error(w, r, http.StatusBadGateway, err, nil)
		return
	}

	// redirect to index page for now; we'll improve this later

	fedcontext.Flash(r, "submitted")
	http.Redirect(w, r, "/", http.StatusFound)
}

// POST /reply
func WebPostReply(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebPostReply()")
	template.Error(w, r, http.StatusNotImplemented, nil, nil)
}

// POST /repeat
func WebPostRepeat(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebPostRepeat()")
	template.Error(w, r, http.StatusNotImplemented, nil, nil)
}

// POST /like
func WebPostLike(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebPostLike()")

	iri, done := getIri(w, r)
	if done {
		return
	}

	client := fedcontext.Context(r).Client
	if client == nil {
		template.Error(w, r, http.StatusUnauthorized, nil, nil)
		return
	}

	if err := client.Like(iri); err != nil {
		template.Error(w, r, http.StatusBadGateway, err, nil)
		return
	}

	fedcontext.Flash(r, "liked")
	http.Redirect(w, r, "/", http.StatusFound)
}

// Try to get the iri_base64 form value from POST request r.
// If it is missing or malformed, this functions writes out
// and error and returns (nil, true).
func getIri(w http.ResponseWriter, r *http.Request) (iri *url.URL, handled bool) {
	payload64, ok := util.FormValue(r, "iri_base64")
	if !ok {
		template.Error(w, r, http.StatusBadRequest, nil, nil)
		return nil, true
	}

	payload, err := util.DecodeBase64ToString(payload64)
	if err != nil {
		template.Error(w, r, http.StatusBadRequest, err, nil)
		return nil, true
	}

	iri, err = url.Parse(payload)
	if err != nil {
		template.Error(w, r, http.StatusBadRequest, err, nil)
		return nil, true
	}

	return iri, false
}
