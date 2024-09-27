package decode

import (
	"io"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/vector"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/types"
)

var columnTypeMap = map[uint8]types.ColumnType{
	0x1: types.ColumnTypeNode,
	0x2: types.ColumnTypeEdge,
	0x3: types.ColumnTypeUnknown,
	0x4: types.ColumnTypeBool,
	0x5: types.ColumnTypeInt8,
	0x6: types.ColumnTypeUint8,
	0x7: types.ColumnTypeInt16,
	0x8: types.ColumnTypeUint16,
	0x9: types.ColumnTypeInt32,
	0xa: types.ColumnTypeUint32,
	0xb: types.ColumnTypeInt64,
	0xc: types.ColumnTypeUint64,
	0xd: types.ColumnTypeFloat32,
	0xe: types.ColumnTypeFloat64,
	// 0xf: ColumnTypeBytes, not support yet
	0x10: types.ColumnTypeString,
	0x11: types.ColumnTypeList,
	0x12: types.ColumnTypePath,
	0x13: types.ColumnTypeRecord,
	0x15: types.ColumnTypeLocalTime,
	0x16: types.ColumnTypeDuration,
	0x17: types.ColumnTypeDate,
	0x18: types.ColumnTypeLocalDatetime,
	0x19: types.ColumnTypeZonedTime,
	0x20: types.ColumnTypeZonedDatetime,
	0xFE: types.ColumnTypeAny,
	0xFF: types.ColumnTypeInvalid,
}

type vectorType uint8

const (
	vectorTypeInvalid vectorType = iota
	vectorTypeConst
	vectorTypeFlat
	vectorTypeParallel
)

type ResultTable struct {
	table                *vector.VectorResultTable
	numBatches           int
	columnNames          []string
	columnTypes          []typeSchema
	ColumnTypeBytes      [][]byte
	batches              []batcher
	batchIndex           uint64
	currentBatch         batcher
	currentBatchRowIndex uint32
	decodeContext        *decodeContext
	decoder              decoder
}

type decodeContext struct {
	timezoneOffset int64 //offset in seconds
	graphsSchema   graphsSchema
}

type batcher interface {
	numRecords() uint32
	getRowByIndex(index uint32) ([]*nebulaValue, error)
}

type batch struct {
	vectors []*vectorWrapper
}

// vectorWrapper is a wrapper for vector.NestedVector
// one column has one vectorWrapper
type vectorWrapper struct {
	vector        *vector.NestedVector
	nullAllSet    bool // allSet is true means all values in the vector are not null
	vectorType    vectorType
	typ           typeSchema
	typBytes      []byte
	elementProps  graphElementProps
	pathMeta      *pathMetaData
	decoder       decoder
	constValue    *nebulaValue
	decodeContext *decodeContext
}

type pathMetaData struct {
	nodeTypeIndex pathElementTypeIndexMap
	edgeTypeIndex pathElementTypeIndexMap
	nodeIndexPair map[pathPairIndex]*pathPair
	edgeIndexPair map[pathPairIndex]*pathPair
}

type pathPairIndex uint16
type pathElementTypeIndexMap map[int32]pathPairIndex
type pathPair struct {
	cur *vector.NestedVector
	adj *vector.NestedVector
}
type pathAdjHeader struct {
	isEnd           bool
	nextIsEdge      bool
	nextVectorIndex uint32
	nextOffset      uint32
}
type props map[string]*vectorProps
type graphElementProps map[int32]props

type vectorProps struct {
	name        string
	typ         typeSchema
	vectorIndex int
}

func newPathAdjHeader(value int64) *pathAdjHeader {
	header := &pathAdjHeader{}
	header.isEnd = ((value >> 63) & 1) == 1
	header.nextIsEdge = ((value >> 62) & 1) == 1
	header.nextVectorIndex = (uint32)((value >> 32) & 0xFFFF)
	header.nextOffset = (uint32)(value & 0xFFFFFFFF)
	return header
}

func NewResultTable(table *vector.VectorResultTable) (*ResultTable, error) {
	if table == nil || table.Meta == nil {
		return nil, nil
	}
	if table.Batch != nil && len(table.Batch) > 0 {
		// check if the length of table.Batch is equal to the length of table.Meta.VectorBatchMetaData
		if int(table.Meta.NumBatches) != len(table.Batch) {
			return nil, errBatchNotEqual
		}
	}
	t := &ResultTable{table: table, decoder: defaultDecoder}
	// construct graph schema
	gsm := t.constructGraphsSchema()
	dctx := &decodeContext{
		timezoneOffset: int64(table.Meta.TimeZoneOffset) * int64(time.Minute/time.Second),
		graphsSchema:   gsm,
	}
	t.decodeContext = dctx

	if table.Meta.RowType != nil {
		t.columnNames = table.Meta.RowType.ColumnNames
		typs := table.Meta.RowType.ColumnTypes
		for _, typ := range typs {
			r := newBytesReader(typ.ValueType)
			columnType, err := newTypeSchema(r)
			if err != nil {
				return nil, err
			}
			t.columnTypes = append(t.columnTypes, columnType)
			t.ColumnTypeBytes = append(t.ColumnTypeBytes, typ.ValueType)
		}
	}
	t.numBatches = len(table.Batch)
	for i := 0; i < t.numBatches; i++ {
		vs := make([]*vectorWrapper, 0, len(table.Batch[i].Vectors))
		for colIndex, v := range table.Batch[i].Vectors {
			typ := t.columnTypes[colIndex]
			wrapper := &vectorWrapper{
				vector:        v,
				typ:           typ,
				decodeContext: dctx,
				decoder:       defaultDecoder,
			}
			if typ.getType() != types.ColumnTypeUnknown {
				vectorTypeCode := v.CommonMetaData.VectorContentType
				vt := uint8(vectorTypeCode & 0xFF)
				nullAllSet := (vectorTypeCode & 1 << 8) != 0
				wrapper.vectorType = vectorType(vt)
				wrapper.nullAllSet = nullAllSet
			}
			vs = append(vs, wrapper)
		}

		t.batches = append(t.batches, &batch{
			vectors: vs,
		})
	}

	return t, nil
}

func (rt *ResultTable) NumRecords() uint64 {
	return rt.table.Meta.NumRecords
}

func (rt *ResultTable) ColumnNames() []string {
	return rt.columnNames
}

func (rt *ResultTable) ColumnTypes() []types.ColumnType {
	types := make([]types.ColumnType, 0, len(rt.columnTypes))
	for _, t := range rt.columnTypes {
		types = append(types, t.getType())
	}
	return types
}

func (rt *ResultTable) constructGraphsSchema() graphsSchema {
	if rt.table.Meta.GraphSchema == nil {
		return nil
	}
	gsm := make(graphsSchema)
	for _, g := range rt.table.Meta.GraphSchema {
		gs := graphSchema{
			name:        string(g.GraphName),
			id:          g.GraphId,
			nodesSchmea: make(map[int32]*elementSchema),
			edgesSchema: make(map[int32]*elementSchema),
		}
		for _, n := range g.NodeType {
			lables := make([]string, 0)
			for _, l := range n.Label {
				lables = append(lables, string(l))
			}
			gs.nodesSchmea[n.NodeTypeId] = &elementSchema{
				typeName: string(n.NodeTypeName),
				typeId:   n.NodeTypeId,
				labels:   lables,
			}
		}
		for _, e := range g.EdgeType {
			lables := make([]string, 0)
			for _, l := range e.Label {
				lables = append(lables, string(l))
			}
			gs.edgesSchema[e.EdgeTypeId] = &elementSchema{
				typeName: string(e.EdgeTypeName),
				typeId:   e.EdgeTypeId,
				labels:   lables,
			}
		}
		gsm[g.GraphId] = &gs
	}
	return gsm
}

func (rt *ResultTable) Next() ([]types.Value, error) {
	if rt.currentBatch == nil {
		if len(rt.batches) == 0 {
			return nil, io.EOF
		} else {
			rt.currentBatch = rt.batches[0]
		}
	}
	if rt.currentBatchRowIndex >= rt.currentBatch.numRecords() {
		rt.batchIndex++
		rt.currentBatchRowIndex = 0
		if rt.batchIndex >= uint64(rt.numBatches) {
			return nil, io.EOF
		}
		rt.currentBatch = rt.batches[rt.batchIndex]
	}
	row, err := rt.currentBatch.getRowByIndex(rt.currentBatchRowIndex)
	if err != nil {
		return nil, err
	}
	rt.currentBatchRowIndex++
	values := make([]types.Value, 0, len(row))
	for _, v := range row {
		values = append(values, v)
	}
	return values, nil
}

func (b *batch) numRecords() uint32 {
	if len(b.vectors) == 0 {
		return 0
	}
	// all vectors in a batch have the same number of records
	v := b.vectors[0]
	return v.vector.CommonMetaData.GetNumRecords()
}

func (b *batch) getRowByIndex(index uint32) ([]*nebulaValue, error) {
	if index >= b.numRecords() {
		return nil, errOutOfRange
	}
	var row []*nebulaValue
	for _, v := range b.vectors {
		nvs, err := v.decodeValue(index)
		if err != nil {
			return nil, err
		}
		row = append(row, nvs)
	}
	return row, nil
}

func (v *vectorWrapper) decodeValue(index uint32) (*nebulaValue, error) {
	if v.typ == nil {
		return nil, errColumnTypeIsNil
	}
	if v.typ.getType() == types.ColumnTypeUnknown {
		return &nebulaValue{data: nil}, nil
	}
	if v.vectorType == vectorTypeConst && v.constValue != nil {
		return v.constValue, nil
	}
	var (
		value *nebulaValue
		err   error
	)
	switch v.vectorType {
	case vectorTypeFlat:
		value, err = v.decoder.decodeFlatValue(v.decodeContext, v.vector, index, v.typ)
	case vectorTypeConst:
		r := newBytesReader(v.vector.VectorData)
		value, err = v.decoder.decodeConstValue(v.decodeContext, r, v.typ.getType())
	default:
		return nil, errInvalidVectorType
	}
	if err != nil {
		return nil, err
	}
	if v.vectorType == vectorTypeConst {
		v.constValue = value
	}
	return value, nil
}
