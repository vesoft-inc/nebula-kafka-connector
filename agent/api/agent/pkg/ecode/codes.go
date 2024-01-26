package ecode

const (
	PlatformCode = 0 // TODO: Please modify it to your own platform code.
)

// Define you error code here
var (
	ErrBadRequest     = newErrCode(CCBadRequest, PlatformCode, 0, "ErrBadRequest")         // 40000000
	ErrParam          = newErrCode(CCBadRequest, PlatformCode, 1, "ErrParam")              // 40000001
	ErrUnauthorized   = newErrCode(CCUnauthorized, PlatformCode, 0, "ErrUnauthorized")     // 40100000
	ErrForbidden      = newErrCode(CCForbidden, PlatformCode, 0, "ErrForbidden")           // 40300000
	ErrNotFound       = newErrCode(CCNotFound, PlatformCode, 0, "ErrNotFound")             // 40400000
	ErrInternalServer = newErrCode(CCInternalServer, PlatformCode, 0, "ErrInternalServer") // 50000000
	ErrNotImplemented = newErrCode(CCNotImplemented, PlatformCode, 0, "ErrNotImplemented") // 50100000
	ErrUnknown        = newErrCode(CCUnknown, PlatformCode, 0, "ErrUnknown")               // 90000000
)
