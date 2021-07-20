// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: apis/hello/v1alpha1/hello.proto

package v1alpha1

import (
	reflect "reflect"
	sync "sync"

	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type IntroRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *IntroRequest) Reset() {
	*x = IntroRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_hello_v1alpha1_hello_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IntroRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IntroRequest) ProtoMessage() {}

func (x *IntroRequest) ProtoReflect() protoreflect.Message {
	mi := &file_apis_hello_v1alpha1_hello_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IntroRequest.ProtoReflect.Descriptor instead.
func (*IntroRequest) Descriptor() ([]byte, []int) {
	return file_apis_hello_v1alpha1_hello_proto_rawDescGZIP(), []int{0}
}

func (x *IntroRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type IntroResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Intro string `protobuf:"bytes,1,opt,name=intro,proto3" json:"intro,omitempty"`
}

func (x *IntroResponse) Reset() {
	*x = IntroResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_hello_v1alpha1_hello_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IntroResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IntroResponse) ProtoMessage() {}

func (x *IntroResponse) ProtoReflect() protoreflect.Message {
	mi := &file_apis_hello_v1alpha1_hello_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IntroResponse.ProtoReflect.Descriptor instead.
func (*IntroResponse) Descriptor() ([]byte, []int) {
	return file_apis_hello_v1alpha1_hello_proto_rawDescGZIP(), []int{1}
}

func (x *IntroResponse) GetIntro() string {
	if x != nil {
		return x.Intro
	}
	return ""
}

var File_apis_hello_v1alpha1_hello_proto protoreflect.FileDescriptor

var file_apis_hello_v1alpha1_hello_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2f, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x2d, 0x76, 0x6f, 0x79, 0x61, 0x67, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x68, 0x2e, 0x64,
	0x65, 0x76, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x61, 0x70, 0x69,
	0x73, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x22,
	0x0a, 0x0c, 0x49, 0x6e, 0x74, 0x72, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x22, 0x25, 0x0a, 0x0d, 0x49, 0x6e, 0x74, 0x72, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x74, 0x72, 0x6f, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x69, 0x6e, 0x74, 0x72, 0x6f, 0x32, 0xc6, 0x02, 0x0a, 0x0c, 0x48, 0x65,
	0x6c, 0x6c, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0xab, 0x01, 0x0a, 0x05, 0x49,
	0x6e, 0x74, 0x72, 0x6f, 0x12, 0x3b, 0x2e, 0x76, 0x6f, 0x79, 0x61, 0x67, 0x65, 0x72, 0x6d, 0x65,
	0x73, 0x68, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x2e, 0x49, 0x6e, 0x74, 0x72, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x3c, 0x2e, 0x76, 0x6f, 0x79, 0x61, 0x67, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x68, 0x2e,
	0x64, 0x65, 0x76, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x61, 0x70,
	0x69, 0x73, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0x2e, 0x49, 0x6e, 0x74, 0x72, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x27, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x21, 0x12, 0x1f, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x68,
	0x65, 0x6c, 0x6c, 0x6f, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x69, 0x6e,
	0x74, 0x72, 0x6f, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x12, 0x87, 0x01, 0x0a, 0x06, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x12, 0x3b, 0x2e, 0x76, 0x6f, 0x79, 0x61, 0x67, 0x65, 0x72, 0x6d, 0x65, 0x73,
	0x68, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x61, 0x70, 0x69, 0x73, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x2e, 0x49, 0x6e, 0x74, 0x72, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x3c, 0x2e, 0x76, 0x6f, 0x79, 0x61, 0x67, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x68, 0x2e, 0x64,
	0x65, 0x76, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x61, 0x70, 0x69,
	0x73, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x2e, 0x49, 0x6e, 0x74, 0x72, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x30, 0x01, 0x42, 0x34, 0x5a, 0x32, 0x76, 0x6f, 0x79, 0x61, 0x67, 0x65, 0x72, 0x6d, 0x65, 0x73,
	0x68, 0x2e, 0x64, 0x65, 0x76, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2d, 0x67, 0x72, 0x70, 0x63,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2f,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_apis_hello_v1alpha1_hello_proto_rawDescOnce sync.Once
	file_apis_hello_v1alpha1_hello_proto_rawDescData = file_apis_hello_v1alpha1_hello_proto_rawDesc
)

func file_apis_hello_v1alpha1_hello_proto_rawDescGZIP() []byte {
	file_apis_hello_v1alpha1_hello_proto_rawDescOnce.Do(func() {
		file_apis_hello_v1alpha1_hello_proto_rawDescData = protoimpl.X.CompressGZIP(file_apis_hello_v1alpha1_hello_proto_rawDescData)
	})
	return file_apis_hello_v1alpha1_hello_proto_rawDescData
}

var file_apis_hello_v1alpha1_hello_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_apis_hello_v1alpha1_hello_proto_goTypes = []interface{}{
	(*IntroRequest)(nil),  // 0: voyagermesh.dev.hellogrpc.apis.hello.v1alpha1.IntroRequest
	(*IntroResponse)(nil), // 1: voyagermesh.dev.hellogrpc.apis.hello.v1alpha1.IntroResponse
}
var file_apis_hello_v1alpha1_hello_proto_depIdxs = []int32{
	0, // 0: voyagermesh.dev.hellogrpc.apis.hello.v1alpha1.HelloService.Intro:input_type -> voyagermesh.dev.hellogrpc.apis.hello.v1alpha1.IntroRequest
	0, // 1: voyagermesh.dev.hellogrpc.apis.hello.v1alpha1.HelloService.Stream:input_type -> voyagermesh.dev.hellogrpc.apis.hello.v1alpha1.IntroRequest
	1, // 2: voyagermesh.dev.hellogrpc.apis.hello.v1alpha1.HelloService.Intro:output_type -> voyagermesh.dev.hellogrpc.apis.hello.v1alpha1.IntroResponse
	1, // 3: voyagermesh.dev.hellogrpc.apis.hello.v1alpha1.HelloService.Stream:output_type -> voyagermesh.dev.hellogrpc.apis.hello.v1alpha1.IntroResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_apis_hello_v1alpha1_hello_proto_init() }
func file_apis_hello_v1alpha1_hello_proto_init() {
	if File_apis_hello_v1alpha1_hello_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_apis_hello_v1alpha1_hello_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IntroRequest); i {
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
		file_apis_hello_v1alpha1_hello_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IntroResponse); i {
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
			RawDescriptor: file_apis_hello_v1alpha1_hello_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_apis_hello_v1alpha1_hello_proto_goTypes,
		DependencyIndexes: file_apis_hello_v1alpha1_hello_proto_depIdxs,
		MessageInfos:      file_apis_hello_v1alpha1_hello_proto_msgTypes,
	}.Build()
	File_apis_hello_v1alpha1_hello_proto = out.File
	file_apis_hello_v1alpha1_hello_proto_rawDesc = nil
	file_apis_hello_v1alpha1_hello_proto_goTypes = nil
	file_apis_hello_v1alpha1_hello_proto_depIdxs = nil
}
