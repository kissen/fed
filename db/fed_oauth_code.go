package db

import "time"

type FedOAuthCode struct {
	Code     string
	Username string
	IssuedOn time.Time
}

// Create a new randomized token for username and put it
// into the database at dst.
func NewFedOAuthCode(dst FedStorage, username string) (*FedOAuthCode, error) {
	oc := &FedOAuthCode{
		Code:     random(),
		Username: username,
		IssuedOn: time.Now().UTC(),
	}

	if err := dst.StoreCode(oc); err != nil {
		return nil, err
	}

	return oc, nil
}

// Return whether this code is expired.
func (c *FedOAuthCode) Expired() bool {
	return false
}
