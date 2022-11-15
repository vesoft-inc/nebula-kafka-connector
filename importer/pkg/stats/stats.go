package stats

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

type (
	Stats struct {
		StartTime      time.Time
		ProcessedBytes int64
		TotalBytes     int64
		FailedBatches  int64
		TotalBatches   int64
		FailedRecords  int64
		TotalRecords   int64
		TotalLatency   time.Duration
		TotalReqTime   time.Duration
	}
)

func (s *Stats) Percentage() float64 {
	if s.TotalBytes == 0 {
		return 0
	}
	return float64(s.ProcessedBytes) / float64(s.TotalBytes) * 100
}

func (s *Stats) String() string {
	var (
		duration   = time.Since(s.StartTime)
		avgLatency time.Duration
		avgReqTime time.Duration
		rps        float64
	)

	if s.TotalRecords > 0 {
		avgLatency = s.TotalLatency / time.Duration(s.TotalBatches)
		avgReqTime = s.TotalReqTime / time.Duration(s.TotalBatches)
		rps = float64(s.TotalRecords) / duration.Seconds()
	}

	return fmt.Sprintf("Time(%s), "+
		"Processed %.2f%%(%s/%s) "+
		"Finished(%d), Failed(%d), "+
		"Latency AVG(%s), Batches Req AVG(%s), Rows AVG(%.2f/s)",
		duration.Truncate(time.Second),
		s.Percentage(), humanize.IBytes(uint64(s.ProcessedBytes)), humanize.IBytes(uint64(s.TotalBytes)),
		s.TotalRecords, s.FailedRecords,
		avgLatency, avgReqTime, rps,
	)
}
