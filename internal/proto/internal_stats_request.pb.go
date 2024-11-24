// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: internal_stats_request.proto

package proto

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type InternalStatsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InternalStatsRequest) Reset() {
	*x = InternalStatsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_stats_request_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InternalStatsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InternalStatsRequest) ProtoMessage() {}

func (x *InternalStatsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_stats_request_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InternalStatsRequest.ProtoReflect.Descriptor instead.
func (*InternalStatsRequest) Descriptor() ([]byte, []int) {
	return file_internal_stats_request_proto_rawDescGZIP(), []int{0}
}

var File_internal_stats_request_proto protoreflect.FileDescriptor

var file_internal_stats_request_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x73,
	0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x16,
	0x0a, 0x14, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x10, 0x5a, 0x0e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_stats_request_proto_rawDescOnce sync.Once
	file_internal_stats_request_proto_rawDescData = file_internal_stats_request_proto_rawDesc
)

func file_internal_stats_request_proto_rawDescGZIP() []byte {
	file_internal_stats_request_proto_rawDescOnce.Do(func() {
		file_internal_stats_request_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_stats_request_proto_rawDescData)
	})
	return file_internal_stats_request_proto_rawDescData
}

var file_internal_stats_request_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_internal_stats_request_proto_goTypes = []any{
	(*InternalStatsRequest)(nil), // 0: InternalStatsRequest
}
var file_internal_stats_request_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_stats_request_proto_init() }
func file_internal_stats_request_proto_init() {
	if File_internal_stats_request_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_stats_request_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*InternalStatsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_stats_request_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_stats_request_proto_goTypes,
		DependencyIndexes: file_internal_stats_request_proto_depIdxs,
		MessageInfos:      file_internal_stats_request_proto_msgTypes,
	}.Build()
	File_internal_stats_request_proto = out.File
	file_internal_stats_request_proto_rawDesc = nil
	file_internal_stats_request_proto_goTypes = nil
	file_internal_stats_request_proto_depIdxs = nil
}
