package db

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
)

// Implements the db/FedStorage interface. Stores all data in volatile
// memory, i.e. after the process is killed, all data is gone. Only
// useful for debugging.
type FedRamStorage struct {
	lastId FedId
	users  map[FedId]*FedUser
	posts  map[FedId]*FedPost
}

func NewFedRamStorage() FedStorage {
	return &FedRamStorage{
		lastId: 0,
		users:  make(map[FedId]*FedUser),
		posts:  make(map[FedId]*FedPost),
	}
}

func (fs *FedRamStorage) nextId() FedId {
	id := fs.lastId
	fs.lastId += 1

	return id
}

func (fs *FedRamStorage) AddUser(username string) (*FedUser, error) {
	id := fs.nextId()

	user := &FedUser{Id: id, Name: username}
	fs.users[id] = user

	log.Printf("created FedUser w/ Id=%v Name='%v'", user.Id, user.Name)

	return user, nil
}

func (fs *FedRamStorage) GetUser(userId FedId) (*FedUser, error) {
	if user, ok := fs.users[userId]; ok {
		return user, nil
	} else {
		return nil, fmt.Errorf("no user with userId=%v", userId)
	}
}

func (fs *FedRamStorage) FindUser(username string) (*FedUser, error) {
	for _, user := range fs.users {
		if user.Name == username {
			return user, nil
		}
	}

	return nil, fmt.Errorf("no user with username=%v", username)
}

func (fs *FedRamStorage) AddPost(userId FedId, content string) (*FedPost, error) {
	if _, err := fs.GetUser(userId); err != nil {
		return nil, errors.Wrapf(err, "unknown author")
	}

	postId := fs.nextId()

	post := &FedPost{Id: postId, Author: userId, Content: content}
	fs.posts[postId] = post

	log.Printf("created FedPost w/ Id=%v Content='%v'", post.Id, post.Content)

	return post, nil
}

func (fs *FedRamStorage) GetPost(postId FedId) (*FedPost, error) {
	if post, ok := fs.posts[postId]; ok {
		return post, nil
	} else {
		return nil, fmt.Errorf("no post with postId=%v", postId)
	}
}

func (fs *FedRamStorage) GetPostsFrom(userId FedId) ([]*FedPost, error) {
	var ret []*FedPost

	for _, post := range fs.posts {
		if post.Author == userId {
			ret = append(ret, post)
		}
	}

	return ret, nil
}
