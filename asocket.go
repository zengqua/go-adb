package adb

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
)

/*
Quote from https://android.googlesource.com/platform/packages/modules/adb

“SMART SOCKET”:
  A smart socket is a simple TCP socket with a smart protocol built on top of it.
This is what Clients connect onto from the Host side. The Client must always
initiate communication via a human readable request but the response format varies.
The smart protocol is documented in SERVICES.TXT. see
https://android.googlesource.com/platform/packages/modules/adb/+/HEAD/SERVICES.TXT
*/

// Asocket smart socket
type Asocket struct {
	conn net.Conn
}

// Dial address format "host:port"
func Dial(address string) (*Asocket, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Asocket{conn: conn}, nil
}

// IO

// SendMessage send string
func (c *Conn) SendMessage(msg string) error {
	return c.Send([]byte(msg))
}

// SendAndCheck check status after send message
func (c *Conn) SendAndCheck(msg string) error {
	if err := c.SendMessage(msg); err != nil {
		return err
	}
	if _, err := c.ReadStatus(); err != nil {
		return err
	}
	return nil
}

// Read bytes
func (c *Conn) Read(n int) ([]byte, error) {
	b := make([]byte, n)
	offset, err := io.ReadFull(c.conn, b)
	if err != nil {
		return nil, fmt.Errorf("incomplete message, expect: %d but get: %d.\n%w", n, offset, err)
	}
	return b, nil
}

// ReadMessage n string
func (c *Conn) ReadMessage(n int) (string, error) {
	b, err := c.Read(n)
	if err != nil {
		return "", err
	}
	return string(bytes.Trim(b, "\x00")), nil
}

// ReadAll string
func (c *Conn) ReadAll() ([]byte, error) {
	data, err := ioutil.ReadAll(c.conn)
	if err != nil {
		return nil, fmt.Errorf("error reading until EOF. %w", err)
	}
	return data, nil
}

// ReadBlock read and parse message length before read all
func (c *Conn) ReadBlock() (string, error) {
	lengthHex, err := c.Read(4)
	if err != nil {
		return "", fmt.Errorf("error read the length of the Block.\n%w", err)
	}
	length, err := strconv.ParseInt(string(lengthHex), 16, 64)
	if err != nil {
		return "", fmt.Errorf("error parse hex length %v.\n%w", lengthHex, err)
	}

	// Clip the length to 255, as per the Google implementation.
	//if length > MaxMessageLength {
	//	length = MaxMessageLength
	//}
	msg, err := c.ReadMessage(int(length))
	if err != nil {
		return "", fmt.Errorf("error read block.\n%w", err)
	}
	return msg, nil
}

// ReadStatus status
func (c *Conn) ReadStatus() (string, error) {
	status, err := c.ReadMessage(4)
	if err != nil {
		return "", fmt.Errorf("error reading status.\n%w", err)
	}
	if isFailureStatus(status) {
		msg, err := c.ReadAll()
		if err != nil {
			return "", fmt.Errorf("error read the reply error message from server.\n%w", err)
		}
		return "", fmt.Errorf("%s", msg)
	}

	return status, nil
}
