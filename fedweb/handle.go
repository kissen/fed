package main

import (
	"github.com/kissen/httpstatus"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
)

// GET /
func GetIndex(w http.ResponseWriter, r *http.Request) {
}

// Get /stream
func GetStream(w http.ResponseWriter, r *http.Request) {
	user, err := Fetch("http://localhost:9999/ap/bob")
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	data := map[string]interface{}{
		"Selected": "Stream",
		"Objects":  []interface{}{user},
	}

	Render(w, "collection.page.tmpl", data, http.StatusOK)
}

// Get /liked
func GetLiked(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Selected": "Liked",
	}

	Render(w, "collection.page.tmpl", data, http.StatusOK)
}

// Get /following
func GetFollowing(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Selected": "Following",
	}

	Render(w, "collection.page.tmpl", data, http.StatusOK)
}

// Get /followers
func GetFollowers(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Selected": "Followers",
	}

	Render(w, "collection.page.tmpl", data, http.StatusOK)
}

// GET /login
func GetLogin(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusNotImplemented, nil)
}

// POST /login
func PostLogin(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusNotImplemented, nil)
}

// GET /static/*
func GetStatic(w http.ResponseWriter, r *http.Request) {
	filename := path.Base(r.URL.Path)

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("opening file failed: %v", err)
		return
	}

	mimetype := mime.TypeByExtension(path.Ext(filename))
	w.Header().Add("Content-Type", mimetype)

	w.WriteHeader(http.StatusOK)

	if _, err := io.WriteString(w, string(content)); err != nil {
		log.Printf("writing file to client failed: %v", err)
	}
}

// Handler for Not Found Errors
func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusNotFound, nil)
}

// Handler for Method Not Allowed Errors
func HandleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusMethodNotAllowed, nil)
}

// Write out the Error template with given status and cause.
// cause may be left nil.
func Error(w http.ResponseWriter, status int, cause error) {
	data := map[string]interface{}{
		"Status":      status,
		"StatusText":  http.StatusText(status),
		"Description": httpstatus.Describe(status),
	}

	if cause != nil {
		data["Cause"] = cause.Error()
	}

	Render(w, "error.page.tmpl", data, status)
}

func Render(w http.ResponseWriter, page string, data map[string]interface{}, status int) {
	// fill in required fiels in data

	required := []string{
		"Selected",
	}

	for _, key := range required {
		if _, found := data[key]; !found {
			data[key] = ""
		}
	}

	// load template files

	templates := []string{
		page, "base.layout.tmpl",
	}

	// compile template

	ts, err := template.ParseFiles(templates...)

	if err != nil {
		log.Printf("parsing templates failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write http response

	w.WriteHeader(status)

	if ts.Execute(w, data); err != nil {
		log.Printf("executing template failed: %v", err)
		return
	}
}
