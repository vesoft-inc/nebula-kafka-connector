package decode

import (
	"encoding/binary"
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/errors"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/types"
)

var order = binary.LittleEndian

type bytesReader struct {
	bs    []byte
	index int
	err   error
}

func newBytesReader(bs []byte) *bytesReader {
	return &bytesReader{
		bs: bs,
	}
}

func (r *bytesReader) readN(n int) []byte {
	if r.index+n > len(r.bs) {
		r.err = errors.Wrap(errOutOfRange, "")
		return nil
	}
	bs := r.bs[r.index : r.index+n]
	r.index += n
	return bs
}

func (r *bytesReader) readPeddingAll() []byte {
	return r.bs[r.index:]
}

func (r *bytesReader) readUtilZero() []byte {
	for i := r.index; i < len(r.bs); i++ {
		if r.bs[i] == 0 {
			bs := r.bs[r.index:i]
			r.index = i + 1
			return bs
		}
	}
	r.err = errNoZeroString
	return nil
}

func (r *bytesReader) error() error {
	return r.err
}

func bytesToInt8(bs []byte) int8 {
	return int8(bs[0])
}

func bytesToInt16(bs []byte) int16 {
	return int16(order.Uint16(bs))
}

func bytesToInt32(bs []byte) int32 {
	return int32(order.Uint32(bs))
}

func bytesToInt64(bs []byte) int64 {
	return int64(order.Uint64(bs))
}

func bytesToUint8(bs []byte) uint8 {
	return uint8(bs[0])
}

func bytesToUint16(bs []byte) uint16 {
	return order.Uint16(bs)
}

func bytesToUint32(bs []byte) uint32 {
	return order.Uint32(bs)
}

func bytesToUint64(bs []byte) uint64 {
	return order.Uint64(bs)
}

func isBasicColumnType(typ types.ColumnType) bool {
	switch typ {
	case types.ColumnTypeBool:
		fallthrough
	case types.ColumnTypeInt8, types.ColumnTypeInt16, types.ColumnTypeInt32, types.ColumnTypeInt64:
		fallthrough
	case types.ColumnTypeUint8, types.ColumnTypeUint16, types.ColumnTypeUint32, types.ColumnTypeUint64:
		fallthrough
	case types.ColumnTypeFloat32, types.ColumnTypeFloat64:
		fallthrough
	case types.ColumnTypeLocalTime, types.ColumnTypeLocalDatetime:
		fallthrough
	case types.ColumnTypeZonedTime, types.ColumnTypeZonedDatetime:
		fallthrough
	case types.ColumnTypeDate, types.ColumnTypeDuration:
		fallthrough
	case types.ColumnTypeDecimal:
		return true
	default:
		return false
	}
}

func getSchemaName(gsm graphsSchema, graphId int32, elementTypeId int32, isNode bool) (graphName, typeName string, labels []string, err error) {
	gs, ok := gsm[graphId]
	if !ok {
		return "", "", nil, errors.Wrap(errGraphNotFound, "")
	}
	var elementsSchema map[int32]*elementSchema
	if isNode {
		elementsSchema = gs.nodesSchmea
	} else {
		elementsSchema = gs.edgesSchema
	}
	es, ok := elementsSchema[elementTypeId]
	if !ok {
		return "", "", nil, errors.Wrap(errElementTypeNotFount, fmt.Sprintf("element id: %d", elementTypeId))
	}
	return gs.name, es.typeName, es.labels, nil
}
