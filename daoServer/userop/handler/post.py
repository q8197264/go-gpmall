import sys
from os.path import dirname, abspath

import grpc
from loguru import logger
from google.protobuf.empty_pb2 import Empty

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from config.config import LOG_PATH
from model.models import leaving_message as Post
from proto import post_pb2,post_pb2_grpc

logger.add(LOG_PATH, rotation="00:00", level="DEBUG")

class PostServicer(post_pb2_grpc.PostServicer):

    @logger.catch
    def AddPost(self, req: post_pb2.UserPostRequest, context) -> Empty:
        try:
            post = Post(
                user_id=req.user_id,
                type=req.type,
                subject=req.subject,
                message=req.message,
                file=req.file)
            post.save(force_insert=True)
        except Post.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在:{e.args}")
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"错误:{e.args}")

        return Empty()

    @logger.catch
    def QueryPostList(self, req: post_pb2.UserPostFilterRequest, context) -> post_pb2.PostListResponse:
        rsp = post_pb2.PostListResponse()
        try:
            if not req.page or req.page<1:
                req.page = 1
            offset = (req.page-1) * req.limit
            
            rows = Post.select().where(Post.user_id == req.user_id).offset(offset).limit(req.limit)
            for row in rows:
                r = post_pb2.PostInfoResponse()
                r.user_id = row.user_id
                r.type = row.type
                r.subject = row.subject
                r.message = row.message
                r.file = row.file
                rsp.data.append(r)
            rsp.total = 0
        except Post.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在:{e.args}")
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"错误: {e.args}")

        return rsp
