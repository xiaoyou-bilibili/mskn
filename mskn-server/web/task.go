package web

import (
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"mskn-server/infra/mongo"
	"mskn-server/model"
	"mskn-server/server"
	"mskn-server/utils"
)

func setRecord(taskId, name, args string) []byte {
	root, err := sonic.Get([]byte(args))
	if err != nil {
		log.Printf("get data err")
		return nil
	}
	id := primitive.NewObjectID()
	root.Set("_rid", ast.NewString(id.Hex()))
	tmp, err := root.MarshalJSON()
	if err != nil {
		log.Printf("marshal data err")
	}
	if err := server.RecordCollection.InsertOne(&model.Record{
		Id:       id.Hex(),
		TaskId:   taskId,
		FuncName: name,
		Args:     args,
	}); err != nil {
		log.Printf("insert data err %v", err)
	}
	return tmp
}

// ExecuteTask 执行任务
func ExecuteTask(c *gin.Context) {
	info := &model.ExecuteTask{}
	c.Bind(info)
	for _, node := range info.NodeId {
		log.Printf("start node %s", node)
		// 获取连接id
		client, err := server.GetClient(node)
		if err != nil {
			noOk(c, err)
			return
		}
		// 获取任务
		task := &model.Task{}
		err = server.TaskCollection.FindOne(mongo.NewFilter().Id(info.TaskId), task)
		if err != nil {
			noOk(c, err)
			return
		}
		log.Printf("get task info %v", task)
		// 判断是否为强制推送
		if info.ForcePush {
			if err = server.PushCode(client, task.TaskName); err != nil {
				noOk(c, err)
				return
			}
		}
		// 任务推送
		if err = client.TaskPush(task.CodeName, setRecord(task.Id, task.CodeName, task.Args)); err != nil {
			noOk(c, err)
			return
		}
	}
	ok(c, nil)
}

// GetDataTask 获取数据
func GetDataTask(c *gin.Context) {
	info := &model.GetDataTask{}
	c.Bind(info)
	// 获取连接id
	client, err := server.GetClient(info.NodeId)
	if err != nil {
		noOk(c, err)
		return
	}
	// 判断是否为强制推送
	if info.ForcePush {
		// 获取代码
		if err = server.PushCode(client, info.DataName); err != nil {
			noOk(c, err)
			return
		}
	}
	// 任务推送
	if data, err := client.DataGet(info.DataName, []byte(info.Args)); err != nil {
		noOk(c, err)
		return
	} else {
		ok(c, utils.Byte2Interface(data))
	}
}

// PushData 数据推送
func PushData(c *gin.Context) {
	info := &model.PushDataReq{}
	c.Bind(info)
	// 获取连接id
	client, err := server.GetClient(info.NodeId)
	if err != nil {
		noOk(c, err)
		return
	}
	// 数据推送
	if err := client.DataPush(info.Topic, []byte(info.Data)); err != nil {
		noOk(c, err)
		return
	} else {
		ok(c, nil)
	}
}

// PushTask 直接执行任务
func PushTask(c *gin.Context) {
	info := &model.PushTaskReq{}
	c.Bind(info)
	// 获取连接id
	client, err := server.GetClient(info.NodeId)
	if err != nil {
		noOk(c, err)
		return
	}
	if info.PushCode {
		// 需要强制推送才推送
		if err = server.PushCode(client, info.FuncName); err != nil {
			noOk(c, err)
			return
		}
	}

	// 推送任务
	if err := client.TaskPush(info.FuncName, setRecord("", info.FuncName, info.Args)); err != nil {
		noOk(c, err)
		return
	} else {
		ok(c, nil)
	}
}
