package template

import (
	"encoding/base64"
	"fmt"
	"github.com/go-fed/activity/streams/vocab"
	"gitlab.cs.fau.de/kissen/fed/errors"
	"gitlab.cs.fau.de/kissen/fed/fetch"
	"gitlab.cs.fau.de/kissen/fed/prop"
	"golang.org/x/sync/errgroup"
	"html/template"
	"log"
	"net/url"
	"time"
)

// Represents a wrapped instance of a vocab.Type object. WebVocab
// provides method for rendering out the underlying vocab object
// to some safe HTML representation.
type WebVocab interface {
	// Render out the HTML representation of the wrapped object.
	Fragment() template.HTML

	// Get the type of the underlying wrapped object.
	Type() string

	// Get the IRI of the underlying wrapped object.
	Id() template.URL
}

// Implementation of WebVocab.
//
// WebVocab has quite a lot more getters than just defined by WebVocab.
// We can use these getters when rendering the HTML fragments.
type webVocab struct {
	target   vocab.Type
	mappings map[string]interface{}
	fragment template.HTML
}

// Return a wrapped version of target.
func New(target vocab.Type) (WebVocab, error) {
	return wrap(target)
}

// Return wrapped versions of all targets. Returns an error
// if at least one of the conversations failed.
func News(targets ...vocab.Type) ([]WebVocab, error) {
	if ws, err := wraps(targets...); err != nil {
		return nil, err
	} else {
		wvs := make([]WebVocab, len(ws))
		for i := range ws {
			wvs[i] = ws[i]
		}
		return wvs, nil
	}
}

// Wrap target into a webVocab. This involves dereferencing
// target if it's just an IRI.
func wrap(target vocab.Type) (*webVocab, error) {
	// do not allow nil arguments
	if target == nil {
		return nil, errors.New("target is nil")
	}

	// serialize the object for quick access; this is a quick but
	// dirty way to get around the verbose go-fed api while hacking
	// on the proof of concept...
	mappings, err := target.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get mappings")
	}

	// create the base struct
	wocab := &webVocab{
		target:   target,
		mappings: mappings,
	}

	// pick the right template
	var page string
	switch v := target.(type) {
	case vocab.ActivityStreamsNote:
		page = "res/note.fragment.tmpl"

	case vocab.ActivityStreamsCreate:
		page = "res/create.fragment.tmpl"

	case vocab.ActivityStreamsPerson:
		page = "res/person.fragment.tmpl"

	case vocab.ActivityStreamsCollection:
		page = "res/collection_page.fragment.tmpl"

	case vocab.ActivityStreamsOrderedCollectionPage:
		page = "res/ordered_collection_page.fragment.tmpl"

	default:
		log.Printf("type=%v not implemented", prop.Type(v))
		page = "res/not_implemented.fragment.tmpl"
	}

	// render the template
	html, err := renderFragement(page, wocab)
	if err != nil {
		return nil, errors.Wrap(err, "cannot generate html")
	}
	wocab.fragment = html

	// return the now fully filled out struct
	return wocab, nil
}

// Wrap all targets into webVocabs. Returns an error if at least
// one of the conversations failed.
func wraps(targets ...vocab.Type) ([]*webVocab, error) {
	group := &errgroup.Group{}
	ws := make([]*webVocab, len(targets))

	for i, target := range targets {
		myi, mytarget := i, target

		group.Go(func() error {
			if w, err := wrap(mytarget); err != nil {
				return err
			} else {
				ws[myi] = w
				return nil
			}
		})
	}

	return ws, group.Wait()
}

// Fetch the ActivityPub object at iri and return a wraped version.
func Fetch(target *url.URL) (WebVocab, error) {
	// do not allow nil arguments

	if target == nil {
		return nil, errors.New("target is nil")
	}

	// get the object from the network

	obj, err := fetch.Fetch(target)
	if err != nil {
		return nil, errors.Wrap(err, "dereference failed")
	}

	// now that we have the object, wrap it like normal

	wocab, err := New(obj)
	if err != nil {
		return nil, errors.Wrap(err, "derference ok, but wrapping failed")
	}

	return wocab, nil
}

// Return the HTML fragment for embedding in HTML pages.
func (v *webVocab) Fragment() template.HTML {
	return v.fragment
}

// Return the type property.
func (v *webVocab) Type() string {
	return v.mapping("type")
}

// Return the id, that is the IRI pointing to the wrapped element.
func (v *webVocab) Id() template.URL {
	id := v.mapping("id")
	return URL(id)
}

// Return the name property.
func (v *webVocab) Name() template.HTML {
	html := v.mapping("name")
	return HTML(html)
}

// Return the content property.
func (v *webVocab) Content() template.HTML {
	html := v.mapping("content")
	return HTML(html)
}

// Return the published timestamp.
func (v *webVocab) Published() string {
	if t, err := time.Parse(time.RFC3339, v.mapping("published")); err != nil {
		return ""
	} else {
		return t.Format("2006/01/02 15:04")
	}
}

// Return a human-readable string that identifies the author of this
// object.
func (v *webVocab) XFrom() string {
	if author, err := v.qualifiedAuthor(); err != nil {
		log.Println(err)
		return "Anonymous"
	} else {
		return author
	}
}

// Return the Id property in base64.
func (v *webVocab) XIdBase64() string {
	return base64.StdEncoding.EncodeToString([]byte(v.Id()))
}

// If this object is some kind of collection, return the individual
// items in this collection as wrapped elements.
func (v *webVocab) XChildren() []*webVocab {
	if cs, err := v.children(); err != nil {
		log.Println(err)
		return nil
	} else {
		return cs
	}
}

func (v *webVocab) XObject() []*webVocab {
	if obj, err := v.object(); err != nil {
		log.Println(err)
		return nil
	} else {
		return obj
	}
}

func (v *webVocab) mapping(key string) string {
	if s, ok := v.mappings[key].(string); !ok {
		return ""
	} else {
		return s
	}
}

func (v *webVocab) qualifiedAuthor() (string, error) {
	// first, get the server name from the id; that's easy and shouldn't
	// fail

	id, err := url.Parse(v.mapping("id"))
	if err != nil {
		return "", errors.Wrap(err, "bad id")
	}

	server := id.Hostname()

	// fetch the actor that is supposed to have authored this object

	actor, err := v.actor()
	if err != nil {
		return "", errors.Wrapf(err, "cannot identify author of id=%v", id)
	}

	// find out how we should call them

	mappings, err := actor.Serialize()
	if err != nil {
		return "", errors.Wrap(err, "bad mappings")
	}

	candidates := []string{
		"preferedUsername", "name",
	}

	for _, candidate := range candidates {
		value := mappings[candidate]

		if s, ok := value.(string); ok {
			return fmt.Sprintf("%v@%v", s, server), nil
		}
	}

	// nothing found

	return "", errors.New("cannot infer username from mappings")
}

func (v *webVocab) actor() (vocab.ActivityStreamsPerson, error) {
	// if the object itself is an author, that's easy

	if person, ok := v.target.(vocab.ActivityStreamsPerson); ok {
		return person, nil
	}

	// iterate through all fiels that might tell us something about
	// the author of this object
	//
	// XXX: you should be using the go-fed accessors here!!!

	var addr *string

	candidates := []string{
		"attributedTo", "actor",
	}

	for _, candidate := range candidates {
		if s := v.mapping(candidate); len(s) != 0 {
			addr = &s
			break
		}
	}

	if addr == nil {
		return nil, errors.New("no known identifying field")
	}

	// look up the actor

	iri, err := url.Parse(*addr)
	if err != nil {
		return nil, errors.Wrap(err, "bad address")
	}

	actor, err := fetch.Fetch(iri)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch identified actor")
	}

	person, ok := actor.(vocab.ActivityStreamsPerson)
	if !ok {
		kind := person.GetJSONLDType()
		return nil, fmt.Errorf("got wrong kind=%v of object", kind)
	}

	return person, nil
}

func (v *webVocab) children() ([]*webVocab, error) {
	return v.wrapAll(v.target)
}

func (v *webVocab) object() ([]*webVocab, error) {
	type objecter interface {
		GetActivityStreamsObject() vocab.ActivityStreamsObjectProperty
	}

	if o, ok := v.target.(objecter); !ok {
		return nil, fmt.Errorf("%T doesn't have object property", v.target)
	} else {
		items := o.GetActivityStreamsObject()
		return v.wrapAll(items)
	}
}

// Iterate over property (with fetch.Begin) and wrap all returned
// objects.
func (v *webVocab) wrapAll(property interface{}) ([]*webVocab, error) {
	it, err := fetch.Begin(property)
	if err != nil {
		return nil, err
	}

	vs, err := fetch.FetchIters(it)
	if err != nil {
		return nil, errors.Wrap(err, "cannot fetch from collection")
	}

	return wraps(vs...)
}
