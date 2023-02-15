// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: request.proto

package proto

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

type GetPhaseReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *GetPhaseReq) Reset() {
	*x = GetPhaseReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_request_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPhaseReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPhaseReq) ProtoMessage() {}

func (x *GetPhaseReq) ProtoReflect() protoreflect.Message {
	mi := &file_request_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPhaseReq.ProtoReflect.Descriptor instead.
func (*GetPhaseReq) Descriptor() ([]byte, []int) {
	return file_request_proto_rawDescGZIP(), []int{0}
}

func (x *GetPhaseReq) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type GetPhaseRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Val string     `protobuf:"bytes,1,opt,name=val,proto3" json:"val,omitempty"`
	Ts  *TimeStamp `protobuf:"bytes,2,opt,name=ts,proto3" json:"ts,omitempty"`
}

func (x *GetPhaseRsp) Reset() {
	*x = GetPhaseRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_request_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPhaseRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPhaseRsp) ProtoMessage() {}

func (x *GetPhaseRsp) ProtoReflect() protoreflect.Message {
	mi := &file_request_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPhaseRsp.ProtoReflect.Descriptor instead.
func (*GetPhaseRsp) Descriptor() ([]byte, []int) {
	return file_request_proto_rawDescGZIP(), []int{1}
}

func (x *GetPhaseRsp) GetVal() string {
	if x != nil {
		return x.Val
	}
	return ""
}

func (x *GetPhaseRsp) GetTs() *TimeStamp {
	if x != nil {
		return x.Ts
	}
	return nil
}

type SetPhaseReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string     `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string     `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Ts    *TimeStamp `protobuf:"bytes,3,opt,name=ts,proto3" json:"ts,omitempty"`
}

func (x *SetPhaseReq) Reset() {
	*x = SetPhaseReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_request_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetPhaseReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetPhaseReq) ProtoMessage() {}

func (x *SetPhaseReq) ProtoReflect() protoreflect.Message {
	mi := &file_request_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetPhaseReq.ProtoReflect.Descriptor instead.
func (*SetPhaseReq) Descriptor() ([]byte, []int) {
	return file_request_proto_rawDescGZIP(), []int{2}
}

func (x *SetPhaseReq) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SetPhaseReq) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *SetPhaseReq) GetTs() *TimeStamp {
	if x != nil {
		return x.Ts
	}
	return nil
}

type SetPhaseRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SetPhaseRsp) Reset() {
	*x = SetPhaseRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_request_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetPhaseRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetPhaseRsp) ProtoMessage() {}

func (x *SetPhaseRsp) ProtoReflect() protoreflect.Message {
	mi := &file_request_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetPhaseRsp.ProtoReflect.Descriptor instead.
func (*SetPhaseRsp) Descriptor() ([]byte, []int) {
	return file_request_proto_rawDescGZIP(), []int{3}
}

type TimeStamp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestNumber uint64 `protobuf:"varint,1,opt,name=requestNumber,proto3" json:"requestNumber,omitempty"`
	ClientID      string `protobuf:"bytes,2,opt,name=clientID,proto3" json:"clientID,omitempty"`
}

func (x *TimeStamp) Reset() {
	*x = TimeStamp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_request_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeStamp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeStamp) ProtoMessage() {}

func (x *TimeStamp) ProtoReflect() protoreflect.Message {
	mi := &file_request_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeStamp.ProtoReflect.Descriptor instead.
func (*TimeStamp) Descriptor() ([]byte, []int) {
	return file_request_proto_rawDescGZIP(), []int{4}
}

func (x *TimeStamp) GetRequestNumber() uint64 {
	if x != nil {
		return x.RequestNumber
	}
	return 0
}

func (x *TimeStamp) GetClientID() string {
	if x != nil {
		return x.ClientID
	}
	return ""
}

var File_request_proto protoreflect.FileDescriptor

var file_request_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x1f, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x52, 0x65, 0x71, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x22, 0x3b, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x52, 0x73, 0x70, 0x12,
	0x10, 0x0a, 0x03, 0x76, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x76, 0x61,
	0x6c, 0x12, 0x1a, 0x0a, 0x02, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x02, 0x74, 0x73, 0x22, 0x51, 0x0a,
	0x0b, 0x53, 0x65, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x1a, 0x0a, 0x02, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0a, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x02, 0x74, 0x73,
	0x22, 0x0d, 0x0a, 0x0b, 0x53, 0x65, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x52, 0x73, 0x70, 0x22,
	0x4d, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x24, 0x0a, 0x0d,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x32, 0x65,
	0x0a, 0x0f, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x73, 0x12, 0x28, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x12, 0x0c, 0x2e,
	0x47, 0x65, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x0c, 0x2e, 0x47, 0x65,
	0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x52, 0x73, 0x70, 0x22, 0x00, 0x12, 0x28, 0x0a, 0x08, 0x53,
	0x65, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x12, 0x0c, 0x2e, 0x53, 0x65, 0x74, 0x50, 0x68, 0x61,
	0x73, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x0c, 0x2e, 0x53, 0x65, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65,
	0x52, 0x73, 0x70, 0x22, 0x00, 0x42, 0x11, 0x5a, 0x0f, 0x2e, 0x2e, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_request_proto_rawDescOnce sync.Once
	file_request_proto_rawDescData = file_request_proto_rawDesc
)

func file_request_proto_rawDescGZIP() []byte {
	file_request_proto_rawDescOnce.Do(func() {
		file_request_proto_rawDescData = protoimpl.X.CompressGZIP(file_request_proto_rawDescData)
	})
	return file_request_proto_rawDescData
}

var file_request_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_request_proto_goTypes = []interface{}{
	(*GetPhaseReq)(nil), // 0: GetPhaseReq
	(*GetPhaseRsp)(nil), // 1: GetPhaseRsp
	(*SetPhaseReq)(nil), // 2: SetPhaseReq
	(*SetPhaseRsp)(nil), // 3: SetPhaseRsp
	(*TimeStamp)(nil),   // 4: TimeStamp
}
var file_request_proto_depIdxs = []int32{
	4, // 0: GetPhaseRsp.ts:type_name -> TimeStamp
	4, // 1: SetPhaseReq.ts:type_name -> TimeStamp
	0, // 2: SharedRegisters.GetPhase:input_type -> GetPhaseReq
	2, // 3: SharedRegisters.SetPhase:input_type -> SetPhaseReq
	1, // 4: SharedRegisters.GetPhase:output_type -> GetPhaseRsp
	3, // 5: SharedRegisters.SetPhase:output_type -> SetPhaseRsp
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_request_proto_init() }
func file_request_proto_init() {
	if File_request_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_request_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPhaseReq); i {
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
		file_request_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPhaseRsp); i {
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
		file_request_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetPhaseReq); i {
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
		file_request_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetPhaseRsp); i {
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
		file_request_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeStamp); i {
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
			RawDescriptor: file_request_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_request_proto_goTypes,
		DependencyIndexes: file_request_proto_depIdxs,
		MessageInfos:      file_request_proto_msgTypes,
	}.Build()
	File_request_proto = out.File
	file_request_proto_rawDesc = nil
	file_request_proto_goTypes = nil
	file_request_proto_depIdxs = nil
}
