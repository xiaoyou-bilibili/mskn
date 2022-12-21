import json
import logging

from core.tcp.message import MessageProcess, MessageType
from core.proto.mnt.ttypes import DataPush, CodeGet
from core.proto.mn import EncryptType
from core.tcp.client import SocketClientThread
from random import choice

_global_node = {}
_global_client = {}


# 获取节点并发送数据
def _get_node_send(message: MessageProcess, nodes: list[str], random: bool) -> list[str]:
    # 需要获取的新节点，避免重复获取
    new_node_get = []
    # 先判断一下节点是否存在
    for node in nodes.copy():
        if node not in _global_node:
            new_node_get.append(node)
    if len(new_node_get) > 0:
        logging.info("need get new node {}".format(new_node_get))
        # 先获取节点信息
        res = message.data_get("core.server.getNode", new_node_get)
        if res is None:
            logging.error("data is nil")
            return []
        # 遍历所有节点
        for info in json.loads(res.data):
            if "name" in info:
                _global_node[info["name"]] = info
    # 判断是否为随机选取
    choose_node = []
    if random:
        choose_node.append(choice(nodes))
    else:
        choose_node = nodes
    # 根据选中的这些节点分别建立连接
    logging.info("choose node {}".format(choose_node))
    return choose_node


# 获取client
def _get_client(message: MessageProcess, nodes: list[str], random: bool) -> list[SocketClientThread]:
    # 先获取节点信息
    client_nodes = _get_node_send(message, nodes, random)
    # 遍历
    client_list = []
    for node in client_nodes:
        if node not in _global_client:
            logging.info("node {} not in client".format(node))
            info = _global_node[node]
            logging.info("node info {}".format(info))
            client = SocketClientThread(info["addr"], EncryptType(info["encrypt_type"]), info["secret"])
            client.start()
            # 进行连接
            client_message = client.get_message()

            # 这里还需要处理数据推送过来时
            def _data_push(process: MessageProcess, tp: MessageType, data: DataPush):
                logging.info("sub client push data{}".format(data.topic))
                # 有数据过来时，直接传递给父节点，这里把地址也透传过来
                message.data_push(data.topic, message.bytes_2_dict(data.data), data.addr)

            # 监听获取代码请求
            def _code_get(process: MessageProcess, tp: MessageType, data: CodeGet):
                # 从父节点获取代码然后返回
                logging.info("sub client get code {}".format(data.name))
                # 从本服务获取代码
                code = message.code_get(data.name)
                # 推送获取到的代码
                process.code_back(code.name, code.code)

            client_message.add_listener(MessageType.data_push, _data_push)
            client_message.add_listener(MessageType.code_get, _code_get)

            if client_message.connect():
                logging.info("connect {} success".format(node))
                _global_client[node] = client
                client_list.append(client)
            else:
                logging.info("connect {} err".format(node))
        else:
            client_list.append(_global_client[node])

    return client_list


# 推送任务
def push_task(message: MessageProcess, nodes: list[str], name: str, args: dict, random=True):
    logging.info("push task nodes {}, name {}".format(nodes, name))
    # 获取连接
    for client in _get_client(message, nodes, random):
        logging.info("client is {}".format(client.get_addr()))
        client.get_message().task_push(name, args)


# 获取数据
def get_data(message: MessageProcess, node: str, name: str, args: dict):
    logging.info("get data node {}, name {}".format(node, name))
    # 获取连接
    for client in _get_client(message, [node], False):
        logging.info("client is {}".format(client.get_addr()))
        return client.get_message().data_get(name, args)

