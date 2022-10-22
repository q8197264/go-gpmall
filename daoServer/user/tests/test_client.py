from os.path import dirname
import sys
import time
from datetime import date

import grpc
from passlib.hash import pbkdf2_sha256
from loguru import logger

sys.path.append(dirname(dirname(__file__)))
from config import config
from proto import user_pb2,user_pb2_grpc


logger.add(f"user/logs/client_{date.today().day}.log", rotation="1 days")


class client:

    @logger.catch
    def __init__(self):
        self.ch = grpc.insecure_channel(f"{config.USER_SERVER_HOST}:{config.USER_SERVER_PORT}")
        self.cli = user_pb2_grpc.UserStub(self.ch)

    
    def __del__(self):
        print("close ch")


    @logger.catch
    def createUser(self, mobile, nick_name, password):
        rsp: user_pb2.UserInfoResponse = self.cli.CreateUser(user_pb2.CreateUserRequest(
            mobile    = mobile,
            nick_name = nick_name,
            password  = password,
            birthday = int(time.time())
        ))
        print(rsp)
        return rsp


    @logger.catch
    def getUserInfo(self, mobile):
        if mobile:
            req = user_pb2.MobileRequest(mobile=mobile)
            rsp: user_pb2.UserInfoResponse = self.cli.GetUserInfo(req)
            print(rsp)
            return rsp


    @logger.catch
    def getUserById(self, uid):
        req = user_pb2.UidRequest(uid=uid)
        rsp: user_pb2.UserInfoResponse = self.cli.GetUserById(req)
        print(rsp)
        return rsp

    
    @logger.catch
    def getUserList(self, page:int, limit:int):
        rsp: user_pb2.UserListResponse = self.cli.GetUserList(user_pb2.PageRequest(
            page = page,
            limit = limit
        ))
        print(rsp)
        return rsp


    @logger.catch
    def updateUserInfo(self, uid, mobile, nick_name, password):
        rsp: user_pb2.UserInfoResponse = self.cli.UpdateUserInfo(user_pb2.UpdateUserRequest(
            uid       = uid,
            data = user_pb2.CreateUserRequest(
                mobile    = mobile,
                nick_name = nick_name,
                password  = password,
                birthday  = int(time.time()),
                # birthday = 1637078198,
                gender = "female",
                role = 1,
                avatar = "https://www.img.cn/img.png",
                desc = "简介：...",
                country = "china",
                provice = "四川省",
                city = "成都",
                area = "武候区",
                address = "康桥新世界5-2-101"
            )
        ))

        return rsp


if "__main__" == __name__:
    # client().createUser("16861291814", "amyly", "toomanysecrets")
    # client().getUserInfo("16861291879")
    # client().getUserById(21)
    client().getUserList(1,10)
    # client().updateUserInfo(40, "16861291818", "amyl12", "passwd123")

    