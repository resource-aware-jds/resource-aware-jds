// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.1
// source: proto/controlplane.proto

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

type WorkerCheckInRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Port        int32  `protobuf:"varint,1,opt,name=port,proto3" json:"port,omitempty"`
	Certificate []byte `protobuf:"bytes,2,opt,name=certificate,proto3" json:"certificate,omitempty"`
}

func (x *WorkerCheckInRequest) Reset() {
	*x = WorkerCheckInRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_controlplane_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkerCheckInRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkerCheckInRequest) ProtoMessage() {}

func (x *WorkerCheckInRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_controlplane_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkerCheckInRequest.ProtoReflect.Descriptor instead.
func (*WorkerCheckInRequest) Descriptor() ([]byte, []int) {
	return file_proto_controlplane_proto_rawDescGZIP(), []int{0}
}

func (x *WorkerCheckInRequest) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *WorkerCheckInRequest) GetCertificate() []byte {
	if x != nil {
		return x.Certificate
	}
	return nil
}

type ComputeNodeRegistrationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip            string `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	Port          int32  `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	NodePublicKey []byte `protobuf:"bytes,3,opt,name=nodePublicKey,proto3" json:"nodePublicKey,omitempty"`
}

func (x *ComputeNodeRegistrationRequest) Reset() {
	*x = ComputeNodeRegistrationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_controlplane_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ComputeNodeRegistrationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ComputeNodeRegistrationRequest) ProtoMessage() {}

func (x *ComputeNodeRegistrationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_controlplane_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ComputeNodeRegistrationRequest.ProtoReflect.Descriptor instead.
func (*ComputeNodeRegistrationRequest) Descriptor() ([]byte, []int) {
	return file_proto_controlplane_proto_rawDescGZIP(), []int{1}
}

func (x *ComputeNodeRegistrationRequest) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *ComputeNodeRegistrationRequest) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *ComputeNodeRegistrationRequest) GetNodePublicKey() []byte {
	if x != nil {
		return x.NodePublicKey
	}
	return nil
}

type ComputeNodeRegistrationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Certificate []byte `protobuf:"bytes,2,opt,name=certificate,proto3" json:"certificate,omitempty"`
}

func (x *ComputeNodeRegistrationResponse) Reset() {
	*x = ComputeNodeRegistrationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_controlplane_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ComputeNodeRegistrationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ComputeNodeRegistrationResponse) ProtoMessage() {}

func (x *ComputeNodeRegistrationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_controlplane_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ComputeNodeRegistrationResponse.ProtoReflect.Descriptor instead.
func (*ComputeNodeRegistrationResponse) Descriptor() ([]byte, []int) {
	return file_proto_controlplane_proto_rawDescGZIP(), []int{2}
}

func (x *ComputeNodeRegistrationResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ComputeNodeRegistrationResponse) GetCertificate() []byte {
	if x != nil {
		return x.Certificate
	}
	return nil
}

type ControlPlaneTask struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID             string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Status         string `protobuf:"bytes,2,opt,name=Status,proto3" json:"Status,omitempty"`
	TaskAttributes []byte `protobuf:"bytes,3,opt,name=TaskAttributes,proto3" json:"TaskAttributes,omitempty"`
}

func (x *ControlPlaneTask) Reset() {
	*x = ControlPlaneTask{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_controlplane_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ControlPlaneTask) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ControlPlaneTask) ProtoMessage() {}

func (x *ControlPlaneTask) ProtoReflect() protoreflect.Message {
	mi := &file_proto_controlplane_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ControlPlaneTask.ProtoReflect.Descriptor instead.
func (*ControlPlaneTask) Descriptor() ([]byte, []int) {
	return file_proto_controlplane_proto_rawDescGZIP(), []int{3}
}

func (x *ControlPlaneTask) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *ControlPlaneTask) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *ControlPlaneTask) GetTaskAttributes() []byte {
	if x != nil {
		return x.TaskAttributes
	}
	return nil
}

type CreateJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ImageURL       string   `protobuf:"bytes,1,opt,name=ImageURL,proto3" json:"ImageURL,omitempty"`
	TaskAttributes [][]byte `protobuf:"bytes,2,rep,name=TaskAttributes,proto3" json:"TaskAttributes,omitempty"`
}

func (x *CreateJobRequest) Reset() {
	*x = CreateJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_controlplane_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateJobRequest) ProtoMessage() {}

func (x *CreateJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_controlplane_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateJobRequest.ProtoReflect.Descriptor instead.
func (*CreateJobRequest) Descriptor() ([]byte, []int) {
	return file_proto_controlplane_proto_rawDescGZIP(), []int{4}
}

func (x *CreateJobRequest) GetImageURL() string {
	if x != nil {
		return x.ImageURL
	}
	return ""
}

func (x *CreateJobRequest) GetTaskAttributes() [][]byte {
	if x != nil {
		return x.TaskAttributes
	}
	return nil
}

type CreateJobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID       string              `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Tasks    []*ControlPlaneTask `protobuf:"bytes,2,rep,name=Tasks,proto3" json:"Tasks,omitempty"`
	Status   string              `protobuf:"bytes,3,opt,name=Status,proto3" json:"Status,omitempty"`
	ImageURL string              `protobuf:"bytes,4,opt,name=ImageURL,proto3" json:"ImageURL,omitempty"`
}

func (x *CreateJobResponse) Reset() {
	*x = CreateJobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_controlplane_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateJobResponse) ProtoMessage() {}

func (x *CreateJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_controlplane_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateJobResponse.ProtoReflect.Descriptor instead.
func (*CreateJobResponse) Descriptor() ([]byte, []int) {
	return file_proto_controlplane_proto_rawDescGZIP(), []int{5}
}

func (x *CreateJobResponse) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *CreateJobResponse) GetTasks() []*ControlPlaneTask {
	if x != nil {
		return x.Tasks
	}
	return nil
}

func (x *CreateJobResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *CreateJobResponse) GetImageURL() string {
	if x != nil {
		return x.ImageURL
	}
	return ""
}

var File_proto_controlplane_proto protoreflect.FileDescriptor

var file_proto_controlplane_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70,
	0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4c, 0x0a, 0x14, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x49, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72,
	0x74, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x65, 0x22, 0x6a, 0x0a, 0x1e, 0x43, 0x6f, 0x6d, 0x70, 0x75, 0x74, 0x65, 0x4e, 0x6f,
	0x64, 0x65, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x6e, 0x6f, 0x64,
	0x65, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x0d, 0x6e, 0x6f, 0x64, 0x65, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x22,
	0x53, 0x0a, 0x1f, 0x43, 0x6f, 0x6d, 0x70, 0x75, 0x74, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x22, 0x62, 0x0a, 0x10, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50,
	0x6c, 0x61, 0x6e, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x26, 0x0a, 0x0e, 0x54, 0x61, 0x73, 0x6b, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0e, 0x54, 0x61, 0x73, 0x6b, 0x41, 0x74,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x22, 0x56, 0x0a, 0x10, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x49, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x52, 0x4c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x49, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x52, 0x4c, 0x12, 0x26, 0x0a, 0x0e, 0x54, 0x61, 0x73, 0x6b,
	0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c,
	0x52, 0x0e, 0x54, 0x61, 0x73, 0x6b, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73,
	0x22, 0x8d, 0x01, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x34, 0x0a, 0x05, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70,
	0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6c, 0x61, 0x6e,
	0x65, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x05, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x12, 0x16, 0x0a, 0x06,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x52, 0x4c,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x52, 0x4c,
	0x32, 0x9c, 0x02, 0x0a, 0x0c, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6c, 0x61, 0x6e,
	0x65, 0x12, 0x71, 0x0a, 0x12, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2c, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x75, 0x74, 0x65, 0x4e, 0x6f,
	0x64, 0x65, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70,
	0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x75, 0x74, 0x65, 0x4e, 0x6f, 0x64, 0x65,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4c, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4a, 0x6f,
	0x62, 0x12, 0x1e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x4b, 0x0a, 0x0d, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x49, 0x6e, 0x12, 0x22, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61,
	0x6e, 0x65, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x49, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42,
	0x42, 0x5a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2d, 0x61, 0x77, 0x61, 0x72, 0x65, 0x2d, 0x6a, 0x64, 0x73,
	0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2d, 0x61, 0x77, 0x61, 0x72, 0x65, 0x2d,
	0x6a, 0x64, 0x73, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_controlplane_proto_rawDescOnce sync.Once
	file_proto_controlplane_proto_rawDescData = file_proto_controlplane_proto_rawDesc
)

func file_proto_controlplane_proto_rawDescGZIP() []byte {
	file_proto_controlplane_proto_rawDescOnce.Do(func() {
		file_proto_controlplane_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_controlplane_proto_rawDescData)
	})
	return file_proto_controlplane_proto_rawDescData
}

var file_proto_controlplane_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_controlplane_proto_goTypes = []interface{}{
	(*WorkerCheckInRequest)(nil),            // 0: controlplane.WorkerCheckInRequest
	(*ComputeNodeRegistrationRequest)(nil),  // 1: controlplane.ComputeNodeRegistrationRequest
	(*ComputeNodeRegistrationResponse)(nil), // 2: controlplane.ComputeNodeRegistrationResponse
	(*ControlPlaneTask)(nil),                // 3: controlplane.ControlPlaneTask
	(*CreateJobRequest)(nil),                // 4: controlplane.CreateJobRequest
	(*CreateJobResponse)(nil),               // 5: controlplane.CreateJobResponse
	(*emptypb.Empty)(nil),                   // 6: google.protobuf.Empty
}
var file_proto_controlplane_proto_depIdxs = []int32{
	3, // 0: controlplane.CreateJobResponse.Tasks:type_name -> controlplane.ControlPlaneTask
	1, // 1: controlplane.ControlPlane.WorkerRegistration:input_type -> controlplane.ComputeNodeRegistrationRequest
	4, // 2: controlplane.ControlPlane.CreateJob:input_type -> controlplane.CreateJobRequest
	0, // 3: controlplane.ControlPlane.WorkerCheckIn:input_type -> controlplane.WorkerCheckInRequest
	2, // 4: controlplane.ControlPlane.WorkerRegistration:output_type -> controlplane.ComputeNodeRegistrationResponse
	5, // 5: controlplane.ControlPlane.CreateJob:output_type -> controlplane.CreateJobResponse
	6, // 6: controlplane.ControlPlane.WorkerCheckIn:output_type -> google.protobuf.Empty
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_controlplane_proto_init() }
func file_proto_controlplane_proto_init() {
	if File_proto_controlplane_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_controlplane_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkerCheckInRequest); i {
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
		file_proto_controlplane_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ComputeNodeRegistrationRequest); i {
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
		file_proto_controlplane_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ComputeNodeRegistrationResponse); i {
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
		file_proto_controlplane_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ControlPlaneTask); i {
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
		file_proto_controlplane_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateJobRequest); i {
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
		file_proto_controlplane_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateJobResponse); i {
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
			RawDescriptor: file_proto_controlplane_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_controlplane_proto_goTypes,
		DependencyIndexes: file_proto_controlplane_proto_depIdxs,
		MessageInfos:      file_proto_controlplane_proto_msgTypes,
	}.Build()
	File_proto_controlplane_proto = out.File
	file_proto_controlplane_proto_rawDesc = nil
	file_proto_controlplane_proto_goTypes = nil
	file_proto_controlplane_proto_depIdxs = nil
}
