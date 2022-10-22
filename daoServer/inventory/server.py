from os.path import dirname,abspath
import os
import sys
import grpc
from concurrent import futures
from datetime import date
from signal import SIGINT,SIGTERM,signal
import argparse
import time
import jaeger_client

from loguru import logger
from uuid import uuid4
from opentracing import Tracer
import opentracing
from rocketmq.client import PushConsumer
from jaeger_client.config import Config

from handler.inventory import InventoryServicer, Reback
from proto import inventory_pb2_grpc
from config.config import cfg
from utils.addr import getFreePort

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from common.consul.consul import Consul
from common.grpc_health.v1 import health_pb2_grpc,health
from common.grpc_opentracing import open_tracing_server_interceptor
from common.grpc_opentracing.grpcext import intercept_server

logger.add(f"logs/{date.today().year}-{date.today().month}-{date.today().day}.log", level="DEBUG", rotation="00:00")

service_id = str(uuid4())

# 库存服务
def main():
    tracer = initialize_tracer()
    tracer_intercept = open_tracing_server_interceptor(tracer)
    srv = grpc.server(futures.ThreadPoolExecutor(max_workers=5))
    srv = intercept_server(srv, tracer_intercept)

    # 注册中心健康检查
    s = Consul(cfg["consul"]["host"],cfg["consul"]["port"])
    s.register(cfg["host"],cfg["port"],cfg["name"],service_id,cfg["consul"]["tags"])
    health_pb2_grpc.add_HealthServicer_to_server(health.HealthServicer(), srv)

    # 手动退出服务
    exit(cfg["name"], s)

    # 主业务
    inventory_pb2_grpc.add_InventoryServicer_to_server(InventoryServicer(), srv)
    srv.add_insecure_port(f"{cfg['host']}:{cfg['port']}")
    logger.info(f"{cfg['name']}服务启动: {cfg['host']}:{getFreePort(cfg['port'])}")
    srv.start()

    # 消息消费 - 库存归还
    p = PushConsumer("transConsumerGroup")
    p.set_name_server_address(f"{cfg['rocketmq']['host']}:{cfg['rocketmq']['port']}")
    p.subscribe("transTopic", Reback)
    p.start()

    # srv.wait_for_termination()
    try:
        while True:
            time.sleep(60 * 60 * 24)
    except KeyboardInterrupt as e:
        print(e.args)
        srv.stop(0)
    tracer.close()

    p.shutdown()


def initialize_tracer():
    conf = Config(
        config={
            "sampler":{
                "type":"const",
                "param":1
            },
            "logging":True
        },
        service_name=cfg['name'],
        validate=True
    )
    tracer = conf.initialize_tracer()
    opentracing.set_global_tracer(tracer)
    return tracer


def comment_args():
    parser = argparse.ArgumentParser()
    parser.add_argumnet('--ip',
        nargs='?',
        type= str,
        default='192.168.8.222',
        help='binding ip',
    )
    parser.add_argument('--port',
        nargs='?',
        type=int,
        default=4345,
        help='binding port',
    )

    args = parser.parse_args()
    args.port = getFreePort(args.port)

    return args


# 获取空闲端口
def getFreePort(port=0)->int:
    if port > 0:
        return port
    import socket
    tcp = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    tcp.bind(("", 0))
    _, port = tcp.getsockname()
    tcp.close()

    return port

def exit(name, s):
    def cb(singal, frame):
        s.deregister(service_id)
        print(f"退出服务:{name}")
        # sys.exit(0)
        os._exit(0)

    signal(SIGINT, cb)
    signal(SIGTERM, cb)

if "__main__" == __name__:
    main()