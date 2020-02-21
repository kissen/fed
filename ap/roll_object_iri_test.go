package ap

import (
	"testing"
)

func TestRollObjectIri(t *testing.T) {
	for i := 0; i < 100; i++ {
		ctx := testSetUpContext()
		iri := rollObjectIri(ctx)

		if iri == nil {
			t.Error("got nil IRI")
		}
	}
}
