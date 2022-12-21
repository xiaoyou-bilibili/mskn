from thrift.transport import TTransport
from thrift.protocol import TBinaryProtocol


# 序列化文件
def serialize(th_obj):
    """ Serialize.
    """
    buf = TTransport.TMemoryBuffer()
    prot = TBinaryProtocol.TBinaryProtocol(buf)
    th_obj.write(prot)
    return buf.getvalue()


# 反序列化文件
def deserialize(val, th_obj_type):
    """ Deserialize.
    """
    th_obj = th_obj_type()
    buf = TTransport.TMemoryBuffer(val)
    prot = TBinaryProtocol.TBinaryProtocol(buf)
    th_obj.read(prot)
    return th_obj
