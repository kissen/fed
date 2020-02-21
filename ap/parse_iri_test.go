package ap

import (
	"testing"
)

func TestParseActorFromIri_ValidIri(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/actor/belano", t)
	actor, err := parseActorFromIri(ctx, url)

	if err != nil {
		t.Errorf("got err=%v for valid url=%v", err.Error(), url)
	}

	if actor != "belano" {
		t.Errorf("got actor=%v for url=%v", actor, url)
	}
}

func TestParseOutboxOwnerFromIri_ValidIri(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/outbox/alice", t)
	owner, err := parseOutboxOwnerFromIri(ctx, url)

	if err != nil {
		t.Errorf("got err=%v for valid url=%v", err.Error(), url)
	}

	if owner != "alice" {
		t.Errorf("got owner=%v for url=%v", owner, url)
	}
}

func TestParseOutboxOwnerFromIri_EmptyUsername(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/outbox/", t)
	_, err := parseOutboxOwnerFromIri(ctx, url)

	if err == nil {
		t.Errorf("got err=nil for empty outbox url=%v", url)
	}
}

func TestParseOutboxOwnerFromIri_BadPrefix(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/inbox/alice", t)
	_, err := parseOutboxOwnerFromIri(ctx, url)

	if err == nil {
		t.Errorf("got err=nil for invalid url=%v", url)
	}
}

func TestParseOutboxOwnerFromIri_TooManyComponents(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/outbox/alice/blahfoo", t)
	_, err := parseOutboxOwnerFromIri(ctx, url)

	if err == nil {
		t.Errorf("got err=nil for invalid url=%v", url)
	}
}

func TestParseInboxOwnerFromIri_ValidIri(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/inbox/alice", t)
	owner, err := parseInboxOwnerFromIri(ctx, url)

	if err != nil {
		t.Errorf("got err=%v for valid url=%v", err.Error(), url)
	}

	if owner != "alice" {
		t.Errorf("got owner=%v for url=%v", owner, url)
	}
}

func TestParseInboxOwnerFromIri_ActivityIri(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/activity/1234", t)
	owner, err := parseInboxOwnerFromIri(ctx, url)

	if err == nil {
		t.Errorf("got err=nil owner=%v for valid url=%v", owner, url)
	}
}

func TestParseInboxOwnerFromIri_EmptyUsername(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/inbox/", t)
	_, err := parseInboxOwnerFromIri(ctx, url)

	if err == nil {
		t.Errorf("got err=nil for empty outbox url=%v", url)
	}
}

func TestParseInboxOwnerFromIri_BadPrefix(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/outbox/alice", t)
	_, err := parseInboxOwnerFromIri(ctx, url)

	if err == nil {
		t.Errorf("got err=nil for invalid url=%v", url)
	}
}

func TestParseInboxOwnerFromIri_TooManyComponents(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/inbox/alice/blahfoo", t)
	_, err := parseInboxOwnerFromIri(ctx, url)

	if err == nil {
		t.Errorf("got err=nil for invalid url=%v", url)
	}
}
