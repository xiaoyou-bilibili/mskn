package server

import (
	"log"
	"mskn-server/core/proto/mnt"
	"mskn-server/core/tcp"
	"mskn-server/infra/mongo"
	"mskn-server/model"
	"strings"
	"time"
)

var clientList = map[string]*tcp.Client{}

// MessageHandle 数据推送监听
func dataPush(client *tcp.Client, data *mnt.DataPush) {
	log.Printf("get data node id %s", client.GetNodeId())
	if strings.Contains(data.Topic, "core.") {
		log.Printf("data container core \n")
		switch data.Topic {
		case "core.server.updateRecord":
			UpdateStatus(client, data)
		}
		return
	}
	info := model.DataPush{
		NodeId: client.GetNodeId(),
		Addr:   data.Addr,
		Topic:  data.Topic,
		Data:   string(data.Data),
		Update: time.Now(),
	}
	res := DataCollection.InsertOne(info)
	log.Printf("insert data res %v \n", res)
}

// 数据获取监听
func dataGet(client *tcp.Client, data *mnt.DataGet) []byte {
	switch data.Name {
	case "core.server.getNode":
		return GetNode(data.Param)
	}

	return nil
}

// GetClient 获取一个连接
func GetClient(nodeId string) (*tcp.Client, error) {
	if client, ok := clientList[nodeId]; ok {
		return client, nil
	}
	// 创建一个新client
	node := &model.Node{}
	if err := NodeCollection.FindOne(mongo.NewFilter().Id(nodeId), node); err != nil {
		return nil, err
	}
	log.Printf("get node info %v", node)
	client, err := tcp.NewTcpClient(nodeId, node.Addr, node.EncryptType)
	if err != nil {
		return nil, err
	}
	// 使用密码进行连接
	if err := client.Connect(node.Secret); err != nil {
		return nil, err
	}
	clientList[nodeId] = client
	// 设置监听器
	client.AddDataPushHandle(dataPush)
	client.AddDataGetHandle(dataGet)
	client.AddCodeGetHandle(codeGet)

	return client, nil
}
