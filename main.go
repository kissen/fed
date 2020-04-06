package main

import (
	"github.com/go-fed/activity/pub"
	"github.com/gorilla/mux"
	"gitlab.cs.fau.de/kissen/fed/ap"
	"gitlab.cs.fau.de/kissen/fed/config"
	"gitlab.cs.fau.de/kissen/fed/db"
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
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
func InstallOAuthHandlers(router *mux.Router) {
	router.HandleFunc("/oauth/authorize", GetOAuthAuthorize).Methods("GET")
	router.HandleFunc("/oauth/authorize", PostOAuthAuthorize).Methods("POST")
	router.HandleFunc("/oauth/token", PostOAuthToken).Methods("POST")
}

// Install the admin handlers. The idea is that we write some admin tool
// that runs on the same machine as fed itself interacts with these admin handlers.
// I think that's easier than bothering with a web gui for configuration for now.
func InstallAdminHandlers(router *mux.Router) {
	router.HandleFunc(`/{username:[A-Za-z]+}`, PutUser).Methods("PUT")
}

// Install the handlers for all /.well-known services.
func InstallWellKnownHandlers(router *mux.Router) {
	router.HandleFunc("/.well-known/nodeinfo", GetNodeInfo).Methods("GET")
	router.HandleFunc("/.well-known/nodeinfo/2.0.json", GetNodeInfo20).Methods("GET")
	router.HandleFunc("/.well-known/webfinger", GetWebfinger).Methods("GET")
	router.HandleFunc("/.well-known/host-meta", GetHostMeta).Methods("GET")
}

// Install the actually interesting handlers. These handlers will differentiate
// between Content-Type/Accept headers and either send out JSON for ActivityPub
// or a gaudy web interface instead.
func InstallSplitHandlers(router *mux.Router) {
	// this is kind of a mess
	router.HandleFunc("/{username:[A-Za-z]+}/outbox", GetPostOutbox).Methods("GET", "POST")
	router.HandleFunc("/{username:[A-Za-z]+}/inbox", GetPostInbox).Methods("GET", "POST")
	router.HandleFunc("/", GetIndex).Methods("GET")
	router.HandleFunc("/stream", WebGetStream).Methods("GET")
	router.HandleFunc("/liked", GetLiked).Methods("GET")
	router.HandleFunc("/following", GetFollowing).Methods("GET")
	router.HandleFunc("/followers", GetFollowers).Methods("GET")
	router.HandleFunc("/login", GetLogin).Methods("GET")
	router.HandleFunc("/login", PostLogin).Methods("POST")
	router.HandleFunc("/logout", PostLogout).Methods("POST")
	router.HandleFunc("/remote/{remotepath:.+}", GetRemote).Methods("GET")
	router.HandleFunc("/static/{.+}", GetStatic).Methods("GET")
	router.HandleFunc("/submit", PostSubmit).Methods("POST")
	router.HandleFunc("/reply", PostReply).Methods("POST")
	router.HandleFunc("/repeat", PostRepeat).Methods("POST")
	router.HandleFunc("/like", PostLike).Methods("POST")

	// catchall for activity pub
	//
	// TODO: remove this; in general, be more specific and get rid of that
	// stupid file split_handlers.go!
	router.PathPrefix("/").HandlerFunc(ApGetPostActivity).Methods("GET", "POST")
}

// Install the different error handler. While the defaults from gorilla are
// reasonable, we can be more specific.
func InstallErrorHandlers(router *mux.Router) {
	router.NotFoundHandler = router.NewRoute().HandlerFunc(NotFound).GetHandler()
	router.MethodNotAllowedHandler = router.NewRoute().HandlerFunc(MethodNotAllowed).GetHandler()
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

	router := mux.NewRouter().StrictSlash(false)
	sr := router.PathPrefix(config.Get().Base.Path).Subrouter()

	InstallOAuthHandlers(sr)
	InstallWellKnownHandlers(sr)
	InstallAdminHandlers(sr)
	InstallSplitHandlers(sr) // includes catchall

	InstallErrorHandlers(router)
	InstallMiddleware(storage, router)

	addr := "[::]:4040"
	log.Printf("starting on addr=%v...", addr)
	Must(http.ListenAndServe(addr, router))
}
