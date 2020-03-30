package db

import (
	"encoding/json"
	"errors"
	"time"
)

const _CODE_LIFETIME = 1 * time.Minute

type FedOAuthCode struct {
	Code     string
	Username string
	IssuedOn time.Time
}

// If username/password are valid credentials, create a new
// code, store it into target and return it.
func NewFedOAuthCode(username, password string, target FedStorage) (*FedOAuthCode, error) {
	if err := CheckCredentials(username, password, target); err != nil {
		return nil, err
	}

	oc := &FedOAuthCode{
		Code:     random(),
		Username: username,
		IssuedOn: time.Now().UTC(),
	}

	if err := target.StoreCode(oc); err != nil {
		return nil, err
	}

	return oc, nil
}

// Return whether this code is expired.
func (c *FedOAuthCode) Expired() bool {
	end := c.IssuedOn.Add(_CODE_LIFETIME)
	return time.Now().UTC().After(end)
}

// Unmarshal override to avoid confusion with FedOAuthToken.
func (c *FedOAuthCode) UnmarshalJSON(data []byte) error {
	type Code struct {
		Code     string
		Username string
		IssuedOn time.Time
	}

	var buf Code

	if err := json.Unmarshal(data, &buf); err != nil {
		return err
	}

	if len(buf.Code) == 0 {
		return errors.New("empty code")
	}

	c.Code = buf.Code
	c.Username = buf.Username
	c.IssuedOn = buf.IssuedOn

	return nil
}
