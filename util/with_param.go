package util

import "net/url"

// Return a deep copy of the URL at u with parameters
// key=value appended.
func WithParam(u *url.URL, key, value string) *url.URL {
	v := *u

	q := v.Query()
	q.Add(key, value)
	v.RawQuery = q.Encode()

	return &v
}
