package file

import (
	"io"
	"net"
	"os"

	"github.com/pkg/errors"
)

//the destination of the file, and the location of the file to send.
type Client struct {
	Dest      string
	Location  string
	Connected bool
}

// Send transfers a file from the local machine to the given destination.
func (c *Client) Send() (bool, error) {
	file, err := os.Open(c.Location)
	if err != nil {
		return false, errors.Wrap(err, "Error opening file")
	}
	defer file.Close()

	fileSize, fileName, err := getFileStats(file)
	if err != nil {
		return false, errors.Wrap(err, "Error getting file info")
	}

	p := os.Getenv("PORT")
	if p == "" {
		p = PORT
	}
	conn, err := net.Dial(Type, c.Dest+":"+p)
	if err != nil {
		return false, errors.Wrap(err, "Error connecting to destination")
	}

	c.Connected = true

	defer conn.Close()

	conn.Write([]byte(fileSize))
	conn.Write([]byte(fileName))

	return c.StreamFile(conn, file)
}

// StreamFile breaks a file down into smaller pieces and sends it to the server
// via TCP.
func (c *Client) StreamFile(conn net.Conn, file *os.File) (bool, error) {
	sendBuffer := make([]byte, BufferSize)
	for {
		_, err := file.Read(sendBuffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return false, err
		}
		conn.Write(sendBuffer)
	}

	c.Connected = false
	return true, nil
}
