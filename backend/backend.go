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
	PersistStorageContainer(s *StorageContainer) error
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

// GetStorageContainer ...
func (cs *CryptStore) GetStorageContainer() (*StorageContainer, error) {
	s := new(StorageContainer)
	pair, err := cs.db.Get(credKey)
	if err != nil {
		if err == store.ErrKeyNotFound {
			// initialize an empty object
			s.CredentialURLs = make([]string, 0)
			s.LastIndex = 0
			return s, nil
		}
		return nil, err
	}
	var ret []string
	if err := json.Unmarshal(pair.Value, &ret); err != nil {
		return nil, err
	}
	s.CredentialURLs = ret
	s.LastIndex = pair.LastIndex
	return s, nil
}

// PersistStorageContainer ...
func (cs *CryptStore) PersistStorageContainer(s *StorageContainer) error {
	var previous *store.KVPair
	if s.LastIndex == 0 {
		previous = nil
	} else {
		previous = &store.KVPair{
			LastIndex: s.LastIndex,
		}
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
