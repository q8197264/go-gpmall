import sys
from os import path
import grpc
import json
from peewee import Model
from jaeger_client.config import Config

sys.path.insert(0, path.dirname(path.abspath(path.dirname(__file__))))
from proto import goods_pb2,goods_pb2_grpc
from config.config import cfg,SRV_HOST,SRV_NAME,SRV_PORT
from common.grpc_opentracing import open_tracing_client_interceptor,ActiveSpanSource
from common.grpc_opentracing.grpcext import intercept_channel
from common.grpc_interceptor import retry


class TestGoods(Model):
    # 测试用例
    def __init__(self):
        cfg = Config(
            config={
                "sampler":{
                    "type":"const",
                    "param":1,
                },
                "logging": True,
            },
            service_name=SRV_NAME,
            validate=True
        )
        tracer = cfg.initialize_tracer()
        intercept = open_tracing_client_interceptor(tracer)
        chan = grpc.insecure_channel(f"{SRV_HOST}:{SRV_PORT}")

        retry_intercept = retry.RetryInterceptor(max_retries=3, retry_codes=[grpc.StatusCode.UNKNOWN, grpc.StatusCode.UNAVAILABLE, grpc.StatusCode.DEADLINE_EXCEEDED])
        chan = grpc.intercept_channel(chan,retry_intercept)
        self.chan = intercept_channel(chan, intercept)
        self.client = goods_pb2_grpc.GoodsStub(self.chan)

    def create(self, req: goods_pb2.GoodsRequest):
        rsp = self.client.CreateGoods(req)
        return rsp

    def getGoodsList(self, page=1, limit=10):
        rsp = self.client.GoodsList(goods_pb2.GoodsFilterRequest(page=page, limit=limit))
        print(rsp)
        return rsp

    def batchGoods(self, ids:list):
        rsp = self.client.BatchGetGoods(goods_pb2.BatchGoodsByIdRequest(id=ids))
        print(rsp)

    def getGoodsDetail(self, id=1):
        rsp = self.client.GetGoodsDetail(goods_pb2.GoodsByIdRequest(id=id))
        print(rsp)
        return rsp

    def createCategory(self, name, parent_id, level):
        rsp = self.client.CreateCategory(goods_pb2.CreateCategoryRequest(
            name=name,
            parent_id=parent_id,
            level=level
        ))
        return rsp
        

    def updateCategory(self, id, name, parent_id, level, is_tab):
        rsp = self.client.UpdateCategory(goods_pb2.CategoryInfoRequest(
            id = id,
            name = name,
            parent_id = parent_id,
            level  = level,
            is_tab = is_tab,
        ))
        return rsp


    def getCategoryList(self):
        rsp = self.client.CategoryList(goods_pb2.CategoryFilterRequest())
        print(json.loads(rsp.JsonData))
        print(rsp.data)
        return rsp

    def categoryBrandsList(self, id):
        rsp = self.client.CategoryBrandList(goods_pb2.CategoryByIdRequest(id=id))
        print(rsp)
        return rsp


    def getCategoryBrandList(self, page, limit):
        rsp = self.client.GetCategoryBrandList(goods_pb2.CategoryBrandFilterRequest(
            page=page, 
            limit=limit
        ))
        print(rsp)
        return rsp
    
    def createBrand(self, name, logo):
        rsp = self.client.CreateBrand(goods_pb2.CreateBrandRequest(name=name, logo=logo))

        return rsp

    def updateBrand(self, id, name, logo):
        rsp = self.client.UpdateBrand(goods_pb2.BrandInfoRequest(id=id, name=name, logo=logo))
        return rsp


    def createCategoryBrand(self, cid, bid):
        rsp = self.client.CreateCategoryBrand(goods_pb2.CategoryBrandRequest(
            category_id=cid, 
            brand_id=bid
        ))
        print(rsp)
        return rsp
    
    def deleteCategoryBrand(self, id):
        rsp = self.client.DeleteCategoryBrand(goods_pb2.CategoryBrandRequest(id=id))
        return rsp

    def updateCategoryBrand(self, id, cid, bid):
        rsp = self.client.UpdateCategoryBrand(goods_pb2.CategoryBrandRequest(
            id=id, 
            category_id=cid, 
            brand_id=bid
        ))
        return rsp

    def createBanner(self, index, image, link):
        rsp = self.client.CreateBanner(goods_pb2.CreateBannerRequest(
            index=index, 
            image=image, 
            link=link
        ))
        return rsp

    def getBannerList(self, page, limit):
        rsp = self.client.GetBannerList(goods_pb2.BannerFilterRequest(page=page, limit=limit))
        print(rsp)
        return rsp

    def updateBanner(self):
        rsp = self.client.UpdateBanner(goods_pb2)
        return rsp

    def delBanner(self, id):
        rsp = self.client.DeleteBanner(goods_pb2.BannerByIdRequest(id=id))
        return rsp

if __name__ == "__main__":
    t = TestGoods()

    # 轮播图
    # t.createBanner(0, "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fimg.zcool.cn%2Fcommunity%2F01e2d558417e53a8012060c8cbbebf.jpg", "https://image.baidu.com")
    # t.getBannerList(0, 10)
    # t.updateBanner()
    # t.delBanner(1)

    # 类目
    # t.createCategory(name="水果", parent_id=1, level=0)
    # t.createCategory(name="李子", parent_id=1, level=1)
    # t.createCategory(name="小李子", parent_id=2, level=2)
    # t.createCategory(name="桃子", parent_id=1, level=1)
    # t.createCategory(name="苹果", parent_id=1, level=1)
    # t.createCategory(name="蔬菜", parent_id=2, level=0)
    # t.createCategory(name="波菜", parent_id=6, level=1)
    # t.createCategory(name="韭菜", parent_id=6, level=1)
    # t.createCategory(name="日用", parent_id=9, level=0)
    # t.updateCategory(id=9, name="日用", parent_id=9, level=0, is_tab=1)

    # t.getCategoryList()
    # t.batchGoods([1,2,3])

    # 品牌
    # t.createBrand("每日优鲜","https://pic.rmb.bdstatic.com/bjh/user/b582c4f33212f92545c637f8515f2d59.jpeg")
    # t.createBrand("美菜", "http://5b0988e595225.cdn.sohucs.com/images/20181210/f3a1fc2c75f34d3cae23e2b0ae46425c.jpg")
    # t.updateBrand(1, "美菜", "http://5b0988e595225.cdn.sohucs.com/images/20181210/f3a1fc2c75f34d3cae23e2b0ae46425c.jpg")
    # t.categoryBrandsList(1)

    # 品牌分类
    # t.getCategoryBrandList(0, 4)
    # t.createCategoryBrand(2,1)
    # t.deleteCategoryBrand(15)
    # t.updateCategoryBrand(14,6,1)
    
    # 商品
    # t.create(goods_pb2.CreateGoodsRequest(
    #     category_id = 1,
    #     brand_id = 1,
    #     name = "棉花0A",
    #     goods_sn = "SN0908978",
    #     subtitle = "简介是是是是是是因为",
    #     market_price = 122.5,
    #     shop_price = 100,
    #     # sint32 sold_num = 8;
    #     # sint32 click_num = 9;
    #     # sint32 fav_num = 10;
    #     # bool on_sale = 11;
    #     # bool is_new = 12;
    #     # bool is_hot = 13;
    #     ship_free = True,
    #     # string front_image = 15;
    #     # string slide_images = 16;
    #     # string desc_images = 17;
    # ))
    # t.getGoodsList(0, 10)
    t.getGoodsDetail(1)
    