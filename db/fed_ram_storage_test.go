package db

import (
	"math/rand"
	"testing"
)

func TestFedStorageAddUser(t *testing.T) {
	storage := NewFedRamStorage()

	user, err := storage.AddUser("alice")

	if err != nil {
		t.Error("inserting into empty storage failed")
	}

	if user.Name != "alice" {
		t.Errorf("created user with wrong username expected=alice actual=%v", user.Name)
	}
}

func TestFedStorageFindUser_ExistingUsers(t *testing.T) {
	storage := NewFedRamStorage()

	// add users

	users := []string{
		"alice", "bob", "celia", "david", "emily", "franz",
	}

	for _, user := range users {
		if _, err := storage.AddUser(user); err != nil {
			t.Errorf("adding user=%v failed", user)
		}
	}

	// retrieve all users in random order

	rand.Shuffle(len(users), func(i, j int) {
		users[i], users[j] = users[j], users[i]
	})

	for _, user := range users {
		retrieved, err := storage.FindUser(user)

		if err != nil {
			t.Errorf("could not find added user=%v err=%v", user, err)
		}

		if retrieved.Name != user {
			t.Errorf("bad retrieved.Name expected=%v actual=%v", user, retrieved.Name)
		}
	}
}
