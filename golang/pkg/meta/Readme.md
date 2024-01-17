# Meta Client

By default, it would retry one more time.
So if the first address is not the meta leader, would return succsessfully
in next request.

```golang
package main

import (
 "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

func main() {
 c, err := meta.NewMetaClient("192.168.8.6:10015")
 if err != nil {
  panic(err)
 }
 req := meta.NewCreateClusterReq("root", 3, nil)
 resp, err := c.CreateCluster(req)
 if err != nil {
  panic(err)
 }
 defer c.Close()
 if resp.GetErrorCode() == meta.ErrorClusterExisted.Code() {
  panic(resp.GetErrorMsg())
 }
}

```
