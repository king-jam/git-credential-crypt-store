package main

import (
	"fmt"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/boltdb"
	"github.com/king-jam/git-credential-crypt-store/crypto"
	"github.com/martinlindhe/inputbox"
)

func init() {
	boltdb.Register()
}

const (
	boltBaseString = "cryptstore"
	credBase       = "creds"
)

type CryptStoreInterface interface {
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

func (cs *CryptStore) LookupCredential(c *Credential) error {
	// creds := []*Credential{}
	// key := c.getCredBase()
	// pairs, err := c.db.List(key)
	// if err != nil {
	// 	return nil, err
	// }
	// for k, v := range pairs {
	//     cred
	// }
	return nil
}

func (cs *CryptStore) StoreCredential(c *Credential) error {
	// ensure we actually have stuff to store
	if !validToStore(c) {
		return fmt.Errorf("Could not store credentials")
	}
	// TODO: Make this better
	password, ok := inputbox.InputBox("git-credential-crypt-store", "Enter encryption passphrase", "")
	if !ok {
		return fmt.Errorf("Could not store credentials")
	}
	match, ok := inputbox.InputBox("git-credential-crypt-store", "Enter encryption passphrase again", "")
	if !ok {
		return fmt.Errorf("Could not store credentials")
	}
	if password != match {
		return fmt.Errorf("Could not store credentials")
	}
	cipher, err := crypto.NewCipher([]byte(password))
	if err != nil {
		return err
	}
	base := getCredBase()
	credKey := ""
	key := base + credKey
	if err != nil {
		return err
	}
	ciphertext, err := cipher.Encrypt([]byte(c.Password))
	if err != nil {
		return err
	}
	current, err := cs.db.Get(key)
	if err != nil {
		if err == store.ErrKeyNotFound {
			_, _, err := cs.db.AtomicPut(key, ciphertext, nil, nil)
			if err != nil {
				return err
			}
		}
		return err
	}
	_, _, err = cs.db.AtomicPut(key, ciphertext, current, nil)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CryptStore) EraseCredential(c *Credential) error {
	return nil
}

func validToStore(c *Credential) bool {
	// ensure we actually have stuff to store
	return c.Protocol == "" || !(c.Host == "" || c.Path == "") || c.Username == "" || c.Password == ""
}

// getCredBase returns the cred base path
func getCredBase() string {
	return credBase + "/"
}
