import sys
from os.path import dirname,abspath
from datetime import datetime

import grpc
from loguru import logger
from peewee import *
from google.protobuf.empty_pb2 import Empty

from order.proto.shopcart_pb2_grpc import ShopCartServicer

sys.path.insert(0,dirname(dirname(abspath(dirname(__file__)))))
from model.models import shop_cart as ShopCart
from proto import shopcart_pb2_grpc,shopcart_pb2

today = datetime.today()
logger.add(f"logs/{today.year}-{today.month}-{today.day}.log", rotation="12:00", level="DEBUG")

# import logging
# log = logging.getLogger("peewee")
# log.setLevel(logging.DEBUG)
# log.addHandler(logging.StreamHandler())

class ShopCartServicer(shopcart_pb2_grpc.ShopCart):

    @logger.catch
    def QueryShopCart(self, req: shopcart_pb2.UserInfoRequest, context)->shopcart_pb2.ShopCartListResponse:
        res = shopcart_pb2.ShopCartListResponse()
        carts = ShopCart.select().where(ShopCart.user_id == req.id)
        res.total = carts.count()
        for item in carts:
            r = shopcart_pb2.ShopCartRequest()
            r.user_id = item.user_id
            r.goods_id = item.goods_id
            r.nums = item.nums
            r.checked = item.checked
            res.data.append(r)
        
        return res


    @logger.catch
    def AddGoodsToShopCart(self, req: shopcart_pb2.ShopCartRequest, context)->Empty:
        try:
            rows = (ShopCart.select()
            .where(ShopCart.user_id==req.user_id , ShopCart.goods_id==req.goods_id)
            .orwhere(
                (ShopCart.is_deleted==1) & 
                (ShopCart.user_id==req.user_id) & (ShopCart.goods_id==req.goods_id)
            ))
            if len(rows) > 0:
                cart = rows[0]
                if cart.is_deleted == 1:
                    cart.nums = req.nums    
                    cart.is_deleted = 0
                else:
                    cart.nums += req.nums
            else:
                cart = ShopCart()
                cart.nums = req.nums
                cart.user_id = req.user_id
                cart.goods_id = req.goods_id
            cart.checked = True
            n = cart.save()
            if n <1:
                raise ValueError("没有变化")
        except Exception as e:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(f"添加失败: {e.args}")
        
        return Empty()


    @logger.catch
    def UpdateShopCart(self, req: shopcart_pb2.ShopCartRequest, context)->Empty:
        try:
            rows = ShopCart.select().where(ShopCart.user_id==req.user_id, ShopCart.goods_id==req.goods_id)
            cart = rows[0]
            cart.nums = req.nums
            cart.checked = req.checked
            cart.save()
        except Exception as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在:{e.args}")

        return Empty()


    @logger.catch
    def DelGoodsInShopCart(self, req: shopcart_pb2.ShopCartRequest, context)->Empty:
        try:
            cart = ShopCart.get(ShopCart.user_id==req.user_id, ShopCart.goods_id==req.goods_id)
            cart.delete_instance()
        except ShopCart.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在:{e.args}")
        
        return Empty()