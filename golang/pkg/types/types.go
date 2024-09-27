package types

import (
	"context"
)

type (
	Result interface {
		Summary() Summary
		Cursor() []byte
		Table
	}

	Table interface {
		RowSize() int
		HasNext() bool
		Next() (Row, error)
		// Scan copies the columns in the current row into the values pointed at by dest.
		// If there's no more row, return io.EOF
		Scan(...any) error
		Columns() []string
		ColumnTypes() []ColumnType //not support yet
	}

	Row interface {
		Values() []Value
		GetValueByName(name string) (Value, error)
		GetValueByIndex(index int) (Value, error)
	}

	Summary interface {
		BuildTimeUs() int64
		OptimizeTimeUs() int64
		SerializeTimeUs() int64
		TotalServerTimeUs() int64
		ExplainType() string
		PlanInfo() PlanInfo
		QueryStats() QueryStats
	}
	PlanInfo interface {
		Id() string
		Name() string
		Details() string
		Columns() []string
		TimeMs() float64
		Rows() int64
		MemoryKib() float64
		BlockedMs() float64
		QueuedMs() float64
		ConsumeMs() float64
		ProduceMs() float64
		FinishMs() float64
		Batches() int64
		Concurrency() int64
		OtherStatsJson() []byte
		Children() []PlanInfo
	}
	QueryStats interface {
		NumAffectedNodes() int64
		NumAffectedEdges() int64
	}

	Client interface {
		Execute(stmt string) (Result, error)
		ExecuteContext(ctx context.Context, stmt string) (Result, error)
		Ping() error // by default, timeout is 1s.
		PingContext(ctx context.Context) error
		IsClosed() bool
		Close() error
		GetSessionId() (int64, error)
	}

	Pool interface {
		GetClient() (Client, error)
		PutClient(Client) error
		Close() error
	}

	ColumnType int
)

const (
	ColumnTypeInvalid ColumnType = iota
	ColumnTypeNode
	ColumnTypeEdge
	ColumnTypePath
	ColumnTypeUnknown
	ColumnTypeBool
	ColumnTypeInt8
	ColumnTypeUint8
	ColumnTypeInt16
	ColumnTypeUint16
	ColumnTypeInt32
	ColumnTypeUint32
	ColumnTypeInt64
	ColumnTypeUint64
	ColumnTypeFloat32
	ColumnTypeFloat64
	ColumnTypeString
	ColumnTypeList
	ColumnTypeRecord
	ColumnTypeLocalTime
	ColumnTypeLocalDatetime
	ColumnTypeZonedTime
	ColumnTypeZonedDatetime
	ColumnTypeDate
	ColumnTypeDuration
	ColumnTypeDecimal
	ColumnTypeAny = 0xFF
)
