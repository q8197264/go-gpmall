import string
import sys
from os.path import dirname,abspath
from datetime import datetime
import time
import json

import grpc
import opentracing
from rocketmq.client import TransactionMQProducer, TransactionStatus, Producer, ConsumeStatus, Message, SendStatus
from loguru import logger
from google.protobuf.empty_pb2 import Empty

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from proto import order_pb2_grpc, order_pb2, goods_pb2,goods_pb2_grpc,inventory_pb2,inventory_pb2_grpc
from model.models import JOIN,shop_cart as ShopCart,order_info as OrderInfo,order_goods as OrderGoods
from config.config import cfg,db,LOG_PATH
from common.consul import consul
from common.grpc_interceptor.retry import RetryInterceptor
from common.grpc_opentracing import grpcext, open_tracing_client_interceptor, ActiveSpanSource

logger.add(LOG_PATH, rotation="00:00", level="DEBUG")

# import logging
# log = logging.getLogger("peewee")
# log.setLevel(logging.DEBUG)
# log.addHandler(logging.StreamHandler())

local_pspan_dict = {}
local_order_dict = {}

class SetActiveSpanSource(ActiveSpanSource):
    def __init__(self, active_span = None):
        self.active_span = active_span
    
    def get_active_span(self):
        return self.active_span


@logger.catch
def cancel_timeout_order(msg):
    #  订单超时返还库存
    try:
        msg = json.loads(msg.body.decode("utf-8"))
        orderSn = msg["order_sn"]
        pspan = local_pspan_dict[msg["span_id"]]
        with opentracing.tracer.start_span("cancel_timeout_order", pspan) as span:
            with db.atomic() as txn:
                try:
                    oi = OrderInfo.get(OrderInfo.order_sn==orderSn)
                    if oi.status != "TRADE_SUCCESS":
                        # 超时未付款: 设置订单状态为 status == CLOSE
                        oi.status = "TRADE_CLOSE"
                        oi.save()

                        # 归还库存 transTopic
                        p = Producer("timeOutProducerGroup")
                        p.set_name_server_address(f"{cfg['rocketmq']['host']}:{cfg['rocketmq']['port']}")
                        p.start()
                        msg = Message("transTopic")
                        msg.set_keys("order")
                        msg.set_tags("reback")
                        msg.set_body(json.dumps({"order_sn":orderSn}))
                        ret = p.send_sync(msg)
                        if ret.status != SendStatus.OK:
                            raise Exception("send message fail: 归还库存失败")
                        p.shutdown()
                except OrderInfo.DoesNotExist as e:
                    txn.rollback()
                    logger.info(f"订单编号 {orderSn} 记录不存在: {e.args[0]}")
                    # return  ConsumeStatus.CONSUME_SUCCESS
                except Exception as e:
                    txn.rollback()
                    logger.info(f"异常(发警告短信给开发人员): {e.args[0]}")
                    # return ConsumeStatus.RECONSUME_LATER
        print("订单已支付 -- ",orderSn)
    except Exception as e:
        logger.info("发警告短信给开发人员 ",e.args)
        # return ConsumeStatus.RECONSUME_LATER

    return ConsumeStatus.CONSUME_SUCCESS


class OrderServicer(order_pb2_grpc.Order):

    def __init__(self) -> None:
        super().__init__()

    @logger.catch
    def _transformToProto(self, item)->order_pb2.OrderGoodsDetailResponse:
        rsp = order_pb2.OrderGoodsDetailResponse()
        rsp.order_id = item.order_id
        rsp.goods_id  = item.goods_id
        rsp.goods_name = item.goods_name
        rsp.market_price = item.market_price
        rsp.shop_price  = item.shop_price
        rsp.nums  = item.nums

        return rsp
    

    @logger.catch
    def _grpcClient(self, srv_name, active_span):
        clientSub = None
        try:
            client = consul.Consul(cfg["consul"]["host"], cfg["consul"]["port"])
            host, port = client.services_filter(f'Service=="{srv_name}"')
            # 三种状态进行重试
            retry_codes = [grpc.StatusCode.UNAVAILABLE, grpc.StatusCode.UNKNOWN, grpc.StatusCode.DEADLINE_EXCEEDED]
            if srv_name in ["goods-dao","inventory-dao"]:
                ch = grpc.insecure_channel(f"{host}:{port}")
                ch = grpc.intercept_channel(ch, RetryInterceptor(max_retries=3, retry_codes=retry_codes))

                span = SetActiveSpanSource(active_span)
                intercept = open_tracing_client_interceptor(tracer=opentracing.global_tracer(),active_span_source=span)
                ch = grpcext.intercept_channel(ch, intercept)
                if srv_name == "goods-dao":
                    clientSub =  goods_pb2_grpc.GoodsStub(ch)
                elif srv_name== "inventory-dao":
                    clientSub = inventory_pb2_grpc.InventoryStub(ch)
        except Exception as e:
            logger.info("grpc client exception:",e.args)
        
        return clientSub


    @staticmethod
    def _createOrderSn(user_id)->string:
        if user_id:
            from random import randint
            rdt = randint(100, 999)
            return f'{datetime.now().strftime("%Y%m%d%H%M%S")}{user_id}{rdt}'
        return None


    @logger.catch
    def CreateOrder(self, req: order_pb2.OrderRequest, context)->Empty:
        try:
            # 设置 parent span
            pspan = context.get_active_span()
            local_pspan_dict[pspan.context.span_id] = pspan
            orderSn = self._createOrderSn(req.user_id)
            p = TransactionMQProducer("transProducerGroup", self.check_callback)
            p.set_name_server_address(f"{cfg['rocketmq']['host']}:{cfg['rocketmq']['port']}")
            p.start()
            msg = Message("transTopic")
            msg.set_tags("shipping, order")
            msg.set_keys(orderSn)
            msg.set_body(json.dumps({
                "order_sn":orderSn,
                "span_id": pspan.context.span_id,
                "user_id":req.user_id,
                "signer_name":req.name,
                "signer_mobile":req.mobile,
                "signer_address":req.address,
                "post":req.post,
            }))
            ret = p.send_message_in_transaction(msg, self.local_execute, user_args=None)
            # SendResult(status=<SendStatus.OK: 0>, msg_id='C0A808DE6CE825BC7D63103C44DC0000', offset=621)
            # print(ret.status)
        except Exception as e:
            logger.info(f"{e.args}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"{e.args}")
            return Empty()

        while True:
            if orderSn in local_order_dict and local_order_dict[orderSn]:
                p.shutdown()
                if local_order_dict[orderSn]["code"] == grpc.StatusCode.OK:
                    context.set_code(local_order_dict[orderSn]["code"])
                    context.set_details(local_order_dict[orderSn]["detail"])
                    return order_pb2.OrderDetailResponse(
                        id = local_order_dict[orderSn]["order"]["id"], 
                        order_sn = orderSn,
                        pay_amount = local_order_dict[orderSn]["order"]["total"],
                        goods=local_order_dict[orderSn]["order"]["goods"],
                        user_id = req.user_id
                    )
                else:
                    return order_pb2.OrderDetailResponse()
            time.sleep(0.1)
        

    @logger.catch
    def local_execute(self, msg, user_args=None):
        """
            事务消息: 半消息
            需要 追踪- 需要 tracer
        """
        gids = []
        buy_nums = dict()
        total = 0
        msg = json.loads(msg.body.decode("utf-8"))

        pspan = local_pspan_dict[msg["span_id"]]
        with opentracing.tracer.start_span("local_execute_shopcart", child_of=pspan) as span:
            orderSn = msg["order_sn"]
            local_order_dict[orderSn] = {}
            try:
                # 获取购买商品与数量
                carts = ShopCart.select().where(ShopCart.user_id==msg["user_id"], ShopCart.checked==True)
                if carts:
                    for item in carts:
                        gids.append(item.goods_id)
                        buy_nums.update({item.goods_id:item.nums})
                else:
                    raise ShopCart.DoesNotExist("记录不存在")
            except ShopCart.DoesNotExist as e:
                local_order_dict[orderSn]["code"] = grpc.StatusCode.NOT_FOUND
                local_order_dict[orderSn]["detail"] = f"没有选择购物车商品"
                return TransactionStatus.ROLLBACK
            except Exception as e:
                local_order_dict[orderSn]["code"] = grpc.StatusCode.UNKNOWN
                local_order_dict[orderSn]["detail"] = f"查询购物车记录失败: {e.args}"
                return TransactionStatus.ROLLBACK

        with opentracing.tracer.start_span("local_execute_goods", child_of=pspan) as span:
            try:
                # grpc 商品价格总计
                client = self._grpcClient(cfg["goods"]["name"], span)
                goods_list = client.BatchGetGoods(goods_pb2.BatchGoodsByIdRequest(id = gids), timeout=1)
                if not goods_list:
                    raise ValueError("商品表内没有对应有的商品")
                for item in goods_list.data:
                    # 计算订单总金额
                    total += buy_nums[item.id] * item.shop_price
            except ValueError as e:
                local_order_dict[orderSn]["code"] = grpc.StatusCode.NOT_FOUND
                local_order_dict[orderSn]["detail"] = f"记录不存在: {e.args}"
                return TransactionStatus.ROLLBACK
            except Exception as e:
                local_order_dict[orderSn]["code"] = grpc.StatusCode.INTERNAL
                local_order_dict[orderSn]["detail"] = f"访问商品服务失败: {e.args}"
                return TransactionStatus.ROLLBACK
        
        with opentracing.tracer.start_span("local_execute_inventory", child_of=pspan) as span:
            try:
                # grpc 减库存
                inv_rsp = inventory_pb2.SellInfo()
                inv_rsp.order_sn = orderSn
                for gid in gids:
                    r = inventory_pb2.GoodsInvInfo()
                    r.goodsId = gid
                    r.num = buy_nums[gid]
                    inv_rsp.data.append(r)
                inv_client = self._grpcClient(cfg["inventory"]["name"],span)
                inv_client.Sell(inv_rsp, timeout=1)
            except grpc.RpcError as e:
                local_order_dict[orderSn]["code"] = grpc.StatusCode.INTERNAL
                local_order_dict[orderSn]["detail"] = f"扣减库存失败: {e.args}"
                err_code = e.code()
                if err_code == grpc.StatusCode.UNKNOWN or err_code == grpc.StatusCode.DEADLINE_EXCEEDED:
                    # 库存不足超时(有可能), 库存服务没有扣减库存, 发送库存归还
                    return TransactionStatus.COMMIT
                else:
                    TransactionStatus.ROLLBACK
            except Exception as e:
                local_order_dict[orderSn]["code"] = grpc.StatusCode.INTERNAL
                local_order_dict[orderSn]["detail"] = f"访问库存服务失败: {e.args}"
                return TransactionStatus.ROLLBACK

        with opentracing.tracer.start_span("local_execute_create_order", child_of=pspan) as span:
            # 开启事务 
            with db.atomic() as txn:
                # 创建订单(order)与订单商品表(order_goods)
                try:
                    if orderSn is None:
                        raise ValueError("build order_sn fail!")
                    fields = {
                        "order_sn": orderSn,
                        "user_id":msg["user_id"],
                        "coupon": 0,
                        "delivery":0,
                        "signer_name":msg["signer_name"],
                        "signer_mobile":msg["signer_mobile"],
                        "signer_address":msg["signer_address"],
                        "post":msg["post"],
                        "amount":total,
                        "pay_amount":total,
                        "status":"WAIT_BUYER_PAY",
                        "pay_mode":"alipay",
                        "pay_time": datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                    }
                    order_id = OrderInfo.insert_many(fields).execute()
                    if not order_id:
                        raise ValueError("order_id 生成失败")
                
                    fields = []
                    for row in goods_list.data:
                        field = {}
                        field["order_id"] = order_id
                        field["goods_id"] = row.id
                        field["goods_name"] = row.name
                        field["market_price"] = row.market_price
                        field["shop_price"] = row.shop_price
                        field["nums"] = buy_nums[row.id]
                        field["goods_image"] = row.front_image
                        fields.append(field)
                    OrderGoods.insert_many(fields).execute()
                
                    # 删除购物车中已购商品
                    ShopCart.delete().where(ShopCart.user_id==msg["user_id"],ShopCart.checked==True).execute()
                
                    # 订单创建成功
                    local_order_dict[orderSn] = {
                        "code": grpc.StatusCode.OK,
                        "detail": f"订单创建成功",
                        "order": {
                            "id": order_id,
                            "order_sn":orderSn,
                            "total":total,
                            "goods":fields,
                            "user_id": msg["user_id"],
                        },
                    }

                    # 发送订单延时消息
                    self.send_delay_order(orderSn, msg["span_id"])

                except Exception as e:
                    local_order_dict[orderSn]["code"] = grpc.StatusCode.INTERNAL
                    local_order_dict[orderSn]["detail"] = f"订单创建失败: {e.args[0]}"
                    txn.rollback()
                    return TransactionStatus.COMMIT # 订单创建失败, 库存归还
        
            # 如果发生异常, 订单创建失败, 发送
            return TransactionStatus.ROLLBACK


    @logger.catch
    def check_callback(self, msg):
        """
            事务消息回查
        """
        msg = json.loads(msg.body.decode("utf-8"))
        pspan = local_pspan_dict[msg["span_id"]]
        with opentracing.tracer.start_span("reduce_inventory_check_callback", child_of=pspan) as span:
            print(msg)
            orderInfo = OrderInfo.select().where(OrderInfo.order_sn==msg["order_sn"])
            if orderInfo:
                return TransactionStatus.ROLLBACK
            else:
                # 确认发送消息, 触发库存归还
                return TransactionStatus.COMMIT


    @logger.catch
    def send_delay_order(self, order_sn, span_id):
        msg = Message("delayOrder")
        msg.set_keys(order_sn)
        msg.set_tags("delay order")
        msg.set_delay_time_level(2)
        msg.set_body(json.dumps({
            "order_sn":order_sn,
            "span_id": span_id,
        }))
        p = Producer("delayOrderProducerGroup")
        p.set_name_server_address(f"{cfg['rocketmq']['host']}:{cfg['rocketmq']['port']}")
        p.start()
        ret = p.send_sync(msg)
        if ret.status == SendStatus.OK:
            print("-- send_delay_order --")
        p.shutdown()

        return ret


    @logger.catch
    def QueryOrderList(self, req: order_pb2.OrderRequest, context)->order_pb2.OrderListResponse:
        rsp = order_pb2.OrderListResponse()
        try:
            orderInfo = (OrderInfo.select(OrderInfo,OrderGoods)
                .join(OrderGoods, JOIN.INNER, on=(OrderInfo.id==OrderGoods.order_id))
            )
            if req.user_id:
                orderInfo = orderInfo.where(OrderInfo.user_id==req.user_id)
                
            exists = {}
            for row in orderInfo.objects():
                if row.id not in exists.keys():
                    r = order_pb2.OrderDetailResponse()
                    r.id = row.id
                    r.order_sn = row.order_sn
                    r.user_id = row.user_id
                    r.coupon = row.coupon
                    r.delivery = row.delivery
                    r.signer_name = row.signer_name
                    r.signer_mobile = row.signer_mobile
                    r.signer_address = row.signer_address
                    r.amount = row.amount
                    r.pay_amount = row.pay_amount
                    r.status = row.status
                    r.pay_time = str(row.pay_time)
                    exists[row.id] = r
                exists[row.id].goods.append(self._transformToProto(row))
                
            for _,item in exists.items():
                rsp.data.append(item)
            rsp.total = len(exists)
        except OrderInfo.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"记录不存在: {e.args}")
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"查询订单列表失败：{str(e)}")

        return rsp
    

    @logger.catch
    def QueryOrderDetail(self, req:order_pb2.OrderRequest, context)->order_pb2.OrderDetailResponse:
        
        # with self.tracer.start_span("select_mysql",child_of=context.get_active_span()) as sel:
        #     time.sleep(1)

        rsp = order_pb2.OrderDetailResponse()
        try:
            orderInfo = (OrderInfo.select(OrderInfo,OrderGoods)
                .join(OrderGoods, JOIN.INNER, on=(OrderInfo.id==OrderGoods.order_id))
            )
            if req.user_id:
                orderInfo = orderInfo.where(OrderInfo.id==req.id, OrderInfo.user_id==req.user_id)
            else:
                orderInfo = orderInfo.where(OrderInfo.id==req.id)
            exists = {}
            for row in orderInfo.objects():
                if row.id not in exists.keys():
                    rsp.id = row.id
                    rsp.order_sn = row.order_sn
                    rsp.user_id = row.user_id
                    rsp.coupon = row.coupon
                    rsp.delivery = row.delivery
                    rsp.signer_name = row.signer_name
                    rsp.signer_mobile = row.signer_mobile
                    rsp.signer_address = row.signer_address
                    rsp.amount = row.amount
                    rsp.pay_amount = row.pay_amount
                    rsp.status = row.status
                    rsp.pay_time = str(row.pay_time)
                    exists[row.id] = {}
                rsp.goods.append(self._transformToProto(row))
        except OrderInfo.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"订单不存在:{str(e)}")
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"查询订单详情失败:{str(e)}")
        
        return rsp

    @logger.catch
    def DelOrder(self, req: order_pb2.OrderRequest, context)->Empty:
        try:
            OrderInfo.delete().where(OrderInfo.id == req.id).execute()
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"删除订单失败:{str(e)}")
        return Empty()


    @logger.catch
    def UpdateOrderStatus(self, req:order_pb2.OrderStatusRequest, context)->Empty:
        try:
            order = OrderInfo.get(OrderInfo.order_sn==req.order_sn)
            order.status= req.status
            order.save()
        except OrderInfo.DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"订单不存在:{str(e)}")
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"更新订单状态失败:{str(e)}")
        
        return Empty()
