package main

import (
	"html/template"
	"net/http"
)

// Path=/static/*
func GetStatic(w http.ResponseWriter, r *http.Request) {
	GetStream(w, r)
}

// Path=/
func GetStream(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("base.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"Title": "My Cool Website",
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Path=/login
func GetLogin(w http.ResponseWriter, r *http.Request) {
}

// Path=/login
func PostLogin(w http.ResponseWriter, r *http.Request) {
}
