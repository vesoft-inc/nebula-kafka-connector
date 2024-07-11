package nebula_ng

import "fmt"

type scanner interface {
	scan(value Value) error
}

type nullable[T any] struct {
	Data  T
	Valid bool
}

type (
	NullString struct {
		nullable[String]
	}
	NullBool struct {
		nullable[Bool]
	}
	NullInt struct {
		nullable[Int64]
	}
	NullInt8 struct {
		nullable[Int8]
	}
	NullInt16 struct {
		nullable[Int16]
	}
	NullInt32 struct {
		nullable[Int32]
	}
	NullInt64 struct {
		nullable[Int64]
	}
	NullUInt struct {
		nullable[UInt64]
	}
	NullUInt8 struct {
		nullable[UInt8]
	}
	NullUInt16 struct {
		nullable[UInt16]
	}
	NullUInt32 struct {
		nullable[UInt32]
	}
	NullUInt64 struct {
		nullable[UInt64]
	}
	NullFloat struct {
		nullable[Float]
	}
	NullDouble struct {
		nullable[Double]
	}
	NullList struct {
		nullable[List]
	}
	NullRecord struct {
		nullable[Record]
	}
	NullDuration struct {
		nullable[Duration]
	}
	NullLocalTime struct {
		nullable[LocalTime]
	}
	NullLocalDatetime struct {
		nullable[LocalDatetime]
	}
	NullDate struct {
		nullable[Date]
	}
	NullZonedDatetime struct {
		nullable[ZonedDatetime]
	}
	NullZonedTime struct {
		nullable[ZonedTime]
	}
	NullNode struct {
		nullable[Node]
	}
	NullEdge struct {
		nullable[Edge]
	}
	NullPath struct {
		nullable[Path]
	}
)

func (n *nullable[T]) scanInternal(value Value, fn func(Value) (T, error)) error {
	if value == nil || value.IsNull() {
		n.Valid = false
		return nil
	}
	d, err := fn(value)
	if err != nil {
		return err
	}
	n.Valid = true
	n.Data = d
	return nil
}

// used for testing
func (n *nullable[T]) getData() any {
	return n.Data
}

// used for testing
func (n *nullable[T]) isValid() bool {
	return n.Valid
}

func (n *NullString) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (String, error) {
		return value.AsString()
	})
}

func (n *NullBool) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (Bool, error) {
		return value.AsBool()
	})
}

func (n *NullInt) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (Int64, error) {
		var (
			i   Int64
			err error
		)
		switch value.GetType() {
		case ValueTypeInt8:
			var d Int8
			d, err = value.AsInt8()
			i = Int64(d)
		case ValueTypeInt16:
			var d Int16
			d, err = value.AsInt16()
			i = Int64(d)
		case ValueTypeInt32:
			var d Int32
			d, err = value.AsInt32()
			i = Int64(d)
		case ValueTypeInt64:
			i, err = value.AsInt64()
		default:
			return 0, errType(fmt.Sprintf("value type not match"))
		}
		if err != nil {
			return 0, err
		}
		return i, nil
	})
}

func (n *NullInt8) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (Int8, error) {
		return value.AsInt8()
	})
}

func (n *NullInt16) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (Int16, error) {
		return value.AsInt16()
	})
}

func (n *NullInt32) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (Int32, error) {
		return value.AsInt32()
	})
}

func (n *NullInt64) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (Int64, error) {
		return value.AsInt64()
	})
}

func (n *NullUInt) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (UInt64, error) {
		var (
			i   UInt64
			err error
		)
		switch value.GetType() {
		case ValueTypeInt8:
			var d UInt8
			d, err = value.AsUInt8()
			i = UInt64(d)
		case ValueTypeInt16:
			var d UInt16
			d, err = value.AsUInt16()
			i = UInt64(d)
		case ValueTypeInt32:
			var d UInt32
			d, err = value.AsUInt32()
			i = UInt64(d)
		case ValueTypeInt64:
			i, err = value.AsUInt64()
		default:
			return 0, errType(fmt.Sprintf("value type not match"))
		}
		if err != nil {
			return 0, err
		}
		return i, nil
	})
}

func (n *NullUInt8) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (UInt8, error) {
		return value.AsUInt8()
	})
}

func (n *NullUInt16) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (UInt16, error) {
		return value.AsUInt16()
	})
}

func (n *NullUInt32) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (UInt32, error) {
		return value.AsUInt32()
	})
}

func (n *NullUInt64) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (UInt64, error) {
		return value.AsUInt64()
	})
}

func (n *NullFloat) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (Float, error) {
		return value.AsFloat()
	})
}
func (n *NullDouble) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (Double, error) {
		var (
			d   Double
			err error
		)
		switch value.GetType() {
		case ValueTypeFloat:
			var f Float
			f, err = value.AsFloat()
			d = Double(f)
		case ValueTypeDouble:
			d, err = value.AsDouble()
		default:
			return 0, errType(fmt.Sprintf("value type not match"))
		}
		if err != nil {
			return 0, err
		}
		return d, nil
	})
}

func (n *NullList) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (List, error) {
		return value.AsList()
	})
}

func (n *NullRecord) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (Record, error) {
		return value.AsRecord()
	})
}

func (n *NullDuration) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (Duration, error) {
		return value.AsDuration()
	})
}

func (n *NullLocalTime) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (LocalTime, error) {
		return value.AsLocalTime()
	})
}

func (n *NullLocalDatetime) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (LocalDatetime, error) {
		return value.AsLocalDatetime()
	})
}

func (n *NullDate) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (Date, error) {
		return value.AsDate()
	})
}

func (n *NullZonedDatetime) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (ZonedDatetime, error) {
		return value.AsZonedDatetime()
	})
}

func (n *NullZonedTime) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (ZonedTime, error) {
		return value.AsZonedTime()
	})
}

func (n *NullNode) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (Node, error) {
		return value.AsNode()
	})
}

func (n *NullEdge) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (Edge, error) {
		return value.AsEdge()
	})
}

func (n *NullPath) scan(value Value) error {
	return n.scanInternal(value, func(value Value) (Path, error) {
		return value.AsPath()
	})
}
