package source

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("localSource", func() {
	It("exists", func() {
		s, err := openLocalFile(&Config{
			Path: "testdata/local.txt",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(s).NotTo(BeNil())

		nBytes, err := s.Size()
		Expect(err).NotTo(HaveOccurred())
		Expect(nBytes).To(Equal(int64(6)))

		var buf [1024]byte
		n, err := s.Read(buf[:])
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(6))

		err = s.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	It("not exists", func() {
		s, err := openLocalFile(&Config{
			Path: "testdata/not-exists.txt",
		})
		Expect(err).To(HaveOccurred())
		Expect(s).To(BeNil())
	})

	It("get size failed", func() {
		s, err := openLocalFile(&Config{
			Path: "testdata/local.txt",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(s).NotTo(BeNil())

		err = s.Close()
		Expect(err).NotTo(HaveOccurred())

		nBytes, err := s.Size()
		Expect(err).To(HaveOccurred())
		Expect(nBytes).To(Equal(int64(0)))
	})
})
