package db

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gitlab.cs.fau.de/kissen/fed/errors"
	"gitlab.cs.fau.de/kissen/fed/util"
	"net/url"
)

// Represents a user registered with the service.
type FedUser struct {
	Name           string
	PasswordSHA256 []byte

	Inbox     []*url.URL
	Outbox    []*url.URL
	Following []*url.URL
	Followers []*url.URL
	Liked     []*url.URL
}

// Return a slice that contains all collections (i.e. Inbox, Outbox,
// Following, Followers and Liked).
func (u *FedUser) Collections() [][]*url.URL {
	return [][]*url.URL{
		u.Inbox, u.Outbox, u.Following, u.Followers, u.Liked,
	}
}

// Hash the plaintext password and assign the result to
// FedUser.PasswordSHA256.
func (u *FedUser) SetPassword(password string) {
	u.PasswordSHA256 = u.hash(password)
}

// Return whether plaintext password, when hashed, matches
// the assigned password.
func (u *FedUser) PasswordOK(password string) bool {
	hash := u.hash(password)
	return bytes.Compare(hash, u.PasswordSHA256) == 0
}

// Returns whether this user is following whatever is at id.
func (u *FedUser) IsFollowing(id *url.URL) bool {
	return util.UrlIn(id, u.Following)
}

// Returns whether id is following this user.
func (u *FedUser) IsFollowedBy(id *url.URL) bool {
	return util.UrlIn(id, u.Followers)
}

// Returns whether this user liked whawtever is at id.
func (u *FedUser) HasLiked(id *url.URL) bool {
	return util.UrlIn(id, u.Liked)
}

func (u *FedUser) String() string {
	return fmt.Sprintf(
		"{Name=%v Inbox=%v Outbox=%v Following=%v Followers=%v Liked=%v}",
		u.Name, u.Inbox, u.Outbox, u.Following, u.Followers, u.Liked,
	)
}

func (u *FedUser) hash(password string) []byte {
	// TODO: salt

	h := sha256.New()
	h.Write([]byte(password))
	return h.Sum(nil)
}

func userToBytes(user *FedUser) ([]byte, error) {
	if bytes, err := json.Marshal(user); err != nil {
		return nil, errors.Wrap(err, "byte marshal from user failed")
	} else {
		return bytes, nil
	}
}

func bytesToUser(bin []byte) (*FedUser, error) {
	var user FedUser

	if err := json.Unmarshal(bin, &user); err != nil {
		return nil, errors.Wrap(err, "byte unmarshal to user failed")
	} else {
		return &user, nil
	}
}
