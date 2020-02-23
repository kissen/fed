package db

import (
	"github.com/PuerkitoBio/purell"
	"net/url"
)

func normalizeIri(source *url.URL) string {
	flags := purell.FlagsUnsafeGreedy
	return purell.NormalizeURL(source, flags)
}
