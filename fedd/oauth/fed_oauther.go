package oauth

import (
	"crypto/rand"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/fedd/db"
	oaerrors "gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"log"
	"net/http"
)

type FedOAuther interface {
	HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request)
	HandleTokenRequest(w http.ResponseWriter, r *http.Request)
}

type fedOAuther struct {
	Manager     *manage.Manager
	ClientStore *store.ClientStore
	Server      *server.Server
}

func NewFedOAuther(storage db.FedStorage) (FedOAuther, error) {
	oa := &fedOAuther{}

	// create manager
	oa.Manager = manage.NewDefaultManager()

	// create token store
	if ts, err := store.NewMemoryTokenStore(); err != nil {
		return nil, errors.Wrap(err, "could not generate token store")
	} else {
		oa.Manager.MustTokenStorage(ts, nil)
	}

	// create client store
	oa.ClientStore = store.NewClientStore()
	oa.Manager.MapClientStorage(oa.ClientStore)

	// create server used in our http handlers
	oa.Server = server.NewDefaultServer(oa.Manager)
	oa.Server.SetAllowGetAccessRequest(true)
	oa.Server.SetClientInfoHandler(server.ClientFormHandler)

	// set error handlers
	oa.Server.SetInternalErrorHandler(func(err error) (re *oaerrors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})
	oa.Server.SetResponseErrorHandler(func(re *oaerrors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	return oa, nil
}

func (oa *fedOAuther) HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleAuthorizeRequest()")

	if err := oa.Server.HandleAuthorizeRequest(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (oa *fedOAuther) HandleTokenRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleTokenRequest()")

	oa.Server.HandleTokenRequest(w, r)
}

// Return a random string that contains 128 bytes of random data.
func (oa *fedOAuther) random() string {
	nbytes := 128
	b := make([]byte, nbytes)

	if _, err := rand.Read(b); err != nil {
		log.Panic("could not generate random data:", err)
	}

	return fmt.Sprintf("%x", b)
}
