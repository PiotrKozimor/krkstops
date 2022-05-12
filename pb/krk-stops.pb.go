// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: pb/krk-stops.proto

package pb

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

type Endpoint int32

const (
	Endpoint_BUS  Endpoint = 0
	Endpoint_TRAM Endpoint = 1
	Endpoint_ALL  Endpoint = 2
)

// Enum value maps for Endpoint.
var (
	Endpoint_name = map[int32]string{
		0: "BUS",
		1: "TRAM",
		2: "ALL",
	}
	Endpoint_value = map[string]int32{
		"BUS":  0,
		"TRAM": 1,
		"ALL":  2,
	}
)

func (x Endpoint) Enum() *Endpoint {
	p := new(Endpoint)
	*p = x
	return p
}

func (x Endpoint) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Endpoint) Descriptor() protoreflect.EnumDescriptor {
	return file_pb_krk_stops_proto_enumTypes[0].Descriptor()
}

func (Endpoint) Type() protoreflect.EnumType {
	return &file_pb_krk_stops_proto_enumTypes[0]
}

func (x Endpoint) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Endpoint.Descriptor instead.
func (Endpoint) EnumDescriptor() ([]byte, []int) {
	return file_pb_krk_stops_proto_rawDescGZIP(), []int{0}
}

type Installation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int32   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Latitude  float32 `protobuf:"fixed32,2,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude float32 `protobuf:"fixed32,3,opt,name=longitude,proto3" json:"longitude,omitempty"`
}

func (x *Installation) Reset() {
	*x = Installation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_krk_stops_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Installation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Installation) ProtoMessage() {}

func (x *Installation) ProtoReflect() protoreflect.Message {
	mi := &file_pb_krk_stops_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Installation.ProtoReflect.Descriptor instead.
func (*Installation) Descriptor() ([]byte, []int) {
	return file_pb_krk_stops_proto_rawDescGZIP(), []int{0}
}

func (x *Installation) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Installation) GetLatitude() float32 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *Installation) GetLongitude() float32 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

type Airly struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Caqi        int32   `protobuf:"varint,1,opt,name=caqi,proto3" json:"caqi,omitempty"`
	ColorStr    string  `protobuf:"bytes,2,opt,name=colorStr,proto3" json:"colorStr,omitempty"`
	Humidity    int32   `protobuf:"varint,3,opt,name=humidity,proto3" json:"humidity,omitempty"`
	Temperature float32 `protobuf:"fixed32,4,opt,name=temperature,proto3" json:"temperature,omitempty"`
	Color       uint32  `protobuf:"varint,5,opt,name=color,proto3" json:"color,omitempty"`
}

func (x *Airly) Reset() {
	*x = Airly{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_krk_stops_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Airly) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Airly) ProtoMessage() {}

func (x *Airly) ProtoReflect() protoreflect.Message {
	mi := &file_pb_krk_stops_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Airly.ProtoReflect.Descriptor instead.
func (*Airly) Descriptor() ([]byte, []int) {
	return file_pb_krk_stops_proto_rawDescGZIP(), []int{1}
}

func (x *Airly) GetCaqi() int32 {
	if x != nil {
		return x.Caqi
	}
	return 0
}

func (x *Airly) GetColorStr() string {
	if x != nil {
		return x.ColorStr
	}
	return ""
}

func (x *Airly) GetHumidity() int32 {
	if x != nil {
		return x.Humidity
	}
	return 0
}

func (x *Airly) GetTemperature() float32 {
	if x != nil {
		return x.Temperature
	}
	return 0
}

func (x *Airly) GetColor() uint32 {
	if x != nil {
		return x.Color
	}
	return 0
}

type InstallationLocation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Latitude  float32 `protobuf:"fixed32,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude float32 `protobuf:"fixed32,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
}

func (x *InstallationLocation) Reset() {
	*x = InstallationLocation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_krk_stops_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstallationLocation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstallationLocation) ProtoMessage() {}

func (x *InstallationLocation) ProtoReflect() protoreflect.Message {
	mi := &file_pb_krk_stops_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstallationLocation.ProtoReflect.Descriptor instead.
func (*InstallationLocation) Descriptor() ([]byte, []int) {
	return file_pb_krk_stops_proto_rawDescGZIP(), []int{2}
}

func (x *InstallationLocation) GetLatitude() float32 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *InstallationLocation) GetLongitude() float32 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

type Departure struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RelativeTime       int32    `protobuf:"varint,1,opt,name=relativeTime,proto3" json:"relativeTime,omitempty"`
	PlannedTime        string   `protobuf:"bytes,2,opt,name=plannedTime,proto3" json:"plannedTime,omitempty"`
	Direction          string   `protobuf:"bytes,3,opt,name=direction,proto3" json:"direction,omitempty"`
	PatternText        string   `protobuf:"bytes,4,opt,name=patternText,proto3" json:"patternText,omitempty"`
	Color              uint32   `protobuf:"varint,5,opt,name=color,proto3" json:"color,omitempty"`
	RelativeTimeParsed string   `protobuf:"bytes,6,opt,name=relativeTimeParsed,proto3" json:"relativeTimeParsed,omitempty"`
	Predicted          bool     `protobuf:"varint,7,opt,name=predicted,proto3" json:"predicted,omitempty"`
	Type               Endpoint `protobuf:"varint,8,opt,name=type,proto3,enum=Endpoint" json:"type,omitempty"`
}

func (x *Departure) Reset() {
	*x = Departure{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_krk_stops_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Departure) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Departure) ProtoMessage() {}

func (x *Departure) ProtoReflect() protoreflect.Message {
	mi := &file_pb_krk_stops_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Departure.ProtoReflect.Descriptor instead.
func (*Departure) Descriptor() ([]byte, []int) {
	return file_pb_krk_stops_proto_rawDescGZIP(), []int{3}
}

func (x *Departure) GetRelativeTime() int32 {
	if x != nil {
		return x.RelativeTime
	}
	return 0
}

func (x *Departure) GetPlannedTime() string {
	if x != nil {
		return x.PlannedTime
	}
	return ""
}

func (x *Departure) GetDirection() string {
	if x != nil {
		return x.Direction
	}
	return ""
}

func (x *Departure) GetPatternText() string {
	if x != nil {
		return x.PatternText
	}
	return ""
}

func (x *Departure) GetColor() uint32 {
	if x != nil {
		return x.Color
	}
	return 0
}

func (x *Departure) GetRelativeTimeParsed() string {
	if x != nil {
		return x.RelativeTimeParsed
	}
	return ""
}

func (x *Departure) GetPredicted() bool {
	if x != nil {
		return x.Predicted
	}
	return false
}

func (x *Departure) GetType() Endpoint {
	if x != nil {
		return x.Type
	}
	return Endpoint_BUS
}

type Stop struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShortName string   `protobuf:"bytes,1,opt,name=shortName,proto3" json:"shortName,omitempty"`
	Name      string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Id        uint32   `protobuf:"varint,3,opt,name=id,proto3" json:"id,omitempty"`
	Type      Endpoint `protobuf:"varint,4,opt,name=type,proto3,enum=Endpoint" json:"type,omitempty"`
}

func (x *Stop) Reset() {
	*x = Stop{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_krk_stops_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stop) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stop) ProtoMessage() {}

func (x *Stop) ProtoReflect() protoreflect.Message {
	mi := &file_pb_krk_stops_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stop.ProtoReflect.Descriptor instead.
func (*Stop) Descriptor() ([]byte, []int) {
	return file_pb_krk_stops_proto_rawDescGZIP(), []int{4}
}

func (x *Stop) GetShortName() string {
	if x != nil {
		return x.ShortName
	}
	return ""
}

func (x *Stop) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Stop) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Stop) GetType() Endpoint {
	if x != nil {
		return x.Type
	}
	return Endpoint_BUS
}

type Stops struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stops []*Stop `protobuf:"bytes,1,rep,name=stops,proto3" json:"stops,omitempty"`
}

func (x *Stops) Reset() {
	*x = Stops{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_krk_stops_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stops) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stops) ProtoMessage() {}

func (x *Stops) ProtoReflect() protoreflect.Message {
	mi := &file_pb_krk_stops_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stops.ProtoReflect.Descriptor instead.
func (*Stops) Descriptor() ([]byte, []int) {
	return file_pb_krk_stops_proto_rawDescGZIP(), []int{5}
}

func (x *Stops) GetStops() []*Stop {
	if x != nil {
		return x.Stops
	}
	return nil
}

type Departures struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Departures []*Departure `protobuf:"bytes,1,rep,name=departures,proto3" json:"departures,omitempty"`
}

func (x *Departures) Reset() {
	*x = Departures{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_krk_stops_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Departures) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Departures) ProtoMessage() {}

func (x *Departures) ProtoReflect() protoreflect.Message {
	mi := &file_pb_krk_stops_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Departures.ProtoReflect.Descriptor instead.
func (*Departures) Descriptor() ([]byte, []int) {
	return file_pb_krk_stops_proto_rawDescGZIP(), []int{6}
}

func (x *Departures) GetDepartures() []*Departure {
	if x != nil {
		return x.Departures
	}
	return nil
}

type StopSearch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Query string `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
}

func (x *StopSearch) Reset() {
	*x = StopSearch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_krk_stops_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopSearch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopSearch) ProtoMessage() {}

func (x *StopSearch) ProtoReflect() protoreflect.Message {
	mi := &file_pb_krk_stops_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopSearch.ProtoReflect.Descriptor instead.
func (*StopSearch) Descriptor() ([]byte, []int) {
	return file_pb_krk_stops_proto_rawDescGZIP(), []int{7}
}

func (x *StopSearch) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

var File_pb_krk_stops_proto protoreflect.FileDescriptor

var file_pb_krk_stops_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x62, 0x2f, 0x6b, 0x72, 0x6b, 0x2d, 0x73, 0x74, 0x6f, 0x70, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x58, 0x0a, 0x0c, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x02, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x22, 0x8b,
	0x01, 0x0a, 0x05, 0x41, 0x69, 0x72, 0x6c, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x61, 0x71, 0x69,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x61, 0x71, 0x69, 0x12, 0x1a, 0x0a, 0x08,
	0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x53, 0x74, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x53, 0x74, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x75, 0x6d, 0x69,
	0x64, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x68, 0x75, 0x6d, 0x69,
	0x64, 0x69, 0x74, 0x79, 0x12, 0x20, 0x0a, 0x0b, 0x74, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0b, 0x74, 0x65, 0x6d, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x22, 0x50, 0x0a, 0x14,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x6f, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x02, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x22, 0x94,
	0x02, 0x0a, 0x09, 0x44, 0x65, 0x70, 0x61, 0x72, 0x74, 0x75, 0x72, 0x65, 0x12, 0x22, 0x0a, 0x0c,
	0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x76, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0c, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x76, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x70, 0x6c, 0x61, 0x6e, 0x6e, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x6c, 0x61, 0x6e, 0x6e, 0x65, 0x64, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x20, 0x0a, 0x0b, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x54, 0x65, 0x78, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x54, 0x65,
	0x78, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x12, 0x2e, 0x0a, 0x12, 0x72, 0x65, 0x6c, 0x61,
	0x74, 0x69, 0x76, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x50, 0x61, 0x72, 0x73, 0x65, 0x64, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x76, 0x65, 0x54, 0x69,
	0x6d, 0x65, 0x50, 0x61, 0x72, 0x73, 0x65, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x65, 0x64,
	0x69, 0x63, 0x74, 0x65, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x70, 0x72, 0x65,
	0x64, 0x69, 0x63, 0x74, 0x65, 0x64, 0x12, 0x1d, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x09, 0x2e, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x67, 0x0a, 0x04, 0x53, 0x74, 0x6f, 0x70, 0x12, 0x1c, 0x0a,
	0x09, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x1d, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x09, 0x2e,
	0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x24,
	0x0a, 0x05, 0x53, 0x74, 0x6f, 0x70, 0x73, 0x12, 0x1b, 0x0a, 0x05, 0x73, 0x74, 0x6f, 0x70, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x52, 0x05, 0x73,
	0x74, 0x6f, 0x70, 0x73, 0x22, 0x38, 0x0a, 0x0a, 0x44, 0x65, 0x70, 0x61, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x73, 0x12, 0x2a, 0x0a, 0x0a, 0x64, 0x65, 0x70, 0x61, 0x72, 0x74, 0x75, 0x72, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x44, 0x65, 0x70, 0x61, 0x72, 0x74, 0x75,
	0x72, 0x65, 0x52, 0x0a, 0x64, 0x65, 0x70, 0x61, 0x72, 0x74, 0x75, 0x72, 0x65, 0x73, 0x22, 0x22,
	0x0a, 0x0a, 0x53, 0x74, 0x6f, 0x70, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x14, 0x0a, 0x05,
	0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65,
	0x72, 0x79, 0x2a, 0x26, 0x0a, 0x08, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x07,
	0x0a, 0x03, 0x42, 0x55, 0x53, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x54, 0x52, 0x41, 0x4d, 0x10,
	0x01, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x4c, 0x4c, 0x10, 0x02, 0x32, 0xfe, 0x01, 0x0a, 0x08, 0x4b,
	0x72, 0x6b, 0x53, 0x74, 0x6f, 0x70, 0x73, 0x12, 0x23, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x41, 0x69,
	0x72, 0x6c, 0x79, 0x12, 0x0d, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x1a, 0x06, 0x2e, 0x41, 0x69, 0x72, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x1c,
	0x46, 0x69, 0x6e, 0x64, 0x4e, 0x65, 0x61, 0x72, 0x65, 0x73, 0x74, 0x41, 0x69, 0x72, 0x6c, 0x79,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x15, 0x2e, 0x49,
	0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x1a, 0x0d, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x41, 0x69, 0x72, 0x6c, 0x79,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0d, 0x2e, 0x49,
	0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x0d, 0x2e, 0x49, 0x6e,
	0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12, 0x26, 0x0a, 0x0e,
	0x47, 0x65, 0x74, 0x44, 0x65, 0x70, 0x61, 0x72, 0x74, 0x75, 0x72, 0x65, 0x73, 0x32, 0x12, 0x05,
	0x2e, 0x53, 0x74, 0x6f, 0x70, 0x1a, 0x0b, 0x2e, 0x44, 0x65, 0x70, 0x61, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x73, 0x22, 0x00, 0x12, 0x25, 0x0a, 0x0c, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x53, 0x74,
	0x6f, 0x70, 0x73, 0x32, 0x12, 0x0b, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x53, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x1a, 0x06, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x73, 0x22, 0x00, 0x42, 0x06, 0x5a, 0x04, 0x2e,
	0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_krk_stops_proto_rawDescOnce sync.Once
	file_pb_krk_stops_proto_rawDescData = file_pb_krk_stops_proto_rawDesc
)

func file_pb_krk_stops_proto_rawDescGZIP() []byte {
	file_pb_krk_stops_proto_rawDescOnce.Do(func() {
		file_pb_krk_stops_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_krk_stops_proto_rawDescData)
	})
	return file_pb_krk_stops_proto_rawDescData
}

var file_pb_krk_stops_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pb_krk_stops_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_pb_krk_stops_proto_goTypes = []interface{}{
	(Endpoint)(0),                // 0: Endpoint
	(*Installation)(nil),         // 1: Installation
	(*Airly)(nil),                // 2: Airly
	(*InstallationLocation)(nil), // 3: InstallationLocation
	(*Departure)(nil),            // 4: Departure
	(*Stop)(nil),                 // 5: Stop
	(*Stops)(nil),                // 6: Stops
	(*Departures)(nil),           // 7: Departures
	(*StopSearch)(nil),           // 8: StopSearch
}
var file_pb_krk_stops_proto_depIdxs = []int32{
	0, // 0: Departure.type:type_name -> Endpoint
	0, // 1: Stop.type:type_name -> Endpoint
	5, // 2: Stops.stops:type_name -> Stop
	4, // 3: Departures.departures:type_name -> Departure
	1, // 4: KrkStops.GetAirly:input_type -> Installation
	3, // 5: KrkStops.FindNearestAirlyInstallation:input_type -> InstallationLocation
	1, // 6: KrkStops.GetAirlyInstallation:input_type -> Installation
	5, // 7: KrkStops.GetDepartures2:input_type -> Stop
	8, // 8: KrkStops.SearchStops2:input_type -> StopSearch
	2, // 9: KrkStops.GetAirly:output_type -> Airly
	1, // 10: KrkStops.FindNearestAirlyInstallation:output_type -> Installation
	1, // 11: KrkStops.GetAirlyInstallation:output_type -> Installation
	7, // 12: KrkStops.GetDepartures2:output_type -> Departures
	6, // 13: KrkStops.SearchStops2:output_type -> Stops
	9, // [9:14] is the sub-list for method output_type
	4, // [4:9] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_pb_krk_stops_proto_init() }
func file_pb_krk_stops_proto_init() {
	if File_pb_krk_stops_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_krk_stops_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Installation); i {
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
		file_pb_krk_stops_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Airly); i {
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
		file_pb_krk_stops_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstallationLocation); i {
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
		file_pb_krk_stops_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Departure); i {
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
		file_pb_krk_stops_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stop); i {
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
		file_pb_krk_stops_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stops); i {
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
		file_pb_krk_stops_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Departures); i {
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
		file_pb_krk_stops_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopSearch); i {
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
			RawDescriptor: file_pb_krk_stops_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_krk_stops_proto_goTypes,
		DependencyIndexes: file_pb_krk_stops_proto_depIdxs,
		EnumInfos:         file_pb_krk_stops_proto_enumTypes,
		MessageInfos:      file_pb_krk_stops_proto_msgTypes,
	}.Build()
	File_pb_krk_stops_proto = out.File
	file_pb_krk_stops_proto_rawDesc = nil
	file_pb_krk_stops_proto_goTypes = nil
	file_pb_krk_stops_proto_depIdxs = nil
}
