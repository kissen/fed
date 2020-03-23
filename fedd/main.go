package main

import (
	"github.com/go-fed/activity/pub"
	"github.com/gorilla/mux"
	"gitlab.cs.fau.de/kissen/fed/fedd/ap"
	"gitlab.cs.fau.de/kissen/fed/fedd/db"
	"gitlab.cs.fau.de/kissen/fed/fedd/oauth"
	"log"
	"net/http"
)

func OpenDatabase() db.FedStorage {
	// open db
	storage := &db.FedEmbeddedStorage{Filepath: Config().Store}
	if err := storage.Open(); err != nil {
		log.Fatal(err)
	}

	// add some users for testing
	for _, username := range []string{"alice", "bob", "celia"} {
		user := &db.FedUser{Name: username}
		user.SetPassword(username)
		Must(storage.StoreUser(user))
	}

	return storage
}

func InstallApHandlers(storage db.FedStorage, router *mux.Router) {
	// set up required structs that implement the various go-fed interfaces
	common := &ap.FedCommonBehavior{}
	socialProtocol := &ap.FedSocialProtocol{}
	fedProtocol := &ap.FedFederatingProtocol{}
	database := &ap.FedDatabase{}
	clock := &ap.FedClock{}

	// set up go-fed proxies
	actor := pub.NewActor(common, socialProtocol, fedProtocol, database, clock)
	handler := pub.NewActivityStreamsHandler(database, clock)

	// create http handler closures
	inboxHandler := newInboxHandler(actor, storage)
	outboxHandler := newOutboxHandler(actor, storage)
	activityHandler := newActivityHandler(handler, storage)

	// in/outbox handlers
	router.HandleFunc(`/{username:[A-Za-z]+}/inbox`, inboxHandler).Methods("GET", "POST")
	router.HandleFunc(`/{username:[A-Za-z]+}/outbox`, outboxHandler).Methods("GET", "POST")

	// everything else
	// TODO: Think about using just a catchall here? Gets in the way of other routes
	// though (admin, oauth &c).
	activityRoutes := []string{
		`/storage/{id:[A-Za-z0-9\-]+}`,
		`/{username:[A-Za-z]+}/followers`,
		`/{username:[A-Za-z]+}/following`,
		`/{username:[A-Za-z]+}/liked`,
		`/{username:[A-Za-z]+}`,
	}
	for _, route := range activityRoutes {
		router.HandleFunc(route, activityHandler).Methods("GET", "POST")
	}

	// middleware that installes the FedContext on all requests
	router.Use(InstallBaseContext(storage))
}

func InstallOAuthHandlers(oa oauth.FedOAuther, router *mux.Router) {
	router.HandleFunc("/oauth/authorize", oa.GetAuthorize).Methods("GET")
	router.HandleFunc("/oauth/authorize", oa.PostAuthorize).Methods("POST")
	router.HandleFunc("/oauth/token", oa.PostToken).Methods("POST")
}

func InstallAdminHandlers(s db.FedStorage, router *mux.Router) {
	a := &ap.FedAdminProtocol{}
	router.HandleFunc(`/{username:[A-Za-z]+}`, newAdminHandler(a, s)).Methods("PUT")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	storage := OpenDatabase()
	defer storage.Close()

	oa := oauth.New(storage)

	router := mux.NewRouter().StrictSlash(false)
	sr := router.PathPrefix(Config().Base.Path).Subrouter()

	InstallOAuthHandlers(oa, sr)
	InstallAdminHandlers(storage, sr)
	InstallApHandlers(storage, sr) // includes catchall

	addr := Config().Base.Host
	log.Printf("starting on addr=%v...", addr)
	Must(http.ListenAndServe(addr, router))
}
