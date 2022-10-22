from os.path import dirname,abspath
import sys
import grpc
import threading

from peewee import Model

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from model.models import *
from proto import inventory_pb2,inventory_pb2_grpc
from config import cfg
from handler.inventory import Reback
class client(Model):
    def __init__(self):
        self.chan = grpc.insecure_channel(f"{cfg['host']}:{cfg['port']}")
        self.client = inventory_pb2_grpc.InventoryStub(self.chan)

    def getInv(self, goods_id)->inventory_pb2.GoodsInvInfo:
        req = inventory_pb2.GoodsInvInfo(
            goodsId=goods_id
        )
        res = self.client.InvDetail(req)
        print(res)
        return res

    def setInv(self, goods_id, num):
        req = inventory_pb2.GoodsInvInfo(
            goodsId=goods_id,
            num=num
        )
        self.client.SetInv(req)

    def sell(self, goods_list):
        g = inventory_pb2.SellInfo()
        for gid,stocks in goods_list:
            row = inventory_pb2.GoodsInvInfo()
            row.goodsId = gid
            row.num = stocks
            g.data.append(row)
        self.client.Sell(g)

    def reback(self, goodlist, n):
        g = inventory_pb2.SellInfo()
        for gid,stocks in goodlist:
            row = inventory_pb2.GoodsInvInfo()
            row.goodsId = gid
            row.num = stocks
            g.data.append(row)
        print("测试 thread:",n)
        self.client.Reback(g)


    # 添加测试数据
    def add_data(self):
        rows = [(11,100),(12,100),(13,100),(14,100),(15,100),(16,100)]
        for gid,stocks in rows:
            inv = Inventory(goods_id=gid, stocks=stocks)
            inv.save()


# def reback():


if __name__ == "__main__":
    c = client()
    # c.add_data()
    # c.setInv(1,2)
    # c.getInv(1)
    # c.sell([(1,2),(3,2)])
    Reback({"body":{"order_sn": "202208160500211759"}})
   
    # t1 = threading.Thread(target=c.sell([(11,2),(12,2),(13,4),(14,1),(15,8),(16,100)]))
    # t2 = threading.Thread(target=c.sell([(11,2),(12,2),(13,4),(14,1),(15,8),(16,100)]))
   
    # args = [(11,10),(12,20),(13,40),(14,10),(15,80),(16,10)]
    # t = threading.Thread(target=c.reback, args=(args,1))
    # t1 = threading.Thread(target=c.reback, args=(args,2))
    # t2 = threading.Thread(target=c.reback, args=(args,3))
    # t.start()
    # t1.start()
    # t2.start()
    # t.join()
    # t1.join()
    # t2.join()
    