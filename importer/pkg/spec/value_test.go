package spec

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Value", func() {
	DescribeTable("IsSupportedValueType",
		func(t ValueType, expectIsSupported bool) {
			Expect(IsSupportedValueType(t)).To(Equal(expectIsSupported))
		},
		EntryDescription("%[1]s -> %[2]t"),
		Entry(nil, ValueTypeInt, true),
		Entry(nil, ValueTypeString, true),
		Entry(nil, ValueTypeDouble, true),
		Entry(nil, ValueTypeDateTime, true),
		Entry(nil, ValueType("int"), true),
		Entry(nil, ValueType("inT"), true),
		Entry(nil, ValueType("iNt"), true),
		Entry(nil, ValueType("InT"), true),
		Entry(nil, ValueType("unsupported"), false),
	)

	DescribeTable(".ValueStatement",
		func(t ValueType, val, expectedVal string, expectedIsError bool) {
			val, err := t.ValueStatement(val)
			if expectedIsError {
				Expect(err).To(HaveOccurred())
				Expect(val).To(BeEmpty())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal(expectedVal))
			}
		},
		EntryDescription("%[1]s Value(%2s) -> %[3]s, %[4]t"),
		Entry(nil, ValueTypeInt, "1", "1", false),
		Entry(nil, ValueTypeString, "str", "\"str\"", false),
		Entry(nil, ValueTypeDouble, "1.1", "1.1", false),
		Entry(nil, ValueTypeDateTime, "2010-02-14T15:32:10", "DATETIME(\"2010-02-14T15:32:10\")", false),
		Entry(nil, ValueType("unsupported"), "1", "", true),
	)

	DescribeTable(".Equal",
		func(t, vt ValueType, expectIsSupported bool) {
			Expect(t.Equal(vt)).To(Equal(expectIsSupported))
		},
		EntryDescription("%[1]s == %[2]s ? %[3]t"),
		Entry(nil, ValueTypeInt, ValueType("int"), true),
		Entry(nil, ValueTypeInt, ValueType("inT"), true),
		Entry(nil, ValueTypeInt, ValueType("iNt"), true),
		Entry(nil, ValueTypeInt, ValueType("InT"), true),
		Entry(nil, ValueTypeInt, ValueType("unsupported"), false),
		Entry(nil, ValueType("unsupported"), ValueTypeInt, false),
		Entry(nil, ValueType("unsupported"), ValueType("unsupported"), false),
	)
})
