package main

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Credential struct {
	Username string
	Password string
	Protocol string
	Host     string
	Path     string
	URL      string
	Quit     int
}

func (c *Credential) ToURL() (*url.URL, error) {
	// if we already have the URL defined, just return it
	if c.URL != "" {
		url, err := url.Parse(c.URL)
		if err != nil {
			return nil, err
		}
		return url, nil
	}
	u := new(url.URL)
	// ensure host is defined
	if c.Host != "" {
		u.Host = c.Host
	}
	// ensure path is defined
	if c.Path != "" {
		u.Path = c.Path
	}
	// ensure protocol is defined
	if c.Protocol != "" {
		u.Scheme = c.Protocol
	}
	// if we have a username and a password
	if c.Username != "" && c.Password != "" {
		u.User = url.UserPassword(c.Username, c.Password)
	}
	// if we just have a username
	if c.Username != "" && c.Password == "" {
		u.User = url.User(c.Username)
	}
	return u, nil
}

func (c *Credential) IsValidToStore() bool {
	// ensure we actually have stuff to store
	return !(c.Protocol == "" || !(c.Host == "" || c.Path == "") || c.Username == "" || c.Password == "")
}

func (c *Credential) PrintToStdOut() {
	fmt.Fprintf(os.Stdout, "username=%s\n", c.Username)
	fmt.Fprintf(os.Stdout, "password=%s\n", c.Password)
}

func CredentialsMatch(want *Credential, have *Credential) bool {
	if want.Protocol != "" {
		if have.Protocol != "" {
			if want.Protocol != have.Protocol {
				return false
			}
		}
	}
	if want.Host != "" {
		if have.Host != "" {
			if want.Host != have.Host {
				return false
			}
		}
	}
	if want.Path != "" {
		if have.Path != "" {
			if want.Path != have.Path {
				return false
			}
		}
	}
	if want.Username != "" {
		if have.Username != "" {
			if want.Username != have.Username {
				return false
			}
		}
	}
	return true
}

func ParseCredentialStdin() (*Credential, error) {
	c := new(Credential)
	reader := bufio.NewReader(os.Stdin)
	for {
		// read a line from Stdin
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return c, nil
			}
			return nil, err
		}
		// if we got a partial line, we pull the rest until we get the full line
		for isPrefix {
			var lineFragment []byte
			lineFragment, isPrefix, err = reader.ReadLine()
			if err != nil {
				return nil, err
			}
			line = append(line, lineFragment...)
		}
		// if the line is empty, we are done getting data from the git credential service
		if len(line) == 0 {
			break
		}
		parts := strings.SplitN(string(line), "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("Invalid Input String")
		}
		key := parts[0]
		value := parts[1]
		switch key {
		case "username":
			c.Username = value
		case "password":
			c.Password = value
		case "protocol":
			c.Protocol = value
		case "host":
			c.Host = value
		case "path":
			c.Path = value
		case "url":
			if err := parseCredentialURL(value, c); err != nil {
				return nil, err
			}
		case "quit":
			if err := parseQuit(value, c); err != nil {
				return nil, err
			}
		default:
			// do nothing
		}
	}
	return nil, nil
}

func parseCredentialURL(rawurl string, c *Credential) error {
	url, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	c.URL = url.String()
	if url.Scheme != "" {
		c.Protocol = url.Scheme
	}
	/*
	       * Match one of:
	   *   (1) proto://<host>/...
	   *   (2) proto://<user>@<host>/...
	   *   (3) proto://<user>:<pass>@<host>/...
	*/
	if url.User != nil {
		username := url.User.Username()
		if username != "" {
			c.Username = url.User.Username()
		}
		password, passwordSet := url.User.Password()
		if passwordSet {
			c.Password = password
		}
	}
	c.Host = url.Host
	c.Path = url.Path
	return nil
}

func parseQuit(value string, c *Credential) error {
	quit, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	c.Quit = quit
	return nil
}
