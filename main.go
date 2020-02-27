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
	r := mux.NewRouter().StrictSlash(true)

	// Activity Pub

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

	r.HandleFunc(`/ap/{username:[A-Za-z]+}`, newActivityHandler(handler, storage)).Methods("GET", "POST")
	r.HandleFunc(`/ap/{username:[A-Za-z]+}/inbox`, newInboxHandler(actor, storage)).Methods("GET", "POST")
	r.HandleFunc(`/ap/{username:[A-Za-z]+}/outbox`, newOutboxHandler(actor, storage)).Methods("GET", "POST")
	r.HandleFunc(`/ap/{username:[A-Za-z]+}/followers`, newActivityHandler(handler, storage)).Methods("GET", "POST")
	r.HandleFunc(`/ap/{username:[A-Za-z]+}/following`, newActivityHandler(handler, storage)).Methods("GET", "POST")
	r.HandleFunc(`/ap/{username:[A-Za-z]+}/liked`, newActivityHandler(handler, storage)).Methods("GET", "POST")
	r.HandleFunc(`/ap/storage/{id:[A-Za-z0-9\-]}`, newActivityHandler(handler, storage)).Methods("GET", "POST")

	// Let's rock!

	addr := ":9999"

	log.Printf("Starting on addr=%v...", addr)
	err := http.ListenAndServe(addr, r)
	log.Fatal(err)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	storage := openDatabase()
	defer storage.Close()

	listenAndAccept(storage)
}
