package fediri

import "github.com/kissen/stringset"

var reserved = stringset.NewWith(
	"storage", "static", "oauth", "stream", "liked",
	"following", "followers", "login", "logout", "remote",
	"submit",
)

// Return whether username is a reserved username, that is a name
// that may not appear as first IRI component because it has other
// functions.
func IsReservedUsername(username string) bool {
	return reserved.Contains(username)
}
