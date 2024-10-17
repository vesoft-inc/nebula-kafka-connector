package main

import (
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

func main() {
	c, err := meta.NewMetaClient("192.168.8.6:10025")
	if err != nil {
		panic(err)
	}
	req := meta.NewCreateServiceGroupReq("root", 3, "root", nil)
	err = c.CreateServiceGroup(req)
	if err != nil {
		panic(err)
	}
	defer c.Close()

}
