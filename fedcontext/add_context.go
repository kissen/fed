package fedcontext

import (
	"context"
	"github.com/go-fed/activity/pub"
	"github.com/kissen/fed/db"
	"github.com/kissen/fed/fediri"
	"github.com/kissen/fed/util"
	"log"
	"net/http"
)

// Return a middleware function that installs a FedContext on the
// request. If the request contained cookies, these are loaded into the
// context.
func AddContext(s db.FedStorage, pa pub.FederatingActor, hf pub.HandlerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(hw http.ResponseWriter, r *http.Request) {
			log.Println("AddContext()")

			// we want to take a look at the resulting status later;
			// because of this we wrap the response writer in a util
			// struct that records the http status

			w := &util.StatusHTTPWriter{
				Target: hw,
			}

			// install the FedContext if necessary; this involves setting
			// default values, starting a database transaction for the request,
			// loading information from cookies and more

			if r.Context().Value(_REQUEST_CONTEXT_KEY) == nil {
				var err error

				// create and install context
				fc := &FedContext{}
				c := context.WithValue(r.Context(), _REQUEST_CONTEXT_KEY, fc)
				r = r.WithContext(c)

				// set default values
				fc.Storage = s
				fc.PubActor = pa
				fc.PubHandler = hf
				fc.Status = http.StatusOK

				// try to load contents from cookie if there is one;
				// fills out the CookieContext fields
				if err = fc.LoadFromCookie(r); err != nil {
					log.Println(err)
				}

				// try to find out the permissions of the request;
				// fills out the Client field
				setClientOn(fc, r)
			}

			// handle the request; this is the core of the application

			next.ServeHTTP(w, r)
		})
	}
}

// Set the Client field on fc. This is done by trying out all
// available authentication schemes one by one.
func setClientOn(fc *FedContext, from *http.Request) {
	if setClientFromBasiAuth(fc, from) {
		return
	}

	if setClientFromTokenCookie(fc, from) {
		return
	}

	if setClientFromTokenParam(fc, from) {
		return
	}
}

// Try to set fc.Client by looking at basic authentication headers
// on the HTTP request.
//
// Basic authentication doesn't seem common on the fediverse, but
// it is very convenient for debugging.
func setClientFromBasiAuth(fc *FedContext, from *http.Request) (authed bool) {
	// get credentials; if none were supplied give up right away

	username, password, ok := from.BasicAuth()
	if !ok {
		return false
	}

	// create token for each basic auth; this is honestly kind of
	// dumb but it will have to do for now

	tm, err := db.NewFedOAuthToken(username, password, fc.Storage)
	if err != nil {
		log.Println(err)
		return false
	}

	if authed := setClientFromToken(fc, tm.Token); !authed {
		log.Println("NOT")
		return false
	}

	log.Printf("authenticated user=%v with basic auth", username)
	return true
}

// Try to set fc.Client by looking at the fc.Token property. fc.Token
// is part of CookieContext and as such persisted by web browsers.
// This is the authentication most users will use when interating
// with the web interface.
func setClientFromTokenCookie(fc *FedContext, from *http.Request) (authed bool) {
	if token := fc.Token; token == nil {
		return false
	} else {
		return setClientFromToken(fc, *token)
	}
}

// Try to set fc.Client by looking at the ?token= parameter in the
// request URI. This is mostly for API calls on the fediverse.
func setClientFromTokenParam(fc *FedContext, from *http.Request) (authed bool) {
	token := from.URL.Query().Get("token")
	return setClientFromToken(fc, token)
}

// Given token, try to set fc.Client accordingly if it is a valid and
// not expired OAuth token.
func setClientFromToken(fc *FedContext, token string) (authed bool) {
	// trim token; if empty don't even bother trying
	tt, ok := util.Trim(&token)
	if !ok {
		return false
	}

	// get the token meta data
	cm, err := fc.Storage.RetrieveToken(tt)
	if err != nil {
		return false
	}

	// build up client
	addr := fediri.ActorIRI(cm.Username).String()
	client, err := NewRemoteClient(addr, tt)
	if err != nil {
		return false
	}

	// set and return
	fc.Client = client
	log.Printf("authenticated user=%v with token auth", cm.Username)
	return true
}
