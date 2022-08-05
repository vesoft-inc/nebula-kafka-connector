// Copyright (c) 2022 vesoft inc. All rights reserved.

package nebula_ng_go

import (
	"fmt"

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

func (res ResultSet) IsSuccess() bool {
	return res.GetStatus() == string("SUCCESS")
}

// Returns all rows
func (res ResultSet) GetRows() []*nebula.RawRecord {
	if res.resp.ExecutionOutcome.Result_ == nil {
		var empty []*nebula.RawRecord
		return empty
	}
	return res.resp.ExecutionOutcome.Result_.Records
}

func (res ResultSet) GetColNames() []string {
	return res.columnNames
}

func (res ResultSet) GetLatency() int64 {
	return res.resp.LatencyInUs
}

type Record struct {
	columnNames     *[]string
	_record         []*ValueWrapper
	colNameIndexMap *map[string]int
	timezoneInfo    timezoneInfo
}

type Node struct {
	rawNode      *nebula.Node
	timezoneInfo timezoneInfo
}

// TODO(Aiee) For demo only
func (node *Node) String() string {
	nodeID := node.rawNode.NodeID
	nodeTypeID := node.rawNode.NodeTypeID
	var kvStr []string
	for key, value := range node.rawNode.Properties {
		kvTemp := fmt.Sprintf("%s: %s", key, ValueWrapper{value, node.timezoneInfo}.String())
		kvStr = append(kvStr, kvTemp)
	}

	// No properties
	if len(kvStr) == 0 {
		return fmt.Sprintf("(nodeID:%d, nodeTypeID: %d)", nodeID, nodeTypeID)
	}
	return fmt.Sprintf("(nodeID:%d, nodeTypeID:%d, props:{%s})", nodeID, nodeTypeID, kvStr)
}

type Edge struct {
	edge         *nebula.Edge
	timezoneInfo timezoneInfo
}

// TODO(Aiee) For demo only
func (edge *Edge) String() string {
	srcID := edge.edge.SrcID
	dstID := edge.edge.DstID
	edgeTypeID := edge.edge.EdgeTypeID
	edgeRank := edge.edge.Rank

	var kvStr []string
	for key, value := range edge.edge.Properties {
		kvTemp := fmt.Sprintf("%s: %s", key, ValueWrapper{value, edge.timezoneInfo}.String())
		kvStr = append(kvStr, kvTemp)
	}

	// No properties
	if len(kvStr) == 0 {
		return fmt.Sprintf("(srcID: %d, dstID: %d, edgeTypeID: %d, edgeRank: %d)",
			srcID, dstID, edgeTypeID, edgeRank)
	}
	return fmt.Sprintf("(srcID: %d, dstID: %d, edgeTypeID: %d, edgeRank: %d, props: %s)",
		srcID, dstID, edgeTypeID, edgeRank, kvStr)
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
		return nil, fmt.Errorf("failed to generate Node: invalid rawNode")
	}

	return &Edge{
		edge:         edge,
		timezoneInfo: timezoneInfo,
	}, nil
}

// func graphvizString(s string) string {
// 	s = strings.Replace(s, "{", "\\{", -1)
// 	s = strings.Replace(s, "}", "\\}", -1)
// 	s = strings.Replace(s, "\"", "\\\"", -1)
// 	s = strings.Replace(s, "[", "\\[", -1)
// 	s = strings.Replace(s, "]", "\\]", -1)
// 	s = strings.Replace(s, "(", "\\(", -1)
// 	s = strings.Replace(s, ")", "\\)", -1)
// 	s = strings.Replace(s, "<", "\\<", -1)
// 	s = strings.Replace(s, ">", "\\>", -1)
// 	return s
// }

// func prettyFormatJsonString(value []byte) string {
// 	var prettyJson bytes.Buffer
// 	if err := json.Indent(&prettyJson, value, "", "  "); err != nil {
// 		return string(value)
// 	}
// 	return prettyJson.String()
// }

// func name(planNodeDesc *graph.PlanNodeDescription) string {
// 	return fmt.Sprintf("%s_%d", planNodeDesc.GetName(), planNodeDesc.GetId())
// }

// func condEdgeLabel(condNode *graph.PlanNodeDescription, doBranch bool) string {
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

// func nodeString(planNodeDesc *graph.PlanNodeDescription, planNodeName string) string {
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

// func edgeString(start, end string) string {
// 	return fmt.Sprintf("\t\"%s\"->\"%s\";\n", start, end)
// }

// func conditionalEdgeString(start, end, label string) string {
// 	return fmt.Sprintf("\t\"%s\"->\"%s\"[label=\"%s\", style=dashed];\n", start, end, label)
// }

// func conditionalNodeString(name string) string {
// 	return fmt.Sprintf("\t\"%s\"[shape=diamond];\n", name)
// }

// func nodeById(p *graph.PlanDescription, nodeId int64) *graph.PlanNodeDescription {
// 	line := p.GetNodeIndexMap()[nodeId]
// 	return p.GetPlanNodeDescs()[line]
// }

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

// // explain/profile format="dot"
// func (res ResultSet) MakeDotGraph() string {
// 	p := res.GetPlanDesc()
// 	planNodeDescs := p.GetPlanNodeDescs()
// 	var builder strings.Builder
// 	builder.WriteString("digraph exec_plan {\n")
// 	builder.WriteString("\trankdir=BT;\n")
// 	for _, planNodeDesc := range planNodeDescs {
// 		planNodeName := name(planNodeDesc)
// 		switch strings.ToLower(string(planNodeDesc.GetName())) {
// 		case "select":
// 			builder.WriteString(conditionalNodeString(planNodeName))
// 			dep := nodeById(p, planNodeDesc.GetDependencies()[0])
// 			// then branch
// 			thenNodeId := findBranchEndNode(p, planNodeDesc.GetId(), true)
// 			builder.WriteString(edgeString(name(nodeById(p, thenNodeId)), name(dep)))
// 			thenStartId := findFirstStartNodeFrom(p, thenNodeId)
// 			builder.WriteString(conditionalEdgeString(name(planNodeDesc), name(nodeById(p, thenStartId)), "Y"))
// 			// else branch
// 			elseNodeId := findBranchEndNode(p, planNodeDesc.GetId(), false)
// 			builder.WriteString(edgeString(name(nodeById(p, elseNodeId)), name(dep)))
// 			elseStartId := findFirstStartNodeFrom(p, elseNodeId)
// 			builder.WriteString(conditionalEdgeString(name(planNodeDesc), name(nodeById(p, elseStartId)), "N"))
// 			// dep
// 			builder.WriteString(edgeString(name(dep), planNodeName))
// 		case "loop":
// 			builder.WriteString(conditionalNodeString(planNodeName))
// 			dep := nodeById(p, planNodeDesc.GetDependencies()[0])
// 			// do branch
// 			doNodeId := findBranchEndNode(p, planNodeDesc.GetId(), true)
// 			builder.WriteString(edgeString(name(nodeById(p, doNodeId)), name(planNodeDesc)))
// 			doStartId := findFirstStartNodeFrom(p, doNodeId)
// 			builder.WriteString(conditionalEdgeString(name(planNodeDesc), name(nodeById(p, doStartId)), "Do"))
// 			// dep
// 			builder.WriteString(edgeString(name(dep), planNodeName))
// 		default:
// 			builder.WriteString(nodeString(planNodeDesc, planNodeName))
// 			if planNodeDesc.IsSetDependencies() {
// 				for _, depId := range planNodeDesc.GetDependencies() {
// 					builder.WriteString(edgeString(name(nodeById(p, depId)), planNodeName))
// 				}
// 			}
// 		}
// 	}
// 	builder.WriteString("}")
// 	return builder.String()
// }

// // explain/profile format="dot:struct"
// func (res ResultSet) MakeDotGraphByStruct() string {
// 	p := res.GetPlanDesc()
// 	planNodeDescs := p.GetPlanNodeDescs()
// 	var builder strings.Builder
// 	builder.WriteString("digraph exec_plan {\n")
// 	builder.WriteString("\trankdir=BT;\n")
// 	for _, planNodeDesc := range planNodeDescs {
// 		planNodeName := name(planNodeDesc)
// 		switch strings.ToLower(string(planNodeDesc.GetName())) {
// 		case "select":
// 			builder.WriteString(conditionalNodeString(planNodeName))
// 		case "loop":
// 			builder.WriteString(conditionalNodeString(planNodeName))
// 		default:
// 			builder.WriteString(nodeString(planNodeDesc, planNodeName))
// 		}

// 		if planNodeDesc.IsSetDependencies() {
// 			for _, depId := range planNodeDesc.GetDependencies() {
// 				dep := nodeById(p, depId)
// 				builder.WriteString(edgeString(name(dep), planNodeName))
// 			}
// 		}

// 		if planNodeDesc.IsSetBranchInfo() {
// 			branchInfo := planNodeDesc.GetBranchInfo()
// 			condNode := nodeById(p, branchInfo.GetConditionNodeID())
// 			label := condEdgeLabel(condNode, branchInfo.GetIsDoBranch())
// 			builder.WriteString(conditionalEdgeString(planNodeName, name(condNode), label))
// 		}
// 	}
// 	builder.WriteString("}")
// 	return builder.String()
// }

// // explain/profile format="row"
// func (res ResultSet) MakePlanByRow() [][]interface{} {
// 	p := res.GetPlanDesc()
// 	planNodeDescs := p.GetPlanNodeDescs()
// 	var rows [][]interface{}
// 	for _, planNodeDesc := range planNodeDescs {
// 		var row []interface{}
// 		row = append(row, planNodeDesc.GetId(), string(planNodeDesc.GetName()))

// 		if planNodeDesc.IsSetDependencies() {
// 			var deps []string
// 			for _, dep := range planNodeDesc.GetDependencies() {
// 				deps = append(deps, fmt.Sprintf("%d", dep))
// 			}
// 			row = append(row, strings.Join(deps, ","))
// 		} else {
// 			row = append(row, "")
// 		}

// 		if planNodeDesc.IsSetProfiles() {
// 			var strArr []string
// 			for i, profile := range planNodeDesc.GetProfiles() {
// 				otherStats := profile.GetOtherStats()
// 				if otherStats != nil {
// 					strArr = append(strArr, "{")
// 				}
// 				s := fmt.Sprintf("ver: %d, rows: %d, execTime: %dus, totalTime: %dus",
// 					i, profile.GetRows(), profile.GetExecDurationInUs(), profile.GetTotalDurationInUs())
// 				strArr = append(strArr, s)

// 				for k, v := range otherStats {
// 					strArr = append(strArr, fmt.Sprintf("%s: %s", k, v))
// 				}
// 				if otherStats != nil {
// 					strArr = append(strArr, "}")
// 				}
// 			}
// 			row = append(row, strings.Join(strArr, "\n"))
// 		} else {
// 			row = append(row, "")
// 		}

// 		var columnInfo []string
// 		if planNodeDesc.IsSetBranchInfo() {
// 			branchInfo := planNodeDesc.GetBranchInfo()
// 			columnInfo = append(columnInfo, fmt.Sprintf("branch: %t, nodeId: %d\n",
// 				branchInfo.GetIsDoBranch(), branchInfo.GetConditionNodeID()))
// 		}

// 		outputVar := fmt.Sprintf("outputVar: %s", prettyFormatJsonString(planNodeDesc.GetOutputVar()))
// 		columnInfo = append(columnInfo, outputVar)

// 		if planNodeDesc.IsSetDescription() {
// 			desc := planNodeDesc.GetDescription()
// 			for _, pair := range desc {
// 				value := prettyFormatJsonString(pair.GetValue())
// 				columnInfo = append(columnInfo, fmt.Sprintf("%s: %s", string(pair.GetKey()), value))
// 			}
// 		}
// 		row = append(row, strings.Join(columnInfo, "\n"))
// 		rows = append(rows, row)
// 	}
// 	return rows
// }
