package db

import (
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	"fmt"
	"go.etcd.io/bbolt"
	"net/url"
)

const _USERS_BUCKET = "Users"
const _DOCUMENTS_BUCKET = "Documents"

type FedEmbeddedStorage struct {
	Filepath   string
	connection *bbolt.DB
}

func (fs *FedEmbeddedStorage) Open() (err error) {
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
	return fs.connection.Close()
}

func (fs *FedEmbeddedStorage) RetrieveUser(username string) (user *FedUser, err error) {
	bucketKey := []byte(_USERS_BUCKET)
	userKey := []byte(username)

	err = fs.connection.View(func(tx *bbolt.Tx) error {
		var bucket *bbolt.Bucket
		var bytes []byte
		var viewErr error

		if bucket = tx.Bucket(bucketKey); bucket == nil {
			return errors.New("could not open users bucket")
		}

		if bytes = bucket.Get(userKey); bytes == nil {
			return fmt.Errorf("no user with username=%v", username)
		}

		user, viewErr = bytesToUser(bytes)
		return viewErr
	})

	return user, err
}

func (fs *FedEmbeddedStorage) StoreUser(user *FedUser) error {
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

func (fs *FedEmbeddedStorage) RetrieveObject(iri *url.URL) (vocab.Type, error) {
	return nil, errors.New("not implemented")
}

func (fs *FedEmbeddedStorage) StoreObject(obj vocab.Type) (*url.URL, error) {
	return nil, errors.New("not implemented")
}

func (fs *FedEmbeddedStorage) StoreObjectAt(iri *url.URL, obj vocab.Type) error {
	return errors.New("not implemented")
}
