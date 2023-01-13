//go:generate mockgen -source=result_set.go -destination result_set_mock.go -package client ResultSet
package client

type ResultSet interface {
	IsSucceed() bool
	GetLatency() int64
	GetError() error
	IsPermanentError() bool
	IsRetryMoreError() bool
}
