import sys
from signal import SIGINT, SIGTERM, signal
from os.path import dirname,abspath
from concurrent import futures 
from datetime import date
import time

import uuid
import grpc
import jaeger_client
from loguru import logger
import opentracing
logger.add(f"logs/server_{date.today().day}.log",  rotation="00:00")

from config import CONSUL_HOST,CONSUL_PORT,SRV_NAME,SRV_HOST,SRV_PORT
from proto import goods_pb2_grpc
from handler.goods import GoodsServicer
from utils.addr import getFreePort

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from common.consul.consul import Consul
from common.grpc_health.v1 import health_pb2_grpc,health
from common.grpc_opentracing  import grpcext, open_tracing_server_interceptor

@logger.catch
def main():
    service_id = str(uuid.uuid4())
    cs = Consul(CONSUL_HOST, CONSUL_PORT)
    
    # listen server
    automicRegister(cs, service_id)

    # quit
    exit(cs, service_id)

    cfg = jaeger_client.Config(
        config={
            "sampler":{
                "type":"const",
                "Param":1
            },
            "logging":False
        },
        service_name = SRV_NAME, 
        validate=True
    )
    tracer = cfg.initialize_tracer()
    opentracing.set_global_tracer(tracer)
    intercept = open_tracing_server_interceptor(tracer)

    # 创建服务
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=5))
    server = grpcext.intercept_server(server, intercept)
    goods_pb2_grpc.add_GoodsServicer_to_server(GoodsServicer(), server)
    health_pb2_grpc.add_HealthServicer_to_server(health.HealthServicer(), server)
    server.add_insecure_port(f"{SRV_HOST}:{getFreePort()}")
    logger.info(f"{SRV_NAME}服务启动: {SRV_HOST}:{getFreePort()}")
    server.start()
    # server.wait_for_termination()

    try:
        while True:
            time.sleep(60 * 60 * 24)
    except KeyboardInterrupt:
        server.stop(0)
    tracer.close()


@logger.catch
def automicRegister(cs: Consul, service_id: str):
    # 注册商品服务
    if not cs.register(SRV_HOST,SRV_PORT, SRV_NAME, service_id, ["gpmall","goods","dao"]):
        return


@logger.catch
def exit(cs: Consul, service_id: str):
    # 退出主程序
    def f(singal, frame):
        cs.deregister(service_id)
        logger.info(f"退出 goods-dao 服务 {SRV_HOST}")
        sys.exit(0)
    
    signal(SIGINT, f)
    signal(SIGTERM, f)


if __name__=="__main__":
    main()