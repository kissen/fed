package main

import (
	"fmt"
	"github.com/go-fed/activity/streams"
	"github.com/gorilla/mux"
	"github.com/kissen/httpstatus"
	"gitlab.cs.fau.de/kissen/fed/fedutil"
	"gitlab.cs.fau.de/kissen/fed/fedweb/fedclient"
	"gitlab.cs.fau.de/kissen/fed/fedweb/wocab"
	"html/template"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// GET /
func GetIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetIndex(%v)", r.URL)

	GetStream(w, r)
}

// GET /stream
func GetStream(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetStream(%v)", r.URL)

	Title(r, "Your Stream")
	Selected(r, "Stream")

	// get client; if we are not logged in /stream does not make any sense

	client := Context(r).Client
	if client == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// render out the collection

	stream, err := client.Stream()
	if err != nil {
		Error(w, r, http.StatusBadGateway, err, nil)
		return
	}

	Iter(w, r, stream)
}

// GET /liked
func GetLiked(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetLiked(%v)", r.URL)

	Title(r, "You Liked")
	Selected(r, "Liked")

	client := Context(r).Client
	if client == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	Remote(w, r, client.LikedIRI())
}

// GET /following
func GetFollowing(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetFollowing(%v)", r.URL)

	Title(r, "Following")
	Selected(r, "Following")

	Error(w, r, http.StatusNotImplemented, nil, nil)
}

// GET /followers
func GetFollowers(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetFollowers(%v)", r.URL)

	Title(r, "Followers")
	Selected(r, "Followers")

	Error(w, r, http.StatusNotImplemented, nil, nil)
}

// GET /remote
func GetRemote(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetRemote(%v)", r.URL)

	// get and sanitize iri

	query := mux.Vars(r)["remotepath"]

	iri, err := url.Parse(query)
	if err != nil {
		Error(w, r, http.StatusBadRequest, err, nil)
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
		Error(w, r, http.StatusInternalServerError, err, nil)
		return
	}

	// let our friend Remote take care of it

	Remote(w, r, iri)
}

// GET /login
func GetLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetLogin(%v)", r.URL)

	// if we are logged in, forward to stream

	if Context(r).Client != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	// we are not logged in; show the login form

	Title(r, "Login")
	Render(w, r, "res/login.page.tmpl", nil)
}

// POST /login
func PostLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("PostLogin(%v)", r.URL)

	// check whether we have valid input

	username := strings.TrimSpace(r.FormValue("username"))
	addr := strings.TrimSpace(r.FormValue("iri"))

	if len(username) == 0 {
		FlashWarning(r, "missing username")
		Status(r, http.StatusBadRequest)
		GetLogin(w, r)
		return
	}

	if len(addr) == 0 {
		FlashWarning(r, "missing actor IRI")
		Status(r, http.StatusBadRequest)
		GetLogin(w, r)
		return
	}

	// try to create a client; this is to check whether the credentials
	// are actually working

	client, err := fedclient.New(addr)
	if err != nil {
		log.Printf(`login username="%v" addr="%v" failed: err="%v"`, username, addr, err)

		FlashWarning(r, "login failed")
		Status(r, http.StatusUnauthorized)
		GetLogin(w, r)
		return
	}

	// success; set context and write cookie for later

	context := Context(r)
	context.Client = client
	context.Username = &username
	context.ActorIRI = &addr

	context.WriteToCookie(w)

	// we are just logged on; forward to stream page for now

	http.Redirect(w, r, "/", http.StatusFound)
}

// GET /static/*
func GetStatic(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetStatic(%v)", r.URL)

	// find requested file on file system and read it into a buffer

	name := path.Base(r.URL.Path)
	relpath := filepath.Join("res", name)

	content, err := ioutil.ReadFile(relpath)
	if err != nil {
		Error(w, r, http.StatusNotFound, err, nil)
		return
	}

	// set the correct mimetype header

	mimetype := mime.TypeByExtension(path.Ext(name))
	w.Header().Add("Content-Type", mimetype)

	// write out the contents

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(content); err != nil {
		log.Printf("writing static file to client failed: %v", err)
	}
}

// POST /submit
func PostSubmit(w http.ResponseWriter, r *http.Request) {
	// check whether we have valid input

	payload := strings.TrimSpace(r.FormValue("postinput"))

	if len(payload) == 0 {
		FlashWarning(r, "cowardly refusing to create an empty note")
		Status(r, http.StatusBadRequest)

		GetStream(w, r)
		return
	}

	if len(payload) > 1024 {
		StatusPayloadTooLarge := 413
		Error(w, r, StatusPayloadTooLarge, nil, nil)
	}

	// retreive the client session

	client := Context(r).Client
	if client == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
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
		Error(w, r, http.StatusBadGateway, err, nil)
		return
	}

	// redirect to index page for now; we'll improve this later

	http.Redirect(w, r, "/", http.StatusFound)
}

// Handler for Not Found Errors
func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	Error(w, r, http.StatusNotFound, nil, nil)
}

// Handler for Method Not Allowed Errors
func HandleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	Error(w, r, http.StatusMethodNotAllowed, nil, nil)
}

// Write out the Error template with given status and cause.
// cause may be left nil.
func Error(w http.ResponseWriter, r *http.Request, status int, cause error, data map[string]interface{}) {
	// set up data for the error handler

	errorData := map[string]interface{}{
		"Status":      status,
		"StatusText":  http.StatusText(status),
		"Description": httpstatus.Describe(status),
	}

	if cause != nil {
		errorData["Cause"] = cause.Error()
	}

	renderData := fedutil.SumMaps(data, errorData)

	// render with correct status

	Status(r, status)
	Render(w, r, "res/error.page.tmpl", renderData)
}

// Write out a page showing remote content at addr.
func Remote(w http.ResponseWriter, r *http.Request, iri *url.URL) {
	// fetch and wrap object

	wrapped, err := wocab.Fetch(iri)
	if err != nil {
		Error(w, r, http.StatusBadGateway, err, nil)
		return
	}

	// set up data dict and render

	data := map[string]interface{}{
		"Items": []wocab.WebVocab{
			wrapped,
		},
	}

	Title(r, iri.String())
	Render(w, r, "res/collection.page.tmpl", data)
}

// Write out a page showing activity pub content accessible via iter.
func Iter(w http.ResponseWriter, r *http.Request, it fedutil.Iter) {
	// fetch objects

	vs, err := fedutil.FetchIter(it)
	if err != nil {
		Error(w, r, http.StatusBadGateway, err, nil)
		return
	}

	// wrap objects

	wrapped, err := wocab.News(vs...)
	if err != nil {
		Error(w, r, http.StatusBadGateway, err, nil)
		return
	}

	// set up data dict and render

	data := map[string]interface{}{
		"Items": wrapped,
	}

	Render(w, r, "res/collection.page.tmpl", data)
}

func Render(w http.ResponseWriter, r *http.Request, page string, data map[string]interface{}) {
	// fill in values that are (almost) always needed

	data = fedutil.SumMaps(data)

	data["Context"] = Context(r)
	data["SubmitPrompt"] = SubmitPrompt()
	data["Flashs"] = Context(r).Flashs
	data["Warnings"] = Context(r).Warnings
	data["Errors"] = Context(r).Errors

	// load template files

	templates := []string{
		page, "res/base.layout.tmpl", "res/card.fragment.tmpl",
		"res/flash.fragment.tmpl",
	}

	// compile template; if this fails it's a programming error

	ts, err := template.ParseFiles(templates...)
	if err != nil {
		log.Printf("parsing templates failed: %v", err)
		return
	}

	// set cookie for next time

	context := Context(r)
	context.ClearFlashes()
	context.WriteToCookie(w)

	// write http status

	status := Context(r).Status
	w.WriteHeader(status)

	// write body

	if err := ts.Execute(w, data); err != nil {
		log.Printf("executing template failed: %v", err)
		return
	}
}
