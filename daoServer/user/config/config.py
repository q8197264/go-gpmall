import sys
from os import path
import json

from playhouse.pool import PooledMySQLDatabase
from playhouse.shortcuts import ReconnectMixin

ROOT_DIR = path.dirname(path.dirname(path.abspath(path.dirname(__file__))))
sys.path.insert(0, ROOT_DIR)
from common import nacos

nacosConfig = {
    "host": "192.168.31.106:8848",
    "port": 8848,
    "namespace": "663ec12a-3811-4949-aa06-f8a7b92ba9be",
    "group": "dev",
    "data_id": "user-srv.json",
    "username": "nacos",
    "password": "nacos"
}
client = nacos.NacosClient(
        nacosConfig["host"], 
        namespace=nacosConfig["namespace"],
        username=nacosConfig["username"],
        password=nacosConfig["password"]
    )
client.set_options(snapshot_base=f"{ROOT_DIR}/common/nacos/nacos-data/snapshot")
cfg = client.get_config(nacosConfig["data_id"], nacosConfig["group"])
cfg = json.loads(cfg)

ENV = "debug"

CONSUL_HOST = cfg['consul']['host']
CONSUL_PORT = cfg['consul']['port']
USER_SERVER_NAME = cfg['name']
USER_SERVER_HOST = cfg['host']
USER_SERVER_PORT = cfg['port']

class ReconnectMysqlDatabase(ReconnectMixin, PooledMySQLDatabase):
    pass

DB = ReconnectMysqlDatabase(
        cfg['mysql']['db'],
        host=cfg['mysql']['host'],
        port=cfg['mysql']['port'],
        user=cfg['mysql']['user'],
        password=cfg['mysql']['password'],
    )

LOG_PATH = path.join(path.dirname(path.dirname(__file__)),"user","logs")
