package torcontroller

import (
	"fmt"
	"net/textproto"
	//	"io/ioutils"
)

type ErrInvalidPassword string

type Client struct {
	c           *textproto.Conn
	addr        string
	isConnected bool
}

func (self ErrInvalidPassword) Error() string {
	return fmt.Sprintf("Invalid password %q. Password MUST BE encoded in hexadecimal OR surrounded by double quotes", self)
}

func (this *Client) Close() error {
	return this.c.Close()
}

func (this *Client) setConnected(connected bool) error {
	this.isConnected = connected
	return nil
}

func (this *Client) readResponse(id uint, expectedCode int) error {
	this.c.StartResponse(id)
	defer this.c.EndResponse(id)

	_, _, err := this.c.ReadCodeLine(expectedCode)

	if err != nil {
		return err
	}
	return nil
}

func (this *Client) Authenticate(pass string) error {
	var id uint
	var err error

	if len(pass) == 0 {
		id, err = this.c.Cmd("AUTHENTICATE")
	} else {
		id, err = this.c.Cmd("AUTHENTICATE %s", pass)
	}

	if err != nil {
		return err
	}

	if err = this.readResponse(id, 250); err != nil {
		this.setConnected(false)
		this.Close()
		return err
	}
	return nil
}

func (this *Client) send(format string, args ...interface{}) (uint, error) {
	return this.c.Cmd(format, args)
}

func (this *Client) ReConnect() (*Client, error) {
	if this.isConnected == false {
		c, err := textproto.Dial("tcp", this.addr)
		if err != nil {
			return nil, err
		}
		this.c = c
		this.setConnected(true)
	}
	return this, nil
}

func NewClient(addr string) (*Client, error) {
	c, err := textproto.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	client := &Client{
		c:    c,
		addr: addr}

	client.setConnected(true)
	return client, nil
}
