// Code generated by protoc-gen-go. DO NOT EDIT.
// source: unit.proto

package pb

/*
compile command:  /export/tools/protobuf-3.1.0/bin/protoc --go_out=. unit.proto
*/

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EVENT_TYPE int32

const (
	EVENT_TYPE_UNKNOWN_EVENT            EVENT_TYPE = 0
	EVENT_TYPE_START_EVENT_V3           EVENT_TYPE = 1
	EVENT_TYPE_QUERY_EVENT              EVENT_TYPE = 2
	EVENT_TYPE_STOP_EVENT               EVENT_TYPE = 3
	EVENT_TYPE_ROTATE_EVENT             EVENT_TYPE = 4
	EVENT_TYPE_INTVAR_EVENT             EVENT_TYPE = 5
	EVENT_TYPE_LOAD_EVENT               EVENT_TYPE = 6
	EVENT_TYPE_SLAVE_EVENT              EVENT_TYPE = 7
	EVENT_TYPE_CREATE_FILE_EVENT        EVENT_TYPE = 8
	EVENT_TYPE_APPEND_BLOCK_EVENT       EVENT_TYPE = 9
	EVENT_TYPE_EXEC_LOAD_EVENT          EVENT_TYPE = 10
	EVENT_TYPE_DELETE_FILE_EVENT        EVENT_TYPE = 11
	EVENT_TYPE_NEW_LOAD_EVENT           EVENT_TYPE = 12
	EVENT_TYPE_RAND_EVENT               EVENT_TYPE = 13
	EVENT_TYPE_USER_VAR_EVENT           EVENT_TYPE = 14
	EVENT_TYPE_FORMAT_DESCRIPTION_EVENT EVENT_TYPE = 15
	EVENT_TYPE_XID_EVENT                EVENT_TYPE = 16
	EVENT_TYPE_BEGIN_LOAD_QUERY_EVENT   EVENT_TYPE = 17
	EVENT_TYPE_EXECUTE_LOAD_QUERY_EVENT EVENT_TYPE = 18
	EVENT_TYPE_TABLE_MAP_EVENT          EVENT_TYPE = 19
	EVENT_TYPE_WRITE_ROWS_EVENTv0       EVENT_TYPE = 20
	EVENT_TYPE_UPDATE_ROWS_EVENTv0      EVENT_TYPE = 21
	EVENT_TYPE_DELETE_ROWS_EVENTv0      EVENT_TYPE = 22
	EVENT_TYPE_WRITE_ROWS_EVENTv1       EVENT_TYPE = 23
	EVENT_TYPE_UPDATE_ROWS_EVENTv1      EVENT_TYPE = 24
	EVENT_TYPE_DELETE_ROWS_EVENTv1      EVENT_TYPE = 25
	EVENT_TYPE_INCIDENT_EVENT           EVENT_TYPE = 26
	EVENT_TYPE_HEARTBEAT_EVENT          EVENT_TYPE = 27
	EVENT_TYPE_IGNORABLE_EVENT          EVENT_TYPE = 28
	EVENT_TYPE_ROWS_QUERY_EVENT         EVENT_TYPE = 29
	EVENT_TYPE_WRITE_ROWS_EVENTv2       EVENT_TYPE = 30
	EVENT_TYPE_UPDATE_ROWS_EVENTv2      EVENT_TYPE = 31
	EVENT_TYPE_DELETE_ROWS_EVENTv2      EVENT_TYPE = 32
	EVENT_TYPE_GTID_EVENT               EVENT_TYPE = 33
	EVENT_TYPE_ANONYMOUS_GTID_EVENT     EVENT_TYPE = 34
	EVENT_TYPE_PREVIOUS_GTIDS_EVENT     EVENT_TYPE = 35
)

var EVENT_TYPE_name = map[int32]string{
	0:  "UNKNOWN_EVENT",
	1:  "START_EVENT_V3",
	2:  "QUERY_EVENT",
	3:  "STOP_EVENT",
	4:  "ROTATE_EVENT",
	5:  "INTVAR_EVENT",
	6:  "LOAD_EVENT",
	7:  "SLAVE_EVENT",
	8:  "CREATE_FILE_EVENT",
	9:  "APPEND_BLOCK_EVENT",
	10: "EXEC_LOAD_EVENT",
	11: "DELETE_FILE_EVENT",
	12: "NEW_LOAD_EVENT",
	13: "RAND_EVENT",
	14: "USER_VAR_EVENT",
	15: "FORMAT_DESCRIPTION_EVENT",
	16: "XID_EVENT",
	17: "BEGIN_LOAD_QUERY_EVENT",
	18: "EXECUTE_LOAD_QUERY_EVENT",
	19: "TABLE_MAP_EVENT",
	20: "WRITE_ROWS_EVENTv0",
	21: "UPDATE_ROWS_EVENTv0",
	22: "DELETE_ROWS_EVENTv0",
	23: "WRITE_ROWS_EVENTv1",
	24: "UPDATE_ROWS_EVENTv1",
	25: "DELETE_ROWS_EVENTv1",
	26: "INCIDENT_EVENT",
	27: "HEARTBEAT_EVENT",
	28: "IGNORABLE_EVENT",
	29: "ROWS_QUERY_EVENT",
	30: "WRITE_ROWS_EVENTv2",
	31: "UPDATE_ROWS_EVENTv2",
	32: "DELETE_ROWS_EVENTv2",
	33: "GTID_EVENT",
	34: "ANONYMOUS_GTID_EVENT",
	35: "PREVIOUS_GTIDS_EVENT",
}
var EVENT_TYPE_value = map[string]int32{
	"UNKNOWN_EVENT":            0,
	"START_EVENT_V3":           1,
	"QUERY_EVENT":              2,
	"STOP_EVENT":               3,
	"ROTATE_EVENT":             4,
	"INTVAR_EVENT":             5,
	"LOAD_EVENT":               6,
	"SLAVE_EVENT":              7,
	"CREATE_FILE_EVENT":        8,
	"APPEND_BLOCK_EVENT":       9,
	"EXEC_LOAD_EVENT":          10,
	"DELETE_FILE_EVENT":        11,
	"NEW_LOAD_EVENT":           12,
	"RAND_EVENT":               13,
	"USER_VAR_EVENT":           14,
	"FORMAT_DESCRIPTION_EVENT": 15,
	"XID_EVENT":                16,
	"BEGIN_LOAD_QUERY_EVENT":   17,
	"EXECUTE_LOAD_QUERY_EVENT": 18,
	"TABLE_MAP_EVENT":          19,
	"WRITE_ROWS_EVENTv0":       20,
	"UPDATE_ROWS_EVENTv0":      21,
	"DELETE_ROWS_EVENTv0":      22,
	"WRITE_ROWS_EVENTv1":       23,
	"UPDATE_ROWS_EVENTv1":      24,
	"DELETE_ROWS_EVENTv1":      25,
	"INCIDENT_EVENT":           26,
	"HEARTBEAT_EVENT":          27,
	"IGNORABLE_EVENT":          28,
	"ROWS_QUERY_EVENT":         29,
	"WRITE_ROWS_EVENTv2":       30,
	"UPDATE_ROWS_EVENTv2":      31,
	"DELETE_ROWS_EVENTv2":      32,
	"GTID_EVENT":               33,
	"ANONYMOUS_GTID_EVENT":     34,
	"PREVIOUS_GTIDS_EVENT":     35,
}

func (x EVENT_TYPE) String() string {
	return proto.EnumName(EVENT_TYPE_name, int32(x))
}
func (EVENT_TYPE) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_unit_47be242f02aa2030, []int{0}
}

type BytePair struct {
	Key                  []byte   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	ColumnBitmap         []byte   `protobuf:"bytes,3,opt,name=column_bitmap,json=columnBitmap,proto3" json:"column_bitmap,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BytePair) Reset()         { *m = BytePair{} }
func (m *BytePair) String() string { return proto.CompactTextString(m) }
func (*BytePair) ProtoMessage()    {}
func (*BytePair) Descriptor() ([]byte, []int) {
	return fileDescriptor_unit_47be242f02aa2030, []int{0}
}
func (m *BytePair) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BytePair.Unmarshal(m, b)
}
func (m *BytePair) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BytePair.Marshal(b, m, deterministic)
}
func (dst *BytePair) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BytePair.Merge(dst, src)
}
func (m *BytePair) XXX_Size() int {
	return xxx_messageInfo_BytePair.Size(m)
}
func (m *BytePair) XXX_DiscardUnknown() {
	xxx_messageInfo_BytePair.DiscardUnknown(m)
}

var xxx_messageInfo_BytePair proto.InternalMessageInfo

func (m *BytePair) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *BytePair) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *BytePair) GetColumnBitmap() []byte {
	if m != nil {
		return m.ColumnBitmap
	}
	return nil
}

type BytesUnit struct {
	Tp                   EVENT_TYPE `protobuf:"varint,1,opt,name=tp,proto3,enum=pb.EVENT_TYPE" json:"tp,omitempty"`
	Key                  []byte     `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Before               *BytePair  `protobuf:"bytes,3,opt,name=before,proto3" json:"before,omitempty"`
	After                *BytePair  `protobuf:"bytes,4,opt,name=after,proto3" json:"after,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *BytesUnit) Reset()         { *m = BytesUnit{} }
func (m *BytesUnit) String() string { return proto.CompactTextString(m) }
func (*BytesUnit) ProtoMessage()    {}
func (*BytesUnit) Descriptor() ([]byte, []int) {
	return fileDescriptor_unit_47be242f02aa2030, []int{1}
}
func (m *BytesUnit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BytesUnit.Unmarshal(m, b)
}
func (m *BytesUnit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BytesUnit.Marshal(b, m, deterministic)
}
func (dst *BytesUnit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BytesUnit.Merge(dst, src)
}
func (m *BytesUnit) XXX_Size() int {
	return xxx_messageInfo_BytesUnit.Size(m)
}
func (m *BytesUnit) XXX_DiscardUnknown() {
	xxx_messageInfo_BytesUnit.DiscardUnknown(m)
}

var xxx_messageInfo_BytesUnit proto.InternalMessageInfo

func (m *BytesUnit) GetTp() EVENT_TYPE {
	if m != nil {
		return m.Tp
	}
	return EVENT_TYPE_UNKNOWN_EVENT
}

func (m *BytesUnit) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *BytesUnit) GetBefore() *BytePair {
	if m != nil {
		return m.Before
	}
	return nil
}

func (m *BytesUnit) GetAfter() *BytePair {
	if m != nil {
		return m.After
	}
	return nil
}

func init() {
	proto.RegisterType((*BytePair)(nil), "pb.BytePair")
	proto.RegisterType((*BytesUnit)(nil), "pb.BytesUnit")
	proto.RegisterEnum("pb.EVENT_TYPE", EVENT_TYPE_name, EVENT_TYPE_value)
}

func init() { proto.RegisterFile("unit.proto", fileDescriptor_unit_47be242f02aa2030) }

var fileDescriptor_unit_47be242f02aa2030 = []byte{
	// 552 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x94, 0xdd, 0x6e, 0xda, 0x40,
	0x10, 0x85, 0x8b, 0x93, 0xd0, 0x64, 0xc2, 0xcf, 0x32, 0x10, 0xe2, 0xa6, 0x69, 0x9a, 0x92, 0x5e,
	0x54, 0xbd, 0x40, 0x85, 0x3c, 0xc1, 0x1a, 0x6f, 0xa8, 0x15, 0xb3, 0x76, 0xd7, 0x6b, 0x08, 0x57,
	0x16, 0x54, 0x8e, 0x84, 0x9a, 0x80, 0x45, 0x4d, 0xa4, 0x3c, 0x41, 0xdf, 0xa2, 0xcf, 0x5a, 0xad,
	0x7f, 0x0a, 0x69, 0x7d, 0xc7, 0x9e, 0x33, 0xfb, 0xcd, 0x99, 0x59, 0x64, 0x80, 0xcd, 0x72, 0x11,
	0x77, 0xa3, 0xf5, 0x2a, 0x5e, 0xa1, 0x16, 0xcd, 0x3b, 0x13, 0x38, 0x34, 0x9e, 0xe3, 0xd0, 0x9d,
	0x2d, 0xd6, 0x48, 0x60, 0xef, 0x47, 0xf8, 0xac, 0x97, 0x2e, 0x4b, 0x9f, 0x2a, 0x42, 0xfd, 0xc4,
	0x16, 0x1c, 0x3c, 0xcd, 0x1e, 0x36, 0xa1, 0xae, 0x25, 0x5a, 0x7a, 0xc0, 0x2b, 0xa8, 0x7e, 0x5f,
	0x3d, 0x6c, 0x1e, 0x97, 0xc1, 0x7c, 0x11, 0x3f, 0xce, 0x22, 0x7d, 0x2f, 0x71, 0x2b, 0xa9, 0x68,
	0x24, 0x5a, 0xe7, 0x57, 0x09, 0x8e, 0x14, 0xf9, 0xa7, 0xbf, 0x5c, 0xc4, 0x78, 0x01, 0x5a, 0x1c,
	0x25, 0xe4, 0x5a, 0xbf, 0xd6, 0x8d, 0xe6, 0x5d, 0x36, 0x66, 0x5c, 0x06, 0x72, 0xea, 0x32, 0xa1,
	0xc5, 0x51, 0xde, 0x5a, 0xdb, 0xb6, 0xfe, 0x08, 0xe5, 0x79, 0x78, 0xbf, 0x5a, 0x87, 0x09, 0xfd,
	0xb8, 0x5f, 0x51, 0xb7, 0xf2, 0xa8, 0x22, 0xf3, 0xb0, 0x03, 0x07, 0xb3, 0xfb, 0x38, 0x5c, 0xeb,
	0xfb, 0x05, 0x45, 0xa9, 0xf5, 0xf9, 0x77, 0x19, 0x60, 0xdb, 0x0e, 0x1b, 0x50, 0xf5, 0xf9, 0x2d,
	0x77, 0x26, 0x3c, 0x48, 0x54, 0xf2, 0x0a, 0x11, 0x6a, 0x9e, 0xa4, 0x42, 0xa6, 0x42, 0x30, 0xbe,
	0x26, 0x25, 0xac, 0xc3, 0xf1, 0x37, 0x9f, 0x89, 0x69, 0x56, 0xa4, 0x61, 0x0d, 0xc0, 0x93, 0x8e,
	0x9b, 0x9d, 0xf7, 0x90, 0x40, 0x45, 0x38, 0x92, 0x4a, 0x96, 0x29, 0xfb, 0x4a, 0xb1, 0xb8, 0x1c,
	0x53, 0x91, 0x29, 0x07, 0xea, 0x8e, 0xed, 0x50, 0x33, 0x3b, 0x97, 0x15, 0xd4, 0xb3, 0xe9, 0x38,
	0xbf, 0xf2, 0x1a, 0x4f, 0xa0, 0x31, 0x10, 0x4c, 0x41, 0x6e, 0x2c, 0x3b, 0x97, 0x0f, 0xb1, 0x0d,
	0x48, 0x5d, 0x97, 0x71, 0x33, 0x30, 0x6c, 0x67, 0x70, 0x9b, 0xe9, 0x47, 0xd8, 0x84, 0x3a, 0xbb,
	0x63, 0x83, 0x60, 0x07, 0x0a, 0x8a, 0x61, 0x32, 0x9b, 0xbd, 0x64, 0x1c, 0xab, 0xa1, 0x38, 0x9b,
	0xec, 0x96, 0x56, 0x54, 0x1e, 0x41, 0x79, 0x7e, 0xae, 0xaa, 0x1a, 0xdf, 0x63, 0x22, 0xd8, 0x66,
	0xae, 0xe1, 0x39, 0xe8, 0x37, 0x8e, 0x18, 0x51, 0x19, 0x98, 0xcc, 0x1b, 0x08, 0xcb, 0x95, 0x96,
	0x93, 0xaf, 0xaa, 0x8e, 0x55, 0x38, 0xba, 0xb3, 0x72, 0x00, 0xc1, 0x33, 0x68, 0x1b, 0x6c, 0x68,
	0xf1, 0xb4, 0xcd, 0xee, 0xc2, 0x1a, 0x0a, 0xa4, 0xc2, 0xfa, 0x92, 0xfd, 0xef, 0xa2, 0x1a, 0x45,
	0x52, 0xc3, 0x66, 0xc1, 0x88, 0xe6, 0x3b, 0x6d, 0xaa, 0xb9, 0x27, 0xc2, 0x92, 0x2c, 0x10, 0xce,
	0xc4, 0x4b, 0xd5, 0xa7, 0x2f, 0xa4, 0x85, 0xa7, 0xd0, 0xf4, 0x5d, 0x93, 0xfe, 0x6b, 0x9c, 0x28,
	0x23, 0x9b, 0xfd, 0x85, 0xd1, 0x2e, 0x24, 0xf5, 0xc8, 0x69, 0x31, 0xa9, 0x47, 0xf4, 0x62, 0x52,
	0x8f, 0xbc, 0x51, 0x3b, 0xb2, 0xf8, 0xc0, 0x32, 0xd5, 0x3f, 0x23, 0xcd, 0x79, 0xa6, 0xc2, 0x7f,
	0x65, 0x54, 0x48, 0x83, 0xd1, 0x5c, 0x7c, 0xab, 0x44, 0x6b, 0xc8, 0x1d, 0x91, 0x4c, 0x95, 0x8a,
	0xe7, 0xd8, 0x02, 0x92, 0xf0, 0x76, 0x87, 0x7f, 0x57, 0x98, 0xae, 0x4f, 0x2e, 0x8a, 0xd3, 0xf5,
	0xc9, 0xfb, 0xe2, 0x74, 0x7d, 0x72, 0xa9, 0x5e, 0x74, 0x28, 0xff, 0x3e, 0xc8, 0x07, 0xd4, 0xa1,
	0x45, 0xb9, 0xc3, 0xa7, 0x23, 0xc7, 0xf7, 0x82, 0x1d, 0xa7, 0xa3, 0x1c, 0x57, 0xb0, 0xb1, 0x95,
	0x1b, 0x19, 0x85, 0x5c, 0xcd, 0xcb, 0xc9, 0xe7, 0xe0, 0xfa, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x1d, 0xfa, 0x07, 0x85, 0x1c, 0x04, 0x00, 0x00,
}
