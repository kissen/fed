package main

import (
	"github.com/go-fed/activity/pub"
	"github.com/gorilla/mux"
	"gitlab.cs.fau.de/kissen/fed/ap"
	"gitlab.cs.fau.de/kissen/fed/db"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func openDatabase() db.FedStorage {
	// open db

	dbPath := filepath.Join(os.TempDir(), "main.db")

	storage := &db.FedEmbeddedStorage{
		Filepath: dbPath,
	}

	storage.Open()

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

	// build up go-fed objects

	inboxHandler := newInboxHandler(actor, storage)
	outboxHandler := newOutboxHandler(actor, storage)
	activityHandler := newActivityHandler(handler, storage)

	// set up routes

	router := mux.NewRouter().StrictSlash(false)

	router.HandleFunc(`/ap/{username:[A-Za-z]+}/inbox`, inboxHandler).Methods("GET", "POST")
	router.HandleFunc(`/ap/{username:[A-Za-z]+}/outbox`, outboxHandler).Methods("GET", "POST")

	activityRoutes := []string{
		`/ap/{username:[A-Za-z]+}`, `/ap/{username:[A-Za-z]+}/followers`,
		`/ap/{username:[A-Za-z]+}/following`, `/ap/{username:[A-Za-z]+}/liked`,
		`/ap/storage/{id:[A-Za-z0-9\-]+}`,
	}

	for _, route := range activityRoutes {
		router.HandleFunc(route, activityHandler).Methods("GET", "POST")
	}

	// let's rock!

	addr := ":9999"

	log.Printf("Starting on addr=%v...", addr)

	err := http.ListenAndServe(addr, router)
	log.Fatal(err)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	storage := openDatabase()
	defer storage.Close()

	listenAndAccept(storage)
}
