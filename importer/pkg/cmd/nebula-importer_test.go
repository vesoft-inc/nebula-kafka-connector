//go:build linux

package cmd

import (
	stderrors "errors"
	"os"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/manager"
)

var _ = Describe("", func() {
	var (
		patches       *gomonkey.Patches
		ctrl          *gomock.Controller
		mockClient    *client.MockClient
		mockResultSet *client.MockResultSet
		mockManager   *manager.MockManager
	)
	BeforeEach(func() {
		patches = gomonkey.NewPatches()
		ctrl = gomock.NewController(GinkgoT())
		mockClient = client.NewMockClient(ctrl)
		mockResultSet = client.NewMockResultSet(ctrl)
		mockManager = manager.NewMockManager(ctrl)
	})
	AfterEach(func() {
		ctrl.Finish()
		patches.Reset()
	})

	It("success", func() {
		patches.ApplyFuncReturn(client.New, mockClient)

		mockClient.EXPECT().Open().AnyTimes().Return(nil)
		mockClient.EXPECT().Execute(gomock.Any()).AnyTimes().Return(mockResultSet, nil)
		mockClient.EXPECT().Close().AnyTimes().Return(nil)

		mockResultSet.EXPECT().IsSucceed().AnyTimes().Return(true)
		mockResultSet.EXPECT().GetLatency().AnyTimes().Return(int64(2))

		command := NewDefaultImporterCommand()
		command.SetArgs([]string{"-c", "testdata/nebula-importer.yaml"})
		err := command.Execute()
		Expect(err).NotTo(HaveOccurred())
	})

	It("parse file failed", func() {
		command := NewDefaultImporterCommand()
		command.SetArgs([]string{"-c", "testdata/not-exists/nebula-importer.yaml"})
		err := command.Execute()
		Expect(err).To(HaveOccurred())
	})

	It("build logger failed", func() {
		command := NewDefaultImporterCommand()
		command.SetArgs([]string{"-c", "testdata/build-logger-failed.yaml"})
		err := command.Execute()
		Expect(err).To(HaveOccurred())
	})

	It("build client failed", func() {
		o := NewImporterOptions(common.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		})
		o.useNopLogger = true
		command := NewImporterCommand(o)
		command.SetArgs([]string{"-c", "testdata/nebula-importer.yaml"})
		err := command.Execute()
		Expect(err).To(HaveOccurred())
	})

	It("build graph failed", func() {
		patches.ApplyFuncReturn(client.New, mockClient)

		mockClient.EXPECT().Open().AnyTimes().Return(nil)
		mockClient.EXPECT().Close().AnyTimes().Return(nil)

		o := NewImporterOptions(common.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		})
		o.useNopLogger = true
		command := NewImporterCommand(o)
		command.SetArgs([]string{"-c", "testdata/build-graph-failed.yaml"})
		err := command.Execute()
		Expect(err).To(HaveOccurred())
	})

	It("build manager failed", func() {
		patches.ApplyFuncReturn(client.New, mockClient)

		mockClient.EXPECT().Open().AnyTimes().Return(nil)
		mockClient.EXPECT().Close().AnyTimes().Return(nil)

		o := NewImporterOptions(common.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		})
		o.useNopLogger = true
		command := NewImporterCommand(o)
		command.SetArgs([]string{"-c", "testdata/build-manager-failed.yaml"})

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
		command.SetArgs([]string{"-c", "testdata/build-manager-failed.yaml"})

		err := command.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(stderrors.New("test error")))
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
		command.SetArgs([]string{"-c", "testdata/build-manager-failed.yaml"})

		err := command.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(stderrors.New("test error")))
	})

	It("manager start failed", func() {
		patches.ApplyFuncReturn(client.New, mockClient)

		mockClient.EXPECT().Open().AnyTimes().Return(nil)
		mockClient.EXPECT().Execute(gomock.Any()).AnyTimes().Return(mockResultSet, nil)
		mockClient.EXPECT().Close().AnyTimes().Return(nil)

		patches.ApplyFuncReturn(manager.NewWithOpts, mockManager)
		mockManager.EXPECT().ImportNode(gomock.Any(), gomock.Any()).Return(nil)
		mockManager.EXPECT().ImportEdge(gomock.Any(), gomock.Any()).Return(nil)
		mockManager.EXPECT().Start().Return(stderrors.New("test error"))

		o := NewImporterOptions(common.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		})

		o.useNopLogger = true
		command := NewImporterCommand(o)
		command.SetArgs([]string{"-c", "testdata/nebula-importer.yaml"})

		err := command.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(stderrors.New("test error")))
	})

	It("manager wait failed", func() {
		patches.ApplyFuncReturn(client.New, mockClient)

		mockClient.EXPECT().Open().AnyTimes().Return(nil)
		mockClient.EXPECT().Execute(gomock.Any()).AnyTimes().Return(mockResultSet, nil)
		mockClient.EXPECT().Close().AnyTimes().Return(nil)

		patches.ApplyFuncReturn(manager.NewWithOpts, mockManager)
		mockManager.EXPECT().ImportNode(gomock.Any(), gomock.Any()).Return(nil)
		mockManager.EXPECT().ImportEdge(gomock.Any(), gomock.Any()).Return(nil)
		mockManager.EXPECT().Start().Return(nil)
		mockManager.EXPECT().Wait().Return(stderrors.New("test error"))

		o := NewImporterOptions(common.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		})

		o.useNopLogger = true
		command := NewImporterCommand(o)
		command.SetArgs([]string{"-c", "testdata/nebula-importer.yaml"})

		err := command.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(stderrors.New("test error")))
	})
})
