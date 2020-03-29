package fedcontext

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/go-fed/activity/pub"
	"gitlab.cs.fau.de/kissen/fed/db"
	"gitlab.cs.fau.de/kissen/fed/errors"
	"gitlab.cs.fau.de/kissen/fed/util"
	"net/http"
)

// The key used in the HTTP request under which we store our
// FedContext.
const _REQUEST_CONTEXT_KEY = "RequestContext"

// They name of the cookie that is used to persist CookieContext
// in the browser.
const _COOKIE_CONTEXT_KEY = "CookieContext"

// A struct that contains all data we attach to each HTTP request.
// Handlers and other functions that work on that request can use
// and update the contents to their liking.
type FedContext struct {
	RequestContext
	WebContext
	CookieContext
}

// Context relevant to every request, be it an request to some API or
// to the web interface.
type RequestContext struct {
	// The connection to the database you can use to read and write
	// metadata and activity pub objects.
	Storage db.FedStorage

	// The federating actor from the go-fed library. You can use it
	// to take care of activity pub requests to inbox and outbox.
	PubActor pub.FederatingActor

	// The HTTP handler function from the go-fed library. You can use
	// it to take care of all sorts of requests that request an
	// ActivityPub document type.
	PubHandler pub.HandlerFunc

	// The HTTP status code to reply with. After being set the
	// first time with function Status(), it will not change anymore.
	// This means a handler can (1) set the status and then (2)
	// just call another handler to take care of the request w/o
	// changing the HTTP status code.
	//
	// Initialized to 200.
	Status int

	// Currently logged in user for this session. Might be nil in
	// which case nobody is logged in.
	Client FedClient
}

// Volatile context of the web interface. This information is only valid
// for one request and is zeroed with every new request.
type WebContext struct {
	// The name of the tab that should be highlighted in the
	// navigation bar. If empty, not tab will be highlighted.
	Selected string

	// The title that should be used.
	Title string
}

// Persistent context of the web interface. This information is restored
// when installing the context.
type CookieContext struct {
	// Token that authenticates this users session. This is identical
	// to OAuth tokens. Might be nil.
	Token *string

	// Flashes to display on top of the page. Might be nil.
	Flashs []string

	// Warning (yellow) flashes to display on the top of the page.
	// Might be nil.
	Warnings []string

	// Error (red) flashes to display on the top of the page.
	// Might be nil.
	Errors []string
}

// Return whether a user is currently logged in.
func (fwc *FedContext) LoggedIn() bool {
	return fwc.Client != nil
}

// Load persisted fields from cookie into context. This function
// is called when installing the FedContext into the HTTP request.
func (cc *CookieContext) LoadFromCookie(r *http.Request) error {
	// check if a cookie is even set

	cookie, err := r.Cookie(_COOKIE_CONTEXT_KEY)
	if err == http.ErrNoCookie {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "error retrieving cookie")
	}

	// convert the base64 value of the cookie to json binary data

	text, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return errors.Wrap(err, "cookie value malformed base64")
	}

	// try to interpret the contents of the cookie

	var buf CookieContext

	if err := json.Unmarshal(text, &buf); err != nil {
		return errors.Wrap(err, "cookie unmarshal failed")
	}

	// set fields; we either want both username and actor iri set
	// or none of them

	cc.Flashs = buf.Flashs
	cc.Warnings = buf.Warnings
	cc.Errors = buf.Errors

	cc.Token = nil

	if tt, ok := util.Trim(buf.Token); ok {
		cc.Token = util.Just(tt)
	}

	return nil
}

// Write persisted fields from context into cookie. This way,
// once the user comes back, the persisted fields can be restored
// with a call to LoadFromCookie.
func (cc *CookieContext) WriteToCookie(w http.ResponseWriter) error {
	// conver buf to json
	text, err := json.Marshal(cc)
	if err != nil {
		return errors.Wrap(err, "cookie marshal failed")
	}

	// convert json to base64; cookies are picky about what
	// they contain after all
	encoded := base64.StdEncoding.EncodeToString(text)

	// build up cookie struct
	cookie := http.Cookie{
		Name:  _COOKIE_CONTEXT_KEY,
		Value: encoded,
		Path:  "/",
	}

	// send out cookie wit the response
	http.SetCookie(w, &cookie)

	return nil
}

// Set all flash slices to nil. It makes sense to call this method after we ensured
// that all flashes were shown to the user.
func (cc *CookieContext) ClearFlashes() {
	cc.Flashs = nil
	cc.Warnings = nil
	cc.Errors = nil
}

// Return the FedContext associated with HTTP request r.
func Context(r *http.Request) (fc *FedContext) {
	return From(r.Context())
}

// Return the FedContext in c.
func From(c context.Context) (fc *FedContext) {
	// if the request does not carry such a context, we forgot
	// to install it earlier; this would be a programming error
	return c.Value(_REQUEST_CONTEXT_KEY).(*FedContext)
}

// Hint the context that the specified tab should be highlighted
// in the web interface.
//
// If the tab was already set for this request, this call is a no-op.
func Selected(r *http.Request, tab string) {
	if fc := Context(r); len(fc.Selected) == 0 {
		fc.Selected = tab
	}
}

// Hint the context that the rendered webpage should have given
// title.
//
// If the title was already set for this request, this call is a
// no-op.
func Title(r *http.Request, title string) {
	if fc := Context(r); len(fc.Title) == 0 {
		fc.Title = title
	}
}

// Hint the context that the HTTP request should reply with given
// HTTP status code.
//
// If the status was already set for this request, this call is a
// no-op.
func Status(r *http.Request, status int) {
	if fc := Context(r); fc.Status == 0 {
		fc.Status = status
	}
}

// Add s to the list of flashes to display on the web interface.
func Flash(r *http.Request, s string) {
	fc := Context(r)
	fc.Flashs = append(fc.Flashs, s)
}

// Add s to the list of warnings to display on the web interface.
func FlashWarning(r *http.Request, s string) {
	fc := Context(r)
	fc.Warnings = append(fc.Warnings, s)
}

// Add s to the list of errors to display on the web interface.
func FlashError(r *http.Request, s string) {
	fc := Context(r)
	fc.Errors = append(fc.Errors, s)
}
