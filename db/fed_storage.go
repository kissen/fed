package db

type FedStorage interface {
	// Add a new user to the database.
	AddUser(username string) (*FedUser, error)

	// Return the metadata of a given user identified by its id.
	GetUser(userId FedId) (*FedUser, error)

	// Return the metadata of a given user identified by its
	// username.
	FindUser(username string) (*FedUser, error)

	// Add a post to the database.
	AddPost(userId FedId, content string) (*FedPost, error)

	// Return the data of a given post identified by its id.
	GetPost(postId FedId) (*FedPost, error)

	// Get all posts from a given user identified by the users id.
	GetPostsFrom(userId FedId) ([]*FedPost, error)
}
