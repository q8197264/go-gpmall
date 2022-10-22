from datetime import date

import grpc
import opentracing
from peewee import *
from passlib.hash import pbkdf2_sha256
from loguru import logger
from google.protobuf import empty_pb2

from config import config
from proto import user_pb2, user_pb2_grpc
from model.models import mc_user as Users, mc_user_info as UsersInfo, mc_user_address as UsersAddress

logger.add(f"{config.LOG_PATH}/server_{date.today().day}.log", rotation="12:00")

# import logging
# log = logging.getLogger("peewee")
# log.setLevel(logging.DEBUG)
# log.addHandler(logging.StreamHandler())

"""
    用户服务数据接口
     创建用户
     获取用户列表
     获取用户信息
     更新用户信息
"""
class User(user_pb2_grpc.UserServicer):

    # Users database data assign to proto
    @logger.catch
    def _assign_from_db(self, user)->user_pb2.UserInfoResponse:
        rsp = user_pb2.UserInfoResponse()
        rsp.id      = user.id
        rsp.mobile  = user.mobile
        rsp.nick_name = user.nick_name
        rsp.password  = user.password
       
        if user.gender:
            rsp.gender = user.gender
        if user.desc:
            rsp.desc = user.desc
        if user.role:
            rsp.role = user.role
        if user.avatar:
            rsp.avatar = user.avatar

        # # 用户信息
        try:
            # 格式不匹配int = date
            if user.birthday:
                rsp.birthday = user.birthday
        except Exception as e:
            logger.debug("user_info表数据获取失败: {}", e.args)

        # 用户地址管理
        try:
            if user.country:
                rsp.country = user.country
            if user.provice:
                rsp.provice = user.provice
            if user.city:
                rsp.city = user.city
            if user.area:
                rsp.area = user.area
            if user.address:
                rsp.address = user.address
        except Exception as e:
            logger.debug("user_address表数据获取失败: {}", e.args)

        return rsp

    """
        获取用户信息
            @param request(
                mobile: str = ?
                nick_name: str = ?,
                password: str  = ?
            )
    """
    @logger.catch
    def CreateUser(self, request: user_pb2.CreateUserRequest, context) -> user_pb2.UserInfoResponse:
        # 插入数据库
        try:
            if not request.mobile:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details('mobile为空')

            Users.get(Users.mobile == request.mobile)
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)
            context.set_details('用户已存在!')

            return user_pb2.UserInfoResponse()
        except Users.DoesNotExist:
            row = {
                Users.mobile: request.mobile,
                Users.nick_name: request.nick_name,
                Users.password: pbkdf2_sha256.hash(request.password),
            }
            uid = Users.insert_many(row).execute()
            try:
                row = {
                    UsersInfo.user_id: uid,
                    UsersInfo.birthday: request.birthday
                }
                UsersInfo.insert_many(row).execute()
            except Exception:
                logger.info("错误：uid=%d ... %s", uid, e)

            try:
                row = {
                    UsersAddress.user_id: uid,
                    UsersAddress.country: request.country
                }
                UsersAddress.insert_many(row).execute()
            except Exception:
                logger.info("错误：uid=%d ... %s", uid, e)

            return self.GetUserById(user_pb2.UidRequest(uid=uid), context)


    """
        检查登陆密码
    """
    @logger.catch
    def CheckPassword(self, request: user_pb2.CheckPasswordRequest, context)->user_pb2.CheckPasswordResponse:
        return user_pb2.CheckPasswordResponse(success=pbkdf2_sha256.verify(request.password,request.encryptedPassword))


    """
        获取用户列表
        @param request(
            page: int = ?,
            limit: int = ?
        )
    """
    @logger.catch
    def GetUserList(self, request: user_pb2.PageRequest, context) -> user_pb2.UserListResponse:
        rsp = user_pb2.UserListResponse()
        page = (request.page-1) * request.limit
        try:
            # count Users record
            res = Users.raw("SELECT COUNT(id) as `total` FROM `mc_user`").execute()[0]
            
            # 连表查
            sql = """
                SELECT 
                    `t1`.`id`, `t1`.`mobile`, `t1`.`password`, `t1`.`nick_name`, `t1`.`avatar`, 
                    `t1`.`gender`, `t1`.`desc`, `t1`.`role`, `t2`.`user_id`, `t2`.`birthday`, 
                    `t3`.`user_id`, `t3`.`country`,`t3`.`provice`, `t3`.`city`, `t3`.`area`, 
                    `t3`.`address` 
                FROM `mc_user` AS `t1` 
                LEFT OUTER JOIN `mc_user_info` AS `t2` ON (`t1`.`id` = `t2`.`user_id`) 
                LEFT OUTER JOIN `mc_user_address` AS `t3` ON (`t1`.`id` = `t3`.`user_id`) 
                LIMIT %s OFFSET %s
            """
            users = Users.raw(sql, request.limit, page).execute()
            for user in users:
                rsp.data.append(self._assign_from_db(user))

            rsp.total = res.total
        except Exception as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在:{e.args}")

        return rsp


    """
        获取用户信息
        @param request(
            mobile: str = ?
        )
    """
    @logger.catch
    def GetUserInfo(self, request: user_pb2.MobileRequest, context) -> user_pb2.UserInfoResponse:
        rsp = user_pb2.UserInfoResponse()
        try:
            user = (Users.select(Users, UsersInfo, UsersAddress)
                .join(UsersInfo, JOIN.LEFT_OUTER, on=(Users.id==UsersInfo.user_id))
                .join(UsersAddress, JOIN.LEFT_OUTER, on=(Users.id==UsersAddress.user_id))
                .where(Users.mobile==request.mobile)
            ).objects()
            if len(user):
                rsp = self._assign_from_db(user[0])
            else:
                raise Users.DoesNotExist()
        except Users.DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("用户不存在!")

        return rsp

    """
        获取用户信息 通过用户id
        @param request(
            uid: int = ?
        )
    """
    @logger.catch
    def GetUserById(self, request: user_pb2.UidRequest, context) -> user_pb2.UserInfoResponse:
        try:
            res = (Users
                .select(Users,UsersInfo,UsersAddress)
                .join(UsersInfo, JOIN.LEFT_OUTER, on=(Users.id==UsersInfo.user_id))
                .join(UsersAddress, JOIN.LEFT_OUTER, on=(Users.id==UsersAddress.user_id))
                .where(Users.id==request.uid)
            ).objects()
            if len(res):
                return self._assign_from_db(res[0])
            else:
                raise Users.DoesNotExist()
        except Users.DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("用户不存在!")


    """
        更新用户信息
        @param request(
            uid = ?,
            mobile = ?,
            nick_name = ?,
            ...
        )

        UPDATE student s LEFT JOIN class c ON s.class_id = c.id SET s.class_name='test22',c.stu_name='test22'
    """
    @logger.catch
    def UpdateUserInfo(self, request: user_pb2.UpdateUserRequest, context) -> empty_pb2.Empty:
        
        # 更新主表 Users
        try:
            # if request.data.password != request.data.repassword:
            #     raise ValueError("您输入的密码不一致")

            set_fields = {}
            set_fields['mobile'] = request.data.mobile
            set_fields['nick_name'] = request.data.nick_name
            # set_fields['password'] = request.data.password
            if hasattr(request.data, 'gender'): 
                set_fields['gender'] = request.data.gender
            if hasattr(request.data, 'desc'): 
                set_fields['desc'] = request.data.desc
            if hasattr(request.data, 'role'): 
                set_fields['role'] = request.data.role
            if hasattr(request.data, 'avatar'): 
                set_fields['avatar'] = request.data.avatar
            row = Users.get(Users.id==request.uid)
            row.id and Users.update(set_fields).where(Users.id==request.uid).execute()
        except ValueError as ve:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(ve.args)
            return
        except Exception as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("user 更新参数异常!")
            logger.warning("update错误：[user] mobile={} <= {}", request.data.mobile, e.args)
            return

        # 更新用户信息表 UsersInfo
        try:
            UsersInfo.insert(
                user_id=request.uid,
                birthday=request.data.birthday
            ).on_conflict(
                preserve=[UsersInfo.user_id],
                update={UsersInfo.birthday:request.data.birthday}
            ).execute()
        except Exception as e:
            logger.warning("on_conflict错误：[info] uid={} <= {}", request.uid, e.args)

        # 更新用户地址表 UsersAddress
        try:
            set_fields = {}
            set_fields['user_id'] = request.uid
            set_fields['country'] = request.data.country
            set_fields['provice'] = request.data.provice
            set_fields['city'] = request.data.city
            set_fields['area'] = request.data.area
            set_fields['address'] = request.data.address
            UsersAddress.insert(set_fields).on_conflict(
                update=set_fields
            ).execute()
        except Exception as e:
            logger.warning("on_conflict错误：[address] uid={} <= {}", request.uid, e.args)

        return empty_pb2.Empty()
