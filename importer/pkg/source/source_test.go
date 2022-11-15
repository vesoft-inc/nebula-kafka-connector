package source

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Source", func() {
	It("should success", func() {
		s, err := Open(&Config{
			Path: "testdata/local.txt",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(s).NotTo(BeNil())

		err = s.Close()
		Expect(err).NotTo(HaveOccurred())
	})
})
