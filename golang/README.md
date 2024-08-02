# nebula-ng-go

Official NebulaGraph Go client which communicates with Nebula service using [gRPC](https://grpc.io/).

## Install & Update

We now only support install by source code.

Could build the source code in `package` workspace.

```bash
workspace=/app/myapp
# get the source code, e.g. nebula-golang-5.0.0.tar.gz
tar zxvf nebula-golang-5.0.0.tar.gz -o /tmp/golang
# in your project, run go mod edit
cd ${workspace}
go mod edit -replace github.com/vesoft-inc/nebula-ng-tools/golang=/tmp/golang
```

for more details, pleader refer: [gomod ref](https://go.dev/doc/modules/gomod-ref#replace)

## Usage example

```golang
package main

import (
 "fmt"

 nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

const (
 address  = "127.0.0.1:10025"
 username = "root"
 password = "NebulaGraph01"
)

// Initialize logger
var log = nebula.DefaultLogger

func basicClient() {
 client, err := nebula.NewNebulaClient(address, username, password)
 if err != nil {
  log.Error(err.Error())
  return
 }
 resp, err := client.Execute("return 1 as a")
 if err != nil {
  log.Error(err.Error())
  return
 }
 log.Info(fmt.Sprintf("columns: %v", resp.Columns()))
 for resp.HasNext() {
  row, err := resp.Next()
  if err != nil {
   log.Error(err.Error())
   return
  }
  v1, err := row.GetValueByIndex(0)
  if err != nil {
   log.Error(err.Error())
   return
  }
  log.Info(v1.String())
  v2, err := row.GetValueByName("a")
  if err != nil {
   log.Error(err.Error())
   return
  }
  log.Info(v2.String())
 }
}

func basicPool() {
 pool, err := nebula.NewNebulaPool(address, username, password)
 if err != nil {
  log.Error(err.Error())
  return
 }
 client, err := pool.GetClient()
 if err != nil {
  log.Error(err.Error())
  return
 }
 resp, err := client.Execute("return 1 as a")
 if err != nil {
  log.Error(err.Error())
  return
 }
 log.Info(fmt.Sprintf("columns: %v", resp.Columns()))
 for resp.HasNext() {
  row, err := resp.Next()
  if err != nil {
   log.Error(err.Error())
   return
  }
  v1, err := row.GetValueByIndex(0)
  if err != nil {
   log.Error(err.Error())
   return
  }
  log.Info(v1.String())
  v2, err := row.GetValueByName("a")
  if err != nil {
   log.Error(err.Error())
   return
  }
  log.Info(v2.String())
 }
}

func main() {
 basicClient()
 basicPool()
}

```
