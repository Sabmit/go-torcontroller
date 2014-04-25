package torcontroller

import (
	"net/textproto"
//	"fmt"
//	"io/ioutils"
)


type Client struct {
	c *textproto.Conn
	addr string
	isConnected bool
}

func (this *Client) Close() error {
	return this.c.Close()
}

func (this *Client) setConnected(connected bool) *Client {
	this.isConnected = connected
	return this
}

func (this *Client) Authenticate(pass string) error {
	var id uint
	var err error

	if (pass == "") {
		id, err = this.c.Cmd("AUTHENTICATE")
	} else {
		id, err = this.c.Cmd("AUTHENTICATE %s", pass)
	}

	if err != nil {
		return err
	}

	this.c.StartResponse(id)
	defer this.c.EndResponse(id)

	code, _, err := this.c.ReadCodeLine(250)

	if err != nil  {
		return err
	} else if code != 250 {
		return nil //new error; isConnected = false
	}
	return nil
}

func (this *Client) send(format string, args ...interface{}) (uint, error) {
	return this.c.Cmd(format, args)
}

func NewClient(addr string) (*Client, error) {
	c, err := textproto.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	client := &Client{
		c: c,
		addr: addr}

	client.setConnected(true)
	return client, nil
}
