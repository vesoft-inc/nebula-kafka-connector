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
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/importer"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/source"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/spec"
)

var _ = Describe("Manager", func() {
	It("New", func() {
		m := New(client.New())
		m1, ok := m.(*defaultManager)
		Expect(ok).To(BeTrue())
		Expect(m1).NotTo(BeNil())
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
			graph           *spec.Graph
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
			graph = spec.NewGraph(
				"graphName",
				spec.WithGraphNodes(
					spec.NewNode(
						"node1",
						spec.WithNodeID(&spec.NodeID{
							Name:  "id",
							Type:  spec.ValueTypeInt,
							Index: 0,
						}),
						spec.WithNodeProps(&spec.Prop{
							Name:  "nodeProp1",
							Type:  spec.ValueTypeString,
							Index: 1,
						}),
					),
					spec.NewNode(
						"node2",
						spec.WithNodeID(&spec.NodeID{
							Name:  "id",
							Type:  spec.ValueTypeInt,
							Index: 0,
						}),
					),
				),
				spec.WithGraphEdges(
					spec.NewEdge(
						"edge1",
						spec.WithEdgeSrc(&spec.EdgeNodeRef{
							Name: "node1",
							ID: &spec.NodeID{
								Name:  "id",
								Type:  spec.ValueTypeInt,
								Index: 0,
							},
						}),
						spec.WithEdgeDst(&spec.EdgeNodeRef{
							Name: "node1",
							ID: &spec.NodeID{
								Name:  "id",
								Type:  spec.ValueTypeInt,
								Index: 1,
							},
						}),
						spec.WithEdgeProps(&spec.Prop{
							Name:  "edgeProp1",
							Type:  spec.ValueTypeString,
							Index: 2,
						}),
					),
					spec.NewEdge(
						"edge1",
						spec.WithEdgeSrc(&spec.EdgeNodeRef{
							Name: "node2",
							ID: &spec.NodeID{
								Name:  "id",
								Type:  spec.ValueTypeInt,
								Index: 0,
							},
						}),
						spec.WithEdgeDst(&spec.EdgeNodeRef{
							Name: "node2",
							ID: &spec.NodeID{
								Name:  "id",
								Type:  spec.ValueTypeInt,
								Index: 1,
							},
						}),
					),
				),
			)
			graph.Complete()
			Expect(graph.Validate()).NotTo(HaveOccurred())
			l, err := logger.New(logger.WithLevel(logger.WarnLevel))
			Expect(err).NotTo(HaveOccurred())
			m = New(
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
				node1, _ := graph.GetNodeByName("node1")
				node2, _ := graph.GetNodeByName("node1")
				err = m.Import(
					&source.Config{Path: nodeFile},
					importer.NewNodeImporter(graph, node1, mockClient),
					importer.NewNodeImporter(graph, node2, mockClient),
				)
				Expect(err).NotTo(HaveOccurred())
				edge1, _ := graph.GetEdgeByName("edge1")
				edge2, _ := graph.GetEdgeByName("edge1")
				err = m.Import(
					&source.Config{Path: edgeFile},
					importer.NewEdgeImporter(graph, edge1, mockClient),
					importer.NewEdgeImporter(graph, edge2, mockClient),
				)
				Expect(err).NotTo(HaveOccurred())
			}

			mockClient.EXPECT().Open().AnyTimes().Return(nil)
			mockClient.EXPECT().Close().AnyTimes().Return(nil)

			var executeFailedTimes int64 = 10
			var currExecuteTimes int64
			fnExecute := func(statement string) (client.ResultSet, error) {
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
			Expect(s.FailedRecords).NotTo(Equal(int64(0)))
			Expect(s.FailedRecords).To(BeNumerically("<=", executeFailedTimes*int64(batch)))
			Expect(s.TotalRecords).To(Equal(int64((nodeRecordCount + edgeRecordCount) * loopCountPreFile)))
			Expect(s.FailedRequest).To(Equal(executeFailedTimes))
			Expect(s.TotalRequest).To(Equal(int64(totalBatches * 2)))
			Expect(s.TotalLatency).To(Equal(2 * time.Microsecond * time.Duration(int64(totalBatches*2)-executeFailedTimes)))
			Expect(s.FailedProcessed).NotTo(Equal(int64(0)))
			Expect(s.FailedRecords).To(BeNumerically("<=", executeFailedTimes*int64(batch)))
			Expect(s.TotalProcessed).To(Equal(int64((nodeRecordCount + edgeRecordCount) * loopCountPreFile * 2)))
		})

		It("Import without importer", func() {
			err := m.Import(&source.Config{Path: nodeFile + "not-exists"})
			Expect(err).NotTo(HaveOccurred())
		})

		It("Import file not exists", func() {
			node1, _ := graph.GetNodeByName("node1")
			err := m.Import(
				&source.Config{Path: nodeFile + "not-exists"},
				importer.NewNodeImporter(graph, node1, mockClient),
			)
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
			mockResultSet.EXPECT().GetError().AnyTimes().Return(stderrors.New("exec failed"))

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

			node1, _ := graph.GetNodeByName("node1")
			err := m.Import(
				&source.Config{Path: nodeFile},
				importer.NewNodeImporter(graph, node1, mockClient),
			)
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

			node1, _ := graph.GetNodeByName("node1")
			err := m.Import(
				&source.Config{Path: nodeFile},
				importer.NewNodeImporter(graph, node1, mockClient),
			)
			Expect(err).NotTo(HaveOccurred())
			edge1, _ := graph.GetEdgeByName("edge1")
			err = m.Import(
				&source.Config{Path: nodeFile},
				importer.NewEdgeImporter(graph, edge1, mockClient),
			)
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

			node1, _ := graph.GetNodeByName("node1")
			err := m.Import(
				&source.Config{Path: nodeFile},
				importer.NewNodeImporter(graph, node1, mockClient),
			)
			Expect(err).NotTo(HaveOccurred())
			edge1, _ := graph.GetEdgeByName("edge1")
			err = m.Import(
				&source.Config{Path: nodeFile},
				importer.NewEdgeImporter(graph, edge1, mockClient),
			)
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

			node1, _ := graph.GetNodeByName("node1")
			err := m.Import(
				&source.Config{Path: nodeFile},
				importer.NewNodeImporter(graph, node1, mockClient),
			)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("test error"))
			edge1, _ := graph.GetEdgeByName("edge1")
			err = m.Import(
				&source.Config{Path: nodeFile},
				importer.NewEdgeImporter(graph, edge1, mockClient),
			)
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

			node1, _ := graph.GetNodeByName("node1")
			err := m.Import(
				&source.Config{Path: nodeFile},
				importer.NewNodeImporter(graph, node1, mockClient),
			)
			Expect(err).NotTo(HaveOccurred())
			edge1, _ := graph.GetEdgeByName("edge1")
			err = m.Import(
				&source.Config{Path: nodeFile},
				importer.NewEdgeImporter(graph, edge1, mockClient),
			)
			Expect(err).NotTo(HaveOccurred())

			err = m.Start()
			Expect(err).NotTo(HaveOccurred())

			err = m.Wait()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
