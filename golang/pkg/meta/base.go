package meta

type (
	headerRequest struct {
		requestType string
		clusterId   int64 // meta admin dont need clusterId
	}
	HeaderResponse struct {
		Code    uint64
		Msg     string
		NewHost string // only if code is leader changed
		NewPort uint32 // only if code is leader changed
	}

	responseHeader interface {
		//getHeader for internal usage, would reconnect if the leader chanage.
		getHeader() *HeaderResponse
		GetErrorCode() uint64
		GetErrorMsg() string
	}
)

func (hresp *HeaderResponse) getHeader() *HeaderResponse {
	return hresp
}

func (hresp *HeaderResponse) GetErrorCode() uint64 {
	return hresp.Code
}

func (hresp *HeaderResponse) GetErrorMsg() string {
	return hresp.Msg
}
