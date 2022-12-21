import json
import logging
import socket
import time

from core.proto.mn import Mn, EncryptType, MessageType
from core.proto.mnt.ttypes import CodePush
from multiprocessing import Queue
from threading import Thread


class MessageProcess:
    def __init__(self, conn: socket.socket, addr: str, encrypt_type: EncryptType, secret: str):
        self._reconnect_handle = None
        self._encoder = 'utf-8'
        self._conn = conn
        self._addr = addr
        self._secret = secret
        self._encrypt_type = encrypt_type
        self._mn = Mn(encrypt_type)
        self._conn_res = Queue(1)
        # 需要监听的消息类型,避免申请太多队列占用空间
        self._listener = {
            MessageType.code_push: [],
            MessageType.code_get: [],
            MessageType.data_push: [],
            MessageType.data_get: [],
            MessageType.task_push: [],
            MessageType.code_get: [],
        }
        # 消息响应队列
        self._message_back = {
            MessageType.code_ack: Queue(1),
            MessageType.data_ack: Queue(1),
            MessageType.task_ack: Queue(1),
            MessageType.data_back: Queue(1),
            MessageType.code_back: Queue(1),
        }
        # 所有队列
        self._queue = {}

    def get_encrypt_type(self) -> EncryptType:
        return self._encrypt_type

    def get_secret(self) -> str:
        return self._secret

    # 添加一个监听器
    def add_listener(self, message_type: MessageType, handle):
        self._listener[message_type].append(handle)

    # 触发监听器
    def handle_listener(self, message_type: MessageType, data: object):
        # 这里需要异步执行，避免阻塞主线程
        for handle in self._listener[message_type]:
            handle(self, message_type, data)

    # 建立连接
    def connect(self) -> bool:
        self.send(self._mn.connect(self._secret))
        # 获取连接结果,超时设置为2s
        return self._conn_res.get(True, 2)

    # 等待队列响应，同时进行重试，顺便计算一下获取结果需要的时间
    def _wait_ack(self, data: bytes, tp: MessageType, retry: int) -> object:
        t1 = time.perf_counter()
        # 先把队列清空，然后发送数据并等待响应
        # if not self._message_back[tp].empty():
        #     self._message_back[tp].get(False)
        self.send(data)
        try:
            # 超时时间为1s，然后如果获取到数据就直接响应
            res = self._message_back[tp].get(True, 1)
        except Exception as e:
            logging.warning("message type {} time out retry {}, err {}".format(tp, retry, e))
            retry -= 1
            # 重试次数达到设置值就返回，否则就继续重试
            if retry < 0:
                logging.error("get data err {}".format(e))
                res = None
            else:
                res = self._wait_ack(data, tp, retry)
        logging.info("wait ack cost {} ms".format((time.perf_counter() - t1) * 1000))
        return res

    def object_2_bytes(self, data: object) -> bytes:
        return json.dumps(data, ensure_ascii=False).encode(self._encoder)

    def bytes_2_dict(self, data: bytes) -> dict:
        return json.loads(data)

    # ping响应
    def ping(self):
        self.send(self._mn.ping())

    # 代码推送
    def code_push(self, name: str, content: str) -> object:
        return self._wait_ack(self._mn.code_push(self._addr, name, content), MessageType.code_ack, 3)

    # 获取代码
    def code_get(self, name: str) -> CodePush:
        return self._wait_ack(self._mn.code_get(self._addr, name), MessageType.code_back, 3)

    # 代码返回
    def code_back(self, name: str, content: str) -> object:
        return self.send(self._mn.code_back(self._addr, name, content))

    # 数据推送
    def data_push(self, topic: str, data: dict, addr=None) -> object:
        if addr is None:
            addr = self._addr
        return self._wait_ack(self._mn.data_push(addr, topic, self.object_2_bytes(data)), MessageType.data_ack, 3)

    # 任务推送
    def task_push(self, name: str, param: dict) -> object:
        return self._wait_ack(self._mn.task_push(self._addr, name, self.object_2_bytes(param)), MessageType.task_ack, 3)

    # 获取数据
    def data_get(self, name: str, param: object) -> object:
        return self._wait_ack(self._mn.data_get(self._addr, name, self.object_2_bytes(param)), MessageType.data_back, 3)

    # 返回数据
    def data_back(self, name: str, content: bytes):
        self.send(self._mn.data_back(self._addr, name, content))

    # 任务状态上报
    def report_status(self, param: dict, status: int, progress: int):
        if "_rid" in param:
            self.data_push("core.server.updateRecord", {"r": param["_rid"], "s": status, "p": progress})

    # 消息发送时的处理
    def handle(self, data: bytes):
        # 每次来新数据都单开一个线程处理
        HandleThread(self, data).start()

    # 连接结果
    def connect_res(self, result: bool):
        self._conn_res.put(result)

    # 发送数据
    def send(self, data: bytes):
        try:
            self._conn.send(data)
        except Exception as e:
            logging.error("send data err {} socket res {}".format(e, getattr(self._conn, '_closed')))
            # 触发重连
            if self._reconnect_handle is not None:
                self._reconnect()

    def set_reconnect(self, handle):
        self._reconnect_handle = handle

    def _reconnect(self):
        self._conn = self._reconnect_handle()

    # 关闭连接
    def close(self):
        self._conn.close()

    def message_back(self, mt: MessageType, data: object):
        self._message_back[mt].put(data)


#  消息处理线程
class HandleThread(Thread):
    def __init__(self, message: MessageProcess, data: bytes):
        Thread.__init__(self)
        self._message = message
        self._data = data
        self._mn = Mn(message.get_encrypt_type())

    def run(self):
        res = self._mn.decode(self._data)
        mt = res.get_message_type()
        logging.info("message type {}".format(mt))
        # print("消息格式", mt)
        # print("消息内容", res.get_data())
        # 根据不同的消息类型走不同的处理逻辑
        if mt == MessageType.connect:
            # 校验密码是否正确
            if res.connect_secret_eq(self._message.get_secret()):
                self._message.send(self._mn.connect_ack())
            else:
                self._message.send(self._mn.connect_refuse())
                self._message.close()
        elif mt == MessageType.connect_ack:
            self._message.connect_res(True)
        elif mt == MessageType.connect_refuse:
            self._message.connect_res(False)
        elif mt == MessageType.ping:
            self._message.send(self._mn.ping_ack())
        elif mt == MessageType.code_push:
            _data = self._mn.get_code_push()
            self._message.send(self._mn.code_ack())
            self._message.handle_listener(MessageType.code_push, _data)
        elif mt == MessageType.code_ack or mt == MessageType.data_ack or mt == MessageType.task_ack:
            self._message.message_back(mt, True)
        elif mt == MessageType.code_get:
            self._message.handle_listener(MessageType.code_get, self._mn.get_code_get())
        elif mt == MessageType.code_back:
            self._message.message_back(MessageType.code_back, self._mn.get_code_back())
        elif mt == MessageType.data_push:
            _data = self._mn.get_data_push()
            self._message.send(self._mn.data_ack())
            self._message.handle_listener(MessageType.data_push, _data)
        elif mt == MessageType.data_get:
            self._message.handle_listener(MessageType.data_get, self._mn.get_data_get())
        elif mt == MessageType.data_back:
            self._message.message_back(MessageType.data_back, self._mn.get_data_back())
        elif mt == MessageType.task_push:
            _data = self._mn.get_task_push()
            self._message.send(self._mn.task_ack())
            self._message.handle_listener(MessageType.task_push, _data)
