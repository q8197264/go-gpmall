from os.path import dirname,abspath
from sys import path
import json

from playhouse.shortcuts import ReconnectMixin
from playhouse.pool import PooledMySQLDatabase
ROOT_DIR = dirname(dirname(abspath(dirname(__file__))))
path.insert(0, ROOT_DIR)

from common import nacos

nacosConfig = {
    "host":"127.0.0.1",
    "port":"8848",
    "namespace":"b71d73b4-bd17-4639-9749-767a6363aff5",
    "username":"nacos",
    "password":"nacos",
    "dataid":"inventory-srv.json",
    "group":"dev"
}

client = nacos.NacosClient(
    nacosConfig["host"],
    namespace=nacosConfig["namespace"],
    username=nacosConfig["username"],
    password=nacosConfig["password"]
)

# 设置从服务端获取到本地的快照配置
client.set_options(snapshot_base=f"{ROOT_DIR}/common/nacos/nacos-data/snapshot")
content = client.get_config(nacosConfig["dataid"], nacosConfig["group"])
cfg = json.loads(content)

# def update_config(args):
#     print(args)

# client.add_config_watcher(nacosConfig["dataid"], nacosConfig["group"], update_config)


SRV_NAME = cfg["name"]
SRV_HOST = cfg["host"]
SRV_PORT = cfg["port"]
CONSUL_HOST = cfg["consul"]["host"]
CONSUL_PORT = cfg["consul"]["port"]
CONSUL_TAGS = cfg["consul"]["tags"]
ROCKETMQ_HOST = cfg["rocketmq"]["host"]
ROCKETMQ_PORT = cfg["rocketmq"]["port"]

class ReconnectMysqlDatabase(PooledMySQLDatabase, ReconnectMixin):
    pass

DB = ReconnectMysqlDatabase(
    cfg["mysql"]["db"], 
    host=cfg["mysql"]["host"],
    port=cfg["mysql"]["port"],
    user=cfg["mysql"]["user"],
    password=cfg["mysql"]["password"]
)
