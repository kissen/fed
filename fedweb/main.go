package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func listenAndServe() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/static/{.+}", GetStatic).Methods("GET")
	router.HandleFunc("/", GetStream).Methods("GET")
	router.HandleFunc("/login", GetLogin).Methods("GET")
	router.HandleFunc("/login", PostLogin).Methods("POST")

	addr := ":8080"
	log.Printf("starting on addr=%v...", addr)
	Must(http.ListenAndServe(addr, router))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	listenAndServe()
}
