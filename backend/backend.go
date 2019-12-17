// Package backend handles overall persistence
package backend

import (
	"encoding/json"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/boltdb"
)

const (
	boltBaseString = "cryptstore"
	credKey        = "creds"
)

// StorageContainer is the top-level struct for persistence
type StorageContainer struct {
	CredentialURLs []string
	LastIndex      uint64
}

// CryptStoreInterface defines the persistence interface exposed to other packages
type CryptStoreInterface interface {
	GetStorageContainer() (*StorageContainer, error)
	PersistStorageContainer(s *StorageContainer) error
}

// CryptStore internally holds all persistence state dependencies
type CryptStore struct {
	db store.Store
}

// OpenCryptStore initializes and opens the persistence file
func OpenCryptStore(storeLocation string) (*CryptStore, error) {
	boltdb.Register()

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

// GetStorageContainer returns the entire storage container
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

// PersistStorageContainer persists the entire storage container
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
