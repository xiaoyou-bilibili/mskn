import importlib
import os.path
import logging

# 所有模块
_all_modules = {}

path = "core/module"


# 模块写入。写入后自动加载
def write_module(name: str, code: str) -> bool:
    with open("{}/{}.py".format(path, name), "w", encoding="utf8") as f:
        f.write(code)
    return load_module(name)


# 加载所有模块
def load_all_module():
    for name in os.listdir(path):
        if name.endswith(".py"):
            res = load_module(name[:-3])
            logging.info("加载{} {}".format(name, res))


# 加载特定模块
def load_module(name: str) -> bool:
    # 先判断文件是否存在
    if not os.path.exists("{}/{}.py".format(path, name)):
        print("模块不存在")
        return False
    name = "{}.{}".format(".".join(path.split("/")), name)
    module = importlib.import_module(name)
    importlib.reload(module)
    _all_modules[name] = module
    return True


# 获取模块
def get_module(name: str) -> object:
    if name in _all_modules:
        return _all_modules[name]
    else:
        return None



