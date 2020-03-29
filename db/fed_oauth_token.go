package db

import "time"

const _TOKEN_LIFETIME = 24 * time.Hour

type FedOAuthToken struct {
	Token    string
	Username string
	IssuedOn time.Time
}

// If username/password are valid credentials, create a new
// token, store it into target and return it.
func NewFedOAuthToken(username, password string, target FedStorage) (*FedOAuthToken, error) {
	if err := CheckCredentials(username, password, target); err != nil {
		return nil, err
	}

	ot := &FedOAuthToken{
		Token:    random(),
		Username: username,
		IssuedOn: time.Now().UTC(),
	}

	if err := target.StoreToken(ot); err != nil {
		return nil, err
	}

	return ot, nil
}

// Create a new token for username, store it into target and return it.
func NewFedOAuthTokenFor(username string, target FedStorage) (*FedOAuthToken, error) {
	ot := &FedOAuthToken{
		Token:    random(),
		Username: username,
		IssuedOn: time.Now().UTC(),
	}

	if err := target.StoreToken(ot); err != nil {
		return nil, err
	}

	return ot, nil
}

// Return whether this token is expired.
func (c *FedOAuthToken) Expired() bool {
	end := c.IssuedOn.Add(_TOKEN_LIFETIME)
	return time.Now().UTC().After(end)
}
