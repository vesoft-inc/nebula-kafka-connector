package nebula_ng

import (
	"strconv"
	"strings"
)

func parseHostPort(address string) (string, int, error) {
	hostPort := strings.Split(address, ":")
	if len(hostPort) != 2 {
		return "", 0, errAddressNotValid(address, "")
	}
	host := hostPort[0]
	port := hostPort[1]
	pInt, err := strconv.Atoi(port)
	if err != nil {
		return "", 0, errAddressNotValid(address, "port is not valid")
	}
	return host, pInt, nil
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
	codes := []string{
		ErrorConnIsBroken,
		ErrorConnConnectTimeout,
		ErrorConnRequestTimeout,
		ErrorConnIsClosed,
	}
	for _, c := range codes {
		if ne.code == c {
			return true
		}
	}
	return false
}
