package template

import (
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
	"net/http"
	"net/url"
)

// Write out a page showing remote content at addr.
func Remote(w http.ResponseWriter, r *http.Request, iri *url.URL) {
	// fetch and wrap object

	fc := fedcontext.Context(r)
	wrapped, err := Fetch(fc, iri)
	if err != nil {
		Error(w, r, http.StatusBadGateway, err, nil)
		return
	}

	// set up data dict and render

	data := map[string]interface{}{
		"Items": []WebVocab{
			wrapped,
		},
	}

	fedcontext.Title(r, iri.String())
	Render(w, r, "res/collection.page.tmpl", data)
}
