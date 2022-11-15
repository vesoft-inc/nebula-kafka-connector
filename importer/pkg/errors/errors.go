package errors

import (
	stderrors "errors"
)

var (
	ErrNoAddresses          = stderrors.New("no addresses")
	ErrInvalidAddress       = stderrors.New("invalid address")
	ErrNoGraphName          = stderrors.New("no graph name")
	ErrNodeNotFound         = stderrors.New("node not found")
	ErrNoNodeName           = stderrors.New("no node name")
	ErrNoNodeID             = stderrors.New("no node id")
	ErrEdgeNotFound         = stderrors.New("edge not found")
	ErrNoEdgeSrc            = stderrors.New("no edge src")
	ErrNoEdgeDst            = stderrors.New("no edge dst")
	ErrNoEdgeName           = stderrors.New("no edge name")
	ErrNoPropName           = stderrors.New("no prop name")
	ErrUnsupportedValueType = stderrors.New("unsupported value type")
	ErrNoRecord             = stderrors.New("no record")
)
