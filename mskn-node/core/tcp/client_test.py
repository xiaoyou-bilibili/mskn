import json
import logging
import time

from core.proto.mn import Mn, EncryptType, MessageType
from core.tcp.client import SocketClientThread

if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s - %(pathname)s[line:%(lineno)d] - %(levelname)s: %(message)s',
                        level=logging.DEBUG)
    print("启动socketserver客户端！")
    t1 = SocketClientThread("127.0.0.1:9000", "xiaoyou")
    t1.start()
    time.sleep(1)
    message = t1.get_message()
    if not message.connect():
        print("连接失败")
    else:
        # 发送代码
        # res = message.code_push("1", "print('hello')")
        # print(message.data_get("test.age", json.dumps({"123":456}).encode()))
        print(message.data_get("demo.age", {"123": 456}))
        # print(message.task_push("test.test", json.dumps({"123": 111}).encode()))

    # server.queue_push(MessageType.code_push, CodePush(_addr,"test", """
    # from core.module.decorate import MnRegister
    #
    # mn = MnRegister("test")
    #
    #
    # @mn.code("test")
    # def handle1(param: dict):
    #     print("测试函数", param)
    #
    # @mn.data("age")
    # def handle2(param: dict):
    #     print("获取数据", param)
    #     return {"age": 123}
    # """))
    # 执行代码
    # print("开始获取数据")
    # # print("获取数据", server.get_data(DataGet("127.0.0.1:9000", "test.age", {"age": 1})))
    # # print("执行函数", server.queue_push(MessageType.task_push, Task(_addr, "test.test", b'{"age": 1}')))
    # print("上报数据", server.queue_push(MessageType.data_push, DataPush("", "123", b'456')))
    # t1.join()
