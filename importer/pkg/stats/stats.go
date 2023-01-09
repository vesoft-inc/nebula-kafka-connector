package stats

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

type (
	Stats struct {
		StartTime       time.Time     // The time to start statistics.
		ProcessedBytes  int64         // The processed bytes.
		TotalBytes      int64         // The total bytes.
		FailedRecords   int64         // The number of records that have failed to be processed.
		TotalRecords    int64         // The number of records that have been processed.
		FailedRequest   int64         // The number of requests that have failed.
		TotalRequest    int64         // The number of requests that have been processed.
		TotalLatency    time.Duration // The cumulative latency.
		TotalReqTime    time.Duration // The cumulative request time.
		FailedProcessed int64         // The number of nodes and edges that have failed to be processed.
		TotalProcessed  int64         // The number of nodes and edges that have been processed.
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
		duration           = time.Since(s.StartTime)
		seconds            = duration.Seconds()
		recordsPreSecond   float64
		avgLatency         time.Duration
		avgReqTime         time.Duration
		requestPreSecond   float64
		processedPreSecond float64
	)

	if s.TotalRecords > 0 {
		recordsPreSecond = float64(s.TotalRecords) / seconds
	}

	if s.TotalRequest > 0 {
		avgLatency = s.TotalLatency / time.Duration(s.TotalRequest)
		avgReqTime = s.TotalReqTime / time.Duration(s.TotalRequest)
		requestPreSecond = float64(s.TotalRequest) / seconds
	}
	if s.TotalProcessed > 0 {
		processedPreSecond = float64(s.TotalProcessed) / seconds
	}

	return fmt.Sprintf("%s "+
		"%.2f%%(%s/%s) "+
		"Records{Finished: %d, Failed: %d, Rate: %.2f/s}, "+
		"Requests{Finished: %d, Failed: %d, Latency: %s/%s, Rate: %.2f/s}, "+
		"Processed{Finished: %d, Failed: %d, Rate: %.2f/s}",
		duration.Truncate(time.Second),
		s.Percentage(), humanize.IBytes(uint64(s.ProcessedBytes)), humanize.IBytes(uint64(s.TotalBytes)),
		s.TotalRecords, s.FailedRecords, recordsPreSecond,
		s.TotalRequest, s.FailedRequest, avgLatency, avgReqTime, requestPreSecond,
		s.TotalProcessed, s.FailedProcessed, processedPreSecond,
	)
}
