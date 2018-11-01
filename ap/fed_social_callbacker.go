package ap

import (
	"context"
	"github.com/go-fed/activity/streams"
)

// Implements the github.com/go-fed/activity/pub/Callbacker interface
// for use in the social (client <-> server) API
type FedSocialCallbacker struct {
}

// Create Activity callback.
func (f *FedSocialCallbacker) Create(c context.Context, s *streams.Create) error {
	// v := s.Raw()
	return nil
}

// Update Activity callback.
func (f *FedSocialCallbacker) Update(c context.Context, s *streams.Update) error {
	// v := s.Raw()
	return nil
}

// Delete Activity callback.
func (f *FedSocialCallbacker) Delete(c context.Context, s *streams.Delete) error {
	// v := s.Raw()
	return nil
}

// Add Activity callback.
func (f *FedSocialCallbacker) Add(c context.Context, s *streams.Add) error {
	// v := s.Raw()
	return nil
}

// Remove Activity callback.
func (f *FedSocialCallbacker) Remove(c context.Context, s *streams.Remove) error {
	// v := s.Raw()
	return nil
}

// Like Activity callback.
func (f *FedSocialCallbacker) Like(c context.Context, s *streams.Like) error {
	// v := s.Raw()
	return nil
}

// Block Activity callback. By default, this implmentation does not
// dictate how blocking should be implemented, so it is up to the
// application to enforce this by implementing the FederateApp
// interface.
func (f *FedSocialCallbacker) Block(c context.Context, s *streams.Block) error {
	// v := s.Raw()
	return nil
}

// Follow Activity callback. In the special case of server-to-server
// delivery of a Follow activity, this implementation supports the
// option of automatically replying with an 'Accept', 'Reject', or
// waiting for human interaction as provided in the FederateApp
// interface.
//
// In the special case that the FederateApp returned AutomaticAccept,
// this library automatically handles adding the 'actor' to the
// 'followers' collection of the 'object'.
func (f *FedSocialCallbacker) Follow(c context.Context, s *streams.Follow) error {
	// v := s.Raw()
	return nil
}

// Undo Activity callback. It is up to the client to provide support
// for all 'Undo' operations; this implementation does not attempt to
// provide a generic implementation.
func (f *FedSocialCallbacker) Undo(c context.Context, s *streams.Undo) error {
	// v := s.Raw()
	return nil
}

// Accept Activity callback. In the special case that this 'Accept'
// activity has an 'object' of 'Follow' type, then the library will
// handle adding the 'actor' to the 'following' collection of the
// original 'actor' who requested the 'Follow'.
func (f *FedSocialCallbacker) Accept(c context.Context, s *streams.Accept) error {
	// v := s.Raw()
	return nil
}

// Reject Activity callback. Note that in the special case that this
// 'Reject' activity has an 'object' of 'Follow' type, then the client
// MUST NOT add the 'actor' to the 'following' collection of the
// original 'actor' who requested the 'Follow'.
func (f *FedSocialCallbacker) Reject(c context.Context, s *streams.Reject) error {
	// v := s.Raw()
	return nil
}

// https://go-fed.org/tutorial#ActivityStreams-Types-and-Properties
// Starting in version 0.3.0, you may define additional methods on
// your Callbacker implementation to handle additional activities not
// required by the Callbacker interface.

func (f *FedSocialCallbacker) Announce(c context.Context, s *streams.Announce) error {
	// v := s.Raw()
	return nil
}

func (f *FedSocialCallbacker) Arrive(c context.Context, s *streams.Arrive) error {
	// v := s.Raw()
	return nil
}

func (f *FedSocialCallbacker) Dislike(c context.Context, s *streams.Dislike) error {
	// v := s.Raw()
	return nil
}

func (f *FedSocialCallbacker) Flag(c context.Context, s *streams.Flag) error {
	// v := s.Raw()
	return nil
}

func (f *FedSocialCallbacker) Ignore(c context.Context, s *streams.Ignore) error {
	// v := s.Raw()
	return nil
}

func (f *FedSocialCallbacker) Invite(c context.Context, s *streams.Invite) error {
	// v := s.Raw()
	return nil
}

func (f *FedSocialCallbacker) Join(c context.Context, s *streams.Join) error {
	// v := s.Raw()
	return nil
}

func (f *FedSocialCallbacker) Leave(c context.Context, s *streams.Leave) error {
	// v := s.Raw()
	return nil
}

func (f *FedSocialCallbacker) Listen(c context.Context, s *streams.Listen) error {
	// v := s.Raw()
	return nil
}

func (f *FedSocialCallbacker) Move(c context.Context, s *streams.Move) error {
	// v := s.Raw()
	return nil
}

func (f *FedSocialCallbacker) Offer(c context.Context, s *streams.Offer) error {
	// v := s.Raw()
	return nil
}

func (f *FedSocialCallbacker) Question(c context.Context, s *streams.Question) error {
	// v := s.Raw()
	return nil
}

func (f *FedSocialCallbacker) Read(c context.Context, s *streams.Read) error {
	// v := s.Raw()
	return nil
}

func (f *FedSocialCallbacker) TentativeAccept(c context.Context, s *streams.TentativeAccept) error {
	// v := s.Raw()
	return nil
}

func (f *FedSocialCallbacker) TentativeReject(c context.Context, s *streams.TentativeReject) error {
	// v := s.Raw()
	return nil
}

func (f *FedSocialCallbacker) Travel(c context.Context, s *streams.Travel) error {
	// v := s.Raw()
	return nil
}

func (f *FedSocialCallbacker) View(c context.Context, s *streams.View) error {
	// v := s.Raw()
	return nil
}
