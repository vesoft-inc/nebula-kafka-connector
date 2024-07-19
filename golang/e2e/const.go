package e2e

import "fmt"

const (
	nebulaHost     = "127.0.0.1"
	nebulaPort     = 9669
	nebulaUser     = "root"
	nebulaPassword = "NebulaGraph01"
)

var nebulaAddress = fmt.Sprintf("%s:%d", nebulaHost, nebulaPort)
