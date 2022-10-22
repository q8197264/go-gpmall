import socket

def get_free_port(port:int = 0)->int:
    if port > 0:
        return port
    tcp = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    tcp.bind(("", 0))
    _, port = tcp.getsockname()
    tcp.close()

    return port