package source

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe(".String", func() {
		It("should success", func() {
			c := Config{
				Path: "path",
			}
			Expect(c.String()).To(Equal("path"))
		})
	})
})
