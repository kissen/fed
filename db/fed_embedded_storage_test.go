package db

import (
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"gitlab.cs.fau.de/kissen/fed/help"
	"net/url"
	"os"
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
			Name:      "alice",
			Liked:     liked,
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

func testNote(t *testing.T) vocab.ActivityStreamsNote {
	name := streams.NewActivityStreamsNameProperty()
	name.AppendXMLSchemaString("Answer July")

	content := streams.NewActivityStreamsContentProperty()
	content.AppendXMLSchemaString("Here – said the Year –")

	note := streams.NewActivityStreamsNote()
	note.SetActivityStreamsName(name)
	note.SetActivityStreamsContent(content)

	return note
}

func TestStoreAndRetrieveNote(t *testing.T) {
	storage := FedEmbeddedStorage{
		Filepath: dbPath(t),
	}

	// create db

	if err := storage.Open(); err != nil {
		t.Fatalf("open failed with err=%v", err)
	}

	defer deleteDbPath(t)

	// put note

	iri := toUrl(t, "https://example.com/poetry/emily/july")
	note := testNote(t)

	if err := storage.StoreObject(iri, note); err != nil {
		t.Errorf("refusing to store object err=%v", err)
	}

	// close db

	if err := storage.Close(); err != nil {
		t.Fatalf("close failed with err=%v", err)
	}

	// re-open db

	if err := storage.Open(); err != nil {
		t.Fatalf("open failed with err=%v", err)
	}

	// get object

	var parsed vocab.ActivityStreamsNote
	var ok bool

	if obj, err := storage.RetrieveObject(iri); err != nil {
		t.Fatalf("getting previously added note failed err=%v", err)
	} else if parsed, ok = obj.(vocab.ActivityStreamsNote); !ok {
		t.Fatalf("retrieved type does not match stored type obj=%v TypeName=", obj)
	}

	origName := help.Name(note)
	parsedName := help.Name(parsed)

	if origName != parsedName {
		t.Errorf("got bad name expected=%v got=%v", origName, parsedName)
	}

	origContent := help.Content(note)
	parsedContent := help.Content(parsed)

	if origContent != parsedContent {
		t.Errorf("got bad content expected=%v got=%v", origName, parsedName)
	}
}
