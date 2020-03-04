package db

import (
	"fmt"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/fedutil"
	"go.etcd.io/bbolt"
	"log"
	"net/url"
)

var _USERS_BUCKET = []byte("Users")
var _DOCUMENTS_BUCKET = []byte("Documents")

type FedEmbeddedStorage struct {
	Filepath   string
	connection *bbolt.DB
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
		if _, err := tx.CreateBucketIfNotExists([]byte(_USERS_BUCKET)); err != nil {
			return err
		}

		if _, err := tx.CreateBucketIfNotExists([]byte(_DOCUMENTS_BUCKET)); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "intializing buckets failed")
	}

	// success

	return nil
}

func (fs *FedEmbeddedStorage) Close() error {
	log.Println("Close()")

	return fs.connection.Close()
}

func (fs *FedEmbeddedStorage) RetrieveUser(username string) (user *FedUser, err error) {
	log.Printf("RetrieveUser(%s)", username)

	userKey := []byte(username)

	var bytes []byte

	err = fs.connection.View(func(tx *bbolt.Tx) error {
		var bucket *bbolt.Bucket

		if bucket = tx.Bucket(_USERS_BUCKET); bucket == nil {
			return errors.New("cannot open users bucket")
		}

		if bytes = bucket.Get(userKey); bytes == nil {
			return fmt.Errorf("no user with username=%v", username)
		}

		return nil
	})

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

	userKey := []byte(user.Name)
	userValue, err := userToBytes(user)
	if err != nil {
		return errors.Wrap(err, "could not serialize user")
	}

	return fs.connection.Update(func(tx *bbolt.Tx) error {
		var bucket *bbolt.Bucket
		var updateErr error

		if bucket = tx.Bucket(_USERS_BUCKET); bucket == nil {
			return errors.New("could not open users bucket")
		}

		if updateErr = bucket.Put(userKey, userValue); updateErr != nil {
			return errors.Wrap(updateErr, "put into bucket failed")
		}

		return nil
	})
}

func (fs *FedEmbeddedStorage) RetrieveObject(iri *url.URL) (obj vocab.Type, err error) {
	log.Printf("RetrieveObject(%v)", iri)

	var bytes []byte

	err = fs.connection.View(func(tx *bbolt.Tx) error {
		var bucket *bbolt.Bucket

		if bucket = tx.Bucket(_DOCUMENTS_BUCKET); bucket == nil {
			return errors.New("could not open documents bucket")
		}

		if bytes = bucket.Get(fs.toKey(iri)); bytes == nil {
			return fmt.Errorf("no entry for iri=%v", iri)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if obj, err = fedutil.BytesToVocab(bytes); err != nil {
		return nil, errors.Wrap(err, "deserializing object failed")
	}

	return obj, err
}

func (fs *FedEmbeddedStorage) StoreObject(iri *url.URL, obj vocab.Type) error {
	log.Printf("StoreObject(%v)", iri)

	documentValue, err := fedutil.VocabToBytes(obj)
	if err != nil {
		return errors.Wrap(err, "could not serialize object")
	}

	return fs.connection.Update(func(tx *bbolt.Tx) error {
		var bucket *bbolt.Bucket
		var updateErr error

		if bucket = tx.Bucket(_DOCUMENTS_BUCKET); bucket == nil {
			return errors.New("could not open documents bucket")
		}

		if updateErr = bucket.Put(fs.toKey(iri), documentValue); updateErr != nil {
			return errors.Wrap(updateErr, "put into bucket failed")
		}

		return nil
	})
}

func (fs *FedEmbeddedStorage) DeleteObject(iri *url.URL) error {
	log.Printf("DeleteObject(%v)", iri)

	return fs.connection.Update(func(tx *bbolt.Tx) error {
		var bucket *bbolt.Bucket
		var updateErr error

		if bucket = tx.Bucket(_DOCUMENTS_BUCKET); bucket == nil {
			return errors.New("could not open documents bucket")
		}

		if updateErr = bucket.Delete(fs.toKey(iri)); updateErr != nil {
			return errors.Wrap(updateErr, "delete from bucket failed")
		}

		return nil
	})
}

// Return a bbolt key that should be associated with iri.
func (fs *FedEmbeddedStorage) toKey(iri *url.URL) []byte {
	var target url.URL

	target.Host = iri.Host
	target.Path = iri.Path

	return []byte(target.String())
}
