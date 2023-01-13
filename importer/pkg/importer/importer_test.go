package importer

import (
	stderrors "errors"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/spec"
)

var _ = Describe("Importer", func() {
	var (
		ctrl          *gomock.Controller
		mockClient    *client.MockClient
		mockResultSet *client.MockResultSet
		node          *spec.Node
		edge          *spec.Edge
		graph         *spec.Graph
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockClient = client.NewMockClient(ctrl)
		mockResultSet = client.NewMockResultSet(ctrl)
		node = spec.NewNode(
			"nodeName",
			spec.WithNodeID(&spec.NodeID{
				Name: "id",
				Type: spec.ValueTypeInt,
			}),
		)
		edge = spec.NewEdge(
			"edgeName",
			spec.WithEdgeSrc(&spec.EdgeNodeRef{
				Name: "nodeName",
				ID: &spec.NodeID{
					Name:  "id",
					Type:  spec.ValueTypeInt,
					Index: 0,
				},
			}),
			spec.WithEdgeDst(&spec.EdgeNodeRef{
				Name: "nodeName",
				ID: &spec.NodeID{
					Name:  "id",
					Type:  spec.ValueTypeInt,
					Index: 1,
				},
			}),
		)
		graph = spec.NewGraph(
			"graphName",
			spec.WithGraphNodes(node),
			spec.WithGraphEdges(edge),
		)
		graph.Complete()
		Expect(graph.Validate()).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("NewNodeImporter", func() {
		var nodeImporter Importer
		BeforeEach(func() {
			nodeImporter = NewNodeImporter(graph, node, mockClient)
			Expect(nodeImporter.Node()).NotTo(BeNil())
			Expect(nodeImporter.Edge()).To(BeNil())
			Expect(nodeImporter.Graph()).NotTo(BeNil())
		})

		It("empty records", func() {
			resp, err := nodeImporter.Import(spec.Record{})
			Expect(err).To(HaveOccurred())
			Expect(stderrors.Is(err, errors.ErrNoRecord)).To(BeTrue())
			Expect(resp).To(BeNil())
		})

		It("execute failed", func() {
			mockClient.EXPECT().Execute(gomock.Any()).Return(nil, stderrors.New("test error"))
			resp, err := nodeImporter.Import(spec.Record{"id"})
			Expect(err).To(HaveOccurred())
			importError, ok := errors.AsImportError(err)
			Expect(ok).To(BeTrue())
			Expect(importError.Statement()).NotTo(BeEmpty())
			Expect(resp).To(BeNil())
		})

		It("execute IsSucceed false", func() {
			mockClient.EXPECT().Execute(gomock.Any()).Times(1).Return(mockResultSet, nil)
			mockResultSet.EXPECT().IsSucceed().Times(1).Return(false)
			mockResultSet.EXPECT().GetError().Times(1).Return(stderrors.New("status failed"))
			resp, err := nodeImporter.Import(spec.Record{"id"})
			Expect(err).To(HaveOccurred())
			importError, ok := errors.AsImportError(err)
			Expect(ok).To(BeTrue())
			Expect(importError.Messages).To(ContainElement(ContainSubstring("status failed")))
			Expect(importError.Statement()).NotTo(BeEmpty())
			Expect(resp).To(BeNil())
		})

		It("execute success", func() {
			mockClient.EXPECT().Execute(gomock.Any()).Times(1).Return(mockResultSet, nil)
			mockResultSet.EXPECT().IsSucceed().Times(1).Return(true)
			mockResultSet.EXPECT().GetLatency().Times(1).Return(int64(10))
			resp, err := nodeImporter.Import(spec.Record{"id"})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.Latency).To(Equal(time.Microsecond * time.Duration(10)))
		})
	})

	Describe("NewEdgeImporter", func() {
		var nodeImporter Importer
		BeforeEach(func() {
			nodeImporter = NewEdgeImporter(graph, edge, mockClient)
			Expect(nodeImporter.Node()).To(BeNil())
			Expect(nodeImporter.Edge()).NotTo(BeNil())
			Expect(nodeImporter.Graph()).NotTo(BeNil())
		})

		It("empty records", func() {
			resp, err := nodeImporter.Import(spec.Record{})
			Expect(err).To(HaveOccurred())
			Expect(stderrors.Is(err, errors.ErrNoRecord)).To(BeTrue())
			Expect(resp).To(BeNil())
		})

		It("execute failed", func() {
			mockClient.EXPECT().Execute(gomock.Any()).Return(nil, stderrors.New("test error"))
			resp, err := nodeImporter.Import(spec.Record{"id1", "id2"})
			Expect(err).To(HaveOccurred())
			importError, ok := errors.AsImportError(err)
			Expect(ok).To(BeTrue())
			Expect(importError.Statement()).NotTo(BeEmpty())
			Expect(resp).To(BeNil())
		})

		It("execute IsSucceed false", func() {
			mockClient.EXPECT().Execute(gomock.Any()).Times(1).Return(mockResultSet, nil)
			mockResultSet.EXPECT().IsSucceed().Times(1).Return(false)
			mockResultSet.EXPECT().GetError().Times(1).Return(stderrors.New("status failed"))
			resp, err := nodeImporter.Import(spec.Record{"id1", "id2"})
			Expect(err).To(HaveOccurred())
			importError, ok := errors.AsImportError(err)
			Expect(ok).To(BeTrue())
			Expect(importError.Messages).To(ContainElement(ContainSubstring("status failed")))
			Expect(importError.Statement()).NotTo(BeEmpty())
			Expect(resp).To(BeNil())
		})

		It("execute success", func() {
			mockClient.EXPECT().Execute(gomock.Any()).Times(1).Return(mockResultSet, nil)
			mockResultSet.EXPECT().IsSucceed().Times(1).Return(true)
			mockResultSet.EXPECT().GetLatency().Times(1).Return(int64(10))
			resp, err := nodeImporter.Import(spec.Record{"id1", "id2"})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.Latency).To(Equal(time.Microsecond * time.Duration(10)))
		})
	})
})
