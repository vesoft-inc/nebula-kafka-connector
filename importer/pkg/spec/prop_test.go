package spec

import (
	stderrors "errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
)

var _ = Describe("Prop", func() {
	Describe(".Complete", func() {
		It("no value type", func() {
			prop := &Prop{}
			prop.Complete()
			Expect(prop.Type).To(Equal(ValueTypeDefault))

			prop = &Prop{
				Type: ValueTypeInt,
			}
			prop.Complete()
			Expect(prop.Type).To(Equal(ValueTypeInt))
		})
		It("have value type", func() {
			prop := &Prop{
				Type: ValueTypeInt,
			}
			prop.Complete()
			Expect(prop.Type).To(Equal(ValueTypeInt))
		})

		It("ignore type", func() {
			prop := &Prop{
				Ignore: true,
			}
			prop.Complete()
			Expect(prop.Type.String()).To(BeEmpty())
		})
	})

	DescribeTable(".Validate",
		func(prop *Prop, expectErr error) {
			err := prop.Validate()
			if expectErr != nil {
				if Expect(err).To(HaveOccurred()) {
					Expect(stderrors.Is(err, expectErr)).To(BeTrue())
					e, ok := errors.AsImportError(err)
					Expect(ok).To(BeTrue())
					Expect(e.Cause()).To(Equal(expectErr))
					Expect(e.PropName()).To(Equal(prop.Name))
				}
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
		},
		Entry("no prop name", &Prop{}, errors.ErrNoPropName),
		Entry("unsupported value type", &Prop{Name: "a", Type: "x"}, errors.ErrUnsupportedValueType),
		Entry("supported value type", &Prop{Name: "a", Type: ValueTypeDefault}, nil),
		Entry("ignore", &Prop{Ignore: true}, nil),
	)

	DescribeTable(".ValueStatement",
		func(p *Prop, record Record, expectStatement string, expectErr error) {
			statement, err := p.ValueStatement(record)
			if expectErr != nil {
				if Expect(err).To(HaveOccurred()) {
					Expect(stderrors.Is(err, expectErr)).To(BeTrue())
					e, ok := errors.AsImportError(err)
					Expect(ok).To(BeTrue())
					Expect(e.Cause()).To(Equal(expectErr))
					Expect(e.PropName()).To(Equal(p.Name))
				}
				Expect(statement).To(Equal(expectStatement))
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(statement).To(Equal(expectStatement))
			}
		},
		Entry("no record empty",
			&Prop{},
			Record([]string{}),
			"",
			errors.ErrNoRecord,
		),
		Entry("no record",
			&Prop{
				Name:   "p1",
				Type:   ValueTypeInt,
				Index:  1,
				Ignore: false,
			},
			Record([]string{"0"}),
			"",
			errors.ErrNoRecord,
		),
		Entry("record int",
			&Prop{
				Name:   "p1",
				Type:   ValueTypeInt,
				Index:  0,
				Ignore: false,
			},
			Record([]string{"1"}),
			"p1: 1",
			nil,
		),
		Entry("record string",
			&Prop{
				Name:   "p1",
				Type:   ValueTypeString,
				Index:  0,
				Ignore: false,
			},
			Record([]string{"str"}),
			"p1: \"str\"",
			nil,
		),
		Entry("record int",
			&Prop{
				Name:   "p1",
				Type:   ValueTypeDouble,
				Index:  0,
				Ignore: false,
			},
			Record([]string{"1.1"}),
			"p1: 1.1",
			nil,
		),
		Entry("record ignore",
			&Prop{
				Name:   "p1",
				Type:   ValueTypeString,
				Index:  0,
				Ignore: true,
			},
			Record([]string{"ignore"}),
			"",
			nil,
		),
		Entry("record ignore without record",
			&Prop{
				Ignore: true,
			},
			nil,
			"",
			nil,
		),
	)
})

var _ = Describe("Props", func() {
	Describe(".Complete", func() {
		It("default value type", func() {
			prop1 := Prop{}
			prop2 := Prop{
				Type: ValueTypeInt,
			}
			prop3 := Prop{
				Type: ValueTypeDouble,
			}
			p1, p2, p3 := prop1, prop2, prop3
			props := Props{&p1, &p2, &p3}
			props.Complete()
			Expect(props).To(HaveLen(3))

			p1.Complete()
			Expect(props[0]).To(Equal(&p1))
			p2.Complete()
			Expect(props[1]).To(Equal(&p2))
			p3.Complete()
			Expect(props[2]).To(Equal(&p3))
		})
	})

	DescribeTable(".Validate",
		func(props Props, failedIndex int) {
			err := props.Validate()
			if failedIndex >= 0 {
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(props[failedIndex].Validate()))
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
		},
		Entry("empty props",
			Props{},
			-1,
		),
		Entry("success",
			Props{
				&Prop{Name: "a", Type: ValueTypeInt},
				&Prop{Name: "b", Type: ValueTypeString},
				&Prop{Name: "c", Type: ValueTypeDouble},
				&Prop{Ignore: true},
			},
			-1,
		),
		Entry("failed at 0",
			Props{
				&Prop{Name: ""},
				&Prop{Name: "a", Type: ValueTypeInt},
				&Prop{Name: "b", Type: ValueTypeString},
				&Prop{Name: "c", Type: ValueTypeDouble},
				&Prop{Ignore: true},
			},
			0,
		),
		Entry("failed at 1",
			Props{
				&Prop{Name: "a", Type: ValueTypeInt},
				&Prop{Name: "failed"},
				&Prop{Name: "b", Type: ValueTypeString},
				&Prop{Name: "c", Type: ValueTypeDouble},
				&Prop{Ignore: true},
			},
			1,
		),
		Entry("failed at end",
			Props{
				&Prop{Name: "a", Type: ValueTypeInt},
				&Prop{Name: "b", Type: ValueTypeString},
				&Prop{Name: "c", Type: ValueTypeDouble},
				&Prop{Ignore: true},
				&Prop{Name: "failed"},
			},
			4,
		),
	)

	DescribeTable(".ValueStatement",
		func(props Props, record Record, expectStatement string, failedIndex int) {
			statement, err := props.ValueStatement(record)
			if failedIndex >= 0 {
				if Expect(err).To(HaveOccurred()) {
					_, expectErr := props[failedIndex].ValueStatement(record)
					Expect(err).To(Equal(expectErr))
				}
				Expect(statement).To(Equal(expectStatement))
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(statement).To(Equal(expectStatement))
			}
		},
		Entry("empty props",
			Props{},
			[]string{"1", "1.1", "str"},
			"{}",
			-1,
		),
		Entry("success",
			Props{
				&Prop{Name: "a", Type: ValueTypeInt, Index: 0},
				&Prop{Name: "b", Type: ValueTypeString, Index: 2},
				&Prop{Name: "c", Type: ValueTypeDouble, Index: 1},
				&Prop{Ignore: true},
			},
			[]string{"1", "1.1", "str"},
			"{a: 1, b: \"str\", c: 1.1}",
			-1,
		),
		Entry("failed",
			Props{
				&Prop{Name: "a", Type: ValueTypeInt, Index: 0},
			},
			nil,
			"",
			0,
		),
	)

	DescribeTable(".Append",
		func(l, r Props) {
			lLen, rLen := len(l), len(r)
			props := l.Append(r...)
			Expect(props).To(HaveLen(lLen + rLen))
		},
		Entry("nil + nil",
			nil,
			nil,
		),
		Entry("nil + non-nil",
			nil,
			Props{&Prop{}},
		),
		Entry("non-nil + nil",
			Props{&Prop{}, &Prop{}},
			nil,
		),
		Entry("non-nil + non-nil",
			Props{&Prop{}, &Prop{}, &Prop{}},
			Props{&Prop{}, &Prop{}, &Prop{}, &Prop{}},
		),
	)
})
