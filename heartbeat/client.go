package heartbeat

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
	p := os.Getenv("PORT")
	if p == "" {
		p = PORT
	}

	conn, err := net.Dial(TYPE, c.Dest+":"+p)
	if err != nil {
		c.Connected = false
		return false, errors.Wrap(err, "Error Connecting to the other side"+c.Dest)
	}
	defer conn.Close()

	enc := gob.NewEncoder(conn)
	hb, err := Create()
	if err != nil {
		c.Connected = false
		return false, err
	}
	c.Connected = true
	enc.Encode(hb)
	return true, nil
}
