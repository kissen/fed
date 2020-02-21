package db

import (
	"fmt"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
	"log"
	"net/url"
)

const _USERS_BUCKET = "Users"
const _DOCUMENTS_BUCKET = "Documents"

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
	log.Println("RetrieveUser()")

	bucketKey := []byte(_USERS_BUCKET)
	userKey := []byte(username)

	err = fs.connection.View(func(tx *bbolt.Tx) error {
		var bucket *bbolt.Bucket
		var bytes []byte
		var viewErr error

		if bucket = tx.Bucket(bucketKey); bucket == nil {
			return errors.New("cannot open users bucket")
		}

		if bytes = bucket.Get(userKey); bytes == nil {
			return fmt.Errorf("no user with username=%v", username)
		}

		if user, viewErr = bytesToUser(bytes); viewErr != nil {
			return errors.Wrap(err, "deserializing user failed, database corrupted?")
		}

		return nil
	})

	return user, err
}

func (fs *FedEmbeddedStorage) StoreUser(user *FedUser) error {
	log.Println("StoreUser()")

	bucketKey := []byte(_USERS_BUCKET)
	userKey := []byte(user.Name)

	userValue, err := userToBytes(user)
	if err != nil {
		return errors.Wrap(err, "could not serialize user")
	}

	return fs.connection.Update(func(tx *bbolt.Tx) error {
		var bucket *bbolt.Bucket
		var updateErr error

		if bucket = tx.Bucket(bucketKey); bucket == nil {
			return errors.New("could not open users bucket")
		}

		if updateErr = bucket.Put(userKey, userValue); updateErr != nil {
			return errors.Wrap(updateErr, "put into bucket failed")
		}

		return nil
	})
}

func (fs *FedEmbeddedStorage) RetrieveObject(iri *url.URL) (obj vocab.Type, err error) {
	log.Println("RetrieveObject()")

	bucketKey := []byte(_DOCUMENTS_BUCKET)
	documentKey := []byte(normalizeIri(iri))

	err = fs.connection.View(func(tx *bbolt.Tx) error {
		var bucket *bbolt.Bucket
		var bytes []byte
		var viewErr error

		if bucket = tx.Bucket(bucketKey); bucket == nil {
			return errors.New("could not open documents bucket")
		}

		if bytes = bucket.Get(documentKey); bytes == nil {
			return fmt.Errorf("no entry for iri=%v", iri)
		}

		if obj, viewErr = bytesToVocab(bytes); viewErr != nil {
			errors.Wrap(err, "deserializing object failed, database corrupted?")
		}

		return nil
	})

	return obj, err
}

func (fs *FedEmbeddedStorage) StoreObject(iri *url.URL, obj vocab.Type) error {
	log.Println("StoreObject()")

	bucketKey := []byte(_DOCUMENTS_BUCKET)
	documentKey := []byte(normalizeIri(iri))

	documentValue, err := vocabToBytes(obj)
	if err != nil {
		return errors.Wrap(err, "could not serialize object")
	}

	return fs.connection.Update(func(tx *bbolt.Tx) error {
		var bucket *bbolt.Bucket
		var updateErr error

		if bucket = tx.Bucket(bucketKey); bucket == nil {
			return errors.New("could not docments bucket")
		}

		if updateErr = bucket.Put(documentKey, documentValue); updateErr != nil {
			return errors.Wrap(updateErr, "put into bucket failed")
		}

		return nil
	})
}
