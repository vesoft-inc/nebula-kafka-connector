package main

import (
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

func main() {
	c, err := meta.NewMetaClient("192.168.8.6:10000")
	if err != nil {
		panic(err)
	}
	req := meta.NewCreateClusterReq("root", 3, nil)
	resp, err := c.CreateCluster(req)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	if resp.GetErrorCode() == nebula.ErrorClusterExisted {
		panic(resp.GetErrorMsg())
	}
}
