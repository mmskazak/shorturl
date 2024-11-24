// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: handle_create_short_url_request.proto

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

type HandleCreateShortURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jwt         *JWT         `protobuf:"bytes,1,opt,name=jwt,proto3" json:"jwt,omitempty"`
	OriginalUrl *OriginalURL `protobuf:"bytes,2,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
}

func (x *HandleCreateShortURLRequest) Reset() {
	*x = HandleCreateShortURLRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_handle_create_short_url_request_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandleCreateShortURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandleCreateShortURLRequest) ProtoMessage() {}

func (x *HandleCreateShortURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_handle_create_short_url_request_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandleCreateShortURLRequest.ProtoReflect.Descriptor instead.
func (*HandleCreateShortURLRequest) Descriptor() ([]byte, []int) {
	return file_handle_create_short_url_request_proto_rawDescGZIP(), []int{0}
}

func (x *HandleCreateShortURLRequest) GetJwt() *JWT {
	if x != nil {
		return x.Jwt
	}
	return nil
}

func (x *HandleCreateShortURLRequest) GetOriginalUrl() *OriginalURL {
	if x != nil {
		return x.OriginalUrl
	}
	return nil
}

var File_handle_create_short_url_request_proto protoreflect.FileDescriptor

var file_handle_create_short_url_request_proto_rawDesc = []byte{
	0x0a, 0x25, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f,
	0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x77, 0x72, 0x61, 0x70, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x66, 0x0a, 0x1b, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x03, 0x6a, 0x77, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x04, 0x2e, 0x4a, 0x57, 0x54, 0x52, 0x03, 0x6a, 0x77, 0x74, 0x12, 0x2f, 0x0a, 0x0c, 0x6f,
	0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0c, 0x2e, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x52,
	0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x72, 0x6c, 0x42, 0x10, 0x5a, 0x0e,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_handle_create_short_url_request_proto_rawDescOnce sync.Once
	file_handle_create_short_url_request_proto_rawDescData = file_handle_create_short_url_request_proto_rawDesc
)

func file_handle_create_short_url_request_proto_rawDescGZIP() []byte {
	file_handle_create_short_url_request_proto_rawDescOnce.Do(func() {
		file_handle_create_short_url_request_proto_rawDescData = protoimpl.X.CompressGZIP(file_handle_create_short_url_request_proto_rawDescData)
	})
	return file_handle_create_short_url_request_proto_rawDescData
}

var file_handle_create_short_url_request_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_handle_create_short_url_request_proto_goTypes = []any{
	(*HandleCreateShortURLRequest)(nil), // 0: HandleCreateShortURLRequest
	(*JWT)(nil),                         // 1: JWT
	(*OriginalURL)(nil),                 // 2: OriginalURL
}
var file_handle_create_short_url_request_proto_depIdxs = []int32{
	1, // 0: HandleCreateShortURLRequest.jwt:type_name -> JWT
	2, // 1: HandleCreateShortURLRequest.original_url:type_name -> OriginalURL
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_handle_create_short_url_request_proto_init() }
func file_handle_create_short_url_request_proto_init() {
	if File_handle_create_short_url_request_proto != nil {
		return
	}
	file_wraps_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_handle_create_short_url_request_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*HandleCreateShortURLRequest); i {
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
			RawDescriptor: file_handle_create_short_url_request_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_handle_create_short_url_request_proto_goTypes,
		DependencyIndexes: file_handle_create_short_url_request_proto_depIdxs,
		MessageInfos:      file_handle_create_short_url_request_proto_msgTypes,
	}.Build()
	File_handle_create_short_url_request_proto = out.File
	file_handle_create_short_url_request_proto_rawDesc = nil
	file_handle_create_short_url_request_proto_goTypes = nil
	file_handle_create_short_url_request_proto_depIdxs = nil
}