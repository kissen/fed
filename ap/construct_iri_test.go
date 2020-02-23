package ap

import (
	"testing"
)

func TestConstructActorIri(t *testing.T) {
	ctx := testSetUpContext()
	owner := "ulises"
	expected := "http://example.com/fed/actor/ulises"

	if actual := constructActorIri(ctx, owner); expected != actual.String() {
		t.Errorf("got bad actor IRI; expected=%v actual=%v", expected, actual)
	}
}

func TestConstructOutboxIri(t *testing.T) {
	ctx := testSetUpContext()
	owner := "ulises"
	expected := "http://example.com/fed/outbox/ulises"

	if actual := constructOutboxIri(ctx, owner); expected != actual.String() {
		t.Errorf("got bad inbox IRI; expected=%v actual=%v", expected, actual)
	}
}

func TestConstructInboxIri(t *testing.T) {
	ctx := testSetUpContext()
	owner := "ulises"
	expected := "http://example.com/fed/inbox/ulises"

	if actual := constructInboxIri(ctx, owner); expected != actual.String() {
		t.Errorf("got bad outbox IRI; expected=%v actual=%v", expected, actual)
	}
}
