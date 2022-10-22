import sys
from os.path import dirname,abspath
from datetime import datetime

from peewee import *

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from config import db

class BaseModel(Model):
    add_time = DateTimeField(default=datetime.now, verbose_name="添加时间")
    update_time = DateTimeField(default=datetime.now, verbose_name="更新时间")
    is_deleted = BooleanField(default=False, verbose_name="是否删除")

    @classmethod
    def select(cls, *fields):
        return super().select(*fields).where(cls.is_deleted==False)

    @classmethod
    def delete(cls, permanently=False):
        if permanently:
            return super().delete()
        else:
            return super().update(is_deleted=True)

    def save(self, *args, **keyargs):
        if self._pk is not None:
            self.update_time = datetime.now()
        return super().save(*args, **keyargs)

    def delete_instance(self, permanently=False):
        if permanently:
            return super().delete_instance().where(self._pk_expr).execute()
        else:
            self.is_deleted = True
            return self.save()

    class Meta:
        database = db


class leaving_message(BaseModel):
    MessageType = (
        (1, "留言"),
        (2, "投诉"),
        (3, "询问"),
        (4, "售后"),
        (5, "求购"),
    )
    id = AutoField(primary_key=True, verbose_name="主键")
    user_id = IntegerField(verbose_name="用户id")
    subject = CharField(max_length=100, default="", verbose_name="主题")
    message = CharField(max_length=255, default="", verbose_name="留言内容", help_text="留言内容")
    type = IntegerField(default=1, choices=MessageType, verbose_name="留言类型", help_text=u"留言内型: 1(留言),2(投诉),3(询问),4(售后),5(求购)")
    file = CharField(max_length=100, verbose_name="上传的文件path", help_text="上传的文件path")

class user_address(BaseModel):
    id = AutoField(primary_key=True, verbose_name="主键")
    user_id = IntegerField(verbose_name="用户id")
    province = CharField(max_length=100, default="", verbose_name="省份")
    city = CharField(max_length=100, default="", verbose_name="城市")
    district = CharField(max_length=100, default="", verbose_name="街区")
    address = CharField(max_length=100, verbose_name="详细地址")
    signer_name = CharField(max_length=100, verbose_name="签收人")
    signer_mobile = CharField(max_length=11, verbose_name="签收人电话")
    is_default = BooleanField(default=False, verbose_name="是否默认收货地址")

class user_favorites(BaseModel):
    user_id = IntegerField(verbose_name="用户id")
    goods_id = IntegerField(verbose_name="商品id")

    class Meta:
        primary_key = CompositeKey("user_id","goods_id")


if "__main__" == __name__:
    db.create_tables([leaving_message, user_address, user_favorites])
    # db.drop_tables([leaving_message, user_address, user_favorites])