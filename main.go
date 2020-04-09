package main

import (
	"github.com/go-fed/activity/pub"
	"github.com/gorilla/mux"
	"gitlab.cs.fau.de/kissen/fed/ap"
	"gitlab.cs.fau.de/kissen/fed/config"
	"gitlab.cs.fau.de/kissen/fed/db"
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
	"gitlab.cs.fau.de/kissen/fed/util"
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
	storage := &db.FedEmbeddedStorage{
		Filepath: config.Get().StorageFile,
	}

	if err := storage.Open(); err != nil {
		log.Fatal(err)
	}

	// add some users for testing
	for _, username := range []string{"alice", "bob", "celia"} {
		user := &db.FedUser{Name: username}
		user.SetPassword(username)
		util.Must(storage.StoreUser(user))
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

// Install the admin handlers. The idea is that we write some admin tool
// that runs on the same machine as fed itself interacts with these admin handlers.
// I think that's easier than bothering with a web gui for configuration for now.
func InstallAdminHandlers(router *mux.Router) {
	router.HandleFunc(`/{username:[A-Za-z]+}`, PutUser).Methods("PUT")
}

// Install the OAuth2 handlers. These handlers take care of authorization
// using codes and tokens are quite important for useful federation.
func InstallOAuthHandlers(router *mux.Router) {
	router.HandleFunc("/oauth/authorize", GetOAuthAuthorize).Methods("GET")
	router.HandleFunc("/oauth/authorize", PostOAuthAuthorize).Methods("POST")
	router.HandleFunc("/oauth/token", PostOAuthToken).Methods("POST")
}

// Install the handlers for all /.well-known/ targets. These are used by
// other software on the fediverse to look up stuff about actors on our
// instance.
func InstallWellKnownHandlers(router *mux.Router) {
	router.HandleFunc("/.well-known/nodeinfo", GetNodeInfo).Methods("GET")
	router.HandleFunc("/.well-known/nodeinfo/2.0.json", GetNodeInfo20).Methods("GET")
	router.HandleFunc("/.well-known/webfinger", GetWebfinger).Methods("GET")
	router.HandleFunc("/.well-known/host-meta", GetHostMeta).Methods("GET")
}

// Install handlers that are really just workaround and redirects
// to deal with other software on the fediverse.
func InstallShimHandlers(router *mux.Router) {
	router.HandleFunc("/shim/ostatus_subscribe", GetOStatusSubscribe).Methods("GET")
}

// Install the handlers that take care of handling requests to Activity
// Pub endpoints.
func InstallApHandlers(router *mux.Router) {
	InstallApHandler(router, ApGetPostOutbox, "/{username:[A-Za-z]+}/outbox")
	InstallApHandler(router, ApGetPostInbox, "/{username:[A-Za-z]+}/inbox")

	InstallApHandler(router, ApGetPostActivity, "/{username:[A-Za-z]+}") // actor endpoint
	InstallApHandler(router, ApGetPostActivity, "/{username:[A-Za-z]+}/following")
	InstallApHandler(router, ApGetPostActivity, "/{username:[A-Za-z]+}/followers")
	InstallApHandler(router, ApGetPostActivity, "/{username:[A-Za-z]+}/liked")
	InstallApHandler(router, ApGetPostActivity, "/storage/{uuid:.+}")
}

// Install activity pub handler h for pattern. This function takes care of
// registering the handler with the correct method and Accept/Content-Type header.
func InstallApHandler(target *mux.Router, h http.HandlerFunc, pattern string) {
	target.HandleFunc(pattern, h).Methods("GET").Headers("Accept", util.AP_TYPE)
	target.HandleFunc(pattern, h).Methods("POST").Headers("Content-Type", util.AP_TYPE)
}

// Install the different error handler. While the defaults from gorilla are
// reasonable, we can be more specific.
func InstallErrorHandlers(router *mux.Router) {
	router.NotFoundHandler = router.NewRoute().HandlerFunc(NotFound).GetHandler()
	router.MethodNotAllowedHandler = router.NewRoute().HandlerFunc(MethodNotAllowed).GetHandler()
}

// Install middleware that runs before every single actual HTTP handler.
func InstallMiddleware(storage db.FedStorage, router *mux.Router) {
	// middleware that signs all responses
	router.Use(SignResponseMiddleware)

	// middleware that installs a FedContext on all requests;
	// it's nicer than dealing with global variables
	pa, hf := CreateProxies()
	router.Use(fedcontext.AddContext(storage, pa, hf))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	storage := OpenDatabase()
	defer storage.Close()

	router := mux.NewRouter().StrictSlash(false)

	InstallAdminHandlers(router)
	InstallOAuthHandlers(router)
	InstallWellKnownHandlers(router)
	InstallShimHandlers(router)
	InstallApHandlers(router)

	InstallErrorHandlers(router)
	InstallMiddleware(storage, router)

	addr := config.Get().ListenAddress
	log.Printf("listening on addr=%v...", addr)
	util.Must(http.ListenAndServe(addr, router))
}
