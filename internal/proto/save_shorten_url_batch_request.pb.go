// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: save_shorten_url_batch_request.proto

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

type SaveShortenURLsBatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jwt      *JWT        `protobuf:"bytes,1,opt,name=jwt,proto3" json:"jwt,omitempty"`
	Incoming []*Incoming `protobuf:"bytes,2,rep,name=incoming,proto3" json:"incoming,omitempty"`
}

func (x *SaveShortenURLsBatchRequest) Reset() {
	*x = SaveShortenURLsBatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_save_shorten_url_batch_request_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveShortenURLsBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveShortenURLsBatchRequest) ProtoMessage() {}

func (x *SaveShortenURLsBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_save_shorten_url_batch_request_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveShortenURLsBatchRequest.ProtoReflect.Descriptor instead.
func (*SaveShortenURLsBatchRequest) Descriptor() ([]byte, []int) {
	return file_save_shorten_url_batch_request_proto_rawDescGZIP(), []int{0}
}

func (x *SaveShortenURLsBatchRequest) GetJwt() *JWT {
	if x != nil {
		return x.Jwt
	}
	return nil
}

func (x *SaveShortenURLsBatchRequest) GetIncoming() []*Incoming {
	if x != nil {
		return x.Incoming
	}
	return nil
}

type Incoming struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CorrelationId *CorrelationID `protobuf:"bytes,1,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	OriginalUrl   *OriginalURL   `protobuf:"bytes,2,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
}

func (x *Incoming) Reset() {
	*x = Incoming{}
	if protoimpl.UnsafeEnabled {
		mi := &file_save_shorten_url_batch_request_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Incoming) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Incoming) ProtoMessage() {}

func (x *Incoming) ProtoReflect() protoreflect.Message {
	mi := &file_save_shorten_url_batch_request_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Incoming.ProtoReflect.Descriptor instead.
func (*Incoming) Descriptor() ([]byte, []int) {
	return file_save_shorten_url_batch_request_proto_rawDescGZIP(), []int{1}
}

func (x *Incoming) GetCorrelationId() *CorrelationID {
	if x != nil {
		return x.CorrelationId
	}
	return nil
}

func (x *Incoming) GetOriginalUrl() *OriginalURL {
	if x != nil {
		return x.OriginalUrl
	}
	return nil
}

var File_save_shorten_url_batch_request_proto protoreflect.FileDescriptor

var file_save_shorten_url_batch_request_proto_rawDesc = []byte{
	0x0a, 0x24, 0x73, 0x61, 0x76, 0x65, 0x5f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x5f, 0x75,
	0x72, 0x6c, 0x5f, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x77, 0x72, 0x61, 0x70, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x5c, 0x0a, 0x1b, 0x53, 0x61, 0x76, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74,
	0x65, 0x6e, 0x55, 0x52, 0x4c, 0x73, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x16, 0x0a, 0x03, 0x6a, 0x77, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x04, 0x2e, 0x4a, 0x57, 0x54, 0x52, 0x03, 0x6a, 0x77, 0x74, 0x12, 0x25, 0x0a, 0x08, 0x69, 0x6e,
	0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x49,
	0x6e, 0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x52, 0x08, 0x69, 0x6e, 0x63, 0x6f, 0x6d, 0x69, 0x6e,
	0x67, 0x22, 0x72, 0x0a, 0x08, 0x49, 0x6e, 0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x12, 0x35, 0x0a,
	0x0e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x43, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x44, 0x52, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x12, 0x2f, 0x0a, 0x0c, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c,
	0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x4f, 0x72, 0x69,
	0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x52, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e,
	0x61, 0x6c, 0x55, 0x72, 0x6c, 0x42, 0x10, 0x5a, 0x0e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_save_shorten_url_batch_request_proto_rawDescOnce sync.Once
	file_save_shorten_url_batch_request_proto_rawDescData = file_save_shorten_url_batch_request_proto_rawDesc
)

func file_save_shorten_url_batch_request_proto_rawDescGZIP() []byte {
	file_save_shorten_url_batch_request_proto_rawDescOnce.Do(func() {
		file_save_shorten_url_batch_request_proto_rawDescData = protoimpl.X.CompressGZIP(file_save_shorten_url_batch_request_proto_rawDescData)
	})
	return file_save_shorten_url_batch_request_proto_rawDescData
}

var file_save_shorten_url_batch_request_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_save_shorten_url_batch_request_proto_goTypes = []any{
	(*SaveShortenURLsBatchRequest)(nil), // 0: SaveShortenURLsBatchRequest
	(*Incoming)(nil),                    // 1: Incoming
	(*JWT)(nil),                         // 2: JWT
	(*CorrelationID)(nil),               // 3: CorrelationID
	(*OriginalURL)(nil),                 // 4: OriginalURL
}
var file_save_shorten_url_batch_request_proto_depIdxs = []int32{
	2, // 0: SaveShortenURLsBatchRequest.jwt:type_name -> JWT
	1, // 1: SaveShortenURLsBatchRequest.incoming:type_name -> Incoming
	3, // 2: Incoming.correlation_id:type_name -> CorrelationID
	4, // 3: Incoming.original_url:type_name -> OriginalURL
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_save_shorten_url_batch_request_proto_init() }
func file_save_shorten_url_batch_request_proto_init() {
	if File_save_shorten_url_batch_request_proto != nil {
		return
	}
	file_wraps_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_save_shorten_url_batch_request_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*SaveShortenURLsBatchRequest); i {
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
		file_save_shorten_url_batch_request_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Incoming); i {
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
			RawDescriptor: file_save_shorten_url_batch_request_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_save_shorten_url_batch_request_proto_goTypes,
		DependencyIndexes: file_save_shorten_url_batch_request_proto_depIdxs,
		MessageInfos:      file_save_shorten_url_batch_request_proto_msgTypes,
	}.Build()
	File_save_shorten_url_batch_request_proto = out.File
	file_save_shorten_url_batch_request_proto_rawDesc = nil
	file_save_shorten_url_batch_request_proto_goTypes = nil
	file_save_shorten_url_batch_request_proto_depIdxs = nil
}
