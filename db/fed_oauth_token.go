package db

import (
	"encoding/json"
	"errors"
	"time"
)

const _TOKEN_LIFETIME = 24 * time.Hour

type FedOAuthToken struct {
	Token    string
	Username string
	IssuedOn time.Time
}

// If username/password are valid credentials, create a new
// token, store it into target and return it.
func NewFedOAuthToken(username, password string, target Tx) (*FedOAuthToken, error) {
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
func NewFedOAuthTokenFor(username string, target Tx) (*FedOAuthToken, error) {
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

// Unmarshal override to avoid confusion with FedOAuthCode.
func (c *FedOAuthToken) UnmarshalJSON(data []byte) error {
	type Token struct {
		Token    string
		Username string
		IssuedOn time.Time
	}

	var buf Token

	if err := json.Unmarshal(data, &buf); err != nil {
		return err
	}

	if len(buf.Token) == 0 {
		return errors.New("empty token")
	}

	c.Token = buf.Token
	c.Username = buf.Username
	c.IssuedOn = buf.IssuedOn

	return nil
}
