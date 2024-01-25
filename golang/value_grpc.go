package nebula_ng

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto"
)

type (
	grpcValue struct {
		data *proto.Value
	}

	grpcLocalTime struct {
		data *proto.LocalTime
	}
	grpcLocalDatetime struct {
		data *proto.LocalDatetime
	}
	grpcDate struct {
		data *proto.Date
	}
	grpcList struct {
		data *proto.List
	}
	grpcRecord struct {
		data *proto.Record
	}
	grpcNode struct {
		data *proto.Node
	}
	grpcEdge struct {
		data *proto.Edge
	}
	grpcPath struct {
		data *proto.Path
	}
)

func (v *grpcValue) String() string {
	switch v.GetType() {
	case ValueTypeNull:
		return "null"
	case ValueTypeBool:
		return fmt.Sprintf("%t", v.data.GetBoolValue())
	case ValueTypeInt8:
		return fmt.Sprintf("%d", v.data.GetInt8Value())
	case ValueTypeInt16:
		return fmt.Sprintf("%d", v.data.GetInt16Value())
	case ValueTypeInt32:
		return fmt.Sprintf("%d", v.data.GetInt32Value())
	case ValueTypeInt64:
		return fmt.Sprintf("%d", v.data.GetInt64Value())
	case ValueTypeUInt8:
		return fmt.Sprintf("%d", v.data.GetUint8Value())
	case ValueTypeUInt16:
		return fmt.Sprintf("%d", v.data.GetUint16Value())
	case ValueTypeUInt32:
		return fmt.Sprintf("%d", v.data.GetUint32Value())
	case ValueTypeUInt64:
		return fmt.Sprintf("%d", v.data.GetUint64Value())
	case ValueTypeFloat:
		fStr := strconv.FormatFloat(float64(v.data.GetFloatValue()), 'g', -1, 32)
		if !strings.Contains(fStr, ".") {
			fStr = fStr + ".0"
		}
		return fStr
	case ValueTypeDouble:
		fStr := strconv.FormatFloat(v.data.GetDoubleValue(), 'g', -1, 64)
		if !strings.Contains(fStr, ".") {
			fStr = fStr + ".0"
		}
		return fStr
	case ValueTypeString:
		s, _ := v.AsString()
		return s.String()
	case ValueTypeDuration:
		d, _ := v.AsDuration()
		return time.Duration(d).String()
	case ValueTypeLocalDateTime:
		dt, _ := v.AsLocalDatetime()
		return dt.String()
	case ValueTypeDate:
		d, _ := v.AsDate()
		return d.String()
	case ValueTypeLocalTime:
		t, _ := v.AsLocalTime()
		return t.String()
	case ValueTypeList:
		l, _ := v.AsList()
		return l.String()
	case ValueTypeRecord:
		r, _ := v.AsRecord()
		return r.String()
	case ValueTypeNode:
		n, _ := v.AsNode()
		return n.String()
	case ValueTypeEdge:
		e, _ := v.AsEdge()
		return e.String()
	case ValueTypePath:
		p, _ := v.AsPath()
		return p.String()

	default:
		return fmt.Sprintf("%v", v.data)
	}
}

func (v *grpcValue) GetType() ValueType {
	if v.data == nil || v.data.Data == nil {
		return ValueTypeNull
	}
	switch v.data.Data.(type) {
	case *proto.Value_BoolValue:
		return ValueTypeBool
	case *proto.Value_Int8Value:
		return ValueTypeInt8
	case *proto.Value_Int16Value:
		return ValueTypeInt16
	case *proto.Value_Int32Value:
		return ValueTypeInt32
	case *proto.Value_Int64Value:
		return ValueTypeInt64
	case *proto.Value_Uint8Value:
		return ValueTypeUInt8
	case *proto.Value_Uint16Value:
		return ValueTypeUInt16
	case *proto.Value_Uint32Value:
		return ValueTypeUInt32
	case *proto.Value_Uint64Value:
		return ValueTypeUInt64
	case *proto.Value_FloatValue:
		return ValueTypeFloat
	case *proto.Value_DoubleValue:
		return ValueTypeDouble
	case *proto.Value_StringValue:
		return ValueTypeString
	case *proto.Value_DurationValue:
		return ValueTypeDuration
	case *proto.Value_LocalTimeValue:
		return ValueTypeLocalTime
	case *proto.Value_LocalDatatimeValue:
		return ValueTypeLocalDateTime
	case *proto.Value_DateValue:
		return ValueTypeDate
	case *proto.Value_ListValue:
		return ValueTypeList
	case *proto.Value_RecordValue:
		return ValueTypeRecord
	case *proto.Value_NodeValue:
		return ValueTypeNode
	case *proto.Value_EdgeValue:
		return ValueTypeEdge
	case *proto.Value_PathValue:
		return ValueTypePath
	default:
		return ValueTypeNull
	}
}

func (v *grpcValue) IsNull() bool {
	return v.GetType() == ValueTypeNull
}

func (v *grpcValue) AsBool() (Bool, error) {
	if v.GetType() != ValueTypeBool {
		return false, errType("value is not bool")
	}
	return Bool(v.data.GetBoolValue()), nil
}

func (v *grpcValue) AsInt8() (Int8, error) {
	if v.GetType() != ValueTypeInt8 {
		return 0, errType("value is not int8")
	}
	return Int8(v.data.GetInt8Value()), nil
}

func (v *grpcValue) AsInt16() (Int16, error) {
	if v.GetType() != ValueTypeInt16 {
		return 0, errType("value is not int16")
	}
	return Int16(v.data.GetInt16Value()), nil
}

func (v *grpcValue) AsInt32() (Int32, error) {
	if v.GetType() != ValueTypeInt32 {
		return 0, errType("value is not int32")
	}
	return Int32(v.data.GetInt32Value()), nil
}

func (v *grpcValue) AsInt64() (Int64, error) {
	if v.GetType() != ValueTypeInt64 {
		return 0, errType("value is not int64")
	}
	return Int64(v.data.GetInt64Value()), nil
}

func (v *grpcValue) AsUInt8() (UInt8, error) {
	if v.GetType() != ValueTypeUInt8 {
		return 0, errType("value is not uint8")
	}
	return UInt8(v.data.GetUint8Value()), nil
}

func (v *grpcValue) AsUInt16() (UInt16, error) {
	if v.GetType() != ValueTypeUInt16 {
		return 0, errType("value is not uint16")
	}
	return UInt16(v.data.GetUint16Value()), nil
}

func (v *grpcValue) AsUInt32() (UInt32, error) {
	if v.GetType() != ValueTypeUInt32 {
		return 0, errType("value is not uint32")
	}
	return UInt32(v.data.GetUint32Value()), nil
}

func (v *grpcValue) AsUInt64() (UInt64, error) {
	if v.GetType() != ValueTypeUInt64 {
		return 0, errType("value is not uint64")
	}
	return UInt64(v.data.GetUint64Value()), nil
}

func (v *grpcValue) AsFloat() (Float, error) {
	if v.GetType() != ValueTypeFloat {
		return 0, errType("value is not float")
	}
	return Float(v.data.GetFloatValue()), nil
}

func (v *grpcValue) AsDouble() (Double, error) {
	if v.GetType() != ValueTypeDouble {
		return 0, errType("value is not double")
	}
	return Double(v.data.GetDoubleValue()), nil
}

func (v *grpcValue) AsString() (String, error) {
	if v.GetType() != ValueTypeString {
		return "", errType("value is not string")
	}
	return String(v.data.GetStringValue()), nil
}

func (v *grpcValue) AsList() (List, error) {
	if v.GetType() != ValueTypeList {
		return nil, errType("value is not list")
	}
	return &grpcList{data: v.data.GetListValue()}, nil
}

func (v *grpcValue) AsRecord() (Record, error) {
	if v.GetType() != ValueTypeRecord {
		return nil, errType("value is not record")
	}
	return &grpcRecord{data: v.data.GetRecordValue()}, nil
}

func (v *grpcValue) AsDuration() (Duration, error) {
	if v.GetType() != ValueTypeDuration {
		return 0, errType("value is not duration")
	}
	var d time.Duration
	d += time.Duration(v.data.GetDurationValue().GetMicroseconds()) * time.Microsecond
	d += time.Duration(v.data.GetDurationValue().GetSeconds()) * time.Second
	//TODO
	// d += time.Duration(v.data.GetDurationValue().GetMinutes()) * time.Minute
	return Duration(d), nil
}

func (v *grpcValue) AsNode() (Node, error) {
	if v.GetType() != ValueTypeNode {
		return nil, errType("value is not node")
	}
	return &grpcNode{data: v.data.GetNodeValue()}, nil
}

func (v *grpcValue) AsEdge() (Edge, error) {
	if v.GetType() != ValueTypeEdge {
		return nil, errType("value is not edge")
	}
	return &grpcEdge{data: v.data.GetEdgeValue()}, nil
}

func (v *grpcValue) AsPath() (Path, error) {
	if v.GetType() != ValueTypePath {
		return nil, errType("value is not path")
	}
	return &grpcPath{data: v.data.GetPathValue()}, nil
}

func (l *grpcList) String() string {
	valuesStr := make([]string, 0, len(l.data.Values))
	for _, v := range l.GetValues() {
		valuesStr = append(valuesStr, v.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(valuesStr, ", "))
}

func (v *grpcValue) AsLocalDatetime() (LocalDatetime, error) {
	if v.GetType() != ValueTypeLocalDateTime {
		return nil, errType("value is not local datetime")
	}
	return &grpcLocalDatetime{data: v.data.GetLocalDatatimeValue()}, nil
}

func (v *grpcValue) AsDate() (Date, error) {
	if v.GetType() != ValueTypeDate {
		return nil, errType("value is not date")
	}
	return &grpcDate{data: v.data.GetDateValue()}, nil
}

func (v *grpcValue) AsLocalTime() (Time, error) {
	if v.GetType() != ValueTypeLocalTime {
		return nil, errType("value is not local time")
	}
	return &grpcLocalTime{data: v.data.GetLocalTimeValue()}, nil
}

func (l *grpcList) GetValues() []Value {
	values := make([]Value, 0, len(l.data.Values))
	for _, v := range l.data.Values {
		values = append(values, &grpcValue{data: v})
	}
	return values
}

func (l *grpcList) Size() int {
	return len(l.data.Values)
}

func (r *grpcRecord) String() string {
	var kvStr []string = make([]string, 0, len(r.data.Values))
	var keys []string = make([]string, 0, len(r.data.Values))
	for key := range r.data.Values {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, key := range keys {
		v := grpcValue{data: r.data.GetValues()[key]}
		kvTemp := fmt.Sprintf(`"%s":%s`, key, v.String())
		kvStr = append(kvStr, kvTemp)
	}
	return fmt.Sprintf("{%s}", strings.Join(kvStr, ","))
}

func (r *grpcRecord) GetValues() map[string]Value {
	values := make(map[string]Value)
	for k, v := range r.data.GetValues() {
		values[k] = &grpcValue{data: v}
	}
	return values
}

func (n *grpcNode) String() string {
	var kvStr []string = make([]string, 0, len(n.data.Properties))
	var keys []string = make([]string, 0, len(n.data.Properties))
	for key := range n.data.Properties {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	for _, key := range keys {
		v := grpcValue{data: n.data.Properties[key]}
		kvTemp := fmt.Sprintf(`"%s":%s`, key, v.String())
		kvStr = append(kvStr, kvTemp)
	}
	//TODO, should verify if we print the internal id
	// return fmt.Sprintf("(%d :{%s})", n.data.NodeID, strings.Join(kvStr, ","))
	return fmt.Sprintf("({%s})", strings.Join(kvStr, ","))
}

func (n *grpcNode) GetProperties() map[string]Value {
	properties := make(map[string]Value)
	for k, v := range n.data.GetProperties() {
		properties[k] = &grpcValue{data: v}
	}
	return properties
}

func (n *grpcNode) GetNodeTypeId() int32 {
	return n.data.NodeTypeId
}

func (e *grpcEdge) String() string {
	var kvStr []string = make([]string, 0, len(e.data.Properties))
	var keys []string = make([]string, 0, len(e.data.Properties))
	for key := range e.data.Properties {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	for _, key := range keys {
		v := grpcValue{data: e.data.Properties[key]}
		kvTemp := fmt.Sprintf(`"%s":%s`, key, v.String())
		kvStr = append(kvStr, kvTemp)
	}
	return fmt.Sprintf("[{%s}]", strings.Join(kvStr, ","))
}

func (e *grpcEdge) GetProperties() map[string]Value {
	properties := make(map[string]Value)
	for k, v := range e.data.GetProperties() {
		properties[k] = &grpcValue{data: v}
	}
	return properties
}

func (e *grpcEdge) GetEdgeTypeId() int32 {
	return e.data.EdgeTypeId
}

func (p *grpcPath) String() string {
	var valuesStr []string
	for _, v := range p.GetValues() {
		valuesStr = append(valuesStr, v.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(valuesStr, " "))
}

func (p *grpcPath) GetValues() []Value {
	values := make([]Value, 0, len(p.data.GetValues()))
	for _, v := range p.data.GetValues() {
		values = append(values, &grpcValue{data: v})
	}
	return values
}

func (l *grpcLocalDatetime) String() string {
	//RFC3339 without timezone
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%06d",
		l.data.Year, l.data.Month, l.data.Day,
		l.data.Hour, l.data.Minute, l.data.Sec, l.data.Microsec)
}

func (l *grpcLocalDatetime) GetYear() uint32 {
	return l.data.Year
}

func (l *grpcLocalDatetime) GetMonth() uint32 {
	return l.data.Month
}

func (l *grpcLocalDatetime) GetDay() uint32 {
	return l.data.Day
}

func (l *grpcLocalDatetime) GetHour() uint32 {
	return l.data.Hour
}

func (l *grpcLocalDatetime) GetMinute() uint32 {
	return l.data.Minute
}

func (l *grpcLocalDatetime) GetSec() uint32 {
	return l.data.Sec
}

func (l *grpcLocalDatetime) GetMicrosec() uint32 {
	return l.data.Microsec
}

func (d *grpcDate) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.data.Year, d.data.Month, d.data.Day)
}

func (d *grpcDate) GetYear() uint32 {
	return d.data.Day
}

func (d *grpcDate) GetMonth() uint32 {
	return d.data.Month
}

func (d *grpcDate) GetDay() uint32 {
	return d.data.Day
}

func (t *grpcLocalTime) String() string {
	return fmt.Sprintf("%02d:%02d:%02d.%06d",
		t.data.Hour, t.data.Minute, t.data.Sec, t.data.Microsec)
}

func (t *grpcLocalTime) GetHour() uint32 {
	return t.data.Hour
}

func (t *grpcLocalTime) GetMinute() uint32 {
	return t.data.Minute
}

func (t *grpcLocalTime) GetSec() uint32 {
	return t.data.Sec
}

func (t *grpcLocalTime) GetMicrosec() uint32 {
	return t.data.Microsec
}
