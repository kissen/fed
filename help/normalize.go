package help

import (
	"github.com/PuerkitoBio/purell"
	"github.com/pkg/errors"
	"net/url"
)

// Return a new URL that is source with normalizations applied to it.
func Normalize(source *url.URL) (*url.URL, error) {
	flags := purell.FlagsUnsafeGreedy
	s := purell.NormalizeURL(source, flags)

	if normalized, err := url.Parse(s); err != nil {
		return nil, errors.Wrapf(err, "cannot normalize source=%v", source)
	} else {
		return normalized, nil
	}
}
