package db

import "time"

type FedOAuthToken struct {
	Token    string
	Username string
	IssuedOn time.Time
}

// Create a new randomized token for username and put it
// into the database at dst.
func NewFedOAuthToken(dst FedStorage, username string) (*FedOAuthToken, error) {
	ot := &FedOAuthToken{
		Token:    random(),
		Username: username,
		IssuedOn: time.Now().UTC(),
	}

	if err := dst.StoreToken(ot); err != nil {
		return nil, err
	}

	return ot, nil
}

// Return whether this token is expired.
func (c *FedOAuthToken) Expired() bool {
	return false
}
