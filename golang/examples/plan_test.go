package main

import (
	"encoding/json"
	"testing"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

func TestValueFormat(t *testing.T) {

	client, err := nebula.NewNebulaClient("192.168.8.145:9669", "nebula", "nebula")

	if err != nil {
		t.Error(err)
	}
	if err != nil {
		panic(err.Error())
	}
	defer client.Close()
	resp, err := client.Execute("profile return 1")
	if err != nil {
		t.Error(err)
	}
	log.Info("Execution response received")
	plan := resp.PlanDesc()
	jsonPlan, err := json.MarshalIndent(plan.GetPlanDesc(), "", "  ")
	if err != nil {
		t.Error(err)
	}
	log.Info(string(jsonPlan))
}
