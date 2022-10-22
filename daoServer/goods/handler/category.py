# -*- coding: UTF-8 -*-
from sys import path
from os.path import dirname,abspath
from datetime import date
import json
from unicodedata import category

import grpc
from peewee import *
from google.protobuf import empty_pb2
from loguru import logger
logger.add(f"logs/category_{date.today().day}.log",  rotation="00:00", level="DEBUG")

path.insert(0, dirname(abspath(dirname(__file__))))
from proto import goods_pb2
from model.models import *
from handler.brand import BrandServicer

# 无限分类
def toTree(arr, parent_id, cache=[], level=0):
    tree = []
    level += 1
    for item in arr:
        if item["id"] in cache:
            continue
        
        if item["parent_id"]==parent_id or (item["parent_id"]==item["id"] and level==parent_id):
            cache.append(item["id"])
            item["child"] = toTree(arr, item["id"], cache, level)
            tree.append(item)
    return tree


class CategoryServicer(BrandServicer):
       
    @logger.catch
    def CreateCategory(self, req: goods_pb2.CategoryRequest, context)->empty_pb2.Empty:
        rsp = goods_pb2.CategoryRequest()
        try:
            c = Category.create(
                name=req.name,
                parent_id= req.parent_id,
                level= req.level
            )
            c.save()
        except Category.DoesNotExist as e:
            logger.warning("新建记录失败:{}", e.args)
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)
            context.set_details(e.args)

        return rsp


    @logger.catch
    def CategoryList(self, req: goods_pb2.CategoryFilterRequest, context) -> goods_pb2.CategoryListResponse:
        rsp = goods_pb2.CategoryListResponse()
        try:
            rows = Category.select()
            if req.id:
                ct = rows.where(Category.id==req.id).get()
                rows = rows.where(Category.parent_id == ct.parent_id)

            res = []
            for item in rows:
                rsp.data.append(goods_pb2.CategoryInfoResponse(
                    name =item.name.strip(),
                    id   =item.id,
                    parent_id = item.parent_id.id,
                    level  =item.level,
                    is_tab =item.is_tab
                ))
                res.append({
                    "name":item.name.strip(),
                    "id":item.id,
                    "parent_id": item.parent_id.id,
                    "level":item.level,
                    "is_tab":item.is_tab
                })
            rsp.JsonData = json.dumps(toTree(res, 1, []))
        except Category.DoesNotExist as e:
            logger.warning("类目记录不存在:{}", e.args)
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(e.args)
        
        return rsp

    @logger.catch
    def DeleteCategory(self, req: goods_pb2.CategoryByIdRequest, context)->empty_pb2.Empty:
        try:
            cates = Category.select()

            # 重构sql结果集
            res = []
            for item in cates:
                res.append({
                    "name":item.name.strip(),
                    "id":item.id,
                    "parent_id": item.parent_id.id,
                    "level":item.level,
                    "is_tab":item.is_tab
                })
            def exportid(rows, parent_id, cache=[]):
                for item in rows:
                    if item["id"] in cache:
                        continue
                    if item["parent_id"]==parent_id:
                        cache.append(item["id"])
                        for i in exportid(rows, item["id"], cache):
                            yield i
                        yield item
            ids = []
            for item in exportid(res,1,[]):
                ids.append(item["id"])
                
            category = Category.update(is_deleted=True).where(Category.id.in_(ids))
            rows = category.execute()
            if rows==0:
                raise Category.DoesNotExist("未更新")
        except Category.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在：{e.args}")
        
        return empty_pb2.Empty()


    @logger.catch
    def UpdateCategory(self, req: goods_pb2.CategoryRequest, context)->empty_pb2.Empty:
        try:
            category = Category.get(req.id)
            category.name = req.name
            if req.parent_id > 0:
                category.parent_id = req.parent_id
            # category.level  = req.level
            category.is_tab = req.is_tab
            rows = category.save()
            if rows == 0:
                raise Category.DoesNotExist("未更新")
        except Category.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在：{e.args}")
        
        return empty_pb2.Empty()