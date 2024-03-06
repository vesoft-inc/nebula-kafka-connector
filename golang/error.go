package nebula_ng

import (
	"fmt"
)

type ErrorCode string

var (
	// Error in client side
	// Tmp error code, would be redefined in future
	ErrorAddressNotValid    ErrorCode = newErrorCode("99", "000")
	ErrorCannotOpen                   = newErrorCode("99", "001")
	ErrorConnIsBroken                 = newErrorCode("99", "002")
	ErrorConnConnectTimeout           = newErrorCode("99", "003")
	ErrorConnRequestTimeout           = newErrorCode("99", "004")
	ErrorConnIsClosed                 = newErrorCode("99", "005")
	ErrorWaitPoolTimeout              = newErrorCode("99", "006")
	ErrorIllegal                      = newErrorCode("99", "007")
	ErrorType                         = newErrorCode("99", "008")
	ErrorClientInternel               = newErrorCode("99", "009")

	// Error in server side
	ErrorSuccessfulCompletion = newErrorCode("00", "000")
	//TODO need to add more error codes
	ErrorLeaderChange          = newErrorCode("ND", "005")
	ErrorClusterExisted        = newErrorCode("NI", "001")
	ErrServiceStaticPortExists = newErrorCode("NM", "019")
)

// TODO add error code in future
type NebulaError struct {
	errorCode   ErrorCode
	errorFormat string
	errorArgs   []interface{}
}

func (e *NebulaError) Error() string {
	msg := fmt.Sprintf(e.errorFormat, e.errorArgs...)
	return fmt.Sprintf("[%s]: %s", e.errorCode, msg)
}

func (e *NebulaError) Code() ErrorCode {
	return e.errorCode
}

func errAddressNotValid(address string, msg string) error {
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
	return &NebulaError{
		errorCode:   ErrorAddressNotValid,
		errorFormat: format,
		errorArgs:   args,
	}
}

func errConnCannotOpen(host string, port int, msg string) error {
	return &NebulaError{
		errorCode:   ErrorCannotOpen,
		errorFormat: "cannot open connection to %s:%d, %s",
		errorArgs:   []interface{}{host, port, msg},
	}
}

func errConnBroken(host string, port int) error {
	return &NebulaError{
		errorCode:   ErrorConnIsBroken,
		errorFormat: "connection to %s:%d is broken",
		errorArgs:   []interface{}{host, port},
	}
}

func errConnConnectTimeout(host string, port int) error {
	return &NebulaError{
		errorCode:   ErrorConnConnectTimeout,
		errorFormat: "connection to %s:%d timeout",
		errorArgs:   []interface{}{host, port},
	}
}

func errConnRequestTimeout(host string, port int) error {
	return &NebulaError{
		errorCode:   ErrorConnRequestTimeout,
		errorFormat: "request to %s:%d timeout",
		errorArgs:   []interface{}{host, port},
	}
}

func errConnIsClosed(host string, port int) error {
	return &NebulaError{
		errorCode:   ErrorConnIsClosed,
		errorFormat: "connection to %s:%d is closed",
		errorArgs:   []interface{}{host, port},
	}
}

func errWaitPoolTimeout() error {
	return &NebulaError{
		errorCode:   ErrorWaitPoolTimeout,
		errorFormat: "get from pool timeout",
		errorArgs:   []interface{}{},
	}
}

func errIllegal(msg string) error {
	return &NebulaError{
		errorCode:   ErrorIllegal, // TODO need to add error code
		errorFormat: "Illegal error, %s",
		errorArgs:   []interface{}{msg},
	}
}

func errType(msg string) error {
	return &NebulaError{
		errorCode:   ErrorType,
		errorFormat: "Type error, %s",
		errorArgs:   []interface{}{msg},
	}
}

// client internel error
// user should not see this error
func errInternel(msg string) error {
	return &NebulaError{
		errorCode:   ErrorClientInternel,
		errorFormat: "Internel error, %s",
		errorArgs:   []interface{}{msg},
	}
}

func errServerResponse(code string, msg string) error {
	// TODO should set the error code
	return &NebulaError{
		errorCode:   ErrorCode(code),
		errorFormat: "%s",
		errorArgs:   []interface{}{msg},
	}
}

func newErrorCode(class, subClass string) ErrorCode {
	if len(class) != 2 || len(subClass) != 3 {
		return ""
	}
	return ErrorCode(class + subClass)
}

func ErrorFromInt(c uint64) ErrorCode {
	encodeSubClass := c >> 16
	encodeClass := c & 0x0000FFFF
	class := string([]byte{byte(encodeClass), byte(encodeClass >> 8)})
	subClass := string([]byte{byte(encodeSubClass), byte(encodeSubClass >> 8),
		byte(encodeSubClass >> 16)})
	return ErrorCode(class + subClass)
}

func (e ErrorCode) codeInt() uint64 {
	class, subClass := e[0:2], e[2:5]
	encodeClass := uint64(class[1])<<8 | uint64(class[0])
	encodeSubClass := uint64(subClass[2])<<16 | uint64(subClass[1])<<8 | uint64(subClass[0])
	return encodeSubClass<<16 | encodeClass
}
