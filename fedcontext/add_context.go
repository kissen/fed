package fedcontext

import (
	"context"
	"github.com/go-fed/activity/pub"
	"gitlab.cs.fau.de/kissen/fed/db"
	"log"
	"net/http"
)

// Return a middleware function that installs a FedContext on the
// request. If the request contained cookies, these are loaded into the
// context.
func AddContext(s db.FedStorage, pa pub.FederatingActor, hf pub.HandlerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value(_REQUEST_CONTEXT_KEY) == nil {
				var err error

				// create context w/ default values
				fc := &FedContext{}

				fc.Storage = s
				fc.PubActor = pa
				fc.PubHandler = hf
				fc.Status = http.StatusOK

				// try to read cookies if there are any any
				if err = fc.LoadFromCookie(r); err != nil {
					log.Println(err)
				}

				// try creating a client
				if fc.Client, err = fc.NewClient(); err != nil {
					log.Println(err)
				}

				// install context into request
				c := context.WithValue(r.Context(), _REQUEST_CONTEXT_KEY, fc)
				r = r.WithContext(c)
			}

			next.ServeHTTP(w, r)
		})
	}
}
