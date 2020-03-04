package ap

import (
	"context"
	"gitlab.cs.fau.de/kissen/fed/fedd/db"
	"log"
	"net/http"
	"sync"
)

// Inspired by go-fed classes like Database and SocialProtocol, this
// struct contains methods for running administrator task on the
// instance.
type FedAdminProtocol struct {
	// admin requests should be rare; to make things easy for us,
	// we only alow one concurrent admin request at a time
	lock sync.Mutex
}

func (f *FedAdminProtocol) AuthenticateAdmin(c context.Context, w http.ResponseWriter, r *http.Request) (out context.Context, authed bool, err error) {
	log.Println("AuthenticateAdmin()")
	return c, true, nil
}

func (f *FedAdminProtocol) CreateUser(c context.Context, username string) (*db.FedUser, error) {
	log.Printf("CreateUser(%v)", username)

	write := db.FedUser{Name: username}
	if err := FromContext(c).Storage.StoreUser(&write); err != nil {
		return nil, err
	}

	return FromContext(c).Storage.RetrieveUser(username)
}

func (f *FedAdminProtocol) Handle(c context.Context, w http.ResponseWriter, r *http.Request) {
	log.Printf("Handle(%v)", r.URL)

	// we only handle one admin request at a time

	f.lock.Lock()
	defer f.lock.Unlock()

	// so far we only support PUT; this works nicely because Activity Pub only
	// uses GET and POST

	if r.Method != "PUT" {
		http.Error(w, "method is not PUT", http.StatusMethodNotAllowed)
		return
	}

	// check authentication

	c, authed, err := f.AuthenticateAdmin(c, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !authed {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	// do basic routing depending on IRI type

	iri := IRI{Context: c, Target: r.URL}

	if username, err := iri.Actor(); err == nil {
		if _, err := f.CreateUser(c, username); err != nil {
			// creating user failed
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			// creating user was successful
			// XXX: it would be nicer to render out the user JSON here
			http.Error(w, "Created", http.StatusCreated)
		}

		return
	}

	// don't know what to do

	http.Error(w, "Bad Admin Request", http.StatusBadRequest)
}
