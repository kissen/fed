package help

import (
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"testing"
)

func TestId_NilArg(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("did not panic")
		}
	}()

	var empty vocab.Type
	Id(empty)
}

func TestId_Note(t *testing.T) {
	note := streams.NewActivityStreamsNote()
	url := toUrl(t, "https://example.com/fed/123456")
	SetIdOn(note, url)

	id := Id(note)

	if id.Scheme != url.Scheme {
		t.Error("bad scheme")
	}

	if id.Host != url.Host {
		t.Error("bad hostname")
	}

	if id.Path != url.Path {
		t.Error("bad hostname")
	}
}
