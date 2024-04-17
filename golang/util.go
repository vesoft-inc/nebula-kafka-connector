package nebula_ng

import (
	"net"
	"strconv"
	"strings"
)

func parseHostPort(address string) (string, int, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", 0, errAddressNotValid(address, err.Error())
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return "", 0, errAddressNotValid(address, err.Error())
	}

	return host, p, nil
}

func parseAddresses(addresses string) ([]*hostAddress, error) {
	var hostAddresses []*hostAddress
	addrs := strings.Split(addresses, ",")
	for _, addr := range addrs {
		if addr == "" {
			continue
		}
		host, port, err := parseHostPort(addr)
		if err != nil {
			return nil, err
		}
		hostAddress := &hostAddress{
			host: host,
			port: port,
		}
		hostAddresses = append(hostAddresses, hostAddress)
	}
	return hostAddresses, nil
}

func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	ne, ok := err.(*NebulaError)
	if !ok {
		return false
	}
	codes := []ErrorCode{
		ERROR_CONN_IS_BROKEN,
		ERROR_CONN_CONNECT_TIMEOUT,
		ERROR_CONN_REQUEST_TIMEOUT,
		ERROR_CONN_IS_CLOSED,
	}
	for _, c := range codes {
		if ne.errorCode == c {
			return true
		}
	}
	return false
}
