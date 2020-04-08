package fediri

import (
	"fmt"
	"github.com/google/uuid"
	"gitlab.cs.fau.de/kissen/fed/config"
	"net/url"
	"path"
	"strings"
)

type IRI struct {
	Target *url.URL
}

// Construct a new IRI with tailing components.
func NewIRI(components ...string) IRI {
	u := config.Get().GlobalURL()
	u.Path = path.Join(u.Path, path.Join(components...))

	target := &url.URL{
		Scheme: "https",
		Host:   u.Host,
		Path:   u.Path,
	}

	return IRI{target}
}

// Generate a new actor IRI.
func ActorIRI(actor string) IRI {
	return NewIRI(actor)
}

// Generate a new inbox IRI.
func InboxIRI(owner string) IRI {
	return NewIRI(owner, "inbox")
}

// Generate a new outbox IRI.
func OutboxIRI(owner string) IRI {
	return NewIRI(owner, "outbox")
}

// Generate a new followers IRI.
func FollowersIRI(owner string) IRI {
	return NewIRI(owner, "followers")
}

// Generate a new following IRI.
func FollowingIRI(owner string) IRI {
	return NewIRI(owner, "following")
}

// Generate a new liked IRI.
func LikedIRI(owner string) IRI {
	return NewIRI(owner, "liked")
}

// Generate a new object IRI with a random UUID used as an object id.
func RollObjectIRI() IRI {
	id := uuid.New().String()
	return NewIRI("storage", id)
}

// Return the owner of the given IRI. The IRI needs to have the form
//
//   */{username}
//
// where the asterix is the placeholder for the base path.
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
// where the asterix is the placeholder for the base path.
func (iri IRI) InboxOwner() (string, error) {
	return iri.owner("inbox")
}

// Return the owner of the given IRI. The IRI needs to have the form
//
//   */{username}/outbox //
// where the asterix is the placeholder for the base path.
func (iri IRI) OutboxOwner() (string, error) {
	return iri.owner("outbox")
}

// Return the owner of the given IRI. The IRI needs to have the form
//
//   */{username}/following
//
// where the asterix is the placeholder for the base path.
func (iri IRI) FollowingOwner() (string, error) {
	return iri.owner("following")
}

// Return the owner of the given IRI. The IRI needs to have the form
//
//   */{username}/followers
//
// where the asterix is the placeholder for the base path.
func (iri IRI) FollowersOwner() (string, error) {
	return iri.owner("followers")
}

// Return the owner of the given IRI. The IRI needs to have the form
//
//   */{username}/liked
//
// where the asterix is the placeholder for the base path.
func (iri IRI) LikedOwner() (string, error) {
	return iri.owner("liked")
}

// Return the object id of the given IRI. The IRI needs to have the form
//
//   */storage/{id}
//
// where the asterix is the placeholder for the base path.
func (iri IRI) Object() (string, error) {
	if dir, id, err := iri.split(); err != nil {
		return "", err
	} else if *dir != "storage" || id == nil {
		return "", fmt.Errorf("Target=%v not an object", iri.Target)
	} else {
		return *id, nil
	}
}

// Return the owner of this IRI.
func (iri IRI) Owner() (string, error) {
	if username, _, err := iri.split(); err != nil {
		return "", err
	} else if IsReservedUsername(*username) {
		return "", fmt.Errorf("reserved username=%v", *username)
	} else {
		return *username, nil
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
// Of course the base path is taken into account.
//
// Returns (nil, nil, *error) on error, (*string, nil, nil) on
// actor IRIs and (*string, *string, nil) on other IRIs.
func (iri IRI) split() (username *string, payload *string, err error) {
	// XXX: we used to have rudimentary support for running under a
	// subdirectory on a server, e.g. example.com/foo/bar/;
	// this is currently not supported; patches welcome :-)

	basePath := config.Get().GlobalURL().Path

	// split up base path and target into path components for easier
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

// Return the owner of the given IRI. The IRI needs to have the form
//
//   */{username}/$tail
//
// where the asterix is the placeholder for the base path.
func (iri IRI) owner(tail string) (string, error) {
	if owner, dir, err := iri.split(); err != nil {
		return "", err
	} else if dir == nil || *dir != tail {
		return "", fmt.Errorf("Target=%v does not have required tail=/%v", iri.Target, tail)
	} else {
		return *owner, nil
	}
}
