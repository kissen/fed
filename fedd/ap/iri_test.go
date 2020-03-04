package ap

import (
	"testing"
)

func Test_RollObjectIRI(t *testing.T) {
	for i := 0; i < 100; i++ {
		ctx := testSetUpContext()
		iri := RollObjectIRI(ctx)

		if _, err := iri.Object(); err != nil {
			t.Errorf("iri=%v not a valid Object", iri)
		}
	}
}

func Test_ActorIRI(t *testing.T) {
	ctx := testSetUpContext()
	iri := ActorIRI(ctx, "diomedes")

	if owner, err := iri.Actor(); err != nil {
		t.Fatal(err)
	} else if owner != "diomedes" {
		t.Fatal("bad owner")
	}
}

func Test_InboxIRI(t *testing.T) {
	ctx := testSetUpContext()
	iri := InboxIRI(ctx, "diomedes")

	if owner, err := iri.InboxOwner(); err != nil {
		t.Fatal(err)
	} else if owner != "diomedes" {
		t.Fatal("bad owner")
	}
}

func Test_OutboxIRI(t *testing.T) {
	ctx := testSetUpContext()
	iri := OutboxIRI(ctx, "diomedes")

	if owner, err := iri.OutboxOwner(); err != nil {
		t.Fatal(err)
	} else if owner != "diomedes" {
		t.Fatal("bad owner")
	}
}

func Test_IRIActor_Valid(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/belano", t)
	iri := IRI{ctx, url}
	actor, err := iri.Actor()

	if err != nil {
		t.Errorf("got err=%v for valid url=%v", err.Error(), url)
	}

	if actor != "belano" {
		t.Errorf("got actor=%v for url=%v", actor, url)
	}
}

func Test_IRIOutboxOwner_Valid(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/alice/outbox", t)
	iri := IRI{ctx, url}
	owner, err := iri.OutboxOwner()

	if err != nil {
		t.Errorf("got err=%v for valid url=%v", err.Error(), url)
	}

	if owner != "alice" {
		t.Errorf("got owner=%v for url=%v", owner, url)
	}
}

func Test_IRIOutboxOwner_EmptyUsername(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/outbox/", t)
	iri := IRI{ctx, url}
	_, err := iri.OutboxOwner()

	if err == nil {
		t.Errorf("got err=nil for empty outbox url=%v", url)
	}
}

func Test_IRIOutboxOwner_BadPrefix(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/alice/inbox", t)
	iri := IRI{ctx, url}
	_, err := iri.OutboxOwner()

	if err == nil {
		t.Errorf("got err=nil for invalid url=%v", url)
	}
}

func Test_IRIOutboxOwner_TooManyComponents(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/alice/outbox/blahfoo", t)
	iri := IRI{ctx, url}
	_, err := iri.OutboxOwner()

	if err == nil {
		t.Errorf("got err=nil for invalid url=%v", url)
	}
}

func Test_IRIInboxOwner_Valid(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/alice/inbox", t)
	iri := IRI{ctx, url}
	owner, err := iri.InboxOwner()

	if err != nil {
		t.Errorf("got err=%v for valid url=%v", err.Error(), url)
	}

	if owner != "alice" {
		t.Errorf("got owner=%v for url=%v", owner, url)
	}
}

func Test_IRIInboxOwner_ActivityIri(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/storage/12345", t)
	iri := IRI{ctx, url}
	owner, err := iri.InboxOwner()

	if err == nil {
		t.Errorf("got err=nil owner=%v for valid url=%v", owner, url)
	}
}

func Test_IRIInboxOwner_EmptyUsername(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/alice/", t)
	iri := IRI{ctx, url}
	_, err := iri.InboxOwner()

	if err == nil {
		t.Errorf("got err=nil for empty outbox url=%v", url)
	}
}

func Test_IRIInboxOwner_BadPrefix(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/alice/outbox", t)
	iri := IRI{ctx, url}
	_, err := iri.InboxOwner()

	if err == nil {
		t.Errorf("got err=nil for invalid url=%v", url)
	}
}

func Test_IRIInboxOwner_TooManyComponents(t *testing.T) {
	ctx := testSetUpContext()
	url := toUrl("http://example.com/fed/alice/inbox/blahfoo", t)
	iri := IRI{ctx, url}
	_, err := iri.InboxOwner()

	if err == nil {
		t.Errorf("got err=nil for invalid url=%v", url)
	}
}
