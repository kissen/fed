package template

import (
	"gitlab.cs.fau.de/kissen/fed/fetch"
	"net/http"
)

// Write out a page showing activity pub content accessible via iter.
func Iter(w http.ResponseWriter, r *http.Request, it fetch.Iter) {
	// fetch objects

	vs, err := fetch.FetchIters(it)
	if err != nil {
		Error(w, r, http.StatusBadGateway, err, nil)
		return
	}

	// wrap objects

	wrapped, err := News(vs...)
	if err != nil {
		Error(w, r, http.StatusBadGateway, err, nil)
		return
	}

	// set up data dict and render

	data := map[string]interface{}{
		"Items": wrapped,
	}

	Render(w, r, "res/collection.page.tmpl", data)
}
