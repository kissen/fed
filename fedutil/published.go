package fedutil

import (
	"fmt"
	"github.com/go-fed/activity/streams/vocab"
	"time"
)

// Try to get the "publish" property from object. On error,
// return the zero value for time.Time and an error explaining
// what happened.
func Published(object vocab.Type) (time.Time, error) {
	type publisher interface {
		GetActivityStreamsPublished() vocab.ActivityStreamsPublishedProperty
	}

	p, ok := object.(publisher)
	if !ok {
		return time.Time{}, fmt.Errorf("%T missing GetActivityStreamsPublished", object)
	}

	published := p.GetActivityStreamsPublished()
	if published == nil {
		return time.Time{}, fmt.Errorf("%T has nil published property", object)
	}

	return published.Get(), nil
}
