package mn

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"mskn-server/core/proto"
	"mskn-server/core/proto/mnt"
	"strings"
)

func NewMn(encryptType EncryptType) *Mn {
	return &Mn{
		encrypt: encryptType,
	}
}

type Mn struct {
	encrypt EncryptType
	msgType uint8  // 1 表示消息格式与加密方式
	msgLen  []byte // 2-4 表示整个数据长度
	data    []byte // 实际数据
}

func (m *Mn) setMessageType(tp MessageType) {
	m.msgType = (uint8(tp) << 3) | (uint8(m.encrypt) & 0b00000111)
}

func (m *Mn) setMessageLen() {
	var size [4]byte
	binary.LittleEndian.PutUint32(size[0:4], uint32(len(m.data)+4))
	m.msgLen = size[0:3]
}

func (m *Mn) getRes() []byte {
	m.setMessageLen()

	// 写入所有字节
	var res bytes.Buffer
	res.WriteByte(m.msgType)
	res.Write(m.msgLen)
	res.Write(m.data)

	return res.Bytes()
}

func (m *Mn) getSha256(secret string) []byte {
	res := sha256.Sum256([]byte(secret))
	return res[:]
}

func (m *Mn) Decode(data []byte) error {
	if len(data) < 4 {
		return errors.New("data len err")
	}
	size := [4]byte{}
	size[0] = data[1]
	size[1] = data[2]
	size[2] = data[3]
	mSize := int(binary.LittleEndian.Uint32(size[:]))
	if len(data) < mSize {
		return errors.New("data size err")
	} else {
		data = data[:mSize]
	}

	m.msgType = data[0]
	m.msgLen = data[1:4]
	if len(data) > 4 {
		m.data = data[4:]
	} else {
		m.data = []byte{}
	}

	return nil
}

func (m *Mn) GetMessageType() MessageType {
	return MessageType(m.msgType >> 3)
}

func (m *Mn) GetEncryptType() EncryptType {
	return EncryptType(m.msgType & 0b00000111)
}

func (m *Mn) GetSize() uint32 {
	return binary.LittleEndian.Uint32(m.msgLen)
}

func (m *Mn) GetData() []byte {
	return m.data
}

func (m *Mn) Connect(secret string) []byte {
	m.setMessageType(MessageTypeConnect)
	m.data = m.getSha256(secret)

	return m.getRes()
}

func (m *Mn) ConnectSecretEq(secret string) bool {
	return bytes.Equal(m.data, m.getSha256(secret))
}

func (m *Mn) ConnectAck() []byte {
	m.setMessageType(MessageTypeConnectAck)
	m.data = []byte{}

	return m.getRes()
}

func (m *Mn) ConnectRefuse() []byte {
	m.setMessageType(MessageTypeConnectRefuse)
	m.data = []byte{}

	return m.getRes()
}

func (m *Mn) Ping() []byte {
	m.setMessageType(MessageTypePing)
	m.data = []byte{}

	return m.getRes()
}

func (m *Mn) PingBack() []byte {
	m.setMessageType(MessageTypePong)
	m.data = []byte{}

	return m.getRes()
}

func (m *Mn) CodePush(addr string, name string, content string) []byte {
	m.setMessageType(MessageTypeCodePush)
	// 把\r\n替换为\n即可
	content = strings.ReplaceAll(content, "\r\n", "\n")
	m.data = proto.Serialize(&mnt.CodePush{
		Addr: addr,
		Name: name,
		Code: content,
	})

	return m.getRes()
}

func (m *Mn) GetCodePush() *mnt.CodePush {
	res := &mnt.CodePush{}
	proto.Deserialize(m.data, res)

	return res
}

func (m *Mn) GetCode() *mnt.CodePush {
	res := &mnt.CodePush{}
	proto.Deserialize(m.data, res)

	return res
}

func (m *Mn) CodeAck() []byte {
	m.setMessageType(MessageTypeCodeAck)
	m.data = []byte{}

	return m.getRes()
}

func (m *Mn) CodeGet(addr string, name string) []byte {
	m.setMessageType(MessageTypeCodeGet)
	m.data = proto.Serialize(&mnt.CodeGet{
		Addr: addr,
		Name: name,
	})

	return m.getRes()
}

func (m *Mn) GetCodeGet() *mnt.CodeGet {
	res := &mnt.CodeGet{}
	proto.Deserialize(m.data, res)

	return res
}

func (m *Mn) CodeBack(addr string, name string, content string) []byte {
	m.setMessageType(MessageTypeCodeBack)
	// 把\r\n替换为\n即可
	content = strings.ReplaceAll(content, "\r\n", "\n")
	m.data = proto.Serialize(&mnt.CodePush{
		Addr: addr,
		Name: name,
		Code: content,
	})

	return m.getRes()
}

func (m *Mn) GetCodeBack() *mnt.CodePush {
	res := &mnt.CodePush{}
	proto.Deserialize(m.data, res)

	return res
}

func (m *Mn) DataPush(addr string, topic string, data []byte) []byte {
	m.setMessageType(MessageTypeDataPush)
	m.data = proto.Serialize(&mnt.DataPush{
		Addr:  addr,
		Topic: topic,
		Data:  data,
	})

	return m.getRes()
}

func (m *Mn) GetDataPush() *mnt.DataPush {
	res := &mnt.DataPush{}
	proto.Deserialize(m.data, res)

	return res
}

func (m *Mn) DataAck() []byte {
	m.setMessageType(MessageTypeDataAck)
	m.data = []byte{}

	return m.getRes()
}

func (m *Mn) DataGet(addr string, name string, params []byte) []byte {
	m.setMessageType(MessageTypeDataGet)
	m.data = proto.Serialize(&mnt.DataGet{
		Addr:  addr,
		Name:  name,
		Param: params,
	})

	return m.getRes()
}

func (m *Mn) GetDataGet() *mnt.DataGet {
	res := &mnt.DataGet{}
	proto.Deserialize(m.data, res)

	return res
}

func (m *Mn) DataBack(addr string, name string, data []byte) []byte {
	m.setMessageType(MessageTypeDataBack)
	m.data = proto.Serialize(&mnt.DataBack{
		Addr: addr,
		Name: name,
		Data: data,
	})

	return m.getRes()
}

func (m *Mn) GetDataBack() *mnt.DataBack {
	res := &mnt.DataBack{}
	proto.Deserialize(m.data, res)

	return res
}

func (m *Mn) TaskPush(addr string, name string, param []byte) []byte {
	m.setMessageType(MessageTypeTaskPush)
	m.data = proto.Serialize(&mnt.Task{
		Addr:  addr,
		Name:  name,
		Param: param,
	})

	return m.getRes()
}

func (m *Mn) GetTaskPush() *mnt.Task {
	res := &mnt.Task{}
	proto.Deserialize(m.data, res)

	return res
}

func (m *Mn) TaskAck() []byte {
	m.setMessageType(MessageTypeTaskAck)
	m.data = []byte{}

	return m.getRes()
}
