from os.path import dirname
import sys
from peewee import *
import time

sys.path.append(dirname(dirname(__file__)))
from config.config import DB

class BaseModel(Model):
    class Meta:
        database = DB

class mc_user(BaseModel):
    GENDER_CHOICES = (
        ("female", "男"),
        ("male","女")
    )

    ROLE_CHOICES = (
        (1,"普通用户"),
        (2,"管理员")
    )

    mobile = CharField(max_length=11, index=True, unique=True, verbose_name="手机号码")
    password = CharField(max_length=100, verbose_name="密码")
    nick_name = CharField(max_length=20, null=True, verbose_name="昵称")
    avatar = CharField(max_length=200, null=True, verbose_name="头像")
    gender = CharField(max_length=6, choices=GENDER_CHOICES, null=True, verbose_name="性别")
    desc = TextField(null=True, verbose_name="个人简介")
    role = IntegerField(default=1, choices=ROLE_CHOICES, verbose_name="用户角色")

class mc_user_address(BaseModel):
    user = ForeignKeyField(mc_user, unique=True, backref="user_address")
    country = CharField(max_length=10, index=True, verbose_name="国家")
    provice = CharField(max_length=10, index=True, verbose_name="省")
    city = CharField(max_length=10, index=True, verbose_name="市")
    area = CharField(max_length=10, index=True, verbose_name="市区")
    address = CharField(max_length=200, verbose_name="详细地址")

class mc_user_info(BaseModel):
    user = ForeignKeyField(mc_user, unique=True, backref="user_info")
    birthday = IntegerField(default=0, verbose_name="生日")


if "__main__" == __name__:
    # DB.drop_tables([mc_user_address, mc_user_info])
    # DB.create_tables([mc_user_address, mc_user_info])
    DB.create_tables([mc_user, mc_user_address, mc_user_info])
    # print(time.time())
