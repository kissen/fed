package db

import (
	"gitlab.cs.fau.de/kissen/fed/errors"
)

// Check credentials in storage s. If they are valid credentials,
// this function returns nil. Otherwise it returns an error telling
// you what's wrong with the credentials.
func CheckCredentials(username, password string, s FedStorage) error {
	user, err := s.RetrieveUser(username)
	if err != nil {
		return errors.Wrap(err, "bad username")
	}

	if !user.PasswordOK(password) {
		return errors.New("bad password")
	}

	return nil
}
