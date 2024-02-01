package nebula_ng

import (
	"testing"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto"
)

func TestGRPCString(t *testing.T) {
	testcases := []struct {
		value  *proto.Value
		expect string
	}{
		{
			value:  nil,
			expect: "null",
		},
		{
			value:  &proto.Value{Data: &proto.Value_BoolValue{true}},
			expect: "true",
		},
		{
			value:  &proto.Value{Data: &proto.Value_BoolValue{false}},
			expect: "false",
		},
		{
			value:  &proto.Value{Data: &proto.Value_Int8Value{8}},
			expect: "8",
		},
		{
			value:  &proto.Value{Data: &proto.Value_Int16Value{16}},
			expect: "16",
		},
		{
			value:  &proto.Value{Data: &proto.Value_FloatValue{16.02}},
			expect: "16.02",
		},
		{
			value:  &proto.Value{Data: &proto.Value_DoubleValue{16}},
			expect: "16.0",
		},
		{
			value:  &proto.Value{Data: &proto.Value_StringValue{[]byte("string")}},
			expect: `"string"`,
		},
		{
			value:  &proto.Value{Data: &proto.Value_StringValue{[]byte("abc\rde")}},
			expect: `"dec"`,
		},
		{
			value: &proto.Value{Data: &proto.Value_ListValue{
				&proto.List{
					Values: []*proto.Value{
						{Data: &proto.Value_Int8Value{8}},
						{Data: &proto.Value_StringValue{[]byte("dec")}},
					}},
			}},
			expect: `[8, "dec"]`,
		},
		{
			value: &proto.Value{Data: &proto.Value_RecordValue{
				&proto.Record{
					Values: map[string]*proto.Value{
						"int": {Data: &proto.Value_Int8Value{8}},
						"str": {Data: &proto.Value_StringValue{[]byte("dec")}},
					},
				},
			}},
			expect: `{"int":8,"str":"dec"}`,
		},
		{
			value: &proto.Value{Data: &proto.Value_NodeValue{
				&proto.Node{
					NodeId:     1,
					NodeTypeId: 123,
					Properties: map[string]*proto.Value{
						"int": {Data: &proto.Value_Int8Value{8}},
						"str": {Data: &proto.Value_StringValue{[]byte("dec")}},
					},
				},
			}},
			expect: `({"int":8,"str":"dec"})`,
		},
		{
			value: &proto.Value{Data: &proto.Value_EdgeValue{
				&proto.Edge{
					SrcId:      1,
					DstId:      2,
					EdgeTypeId: 123,
					Properties: map[string]*proto.Value{
						"int": {Data: &proto.Value_Int8Value{8}},
						"str": {Data: &proto.Value_StringValue{[]byte("dec")}},
					},
				},
			}},
			expect: `[{"int":8,"str":"dec"}]`,
		},
		{
			value: &proto.Value{Data: &proto.Value_PathValue{
				&proto.Path{
					Values: []*proto.Value{
						{Data: &proto.Value_NodeValue{
							&proto.Node{
								NodeId:     1,
								NodeTypeId: 123,
								Properties: map[string]*proto.Value{
									"int": {Data: &proto.Value_Int8Value{8}},
									"str": {Data: &proto.Value_StringValue{[]byte("dec")}},
								}},
						}},
						{Data: &proto.Value_EdgeValue{
							&proto.Edge{
								SrcId:      1,
								DstId:      2,
								EdgeTypeId: 123,
								Properties: map[string]*proto.Value{
									"int": {Data: &proto.Value_Int8Value{8}},
									"str": {Data: &proto.Value_StringValue{[]byte("dec")}},
								}},
						}},
						{Data: &proto.Value_NodeValue{
							&proto.Node{
								NodeId:     2,
								NodeTypeId: 456,
								Properties: map[string]*proto.Value{
									"int": {Data: &proto.Value_Int8Value{9}},
									"str": {Data: &proto.Value_StringValue{[]byte("abc")}},
								}},
						}},
					}}}},
			expect: `[({"int":8,"str":"dec"}) [{"int":8,"str":"dec"}] ({"int":9,"str":"abc"})]`,
		},
		{
			value: &proto.Value{Data: &proto.Value_LocalTimeValue{
				&proto.LocalTime{
					Hour: 23, Minute: 59, Sec: 59, Microsec: 999999,
				},
			}},
			expect: `23:59:59.999999`,
		},
		{
			value: &proto.Value{Data: &proto.Value_LocalDatatimeValue{
				&proto.LocalDatetime{
					Year: 2024, Month: 12, Day: 31, Hour: 23, Minute: 59, Sec: 59, Microsec: 999999,
				},
			}},
			expect: `2024-12-31T23:59:59.999999`,
		},
		{
			value: &proto.Value{Data: &proto.Value_ZonedTimeValue{
				&proto.ZonedTime{
					Hour: 23, Minute: 59, Sec: 59, Microsec: 999999,
				},
			}},
			expect: `23:59:59.999999Z`,
		},
		{
			value: &proto.Value{Data: &proto.Value_ZonedDatatimeValue{
				&proto.ZonedDatetime{
					Year: 2024, Month: 12, Day: 31, Hour: 23, Minute: 59, Sec: 59, Microsec: 999999,
				},
			}},
			expect: `2024-12-31T23:59:59.999999Z`,
		},
		{
			value: &proto.Value{Data: &proto.Value_DurationValue{
				&proto.Duration{
					Months: 12, Seconds: 23, Microseconds: 999999,
				},
			}},
			expect: `P1YT23.999999S`,
		},
		{
			value: &proto.Value{Data: &proto.Value_DurationValue{
				&proto.Duration{
					Months: 13, Seconds: 3600, Microseconds: 0,
				},
			}},
			expect: `P1Y1MT1H`,
		},
		{
			value: &proto.Value{Data: &proto.Value_DurationValue{
				&proto.Duration{
					Months: 0, Seconds: 3661, Microseconds: 0,
				},
			}},
			expect: `PT1H1M1S`,
		},
	}
	for _, c := range testcases {
		v := grpcValue{data: c.value}
		if v.String() != c.expect {
			t.Fatalf("expect %s, got %s", c.expect, v.String())
		}
	}
}
