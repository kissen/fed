package main

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
	"log"
	"net/http"
)

// Functions provided by the pub.FederatingActor provide methods that
// match this signature. These methods take care of handling inbox and
// outbox requests for ActivityPub.
type BoxHandler func(context.Context, http.ResponseWriter, *http.Request) (bool, error)

func OutboxHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("OutboxHandler(%v)", r.URL)

	if done, err := HandleOutboxWithPubActor(w, r); err != nil {
		ApiError(w, r, "go-fed error", err, http.StatusInternalServerError)
		return
	} else if done {
		return
	}

	ApiError(w, r, "check https://www.w3.org/TR/activitypub/#retrieving-objects", nil, http.StatusNotAcceptable)
}

func InboxHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("InboxHandler(%v)", r.URL)

	if done, err := HandleInboxWithPubActor(w, r); err != nil {
		ApiError(w, r, "go-fed error", err, http.StatusInternalServerError)
		return
	} else if done {
		return
	}

	ApiError(w, r, "check https://www.w3.org/TR/activitypub/#retrieving-objects", nil, http.StatusNotAcceptable)
}

func ActivityHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ActivityHandler(%v)", r.URL)

	if done, err := HandleWithHandleFunc(w, r); err != nil {
		ApiError(w, r, "go-fed error", err, http.StatusInternalServerError)
		return
	} else if done {
		return
	}

	ApiError(w, r, "check https://www.w3.org/TR/activitypub/#retrieving-objects", nil, http.StatusNotAcceptable)
}

// Try to handle the given HTTP request for an outbox with go-fed.
func HandleOutboxWithPubActor(w http.ResponseWriter, r *http.Request) (handled bool, err error) {
	pa := fedcontext.Context(r).PubActor
	bhs := []BoxHandler{pa.PostOutbox, pa.GetOutbox}
	return HandleWithPubActor(w, r, bhs...)
}

// Try to handle the given HTTP request for an inbox with go-fed.
func HandleInboxWithPubActor(w http.ResponseWriter, r *http.Request) (handled bool, err error) {
	pa := fedcontext.Context(r).PubActor
	bhs := []BoxHandler{pa.PostInbox, pa.GetInbox}
	return HandleWithPubActor(w, r, bhs...)
}

// Given box handlers, try out all of them. If any of them succeeds in
// handling this request, return (true, nil). If an error gets
// returned by any function in bhs, return that error. If none of the
// handlers bhs took care of the request, but no errors occured, return
// (false, nil).
func HandleWithPubActor(w http.ResponseWriter, r *http.Request, bhs ...BoxHandler) (handled bool, err error) {
	for _, bh := range bhs {
		handled, err = bh(r.Context(), w, r)
		if err != nil {
			return false, errors.Wrap(err, "error in go-fed actor")
		}
		if handled {
			return true, nil
		}
	}

	return false, nil
}

// Try to handle the given HTTP request for some objects with go-fed.
func HandleWithHandleFunc(w http.ResponseWriter, r *http.Request) (handled bool, err error) {
	ph := fedcontext.Context(r).PubHandler

	if handled, err = ph(r.Context(), w, r); err != nil {
		return false, errors.Wrap(err, "error in go-fed handle func")
	}

	return handled, nil
}
