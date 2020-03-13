package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kissen/httpstatus"
	"gitlab.cs.fau.de/kissen/fed/fedutil"
	"gitlab.cs.fau.de/kissen/fed/fedweb/wocab"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

// GET /
func GetIndex(w http.ResponseWriter, r *http.Request) {
	http.RedirectHandler("stream", http.StatusFound).ServeHTTP(w, r)
}

// Get /stream
func GetStream(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Selected": "Stream",
	}

	Render(w, "res/collection.page.tmpl", data, http.StatusOK)
}

// Get /liked
func GetLiked(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Selected": "Liked",
	}

	Render(w, "res/collection.page.tmpl", data, http.StatusOK)
}

// Get /following
func GetFollowing(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Selected": "Following",
	}

	Render(w, "res/collection.page.tmpl", data, http.StatusOK)
}

// Get /followers
func GetFollowers(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Selected": "Followers",
	}

	Render(w, "res/collection.page.tmpl", data, http.StatusOK)
}

func GetRemote(w http.ResponseWriter, r *http.Request) {
	// get and sanitize iri

	query := mux.Vars(r)["remotepath"]

	iri, err := url.Parse(query)
	if err != nil {
		Error(w, http.StatusBadRequest, err, nil)
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
		Error(w, http.StatusInternalServerError, err, nil)
		return
	}

	// fetch and wrap object

	wrapped, err := wocab.Fetch(iri)
	if err != nil {
		Error(w, http.StatusInternalServerError, err, nil)
		return
	}

	// set up data dict and render

	data := map[string]interface{}{
		"Item": wrapped,
	}

	Render(w, "res/collection.page.tmpl", data, http.StatusOK)
}

// GET /login
func GetLogin(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusNotImplemented, nil, nil)
}

// POST /login
func PostLogin(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusNotImplemented, nil, nil)
}

// GET /static/*
func GetStatic(w http.ResponseWriter, r *http.Request) {
	name := path.Base(r.URL.Path)
	relpath := filepath.Join("res", name)

	content, err := ioutil.ReadFile(relpath)
	if err != nil {
		log.Printf("opening file failed: %v", err)
		return
	}

	mimetype := mime.TypeByExtension(path.Ext(name))
	w.Header().Add("Content-Type", mimetype)

	w.WriteHeader(http.StatusOK)

	if _, err := io.WriteString(w, string(content)); err != nil {
		log.Printf("writing file to client failed: %v", err)
	}
}

// Handler for Not Found Errors
func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusNotFound, nil, nil)
}

// Handler for Method Not Allowed Errors
func HandleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusMethodNotAllowed, nil, nil)
}

// Write out the Error template with given status and cause.
// cause may be left nil.
func Error(w http.ResponseWriter, status int, cause error, data map[string]interface{}) {
	// set up data for the error handler

	errorData := map[string]interface{}{
		"Status":      status,
		"StatusText":  http.StatusText(status),
		"Description": httpstatus.Describe(status),
	}

	if cause != nil {
		errorData["Cause"] = cause.Error()
	}

	// join with other generic keys; render

	renderData := fedutil.SumMaps(data, errorData)
	Render(w, "res/error.page.tmpl", renderData, status)
}

func Render(w http.ResponseWriter, page string, data map[string]interface{}, status int) {
	// fill in required fields that need to have some well defined values

	required := []string{
		"Selected",
	}

	for _, key := range required {
		if _, found := data[key]; !found {
			data[key] = ""
		}
	}

	// fill in values that are (almost) always needed

	data["SubmitPrompt"] = SubmitPrompt()

	// load template files

	templates := []string{
		page, "res/base.layout.tmpl", "res/card.fragment.tmpl",
	}

	// compile template

	ts, err := template.ParseFiles(templates...)
	if err != nil {
		log.Printf("parsing templates failed: %v", err)
		return
	}

	// write http response

	w.WriteHeader(status)

	if err := ts.Execute(w, data); err != nil {
		log.Printf("executing template failed: %v", err)
		return
	}
}
