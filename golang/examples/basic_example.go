package main

import (
	"fmt"
	"time"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

const (
	address  = "127.0.0.1"
	port     = 29562
	username = "root"
	password = "nebula"
)

// Initialize logger
var log = nebula.DefaultLogger{}

func printResult(res *nebula.ResultSet) {
	tableStr := res.AsStringTable()

	for _, row := range tableStr {
		fmt.Printf("|")
		for _, col := range row {
			fmt.Printf("%s ", col)
		}
		fmt.Println("|")
	}
}

func runQuery(session *nebula.Session, query string) {
	resp, err := session.Execute(query)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to execute, error: %s", err.Error()))
	}
	log.Info("Execution response received")

	if !resp.IsSucceed() {
		log.Fatal(fmt.Sprintf("execute failed, error: %s", string(resp.GetStatus())))
	}

	if resp.GetRows() != nil {
		printResult(resp)
	}

	log.Info("query: " + query + " was executed successfully")
}

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
	if string(authResp.GetGqlStatus().Status) != "SUCCESS" {
		log.Fatal(fmt.Sprintf("authentication failed, error: %s", string(authResp.GetGqlStatus().Status)))
	}
	log.Info(fmt.Sprintf("Authentication with Identifier: %d succeed", authResp.GetIdentifier()))

	// Build session
	session := nebula.NewSession(authResp.GetIdentifier(), connection, log)
	defer session.Release()

	// Execute a query
	log.Info("Create graph type...")
	runQuery(session, `CREATE GRAPH TYPE graph_type IF NOT EXISTS AS GRAPH TYPE {(node_type(id) LABEL player {id INT, name STRING}),(node_type)-[edge_type LABEL follow {followness INT}]->(node_type)}`)
	time.Sleep(3 * time.Second)

	log.Info("Create graph nba...")
	runQuery(session, `CREATE GRAPH nba IF NOT EXISTS OF graph_type`)
	time.Sleep(5 * time.Second)

	log.Info("Insert a node...")
	runQuery(session, `USE nba INSERT NODE node_type ({id:1, name:"Tim"}),({id:2, name:"Jerry"}),({id:3, name:"Kyle"})`)

	log.Info("Insert an edge...")
	runQuery(session, `USE nba INSERT EDGE edge_type ({id:1})-[{followness:90}]->({id:2}),({id:2})-[{followness:100}]->({id:3})`)

	log.Info("Run a simple query...")
	runQuery(session, `USE nba RETURN 1 + 1`)

	log.Info("Run a MATCH query...")
	runQuery(session, `FROM nba MATCH (v) RETURN v`)

	log.Info("Execution response received")

	fmt.Println("Basic example finished")
}
