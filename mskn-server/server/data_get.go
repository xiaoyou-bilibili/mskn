package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mskn-server/core/proto/mnt"
	"mskn-server/core/tcp"
	"mskn-server/infra/mongo"
	"mskn-server/model"
	"mskn-server/utils"
	"strings"
)

func GetNode(args []byte) []byte {
	var data []string
	json.Unmarshal(args, &data)
	log.Printf("get nodes %v", data)
	nodes := []*model.Node{}
	err := NodeCollection.FindMany(mongo.NewFilter().FieldIn("name", data), &nodes)
	if err != nil {
		log.Printf("get node err %v", err)
	}
	log.Printf("node list %v", nodes)
	return utils.Interface2byte(nodes)
}

func getCode(name string) (*model.Code, error) {
	// 获取代码
	codeName := strings.Split(name, ".")
	if len(codeName) > 2 {
		return nil, errors.New("代码长度有误")
	}
	// 获取代码
	code := &model.Code{}
	if err := CodeCollection.FindOne(mongo.NewFilter().FieldEq("name", codeName[0]), code); err != nil {
		return nil, fmt.Errorf("get code err %v", err)
	}
	log.Printf("get code %v", code.Id)

	return code, nil
}

func PushCode(client *tcp.Client, name string) error {
	code, err := getCode(name)
	if err != nil {
		return err
	}
	// 代码推送
	if err := client.CodePush(code.Name, code.Content); err != nil {
		return fmt.Errorf("push code err %v", err)
	}
	return nil
}

func codeGet(client *tcp.Client, data *mnt.CodeGet) string {
	code, err := getCode(data.Name)
	if err != nil {
		log.Printf("get code err %v", err)
		return ""
	}
	return code.Content
}
