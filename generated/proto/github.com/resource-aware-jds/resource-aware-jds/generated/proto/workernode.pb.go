// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: proto/workernode.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RecievedTask struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID             string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	TaskAttributes []byte `protobuf:"bytes,2,opt,name=TaskAttributes,proto3" json:"TaskAttributes,omitempty"`
	DockerImage    string `protobuf:"bytes,3,opt,name=DockerImage,proto3" json:"DockerImage,omitempty"`
}

func (x *RecievedTask) Reset() {
	*x = RecievedTask{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_workernode_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecievedTask) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecievedTask) ProtoMessage() {}

func (x *RecievedTask) ProtoReflect() protoreflect.Message {
	mi := &file_proto_workernode_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecievedTask.ProtoReflect.Descriptor instead.
func (*RecievedTask) Descriptor() ([]byte, []int) {
	return file_proto_workernode_proto_rawDescGZIP(), []int{0}
}

func (x *RecievedTask) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *RecievedTask) GetTaskAttributes() []byte {
	if x != nil {
		return x.TaskAttributes
	}
	return nil
}

func (x *RecievedTask) GetDockerImage() string {
	if x != nil {
		return x.DockerImage
	}
	return ""
}

type Resource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CpuCores               int64   `protobuf:"varint,1,opt,name=CpuCores,proto3" json:"CpuCores,omitempty"`
	AvailableCpuPercentage float32 `protobuf:"fixed32,2,opt,name=AvailableCpuPercentage,proto3" json:"AvailableCpuPercentage,omitempty"`
	AvailableMemory        string  `protobuf:"bytes,3,opt,name=AvailableMemory,proto3" json:"AvailableMemory,omitempty"`
}

func (x *Resource) Reset() {
	*x = Resource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_workernode_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Resource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resource) ProtoMessage() {}

func (x *Resource) ProtoReflect() protoreflect.Message {
	mi := &file_proto_workernode_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Resource.ProtoReflect.Descriptor instead.
func (*Resource) Descriptor() ([]byte, []int) {
	return file_proto_workernode_proto_rawDescGZIP(), []int{1}
}

func (x *Resource) GetCpuCores() int64 {
	if x != nil {
		return x.CpuCores
	}
	return 0
}

func (x *Resource) GetAvailableCpuPercentage() float32 {
	if x != nil {
		return x.AvailableCpuPercentage
	}
	return 0
}

func (x *Resource) GetAvailableMemory() string {
	if x != nil {
		return x.AvailableMemory
	}
	return ""
}

var File_proto_workernode_proto protoreflect.FileDescriptor

var file_proto_workernode_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x6e, 0x6f,
	0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72,
	0x4e, 0x6f, 0x64, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x68, 0x0a, 0x0c, 0x52, 0x65, 0x63, 0x69, 0x65, 0x76, 0x65, 0x64, 0x54, 0x61, 0x73,
	0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49,
	0x44, 0x12, 0x26, 0x0a, 0x0e, 0x54, 0x61, 0x73, 0x6b, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75,
	0x74, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0e, 0x54, 0x61, 0x73, 0x6b, 0x41,
	0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x6f, 0x63,
	0x6b, 0x65, 0x72, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x44, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x22, 0x88, 0x01, 0x0a, 0x08,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x70, 0x75, 0x43,
	0x6f, 0x72, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x43, 0x70, 0x75, 0x43,
	0x6f, 0x72, 0x65, 0x73, 0x12, 0x36, 0x0a, 0x16, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c,
	0x65, 0x43, 0x70, 0x75, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x61, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x02, 0x52, 0x16, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x43,
	0x70, 0x75, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x61, 0x67, 0x65, 0x12, 0x28, 0x0a, 0x0f,
	0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65,
	0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x32, 0x8b, 0x01, 0x0a, 0x0a, 0x57, 0x6f, 0x72, 0x6b, 0x65,
	0x72, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x3e, 0x0a, 0x08, 0x53, 0x65, 0x6e, 0x64, 0x54, 0x61, 0x73,
	0x6b, 0x12, 0x18, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x2e, 0x52,
	0x65, 0x63, 0x69, 0x65, 0x76, 0x65, 0x64, 0x54, 0x61, 0x73, 0x6b, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3d, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x14, 0x2e, 0x57,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x22, 0x00, 0x42, 0x42, 0x5a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2d, 0x61, 0x77, 0x61, 0x72,
	0x65, 0x2d, 0x6a, 0x64, 0x73, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2d, 0x61,
	0x77, 0x61, 0x72, 0x65, 0x2d, 0x6a, 0x64, 0x73, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74,
	0x65, 0x64, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_workernode_proto_rawDescOnce sync.Once
	file_proto_workernode_proto_rawDescData = file_proto_workernode_proto_rawDesc
)

func file_proto_workernode_proto_rawDescGZIP() []byte {
	file_proto_workernode_proto_rawDescOnce.Do(func() {
		file_proto_workernode_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_workernode_proto_rawDescData)
	})
	return file_proto_workernode_proto_rawDescData
}

var file_proto_workernode_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_workernode_proto_goTypes = []interface{}{
	(*RecievedTask)(nil),  // 0: WorkerNode.RecievedTask
	(*Resource)(nil),      // 1: WorkerNode.Resource
	(*emptypb.Empty)(nil), // 2: google.protobuf.Empty
}
var file_proto_workernode_proto_depIdxs = []int32{
	0, // 0: WorkerNode.WorkerNode.SendTask:input_type -> WorkerNode.RecievedTask
	2, // 1: WorkerNode.WorkerNode.HealthCheck:input_type -> google.protobuf.Empty
	2, // 2: WorkerNode.WorkerNode.SendTask:output_type -> google.protobuf.Empty
	1, // 3: WorkerNode.WorkerNode.HealthCheck:output_type -> WorkerNode.Resource
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_workernode_proto_init() }
func file_proto_workernode_proto_init() {
	if File_proto_workernode_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_workernode_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecievedTask); i {
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
		file_proto_workernode_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Resource); i {
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
			RawDescriptor: file_proto_workernode_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_workernode_proto_goTypes,
		DependencyIndexes: file_proto_workernode_proto_depIdxs,
		MessageInfos:      file_proto_workernode_proto_msgTypes,
	}.Build()
	File_proto_workernode_proto = out.File
	file_proto_workernode_proto_rawDesc = nil
	file_proto_workernode_proto_goTypes = nil
	file_proto_workernode_proto_depIdxs = nil
}
