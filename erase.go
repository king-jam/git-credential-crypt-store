package main

import "github.com/king-jam/git-credential-crypt-store/backend"

func removeCredentials(db backend.CryptStoreInterface, credentials *Credential) error {
	s, err := db.GetStorageContainer()
	if err != nil {
		return err
	}
	// iterate to see if we already have these credentials stored
	for idx, elem := range s.CredentialURLs {
		c := new(Credential)
		err := parseCredentialURL(elem, c)
		if err != nil {
			return err
		}

		if CredentialsMatch(credentials, c) {
			// delete the current index
			copy(s.CredentialURLs[idx:], s.CredentialURLs[idx+1:])
			s.CredentialURLs[len(s.CredentialURLs)-1] = ""
			s.CredentialURLs = s.CredentialURLs[:len(s.CredentialURLs)-1]
			// persist the updated copy
			err = db.PersistStorageContainer(s)
			if err != nil {
				return err
			}

			break
		}
	}

	return nil
}
