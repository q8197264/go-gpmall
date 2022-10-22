import sys
from os.path import dirname,abspath
import json

# playhouse 模块不用单独装，装完 peewee 就有了
from playhouse.shortcuts import ReconnectMixin
from playhouse.pool import PooledMySQLDatabase

ROOT_DIR = dirname(dirname(abspath(dirname(__file__))))
sys.path.insert(0, ROOT_DIR)
from common import nacos

nacosConfig = {
    "host": "127.0.0.1:8848",
    "port": 8848,
    "namespace": "bfc201da-e8ad-41b3-a7f5-4e06bf616d47",
    "group": "dev",
    "data_id": "goods-srv.json",
    "username": "nacos",
    "password": "nacos"
}

client = nacos.NacosClient(
    nacosConfig["host"], 
    namespace=nacosConfig["namespace"],
    username=nacosConfig["username"],
    password=nacosConfig["password"]
)
# 设置从服务端获取到本地的快照配置
client.set_options(snapshot_base=f"{ROOT_DIR}/common/nacos/nacos-data/snapshot")
cfg = client.get_config(nacosConfig["data_id"], nacosConfig["group"])
cfg = json.loads(cfg)

def update_cfg(args):
    print("配置产生变化")
    print(args)

# client.add_config_watcher(nacosConfig["data_id"], nacosConfig["group"], update_cfg)
# print(config.client)

SRV_NAME = cfg["name"]
SRV_HOST = cfg["host"]
SRV_PORT = cfg["port"]
CONSUL_HOST = cfg["consul"]["host"]
CONSUL_PORT = cfg["consul"]["port"]
CONSUL_TAGS = cfg["consul"]["tags"]

class ReconnectMysqlDatabase(ReconnectMixin, PooledMySQLDatabase):
    pass

DB = ReconnectMysqlDatabase (
    cfg["mysql"]["db"], 
    host=cfg["mysql"]["host"],
    port=cfg["mysql"]["port"],
    user=cfg["mysql"]["user"],
    password=cfg["mysql"]["password"]
)

if __name__ == "__main__":
    print(cfg)
