package db

import (
	"errors"
	"github.com/go-fed/activity/streams/vocab"
	"net/url"
)

// Implements FedStorage, but never returns any errors or actual content.
// Can be used in tests where we do not want to retrieve anything from
// storage.
type FedEmptyStorage struct{}

func (f *FedEmptyStorage) Open() error {
	return nil
}

func (f *FedEmptyStorage) Close() error {
	return nil
}

func (f *FedEmptyStorage) RetrieveUser(username string) (*FedUser, error) {
	return nil, errors.New("not found (simulated)")
}

func (f *FedEmptyStorage) StoreUser(user *FedUser) error {
	return nil
}

func (f *FedEmptyStorage) RetrieveObject(iri *url.URL) (vocab.Type, error) {
	return nil, errors.New("not found (simulated)")
}

func (f *FedEmptyStorage) StoreObject(iri *url.URL, obj vocab.Type) error {
	return nil
}

func (f *FedEmptyStorage) DeleteObject(iri *url.URL) error {
	return nil
}
