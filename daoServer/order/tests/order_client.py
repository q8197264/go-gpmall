import sys
import os
from os.path import dirname,abspath
from async_timeout import timeout
import grpc
from signal import signal, SIGINT, SIGTERM

from jaeger_client import Config

sys.path.insert(0, dirname(dirname(abspath(dirname(__file__)))))
from order.proto import order_pb2,order_pb2_grpc
from order.config import cfg,ROOT_DIR

from common.grpc_interceptor.retry import RetryInterceptor

from common.grpc_opentracing import open_tracing_client_interceptor
from common.grpc_opentracing.grpcext import intercept_channel

import logging

from loguru import logger
file_handler = logging.handlers.RotatingFileHandler(f"{ROOT_DIR}/order/logs/1.log", encoding="utf-8")
logger.add(file_handler, level=0)

class client():
    def __init__(self):
        self.exit()

        log_level = logging.DEBUG
        logging.getLogger('').handlers = []
        logging.basicConfig(format='>> %(asctime)s %(message)s', level=log_level)
        config = Config(
            config={ # usually read from some yaml config
                'sampler': {
                    'type': 'const',
                    'param': 1,
                },
                'logging': False,
            },
            service_name='order-client',
            validate=True,
        )
        self.tracer = config.initialize_tracer()
        tracer_interceptor  = open_tracing_client_interceptor(self.tracer)
       
        #  三种情况下重试
        retry_codes = [grpc.StatusCode.UNKNOWN, grpc.StatusCode.DEADLINE_EXCEEDED, grpc.StatusCode.UNAVAILABLE]
        ch = grpc.insecure_channel(f"{cfg['host']}:{cfg['port']}")
        ch = grpc.intercept_channel(ch, RetryInterceptor(max_retries=3,retry_codes=retry_codes))

        ch = intercept_channel(ch, tracer_interceptor)
        self.client = order_pb2_grpc.OrderStub(ch)

        # tracer.close()

    def Close(self):
        import time
        time.sleep(2)
        self.tracer.close()

    def CreateOrder(self, user_id):
        rsp = self.client.CreateOrder(order_pb2.OrderRequest(
            user_id=user_id,
            name = "段超123",
            mobile = "18612917508",
            address = "华中1code",
            post = "",
        ), timeout=3)
        print(rsp)

    @logger.catch
    def QueryOrderList(self):
        res = self.client.QueryOrderList(order_pb2.OrderRequest(
            userId=1
        ), timeout=1)
        print(res)

    @logger.catch
    def QueryOrderDetail(self, order_id):
        res = c.client.QueryOrderDetail(order_pb2.OrderRequest(
            id=order_id
        ), timeout=3)
        print(res)

        # import time
        # time.sleep(2)
        # self.tracer.close()


    def UpdateOrderStatus(self, order_sn, status):
        c.client.UpdateOrderStatus(order_pb2.OrderStatusRequest(
            orderSn=order_sn,
            status=status
        ))
    
    def DelOrder(self, order_id):
        self.client.DelOrder(order_pb2.OrderRequest(
            id=order_id
        ))

    def exit(self):
        def f(signalnum, frame):
            print("退出订单测试 client ...")
            os._exit(0)
        signal(SIGTERM, f)
        signal(SIGINT, f)



if "__main__" == __name__:
    c = client()
    c.CreateOrder(1)
    # c.DelOrder(3)

    # c.QueryOrderList()
    # c.QueryOrderDetail(4)

    # c.UpdateOrderStatus("202204151713271883", "WAIT_BUYER_PAY")
    c.Close()