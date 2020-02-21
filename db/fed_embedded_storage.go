package db

type FedEmbeddedStorage struct {
}

func (fs *FedEmbeddedStorage) RetrieveUser(username string) (*FedUser, error) {
}

func (fs *FedEmbeddedStorage) StoreUser(username string) (*FedUser, error) {
}

func (fs *FedEmbeddedStorage) RetrieveObject(iri *url.URL) (vocab.Type, error) {
}

func (fs *FedEmbeddedStorage) StoreObject(obj vocab.Type) (*url.URL, error) {
}

func (fs *FedEmbeddedStorage) StoreObjectAt(iri *url.URL, obj vocab.Type) error {
}
