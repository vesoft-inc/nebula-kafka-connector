package logger

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Option", func() {
	It("WithLevel", func() {
		o := options{
			level: DebugLevel,
		}
		WithLevel(WarnLevel)(&o)
		Expect(o.level).To(Equal(WarnLevel))
	})

	It("WithLevelText", func() {
		o := options{
			level: DebugLevel,
		}
		WithLevelText("")(&o)
		Expect(o.level).To(Equal(InfoLevel))
	})

	It("WithFields", func() {
		fields := Fields{
			{Key: "i", Value: 1},
			{Key: "f", Value: 1.1},
			{Key: "s", Value: "str"},
		}
		o := options{}
		WithFields(fields...)(&o)
		Expect(o.fields).To(Equal(fields))
	})

	It("nopLogger", func() {
		files := []string{"f1", "f2"}
		o := options{}
		WithFiles(files...)(&o)
		Expect(o.files).To(Equal(files))
	})

	It("nopLogger", func() {
		o := options{}
		WithConsole(true)(&o)
		Expect(o.console).To(BeTrue())
	})
})
