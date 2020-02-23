package ap

import "net/url"

func urlEq(a, b *url.URL) bool {
	// if both are nil, they are the same

	if a == nil && b == nil {
		return true
	}

	// if just one is nil, they are not equal

	if a == nil && b != nil {
		return false
	}

	if a != nil && b == nil {
		return false
	}

	// we have two non-nil URLs; compare the relevant members

	return (a.Scheme == b.Scheme) && (a.Host == b.Host) && (a.Path == b.Path)
}
