package db

import "time"

type FedOAuthToken struct {
	Token    string
	Username string
	IssuedOn time.Time
}

func (c *FedOAuthToken) Expired() bool {
	return false
}
