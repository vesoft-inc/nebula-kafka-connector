package config

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/manager"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/source"
)

var _ = Describe("Config", func() {
	var (
		tmpdir string
	)

	BeforeEach(func() {
		var err error
		tmpdir, err = os.MkdirTemp("", "test")
		Expect(err).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		err := os.RemoveAll(tmpdir)
		Expect(err).NotTo(HaveOccurred())
	})

	It("parse", func() {
		c1 := &Config{}
		err := c1.FromFile("testdata/nebula-importer.yaml")
		Expect(err).NotTo(HaveOccurred())

		content, err := c1.Yaml()
		Expect(err).NotTo(HaveOccurred())
		Expect(content).NotTo(BeEmpty())

		c2 := &Config{}
		Expect(err).NotTo(HaveOccurred())
		err = c2.FromBytes(content)
		Expect(err).NotTo(HaveOccurred())
		Expect(c2).To(Equal(c1))

		c3 := &Config{}
		err = c3.FromReader(bytes.NewReader(content))
		Expect(err).NotTo(HaveOccurred())
	})

	It("parse file not exists", func() {
		c := &Config{}
		err := c.FromFile("testdata/not-exists.yaml")
		Expect(err).To(HaveOccurred())
	})

	It("BuildLogger failed", func() {
		c := &Config{}
		err := c.FromFile("testdata/nebula-importer.yaml")
		Expect(err).NotTo(HaveOccurred())
		Expect(c.Log.Files).NotTo(BeEmpty())

		c.Log.Files = []string{filepath.Join(tmpdir, "not-exists", "1.log")}
		l, err := c.BuildLogger()
		Expect(err).To(HaveOccurred())
		Expect(l).To(BeNil())
	})

	It("BuildLogger success", func() {
		c := &Config{}
		err := c.FromFile("testdata/nebula-importer.yaml")
		Expect(err).NotTo(HaveOccurred())
		Expect(c.Log.Files).NotTo(BeEmpty())

		c.Log.Files = []string{filepath.Join(tmpdir, "1.log")}
		l, err := c.BuildLogger()
		Expect(err).NotTo(HaveOccurred())
		defer l.Close()
		Expect(l).NotTo(BeNil())
	})

	It("BuildClient failed", func() {
		c := &Config{}
		err := c.FromFile("testdata/nebula-importer.yaml")
		Expect(err).NotTo(HaveOccurred())
		Expect(c.Log.Files).NotTo(BeEmpty())

		cli, err := c.BuildClient()
		Expect(err).To(HaveOccurred())
		Expect(cli).To(BeNil())
	})

	It("BuildLogger success", func() {
		c := &Config{}
		err := c.FromFile("testdata/nebula-importer.yaml")
		Expect(err).NotTo(HaveOccurred())
		Expect(c.Log.Files).NotTo(BeEmpty())

		ctrl := gomock.NewController(GinkgoT())
		defer ctrl.Finish()
		mockNebulaSession := client.NewMockNebulaSession(ctrl)
		mockNebulaSession.EXPECT().Open().AnyTimes().Return(nil)
		mockNebulaSession.EXPECT().Close().AnyTimes().Return(nil)
		cli, err := c.BuildClient(client.WithNewNebulaSessionFunc(func(_ nebula.HostAddress) client.NebulaSession {
			return mockNebulaSession
		}))
		Expect(err).NotTo(HaveOccurred())
		Expect(cli).NotTo(BeNil())
		cli.Close()
	})

	It("BuildGraph failed", func() {
		c := &Config{}
		err := c.FromFile("testdata/nebula-importer.yaml")
		Expect(err).NotTo(HaveOccurred())
		Expect(c.Log.Files).NotTo(BeEmpty())

		c.Nodes[0].Name = ""
		graph, err := c.BuildGraph()
		Expect(err).To(HaveOccurred())
		Expect(graph).To(BeNil())
	})

	It("BuildGraph success", func() {
		c := &Config{}
		err := c.FromFile("testdata/nebula-importer.yaml")
		Expect(err).NotTo(HaveOccurred())
		Expect(c.Log.Files).NotTo(BeEmpty())

		graph, err := c.BuildGraph()
		Expect(err).NotTo(HaveOccurred())
		Expect(graph).NotTo(BeNil())
	})

	It("BuildManager success", func() {
		c := &Config{}
		err := c.FromFile("testdata/nebula-importer.yaml")
		Expect(err).NotTo(HaveOccurred())
		Expect(c.Log.Files).NotTo(BeEmpty())

		graph, err := c.BuildGraph()
		Expect(err).NotTo(HaveOccurred())
		Expect(graph).NotTo(BeNil())

		m, err := c.BuildManager(manager.WithGraph(graph))
		Expect(err).NotTo(HaveOccurred())
		Expect(m).NotTo(BeNil())
	})

	It("BuildManager failed at node file", func() {
		c := &Config{}
		err := c.FromFile("testdata/nebula-importer.yaml")
		Expect(err).NotTo(HaveOccurred())
		Expect(c.Log.Files).NotTo(BeEmpty())

		c.Graph.Nodes[0].SourceConfigs = []*source.Config{
			{
				Path: filepath.Join(tmpdir, "not-exists.yaml"),
			},
		}
		graph, err := c.BuildGraph()
		Expect(err).NotTo(HaveOccurred())
		Expect(graph).NotTo(BeNil())

		m, err := c.BuildManager(manager.WithGraph(graph))
		Expect(err).To(HaveOccurred())
		Expect(m).To(BeNil())
	})

	It("BuildManager failed at edge file", func() {
		c := &Config{}
		err := c.FromFile("testdata/nebula-importer.yaml")
		Expect(err).NotTo(HaveOccurred())
		Expect(c.Log.Files).NotTo(BeEmpty())

		c.Graph.Edges[0].SourceConfigs = []*source.Config{
			{
				Path: filepath.Join(tmpdir, "not-exists.yaml"),
			},
		}
		graph, err := c.BuildGraph()
		Expect(err).NotTo(HaveOccurred())
		Expect(graph).NotTo(BeNil())

		m, err := c.BuildManager(manager.WithGraph(graph))
		Expect(err).To(HaveOccurred())
		Expect(m).To(BeNil())
	})
})
