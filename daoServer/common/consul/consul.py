import requests
from common.consul.base import Base

class Consul(Base):
    def __init__(self, host, port) -> None:
        self.baseUrl = f"http://{host}:{port}/v1/agent/service" 
        self.headers = {
            "Content-Type":"application/json"
        }

    def register(self, host, port, name, service_id, tags=[], check=True)->bool:
        params = {
            "Name":name,
            "ID":service_id,
            "tags":tags,
            "Port":port,
            "Address":host
        }
        if check:
            params["Check"] = {
                "GRPC":f"{host}:{port}",
                "DeregisterCriticalServiceAfter": "5s",
                "Interval": "10s",
                "Timeout": "5s"

            }
        rsp = requests.put(f"{self.baseUrl}/register", headers=self.headers, json=params)
        if rsp.status_code == 200:
            # print(f"{name}服务注册成功")
            return True
        else:
            # print(f"{name}服务注册失败: %d %s" % (rsp.status_code, rsp.reason))
            return False

    def deregister(self, service_id):
        rsp = requests.put(f"{self.baseUrl}/deregister/{service_id}")
        if rsp.status_code == 200:
            # print("服务注销成功")
            return True
        else:
            # print(f"服务注销失败:{rsp.status_code} {rsp.reason}")
            return False


    def services(self):
        rsp = requests.get(f"{self.baseUrl}s")
        if rsp.status_code == 200:
            for k,v in rsp.json().items():
                print(k,v)

        else:
            print(f"fail:{rsp.status_code} {rsp.reason}")


    def services_filter(self, filter):
        params = {
            "filter":filter
        }
        rsp = requests.get(f"{self.baseUrl}s", params=params)
        if rsp.status_code == 200:
            for k,v in rsp.json().items():
                return v["Address"], v["Port"]
        else:
            print(f"fail:{rsp.status_code} {rsp.reason}")
            return None, None
    

if "__main__" == __name__:
    c = Consul("192.168.1.122",8500)
    c.deregister("user-dao_1")
    c.register("192.168.1.122",4343,"user-dao","user-dao_1")
    # c.services()
    # c.services_filter('Service == "user-web"')

    