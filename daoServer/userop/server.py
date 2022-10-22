import sys
from uuid import uuid4
import grpc
from concurrent import futures
from signal import signal,SIGINT,SIGTERM
import argparse

from loguru import logger
from peewee import *
import uuid

from handler import fav,address,post
from config.config import cfg,LOG_PATH,ROOTDIR
from proto import favorites_pb2_grpc,address_pb2_grpc,post_pb2_grpc

sys.path.insert(0, ROOTDIR)
from common.consul.consul import Consul
from common.grpc_health.v1 import health,health_pb2_grpc

logger.add(LOG_PATH, rotation="0:0", level="DEBUG")

def main():
    service_id = str(uuid.uuid4())
    # 
    ip, port = args_comment()

    # consul register service
    sc = Consul(cfg["consul"]["host"], cfg["consul"]["port"])
    if sc.register(
        host=ip,
        port=port, 
        name=cfg["name"], 
        tags=cfg["consul"]["tags"],
        service_id=service_id
    ):
        print(f"{cfg['name']}服务注册成功")
    else:
        print(f"{cfg['name']}服务注册失败")

    # 优雅退出
    wait(sc, service_id)

    serv = grpc.server(futures.ThreadPoolExecutor(max_workers=5))
    health_pb2_grpc.add_HealthServicer_to_server(health.HealthServicer(), serv)
    favorites_pb2_grpc.add_FavoritesServicer_to_server(fav.FavServicer(), serv)
    address_pb2_grpc.add_AddressServicer_to_server(address.AddressServicer(), serv)
    post_pb2_grpc.add_PostServicer_to_server(post.PostServicer(), serv)
    serv.add_insecure_port(f"{ip}:{port}")

    print(f"开启服务...{ip}:{port}")
    serv.start()
    serv.wait_for_termination()
    
@logger.catch
def wait(sc, service_id):
    def f(signalnum, frame):
        if sc.deregister(service_id):
            print("服务注销成功")
        else:
            print("服务注销失败")
        sys.exit(0)

    signal(SIGINT, f)
    signal(SIGTERM, f)

def args_comment():
    parse = argparse.ArgumentParser()
    parse.add_argument(
        "--ip",
        nargs="?",
        type=str,
        default=cfg["host"],
        help="binding ip"
    )
    parse.add_argument(
        "--port",
        nargs="?",
        type=int,
        default=cfg["port"],
        help="binding port"
    )
    args = parse.parse_args()

    return args.ip,args.port


if __name__ == "__main__":
    main()