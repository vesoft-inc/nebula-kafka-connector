package decode

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/types"
)

type (
	typeSchema interface {
		getType() types.ColumnType
	}

	columnTypeSchemaBasic struct {
		typ types.ColumnType
	}

	columnTypeSchemaList struct {
		typ       types.ColumnType
		subSchema typeSchema
	}
	columnTypeSchemaRecord struct {
		typ         types.ColumnType
		propSchemas map[string]typeSchema
	}

	columnTypeSchemaElement struct {
		typ           types.ColumnType
		elemenetProps graphElementProps
	}

	columnTypeSchemaPath struct {
		typ        types.ColumnType
		nodeSchema *columnTypeSchemaElement
		edgeSchema *columnTypeSchemaElement
		meta       *pathMetaData
	}

	graphsSchema map[int32]*graphSchema
	graphSchema  struct {
		name        string
		id          int32
		nodesSchmea map[int32]*elementSchema
		edgesSchema map[int32]*elementSchema
	}
	elementSchema struct {
		typeName string
		typeId   int32
		labels   []string
	}
)

func newTypeSchema(r *bytesReader) (typeSchema, error) {
	b := r.readN(1)
	if r.error() != nil {
		return nil, r.error()
	}
	t, ok := columnTypeMap[b[0]]
	if !ok {
		//TODO
		return nil, fmt.Errorf("unknown column type %d", b[0])
	}
	switch t {
	case types.ColumnTypeList:
		subSchema, err := newTypeSchema(r)
		if err != nil {
			return nil, err
		}
		typ := columnTypeSchemaList{
			typ:       t,
			subSchema: subSchema,
		}
		return &typ, nil
	case types.ColumnTypeRecord:
		// filed num + [name + "0" + type]
		schema := make(map[string]typeSchema)
		numFieldsBytes := r.readN(4)
		if r.error() != nil {
			return nil, r.error()
		}
		numFields := bytesToInt32(numFieldsBytes)
		for i := int32(0); i < numFields; i++ {
			nameBytes := r.readUtilZero()
			if r.error() != nil {
				return nil, r.error()
			}
			propSchema, err := newTypeSchema(r)
			if err != nil {
				return nil, err
			}
			schema[string(nameBytes)] = propSchema
		}
		typ := columnTypeSchemaRecord{
			typ:         t,
			propSchemas: schema,
		}
		return &typ, nil
	case types.ColumnTypeNode:
		pp, err := decodeElementTypes(r, true)
		if err != nil {
			return nil, err
		}
		typ := columnTypeSchemaElement{
			typ:           t,
			elemenetProps: pp,
		}
		return &typ, nil
	case types.ColumnTypeEdge:
		pp, err := decodeElementTypes(r, false)
		if err != nil {
			return nil, err
		}
		typ := columnTypeSchemaElement{
			typ:           t,
			elemenetProps: pp,
		}
		return &typ, nil
	case types.ColumnTypePath:
		// num of elements + [element type]
		typ := columnTypeSchemaPath{
			typ: t,
			nodeSchema: &columnTypeSchemaElement{
				typ:           types.ColumnTypeNode,
				elemenetProps: make(graphElementProps),
			},
			edgeSchema: &columnTypeSchemaElement{
				typ:           types.ColumnTypeEdge,
				elemenetProps: make(graphElementProps),
			},
		}
		elementNumBytes := r.readN(4)
		if r.error() != nil {
			return nil, r.error()
		}
		elementNum := bytesToInt32(elementNumBytes)
		for i := int32(0); i < elementNum; i++ {
			elementSchema, err := newTypeSchema(r)
			if err != nil {
				return nil, err
			}
			if elementSchema.getType() == types.ColumnTypeNode {
				s := elementSchema.(*columnTypeSchemaElement)
				for k, v := range s.elemenetProps {
					typ.nodeSchema.elemenetProps[k] = v
				}
			} else if elementSchema.getType() == types.ColumnTypeEdge {
				s := elementSchema.(*columnTypeSchemaElement)
				for k, v := range s.elemenetProps {
					typ.edgeSchema.elemenetProps[k] = v
				}
			} else {
				return nil, errInvalidColumnType
			}
		}
		return &typ, nil
	case types.ColumnTypeDecimal:
		// precision + scale
		_ = r.readN(2)
		_ = r.readN(2)
		if r.error() != nil {
			return nil, r.error()
		}
		return &columnTypeSchemaBasic{
			typ: t,
		}, nil
	default:
		return &columnTypeSchemaBasic{
			typ: t,
		}, nil
	}
}

func (s *columnTypeSchemaBasic) getType() types.ColumnType {
	return s.typ
}

func (s *columnTypeSchemaList) getType() types.ColumnType {
	return s.typ
}

func (s *columnTypeSchemaRecord) getType() types.ColumnType {
	return s.typ
}

func (s *columnTypeSchemaElement) getType() types.ColumnType {
	return s.typ
}

func (s *columnTypeSchemaElement) getElementProps() graphElementProps {
	return s.elemenetProps
}

func (s *columnTypeSchemaPath) getType() types.ColumnType {
	return s.typ
}

func (s *columnTypeSchemaPath) getNodeSchema() *columnTypeSchemaElement {
	return s.nodeSchema
}

func (s *columnTypeSchemaPath) getEdgeSchema() *columnTypeSchemaElement {
	return s.edgeSchema
}

func decodeElementTypes(r *bytesReader, isNode bool) (graphElementProps, error) {
	// num of element type + [element type id + num of props + [prop name + prop type]]
	// ignore the first byte nodeType
	elementTypeSize := 2
	if !isNode {
		elementTypeSize = 4
	}
	numElementTypeBytes := r.readN(4)
	if r.error() != nil {
		return nil, r.error()
	}
	numElementType := bytesToInt32(numElementTypeBytes)
	elementTypes := make(graphElementProps)
	for i := 0; i < int(numElementType); i++ {
		var elementTypeId int32
		elementTypeIdBytes := r.readN(elementTypeSize)
		if r.error() != nil {
			return nil, r.error()
		}
		if isNode {
			elementTypeId = int32(bytesToInt16(elementTypeIdBytes))
		} else {
			elementTypeId = bytesToInt32(elementTypeIdBytes)
		}
		nt := make(map[string]*vectorProps)
		numPropsBytes := r.readN(4)
		if r.error() != nil {
			return nil, r.error()
		}
		numProps := bytesToInt32(numPropsBytes)
		for j := 0; j < int(numProps); j++ {
			propNameBytes := r.readUtilZero()
			if r.error() != nil {
				return nil, r.error()
			}
			prop, err := decodePropNameAndType(propNameBytes, r)
			if err != nil {
				return nil, err
			}
			nt[prop.name] = prop
		}
		elementTypes[elementTypeId] = nt
	}
	return elementTypes, nil
}

func decodePropNameAndType(name []byte, r *bytesReader) (*vectorProps, error) {
	var prps vectorProps

	prps.name = string(name)
	tt, err := newTypeSchema(r)
	if err != nil {
		return nil, err
	}
	prps.typ = tt
	return &prps, nil
}
