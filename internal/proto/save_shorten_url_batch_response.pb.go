// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: save_shorten_url_batch_response.proto

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

type SaveShortenURLsBatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Output []*Output               `protobuf:"bytes,1,rep,name=output,proto3" json:"output,omitempty"`
	Jwt    *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=jwt,proto3" json:"jwt,omitempty"`
}

func (x *SaveShortenURLsBatchResponse) Reset() {
	*x = SaveShortenURLsBatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_save_shorten_url_batch_response_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveShortenURLsBatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveShortenURLsBatchResponse) ProtoMessage() {}

func (x *SaveShortenURLsBatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_save_shorten_url_batch_response_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveShortenURLsBatchResponse.ProtoReflect.Descriptor instead.
func (*SaveShortenURLsBatchResponse) Descriptor() ([]byte, []int) {
	return file_save_shorten_url_batch_response_proto_rawDescGZIP(), []int{0}
}

func (x *SaveShortenURLsBatchResponse) GetOutput() []*Output {
	if x != nil {
		return x.Output
	}
	return nil
}

func (x *SaveShortenURLsBatchResponse) GetJwt() *wrapperspb.StringValue {
	if x != nil {
		return x.Jwt
	}
	return nil
}

type Output struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CorrelationId *wrapperspb.StringValue `protobuf:"bytes,1,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	ShortUrl      *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
}

func (x *Output) Reset() {
	*x = Output{}
	if protoimpl.UnsafeEnabled {
		mi := &file_save_shorten_url_batch_response_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Output) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Output) ProtoMessage() {}

func (x *Output) ProtoReflect() protoreflect.Message {
	mi := &file_save_shorten_url_batch_response_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Output.ProtoReflect.Descriptor instead.
func (*Output) Descriptor() ([]byte, []int) {
	return file_save_shorten_url_batch_response_proto_rawDescGZIP(), []int{1}
}

func (x *Output) GetCorrelationId() *wrapperspb.StringValue {
	if x != nil {
		return x.CorrelationId
	}
	return nil
}

func (x *Output) GetShortUrl() *wrapperspb.StringValue {
	if x != nil {
		return x.ShortUrl
	}
	return nil
}

var File_save_shorten_url_batch_response_proto protoreflect.FileDescriptor

var file_save_shorten_url_batch_response_proto_rawDesc = []byte{
	0x0a, 0x25, 0x73, 0x61, 0x76, 0x65, 0x5f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x5f, 0x75,
	0x72, 0x6c, 0x5f, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6f, 0x0a, 0x1c, 0x53, 0x61, 0x76, 0x65, 0x53,
	0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55, 0x52, 0x4c, 0x73, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75,
	0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x52, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x2e, 0x0a, 0x03, 0x6a, 0x77, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x52, 0x03, 0x6a, 0x77, 0x74, 0x22, 0x88, 0x01, 0x0a, 0x06, 0x4f, 0x75, 0x74,
	0x70, 0x75, 0x74, 0x12, 0x43, 0x0a, 0x0e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x39, 0x0a, 0x09, 0x73, 0x68, 0x6f, 0x72,
	0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74,
	0x55, 0x72, 0x6c, 0x42, 0x10, 0x5a, 0x0e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_save_shorten_url_batch_response_proto_rawDescOnce sync.Once
	file_save_shorten_url_batch_response_proto_rawDescData = file_save_shorten_url_batch_response_proto_rawDesc
)

func file_save_shorten_url_batch_response_proto_rawDescGZIP() []byte {
	file_save_shorten_url_batch_response_proto_rawDescOnce.Do(func() {
		file_save_shorten_url_batch_response_proto_rawDescData = protoimpl.X.CompressGZIP(file_save_shorten_url_batch_response_proto_rawDescData)
	})
	return file_save_shorten_url_batch_response_proto_rawDescData
}

var file_save_shorten_url_batch_response_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_save_shorten_url_batch_response_proto_goTypes = []any{
	(*SaveShortenURLsBatchResponse)(nil), // 0: SaveShortenURLsBatchResponse
	(*Output)(nil),                       // 1: Output
	(*wrapperspb.StringValue)(nil),       // 2: google.protobuf.StringValue
}
var file_save_shorten_url_batch_response_proto_depIdxs = []int32{
	1, // 0: SaveShortenURLsBatchResponse.output:type_name -> Output
	2, // 1: SaveShortenURLsBatchResponse.jwt:type_name -> google.protobuf.StringValue
	2, // 2: Output.correlation_id:type_name -> google.protobuf.StringValue
	2, // 3: Output.short_url:type_name -> google.protobuf.StringValue
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_save_shorten_url_batch_response_proto_init() }
func file_save_shorten_url_batch_response_proto_init() {
	if File_save_shorten_url_batch_response_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_save_shorten_url_batch_response_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*SaveShortenURLsBatchResponse); i {
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
		file_save_shorten_url_batch_response_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Output); i {
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
			RawDescriptor: file_save_shorten_url_batch_response_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_save_shorten_url_batch_response_proto_goTypes,
		DependencyIndexes: file_save_shorten_url_batch_response_proto_depIdxs,
		MessageInfos:      file_save_shorten_url_batch_response_proto_msgTypes,
	}.Build()
	File_save_shorten_url_batch_response_proto = out.File
	file_save_shorten_url_batch_response_proto_rawDesc = nil
	file_save_shorten_url_batch_response_proto_goTypes = nil
	file_save_shorten_url_batch_response_proto_depIdxs = nil
}
