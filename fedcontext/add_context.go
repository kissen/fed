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

				// create context with default values
				fc := &FedContext{}
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
				c := context.WithValue(r.Context(), _REQUEST_CONTEXT_KEY, fc)
				r = r.WithContext(c)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func setPermissionsOn(fc *FedContext, from *http.Request) {
	if setPermissionsFromBasiAuth(fc, from) {
		return
	}

	if setPermissionsFromCode(fc, from) {
		return
	}

	if setPermissionsFromToken(fc, from) {
		return
	}
}

func setPermissionsFromBasiAuth(fc *FedContext, from *http.Request) (authed bool) {
	username, password, ok := from.BasicAuth()
	if !ok {
		return false
	}

	user, err := fc.Storage.RetrieveUser(username)
	if err != nil {
		log.Println("basic auth:", err)
		return false
	}

	if !user.PasswordOK(password) {
		log.Printf("basic auth: bad password=%v for user=%v", username, user.Name)
		return false
	}

	fc.Perms = &Permissions{
		User:   *user,
		Create: true,
		Like:   true,
	}

	log.Printf("authenticated user=%v with basic auth", fc.Perms.User.Name)
	return true
}

func setPermissionsFromCode(fc *FedContext, from *http.Request) (authed bool) {
	code := from.URL.Query().Get("code")
	if len(code) == 0 {
		return false
	}

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

func setPermissionsFromToken(fc *FedContext, from *http.Request) (authed bool) {
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
