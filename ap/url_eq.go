package ap

import "net/url"

// Return whether a and b are equal.
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

//
// scary O(n^2) runtimes ahead
//

// Return whether haystack contains an URL that we consider equal
// to needle.
func urlIn(needle *url.URL, haystack []*url.URL) bool {
	for _, hay := range haystack {
		if urlEq(hay, needle) {
			return true
		}
	}

	return false
}

// Return whether any of the slices in haystacks contains an URL that
// we consider equal to needle.
func urlInAny(needle *url.URL, haystacks [][]*url.URL) bool {
	for _, haystack := range haystacks {
		if urlIn(needle, haystack) {
			return true
		}
	}

	return false
}
