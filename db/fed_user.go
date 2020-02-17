package db

// Represents a user registered with the service.
type FedUser struct {
	// The unique identifier of this record.
	Id FedId

	// The name of the user. Currently only alphanumeric ASCII
	// will work.
	Name string
}
