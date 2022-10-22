# -*- coding=utf-8 -*-
import argparse
import os
from signal import SIGINT,SIGTERM,signal
from datetime import date

from loguru import logger

from consumer.general import GeneralConsumer
from consumer.orderly import OrderlyConsumer
from consumer.transaction import TransConsumer
from config import ROOT_DIR,GENERAL_TOPIC,GENERAL_GROUP,ORDERLY_TOPIC,ORDERLY_GROUP,TRANS_TOPIC,TRANS_GROUP,HOST,BROKER_ADDR

logger.add(f"logs/{date.today().day}.log")

# mac  没有librocketmq.so
@logger.catch
def main():
    exit()

    topic, group, host, broker_addr = commonline_args()
    
    if topic == "transTopic":
        # python3 server.py  --topic=transTopic --group=transConsumerGroup
        g = TransConsumer()
        g.start_transaction_consumer(topic,group,host,broker_addr)
    elif topic == "orderlyTopic":
        # python3 server.py --topic=orderlyTopic --group=orderlyConsumerGroup
        g = OrderlyConsumer()
        g.start_orderly_consumer(topic,group,host,broker_addr)
    else:
        g = GeneralConsumer()
        g.start_general_consumer(topic,group,host,broker_addr)
        

@logger.catch
def exit():
    def f(singal, frame):
        yn = input("是否关闭 y/n:")
        if yn == "y":
            # sys.exit(0)
            os._exit(0)
        else:
            print("继续...")

    signal(SIGINT, f)
    signal(SIGTERM, f)


@logger.catch
def commonline_args()->tuple[str, str, str, str]:
    parse = argparse.ArgumentParser()
    parse.add_argument(
        "--topic",
        nargs="?",
        type=str,
        default=GENERAL_TOPIC,
        help="bind topic"
    )
    parse.add_argument(
        "--group",
        nargs="?",
        type=str,
        default=GENERAL_GROUP,
        help="bind group"
    )
    parse.add_argument(
        "--host",
        nargs="?",
        type=str,
        default=HOST,
        help="bind host"
    )
    parse.add_argument(
        "--broker_addr",
        nargs="?",
        type=str,
        default=BROKER_ADDR,
        help="bind broker_addr"
    )
    args = parse.parse_args()

    return args.topic, args.group, args.host, args.broker_addr


if "__main__" == __name__:
    main()