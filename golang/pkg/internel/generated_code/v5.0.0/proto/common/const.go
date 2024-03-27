package common

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

var PROTOCOL_VERSION []byte

func init() {
	File_common_proto.Options().ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if fd.FullName() == protoreflect.FullName(E_ProtocolVersion.Name) {
			f := File_common_proto.Options().ProtoReflect().Get(fd)
			i := f.Interface()
			if _, ok := i.([]byte); ok {
				PROTOCOL_VERSION = i.([]byte)
			}
		}
		return true
	})
}
