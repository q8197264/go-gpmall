# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: address.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='address.proto',
  package='',
  syntax='proto3',
  serialized_options=b'Z\007.;proto',
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n\raddress.proto\x1a\x1bgoogle/protobuf/empty.proto\"\xb0\x01\n\x0e\x41\x64\x64ressRequest\x12\n\n\x02id\x18\x01 \x01(\x05\x12\x0f\n\x07user_id\x18\x02 \x01(\x05\x12\x10\n\x08province\x18\x03 \x01(\t\x12\x0c\n\x04\x63ity\x18\x04 \x01(\t\x12\x10\n\x08\x64istrict\x18\x05 \x01(\t\x12\x0f\n\x07\x61\x64\x64ress\x18\x06 \x01(\t\x12\x13\n\x0bsigner_name\x18\x07 \x01(\t\x12\x15\n\rsigner_mobile\x18\x08 \x01(\t\x12\x12\n\nis_default\x18\t \x01(\x08\"7\n\x15\x41\x64\x64ressStatusResponse\x12\n\n\x02id\x18\x01 \x01(\x05\x12\x12\n\nis_default\x18\x02 \x01(\x08\"\xa6\x01\n\x15\x41\x64\x64ressDetailResponse\x12\n\n\x02id\x18\x01 \x01(\x05\x12\x10\n\x08province\x18\x02 \x01(\t\x12\x0c\n\x04\x63ity\x18\x03 \x01(\t\x12\x10\n\x08\x64istrict\x18\x04 \x01(\t\x12\x0f\n\x07\x61\x64\x64ress\x18\x05 \x01(\t\x12\x13\n\x0bsigner_name\x18\x06 \x01(\t\x12\x15\n\rsigner_mobile\x18\x07 \x01(\t\x12\x12\n\nis_default\x18\x08 \x01(\x08\"J\n\x13\x41\x64\x64ressListResponse\x12\r\n\x05total\x18\x01 \x01(\x05\x12$\n\x04\x64\x61ta\x18\x02 \x03(\x0b\x32\x16.AddressDetailResponse2\xef\x01\n\x07\x41\x64\x64ress\x12\x39\n\x10QueryAddressList\x12\x0f.AddressRequest\x1a\x14.AddressListResponse\x12\x35\n\nAddAddress\x12\x0f.AddressRequest\x1a\x16.google.protobuf.Empty\x12\x38\n\rUpdateAddress\x12\x0f.AddressRequest\x1a\x16.google.protobuf.Empty\x12\x38\n\rDeleteAddress\x12\x0f.AddressRequest\x1a\x16.google.protobuf.EmptyB\tZ\x07.;protob\x06proto3'
  ,
  dependencies=[google_dot_protobuf_dot_empty__pb2.DESCRIPTOR,])




_ADDRESSREQUEST = _descriptor.Descriptor(
  name='AddressRequest',
  full_name='AddressRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='AddressRequest.id', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='user_id', full_name='AddressRequest.user_id', index=1,
      number=2, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='province', full_name='AddressRequest.province', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='city', full_name='AddressRequest.city', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='district', full_name='AddressRequest.district', index=4,
      number=5, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='address', full_name='AddressRequest.address', index=5,
      number=6, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='signer_name', full_name='AddressRequest.signer_name', index=6,
      number=7, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='signer_mobile', full_name='AddressRequest.signer_mobile', index=7,
      number=8, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='is_default', full_name='AddressRequest.is_default', index=8,
      number=9, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
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
  serialized_start=47,
  serialized_end=223,
)


_ADDRESSSTATUSRESPONSE = _descriptor.Descriptor(
  name='AddressStatusResponse',
  full_name='AddressStatusResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='AddressStatusResponse.id', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='is_default', full_name='AddressStatusResponse.is_default', index=1,
      number=2, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
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
  serialized_start=225,
  serialized_end=280,
)


_ADDRESSDETAILRESPONSE = _descriptor.Descriptor(
  name='AddressDetailResponse',
  full_name='AddressDetailResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='AddressDetailResponse.id', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='province', full_name='AddressDetailResponse.province', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='city', full_name='AddressDetailResponse.city', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='district', full_name='AddressDetailResponse.district', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='address', full_name='AddressDetailResponse.address', index=4,
      number=5, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='signer_name', full_name='AddressDetailResponse.signer_name', index=5,
      number=6, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='signer_mobile', full_name='AddressDetailResponse.signer_mobile', index=6,
      number=7, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='is_default', full_name='AddressDetailResponse.is_default', index=7,
      number=8, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
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
  serialized_start=283,
  serialized_end=449,
)


_ADDRESSLISTRESPONSE = _descriptor.Descriptor(
  name='AddressListResponse',
  full_name='AddressListResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='total', full_name='AddressListResponse.total', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='data', full_name='AddressListResponse.data', index=1,
      number=2, type=11, cpp_type=10, label=3,
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
  serialized_start=451,
  serialized_end=525,
)

_ADDRESSLISTRESPONSE.fields_by_name['data'].message_type = _ADDRESSDETAILRESPONSE
DESCRIPTOR.message_types_by_name['AddressRequest'] = _ADDRESSREQUEST
DESCRIPTOR.message_types_by_name['AddressStatusResponse'] = _ADDRESSSTATUSRESPONSE
DESCRIPTOR.message_types_by_name['AddressDetailResponse'] = _ADDRESSDETAILRESPONSE
DESCRIPTOR.message_types_by_name['AddressListResponse'] = _ADDRESSLISTRESPONSE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

AddressRequest = _reflection.GeneratedProtocolMessageType('AddressRequest', (_message.Message,), {
  'DESCRIPTOR' : _ADDRESSREQUEST,
  '__module__' : 'address_pb2'
  # @@protoc_insertion_point(class_scope:AddressRequest)
  })
_sym_db.RegisterMessage(AddressRequest)

AddressStatusResponse = _reflection.GeneratedProtocolMessageType('AddressStatusResponse', (_message.Message,), {
  'DESCRIPTOR' : _ADDRESSSTATUSRESPONSE,
  '__module__' : 'address_pb2'
  # @@protoc_insertion_point(class_scope:AddressStatusResponse)
  })
_sym_db.RegisterMessage(AddressStatusResponse)

AddressDetailResponse = _reflection.GeneratedProtocolMessageType('AddressDetailResponse', (_message.Message,), {
  'DESCRIPTOR' : _ADDRESSDETAILRESPONSE,
  '__module__' : 'address_pb2'
  # @@protoc_insertion_point(class_scope:AddressDetailResponse)
  })
_sym_db.RegisterMessage(AddressDetailResponse)

AddressListResponse = _reflection.GeneratedProtocolMessageType('AddressListResponse', (_message.Message,), {
  'DESCRIPTOR' : _ADDRESSLISTRESPONSE,
  '__module__' : 'address_pb2'
  # @@protoc_insertion_point(class_scope:AddressListResponse)
  })
_sym_db.RegisterMessage(AddressListResponse)


DESCRIPTOR._options = None

_ADDRESS = _descriptor.ServiceDescriptor(
  name='Address',
  full_name='Address',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=528,
  serialized_end=767,
  methods=[
  _descriptor.MethodDescriptor(
    name='QueryAddressList',
    full_name='Address.QueryAddressList',
    index=0,
    containing_service=None,
    input_type=_ADDRESSREQUEST,
    output_type=_ADDRESSLISTRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='AddAddress',
    full_name='Address.AddAddress',
    index=1,
    containing_service=None,
    input_type=_ADDRESSREQUEST,
    output_type=google_dot_protobuf_dot_empty__pb2._EMPTY,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='UpdateAddress',
    full_name='Address.UpdateAddress',
    index=2,
    containing_service=None,
    input_type=_ADDRESSREQUEST,
    output_type=google_dot_protobuf_dot_empty__pb2._EMPTY,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='DeleteAddress',
    full_name='Address.DeleteAddress',
    index=3,
    containing_service=None,
    input_type=_ADDRESSREQUEST,
    output_type=google_dot_protobuf_dot_empty__pb2._EMPTY,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_ADDRESS)

DESCRIPTOR.services_by_name['Address'] = _ADDRESS

# @@protoc_insertion_point(module_scope)