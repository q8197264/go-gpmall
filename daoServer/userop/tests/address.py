from audioop import add
import sys
from os.path import dirname,abspath
import grpc

from google.protobuf.empty_pb2 import Empty

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from config.config import cfg
from proto import address_pb2,address_pb2_grpc
from model.models import user_address as Address


class AddressClient(address_pb2_grpc.AddressServicer):
    def __init__(self):
        ch = grpc.insecure_channel(f"{cfg['host']}:{cfg['port']}")
        self.client = address_pb2_grpc.AddressStub(ch)

    def addAddress(self, user_id, province, city, district, address, signer_name, signer_mobile):
        rsp = self.client.AddAddress(address_pb2.AddressRequest(
            user_id = user_id,
            province = province,
            city = city,
            district = district,
            address = address,
            signer_name = signer_name,
            signer_mobile = signer_mobile,
        ))
        print(rsp)


    def queryAddressList(self):
        rsp = self.client.QueryAddressList(address_pb2.AddressRequest(
            user_id=1
        ))
        print(rsp)


    def updateAddress(self, id, user_id, province, city, district, address, signer_name, signer_mobile):
        self.client.UpdateAddress(address_pb2.AddressRequest(
            id = id,
            user_id = user_id,
            province = province,
            city = city,
            district = district,
            address = address,
            signer_name = signer_name,
            signer_mobile = signer_mobile,
        ))


    def deleteAddress(self,id):
        self.client.DeleteAddress(address_pb2.AddressRequest(
            id=id
        ))

if "__main__"==__name__:
    c = AddressClient()
    # c.addAddress(1, "四川省", "泸州市", "朝阳区", "福天路春天小区5-2-1708", "村里", "19216877654")
    # c.updateAddress(3, 1, "四川省", "泸州市", "朝阳区", "福天路春天小区5-2-1708", "川岛", "19216877655")
    # c.deleteAddress(2)
    c.queryAddressList()