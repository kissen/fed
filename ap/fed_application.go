package ap

import (
	"context"
	"crypto"
	"fmt"
	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/go-fed/httpsig"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/db"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

// Implements the
//
//   github.com/go-fed/activity/pub/Application
//   github.com/go-fed/activity/pub/SocialAPI
//   github.com/go-fed/activity/pub/SocialApplication
//
// interfaces
type FedApplication struct {
	BaseIRI string
	Storage db.FedStorer
}

/*
 * Implementation of Application methods
 */

// Determines whether the application owns an IRI
func (f *FedApplication) Owns(c context.Context, id *url.URL) bool {
	log.Printf("Owns(id=%v)\n", id)

	if _, err := getActivityIdFromIRI(f.BaseIRI, id); err != nil {
		return true
	}

	if _, err := getInboxUsernameFromIRI(f.BaseIRI, id); err != nil {
		return true
	}

	if _, err := getOutboxUsernameFromIRI(f.BaseIRI, id); err != nil {
		return true
	}

	return false
}

// Gets ActivityStream content
func (f *FedApplication) Get(c context.Context, id *url.URL, rw pub.RWType) (pub.PubObject, error) {
	log.Printf("Get(id=%v rw=%v)\n", id, rw)
	log.Printf("Ignoring rw=%v in Get()\n", rw)

	postId, err := getActivityIdFromIRI(f.BaseIRI, id)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get integer id for id=%v", id)
	}

	fedPost := f.Storage.GetPost(postId)
	if fedPost == nil {
		return nil, fmt.Errorf("no activity found for postId=%v", postId)
	}

	return postToNote(fedPost), nil
}

// GetAsVerifiedUser fetches the ActivityStream representation of the
// given id with the provided IRI representing the authenticated user
// making the request.
func (f *FedApplication) GetAsVerifiedUser(c context.Context, id, authdUser *url.URL, rw pub.RWType) (pub.PubObject, error) {
	log.Printf("GetAsVerifiedUser(id=%v rw=%v)\n", id, rw)
	log.Printf("Using Get(id=%v) to implement GetAsverifieduser()\n", id)

	return f.Get(c, id, rw)
}

// Determines if it has ActivityStream data at the IRI
func (f *FedApplication) Has(c context.Context, id *url.URL) (bool, error) {
	log.Printf("Has(id=%v)\n", id)

	postId, err := getActivityIdFromIRI(f.BaseIRI, id)
	if err != nil {
		return false, errors.Wrapf(err, "could not get integer id for id=%v", id)
	}

	haveActivity := f.Storage.GetPost(postId) != nil
	return haveActivity, nil
}

// Setting ActivityStream data
func (f *FedApplication) Set(c context.Context, o pub.PubObject) error {
	log.Printf("Set()\n")

	return errors.New("setting activities not supported")
}

// Getting an actor's outbox or inbox
func (f *FedApplication) GetInbox(c context.Context, r *http.Request, rw pub.RWType) (vocab.OrderedCollectionType, error) {
	log.Printf("GetInbox(r.URL=%v rw=%v)\n", r.URL, rw)

	id := r.URL

	log.Printf("Ignoring rw=%v in GetInbox()\n", rw)
	log.Printf("Returning empty inbox on GetInbox(id=%v)\n", id)

	empty := &vocab.OrderedCollection{}
	return empty, nil
}

func (f *FedApplication) GetOutbox(c context.Context, r *http.Request, rw pub.RWType) (vocab.OrderedCollectionType, error) {
	log.Printf("GetOutbox(r.URL=%v, rw=%v)\n", r.URL, rw)
	log.Printf("Ignoring rw=%v in GetOutbox()\n", rw)

	id := r.URL

	username, err := getInboxUsernameFromIRI(f.BaseIRI, id)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get username for id=%v", id)
	}

	user := f.Storage.FindUser(username)
	if user == nil {
		return nil, fmt.Errorf("could not get user with username=%v", user)
	}

	posts := f.Storage.GetPostsFrom(user.Id)
	return postsToOutbox(user, posts), nil
}

// Creating new ids
func (f *FedApplication) NewId(c context.Context, t pub.Typer) *url.URL {
	log.Printf("NewId()\n")
	log.Printf("Returning dummy id in NewId()\n")

	id, err := url.Parse("https://localhost/activity/1337")
	if err != nil {
		panic(err)
	}

	return id
}

// Obtaining the public key for another user for verification purposes
func (f *FedApplication) GetPublicKey(c context.Context, publicKeyId string) (pubKey crypto.PublicKey, algo httpsig.Algorithm, user *url.URL, err error) {
	log.Printf("GetPublicKey(publicKeyId=%v)\n", publicKeyId)

	return nil, httpsig.RSA_SHA256, nil, errors.New("public keys not supported")
}

// Whether adding/removing is permitted
func (f *FedApplication) CanAdd(c context.Context, o vocab.ObjectType, t vocab.ObjectType) bool {
	log.Printf("CanAdd()\n")

	return true
}

func (f *FedApplication) CanRemove(c context.Context, o vocab.ObjectType, t vocab.ObjectType) bool {
	log.Printf("CanRemove()\n")

	return false
}

/*
 * Implementation of SocialAPI methods
 */

// ActorIRI returns the actor's IRI associated with the given request.
func (f *FedApplication) ActorIRI(c context.Context, r *http.Request) (*url.URL, error) {
	return r.URL, nil
}

// GetSocialAPIVerifier returns the authentication mechanism used for
// incoming ActivityPub client requests. It is optional and allowed to
// return null.
//
// Note that regardless of what this implementation returns, HTTP
// Signatures is supported natively as a fallback.
func (f *FedApplication) GetSocialAPIVerifier(c context.Context) pub.SocialAPIVerifier {
	// only required for oauth2
	return nil
}

// GetPublicKeyForOutbox fetches the public key for a user based on the
// public key id. It also determines which algorithm to use to verify
// the signature.
//
// Note that a key difference from Application's GetPublicKey is that
// this function must make sure that the actor whose boxIRI is passed in
// matches the public key id that is requested, or return an error.
func (f *FedApplication) GetPublicKeyForOutbox(c context.Context, publicKeyId string, boxIRI *url.URL) (crypto.PublicKey, httpsig.Algorithm, error) {
	return nil, httpsig.RSA_SHA256, errors.New("not impelemented")
}

/*
 * Helpers
 */

// panics on error
func postToNote(post *db.FedPost) *vocab.Note {
	if post == nil {
		panic("post musn't be nil")
	}

	note := &vocab.Note{}
	note.AppendNameString(post.Content)

	return note
}

func getActivityIdFromIRI(baseIRI string, id *url.URL) (uint64, error) {
	activityPattern := path.Join(baseIRI, "activity") + "/*"

	// path.Match only fails on invalid patterns; that should not
	// happen

	patternMatches, err := path.Match(activityPattern, id.Path)

	if err != nil {
		panic(err)
	}

	if !patternMatches {
		return 0, fmt.Errorf("cannot infer activity id on id=%v", id)
	}

	// now parse out the id

	_, filename := path.Split(id.Path)

	activityId, err := strconv.ParseUint(filename, 10, 64)

	if err != nil {
		return 0, errors.Wrap(err, "cannot infer activity id")
	}

	return activityId, nil
}

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func getInboxUsernameFromIRI(baseIRI string, id *url.URL) (string, error) {
	activityPattern := path.Join(baseIRI, "inbox") + "/*"

	// path.Match only fails on invalid patterns; that should not
	// happen

	patternMatches, err := path.Match(activityPattern, id.Path)

	if err != nil {
		panic(err)
	}

	if !patternMatches {
		return "", fmt.Errorf("cannot infer inbox username on id=%v", id)
	}

	// now parse out the username

	_, username := path.Split(id.Path)

	if isEmpty(username) {
		return "", errors.New("username for inbox is empty")
	}

	return username, nil
}

func getOutboxUsernameFromIRI(baseIRI string, id *url.URL) (string, error) {
	activityPattern := path.Join(baseIRI, "outbox") + "/*"

	// path.Match only fails on invalid patterns; that should not
	// happen

	patternMatches, err := path.Match(activityPattern, id.Path)

	if err != nil {
		panic(err)
	}

	if !patternMatches {
		return "", fmt.Errorf("cannot infer outbox username on id=%v", id)
	}

	// now parse out the username

	_, username := path.Split(id.Path)

	if isEmpty(username) {
		return "", errors.New("username for outbox is empty")
	}

	return username, nil
}


// panics on error
func postsToOutbox(user *db.FedUser, posts []*db.FedPost) vocab.OrderedCollectionType {
	outbox := &vocab.OrderedCollection{}

	outbox.AppendSummaryString("Outbox of user=" + user.Name)

	for _, post := range posts {
		note := postToNote(post)
		outbox.AppendOrderedItemsObject(note)
	}

	return outbox
}
