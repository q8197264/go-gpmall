import sys
from os.path import dirname,abspath
import grpc

from loguru import logger
from google.protobuf.empty_pb2 import Empty

sys.path.insert(0, dirname(dirname(abspath(dirname(__file__)))))
from userop.config.config import cfg
from common.consul.consul import Consul
from userop.proto import favorites_pb2,favorites_pb2_grpc,goods_pb2,goods_pb2_grpc
from userop.model.models import user_favorites as Favorites


class FavServicer(favorites_pb2_grpc.FavoritesServicer):
    _goods_client = None
    _host = None
    _port = None

    def __init__(self):
        self._init_goods_client()

    def _init_goods_client(self):
        if self._goods_client is None:
            cs = Consul(cfg['consul']['host'], cfg['consul']['port'])
            self._host,self._port = cs.services_filter('Service=="{}"'.format(cfg["goods"]["name"]))
            if self._host is not None and self._port is not None:
                ch = grpc.insecure_channel(f"{self._host}:{self._port}")
                self._goods_client = goods_pb2_grpc.GoodsStub(ch)
            
        return self._goods_client


    @logger.catch
    def QueryFav(self, req: favorites_pb2.UserFavRequest, context)->favorites_pb2.FavResponse:
        rsp = favorites_pb2.FavResponse()
        try:
            favorites = Favorites.select().where(Favorites.user_id==req.user_id, Favorites.goods_id==req.goods_id)    
            for row in favorites:
                rsp.goods_id = row.goods_id
                rsp.user_id = row.user_id
                res = self._goods_client.GetGoodsDetail(goods_pb2.GoodsByIdRequest(id=row.goods_id))
                rsp.title = res.name
                rsp.shop_price = res.shop_price
                rsp.on_sale = res.on_sale
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"srv错误：{e.args}")

        return rsp

    @logger.catch
    def QueryFavList(self, req: favorites_pb2.UserFavRequest, context)->favorites_pb2.FavListResponse:
        rsp = favorites_pb2.FavListResponse()
        try:
            favorites = Favorites.select()
            if req.user_id:
                favorites = favorites.where(Favorites.user_id==req.user_id)

            favs = dict()
            for row in favorites:
                favs.update({row.goods_id:row.user_id})

            if len(favs) <= 0:
                raise Favorites.DoesNotExist("没有收藏")

            goodslist = self._goods_client.BatchGetGoods(goods_pb2.BatchGoodsByIdRequest(id=list(favs.keys())))
            for goods in goodslist.data:
                if goods.id in favs.keys():
                    r = favorites_pb2.FavResponse()
                    r.goods_id = goods.id
                    r.user_id = favs[goods.id]
                    r.title = goods.name
                    r.shop_price = goods.shop_price
                    r.on_sale = goods.on_sale
                    rsp.data.append(r)
            rsp.total = len(rsp.data)

        except Favorites.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在:{e.args}")
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"错误：{e.args}")
        
        return rsp

    @logger.catch
    def AddFav(self, req: favorites_pb2.UserFavRequest, context)->Empty:
        try:
            Favorites.insert(
                user_id = req.user_id,
                goods_id = req.goods_id,
                is_deleted=0
            ).on_conflict(
                preserve=[Favorites.user_id,Favorites.goods_id,Favorites.is_deleted]
            ).execute()
        except Favorites.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"{e.args}")
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"{e.args}")
        
        return Empty()


    @logger.catch
    def DeleteFav(self, req: favorites_pb2.UserFavRequest, context)->Empty:
        try:
            Favorites.delete().where(Favorites.user_id==req.user_id,Favorites.goods_id==req.goods_id).execute()
        except Favorites.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"错误:{e.args}")
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"错误:{e.args}")
        
        return Empty()