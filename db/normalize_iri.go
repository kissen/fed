package db

import (
	"net/url"
	"github.com/PuerkitoBio/purell"
)

func normalizeIri(source *url.URL) string {
	flags := purell.FlagsUnsafeGreedy
	return purell.NormalizeURL(source, flags)
}
