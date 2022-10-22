from concurrent import futures
import sys
import os
from signal import signal, SIGINT, SIGTERM
from datetime import date
from uuid import uuid4
import argparse
import time

import grpc
from rocketmq.client import PushConsumer
from jaeger_client import Config
from loguru import logger

from config import ROOT_DIR,cfg
from proto import order_pb2_grpc,shopcart_pb2_grpc
from handler import order,cart
from utils.addr import get_free_port

sys.path.insert(0, ROOT_DIR)
from common.consul import consul
from common.grpc_health.v1 import health_pb2_grpc,health
from common.grpc_opentracing import open_tracing_server_interceptor
from common.grpc_opentracing.grpcext import intercept_server

today = date.today()
logger.add(f"logs/{today.year}-{today.month}-{today.day}.log", rotation="12:00", level="DEBUG")

def run():
    # cfg["host"],cfg["port"] = args_comment() #不能跨包
    
    # 注册中心 consul
    service_id= str(uuid4())
    client = consul.Consul(cfg["consul"]["host"], cfg["consul"]["port"])
    client.register(cfg["host"], cfg["port"], cfg["name"], tags=cfg["consul"]["tags"], service_id=service_id)

    # 退出
    exit(client, service_id)

    """
        把 tracer 放 config 共享
    """
    # jaeger 链路追踪
    # import logging
    # log_level = logging.DEBUG
    # logging.getLogger('').handlers = []
    # logging.basicConfig(format='%(asctime)s %(message)s', level=log_level)
    config = Config(
        config={ # usually read from some yaml config
            'sampler': {
                'type': 'const',
                'param': 1,
            },
            'logging': False,
        },
        service_name=cfg["name"],
        validate=True,
    )
    tracer = config.initialize_tracer()
    tracer_interceptor = open_tracing_server_interceptor(tracer)

    srv = grpc.server(futures.ThreadPoolExecutor(max_workers=5))
    health_pb2_grpc.add_HealthServicer_to_server(health.HealthServicer(), srv)

    srv = intercept_server(srv, tracer_interceptor)
    shopcart_pb2_grpc.add_ShopCartServicer_to_server(cart.ShopCartServicer(), srv)
    order_pb2_grpc.add_OrderServicer_to_server(order.OrderServicer(), srv)
    srv.add_insecure_port("%s:%d" % (cfg["host"], cfg["port"]))
    logger.info(f"service start...: {cfg['host']}:{cfg['port']}")
    srv.start()

    # 取消超时订单
    c = PushConsumer("delayOrderConsumerGroup")
    c.set_name_server_address(f"{cfg['rocketmq']['host']}:{cfg['rocketmq']['port']}")
    c.subscribe("delayOrder", order.cancel_timeout_order)
    c.start()

    # srv.wait_for_termination()

    try:
        while True:
            time.sleep(60 * 60 * 24)
    except KeyboardInterrupt:
        srv.stop(0)
    tracer.close()

    c.shutdown()

def args_comment():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--ip",
        nargs="?",
        type=str,
        default=cfg["host"],
        help="binding ip"
    )
    parser.add_argument(
        "--port",
        nargs="?",
        type=int,
        default=cfg["port"],
        help="binding port"
    )
    args = parser.parse_args()
    args.port = get_free_port(0)

    return args.ip, args.port


def exit(client: consul.Consul, service_id):
    def f(signalnum, frame):
        client.deregister(service_id)
        logger.info(f"退出订单服务 service stop...:{cfg['host']}:{cfg['port']}")
        # sys.exit(0)
        os._exit(0)

    signal(SIGINT, f)
    signal(SIGTERM, f)

if "__main__" == __name__:
    run()