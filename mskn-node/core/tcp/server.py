import logging
import socket
import time

from core.tcp.message import MessageProcess
from core.proto.mn import Mn, EncryptType
from threading import Thread
from core.mskn.handle import server_message_handle

# 全局消息处理队列
_global_message = {}


def buffer_process(client: socket.socket, message: MessageProcess, server: bool):
    _buffer = bytes()
    _header_size = 4
    while True:
        try:
            # 一次只获取1024字节数据
            data = client.recv(1024)
            if data:
                # 把数据存入缓冲区
                _buffer += data
                while True:
                    # 如果缓冲区小于消息长度，就跳出循环
                    if len(_buffer) < _header_size:
                        logging.info("data less than header data len {}".format(len(_buffer)))
                        break

                    # 获取一下消息的长度
                    _size = int.from_bytes(_buffer[1:4], 'little')
                    # 如果buffer小于消息长度，那么久继续接收数据
                    if len(_buffer) < _size:
                        logging.info(
                            "data less than message size data len {} size len {}".format(len(_buffer), _size))
                        break
                    # 获取实际消息内容，并让客户端进行处理
                    message.handle(_buffer[:_size])
                    # 缓冲区删除已处理的数据
                    _buffer = _buffer[_size:]
        except Exception as e:
            logging.warning("recv data err {}, server {}".format(e, server))
            if server:  # 服务端直接return，客户端则需要不断重试
                return
            else:
                # 把buffer清空
                _buffer = []
                time.sleep(1)


# 单个连接处理线程
class SocketHandleThread(Thread):
    def __init__(self, conn: socket.socket, message: MessageProcess):
        Thread.__init__(self)
        self._conn = conn
        self._message = message

    def run(self):
        buffer_process(self._conn, self._message, True)


# socket 服务线程
class SocketServerThread(Thread):
    def __init__(self, port: int, encrypt_type: EncryptType, secret: str):
        logging.info("start listener 0.0.0.0:{}, secret {}".format(port, secret))
        Thread.__init__(self)
        self._server = socket.socket()
        self._server.bind(('0.0.0.0', port))
        # 最大只能连10个服务
        self._server.listen(10)
        self._encrypt_type = encrypt_type
        self._secret = secret

    def run(self):
        while True:
            # 等待一个连接，此处自动阻塞
            conn, address = self._server.accept()
            # 启动一个新线程来处理消息
            addr = "{}:{}".format(address[0], address[1])
            process = MessageProcess(conn, addr, self._encrypt_type, self._secret)
            SocketHandleThread(conn, process).start()
            # 设置自己的监听器
            server_message_handle(process)
            _global_message[addr] = process
