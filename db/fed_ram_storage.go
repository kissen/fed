package db

import "errors"

// Implements the db/FedStorer interface. Stores all data in volatile
// memory, i.e. after the process is killed, all data is gone. Only
// useful for debugging.
type FedRamStorage struct {
	lastId uint64
	users  map[uint64]*FedUser
	posts  map[uint64]*FedPost
}

func NewFedRamStorage() FedStorer {
	return &FedRamStorage{
		lastId: 0,
		users:  make(map[uint64]*FedUser),
		posts:  make(map[uint64]*FedPost),
	}
}

func (fs *FedRamStorage) AddUser(username string) *FedUser {
	id := fs.lastId
	fs.lastId += 1

	u := &FedUser{Id: id, Name: username}
	fs.users[id] = u

	return u
}

func (fs *FedRamStorage) GetUser(userId uint64) *FedUser {
	return fs.users[userId]
}

func (fs *FedRamStorage) FindUser(username string) *FedUser {
	for _, user := range fs.users {
		if user.Name == username {
			return user
		}
	}

	return nil
}

func (fs *FedRamStorage) AddPost(userId uint64, content string) (*FedPost, error) {
	if fs.GetUser(userId) == nil {
		err := errors.New("no such user")
		return nil, err
	}

	postId := fs.lastId
	fs.lastId += 1

	post := &FedPost{UserId: userId, Content: content}
	fs.posts[postId] = post

	return post, nil
}

func (fs *FedRamStorage) GetPost(postId uint64) *FedPost {
	return fs.posts[postId]
}

func (fs *FedRamStorage) GetPostsFrom(userId uint64) []*FedPost {
	var ret []*FedPost

	for _, post := range fs.posts {
		if post.UserId == userId {
			ret = append(ret, post)
		}
	}

	return ret
}
