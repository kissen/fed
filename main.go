package main

import (
	"github.com/go-fed/activity/pub"
	"github.com/gorilla/mux"
	"gitlab.cs.fau.de/kissen/fed/ap"
	"gitlab.cs.fau.de/kissen/fed/db"
	"log"
	"net/http"
)

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

	r.HandleFunc("/ap/inbox/{username:[A-Za-z]+}", newInboxHandler(actor, storage)).Methods("GET", "POST")
	r.HandleFunc("/ap/outbox/{username:[A-Za-z]+}", newOutboxHandler(actor, storage)).Methods("GET", "POST")
	r.HandleFunc("/ap/activity/{id:[0-9]+}", newActivityHandler(handler, storage)).Methods("GET", "POST")

	// Let's rock!

	addr := ":9999"

	log.Printf("Starting on addr=%v...", addr)
	err := http.ListenAndServe(addr, r)
	log.Fatal(err)
}

func main() {
	// configure logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// set up database
	storage := db.FedEmbeddedStorage{
		Filepath: "/tmp/fed.bbolt",
	}

	// start http server
	listenAndAccept(storage)
}
