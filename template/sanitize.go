package template

import (
	"github.com/microcosm-cc/bluemonday"
	"html/template"
)

// Given a string containing arbitrary HTML, return a sanitized
// version safe for embedding into output HTML.
func sanitize(html string) string {
	policy := bluemonday.NewPolicy()

	policy.AllowStandardURLs()
	policy.AllowAttrs("href").OnElements("a")
	policy.AllowElements("p")

	return policy.Sanitize(html)
}

// Sanitize s and return the results as type template.URL.
func URL(s string) template.URL {
	sanitized := sanitize(s)
	return template.URL(sanitized)
}

// Sanitize s and return the results as type template.HTML.
func HTML(s string) template.HTML {
	sanitized := sanitize(s)
	return template.HTML(sanitized)
}
