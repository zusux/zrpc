// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.2
// source: info_value.proto

package pb_value

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type InfoValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Kind        *string `protobuf:"bytes,1,opt,name=kind,proto3,oneof" json:"kind,omitempty"`                                   //grpc,http
	Ip          *string `protobuf:"bytes,2,opt,name=ip,proto3,oneof" json:"ip,omitempty"`                                       //服务注册ip
	Port        *uint32 `protobuf:"varint,3,opt,name=port,proto3,oneof" json:"port,omitempty"`                                  //服务注册端口
	Status      *uint32 `protobuf:"varint,4,opt,name=status,proto3,oneof" json:"status,omitempty"`                              //流量统计
	RequestFlow *uint32 `protobuf:"varint,5,opt,name=request_flow,json=requestFlow,proto3,oneof" json:"request_flow,omitempty"` //流量统计
	UpdatedAt   *uint32 `protobuf:"varint,6,opt,name=UpdatedAt,proto3,oneof" json:"UpdatedAt,omitempty"`                        //更新时间
}

func (x *InfoValue) Reset() {
	*x = InfoValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_info_value_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InfoValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InfoValue) ProtoMessage() {}

func (x *InfoValue) ProtoReflect() protoreflect.Message {
	mi := &file_info_value_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InfoValue.ProtoReflect.Descriptor instead.
func (*InfoValue) Descriptor() ([]byte, []int) {
	return file_info_value_proto_rawDescGZIP(), []int{0}
}

func (x *InfoValue) GetKind() string {
	if x != nil && x.Kind != nil {
		return *x.Kind
	}
	return ""
}

func (x *InfoValue) GetIp() string {
	if x != nil && x.Ip != nil {
		return *x.Ip
	}
	return ""
}

func (x *InfoValue) GetPort() uint32 {
	if x != nil && x.Port != nil {
		return *x.Port
	}
	return 0
}

func (x *InfoValue) GetStatus() uint32 {
	if x != nil && x.Status != nil {
		return *x.Status
	}
	return 0
}

func (x *InfoValue) GetRequestFlow() uint32 {
	if x != nil && x.RequestFlow != nil {
		return *x.RequestFlow
	}
	return 0
}

func (x *InfoValue) GetUpdatedAt() uint32 {
	if x != nil && x.UpdatedAt != nil {
		return *x.UpdatedAt
	}
	return 0
}

var File_info_value_proto protoreflect.FileDescriptor

var file_info_value_proto_rawDesc = []byte{
	0x0a, 0x10, 0x69, 0x6e, 0x66, 0x6f, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0xfd, 0x01, 0x0a,
	0x09, 0x49, 0x6e, 0x66, 0x6f, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x17, 0x0a, 0x04, 0x6b, 0x69,
	0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64,
	0x88, 0x01, 0x01, 0x12, 0x13, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x01, 0x52, 0x02, 0x69, 0x70, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x02, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x88, 0x01,
	0x01, 0x12, 0x1b, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0d, 0x48, 0x03, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x88, 0x01, 0x01, 0x12, 0x26,
	0x0a, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x66, 0x6c, 0x6f, 0x77, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0d, 0x48, 0x04, 0x52, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x46,
	0x6c, 0x6f, 0x77, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x05, 0x52, 0x09, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6b, 0x69,
	0x6e, 0x64, 0x42, 0x05, 0x0a, 0x03, 0x5f, 0x69, 0x70, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x70, 0x6f,
	0x72, 0x74, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x0f, 0x0a,
	0x0d, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x66, 0x6c, 0x6f, 0x77, 0x42, 0x0c,
	0x0a, 0x0a, 0x5f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x42, 0x18, 0x5a, 0x16,
	0x7a, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3b, 0x70, 0x62,
	0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_info_value_proto_rawDescOnce sync.Once
	file_info_value_proto_rawDescData = file_info_value_proto_rawDesc
)

func file_info_value_proto_rawDescGZIP() []byte {
	file_info_value_proto_rawDescOnce.Do(func() {
		file_info_value_proto_rawDescData = protoimpl.X.CompressGZIP(file_info_value_proto_rawDescData)
	})
	return file_info_value_proto_rawDescData
}

var file_info_value_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_info_value_proto_goTypes = []interface{}{
	(*InfoValue)(nil), // 0: protobuf.InfoValue
}
var file_info_value_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_info_value_proto_init() }
func file_info_value_proto_init() {
	if File_info_value_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_info_value_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InfoValue); i {
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
	file_info_value_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_info_value_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_info_value_proto_goTypes,
		DependencyIndexes: file_info_value_proto_depIdxs,
		MessageInfos:      file_info_value_proto_msgTypes,
	}.Build()
	File_info_value_proto = out.File
	file_info_value_proto_rawDesc = nil
	file_info_value_proto_goTypes = nil
	file_info_value_proto_depIdxs = nil
}
