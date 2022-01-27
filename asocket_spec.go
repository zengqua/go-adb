package adb

import (
	"fmt"
	"strconv"
	"strings"
)

func (s *Asocket) Connect(address string, port int, serial string) {
	if strings.HasPrefix(address, "tcp:") {
	}
}

func parseTcpSocketSpec(spec string) (string, int, error) {
	if !strings.HasPrefix(spec, "tcp:") {
		return "", -1, fmt.Errorf("specification is not tcp: %s", spec)
	}
	var hostname string
	var port int64

	addr := spec[4:]
	hostAndPort := strings.Split(addr, ":")
	hostname = hostAndPort[0]
	if hostname == "" {
		return "", -1, fmt.Errorf("hostname is empty")
	}
	if len(hostAndPort) == 1 {
		port = defaultAdbLocalTransportPort
	} else if len(hostAndPort) == 2 {
		port, err := strconv.ParseInt(spec[4:], 10, 32)
		if err != nil {
			return "", -1, fmt.Errorf("coudn't parse port: %w", err)
		}
		if port < 0 || port > 65535 {
			return "", -1, fmt.Errorf("bad port number %d", port)
		}
	} else {
		return "", -1, fmt.Errorf("coundn't parse spec")
	}
	return hostname, int(port), nil
}
