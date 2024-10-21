package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

// #bug https://github.com/vesoft-inc/nebula-ng-tools/issues/553
func TestConfigAgentPort(t *testing.T) {
	config := `
version: 1.0
kind: NebulaCluster
rollback: true
installPath: "/data2/luban/knife/nebula-ng"
certFile: certs/client.crt
keyFile: certs/client.key
caFile: certs/ca.crt
spec:
  metad:
    packagePath: "nebula-graph-2024.07.16-nightly-x86_64-glibc-2.31.sh"
    hosts:
      - host:
        ip: "192.168.15.31"
        port: 9559
        agentPort: 6688
        agents:
          host: "192.168.15.31:6688"
    serviceGroups:
      - zoneList: [zone1,zone2,zone3]
        name: svcgrp_test
        replica: 1
        graphd:
          hosts:
            - host:
              ip: "192.168.15.31"
              port: 9669
              agentPort: 6688
              agents:
                host: "192.168.15.31:6688"
        storaged:
          hosts:
            - host:
              ip: "192.168.15.31"
              port: 9779
              agentPort: 6688
              agents:
                host: "192.168.15.31:6688"
`
	err := yaml.Unmarshal([]byte(config), &ConfigSpec)
	if err != nil {
		t.Fatalf("unmarshal config failed: %v", err)
	}
	hostList, err := DeriveHostList("192.168.15.31", "svcgrp_test", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hostList[0].AgentPort)
	assert.Equal(t, "6688", hostList[0].AgentPort)
}
