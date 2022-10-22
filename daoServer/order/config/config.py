import sys
from os.path import dirname,abspath
import json
from datetime import datetime

from playhouse.shortcuts import ReconnectMixin
from playhouse.pool import PooledMySQLDatabase

ROOT_DIR = dirname(dirname(abspath(dirname(__file__))))

sys.path.insert(0, ROOT_DIR)
from common.nacos import NacosClient

today = datetime.today()
LOG_PATH = f"{dirname(abspath(dirname(__file__)))}/logs/{today.year}-{today.month}-{today.day}.log"

nacos = {
    "host":"192.168.31.106",
    "username":"nacos",
    "password":"nacos",
    "namespace":"0c544bf3-c29c-43da-abbf-ce332438e1cb",
    "dataid":"order-srv.json",
    "group":"dev",
}

client = NacosClient(nacos["host"],namespace=nacos["namespace"], username=nacos["username"], password=nacos["password"])
# 设置从服务端获取到本地的快照配置
client.set_options(snapshot_base=f"{ROOT_DIR}/common/nacos/nacos-data/snapshot")
content = client.get_config(nacos["dataid"], nacos["group"])
cfg = json.loads(content)



class ReconnectMysqlDatabase(ReconnectMixin, PooledMySQLDatabase):
    pass

db = ReconnectMysqlDatabase(
    cfg["mysql"]["db"],
    host = cfg["mysql"]["host"],
    port = cfg["mysql"]["port"],
    user = cfg["mysql"]["user"],
    password = cfg["mysql"]["password"]
)
