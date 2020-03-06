package main

import (
	"fmt"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	"html/template"
	"log"
	"net/url"
)

// WebVocab wraps a vocab.Type and provides various getters
// that are helpful when rendering that vocab.Type object as
// HTML onto a webpage.
type WebVocab interface {
	// plain getters
	Type() string
	Id() template.URL
	AttributedTo() template.URL
	Updated() string
	Name() template.HTML
	Content() template.HTML

	// helpers
	XFrom() string
}

type webVocab struct {
	target   vocab.Type
	mappings map[string]interface{}
}

func NewWebVocab(target vocab.Type) (WebVocab, error) {
	if target == nil {
		return nil, errors.New("target is nil")
	}

	mappings, err := target.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get mappings")
	}

	wocab := &webVocab{
		target:   target,
		mappings: mappings,
	}

	return wocab, nil
}

func NewWebVocabOnline(iri *url.URL) (WebVocab, error) {
	if iri == nil {
		return nil, errors.New("iri is nil")
	}

	obj, err := FetchIRI(iri)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read mappings")
	}

	wocab, err := NewWebVocab(obj)
	if err != nil {
		return nil, errors.Wrap(err, "derference ok, but wrapping failed")
	}

	return wocab, nil
}

func NewWebVocabs(collection vocab.ActivityStreamsCollection) ([]WebVocab, error) {
	items := collection.GetActivityStreamsItems()

	var ws []WebVocab

	for it := items.Begin(); it != items.End(); it = it.Next() {
		if it.IsIRI() {
			if w, err := NewWebVocabOnline(it.GetIRI()); err != nil {
				return nil, errors.Wrap(err, "wrapping IRI failed")
			} else {
				ws = append(ws, w)
			}
		} else {
			if w, err := NewWebVocab(it.GetType()); err != nil {
				return nil, errors.Wrap(err, "wrapping object failed")
			} else {
				ws = append(ws, w)
			}
		}
	}

	return ws, nil
}

func (v *webVocab) Type() string {
	return v.mapping("type")
}

func (v *webVocab) Id() template.URL {
	id := v.mapping("id")
	return template.URL(id)
}

func (v *webVocab) AttributedTo() template.URL {
	author := v.mapping("attributedTo")
	return template.URL(author)
}

func (v *webVocab) Updated() string {
	return v.mapping("updated")
}

func (v *webVocab) Name() template.HTML {
	html := v.mapping("name")
	return template.HTML(html)
}

func (v *webVocab) Content() template.HTML {
	html := v.mapping("content")
	return template.HTML(html)
}

func (v *webVocab) XFrom() string {
	if author, err := v.from(); err != nil {
		log.Println("could not determine author:", err)
		return "Anonymous"
	} else {
		return author
	}
}

func (v *webVocab) from() (string, error) {
	// first, get the server name from the id; that's easy and shouldn't
	// fail

	id, err := url.Parse(v.mapping("id"))
	if err != nil {
		return "", errors.Wrap(err, "bad id or hostname")
	}

	server := id.Hostname()

	// fetch the actor that is supposed to have authored this object

	actor, err := v.author()
	if err != nil {
		return "", errors.Wrap(err, "cannot identify author actor")
	}

	mappings, err := actor.Serialize()
	if err != nil {
		return "", errors.Wrap(err, "bad mappings")
	}

	// find out how we should call them

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

func (v *webVocab) author() (vocab.ActivityStreamsPerson, error) {
	// if the object itself is an author, that's easy

	if person, ok := v.target.(vocab.ActivityStreamsPerson); ok {
		return person, nil
	}

	// iterate through all fiels that might tell us something about
	// the author of this object

	var addr *string

	candidates := []string{
		"attributedTo", "actor",
	}

	for _, candidate := range candidates {
		if value := v.mapping(candidate); len(value) != 0 {
			addr = &value
			break
		}
	}

	if addr == nil {
		return nil, errors.New("cannot determine author")
	}

	// look up the actor

	iri, err := url.Parse(*addr)
	if err != nil {
		return nil, errors.Wrap(err, "bad actor IRI")
	}

	obj, err := FetchIRI(iri)
	if err != nil {
		return nil, errors.Wrap(err, "do not understand actor json")
	}

	person, ok := obj.(vocab.ActivityStreamsPerson)
	if !ok {
		kind := person.GetJSONLDType()
		return nil, fmt.Errorf("got wrong kind=%v of object", kind)
	}

	return person, nil
}

// Return the matching mapping if it is an atomic entry.
// Returns an empty string on error.
func (v *webVocab) mapping(key string) string {
	if s, ok := v.mappings[key].(string); !ok {
		return ""
	} else {
		return s
	}
}
