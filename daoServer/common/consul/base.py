from abc import ABC,abstractmethod

class Base(ABC):
    @abstractmethod
    def register(self, host, port, name, service_id):
        pass

    @abstractmethod
    def deregister(self, service_id):
        pass

    @abstractmethod
    def services(self):
        pass

    @abstractmethod
    def services_filter(self, filter):
        pass