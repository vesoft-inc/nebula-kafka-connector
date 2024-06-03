package printer

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

type stringFunc func(nebula.Value) string

type valueStringer interface {
	String(nebula.Value) string
}

type defaultStringer struct{}
type tckValueStringer struct {
	defaultStringer
}

func newValueStringer(format string) valueStringer {
	switch format {
	case "tck":
		return &tckValueStringer{}
	default:
		return &defaultStringer{}
	}
}

// used for record value
// if the value is string and contain special character, should quote it
func (s *defaultStringer) strQuoteSpecial(str string) string {
	subStr := []string{"\"", "\a", "\b", "\f", "\n", "\r", "\t", "\v", " "}
	for _, c := range subStr {
		if strings.Contains(str, c) {
			return strconv.Quote(str)
		}
	}
	return str
}

// if print string without quote, would quote key and string value when the string
// value contains special chars. e.g. {"ab c":hello}
// if print string with quote, quote any string value. e.g. {"ab c":"hello"}
func (s *defaultStringer) mapString(m map[string]nebula.Value, strWithQuote bool, strFn stringFunc) string {
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
		var (
			key   string
			value string
		)
		if strWithQuote {
			key = strconv.Quote(k)
		} else {
			key = s.strQuoteSpecial(k)
		}
		if v.GetType() == nebula.ValueTypeString {
			if strWithQuote {
				value = strconv.Quote(v.String())
			} else {
				value = s.strQuoteSpecial(v.String())
			}
		} else {
			value = strFn(v)
		}
		kvStr = append(kvStr, fmt.Sprintf("%s:%s", key, value))
		keys = append(keys, key)
	}
	return strings.Join(kvStr, ",")
}

func (s *defaultStringer) String(v nebula.Value) string {
	switch v.GetType() {
	case nebula.ValueTypeString:
		d, _ := v.AsString()
		return d.String()
	case nebula.ValueTypeList, nebula.ValueTypeRecord, nebula.ValueTypeNode, nebula.ValueTypeEdge, nebula.ValueTypePath:
		return s.complexString(v, false, s.String)
	default:
		return v.String()
	}
}

func (s *defaultStringer) complexString(v nebula.Value, strWithQuote bool, strFn stringFunc) string {
	switch v.GetType() {
	case nebula.ValueTypeList:
		d, _ := v.AsList()
		buf := make([]string, 0)
		for _, v := range d.GetValues() {
			if v.GetType() == nebula.ValueTypeString {
				if strWithQuote {
					buf = append(buf, strconv.Quote(v.String()))
				} else {
					buf = append(buf, s.strQuoteSpecial(v.String()))
				}

			} else {
				buf = append(buf, strFn(v))
			}
		}
		return fmt.Sprintf("[%s]", strings.Join(buf, ","))
	case nebula.ValueTypeRecord:
		d, _ := v.AsRecord()
		values := d.GetValues()
		return fmt.Sprintf("{%s}", s.mapString(values, strWithQuote, strFn))
	case nebula.ValueTypeNode:
		d, _ := v.AsNode()
		properties := d.GetProperties()
		return fmt.Sprintf("(%d@%s:%s{%s})",
			d.GetId(),
			d.GetType(),
			strings.Join(d.GetLabels(), "&"),
			s.mapString(properties, strWithQuote, strFn),
		)
	case nebula.ValueTypeEdge:
		d, _ := v.AsEdge()
		properties := d.GetProperties()
		var (
			leftBracket  string
			rightBracket string
		)
		if d.IsDirected() {
			leftBracket = "-"
			rightBracket = "->"
		} else {
			leftBracket = "~"
			rightBracket = "~"
		}

		return fmt.Sprintf("(%d)%s[%d@%s:%s{%s}]%s(%d)",
			d.GetSrcId(),
			leftBracket,
			d.GetRank(),
			d.GetType(),
			strings.Join(d.GetLabels(), "&"),
			s.mapString(properties, false, strFn),
			rightBracket,
			d.GetDstId(),
		)
	case nebula.ValueTypePath:
		p, _ := v.AsPath()
		values := p.GetValues()
		buf := bytes.NewBuffer(nil)
		defer buf.Reset()
		var preN nebula.Node
		for _, v := range values {
			if v.GetType() == nebula.ValueTypeNode {
				n, _ := v.AsNode()
				preN = n
				buf.WriteString(n.String())
			} else if v.GetType() == nebula.ValueTypeEdge {
				e, _ := v.AsEdge()
				p := e.GetProperties()
				estr := fmt.Sprintf("[%d@%s:%s{%s}]",
					e.GetRank(),
					e.GetType(),
					strings.Join(e.GetLabels(), "&"),
					s.mapString(p, false, strFn),
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

	default:
		return v.String()
	}
}

func (s *tckValueStringer) String(v nebula.Value) string {
	switch v.GetType() {
	case nebula.ValueTypeString:
		return fmt.Sprintf("\"%s\"", v.String())
	case nebula.ValueTypeDate:
		return fmt.Sprintf("DATE \"%s\"", v.String())
	case nebula.ValueTypeLocalTime:
		return fmt.Sprintf("TIME \"%s\"", v.String())
	case nebula.ValueTypeLocalDateTime:
		return fmt.Sprintf("DATETIME \"%s\"", v.String())
	case nebula.ValueTypeZonedTime:
		zt, _ := v.AsZonedTime()
		loc := time.FixedZone("", zt.GetOffset())
		data := time.Date(
			0, 0, 0,
			int(zt.GetHour()), int(zt.GetMinute()), int(zt.GetSec()),
			int(zt.GetMicrosec())*int(time.Microsecond), loc,
		)
		tm := data.UTC()
		return fmt.Sprintf("ZONED TIME \"%s\"", tm.Format("15:04:05.000000"))
	case nebula.ValueTypeZonedDateTime:
		zt, _ := v.AsZonedDatetime()
		loc := time.FixedZone("", zt.GetOffset())
		data := time.Date(
			int(zt.GetYear()), time.Month(zt.GetMonth()), int(zt.GetDay()),
			int(zt.GetHour()), int(zt.GetMinute()), int(zt.GetSec()),
			int(zt.GetMicrosec())*int(time.Microsecond), loc,
		)
		tm := data.UTC()
		return fmt.Sprintf("ZONED DATETIME \"%s\"", tm.Format("2006-01-02T15:04:05.000000"))
	case nebula.ValueTypeList:
		return fmt.Sprintf("LIST %s", s.complexString(v, true, s.String))
	case nebula.ValueTypeRecord:
		return fmt.Sprintf("RECORD %s", s.complexString(v, true, s.String))
	case nebula.ValueTypeNode, nebula.ValueTypeEdge, nebula.ValueTypePath:
		return s.complexString(v, true, s.String)
	default:
		return s.defaultStringer.String(v)
	}
}
