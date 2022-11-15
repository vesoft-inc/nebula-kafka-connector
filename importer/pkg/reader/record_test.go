package reader

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/source"
)

var _ = Describe("RecordReader", func() {
	var s source.Source
	BeforeEach(func() {
		var err error
		s, err = source.Open(&source.Config{
			Path: "testdata/local.csv",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(s).NotTo(BeNil())
	})
	AfterEach(func() {
		err := s.Close()
		Expect(err).NotTo(HaveOccurred())
	})
	It("should success", func() {
		r := NewRecordReader(s)
		Expect(r).NotTo(BeNil())
	})
})
