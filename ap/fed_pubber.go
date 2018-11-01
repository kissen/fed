package ap

import (
	"context"
	"net/http"
)

// Implements the github.com/go-fed/activity/pub/Pubber interface
type FedPubber struct {
}

// PostInbox returns true if the request was handled as an ActivityPub
// POST to an actor's inbox. If false, the request was not an
// ActivityPub request.
//
// If the error is nil, then the ResponseWriter's headers and response
// has already been written. If a non-nil error is returned, then no
// response has been written.
func (f *FedPubber) PostInbox(c context.Context, w http.ResponseWriter, r *http.Request) (bool, error) {
	return false, nil
}

// GetInbox returns true if the request was handled as an ActivityPub
// GET to an actor's inbox. If false, the request was not an ActivityPub
// request.
//
// If the error is nil, then the ResponseWriter's headers and response
// has already been written. If a non-nil error is returned, then no
// response has been written.
func (f *FedPubber) GetInbox(c context.Context, w http.ResponseWriter, r *http.Request) (bool, error) {
	return false, nil
}

// PostOutbox returns true if the request was handled as an ActivityPub
// POST to an actor's outbox. If false, the request was not an
// ActivityPub request.
//
// If the error is nil, then the ResponseWriter's headers and response
// has already been written. If a non-nil error is returned, then no
// response has been written.
func (f *FedPubber) PostOutbox(c context.Context, w http.ResponseWriter, r *http.Request) (bool, error) {
	return false, nil
}

// GetOutbox returns true if the request was handled as an ActivityPub
// GET to an actor's outbox. If false, the request was not an
// ActivityPub request.
//
// If the error is nil, then the ResponseWriter's headers and response
// has already been written. If a non-nil error is returned, then no
// response has been written.
func (f *FedPubber) GetOutbox(c context.Context, w http.ResponseWriter, r *http.Request) (bool, error) {
	return false, nil
}
