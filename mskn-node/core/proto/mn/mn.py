import hashlib, json
import logging

import core.proto as proto
from core.proto.mn import EncryptType, MessageType, MnFormatErr, MnSizeErr
import core.proto.mnt.ttypes as mnt


def print_bytes(data: bytes):
    print("---")
    print(" ".join([bin(int(i)) for i in data]))
    print(" ".join([str(int(i)) for i in data]))


class Mn:
    _encrypt_type = EncryptType.none
    _byteorder = 'little'
    _encoder = "utf-8"
    _message_type = bytes()  # 1 表示消息格式与加密方式
    _message_len = bytes()  # 2-4 表示整个数据长度
    _data = bytes()  # 后面就是具体的数据了

    def __init__(self, encrypt_type: EncryptType):
        self._encrypt_type = encrypt_type

    def _set_message_type(self, message_type: MessageType):
        # 消息类型左移三位，加密方式为最后一位
        self._message_type = ((message_type.value << 3) | (self._encrypt_type.value & 0b00000111)).to_bytes(1,
                                                                                                            self._byteorder)

    def _set_byte_size(self):
        # 消息长度为数据长度加上6个固定长度
        self._message_len = (len(self._data) + 4).to_bytes(3, self._byteorder)

    def _get_res(self) -> bytes:
        # 先设置好对应的值
        self._set_byte_size()
        # 初始化一个可变数据，然后拼接各个字节
        res = bytearray()
        res += self._message_type
        res += self._message_len
        res += self._data
        return res

    def get_json_bytes(self, data: dict) -> bytes:
        return json.dumps(data, ensure_ascii=False).encode(self._encoder)

    def _get_sha256(self, secret: str) -> bytes:
        return hashlib.sha256(secret.encode(self._encoder)).digest()

    # 对消息进行解码
    def decode(self, byte_data: bytes):
        if len(byte_data) < 4:
            raise MnFormatErr
        size = byte_data[1:4]
        # 判断size是否一致
        m_size = int.from_bytes(size, self._byteorder)
        logging.info("data size {} msg size {}".format(len(byte_data), m_size))
        if len(byte_data) < m_size:
            raise MnSizeErr
        else:
            # 如果大于，那么就把后面的全部截断
            byte_data = byte_data[:m_size]
        # 获取各个字节数据
        self._message_type = byte_data[0]
        self._message_len = size
        self._data = byte_data[4:]
        return self

    def get_message_type(self) -> MessageType:
        return MessageType(self._message_type >> 3)

    def get_encrypt_type(self) -> EncryptType:
        return EncryptType(self._message_type & 0b00000111)

    def get_size(self) -> int:
        return int.from_bytes(self._message_len, self._byteorder)

    def get_data(self) -> bytes:
        return self._data

    # 连接请求
    def connect(self, secret: str) -> bytes:
        self._set_message_type(MessageType.connect)
        # 连接的内容需要对密码进行sha256加密
        self._data = self._get_sha256(secret)
        return self._get_res()

    # 判断两个秘钥加密是否相等
    def connect_secret_eq(self, secret: str) -> bool:
        return self._data == self._get_sha256(secret)

    # 同意连接
    def connect_ack(self) -> bytes:
        self._set_message_type(MessageType.connect_ack)
        self._data = bytes()
        return self._get_res()

    # 拒绝连接
    def connect_refuse(self) -> bytes:
        self._set_message_type(MessageType.connect_refuse)
        self._data = bytes()
        return self._get_res()

    # ping消息
    def ping(self) -> bytes:
        self._set_message_type(MessageType.ping)
        self._data = bytes()
        return self._get_res()

    # ping回复
    def ping_ack(self) -> bytes:
        self._set_message_type(MessageType.ping_ack)
        self._data = bytes()
        return self._get_res()

    # 代码推送
    def code_push(self, addr: str, code_name: str, code_content: str) -> bytes:
        self._set_message_type(MessageType.code_push)
        self._data = proto.serialize(mnt.CodePush(addr, code_name, code_content))
        return self._get_res()

    # 获取代码
    def get_code_push(self) -> mnt.CodePush:
        return proto.deserialize(self._data, mnt.CodePush)

    # 收到代码
    def code_ack(self) -> bytes:
        self._set_message_type(MessageType.code_ack)
        self._data = bytes()
        return self._get_res()

    # 获取代码
    def code_get(self, addr: str, code_name: str) -> bytes:
        self._set_message_type(MessageType.code_get)
        self._data = proto.serialize(mnt.CodeGet(addr, code_name))
        return self._get_res()

    # 获取代码
    def get_code_get(self) -> mnt.CodeGet:
        return proto.deserialize(self._data, mnt.CodeGet)

    # 返回代码
    def code_back(self, addr: str, code_name: str, code_content: str) -> bytes:
        self._set_message_type(MessageType.code_back)
        self._data = proto.serialize(mnt.CodePush(addr, code_name, code_content))
        return self._get_res()

    # 获取返回的代码
    def get_code_back(self) -> mnt.CodePush:
        return proto.deserialize(self._data, mnt.CodePush)

    # 数据推送
    def data_push(self, addr: str, data_topic: str, data: bytes) -> bytes:
        self._set_message_type(MessageType.data_push)
        self._data = proto.serialize(mnt.DataPush(addr, data_topic, data))
        return self._get_res()

    # 获取推送的数据
    def get_data_push(self) -> mnt.DataPush:
        return proto.deserialize(self._data, mnt.DataPush)

    # 收到数据
    def data_ack(self) -> bytes:
        self._set_message_type(MessageType.data_ack)
        self._data = bytes()
        return self._get_res()

    # 数据获取
    def data_get(self, addr: str, name: str, args: bytes) -> bytes:
        self._set_message_type(MessageType.data_get)
        self._data = proto.serialize(mnt.DataGet(addr, name, args))
        return self._get_res()

    # 获取数据获取的数据
    def get_data_get(self) -> mnt.DataGet:
        return proto.deserialize(self._data, mnt.DataGet)

    # 数据返回
    def data_back(self, addr: str, name: str, data: bytes) -> bytes:
        self._set_message_type(MessageType.data_back)
        self._data = proto.serialize(mnt.DataBack(addr, name, data))
        return self._get_res()

    # 获取数据获取的数据
    def get_data_back(self) -> mnt.DataBack:
        return proto.deserialize(self._data, mnt.DataBack)

    #  执行任务
    def task_push(self, addr: str, name: str, args: bytes) -> bytes:
        self._set_message_type(MessageType.task_push)
        self._data = proto.serialize(mnt.Task(addr, name, args))
        return self._get_res()

    # 获取任务
    def get_task_push(self) -> bytes:
        return proto.deserialize(self._data, mnt.Task)

    # 收到数据
    def task_ack(self) -> bytes:
        self._set_message_type(MessageType.task_ack)
        self._data = bytes()
        return self._get_res()
