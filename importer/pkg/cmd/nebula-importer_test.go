//go:build linux

package cmd

import (
	stderrors "errors"
	"os"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/manager"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ImporterCommand", func() {
	var (
		patches        *gomonkey.Patches
		ctrl           *gomock.Controller
		mockClient     *client.MockClient
		mockClientPool *client.MockPool
		mockResultSet  *client.MockResultSet
		mockManager    *manager.MockManager
	)
	BeforeEach(func() {
		patches = gomonkey.NewPatches()
		ctrl = gomock.NewController(GinkgoT())
		mockClient = client.NewMockClient(ctrl)
		mockClientPool = client.NewMockPool(ctrl)
		mockResultSet = client.NewMockResultSet(ctrl)
		mockManager = manager.NewMockManager(ctrl)
	})
	AfterEach(func() {
		ctrl.Finish()
		patches.Reset()
	})

	It("successfully", func() {
		patches.ApplyFuncReturn(client.NewPool, mockClientPool)

		mockClientPool.EXPECT().GetClient(gomock.Any()).AnyTimes().Return(mockClient, nil)
		mockClientPool.EXPECT().Open().AnyTimes().Return(nil)
		mockClientPool.EXPECT().Execute(gomock.Any()).AnyTimes().Return(mockResultSet, nil)
		mockClientPool.EXPECT().Close().AnyTimes().Return(nil)

		mockClient.EXPECT().Open().AnyTimes().Return(nil)
		mockClient.EXPECT().Execute(gomock.Any()).AnyTimes().Return(mockResultSet, nil)
		mockClient.EXPECT().Close().AnyTimes().Return(nil)

		mockResultSet.EXPECT().IsSucceed().AnyTimes().Return(true)
		mockResultSet.EXPECT().GetLatency().AnyTimes().Return(int64(2))

		for _, f := range []string{
			"testdata/nebula-importer.v3.yaml",
			"testdata/nebula-importer.v5.yaml",
		} {
			command := NewDefaultImporterCommand()
			command.SetArgs([]string{"-c", f})
			err := command.Execute()
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("parse file failed", func() {
		command := NewDefaultImporterCommand()
		command.SetArgs([]string{"-c", "testdata/not-exists/nebula-importer.yaml"})
		err := command.Execute()
		Expect(err).To(HaveOccurred())
	})

	It("optimize failed", func() {
		command := NewDefaultImporterCommand()
		command.SetArgs([]string{"-c", "testdata/optimize-failed.yaml"})
		err := command.Execute()
		Expect(err).To(HaveOccurred())
	})

	It("build failed", func() {
		command := NewDefaultImporterCommand()
		command.SetArgs([]string{"-c", "testdata/build-failed.yaml"})
		err := command.Execute()
		Expect(err).To(HaveOccurred())
	})

	It("complete failed", func() {
		o := NewImporterOptions(common.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		})

		patches.ApplyMethodReturn(o, "Complete", stderrors.New("test error"))

		o.useNopLogger = true
		command := NewImporterCommand(o)
		command.SetArgs([]string{"-c", "testdata/nebula-importer.v5.yaml"})

		err := command.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(stderrors.New("test error")))
	})

	It("manager start failed", func() {
		patches.ApplyFuncReturn(manager.NewWithOpts, mockManager)
		mockManager.EXPECT().Import(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
		mockManager.EXPECT().Start().Return(stderrors.New("test error"))

		o := NewImporterOptions(common.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		})

		o.useNopLogger = true
		command := NewImporterCommand(o)
		command.SetArgs([]string{"-c", "testdata/nebula-importer.v5.yaml"})

		err := command.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(stderrors.New("test error")))
	})

	It("manager wait failed", func() {
		patches.ApplyFuncReturn(manager.NewWithOpts, mockManager)
		mockManager.EXPECT().Import(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
		mockManager.EXPECT().Start().Return(nil)
		mockManager.EXPECT().Wait().Return(stderrors.New("test error"))

		o := NewImporterOptions(common.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		})

		o.useNopLogger = true
		command := NewImporterCommand(o)
		command.SetArgs([]string{"-c", "testdata/nebula-importer.v5.yaml"})

		err := command.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(stderrors.New("test error")))
	})
})
