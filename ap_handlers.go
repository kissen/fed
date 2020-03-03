package main

import (
	"context"
	"github.com/go-fed/activity/pub"
	"gitlab.cs.fau.de/kissen/fed/ap"
	"gitlab.cs.fau.de/kissen/fed/db"
	"log"
	"net/http"
)

func baseContext(store db.FedStorage) context.Context {
	ctx := ap.WithFedContext(context.Background())
	fc := ap.FromContext(ctx)

	fc.Scheme = ap.Just("http")
	fc.Host = ap.Just("localhost:9999")
	fc.BasePath = ap.Just("/ap/")

	fc.Storage = store

	return ctx
}

// Handlers adapted from tutorial on https://go-fed.org/tutorial#ActivityStreams-Types-and-Properties

func newOutboxHandler(actor pub.Actor, store db.FedStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("OutboxHandler(%v)", r.URL)

		// populate c with application specific information here
		c := baseContext(store)

		// try POST w/ Activity Pub
		if handled, err := actor.PostOutbox(c, w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if handled {
			return
		}

		// try GET w/ Activity Pub
		if handled, err := actor.GetOutbox(c, w, r); err != nil {
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

		// populate c with application specific information here
		c := baseContext(store)

		// try POST
		if handled, err := actor.PostInbox(c, w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if handled {
			return
		}

		// try GET
		if handled, err := actor.GetInbox(c, w, r); err != nil {
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

		// populate c with application specific information here
		c := baseContext(store)

		// or myCustomVerifiedPubHandler
		if handled, err := handler(c, w, r); err != nil {
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

func newAdminHandler(admin *ap.FedAdminProtocol, store db.FedStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AdminHandler(%v)", r.URL)

		// populate c with application specific information here
		c := baseContext(store)

		// AdminProtocol takes care of everything :)
		admin.Handle(c, w, r)
	}
}
