from rocketmq.client import PushConsumer, ConsumeStatus
import time

class OrderlyConsumer(object):

    @staticmethod
    def callback(msg):
        print(msg.id, msg.tags.decode("utf-8"), msg.body.decode("utf-8"), msg.get_property("user").decode("utf-8"), msg.get_property("age"))
        return ConsumeStatus.CONSUME_SUCCESS

    def start_orderly_consumer(self,topic,group,host,broker_addr):
        c = PushConsumer(group)
        c.set_name_server_address(host)
        c.subscribe(topic, self.callback)

        print("消费开始...",topic)
        c.start()

        while True:
            time.sleep(3600)

        c.shutdown()

if __name__ == "__main__":
    o = OrderlyConsumer()
    o.start_orderly_consumer()

