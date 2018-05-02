package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var storeLocation string

func init() {
	flag.StringVar(&storeLocation, "file", "$HOME/.git-credential-crypt-store", "Location to store the credentials.")
}

func main() {
	// define a quick helper function for usage so we can let people know
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage:\n")
		fmt.Fprint(os.Stderr, "  git-credential-crypt-store [OPTIONS] [CMD]\n\n")

		title := "git credential helper to store passwords encrypted to enable usage of access tokens with 2FA."
		fmt.Fprint(os.Stderr, title+"\n\n")
		flag.PrintDefaults()
		return
	}
	// parse the flags, we will use the default if nothing is configured
	flag.Parse()
	// if we don't get anything after the program, just give an error back
	if len(os.Args[1:]) == 0 {
		flag.Usage()
	}
	// if we got here then the input arguments are at least correct
	// open up the credential storage
	creds, err := parseCredentialStdin()
	log.Printf("%+v", creds)
	if err != nil {
		os.Exit(1)
	}
	// just grab the last argument at this point, if it isn't a command, just ignore it
	switch op := os.Args[len(os.Args)-1]; op {
	case "get":
		log.Print("GET")
		lookupCredentials(storeLocation, creds)
	case "store":
		log.Print("STORE")
		storeCredentials(storeLocation, creds)
	case "erase":
		log.Print("ERASE")
		removeCredentials(storeLocation, creds)
	default:
		// ignore unknown operation
	}
}

func lookupCredentials(storeLocation string, credentials *Credential) {
	fmt.Fprintln(os.Stdout, "username=king-jam")
	fmt.Fprintln(os.Stdout, "password=password")
	return
}

func storeCredentials(storeLocation string, credentials *Credential) {
	return
}

func removeCredentials(storeLocation string, credentials *Credential) {
	return
}
