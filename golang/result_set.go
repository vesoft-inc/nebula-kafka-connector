// Copyright (c) 2022 vesoft inc. All rights reserved.

package nebula_ng

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/generated_code/v5.0.0/nebula"
)

type ResultSet struct {
	resp            nebula.ExecutionResponse
	columnNames     []string
	colNameIndexMap map[string]int
	timezoneInfo    timezoneInfo
	planDesc        map[string]interface{}
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
func (res ResultSet) GetRows() []nebula.Row {
	if res.resp.ExecutionOutcome.Result_ == nil {
		var empty []nebula.Row
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
	if _, ok := list.([]nebula.Row); ok {
		if index < 0 || index >= len(list.([]nebula.Row)) {
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
		fmt.Printf("Failed at checking the index %d: %s\n", index, err.Error())
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
	return len(res.planDesc) != 0
}

func (res ResultSet) GetPlanDesc() map[string]interface{} {
	return res.planDesc
}

func (res ResultSet) GetOptimizeTimeInUs() int64 {
	if !res.IsSetPlanDesc() {
		return 0
	}
	return parseInt64(res.planDesc, "optimizeTimeInUs")
}

func (res ResultSet) GetPreamble() string {
	if !res.IsSetPlanDesc() {
		return ""
	}
	return parseString(res.planDesc, "preamble")
}

func (res ResultSet) GetHeader() []string {
	if !res.IsSetPlanDesc() {
		return nil
	}
	var header []string
	for _, col := range res.planDesc["header"].([]interface{}) {
		header = append(header, col.(string))
	}
	return header
}

func (res ResultSet) GetPlanPrintFormat() string {
	// TODO(yee): support more formats
	return "row"
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
	rawNode      nebula.Node
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
	rawEdge      nebula.Edge
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
	time         nebula.LocalTime
	timezoneInfo timezoneInfo
}

func genLocalTimeWrapper(time nebula.LocalTime, timezoneInfo timezoneInfo) (*LocalTimeWrapper, error) {
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
func (t LocalTimeWrapper) getRawTime() nebula.LocalTime {
	return t.time
}

func (t1 LocalTimeWrapper) IsEqualTo(t2 LocalTimeWrapper) bool {
	return t1.getHour() == t2.getHour() &&
		t1.getSecond() == t2.getSecond() &&
		t1.getSecond() == t2.getSecond() &&
		t1.getMicrosec() == t2.getMicrosec()
}

type DateWrapper struct {
	date nebula.Date
}

func genDateWrapper(date nebula.Date) (*DateWrapper, error) {
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
func (d DateWrapper) getRawDate() nebula.Date {
	return d.date
}

func (d1 DateWrapper) IsEqualTo(d2 DateWrapper) bool {
	return d1.getYear() == d2.getYear() &&
		d1.getMonth() == d2.getMonth() &&
		d1.getDay() == d2.getDay()
}

type LocalDatetimeWrapper struct {
	dateTime     nebula.LocalDatetime
	timezoneInfo timezoneInfo
}

func genLocalDatetimeWrapper(datetime nebula.LocalDatetime, timezoneInfo timezoneInfo) (*LocalDatetimeWrapper, error) {
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
func (dt LocalDatetimeWrapper) getRawDateTime() nebula.LocalDatetime {
	return dt.dateTime
}

func genResultSet(resp nebula.ExecutionResponse, timezoneInfo timezoneInfo) (*ResultSet, error) {
	var colNames []string
	var colNameIndexMap = make(map[string]int)
	var planDesc = make(map[string]interface{})

	if resp.ExecutionOutcome.PlanDesc != nil {
		json.Unmarshal(resp.ExecutionOutcome.PlanDesc, &planDesc)
	}

	if resp.ExecutionOutcome.Result_ == nil { // if resp.Data != nil then resp.Data.row and resp.Data.colNames wont be nil
		return &ResultSet{
			resp:            resp,
			columnNames:     colNames,
			colNameIndexMap: colNameIndexMap,
			planDesc:        planDesc,
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
		planDesc:        planDesc,
	}, nil
}

func genValWraps(row nebula.Row, timezoneInfo timezoneInfo) ([]*ValueWrapper, error) {
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

func genNode(rawNode nebula.Node, timezoneInfo timezoneInfo) (*Node, error) {
	if rawNode == nil {
		return nil, fmt.Errorf("failed to generate Node: invalid rawNode")
	}

	return &Node{
		rawNode:      rawNode,
		timezoneInfo: timezoneInfo,
	}, nil
}

func genEdge(edge nebula.Edge, timezoneInfo timezoneInfo) (*Edge, error) {
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

// func name(planNodeDesc PlanNodeDescription) string {
// 	return fmt.Sprintf("%s_%d", planNodeDesc.GetName(), planNodeDesc.GetId())
// }

// func condEdgeLabel(condNode PlanNodeDescription, doBranch bool) string {
// 	name := strings.ToLower(string(condNode.GetName()))
// 	if strings.HasPrefix(name, "select") {
// 		if doBranch {
// 			return "Y"
// 		}
// 		return "N"
// 	}
// 	if strings.HasPrefix(name, "loop") {
// 		if doBranch {
// 			return "Do"
// 		}
// 	}
// 	return ""
// }

// func nodeString(planNodeDesc PlanNodeDescription, planNodeName string) string {
// 	var outputVar = graphvizString(string(planNodeDesc.GetOutputVar()))
// 	var inputVar string
// 	if planNodeDesc.IsSetDescription() {
// 		desc := planNodeDesc.GetDescription()
// 		for _, pair := range desc {
// 			key := string(pair.GetKey())
// 			if key == "inputVar" {
// 				inputVar = graphvizString(string(pair.GetValue()))
// 			}
// 		}
// 	}
// 	return fmt.Sprintf("\t\"%s\"[label=\"{%s|outputVar: %s|inputVar: %s}\", shape=Mrecord];\n",
// 		planNodeName, planNodeName, outputVar, inputVar)
// }

func edgeString(start, end string) string {
	return fmt.Sprintf("\t\"%s\"->\"%s\";\n", start, end)
}

func conditionalEdgeString(start, end, label string) string {
	return fmt.Sprintf("\t\"%s\"->\"%s\"[label=\"%s\", style=dashed];\n", start, end, label)
}

func conditionalNodeString(name string) string {
	return fmt.Sprintf("\t\"%s\"[shape=diamond];\n", name)
}

// func nodeById(p PlanDescription, nodeId int64) PlanNodeDescription {
// 	line := p.GetNodeIndexMap()[nodeId]
// 	return p.GetPlanNodeDescs()[line]
// }

// func findBranchEndNode(p *PlanDescription, condNodeId int64, isDoBranch bool) int64 {
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

// func findFirstStartNodeFrom(p *PlanDescription, nodeId int64) int64 {
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
	// p := res.GetPlanDesc()
	// planNodeDescs := p.GetPlanNodeDescs()
	// var builder strings.Builder
	// builder.WriteString("digraph exec_plan {\n")
	// builder.WriteString("\trankdir=BT;\n")
	// for _, planNodeDesc := range planNodeDescs {
	// 	planNodeName := name(planNodeDesc)
	// 	switch strings.ToLower(string(planNodeDesc.GetName())) {
	// 	default:
	// 		builder.WriteString(nodeString(planNodeDesc, planNodeName))
	// 		if planNodeDesc.IsSetDependencies() {
	// 			for _, depId := range planNodeDesc.GetDependencies() {
	// 				builder.WriteString(edgeString(name(nodeById(p, depId)), planNodeName))
	// 			}
	// 		}
	// 	}
	// }
	// builder.WriteString("}")
	// return builder.String()

	return ""
}

// // explain/profile format="dot:struct"
func (res ResultSet) MakeDotGraphByStruct() string {
	// p := res.GetPlanDesc()
	// planNodeDescs := p.GetPlanNodeDescs()
	// var builder strings.Builder
	// builder.WriteString("digraph exec_plan {\n")
	// builder.WriteString("\trankdir=BT;\n")
	// for _, planNodeDesc := range planNodeDescs {
	// 	planNodeName := name(planNodeDesc)
	// 	switch strings.ToLower(string(planNodeDesc.GetName())) {
	// 	case "select":
	// 		builder.WriteString(conditionalNodeString(planNodeName))
	// 	case "loop":
	// 		builder.WriteString(conditionalNodeString(planNodeName))
	// 	default:
	// 		builder.WriteString(nodeString(planNodeDesc, planNodeName))
	// 	}

	// 	if planNodeDesc.IsSetDependencies() {
	// 		for _, depId := range planNodeDesc.GetDependencies() {
	// 			dep := nodeById(p, depId)
	// 			builder.WriteString(edgeString(name(dep), planNodeName))
	// 		}
	// 	}
	// }
	// builder.WriteString("}")
	// return builder.String()

	return ""
}

type OperatorUniqueId struct {
	pipelineId int32
	operatorId int32
	planNodeId int64
	inStorage  bool
}

var InvalidOperatorUniqueId = OperatorUniqueId{
	pipelineId: -1,
	operatorId: -1,
	planNodeId: -1,
	inStorage:  false,
}

func (i OperatorUniqueId) Equals(other OperatorUniqueId) bool {
	return i.pipelineId == other.pipelineId &&
		i.operatorId == other.operatorId &&
		i.planNodeId == other.planNodeId &&
		i.inStorage == other.inStorage
}

type KeyType interface {
	int64 | OperatorUniqueId
}

type GetKey[T KeyType] func(val interface{}) T

func parseFloat(m interface{}, name string) float64 {
	return m.(map[string]interface{})[name].(float64)
}

func parseInt64(m interface{}, name string) int64 {
	return int64(parseFloat(m, name))
}

func parseInt32(m interface{}, name string) int32 {
	return int32(parseFloat(m, name))
}

func parseString(m interface{}, name string) string {
	if tmp, ok := m.(map[string]interface{}); ok {
		if v, ok := tmp[name]; ok {
			return v.(string)
		}
	}
	return ""
}

func parseBool(m interface{}, name string) bool {
	return m.(map[string]interface{})[name].(bool)
}

// explain/profile format="row"
func (res ResultSet) MakePlanByRow() (rightSepToTailWidth []int, rows [][]interface{}) {
	if !res.IsSetPlanDesc() {
		return
	}

	header := res.GetHeader()
	if len(header) == 0 {
		return
	}

	for i := 0; i < len(header); i++ {
		rightSepToTailWidth = append(rightSepToTailWidth, 0)
	}

	switch res.GetPreamble() {
	case "explain":
		var idToPlanNodeMap = make(map[int64]map[string]interface{})
		var planNodeDescs []interface{}
		if header0, ok := res.planDesc[header[0]]; ok {
			planNodeDescs = header0.([]interface{})
		} else {
			return
		}

		for _, planNodeDesc := range planNodeDescs {
			planNode := planNodeDesc.(map[string]interface{})
			id := parseInt64(planNode, "id")
			idToPlanNodeMap[id] = planNode
		}

		rootNode := planNodeDescs[0].(map[string]interface{})
		makeAsciiPlanTreeText(
			res,
			idToPlanNodeMap,
			header,
			[]map[string]interface{}{rootNode},
			true,
			"",
			func(val interface{}) int64 { return int64(val.(float64)) },
			&rightSepToTailWidth,
			&rows,
		)
	case "profile":
		var idToOperatorMap = make(map[OperatorUniqueId]map[string]interface{})

		var operatorObject map[string]interface{}
		if header0, ok := res.planDesc[header[0]]; ok {
			operatorObject = header0.(map[string]interface{})
		} else {
			return
		}

		var pipelines []interface{}
		if pipelinesObj, ok := operatorObject["pipelines"]; ok {
			pipelines = pipelinesObj.([]interface{})
		} else {
			return
		}

		var rootOperator map[string]interface{}

		addChild := func(parent OperatorUniqueId, child OperatorUniqueId) {
			if idToOperatorMap[parent] != nil {
				parentOp := idToOperatorMap[parent]
				if childrenObj, ok := parentOp["children"]; ok {
					children := childrenObj.([]interface{})
					children = append(children, child)
					parentOp["children"] = children
				} else {
					parentOp["children"] = []interface{}{child}
				}
			}
		}

		getPipeId := func(val interface{}) OperatorUniqueId {
			obj := val.(map[string]interface{})
			return OperatorUniqueId{
				pipelineId: parseInt32(obj, "pipelineId"),
				operatorId: parseInt32(obj, "operatorId"),
				planNodeId: parseInt64(obj, "planNodeId"),
				inStorage:  parseBool(obj, "inStorage"),
			}
		}

		for _, pipeline := range pipelines {
			firstTime := false
			var prevOperatorId OperatorUniqueId

			pipeObj := pipeline.(map[string]interface{})
			consumerOperatorId := InvalidOperatorUniqueId
			if consumerOpIdObj, ok := pipeObj["consumerOperatorId"]; ok {
				consumerOperatorId = getPipeId(consumerOpIdObj)
			}

			var operatorArray []interface{}
			if operatorArrayObj, ok := pipeObj["operators"]; ok {
				operatorArray = operatorArrayObj.([]interface{})
			} else {
				continue
			}

			for i := len(operatorArray) - 1; i >= 0; i-- {
				opObj := operatorArray[i].(map[string]interface{})
				id := getPipeId(opObj["id"])
				idToOperatorMap[id] = opObj

				if i == len(operatorArray)-1 && consumerOperatorId != InvalidOperatorUniqueId {
					addChild(consumerOperatorId, id)
				}

				if !firstTime {
					if len(rootOperator) == 0 {
						rootOperator = opObj
					}
					firstTime = true
				} else {
					addChild(prevOperatorId, id)
				}
				prevOperatorId = id
			}
		}

		makeAsciiPlanTreeText(
			res,
			idToOperatorMap,
			header,
			[]map[string]interface{}{rootOperator},
			true,
			"",
			func(val interface{}) OperatorUniqueId { return val.(OperatorUniqueId) },
			&rightSepToTailWidth,
			&rows,
		)
	}

	return
}

const (
	dash       string = `├─`
	spacer     string = `│ `
	dashLast   string = `└─`
	spacerLast string = `  `
)

func makeAsciiPlanTreeText[T KeyType](
	res ResultSet,
	idToPlanNodeDescMap map[T]map[string]interface{},
	header []string,
	planNodes []map[string]interface{},
	isRoot bool,
	prefix string,
	getKey GetKey[T],
	rightSepToTailWidth *[]int,
	rows *[][]interface{}) {
	for i, planNode := range planNodes {
		lineBuffer := prefix
		childPrefix := prefix
		if !isRoot {
			if children, ok := planNode["children"]; ok && len(children.([]interface{})) > 0 {
				if i == len(planNodes)-1 {
					lineBuffer += dashLast
					childPrefix += spacerLast
				} else {
					lineBuffer += dash
					childPrefix += spacer
				}
			} else {
				if i == len(planNodes)-1 {
					lineBuffer += dashLast
				} else {
					lineBuffer += dash
				}
				childPrefix += spacerLast
			}
		}
		lineBuffer += parseString(planNode, "name")
		row := []interface{}{lineBuffer}
		for i := 1; i < len(header); i++ {
			width := 0
			if col, ok := planNode[header[i]]; ok {
				if str, ok := col.(string); ok && strings.Contains(header[i], "/") {
					if pos := strings.LastIndex(str, "/"); pos > 0 {
						width = len(str) - pos
					}
				}
				row = append(row, col)
			} else {
				row = append(row, "")
			}
			if width > (*rightSepToTailWidth)[i] {
				(*rightSepToTailWidth)[i] = width
			}
		}
		*rows = append(*rows, row)

		var childPlanNodes []map[string]interface{}
		if children, ok := planNode["children"]; ok {
			for _, child := range children.([]interface{}) {
				key := getKey(child)
				childPlanNodes = append(childPlanNodes, idToPlanNodeDescMap[key])
			}
		}

		makeAsciiPlanTreeText(
			res,
			idToPlanNodeDescMap,
			header,
			childPlanNodes,
			false,
			childPrefix,
			getKey,
			rightSepToTailWidth,
			rows,
		)
	}
}
