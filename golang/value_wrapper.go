// Copyright (c) 2022 vesoft inc. All rights reserved.

package nebula_ng_go

import (
	"fmt"
	"strconv"
	"strings"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/generated_code/v5.0.0/nebula"
)

type ValueWrapper struct {
	value        nebula.Value
	timezoneInfo timezoneInfo
}

func (valWrap ValueWrapper) IsNull() bool {
	return valWrap.GetType() == "Null"
}

func (valWrap ValueWrapper) IsBool() bool {
	return valWrap.value.IsSetBoolVal()
}

func (valWrap ValueWrapper) IsInt8() bool {
	return valWrap.value.IsSetInt8Val()
}

func (valWrap ValueWrapper) IsInt16() bool {
	return valWrap.value.IsSetInt16Val()
}

func (valWrap ValueWrapper) IsInt32() bool {
	return valWrap.value.IsSetInt32Val()
}

func (valWrap ValueWrapper) IsInt64() bool {
	return valWrap.value.IsSetInt64Val()
}

func (valWrap ValueWrapper) IsFloat() bool {
	return valWrap.value.IsSetFloatVal()
}

func (valWrap ValueWrapper) IsString() bool {
	return valWrap.value.IsSetStringVal()
}

func (valWrap ValueWrapper) IsList() bool {
	return valWrap.value.IsSetListVal()
}

func (valWrap ValueWrapper) IsRecord() bool {
	return valWrap.value.IsSetRecordVal()
}

func (valWrap ValueWrapper) IsNode() bool {
	return valWrap.value.IsSetNodeVal()
}

func (valWrap ValueWrapper) IsEdge() bool {
	return valWrap.value.IsSetEdgeVal()
}

func (valWrap ValueWrapper) IsLocalTime() bool {
	return valWrap.value.IsSetLocalTimeVal()
}

func (valWrap ValueWrapper) IsDate() bool {
	return valWrap.value.IsSetDateVal()
}

func (valWrap ValueWrapper) IsLocalDatetime() bool {
	return valWrap.value.IsSetLocalDatetimeVal()
}

func (valWrap ValueWrapper) IsDuration() bool {
	return valWrap.value.IsSetDurationVal()
}

// AsBool converts the ValueWrapper to a boolean value
func (valWrap ValueWrapper) AsBool() (bool, error) {
	if valWrap.value.IsSetBoolVal() {
		return valWrap.value.GetBoolVal(), nil
	}
	return false, fmt.Errorf("failed to convert value %s to bool", valWrap.GetType())
}

// AsInt8 converts the ValueWrapper to an int64
func (valWrap ValueWrapper) AsInt8() (int8, error) {
	if valWrap.value.IsSetInt8Val() {
		return valWrap.value.GetInt8Val(), nil
	}
	return -1, fmt.Errorf("failed to convert value %s to int", valWrap.GetType())
}

// AsInt16 converts the ValueWrapper to an int16
func (valWrap ValueWrapper) AsInt16() (int16, error) {
	if valWrap.value.IsSetInt16Val() {
		return valWrap.value.GetInt16Val(), nil
	}
	return -1, fmt.Errorf("failed to convert value %s to int", valWrap.GetType())
}

// AsInt32 converts the ValueWrapper to an int32
func (valWrap ValueWrapper) AsInt32() (int32, error) {
	if valWrap.value.IsSetInt32Val() {
		return valWrap.value.GetInt32Val(), nil
	}
	return -1, fmt.Errorf("failed to convert value %s to int", valWrap.GetType())
}

// AsInt64 converts the ValueWrapper to an int64
func (valWrap ValueWrapper) AsInt64() (int64, error) {
	if valWrap.value.IsSetInt64Val() {
		return valWrap.value.GetInt64Val(), nil
	}
	return -1, fmt.Errorf("failed to convert value %s to int", valWrap.GetType())
}

// AsFloat converts the ValueWrapper to a float64
func (valWrap ValueWrapper) AsFloat() (float64, error) {
	if valWrap.value.IsSetFloatVal() {
		return valWrap.value.GetFloatVal(), nil
	}
	return -1, fmt.Errorf("failed to convert value %s to float", valWrap.GetType())
}

// AsString converts the ValueWrapper to a String
func (valWrap ValueWrapper) AsString() (string, error) {
	if valWrap.value.IsSetStringVal() {
		return string(valWrap.value.GetStringVal()), nil
	}
	return "", fmt.Errorf("failed to convert value %s to string", valWrap.GetType())
}

// AsList converts the ValueWrapper to a slice of ValueWrapper
func (valWrap ValueWrapper) AsList() ([]ValueWrapper, error) {
	if valWrap.value.IsSetListVal() {
		var varList []ValueWrapper
		vals := valWrap.value.GetListVal().Values
		for _, val := range vals {
			varList = append(varList, ValueWrapper{val, valWrap.timezoneInfo})
		}
		return varList, nil
	}
	return nil, fmt.Errorf("failed to convert value %s to List", valWrap.GetType())
}

// AsDate converts the ValueWrapper to a nebula.Date
func (valWrap ValueWrapper) AsDate() (nebula.Date, error) {
	if valWrap.value.IsSetDateVal() {
		return valWrap.value.GetDateVal(), nil
	}
	return nil, fmt.Errorf("failed to convert value %s to Date", valWrap.GetType())
}

// AsLocalTime converts the ValueWrapper to a nebula.LocalTime
func (valWrap ValueWrapper) AsLocalTime() (*LocalTimeWrapper, error) {
	if valWrap.value.IsSetLocalTimeVal() {
		rawTime := valWrap.value.GetLocalTimeVal()
		time, err := genLocalTimeWrapper(rawTime, valWrap.timezoneInfo)
		if err != nil {
			return nil, err
		}
		return time, nil
	}
	return nil, fmt.Errorf("failed to convert value %s to LocalTime", valWrap.GetType())
}

// AsLocalDatetime converts the ValueWrapper to a nebula.LocalDatetime
func (valWrap ValueWrapper) AsLocalDatetime() (nebula.LocalDatetime, error) {
	if valWrap.value.IsSetLocalDatetimeVal() {
		return valWrap.value.GetLocalDatetimeVal(), nil
	}
	return nil, fmt.Errorf("failed to convert value %s to LocalDatetime", valWrap.GetType())
}

// AsDuration converts the ValueWrapper to a nebula.Duration
func (valWrap ValueWrapper) AsDuration() (nebula.Duration, error) {
	if valWrap.value.IsSetDurationVal() {
		return valWrap.value.GetDurationVal(), nil
	}
	return nil, fmt.Errorf("failed to convert value %s to Duration", valWrap.GetType())
}

// // AsMap converts the ValueWrapper to a map of string and ValueWrapper
// func (valWrap ValueWrapper) AsMap() (map[string]ValueWrapper, error) {
// 	if valWrap.value.IsSetMapVal() {
// 		newMap := make(map[string]ValueWrapper)

// 		kvs := valWrap.value.GetMapVal().values
// 		for key, val := range kvs {
// 			newMap[key] = ValueWrapper{val, valWrap.timezoneInfo}
// 		}
// 		return newMap, nil
// 	}
// 	return nil, fmt.Errorf("failed to convert value %s to Map", valWrap.GetType())
// }

// AsNode converts the ValueWrapper to a Node
func (valWrap ValueWrapper) AsNode() (*Node, error) {
	if !valWrap.value.IsSetNodeVal() {
		return nil, fmt.Errorf("failed to convert value %s to Node, value is not an vertex", valWrap.GetType())
	}
	vertex := valWrap.value.NodeVal
	node, err := genNode(vertex, valWrap.timezoneInfo)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// AsEdge converts the ValueWrapper to an Edge
func (valWrap ValueWrapper) AsEdge() (*Edge, error) {
	if !valWrap.value.IsSetEdgeVal() {
		return nil, fmt.Errorf("failed to convert value %s to Edge, value is not an edge", valWrap.GetType())
	}
	edge := valWrap.value.EdgeVal
	return genEdge(edge, valWrap.timezoneInfo)
}

// GetType returns the value type of value in the valWrap as a string
func (valWrap ValueWrapper) GetType() string {
	if valWrap.value.IsSetBoolVal() {
		return "bool"
	} else if valWrap.value.IsSetInt8Val() {
		return "int8"
	} else if valWrap.value.IsSetInt16Val() {
		return "int16"
	} else if valWrap.value.IsSetInt32Val() {
		return "int32"
	} else if valWrap.value.IsSetInt64Val() {
		return "int64"
	} else if valWrap.value.IsSetFloatVal() {
		return "float"
	} else if valWrap.value.IsSetDoubleVal() {
		return "double"
	} else if valWrap.value.IsSetStringVal() {
		return "string"
	} else if valWrap.value.IsSetListVal() {
		return "list"
	} else if valWrap.value.IsSetRecordVal() {
		return "record"
	} else if valWrap.value.IsSetNodeVal() {
		return "node"
	} else if valWrap.value.IsSetEdgeVal() {
		return "edge"
	} else if valWrap.value.IsSetDateVal() {
		return "date"
	} else if valWrap.value.IsSetLocalTimeVal() {
		return "localDime"
	} else if valWrap.value.IsSetLocalDatetimeVal() {
		return "localDatetime"
	} else if valWrap.value.IsSetDurationVal() {
		return "duration"
	} else {
		return "Null"
	}
}

// TODO(Aiee) Add String() method to ValueWrapper, now we use the original string representation of the value
// String() returns the value in the ValueWrapper as a string.
func (valWrap ValueWrapper) String() string {
	value := valWrap.value
	if value.IsSetBoolVal() {
		return fmt.Sprintf("%t", value.GetBoolVal())
	} else if value.IsSetInt8Val() { // Integer
		return fmt.Sprintf("%d", value.GetInt8Val())
	} else if value.IsSetInt16Val() {
		return fmt.Sprintf("%d", value.GetInt16Val())
	} else if value.IsSetInt32Val() {
		return fmt.Sprintf("%d", value.GetInt32Val())
	} else if value.IsSetInt64Val() {
		return fmt.Sprintf("%d", value.GetInt64Val())
	} else if value.IsSetDoubleVal() { // Double
		fStr := strconv.FormatFloat(value.GetDoubleVal(), 'g', -1, 64)
		if !strings.Contains(fStr, ".") {
			fStr = fStr + ".0"
		}
		return fStr
	} else if value.IsSetFloatVal() { // Float TODO(Aiee) Check if this is correct
		fStr := strconv.FormatFloat(value.GetFloatVal(), 'g', -1, 64)
		if !strings.Contains(fStr, ".") {
			fStr = fStr + ".0"
		}
		return fStr
	} else if value.IsSetStringVal() {
		return `"` + string(value.GetStringVal()) + `"`
	} else if value.IsSetNodeVal() { // Node
		rawNode := value.GetNodeVal()
		node, _ := genNode(rawNode, valWrap.timezoneInfo)
		return node.String()
	} else if value.IsSetEdgeVal() { // Edge
		rawEdge := value.GetEdgeVal()
		edge, _ := genEdge(rawEdge, valWrap.timezoneInfo)
		return edge.String()
	} else if value.IsSetListVal() { // List
		lval := value.GetListVal()
		var strs []string
		for _, val := range lval.Values {
			strs = append(strs, ValueWrapper{val, valWrap.timezoneInfo}.String())
		}
		return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
	} else if value.IsSetDateVal() { // Date yyyy-mm-dd
		date := value.GetDateVal()
		dateWrapper, _ := genDateWrapper(date)
		return fmt.Sprintf("%04d-%02d-%02d",
			dateWrapper.getYear(),
			dateWrapper.getMonth(),
			dateWrapper.getDay())
	} else if value.IsSetLocalTimeVal() { // Time HH:MM:SS.MSMSMS
		rawTime := value.GetLocalTimeVal()
		localTime, _ := genLocalTimeWrapper(rawTime, valWrap.timezoneInfo)
		return fmt.Sprintf("%02d:%02d:%02d.%06d",
			localTime.getHour(),
			localTime.getMinute(),
			localTime.getSecond(),
			localTime.getMicrosec())
	} else if value.IsSetLocalDatetimeVal() { // DateTime yyyy-mm-ddTHH:MM:SS.MSMSMS
		rawLocalDateTime := value.GetLocalDatetimeVal()
		localDateTime, _ := genLocalDatetimeWrapper(rawLocalDateTime, valWrap.timezoneInfo)
		return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.%06d",
			localDateTime.getYear(),
			localDateTime.getMonth(),
			localDateTime.getDay(),
			localDateTime.getHour(),
			localDateTime.getMinute(),
			localDateTime.getSecond(),
			localDateTime.getMicrosec())
	} else if value.IsSetDurationVal() { // Duration PnYnMnDTnHnMnS
		return "duration type string unimplemented"
	} else if value.IsSetRecordVal() { // Map TODO(Aiee) Unimplemented
		// // {k0: v0, k1: v1}
		// mval := value.GetMapVal()
		// var keyList []string
		// var output []string
		// kvs := mval.Kvs
		// for k := range kvs {
		// 	keyList = append(keyList, k)
		// }
		// sort.Strings(keyList)
		// for _, k := range keyList {
		// 	output = append(output, fmt.Sprintf("%s: %s", k, ValueWrapper{kvs[k], valWrap.timezoneInfo}.String()))
		// }
		// return fmt.Sprintf("{%s}", strings.Join(output, ", "))
		return "map type string unimplemented"
	} else { // Null
		return "__NULL__"
	}
}
