package fedcontext

import (
	"gitlab.cs.fau.de/kissen/fed/db"
	"gitlab.cs.fau.de/kissen/fed/errors"
	"net/http"
)

type Permissions struct {
	// The user as which the request is signed in.
	User db.FedUser

	// Whether these permissions allow the request to
	// issue Create on behalf of User.
	Create bool

	// Whether these permissions allow the request to
	// issue Like on behalf of User.
	Like bool

	// tbc
}

// Given a username and password, try to look up that user in the database.
// If a correct password was supplied, return the permissions that user has.
func PermissionsFrom(r *http.Request, username, password string) (*Permissions, error) {
	storage := Context(r).Storage

	user, err := storage.RetrieveUser(username)
	if err != nil {
		return nil, errors.NewWith(http.StatusUnauthorized, "bad username")
	}

	if !user.PasswordOK(password) {
		return nil, errors.NewWith(http.StatusUnauthorized, "bad password")
	}

	p := &Permissions{
		User:   *user,
		Create: true,
		Like:   true,
	}

	return p, nil
}
