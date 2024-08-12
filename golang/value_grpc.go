package nebula_ng

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/common"
)

type (
	grpcValue struct {
		data *common.Value
	}
	grpcLocalTime struct {
		data *common.LocalTime
	}
	grpcLocalDatetime struct {
		data *common.LocalDatetime
	}
	grpcZonedTime struct {
		data *common.ZonedTime
	}
	grpcZonedDatetime struct {
		data *common.ZonedDatetime
	}
	grpcDuration struct {
		data *common.Duration
	}
	grpcDate struct {
		data *common.Date
	}
	grpcList struct {
		data *common.List
	}
	grpcRecord struct {
		data *common.Record
	}
	grpcNode struct {
		data *common.Node
	}
	grpcEdge struct {
		data *common.Edge
	}
	grpcPath struct {
		data *common.Path
	}
	grpcDecimal struct {
		data *common.Decimal
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
		return d.String()
	case ValueTypeDate:
		d, _ := v.AsDate()
		return d.String()
	case ValueTypeLocalDateTime:
		dt, _ := v.AsLocalDatetime()
		return dt.String()
	case ValueTypeLocalTime:
		t, _ := v.AsLocalTime()
		return t.String()
	case ValueTypeZonedTime:
		t, _ := v.AsZonedTime()
		return t.String()
	case ValueTypeZonedDateTime:
		dt, _ := v.AsZonedDatetime()
		return dt.String()
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
	case ValueTypeDecimal:
		d, _ := v.AsDecimal()
		return d.String()
	default:
		return fmt.Sprintf("%v", v.data)
	}
}

func (v *grpcValue) GetType() ValueType {
	if v.data == nil || v.data.Data == nil {
		return ValueTypeNull
	}
	switch v.data.Data.(type) {
	case *common.Value_BoolValue:
		return ValueTypeBool
	case *common.Value_Int8Value:
		return ValueTypeInt8
	case *common.Value_Int16Value:
		return ValueTypeInt16
	case *common.Value_Int32Value:
		return ValueTypeInt32
	case *common.Value_Int64Value:
		return ValueTypeInt64
	case *common.Value_Uint8Value:
		return ValueTypeUInt8
	case *common.Value_Uint16Value:
		return ValueTypeUInt16
	case *common.Value_Uint32Value:
		return ValueTypeUInt32
	case *common.Value_Uint64Value:
		return ValueTypeUInt64
	case *common.Value_FloatValue:
		return ValueTypeFloat
	case *common.Value_DoubleValue:
		return ValueTypeDouble
	case *common.Value_StringValue:
		return ValueTypeString
	case *common.Value_DurationValue:
		return ValueTypeDuration
	case *common.Value_LocalTimeValue:
		return ValueTypeLocalTime
	case *common.Value_LocalDatetimeValue:
		return ValueTypeLocalDateTime
	case *common.Value_ZonedTimeValue:
		return ValueTypeZonedTime
	case *common.Value_ZonedDatetimeValue:
		return ValueTypeZonedDateTime
	case *common.Value_DateValue:
		return ValueTypeDate
	case *common.Value_ListValue:
		return ValueTypeList
	case *common.Value_RecordValue:
		return ValueTypeRecord
	case *common.Value_NodeValue:
		return ValueTypeNode
	case *common.Value_EdgeValue:
		return ValueTypeEdge
	case *common.Value_PathValue:
		return ValueTypePath
	case *common.Value_DecimalValue:
		return ValueTypeDecimal
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
		return nil, errType("value is not duration")
	}
	return &grpcDuration{data: v.data.GetDurationValue()}, nil
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

func (v *grpcValue) AsDecimal() (Decimal, error) {
	if v.GetType() != ValueTypeDecimal {
		return nil, errType("value is not decimal")
	}
	return &grpcDecimal{data: v.data.GetDecimalValue()}, nil
}

func (l *grpcList) String() string {
	valuesStr := make([]string, 0, len(l.data.Values))
	for _, v := range l.GetValues() {
		valuesStr = append(valuesStr, v.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(valuesStr, ","))
}

func (v *grpcValue) AsLocalDatetime() (LocalDatetime, error) {
	if v.GetType() != ValueTypeLocalDateTime {
		return nil, errType("value is not local datetime")
	}
	return &grpcLocalDatetime{data: v.data.GetLocalDatetimeValue()}, nil
}

func (v *grpcValue) AsDate() (Date, error) {
	if v.GetType() != ValueTypeDate {
		return nil, errType("value is not date")
	}
	return &grpcDate{data: v.data.GetDateValue()}, nil
}

func (v *grpcValue) AsLocalTime() (LocalTime, error) {
	if v.GetType() != ValueTypeLocalTime {
		return nil, errType("value is not local time")
	}
	return &grpcLocalTime{data: v.data.GetLocalTimeValue()}, nil
}

func (v *grpcValue) AsZonedTime() (ZonedTime, error) {
	if v.GetType() != ValueTypeZonedTime {
		return nil, errType("value is not zoned time")
	}
	return &grpcZonedTime{data: v.data.GetZonedTimeValue()}, nil
}

func (v *grpcValue) AsZonedDatetime() (ZonedDatetime, error) {
	if v.GetType() != ValueTypeZonedDateTime {
		return nil, errType("value is not zoned time")
	}
	return &grpcZonedDatetime{data: v.data.GetZonedDatetimeValue()}, nil
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
	mv := mapValue(r.GetValues())
	return fmt.Sprintf("{%s}", mv.string())
}

func (r *grpcRecord) GetValues() map[string]Value {
	values := make(map[string]Value)
	for k, v := range r.data.GetValues() {
		values[k] = &grpcValue{data: v}
	}
	return values
}

func (g *grpcDecimal) String() string {
	return g.data.Sval
}

// (288314845273522179@City:City&Place{id:32,name:Norway,url:http://dbpedia.org/resource/Norway})
func (n *grpcNode) String() string {
	mv := mapValue(n.GetProperties())
	return fmt.Sprintf("(%d@%s:%s{%s})",
		n.GetId(),
		n.GetType(),
		strings.Join(n.GetLabels(), "&"),
		mv.string(),
	)
}

func (n *grpcNode) GetProperties() map[string]Value {
	properties := make(map[string]Value)
	for k, v := range n.data.GetProperties() {
		properties[k] = &grpcValue{data: v}
	}
	return properties
}

func (n *grpcNode) GetId() int64 {
	return n.data.NodeId
}

func (n *grpcNode) GetGraph() string {
	return n.data.Graph
}

func (n *grpcNode) GetType() string {
	return n.data.Type
}

func (n *grpcNode) GetLabels() []string {
	return n.data.Labels
}

// (288314845273522179)<-[288314845273522179@connected_Sub_Load:
// connected_Sub_Load{cid:115967690673232363}]-(288314982712475649)
func (e *grpcEdge) String() string {
	mv := mapValue(e.GetProperties())
	var (
		leftBracket  string
		rightBracket string
	)
	if e.IsDirected() {
		leftBracket = "-"
		rightBracket = "->"
	} else {
		leftBracket = "~"
		rightBracket = "~"
	}

	return fmt.Sprintf("(%d)%s[%d@%s:%s{%s}]%s(%d)",
		e.GetSrcId(),
		leftBracket,
		e.GetRank(),
		e.GetType(),
		strings.Join(e.GetLabels(), "&"),
		mv.string(),
		rightBracket,
		e.GetDstId(),
	)
}

func (e *grpcEdge) GetProperties() map[string]Value {
	properties := make(map[string]Value)
	for k, v := range e.data.GetProperties() {
		properties[k] = &grpcValue{data: v}
	}
	return properties
}

func (e *grpcEdge) GetSrcId() int64 {
	return e.data.SrcId
}

func (e *grpcEdge) GetDstId() int64 {
	return e.data.DstId
}

func (e *grpcEdge) GetGraph() string {
	return e.data.Graph
}

func (e *grpcEdge) GetType() string {
	return e.data.Type
}

func (e *grpcEdge) GetLabels() []string {
	return e.data.Labels
}

func (e *grpcEdge) GetRank() int64 {
	return e.data.Rank
}

func (e *grpcEdge) IsDirected() bool {
	return e.data.Direction == common.Edge_DIRECTED
}

// (288314845273522179@City:City&Place{id:3,kind:3,name:org3,url:https://org3.com})-[288314845273522179@connected_Sub_Load:connected_Sub_Load{}]-(288315214640709633@City:City&Place{id:3,kind:city,name:Hangzhou,url:https://hangzhou.com})
// -[288315214640709633@connected_Sub_Load:connected_Sub_Load{}]
// -(288315309129990145@City:City&Place{id:6,kind:city,name:Shenzhen,url:https://shenzhen.com})
func (p *grpcPath) String() string {
	values := p.GetValues()
	buf := bytes.NewBuffer(nil)
	defer buf.Reset()
	var preN Node
	for _, v := range values {
		if v.GetType() == ValueTypeNode {
			n, _ := v.AsNode()
			preN = n
			buf.WriteString(n.String())
		} else if v.GetType() == ValueTypeEdge {
			e, _ := v.AsEdge()
			mv := mapValue(e.GetProperties())
			estr := fmt.Sprintf("[%d@%s:%s{%s}]",
				e.GetRank(),
				e.GetType(),
				strings.Join(e.GetLabels(), "&"),
				mv.string(),
			)

			if e.IsDirected() {
				if e.GetSrcId() == preN.GetId() {
					buf.WriteString(fmt.Sprintf("-%s->", estr))
				} else {
					buf.WriteString(fmt.Sprintf("<-%s-", estr))
				}
			} else {
				buf.WriteString(fmt.Sprintf("~%s~", estr))
			}
		} else {
			//no other type
		}
	}
	return buf.String()

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

func (l *grpcLocalDatetime) GetYear() int32 {
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

func (d *grpcDate) GetYear() int32 {
	return d.data.Year
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

func (d *grpcDuration) String() string {
	var prefix string = "P"
	if d.IsMonthBased() {
		if d.GetYear() != 0 {
			prefix += fmt.Sprintf("%dY", d.GetYear())
		}
		if d.GetMonth() != 0 {
			prefix += fmt.Sprintf("%dM", d.GetMonth())
		}
	} else {
		if d.GetDay() != 0 {
			prefix += fmt.Sprintf("%dD", d.GetDay())
		}
		if d.GetHour() != 0 || d.GetMinute() != 0 || d.GetSecond() != 0 || d.GetMicrosecond() != 0 {
			prefix += "T"
		}
		if d.GetHour() != 0 {
			prefix += fmt.Sprintf("%dH", d.GetHour())
		}
		if d.GetMinute() != 0 {
			prefix += fmt.Sprintf("%dM", d.GetMinute())
		}
		if d.GetSecond() != 0 || d.GetMicrosecond() != 0 {
			if d.GetMicrosecond() == 0 {
				prefix += fmt.Sprintf("%dS", d.GetSecond())
			} else {
				ms := d.GetSecond()*1e6 + d.GetMicrosecond()
				isMinus := d.GetSecond() < 0 || d.GetMicrosecond() < 0
				if isMinus {
					ms = -ms
				}
				s, ss := ms/1e6, ms%1e6
				if isMinus {
					prefix += fmt.Sprintf("-%d.%06d", s, ss)
				} else {
					prefix += fmt.Sprintf("%d.%06d", s, ss)
				}
				prefix = strings.TrimRight(prefix, "0")
				prefix += "S"
			}
		}
	}
	return prefix
}

func (d *grpcDuration) IsMonthBased() bool {
	return d.data.IsMonthBased
}

func (d *grpcDuration) GetYear() int32 {
	return d.data.GetYear()
}

func (d *grpcDuration) GetMonth() int32 {
	return d.data.GetMonth()
}

func (d *grpcDuration) GetDay() int32 {
	return d.data.GetDay()
}

func (d *grpcDuration) GetHour() int32 {
	return d.data.GetHour()
}

func (d *grpcDuration) GetMinute() int32 {
	return d.data.GetMinute()
}

func (d *grpcDuration) GetSecond() int32 {
	return d.data.GetSec()
}

func (d *grpcDuration) GetMicrosecond() int32 {
	return d.data.GetMicrosec()
}

func (zt *grpcZonedTime) String() string {
	//TODO server would return offset with seconds
	offset := zt.GetOffset()
	var zone string
	if offset < 0 {
		zone = fmt.Sprintf("-%02d:%02d", -offset/3600, (-offset%3600)/60)
	} else if offset > 0 {
		zone = fmt.Sprintf("+%02d:%02d", offset/3600, (offset%3600)/60)
	} else {
		zone = "Z"
	}

	return fmt.Sprintf("%02d:%02d:%02d.%06d%s",
		zt.data.Hour,
		zt.data.Minute,
		zt.data.Sec,
		zt.data.Microsec,
		zone,
	)
}

func (zt *grpcZonedTime) GetHour() uint32 {
	return zt.data.Hour
}

func (zt *grpcZonedTime) GetMinute() uint32 {
	return zt.data.Minute
}

func (zt *grpcZonedTime) GetSec() uint32 {
	return zt.data.Sec
}

func (zt *grpcZonedTime) GetMicrosec() uint32 {
	return zt.data.Microsec
}

func (zt *grpcZonedTime) GetOffset() int {
	return int(zt.data.GetOffset())
}

func (zdt *grpcZonedDatetime) String() string {
	offset := zdt.GetOffset()
	var zone string
	if offset < 0 {
		zone = fmt.Sprintf("-%02d:%02d", -offset/3600, (-offset%3600)/60)
	} else if offset > 0 {
		zone = fmt.Sprintf("+%02d:%02d", offset/3600, (offset%3600)/60)
	} else {
		zone = "Z"
	}
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%06d%s",
		zdt.data.Year, zdt.data.Month, zdt.data.Day,
		zdt.data.Hour, zdt.data.Minute, zdt.data.Sec, zdt.data.Microsec,
		zone)
}

func (zdt *grpcZonedDatetime) GetOffset() int {
	return int(zdt.data.GetOffset())
}

func (zdt *grpcZonedDatetime) Time() *time.Time {
	return proto.ConvertZonedTime(zdt.data)
}

func (zdt *grpcZonedDatetime) GetYear() int32 {
	return zdt.data.Year
}

func (zdt *grpcZonedDatetime) GetMonth() uint32 {
	return zdt.data.Month
}

func (zdt *grpcZonedDatetime) GetDay() uint32 {
	return zdt.data.Day
}

func (zdt *grpcZonedDatetime) GetHour() uint32 {
	return zdt.data.Hour
}

func (zdt *grpcZonedDatetime) GetMinute() uint32 {
	return zdt.data.Minute
}

func (zdt *grpcZonedDatetime) GetSec() uint32 {
	return zdt.data.Sec
}

func (zdt *grpcZonedDatetime) GetMicrosec() uint32 {
	return zdt.data.Microsec
}
