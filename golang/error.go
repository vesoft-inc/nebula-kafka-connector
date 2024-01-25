package nebula_ng

import (
	"fmt"
)

var (
	// Error in client side
	// Tmp error code, would be redefined in future
	ErrorAddressNotValid    string = newErrorCode("99", "000")
	ErrorCannotOpen                = newErrorCode("99", "001")
	ErrorConnIsBroken              = newErrorCode("99", "002")
	ErrorConnConnectTimeout        = newErrorCode("99", "003")
	ErrorConnRequestTimeout        = newErrorCode("99", "004")
	ErrorConnIsClosed              = newErrorCode("99", "005")
	ErrorWaitPoolTimeout           = newErrorCode("99", "006")
	ErrorIllegal                   = newErrorCode("99", "007")
	ErrorType                      = newErrorCode("99", "008")
	ErrorClientInternel            = newErrorCode("99", "009")

	// Error in server side
	ErrorSuccessfulCompletion = newErrorCode("00", "000")
	//TODO need to add more error codes
	ErrorLeaderChange   = newErrorCode("ND", "004")
	ErrorClusterExisted = newErrorCode("NI", "001")
)

// TODO add error code in future
type NebulaError struct {
	code        string
	errorFormat string
	errorArgs   []interface{}
}

func (e *NebulaError) Error() string {
	return fmt.Sprintf(e.errorFormat, e.errorArgs...)
}

func (e *NebulaError) Code() string {
	return e.code
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
		code:        ErrorAddressNotValid,
		errorFormat: format,
		errorArgs:   args,
	}
}

func errConnCannotOpen(host string, port int, msg string) error {
	return &NebulaError{
		code:        ErrorCannotOpen,
		errorFormat: "cannot open connection to %s:%d, %s",
		errorArgs:   []interface{}{host, port, msg},
	}
}

func errConnBroken(host string, port int) error {
	return &NebulaError{
		code:        ErrorConnIsBroken,
		errorFormat: "connection to %s:%d is broken",
		errorArgs:   []interface{}{host, port},
	}
}

func errConnConnectTimeout(host string, port int) error {
	return &NebulaError{
		code:        ErrorConnConnectTimeout,
		errorFormat: "connection to %s:%d timeout",
		errorArgs:   []interface{}{host, port},
	}
}

func errConnRequestTimeout(host string, port int) error {
	return &NebulaError{
		code:        ErrorConnRequestTimeout,
		errorFormat: "request to %s:%d timeout",
		errorArgs:   []interface{}{host, port},
	}
}

func errConnIsClosed(host string, port int) error {
	return &NebulaError{
		code:        ErrorConnIsClosed,
		errorFormat: "connection to %s:%d is closed",
		errorArgs:   []interface{}{host, port},
	}
}

func errWaitPoolTimeout() error {
	return &NebulaError{
		code:        ErrorWaitPoolTimeout,
		errorFormat: "get from pool timeout",
		errorArgs:   []interface{}{},
	}
}

func errIllegal(msg string) error {
	return &NebulaError{
		code:        ErrorIllegal, // TODO need to add error code
		errorFormat: "Illegal error, %s",
		errorArgs:   []interface{}{msg},
	}
}

func errType(msg string) error {
	return &NebulaError{
		code:        ErrorType,
		errorFormat: "Type error, %s",
		errorArgs:   []interface{}{msg},
	}
}

// client internel error
// user should not see this error
func errInternel(msg string) error {
	return &NebulaError{
		code:        ErrorClientInternel,
		errorFormat: "Internel error, %s",
		errorArgs:   []interface{}{msg},
	}
}

func errServerResponse(code string, msg string) error {
	// TODO should set the error code
	return &NebulaError{
		code:        code,
		errorFormat: "%s",
		errorArgs:   []interface{}{msg},
	}
}

func newErrorCode(class, subClass string) string {
	if len(class) != 2 || len(subClass) != 3 {
		return ""
	}
	return class + subClass
}
