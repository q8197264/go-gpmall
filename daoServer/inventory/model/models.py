from os.path import dirname,abspath
import sys
from datetime import datetime

from peewee import *
from playhouse.mysql_ext import JSONField

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from config.config import DB,client,nacosConfig

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

    class Meta:
        database = DB


class inventory(BaseModel):
    # 商品库存表
    goods_id = IntegerField(verbose_name="商品id", unique=True)
    stocks = IntegerField(verbose_name="库存数量", default=0)
    version = IntegerField(verbose_name="版本号", default=0)

class inventory_history(BaseModel):
    # 出库历史表
    order_sn = CharField(verbose_name="订单编号", max_length=20, unique=True)
    order_inv_detail = CharField(verbose_name="订单详情", max_length=200)
    status = IntegerField(choices=((1, "已扣减"),(2, "已归还")), default=1)
    

if "__main__" == __name__:
    def update_config(args):
        # print(type(args["raw_content"]))
        print(args)

    # client.add_config_watchers(nacosConfig["dataid"], nacosConfig["group"], [update_config])

    # DB.drop_tables([Inventory])
    DB.create_tables([inventory, inventory_history])
    # DB.drop_tables([Inventory])
    # DB.create_tables([Inventory])
    