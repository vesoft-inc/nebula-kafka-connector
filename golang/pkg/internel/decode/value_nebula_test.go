package decode

import (
	"testing"
)

func TestValueString(t *testing.T) {
	testcases := []struct {
		value  *nebulaValue
		expect string
	}{
		{
			value:  &nebulaValue{data: nil},
			expect: "null",
		},
		{
			value:  &nebulaValue{data: &NebulaBool{Value: true}},
			expect: "true",
		},
		{
			value:  &nebulaValue{data: &NebulaBool{Value: false}},
			expect: "false",
		},
		{
			value:  &nebulaValue{data: &NebulaInt8{Value: 8}},
			expect: "8",
		},
		{
			value:  &nebulaValue{data: &NebulaInt16{Value: 16}},
			expect: "16",
		},
		{
			value:  &nebulaValue{data: &NebulaFloat{Value: 16.02}},
			expect: "16.02",
		},
		{
			value:  &nebulaValue{data: &NebulaDouble{Value: 16}},
			expect: "16.0",
		},
		{
			value:  &nebulaValue{data: &NebulaString{Value: "string"}},
			expect: `string`,
		},
		{
			value:  &nebulaValue{data: &NebulaString{Value: "abc\rde"}},
			expect: `dec`,
		},
		{
			value: &nebulaValue{data: &NebulaList{
				Values: []*nebulaValue{
					{data: &NebulaInt8{Value: 8}},
					{data: &NebulaString{Value: "dec"}},
				},
			}},
			expect: `[8,dec]`,
		},
		{
			value: &nebulaValue{data: &NebulaRecord{
				Values: map[string]*nebulaValue{
					"int": {data: &NebulaInt8{Value: 8}},
					"str": {data: &NebulaString{Value: "dec"}},
				},
			}},
			expect: `{int:8,str:dec}`,
		},
		{
			value: &nebulaValue{data: &NebulaNode{
				NodeId: 1,
				Properties: map[string]*nebulaValue{
					"int": {data: &NebulaInt8{Value: 8}},
					"str": {data: &NebulaString{Value: "dec"}},
				},
				Type:   "node",
				Labels: []string{"label1", "label2"},
			}},
			expect: `(1@node:label1&label2{int:8,str:dec})`,
		},
		{
			value: &nebulaValue{data: &NebulaEdge{
				SrcId: 1,
				DstId: 2,
				Properties: map[string]*nebulaValue{
					"int": {data: &NebulaInt8{Value: 8}},
					"str": {data: &NebulaString{Value: "dec"}},
				},
				Type:   "edge",
				Labels: []string{"label3", "label4"},
				Rank:   123456,
			}},
			expect: `(1)-[123456@edge:label3&label4{int:8,str:dec}]->(2)`,
		},
		{
			value: &nebulaValue{data: &NebulaEdge{
				SrcId: 1,
				DstId: 2,
				Properties: map[string]*nebulaValue{
					"int": {data: &NebulaInt8{Value: 8}},
					"str": {data: &NebulaString{Value: "dec"}},
				},
				Type:      "edge",
				Labels:    []string{"label3", "label4"},
				Rank:      123456,
				Direction: edgeNoDirection,
			}},
			expect: `(1)~[123456@edge:label3&label4{int:8,str:dec}]~(2)`,
		},
		{
			value: &nebulaValue{data: &NebulaPath{
				Values: []*nebulaValue{
					{data: &NebulaNode{
						NodeId: 1,
						Properties: map[string]*nebulaValue{
							"int": {data: &NebulaInt8{Value: 8}},
							"str": {data: &NebulaString{Value: "dec"}},
						},
						Type:   "node1",
						Labels: []string{"label1", "label2"},
					}},
					{data: &NebulaEdge{
						SrcId: 1,
						DstId: 2,
						Properties: map[string]*nebulaValue{
							"int": {data: &NebulaInt8{Value: 8}},
							"str": {data: &NebulaString{Value: "dec"}},
						},
						Rank:   123456,
						Type:   "edge",
						Labels: []string{"label3", "label4"},
					}},
					{data: &NebulaNode{
						NodeId: 2,
						Properties: map[string]*nebulaValue{
							"int": {data: &NebulaInt8{Value: 9}},
							"str": {data: &NebulaString{Value: "abc"}},
						},
						Type:   "node2",
						Labels: []string{"label5", "label6"},
					},
					}}}},
			expect: `(1@node1:label1&label2{int:8,str:dec})-[123456@edge:label3&label4{int:8,str:dec}]->(2@node2:label5&label6{int:9,str:abc})`,
		},
		{
			value: &nebulaValue{data: &NebulaPath{
				Values: []*nebulaValue{
					{data: &NebulaNode{
						NodeId: 1,
						Type:   "node1",
						Labels: []string{"label1", "label2"},
					},
					},
					{data: &NebulaEdge{
						SrcId:  1,
						DstId:  2,
						Rank:   123456,
						Type:   "edge",
						Labels: []string{"label3", "label4"},
					},
					},
					{data: &NebulaNode{
						NodeId: 2,
						Type:   "node2",
						Labels: []string{"label5", "label6"},
					},
					},
					{data: &NebulaEdge{
						SrcId:  1,
						DstId:  2,
						Rank:   123456,
						Type:   "edge",
						Labels: []string{"label3", "label4"},
					},
					},
					{data: &NebulaNode{
						NodeId: 1,
						Type:   "node1",
						Labels: []string{"label1", "label2"},
					},
					},
				}}},
			expect: `(1@node1:label1&label2{})-[123456@edge:label3&label4{}]->(2@node2:label5&label6{})<-[123456@edge:label3&label4{}]-(1@node1:label1&label2{})`,
		},
		{
			value: &nebulaValue{data: &NebulaPath{
				Values: []*nebulaValue{
					{data: &NebulaNode{
						NodeId: 2,
						Properties: map[string]*nebulaValue{
							"int": {data: &NebulaInt8{Value: 8}},
							"str": {data: &NebulaString{Value: "dec"}},
						},
						Type:   "node1",
						Labels: []string{"label1", "label2"},
					}},
					{data: &NebulaEdge{
						SrcId: 1,
						DstId: 2,
						Properties: map[string]*nebulaValue{
							"int": {data: &NebulaInt8{Value: 8}},
							"str": {data: &NebulaString{Value: "dec"}},
						},
						Rank:      123456,
						Type:      "edge",
						Labels:    []string{"label3", "label4"},
						Direction: edgeInComingDirection,
					}},
					{data: &NebulaNode{
						NodeId: 1,
						Properties: map[string]*nebulaValue{
							"int": {data: &NebulaInt8{Value: 9}},
							"str": {data: &NebulaString{Value: "abc"}},
						},
						Type:   "node2",
						Labels: []string{"label5", "label6"},
					},
					}}}},
			expect: `(2@node1:label1&label2{int:8,str:dec})<-[123456@edge:label3&label4{int:8,str:dec}]-(1@node2:label5&label6{int:9,str:abc})`,
		},
		{
			value: &nebulaValue{data: &NebulaLocalTime{
				Hour: 23, Minute: 59, Sec: 59, Microsec: 999999,
			},
			},
			expect: `23:59:59.999999`,
		},
		{
			value: &nebulaValue{data: &NebulaLocalDatetime{
				Year: 2024, Month: 12, Day: 31, Hour: 23, Minute: 59, Sec: 59, Microsec: 999999,
			}},
			expect: `2024-12-31T23:59:59.999999`,
		},
		{
			value: &nebulaValue{data: &NebulaZonedTime{
				Hour: 23, Minute: 59, Sec: 59, Microsec: 999999,
			}},
			expect: `23:59:59.999999Z`,
		},
		{
			value: &nebulaValue{data: &NebulaZonedDatetime{
				Year: 2024, Month: 12, Day: 31, Hour: 23, Minute: 59, Sec: 59, Microsec: 999999,
			}},
			expect: `2024-12-31T23:59:59.999999Z`,
		},
		{
			value: &nebulaValue{data: &NebulaDuration{
				Sec: 23, Microsec: 999999,
			}},
			expect: `PT23.999999S`,
		},
		{
			value: &nebulaValue{data: &NebulaDuration{
				Year: 1, Month: 1, isMonthBased: true,
			}},
			expect: `P1Y1M`,
		},
		{
			value: &nebulaValue{data: &NebulaDuration{
				Day: 1,
			}},
			expect: `P1D`,
		},
		{
			value: &nebulaValue{data: &NebulaDuration{
				Sec: 59, Microsec: 0,
			}},
			expect: `PT59S`,
		},
		{
			value: &nebulaValue{data: &NebulaDuration{
				Sec: -59, Microsec: -99000,
			}},
			expect: `PT-59.099S`,
		},
		{
			value: &nebulaValue{data: &NebulaDuration{
				Sec: 0, Microsec: -99,
			}},
			expect: `PT-0.000099S`,
		},
		{
			value: &nebulaValue{data: &NebulaDuration{
				Sec: 0, Microsec: 99,
			}},
			expect: `PT0.000099S`,
		},
	}
	for _, c := range testcases {
		v := c.value
		if v.String() != c.expect {
			t.Fatalf("expect %s, got %s", c.expect, v.String())
		}
	}
}
