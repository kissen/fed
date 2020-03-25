package main

import (
	"github.com/go-fed/activity/pub"
	"github.com/gorilla/mux"
	"gitlab.cs.fau.de/kissen/fed/ap"
	"gitlab.cs.fau.de/kissen/fed/config"
	"gitlab.cs.fau.de/kissen/fed/db"
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
	"gitlab.cs.fau.de/kissen/fed/oauth"
	"log"
	"net/http"
)

// Given the servers configuration, open a connection to the selected
// storage engine.
//
// This function panics on failure. If we can't open storage, there is
// nothing we can do.
func OpenDatabase() db.FedStorage {
	// open db
	storage := &db.FedEmbeddedStorage{Filepath: config.Get().Store}
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

// Create the go-fed handler objects that HTTP handlers can use to take
// care of ActivityPub requests.
func CreateProxies() (pub.FederatingActor, pub.HandlerFunc) {
	common := &ap.FedCommonBehavior{}
	socialProtocol := &ap.FedSocialProtocol{}
	fedProtocol := &ap.FedFederatingProtocol{}
	database := &ap.FedDatabase{}
	clock := &ap.FedClock{}

	actor := pub.NewActor(common, socialProtocol, fedProtocol, database, clock)
	handler := pub.NewActivityStreamsHandler(database, clock)

	return actor, handler
}

// Install the OAuth2 handlers. These handlers take care of authorization
// using codes and tokens are quite important for useful federation.
func InstallOAuthHandlers(oa oauth.FedOAuther, router *mux.Router) {
	router.HandleFunc("/oauth/authorize", oa.GetAuthorize).Methods("GET")
	router.HandleFunc("/oauth/authorize", oa.PostAuthorize).Methods("POST")
	router.HandleFunc("/oauth/token", oa.PostToken).Methods("POST")
}

// Install the admin handlers. The idea is that we write some admin tool
// that runs on the same machine as fed itself interacts with these admin handlers.
// I think that's easier than bothering with a web gui for configuration for now.
func InstallAdminHandlers(s db.FedStorage, router *mux.Router) {
	a := &ap.FedAdminProtocol{}
	router.HandleFunc(`/{username:[A-Za-z]+}`, newAdminHandler(a, s)).Methods("PUT")
}

// Install the actually interesting handlers. These handlers will differentiate
// between Content-Type/Accept headers and either send out JSON for ActivityPub
// or a gaudy web interface instead.
func InstallApHandlers(storage db.FedStorage, router *mux.Router) {
	router.HandleFunc("/{username:[A-Za-z]+}/outbox", OutboxHandler).Methods("GET", "POST")
	router.HandleFunc("/{username:[A-Za-z]+}/inbox", InboxHandler).Methods("GET", "POST")

	// catchall
	router.PathPrefix("/").HandlerFunc(ActivityHandler).Methods("GET", "POST")
}

// Install the different error handler. While the defaults from gorilla are
// reasonable, we can be more specific.
func InstallErrorHandlers(router *mux.Router) {
	// TODO: differentiate by Accept header
	router.NotFoundHandler = router.NewRoute().HandlerFunc(ApiNotFound).GetHandler()
	router.MethodNotAllowedHandler = router.NewRoute().HandlerFunc(ApiNotAllowed).GetHandler()
}

// Install middleware that runs before every single actual HTTP handler.
func InstallMiddleware(storage db.FedStorage, router *mux.Router) {
	pa, hf := CreateProxies()
	router.Use(fedcontext.AddContext(storage, pa, hf))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	storage := OpenDatabase()
	defer storage.Close()

	oa := oauth.New(storage)

	router := mux.NewRouter().StrictSlash(false)
	sr := router.PathPrefix(config.Get().Base.Path).Subrouter()

	InstallOAuthHandlers(oa, sr)
	InstallAdminHandlers(storage, sr)
	InstallApHandlers(storage, sr) // includes catchall
	InstallErrorHandlers(router)
	InstallMiddleware(storage, router)

	addr := config.Get().Base.Host
	log.Printf("starting on addr=%v...", addr)
	Must(http.ListenAndServe(addr, router))
}
