package main

import (
	"github.com/king-jam/git-credential-crypt-store/backend"
	"github.com/king-jam/git-credential-crypt-store/crypto"
	"github.com/king-jam/git-credential-crypt-store/dialogs"
)

func lookupCredentials(db backend.CryptStoreInterface, credentials *Credential) error {
	var s *backend.StorageContainer
	s, err := db.GetStorageContainer()
	if err != nil {
		return err
	}
	// iterate to see if we already have these credentials stored
	for _, elem := range s.CredentialURLs {
		c := new(Credential)
		err := parseCredentialURL(elem, c)
		if err != nil {
			return err
		}

		if CredentialsMatch(credentials, c) {
			password, err := dialogs.PasswordBox(c.Username)
			if err != nil {
				return err
			}
			cipher, err := crypto.NewCipher(password)
			if err != nil {
				return err
			}
			decryptedPassword, err := cipher.Decrypt([]byte(c.Password))
			if err != nil {
				return err
			}
			c.Password = string(decryptedPassword)
			c.PrintToStdOut()
			return nil
		}
	}
	return nil
}
