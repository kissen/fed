package main

import (
	"context"
	"log"
	"errors"
	"fmt"
	"github.com/go-fed/activity/pub"
	"io"
	"net/http"
)

// Write out the error to an HTTP connection
func replyWithError(w http.ResponseWriter, statusCode int, cause error) {
	// https://golang.org/pkg/net/http/#ResponseWriter

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	statusStr := http.StatusText(statusCode)
	errStr := cause.Error()
	content := fmt.Sprintf("%d %s: %s\n", statusCode, statusStr, errStr)

	io.WriteString(w, content)
}

// Handlers adapted from tutorial on https://go-fed.org/tutorial#ActivityStreams-Types-and-Properties

func newOutboxHandler(pubber pub.Pubber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := context.Background()

		// Populate c with application specific information here
		// ...

		// Try POST w/ Activity Pub

		handled, err := pubber.PostOutbox(c, w, r)

		if err != nil {
			replyWithError(w, http.StatusInternalServerError, err)
			return
		}

		if handled {
			return
		}

		// Try GET w/ Activity Pub

		handled, err = pubber.GetOutbox(c, w, r)

		if err != nil {
			replyWithError(w, http.StatusInternalServerError, err)
			return
		}

		if handled {
			return
		}

		// Handle non-ActivityPub request, such as responding with a HTML
		// representation with correct view permissions.

		replyWithError(w, http.StatusNotImplemented, errors.New("only ActivityPub may get the outbox"))
	}
}

func newInboxHandler(pubber pub.Pubber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := context.Background()

		// Populate c with application specific information here
		// ...

		// Try POST

		handled, err := pubber.PostInbox(c, w, r)

		if err != nil {
			replyWithError(w, http.StatusInternalServerError, err)
			return
		}

		if handled {
			return
		}

		// Try GET

		handled, err = pubber.GetInbox(c, w, r)

		if err != nil {
			replyWithError(w, http.StatusInternalServerError, err)
			return
		}

		if handled {
			return
		}

		// Handle non-ActivityPub request, such as responding with a HTML
		// representation with correct view permissions.

		replyWithError(w, http.StatusNotImplemented, errors.New("only ActivityPub may get the inbox"))
	}
}

func newActivityHandler(handler pub.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := context.Background()

		// Populate c with application specific information here
		// ...

		log.Println(r.URL)

		handled, err := handler(c, w, r) // or myCustomVerifiedPubHandler

		if err != nil {
			replyWithError(w, http.StatusInternalServerError, err)
			return
		}

		if handled {
			return
		}

		// Handle non-ActivityPub request, such as responding with a HTML
		// representation with correct view permissions.

		replyWithError(w, http.StatusNotImplemented, errors.New("only ActivityPub may return activties"))
	}
}
