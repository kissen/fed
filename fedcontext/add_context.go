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
			log.Println("AddContext()")

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
				// fills out the Perms field
				setPermissionsOn(fc, r)

				// try creating a client; this is for the web interface
				if fc.Client, err = fc.NewClient(); err != nil {
					log.Println(err)
				}

				// install context into request
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Set the Perms field on fc. This is done by trying out all
// available authentication schemes one by one.
func setPermissionsOn(fc *FedContext, from *http.Request) {
	if setPermissionsFromBasiAuth(fc, from) {
		return
	}

	if setPermissionsFromCodeParam(fc, from) {
		return
	}

	if setPermissionsFromCodeCookie(fc, from) {
		return
	}

	if setPermissionsFromTokenParam(fc, from) {
		return
	}
}

// Try to set fc.Perms by looking at basic authentication headers
// on the HTTP request.
//
// Basic authentication doesn't seem common on the fediverse, but
// it is very convenient for debugging.
func setPermissionsFromBasiAuth(fc *FedContext, from *http.Request) (authed bool) {
	username, password, ok := from.BasicAuth()
	if !ok {
		return false
	}

	permissions, err := PermissionsFrom(from, username, password)
	if err != nil {
		log.Println("basic auth:", err)
		return false
	}

	fc.Perms = permissions

	log.Printf("authenticated user=%v with basic auth", fc.Perms.User.Name)

	return true
}

// Try to set fc.Perms by looking at the ?code= parameter in the
// request URI. This is mostly for API calls on the fediverse.
func setPermissionsFromCodeParam(fc *FedContext, from *http.Request) (authed bool) {
	code := from.URL.Query().Get("code")
	return setPermissionsFromCode(fc, code)
}

// Try to set fc.Perms by looking at the fc.Code property. fc.Code
// is part of CookieContext and as such persisted by web browsers.
// This is the authentication most users will use when interating
// with the web interface.
func setPermissionsFromCodeCookie(fc *FedContext, from *http.Request) (authed bool) {
	if code := fc.Code; code == nil {
		return false
	} else {
		return setPermissionsFromCode(fc, *code)
	}
}

// Given code, try to set fc.Perms accordingly if it is a valid and
// not expired OAuth code.
func setPermissionsFromCode(fc *FedContext, code string) (authed bool) {
	cm, err := fc.Storage.RetrieveCode(code)
	if err != nil {
		log.Printf("code auth: rejecting code=%v: %v", code, err)
		return false
	}

	user, err := fc.Storage.RetrieveUser(cm.Username)
	if err != nil {
		log.Printf("code auth: bad username=%v in code metadata: %v", cm.Username, err)
		return false
	}

	fc.Perms = &Permissions{
		User:   *user,
		Create: true,
		Like:   true,
	}

	log.Printf("authenticated user=%v with code auth", fc.Perms.User.Name)

	return true
}

// Try to set fc.Perms by looking at the ?token= parameter in the
// request URI. This is mostly for API calls on the fediverse.
func setPermissionsFromTokenParam(fc *FedContext, from *http.Request) (authed bool) {
	token := from.URL.Query().Get("token")
	if len(token) == 0 {
		return false
	}

	cm, err := fc.Storage.RetrieveToken(token)
	if err != nil {
		log.Printf("token auth: rejecting token=%v: %v", token, err)
		return false
	}

	user, err := fc.Storage.RetrieveUser(cm.Username)
	if err != nil {
		log.Printf("token auth: bad username=%v in token metadata: %v", cm.Username, err)
		return false
	}

	fc.Perms = &Permissions{
		User:   *user,
		Create: true,
		Like:   true,
	}

	log.Printf("authenticated user=%v with token auth", fc.Perms.User.Name)

	return true
}
