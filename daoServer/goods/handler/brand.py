from sys import path
from os.path import dirname,abspath
from datetime import date

import grpc
from peewee import *
from google.protobuf import empty_pb2
from loguru import logger
logger.add(f"logs/goods_{date.today().day}.log",  rotation="00:00", level="DEBUG")

path.insert(0, dirname(abspath(dirname(__file__))))
from proto import goods_pb2
from model.models import *
from handler.banner import BannerServicer

class BrandServicer(BannerServicer):

    @logger.catch
    def GetBrandList(self, req: goods_pb2.BrandFilterRquest, context)->goods_pb2.BrandListResponse:
        rsp = goods_pb2.BrandListResponse()
        offset = (req.page-1) * req.limit
        try:
            brands = Brands.select().limit(req.limit).offset(offset)
            rsp.total = brands.count()
            for brand in brands:
                brand_rsp = goods_pb2.BrandInfoResponse()
                brand_rsp.id = brand.id
                brand_rsp.name = brand.name
                brand_rsp.logo = brand.logo
                rsp.data.append(brand_rsp)
        except Brands.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(e.args)

        return rsp
        

    @logger.catch
    def CreateBrand(self, req: goods_pb2.CreateBrandRequest, context)->goods_pb2.BrandInfoResponse:
        rsp = goods_pb2.BrandInfoResponse()
        brands = Brands.select().where(Brands.name == req.name)
        if brands:
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)
            context.set_details("记录已存在")
            return rsp
            
        try:
            brands = Brands()
            brands.name = req.name
            brands.logo = req.logo
            rsp.id = brands.save()
            rsp.name = req.name
            rsp.logo = req.logo
        except Brands.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("记录不存在")

        return rsp


    @logger.catch
    def DeleteBrand(self, req: goods_pb2.BrandByIdRequest, context)->empty_pb2.Empty:
        try:
            brand = Brands.get(req.id)
            brand.delete_instance()
        except Brands.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("记录不存在")
        
        return empty_pb2.Empty()


    @logger.catch
    def UpdateBrand(self, req: goods_pb2.BrandInfoRequest, context)->empty_pb2.Empty:
        try:
            brand = Brands.get(req.id)
            if req.name:
                brand.name = req.name
            if req.logo:
                brand.logo = req.logo
            brand.save()
        except Brands.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("记录不存在")

        return empty_pb2.Empty()


    # 获取 category 下 brands
    @logger.catch
    def GetBrandsByCategory(self, req: goods_pb2.CategoryByIdRequest, context)->goods_pb2.BrandListResponse:
        rsp = goods_pb2.BrandListResponse()
        try:
            res = (Category_Bind_Brands.select(Category, Category_Bind_Brands, Brands)
            .join(Category, JOIN.LEFT_OUTER, on=(Category_Bind_Brands.category_id==Category.id))
            .join(Brands, JOIN.LEFT_OUTER, on=(Category_Bind_Brands.brand_id==Brands.id))
            .where(Category_Bind_Brands.category_id == req.id)
            ).objects()
            rsp.total = res.count()
            for item in res:
                r = goods_pb2.BrandInfoResponse()
                r.id = item.brand.id
                r.name = item.brand.name
                r.logo = item.brand.logo
                rsp.data.append(r)
        except Category_Bind_Brands.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(e.args)

        return rsp


    # 所有分类品牌 查询category_bind_brand 表
    @logger.catch
    def CategoryBrandList(self, req: goods_pb2.CategoryBrandFilterRequest, context)->goods_pb2.CategoryBrandListResponse:
        rsp = goods_pb2.CategoryBrandListResponse()
        offset = req.limit * (req.page -1)
        try:
            category_brands = (Category_Bind_Brands.select(Category_Bind_Brands, Category, Brands)
                .join(Category, JOIN.LEFT_OUTER, on=(Category_Bind_Brands.category_id == Category.id))
                .join(Brands, JOIN.LEFT_OUTER, on=(Category_Bind_Brands.brand_id == Brands.id))
            )
            rsp.total = category_brands.count()
            category_brands = category_brands.limit(req.limit).offset(offset)
            for item in category_brands:
                r = goods_pb2.CategoryBrandResponse()
                r.id = item.id
                r.category.id = item.category.id
                r.category.name = item.category.name
                r.category.parent_id = item.category.parent_id.id
                r.category.level = item.category.level
                r.category.is_tab = item.category.is_tab
                r.brand.id = item.brand.id
                r.brand.name = item.brand.name
                r.brand.logo = item.brand.logo
                rsp.data.append(r)
        except Category_Bind_Brands.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在:{e.args}")
        except Exception as err:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(f"记录找不到:{err.args}")

        return rsp


    @logger.catch
    def UpdateCategoryBrand(self, req: goods_pb2.CategoryBrandRequest, context)->goods_pb2.CategoryBrandResponse:
        rsp = goods_pb2.CategoryBrandResponse()
        try:
            cbb = Category_Bind_Brands.get(req.id)
            brand = Brands.get(req.brand_id)
            cbb.brand = brand
            category = Category.get(req.category_id)
            cbb.category = category
            cbb.save()
        except Category_Bind_Brands.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在：{e.args}")

        return rsp


    @logger.catch
    def CreateCategoryBrand(self, req: goods_pb2.CategoryBrandRequest, context)->goods_pb2.CategoryBrandResponse:
        rsp = goods_pb2.CategoryBrandResponse()
        cbb = Category_Bind_Brands()
        try:
            category = Category.get(Category.id==req.category_id)
            cbb.category_id = category.id
        except Category.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"类目不存在：{e.args}")
            return rsp

        try:
            brand = Brands.get(req.brand_id)
            cbb.brand = brand
            cbb.save()
            rsp.id = cbb.id
        except Brands.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"品牌不存在:{e.args}")

        return rsp


    @logger.catch
    def DeleteCategoryBrand(self, req: goods_pb2.CategoryBrandRequest, context)->empty_pb2.Empty:
        try:
            cbb = Category_Bind_Brands.get(req.id)
            rows = cbb.delete_instance()
        except Category_Bind_Brands.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在:{e.args}")
        
        return empty_pb2.Empty()
