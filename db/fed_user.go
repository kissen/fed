package db

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/url"
)

// Represents a user registered with the service.
type FedUser struct {
	Name string

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
