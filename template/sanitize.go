package template

import (
	"github.com/microcosm-cc/bluemonday"
	"html/template"
)

func sanitize(html string) string {
	policy := bluemonday.NewPolicy()

	policy.AllowStandardURLs()
	policy.AllowAttrs("href").OnElements("a")
	policy.AllowElements("p")

	return policy.Sanitize(html)
}

func URL(s string) template.URL {
	sanitized := sanitize(s)
	return template.URL(sanitized)
}

func HTML(s string) template.HTML {
	sanitized := sanitize(s)
	return template.HTML(sanitized)
}
