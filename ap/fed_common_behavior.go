package ap

import (
	"context"
	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams"
	"gitlab.cs.fau.de/kissen/fed/db"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"net/url"
)

// Implements the go-fed/activity/pub/Datbase interface (version 1.0)
type FedCommonBehavior struct{}

// AuthenticateGetInbox delegates the authentication of a GET to an
// inbox.
//
// Always called, regardless whether the Federated Protocol or Social
// API is enabled.
//
// If an error is returned, it is passed back to the caller of
// GetInbox. In this case, the implementation must not write a
// response to the ResponseWriter as is expected that the client will
// do so when handling the error. The 'authenticated' is ignored.
//
// If no error is returned, but authentication or authorization fails,
// then authenticated must be false and error nil. It is expected that
// the implementation handles writing to the ResponseWriter in this
// case.
//
// Finally, if the authentication and authorization succeeds, then
// authenticated must be true and error nil. The request will continue
// to be processed.
func (f *FedCommonBehavior) AuthenticateGetInbox(c context.Context, w http.ResponseWriter, r *http.Request) (out context.Context, authed bool, err error) {
	log.Println("AuthenticateGetInbox()")

	if username, _, ok := r.BasicAuth(); ok {
		FromContext(c).Username = Just(username)
	}

	return c, true, nil
}

// AuthenticateGetOutbox delegates the authentication of a GET to an
// outbox.
//
// Always called, regardless whether the Federated Protocol or Social
// API is enabled.
//
// If an error is returned, it is passed back to the caller of
// GetOutbox. In this case, the implementation must not write a
// response to the ResponseWriter as is expected that the client will
// do so when handling the error. The 'authenticated' is ignored.
//
// If no error is returned, but authentication or authorization fails,
// then authenticated must be false and error nil. It is expected that
// the implementation handles writing to the ResponseWriter in this
// case.
//
// Finally, if the authentication and authorization succeeds, then
// authenticated must be true and error nil. The request will continue
// to be processed.
func (f *FedCommonBehavior) AuthenticateGetOutbox(c context.Context, w http.ResponseWriter, r *http.Request) (out context.Context, authed bool, err error) {
	log.Println("AuthenticateGetOutbox()")

	if username, _, ok := r.BasicAuth(); ok {
		FromContext(c).Username = Just(username)
	}

	return c, true, nil
}

func (f *FedCommonBehavior) userForOutbox(c context.Context, r *http.Request) (user *db.FedUser, err error) {
	var username string
	var iri *url.URL = r.URL

	if username, err = parseOutboxOwnerFromIri(c, iri); err != nil {
		return nil, errors.Wrapf(err, "cannot determine owner of iri=%v", iri)
	}

	if user, err = FromContext(c).Storage.RetrieveUser(username); err != nil {
		return nil, errors.Wrapf(err, "no user found for username=%v", username)
	}

	return user, err
}

// GetOutbox returns the OrderedCollection inbox of the actor for this
// context. It is up to the implementation to provide the correct
// collection for the kind of authorization given in the request.
//
// AuthenticateGetOutbox will be called prior to this.
//
// Always called, regardless whether the Federated Protocol or Social
// API is enabled.
func (f *FedCommonBehavior) GetOutbox(c context.Context, r *http.Request) (page vocab.ActivityStreamsOrderedCollectionPage, err error) {
	log.Println("GetOutbox()")

	// fetch user meta data

	var user *db.FedUser

	if user, err = f.userForOutbox(c, r); err != nil {
		return nil, err
	}

	// build up collection

	collection := streams.NewActivityStreamsOrderedItemsProperty()

	for _, iri := range user.Outbox {
		if obj, err := FromContext(c).Storage.RetrieveObject(iri); err != nil {
			return nil, errors.Wrapf(err, "missing iri=%v in database", iri)
		} else if err := collection.AppendType(obj); err != nil {
			return nil, errors.Wrapf(err, "cannot add iri=%v to collection", iri)
		}
	}

	// send out reply

	inbox := streams.NewActivityStreamsOrderedCollectionPage()
	inbox.SetActivityStreamsOrderedItems(collection)

	return inbox, nil
}

// NewTransport returns a new Transport on behalf of a specific actor.
//
// The actorBoxIRI will be either the inbox or outbox of an actor who is
// attempting to do the dereferencing or delivery. Any authentication
// scheme applied on the request must be based on this actor. The
// request must contain some sort of credential of the user, such as a
// HTTP Signature.
//
// The gofedAgent passed in should be used by the Transport
// implementation in the User-Agent, as well as the application-specific
// user agent string. The gofedAgent will indicate this library's use as
// well as the library's version number.
//
// Any server-wide rate-limiting that needs to occur should happen in a
// Transport implementation. This factory function allows this to be
// created, so peer servers are not DOS'd.
//
// Any retry logic should also be handled by the Transport
// implementation.
//
// Note that the library will not maintain a long-lived pointer to the
// returned Transport so that any private credentials are able to be
// garbage collected.
func (f *FedCommonBehavior) NewTransport(c context.Context, actorBoxIRI *url.URL, gofedAgent string) (pub.Transport, error) {
	log.Println("NewTransport()")

	transport := &FedTransport{
		Context:   c,
		UserAgent: gofedAgent,
		Target:    actorBoxIRI,
	}

	return transport, nil
}
