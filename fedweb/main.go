package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func listenAndServe() {
	router := mux.NewRouter().StrictSlash(false)

	router.HandleFunc("/", GetIndex).Methods("GET")
	router.HandleFunc("/stream", GetStream).Methods("GET")
	router.HandleFunc("/liked", GetLiked).Methods("GET")
	router.HandleFunc("/following", GetFollowing).Methods("GET")
	router.HandleFunc("/followers", GetFollowers).Methods("GET")
	router.HandleFunc("/login", GetLogin).Methods("GET")
	router.HandleFunc("/login", PostLogin).Methods("POST")
	router.HandleFunc("/logout", PostLogout).Methods("POST")
	router.HandleFunc("/remote/{remotepath:.+}", GetRemote).Methods("GET")
	router.HandleFunc("/static/{.+}", GetStatic).Methods("GET")
	router.HandleFunc("/submit", PostSubmit).Methods("POST")

	// https://github.com/gorilla/mux/issues/416
	router.NotFoundHandler = router.NewRoute().HandlerFunc(HandleNotFound).GetHandler()
	router.MethodNotAllowedHandler = router.NewRoute().HandlerFunc(HandleMethodNotAllowed).GetHandler()

	router.Use(AddContext)

	addr := Config().Base.Host
	log.Printf("starting on addr=%v...", addr)
	Must(http.ListenAndServe(addr, router))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	listenAndServe()
}
