import sys
from os import path,environ
import socket

sys.path.insert(0, path.dirname(path.abspath(path.dirname(__file__))))
from config import config

def getFreePort()->int:
    debug = environ.get("GO_WEBSERVER_DEBUG_CONFIG")
    if debug:
        return config.SRV_PORT
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.bind(("",0))
    _, port = s.getsockname()
    s.close()
    
    return port
