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
type EmptyTx struct{}

func (f FedEmptyStorage) Open() error {
	return nil
}

func (f FedEmptyStorage) Close() error {
	return nil
}

func (f FedEmptyStorage) Begin() (Tx, error) {
	return EmptyTx{}, nil
}

func (f EmptyTx) Commit() error {
	return nil
}

func (f EmptyTx) Rollback() error {
	return nil
}

func (f EmptyTx) RetrieveUser(username string) (*FedUser, error) {
	return nil, errors.New("not found (simulated)")
}

func (f EmptyTx) StoreUser(user *FedUser) error {
	return nil
}

func (f EmptyTx) RetrieveCode(code string) (*FedOAuthCode, error) {
	return nil, nil
}

func (f EmptyTx) StoreCode(code *FedOAuthCode) error {
	return nil
}

func (f EmptyTx) RetrieveToken(token string) (*FedOAuthToken, error) {
	return nil, nil
}

func (f EmptyTx) StoreToken(token *FedOAuthToken) error {
	return nil
}

func (f EmptyTx) RetrieveObject(iri *url.URL) (vocab.Type, error) {
	return nil, errors.New("not found (simulated)")
}

func (f EmptyTx) StoreObject(iri *url.URL, obj vocab.Type) error {
	return nil
}

func (f EmptyTx) DeleteObject(iri *url.URL) error {
	return nil
}
