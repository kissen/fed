package main

import (
	"github.com/go-fed/activity/pub"
	"gitlab.cs.fau.de/kissen/fed/fedd/db"
	"log"
	"net/http"
)

// Handlers adapted from tutorial on https://go-fed.org/tutorial#ActivityStreams-Types-and-Properties

func newOutboxHandler(actor pub.Actor, store db.FedStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("OutboxHandler(%v)", r.URL)

		// try POST w/ Activity Pub
		if handled, err := actor.PostOutbox(r.Context(), w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if handled {
			return
		}

		// try GET w/ Activity Pub
		if handled, err := actor.GetOutbox(r.Context(), w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if handled {
			return
		}

		// handle non-ActivityPub request, such as responding with a HTML
		// representation with correct view permissions.
		http.Error(
			w,
			"not acceptable; check https://www.w3.org/TR/activitypub/#retrieving-objects",
			http.StatusNotAcceptable,
		)
	}
}

func newInboxHandler(actor pub.Actor, store db.FedStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("InboxHandler(%v)", r.URL)

		// try POST
		if handled, err := actor.PostInbox(r.Context(), w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if handled {
			return
		}

		// try GET
		if handled, err := actor.GetInbox(r.Context(), w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if handled {
			return
		}

		// handle non-ActivityPub request, such as responding with a HTML
		// representation with correct view permissions.
		http.Error(
			w,
			"not acceptable; check https://www.w3.org/TR/activitypub/#retrieving-objects",
			http.StatusNotAcceptable,
		)
	}
}

func newActivityHandler(handler pub.HandlerFunc, store db.FedStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("ActivityHandler(%v)", r.URL)

		// or myCustomVerifiedPubHandler
		if handled, err := handler(r.Context(), w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if handled {
			return
		}

		// handle non-ActivityPub request, such as responding with a HTML
		// representation with correct view permissions.
		http.Error(
			w,
			"not acceptable; check https://www.w3.org/TR/activitypub/#retrieving-objects",
			http.StatusNotAcceptable,
		)
	}
}
