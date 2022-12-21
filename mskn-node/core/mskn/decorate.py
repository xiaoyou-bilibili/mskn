_all_decorate = {}


def get_decorate(module: str, decorate_type: str, fun_name: str) -> object:
    if module in _all_decorate:
        funcs = _all_decorate[module][decorate_type]
        # 再判断对应函数是否存在
        if fun_name in funcs:
            return funcs[fun_name]
    return None


# 自定义函数装饰器
class MnRegister:
    # 注册时传入
    def __init__(self, name: str):
        self.module = {
            "name": name,  # 模块名字
            "funcs": {},  # 模块里面所有的函数
            "datas": {}  # 模块里面所有获取数据的接口
        }

    # 执行代码，无返回
    def code(self, name):
        def wrapper(func):
            # 模块加载时自动把相关信息注入到全局变量中
            self.module["funcs"][name] = {
                "func": func
            }
            _all_decorate[self.module["name"]] = self.module

            def decorate(*args, **kw):
                # 当函数被调用时，会触发该方法，可以在这里执行一些操作
                return func(*args, **kw)

            return decorate

        return wrapper

    # 获取数据，有交互
    def data(self, name):
        def wrapper(func):
            self.module["datas"][name] = {
                "name": func.__name__,
                "module": func.__module__,
                "func": func
            }
            _all_decorate[self.module["name"]] = self.module

            def decorate(*args, **kw):
                return func(*args, **kw)

            return decorate

        return wrapper
