import sys
from os.path import dirname,abspath
from datetime import datetime

from peewee import *
from playhouse.mysql_ext import JSONField

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from config.config import db

class BaseModel(Model):
    add_time = DateTimeField(default=datetime.now, verbose_name="添加时间")
    update_time = DateTimeField(default=datetime.now, verbose_name="更新时间")
    is_deleted = BooleanField(default=False, verbose_name="是否删除")

    @classmethod
    def select(cls, *fields):
        return super().select(*fields).where(cls.is_deleted==False)

    @classmethod
    def delete(cls, permanently = False):
        if permanently:
            return super().delete()
        else:
            return super().update(is_deleted=True)

    def save(self, *args, **keyargs):
        if self._pk is not None:
            self.update_time = datetime.now()
        return super().save(*args, **keyargs)

    def delete_instance(self, peramently=False):
        if peramently:
            return super().delete_instance().where(self._pk_expr).execute()
        else:
            self.is_deleted = True
            return self.save()

    class Meta:
        database = db

class shop_cart(BaseModel):
    id = AutoField(primary_key=True, verbose_name="购物车id")
    user_id = IntegerField(index=True, verbose_name="用户id")
    goods_id = IntegerField(index=True, verbose_name="商品id")
    nums = IntegerField(verbose_name="购买数量")
    checked = BooleanField(index=True, default=False, verbose_name="确认选中")


class order_info(BaseModel):
    STATUS = (
        ("WAIT_BUYER_PAY", "交易创建，等待买家付款"),
        ("TRADE_CLOSED", "未付款交易超时关闭，或支付完成后全额退款"),
        ("TRADE_SUCCESS", "交易支付成功"),
        ("TRADE_FINISHED", "交易结束，不可退款")
    )
    id = AutoField(primary_key=True, verbose_name="订单id")
    order_sn = CharField(max_length=20, index=True, unique=True, verbose_name="订单编号")
    user_id = IntegerField(index=True, verbose_name="用户id")
    coupon = DecimalField(max_digits=10, decimal_places=2, verbose_name="优惠券")
    delivery = IntegerField(default=0, verbose_name="配送方式")
    signer_address = CharField(verbose_name="配送地址")
    signer_mobile = CharField(max_length=11, verbose_name="收货人电话")
    signer_name = CharField(max_length=2, verbose_name="收货人")
    post = CharField(verbose_name="留言")
    amount = DecimalField(max_digits=10, decimal_places=2, verbose_name="总金额")
    pay_amount = DecimalField(max_digits=10, decimal_places=2, verbose_name="实际应付金额")
    status = CharField(choices=STATUS,verbose_name="订单状态")
    pay_mode = CharField(default="", verbose_name="支付方式")
    pay_time = DateTimeField(default=datetime.now, verbose_name="付款时间")


class order_goods(BaseModel):
    id = AutoField(primary_key=True, verbose_name="订单商品id")
    order_id = IntegerField(index=True, verbose_name="订单id")
    goods_id = IntegerField(index=True, verbose_name="商品id")
    goods_name = CharField(verbose_name="商品名")
    goods_image = JSONField(verbose_name="商品图")
    nums = IntegerField(default=0, verbose_name="购买数量")
    market_price = DecimalField(max_digits=10, decimal_places=2, verbose_name="市场价")
    shop_price = DecimalField(max_digits=10, decimal_places=2, verbose_name="商品价")
    

if "__main__" == __name__:
    # db.drop_tables([shop_cart, order_info, order_goods])
    db.create_tables([shop_cart, order_info, order_goods])