//go:build linux

package client

import (
	"github.com/agiledragon/gomonkey/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	nebula "github.com/vesoft-inc/nebula-go/v3"
)

var _ = Describe("ResultSetV3", func() {
	It("newResultSetV3", func() {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		nRS := nebula.ResultSet{}
		rs := newResultSetV3(&nRS)

		patches.ApplyMethodReturn(nRS, "GetErrorCode", nebula.ErrorCode_SUCCEEDED)
		patches.ApplyMethodReturn(nRS, "GetErrorMsg", "")

		err := rs.GetError()
		Expect(err).NotTo(HaveOccurred())

		patches.Reset()

		patches.ApplyMethodReturn(nRS, "GetErrorCode", nebula.ErrorCode_E_DISCONNECTED)
		patches.ApplyMethodReturn(nRS, "GetErrorMsg", "test msg")

		err = rs.GetError()
		Expect(err).To(HaveOccurred())

		Expect(rs.IsPermanentError()).To(BeFalse())
		Expect(rs.IsRetryMoreError()).To(BeFalse())
	})

	DescribeTable("IsPermanentError",
		func(errorCode nebula.ErrorCode, isPermanentError bool) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			nRS := nebula.ResultSet{}
			rs := newResultSetV3(&nRS)

			patches.ApplyMethodReturn(nRS, "GetErrorCode", errorCode)

			Expect(rs.IsPermanentError()).To(Equal(isPermanentError))
		},
		EntryDescription("%[1]s -> %[2]t"),
		Entry(nil, nebula.ErrorCode_E_SYNTAX_ERROR, true),
		Entry(nil, nebula.ErrorCode_E_SEMANTIC_ERROR, true),
		Entry(nil, nebula.ErrorCode_E_DISCONNECTED, false),
	)

	DescribeTable("IsPermanentError",
		func(errorMsg string, isPermanentError bool) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			nRS := nebula.ResultSet{}
			rs := newResultSetV3(&nRS)

			patches.ApplyMethodReturn(nRS, "GetErrorMsg", errorMsg)

			Expect(rs.IsRetryMoreError()).To(Equal(isPermanentError))
		},
		EntryDescription("%[1]s -> %[2]t"),
		Entry(nil, "x raft buffer is full x", true),
		Entry(nil, "x x", false),
	)
})
