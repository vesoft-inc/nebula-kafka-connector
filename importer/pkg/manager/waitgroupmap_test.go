package manager

import (
	stderrors "errors"
	"fmt"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("WaitGroupMap", func() {
	It("concurrency", func() {
		wgMap := NewWaitGroups()

		keyNum := 25
		concurrencyPreKey := 200

		var wgWaitAllKeys sync.WaitGroup
		wgWaitAllKeys.Add(keyNum)

		finish := make(chan struct{})
		go func() {
			wgWaitAllKeys.Wait()
			close(finish)
		}()

		for i := 0; i < keyNum; i++ {
			key := fmt.Sprintf("key%d", i)

			var wgAddKeys sync.WaitGroup
			wgAddKeys.Add(concurrencyPreKey)

			// add
			go func(key string) {
				for i := 0; i < concurrencyPreKey; i++ {
					go func() {
						wgMap.Add(key, 1)
						wgAddKeys.Done()
					}()
				}
			}(key)

			// done
			go func(key string) {
				wgAddKeys.Wait()
				for i := 0; i < concurrencyPreKey; i++ {
					go func() {
						wgMap.Done(key)
					}()
				}
			}(key)

			// wait
			go func(key string) {
				wgAddKeys.Wait()
				wgMap.Wait(key)
				wgWaitAllKeys.Done()
			}(key)
		}

		select {
		case <-finish:
		case <-time.After(time.Second * 10):
			Expect(stderrors.New("timeout")).NotTo(HaveOccurred())
		}
	})
})
