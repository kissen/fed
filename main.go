package main

import (
	"github.com/go-fed/activity/pub"
	"github.com/gorilla/mux"
	"gitlab.cs.fau.de/kissen/fed/db"
	"gitlab.cs.fau.de/kissen/fed/ap"
	"log"
	"net/http"
)

func listenAndAccept(storage db.FedStorer) {
	r := mux.NewRouter().StrictSlash(true)

	// Custom Shit

	r.HandleFunc("/custom/", makeHandleRootGet(storage)).Methods("GET")
	r.HandleFunc("/custom/inst/{userId:[0-9]+}/{postId:[0-9]+}", makeHandleUserPostGet(storage)).Methods("GET")
	r.HandleFunc("/custom/inst/{userId:[0-9]+}", makeHandleUserPostPut(storage)).Methods("PUT")

	// Activity Pub

	application := &ap.FedApplication{
		BaseIRI: "/ap/",
		Storage: storage,
	}

	clock := &ap.FedClock{}

	socialCallbacker := &ap.FedSocialCallbacker{}

	socialPubber := pub.NewSocialPubber(clock, application, socialCallbacker)
	activityHandler := pub.ServeActivityPubObject(application, clock)

	r.HandleFunc("/ap/inbox/{username:[A-Za-z]+}", newInboxHandler(socialPubber)).Methods("GET", "POST")
	r.HandleFunc("/ap/outbox/{username:[A-Za-z]+}", newOutboxHandler(socialPubber)).Methods("GET", "POST")
	r.HandleFunc("/ap/activity/{id:[0-9]+}", newActivityHandler(activityHandler)).Methods("GET", "POST")

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
	storage := db.NewFedRamStorage()
	john := storage.AddUser("John")
	storage.AddPost(john.Id, "Hallo, Welt!")

	// start http server
	listenAndAccept(storage)
}
