package main

import (
	"gitlab.cs.fau.de/kissen/fed/template"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
	"path/filepath"
)

// GET /static/*
func GetStatic(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetStatic(%v)", r.URL)

	// find requested file on file system and read it into a buffer

	name := path.Base(r.URL.Path)
	relpath := filepath.Join("res", name)

	content, err := ioutil.ReadFile(relpath)
	if err != nil {
		template.Error(w, r, http.StatusNotFound, err, nil)
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
