package main

import (
	"fmt"
	"time"

	nebula "github.com/vesoft-inc/nebula-ng/go/client"
)

const (
	address  = "127.0.0.1"
	port     = 9669
	username = "root"
	password = "nebula"
)

// Initialize logger
var log = nebula.DefaultLogger{}

func main() {
	fmt.Println("Basic example starts ...")
	hostAddr := nebula.HostAddress{Host: address, Port: port}

	// Build connection
	connection := nebula.NewConnection(hostAddr)
	err := connection.Open(hostAddr, 1000*time.Millisecond, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Authenticate to get the identifier
	authResp, err := connection.Authenticate(username, password)
	if err != nil {
		log.Fatal(err.Error())
	}
	if string(authResp.GetGqlStatus().Status) != "Success" {
		log.Fatal(fmt.Sprintf("authentication failed, error: %s", string(authResp.GetGqlStatus().Status)))
	}
	log.Info(fmt.Sprintf("Authentication with Identifier: %d succeed", authResp.GetIdentifier()))

	// Build session
	session := nebula.NewSession(authResp.GetIdentifier(), connection, log)
	defer session.Release()

	// Execute a query
	resp, err := session.Execute("INSERT NODE node_type ({id:1, age:1})")
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to execute, error: %s", err.Error()))
	}
	log.Info("Execution response received")

	if string(resp.GetExecutionOutcome().GetGqlStatus().Status) != "Success" {
		log.Fatal(fmt.Sprintf("execute failed, error: %s", string(resp.GetExecutionOutcome().GetGqlStatus().Status)))
	}

	fmt.Println("Basic example finished")
}
