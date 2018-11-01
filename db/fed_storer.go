package db

type FedStorer interface {
	// Add a new user to the database.
	AddUser(username string) *FedUser

	// Return the metadata of a given user identified by its id.
	// Return nil on error.
	GetUser(userId uint64) *FedUser

	// Return the metadata of a given user identified by its
	// username. Returns nil on error.
	FindUser(username string) *FedUser

	// Add a post to the database.
	AddPost(userId uint64, content string) (*FedPost, error)

	// Return the data of a given post identified by its id.
	GetPost(postId uint64) *FedPost

	// Get all posts from a given user identified by the users id.
	GetPostsFrom(userId uint64) []*FedPost
}
