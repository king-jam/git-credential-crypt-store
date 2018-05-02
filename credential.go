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

func parseCredentialStdin() (*Credential, error) {
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
