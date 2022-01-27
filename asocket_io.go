package adb

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

const (
	maxPayload = 1024 * 1024

	StatusSuccess  string = "OKAY"
	StatusFailure         = "FAIL"
	StatusSyncData        = "DATA"
	StatusSyncDone        = "DONE"
	StatusDent            = "DENT"
	StatusNone            = ""
)

// SendProtocolString send a protocol-format string;
// a four hex digit length followed by the string data.
func (s *Asocket) SendProtocolString(msg string) error {
	length := len(msg)
	if length > maxPayload-4 {
		return errors.New("message too long")
	}
	str := fmt.Sprintf("%04x%s", length, msg)
	if err := s.WriteConnExactly([]byte(str)); err != nil {
		return err
	}
	return nil
}

// ReadProtocolString Reads a protocol-format string;
// a four hex digit length followed by the string data.
func (s *Asocket) ReadProtocolString() (string, error) {
	buf, err := s.ReadConnExactly(4)
	if err != nil {
		return "", fmt.Errorf("protocol fault (couldn't read status length)\n%w", err)
	}
	length, err := strconv.ParseInt(string(buf), 16, 64)
	if err != nil {
		return "", fmt.Errorf("protocol fault (couldn't parse status length\n%w", err)
	}
	buf, err = s.ReadConnExactly(int(length))
	if err != nil {
		return "", fmt.Errorf("protocol fault (couldn't read status message)\n%w", err)
	}
	return string(buf), nil
}

// WriteConnExactly writes exactly len bytes to conn.
func (s *Asocket) WriteConnExactly(buf []byte) error {
	p := 0
	for p < len(buf) {
		n, err := s.conn.Write(buf[p:])
		if err != nil {
			return err
		}
		p += n
	}
	return nil
}

// ReadConnExactly reads exactly len bytes from conn.
func (s *Asocket) ReadConnExactly(length int) ([]byte, error) {
	buf := make([]byte, length)
	p := 0
	for p < len(buf) {
		n, err := s.conn.Read(buf[p:])
		if err != nil {
			return nil, err
		}
		p += n
	}
	return buf, nil
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
