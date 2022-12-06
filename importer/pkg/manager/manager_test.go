package manager

import (
	stderrors "errors"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/source"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/spec"
)

var _ = Describe("WaitGroupMap", func() {
	It("New", func() {
		m := New(spec.NewGraph(""), client.New())
		m1, ok := m.(*defaultManager)
		Expect(ok).To(BeTrue())
		Expect(m1).NotTo(BeNil())
		Expect(m1.graph).NotTo(BeNil())
		Expect(m1.c).NotTo(BeNil())
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
			WithGraph(spec.NewGraph("")),
			WithClient(client.New()),
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
		Expect(m1.graph).NotTo(BeNil())
		Expect(m1.c).NotTo(BeNil())
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
			mockResultSet   *client.MockResultSet
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
			mockResultSet = client.NewMockResultSet(ctrl)
			graph := spec.NewGraph(
				"graphName",
				spec.WithGraphNodes(
					spec.NewNode(
						"node1",
						spec.WithNodeID(&spec.NodeID{
							Prop: &spec.Prop{
								Name:  "id",
								Type:  spec.ValueTypeInt,
								Index: 0,
							},
						}),
						spec.WithNodeProps(&spec.Prop{
							Name:  "nodeProp1",
							Type:  spec.ValueTypeString,
							Index: 1,
						}),
					),
				),
				spec.WithGraphEdges(
					spec.NewEdge(
						"edge1",
						spec.WithEdgeSrc(&spec.EdgeNodeRef{
							Name: "node1",
							ID: &spec.NodeID{
								Prop: &spec.Prop{
									Name:  "id",
									Type:  spec.ValueTypeInt,
									Index: 0,
								},
							},
						}),
						spec.WithEdgeDst(&spec.EdgeNodeRef{
							Name: "node1",
							ID: &spec.NodeID{
								Prop: &spec.Prop{
									Name:  "id",
									Type:  spec.ValueTypeInt,
									Index: 1,
								},
							},
						}),
						spec.WithEdgeProps(&spec.Prop{
							Name:  "edgeProp1",
							Type:  spec.ValueTypeString,
							Index: 2,
						}),
					),
				),
			)
			l, err := logger.New(logger.WithLevel(logger.WarnLevel))
			Expect(err).NotTo(HaveOccurred())
			m = New(
				graph,
				mockClient,
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

		It("concurrency success", func() {
			var err error
			loopCountPreFile := 10
			for i := 0; i < loopCountPreFile; i++ {
				err = m.ImportNode("node1", &source.Config{Path: nodeFile})
				Expect(err).NotTo(HaveOccurred())
				err = m.ImportEdge("edge1", &source.Config{Path: edgeFile})
				Expect(err).NotTo(HaveOccurred())
			}

			mockClient.EXPECT().Open().AnyTimes().Return(nil)
			mockClient.EXPECT().Close().AnyTimes().Return(nil)

			var executeFailedTimes int64 = 10
			var currExecuteTimes int64
			fnExecute := func(_ string) (client.ResultSet, error) {
				curr := atomic.AddInt64(&currExecuteTimes, 1)
				if curr%100 == 0 && curr/100 <= executeFailedTimes {
					return nil, stderrors.New("execute failed")
				}
				return mockResultSet, nil
			}

			mockClient.EXPECT().Execute(gomock.Any()).AnyTimes().DoAndReturn(fnExecute)

			mockResultSet.EXPECT().IsSucceed().AnyTimes().Return(true)
			mockResultSet.EXPECT().GetLatency().AnyTimes().Return(int64(2))

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
			Expect(s.FailedBatches).To(Equal(executeFailedTimes))
			Expect(s.TotalBatches).To(Equal(int64(totalBatches)))
			Expect(s.FailedRecords).NotTo(Equal(int64(0)))
			Expect(s.FailedRecords).To(BeNumerically("<=", executeFailedTimes*int64(batch)))
			Expect(s.TotalRecords).To(Equal(int64((nodeRecordCount + edgeRecordCount) * loopCountPreFile)))
			Expect(s.TotalLatency).To(Equal(2 * time.Microsecond * time.Duration(int64(totalBatches)-executeFailedTimes)))
		})

		It("ImportNode node not exists", func() {
			err := m.ImportNode("not-exists", &source.Config{Path: nodeFile})
			Expect(err).To(HaveOccurred())
			Expect(stderrors.Is(err, errors.ErrNodeNotFound)).To(BeTrue())
		})

		It("ImportNode file not exists", func() {
			err := m.ImportNode("node1", &source.Config{Path: nodeFile + "not-exists"})
			Expect(err).To(HaveOccurred())
			Expect(stderrors.Is(err, os.ErrNotExist)).To(BeTrue())
		})

		It("ImportEdge edge not exists", func() {
			err := m.ImportEdge("not-exists", &source.Config{Path: edgeFile})
			Expect(err).To(HaveOccurred())
			Expect(stderrors.Is(err, errors.ErrEdgeNotFound)).To(BeTrue())
		})

		It("ImportEdge file not exists", func() {
			err := m.ImportEdge("edge1", &source.Config{Path: nodeFile + "not-exists"})
			Expect(err).To(HaveOccurred())
			Expect(stderrors.Is(err, os.ErrNotExist)).To(BeTrue())
		})

		It("exec before failed", func() {
			mockClient.EXPECT().Execute(gomock.Any()).AnyTimes().Return(nil, stderrors.New("test error"))

			err := m.Start()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("test error"))
		})

		It("exec after failed", func() {
			mockClient.EXPECT().Execute("before statement").AnyTimes().Return(mockResultSet, nil)
			mockClient.EXPECT().Execute("after statement").AnyTimes().Return(nil, stderrors.New("test error"))

			mockResultSet.EXPECT().IsSucceed().AnyTimes().Return(true)
			mockResultSet.EXPECT().GetLatency().AnyTimes().Return(int64(2))

			err := m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Wait()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("test error"))
		})

		It("stop successfully", func() {
			mockClient.EXPECT().Execute("before statement").AnyTimes().Return(mockResultSet, nil)
			mockClient.EXPECT().Execute("after statement").AnyTimes().Return(mockResultSet, nil)

			mockResultSet.EXPECT().IsSucceed().AnyTimes().Return(true)
			mockResultSet.EXPECT().GetLatency().AnyTimes().Return(int64(2))

			err := m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Stop()
			Expect(err).NotTo(HaveOccurred())
		})

		It("stop failed", func() {
			mockClient.EXPECT().Execute("before statement").AnyTimes().Return(mockResultSet, nil)
			mockClient.EXPECT().Execute("after statement").AnyTimes().Return(mockResultSet, nil)

			gomock.InOrder(
				mockResultSet.EXPECT().IsSucceed().Return(true),
				mockResultSet.EXPECT().IsSucceed().Return(false),
			)
			mockResultSet.EXPECT().GetLatency().AnyTimes().Return(int64(2))
			mockResultSet.EXPECT().GetStatus().AnyTimes().Return("exec failed")

			err := m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Stop()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("exec failed"))
		})

		It("stop without read finished", func() {
			patches.ApplyGlobalVar(&sourceOpen, func(_ *source.Config) (source.Source, error) {
				return mockSource, nil
			})

			mockResultSet.EXPECT().IsSucceed().AnyTimes().Return(true)
			mockResultSet.EXPECT().GetLatency().AnyTimes().Return(int64(2))

			mockSource.EXPECT().Size().Return(int64(1024*1024*1024*1024), nil)
			mockSource.EXPECT().Config().Return(nil)
			mockSource.EXPECT().Read(gomock.Any()).AnyTimes().DoAndReturn(func(p []byte) (int, error) {
				n := copy(p, "1,np1\n")
				return n, nil
			})
			mockSource.EXPECT().Close().Return(nil)

			mockClient.EXPECT().Execute(gomock.Any()).AnyTimes().Return(mockResultSet, nil)
			mockClient.EXPECT().Execute(gomock.Any()).AnyTimes().Return(mockResultSet, nil)
			mockResultSet.EXPECT().IsSucceed().AnyTimes().Return(true)
			mockResultSet.EXPECT().GetLatency().AnyTimes().Return(int64(2))

			err := m.ImportNode("node1", &source.Config{Path: nodeFile})
			Expect(err).NotTo(HaveOccurred())

			err = m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Stop()
			Expect(err).NotTo(HaveOccurred())
		})

		It("no hooks", func() {
			m.(*defaultManager).hooks.Before = nil
			m.(*defaultManager).hooks.After = nil
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
			err := m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Wait()
			Expect(err).NotTo(HaveOccurred())
		})

		It("disable stats interval", func() {
			m.(*defaultManager).hooks.Before = nil
			m.(*defaultManager).hooks.After = nil
			m.(*defaultManager).statsInterval = 0
			err := m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Wait()
			Expect(err).NotTo(HaveOccurred())
		})

		It("stats interval print", func() {
			m.(*defaultManager).hooks.Before = nil
			m.(*defaultManager).hooks.After = nil
			m.(*defaultManager).statsInterval = time.Microsecond / 2
			err := m.Start()
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(time.Microsecond)

			err = m.Wait()
			Expect(err).NotTo(HaveOccurred())
		})

		It("submit reader failed", func() {
			m.(*defaultManager).readerPool.Release()

			err := m.ImportNode("node1", &source.Config{Path: nodeFile})
			Expect(err).NotTo(HaveOccurred())
			err = m.ImportEdge("edge1", &source.Config{Path: edgeFile})
			Expect(err).NotTo(HaveOccurred())

			mockClient.EXPECT().Open().AnyTimes().Return(nil)
			mockClient.EXPECT().Close().AnyTimes().Return(nil)

			mockClient.EXPECT().Execute(gomock.Any()).AnyTimes().Return(mockResultSet, nil)

			mockResultSet.EXPECT().IsSucceed().AnyTimes().Return(true)
			mockResultSet.EXPECT().GetLatency().AnyTimes().Return(int64(2))

			err = m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Wait()
			Expect(err).NotTo(HaveOccurred())
		})

		It("submit importer failed", func() {
			m.(*defaultManager).importerPool.Release()

			err := m.ImportNode("node1", &source.Config{Path: nodeFile})
			Expect(err).NotTo(HaveOccurred())
			err = m.ImportEdge("edge1", &source.Config{Path: edgeFile})
			Expect(err).NotTo(HaveOccurred())

			mockClient.EXPECT().Open().AnyTimes().Return(nil)
			mockClient.EXPECT().Close().AnyTimes().Return(nil)

			mockClient.EXPECT().Execute(gomock.Any()).AnyTimes().Return(mockResultSet, nil)

			mockResultSet.EXPECT().IsSucceed().AnyTimes().Return(true)
			mockResultSet.EXPECT().GetLatency().AnyTimes().Return(int64(2))

			err = m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Wait()
			Expect(err).NotTo(HaveOccurred())
		})

		It("get size failed", func() {
			patches.ApplyGlobalVar(&sourceOpen, func(_ *source.Config) (source.Source, error) {
				return mockSource, nil
			})

			mockSource.EXPECT().Size().Times(2).Return(int64(0), stderrors.New("test error"))
			mockSource.EXPECT().Close().Times(2).Return(nil)

			err := m.ImportNode("node1", &source.Config{Path: nodeFile})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("test error"))
			err = m.ImportEdge("edge1", &source.Config{Path: edgeFile})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("test error"))
		})

		It("read failed", func() {
			patches.ApplyGlobalVar(&sourceOpen, func(_ *source.Config) (source.Source, error) {
				return mockSource, nil
			})

			mockClient.EXPECT().Execute("before statement").AnyTimes().Return(mockResultSet, nil)
			mockClient.EXPECT().Execute("after statement").AnyTimes().Return(mockResultSet, nil)
			mockResultSet.EXPECT().IsSucceed().AnyTimes().Return(true)
			mockResultSet.EXPECT().GetLatency().AnyTimes().Return(int64(2))

			mockSource.EXPECT().Size().Times(2).Return(int64(1024), nil)
			mockSource.EXPECT().Config().Times(2).Return(nil)
			mockSource.EXPECT().Read(gomock.Any()).Times(2).Return(0, stderrors.New("test error"))
			mockSource.EXPECT().Close().Times(2).Return(nil)

			err := m.ImportNode("node1", &source.Config{Path: nodeFile})
			Expect(err).NotTo(HaveOccurred())
			err = m.ImportEdge("edge1", &source.Config{Path: edgeFile})
			Expect(err).NotTo(HaveOccurred())

			err = m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Wait()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
