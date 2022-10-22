from rocketmq.client import  PushConsumer, ConsumeStatus, dll
import time
import logging

class TransConsumer(object):
    
    def __init__(self):
        # logging.basicConfig(level=logging.CRITICAL, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
        # self.logger = logging.getLogger(__name__)   
        # self.consumer = PushConsumer("PID-XXX")
        # self.consumer.set_namesrv_addr("XX.XX.XX.XX:XXXX")
        # self.topic_name = "xxx"
        # 减少日志输出
        # dll.SetPushConsumerLogLevel(namesrv_addr.encode('utf-8'), 1)
        pass
    

    @staticmethod
    def callback(msg):
        print(msg.id, msg.body.decode('utf-8'))
        return ConsumeStatus.CONSUME_SUCCESS
    
    def start_transaction_consumer(self, topic, group, host, broker_addr):
        c = PushConsumer(group)
        c.set_name_server_address(host)
        c.subscribe(topic, self.callback)
        c.start()
        print("事务消费者开始 ...")
        while True:
            time.sleep(3600)

        c.shutdown()
    
if __name__ == "__main__":
    t = TransConsumer()
    t.start_transaction_consumer()