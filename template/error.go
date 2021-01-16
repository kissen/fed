package template

import (
	"github.com/kissen/httpstatus"
	"github.com/kissen/fed/fedcontext"
	"net/http"
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
