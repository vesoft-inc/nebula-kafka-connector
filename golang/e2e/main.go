package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/types"
)

const (
	nebulaHost     = "127.0.0.1"
	nebulaPort     = 9669
	nebulaUser     = "root"
	nebulaPassword = "NebulaGraph01"
)

var nebulaAddress = fmt.Sprintf("%s:%d", nebulaHost, nebulaPort)

func prepareData(c types.Client, scanner *bufio.Scanner) error {
	var stmt string
	var multiLine bool
	var execute bool
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == `"""` {
			if multiLine {
				multiLine = false
				execute = true
			} else {
				multiLine = true
			}
		} else {
			if multiLine {
				stmt += line + "\n"
			} else {
				stmt = line
				execute = true
			}
		}
		if execute {
			execute = false
			if _, err := c.Execute(stmt); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	file, err := os.Open("movie.ngql")
	if err != nil {
		log.Fatalf("Open movie.ngql failed, %s", err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	c, err := nebula.NewNebulaClient(nebulaAddress, nebulaUser, nebulaPassword)
	if err != nil {
		log.Fatalf("NewNebulaClient failed, %s", err.Error())
	}
	defer c.Close()
	if err := prepareData(c, scanner); err != nil {
		log.Fatalf("prepareData failed, %s", err.Error())
	}
	_, err = c.Execute("create schema \"/test_schema\"")
	if err != nil {
		log.Fatalf("create schema failed, %s", err.Error())
	}
	_, err = c.Execute("session set schema \"/test_schema\"")
	if err != nil {
		log.Fatalf("session set schema failed, %s", err.Error())
	}
	_, err = c.Execute(`
	CREATE GRAPH TYPE test_graph_type AS {
		Node Actor (:Person {id INT PRIMARY KEY, name STRING, birthDate Date})
	}`)
	if err != nil {
		log.Fatalf("execute failed, %s", err.Error())
	}
	_, err = c.Execute(`CREATE GRAPH test_graph TYPED test_graph_type`)
	if err != nil {
		log.Fatalf("execute failed, %s", err.Error())
	}

}
