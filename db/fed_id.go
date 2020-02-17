package db

// In our system, each piece of data, be it a user, a post, a picture
// and so on is identified by one global FedId which is just an
// unsigned integer.
type FedId uint64
