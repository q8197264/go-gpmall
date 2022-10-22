import sys
import os
from rocketmq.client import Message, TransactionMQProducer, TransactionStatus
import time
from signal import signal,SIGINT, SIGTERM


sys.path.insert(0, os.path.dirname(os.path.abspath(os.path.dirname(__file__))))
from config import TRANS_TOPIC,TRANS_GROUP,HOST,BROKER_ADDR


def create_message():
    msg = Message(TRANS_TOPIC)
    msg.set_tags("tagA")
    msg.set_keys("key")
    msg.set_property("user","234")
    msg.set_body("事务消息内容test1")

    return msg


def check_callback(msg):
    print("检查回调", msg.body.decode('utf-8'))
    return TransactionStatus.COMMIT


def local_execute(msg, user_args):
    print("执行本地业务:", msg.body.decode('utf-8'))

    return TransactionStatus.UNKNOWN


def send_transaction_message(nums):
    exit()

    p = TransactionMQProducer("transProducerGroup", check_callback)
    p.set_name_server_address(HOST)
    p.set_max_message_size(1024*4)
    p.start()
    # for i in range(0, nums):
    msg = create_message()
    try:
        res = p.send_message_in_transaction(msg, local_execute, None)
        print('send message status: ' + str(res.status) + ' msgId: ' + res.msg_id)
        print("消息发送完毕")
    except Exception as e:
        print(e)

    while True:
        time.sleep(1)

    # p.shutdown()


# 优雅退出
def exit():
    def f(signal, frame):
        print("producer关闭成功")
        os._exit(0)

    signal(SIGINT, f)
    signal(SIGTERM, f)

if __name__ == "__main__":
    send_transaction_message(5)