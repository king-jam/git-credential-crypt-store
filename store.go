package main

import (
	"fmt"

	"github.com/king-jam/git-credential-crypt-store/backend"
	"github.com/king-jam/git-credential-crypt-store/crypto"
	"github.com/king-jam/git-credential-crypt-store/dialogs"
)

func storeCredentials(db backend.CryptStoreInterface, credentials *Credential) error {
	// check that they are valid to store
	if !credentials.IsValidToStore() {
		return fmt.Errorf("Invalid Credential Storage Format")
	}
	s, err := db.GetStorageContainer()
	if err != nil {
		return err
	}
	// declare index to -1 for a later check
	// iterate to see if we already have these credentials stored
	for _, elem := range s.CredentialURLs {
		c := new(Credential)
		err := parseCredentialURL(elem, c)
		if err != nil {
			return err
		}
		if CredentialsMatch(credentials, c) {
			return nil
		}
	}
	// if we don't have an entry, create it
	password, err := dialogs.PasswordCreationBox(credentials.Username)
	if err != nil {
		return err
	}
	cipher, err := crypto.NewCipher(password)
	if err != nil {
		return err
	}
	ciphertext, err := cipher.Encrypt([]byte(credentials.Password))
	if err != nil {
		return err
	}
	credentials.Password = string(ciphertext)
	credAsURL, err := credentials.ToURL()
	if err != nil {
		return err
	}

	s.CredentialURLs = append(s.CredentialURLs, credAsURL.String())
	err = db.PersistStorageContainer(s)
	if err != nil {
		return err
	}
	return nil
}
