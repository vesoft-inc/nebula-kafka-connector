# Example

```golang
package main

import (
	"fmt"
	"time"
	nrpc "github.com/vesoft-inc/nebula-ng-tools/golang/nrpc"
)

func main() {
	client := nrpc.NewClient("localhost:52104")

	req := []byte("Hello")
	resp, err := client.Send(req, 3 * time.Millisecond)
	if err != nil {
		nerr, ok := err.(nrpc.Error)
		if ok {
			if nerr.Timeout() {
				fmt.Printf("Timeout: ")
				// `client' is still healthy to use
			} else if nerr.BadChannel() {
                // `client' holds a bad TCP connection,
                // better to keep invoking `Reconnect' until success to guarrantee
                // following `Send' will success
				fmt.Printf("BadChannel: ")
                fmt.Println("Reconnecting...")
                if err = client.Reconnect(3, time.Second); err != nil {
                    fmt.Println(err)
                }
			} else {
                // Other cases, e.g. ErrBadArg
            }
		}
		fmt.Println(err)
		return
	}
	fmt.Println(string(resp))
}
```
