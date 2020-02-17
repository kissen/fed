package ap

import (
	"context"
	"errors"
	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams/vocab"
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
	return nil, errors.New("not implemented")
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
	return nil, false, errors.New("not implemented")
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
	return false, errors.New("not implemented")
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
	// Create handles additional side effects for the Create ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping callback for the Federating Protocol ensures the
	// 'object' property is created in the database.
	//
	// Create calls Create for each object in the federated Activity.
	wrapped.Create = func(context.Context, vocab.ActivityStreamsCreate) error {
		return errors.New("not implemented")
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
		return errors.New("not implemented")
	}

	// Delete handles additional side effects for the Delete ActivityStreams
	// type, specific to the application using go-fed.
	//
	// Delete removes the federated entry from the database.
	wrapped.Delete = func(context.Context, vocab.ActivityStreamsDelete) error {
		return errors.New("not implemented")
	}

	// Follow handles additional side effects for the Follow ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping function can have one of several default behaviors,
	// depending on the value of the OnFollow setting.
	wrapped.Follow = func(context.Context, vocab.ActivityStreamsFollow) error {
		return errors.New("not implemented")
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
		return errors.New("not implemented")
	}

	// Reject handles additional side effects for the Reject ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping function has no default side effects. However, if this
	// 'Reject' is in response to a 'Follow' then the client MUST NOT go
	// forward with adding the 'actor' to the original 'actor's 'following'
	// collection by the client application.
	wrapped.Reject = func(context.Context, vocab.ActivityStreamsReject) error {
		return errors.New("not implemented")
	}

	// Add handles additional side effects for the Add ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping function will add the 'object' IRIs to a specific
	// 'target' collection if the 'target' collection(s) live on this
	// server.
	wrapped.Add = func(context.Context, vocab.ActivityStreamsAdd) error {
		return errors.New("not implemented")
	}

	// Remove handles additional side effects for the Remove ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping function will remove all 'object' IRIs from a specific
	// 'target' collection if the 'target' collection(s) live on this
	// server.
	wrapped.Remove = func(context.Context, vocab.ActivityStreamsRemove) error {
		return errors.New("not implemented")
	}

	// Like handles additional side effects for the Like ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping function will add the activity to the "likes" collection
	// on all 'object' targets owned by this server.
	wrapped.Like = func(context.Context, vocab.ActivityStreamsLike) error {
		return errors.New("not implemented")
	}

	// Announce handles additional side effects for the Announce
	// ActivityStreams type, specific to the application using go-fed.
	//
	// The wrapping function will add the activity to the "shares"
	// collection on all 'object' targets owned by this server.
	wrapped.Announce = func(context.Context, vocab.ActivityStreamsAnnounce) error {
		return errors.New("not implemented")
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
		return errors.New("not implemented")
	}

	// Block handles additional side effects for the Block ActivityStreams
	// type, specific to the application using go-fed.
	//
	// The wrapping function provides no default side effects. It simply
	// calls the wrapped function. However, note that Blocks should not be
	// received from a federated peer, as delivering Blocks explicitly
	// deviates from the original ActivityPub specification.
	wrapped.Block = func(context.Context, vocab.ActivityStreamsBlock) error {
		return errors.New("not implemented")
	}

	return wrapped, nil, errors.New("not implemented")
}

// DefaultCallback is called for types that go-fed can deserialize but
// are not handled by the application's callbacks returned in the
// Callbacks method.
//
// Applications are not expected to handle every single ActivityStreams
// type and extension, so the unhandled ones are passed to
// DefaultCallback.
func (f *FedFederatingProtocol) DefaultCallback(c context.Context, activity pub.Activity) error {
	return errors.New("not implemented")
}

// MaxInboxForwardingRecursionDepth determines how deep to search within
// an activity to determine if inbox forwarding needs to occur.
//
// Zero or negative numbers indicate infinite recursion.
func (f *FedFederatingProtocol) MaxInboxForwardingRecursionDepth(c context.Context) int {
	return -1
}

// MaxDeliveryRecursionDepth determines how deep to search within
// collections owned by peers when they are targeted to receive a
// delivery.
//
// Zero or negative numbers indicate infinite recursion.
func (f *FedFederatingProtocol) MaxDeliveryRecursionDepth(c context.Context) int {
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
	return nil, errors.New("not implemented")
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
	return nil, errors.New("not implemented")
}
