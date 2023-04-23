//go:build linux

package client

import (
	stderrors "errors"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/generated_code/v5.0.0/nebula/graph"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"

	"github.com/agiledragon/gomonkey/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SessionV5", func() {
	It("success", func() {
		session := newSessionV5(HostAddress{}, "user", "password", nil)
		connection := nebula.NewConnection(nebula.HostAddress{})
		nSession := &nebula.Session{}
		id := int64(1)
		authResp := &graph.XAuthResponse_{
			Identifier: &id,
		}

		patches := gomonkey.NewPatches()
		defer patches.Reset()

		patches.ApplyFuncReturn(nebula.NewConnection, connection)
		patches.ApplyMethodReturn(connection, "Open", nil)
		patches.ApplyMethodReturn(connection, "Authenticate", authResp, nil)

		patches.ApplyFuncReturn(nebula.NewSession, nSession)

		patches.ApplyMethodReturn(nSession, "Execute", &nebula.ResultSet{}, nil)
		patches.ApplyMethodReturn(nSession, "Release")

		err := session.Open()
		Expect(err).NotTo(HaveOccurred())
		resp, err := session.Execute("")
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).NotTo(BeNil())

		err = session.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	It("failed", func() {
		session := newSessionV5(HostAddress{}, "user", "password", logger.NopLogger)
		connection := nebula.NewConnection(nebula.HostAddress{})
		nSession := &nebula.Session{}
		id := int64(1)
		authResp := &graph.XAuthResponse_{
			Identifier: &id,
		}

		patches := gomonkey.NewPatches()
		defer patches.Reset()

		patches.ApplyFuncReturn(nebula.NewConnection, connection)
		patches.ApplyMethodReturn(connection, "Open", stderrors.New("open failed"))

		err := session.Open()
		Expect(err).To(HaveOccurred())

		patches.Reset()

		patches.ApplyFuncReturn(nebula.NewConnection, connection)
		patches.ApplyMethodReturn(connection, "Open", nil)
		patches.ApplyMethodReturn(connection, "Authenticate", authResp, stderrors.New("authenticate failed"))

		err = session.Open()
		Expect(err).To(HaveOccurred())

		patches.Reset()

		patches.ApplyFuncReturn(nebula.NewConnection, connection)
		patches.ApplyMethodReturn(connection, "Open", nil)
		patches.ApplyMethodReturn(connection, "Authenticate", authResp, nil)

		patches.ApplyFuncReturn(nebula.NewSession, nSession)

		patches.ApplyMethodReturn(nSession, "Execute", nil, stderrors.New("execute failed"))
		patches.ApplyMethodReturn(nSession, "Release")

		err = session.Open()
		Expect(err).NotTo(HaveOccurred())
		resp, err := session.Execute("")
		Expect(err).To(HaveOccurred())
		Expect(resp).To(BeNil())

		err = session.Close()
		Expect(err).NotTo(HaveOccurred())
	})
})
