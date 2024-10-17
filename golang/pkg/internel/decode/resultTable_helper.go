package decode

import (
	"fmt"
	"math"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/errors"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/vector"
	internal_error "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/internal_error"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/types"
)

var (
	errTypeAssertion       = internal_error.ErrInternel("type assertion failed")
	errOutOfRange          = internal_error.ErrInternel("index out of range")
	errInvalidColumnType   = internal_error.ErrInternel("invalid type")
	errInvalidVectorType   = internal_error.ErrInternel("invalid vector type")
	errNoZeroString        = internal_error.ErrInternel("no zero found in string")
	errElementTypeNotFount = internal_error.ErrInternel("element type not found")
	errGraphNotFound       = internal_error.ErrInternel("graph not found")
	errPropNotFound        = internal_error.ErrInternel("prop not found")
	errBatchNotEqual       = internal_error.ErrInternel("batch not equal")
	errColumnTypeIsNil     = internal_error.ErrInternel("column type is nil")
)

const (
	kMicrosecondsOfSecond = 1000000
	kMicrosecondsOfMinute = 60 * kMicrosecondsOfSecond
	kMicrosecondsOfHour   = 60 * kMicrosecondsOfMinute
	kMicrosecondsOfDay    = 24 * kMicrosecondsOfHour
)

var sizeMap = map[types.ColumnType]int{
	types.ColumnTypeBool:          1,
	types.ColumnTypeInt8:          1,
	types.ColumnTypeUint8:         1,
	types.ColumnTypeInt16:         2,
	types.ColumnTypeUint16:        2,
	types.ColumnTypeInt32:         4,
	types.ColumnTypeUint32:        4,
	types.ColumnTypeInt64:         8,
	types.ColumnTypeUint64:        8,
	types.ColumnTypeFloat32:       4,
	types.ColumnTypeFloat64:       8,
	types.ColumnTypeDate:          4,
	types.ColumnTypeLocalTime:     8,
	types.ColumnTypeZonedTime:     8,
	types.ColumnTypeLocalDatetime: 8,
	types.ColumnTypeZonedDatetime: 8,
	types.ColumnTypeDuration:      8,
}

var kOneBitmasks = []byte{
	1 << 0,
	1 << 1,
	1 << 2,
	1 << 3,
	1 << 4,
	1 << 5,
	1 << 6,
	1 << 7,
}

type decodeFlatFn func(dctx *decodeContext, v *vector.NestedVector, index uint32, columnType typeSchema) (*nebulaValue, error)

type decoder interface {
	// decodeValue decodes the value from the vector
	decodeValue(dctx *decodeContext, v *vector.NestedVector, t vectorType, index uint32, columnType typeSchema) (*nebulaValue, error)
}
type vectorDecoder struct {
	decodeFlatFns map[types.ColumnType]decodeFlatFn
}

var defaultDecoder decoder = &vectorDecoder{}

// for common value, the bytes size is fixed
// use a function to get the bytes
func (c *vectorDecoder) getCommonBytes(v *vector.NestedVector, typ types.ColumnType, index uint32) ([]byte, error) {
	size, ok := sizeMap[typ]
	if !ok {
		return nil, errors.Wrap(errInvalidColumnType, "")
	}
	if len(v.VectorData) < int(index)*size+size {
		return nil, errors.Wrap(errOutOfRange, "")
	}
	return v.VectorData[int(index)*size : int(index)*size+size], nil
}

func (c *vectorDecoder) decodeValue(dctx *decodeContext, v *vector.NestedVector, t vectorType,
	index uint32, columnType typeSchema) (*nebulaValue, error) {
	if v.NullBitMap != nil {
		if v.NullBitMap[index/8]&kOneBitmasks[index%8] == 0 {
			return &nebulaValue{data: nil}, nil
		}
	}
	switch t {
	case vectorTypeFlat:
		return c.decodeFlatValue(dctx, v, index, columnType)
	case vectorTypeConst:
		r := newBytesReader(v.VectorData)
		return c.decodeConstValue(dctx, r, columnType.getType())
	default:
		return nil, errInvalidVectorType
	}
}

func (c *vectorDecoder) decodeFlatValue(dctx *decodeContext, v *vector.NestedVector,
	index uint32, columnType typeSchema) (*nebulaValue, error) {
	typ := columnType.getType()
	if fn, ok := c.decodeFlatFns[typ]; ok {
		return fn(dctx, v, index, columnType)
	} else {
		return nil, errors.Wrap(errInvalidColumnType, "")
	}
}

func (c *vectorDecoder) decodeConstValue(dctx *decodeContext, r *bytesReader, columnType types.ColumnType) (*nebulaValue, error) {
	typ := columnType
	offset := dctx.timezoneOffset
	switch typ {
	case types.ColumnTypeUnknown:
		return &nebulaValue{data: nil}, nil
	case types.ColumnTypeBool:
		fallthrough
	case types.ColumnTypeInt8, types.ColumnTypeUint8, types.ColumnTypeInt16, types.ColumnTypeUint16, types.ColumnTypeInt32, types.ColumnTypeUint32, types.ColumnTypeInt64, types.ColumnTypeUint64:
		fallthrough
	case types.ColumnTypeFloat32, types.ColumnTypeFloat64:
		fallthrough
	case types.ColumnTypeDate, types.ColumnTypeLocalTime, types.ColumnTypeZonedTime, types.ColumnTypeLocalDatetime, types.ColumnTypeZonedDatetime, types.ColumnTypeDuration:
		bs := r.readN(sizeMap[typ])
		if r.error() != nil {
			return nil, r.error()
		}
		return &nebulaValue{
			data: decodeBasicValue(bs, typ, offset),
		}, nil
	case types.ColumnTypeString, types.ColumnTypeDecimal:
		var bs []byte
		bs = r.readUtilZero()
		// no zero fount, return all
		if r.error() != nil {
			bs = r.readPeddingAll()
		}
		if typ == types.ColumnTypeDecimal {
			return &nebulaValue{
				data: &NebulaDecimal{Sval: string(bs)},
			}, nil
		} else {
			return &nebulaValue{
				data: decodeBasicValue(bs, types.ColumnTypeString, offset),
			}, nil
		}
	case types.ColumnTypeList:
		// size + [value]
		bs := r.readN(4)
		if r.error() != nil {
			return nil, r.error()
		}
		listSize := bytesToInt32(bs)
		ll := &NebulaList{
			Values: make([]*nebulaValue, 0, listSize),
		}
		for i := int32(0); i < listSize; i++ {
			bs := r.readN(1)
			if r.error() != nil {
				return nil, r.error()
			}
			subColumnType, ok := columnTypeMap[bs[0]]
			if !ok {
				return nil, errors.Wrap(errInvalidColumnType, "")
			}
			v, err := c.decodeConstValue(dctx, r, subColumnType)
			if err != nil {
				return nil, err
			}
			ll.Values = append(ll.Values, v)
		}
		return &nebulaValue{
			data: ll,
		}, nil
	case types.ColumnTypeRecord:
		// prop num + [prop name] + [type + value]
		bs := r.readN(4)
		if r.error() != nil {
			return nil, r.error()
		}
		propNum := bytesToInt32(bs)
		nameIndex := make([]string, 0, propNum)
		mm := &NebulaRecord{
			Values: make(map[string]*nebulaValue),
		}
		for i := int32(0); i < propNum; i++ {
			name := r.readUtilZero()
			if r.error() != nil {
				return nil, r.error()
			}
			nameIndex = append(nameIndex, string(name))
		}
		for i := int32(0); i < propNum; i++ {
			bs := r.readN(1)
			if r.error() != nil {
				return nil, r.error()
			}
			subColumnType, ok := columnTypeMap[bs[0]]
			if !ok {
				return nil, errors.Wrap(errInvalidColumnType, "")
			}

			v, err := c.decodeConstValue(dctx, r, subColumnType)
			if err != nil {
				return nil, err
			}
			mm.Values[nameIndex[i]] = v
		}
		return &nebulaValue{
			data: mm,
		}, nil
	default:
		return nil, errors.Wrap(errInvalidColumnType, fmt.Sprintf("%d", typ))
	}
}

func (c *vectorDecoder) decodeStringValue() decodeFlatFn {
	return func(dctx *decodeContext, v *vector.NestedVector, index uint32, columnType typeSchema) (*nebulaValue, error) {
		// uint32 length + prefix string + uint32 chunk offset + uint32 chunk index
		// if the string is less then 12 bytes, no need to get data from chunk
		length := 4 + 4 + 4 + 4
		header := v.VectorData[index*uint32(length) : index*(uint32(length))+uint32(length)]
		strLen := bytesToUint32(header[:4])

		if strLen <= 12 {
			return &nebulaValue{
				data: &NebulaString{Value: string(header[4 : 4+strLen])},
			}, nil
		}
		chunkIndex := bytesToUint32(header[12:16])
		chunkOffset := bytesToUint32(header[8:12])
		chunk := v.NestedVectors[chunkIndex]
		data := chunk.VectorData[chunkOffset : chunkOffset+strLen]
		offset := dctx.timezoneOffset
		return &nebulaValue{
			data: &NebulaString{Value: decodeBasicValue(data, types.ColumnTypeString, offset).(*NebulaString).Value},
		}, nil
	}
}

func (c *vectorDecoder) decodeDecimalValue() decodeFlatFn {
	// same with string
	return func(dctx *decodeContext, v *vector.NestedVector, index uint32, columnType typeSchema) (*nebulaValue, error) {
		// uint32 length + prefix string + uint32 chunk offset + uint32 chunk index
		// if the string is less then 12 bytes, no need to get data from chunk
		length := 4 + 4 + 4 + 4
		header := v.VectorData[index*uint32(length) : index*(uint32(length))+uint32(length)]
		strLen := bytesToUint32(header[:4])

		if strLen <= 12 {
			return &nebulaValue{
				data: &NebulaDecimal{Sval: string(header[4 : 4+strLen])},
			}, nil
		}
		chunkIndex := bytesToUint32(header[12:16])
		chunkOffset := bytesToUint32(header[8:12])
		chunk := v.NestedVectors[chunkIndex]
		data := chunk.VectorData[chunkOffset : chunkOffset+strLen]
		offset := dctx.timezoneOffset
		return &nebulaValue{
			data: &NebulaDecimal{Sval: decodeBasicValue(data, types.ColumnTypeString, offset).(*NebulaString).Value},
		}, nil
	}
}

func (c *vectorDecoder) decodeBasiceValue(t types.ColumnType) decodeFlatFn {
	return func(dctx *decodeContext, v *vector.NestedVector, index uint32, columnType typeSchema) (*nebulaValue, error) {
		bs, err := c.getCommonBytes(v, t, index)
		offset := dctx.timezoneOffset
		if err != nil {
			return nil, err
		}
		return &nebulaValue{
			data: decodeBasicValue(bs, t, offset),
		}, nil
	}
}

func (c *vectorDecoder) decodeListValue() decodeFlatFn {
	return func(dctx *decodeContext, v *vector.NestedVector, index uint32, columnType typeSchema) (*nebulaValue, error) {
		// offset + size
		length := 4 + 4
		header := v.VectorData[index*uint32(length) : index*(uint32(length))+uint32(length)]
		offset := bytesToUint32(header[:4])
		size := bytesToUint32(header[4:8])
		l := &NebulaList{
			Values: make([]*nebulaValue, 0, size),
		}
		schema, ok := columnType.(*columnTypeSchemaList)
		if !ok {
			return nil, errTypeAssertion
		}
		for i := uint32(0); i < size; i++ {
			vv, err := c.decodeValue(dctx, v.NestedVectors[0], vectorTypeFlat, offset+i, schema.subSchema)
			if err != nil {
				return nil, err
			}
			l.Values = append(l.Values, vv)
		}
		return &nebulaValue{
			data: l,
		}, nil
	}
}

func (c *vectorDecoder) decodeRecordValue() decodeFlatFn {
	return func(dctx *decodeContext, v *vector.NestedVector, index uint32, columnType typeSchema) (*nebulaValue, error) {
		r := newBytesReader(v.SpecialMetaData)
		schema, ok := columnType.(*columnTypeSchemaRecord)
		if !ok {
			return nil, errTypeAssertion
		}
		propSchemas := schema.propSchemas
		nameIndexs := make([]string, 0, len(propSchemas))
		for i := 0; i < len(propSchemas); i++ {
			bs := r.readUtilZero()
			if r.error() != nil {
				return nil, r.error()
			}
			nameIndexs = append(nameIndexs, string(bs))
		}
		record := &NebulaRecord{
			Values: make(map[string]*nebulaValue),
		}
		for i := 0; i < len(propSchemas); i++ {
			n := nameIndexs[i]
			s, ok := propSchemas[n]
			if !ok {
				return nil, errPropNotFound
			}
			value, err := c.decodeValue(dctx, v.NestedVectors[i], vectorTypeFlat, index, s)
			if err != nil {
				return nil, err
			}
			record.Values[n] = value
		}
		return &nebulaValue{
			data: record,
		}, nil

	}
}

func (c *vectorDecoder) decodeNullValue() decodeFlatFn {
	return func(dctx *decodeContext, v *vector.NestedVector, index uint32, columnType typeSchema) (*nebulaValue, error) {
		return &nebulaValue{data: nil}, nil
	}
}

func (c *vectorDecoder) decodeNodeValue() decodeFlatFn {
	return func(dctx *decodeContext, v *vector.NestedVector, index uint32, columnType typeSchema) (*nebulaValue, error) {
		nodeSchema, ok := columnType.(*columnTypeSchemaElement)
		if !ok {
			return nil, errTypeAssertion
		}
		allNodeProps := nodeSchema.elemenetProps
		if err := decodePropVectorIndex(allNodeProps, v.SpecialMetaData, true); err != nil {
			return nil, err
		}
		// header = nodeID + graphId + padding
		length := 8 + 4 + 4
		header := v.VectorData[index*uint32(length) : index*(uint32(length))+uint32(length)]

		nodeID := bytesToInt64(header[:8])
		graphID := bytesToInt32(header[8:12])
		nodeTypeID := (int32)(nodeID >> 48)
		nodeProps, ok := allNodeProps[nodeTypeID]
		if !ok {
			return nil, errors.Wrap(errElementTypeNotFount, "")
		}
		gsm := dctx.graphsSchema
		graphName, typeName, labels, err := getSchemaName(gsm, graphID, nodeTypeID, true)
		if err != nil {
			return nil, err
		}
		node := &NebulaNode{
			NodeId: nodeID,
			Graph:  graphName,
			Type:   typeName,
			Labels: labels,
		}
		node.Properties = make(map[string]*nebulaValue)
		for _, prop := range nodeProps {
			vw := v.NestedVectors[prop.vectorIndex]
			value, err := c.decodeValue(dctx, vw, vectorTypeFlat, index, prop.typ)
			if err != nil {
				return nil, err
			}
			// TODO: check if the prop name is valid
			node.Properties[prop.name] = value
		}
		return &nebulaValue{
			data: node,
		}, nil

	}
}

func (c *vectorDecoder) decodeEdgeValue() decodeFlatFn {
	return func(dctx *decodeContext, v *vector.NestedVector, index uint32, columnType typeSchema) (*nebulaValue, error) {
		schema, ok := columnType.(*columnTypeSchemaElement)
		if !ok {
			return nil, errTypeAssertion
		}
		allGraphElementProps := schema.elemenetProps
		if err := decodePropVectorIndex(allGraphElementProps, v.SpecialMetaData, false); err != nil {
			return nil, err
		}
		// header = src id + dst id + rank + graph id + edge type id
		length := 8 + 8 + 8 + 4 + 4
		header := v.VectorData[index*uint32(length) : index*(uint32(length))+uint32(length)]

		srcID := bytesToInt64(header[:8])
		dstID := bytesToInt64(header[8:16])
		rank := bytesToInt64(header[16:24])
		graphID := bytesToInt32(header[24:28])
		edgeTypeID := bytesToInt32(header[28:32])
		noDirectType := edgeTypeID & 0x3FFFFFFF
		direction := getEdgeDirection(uint8(edgeTypeID >> 30))
		props, ok := allGraphElementProps[noDirectType]
		if !ok {
			return nil, errors.Wrap(errElementTypeNotFount, "")
		}
		gsm := dctx.graphsSchema
		graphName, typeName, labels, err := getSchemaName(gsm, graphID, noDirectType, false)
		if err != nil {
			return nil, err
		}
		e := &NebulaEdge{
			Rank:      rank,
			Graph:     graphName,
			Type:      typeName,
			Labels:    labels,
			Direction: direction,
		}
		switch direction {
		case EdgeInComingDirection:
			e.SrcId = dstID
			e.DstId = srcID
		default:
			e.SrcId = srcID
			e.DstId = dstID
		}
		e.Properties = make(map[string]*nebulaValue)
		for _, prop := range props {
			vw := v.NestedVectors[prop.vectorIndex]
			value, err := c.decodeValue(dctx, vw, vectorTypeFlat, index, prop.typ)
			if err != nil {
				return nil, err
			}
			// TODO: check if the prop name is valid
			e.Properties[prop.name] = value
		}
		return &nebulaValue{
			data: e,
		}, nil
	}
}

func (c *vectorDecoder) decodePathValue() decodeFlatFn {
	return func(dctx *decodeContext, v *vector.NestedVector, index uint32, columnType typeSchema) (*nebulaValue, error) {
		schema, ok := columnType.(*columnTypeSchemaPath)
		if !ok {
			return nil, errTypeAssertion
		}

		meta := &pathMetaData{
			nodeTypeIndex: make(pathElementTypeIndexMap),
			edgeTypeIndex: make(pathElementTypeIndexMap),
			nodeIndexPair: make(map[pathPairIndex]*pathPair),
			edgeIndexPair: make(map[pathPairIndex]*pathPair),
		}
		if err := decodePathSpecialData(v, meta); err != nil {
			return nil, err
		}
		nodeSchema := schema.nodeSchema
		edgeSchema := schema.edgeSchema
		// header = headNodeId + tailNodeId + totalNum + edgeNum + headoffset + tailoffset
		length := 8 + 8 + 4 + 4 + 4 + 4
		header := v.VectorData[index*uint32(length) : index*(uint32(length))+uint32(length)]
		headNodeID := bytesToInt64(header[:8])
		tailNodeID := bytesToInt64(header[8:16])
		totalNum := bytesToInt32(header[16:20])
		edgeNum := bytesToInt32(header[20:24])
		headOffset := bytesToUint32(header[24:28])
		tailOffset := bytesToInt32(header[28:32])
		_, _, _ = edgeNum, tailNodeID, tailOffset

		path := &NebulaPath{
			Values: make([]*nebulaValue, 0, totalNum),
		}
		if totalNum == 0 {
			return &nebulaValue{
				data: path,
			}, nil
		}
		var (
			dataVector           *vector.NestedVector
			adjVector            *vector.NestedVector
			pathHeaderValue      *nebulaValue
			pathHeaderValueInt64 types.Int64
			pathHeader           *pathAdjHeader
			err                  error
		)
		headNodeType := int32(headNodeID >> 48)
		pairIndex, ok := meta.nodeTypeIndex[headNodeType]
		if !ok {
			return nil, errors.Wrap(errElementTypeNotFount, "")
		}
		sentinelPair, ok := meta.nodeIndexPair[pairIndex]
		if !ok {
			return nil, errors.Wrap(errElementTypeNotFount, "")
		}
		sentinelType := types.ColumnTypeNode
		sentinelOffset := headOffset
		pathHeaderSchema := &columnTypeSchemaBasic{
			typ: types.ColumnTypeInt64,
		}
		for i := int32(0); i < totalNum; i++ {
			dataVector = sentinelPair.cur
			adjVector = sentinelPair.adj

			var data *nebulaValue
			if sentinelType == types.ColumnTypeNode {
				data, err = c.decodeValue(dctx, dataVector, vectorTypeFlat, sentinelOffset, nodeSchema)
			} else {
				data, err = c.decodeValue(dctx, dataVector, vectorTypeFlat, sentinelOffset, edgeSchema)
			}
			if err != nil {
				return nil, err
			}
			path.Values = append(path.Values, data)
			if len(path.Values) == int(totalNum) {
				break
			}
			pathHeaderValue, err = c.decodeValue(dctx, adjVector, vectorTypeFlat, sentinelOffset, pathHeaderSchema)
			if err != nil {
				return nil, err
			}
			pathHeaderValueInt64, err = pathHeaderValue.AsInt64()
			if err != nil {
				return nil, err
			}
			pathHeader = newPathAdjHeader(int64(pathHeaderValueInt64))
			sentinelOffset = pathHeader.nextOffset
			if pathHeader.nextIsEdge {
				sentinelType = types.ColumnTypeEdge
			} else {
				sentinelType = types.ColumnTypeNode
			}
			if sentinelType == types.ColumnTypeNode {
				sentinelPair, ok = meta.nodeIndexPair[pathPairIndex(pathHeader.nextVectorIndex)]
				if !ok {
					return nil, errors.Wrap(errElementTypeNotFount, "")
				}
			} else {
				sentinelPair, ok = meta.edgeIndexPair[pathPairIndex(pathHeader.nextVectorIndex)]
				if !ok {
					return nil, errors.Wrap(errElementTypeNotFount, "")
				}
			}
		}
		return &nebulaValue{
			data: path,
		}, nil
	}
}

func decodePropVectorIndex(elemmentTypes graphElementProps, bs []byte, isNode bool) error {
	// properties num + [prop names] + element type num + [element type id + prop num + [vector index]]
	// vector index is same as the order of prop names
	elementTypeSize := 2
	if !isNode {
		elementTypeSize = 4
	}
	var propList = make([]string, 0)
	r := newBytesReader(bs)
	propsNumBytes := r.readN(4)
	if r.error() != nil {
		return r.error()
	}
	propsNum := bytesToInt32(propsNumBytes)
	for i := 0; i < int(propsNum); i++ {
		propNameBytes := r.readUtilZero()
		if r.error() != nil {
			return r.error()
		}
		propList = append(propList, string(propNameBytes))
	}
	nodeTypeNumBytes := r.readN(4)
	if r.error() != nil {
		return r.error()
	}
	nodeTypeNum := bytesToInt32(nodeTypeNumBytes)
	for i := 0; i < int(nodeTypeNum); i++ {
		elementTypeIdBytes := r.readN(elementTypeSize)

		propNumBytes := r.readN(4)
		if r.error() != nil {
			return r.error()
		}
		var elementId int32
		if isNode {
			elementId = int32(bytesToInt16(elementTypeIdBytes))
		} else {
			elementId = bytesToInt32(elementTypeIdBytes)
		}
		propNum := int(bytesToInt32(propNumBytes))
		for j := 0; j < propNum; j++ {
			vectorIndexBytes := r.readN(4)
			if r.error() != nil {
				return r.error()
			}
			vectorIndex := bytesToInt32(vectorIndexBytes)
			nodeType, ok := elemmentTypes[elementId]
			if !ok {
				return errors.Wrap(errElementTypeNotFount, "")
			}
			prop, ok := nodeType[propList[vectorIndex]]
			if !ok {
				return errPropNotFound
			}
			prop.vectorIndex = int(vectorIndex)
		}
	}
	return nil
}

func decodePathSpecialData(v *vector.NestedVector, meta *pathMetaData) error {
	// special meta data
	// num of node type + [node type id + pair index] + num of edge type + [edge type id + pair index]
	r := newBytesReader(v.SpecialMetaData)
	nodeTypeNumBytes := r.readN(4)
	if r.error() != nil {
		return r.error()
	}
	nodeTypeNum := bytesToInt32(nodeTypeNumBytes)
	nestedVectorIndex := 0
	for i := 0; i < int(nodeTypeNum); i++ {
		nodeTypeIdBytes := r.readN(2)
		pairIndexBytes := r.readN(2)
		if r.error() != nil {
			return r.error()
		}
		nodeTypeId := int32(bytesToUint16(nodeTypeIdBytes))
		pairIndex := bytesToUint16(pairIndexBytes)
		meta.nodeTypeIndex[nodeTypeId] = pathPairIndex(pairIndex)
		meta.nodeIndexPair[pathPairIndex(pairIndex)] = &pathPair{
			cur: v.NestedVectors[nestedVectorIndex],
			adj: v.NestedVectors[nestedVectorIndex+1],
		}
		nestedVectorIndex += 2
	}
	edgeTypeNumBytes := r.readN(4)
	if r.error() != nil {
		return r.error()
	}
	edgeTypeNum := bytesToInt32(edgeTypeNumBytes)
	for i := 0; i < int(edgeTypeNum); i++ {
		edgeTypeIdBytes := r.readN(4)
		pairIndexBytes := r.readN(2)
		if r.error() != nil {
			return r.error()
		}
		edgeTypeId := bytesToInt32(edgeTypeIdBytes)
		pairIndex := bytesToInt16(pairIndexBytes)
		meta.edgeTypeIndex[edgeTypeId] = pathPairIndex(pairIndex)
		meta.edgeIndexPair[pathPairIndex(pairIndex)] = &pathPair{
			cur: v.NestedVectors[nestedVectorIndex],
			adj: v.NestedVectors[nestedVectorIndex+1],
		}
		nestedVectorIndex += 2
	}
	return nil
}

func (c *vectorDecoder) decodeAnyValue() decodeFlatFn {
	return func(dctx *decodeContext, v *vector.NestedVector, index uint32, columnType typeSchema) (*nebulaValue, error) {
		dataTypeVector := v.NestedVectors[0]
		typeLength := uint32(1)
		dateHeaderLength := uint32(8)
		typeBytes := dataTypeVector.VectorData[typeLength*index : typeLength*(index+1)]
		dataType, ok := columnTypeMap[typeBytes[0]]
		offset := dctx.timezoneOffset
		if !ok {
			return nil, errors.Wrap(errInvalidColumnType, "")
		}
		dataBytes := v.VectorData[dateHeaderLength*index : dateHeaderLength*(index+1)]
		if isBasicColumnType(dataType) {
			v := decodeBasicValue(dataBytes, dataType, offset)
			return &nebulaValue{
				data: v,
			}, nil
		}
		chunkIndexBytes, chunkOffsetBytes := dataBytes[:4], dataBytes[4:8]
		chunkIndex := bytesToUint32(chunkIndexBytes)
		chunkOffset := bytesToUint32(chunkOffsetBytes)
		r := newBytesReader(v.NestedVectors[chunkIndex+1].VectorData)
		r.index = int(chunkOffset)

		return decodeAnyCompositeValue(dctx, r, dataType, true)
	}
}

func decodeBasicValue(bs []byte, typ types.ColumnType, offset int64) valuer {
	switch typ {
	case types.ColumnTypeBool:
		return &NebulaBool{Value: bs[0] == 1}
	case types.ColumnTypeInt8:
		return &NebulaInt8{Value: bytesToInt8(bs[:1])}
	case types.ColumnTypeInt16:
		return &NebulaInt16{Value: bytesToInt16(bs[:2])}
	case types.ColumnTypeInt32:
		return &NebulaInt32{Value: bytesToInt32(bs[:4])}
	case types.ColumnTypeInt64:
		return &NebulaInt64{Value: bytesToInt64(bs[:8])}
	case types.ColumnTypeUint8:
		return &NebulaUint8{Value: bytesToUint8(bs[:1])}
	case types.ColumnTypeUint16:
		return &NebulaUint16{Value: bytesToUint16(bs[:2])}
	case types.ColumnTypeUint32:
		return &NebulaUint32{Value: bytesToUint32(bs[:4])}
	case types.ColumnTypeUint64:
		return &NebulaUint64{Value: bytesToUint64(bs[:8])}
	case types.ColumnTypeFloat32:
		return &NebulaFloat{Value: math.Float32frombits(order.Uint32(bs[:4]))}
	case types.ColumnTypeFloat64:
		return &NebulaDouble{Value: math.Float64frombits(order.Uint64(bs[:8]))}
	case types.ColumnTypeString:
		return &NebulaString{Value: string(bs)}
	case types.ColumnTypeDecimal:
		return &NebulaDecimal{Sval: string(bs)}
	case types.ColumnTypeDate:
		year := int16(bytesToInt16(bs[:2]))
		month := bytesToInt8(bs[2:3])
		day := bytesToInt8(bs[3:4])
		return &NebulaDate{
			Year:  year,
			Month: month,
			Day:   day,
		}
	case types.ColumnTypeLocalTime:
		hour := bytesToInt8(bs[:1])
		minute := bytesToInt8(bs[1:2])
		second := bytesToInt8(bs[2:3])
		// pedding
		microsecond := bytesToInt32(bs[4:8])
		return &NebulaLocalTime{
			Hour:     hour,
			Minute:   minute,
			Sec:      second,
			Microsec: microsecond,
		}
	case types.ColumnTypeZonedTime:
		hour := bytesToInt8(bs[:1])
		minute := bytesToInt8(bs[1:2])
		second := bytesToInt8(bs[2:3])
		// pedding
		microsecond := bytesToInt32(bs[4:8])
		n := time.Now()
		var h int
		if hour < 0 {
			h = -int(hour)
		} else {
			h = int(hour)
		}
		t := time.Date(n.Year(), n.Month(), n.Day(), h, int(minute), int(second), int(microsecond)*int(time.Microsecond), time.UTC)
		zonedT := t.In(time.FixedZone("", int(offset)))
		return &NebulaZonedTime{
			Hour:     int8(zonedT.Hour()),
			Minute:   int8(zonedT.Minute()),
			Sec:      int8(zonedT.Second()),
			Microsec: int32(zonedT.Nanosecond() / int(time.Microsecond)),
			Offset:   int32(offset),
		}
	case types.ColumnTypeLocalDatetime:
		qword := bytesToInt64(bs)
		year := int16(qword & 0xffff)
		month := int8(qword >> 16 & 0xf)
		day := int8(qword >> 20 & 0x1f)
		hour := int8(qword >> 25 & 0x1f)
		minute := int8(qword >> 30 & 0x3f)
		second := int8(qword >> 36 & 0x3f)
		microsecond := int32(qword >> 42 & 0x3ffffff)
		return &NebulaLocalDatetime{
			Year:     year,
			Month:    month,
			Day:      day,
			Hour:     hour,
			Minute:   minute,
			Sec:      second,
			Microsec: microsecond,
		}
	case types.ColumnTypeZonedDatetime:
		qword := bytesToInt64(bs)
		year := int16(qword & 0xffff)
		month := int8(qword >> 16 & 0xf)
		day := int8(qword >> 20 & 0x1f)
		hour := int8(qword >> 25 & 0x1f)
		minute := int8(qword >> 30 & 0x3f)
		second := int8(qword >> 36 & 0x3f)
		microsecond := int32(qword >> 42 & 0x3ffffff)
		t := time.Date(int(year), time.Month(month), int(day), int(hour), int(minute),
			int(second), int(microsecond)*int(time.Microsecond), time.UTC)
		zonedT := t.In(time.FixedZone("", int(offset)))
		return &NebulaZonedDatetime{
			Year:     int16(zonedT.Year()),
			Month:    int8(zonedT.Month()),
			Day:      int8(zonedT.Day()),
			Hour:     int8(zonedT.Hour()),
			Minute:   int8(zonedT.Minute()),
			Sec:      int8(zonedT.Second()),
			Microsec: int32(zonedT.Nanosecond() / int(time.Microsecond)),
			Offset:   int32(offset),
		}
	case types.ColumnTypeDuration:
		value := bytesToInt64(bs)
		isMonthBased := value&0x1 == 1
		value >>= 1
		var (
			year        int64
			month       int8
			day         int32
			hour        int8
			minute      int8
			second      int8
			microsecond int32
		)
		if isMonthBased {
			year = value / 12
			month = int8(value % 12)
		} else {
			day = int32(value / kMicrosecondsOfDay)
			hour = int8(value % kMicrosecondsOfDay / kMicrosecondsOfHour)
			minute = int8(value % kMicrosecondsOfHour / kMicrosecondsOfMinute)
			second = int8(value % kMicrosecondsOfMinute / kMicrosecondsOfSecond)
			microsecond = int32(value % kMicrosecondsOfSecond)
		}
		return &NebulaDuration{
			isMonthBased: isMonthBased,
			Year:         year,
			Month:        month,
			Day:          day,
			Hour:         hour,
			Minute:       minute,
			Sec:          second,
			Microsec:     microsecond,
		}
	default:
		return nil
	}
}

func decodeAnyCompositeValue(dctx *decodeContext, r *bytesReader, typ types.ColumnType, first bool) (*nebulaValue, error) {
	if first {
		// first byte is the type
		if isBasicColumnType(typ) {
			return nil, errors.Wrap(errInvalidColumnType, "")
		}
	}
	offset := dctx.timezoneOffset
	gsm := dctx.graphsSchema
	if isBasicColumnType(typ) {
		size, ok := sizeMap[typ]
		if !ok {
			return nil, errors.Wrap(errInvalidColumnType, "")
		}
		bs := r.readN(size)
		if r.error() != nil {
			return nil, r.error()
		}
		v := decodeBasicValue(bs, typ, offset)
		return &nebulaValue{
			data: v,
		}, nil
	}
	switch typ {
	case types.ColumnTypeUnknown:
		return &nebulaValue{data: nil}, nil
	case types.ColumnTypeString, types.ColumnTypeDecimal:
		sizeBytes := r.readN(2)
		if r.error() != nil {
			return nil, r.error()
		}
		size := int(bytesToInt16(sizeBytes))
		bs := r.readN(size)
		if r.error() != nil {
			return nil, r.error()
		}
		if typ == types.ColumnTypeDecimal {
			return &nebulaValue{
				data: &NebulaDecimal{Sval: string(bs)},
			}, nil
		} else {
			return &nebulaValue{
				data: &NebulaString{Value: string(bs)},
			}, nil
		}
	case types.ColumnTypeList:
		typeBytes := r.readN(1)
		sizeBytes := r.readN(2)
		if r.error() != nil {
			return nil, r.error()
		}
		subType, ok := columnTypeMap[typeBytes[0]]
		if !ok {
			return nil, errors.Wrap(errInvalidColumnType, "")
		}
		size := int(bytesToInt16(sizeBytes))
		l := make([]*nebulaValue, 0, size)
		var bitSize int
		if size%8 != 0 {
			bitSize = size/8 + 1
		} else {
			bitSize = size / 8
		}
		nullBitByte := r.readN(bitSize)
		if r.error() != nil {
			return nil, r.error()
		}
		if r.error() != nil {
			return nil, r.error()
		}
		for i := 0; i < size; i++ {
			if nullBitByte[i/8]&(1<<(i%8)) == 0 {
				l = append(l, &nebulaValue{data: nil})
			} else {
				v, err := decodeAnyCompositeValue(dctx, r, subType, false)
				if err != nil {
					return nil, err
				}
				l = append(l, v)
			}
		}
		return &nebulaValue{
			data: &NebulaList{
				Values: l,
			},
		}, nil
	case types.ColumnTypeRecord:
		sizeBytes := r.readN(2)
		size := int(bytesToInt16(sizeBytes))
		m := make(map[string]*nebulaValue, 0)
		for i := 0; i < size; i++ {
			bs := r.readN(4)
			if r.error() != nil {
				return nil, r.error()
			}
			nameLength := int(bytesToInt32(bs))
			nameBytes := r.readN(nameLength)
			typeBytes := r.readN(1)
			if r.error() != nil {
				return nil, r.error()
			}
			subType, ok := columnTypeMap[typeBytes[0]]
			if !ok {
				return nil, errors.Wrap(errInvalidColumnType, "")
			}
			v, err := decodeAnyCompositeValue(dctx, r, subType, false)
			if err != nil {
				return nil, err
			}
			m[string(nameBytes)] = v
		}
		return &nebulaValue{
			data: &NebulaRecord{
				Values: m,
			},
		}, nil
	case types.ColumnTypeNode:
		// nodeID 8B + graphId 4B + prop_Size 2B
		nodeIdBytes, graphIdBytes, propSizeBytes := r.readN(8), r.readN(4), r.readN(2)
		if r.error() != nil {
			return nil, r.error()
		}
		nodeId := bytesToInt64(nodeIdBytes)
		nodeTypeId := int32(nodeId >> 48)
		graphId := bytesToInt32(graphIdBytes)
		_ = graphId
		propSize := bytesToInt16(propSizeBytes)
		props := make(map[string]*nebulaValue, 0)
		for i := 0; i < int(propSize); i++ {
			nameSizeBytes := r.readN(2)
			if r.error() != nil {
				return nil, r.error()
			}
			nameSize := int(bytesToInt16(nameSizeBytes))
			nameBytes := r.readN(nameSize)
			typeBytes := r.readN(1)
			if r.error() != nil {
				return nil, r.error()
			}
			subType, ok := columnTypeMap[typeBytes[0]]
			if !ok {
				return nil, errors.Wrap(errInvalidColumnType, "")
			}
			v, err := decodeAnyCompositeValue(dctx, r, subType, false)
			if err != nil {
				return nil, err
			}
			props[string(nameBytes)] = v
		}
		graphName, typeName, labels, err := getSchemaName(gsm, graphId, nodeTypeId, true)
		if err != nil {
			return nil, err
		}
		return &nebulaValue{
			data: &NebulaNode{
				NodeId:     nodeId,
				Graph:      graphName,
				Type:       typeName,
				Labels:     labels,
				Properties: props,
			},
		}, nil
	case types.ColumnTypeEdge:
		// src nodeID 8B + dst nodeID 8B + graphId 4B + edge type ID 4B + edge rank 8B + prop_size 2B
		srcNodeIDBytes, dstNodeIDBytes, graphIdBytes, edgeTypeIdBytes,
			edgeRankBytes, propSizeBytes := r.readN(8), r.readN(8), r.readN(4), r.readN(4), r.readN(8), r.readN(2)
		if r.error() != nil {
			return nil, r.error()
		}
		srcNodeId := bytesToInt64(srcNodeIDBytes)
		dstNodeId := bytesToInt64(dstNodeIDBytes)
		graphId := bytesToInt32(graphIdBytes)
		edgeTypeID := bytesToInt32(edgeTypeIdBytes)
		edgeRank := bytesToInt64(edgeRankBytes)
		propSize := bytesToInt16(propSizeBytes)
		props := make(map[string]*nebulaValue, 0)
		noDirectType := edgeTypeID & 0x3FFFFFFF
		direction := getEdgeDirection(uint8(edgeTypeID >> 30))
		for i := 0; i < int(propSize); i++ {
			nameSizeBytes := r.readN(4)
			if r.error() != nil {
				return nil, r.error()
			}
			nameSize := int(bytesToInt32(nameSizeBytes))
			nameBytes := r.readN(nameSize)
			typeBytes := r.readN(1)
			if r.error() != nil {
				return nil, r.error()
			}
			subType, ok := columnTypeMap[typeBytes[0]]
			if !ok {
				return nil, errors.Wrap(errInvalidColumnType, "")
			}
			v, err := decodeAnyCompositeValue(dctx, r, subType, false)
			if err != nil {
				return nil, err
			}
			props[string(nameBytes)] = v
		}
		graphName, typeName, labels, err := getSchemaName(gsm, graphId, noDirectType, false)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		e := &NebulaEdge{
			Rank:       edgeRank,
			Graph:      graphName,
			Type:       typeName,
			Labels:     labels,
			Properties: props,
			Direction:  direction,
		}
		switch direction {
		case EdgeInComingDirection:
			e.SrcId = dstNodeId
			e.DstId = srcNodeId
		default:
			e.SrcId = srcNodeId
			e.DstId = dstNodeId
		}
		return &nebulaValue{
			data: e,
		}, nil
	case types.ColumnTypePath:
		elementNumBytes := r.readN(2)
		if r.error() != nil {
			return nil, r.error()
		}
		elementNum := int(bytesToInt16(elementNumBytes))
		p := &NebulaPath{
			Values: make([]*nebulaValue, 0, elementNum),
		}
		for i := 0; i < elementNum; i++ {
			subTypeBytes := r.readN(1)
			if r.error() != nil {
				return nil, r.error()
			}
			subType, ok := columnTypeMap[subTypeBytes[0]]
			if !ok {
				return nil, errors.Wrap(errInvalidColumnType, "")
			}
			element, err := decodeAnyCompositeValue(dctx, r, subType, false)
			if err != nil {
				return nil, err
			}
			p.Values = append(p.Values, element)
		}
		return &nebulaValue{
			data: p,
		}, nil
	default:
		return nil, errors.Wrap(errInvalidColumnType, "")
	}
}

func init() {
	d := defaultDecoder.(*vectorDecoder)
	d.decodeFlatFns = make(map[types.ColumnType]decodeFlatFn)
	d.decodeFlatFns = map[types.ColumnType]decodeFlatFn{
		types.ColumnTypeBool:          d.decodeBasiceValue(types.ColumnTypeBool),
		types.ColumnTypeInt8:          d.decodeBasiceValue(types.ColumnTypeInt8),
		types.ColumnTypeInt16:         d.decodeBasiceValue(types.ColumnTypeInt16),
		types.ColumnTypeInt32:         d.decodeBasiceValue(types.ColumnTypeInt32),
		types.ColumnTypeInt64:         d.decodeBasiceValue(types.ColumnTypeInt64),
		types.ColumnTypeUint8:         d.decodeBasiceValue(types.ColumnTypeUint8),
		types.ColumnTypeUint16:        d.decodeBasiceValue(types.ColumnTypeUint16),
		types.ColumnTypeUint32:        d.decodeBasiceValue(types.ColumnTypeUint32),
		types.ColumnTypeUint64:        d.decodeBasiceValue(types.ColumnTypeUint64),
		types.ColumnTypeFloat32:       d.decodeBasiceValue(types.ColumnTypeFloat32),
		types.ColumnTypeFloat64:       d.decodeBasiceValue(types.ColumnTypeFloat64),
		types.ColumnTypeString:        d.decodeStringValue(),
		types.ColumnTypeNode:          d.decodeNodeValue(),
		types.ColumnTypeEdge:          d.decodeEdgeValue(),
		types.ColumnTypePath:          d.decodePathValue(),
		types.ColumnTypeUnknown:       d.decodeNullValue(),
		types.ColumnTypeList:          d.decodeListValue(),
		types.ColumnTypeRecord:        d.decodeRecordValue(),
		types.ColumnTypeLocalTime:     d.decodeBasiceValue(types.ColumnTypeLocalTime),
		types.ColumnTypeLocalDatetime: d.decodeBasiceValue(types.ColumnTypeLocalDatetime),
		types.ColumnTypeZonedTime:     d.decodeBasiceValue(types.ColumnTypeZonedTime),
		types.ColumnTypeZonedDatetime: d.decodeBasiceValue(types.ColumnTypeZonedDatetime),
		types.ColumnTypeDate:          d.decodeBasiceValue(types.ColumnTypeDate),
		types.ColumnTypeDuration:      d.decodeBasiceValue(types.ColumnTypeDuration),
		types.ColumnTypeDecimal:       d.decodeDecimalValue(),
		types.ColumnTypeAny:           d.decodeAnyValue(),
	}
}
