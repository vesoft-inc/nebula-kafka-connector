package main

import (
	"fmt"
	"time"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

const (
	address  = "127.0.0.1"
	port     = 16720
	username = "root"
	password = "nebula"
)

// Initialize logger
var log = nebula.DefaultLogger

func printResult(res nebula.Result) error {
	columns := res.Columns()
	for _, col := range columns {
		fmt.Printf("|")
		fmt.Printf("%s ", col)
		fmt.Printf("|\n")
	}

	for res.HasNext() {
		row, err := res.Next()
		if err != nil {
			return err
		}
		values := row.Values()
		for _, v := range values {
			fmt.Printf("|")
			fmt.Printf("%s ", v.String())
			fmt.Printf("|")
		}
		fmt.Println()
	}
	return nil
}

func runQuery(client nebula.Conn, query string) error {
	resp, err := client.Execute(query)
	if err != nil {
		return err
	}
	log.Info("Execution response received")
	if err := printResult(resp); err != nil {
		return err
	}
	log.Info("query: " + query + " was executed successfully")
	return nil
}

func main() {
	fmt.Println("Basic example starts ...")
	addresses := fmt.Sprintf("%s:%d", address, port)
	client, err := nebula.NewNebulaClient(addresses, username, password)
	if err != nil {
		panic(err.Error())
	}
	defer client.Close()

	// Execute a query
	log.Info("Create graph type...")
	runQuery(client, `CREATE GRAPH TYPE graph_type IF NOT EXISTS AS GRAPH TYPE {(node_type(id) LABEL player {id INT, name STRING}),(node_type)-[edge_type LABEL follow {followness INT}]->(node_type)}`)
	time.Sleep(3 * time.Second)

	log.Info("Create graph nba...")
	runQuery(client, `CREATE GRAPH nba IF NOT EXISTS OF graph_type`)
	time.Sleep(5 * time.Second)

	log.Info("Insert a node...")
	runQuery(client, `USE nba INSERT NODE node_type ({id:1, name:"Tim"}),({id:2, name:"Jerry"}),({id:3, name:"Kyle"})`)

	log.Info("Insert an edge...")
	runQuery(client, `USE nba INSERT EDGE edge_type ({id:1})-[{followness:90}]->({id:2}),({id:2})-[{followness:100}]->({id:3})`)

	log.Info("Run a simple query...")
	runQuery(client, `USE nba RETURN 1 + 1`)

	log.Info("Run a MATCH query...")
	runQuery(client, `USE nba MATCH (v) RETURN v`)

	log.Info("Execution response received")

	// run from pool
	pool, err := nebula.NewNebulaPool(addresses, username, password)
	defer pool.Close()
	pool.OnOpenClient(func(conn nebula.ConnSetter) {
		conn.SetRequestTimeout(100 * time.Second)
		conn.SetConnectTimeout(100 * time.Second)
	})
	client2, err := pool.GetClient()
	defer pool.PutClient(client2)
	if err != nil {
		panic(err.Error())
	}
	log.Info("Run a MATCH query...")
	runQuery(client2, `USE nba MATCH (v) RETURN v`)
	log.Info("Execution response received")

	fmt.Println("Basic example finished")
}
