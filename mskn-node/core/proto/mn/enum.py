from enum import Enum


# 加密方式
class EncryptType(Enum):
    none = 0
    b64 = 1
    aesCbcPack7 = 2


# 消息类型
class MessageType(Enum):
    connect = 1  # 连接请求
    connect_ack = 2  # 同意连接
    connect_refuse = 3  # 拒绝连接
    ping = 4  # ping消息
    ping_ack = 5  # ping回复
    code_push = 6  # 代码推送
    code_ack = 7  # 代码响应
    code_get = 8  # 代码获取
    code_back = 9  # 代码返回
    data_push = 10  # 数据推送
    data_ack = 11  # 收到数据
    data_get = 12  # 数据获取
    data_back = 13  # 数据返回
    task_push = 14  # 执行任务
    task_ack = 15  # 收到任务
