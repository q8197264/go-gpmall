import time
import json
from rocketmq.client import PushConsumer, ConsumeStatus


class GeneralConsumer(object):

    def callback(self, msg):
        print("callback",msg.id, msg.body.decode("utf-8"), msg.get_property('name').decode("utf-8"))
        return ConsumeStatus.CONSUME_SUCCESS

    def start_general_consumer(self, topic: str, group:str, host:str, broker_addr:str):
        p = PushConsumer(group)
        p.set_name_server_address(host)
        p.subscribe(topic, self.callback)

        print("消费开始...",topic)
        p.start()

        while True:
            time.sleep(3600)
