package db

import (
	"os"
	"net/url"
	//"github.com/pkg/errors"
	"path/filepath"
	"testing"
)

func dbPath(t *testing.T) string {
	cachedir, err := os.UserCacheDir()

	if err != nil {
		t.Fatalf("cannot determine cache directory err=%v", err)
	}

	return filepath.Join(cachedir, "fed.db")
}

func deleteDbPath(t *testing.T) {
	if err := os.Remove(dbPath(t)); err != nil {
		t.Fatalf("cannot remove file err=%v", err)
	}
}

func TestOpenAndClose(t *testing.T) {
	storage := FedEmbeddedStorage{
		Filepath: dbPath(t),
	}

	// create db

	if err := storage.Open(); err != nil {
		t.Fatalf("open failed with err=%v", err)
	}

	defer deleteDbPath(t)

	// finish

	if err := storage.Close(); err != nil {
		t.Fatalf("close failed with err=%v", err)
	}
}

func TestUserBucket(t *testing.T) {
	storage := FedEmbeddedStorage{
		Filepath: dbPath(t),
	}

	// create db

	if err := storage.Open(); err != nil {
		t.Fatalf("open failed with err=%v", err)
	}

	defer deleteDbPath(t)

	// put user

	follower := "https://example.com/ap/lewis.caroll"

	{
		liked := []*url.URL{
			toUrl(t, "https://example.com/ice"),
			toUrl(t, "https://hacks.moe/privacy"),
		}

		followers := []*url.URL{
			toUrl(t, follower),
		}

		user := FedUser{
			Name: "alice",
			Liked: liked,
			Followers: followers,
		}

		if err := storage.StoreUser(&user); err != nil {
			t.Fatalf("storing new user failed err=%v", err)
		}
	}

	// get user

	{
		user, err := storage.RetrieveUser("alice")

		if err != nil {
			t.Fatalf("retrieving previously added user failed err=%v", err)
		}

		if user.Name != "alice" {
			t.Errorf("got wrong username expected=alice got=%v", user.Name)
		}

		if count := len(user.Inbox); count != 0 {
			t.Errorf("bad number of items in inbox expected=0 got=%v", count)
		}

		if count := len(user.Liked); count != 2 {
			t.Errorf("bad number of liked objects expected=2 got=%v", count)
		}

		if count := len(user.Followers); count != 1 {
			t.Errorf("bad number of followers expected=1 got=%v", count)
		} else {
			followerOrig := follower
			followerDeserialized := user.Followers[0].String()

			if followerOrig != followerDeserialized {
				t.Errorf("bad follower expected=%v got=%v", followerOrig, followerDeserialized)
			}
		}
	}

	// finish

	if err := storage.Close(); err != nil {
		t.Fatalf("close failed with err=%v", err)
	}
}
