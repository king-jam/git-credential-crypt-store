package main

import (
	"log"

	"github.com/king-jam/git-credential-crypt-store/backend"
	"github.com/king-jam/git-credential-crypt-store/crypto"
	"github.com/king-jam/git-credential-crypt-store/dialogs"
)

func storeCredentials(db backend.CryptStoreInterface, credentials *Credential) error {
	var s *backend.StorageContainer
	password, err := dialogs.PasswordCreationBox()
	if err != nil {
		return err
	}
	cipher, err := crypto.NewCipher(password)
	if err != nil {
		return err
	}
	s, err = db.GetStorageContainer()
	if err != nil {
		if err == backend.ErrKeyNotFound {
			// initialize our Credentials
			s = &backend.StorageContainer{
				CredentialURLs: make([]string, 0),
				LastIndex:      0,
			}
		} else {
			return err
		}
	}
	// declare index to -1 for a later check
	matchIndex := -1
	// iterate to see if we already have these credentials stored
	for idx, elem := range s.CredentialURLs {
		c := new(Credential)
		err := parseCredentialURL(elem, c)
		if err != nil {
			return err
		}
		if CredentialsMatch(credentials, c) {
			matchIndex = idx
			break
		}
	}
	// if we already have an entry for this, overwrite it
	if matchIndex != -1 {
		// delete the current index
		copy(s.CredentialURLs[matchIndex:], s.CredentialURLs[matchIndex+1:])
		s.CredentialURLs[len(s.CredentialURLs)-1] = ""
		s.CredentialURLs = s.CredentialURLs[:len(s.CredentialURLs)-1]
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
	err = db.AtomicPutStorageContainer(s)
	if err != nil {
		return err
	}
	log.Printf("%+v", s)
	for _, t := range s.CredentialURLs {
		log.Printf("%s\n", t)
	}
	return nil
}
