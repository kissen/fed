package ap

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/fedd/db"
	"net/url"
	"path"
	"strings"
)

type IRI struct {
	Context context.Context
	Target  *url.URL
}

// Construct a new IRI with tailing components.
func NewIRI(c context.Context, components ...string) IRI {
	base := *FromContext(c).BasePath
	payload := path.Join(base, path.Join(components...))

	target := &url.URL{
		Scheme: *FromContext(c).Scheme,
		Host:   *FromContext(c).Host,
		Path:   payload,
	}

	return IRI{
		Context: c,
		Target:  target,
	}
}

// Generate a new actor IRI.
func ActorIRI(c context.Context, actor string) IRI {
	return NewIRI(c, actor)
}

// Generate a new inbox IRI.
func InboxIRI(c context.Context, owner string) IRI {
	return NewIRI(c, owner, "inbox")
}

// Generate a new outbox IRI.
func OutboxIRI(c context.Context, owner string) IRI {
	return NewIRI(c, owner, "outbox")
}

// Generate a new followers IRI.
func FollowersIRI(c context.Context, owner string) IRI {
	return NewIRI(c, owner, "followers")
}

// Generate a new following IRI.
func FollowingIRI(c context.Context, owner string) IRI {
	return NewIRI(c, owner, "following")
}

// Generate a new liked IRI.
func LikedIRI(c context.Context, owner string) IRI {
	return NewIRI(c, owner, "liked")
}

// Generate a new object IRI with a random UUID used as an object id.
func RollObjectIRI(c context.Context) IRI {
	id := uuid.New().String()
	return NewIRI(c, "storage", id)
}

// Return the owner of the given IRI. The IRI needs to have the form
//
//   */{username}
//
// where the asterix is the placeholder for the base path as defined
// in iri.Context.
func (iri IRI) Actor() (string, error) {
	if owner, dir, err := iri.split(); err != nil {
		return "", err
	} else if dir != nil {
		return "", fmt.Errorf("Target=%v not an actor", iri.Target)
	} else {
		return *owner, nil
	}
}

// Return the owner of the given IRI. The IRI needs to have the form
//
//   */{username}/inbox
//
// where the asterix is the placeholder for the base path as defined
// in iri.Context.
func (iri IRI) InboxOwner() (string, error) {
	if owner, dir, err := iri.split(); err != nil {
		return "", err
	} else if dir == nil || *dir != "inbox" {
		return "", fmt.Errorf("Target=%v not an inbox", iri.Target)
	} else {
		return *owner, nil
	}
}

// Return the owner of the given IRI. The IRI needs to have the form
//
//   */{username}/outbox //
// where the asterix is the placeholder for the base path as defined
// in iri.Context.
func (iri IRI) OutboxOwner() (string, error) {
	if owner, dir, err := iri.split(); err != nil {
		return "", err
	} else if dir == nil || *dir != "outbox" {
		return "", fmt.Errorf("Target=%v not an outbox", iri.Target)
	} else {
		return *owner, nil
	}
}

// Return the object id of the given IRI. The IRI needs to have the form
//
//   */storage/{id}
//
// where the asterix is the placeholder for the base path as defined
// in iri.Context.
func (iri IRI) Object() (string, error) {
	if dir, id, err := iri.split(); err != nil {
		return "", err
	} else if *dir != "storage" || id == nil {
		return "", fmt.Errorf("Target=%v not an object", iri.Target)
	} else {
		return *id, nil
	}
}

// Return the owner (username) of this IRI.
func (iri IRI) RetrieveOwner() (*db.FedUser, error) {
	// the owner of an IRI, in the easy case, is the first
	// path component; we do not support getting the owner
	// of object IRIs yet

	if username, _, err := iri.split(); err != nil {
		return nil, err
	} else if *username == "storage" {
		return nil, errors.New("reserved username")
	} else {
		return FromContext(iri.Context).Storage.RetrieveUser(*username)
	}
}

func (iri IRI) String() string {
	return iri.Target.String()
}

// Return the underlying URL for use in other functions.
func (iri IRI) URL() *url.URL {
	return iri.Target
}

// Split up the IRI and return the last two components. Here
// are some examples.
//
//   /alice        -> username=alice payload=nil
//   /alice/inbox  -> username=alice payload=inbox
//
// Of course the base path in iri.Context is taken into account.
//
// Returns (nil, nil, *error) on error, (*string, nil, nil) on
// actor IRIs and (*string, *string, nil) on other IRIs.
func (iri IRI) split() (username *string, payload *string, err error) {
	basePath := *FromContext(iri.Context).BasePath

	// split up base path and target into path components fo easier
	// handling

	base := iri.splitPath(basePath)
	target := iri.splitPath(iri.Target.Path)

	// check the length of target; a valid IRI has either one
	// (actor) or two (inbox, object) more components than the
	// base path

	actorlen := len(base) + 1
	contentlen := len(base) + 2

	if len(target) != actorlen && len(target) != contentlen {
		return nil, nil, fmt.Errorf("Target=%v does not match basePath=%v", iri.Target, basePath)
	}

	// match the individual components of the base path

	for i := range base {
		if target[i] != base[i] {
			return nil, nil, fmt.Errorf("Target=%v does not match basePath=%v", iri.Target, basePath)
		}
	}

	// we are golden; return the components we have

	if len(target) == actorlen {
		return iri.last(target), nil, nil
	} else {
		return iri.secondToLast(target), iri.last(target), nil
	}
}

// Split path into the individual path components. Never returns
// comonents that are empty.
func (iri IRI) splitPath(path string) []string {
	// split up into components

	components := strings.Split(path, "/")

	// remove empty components

	var trimmed []string

	for _, s := range components {
		t := strings.TrimSpace(s)

		if len(s) > 0 {
			trimmed = append(trimmed, t)
		}
	}

	return trimmed
}

// Return the last entryin the non-nil slice ss.
func (iri IRI) last(ss []string) *string {
	idx := len(ss) - 1
	return &ss[idx]
}

// Return the second to last entry in slice ss which has to
// contain at least two elements.
func (iri IRI) secondToLast(ss []string) *string {
	idx := len(ss) - 2
	return &ss[idx]
}
