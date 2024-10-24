package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// #bug https://github.com/vesoft-inc/nebula-ng-tools/issues/553
func TestConfigAgentPort(t *testing.T) {
	config := `
version: 1.0
installPath: "/data2/luban/knife/nebula-ng"
packagePath: "nebula-graph-2024.07.16-nightly-x86_64-glibc-2.31.sh"
certFile: certs/client.crt
keyFile: certs/client.key
caFile: certs/ca.crt
spec:
  metad:
    config:
      v: 1
    instances:
      - host: "192.168.15.31"
        port: 9559
        agentPort: 6688
  serviceGroups:
    - name: svcgrp_test
      replicaFactor: 1
      graphd:
        config:
          v: 1
        instances:
          - host: "192.168.15.31"
            port: 9669
            agentPort: 6688
          - host: "192.168.15.32"
            port: 9669
            agentPort: 6699
      storaged:
        config:
          v: 1
        instances:
          - host: "192.168.15.31"
            port: 9669
            agentPort: 6688
`
	c, err := NewConfigFromBytes([]byte(config))
	if err != nil {
		t.Fatal(err)
	}
	instances, err := c.GetServiceInstances("svcgrp_test", "graphd", "")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(instances[1].AgentPort)
	assert.Equal(t, 6699, instances[1].AgentPort)
}

func TestConfigValidate(t *testing.T) {
	config := `
version: 1.0
packagePath: "nebula-graph-2024.07.16-nightly-x86_64-glibc-2.31.sh"
installPath: "/data2/luban/knife/nebula-ng"
#certFile: certs/client.crt
#keyFile: certs/client.key
#caFile: certs/ca.crt
spec:
  metad:
    config:
      v: 1
    instances:
      - host: "192.168.15.31"
        port: 9559
        agentPort: 6688
  serviceGroups:
    - name: svcgrp_test
      replicaFactor: 1
      graphd:
        config:
          v: 1
        instances:
          - host: "192.168.15.31"
            port: 9669
            agentPort: 6688
      storaged:
        config:
          v: 1
        instances:
          - host: "192.168.15.31"
            port: 9669
            agentPort: 6688
`
	_, err := NewConfigFromBytes([]byte(config))
	assert.Equal(t, "certFile is empty, keyFile is empty, caFile is empty", err.Error())
}
