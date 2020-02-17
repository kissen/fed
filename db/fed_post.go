package db

// Represents a post made by a registered user.
type FedPost struct {
	// The unique identifier of this record.
	Id FedId

	// The Id of the FedUser that created authored this post.
	Author FedId

	// The conent of the post.
	Content string
}
