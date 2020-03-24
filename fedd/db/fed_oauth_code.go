package db

import "time"

type FedOAuthCode struct {
	Code     string
	Username string
	IssuedOn time.Time
}

func (c *FedOAuthCode) Expired() bool {
	return false
}
