package stats

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stats", func() {
	Describe(".String", func() {
		It("TotalRecords is zero", func() {
			s := &Stats{
				StartTime: time.Now(),
			}
			Expect(s.String()).Should(Equal("Time(0s), Processed 0.00%(0 B/0 B) Finished(0), Failed(0), Latency AVG(0s), Batches Req AVG(0s), Rows AVG(0.00/s)"))
		})
		It("TotalRecords is not zero", func() {
			s := &Stats{
				StartTime:      time.Now().Add(-time.Second * 10),
				ProcessedBytes: 100 * 1024,
				TotalBytes:     300 * 1024,
				FailedBatches:  1,
				TotalBatches:   12,
				FailedRecords:  23,
				TotalRecords:   1234,
				TotalLatency:   time.Second * 12,
				TotalReqTime:   2 * time.Second * 12,
			}
			Expect(s.String()).Should(Equal("Time(10s), Processed 33.33%(100 KiB/300 KiB) Finished(1234), Failed(23), Latency AVG(1s), Batches Req AVG(2s), Rows AVG(123.40/s)"))
		})
	})
})
