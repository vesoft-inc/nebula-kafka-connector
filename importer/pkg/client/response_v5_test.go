//go:build linux

package client

import (
	"time"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"

	"github.com/agiledragon/gomonkey/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("defaultResponseV5", func() {
	It("newResponseV5", func() {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		rs := nebula.ResultSet{}
		resp := newResponseV5(&rs, time.Second)

		patches.ApplyMethodReturn(rs, "IsSucceed", true)

		err := resp.GetError()
		Expect(err).NotTo(HaveOccurred())

		patches.Reset()

		patches.ApplyMethodReturn(rs, "GetLatency", int64(1))
		patches.ApplyMethodReturn(rs, "IsSucceed", false)
		patches.ApplyMethodReturn(rs, "GetStatus", "test status")

		err = resp.GetError()
		Expect(err).To(HaveOccurred())

		Expect(resp.GetLatency()).To(Equal(time.Microsecond))
		Expect(resp.GetRespTime()).To(Equal(time.Second))
		Expect(resp.IsPermanentError()).To(BeFalse())
		Expect(resp.IsRetryMoreError()).To(BeFalse())
	})
})
