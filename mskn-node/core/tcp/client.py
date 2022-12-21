import logging
import time
from threading import Thread

import socket

from core.tcp.message import MessageProcess
from core.proto.mnt.ttypes import DataPush, DataGet
from core.proto.mn import EncryptType
from core.tcp.server import buffer_process


class SocketClientThread(Thread):
    def __init__(self, addr: str, encrypt_type: EncryptType, secret: str):
        Thread.__init__(self)
        _addr_info = addr.split(":")
        if len(_addr_info) != 2:
            raise "addr err"
        self._addr = addr
        self._client = socket.socket()
        self._client.connect((_addr_info[0], int(_addr_info[1])))
        self._message = MessageProcess(self._client, self._addr, encrypt_type, secret)

        def reconnect():
            logging.info("reconnect...")
            self._client = socket.socket()
            self._client.connect((_addr_info[0], int(_addr_info[1])))
            return self._client

        self._message.set_reconnect(reconnect)

    def get_message(self) -> MessageProcess:
        return self._message

    def get_addr(self) -> str:
        return self._addr

    def run(self):
        buffer_process(self._client, self._message, False)
