package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/version"
)

var PROTOCOL_VERSION []byte

func init() {
	version.File_version_proto.Options().ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if fd.FullName() == protoreflect.FullName(version.E_ProtocolVersion.Name) {
			f := version.File_version_proto.Options().ProtoReflect().Get(fd)
			i := f.Interface()
			if _, ok := i.([]byte); ok {
				PROTOCOL_VERSION = i.([]byte)
			}
		}
		return true
	})
}
