# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: oomd.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from . import status_pb2 as status__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='oomd.proto',
  package='oomd',
  syntax='proto3',
  serialized_options=b'Z\010/codegen',
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n\noomd.proto\x12\x04oomd\x1a\x0cstatus.proto\x1a\x19google/protobuf/any.proto\"=\n\x10OnlineGetRequest\x12\x12\n\nentity_key\x18\x01 \x01(\t\x12\x15\n\rfeature_names\x18\x02 \x03(\t\"\x80\x01\n\x0f\x46\x65\x61tureValueMap\x12+\n\x03map\x18\x01 \x03(\x0b\x32\x1e.oomd.FeatureValueMap.MapEntry\x1a@\n\x08MapEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12#\n\x05value\x18\x02 \x01(\x0b\x32\x14.google.protobuf.Any:\x02\x38\x01\"^\n\x11OnlineGetResponse\x12\"\n\x06status\x18\x01 \x01(\x0b\x32\x12.google.rpc.Status\x12%\n\x06result\x18\x02 \x01(\x0b\x32\x15.oomd.FeatureValueMap\"C\n\x15OnlineMultiGetRequest\x12\x13\n\x0b\x65ntity_keys\x18\x01 \x03(\t\x12\x15\n\rfeature_names\x18\x02 \x03(\t\"\xbc\x01\n\x16OnlineMultiGetResponse\x12\"\n\x06status\x18\x01 \x01(\x0b\x32\x12.google.rpc.Status\x12\x38\n\x06result\x18\x02 \x03(\x0b\x32(.oomd.OnlineMultiGetResponse.ResultEntry\x1a\x44\n\x0bResultEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12$\n\x05value\x18\x02 \x01(\x0b\x32\x15.oomd.FeatureValueMap:\x02\x38\x01\"\"\n\x0bSyncRequest\x12\x13\n\x0brevision_id\x18\x01 \x01(\x05\"2\n\x0cSyncResponse\x12\"\n\x06status\x18\x01 \x01(\x0b\x32\x12.google.rpc.Status\"\x7f\n\rImportRequest\x12\x12\n\ngroup_name\x18\x01 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x02 \x01(\t\x12\x15\n\x08revision\x18\x03 \x01(\x03H\x00\x88\x01\x01\x12!\n\x03row\x18\x04 \x03(\x0b\x32\x14.google.protobuf.AnyB\x0b\n\t_revision\"I\n\x0eImportResponse\x12\"\n\x06status\x18\x01 \x01(\x0b\x32\x12.google.rpc.Status\x12\x13\n\x0brevision_id\x18\x02 \x01(\x03\"\x8e\x01\n\x13ImportByFileRequest\x12\x12\n\ngroup_name\x18\x01 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x02 \x01(\t\x12\x15\n\x08revision\x18\x03 \x01(\x03H\x00\x88\x01\x01\x12\x17\n\x0finput_file_path\x18\x04 \x01(\t\x12\x11\n\tdelimiter\x18\x05 \x01(\tB\x0b\n\t_revision\"2\n\tEntityRow\x12\x12\n\nentity_key\x18\x01 \x01(\t\x12\x11\n\tunix_time\x18\x02 \x01(\x03\"I\n\x0bJoinRequest\x12\x15\n\rfeature_names\x18\x01 \x03(\t\x12#\n\nentity_row\x18\x02 \x01(\x0b\x32\x0f.oomd.EntityRow\"l\n\x0cJoinResponse\x12\"\n\x06status\x18\x01 \x01(\x0b\x32\x12.google.rpc.Status\x12\x0e\n\x06header\x18\x02 \x03(\t\x12(\n\njoined_row\x18\x03 \x03(\x0b\x32\x14.google.protobuf.Any\"]\n\x11JoinByFileRequest\x12\x15\n\rfeature_names\x18\x01 \x03(\t\x12\x17\n\x0finput_file_path\x18\x02 \x01(\t\x12\x18\n\x10output_file_path\x18\x03 \x01(\t\"8\n\x12JoinByFileResponse\x12\"\n\x06status\x18\x01 \x01(\x0b\x32\x12.google.rpc.Status2\xba\x03\n\x04OomD\x12>\n\tOnlineGet\x12\x16.oomd.OnlineGetRequest\x1a\x17.oomd.OnlineGetResponse\"\x00\x12M\n\x0eOnlineMultiGet\x12\x1b.oomd.OnlineMultiGetRequest\x1a\x1c.oomd.OnlineMultiGetResponse\"\x00\x12/\n\x04Sync\x12\x11.oomd.SyncRequest\x1a\x12.oomd.SyncResponse\"\x00\x12\x37\n\x06Import\x12\x13.oomd.ImportRequest\x1a\x14.oomd.ImportResponse\"\x00(\x01\x12\x33\n\x04Join\x12\x11.oomd.JoinRequest\x1a\x12.oomd.JoinResponse\"\x00(\x01\x30\x01\x12\x41\n\x0cImportByFile\x12\x19.oomd.ImportByFileRequest\x1a\x14.oomd.ImportResponse\"\x00\x12\x41\n\nJoinByFile\x12\x17.oomd.JoinByFileRequest\x1a\x18.oomd.JoinByFileResponse\"\x00\x42\nZ\x08/codegenb\x06proto3'
  ,
  dependencies=[status__pb2.DESCRIPTOR,google_dot_protobuf_dot_any__pb2.DESCRIPTOR,])




_ONLINEGETREQUEST = _descriptor.Descriptor(
  name='OnlineGetRequest',
  full_name='oomd.OnlineGetRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='entity_key', full_name='oomd.OnlineGetRequest.entity_key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='feature_names', full_name='oomd.OnlineGetRequest.feature_names', index=1,
      number=2, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=61,
  serialized_end=122,
)


_FEATUREVALUEMAP_MAPENTRY = _descriptor.Descriptor(
  name='MapEntry',
  full_name='oomd.FeatureValueMap.MapEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='oomd.FeatureValueMap.MapEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='value', full_name='oomd.FeatureValueMap.MapEntry.value', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=189,
  serialized_end=253,
)

_FEATUREVALUEMAP = _descriptor.Descriptor(
  name='FeatureValueMap',
  full_name='oomd.FeatureValueMap',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='map', full_name='oomd.FeatureValueMap.map', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[_FEATUREVALUEMAP_MAPENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=125,
  serialized_end=253,
)


_ONLINEGETRESPONSE = _descriptor.Descriptor(
  name='OnlineGetResponse',
  full_name='oomd.OnlineGetResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='status', full_name='oomd.OnlineGetResponse.status', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='result', full_name='oomd.OnlineGetResponse.result', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=255,
  serialized_end=349,
)


_ONLINEMULTIGETREQUEST = _descriptor.Descriptor(
  name='OnlineMultiGetRequest',
  full_name='oomd.OnlineMultiGetRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='entity_keys', full_name='oomd.OnlineMultiGetRequest.entity_keys', index=0,
      number=1, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='feature_names', full_name='oomd.OnlineMultiGetRequest.feature_names', index=1,
      number=2, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=351,
  serialized_end=418,
)


_ONLINEMULTIGETRESPONSE_RESULTENTRY = _descriptor.Descriptor(
  name='ResultEntry',
  full_name='oomd.OnlineMultiGetResponse.ResultEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='oomd.OnlineMultiGetResponse.ResultEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='value', full_name='oomd.OnlineMultiGetResponse.ResultEntry.value', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=541,
  serialized_end=609,
)

_ONLINEMULTIGETRESPONSE = _descriptor.Descriptor(
  name='OnlineMultiGetResponse',
  full_name='oomd.OnlineMultiGetResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='status', full_name='oomd.OnlineMultiGetResponse.status', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='result', full_name='oomd.OnlineMultiGetResponse.result', index=1,
      number=2, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[_ONLINEMULTIGETRESPONSE_RESULTENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=421,
  serialized_end=609,
)


_SYNCREQUEST = _descriptor.Descriptor(
  name='SyncRequest',
  full_name='oomd.SyncRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='revision_id', full_name='oomd.SyncRequest.revision_id', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=611,
  serialized_end=645,
)


_SYNCRESPONSE = _descriptor.Descriptor(
  name='SyncResponse',
  full_name='oomd.SyncResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='status', full_name='oomd.SyncResponse.status', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=647,
  serialized_end=697,
)


_IMPORTREQUEST = _descriptor.Descriptor(
  name='ImportRequest',
  full_name='oomd.ImportRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='group_name', full_name='oomd.ImportRequest.group_name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='description', full_name='oomd.ImportRequest.description', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='revision', full_name='oomd.ImportRequest.revision', index=2,
      number=3, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='row', full_name='oomd.ImportRequest.row', index=3,
      number=4, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
    _descriptor.OneofDescriptor(
      name='_revision', full_name='oomd.ImportRequest._revision',
      index=0, containing_type=None,
      create_key=_descriptor._internal_create_key,
    fields=[]),
  ],
  serialized_start=699,
  serialized_end=826,
)


_IMPORTRESPONSE = _descriptor.Descriptor(
  name='ImportResponse',
  full_name='oomd.ImportResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='status', full_name='oomd.ImportResponse.status', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='revision_id', full_name='oomd.ImportResponse.revision_id', index=1,
      number=2, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=828,
  serialized_end=901,
)


_IMPORTBYFILEREQUEST = _descriptor.Descriptor(
  name='ImportByFileRequest',
  full_name='oomd.ImportByFileRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='group_name', full_name='oomd.ImportByFileRequest.group_name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='description', full_name='oomd.ImportByFileRequest.description', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='revision', full_name='oomd.ImportByFileRequest.revision', index=2,
      number=3, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='input_file_path', full_name='oomd.ImportByFileRequest.input_file_path', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='delimiter', full_name='oomd.ImportByFileRequest.delimiter', index=4,
      number=5, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
    _descriptor.OneofDescriptor(
      name='_revision', full_name='oomd.ImportByFileRequest._revision',
      index=0, containing_type=None,
      create_key=_descriptor._internal_create_key,
    fields=[]),
  ],
  serialized_start=904,
  serialized_end=1046,
)


_ENTITYROW = _descriptor.Descriptor(
  name='EntityRow',
  full_name='oomd.EntityRow',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='entity_key', full_name='oomd.EntityRow.entity_key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='unix_time', full_name='oomd.EntityRow.unix_time', index=1,
      number=2, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1048,
  serialized_end=1098,
)


_JOINREQUEST = _descriptor.Descriptor(
  name='JoinRequest',
  full_name='oomd.JoinRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='feature_names', full_name='oomd.JoinRequest.feature_names', index=0,
      number=1, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='entity_row', full_name='oomd.JoinRequest.entity_row', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1100,
  serialized_end=1173,
)


_JOINRESPONSE = _descriptor.Descriptor(
  name='JoinResponse',
  full_name='oomd.JoinResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='status', full_name='oomd.JoinResponse.status', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='header', full_name='oomd.JoinResponse.header', index=1,
      number=2, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='joined_row', full_name='oomd.JoinResponse.joined_row', index=2,
      number=3, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1175,
  serialized_end=1283,
)


_JOINBYFILEREQUEST = _descriptor.Descriptor(
  name='JoinByFileRequest',
  full_name='oomd.JoinByFileRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='feature_names', full_name='oomd.JoinByFileRequest.feature_names', index=0,
      number=1, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='input_file_path', full_name='oomd.JoinByFileRequest.input_file_path', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='output_file_path', full_name='oomd.JoinByFileRequest.output_file_path', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1285,
  serialized_end=1378,
)


_JOINBYFILERESPONSE = _descriptor.Descriptor(
  name='JoinByFileResponse',
  full_name='oomd.JoinByFileResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='status', full_name='oomd.JoinByFileResponse.status', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1380,
  serialized_end=1436,
)

_FEATUREVALUEMAP_MAPENTRY.fields_by_name['value'].message_type = google_dot_protobuf_dot_any__pb2._ANY
_FEATUREVALUEMAP_MAPENTRY.containing_type = _FEATUREVALUEMAP
_FEATUREVALUEMAP.fields_by_name['map'].message_type = _FEATUREVALUEMAP_MAPENTRY
_ONLINEGETRESPONSE.fields_by_name['status'].message_type = status__pb2._STATUS
_ONLINEGETRESPONSE.fields_by_name['result'].message_type = _FEATUREVALUEMAP
_ONLINEMULTIGETRESPONSE_RESULTENTRY.fields_by_name['value'].message_type = _FEATUREVALUEMAP
_ONLINEMULTIGETRESPONSE_RESULTENTRY.containing_type = _ONLINEMULTIGETRESPONSE
_ONLINEMULTIGETRESPONSE.fields_by_name['status'].message_type = status__pb2._STATUS
_ONLINEMULTIGETRESPONSE.fields_by_name['result'].message_type = _ONLINEMULTIGETRESPONSE_RESULTENTRY
_SYNCRESPONSE.fields_by_name['status'].message_type = status__pb2._STATUS
_IMPORTREQUEST.fields_by_name['row'].message_type = google_dot_protobuf_dot_any__pb2._ANY
_IMPORTREQUEST.oneofs_by_name['_revision'].fields.append(
  _IMPORTREQUEST.fields_by_name['revision'])
_IMPORTREQUEST.fields_by_name['revision'].containing_oneof = _IMPORTREQUEST.oneofs_by_name['_revision']
_IMPORTRESPONSE.fields_by_name['status'].message_type = status__pb2._STATUS
_IMPORTBYFILEREQUEST.oneofs_by_name['_revision'].fields.append(
  _IMPORTBYFILEREQUEST.fields_by_name['revision'])
_IMPORTBYFILEREQUEST.fields_by_name['revision'].containing_oneof = _IMPORTBYFILEREQUEST.oneofs_by_name['_revision']
_JOINREQUEST.fields_by_name['entity_row'].message_type = _ENTITYROW
_JOINRESPONSE.fields_by_name['status'].message_type = status__pb2._STATUS
_JOINRESPONSE.fields_by_name['joined_row'].message_type = google_dot_protobuf_dot_any__pb2._ANY
_JOINBYFILERESPONSE.fields_by_name['status'].message_type = status__pb2._STATUS
DESCRIPTOR.message_types_by_name['OnlineGetRequest'] = _ONLINEGETREQUEST
DESCRIPTOR.message_types_by_name['FeatureValueMap'] = _FEATUREVALUEMAP
DESCRIPTOR.message_types_by_name['OnlineGetResponse'] = _ONLINEGETRESPONSE
DESCRIPTOR.message_types_by_name['OnlineMultiGetRequest'] = _ONLINEMULTIGETREQUEST
DESCRIPTOR.message_types_by_name['OnlineMultiGetResponse'] = _ONLINEMULTIGETRESPONSE
DESCRIPTOR.message_types_by_name['SyncRequest'] = _SYNCREQUEST
DESCRIPTOR.message_types_by_name['SyncResponse'] = _SYNCRESPONSE
DESCRIPTOR.message_types_by_name['ImportRequest'] = _IMPORTREQUEST
DESCRIPTOR.message_types_by_name['ImportResponse'] = _IMPORTRESPONSE
DESCRIPTOR.message_types_by_name['ImportByFileRequest'] = _IMPORTBYFILEREQUEST
DESCRIPTOR.message_types_by_name['EntityRow'] = _ENTITYROW
DESCRIPTOR.message_types_by_name['JoinRequest'] = _JOINREQUEST
DESCRIPTOR.message_types_by_name['JoinResponse'] = _JOINRESPONSE
DESCRIPTOR.message_types_by_name['JoinByFileRequest'] = _JOINBYFILEREQUEST
DESCRIPTOR.message_types_by_name['JoinByFileResponse'] = _JOINBYFILERESPONSE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

OnlineGetRequest = _reflection.GeneratedProtocolMessageType('OnlineGetRequest', (_message.Message,), {
  'DESCRIPTOR' : _ONLINEGETREQUEST,
  '__module__' : 'oomd_pb2'
  # @@protoc_insertion_point(class_scope:oomd.OnlineGetRequest)
  })
_sym_db.RegisterMessage(OnlineGetRequest)

FeatureValueMap = _reflection.GeneratedProtocolMessageType('FeatureValueMap', (_message.Message,), {

  'MapEntry' : _reflection.GeneratedProtocolMessageType('MapEntry', (_message.Message,), {
    'DESCRIPTOR' : _FEATUREVALUEMAP_MAPENTRY,
    '__module__' : 'oomd_pb2'
    # @@protoc_insertion_point(class_scope:oomd.FeatureValueMap.MapEntry)
    })
  ,
  'DESCRIPTOR' : _FEATUREVALUEMAP,
  '__module__' : 'oomd_pb2'
  # @@protoc_insertion_point(class_scope:oomd.FeatureValueMap)
  })
_sym_db.RegisterMessage(FeatureValueMap)
_sym_db.RegisterMessage(FeatureValueMap.MapEntry)

OnlineGetResponse = _reflection.GeneratedProtocolMessageType('OnlineGetResponse', (_message.Message,), {
  'DESCRIPTOR' : _ONLINEGETRESPONSE,
  '__module__' : 'oomd_pb2'
  # @@protoc_insertion_point(class_scope:oomd.OnlineGetResponse)
  })
_sym_db.RegisterMessage(OnlineGetResponse)

OnlineMultiGetRequest = _reflection.GeneratedProtocolMessageType('OnlineMultiGetRequest', (_message.Message,), {
  'DESCRIPTOR' : _ONLINEMULTIGETREQUEST,
  '__module__' : 'oomd_pb2'
  # @@protoc_insertion_point(class_scope:oomd.OnlineMultiGetRequest)
  })
_sym_db.RegisterMessage(OnlineMultiGetRequest)

OnlineMultiGetResponse = _reflection.GeneratedProtocolMessageType('OnlineMultiGetResponse', (_message.Message,), {

  'ResultEntry' : _reflection.GeneratedProtocolMessageType('ResultEntry', (_message.Message,), {
    'DESCRIPTOR' : _ONLINEMULTIGETRESPONSE_RESULTENTRY,
    '__module__' : 'oomd_pb2'
    # @@protoc_insertion_point(class_scope:oomd.OnlineMultiGetResponse.ResultEntry)
    })
  ,
  'DESCRIPTOR' : _ONLINEMULTIGETRESPONSE,
  '__module__' : 'oomd_pb2'
  # @@protoc_insertion_point(class_scope:oomd.OnlineMultiGetResponse)
  })
_sym_db.RegisterMessage(OnlineMultiGetResponse)
_sym_db.RegisterMessage(OnlineMultiGetResponse.ResultEntry)

SyncRequest = _reflection.GeneratedProtocolMessageType('SyncRequest', (_message.Message,), {
  'DESCRIPTOR' : _SYNCREQUEST,
  '__module__' : 'oomd_pb2'
  # @@protoc_insertion_point(class_scope:oomd.SyncRequest)
  })
_sym_db.RegisterMessage(SyncRequest)

SyncResponse = _reflection.GeneratedProtocolMessageType('SyncResponse', (_message.Message,), {
  'DESCRIPTOR' : _SYNCRESPONSE,
  '__module__' : 'oomd_pb2'
  # @@protoc_insertion_point(class_scope:oomd.SyncResponse)
  })
_sym_db.RegisterMessage(SyncResponse)

ImportRequest = _reflection.GeneratedProtocolMessageType('ImportRequest', (_message.Message,), {
  'DESCRIPTOR' : _IMPORTREQUEST,
  '__module__' : 'oomd_pb2'
  # @@protoc_insertion_point(class_scope:oomd.ImportRequest)
  })
_sym_db.RegisterMessage(ImportRequest)

ImportResponse = _reflection.GeneratedProtocolMessageType('ImportResponse', (_message.Message,), {
  'DESCRIPTOR' : _IMPORTRESPONSE,
  '__module__' : 'oomd_pb2'
  # @@protoc_insertion_point(class_scope:oomd.ImportResponse)
  })
_sym_db.RegisterMessage(ImportResponse)

ImportByFileRequest = _reflection.GeneratedProtocolMessageType('ImportByFileRequest', (_message.Message,), {
  'DESCRIPTOR' : _IMPORTBYFILEREQUEST,
  '__module__' : 'oomd_pb2'
  # @@protoc_insertion_point(class_scope:oomd.ImportByFileRequest)
  })
_sym_db.RegisterMessage(ImportByFileRequest)

EntityRow = _reflection.GeneratedProtocolMessageType('EntityRow', (_message.Message,), {
  'DESCRIPTOR' : _ENTITYROW,
  '__module__' : 'oomd_pb2'
  # @@protoc_insertion_point(class_scope:oomd.EntityRow)
  })
_sym_db.RegisterMessage(EntityRow)

JoinRequest = _reflection.GeneratedProtocolMessageType('JoinRequest', (_message.Message,), {
  'DESCRIPTOR' : _JOINREQUEST,
  '__module__' : 'oomd_pb2'
  # @@protoc_insertion_point(class_scope:oomd.JoinRequest)
  })
_sym_db.RegisterMessage(JoinRequest)

JoinResponse = _reflection.GeneratedProtocolMessageType('JoinResponse', (_message.Message,), {
  'DESCRIPTOR' : _JOINRESPONSE,
  '__module__' : 'oomd_pb2'
  # @@protoc_insertion_point(class_scope:oomd.JoinResponse)
  })
_sym_db.RegisterMessage(JoinResponse)

JoinByFileRequest = _reflection.GeneratedProtocolMessageType('JoinByFileRequest', (_message.Message,), {
  'DESCRIPTOR' : _JOINBYFILEREQUEST,
  '__module__' : 'oomd_pb2'
  # @@protoc_insertion_point(class_scope:oomd.JoinByFileRequest)
  })
_sym_db.RegisterMessage(JoinByFileRequest)

JoinByFileResponse = _reflection.GeneratedProtocolMessageType('JoinByFileResponse', (_message.Message,), {
  'DESCRIPTOR' : _JOINBYFILERESPONSE,
  '__module__' : 'oomd_pb2'
  # @@protoc_insertion_point(class_scope:oomd.JoinByFileResponse)
  })
_sym_db.RegisterMessage(JoinByFileResponse)


DESCRIPTOR._options = None
_FEATUREVALUEMAP_MAPENTRY._options = None
_ONLINEMULTIGETRESPONSE_RESULTENTRY._options = None

_OOMD = _descriptor.ServiceDescriptor(
  name='OomD',
  full_name='oomd.OomD',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=1439,
  serialized_end=1881,
  methods=[
  _descriptor.MethodDescriptor(
    name='OnlineGet',
    full_name='oomd.OomD.OnlineGet',
    index=0,
    containing_service=None,
    input_type=_ONLINEGETREQUEST,
    output_type=_ONLINEGETRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='OnlineMultiGet',
    full_name='oomd.OomD.OnlineMultiGet',
    index=1,
    containing_service=None,
    input_type=_ONLINEMULTIGETREQUEST,
    output_type=_ONLINEMULTIGETRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Sync',
    full_name='oomd.OomD.Sync',
    index=2,
    containing_service=None,
    input_type=_SYNCREQUEST,
    output_type=_SYNCRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Import',
    full_name='oomd.OomD.Import',
    index=3,
    containing_service=None,
    input_type=_IMPORTREQUEST,
    output_type=_IMPORTRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Join',
    full_name='oomd.OomD.Join',
    index=4,
    containing_service=None,
    input_type=_JOINREQUEST,
    output_type=_JOINRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='ImportByFile',
    full_name='oomd.OomD.ImportByFile',
    index=5,
    containing_service=None,
    input_type=_IMPORTBYFILEREQUEST,
    output_type=_IMPORTRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='JoinByFile',
    full_name='oomd.OomD.JoinByFile',
    index=6,
    containing_service=None,
    input_type=_JOINBYFILEREQUEST,
    output_type=_JOINBYFILERESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_OOMD)

DESCRIPTOR.services_by_name['OomD'] = _OOMD

# @@protoc_insertion_point(module_scope)