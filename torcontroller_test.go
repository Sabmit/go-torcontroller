package torcontroller

import (
	//	"fmt"
	"testing"
	//	"io/ioutils"
)

const (
	addr = "127.0.0.1:9051"
	// plain password MUST BE put in double quotes
	password = "\"DrAwatRubexA3e=\""
	// or, you can encode it in hexadicimal
	passwordHex = "44724177617452756265784133653d"
)

func setupClient(t *testing.T) *Client {
	c, err := NewClient(addr)

	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestAuthenticate(t *testing.T) {
	c := setupClient(t)
	defer c.Close()

	err := c.Authenticate("A bad password")
	if err == nil {
		t.Fatal("Error, expected error code 515, got", 250)
	}

	c.ReConnect() // Bad authentificate causes the server to close de socket

	err = c.Authenticate(passwordHex)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Authenticate(password)
	if err != nil {
		t.Fatal(err)
	}

}

func TestNewClient(t *testing.T) {
	c, err := NewClient(addr)
	cSec, errSec := NewClient("4242:4242")

	if err != nil {
		t.Fatal(err)
	} else if errSec == nil {
		cSec.Close()
		t.Fatal("Error must be, but is not returned by NewClient")
	}
	c.Close()
}
