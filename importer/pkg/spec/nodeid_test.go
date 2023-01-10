package spec

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NodeID", func() {
	Describe(".Complete", func() {
		It("no prop", func() {
			nodeID := &NodeID{}
			nodeID.Complete()
			prop := &Prop{}
			prop.Complete()
			Expect(nodeID.Prop).To(Equal(prop))
		})

		It("empty prop", func() {
			nodeID := &NodeID{
				Prop: &Prop{},
			}
			nodeID.Complete()
			prop := &Prop{}
			prop.Complete()
			Expect(nodeID.Prop).To(Equal(prop))
		})
	})

	Describe(".Validate", func() {
		It("failed", func() {
			nodeID := &NodeID{Prop: &Prop{}}
			err := nodeID.Validate()
			prop := &Prop{}
			Expect(err).To(Equal(prop.Validate()))
		})

		It("success", func() {
			nodeID := &NodeID{Prop: &Prop{Name: "n", Type: ValueTypeDefault}}
			err := nodeID.Validate()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe(".ValueStatement", func() {
		It("failed", func() {
			nodeID := &NodeID{Prop: &Prop{Name: "id"}}
			nodeID.Complete()
			Expect(nodeID.Validate()).NotTo(HaveOccurred())
			statement, err := nodeID.ValueStatement(nil)
			prop := nodeID.Prop
			statement1, err1 := prop.ValueStatement(nil)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(err1))
			Expect(statement).To(Equal(statement1))
		})

		It("success", func() {
			nodeID := &NodeID{Prop: &Prop{Name: "n", Type: ValueTypeDefault}}
			err := nodeID.Validate()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	DescribeTable(".ValueStatement",
		func(nodeID *NodeID, record Record, exceptStatement string) {
			nodeID.Complete()
			Expect(nodeID.Validate()).NotTo(HaveOccurred())
			statement, err := nodeID.ValueStatement(record)
			props := Props{nodeID.Prop}
			statement1, err1 := props.ValueStatement(record)
			if err1 != nil {
				Expect(err).To(Equal(err1))
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
			Expect(statement).To(Equal(statement1))
			Expect(statement).To(Equal(exceptStatement))
		},
		Entry("no record empty",
			&NodeID{Prop: &Prop{Name: "id"}},
			Record([]string{}),
			"",
		),
		Entry("value type int",
			&NodeID{Prop: &Prop{Name: "id", Type: ValueTypeInt, Index: 0}},
			Record([]string{"1"}),
			"{id: 1}",
		),
		Entry("value type string",
			&NodeID{Prop: &Prop{Name: "id", Type: ValueTypeString, Index: 0}},
			Record([]string{"1"}),
			"{id: \"1\"}",
		),
	)
})
