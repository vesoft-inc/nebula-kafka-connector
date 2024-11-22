package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"

	"github.com/vesoft-inc/nebula-ng-tools/golang/internal/generated_code/v5.0.0/proto/common"
)

var PROTOCOL_VERSION []byte

func init() {
	common.File_common_proto.Options().ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if fd.FullName() == protoreflect.FullName(common.E_ProtocolVersion.Name) {
			f := common.File_common_proto.Options().ProtoReflect().Get(fd)
			i := f.Interface()
			if _, ok := i.([]byte); ok {
				PROTOCOL_VERSION = i.([]byte)
			}
		}
		return true
	})
}
