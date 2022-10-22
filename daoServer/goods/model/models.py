from sys import path
from os.path import dirname,abspath
from datetime import datetime

from peewee import *
from playhouse.mysql_ext import JSONField

path.insert(0, dirname(abspath(dirname(__file__))))
from config.config import DB


class BaseModel(Model):
    add_time = DateTimeField(default=datetime.now, verbose_name="添加时间")
    update_time = DateTimeField(default=datetime.now, verbose_name="更新时间")
    is_deleted = BooleanField(default=False, verbose_name="是否删除")

    def save(self, *args, **keyargs):
        if self._pk is not None:
            self.update_time = datetime.now()
        return super().save(*args, **keyargs)

    @classmethod
    def select(cls, *fields):
        return super().select(*fields).where(cls.is_deleted==False)

    @classmethod
    def delete(cls, permanently = False):
        if permanently:
            return super().delete()
        else:
            return super().update(is_deleted=True)

    # 需要执行sql
    def delete_instance(self, permanently = False):
        if permanently:
            return super().delete_instance().where(self._pk_expr).execute()
        else:
            self.is_deleted = True
            return self.save()

    # 公共字段
    class Meta:
        database = DB


class Category(BaseModel):
    # 类目
    id = AutoField(primary_key=True, verbose_name="id")
    name = CharField(max_length=20, verbose_name="类目名")
    parent_id = ForeignKeyField("self", verbose_name="父类id")
    level = IntegerField(default=1, verbose_name="级别")
    is_tab = BooleanField(default=False, verbose_name="是否显示在首页")

class Brands(BaseModel):
    # 品牌
    id = AutoField(primary_key=True, verbose_name="id")
    name = CharField(max_length=20, index=True, verbose_name="品牌名")
    logo = CharField(max_length=200, null=True, verbose_name="logo图标url")


class Goods(BaseModel):
    # 商品信息
    id = AutoField(primary_key=True, verbose_name="id")
    category = ForeignKeyField(Category, verbose_name="类目id", on_delete="CASCADE")
    brand = ForeignKeyField(Brands, verbose_name="品牌id", on_delete="CASCADE")
    name = CharField(max_length=100, verbose_name="商品名")
    goods_sn = CharField(max_length=50, default="", verbose_name="商品序号")
    subtitle = CharField(max_length=200, verbose_name="副标题")
    market_price = FloatField(default=0, verbose_name="市场价")
    shop_price = FloatField(default=0, verbose_name="本店价")
    sold_num = IntegerField(default=0, verbose_name="销量")
    click_num = IntegerField(default=0, verbose_name="点击数")
    fav_num = IntegerField(default=0, verbose_name="收藏数")
    on_sale = BooleanField(default=False, verbose_name="是否上架")
    is_new = BooleanField(default=True, verbose_name="是否新品")
    is_hot = BooleanField(default=False, verbose_name="是否热销")
    ship_free = BooleanField(default=True, verbose_name="是否承担运费")
    front_image = CharField(max_length=200, default="", verbose_name="封面图")
    slide_images = JSONField(verbose_name="轮播图")
    desc_images = JSONField(verbose_name="详情图")


class Category_Bind_Brands(BaseModel):
    # 类目与品牌关系表
    id = AutoField(primary_key=True, verbose_name="id")
    category = ForeignKeyField(Category, verbose_name="类目id", on_delete="CASCADE")
    brand = ForeignKeyField(Brands, verbose_name="品牌id", on_delete="CASCADE")

    class Meta:
        indexes={
            (("category","brand"),True)
        }


class Banner(BaseModel):
    id = AutoField(primary_key=True, verbose_name="id")
    index = IntegerField(default=0, verbose_name="轮播顺序")
    image = CharField(max_length=200, default="", verbose_name="图片url")
    url = CharField(max_length=200, default="", verbose_name="跳转url")


if __name__ == "__main__":
    # DB.drop_tables([Category, Goods, Brands, Banner, Category_Bind_Brands])
    DB.create_tables([Category, Goods, Brands, Banner, Category_Bind_Brands])
    # DB.drop_tables([Banner])
    # DB.create_tables([Banner])