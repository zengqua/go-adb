package adb

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

const (
	defaultAdbLocalTransportPort = 5555
)

type Adb struct {
	Host string
	Port int
}

func (a *Adb) Connect(service string) (*Asocket, error) {
	log.Debug(fmt.Sprintf("adb_connect: service: %s", service))
	if service == "" || len(service) > maxPayload {
		return nil, fmt.Errorf("bad service name length (%d)", len(service))
	}

}

func (c *Adb) Dial() (*Conn, error) {
	address := fmt.Sprintf("%s:%d", c.dopts.host, c.dopts.port)
	conn, err := Dial(address)
	if err != nil {
		// Attempt to start the server and try again.
		if err = c.StartServer(); err != nil {
			return nil, fmt.Errorf("error starting server for dial: %s", err)
		}

		conn, err = Dial(address)
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}

func (c *Adb) StartServer() error {
	cmd := exec.Command(c.dopts.adbExecPath,
		"-L", fmt.Sprintf("tcp:%s:%d", c.dopts.host, c.dopts.port),
		"start-server")
	if output, err := cmd.CombinedOutput(); err != nil {
		stdoutStderr := strings.TrimSpace(string(output))
		return fmt.Errorf("error starting server: %s\noutput:\n%s", err, stdoutStderr)
	}
	return nil
}

func (c *Adb) Device(serial string) *Device {
	return &Device{
		adb:    c,
		serial: serial,
	}
}

// Command exec through adb in the local machine
func (c *Adb) Command(arg ...string) (string, error) {
	cmd := exec.Command(c.dopts.adbExecPath, arg...)
	output, err := cmd.CombinedOutput()
	stdoutStderr := strings.TrimSpace(string(output))
	if err != nil {
		return stdoutStderr, fmt.Errorf("error run command: %w\noutput:\n%s", err, stdoutStderr)
	}
	return stdoutStderr, nil
}
