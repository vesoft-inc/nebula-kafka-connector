package nebula_ng

import "strings"

const (
	dash       string = `├─`
	spacer     string = `│ `
	dashLast   string = `└─`
	spacerLast string = `  `
)

type plan struct {
	planDesc map[string]interface{}
}

func (p *plan) GetHeader() []string {
	var header []string
	for _, col := range p.planDesc["header"].([]interface{}) {
		header = append(header, col.(string))
	}
	return header
}

func (p *plan) GetPlanPrintFormat() string {
	return "row"
}

func (p *plan) GetPlanDesc() map[string]interface{} {
	return p.planDesc
}

func (p *plan) MakePlanByRow() (rightSepToTailWidth []int, rows [][]interface{}) {
	header := p.GetHeader()
	if len(header) == 0 {
		return
	}

	for i := 0; i < len(header); i++ {
		rightSepToTailWidth = append(rightSepToTailWidth, 0)
	}

	switch p.getPreamble() {
	case "explain":
		var idToPlanNodeMap = make(map[int64]map[string]interface{})
		var planNodeDescs []interface{}
		if header0, ok := p.planDesc[header[0]]; ok {
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
		var idToOperatorMap = make(map[operatorUniqueId]map[string]interface{})

		var operatorObject map[string]interface{}
		if header0, ok := p.planDesc[header[0]]; ok {
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

		addChild := func(parent operatorUniqueId, child operatorUniqueId) {
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

		getPipeId := func(val interface{}) operatorUniqueId {
			obj := val.(map[string]interface{})
			return operatorUniqueId{
				pipelineId: parseInt32(obj, "pipelineId"),
				operatorId: parseInt32(obj, "operatorId"),
				planNodeId: parseInt64(obj, "planNodeId"),
				inStorage:  parseBool(obj, "inStorage"),
			}
		}

		for _, pipeline := range pipelines {
			firstTime := false
			var prevOperatorId operatorUniqueId

			pipeObj := pipeline.(map[string]interface{})
			consumerOperatorId := invalidOperatorUniqueId
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

				if i == len(operatorArray)-1 && consumerOperatorId != invalidOperatorUniqueId {
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
			idToOperatorMap,
			header,
			[]map[string]interface{}{rootOperator},
			true,
			"",
			func(val interface{}) operatorUniqueId { return val.(operatorUniqueId) },
			&rightSepToTailWidth,
			&rows,
		)
	}

	return
}
func (p *plan) GetBuildTimeInUs() int64 {
	return parseInt64(p.planDesc, "buildTimeInUs")
}
func (p *plan) GetOptimizeTimeInUs() int64 {
	return parseInt64(p.planDesc, "optimizeTimeInUs")
}

func (p *plan) getPreamble() string {
	return parseString(p.planDesc, "preamble")
}

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

type operatorUniqueId struct {
	pipelineId int32
	operatorId int32
	planNodeId int64
	inStorage  bool
}

var invalidOperatorUniqueId = operatorUniqueId{
	pipelineId: -1,
	operatorId: -1,
	planNodeId: -1,
	inStorage:  false,
}

func (i operatorUniqueId) Equals(other operatorUniqueId) bool {
	return i.pipelineId == other.pipelineId &&
		i.operatorId == other.operatorId &&
		i.planNodeId == other.planNodeId &&
		i.inStorage == other.inStorage
}

type keyType interface {
	int64 | operatorUniqueId
}

type getKey[T keyType] func(val interface{}) T

func makeAsciiPlanTreeText[T keyType](
	idToPlanNodeDescMap map[T]map[string]interface{},
	header []string,
	planNodes []map[string]interface{},
	isRoot bool,
	prefix string,
	getKey getKey[T],
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
