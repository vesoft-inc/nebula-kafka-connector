//go:build linux

package client

import (
	"github.com/agiledragon/gomonkey/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

var _ = Describe("ResultSetV5", func() {
	It("newResultSetV5", func() {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		nRS := nebula.ResultSet{}
		rs := newResultSetV5(&nRS)

		patches.ApplyMethodReturn(nRS, "IsSucceed", true)

		err := rs.GetError()
		Expect(err).NotTo(HaveOccurred())

		patches.Reset()

		patches.ApplyMethodReturn(nRS, "IsSucceed", false)
		patches.ApplyMethodReturn(nRS, "GetStatus", "test status")

		err = rs.GetError()
		Expect(err).To(HaveOccurred())

		Expect(rs.IsPermanentError()).To(BeFalse())
		Expect(rs.IsRetryMoreError()).To(BeFalse())
	})
})
