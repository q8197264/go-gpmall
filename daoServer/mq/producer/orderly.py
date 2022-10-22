from rocketmq.client import Producer, Message
import sys
from os.path import dirname, abspath

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from config import ROOT_DIR,ORDERLY_TOPIC,ORDERLY_GROUP,HOST,BROKER_ADDR

def create_message(body:str, keys:str, tags:str="", properties:dict={}):
    msg = Message(ORDERLY_TOPIC)
    msg.set_tags(tags)
    msg.set_keys(keys)
    for k in properties:
        msg.set_property(k, properties[k])
    msg.set_body(body)
    return msg
    
def send_orderly_message():
    p = Producer("orderly_producer",max_message_size=1024*1024)
    p.set_name_server_address(HOST)
    p.start()

    users = [
        {"uid":"223","name":"lucy", "age":"1"},
        {"uid":"224","name":"lily", "age":"1"},
        {"uid":"225","name":"hanmeimei", "age":"1"},
        {"uid":"226","name":"lilei", "age":"1"},
        {"uid":"225","name":"hanmeimei", "age":"2"},
        {"uid":"226","name":"lilei", "age":"2"},
        {"uid":"224","name":"lily", "age":"2"},
        {"uid":"224","name":"lily", "age":"3"},
        {"uid":"226","name":"lilei", "age":"3"},
        {"uid":"224","name":"lily", "age":"4"},
        {"uid":"223","name":"lucy", "age":"2"},
        {"uid":"225","name":"hanmeimei", "age":"3"},
        {"uid":"226","name":"lilei", "age":"4"},
        {"uid":"225","name":"hanmeimei", "age":"4"},
        {"uid":"223","name":"lucy", "age":"3"},
        {"uid":"223","name":"lucy", "age":"4"},
    ]
    for item in users:
        msg = create_message("顺序消息:%s"%item["uid"], "key:%s"%item["uid"], "tagA,tagB", {"user":item["name"],"age":item["age"]})
        res = p.send_orderly_with_sharding_key(msg, str(item["uid"]))
        print(res.status, res.msg_id)

    p.shutdown()


if __name__ == "__main__":
    send_orderly_message()