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
		Entry(nil, ValueType("int"), true),
		Entry(nil, ValueType("inT"), true),
		Entry(nil, ValueType("iNt"), true),
		Entry(nil, ValueType("InT"), true),
		Entry(nil, ValueType("aa"), false),
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
		Entry(nil, ValueTypeInt, ValueType("aa"), false),
		Entry(nil, ValueType("aa"), ValueTypeInt, false),
		Entry(nil, ValueType("aa"), ValueType("aa"), false),
	)
})
