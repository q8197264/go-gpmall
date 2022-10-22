# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2
from . import order_pb2 as order__pb2


class OrderStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.CreateOrder = channel.unary_unary(
                '/Order/CreateOrder',
                request_serializer=order__pb2.OrderRequest.SerializeToString,
                response_deserializer=order__pb2.OrderDetailResponse.FromString,
                )
        self.QueryOrderList = channel.unary_unary(
                '/Order/QueryOrderList',
                request_serializer=order__pb2.OrderRequest.SerializeToString,
                response_deserializer=order__pb2.OrderListResponse.FromString,
                )
        self.QueryOrderDetail = channel.unary_unary(
                '/Order/QueryOrderDetail',
                request_serializer=order__pb2.OrderRequest.SerializeToString,
                response_deserializer=order__pb2.OrderDetailResponse.FromString,
                )
        self.DelOrder = channel.unary_unary(
                '/Order/DelOrder',
                request_serializer=order__pb2.OrderRequest.SerializeToString,
                response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
                )
        self.UpdateOrderStatus = channel.unary_unary(
                '/Order/UpdateOrderStatus',
                request_serializer=order__pb2.OrderStatusRequest.SerializeToString,
                response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
                )


class OrderServicer(object):
    """Missing associated documentation comment in .proto file."""

    def CreateOrder(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def QueryOrderList(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def QueryOrderDetail(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def DelOrder(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def UpdateOrderStatus(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_OrderServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'CreateOrder': grpc.unary_unary_rpc_method_handler(
                    servicer.CreateOrder,
                    request_deserializer=order__pb2.OrderRequest.FromString,
                    response_serializer=order__pb2.OrderDetailResponse.SerializeToString,
            ),
            'QueryOrderList': grpc.unary_unary_rpc_method_handler(
                    servicer.QueryOrderList,
                    request_deserializer=order__pb2.OrderRequest.FromString,
                    response_serializer=order__pb2.OrderListResponse.SerializeToString,
            ),
            'QueryOrderDetail': grpc.unary_unary_rpc_method_handler(
                    servicer.QueryOrderDetail,
                    request_deserializer=order__pb2.OrderRequest.FromString,
                    response_serializer=order__pb2.OrderDetailResponse.SerializeToString,
            ),
            'DelOrder': grpc.unary_unary_rpc_method_handler(
                    servicer.DelOrder,
                    request_deserializer=order__pb2.OrderRequest.FromString,
                    response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
            ),
            'UpdateOrderStatus': grpc.unary_unary_rpc_method_handler(
                    servicer.UpdateOrderStatus,
                    request_deserializer=order__pb2.OrderStatusRequest.FromString,
                    response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'Order', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Order(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def CreateOrder(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/Order/CreateOrder',
            order__pb2.OrderRequest.SerializeToString,
            order__pb2.OrderDetailResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def QueryOrderList(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/Order/QueryOrderList',
            order__pb2.OrderRequest.SerializeToString,
            order__pb2.OrderListResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def QueryOrderDetail(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/Order/QueryOrderDetail',
            order__pb2.OrderRequest.SerializeToString,
            order__pb2.OrderDetailResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def DelOrder(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/Order/DelOrder',
            order__pb2.OrderRequest.SerializeToString,
            google_dot_protobuf_dot_empty__pb2.Empty.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def UpdateOrderStatus(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/Order/UpdateOrderStatus',
            order__pb2.OrderStatusRequest.SerializeToString,
            google_dot_protobuf_dot_empty__pb2.Empty.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
