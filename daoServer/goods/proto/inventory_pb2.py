# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: inventory.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='inventory.proto',
  package='',
  syntax='proto3',
  serialized_options=b'Z\007.;proto',
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n\x0finventory.proto\x1a\x1bgoogle/protobuf/empty.proto\",\n\x0cGoodsInvInfo\x12\x0f\n\x07goodsId\x18\x01 \x01(\x05\x12\x0b\n\x03num\x18\x02 \x01(\x05\"9\n\x08SellInfo\x12\x10\n\x08order_sn\x18\x01 \x01(\t\x12\x1b\n\x04\x64\x61ta\x18\x02 \x03(\x0b\x32\r.GoodsInvInfo2\xbf\x01\n\tInventory\x12/\n\x06SetInv\x12\r.GoodsInvInfo\x1a\x16.google.protobuf.Empty\x12+\n\x06Reback\x12\t.SellInfo\x1a\x16.google.protobuf.Empty\x12)\n\x04Sell\x12\t.SellInfo\x1a\x16.google.protobuf.Empty\x12)\n\tInvDetail\x12\r.GoodsInvInfo\x1a\r.GoodsInvInfoB\tZ\x07.;protob\x06proto3'
  ,
  dependencies=[google_dot_protobuf_dot_empty__pb2.DESCRIPTOR,])




_GOODSINVINFO = _descriptor.Descriptor(
  name='GoodsInvInfo',
  full_name='GoodsInvInfo',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='goodsId', full_name='GoodsInvInfo.goodsId', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='num', full_name='GoodsInvInfo.num', index=1,
      number=2, type=5, cpp_type=1, label=1,
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
  serialized_start=48,
  serialized_end=92,
)


_SELLINFO = _descriptor.Descriptor(
  name='SellInfo',
  full_name='SellInfo',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='order_sn', full_name='SellInfo.order_sn', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='data', full_name='SellInfo.data', index=1,
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
  serialized_start=94,
  serialized_end=151,
)

_SELLINFO.fields_by_name['data'].message_type = _GOODSINVINFO
DESCRIPTOR.message_types_by_name['GoodsInvInfo'] = _GOODSINVINFO
DESCRIPTOR.message_types_by_name['SellInfo'] = _SELLINFO
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

GoodsInvInfo = _reflection.GeneratedProtocolMessageType('GoodsInvInfo', (_message.Message,), {
  'DESCRIPTOR' : _GOODSINVINFO,
  '__module__' : 'inventory_pb2'
  # @@protoc_insertion_point(class_scope:GoodsInvInfo)
  })
_sym_db.RegisterMessage(GoodsInvInfo)

SellInfo = _reflection.GeneratedProtocolMessageType('SellInfo', (_message.Message,), {
  'DESCRIPTOR' : _SELLINFO,
  '__module__' : 'inventory_pb2'
  # @@protoc_insertion_point(class_scope:SellInfo)
  })
_sym_db.RegisterMessage(SellInfo)


DESCRIPTOR._options = None

_INVENTORY = _descriptor.ServiceDescriptor(
  name='Inventory',
  full_name='Inventory',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=154,
  serialized_end=345,
  methods=[
  _descriptor.MethodDescriptor(
    name='SetInv',
    full_name='Inventory.SetInv',
    index=0,
    containing_service=None,
    input_type=_GOODSINVINFO,
    output_type=google_dot_protobuf_dot_empty__pb2._EMPTY,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Reback',
    full_name='Inventory.Reback',
    index=1,
    containing_service=None,
    input_type=_SELLINFO,
    output_type=google_dot_protobuf_dot_empty__pb2._EMPTY,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Sell',
    full_name='Inventory.Sell',
    index=2,
    containing_service=None,
    input_type=_SELLINFO,
    output_type=google_dot_protobuf_dot_empty__pb2._EMPTY,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='InvDetail',
    full_name='Inventory.InvDetail',
    index=3,
    containing_service=None,
    input_type=_GOODSINVINFO,
    output_type=_GOODSINVINFO,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_INVENTORY)

DESCRIPTOR.services_by_name['Inventory'] = _INVENTORY

# @@protoc_insertion_point(module_scope)
