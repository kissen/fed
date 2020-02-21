package db

import "net/url"

// Represents a user registered with the service.
type FedUser struct {
	Name string

	Inbox     []*url.URL
	Outbox    []*url.URL
	Following []*url.URL
	Folowers  []*url.URL
	Liked     []*url.URL
}
