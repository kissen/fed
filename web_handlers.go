package main

import (
	"fmt"
	"github.com/go-fed/activity/streams"
	"github.com/gorilla/mux"
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
	"gitlab.cs.fau.de/kissen/fed/template"
	"gitlab.cs.fau.de/kissen/fed/template/fedclient"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func WebGetIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebGetIndex(%v)", r.URL)
	WebGetStream(w, r)
}

func WebGetStream(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebGetStream(%v)", r.URL)

	fedcontext.Title(r, "Your Stream")
	fedcontext.Selected(r, "Stream")

	// get client; if we are not logged in /stream does not make any sense

	client := fedcontext.Context(r).Client
	if client == nil {
		fedcontext.FlashWarning(r, "authorization required")
		RedirectLocal(w, r, "/login")
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

func WebGetLiked(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebGetLiked(%v)", r.URL)

	fedcontext.Title(r, "You Liked")
	fedcontext.Selected(r, "Liked")

	client := fedcontext.Context(r).Client
	if client == nil {
		fedcontext.FlashWarning(r, "authorization requried")
		RedirectLocal(w, r, "/login")
		return
	}

	template.Remote(w, r, client.LikedIRI())
}

func WebGetFollowing(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebGetFollowing(%v)", r.URL)

	fedcontext.Title(r, "Following")
	fedcontext.Selected(r, "Following")

	template.Error(w, r, http.StatusNotImplemented, nil, nil)
}

func WebGetFollowers(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebGetFollowers(%v)", r.URL)

	fedcontext.Title(r, "Followers")
	fedcontext.Selected(r, "Followers")

	template.Error(w, r, http.StatusNotImplemented, nil, nil)
}

func WebGetRemote(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebGetRemote(%v)", r.URL)

	// get and sanitize iri

	query := mux.Vars(r)["remotepath"]

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

func WebGetLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebGetLogin(%v)", r.URL)

	// if we are logged in, forward to stream

	if fedcontext.Context(r).Client != nil {
		fedcontext.Flash(r, "already logged in")
		RedirectLocal(w, r, "/stream")
		return
	}

	// we are not logged in; show the login form

	fedcontext.Title(r, "Login")
	template.Render(w, r, "res/login.page.tmpl", nil)
}

func WebPostLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebPostLogin(%v)", r.URL)

	// check whether we have valid input

	username := strings.TrimSpace(r.FormValue("username"))
	addr := strings.TrimSpace(r.FormValue("iri"))

	if len(username) == 0 {
		fedcontext.FlashWarning(r, "missing username")
		fedcontext.Status(r, http.StatusBadRequest)
		WebGetLogin(w, r)
		return
	}

	if len(addr) == 0 {
		fedcontext.FlashWarning(r, "missing actor IRI")
		fedcontext.Status(r, http.StatusBadRequest)
		WebGetLogin(w, r)
		return
	}

	// try to create a client; this is to check whether the credentials
	// are actually working

	client, err := fedclient.New(addr)
	if err != nil {
		log.Printf(`login username="%v" addr="%v" failed: err="%v"`, username, addr, err)

		fedcontext.FlashWarning(r, "login failed")
		fedcontext.Status(r, http.StatusUnauthorized)
		WebGetLogin(w, r)
		return
	}

	// success; set context and write cookie for later

	context := fedcontext.Context(r)
	context.Client = client
	context.Username = &username
	context.ActorIRI = &addr

	// we are just logged on; forward to stream page for now

	fedcontext.Flash(r, "successfully logged in")
	RedirectLocal(w, r, "/stream")
}

func WebPostLogout(w http.ResponseWriter, r *http.Request) {
	// remove credentials and client from context

	context := fedcontext.Context(r)

	context.Username = nil
	context.ActorIRI = nil
	context.Client = nil

	// w/o login we do nothing

	fedcontext.Flash(r, "logged out")
	RedirectLocal(w, r, "/login")
}

func WebPostSubmit(w http.ResponseWriter, r *http.Request) {
	// check whether we have valid input

	payload := strings.TrimSpace(r.FormValue("postinput"))

	if len(payload) == 0 {
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
		RedirectLocal(w, r, "/login")
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
	RedirectLocal(w, r, "/")
}
