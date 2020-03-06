package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func listenAndServe() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", GetIndex).Methods("GET")
	router.HandleFunc("/stream", GetStream).Methods("GET")
	router.HandleFunc("/liked", GetLiked).Methods("GET")
	router.HandleFunc("/following", GetFollowing).Methods("GET")
	router.HandleFunc("/followers", GetFollowers).Methods("GET")
	router.HandleFunc("/login", GetLogin).Methods("GET")
	router.HandleFunc("/login", PostLogin).Methods("POST")
	router.HandleFunc("/remote/{remotepath:.+}", GetRemote).Methods("GET")
	router.HandleFunc("/static/{.+}", GetStatic).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(HandleNotFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(HandleMethodNotAllowed)

	addr := ":8080"
	log.Printf("starting on addr=%v...", addr)
	Must(http.ListenAndServe(addr, router))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	listenAndServe()
}
