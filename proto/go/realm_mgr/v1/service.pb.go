// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: realm_mgr/v1/service.proto

package realm_mgr_v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_realm_mgr_v1_service_proto protoreflect.FileDescriptor

var file_realm_mgr_v1_service_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x5f, 0x6d, 0x67, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x72, 0x65,
	0x61, 0x6c, 0x6d, 0x5f, 0x6d, 0x67, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x18, 0x72, 0x65, 0x61, 0x6c,
	0x6d, 0x5f, 0x6d, 0x67, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x32, 0xe7, 0x02, 0x0a, 0x13, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4b, 0x0a, 0x08,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x12, 0x1d, 0x2e, 0x72, 0x65, 0x61, 0x6c, 0x6d,
	0x5f, 0x6d, 0x67, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x5f,
	0x6d, 0x67, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x54, 0x0a, 0x0b, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x12, 0x20, 0x2e, 0x72, 0x65, 0x61, 0x6c, 0x6d,
	0x5f, 0x6d, 0x67, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x61, 0x6c, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x72, 0x65, 0x61,
	0x6c, 0x6d, 0x5f, 0x6d, 0x67, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x57, 0x0a, 0x0c, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x12,
	0x21, 0x2e, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x5f, 0x6d, 0x67, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x52,
	0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x22, 0x2e, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x5f, 0x6d, 0x67, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x54, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x12, 0x20, 0x2e, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x5f,
	0x6d, 0x67, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x61,
	0x6c, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x72, 0x65, 0x61, 0x6c,
	0x6d, 0x5f, 0x6d, 0x67, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x61, 0x6c, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x1b,
	0x5a, 0x19, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x5f, 0x6d, 0x67, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x72,
	0x65, 0x61, 0x6c, 0x6d, 0x5f, 0x6d, 0x67, 0x72, 0x5f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var file_realm_mgr_v1_service_proto_goTypes = []interface{}{
	(*GetRealmRequest)(nil),      // 0: realm_mgr.v1.GetRealmRequest
	(*CreateRealmRequest)(nil),   // 1: realm_mgr.v1.CreateRealmRequest
	(*ReleaseRealmRequest)(nil),  // 2: realm_mgr.v1.ReleaseRealmRequest
	(*UpdateRealmRequest)(nil),   // 3: realm_mgr.v1.UpdateRealmRequest
	(*GetRealmResponse)(nil),     // 4: realm_mgr.v1.GetRealmResponse
	(*CreateRealmResponse)(nil),  // 5: realm_mgr.v1.CreateRealmResponse
	(*ReleaseRealmResponse)(nil), // 6: realm_mgr.v1.ReleaseRealmResponse
	(*UpdateRealmResponse)(nil),  // 7: realm_mgr.v1.UpdateRealmResponse
}
var file_realm_mgr_v1_service_proto_depIdxs = []int32{
	0, // 0: realm_mgr.v1.RealmManagerService.GetRealm:input_type -> realm_mgr.v1.GetRealmRequest
	1, // 1: realm_mgr.v1.RealmManagerService.CreateRealm:input_type -> realm_mgr.v1.CreateRealmRequest
	2, // 2: realm_mgr.v1.RealmManagerService.ReleaseRealm:input_type -> realm_mgr.v1.ReleaseRealmRequest
	3, // 3: realm_mgr.v1.RealmManagerService.UpdateRealm:input_type -> realm_mgr.v1.UpdateRealmRequest
	4, // 4: realm_mgr.v1.RealmManagerService.GetRealm:output_type -> realm_mgr.v1.GetRealmResponse
	5, // 5: realm_mgr.v1.RealmManagerService.CreateRealm:output_type -> realm_mgr.v1.CreateRealmResponse
	6, // 6: realm_mgr.v1.RealmManagerService.ReleaseRealm:output_type -> realm_mgr.v1.ReleaseRealmResponse
	7, // 7: realm_mgr.v1.RealmManagerService.UpdateRealm:output_type -> realm_mgr.v1.UpdateRealmResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_realm_mgr_v1_service_proto_init() }
func file_realm_mgr_v1_service_proto_init() {
	if File_realm_mgr_v1_service_proto != nil {
		return
	}
	file_realm_mgr_v1_realm_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_realm_mgr_v1_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_realm_mgr_v1_service_proto_goTypes,
		DependencyIndexes: file_realm_mgr_v1_service_proto_depIdxs,
	}.Build()
	File_realm_mgr_v1_service_proto = out.File
	file_realm_mgr_v1_service_proto_rawDesc = nil
	file_realm_mgr_v1_service_proto_goTypes = nil
	file_realm_mgr_v1_service_proto_depIdxs = nil
}
