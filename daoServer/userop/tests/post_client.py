import sys
from os.path import dirname,abspath

import grpc

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from config.config import cfg
from proto import post_pb2,post_pb2_grpc

class PostClient():
    def __init__(self):
        ch = grpc.insecure_channel(f"{cfg['host']}:{cfg['port']}")
        self.client = post_pb2_grpc.PostStub(ch)

    def addPost(self, user_id, type, subject, message, file):
        self.client.AddPost(post_pb2.UserPostRequest(
            user_id = user_id,
            type = type,
            subject = subject,
            message = message,
            file = file
        ))

    def queryPost(self, uid, page, limit):
        rsp = self.client.QueryPostList(post_pb2.UserPostFilterRequest(
            user_id=uid,
            page = page,
            limit = limit
        ))
        print(rsp)

if "__main__"==__name__:
    p = PostClient()
    # p.addPost(1, 1, "subject1", "message123", "https://www.ccc/sdf.jpg")
    p.queryPost(1, 1, 10)