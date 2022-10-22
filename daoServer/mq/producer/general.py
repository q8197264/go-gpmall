from rocketmq.client import Producer, Message
from os.path import dirname, abspath
import sys
import time

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from config import ROOT_DIR,GENERAL_TOPIC,GENERAL_GROUP,HOST,BROKER_ADDR


def create_message(topic:str, keys:str, tags:list=[], properties:dict={}, body:str=""):
    msg = Message(topic)
    for v in tags:
        msg.set_tags(v)
    for k in properties:
        msg.set_property(k, properties[k])
    msg.set_keys(keys)
    msg.set_body(body)

    return msg


def send_sync():
    producer = Producer(GENERAL_GROUP)
    producer.set_name_server_address(HOST)
    producer.start()

    for i in range(0, 5):
        msg = create_message(GENERAL_TOPIC, "key:1", ['tagA','tagB'], {'name':'sai','age':'18'}, "消息内容")
        
        try:
            ret = producer.send_sync(msg)
            print(ret.status, ret.msg_id, ret.offset)
        except Exception as e:
            print(">> ",e," <<")
    
    time.sleep(1)
    producer.shutdown()

if __name__ == "__main__":
    send_sync()