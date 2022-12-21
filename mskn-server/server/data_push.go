package server

import (
	"github.com/bytedance/sonic"
	"log"
	"mskn-server/core/proto/mnt"
	"mskn-server/core/tcp"
	"mskn-server/infra/mongo"
	"time"
)

func UpdateStatus(client *tcp.Client, data *mnt.DataPush) {
	// 获取数据
	root, err := sonic.Get(data.Data)
	if err != nil {
		log.Printf("marsha data err %v", err)
		return
	}
	log.Printf("update info %s", string(data.Data))
	update := mongo.NewUpdate()
	recordId, err := root.Get("r").String()
	if err != nil {
		log.Printf("no record found")
		return
	}
	status, err := root.Get("s").Int64()
	if err == nil {
		if status == 2 {
			update.SetField("finish", time.Now())
		} else if status == 1 {
			update.SetField("create", time.Now())
		}
		update.SetField("status", status)
	}
	progress, err := root.Get("p").Int64()
	if err == nil {
		update.SetField("progress", progress)
	}
	err = RecordCollection.UpdateOne(mongo.NewFilter().FieldEq("_id", recordId), update)
	if err != nil {
		log.Printf("update record err %v", err)
	}
}
