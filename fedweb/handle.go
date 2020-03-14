package main

import (
	"errors"
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
	log.Printf("GetIndex(%v)", r.URL)

	GetStream(w, r)
}

// GET /stream
func GetStream(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetStream(%v)", r.URL)

	data := map[string]interface{}{
		"Selected": "Stream",
	}

	Render(w, r, "res/collection.page.tmpl", data)
}

// GET /liked
func GetLiked(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetLiked(%v)", r.URL)

	Flash(r, "yes")

	data := map[string]interface{}{
		"Selected": "Liked",
	}

	Render(w, r, "res/collection.page.tmpl", data)
}

// GET /following
func GetFollowing(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetFollowing(%v)", r.URL)

	data := map[string]interface{}{
		"Selected": "Following",
	}

	Render(w, r, "res/collection.page.tmpl", data)
}

// GET /followers
func GetFollowers(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetFollowers(%v)", r.URL)

	data := map[string]interface{}{
		"Selected": "Followers",
	}

	Render(w, r, "res/collection.page.tmpl", data)
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

	// fetch and wrap object

	wrapped, err := wocab.Fetch(iri)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err, nil)
		return
	}

	// set up data dict and render

	data := map[string]interface{}{
		"Item": wrapped,
	}

	Render(w, r, "res/collection.page.tmpl", data)
}

// GET /login
func GetLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetLogin(%v)", r.URL)

	Error(w, r, http.StatusNotImplemented, nil, nil)
}

// POST /login
func PostLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("PostLogin(%v)", r.URL)

	Error(w, r, http.StatusNotImplemented, nil, nil)
}

// GET /static/*
func GetStatic(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetStatic(%v)", r.URL)

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

// POST /submit
func PostSubmit(w http.ResponseWriter, r *http.Request) {
	contents := strings.TrimSpace(r.FormValue("postinput"))

	if len(contents) == 0 {
		Error(w, r, http.StatusBadRequest, errors.New("empty note"), nil)
		return
	}

	Error(w, r, http.StatusNotImplemented, nil, nil)
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

func Render(w http.ResponseWriter, r *http.Request, page string, data map[string]interface{}) {
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
	data["Context"] = Context(r)

	// load template files

	templates := []string{
		page, "res/base.layout.tmpl", "res/card.fragment.tmpl",
		"res/flash.fragment.tmpl",
	}

	// compile template

	ts, err := template.ParseFiles(templates...)
	if err != nil {
		log.Printf("parsing templates failed: %v", err)
		return
	}

	// write http status

	status := Context(r).Status
	w.WriteHeader(status)

	// write body

	if err := ts.Execute(w, data); err != nil {
		log.Printf("executing template failed: %v", err)
		return
	}
}
