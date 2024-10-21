package internal_error

import (
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/errors"
)

func ErrAddressNotValid(address string, msg string) error {
	var (
		format string
		args   []interface{}
	)

	if msg == "" {
		format = "address %s is not valid"
		args = []interface{}{address}
	} else {
		format = "address %s is not valid, %s"
		args = []interface{}{address, msg}
	}
	return errors.NewNebulaError(errors.ERROR_ADDRESS_NOT_VALID, format, args...)
}

func ErrConnCannotOpen(host string, port int, msg string) error {
	return errors.NewNebulaError(errors.ERROR_CANNOT_OPEN, "cannot open connection to %s:%d, %s", host, port, msg)
}

func ErrConnUnavailable(host string, port int) error {
	return errors.NewNebulaError(errors.ERROR_CONN_IS_BROKEN, "connection to %s:%d is unavailable", host, port)
}

func ErrConnBroken(host string, port int) error {
	return errors.NewNebulaError(errors.ERROR_CONN_IS_BROKEN, "connection to %s:%d is broken", host, port)
}

func ErrConnConnectTimeout(host string, port int) error {
	return errors.NewNebulaError(errors.ERROR_CONN_CONNECT_TIMEOUT, "connection to %s:%d timeout", host, port)

}

func ErrConnRequestTimeout(host string, port int) error {
	return errors.NewNebulaError(errors.ERROR_CONN_REQUEST_TIMEOUT, "request to %s:%d timeout", host, port)
}

func ErrConnIsClosed(host string, port int) error {
	return errors.NewNebulaError(errors.ERROR_CONN_IS_CLOSED, "connection to %s:%d is closed", host, port)
}

func ErrWaitPoolTimeout() error {
	return errors.NewNebulaError(errors.ERROR_WAIT_POOL_TIMEOUT, "get from pool timeout")

}

func ErrIllegal(msg string) error {
	return errors.NewNebulaError(errors.ERROR_ILLEGAL, "illegal, %s", msg)
}

func ErrType(msg string) error {
	return errors.NewNebulaError(errors.ERROR_TYPE, "Type error, %s", msg)
}

// client internel error
// user should not see this error
func ErrInternel(msg string) error {
	return errors.NewNebulaError(errors.ERROR_CLIENT_INTERNEL, "Internel error, %s", msg)
}

func ErrServerResponse(code string, msg string) error {
	return errors.NewNebulaError(errors.ErrorCode(code), "%s", msg)
}

func ErrorFromBytes(c []byte) errors.ErrorCode {
	return errors.ErrorCode(c)
}

func ErrTLS(msg string) error {
	return errors.NewNebulaError(errors.ERROR_TLS_ERROR, "TLS error, %s", msg)
}

func ErrUnknownColumnType(columnType int) error {
	return errors.NewNebulaError(errors.ERROR_UNKNOWN_COLUMN_TYPE, "unknown column type %d", columnType)
}
