package configv5

import (
	"path/filepath"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/source"
	specv5 "github.com/vesoft-inc/nebula-ng-tools/importer/pkg/spec/v5"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manager", func() {
	Describe(".BuildManager", func() {
		var c Config
		BeforeEach(func() {
			c = Config{
				Manager: Manager{
					GraphName: "graphName",
				},
				Sources: Sources{
					{
						SourceConfig: source.Config{
							Path: filepath.Join("testdata", "file10"),
						},
						Nodes: specv5.Nodes{
							&specv5.Node{
								Name: "n1",
								ID: &specv5.NodeID{
									Name:  "id",
									Type:  specv5.ValueTypeString,
									Index: 0,
								},
							},
						},
					},
				},
			}
		})

		It("BuildImporters failed", func() {
			c.Manager.GraphName = ""
			Expect(c.Build()).To(HaveOccurred())
		})

		It("Importer failed", func() {
			c.Sources[0].SourceConfig.Path = filepath.Join("testdata", "not-exists.csv")
			Expect(c.Build()).To(HaveOccurred())
		})

		It("successfully", func() {
			Expect(c.Build()).NotTo(HaveOccurred())
		})
	})
})
