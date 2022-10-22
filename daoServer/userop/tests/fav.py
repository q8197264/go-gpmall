import sys
from os.path import dirname,abspath
import grpc

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from config import cfg
from proto import favorites_pb2,favorites_pb2_grpc

class client():
    def __init__(self):
        ch = grpc.insecure_channel(f"{cfg['host']}:{cfg['port']}")
        self.client = favorites_pb2_grpc.FavoritesStub(ch)

    def addFav(self, uid, gid):
        self.client.AddFav(favorites_pb2.UserFavRequest(
            user_id=uid,
            goods_id=gid,
        ))

    def queryFav(self, uid, gid):
        rsp = self.client.QueryFav(favorites_pb2.UserFavRequest(
            user_id=uid,
            goods_id=gid,
        ))
        print(rsp)
    
    def queryFavs(self, uid):
        rsp = self.client.QueryFavList(favorites_pb2.UserFavRequest(
            user_id = uid,
        ))
        print(rsp)
    
    def deleteFav(self, uid, gid):
        rsp = self.client.DeleteFav(favorites_pb2.UserFavRequest(
            user_id=uid,
            goods_id=gid
        ))
        print(rsp)


if __name__ == "__main__":
    c = client()
    # c.addFav(1,2)
    # c.addFav(1,3)
    c.queryFav(1,3)
    # c.queryFavs(1)
    # c.deleteFav(1,3)