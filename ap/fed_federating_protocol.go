package ap

import (
	"context"
	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams/vocab"
	"gitlab.cs.fau.de/kissen/fed/errors"
	"gitlab.cs.fau.de/kissen/fed/prop"
	"log"
	"net/http"
	"net/url"
)

// Implements the go-fed/activity/pub/FederatingProtocol interface (version 1.0)
type FedFederatingProtocol struct{}

// Hook callback after parsing the request body for a federated request
// to the Actor's inbox.
//
// Can be used to set contextual information based on the Activity
// received.
//
// Only called if the Federated Protocol is enabled.
//
// Warning: Neither authentication nor authorization has taken place at
// this time. Doing anything beyond setting contextual information is
// strongly discouraged.
//
// If an error is returned, it is passed back to the caller of
// PostInbox. In this case, the DelegateActor implementation must not
// write a response to the ResponseWriter as is expected that the caller
// to PostInbox will do so when handling the error.
func (f *FedFederatingProtocol) PostInboxRequestBodyHook(c context.Context, r *http.Request, activity pub.Activity) (context.Context, error) {
	log.Printf("PostInboxRequestBodyHook(%v)", r.URL)
	return c, nil
}

// AuthenticatePostInbox delegates the authentication of a POST to an
// inbox.
//
// If an error is returned, it is passed back to the caller of
// PostInbox. In this case, the implementation must not write a
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
func (f *FedFederatingProtocol) AuthenticatePostInbox(c context.Context, w http.ResponseWriter, r *http.Request) (out context.Context, authed bool, err error) {
	log.Printf("AuthenticatePostInbox(%v)", r.URL)
	return c, true, nil
}

// Blocked should determine whether to permit a set of actors given by
// their ids are able to interact with this particular end user due to
// being blocked or other application-specific logic.
//
// If an error is returned, it is passed back to the caller of
// PostInbox.
//
// If no error is returned, but authentication or authorization fails,
// then blocked must be true and error nil. An http.StatusForbidden
// will be written in the wresponse.
//
// Finally, if the authentication and authorization succeeds, then
// blocked must be false and error nil. The request will continue
// to be processed.
func (f *FedFederatingProtocol) Blocked(c context.Context, actorIRIs []*url.URL) (blocked bool, err error) {
	log.Printf("Blocked(%v)", actorIRIs)
	return false, nil
}

// Callbacks returns the application logic that handles ActivityStreams
// received from federating peers.
//
// Note that certain types of callbacks will be 'wrapped' with default
// behaviors supported natively by the library. Other callbacks
// compatible with streams.TypeResolver can be specified by 'other'.
//
// For example, setting the 'Create' field in the
// FederatingWrappedCallbacks lets an application dependency inject
// additional behaviors they want to take place, including the default
// behavior supplied by this library. This is guaranteed to be compliant
// with the ActivityPub Social protocol.
//
// To override the default behavior, instead supply the function in
// 'other', which does not guarantee the application will be compliant
// with the ActivityPub Social Protocol.
//
// Applications are not expected to handle every single ActivityStreams
// type and extension. The unhandled ones are passed to DefaultCallback.
func (f *FedFederatingProtocol) Callbacks(c context.Context) (wrapped pub.FederatingWrappedCallbacks, other []interface{}, err error) {
	log.Println("Callbacks()")

	// Create handles additional side effects for the Create ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping callback for the Federating Protocol ensures the
	// 'object' property is created in the database.
	//
	// Create calls Create for each object in the federated Activity.
	wrapped.Create = func(context.Context, vocab.ActivityStreamsCreate) error {
		log.Println("Create()")
		return nil
	}

	// Update handles additional side effects for the Update ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping callback for the Federating Protocol ensures the
	// 'object' property is updated in the database.
	//
	// Update calls Update on the federated entry from the database, with a
	// new value.
	wrapped.Update = func(context.Context, vocab.ActivityStreamsUpdate) error {
		log.Println("Update()")
		return nil
	}

	// Delete handles additional side effects for the Delete ActivityStreams
	// type, specific to the application using go-fed.
	//
	// Delete removes the federated entry from the database.
	wrapped.Delete = func(context.Context, vocab.ActivityStreamsDelete) error {
		log.Println("Delete()")
		return nil
	}

	// Follow handles additional side effects for the Follow ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping function can have one of several default behaviors,
	// depending on the value of the OnFollow setting.
	wrapped.Follow = func(context.Context, vocab.ActivityStreamsFollow) error {
		log.Println("Follow()")
		return nil
	}

	// OnFollow determines what action to take for this particular callback
	// if a Follow Activity is handled.
	wrapped.OnFollow = pub.OnFollowAutomaticallyAccept

	// Accept handles additional side effects for the Accept ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping function determines if this 'Accept' is in response to a
	// 'Follow'. If so, then the 'actor' is added to the original 'actor's
	// 'following' collection.
	//
	// Otherwise, no side effects are done by go-fed.
	wrapped.Accept = func(context.Context, vocab.ActivityStreamsAccept) error {
		log.Println("Accept()")
		return nil
	}

	// Reject handles additional side effects for the Reject ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping function has no default side effects. However, if this
	// 'Reject' is in response to a 'Follow' then the client MUST NOT go
	// forward with adding the 'actor' to the original 'actor's 'following'
	// collection by the client application.
	wrapped.Reject = func(context.Context, vocab.ActivityStreamsReject) error {
		log.Println("Reject()")
		return nil
	}

	// Add handles additional side effects for the Add ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping function will add the 'object' IRIs to a specific
	// 'target' collection if the 'target' collection(s) live on this
	// server.
	wrapped.Add = func(context.Context, vocab.ActivityStreamsAdd) error {
		log.Println("Add()")
		return nil
	}

	// Remove handles additional side effects for the Remove ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping function will remove all 'object' IRIs from a specific
	// 'target' collection if the 'target' collection(s) live on this
	// server.
	wrapped.Remove = func(context.Context, vocab.ActivityStreamsRemove) error {
		log.Println("Remove()")
		return nil
	}

	// Like handles additional side effects for the Like ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping function will add the activity to the "likes" collection
	// on all 'object' targets owned by this server.
	wrapped.Like = func(context.Context, vocab.ActivityStreamsLike) error {
		log.Println("Like()")
		return nil
	}

	// Announce handles additional side effects for the Announce
	// ActivityStreams type, specific to the application using go-fed.
	//
	// The wrapping function will add the activity to the "shares"
	// collection on all 'object' targets owned by this server.
	wrapped.Announce = func(context.Context, vocab.ActivityStreamsAnnounce) error {
		log.Println("Announce()")
		return nil
	}

	// Undo handles additional side effects for the Undo ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping function ensures the 'actor' on the 'Undo'
	// is be the same as the 'actor' on all Activities being undone.
	// It enforces that the actors on the Undo must correspond to all of the
	// 'object' actors in some manner.
	//
	// It is expected that the application will implement the proper
	// reversal of activities that are being undone.
	wrapped.Undo = func(context.Context, vocab.ActivityStreamsUndo) error {
		log.Println("Undo()")
		return nil
	}

	// Block handles additional side effects for the Block ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping function provides no default side effects. It simply
	// calls the wrapped function. However, note that Blocks should not be
	// received from a federated peer, as delivering Blocks explicitly
	// deviates from the original ActivityPub specification.
	wrapped.Block = func(context.Context, vocab.ActivityStreamsBlock) error {
		log.Println("Block()")
		return nil
	}

	return wrapped, nil, nil
}

// DefaultCallback is called for types that go-fed can deserialize but
// are not handled by the application's callbacks returned in the
// Callbacks method.
//
// Applications are not expected to handle every single ActivityStreams
// type and extension, so the unhandled ones are passed to
// DefaultCallback.
func (f *FedFederatingProtocol) DefaultCallback(c context.Context, activity pub.Activity) error {
	log.Println("DefaultCallback()")
	return nil
}

// MaxInboxForwardingRecursionDepth determines how deep to search within
// an activity to determine if inbox forwarding needs to occur.
//
// Zero or negative numbers indicate infinite recursion.
func (f *FedFederatingProtocol) MaxInboxForwardingRecursionDepth(c context.Context) int {
	log.Println("MaxInboxForwardingRecursionDepth()")
	return -1
}

// MaxDeliveryRecursionDepth determines how deep to search within
// collections owned by peers when they are targeted to receive a
// delivery.
//
// Zero or negative numbers indicate infinite recursion.
func (f *FedFederatingProtocol) MaxDeliveryRecursionDepth(c context.Context) int {
	log.Println("MaxDeliveryRecursionDepth()")
	return -1
}

// FilterForwarding allows the implementation to apply business logic
// such as blocks, spam filtering, and so on to a list of potential
// Collections and OrderedCollections of recipients when inbox
// forwarding has been triggered.
//
// The activity is provided as a reference for more intelligent
// logic to be used, but the implementation must not modify it.
func (f *FedFederatingProtocol) FilterForwarding(c context.Context, potentialRecipients []*url.URL, a pub.Activity) (filteredRecipients []*url.URL, err error) {
	log.Println("FilterForwarding()")
	return nil, nil
}

// GetInbox returns the OrderedCollection inbox of the actor for this
// context. It is up to the implementation to provide the correct
// collection for the kind of authorization given in the request.
//
// AuthenticateGetInbox will be called prior to this.
//
// Always called, regardless whether the Federated Protocol or Social
// API is enabled.
func (f *FedFederatingProtocol) GetInbox(c context.Context, r *http.Request) (vocab.ActivityStreamsOrderedCollectionPage, error) {
	log.Printf("GetInbox(%v)", r.URL)

	iri := IRI{c, r.URL}

	if user, err := iri.RetrieveOwner(); err != nil {
		return nil, err
	} else if page, err := collectPage(c, user.Inbox); err != nil {
		return nil, errors.Wrap(err, "collect failed")
	} else {
		prop.SetIdOn(page, iri.URL())
		return page, nil
	}
}
