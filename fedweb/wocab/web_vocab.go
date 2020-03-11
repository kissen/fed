package wocab

import (
	"fmt"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/fedutil"
	"html/template"
	"log"
	"net/url"
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

	case vocab.ActivityStreamsPerson:
		page = "res/person.fragment.tmpl"

	case vocab.ActivityStreamsOrderedCollectionPage:
		page = "res/ordered_collection_page.fragment.tmpl"

	default:
		log.Printf("type=%v not implemented", fedutil.Type(v))
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

// Fetch the ActivityPub object at iri and return a wraped version.
func Fetch(target *url.URL) (WebVocab, error) {
	// do not allow nil arguments

	if target == nil {
		return nil, errors.New("target is nil")
	}

	// get the object from the network

	obj, err := fedutil.Fetch(target)
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

func (v *webVocab) Fragment() template.HTML {
	return v.fragment
}

func (v *webVocab) Type() string {
	return v.mapping("type")
}

func (v *webVocab) Id() template.URL {
	id := v.mapping("id")
	return URL(id)
}

// Return a human-readable string that identifies the author of this
// object.
func (v *webVocab) XFrom() string {
	if author, err := v.qualifiedAuthor(); err != nil {
		return "Anonymous"
	} else {
		return author
	}
}

// If this object is some kind of collection, return the individual
// items in this collection as wrapped elements.
func (v *webVocab) XChildren() []*webVocab {
	if cs, err := v.children(); err != nil {
		return nil
	} else {
		return cs
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

	actor, err := fedutil.FetchString(*addr)
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
	// make sure we are dealing with an ap collection

	var c vocab.ActivityStreamsOrderedCollectionPage
	var ok bool

	if c, ok = v.target.(vocab.ActivityStreamsOrderedCollectionPage); !ok {
		return nil, fmt.Errorf("type=%v not a collection", fedutil.Type(v.target))
	}

	// get the underlying items

	items, err := fedutil.FetchOrGet(c)
	if err != nil {
		return nil, errors.Wrap(err, "cannot dereference list")
	}

	// wrap all underlying objs

	var ws []*webVocab

	for _, item := range items {
		if w, err := New(item); err != nil {
			return nil, errors.Wrap(err, "cannot wrap retrieved list entry")
		} else {
			ws = append(ws, w.(*webVocab))
		}
	}

	return ws, nil
}
