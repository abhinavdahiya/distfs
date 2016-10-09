package main

import (
	"encoding/gob"
	"net"
	"os"

	"github.com/pkg/errors"
)

const (
	PORT = "5555"
	TYPE = "tcp"
)

type Client struct {
	Dest      string
	Connected bool
}

func (c *Client) Send() (bool, error) {
	defer func() {
		c.Connected = false
	}()

	p := os.Getenv("PORT")
	if p == "" {
		p = PORT
	}

	conn, err := net.Dial(TYPE, c.Dest+":"+p)
	if err != nil {
		return false, errors.Wrap(err, "Error Connecting to the other side"+c.Dest)
	}
	defer conn.Close()

	enc := gob.NewEncoder(conn)
	hb, err := Create()
	if err != nil {
		return false, err
	}

	enc.Encode(hb)
	return true, nil
}
