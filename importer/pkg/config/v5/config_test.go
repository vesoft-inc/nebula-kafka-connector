package configv5

import (
	"path/filepath"

	configbase "github.com/vesoft-inc/nebula-ng-tools/importer/pkg/config/base"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/source"
	specv5 "github.com/vesoft-inc/nebula-ng-tools/importer/pkg/spec/v5"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe(".Optimize", func() {
		It("c.Sources.OptimizePathWildCard failed", func() {
			c := &Config{
				Sources: Sources{
					Source{
						Source: configbase.Source{
							SourceConfig: source.Config{
								Local: &source.LocalConfig{
									Path: "[a-b",
								},
							},
						},
					},
				},
			}
			Expect(c.Optimize(".")).To(HaveOccurred())
		})

		It("successfully", func() {
			c := &Config{
				Sources: Sources{
					Source{
						Source: configbase.Source{
							SourceConfig: source.Config{
								Local: &source.LocalConfig{
									Path: filepath.Join("testdata", "file*"),
								},
							},
						},
					},
				},
			}
			Expect(c.Optimize(".")).NotTo(HaveOccurred())
		})
	})

	Describe(".Build", func() {
		var c Config
		BeforeEach(func() {
			c = Config{
				Manager: Manager{
					GraphName: "graphName",
				},
				Sources: Sources{
					{
						Source: configbase.Source{
							SourceConfig: source.Config{
								Local: &source.LocalConfig{
									Path: filepath.Join("testdata", "file10"),
								},
							},
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

		It("BuildLogger failed", func() {
			c.Log = &Log{
				Files: []string{filepath.Join("testdata", "not-exists", "1.log")},
			}
			Expect(c.Build()).To(HaveOccurred())
		})

		It("BuildClientPool failed", func() {
			c.Client.Version = "v"
			Expect(c.Build()).To(HaveOccurred())
		})

		It("BuildManager failed", func() {
			c.Manager.GraphName = ""
			Expect(c.Build()).To(HaveOccurred())
		})

		It("successfully", func() {
			Expect(c.Build()).NotTo(HaveOccurred())
			Expect(c.GetLogger()).NotTo(BeNil())
			Expect(c.GetClientPool()).NotTo(BeNil())
			Expect(c.GetManager()).NotTo(BeNil())
		})
	})
})
