package main

import (
	"fmt"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

func main() {
	c, err := meta.NewMetaClient("192.168.8.6:10025")
	if err != nil {
		panic(err)
	}
	req := meta.NewCreateClusterReq("root", 3, nil)
	resp, err := c.CreateCluster(req)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	if resp.Code == nebula.ErrorClusterExisted {
		panic(fmt.Sprintf("msg:%s, code:%s", resp.GetErrorMsg(), resp.GetErrorCode()))
	}
}
