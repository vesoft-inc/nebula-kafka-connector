package errors

import (
	stderrors "errors"
)

var (
	ErrNoAddresses            = stderrors.New("no addresses")
	ErrInvalidAddress         = stderrors.New("invalid address")
	ErrNoGraphName            = stderrors.New("no graph name")
	ErrNoNodeName             = stderrors.New("no node name")
	ErrNoNodeID               = stderrors.New("no node id")
	ErrNoEdgeSrc              = stderrors.New("no edge src")
	ErrNoEdgeDst              = stderrors.New("no edge dst")
	ErrNoEdgeName             = stderrors.New("no edge name")
	ErrNoPropName             = stderrors.New("no prop name")
	ErrUnsupportedValueType   = stderrors.New("unsupported value type")
	ErrNoRecord               = stderrors.New("no record")
	ErrNoIndicesOrConcatItems = stderrors.New("no indices or concat items")
)
