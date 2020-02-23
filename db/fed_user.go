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
