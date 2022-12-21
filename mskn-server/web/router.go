package web

import (
	"github.com/gin-gonic/gin"
	"mskn-server/model"
	"mskn-server/server"
)

func RegisterRouter(g *gin.Engine) {
	// 代码相关
	g.POST("/api/code", addData(server.CodeCollection, &model.Code{}))
	g.GET("/api/code", getList(server.CodeCollection, []*model.CodeList{}))
	g.GET("/api/code/:id", getData(server.CodeCollection, &model.Code{}))
	g.PUT("/api/code/:id", updateData(server.CodeCollection, &model.Code{}))
	g.DELETE("/api/code/:id", deleteData(server.CodeCollection))
	// 任务相关
	g.POST("/api/task", addData(server.TaskCollection, &model.Task{}))
	g.GET("/api/task", getList(server.TaskCollection, []*model.Task{}))
	g.DELETE("/api/task/:id", deleteData(server.TaskCollection))
	// 执行任务
	g.POST("/api/task/execute", ExecuteTask)
	// 获取数据
	g.POST("/api/task/data", GetDataTask)
	// 设置数据
	g.POST("/api/task/data/push", PushData)
	// 执行任务
	g.POST("/api/task/push", PushTask)

	// 节点相关
	g.POST("/api/node", addData(server.NodeCollection, &model.Node{}))
	g.GET("/api/node", getList(server.NodeCollection, []*model.Node{}))
	g.PUT("/api/node/:id", updateData(server.NodeCollection, &model.Node{}))
	g.DELETE("/api/node/:id", deleteData(server.NodeCollection))

	// 数据相关
	g.GET("/api/data", getList(server.DataCollection, []*model.DataPush{}))

	// 任务执行记录
	g.POST("/api/record", addData(server.RecordCollection, &model.Record{}))
	g.GET("/api/record", getList(server.RecordCollection, []*model.Record{}))
	g.PUT("/api/record/:id", updateData(server.RecordCollection, &model.Record{}))

	// 获取首页统计
	g.GET("/api/status", func(context *gin.Context) {
		node, _ := server.NodeCollection.Count(nil)
		task, _ := server.TaskCollection.Count(nil)
		data, _ := server.DataCollection.Count(nil)
		record, _ := server.RecordCollection.Count(nil)
		ok(context, map[string]interface{}{
			"node":   node,
			"task":   task,
			"data":   data,
			"record": record,
		})
	})
}
