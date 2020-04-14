package main

import (
	"github.com/gobuffalo/packr"
	"log"
	"mime"
	"net/http"
	"path"
)

// GET /favicon.ico
func GetFavicon(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/inbox.svg", http.StatusPermanentRedirect)
}

// GET /static/*
func GetStatic(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetStatic(%v)", r.URL)

	box := packr.NewBox("static")
	filename := path.Base(r.URL.Path)

	if !box.Has(filename) {
		DoError(w, r, "no resource with that name", http.StatusNotFound)
		return
	}

	mimetype := mime.TypeByExtension(path.Ext(filename))
	w.Header().Add("Content-Type", mimetype)

	if _, err := w.Write(box.Bytes(filename)); err != nil {
		log.Printf(`serving static filename="%v" failed with err="%v"`, filename, err)
	}
}
