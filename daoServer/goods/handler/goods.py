from sys import path
from datetime import date
import json

import grpc
import opentracing
from peewee import *
from loguru import logger
logger.add(f"logs/goods_{date.today().day}.log",  rotation="00:00", level="DEBUG",retention=1)

from google.protobuf import empty_pb2
from config import cfg,ROOT_DIR,SRV_NAME,SRV_HOST,SRV_PORT,CONSUL_HOST,CONSUL_PORT,CONSUL_TAGS
from proto import goods_pb2,goods_pb2_grpc,inventory_pb2,inventory_pb2_grpc
from model.models import *
from handler.category import CategoryServicer

path.insert(0, ROOT_DIR)
from common.consul.consul import Consul
from common.grpc_opentracing import open_tracing_client_interceptor,ActiveSpanSource
from common.grpc_opentracing.grpcext import intercept_channel

# import logging
# log = logging.getLogger("peewee")
# log.setLevel(logging.DEBUG)
# log.addHandler(logging.StreamHandler())

class SetActiveSpanSource(ActiveSpanSource):

    def __init__(self, active_span):
        self.active_span = active_span
    
    def get_active_span(self):
        return self.active_span

class GoodsServicer(CategoryServicer):

    @logger.catch
    def _get_proto_stub(self, span):
        stub = None
        try:
            c = Consul(CONSUL_HOST, CONSUL_PORT)
            target = cfg['inventory']['name']
            host,port = c.services_filter(f'Service=="{target}"')

            # 传入当前span
            span = SetActiveSpanSource(span)
            intercept = open_tracing_client_interceptor(opentracing.global_tracer(), span)
            ch = grpc.insecure_channel(f"{host}:{port}")
            ch = intercept_channel(ch, intercept)
            stub = inventory_pb2_grpc.InventoryStub(ch)
        except Exception as e:
            logger.info(f"{e.args}")

        return stub

    @logger.catch
    def _assgin_from_db(self, row)->goods_pb2.GoodsDetailResponse:
        rsp = goods_pb2.GoodsDetailResponse()
        rsp.id = row.id
        rsp.category_id = row.category_id
        rsp.brand_id = row.brand_id
        rsp.name = row.name
        rsp.goods_sn = row.goods_sn
        rsp.subtitle = row.subtitle
        rsp.market_price = row.market_price
        rsp.shop_price = row.shop_price
        rsp.sold_num = row.sold_num
        rsp.click_num = row.click_num
        rsp.fav_num = row.fav_num
        rsp.on_sale = row.on_sale
        rsp.is_new = row.is_new
        rsp.is_hot = row.is_hot
        rsp.ship_free = row.ship_free
        rsp.front_image = row.front_image
        if row.slide_images:
            for item in json.loads(row.slide_images):
                rsp.images.append(item)
        if row.desc_images:
            for path in json.loads(row.desc_images):
                rsp.desc_images.append(path)
            
        return rsp


    @logger.catch
    def CreateGoods(self, req: goods_pb2.GoodsRequest, context)->goods_pb2.GoodsDetailResponse:
        rsp = goods_pb2.GoodsDetailResponse()
        try:
            # Category.get(req.category_id)
            # Brands.get(req.brand_id)
            row = {
                Goods.name: req.name,
                Goods.goods_sn: req.goods_sn,
                Goods.subtitle: req.subtitle,
                Goods.category_id: req.category_id,
                Goods.brand_id: req.brand_id,
                Goods.market_price: req.market_price,
                Goods.shop_price: req.shop_price,
                # Goods.stocks:req.stocks,
                Goods.ship_free: req.ship_free,
                Goods.front_image: req.front_image,
                Goods.slide_images: json.dumps(list(req.images)),
                Goods.desc_images: json.dumps(list(req.desc_images)),
            }
            id = Goods.insert(row).execute()
            if id>0:
                rsp = goods_pb2.GoodsDetailResponse(
                    id=id,
                )
        except Exception as e:
            logger.warning("mysql insert err: {}",e.args)
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)
            context.set_details(e.args)

        return rsp


    @logger.catch
    def GoodsList(self, req: goods_pb2.GoodsFilterRequest, context)->goods_pb2.GoodsListResponse:
        rsp = goods_pb2.GoodsListResponse()
        offset = req.page * req.limit
        try:
            goods = Goods.select()
            if req.keywords:
                goods = goods.where(Goods.name.contains(req.keywords))
            if req.is_hot:
                goods = goods.filter(Goods.is_hot == req.is_hot)
            if req.is_new:
                goods = goods.filter(Goods.is_new == req.is_new)
            if req.price_min:
                goods = goods.filter(Goods.shop_price >= req.price_min)
            if req.price_max:
                goods = goods.filter(Goods.shop_price <= req.price_max)
            if req.brand_id:
                goods = goods.filter(Goods.brand_id == req.brand_id)
            if req.category_id:
                pass
            goods = goods.limit(req.limit).offset(offset)
            for item in goods:
                r = self._assgin_from_db(item)
                rsp.data.append(r)
            rsp.total = goods.count()
        except Goods.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在:{e.args}")

        # TODO: 此处待完善库存

        return rsp


    @logger.catch
    def BatchGetGoods(self, req: goods_pb2.BatchGoodsByIdRequest, context)->goods_pb2.GoodsListResponse:
        #TODO: 此处待完善库存
        rsp = goods_pb2.GoodsListResponse()
        try:
            rows = Goods.select().where(Goods.id.in_(list(req.id)))
            for row in rows:
                rsp.data.append(self._assgin_from_db(row))
            rsp.total = rows.count()
        except Goods.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在:{str(e)}")
        except Exception as e:
            context.set_code(grpc.StatusCode.UNKNOWN)
            context.set_details(f"错误:{str(e)}")
        
        return rsp


    @logger.catch
    def GetGoodsDetail(self, req: goods_pb2.GoodsByIdRequest, context)->goods_pb2.GoodsDetailResponse:
        rsp = goods_pb2.GoodsDetailResponse()
        parent_span = context.get_active_span()
        with opentracing.global_tracer().start_span("mysql_goods", child_of=parent_span) as span:
            try:
                row = Goods.get(Goods.id==req.id)
                rsp = self._assgin_from_db(row)
            except Goods.DoesNotExist as e:
                # logger.warning(f"goods not exist: {e.args}")
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details(f"商品记录不存在：{e.args}")
            except Exception as e:
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(f"商品查询错误:{e.args}")
        
        with opentracing.global_tracer().start_span("grpc_inventory", child_of=parent_span) as span2:
            #库存
            try:
                stocks = self._get_proto_stub(span2).InvDetail(inventory_pb2.GoodsInvInfo(goodsId=req.id))
            except Exception as e:
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(f"库存查询错误:{e.args}")

        return rsp


    @logger.catch
    def DeleteGoods(self, req: goods_pb2.GoodsByIdRequest, context)->empty_pb2.Empty:
        try:
            goods = Goods.get(req.id)
            goods.delete_instance()
        except Goods.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在:{e.args}")
        
        return empty_pb2.Empty()


    @logger.catch
    def UpdateGoods(self, req: goods_pb2.GoodsRequest, context)->empty_pb2.Empty:
        try:
            goods = Goods.get(req.id)
            if req.name:
                goods.name = req.name
            goods.goods_sn = req.goods_sn
            goods.subtitle = req.subtitle
            if req.category_id:
                goods.category_id = req.category_id
            if req.brand_id:
                goods.brand_id = req.brand_id
            if req.market_price:
                goods.market_price = req.market_price
            if req.shop_price:
                goods.shop_price = req.shop_price
            goods.on_sale = req.on_sale
            goods.is_new = req.is_new
            goods.is_hot = req.is_hot
            goods.ship_free = req.ship_free
            goods.front_image = req.front_image
            goods.slide_images = json.dumps(list(req.images))
            goods.desc_images = json.dumps(list(req.desc_images))
            rows = goods.save()
            if rows == 0:
                raise Goods.DoesNotExist("未更新")
        except Goods.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在：{e.args}")

        return empty_pb2.Empty()
 

    @logger.catch
    def UpdateStatus(self, req: goods_pb2.GoodsRequest, context)->empty_pb2.Empty:
        try:
            fields = {
                "on_sale" : req.on_sale,
                "is_new" : req.is_new,
                "is_hot" : req.is_hot
            }
            rows = Goods.update(fields).where(Goods.id==req.id).execute()
            if rows == 0:
                raise Goods.DoesNotExist("未更新")
        except Goods.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在：{e.args}")

        return empty_pb2.Empty()
