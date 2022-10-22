from sys import path
from os.path import dirname,abspath
from datetime import date

import grpc
from peewee import *
from google.protobuf import empty_pb2
from loguru import logger
logger.add(f"logs/goods_{date.today().day}.log",  rotation="00:00")

path.insert(0, dirname(abspath(dirname(__file__))))
from proto import goods_pb2,goods_pb2_grpc
from model.models import *

class BannerServicer(goods_pb2_grpc.Goods):

    @logger.catch
    def GetBannerList(self, req: goods_pb2.BannerFilterRequest, context)->goods_pb2.BannerListResponse:
        rsp = goods_pb2.BannerListResponse()
        offset = req.limit * (req.page-1)
        try:
            banners = Banner.select().limit(req.limit).offset(offset)
            rsp.total = banners.count()
            for item in banners:
                r = goods_pb2.BannerInfoResponse()
                r.id    = item.id
                r.index = item.index
                r.image = item.image
                r.url  = item.url
                rsp.data.append(r)
        except Banner.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在：{e.args}")
        
        return rsp


    @logger.catch
    def CreateBanner(self, req: goods_pb2.BannerRequest, context)->goods_pb2.BannerInfoResponse:
        rsp = goods_pb2.BannerInfoResponse()
        try:
            banner = Banner()
            banner.index = req.index
            banner.image = req.image
            banner.url = req.url
            banner.save()
            rsp.id = banner.id
        except Banner.DoesNotExist as e:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(f"参数有误:{e.args}")
        
        return rsp

    @logger.catch
    def UpdateBanner(self, req: goods_pb2.BannerByIdRequest, context)->empty_pb2.Empty:
        try:
            banner = Banner.get(req.id)
            banner.index = req.index
            banner.image = req.image
            banner.url = req.url
            rows = banner.save()
            if rows == 0:
                raise Banner.DoesNotExist("未更新")
        except Banner.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在：{e.args}")
        
        return empty_pb2.Empty()


    @logger.catch
    def DeleteBanner(self, req: goods_pb2.BannerByIdRequest, context)->empty_pb2.Empty:
        try:
            banner = Banner.get(req.id)
            banner.delete_instance()
        except Banner.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在:{e.args}")
        
        return empty_pb2.Empty()
