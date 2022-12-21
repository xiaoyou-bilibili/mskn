import json
import logging

from core.mskn.module import write_module
from core.mskn.execute import get_data, exec_code
from core.proto.mnt.ttypes import CodePush, DataGet, Task
from core.proto.mn import MessageType
from core.tcp.message import MessageProcess


# 代码推送逻辑
def _code_push(process: MessageProcess, tp: MessageType, data: CodePush):
    logging.info("[listener] code push {}".format(data))
    write_module(data.name, data.code)


# 代码执行逻辑
def _task_push(process: MessageProcess, tp: MessageType, data: Task):
    logging.info("[listener] task push {}".format(data))
    res = exec_code(process, data.name, json.loads(data.param))
    if not res:
        logging.info("[listener] need get code {}".format(data))
        _code = process.code_get(data.name)
        write_module(str(_code.name).split(".")[0], _code.code)
        res = exec_code(process, data.name, json.loads(data.param))
    logging.info("[listener] task exec res {}".format(res))


# 服务消息处理
def server_message_handle(message: MessageProcess):
    # 获取数据
    def data_get(process: MessageProcess, tp: MessageType, data: DataGet):
        logging.info("[listener] data get req {}".format(data))
        res = get_data(process, data.name, json.loads(data.param))
        message.data_back(data.name, json.dumps(res).encode())
        logging.info("[listener] data get res {}".format(res))
        # 代码推送

    message.add_listener(MessageType.code_push, _code_push)
    message.add_listener(MessageType.task_push, _task_push)
    message.add_listener(MessageType.data_get, data_get)
