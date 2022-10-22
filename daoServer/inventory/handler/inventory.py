from os.path import dirname,abspath
from sys import path
from datetime import date
import json

import grpc
from loguru import logger
from rocketmq.client import ConsumeStatus
from model.models import inventory as Inventory, inventory_history as InventoryHistory
from google.protobuf import empty_pb2

from config import DB

path.insert(0, dirname(abspath(dirname(__file__))))
from proto import inventory_pb2,inventory_pb2_grpc
from model.models import *

from common.lock.redis_lock import RLock

logger.add(f"logs/{date.today().year}-{date.today().month}-{date.today().day}.log", rotation="00:00",level="DEBUG")

@logger.catch
class InventoryServicer(inventory_pb2_grpc.Inventory):

    @logger.catch
    def SetInv(self, req: inventory_pb2.GoodsInvInfo, context)->empty_pb2.Empty:
        try:
            inventory = Inventory.get(req.goodsId)
            inventory.stocks = req.num
            inventory.save()
        except Inventory.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("库存不存在")
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"设置库存失败：{str(e)}")
        
        return empty_pb2.Empty()


    @logger.catch
    def InvDetail(self, req: inventory_pb2.GoodsInvInfo, context)->inventory_pb2.GoodsInvInfo:
        res = inventory_pb2.GoodsInvInfo()
        try:
            inventory = Inventory.get(req.goodsId)
            res.goodsId = inventory.goods_id
            res.num = inventory.stocks
        except Inventory.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"数据不存在")
        
        return res


    @logger.catch
    def BatchInvDetail(seslf, req: inventory_pb2.GoodsInvInfo, context):
        pass


    @logger.catch
    def Sell(self, req: inventory_pb2.SellInfo, context)->empty_pb2.Empty:
        """
        销售扣减库存
        Mysql乐观锁实验
        """
        with DB.atomic() as txn:
            fields = []
            for item in req.data:
                try:
                    while True:#模拟乐观锁: 无限读 + 写锁
                        inventory = Inventory.get(item.goodsId)

                        # TODO: 这里判断可能有超卖现象
                        if inventory.stocks >= item.num:
                            # TODO:mysql乐观锁解决 [用版本号实现并发锁机制]
                            ok = Inventory.update(
                                stocks=Inventory.stocks-item.num, 
                                version=Inventory.version+1
                            ).where(
                                Inventory.goods_id==item.goodsId, 
                                Inventory.version==inventory.version
                            ).execute()
                            fields.append({
                                "goods_id": item.goodsId,
                                "nums":item.num
                            })
                            # if ok:
                            #     print("更新成功",item.goodsId)
                            # else:
                            #     print("更新失败", item.goodsId)
                            break
                        else:
                            context.set_code(grpc.StatusCode.RESOURCE_EXHAUSTED)
                            context.set_details(f"库存不足:-{inventory.stocks}")
                            txn.rollback()
                            break
                except Inventory.DoesNotExist as e:
                    context.set_code(grpc.StatusCode.NOT_FOUND)
                    context.set_details(f"记录不存在")
                    txn.rollback()
                    break

            # 添加库存记录
            try:
                id = InventoryHistory.insert({
                    "order_sn":req.order_sn,
                    "order_inv_detail": json.dumps(fields),
                    "status":1
                }).execute()
                if id < 1:
                    raise ValueError(f"插入库存记录失败: {id}")
            except ValueError as e:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details(f"{e.args}")
                txn.rollback()

        return empty_pb2.Empty()


    # 无用代码, 可删除
    @logger.catch
    def Reback_test(self, req:inventory_pb2.SellInfo, context)->empty_pb2.Empty:
        """
         库存归还(加库存)超量
         Redis 分布式锁实验(推荐)
        """
        import threading
        import uuid

        id = "%s:%d" % (str(uuid.uuid4()),threading.currentThread().ident)
        rlock = RLock(id=id, expire=5, auto_renewal=True)
        # for goods in req.goodsInfo:
        #     rlock._test(True, goods.goodsId)

        rlock.acquire()
        with DB.atomic() as txn:
            for goods in req.data:
                try:
                    inventory = Inventory.get(Inventory.goods_id==goods.goodsId)
                except Inventory.DoesNotExist as e:
                    context.set_code(grpc.StatusCode.NOT_FOUND)
                    context.set_details(f"记录不存在")
                    txn.rollback()
                    rlock.release()
                    return empty_pb2.Empty()

                # 模拟延时 begin
                # import time
                # from random import randint
                # time.sleep(randint(0,3))
                # print("更新：", inventory.goods_id)
                # end

                #TODO: 这里可能引起不一致 分布式锁
                inventory.stocks += goods.num
                inventory.save()
        rlock.release()

        return empty_pb2.Empty()


@logger.catch
def Reback(msg):
    msg = json.loads(msg.body.decode("utf-8"))
    orderSn = msg["order_sn"]
    """
        库存归还(加库存)超量
        Redis 分布式锁实验(推荐)
    """
    import threading
    import uuid

    # 分布式锁
    id = "%s:%d" % (str(uuid.uuid4()),threading.currentThread().ident)
    rlock = RLock(id=id, expire=5, auto_renewal=True)
    # for goods in req.goodsInfo:
    #     rlock._test(True, goods.goodsId)

    rlock.acquire()
    with DB.atomic() as txn:
        try:
            inv_history = InventoryHistory.get(InventoryHistory.order_sn == orderSn, InventoryHistory.status == 1)
            goods_list = json.loads(inv_history.order_inv_detail)
            for item in goods_list:
                # 这里可能超卖 所用以分布式锁
                Inventory.update(stocks = inventory.stocks+item["nums"]).where(Inventory.goods_id==item["goods_id"]).execute()
                
                # 模拟延时 begin
                # import time
                # from random import randint
                # time.sleep(randint(0,3))
                # print("更新：", inventory.goods_id)
                # end

            inv_history.status = 2
            inv_history.save()
        except DoesNotExist as e:
            txn.rollback()
            logger.info(f"库存记录不存在: {e.args}")
        except Exception as e:
            txn.rollback()
            rlock.release()
            logger.info(f"库存归还出现异常: {e.args}")
            return ConsumeStatus.RECONSUME_LATER

    rlock.release()

    return ConsumeStatus.CONSUME_SUCCESS
