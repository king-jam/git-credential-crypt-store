package main

import (
	"fmt"
	"os"

	"github.com/king-jam/git-credential-crypt-store/backend"
)

func lookupCredentials(backend backend.CryptStoreInterface, credentials *Credential) {
	fmt.Fprintln(os.Stdout, "username=king-jam")
	fmt.Fprintln(os.Stdout, "password=password")
	return
}
