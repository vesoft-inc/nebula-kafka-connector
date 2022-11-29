// Copyright (c) 2022 vesoft inc. All rights reserved.

package nebula_ng_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/generated_code/v5.0.0/nebula"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/generated_code/v5.0.0/nebula/graph"
)

type ResultSet struct {
	resp            *graph.ExecutionResponse
	columnNames     []string
	colNameIndexMap map[string]int
	timezoneInfo    timezoneInfo
}

// Returns a 2D array of strings representing the query result
// If resultSet.resp.data is nil, returns an empty 2D array
func (res ResultSet) AsStringTable() [][]string {
	var resTable [][]string
	colNames := res.GetColNames()
	resTable = append(resTable, colNames)
	rows := res.GetRows()
	for _, row := range rows {
		var tempRow []string
		for _, val := range row.Values {
			tempRow = append(tempRow, ValueWrapper{val, res.timezoneInfo}.String())
		}
		resTable = append(resTable, tempRow)
	}
	return resTable
}

func (res ResultSet) GetStatus() string {
	return string(res.resp.ExecutionOutcome.GqlStatus.Status)
}

func (res ResultSet) IsSucceed() bool {
	return res.GetStatus() == string("SUCCESS")
}

func (res ResultSet) IsSetData() bool {
	return res.resp.ExecutionOutcome.Result_ != nil
}

// Returns all rows
func (res ResultSet) GetRows() []*nebula.RawRecord {
	if res.resp.ExecutionOutcome.Result_ == nil {
		var empty []*nebula.RawRecord
		return empty
	}
	return res.resp.ExecutionOutcome.Result_.Records
}

func (res ResultSet) GetRowSize() int {
	return len(res.GetRows())
}

func (res ResultSet) GetColNames() []string {
	return res.columnNames
}

func (res ResultSet) GetColSize() int {
	return len(res.GetColNames())
}

func (res ResultSet) GetLatency() int64 {
	return res.resp.LatencyInUs
}

func checkIndex(index int, list interface{}) error {
	if _, ok := list.([]*nebula.RawRecord); ok {
		if index < 0 || index >= len(list.([]*nebula.RawRecord)) {
			return fmt.Errorf("failed to get Value, the index is out of range")
		}
		return nil
	} else if _, ok := list.([]*ValueWrapper); ok {
		if index < 0 || index >= len(list.([]*ValueWrapper)) {
			return fmt.Errorf("failed to get Value, the index is out of range")
		}
		return nil
	}
	return fmt.Errorf("given list type is invalid")
}

// Returns all values in the row at given index
func (res ResultSet) GetRowValuesByIndex(index int) (*Record, error) {
	if err := checkIndex(index, res.resp.ExecutionOutcome.Result_.Records); err != nil {
		return nil, err
	}
	valWrap, err := genValWraps(res.resp.ExecutionOutcome.Result_.Records[index], res.timezoneInfo)
	if err != nil {
		return nil, err
	}
	return &Record{
		columnNames:     &res.columnNames,
		_record:         valWrap,
		colNameIndexMap: &res.colNameIndexMap,
		timezoneInfo:    res.timezoneInfo,
	}, nil
}

func (res ResultSet) IsSetPlanDesc() bool {
	return res.resp.ExecutionOutcome.PlanDesc != nil
}

func (res ResultSet) GetPlanDesc() *graph.PlanDescription {
	return res.resp.ExecutionOutcome.PlanDesc
}

type Record struct {
	columnNames     *[]string
	_record         []*ValueWrapper
	colNameIndexMap *map[string]int
	timezoneInfo    timezoneInfo
}

// Returns value in the record at given column index
func (record Record) GetValueByIndex(index int) (*ValueWrapper, error) {
	if err := checkIndex(index, record._record); err != nil {
		return nil, err
	}
	return record._record[index], nil
}

type Node struct {
	rawNode      *nebula.Node
	timezoneInfo timezoneInfo
}

// TODO(Aiee) For demo only
// The properties of the node are sorted by key
func (node *Node) String() string {
	var keyList []string
	var kvStr []string
	for key, _ := range node.rawNode.Properties {
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)
	for _, key := range keyList {
		kvTemp := fmt.Sprintf(`%s:%s`,
			key, ValueWrapper{node.rawNode.Properties[key], node.timezoneInfo}.String())
		kvStr = append(kvStr, kvTemp)
	}

	// vid and tag are internal values now, so do not print
	return fmt.Sprintf("({%s})", strings.Join(kvStr, ","))
}

type Edge struct {
	rawEdge      *nebula.Edge
	timezoneInfo timezoneInfo
}

// TODO(Aiee) For demo only
// The properties of the edge are sorted by key
func (edge *Edge) String() string {
	var keyList []string
	var kvStr []string
	for key, _ := range edge.rawEdge.Properties {
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)
	for _, key := range keyList {
		kvTemp := fmt.Sprintf(`%s:%s`,
			key, ValueWrapper{edge.rawEdge.Properties[key], edge.timezoneInfo}.String())
		kvStr = append(kvStr, kvTemp)
	}

	return fmt.Sprintf("[{%s}]", strings.Join(kvStr, ","))
}

type ErrorCode int64

const (
	ErrorCode_SUCCEEDED               ErrorCode = ErrorCode(nebula.ErrorCode_SUCCEEDED)
	ErrorCode_E_DISCONNECTED          ErrorCode = ErrorCode(nebula.ErrorCode_E_DISCONNECTED)
	ErrorCode_E_FAIL_TO_CONNECT       ErrorCode = ErrorCode(nebula.ErrorCode_E_FAIL_TO_CONNECT)
	ErrorCode_E_RPC_FAILURE           ErrorCode = ErrorCode(nebula.ErrorCode_E_RPC_FAILURE)
	ErrorCode_E_BAD_USERNAME_PASSWORD ErrorCode = ErrorCode(nebula.ErrorCode_E_BAD_USERNAME_PASSWORD)
	ErrorCode_E_SESSION_INVALID       ErrorCode = ErrorCode(nebula.ErrorCode_E_SESSION_INVALID)
	ErrorCode_E_SESSION_TIMEOUT       ErrorCode = ErrorCode(nebula.ErrorCode_E_SESSION_TIMEOUT)
	ErrorCode_E_SYNTAX_ERROR          ErrorCode = ErrorCode(nebula.ErrorCode_E_SYNTAX_ERROR)
	ErrorCode_E_EXECUTION_ERROR       ErrorCode = ErrorCode(nebula.ErrorCode_E_EXECUTION_ERROR)
	ErrorCode_E_STATEMENT_EMPTY       ErrorCode = ErrorCode(nebula.ErrorCode_E_STATEMENT_EMPTY)
	ErrorCode_E_USER_NOT_FOUND        ErrorCode = ErrorCode(nebula.ErrorCode_E_USER_NOT_FOUND)
	ErrorCode_E_BAD_PERMISSION        ErrorCode = ErrorCode(nebula.ErrorCode_E_BAD_PERMISSION)
	ErrorCode_E_SEMANTIC_ERROR        ErrorCode = ErrorCode(nebula.ErrorCode_E_SEMANTIC_ERROR)
	ErrorCode_E_PARTIAL_SUCCEEDED     ErrorCode = ErrorCode(nebula.ErrorCode_E_PARTIAL_SUCCEEDED)
)

type LocalTimeWrapper struct {
	time         *nebula.LocalTime
	timezoneInfo timezoneInfo
}

func genLocalTimeWrapper(time *nebula.LocalTime, timezoneInfo timezoneInfo) (*LocalTimeWrapper, error) {
	if time == nil {
		return nil, fmt.Errorf("failed to generate Time: invalid Time")
	}

	return &LocalTimeWrapper{
		time:         time,
		timezoneInfo: timezoneInfo,
	}, nil
}

// getHour returns the hour in UTC
func (t LocalTimeWrapper) getHour() int8 {
	return t.time.Hour
}

// getHour returns the minute in UTC
func (t LocalTimeWrapper) getMinute() int8 {
	return t.time.Minute
}

// getHour returns the second in UTC
func (t LocalTimeWrapper) getSecond() int8 {
	return t.time.Sec
}

func (t LocalTimeWrapper) getMicrosec() int32 {
	return t.time.Microsec
}

// getRawTime returns a nebula.Time object.
func (t LocalTimeWrapper) getRawTime() *nebula.LocalTime {
	return t.time
}

func (t1 LocalTimeWrapper) IsEqualTo(t2 LocalTimeWrapper) bool {
	return t1.getHour() == t2.getHour() &&
		t1.getSecond() == t2.getSecond() &&
		t1.getSecond() == t2.getSecond() &&
		t1.getMicrosec() == t2.getMicrosec()
}

type DateWrapper struct {
	date *nebula.Date
}

func genDateWrapper(date *nebula.Date) (*DateWrapper, error) {
	if date == nil {
		return nil, fmt.Errorf("failed to generate date: invalid date")
	}
	return &DateWrapper{
		date: date,
	}, nil
}

func (d DateWrapper) getYear() int16 {
	return d.date.Year
}

func (d DateWrapper) getMonth() int8 {
	return d.date.Month
}

func (d DateWrapper) getDay() int8 {
	return d.date.Day
}

// getRawDate returns a nebula.Date object.
func (d DateWrapper) getRawDate() *nebula.Date {
	return d.date
}

func (d1 DateWrapper) IsEqualTo(d2 DateWrapper) bool {
	return d1.getYear() == d2.getYear() &&
		d1.getMonth() == d2.getMonth() &&
		d1.getDay() == d2.getDay()
}

type LocalDatetimeWrapper struct {
	dateTime     *nebula.LocalDatetime
	timezoneInfo timezoneInfo
}

func genLocalDatetimeWrapper(datetime *nebula.LocalDatetime, timezoneInfo timezoneInfo) (*LocalDatetimeWrapper, error) {
	if datetime == nil {
		return nil, fmt.Errorf("failed to generate datetime: invalid datetime")
	}
	return &LocalDatetimeWrapper{
		dateTime:     datetime,
		timezoneInfo: timezoneInfo,
	}, nil
}

func (dt LocalDatetimeWrapper) getYear() int16 {
	return dt.dateTime.Year
}

func (dt LocalDatetimeWrapper) getMonth() int8 {
	return dt.dateTime.Month
}

func (dt LocalDatetimeWrapper) getDay() int8 {
	return dt.dateTime.Day
}

func (dt LocalDatetimeWrapper) getHour() int8 {
	return dt.dateTime.Hour
}

func (dt LocalDatetimeWrapper) getMinute() int8 {
	return dt.dateTime.Minute
}

func (dt LocalDatetimeWrapper) getSecond() int8 {
	return dt.dateTime.Sec
}

func (dt LocalDatetimeWrapper) getMicrosec() int32 {
	return dt.dateTime.Microsec
}

func (dt1 LocalDatetimeWrapper) IsEqualTo(dt2 LocalDatetimeWrapper) bool {
	return dt1.getYear() == dt2.getYear() &&
		dt1.getMonth() == dt2.getMonth() &&
		dt1.getDay() == dt2.getDay() &&
		dt1.getHour() == dt2.getHour() &&
		dt1.getSecond() == dt2.getSecond() &&
		dt1.getSecond() == dt2.getSecond() &&
		dt1.getMicrosec() == dt2.getMicrosec()
}

// getRawDateTime returns a nebula.DateTime object representing local dateTime in UTC.
func (dt LocalDatetimeWrapper) getRawDateTime() *nebula.LocalDatetime {
	return dt.dateTime
}

func genResultSet(resp *graph.ExecutionResponse, timezoneInfo timezoneInfo) (*ResultSet, error) {
	var colNames []string
	var colNameIndexMap = make(map[string]int)

	if resp.ExecutionOutcome.Result_ == nil { // if resp.Data != nil then resp.Data.row and resp.Data.colNames wont be nil
		return &ResultSet{
			resp:            resp,
			columnNames:     colNames,
			colNameIndexMap: colNameIndexMap,
		}, nil
	}

	// TODO(Aiee) Add binding table descriptor in thrift
	for i, name := range resp.ExecutionOutcome.Result_.ColumnNames {
		colNames = append(colNames, string(name))
		colNameIndexMap[string(name)] = i
	}

	return &ResultSet{
		resp:            resp,
		columnNames:     colNames,
		colNameIndexMap: colNameIndexMap,
		timezoneInfo:    timezoneInfo,
	}, nil
}

func genValWraps(row *nebula.RawRecord, timezoneInfo timezoneInfo) ([]*ValueWrapper, error) {
	if row == nil {
		return nil, fmt.Errorf("failed to generate valueWrapper: invalid row")
	}
	var valWraps []*ValueWrapper
	for _, val := range row.Values {
		if val == nil {
			return nil, fmt.Errorf("failed to generate valueWrapper: value is nil")
		}
		valWraps = append(valWraps, &ValueWrapper{val, timezoneInfo})
	}
	return valWraps, nil
}

func genNode(rawNode *nebula.Node, timezoneInfo timezoneInfo) (*Node, error) {
	if rawNode == nil {
		return nil, fmt.Errorf("failed to generate Node: invalid rawNode")
	}

	return &Node{
		rawNode:      rawNode,
		timezoneInfo: timezoneInfo,
	}, nil
}

func genEdge(edge *nebula.Edge, timezoneInfo timezoneInfo) (*Edge, error) {
	if edge == nil {
		return nil, fmt.Errorf("failed to generate Node: invalid rawEdge")
	}

	return &Edge{
		rawEdge:      edge,
		timezoneInfo: timezoneInfo,
	}, nil
}

// Used for printing the execution plan. Consider move to a separate module
func graphvizString(s string) string {
	s = strings.Replace(s, "{", "\\{", -1)
	s = strings.Replace(s, "}", "\\}", -1)
	s = strings.Replace(s, "\"", "\\\"", -1)
	s = strings.Replace(s, "[", "\\[", -1)
	s = strings.Replace(s, "]", "\\]", -1)
	s = strings.Replace(s, "(", "\\(", -1)
	s = strings.Replace(s, ")", "\\)", -1)
	s = strings.Replace(s, "<", "\\<", -1)
	s = strings.Replace(s, ">", "\\>", -1)
	return s
}

func prettyFormatJsonString(value []byte) string {
	var prettyJson bytes.Buffer
	if err := json.Indent(&prettyJson, value, "", "  "); err != nil {
		return string(value)
	}
	return prettyJson.String()
}

func name(planNodeDesc *graph.PlanNodeDescription) string {
	return fmt.Sprintf("%s_%d", planNodeDesc.GetName(), planNodeDesc.GetId())
}

func condEdgeLabel(condNode *graph.PlanNodeDescription, doBranch bool) string {
	name := strings.ToLower(string(condNode.GetName()))
	if strings.HasPrefix(name, "select") {
		if doBranch {
			return "Y"
		}
		return "N"
	}
	if strings.HasPrefix(name, "loop") {
		if doBranch {
			return "Do"
		}
	}
	return ""
}

func nodeString(planNodeDesc *graph.PlanNodeDescription, planNodeName string) string {
	var outputVar = graphvizString(string(planNodeDesc.GetOutputVar()))
	var inputVar string
	if planNodeDesc.IsSetDescription() {
		desc := planNodeDesc.GetDescription()
		for _, pair := range desc {
			key := string(pair.GetKey())
			if key == "inputVar" {
				inputVar = graphvizString(string(pair.GetValue()))
			}
		}
	}
	return fmt.Sprintf("\t\"%s\"[label=\"{%s|outputVar: %s|inputVar: %s}\", shape=Mrecord];\n",
		planNodeName, planNodeName, outputVar, inputVar)
}

func edgeString(start, end string) string {
	return fmt.Sprintf("\t\"%s\"->\"%s\";\n", start, end)
}

func conditionalEdgeString(start, end, label string) string {
	return fmt.Sprintf("\t\"%s\"->\"%s\"[label=\"%s\", style=dashed];\n", start, end, label)
}

func conditionalNodeString(name string) string {
	return fmt.Sprintf("\t\"%s\"[shape=diamond];\n", name)
}

func nodeById(p *graph.PlanDescription, nodeId int64) *graph.PlanNodeDescription {
	line := p.GetNodeIndexMap()[nodeId]
	return p.GetPlanNodeDescs()[line]
}

// func findBranchEndNode(p *graph.PlanDescription, condNodeId int64, isDoBranch bool) int64 {
// 	for _, node := range p.GetPlanNodeDescs() {
// 		if node.IsSetBranchInfo() {
// 			bInfo := node.GetBranchInfo()
// 			if bInfo.GetConditionNodeID() == condNodeId && bInfo.GetIsDoBranch() == isDoBranch {
// 				return node.GetId()
// 			}
// 		}
// 	}
// 	return -1
// }

// func findFirstStartNodeFrom(p *graph.PlanDescription, nodeId int64) int64 {
// 	node := nodeById(p, nodeId)
// 	for {
// 		deps := node.GetDependencies()
// 		if len(deps) == 0 {
// 			if strings.ToLower(string(node.GetName())) != "start" {
// 				return -1
// 			}
// 			return node.GetId()
// 		}
// 		node = nodeById(p, deps[0])
// 	}
// }

// explain/profile format="dot"
func (res ResultSet) MakeDotGraph() string {
	p := res.GetPlanDesc()
	planNodeDescs := p.GetPlanNodeDescs()
	var builder strings.Builder
	builder.WriteString("digraph exec_plan {\n")
	builder.WriteString("\trankdir=BT;\n")
	for _, planNodeDesc := range planNodeDescs {
		planNodeName := name(planNodeDesc)
		switch strings.ToLower(string(planNodeDesc.GetName())) {
		default:
			builder.WriteString(nodeString(planNodeDesc, planNodeName))
			if planNodeDesc.IsSetDependencies() {
				for _, depId := range planNodeDesc.GetDependencies() {
					builder.WriteString(edgeString(name(nodeById(p, depId)), planNodeName))
				}
			}
		}
	}
	builder.WriteString("}")
	return builder.String()
}

// // explain/profile format="dot:struct"
func (res ResultSet) MakeDotGraphByStruct() string {
	p := res.GetPlanDesc()
	planNodeDescs := p.GetPlanNodeDescs()
	var builder strings.Builder
	builder.WriteString("digraph exec_plan {\n")
	builder.WriteString("\trankdir=BT;\n")
	for _, planNodeDesc := range planNodeDescs {
		planNodeName := name(planNodeDesc)
		switch strings.ToLower(string(planNodeDesc.GetName())) {
		case "select":
			builder.WriteString(conditionalNodeString(planNodeName))
		case "loop":
			builder.WriteString(conditionalNodeString(planNodeName))
		default:
			builder.WriteString(nodeString(planNodeDesc, planNodeName))
		}

		if planNodeDesc.IsSetDependencies() {
			for _, depId := range planNodeDesc.GetDependencies() {
				dep := nodeById(p, depId)
				builder.WriteString(edgeString(name(dep), planNodeName))
			}
		}
	}
	builder.WriteString("}")
	return builder.String()
}

// explain/profile format="row"
func (res ResultSet) MakePlanByRow() [][]interface{} {
	p := res.GetPlanDesc()
	planNodeDescs := p.GetPlanNodeDescs()
	var rows [][]interface{}
	for _, planNodeDesc := range planNodeDescs {
		var row []interface{}
		row = append(row, planNodeDesc.GetId(), string(planNodeDesc.GetName()))

		if planNodeDesc.IsSetDependencies() {
			var deps []string
			for _, dep := range planNodeDesc.GetDependencies() {
				deps = append(deps, fmt.Sprintf("%d", dep))
			}
			row = append(row, strings.Join(deps, ","))
		} else {
			row = append(row, "")
		}

		if planNodeDesc.IsSetProfiles() {
			var strArr []string
			for i, profile := range planNodeDesc.GetProfiles() {
				otherStats := profile.GetOtherStats()
				if otherStats != nil {
					strArr = append(strArr, "{")
				}
				s := fmt.Sprintf("ver: %d, rows: %d, execTime: %dus, totalTime: %dus",
					i, profile.GetRows(), profile.GetExecDurationInUs(), profile.GetTotalDurationInUs())
				strArr = append(strArr, s)

				for k, v := range otherStats {
					strArr = append(strArr, fmt.Sprintf("%s: %s", k, v))
				}
				if otherStats != nil {
					strArr = append(strArr, "}")
				}
			}
			row = append(row, strings.Join(strArr, "\n"))
		} else {
			row = append(row, "")
		}

		var columnInfo []string

		outputVar := fmt.Sprintf("outputVar: %s", prettyFormatJsonString(planNodeDesc.GetOutputVar()))
		columnInfo = append(columnInfo, outputVar)

		if planNodeDesc.IsSetDescription() {
			desc := planNodeDesc.GetDescription()
			for _, pair := range desc {
				value := prettyFormatJsonString(pair.GetValue())
				columnInfo = append(columnInfo, fmt.Sprintf("%s: %s", string(pair.GetKey()), value))
			}
		}
		row = append(row, strings.Join(columnInfo, "\n"))
		rows = append(rows, row)
	}
	return rows
}
