package ap

import (
	"testing"
)

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

func TestParseActivityIdFromIri_ValidIri(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/activity/1592", t)
	id, err := parseActivityIdFromIri(ctx, url)

	if err != nil {
		t.Errorf("got err=%v for valid url=%v", err.Error(), url)
	}

	if id != 1592 {
		t.Errorf("got id=%v for url=%v", id, url)
	}
}

func TestParseActivityIdFromIri_EmptyId(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/activity/", t)
	_, err := parseActivityIdFromIri(ctx, url)

	if err == nil {
		t.Errorf("got err=nil for empty url=%v", url)
	}
}

func TestParseActivityIdFromIri_BadPrefix(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/ytivitca/1592", t)
	_, err := parseActivityIdFromIri(ctx, url)

	if err == nil {
		t.Errorf("got err=nil for invalid url=%v", url)
	}
}

func TestParseActivityIdFromIri_TooManyComponents(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/activity/1592/fooblah", t)
	_, err := parseActivityIdFromIri(ctx, url)

	if err == nil {
		t.Errorf("got err=nil for invalid url=%v", url)
	}
}

func TestParseActivityIdFromIri_NegativeId(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/activity/-13594", t)
	_, err := parseActivityIdFromIri(ctx, url)

	if err == nil {
		t.Errorf("got err=nil for invalid url=%v", url)
	}
}

func TestParseActivityIdFromIri_OutOfRange(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/activity/123456789012345678901234567890", t)
	_, err := parseActivityIdFromIri(ctx, url)

	if err == nil {
		t.Errorf("got err=nil for invalid url=%v", url)
	}
}
