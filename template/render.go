package template

import (
	"github.com/kissen/fed/fedcontext"
	"html/template"
	"log"
	"net/http"
)

func Render(w http.ResponseWriter, r *http.Request, page string, data map[string]interface{}) {
	// fill in values that are (almost) always needed

	data = sumMaps(data)

	data["Context"] = fedcontext.Context(r)
	data["SubmitPrompt"] = submitPrompt()
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
