package fedcontext

import (
	"gitlab.cs.fau.de/kissen/fed/db"
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
