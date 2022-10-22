import sys
from os.path import dirname, abspath
from rocketmq.client import Producer, Message
import time

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from config import ROOT_DIR,GENERAL_TOPIC,GENERAL_GROUP,ORDERLY_TOPIC,ORDERLY_GROUP,TRANS_TOPIC,TRANS_GROUP,HOST,BROKER_ADDR


def create_message(topic, body, keys:str, tags:list=[], properties:dict={}):
    msg = Message(topic)
    for tag in tags:
        msg.set_tags(tag)
    for k in properties:
        msg.set_property(k, properties[k])
    msg.set_keys(keys)
    msg.set_body(body)

    return msg

def send_delay_message():
    p = Producer(GENERAL_GROUP)
    p.set_name_server_address(HOST)
    p.start()
    time.time()
    msg = create_message(GENERAL_TOPIC, "延迟消息"+time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(int(round((time.time()+60))))), "key:delay", ["tagA","tagB"], {"user":"sai+"},)
    # msg.set_property('__STARTDELIVERTIME', str(int(round(time.time()+60))))
    msg.set_delay_time_level(2) #1s 5s 10s 30s 1m 2m 3m ... 10m 20m 30m 1h 2h
    ret = p.send_sync(msg)
    print('send delay message status: '+str(ret.status)+ ' msgId: ' + ret.msg_id +' time: '+str(int(round(time.time()+60))))

    p.shutdown()


if __name__ == "__main__":
    send_delay_message()