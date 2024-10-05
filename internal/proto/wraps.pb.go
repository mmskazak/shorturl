// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: wraps.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type JWT struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jwt *wrapperspb.StringValue `protobuf:"bytes,1,opt,name=jwt,proto3" json:"jwt,omitempty"`
}

func (x *JWT) Reset() {
	*x = JWT{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wraps_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JWT) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JWT) ProtoMessage() {}

func (x *JWT) ProtoReflect() protoreflect.Message {
	mi := &file_wraps_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JWT.ProtoReflect.Descriptor instead.
func (*JWT) Descriptor() ([]byte, []int) {
	return file_wraps_proto_rawDescGZIP(), []int{0}
}

func (x *JWT) GetJwt() *wrapperspb.StringValue {
	if x != nil {
		return x.Jwt
	}
	return nil
}

type CorrelationID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CorrelationId *wrapperspb.StringValue `protobuf:"bytes,1,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
}

func (x *CorrelationID) Reset() {
	*x = CorrelationID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wraps_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CorrelationID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CorrelationID) ProtoMessage() {}

func (x *CorrelationID) ProtoReflect() protoreflect.Message {
	mi := &file_wraps_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CorrelationID.ProtoReflect.Descriptor instead.
func (*CorrelationID) Descriptor() ([]byte, []int) {
	return file_wraps_proto_rawDescGZIP(), []int{1}
}

func (x *CorrelationID) GetCorrelationId() *wrapperspb.StringValue {
	if x != nil {
		return x.CorrelationId
	}
	return nil
}

type ShortURL struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShortUrl *wrapperspb.StringValue `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
}

func (x *ShortURL) Reset() {
	*x = ShortURL{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wraps_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShortURL) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortURL) ProtoMessage() {}

func (x *ShortURL) ProtoReflect() protoreflect.Message {
	mi := &file_wraps_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortURL.ProtoReflect.Descriptor instead.
func (*ShortURL) Descriptor() ([]byte, []int) {
	return file_wraps_proto_rawDescGZIP(), []int{2}
}

func (x *ShortURL) GetShortUrl() *wrapperspb.StringValue {
	if x != nil {
		return x.ShortUrl
	}
	return nil
}

type OriginalURL struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OriginalUrl *wrapperspb.StringValue `protobuf:"bytes,1,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
}

func (x *OriginalURL) Reset() {
	*x = OriginalURL{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wraps_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OriginalURL) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OriginalURL) ProtoMessage() {}

func (x *OriginalURL) ProtoReflect() protoreflect.Message {
	mi := &file_wraps_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OriginalURL.ProtoReflect.Descriptor instead.
func (*OriginalURL) Descriptor() ([]byte, []int) {
	return file_wraps_proto_rawDescGZIP(), []int{3}
}

func (x *OriginalURL) GetOriginalUrl() *wrapperspb.StringValue {
	if x != nil {
		return x.OriginalUrl
	}
	return nil
}

var File_wraps_proto protoreflect.FileDescriptor

var file_wraps_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x77, 0x72, 0x61, 0x70, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77,
	0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x35, 0x0a,
	0x03, 0x4a, 0x57, 0x54, 0x12, 0x2e, 0x0a, 0x03, 0x6a, 0x77, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x03, 0x6a, 0x77, 0x74, 0x22, 0x54, 0x0a, 0x0d, 0x43, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x43, 0x0a, 0x0e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0d, 0x63, 0x6f, 0x72,
	0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x45, 0x0a, 0x08, 0x53, 0x68,
	0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x12, 0x39, 0x0a, 0x09, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f,
	0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x72,
	0x6c, 0x22, 0x4e, 0x0a, 0x0b, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c,
	0x12, 0x3f, 0x0a, 0x0c, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x72,
	0x6c, 0x42, 0x10, 0x5a, 0x0e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_wraps_proto_rawDescOnce sync.Once
	file_wraps_proto_rawDescData = file_wraps_proto_rawDesc
)

func file_wraps_proto_rawDescGZIP() []byte {
	file_wraps_proto_rawDescOnce.Do(func() {
		file_wraps_proto_rawDescData = protoimpl.X.CompressGZIP(file_wraps_proto_rawDescData)
	})
	return file_wraps_proto_rawDescData
}

var file_wraps_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_wraps_proto_goTypes = []any{
	(*JWT)(nil),                    // 0: JWT
	(*CorrelationID)(nil),          // 1: CorrelationID
	(*ShortURL)(nil),               // 2: ShortURL
	(*OriginalURL)(nil),            // 3: OriginalURL
	(*wrapperspb.StringValue)(nil), // 4: google.protobuf.StringValue
}
var file_wraps_proto_depIdxs = []int32{
	4, // 0: JWT.jwt:type_name -> google.protobuf.StringValue
	4, // 1: CorrelationID.correlation_id:type_name -> google.protobuf.StringValue
	4, // 2: ShortURL.short_url:type_name -> google.protobuf.StringValue
	4, // 3: OriginalURL.original_url:type_name -> google.protobuf.StringValue
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_wraps_proto_init() }
func file_wraps_proto_init() {
	if File_wraps_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_wraps_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*JWT); i {
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
		file_wraps_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CorrelationID); i {
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
		file_wraps_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*ShortURL); i {
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
		file_wraps_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*OriginalURL); i {
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
			RawDescriptor: file_wraps_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_wraps_proto_goTypes,
		DependencyIndexes: file_wraps_proto_depIdxs,
		MessageInfos:      file_wraps_proto_msgTypes,
	}.Build()
	File_wraps_proto = out.File
	file_wraps_proto_rawDesc = nil
	file_wraps_proto_goTypes = nil
	file_wraps_proto_depIdxs = nil
}
