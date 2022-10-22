import threading
import weakref
from redis import StrictRedis
import redis
from typing import Union

from loguru import logger

from common.lock.py_redis_lock import UNLOCK_SCRIPT

# refresh lock
EXTEND_LOCK_SCRIPT = b"""
    if redis.call("get", KEYS[1]) ~= ARGV[1] then
        return 1
    elseif redis.call("ttl", KEYS[1]) < 0 then
        return 2
    else
        redis.call("expire", KEYS[1], ARGV[2])
        return 0
    end
"""

# delete lock
UNLOCK_SCRIPT = b"""
    if redis.call("get", KEYS[1]) ~= ARGV[1] then
        return 1
    else
        redis.call("del", KEYS[2])
        redis.call("lpush", KEYS[2], 1)
        redis.call("pexpire", KEYS[2], ARGV[2])
        redis.call("del", KEYS[1])
    end
"""

class RLock(object):
    unlock_script = None
    extend_script = None

    _lock_renewal_interval: float
    # 类型 Union[int, str] 表示既可以是 int，也可以是 str
    _lock_renewval_thread: Union[threading.Thread, None]

    def __init__(self, redis_client = None, name = "test", id = None, expire = None, auto_renewal = False, strict=True, signal_expire=1000):
        pool = redis.ConnectionPool(host="127.0.0.1", port=6379, db=1, max_connections=5)
        redis_client = redis.StrictRedis(connection_pool=pool)

        if strict and not isinstance(redis_client, StrictRedis):
            raise ValueError("redis_client must be instance of StrictRedis. Use strict=False if you know what you're doing.")
        if auto_renewal and expire is None:
            raise ValueError("Expire may not be None when auto_renewal is set")

        self._redis_client = redis_client
        self._name = "lock:%s" % name

        if expire:
            expire = int(expire)
            if expire < 0:
                raise RuntimeError("runtimeError: expire")

        self._expire = expire
        self._signal_expire = signal_expire
        self._signal = "lock-signal:" + name
        self._id = id
        self._lock_renewal_interval = float(expire) * 2/3 if auto_renewal else None
        self._lock_renewval_thread = None

        self.register_lua_script(redis_client)

    @classmethod
    def register_lua_script(cls, redis_client):
        cls.unlock_script = redis_client.register_script(UNLOCK_SCRIPT)
        cls.extend_script = redis_client.register_script(EXTEND_LOCK_SCRIPT)

    def extend(self, expire = None):
        self.extend_script(client=self._redis_client, keys=(self._name,), args=(self._id, self._expire))
        print("script lua: ttl = ", expire)

    @staticmethod
    def _lock_signal_renewer(name, lockref, interval, stop):
        """
        """
        # print("Event wait...", stop.is_set())
        # # stop.set()
        # w=not stop.wait(timeout=interval)
        # lock1: RLock = lockref()
        # print(w,'||',lock)
            
        while not stop.wait(timeout=interval):
            lock: "RLock" = lockref()
            print("Refreshing Lock(%r).", name)
            if lock is None:
                print("Stopping loop because Lock(%r) garbage collected." % (name))
                break
            lock.extend(expire=lock._expire)
            del lock
            

    def _start_lock_signal_renewer(self):
        """
        Starts the lock refresher thread.
        开辟看门狗新线程
        """
        print("renewer start...")

        if self._lock_renewval_thread is not None:
            raise RuntimeError("Lock refresh thread already started")

        self._lock_renewval_event = threading.Event()
        self._lock_renewval_thread = threading.Thread(
            group=None,
            target=self._lock_signal_renewer,
            kwargs={
                "name": self._name,
                "lockref": weakref.ref(self),
                "interval": self._lock_renewal_interval,
                "stop": self._lock_renewval_event
            }
        )
        # 设置为守护线程 
        self._lock_renewval_thread.daemon = True
        self._lock_renewval_thread.start()


    def acquire(self, blocking = True, timeout = None):
        """
        Add Lock
        """
        if not blocking and timeout is not None:
            raise RuntimeError("Timeout cannot used if blocking=False")

        if timeout:
            timeout = int(timeout)
            if timeout<0:
                raise RuntimeError("runtimeError:timeout")
            if self._expire and not self._lock_renewal_interval and self._expire < timeout:
                raise RuntimeError("runtimeError: self._expire")

        blocking_timeout = timeout or self._expire or 0
        time_out = False
        busy = True
        while busy:
            print(self._name)
            busy = not self._redis_client.set(self._name, self._id, ex = self._expire, nx = True)
            if busy:
                if time_out:
                    return False
                elif blocking:
                    time_out = not self._redis_client.blpop(self._signal, blocking_timeout) and timeout
                else:
                    print(f"Fail to acquire Lock({self._name}).")
                    return False
        print(f"Acquire Lock({self._name}).")
        if self._lock_renewal_interval is not None:
            self._start_lock_signal_renewer()

        return True


    def release(self):
        """
        Delete Lock
        """
        if self._lock_renewal_interval is not None:
            # 关闭看门狗（不阻塞）
            self._stop_lock_signal_renewer()

        # delete lock
        error = self.unlock_script(client=self._redis_client, keys=(self._name, self._signal), args=(self._id, self._signal_expire))
        if error == 1:
            raise RuntimeError(f"Lock({self._name}) is not captrued or it already expired.")
        elif error:
            raise RuntimeError(f"Unsupported error code({error}) from EXTEND script. ")
    
    
    def _stop_lock_signal_renewer(self):
        self._lock_renewval_event.set()
        self._lock_renewval_thread.join()
        self._lock_renewval_thread = None


    def locked(self):
        # Lock exists
        return self._redis_client.exists(self._name) == 1

    
    def reset(self):
        pass
