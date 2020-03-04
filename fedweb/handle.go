package main

import (
	"html/template"
	"log"
	"net/http"
)

// Path=/static/*
func GetStatic(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusNotImplemented, nil)
}

// Path=/
func GetStream(w http.ResponseWriter, r *http.Request) {
	Render(w, "stream.page.tmpl", nil, http.StatusOK)
}

// Path=/login
func GetLogin(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusNotImplemented, nil)
}

// Path=/login
func PostLogin(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusNotImplemented, nil)
}

// Error Handler
func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusNotFound, nil)
}

// Error Handler
func HandleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusMethodNotAllowed, nil)
}

// Write out the Error template with given status and cause.
// cause may be left nil.
func Error(w http.ResponseWriter, status int, cause error) {
	data := map[string]interface{}{
		"Status":     status,
		"StatusText": http.StatusText(status),
	}

	if cause != nil {
		data["Message"] = cause.Error()
	} else {
		data["Message"] = ""
	}

	Render(w, "error.page.tmpl", data, status)
}

func Render(w http.ResponseWriter, page string, data interface{}, status int) {
	templates := []string{
		page, "base.layout.tmpl",
	}

	ts, err := template.ParseFiles(templates...)

	if err != nil {
		log.Printf("parsing templates failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	if ts.Execute(w, data); err != nil {
		log.Printf("executing template failed: %v", err)
		return
	}
}
