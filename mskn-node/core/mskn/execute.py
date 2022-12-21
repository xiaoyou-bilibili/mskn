import logging

from core.mskn.decorate import get_decorate
from core.mskn.module import get_module
from core.tcp.message import MessageProcess
from threading import Thread
import time


class HandleThread(Thread):
    def __init__(self, message: MessageProcess, handle, name: str, param: dict, ):
        Thread.__init__(self)
        self._message = message
        self._name = name
        self._param = param
        self._handle = handle

    def run(self):
        # 上报进度
        t1 = time.perf_counter()
        self._message.report_status(self._param, 1, 0)
        self._handle(self._message, self._param)
        self._message.report_status(self._param, 2, 100)
        logging.info("func {} execute cost {} ms".format(self._name, (time.perf_counter() - t1) * 1000))


# 执行具体函数
def exec_code(message: MessageProcess, name: str, param: dict) -> bool:
    func = _get_handle(name, "funcs")
    if func is not None:
        # 函数执行是异步的，新开一个线程去执行
        HandleThread(message, func, name, param).start()
        return True
    return False


# 获取数据
def get_data(message: MessageProcess, name: str, param: dict) -> object:
    func = _get_handle(name, "datas")
    if func is not None:
        return func(message, param)
    return None


# 获取处理函数
def _get_handle(name: str, decorate_type: str):
    name = name.split(".")
    # 目前只能有两个，第一个是模块名字，第二个是函数名字
    if len(name) != 2:
        return None
    module, func = name[0], name[1]
    # 判断名字是否在装饰器中
    data = get_decorate(module, decorate_type, func)
    if data is not None:
        handle = data["func"]
        return handle
        # if handle is not None:
        #     # 从模块中获取对于的函数，然后执行其中的内容
        #     return h
        # return getattr(mo, funcs[f]["name"])
    return None
