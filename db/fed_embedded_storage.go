package db

import (
	"encoding/json"
	"fmt"
	"github.com/go-fed/activity/streams/vocab"
	"gitlab.cs.fau.de/kissen/fed/errors"
	"gitlab.cs.fau.de/kissen/fed/marshal"
	"go.etcd.io/bbolt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const _GARBAGE_COLLECTION_WAIT = 1 * time.Minute

var _USERS_BUCKET = []byte("Users")
var _CODES_BUCKET = []byte("OAuth/Codes")
var _TOKENS_BUCKET = []byte("OAuth/Tokens")
var _DOCUMENTS_BUCKET = []byte("Documents")

type FedEmbeddedStorage struct {
	Filepath   string
	connection *bbolt.DB
	closed     bool
}

func (fs *FedEmbeddedStorage) Open() (err error) {
	log.Println("Open()")

	// open db

	fs.connection, err = bbolt.Open(fs.Filepath, 0600, nil)
	if err != nil {
		return errors.Wrapf(err, "open db at Filepath=%v failed", fs.Filepath)
	}

	// create buckets

	err = fs.connection.Update(func(tx *bbolt.Tx) error {
		buckets := [][]byte{
			_USERS_BUCKET, _CODES_BUCKET, _TOKENS_BUCKET, _DOCUMENTS_BUCKET,
		}

		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		fs.connection.Close()
		return errors.Wrap(err, "intializing buckets failed")
	}

	// start garbage collection; it will run until Close

	go fs.gcLoop()

	// success

	return nil
}

func (fs *FedEmbeddedStorage) Close() error {
	log.Println("Close()")

	if fs.closed {
		return errors.New("database was already closed")
	}

	fs.closed = true
	return fs.connection.Close()
}

func (fs *FedEmbeddedStorage) RetrieveUser(username string) (user *FedUser, err error) {
	log.Printf("RetrieveUser(%s)", username)

	bytes, err := fs.retrieve(_USERS_BUCKET, username)
	if err != nil {
		return nil, err
	}

	if user, err = bytesToUser(bytes); err != nil {
		return nil, errors.Wrap(err, "deserializing user failed")
	}

	return user, err
}

func (fs *FedEmbeddedStorage) StoreUser(user *FedUser) error {
	log.Printf("StoreUser(Name=%v #Inbox=%v #Outbox=%v)", user.Name, len(user.Inbox), len(user.Outbox))

	bytes, err := userToBytes(user)
	if err != nil {
		return errors.Wrap(err, "could not serialize user")
	}

	return fs.store(_USERS_BUCKET, user.Name, bytes)
}

func (fs *FedEmbeddedStorage) RetrieveCode(code string) (*FedOAuthCode, error) {
	log.Printf("RetrieveCode(%s)", code)

	bs, err := fs.retrieve(_CODES_BUCKET, code)
	if err != nil {
		return nil, err
	}

	var c FedOAuthCode
	if err := json.Unmarshal(bs, &c); err != nil {
		return nil, errors.Wrap(err, "deserializing code failed")
	}

	return &c, nil
}

func (fs *FedEmbeddedStorage) StoreCode(code *FedOAuthCode) error {
	log.Printf("StoreCode(Code=%v)", code.Code)

	bs, err := json.Marshal(code)
	if err != nil {
		return errors.Wrap(err, "serializing code failed")
	}

	return fs.store(_CODES_BUCKET, code.Code, bs)
}

func (fs *FedEmbeddedStorage) RetrieveToken(token string) (*FedOAuthToken, error) {
	log.Printf("RetrieveToken(%s)", token)

	bs, err := fs.retrieve(_TOKENS_BUCKET, token)
	if err != nil {
		return nil, err
	}

	var c FedOAuthToken
	if err := json.Unmarshal(bs, &c); err != nil {
		return nil, errors.Wrap(err, "deserializing token failed")
	}

	return &c, nil
}

func (fs *FedEmbeddedStorage) StoreToken(token *FedOAuthToken) error {
	log.Printf("StoreToken(Token=%v)", token.Token)

	bs, err := json.Marshal(token)
	if err != nil {
		return errors.Wrap(err, "serializing token failed")
	}

	return fs.store(_TOKENS_BUCKET, token.Token, bs)
}

func (fs *FedEmbeddedStorage) RetrieveObject(iri *url.URL) (obj vocab.Type, err error) {
	log.Printf("RetrieveObject(%v)", iri)

	bytes, err := fs.retrieve(_DOCUMENTS_BUCKET, fs.toKey(iri))
	if err != nil {
		return nil, err
	}

	if obj, err = marshal.BytesToVocab(bytes); err != nil {
		return nil, errors.Wrap(err, "deserializing object failed")
	}

	return obj, err
}

func (fs *FedEmbeddedStorage) StoreObject(iri *url.URL, obj vocab.Type) error {
	log.Printf("StoreObject(%v)", iri)

	bytes, err := marshal.VocabToBytes(obj)
	if err != nil {
		return errors.Wrap(err, "could not serialize object")
	}

	return fs.store(_DOCUMENTS_BUCKET, fs.toKey(iri), bytes)
}

func (fs *FedEmbeddedStorage) DeleteObject(iri *url.URL) error {
	log.Printf("DeleteObject(%v)", iri)

	return fs.connection.Update(func(tx *bbolt.Tx) error {
		var bucket *bbolt.Bucket
		var updateErr error

		if bucket = tx.Bucket(_DOCUMENTS_BUCKET); bucket == nil {
			return errors.New("could not open documents bucket")
		}

		key := []byte(fs.toKey(iri))

		if updateErr = bucket.Delete(key); updateErr != nil {
			return errors.Wrap(updateErr, "delete from bucket failed")
		}

		return nil
	})
}

// Return a bbolt key that should be associated with iri.
func (fs *FedEmbeddedStorage) toKey(iri *url.URL) string {
	var target url.URL

	target.Host = iri.Host
	target.Path = iri.Path

	return target.String()
}

func (fs *FedEmbeddedStorage) retrieve(bucket []byte, key string) ([]byte, error) {
	var bytes []byte

	err := fs.connection.View(func(tx *bbolt.Tx) error {
		var b *bbolt.Bucket

		if b = tx.Bucket(bucket); b == nil {
			return fmt.Errorf("cannot open bucket=%v", string(bucket))
		}

		if bytes = b.Get([]byte(key)); bytes == nil {
			return errors.NewfWith(http.StatusNotFound, "no entry for key=%v in bucket=%v", key, string(bucket))
		}

		return nil
	})

	return bytes, err
}

func (fs *FedEmbeddedStorage) store(bucket []byte, key string, value []byte) error {
	return fs.connection.Update(func(tx *bbolt.Tx) error {
		var b *bbolt.Bucket

		if b = tx.Bucket(bucket); b == nil {
			return fmt.Errorf("cannot open bucket=%v", string(bucket))
		}

		if err := b.Put([]byte(key), value); err != nil {
			return errors.Wrapf(err, "put key=%v into bucket=%v failed", key, string(bucket))
		}

		return nil
	})
}

func (fs *FedEmbeddedStorage) gcLoop() {
	for !fs.closed {
		if err := fs.gc(_CODES_BUCKET); err != nil {
			log.Printf("code garbage collection failed: %v", err)
		}

		if err := fs.gc(_TOKENS_BUCKET); err != nil {
			log.Printf("token garbage collection failed: %v", err)
		}

		time.Sleep(_GARBAGE_COLLECTION_WAIT)
	}
}

func (fs *FedEmbeddedStorage) gc(bucket []byte) (err error) {
	return fs.connection.Update(func(tx *bbolt.Tx) error {
		// open bucket for update

		var b *bbolt.Bucket

		if b = tx.Bucket(bucket); b == nil {
			return fmt.Errorf("cannot open bucket=%v", string(bucket))
		}

		// iterate over bucket; find all keys which contain expired codes
		// or tokens

		var expiredKeys [][]byte

		err = b.ForEach(func(key, value []byte) error {
			type expirer interface {
				Expired() bool
			}

			var (
				e  expirer
				oc FedOAuthCode
				ot FedOAuthToken
			)

			if err := json.Unmarshal(value, &oc); err == nil {
				e = &oc
				goto found
			}

			if err := json.Unmarshal(value, &ot); err == nil {
				e = &ot
				goto found
			}

		found:
			if e == nil {
				return errors.Newf("unexpected value=%v", string(value))
			}

			if e.Expired() {
				expiredKeys = append(expiredKeys, key)
			}

			return nil
		})

		if err != nil {
			return errors.Wrap(err, "error while trying to detrmine expired keys")
		}

		// now that we have all expired keys, we can delete those entries

		for _, key := range expiredKeys {
			if err := b.Delete(key); err != nil {
				return errors.Wrapf(err, "error deleting key=%v", string(key))
			}
		}

		return nil
	})
}
