import sys
import os
import signal
import argparse
from concurrent import futures
import socket
from datetime import date
import time
import uuid

import grpc
from grpc_opentracing import open_tracing_server_interceptor
from loguru import logger
from jaeger_client.config import Config

from config import config
from proto import user_pb2_grpc
from handler.user import User

sys.path.insert(0, os.path.dirname(os.path.abspath(os.path.dirname(__file__))))
from common.grpc_health.v1 import health_pb2_grpc,health
from common.consul.consul import Consul
from common.grpc_opentracing import grpcext

logger.add(f"logs/server_{date.today().day}.log", rotation="12:00")


@logger.catch
def server():

    args = command_args()

    # 注册服务to注册中心
    service_id = str(uuid.uuid4())
    c = Consul(config.CONSUL_HOST, config.CONSUL_PORT)
    c.register(config.USER_SERVER_HOST, config.USER_SERVER_PORT, config.USER_SERVER_NAME, service_id, )
    exit(c, service_id)

    serv = grpc.server(futures.ThreadPoolExecutor(max_workers=5))
    health_pb2_grpc.add_HealthServicer_to_server(health.HealthServicer(), serv)
    cfg =  Config(
        config={
            "sampler":{
                "type":"const",
                "param":1
            },
            "logging":False
        },
        service_name=config.cfg['name'],
        validate=True
    )
    tracer = cfg.initialize_tracer()
    intercept = open_tracing_server_interceptor(tracer)
    serv = grpcext.intercept_server(serv,intercept)

    user_pb2_grpc.add_UserServicer_to_server(User(), serv)
    serv.add_insecure_port(f"{args.ip}:{args.port}")
    logger.info(f"{config.USER_SERVER_NAME} 服务启动：{args.ip}:{args.port} ...")
    serv.start()
    # serv.wait_for_termination()

    try:
        while True:
            time.sleep(24*60*60)
    except KeyboardInterrupt:
        serv.stop(0)
    
    tracer.close()
    


# 命命行添加参数
def command_args():
    parser = argparse.ArgumentParser()
    parser.add_argument('--ip',
                        nargs="?",
                        type = str,
                        default = config.USER_SERVER_HOST,
                        help = "binding ip"
    )
    parser.add_argument('--port',
                        nargs="?",
                        type =int,
                        default = config.USER_SERVER_PORT,
                        help = "the listening port"
    )
    args = parser.parse_args()

    if config.ENV != "debug":
        args.port = getFreePort()
    
    return args


# 获取空闲端口
def getFreePort()->int:
    tcp = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    tcp.bind(("", 0))
    _, port = tcp.getsockname()
    tcp.close()

    return port


"""
 优雅退出
    window 下支持的信号是有限的:
        SIGINT  ctrl+c终止
        SIGTERM kill发出的软件终止
"""
def exit(c: Consul, service_id:uuid.UUID):
    def f(signo, frame):
        print("注销注册中心的服务", service_id)
        c.deregister(service_id)
        print("退出服务 ...")
        sys.exit(0)
    signal.signal(signal.SIGINT, f)
    signal.signal(signal.SIGTERM, f)


if "__main__" == __name__:
    server()
