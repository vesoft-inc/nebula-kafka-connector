package client

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"
)

var _ = Describe("nebulaLogger", func() {
	It("newNebulaLogger", func() {
		l := newNebulaLogger(logger.NopLogger)
		l.Info("")
		l.Warn("")
		l.Error("")
		l.Fatal("")
	})
})
