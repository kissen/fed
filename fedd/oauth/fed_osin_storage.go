package oauth

import (
	"log"
	"github.com/RangelReale/osin"
	"gitlab.cs.fau.de/kissen/fed/fedd/db"
	"gitlab.cs.fau.de/kissen/fed/fedutil"
)

type FedOsinStorage struct {
	Storage db.FedStorage
}

func (os *FedOsinStorage) Clone() osin.Storage {
	log.Println("Clone")

	return &FedOsinStorage{
		Storage: os.Storage,
	}
}

func (os *FedOsinStorage) Close() {
	log.Println("Close")
}

// GetClient loads the client by id (client_id)
func (os *FedOsinStorage) GetClient(id string) (osin.Client, error) {
	log.Printf("GetClient(%v)", id)
}

// SaveAuthorize saves authorize data.
func (os *FedOsinStorage) SaveAuthorize(*osin.AuthorizeData) error {
	log.Printf("SaveAuthorize()")
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (os *FedOsinStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	log.Printf("LoadAuthorize(%v)", code)
}

// RemoveAuthorize revokes or deletes the authorization code.
func (os *FedOsinStorage) RemoveAuthorize(code string) error {
	log.Printf("RemoveAuthorize(%v)", code)
}

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (os *FedOsinStorage) SaveAccess(*osin.AccessData) error {
	log.Printf("SaveAccess")
}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (os *FedOsinStorage) LoadAccess(token string) (*osin.AccessData, error) {
	log.Printf("LoadAccess(%v)", token)
}

// RemoveAccess revokes or deletes an AccessData.
func (os *FedOsinStorage) RemoveAccess(token string) error {
	log.Printf("RemoveAccess(%v)", token)
}

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (os *FedOsinStorage) LoadRefresh(token string) (*osin.AccessData, error) {
	log.Printf("LoadRefresh(%v)", token)
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (os *FedOsinStorage) RemoveRefresh(token string) error {
	log.Printf("RemoveRefresh(%v)", token)
}
