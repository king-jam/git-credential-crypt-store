package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/king-jam/git-credential-crypt-store/backend"
	homedir "github.com/mitchellh/go-homedir"
)

const storeFileName = ".git-credential-crypt-store"

const storeLocationDefault = "~/" + storeFileName

func main() {
	var storeLocation string

	flag.StringVar(&storeLocation, "file", storeLocationDefault, "Location to store the credentials.")
	// define a quick helper function for usage so we can let people know
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage:\n")
		fmt.Fprint(os.Stderr, "  git-credential-crypt-store [OPTIONS] [CMD]\n\n")

		title := "git credential helper to store passwords encrypted to enable usage of access tokens with 2FA."
		fmt.Fprint(os.Stderr, title+"\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	// parse the flags, we will use the default if nothing is configured
	flag.Parse()

	if storeLocation == storeLocationDefault {
		home, err := homedir.Dir()
		if err != nil {
			os.Exit(1)
		}
		storeLocation = home + "/.git-credential-crypt-store"
	}
	// if we don't get anything after the program, just give an error back
	if len(os.Args[1:]) == 0 {
		flag.Usage()
	}
	// open up the credential storage
	cs, err := backend.OpenCryptStore(storeLocation)
	if err != nil {
		os.Exit(1)
	}
	// if we got here then the input arguments are at least correct
	// parse in the credentials
	creds, err := ParseCredentialStdin()
	if err != nil {
		os.Exit(1)
	}
	// just grab the last argument at this point, if it isn't a command, just ignore it
	switch op := os.Args[len(os.Args)-1]; op {
	case "get":
		if err := lookupCredentials(cs, creds); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "store":
		if err := storeCredentials(cs, creds); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "erase":
		if err := removeCredentials(cs, creds); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		// ignore unknown operation
	}

	os.Exit(0)
}
