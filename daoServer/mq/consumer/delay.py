from rocketmq.client import PushConsumer,ComsumeStatus
import sys
from os.path import dirname, abspath

sys.path.insert(0, dirname(abspath(dirname(__file__))))
from config import ROOT_DIR,GENERAL_TOPIC,GENERAL_GROUP,ORDERLY_TOPIC,ORDERLY_GROUP,TRANS_TOPIC,TRANS_GROUP,HOST,BROKER_ADDR

class DelayConsumer(object):

    @staticmethod
    def callback(msg):
        print("消费:", msg.body, msg.properties)
        return ConsumeStatus.CONSUME_SUCCESS

    def start_send_delay_consumer():
        c = PushConsumer(GENERAL_GROUP)
        c.set_name_server_address(HOST)
        c.subscribe(GENERAL_TOPIC, callback)
        c.start()

        # while True:
        #     time.sleep(3600)