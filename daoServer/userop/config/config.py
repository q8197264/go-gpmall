import sys
from os.path import dirname,abspath
import json
from datetime import datetime

from playhouse.shortcuts import ReconnectMixin 
from playhouse.pool import PooledMySQLDatabase

ROOTDIR = dirname(dirname(abspath(dirname(__file__))))
sys.path.insert(0, ROOTDIR)
from common.nacos import client

today = datetime.today()
LOG_PATH = f"{dirname(abspath(dirname(__file__)))}/logs/{today.year}-{today.month}-{today.day}.log"

nacos = {
    "host":"192.168.8.222",
    "port":8848,
    "username":"nacos",
    "password":"nacos",
    "group":"dev",
    "namespace":"c8b504c3-686d-4c9a-9e6d-2e425d15d9f4",
    "dataid":"userop-srv.json",
}

sr = client.NacosClient(f"{nacos['host']}:{nacos['port']}",namespace=nacos["namespace"],username=nacos["username"],password=nacos["password"])
content = sr.get_config(data_id=nacos['dataid'],group=nacos['group'])
cfg = json.loads(content)


class ReconnectMysqlDatabase(ReconnectMixin,PooledMySQLDatabase):
    pass

db = ReconnectMysqlDatabase(
    cfg["mysql"]["db"],
    host=cfg["mysql"]["host"],
    port=cfg["mysql"]["port"],
    user=cfg["mysql"]["user"],
    password=cfg["mysql"]["password"]
)

if "__main__" == __name__:
    print(db)
