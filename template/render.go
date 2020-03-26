package template

import (
	"github.com/kissen/httpstatus"
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
	"gitlab.cs.fau.de/kissen/fed/fetch"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

// Write out the Error template with given status and cause.
// cause may be left nil.
func Error(w http.ResponseWriter, r *http.Request, status int, cause error, data map[string]interface{}) {
	// set up data for the error handler

	errorData := map[string]interface{}{
		"Status":      status,
		"StatusText":  http.StatusText(status),
		"Description": httpstatus.Describe(status),
	}

	if cause != nil {
		errorData["Cause"] = cause.Error()
	}

	renderData := sumMaps(data, errorData)

	// render with correct status

	fedcontext.Status(r, status)
	Render(w, r, "res/error.page.tmpl", renderData)
}

// Write out a page showing remote content at addr.
func Remote(w http.ResponseWriter, r *http.Request, iri *url.URL) {
	// fetch and wrap object

	wrapped, err := Fetch(iri)
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

func Render(w http.ResponseWriter, r *http.Request, page string, data map[string]interface{}) {
	// fill in values that are (almost) always needed

	data = sumMaps(data)

	data["Context"] = fedcontext.Context(r)
	data["SubmitPrompt"] = SubmitPrompt()
	data["Flashs"] = fedcontext.Context(r).Flashs
	data["Warnings"] = fedcontext.Context(r).Warnings
	data["Errors"] = fedcontext.Context(r).Errors

	// load template files

	templates := []string{
		page, "res/base.layout.tmpl", "res/card.fragment.tmpl",
		"res/flash.fragment.tmpl",
	}

	// compile template; if this fails it's a programming error

	ts, err := template.ParseFiles(templates...)
	if err != nil {
		log.Printf("parsing templates failed: %v", err)
		return
	}

	// set cookie for next time

	context := fedcontext.Context(r)
	context.ClearFlashes()
	context.WriteToCookie(w)

	// write http status

	status := fedcontext.Context(r).Status
	w.WriteHeader(status)

	// write body

	if err := ts.Execute(w, data); err != nil {
		log.Printf("executing template failed: %v", err)
		return
	}
}
