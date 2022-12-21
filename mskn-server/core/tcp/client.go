package tcp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"mskn-server/core/proto/mn"
	"mskn-server/core/proto/mnt"
	"net"
	"time"
)

const retryTime = 3 // 默认值重试3次

type codeGetHandle = func(client *Client, data *mnt.CodeGet) string
type dataPushHandle = func(client *Client, data *mnt.DataPush)
type dataGetHandle = func(client *Client, data *mnt.DataGet) []byte

func NewTcpClient(id string, addr string, encryptType mn.EncryptType) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	client := &Client{
		id:              id,
		client:          conn,
		mn:              mn.NewMn(encryptType),
		encryptType:     encryptType,
		connRes:         make(chan bool),
		addr:            addr,
		dataBackQueue:   map[string]chan []byte{},
		messageBackChan: map[mn.MessageType]chan bool{},
		pingExit:        make(chan bool),
	}
	// 单独开一个线程去监听消息
	go client.startListener()
	return client, nil
}

type Client struct {
	id              string
	client          net.Conn
	mn              *mn.Mn
	addr            string
	secret          string
	encryptType     mn.EncryptType
	connRes         chan bool // 连接结果
	dataBackQueue   map[string]chan []byte
	messageBackChan map[mn.MessageType]chan bool // 是否有消息响应
	codeGetHandle   []codeGetHandle              // 代码获取监听器
	dataPushHandle  []dataPushHandle             // 数据推送监听器
	dataGetHandle   []dataGetHandle              // 数据推送监听器
	pingExit        chan bool
}

func (c *Client) GetNodeId() string {
	return c.id
}

func (c *Client) startListener() {
	defer c.client.Close() // 关闭TCP连接
	// 一分钟发送一个ping数据，保持连接
	go func() {
		for {
			select {
			case <-c.pingExit:
				log.Printf("get exit ping")
				return
			case <-time.After(time.Minute):
				c.client.Write(c.mn.Ping())
			}
		}
	}()
	var buffer []byte
	for {
		// 一次性读1024字节
		tmp := [1024]byte{}
		n, err := c.client.Read(tmp[:])
		if err != nil {
			fmt.Println(err)
			if err != io.EOF {
				log.Printf("recv failed, err: %v", err)
				// 退出旧循环
				c.pingExit <- true
				return
			}
			continue
		}
		// 把数据写到buffer中
		buffer = append(buffer, tmp[:n]...)
		// 判断一下buffer长度，小于4就先跳过本轮循环
		if len(buffer) < 4 {
			continue
		}
		// 获取消息长度
		size := [4]byte{}
		size[0] = buffer[1]
		size[1] = buffer[2]
		size[2] = buffer[3]
		mSize := int(binary.LittleEndian.Uint32(size[:]))
		// 如果长度不一样，也不管
		if len(buffer) < mSize {
			continue
		}
		// 对数据进行拷贝，避免丢失
		data := make([]byte, mSize)
		copy(data, buffer[:mSize])
		// 对消息进行解码，这里新开一个线程来处理
		go func() {
			tmpMn := mn.NewMn(c.encryptType)
			if err := tmpMn.Decode(data); err != nil {
				log.Printf("decode data err %v", err)
				return
			}
			log.Printf("消息类型 %v", tmpMn.GetMessageType())
			switch tmpMn.GetMessageType() {
			case mn.MessageTypeConnectAck:
				c.connRes <- true
			case mn.MessageTypePing:
				c.write(c.mn.PingBack())
			case mn.MessageTypeConnectRefuse:
				c.connRes <- false
			case mn.MessageTypeCodeAck, mn.MessageTypeTaskAck, mn.MessageTypeDataAck:
				if ch, ok := c.messageBackChan[tmpMn.GetMessageType()]; ok {
					ch <- true
				}
			case mn.MessageTypeCodeGet:
				// 遍历所有监听器，然后直接把消息给发送过去
				for _, handle := range c.codeGetHandle {
					if data := handle(c, tmpMn.GetCodeGet()); data != "" {
						c.write(c.mn.CodeBack(c.addr, tmpMn.GetDataGet().Name, data))
					}
				}
			case mn.MessageTypeDataGet:
				// 遍历所有监听器，然后直接把消息给发送过去
				for _, handle := range c.dataGetHandle {
					if data := handle(c, tmpMn.GetDataGet()); data != nil {
						fmt.Printf("return data %s \n", string(data))
						c.write(c.mn.DataBack(c.addr, tmpMn.GetDataGet().Name, data))
					}
				}
			case mn.MessageTypeDataPush:
				c.write(c.mn.DataAck())
				// 变量所有的数据监听器
				for _, handle := range c.dataPushHandle {
					handle(c, tmpMn.GetDataPush())
				}
			case mn.MessageTypeDataBack: // 获取数据返回
				back := tmpMn.GetDataBack()
				if ch, ok := c.dataBackQueue[back.Name]; ok {
					ch <- back.Data
				}
			}
		}()
		buffer = buffer[mSize:]
	}
}

// Reconnect 重连
func (c *Client) Reconnect() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}
	c.client = conn
	go c.startListener()
	err = c.Connect(c.secret)
	if err != nil {
		return err
	}
	log.Printf("reconnect success")
	return nil
}

// 写入数据
func (c *Client) write(data []byte) error {
	_, err := c.client.Write(data)
	if err != nil {
		log.Printf("write dta err %v", err)
		// 尝试重连
		if err := c.Reconnect(); err != nil {
			return fmt.Errorf("reconnect err %v", err)
		}
		// 重新发送
		_, err = c.client.Write(data)
	}
	return err
}

// 等待响应才算完成
func (c *Client) waitAck(data []byte, tp mn.MessageType, retry int) error {
	// 先写入数据
	if err := c.write(data); err != nil {
		return err
	}
	// 创建相关管道
	if _, ok := c.messageBackChan[tp]; !ok {
		c.messageBackChan[tp] = make(chan bool)
	}
	// 等待对应通道响应，最长响应时间为1s
	select {
	case <-c.messageBackChan[tp]:
		return nil
	case <-time.After(time.Second):
		retry--
		log.Printf("time out retey %d", retry)
		if retry < 0 {
			return errors.New("响应超时")
		} else {
			// 开始下次重试
			return c.waitAck(data, tp, retry)
		}
	}
}

// Connect 启动连接
func (c *Client) Connect(secret string) error {
	c.secret = secret
	if err := c.write(c.mn.Connect(secret)); err != nil {
		return err
	}
	// 等待连接结果，超时时间为1s
	select {
	case res := <-c.connRes:
		if res {
			return nil
		} else {
			c.client.Close()
			return errors.New("连接被拒绝，请检查秘钥是否正确")
		}
	case <-time.After(time.Second):
		return errors.New("连接超时")
	}
}

// CodePush 代码推送
func (c *Client) CodePush(name string, content string) error {
	return c.waitAck(c.mn.CodePush(c.addr, name, content), mn.MessageTypeCodeAck, retryTime)
}

// AddCodeGetHandle 添加一个代码获取监听器
func (c *Client) AddCodeGetHandle(handle codeGetHandle) {
	c.codeGetHandle = append(c.codeGetHandle, handle)
}

// AddDataPushHandle 添加一个数据推送监听器
func (c *Client) AddDataPushHandle(handle dataPushHandle) {
	c.dataPushHandle = append(c.dataPushHandle, handle)
}

// AddDataGetHandle 添加一个数据获取监听器
func (c *Client) AddDataGetHandle(handle dataGetHandle) {
	c.dataGetHandle = append(c.dataGetHandle, handle)
}

// DataGet 获取数据
func (c *Client) DataGet(name string, param []byte) ([]byte, error) {
	// 先判断队列是否存在。不存在就新建
	if _, ok := c.dataBackQueue[name]; !ok {
		c.dataBackQueue[name] = make(chan []byte)
	}
	// 发送数据获取命令
	if err := c.write(c.mn.DataGet(c.addr, name, param)); err != nil {
		return nil, err
	}
	// 监听管道获取数据，5秒超时
	select {
	case data := <-c.dataBackQueue[name]:
		return data, nil
	case <-time.After(time.Second * 5):
		return nil, errors.New("返回超时")
	}
}

// DataPush 数据推送
func (c *Client) DataPush(topic string, data []byte) error {
	return c.waitAck(c.mn.DataPush(c.addr, topic, data), mn.MessageTypeDataAck, retryTime)
}

// TaskPush 任务推送
func (c *Client) TaskPush(name string, param []byte) error {
	return c.waitAck(c.mn.TaskPush(c.addr, name, param), mn.MessageTypeTaskAck, retryTime)
}
