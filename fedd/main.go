package main

import (
	"github.com/go-fed/activity/pub"
	"github.com/gorilla/mux"
	"gitlab.cs.fau.de/kissen/fed/fedd/ap"
	"gitlab.cs.fau.de/kissen/fed/fedd/db"
	"log"
	"net/http"
)

func openDatabase() db.FedStorage {
	// open db

	storage := &db.FedEmbeddedStorage{
		Filepath: Config().Store,
	}

	if err := storage.Open(); err != nil {
		log.Fatal(err)
	}

	// add some users

	usernames := []string{
		"alice", "bob", "celia", "daniel", "emily", "frank",
	}

	for _, username := range usernames {
		user := &db.FedUser{Name: username}
		Must(storage.StoreUser(user))
	}

	// done

	return storage
}

func listenAndAccept(storage db.FedStorage) {
	// set up ap/ structs

	common := &ap.FedCommonBehavior{}
	socialProtocol := &ap.FedSocialProtocol{}
	fedProtocol := &ap.FedFederatingProtocol{}
	database := &ap.FedDatabase{}
	clock := &ap.FedClock{}

	actor := pub.NewActor(
		common, socialProtocol, fedProtocol, database, clock,
	)

	handler := pub.NewActivityStreamsHandler(
		database, clock,
	)

	admin := &ap.FedAdminProtocol{}

	// build up http handlers

	inboxHandler := newInboxHandler(actor, storage)
	outboxHandler := newOutboxHandler(actor, storage)
	activityHandler := newActivityHandler(handler, storage)
	adminHandler := newAdminHandler(admin, storage)

	// configure routes

	router := mux.NewRouter().StrictSlash(false)
	sr := router.PathPrefix(Config().Base.Path).Subrouter()

	sr.HandleFunc(`/{username:[A-Za-z]+}/inbox`, inboxHandler).Methods("GET", "POST")
	sr.HandleFunc(`/{username:[A-Za-z]+}/outbox`, outboxHandler).Methods("GET", "POST")

	activityRoutes := []string{
		`/{username:[A-Za-z]+}`, `/{username:[A-Za-z]+}/followers`,
		`/{username:[A-Za-z]+}/following`, `/{username:[A-Za-z]+}/liked`,
		`/storage/{id:[A-Za-z0-9\-]+}`,
	}

	for _, route := range activityRoutes {
		sr.HandleFunc(route, activityHandler).Methods("GET", "POST")
	}

	// build up admin routes

	sr.HandleFunc(`/{username:[A-Za-z]+}`, adminHandler).Methods("PUT")

	// install midleware

	sr.Use(InstallBaseContext(storage))

	// let's rock!

	addr := Config().Base.Host
	log.Printf("starting on addr=%v...", addr)
	Must(http.ListenAndServe(addr, router))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	storage := openDatabase()
	defer storage.Close()

	listenAndAccept(storage)
}
