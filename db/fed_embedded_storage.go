package db

import (
	"encoding/json"
	"fmt"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/kissen/fed/errors"
	"github.com/kissen/fed/marshal"
	"github.com/kissen/fed/util"
	"go.etcd.io/bbolt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const _GARBAGE_COLLECTION_WAIT = 1 * time.Minute
const _READ_ONLY = false
const _READ_WRITE = true

var _USERS_BUCKET = []byte("Users")
var _CODES_BUCKET = []byte("OAuth/Codes")
var _TOKENS_BUCKET = []byte("OAuth/Tokens")
var _DOCUMENTS_BUCKET = []byte("Documents")

type FedEmbeddedStorage struct {
	Filepath   string
	connection *bbolt.DB
	closed     bool
	rwlock     util.RWLock
}

type fedembeddedtx struct {
	// Whoever created this tx
	parent *FedEmbeddedStorage

	// btx is the underlying bbolt transaction; when haveWriteLock
	// is false, it is a read-only transaction, when haveWriteLock
	// is true, btx is an rw transaction; do not use btx directly,
	// instead call view or update
	btx           *bbolt.Tx
	haveWriteLock bool

	// Whether Commit or Rollback has been called before.
	commited bool

	// The error returned by the first call to Commit or Rollback.
	commitedError error
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

	fs.rwlock.Lock()
	defer fs.rwlock.Unlock()

	fs.closed = true
	return fs.connection.Close()
}

func (fs *FedEmbeddedStorage) Begin() (Tx, error) {
	log.Println("Begin()")

	// commited/rollbacked by tx (below) in Commit or Rollback
	btx, err := fs.connection.Begin(_READ_ONLY)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create transaction")
	}

	// released by tx (below) in Commit or Rollback
	fs.rwlock.RLock()

	tx := &fedembeddedtx{
		parent:        fs,
		btx:           btx,
		haveWriteLock: false,
	}

	return tx, nil
}

func (fs *FedEmbeddedStorage) RetrieveUser(username string) (user *FedUser, err error) {
	if tx, err := fs.Begin(); err != nil {
		return nil, err
	} else if user, err := tx.RetrieveUser(username); err != nil {
		tx.Commit()
		return nil, err
	} else {
		return user, tx.Commit()
	}
}

func (fs *FedEmbeddedStorage) StoreUser(user *FedUser) error {
	if tx, err := fs.Begin(); err != nil {
		return err
	} else if err := tx.StoreUser(user); err != nil {
		tx.Commit()
		return err
	} else {
		return tx.Commit()
	}
}

func (fs *FedEmbeddedStorage) RetrieveCode(code string) (*FedOAuthCode, error) {
	if tx, err := fs.Begin(); err != nil {
		return nil, err
	} else if oc, err := tx.RetrieveCode(code); err != nil {
		tx.Commit()
		return nil, err
	} else {
		return oc, tx.Commit()
	}
}

func (fs *FedEmbeddedStorage) StoreCode(code *FedOAuthCode) error {
	if tx, err := fs.Begin(); err != nil {
		return err
	} else if err := tx.StoreCode(code); err != nil {
		tx.Commit()
		return err
	} else {
		return tx.Commit()
	}
}

func (fs *FedEmbeddedStorage) RetrieveToken(token string) (*FedOAuthToken, error) {
	if tx, err := fs.Begin(); err != nil {
		return nil, err
	} else if ot, err := tx.RetrieveToken(token); err != nil {
		tx.Commit()
		return nil, err
	} else {
		return ot, tx.Commit()
	}
}

func (fs *FedEmbeddedStorage) StoreToken(token *FedOAuthToken) error {
	if tx, err := fs.Begin(); err != nil {
		return err
	} else if err := tx.StoreToken(token); err != nil {
		tx.Commit()
		return err
	} else {
		return tx.Commit()
	}
}

func (fs *FedEmbeddedStorage) RetrieveObject(iri *url.URL) (obj vocab.Type, err error) {
	if tx, err := fs.Begin(); err != nil {
		return nil, err
	} else if obj, err := tx.RetrieveObject(iri); err != nil {
		tx.Commit()
		return nil, err
	} else {
		return obj, tx.Commit()
	}
}

func (fs *FedEmbeddedStorage) StoreObject(iri *url.URL, obj vocab.Type) error {
	if tx, err := fs.Begin(); err != nil {
		return err
	} else if err := tx.StoreObject(iri, obj); err != nil {
		tx.Commit()
		return err
	} else {
		return tx.Commit()
	}
}

func (fs *FedEmbeddedStorage) DeleteObject(iri *url.URL) error {
	if tx, err := fs.Begin(); err != nil {
		return err
	} else if err := tx.DeleteObject(iri); err != nil {
		tx.Commit()
		return err
	} else {
		return tx.Commit()
	}
}

// Keep garbage collecting the database.
func (fs *FedEmbeddedStorage) gcLoop() {
	for !fs.closed {
		if err := fs.gc(); err != nil {
			log.Println("garbage collection failed:", err)
		}

		time.Sleep(_GARBAGE_COLLECTION_WAIT)
	}
}

// Garbage collect everything there is to garbage collect. This is one
// long running operation that blocks everything else. Don't run it too often.
func (fs *FedEmbeddedStorage) gc() (err error) {
	tx, err := fs.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Commit(); err != nil {
			log.Println("garbage collection commit failed:", err)
		}
	}()

	etx := tx.(*fedembeddedtx)

	if err := fs.gcBucket(etx, _CODES_BUCKET); err != nil {
		err = errors.Wrap(err, "code garbage collection failed")
	}

	if err := fs.gcBucket(etx, _TOKENS_BUCKET); err != nil {
		err = errors.Wrap(err, "token garbage collection failed")
	}

	return err
}

// Garbage collection bucket. Only call this method if you are holding
// the write lock.
func (fs *FedEmbeddedStorage) gcBucket(tx *fedembeddedtx, bucket []byte) (err error) {
	return tx.update(func(tx *bbolt.Tx) error {
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

func (fs *fedembeddedtx) Commit() (err error) {
	log.Println("Commit()")

	if !fs.commited {
		if fs.haveWriteLock {
			fs.parent.rwlock.RUnlock()
			fs.parent.rwlock.Unlock()
			err = fs.btx.Commit()
		} else {
			fs.parent.rwlock.RUnlock()
			err = fs.btx.Rollback()
		}

		fs.commited = true

		if err != nil {
			fs.commitedError = errors.Wrap(err, "previous Commit failed")
		}

		return err
	}

	return fs.commitedError
}

func (fs *fedembeddedtx) Rollback() (err error) {
	log.Println("Rollback()")

	if !fs.commited {
		if fs.haveWriteLock {
			fs.parent.rwlock.RUnlock()
			fs.parent.rwlock.Unlock()
			err = fs.btx.Rollback()
		} else {
			fs.parent.rwlock.RUnlock()
			err = fs.btx.Rollback()
		}

		fs.commited = true

		if err != nil {
			fs.commitedError = errors.Wrap(err, "previous Rollback failed")
		}

		return err
	}

	return fs.commitedError
}

func (fs *fedembeddedtx) RetrieveUser(username string) (user *FedUser, err error) {
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

func (fs *fedembeddedtx) StoreUser(user *FedUser) error {
	log.Printf("StoreUser(Name=%v #Inbox=%v #Outbox=%v)", user.Name, len(user.Inbox), len(user.Outbox))

	bytes, err := userToBytes(user)
	if err != nil {
		return errors.Wrap(err, "could not serialize user")
	}

	return fs.store(_USERS_BUCKET, user.Name, bytes)
}

func (fs *fedembeddedtx) RetrieveCode(code string) (*FedOAuthCode, error) {
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

func (fs *fedembeddedtx) StoreCode(code *FedOAuthCode) error {
	log.Printf("StoreCode(Code=%v)", code.Code)

	bs, err := json.Marshal(code)
	if err != nil {
		return errors.Wrap(err, "serializing code failed")
	}

	return fs.store(_CODES_BUCKET, code.Code, bs)
}

func (fs *fedembeddedtx) RetrieveToken(token string) (*FedOAuthToken, error) {
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

func (fs *fedembeddedtx) StoreToken(token *FedOAuthToken) error {
	log.Printf("StoreToken(Token=%v)", token.Token)

	bs, err := json.Marshal(token)
	if err != nil {
		return errors.Wrap(err, "serializing token failed")
	}

	return fs.store(_TOKENS_BUCKET, token.Token, bs)
}

func (fs *fedembeddedtx) RetrieveObject(iri *url.URL) (obj vocab.Type, err error) {
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

func (fs *fedembeddedtx) StoreObject(iri *url.URL, obj vocab.Type) error {
	log.Printf("StoreObject(%v)", iri)

	bytes, err := marshal.VocabToBytes(obj)
	if err != nil {
		return errors.Wrap(err, "could not serialize object")
	}

	return fs.store(_DOCUMENTS_BUCKET, fs.toKey(iri), bytes)
}

func (fs *fedembeddedtx) DeleteObject(iri *url.URL) error {
	log.Printf("DeleteObject(%v)", iri)

	return fs.update(func(tx *bbolt.Tx) error {
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
func (fs *fedembeddedtx) toKey(iri *url.URL) string {
	var target url.URL

	target.Host = iri.Host
	target.Path = iri.Path

	return target.String()
}

// Retreive bytes from bucket.
func (fs *fedembeddedtx) retrieve(bucket []byte, key string) ([]byte, error) {
	var bytes []byte

	err := fs.view(func(tx *bbolt.Tx) error {
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

func (fs *fedembeddedtx) store(bucket []byte, key string, value []byte) error {
	return fs.update(func(tx *bbolt.Tx) error {
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

func (fs *fedembeddedtx) view(operation func(tx *bbolt.Tx) error) error {
	return operation(fs.btx)
}

func (fs *fedembeddedtx) update(operation func(tx *bbolt.Tx) error) (err error) {
	if !fs.haveWriteLock {
		fs.parent.rwlock.Lock()
		fs.haveWriteLock = true

		if err := fs.btx.Rollback(); err != nil {
			return errors.Wrap(err, "cannot close read transaction")
		}

		fs.btx, err = fs.parent.connection.Begin(_READ_WRITE)
		if err != nil {
			return errors.Wrap(err, "cannot create write transaction")
		}
	}

	return operation(fs.btx)
}
