from os.path import dirname,abspath

ROOT_DIR = dirname(abspath(dirname(__file__)))

GENERAL_TOPIC = "generalTopic"
GENERAL_GROUP = "generalConsumerGroup"

ORDERLY_TOPIC = "orderlyTopic"
ORDERLY_GROUP = "orderlyConsumerGroup"

TRANS_TOPIC = "transTopic"
TRANS_GROUP = "transConsumerGroup"

HOST = "192.168.8.222:9876"
BROKER_ADDR = "192.168.8.222:10911"