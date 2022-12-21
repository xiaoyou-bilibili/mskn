from core.mskn.module import load_all_module
from core.tcp import SocketServerThread
from core.proto.mn import EncryptType
import logging

logging.basicConfig(format='%(asctime)s - %(pathname)s[line:%(lineno)d] - %(levelname)s: %(message)s',
                    level=logging.DEBUG)

if __name__ == '__main__':
    # 首先加载全部模块
    load_all_module()
    # 然后启动socket服务
    t1 = SocketServerThread(9000, EncryptType.none, "xiaoyou")
    t1.start()
    t1.join()
