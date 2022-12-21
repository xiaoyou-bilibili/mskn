package tcp

import (
	"fmt"
	"testing"
	"time"
)

func TestClient_Connect(t *testing.T) {
	client, err := NewTcpClient("127.0.0.1:9000")
	if err != nil {
		panic(err)
	}

	err = client.Connect("xiaoyou")
	if err != nil {
		panic(err)
	}
	now := time.Now()
	// 测试一下获取数据
	//err = client.CodePush("1", "测试代码啊啊啊啊")
	data, err := client.DataGet("test.age", []byte(`{"age": 1}`))
	fmt.Println(time.Now().Sub(now))
	fmt.Println("返回数据", string(data), err)

	time.Sleep(time.Second * 5)
}
