package ap

import (
	"context"
	"errors"
	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams/vocab"
	"net/http"
)

// Implements the go-fed/activity/pub/SocialProtocol interface (version 1.0)
type FedSocialProtocol struct{}

// Hook callback after parsing the request body for a client request
// to the Actor's outbox.
//
// Can be used to set contextual information based on the
// ActivityStreams object received.
//
// Only called if the Social API is enabled.
//
// Warning: Neither authentication nor authorization has taken place at
// this time. Doing anything beyond setting contextual information is
// strongly discouraged.
//
// If an error is returned, it is passed back to the caller of
// PostOutbox. In this case, the DelegateActor implementation must not
// write a response to the ResponseWriter as is expected that the caller
// to PostOutbox will do so when handling the error.
func (f *FedSocialProtocol) PostOutboxRequestBodyHook(c context.Context, r *http.Request, data vocab.Type) (context.Context, error) {
	return nil, errors.New("not implemented")
}

// AuthenticatePostOutbox delegates the authentication of a POST to an
// outbox.
//
// Only called if the Social API is enabled.
//
// If an error is returned, it is passed back to the caller of
// PostOutbox. In this case, the implementation must not write a
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
func (f *FedSocialProtocol) AuthenticatePostOutbox(c context.Context, w http.ResponseWriter, r *http.Request) (out context.Context, authed bool, err error) {
	return nil, false, errors.New("not implemented")
}

// Callbacks returns the application logic that handles ActivityStreams
// received from C2S clients.
//
// Note that certain types of callbacks will be 'wrapped' with default
// behaviors supported natively by the library. Other callbacks
// compatible with streams.TypeResolver can be specified by 'other'.
//
// For example, setting the 'Create' field in the SocialWrappedCallbacks
// lets an application dependency inject additional behaviors they want
// to take place, including the default behavior supplied by this
// library. This is guaranteed to be compliant with the ActivityPub
// Social protocol.
//
// To override the default behavior, instead supply the function in
// 'other', which does not guarantee the application will be compliant
// with the ActivityPub Social Protocol.
//
// Applications are not expected to handle every single ActivityStreams
// type and extension. The unhandled ones are passed to DefaultCallback.
func (f *FedSocialProtocol) Callbacks(c context.Context) (wrapped pub.SocialWrappedCallbacks, other []interface{}, err error) {
	// Create handles additional side effects for the Create ActivityStreams
	// type.
	//
	// The wrapping callback copies the actor(s) to the 'attributedTo'
	// property and copies recipients between the Create activity and all
	// objects. It then saves the entry in the database.
	wrapped.Create = func(context.Context, vocab.ActivityStreamsCreate) error {
		return errors.New("not implemented")
	}

	// Update handles additional side effects for the Update ActivityStreams
	// type.
	//
	// The wrapping callback applies new top-level values on an object to
	// the stored objects. Any top-level null literals will be deleted on
	// the stored objects as well.
	wrapped.Update = func(context.Context, vocab.ActivityStreamsUpdate) error {
		return errors.New("not implemented")
	}

	// Delete handles additional side effects for the Delete ActivityStreams
	// type.
	//
	// The wrapping callback replaces the object(s) with tombstones in the
	// database.
	wrapped.Delete = func(context.Context, vocab.ActivityStreamsDelete) error {
		return errors.New("not implemented")
	}

	// Follow handles additional side effects for the Follow ActivityStreams
	// type.
	//
	// The wrapping callback only ensures the 'Follow' has at least one
	// 'object' entry, but otherwise has no default side effect.
	wrapped.Follow = func(context.Context, vocab.ActivityStreamsFollow) error {
		return errors.New("not implemented")
	}

	// Add handles additional side effects for the Add ActivityStreams
	// type.
	//
	//
	// The wrapping function will add the 'object' IRIs to a specific
	// 'target' collection if the 'target' collection(s) live on this
	// server.
	wrapped.Add = func(context.Context, vocab.ActivityStreamsAdd) error {
		return errors.New("not implemented")
	}

	// Remove handles additional side effects for the Remove ActivityStreams
	// type.
	//
	// The wrapping function will remove all 'object' IRIs from a specific
	// 'target' collection if the 'target' collection(s) live on this
	// server.
	wrapped.Remove = func(context.Context, vocab.ActivityStreamsRemove) error {
		return errors.New("not implemented")
	}

	// Like handles additional side effects for the Like ActivityStreams
	// type.
	//
	// The wrapping function will add the objects on the activity to the
	// "liked" collection of this actor.
	wrapped.Like = func(context.Context, vocab.ActivityStreamsLike) error {
		return errors.New("not implemented")
	}

	// Undo handles additional side effects for the Undo ActivityStreams
	// type.
	//
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
	// type.
	//
	// The wrapping callback only ensures the 'Block' has at least one
	// 'object' entry, but otherwise has no default side effect. It is up
	// to the wrapped application function to properly enforce the new
	// blocking behavior.
	//
	// Note that go-fed does not federate 'Block' activities received in the
	// Social Protocol.
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
func (f *FedSocialProtocol) DefaultCallback(c context.Context, activity pub.Activity) error {
	return errors.New("not implemented")
}
