package fedutil

import (
	"github.com/go-fed/activity/streams/vocab"
	"log"
	"net/url"
	"sync"
	"time"
)

const _EXPIRATION_TIME = 30 * time.Second
const _GC_INTERVAL = 3 * time.Minute

// VocabCache provides a cache for vocab.Type objects retrieved via HTTP.
type VocabCache interface {
	// Get the object at iri. If it has recently been fetchted, a cached
	// version is returned instead, saving on network traffic.
	Get(iri *url.URL) (vocab.Type, error)

	// Close the cache, that is stop cleaning up old elements from
	// the cache.
	Close() error
}

type vocabCache struct {
	sync.Mutex
	contents map[string]*entry
	closed   bool
}

// Create a new cache. Calling this function also starts a goroutine
// that regularry cleans up the cache. If you do not want to use this
// cache anymore and get rid of that goroutine, call Close().
func NewVocabCache() VocabCache {
	ch := &vocabCache{contents: make(map[string]*entry)}
	go ch.collectGarbage()
	return ch
}

func (vc *vocabCache) Get(iri *url.URL) (vocab.Type, error) {
	key := iri.String()

	// only one concurrent acces of the map at a given time

	vc.Lock()

	// get entry from map; maybe we can just return the value right away

	ent, ok := vc.contents[key]
	if ok && !ent.Expired() {
		vc.Unlock()
		return ent.Get()
	}

	// no entry or only an expired one; create a new one

	ent = &entry{
		object:      nil,
		err:         nil,
		present:     false,
		requestedOn: time.Now(),
	}

	ent.cond = sync.NewCond(&sync.Mutex{})

	// add entry to map for retrival by other goroutines

	vc.contents[key] = ent

	// start the http request

	go ent.Fill(iri)

	// let others touch the map while we wait for the result

	vc.Unlock()
	return ent.Get()
}

func (vc *vocabCache) Close() error {
	vc.closed = true
	return nil
}

// Regularlly remove objects from the cache we consider expired.
func (vc *vocabCache) collectGarbage() {
	for !vc.closed {
		vc.Lock()

		for key, value := range vc.contents {
			if value.Expired() {
				log.Println("Evicting ", key)
				delete(vc.contents, key)
			}
		}

		vc.Unlock()

		time.Sleep(_GC_INTERVAL)
	}
}

type entry struct {
	// result form fetching
	object vocab.Type
	err    error

	// true once fetching was complted
	present bool

	// time on which fetching this resource started;
	// used to determine whether an entry is expired
	requestedOn time.Time

	// synchronization aids that allow goroutines to wait
	// on the result
	cond *sync.Cond
}

// Go and fill the underlying entry with iri.
func (e *entry) Fill(iri *url.URL) {
	// actually fetch the object from the network

	if raw, err := Get(iri); err != nil {
		e.object = nil
		e.err = err
	} else if obj, err := BytesToVocab(raw); err != nil {
		e.object = nil
		e.err = err
	} else {
		e.object = obj
		e.err = nil
	}

	e.cond.L.Lock()
	defer e.cond.L.Unlock()

	e.present = true
	e.cond.Broadcast()
}

// Return the underlying object/error or wait for it first
// if it isn't available yet.
//
// This is essentially a future-style object.
func (e *entry) Get() (vocab.Type, error) {
	e.cond.L.Lock()

	for !e.present {
		e.cond.Wait()
	}

	e.cond.L.Unlock()

	return e.object, e.err
}

// Returns whether this entry is expired.
func (e *entry) Expired() bool {
	expiredOn := e.requestedOn.Add(_EXPIRATION_TIME)
	return time.Now().After(expiredOn)
}
