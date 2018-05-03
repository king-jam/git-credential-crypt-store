package backend

import (
	"encoding/json"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/boltdb"
)

func init() {
	boltdb.Register()
}

const (
	boltBaseString = "cryptstore"
	credKey        = "creds"
)

type StorageContainer struct {
	CredentialURLs []string
	LastIndex      uint64
}

type CryptStoreInterface interface {
	GetStorageContainer() (*StorageContainer, error)
	AtomicPutStorageContainer(s *StorageContainer, prevIndex uint64) error
}

type CryptStore struct {
	db store.Store
}

func OpenCryptStore(storeLocation string) (*CryptStore, error) {
	kv, err := libkv.NewStore(
		store.BOLTDB,
		[]string{storeLocation},
		&store.Config{
			Bucket: boltBaseString,
		},
	)
	if err != nil {
		return nil, err
	}
	return &CryptStore{
		db: kv,
	}, nil
}

// getCredentialList ...
func (cs *CryptStore) GetStorageContainer() (*StorageContainer, error) {
	pair, err := cs.db.Get(credKey)
	if err != nil {
		return nil, err
	}
	var ret []string
	if err := json.Unmarshal(pair.Value, ret); err != nil {
		return nil, err
	}
	return &StorageContainer{
		CredentialURLs: ret,
		LastIndex:      pair.LastIndex,
	}, nil
}

// atomicStoreCredentialList ...
func (cs *CryptStore) AtomicPutStorageContainer(s *StorageContainer, prevIndex uint64) error {
	previous := &store.KVPair{
		LastIndex: prevIndex,
	}
	data, err := json.Marshal(s.CredentialURLs)
	if err != nil {
		return err
	}
	_, _, err = cs.db.AtomicPut(credKey, data, previous, nil)
	if err != nil {
		return err
	}
	return nil
}
