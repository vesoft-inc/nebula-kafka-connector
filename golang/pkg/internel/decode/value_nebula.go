package decode

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	internel_error "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/internal_error"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/types"
)

type valuer interface {
	GetValue()
}

type (
	mapValue    map[string]types.Value
	nebulaValue struct {
		data valuer
	}
	commonValue struct{}
	NebulaBool  struct {
		commonValue
		Value bool
	}
	NebulaInt8 struct {
		commonValue
		Value int8
	}
	NebulaInt16 struct {
		commonValue
		Value int16
	}
	NebulaInt32 struct {
		commonValue
		Value int32
	}
	NebulaInt64 struct {
		commonValue
		Value int64
	}
	NebulaUint8 struct {
		commonValue
		Value uint8
	}
	NebulaUint16 struct {
		commonValue
		Value uint16
	}
	NebulaUint32 struct {
		commonValue
		Value uint32
	}
	NebulaUint64 struct {
		commonValue
		Value uint64
	}
	NebulaFloat struct {
		commonValue
		Value float32
	}
	NebulaDouble struct {
		commonValue
		Value float64
	}
	NebulaString struct {
		commonValue
		Value string
	}

	NebulaLocalTime struct {
		commonValue
		Hour     int8
		Minute   int8
		Sec      int8
		Microsec int32
	}
	NebulaLocalDatetime struct {
		commonValue
		Year     int16
		Month    int8
		Day      int8
		Hour     int8
		Minute   int8
		Sec      int8
		Microsec int32
	}
	NebulaZonedTime struct {
		commonValue
		Hour     int8
		Minute   int8
		Sec      int8
		Microsec int32
		Offset   int32
	}
	NebulaZonedDatetime struct {
		commonValue
		Year     int16
		Month    int8
		Day      int8
		Hour     int8
		Minute   int8
		Sec      int8
		Microsec int32
		Offset   int32
	}
	NebulaDuration struct {
		commonValue
		isMonthBased bool
		Year         int64
		Month        int8
		Day          int32
		Hour         int8
		Minute       int8
		Sec          int8
		Microsec     int32
	}
	NebulaDate struct {
		commonValue
		Year  int16
		Month int8
		Day   int8
	}
	NebulaList struct {
		commonValue
		Values []*nebulaValue
	}
	NebulaRecord struct {
		commonValue
		Values map[string]*nebulaValue
	}
	NebulaNode struct {
		commonValue
		NodeId     int64
		Graph      string
		Type       string
		Labels     []string
		Properties map[string]*nebulaValue
	}
	NebulaEdge struct {
		commonValue
		SrcId      int64
		DstId      int64
		Direction  int
		Graph      string
		Type       string
		Labels     []string
		Rank       int64
		Properties map[string]*nebulaValue
	}
	NebulaPath struct {
		commonValue
		Values []*nebulaValue
	}
	NebulaDecimal struct {
		commonValue
		Sval string
	}
)

func (v *nebulaValue) GetValue() {}
func (v *commonValue) GetValue() {}

func (v *nebulaValue) String() string {
	switch v.GetType() {
	case types.ValueTypeNull:
		return "null"
	case types.ValueTypeBool:
		d, _ := v.AsBool()
		return fmt.Sprintf("%t", d)
	case types.ValueTypeInt8:
		d, _ := v.AsInt8()
		return fmt.Sprintf("%d", d)
	case types.ValueTypeInt16:
		d, _ := v.AsInt16()
		return fmt.Sprintf("%d", d)
	case types.ValueTypeInt32:
		d, _ := v.AsInt32()
		return fmt.Sprintf("%d", d)
	case types.ValueTypeInt64:
		d, _ := v.AsInt64()
		return fmt.Sprintf("%d", d)
	case types.ValueTypeUInt8:
		d, _ := v.AsUInt8()
		return fmt.Sprintf("%d", d)
	case types.ValueTypeUInt16:
		d, _ := v.AsUInt16()
		return fmt.Sprintf("%d", d)
	case types.ValueTypeUInt32:
		d, _ := v.AsUInt32()
		return fmt.Sprintf("%d", d)
	case types.ValueTypeUInt64:
		d, _ := v.AsUInt64()
		return fmt.Sprintf("%d", d)
	case types.ValueTypeFloat:
		d, _ := v.AsFloat()
		fStr := strconv.FormatFloat(float64(d), 'g', -1, 32)
		if !strings.Contains(fStr, ".") {
			fStr = fStr + ".0"
		}
		return fStr
	case types.ValueTypeDouble:
		d, _ := v.AsDouble()
		fStr := strconv.FormatFloat(float64(d), 'g', -1, 64)
		if !strings.Contains(fStr, ".") {
			fStr = fStr + ".0"
		}
		return fStr
	case types.ValueTypeString:
		s, _ := v.AsString()
		return s.String()
	case types.ValueTypeDuration:
		d, _ := v.AsDuration()
		return d.String()
	case types.ValueTypeDate:
		d, _ := v.AsDate()
		return d.String()
	case types.ValueTypeLocalDateTime:
		dt, _ := v.AsLocalDatetime()
		return dt.String()
	case types.ValueTypeLocalTime:
		t, _ := v.AsLocalTime()
		return t.String()
	case types.ValueTypeZonedTime:
		t, _ := v.AsZonedTime()
		return t.String()
	case types.ValueTypeZonedDateTime:
		dt, _ := v.AsZonedDatetime()
		return dt.String()
	case types.ValueTypeList:
		l, _ := v.AsList()
		return l.String()
	case types.ValueTypeRecord:
		r, _ := v.AsRecord()
		return r.String()
	case types.ValueTypeNode:
		n, _ := v.AsNode()
		return n.String()
	case types.ValueTypeEdge:
		e, _ := v.AsEdge()
		return e.String()
	case types.ValueTypePath:
		p, _ := v.AsPath()
		return p.String()
	case types.ValueTypeDecimal:
		d, _ := v.AsDecimal()
		return d.String()
	default:
		return fmt.Sprintf("%v", v.data)
	}
}

func (v *nebulaValue) GetType() types.ValueType {
	if v.data == nil {
		return types.ValueTypeNull
	}
	switch v.data.(type) {
	case *NebulaBool:
		return types.ValueTypeBool
	case *NebulaInt8:
		return types.ValueTypeInt8
	case *NebulaInt16:
		return types.ValueTypeInt16
	case *NebulaInt32:
		return types.ValueTypeInt32
	case *NebulaInt64:
		return types.ValueTypeInt64
	case *NebulaUint8:
		return types.ValueTypeUInt8
	case *NebulaUint16:
		return types.ValueTypeUInt16
	case *NebulaUint32:
		return types.ValueTypeUInt32
	case *NebulaUint64:
		return types.ValueTypeUInt64
	case *NebulaFloat:
		return types.ValueTypeFloat
	case *NebulaDouble:
		return types.ValueTypeDouble
	case *NebulaString:
		return types.ValueTypeString
	case *NebulaDuration:
		return types.ValueTypeDuration
	case *NebulaLocalTime:
		return types.ValueTypeLocalTime
	case *NebulaLocalDatetime:
		return types.ValueTypeLocalDateTime
	case *NebulaZonedTime:
		return types.ValueTypeZonedTime
	case *NebulaZonedDatetime:
		return types.ValueTypeZonedDateTime
	case *NebulaDate:
		return types.ValueTypeDate
	case *NebulaList:
		return types.ValueTypeList
	case *NebulaRecord:
		return types.ValueTypeRecord
	case *NebulaNode:
		return types.ValueTypeNode
	case *NebulaEdge:
		return types.ValueTypeEdge
	case *NebulaPath:
		return types.ValueTypePath
	case *NebulaDecimal:
		return types.ValueTypeDecimal
	default:
		return types.ValueUnSupport
	}
}

func (v *nebulaValue) IsNull() bool {
	return v.data == nil
}

func asValue[T valuer](v *nebulaValue, valueTyp types.ValueType) (T, error) {
	var t T
	if v.GetType() != valueTyp {
		return t, internel_error.ErrType("value is not " + valueTyp.String())
	}
	data, ok := v.data.(T)
	if !ok {
		return t, internel_error.ErrType("value is not " + valueTyp.String())
	}
	return data, nil
}

func (v *nebulaValue) AsBool() (types.Bool, error) {
	d, err := asValue[*NebulaBool](v, types.ValueTypeBool)
	if err != nil {
		return false, err
	}
	return types.Bool(d.Value), nil
}

func (v *nebulaValue) AsInt8() (types.Int8, error) {
	d, err := asValue[*NebulaInt8](v, types.ValueTypeInt8)
	if err != nil {
		return 0, err
	}
	return types.Int8(d.Value), nil
}

func (v *nebulaValue) AsInt16() (types.Int16, error) {
	d, err := asValue[*NebulaInt16](v, types.ValueTypeInt16)
	if err != nil {
		return 0, err
	}
	return types.Int16(d.Value), nil
}

func (v *nebulaValue) AsInt32() (types.Int32, error) {
	d, err := asValue[*NebulaInt32](v, types.ValueTypeInt32)
	if err != nil {
		return 0, err
	}
	return types.Int32(d.Value), nil
}

func (v *nebulaValue) AsInt64() (types.Int64, error) {
	d, err := asValue[*NebulaInt64](v, types.ValueTypeInt64)
	if err != nil {
		return 0, err
	}
	return types.Int64(d.Value), nil
}

func (v *nebulaValue) AsUInt8() (types.UInt8, error) {
	d, err := asValue[*NebulaUint8](v, types.ValueTypeUInt8)
	if err != nil {
		return 0, err
	}
	return types.UInt8(d.Value), nil
}

func (v *nebulaValue) AsUInt16() (types.UInt16, error) {
	d, err := asValue[*NebulaUint16](v, types.ValueTypeUInt16)
	if err != nil {
		return 0, err
	}
	return types.UInt16(d.Value), nil
}

func (v *nebulaValue) AsUInt32() (types.UInt32, error) {
	d, err := asValue[*NebulaUint32](v, types.ValueTypeUInt32)
	if err != nil {
		return 0, err
	}
	return types.UInt32(d.Value), nil
}

func (v *nebulaValue) AsUInt64() (types.UInt64, error) {
	d, err := asValue[*NebulaUint64](v, types.ValueTypeUInt64)
	if err != nil {
		return 0, err
	}
	return types.UInt64(d.Value), nil
}

func (v *nebulaValue) AsFloat() (types.Float, error) {
	d, err := asValue[*NebulaFloat](v, types.ValueTypeFloat)
	if err != nil {
		return 0, err
	}
	return types.Float(d.Value), nil
}

func (v *nebulaValue) AsDouble() (types.Double, error) {
	d, err := asValue[*NebulaDouble](v, types.ValueTypeDouble)
	if err != nil {
		return 0, err
	}
	return types.Double(d.Value), nil
}

func (v *nebulaValue) AsString() (types.String, error) {
	d, err := asValue[*NebulaString](v, types.ValueTypeString)
	if err != nil {
		return "", err
	}
	return types.String(d.Value), nil
}

func (v *nebulaValue) AsList() (types.List, error) {
	return asValue[*NebulaList](v, types.ValueTypeList)
}

func (v *nebulaValue) AsRecord() (types.Record, error) {
	return asValue[*NebulaRecord](v, types.ValueTypeRecord)
}

func (v *nebulaValue) AsDuration() (types.Duration, error) {
	return asValue[*NebulaDuration](v, types.ValueTypeDuration)
}

func (v *nebulaValue) AsNode() (types.Node, error) {
	return asValue[*NebulaNode](v, types.ValueTypeNode)
}

func (v *nebulaValue) AsEdge() (types.Edge, error) {
	return asValue[*NebulaEdge](v, types.ValueTypeEdge)
}

func (v *nebulaValue) AsPath() (types.Path, error) {
	return asValue[*NebulaPath](v, types.ValueTypePath)
}

func (v *nebulaValue) AsDecimal() (types.Decimal, error) {
	return asValue[*NebulaDecimal](v, types.ValueTypeDecimal)
}

func (l *NebulaList) String() string {
	valuesStr := make([]string, 0, len(l.Values))
	for _, v := range l.GetValues() {
		valuesStr = append(valuesStr, v.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(valuesStr, ","))
}

func (v *nebulaValue) AsLocalDatetime() (types.LocalDatetime, error) {
	return asValue[*NebulaLocalDatetime](v, types.ValueTypeLocalDateTime)
}

func (v *nebulaValue) AsDate() (types.Date, error) {
	return asValue[*NebulaDate](v, types.ValueTypeDate)
}

func (v *nebulaValue) AsLocalTime() (types.LocalTime, error) {
	return asValue[*NebulaLocalTime](v, types.ValueTypeLocalTime)
}

func (v *nebulaValue) AsZonedTime() (types.ZonedTime, error) {
	return asValue[*NebulaZonedTime](v, types.ValueTypeZonedTime)
}

func (v *nebulaValue) AsZonedDatetime() (types.ZonedDatetime, error) {
	return asValue[*NebulaZonedDatetime](v, types.ValueTypeZonedDateTime)
}

func (l *NebulaList) GetValues() []types.Value {
	values := make([]types.Value, 0, len(l.Values))
	for _, v := range l.Values {
		values = append(values, v)
	}
	return values
}

func (l *NebulaList) Size() int {
	return len(l.Values)
}

func (r *NebulaRecord) String() string {
	mv := mapValue(r.GetValues())
	return fmt.Sprintf("{%s}", mv.string())
}

func (r *NebulaRecord) GetValues() map[string]types.Value {
	values := make(map[string]types.Value)
	for k, v := range r.Values {
		values[k] = v
	}
	return values
}

func (g *NebulaDecimal) String() string {
	return g.Sval
}

// (288314845273522179@City:City&Place{id:32,name:Norway,url:http://dbpedia.org/resource/Norway})
func (n *NebulaNode) String() string {
	mv := mapValue(n.GetProperties())
	return fmt.Sprintf("(%d@%s:%s{%s})",
		n.GetId(),
		n.GetType(),
		strings.Join(n.GetLabels(), "&"),
		mv.string(),
	)
}

func (n *NebulaNode) GetProperties() map[string]types.Value {
	properties := make(map[string]types.Value)
	for k, v := range n.Properties {
		properties[k] = v
	}
	return properties
}

func (n *NebulaNode) GetId() int {
	return int(n.NodeId)
}

func (n *NebulaNode) GetGraph() string {
	return n.Graph
}

func (n *NebulaNode) GetType() string {
	return n.Type
}

func (n *NebulaNode) GetLabels() []string {
	return n.Labels
}

// (288314845273522179)<-[288314845273522179@connected_Sub_Load:
// connected_Sub_Load{cid:115967690673232363}]-(288314982712475649)
func (e *NebulaEdge) String() string {
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

func (e *NebulaEdge) GetProperties() map[string]types.Value {
	properties := make(map[string]types.Value)
	for k, v := range e.Properties {
		properties[k] = v
	}
	return properties
}

func (e *NebulaEdge) GetSrcId() int {
	return int(e.SrcId)
}

func (e *NebulaEdge) GetDstId() int {
	return int(e.DstId)
}

func (e *NebulaEdge) GetGraph() string {
	return e.Graph
}

func (e *NebulaEdge) GetType() string {
	return e.Type
}

func (e *NebulaEdge) GetLabels() []string {
	return e.Labels
}

func (e *NebulaEdge) GetRank() int {
	return int(e.Rank)
}

func (e *NebulaEdge) IsDirected() bool {
	return e.Direction == 0
}

// (288314845273522179@City:City&Place{id:3,kind:3,name:org3,url:https://org3.com})-[288314845273522179@connected_Sub_Load:connected_Sub_Load{}]-(288315214640709633@City:City&Place{id:3,kind:city,name:Hangzhou,url:https://hangzhou.com})
// -[288315214640709633@connected_Sub_Load:connected_Sub_Load{}]
// -(288315309129990145@City:City&Place{id:6,kind:city,name:Shenzhen,url:https://shenzhen.com})
func (p *NebulaPath) String() string {
	values := p.GetValues()
	buf := bytes.NewBuffer(nil)
	defer buf.Reset()
	var preN types.Node
	for _, v := range values {
		if v.GetType() == types.ValueTypeNode {
			n, _ := v.AsNode()
			preN = n
			buf.WriteString(n.String())
		} else if v.GetType() == types.ValueTypeEdge {
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

func (p *NebulaPath) GetValues() []types.Value {
	values := make([]types.Value, 0, len(p.Values))
	for _, v := range p.Values {
		values = append(values, v)
	}
	return values
}

func (l *NebulaLocalDatetime) String() string {
	//RFC3339 without timezone
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%06d",
		l.Year, l.Month, l.Day,
		l.Hour, l.Minute, l.Sec, l.Microsec)
}

func (l *NebulaLocalDatetime) GetYear() int {
	return int(l.Year)
}

func (l *NebulaLocalDatetime) GetMonth() int {
	return int(l.Month)
}

func (l *NebulaLocalDatetime) GetDay() int {
	return int(l.Day)
}

func (l *NebulaLocalDatetime) GetHour() int {
	return int(l.Hour)
}

func (l *NebulaLocalDatetime) GetMinute() int {
	return int(l.Minute)
}

func (l *NebulaLocalDatetime) GetSec() int {
	return int(l.Sec)
}

func (l *NebulaLocalDatetime) GetMicrosec() int {
	return int(l.Microsec)
}

func (d *NebulaDate) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

func (d *NebulaDate) GetYear() int {
	return int(d.Year)
}

func (d *NebulaDate) GetMonth() int {
	return int(d.Month)
}

func (d *NebulaDate) GetDay() int {
	return int(d.Day)
}

func (t *NebulaLocalTime) String() string {
	return fmt.Sprintf("%02d:%02d:%02d.%06d",
		t.Hour, t.Minute, t.Sec, t.Microsec)
}

func (t *NebulaLocalTime) GetHour() int {
	return int(t.Hour)
}

func (t *NebulaLocalTime) GetMinute() int {
	return int(t.Minute)
}

func (t *NebulaLocalTime) GetSec() int {
	return int(t.Sec)
}

func (t *NebulaLocalTime) GetMicrosec() int {
	return int(t.Microsec)
}

func (d *NebulaDuration) String() string {
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

func (d *NebulaDuration) IsMonthBased() bool {
	return d.isMonthBased
}

func (d *NebulaDuration) GetYear() int {
	return int(d.Year)
}

func (d *NebulaDuration) GetMonth() int {
	return int(d.Month)
}

func (d *NebulaDuration) GetDay() int {
	return int(d.Day)
}

func (d *NebulaDuration) GetHour() int {
	return int(d.Hour)
}

func (d *NebulaDuration) GetMinute() int {
	return int(d.Minute)
}

func (d *NebulaDuration) GetSecond() int {
	return int(d.Sec)
}

func (d *NebulaDuration) GetMicrosecond() int {
	return int(d.Microsec)
}

func (zt *NebulaZonedTime) String() string {
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
		zt.Hour,
		zt.Minute,
		zt.Sec,
		zt.Microsec,
		zone,
	)
}

func (zt *NebulaZonedTime) GetHour() int {
	return int(zt.Hour)
}

func (zt *NebulaZonedTime) GetMinute() int {
	return int(zt.Minute)
}

func (zt *NebulaZonedTime) GetSec() int {
	return int(zt.Sec)
}

func (zt *NebulaZonedTime) GetMicrosec() int {
	return int(zt.Microsec)
}

func (zt *NebulaZonedTime) GetOffset() int {
	return int(zt.Offset)
}

func (zdt *NebulaZonedDatetime) String() string {
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
		zdt.Year, zdt.Month, zdt.Day,
		zdt.Hour, zdt.Minute, zdt.Sec, zdt.Microsec,
		zone)
}

func (zdt *NebulaZonedDatetime) GetOffset() int {
	return int(zdt.Offset)
}

func (zdt *NebulaZonedDatetime) Time() *time.Time {
	if zdt == nil {
		return nil
	}
	timezone := time.FixedZone("", int(zdt.GetOffset()))
	t := time.Date(
		int(zdt.GetYear()),
		time.Month(zdt.GetMonth()),
		int(zdt.GetDay()),
		int(zdt.GetHour()),
		int(zdt.GetMinute()),
		int(zdt.GetSec()),
		int(zdt.GetMicrosec())*int(time.Microsecond),
		timezone,
	)
	return &t
}

func (zdt *NebulaZonedDatetime) GetYear() int {
	return int(zdt.Year)
}

func (zdt *NebulaZonedDatetime) GetMonth() int {
	return int(zdt.Month)
}

func (zdt *NebulaZonedDatetime) GetDay() int {
	return int(zdt.Day)
}

func (zdt *NebulaZonedDatetime) GetHour() int {
	return int(zdt.Hour)
}

func (zdt *NebulaZonedDatetime) GetMinute() int {
	return int(zdt.Minute)
}

func (zdt *NebulaZonedDatetime) GetSec() int {
	return int(zdt.Sec)
}

func (zdt *NebulaZonedDatetime) GetMicrosec() int {
	return int(zdt.Microsec)
}

func (m mapValue) string() string {
	var kvStr []string = make([]string, 0, len(m))
	var keys []string = make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	for _, k := range keys {
		v := m[k]
		kvTemp := fmt.Sprintf(`%s:%s`, k, v)
		kvStr = append(kvStr, kvTemp)
	}
	return strings.Join(kvStr, ",")
}
