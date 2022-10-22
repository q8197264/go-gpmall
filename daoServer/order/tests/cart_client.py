import sys
from os.path import dirname,abspath
from datetime import datetime

import grpc
from peewee import Model
from loguru import logger

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from proto import shopcart_pb2_grpc,order_pb2_grpc,shopcart_pb2,order_pb2
from config.config import cfg

today = datetime.today()
logger.add(f"logs/{today.year}-{today.month}-{today.day}.log", rotation="12:00", level="DEBUG")

class client(Model):
    def __init__(self):
        ch = grpc.insecure_channel(f"{cfg['host']}:{cfg['port']}")
        self.client = shopcart_pb2_grpc.ShopCartStub(ch)

    def QueryShopCart(self, uid):
        res = self.client.QueryShopCart(shopcart_pb2.UserInfoRequest(id=uid))
        print(res.total)
        for item in res.data:
            print(item)
    
    def AddShopCart(self, uid, gid, nums):
        self.client.AddGoodsToShopCart(shopcart_pb2.ShopCartRequest(
            user_id = uid,
            goods_id = gid,
            nums = nums,
            checked = True
        ))

    def UpdateShopCart(self, uid, gid, nums, checked=True):
        self.client.UpdateShopCart(shopcart_pb2.ShopCartRequest(
            user_id = uid,
            goods_id = gid,
            nums = nums,
            checked = checked
        ))
    
    def DelGoodsFromShopCart(self, uid, gid):
        self.client.DelGoodsInShopCart(shopcart_pb2.ShopCartRequest(
            user_id = uid,
            goods_id = gid
        ))

if "__main__" == __name__:
    c = client()
    c.AddShopCart(1, 14, 1)
    # c.UpdateShopCart(1, 1, 2, False)
    # c.DelGoodsFromShopCart(1, 3)
    # c.QueryShopCart(1)