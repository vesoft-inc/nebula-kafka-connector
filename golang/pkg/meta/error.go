package meta

// the error code of metad
type ErrorCode struct {
	class    string
	subClass string
}

func (ec *ErrorCode) Code() uint64 {
	if len(ec.class) != 3 || len(ec.subClass) != 3 {
		return 0
	}
	c1, c2, c3 := uint32(ec.class[2]), uint32(ec.class[1]), uint32(ec.class[0])
	sc1, sc2, sc3 := uint32(ec.subClass[2]), uint32(ec.subClass[1]), uint32(ec.subClass[0])

	encodeClass := c1<<16 | c2<<8 | c3
	encodeSubClass := sc1<<16 | sc2<<8 | sc3
	return (uint64(encodeSubClass) << 24) | uint64(encodeClass)
}

func newErrorCode(class, subClass string) *ErrorCode {
	return &ErrorCode{
		class:    class,
		subClass: subClass,
	}
}

var (
	//TODO need to add more error codes
	ErrorLeaderChange   = newErrorCode("XND", "004")
	ErrorClusterExisted = newErrorCode("XNP", "001")
)
