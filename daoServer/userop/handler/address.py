import sys
from os.path import dirname, abspath
import grpc

from google.protobuf.empty_pb2 import Empty
from loguru import logger

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from config.config import LOG_PATH
from proto import address_pb2, address_pb2_grpc
from model.models import user_address as Address

logger.add(LOG_PATH, rotation="00:00", level="DEBUG")

class AddressServicer(address_pb2_grpc.AddressServicer):

    @logger.catch
    def QueryAddressList(self, req: address_pb2.AddressRequest, context) -> address_pb2.AddressListResponse:
        rsp = address_pb2.AddressListResponse()
        try:
            address = Address.select().where(Address.user_id==req.user_id)
            rsp.total = address.count()
            for item in address:
                r = address_pb2.AddressDetailResponse()
                r.id = item.id
                r.province = item.province
                r.city = item.city
                r.district = item.district
                r.address = item.address
                r.signer_name = item.signer_name
                r.signer_mobile = item.signer_mobile
                r.is_default = item.is_default
                rsp.data.append(r)
        except Address.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在:{e.args}")
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"错误：{e.args}")

        return rsp


    @logger.catch
    def AddAddress(self, req: address_pb2.AddressRequest, context) -> Empty:
        try:
            if not req.is_default:
                req.is_default = False
            address = Address(
                user_id=req.user_id,
                province=req.province,
                city=req.city,
                district=req.district,
                address=req.address,
                signer_name=req.signer_name,
                signer_mobile=req.signer_mobile,
                is_default = req.is_default
            )
            address.save()
        except Address.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在:{e.args}")
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"错误：{e.args}")

        return Empty()


    @logger.catch
    def UpdateAddress(self, req: address_pb2.AddressRequest, context)->Empty:
        try:
            address = Address.get(Address.id == req.id)
            address.user_id = req.user_id
            address.province = req.province
            address.city = req.city
            address.district = req.district
            address.address = req.address
            address.signer_name = req.signer_name
            address.signer_mobile = req.signer_mobile
            address.is_default = req.is_default
            address.save()
        except Address.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在:{e.args}")
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"错误:{e.args}")

        return Empty()


    @logger.catch
    def DeleteAddress(self, req: address_pb2.AddressRequest, context)->Empty:
        try:
            Address.delete().where(Address.id==req.id, Address.user_id==req.user_id).execute()
        except Address.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在:{e.args}")
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"错误:{e.args}")
        
        return Empty()
