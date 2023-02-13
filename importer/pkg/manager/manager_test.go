package manager

import (
	stderrors "errors"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/importer"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/source"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/spec"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manager", func() {
	It("New", func() {
		m := New(client.NewPool())
		m1, ok := m.(*defaultManager)
		Expect(ok).To(BeTrue())
		Expect(m1).NotTo(BeNil())
		Expect(m1.pool).NotTo(BeNil())
		Expect(m1.getClientOptions).To(BeNil())
		Expect(m1.batch).To(Equal(0))
		Expect(m1.readerConcurrency).To(Equal(DefaultReaderConcurrency))
		Expect(m1.readerPool).NotTo(BeNil())
		Expect(m1.importerConcurrency).To(Equal(DefaultImporterConcurrency))
		Expect(m1.importerPool).NotTo(BeNil())
		Expect(m1.statsInterval).To(Equal(DefaultStatsInterval))
		Expect(m1.hooks.Before).To(BeEmpty())
		Expect(m1.hooks.After).To(BeEmpty())
		Expect(m1.logger).NotTo(BeNil())
	})

	It("NewWithOpts", func() {
		m := NewWithOpts(
			WithGraphName("graphName"),
			WithClientPool(client.NewPool()),
			WithGetClientOptions(client.WithClientInitFunc(nil)),
			WithBatch(1),
			WithReaderConcurrency(DefaultReaderConcurrency+1),
			WithImporterConcurrency(DefaultImporterConcurrency+1),
			WithStatsInterval(DefaultStatsInterval+1),
			WithBeforeHooks(&Hook{
				Statements: []string{"before statements1"},
				Wait:       time.Second,
			}),
			WithAfterHooks(&Hook{
				Statements: []string{"after statements"},
				Wait:       time.Second,
			}),
			WithLogger(logger.NopLogger),
		)
		m1, ok := m.(*defaultManager)
		Expect(ok).To(BeTrue())
		Expect(m1).NotTo(BeNil())
		Expect(m1.pool).NotTo(BeNil())
		Expect(m1.getClientOptions).NotTo(BeNil())
		Expect(m1.batch).To(Equal(1))
		Expect(m1.readerConcurrency).To(Equal(DefaultReaderConcurrency + 1))
		Expect(m1.readerPool).NotTo(BeNil())
		Expect(m1.importerConcurrency).To(Equal(DefaultImporterConcurrency + 1))
		Expect(m1.importerPool).NotTo(BeNil())
		Expect(m1.statsInterval).To(Equal(DefaultStatsInterval + 1))
		Expect(m1.hooks.Before).To(HaveLen(1))
		Expect(m1.hooks.After).To(HaveLen(1))
		Expect(m1.logger).NotTo(BeNil())
	})

	Describe("Run", func() {
		var (
			tmpdir          string
			nodeFile        string
			nodeSize        int64
			edgeFile        string
			edgeSize        int64
			patches         *gomonkey.Patches
			ctrl            *gomock.Controller
			mockSource      *source.MockSource
			mockClient      *client.MockClient
			mockClientPool  *client.MockPool
			mockResponse    *client.MockResponse
			mockImporter    *importer.MockImporter
			m               Manager
			batch           = 10
			nodeRecordCount = 1005
			edgeRecordCount = 2006
		)
		BeforeEach(func() {
			var err error
			tmpdir, err = os.MkdirTemp("", "test")
			Expect(err).NotTo(HaveOccurred())

			nodeFile = filepath.Join(tmpdir, "node1.csv")
			edgeFile = filepath.Join(tmpdir, "edge1.csv")

			patches = gomonkey.NewPatches()
			ctrl = gomock.NewController(GinkgoT())
			mockSource = source.NewMockSource(ctrl)
			mockClient = client.NewMockClient(ctrl)
			mockClientPool = client.NewMockPool(ctrl)
			mockResponse = client.NewMockResponse(ctrl)
			mockImporter = importer.NewMockImporter(ctrl)

			l, err := logger.New(logger.WithLevel(logger.WarnLevel))
			Expect(err).NotTo(HaveOccurred())
			m = New(
				mockClientPool,
				WithBatch(batch),
				WithLogger(l),
				WithBeforeHooks(&Hook{
					Statements: []string{"before statement"},
					Wait:       time.Second,
				}),
				WithAfterHooks(&Hook{
					Statements: []string{"after statement"},
					Wait:       time.Second,
				}),
			)

			fNode, err := os.Create(nodeFile)
			Expect(err).NotTo(HaveOccurred())
			defer fNode.Close()
			for i := 0; i < nodeRecordCount; i++ {
				fNode.WriteString(fmt.Sprintf("%d,np%d\n", i, i))
			}
			fiNode, err := fNode.Stat()
			Expect(err).NotTo(HaveOccurred())
			nodeSize = fiNode.Size()

			fEdge, err := os.Create(edgeFile)
			Expect(err).NotTo(HaveOccurred())
			defer fEdge.Close()
			for i := 0; i < edgeRecordCount; i++ {
				fEdge.WriteString(fmt.Sprintf("%d,%d,ep%d\n", i, i, i))
			}
			fiNode, err = fEdge.Stat()
			Expect(err).NotTo(HaveOccurred())
			edgeSize = fiNode.Size()
		})

		AfterEach(func() {
			ctrl.Finish()
			patches.Reset()
			err := os.RemoveAll(tmpdir)
			Expect(err).NotTo(HaveOccurred())
		})

		It("concurrency successfully", func() {
			var err error
			loopCountPreFile := 10
			for i := 0; i < loopCountPreFile; i++ {
				err = m.Import(
					&source.Config{Path: nodeFile},
					mockImporter,
					mockImporter,
				)
				Expect(err).NotTo(HaveOccurred())
				err = m.Import(
					&source.Config{Path: edgeFile},
					mockImporter,
					mockImporter,
				)
				Expect(err).NotTo(HaveOccurred())
			}

			gomock.InOrder(
				mockClientPool.EXPECT().GetClient(gomock.Any()).Return(mockClient, nil),
				mockClient.EXPECT().Execute("before statement").Times(1).Return(mockResponse, nil),
				mockResponse.EXPECT().IsSucceed().Return(true),

				mockClientPool.EXPECT().Open().Return(nil),

				mockClientPool.EXPECT().GetClient(gomock.Any()).Return(mockClient, nil),
				mockClient.EXPECT().Execute("after statement").Times(1).Return(mockResponse, nil),
				mockResponse.EXPECT().IsSucceed().Return(true),
			)

			var executeFailedTimes int64 = 10
			var currExecuteTimes int64
			fnImport := func(records ...spec.Record) (*importer.ImportResp, error) {
				curr := atomic.AddInt64(&currExecuteTimes, 1)
				if curr%100 == 0 && curr/100 <= executeFailedTimes {
					return nil, stderrors.New("import failed")
				}
				return &importer.ImportResp{
					Latency:  2 * time.Microsecond,
					RespTime: 3 * time.Microsecond,
				}, nil
			}
			mockImporter.EXPECT().Wait().Times(loopCountPreFile * 4)
			mockImporter.EXPECT().Import(gomock.Any()).AnyTimes().DoAndReturn(fnImport)
			mockImporter.EXPECT().Done().Times(loopCountPreFile * 4)

			err = m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Wait()
			Expect(err).NotTo(HaveOccurred())
			s := m.Stats()

			getBatches := func(recordCount int) int {
				batches := recordCount / batch
				if recordCount%batch != 0 {
					batches++
				}
				return batches
			}
			totalBatches := (getBatches(nodeRecordCount) + getBatches(edgeRecordCount)) * loopCountPreFile

			Expect(s.StartTime.IsZero()).To(BeFalse())
			Expect(s.ProcessedBytes).To(Equal((nodeSize + edgeSize) * int64(loopCountPreFile)))
			Expect(s.TotalBytes).To(Equal((nodeSize + edgeSize) * int64(loopCountPreFile)))
			Expect(s.FailedRecords).NotTo(Equal(int64(0)))
			Expect(s.FailedRecords).To(BeNumerically("<=", executeFailedTimes*int64(batch)))
			Expect(s.TotalRecords).To(Equal(int64((nodeRecordCount + edgeRecordCount) * loopCountPreFile)))
			Expect(s.FailedRequest).To(Equal(executeFailedTimes))
			Expect(s.TotalRequest).To(Equal(int64(totalBatches * 2)))
			Expect(s.TotalLatency).To(Equal(2 * time.Microsecond * time.Duration(int64(totalBatches*2)-executeFailedTimes)))
			Expect(s.TotalRespTime).To(Equal(3 * time.Microsecond * time.Duration(int64(totalBatches*2)-executeFailedTimes)))
			Expect(s.FailedProcessed).NotTo(Equal(int64(0)))
			Expect(s.FailedRecords).To(BeNumerically("<=", executeFailedTimes*int64(batch)))
			Expect(s.TotalProcessed).To(Equal(int64((nodeRecordCount + edgeRecordCount) * loopCountPreFile * 2)))
		})

		It("Import without importer", func() {
			err := m.Import(&source.Config{Path: nodeFile + "not-exists"})
			Expect(err).NotTo(HaveOccurred())
		})

		It("Import file not exists", func() {
			err := m.Import(
				&source.Config{Path: nodeFile + "not-exists"},
				mockImporter,
			)
			Expect(err).To(HaveOccurred())
			Expect(stderrors.Is(err, os.ErrNotExist)).To(BeTrue())
		})

		It("get client failed", func() {
			mockClientPool.EXPECT().GetClient(gomock.Any()).Return(nil, stderrors.New("test error"))

			err := m.Start()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("test error"))
		})

		It("exec before failed", func() {
			gomock.InOrder(
				mockClientPool.EXPECT().GetClient(gomock.Any()).Return(mockClient, nil),
				mockClient.EXPECT().Execute("before statement").Times(1).Return(nil, stderrors.New("test error")),
			)

			err := m.Start()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("test error"))
		})

		It("client pool open failed", func() {
			gomock.InOrder(
				mockClientPool.EXPECT().GetClient(gomock.Any()).Return(mockClient, nil),
				mockClient.EXPECT().Execute("before statement").Times(1).Return(mockResponse, nil),
				mockResponse.EXPECT().IsSucceed().Return(true),

				mockClientPool.EXPECT().Open().Return(stderrors.New("test error")),
			)

			err := m.Start()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("test error"))
		})

		It("exec after failed", func() {
			gomock.InOrder(
				mockClientPool.EXPECT().GetClient(gomock.Any()).Return(mockClient, nil),
				mockClient.EXPECT().Execute("before statement").Times(1).Return(mockResponse, nil),
				mockResponse.EXPECT().IsSucceed().Return(true),

				mockClientPool.EXPECT().Open().Return(nil),

				mockClientPool.EXPECT().GetClient(gomock.Any()).Return(mockClient, nil),
				mockClient.EXPECT().Execute("after statement").Times(1).Return(nil, stderrors.New("test error")),
			)

			err := m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Wait()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("test error"))
		})

		It("stop successfully", func() {
			gomock.InOrder(
				mockClientPool.EXPECT().GetClient(gomock.Any()).Return(mockClient, nil),
				mockClient.EXPECT().Execute("before statement").Times(1).Return(mockResponse, nil),
				mockResponse.EXPECT().IsSucceed().Return(true),

				mockClientPool.EXPECT().Open().Return(nil),

				mockClientPool.EXPECT().GetClient(gomock.Any()).Return(mockClient, nil),
				mockClient.EXPECT().Execute("after statement").Times(1).Return(mockResponse, nil),
				mockResponse.EXPECT().IsSucceed().Return(true),
			)

			err := m.Start()
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			err = m.Stop()
			Expect(err).NotTo(HaveOccurred())
		})

		It("stop failed", func() {
			gomock.InOrder(
				mockClientPool.EXPECT().GetClient(gomock.Any()).Return(mockClient, nil),
				mockClient.EXPECT().Execute("before statement").Times(1).Return(mockResponse, nil),
				mockResponse.EXPECT().IsSucceed().Return(true),

				mockClientPool.EXPECT().Open().Return(nil),

				mockClientPool.EXPECT().GetClient(gomock.Any()).Return(mockClient, nil),
				mockClient.EXPECT().Execute("after statement").Times(1).Return(mockResponse, nil),
				mockResponse.EXPECT().IsSucceed().Return(false),
				mockResponse.EXPECT().GetError().Times(1).Return(stderrors.New("exec failed")),
			)

			err := m.Start()
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			err = m.Stop()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("exec failed"))
		})

		It("stop without read finished", func() {
			patches.ApplyGlobalVar(&sourceOpen, func(_ *source.Config) (source.Source, error) {
				return mockSource, nil
			})

			gomock.InOrder(
				mockClientPool.EXPECT().GetClient(gomock.Any()).Return(mockClient, nil),
				mockClient.EXPECT().Execute("before statement").Times(1).Return(mockResponse, nil),
				mockResponse.EXPECT().IsSucceed().Return(true),

				mockClientPool.EXPECT().Open().Return(nil),

				mockClientPool.EXPECT().GetClient(gomock.Any()).Return(mockClient, nil),
				mockClient.EXPECT().Execute("after statement").Times(1).Return(mockResponse, nil),
				mockResponse.EXPECT().IsSucceed().Return(true),
			)

			mockSource.EXPECT().Size().Return(int64(1024*1024*1024*1024), nil)
			mockSource.EXPECT().Config().Return(&source.Config{})
			mockSource.EXPECT().Read(gomock.Any()).AnyTimes().DoAndReturn(func(p []byte) (int, error) {
				n := copy(p, "1,np1\n")
				return n, nil
			})
			mockSource.EXPECT().Close().Return(nil)

			mockImporter.EXPECT().Wait().Times(1)
			mockImporter.EXPECT().Import(gomock.Any()).AnyTimes().Return(&importer.ImportResp{}, nil)
			mockImporter.EXPECT().Done().Times(1)

			err := m.Import(
				&source.Config{Path: nodeFile},
				mockImporter,
			)
			Expect(err).NotTo(HaveOccurred())

			err = m.Start()
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			err = m.Stop()
			Expect(err).NotTo(HaveOccurred())
		})

		It("no hooks", func() {
			m.(*defaultManager).hooks.Before = nil
			m.(*defaultManager).hooks.After = nil

			mockClientPool.EXPECT().Open().Return(nil)

			err := m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Wait()
			Expect(err).NotTo(HaveOccurred())
		})

		It("nil or empty hooks", func() {
			m.(*defaultManager).hooks.Before = []*Hook{
				nil,
				{Statements: []string{""}},
			}
			m.(*defaultManager).hooks.After = []*Hook{
				{Statements: []string{""}},
				nil,
			}

			mockClientPool.EXPECT().Open().Return(nil)

			err := m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Wait()
			Expect(err).NotTo(HaveOccurred())
		})

		It("disable stats interval", func() {
			m.(*defaultManager).hooks.Before = nil
			m.(*defaultManager).hooks.After = nil
			m.(*defaultManager).statsInterval = 0

			mockClientPool.EXPECT().Open().Return(nil)

			err := m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Wait()
			Expect(err).NotTo(HaveOccurred())
		})

		It("stats interval print", func() {
			m.(*defaultManager).hooks.Before = nil
			m.(*defaultManager).hooks.After = nil
			m.(*defaultManager).statsInterval = 10 * time.Microsecond

			mockClientPool.EXPECT().Open().Return(nil)

			err := m.Start()
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			err = m.Wait()
			Expect(err).NotTo(HaveOccurred())
		})

		It("submit reader failed", func() {
			m.(*defaultManager).hooks.Before = nil
			m.(*defaultManager).hooks.After = nil
			m.(*defaultManager).readerPool.Release()

			mockClientPool.EXPECT().Open().Return(nil)

			mockImporter.EXPECT().Done().Times(2)

			err := m.Import(
				&source.Config{Path: nodeFile},
				mockImporter,
			)
			Expect(err).NotTo(HaveOccurred())
			err = m.Import(
				&source.Config{Path: nodeFile},
				mockImporter,
			)
			Expect(err).NotTo(HaveOccurred())

			err = m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Wait()
			Expect(err).NotTo(HaveOccurred())
		})

		It("submit importer failed", func() {
			m.(*defaultManager).hooks.Before = nil
			m.(*defaultManager).hooks.After = nil
			m.(*defaultManager).importerPool.Release()

			mockClientPool.EXPECT().Open().Return(nil)

			mockImporter.EXPECT().Wait().Times(2)
			mockImporter.EXPECT().Done().Times(2)

			err := m.Import(
				&source.Config{Path: nodeFile},
				mockImporter,
			)
			Expect(err).NotTo(HaveOccurred())
			err = m.Import(
				&source.Config{Path: nodeFile},
				mockImporter,
			)
			Expect(err).NotTo(HaveOccurred())

			err = m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Wait()
			Expect(err).NotTo(HaveOccurred())
		})

		It("get size failed", func() {
			m.(*defaultManager).hooks.Before = nil
			m.(*defaultManager).hooks.After = nil
			patches.ApplyGlobalVar(&sourceOpen, func(_ *source.Config) (source.Source, error) {
				return mockSource, nil
			})

			mockSource.EXPECT().Size().Times(2).Return(int64(0), stderrors.New("test error"))
			mockSource.EXPECT().Close().Times(2).Return(nil)

			err := m.Import(
				&source.Config{Path: nodeFile},
				mockImporter,
			)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("test error"))

			err = m.Import(
				&source.Config{Path: nodeFile}, mockImporter,
			)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("test error"))
		})

		It("read failed", func() {
			m.(*defaultManager).hooks.Before = nil
			m.(*defaultManager).hooks.After = nil
			patches.ApplyGlobalVar(&sourceOpen, func(_ *source.Config) (source.Source, error) {
				return mockSource, nil
			})

			mockClientPool.EXPECT().Open().Return(nil)
			mockSource.EXPECT().Size().Times(2).Return(int64(1024), nil)
			mockSource.EXPECT().Config().Times(2).Return(nil)
			mockSource.EXPECT().Read(gomock.Any()).Times(2).Return(0, stderrors.New("test error"))
			mockSource.EXPECT().Close().Times(2).Return(nil)

			mockImporter.EXPECT().Wait().Times(2)
			mockImporter.EXPECT().Done().Times(2)

			err := m.Import(
				&source.Config{Path: nodeFile},
				mockImporter,
			)
			Expect(err).NotTo(HaveOccurred())
			err = m.Import(
				&source.Config{Path: nodeFile},
				mockImporter,
			)
			Expect(err).NotTo(HaveOccurred())

			err = m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Wait()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
