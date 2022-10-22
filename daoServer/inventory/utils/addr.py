from os import environ
import socket

def getFreePort(default_port=0)->int:
    # debug = environ.get("GO_WEBSERVER_DEBUG_CONFIG")
    if default_port > 0:
        return default_port
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.bind(("",0))
    _, port = s.getsockname()
    s.close()
    
    return port

if "__main__" == __name__:
    print(getFreePort())