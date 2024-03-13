package nebula_ng

import (
	"testing"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/graph"
)

func TestGRPCString(t *testing.T) {
	testcases := []struct {
		value  *graph.Value
		expect string
	}{
		{
			value:  nil,
			expect: "null",
		},
		{
			value:  &graph.Value{Data: &graph.Value_BoolValue{true}},
			expect: "true",
		},
		{
			value:  &graph.Value{Data: &graph.Value_BoolValue{false}},
			expect: "false",
		},
		{
			value:  &graph.Value{Data: &graph.Value_Int8Value{8}},
			expect: "8",
		},
		{
			value:  &graph.Value{Data: &graph.Value_Int16Value{16}},
			expect: "16",
		},
		{
			value:  &graph.Value{Data: &graph.Value_FloatValue{16.02}},
			expect: "16.02",
		},
		{
			value:  &graph.Value{Data: &graph.Value_DoubleValue{16}},
			expect: "16.0",
		},
		{
			value:  &graph.Value{Data: &graph.Value_StringValue{[]byte("string")}},
			expect: `"string"`,
		},
		{
			value:  &graph.Value{Data: &graph.Value_StringValue{[]byte("abc\rde")}},
			expect: `"dec"`,
		},
		{
			value: &graph.Value{Data: &graph.Value_ListValue{
				&graph.List{
					Values: []*graph.Value{
						{Data: &graph.Value_Int8Value{8}},
						{Data: &graph.Value_StringValue{[]byte("dec")}},
					}},
			}},
			expect: `[8, "dec"]`,
		},
		{
			value: &graph.Value{Data: &graph.Value_RecordValue{
				&graph.Record{
					Values: map[string]*graph.Value{
						"int": {Data: &graph.Value_Int8Value{8}},
						"str": {Data: &graph.Value_StringValue{[]byte("dec")}},
					},
				},
			}},
			expect: `{"int":8,"str":"dec"}`,
		},
		{
			value: &graph.Value{Data: &graph.Value_NodeValue{
				&graph.Node{
					NodeId: 1,
					Properties: map[string]*graph.Value{
						"int": {Data: &graph.Value_Int8Value{8}},
						"str": {Data: &graph.Value_StringValue{[]byte("dec")}},
					},
				},
			}},
			expect: `({"int":8,"str":"dec"})`,
		},
		{
			value: &graph.Value{Data: &graph.Value_EdgeValue{
				&graph.Edge{
					SrcId: 1,
					DstId: 2,
					Properties: map[string]*graph.Value{
						"int": {Data: &graph.Value_Int8Value{8}},
						"str": {Data: &graph.Value_StringValue{[]byte("dec")}},
					},
				},
			}},
			expect: `[{"int":8,"str":"dec"}]`,
		},
		{
			value: &graph.Value{Data: &graph.Value_PathValue{
				&graph.Path{
					Values: []*graph.Value{
						{Data: &graph.Value_NodeValue{
							&graph.Node{
								NodeId: 1,
								Properties: map[string]*graph.Value{
									"int": {Data: &graph.Value_Int8Value{8}},
									"str": {Data: &graph.Value_StringValue{[]byte("dec")}},
								}},
						}},
						{Data: &graph.Value_EdgeValue{
							&graph.Edge{
								SrcId: 1,
								DstId: 2,
								Properties: map[string]*graph.Value{
									"int": {Data: &graph.Value_Int8Value{8}},
									"str": {Data: &graph.Value_StringValue{[]byte("dec")}},
								}},
						}},
						{Data: &graph.Value_NodeValue{
							&graph.Node{
								NodeId: 2,
								Type:   "name",
								Properties: map[string]*graph.Value{
									"int": {Data: &graph.Value_Int8Value{9}},
									"str": {Data: &graph.Value_StringValue{[]byte("abc")}},
								}},
						}},
					}}}},
			expect: `[({"int":8,"str":"dec"}) [{"int":8,"str":"dec"}] ({"int":9,"str":"abc"})]`,
		},
		{
			value: &graph.Value{Data: &graph.Value_LocalTimeValue{
				&graph.LocalTime{
					Hour: 23, Minute: 59, Sec: 59, Microsec: 999999,
				},
			}},
			expect: `23:59:59.999999`,
		},
		{
			value: &graph.Value{Data: &graph.Value_LocalDatetimeValue{
				&graph.LocalDatetime{
					Year: 2024, Month: 12, Day: 31, Hour: 23, Minute: 59, Sec: 59, Microsec: 999999,
				},
			}},
			expect: `2024-12-31T23:59:59.999999`,
		},
		{
			value: &graph.Value{Data: &graph.Value_ZonedTimeValue{
				&graph.ZonedTime{
					Hour: 23, Minute: 59, Sec: 59, Microsec: 999999,
				},
			}},
			expect: `23:59:59.999999Z`,
		},
		{
			value: &graph.Value{Data: &graph.Value_ZonedDatetimeValue{
				&graph.ZonedDatetime{
					Year: 2024, Month: 12, Day: 31, Hour: 23, Minute: 59, Sec: 59, Microsec: 999999,
				},
			}},
			expect: `2024-12-31T23:59:59.999999Z`,
		},
		{
			value: &graph.Value{Data: &graph.Value_DurationValue{
				&graph.Duration{
					Months: 12, Seconds: 23, Microseconds: 999999,
				},
			}},
			expect: `P1YT23.999999S`,
		},
		{
			value: &graph.Value{Data: &graph.Value_DurationValue{
				&graph.Duration{
					Months: 13, Seconds: 3600, Microseconds: 0,
				},
			}},
			expect: `P1Y1MT1H`,
		},
		{
			value: &graph.Value{Data: &graph.Value_DurationValue{
				&graph.Duration{
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
