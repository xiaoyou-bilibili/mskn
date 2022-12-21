package web

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"mskn-server/infra/mongo"
	"mskn-server/model"
	"mskn-server/utils"
)

func ok(c *gin.Context, data interface{}) {
	c.JSON(200, &model.Resp{
		Code: 200,
		Msg:  "ok",
		Data: data,
	})
}

func noOk(c *gin.Context, err error) {
	c.JSON(200, &model.Resp{
		Code: 500,
		Msg:  err.Error(),
		Data: nil,
	})
}

// 新增数据
func addData(collection *mongo.Collection, data interface{}) func(context *gin.Context) {
	return func(c *gin.Context) {
		c.Bind(data)
		if err := collection.InsertOne(data); err != nil {
			noOk(c, err)
			return
		}
		ok(c, nil)
	}
}

// 获取列表
func getList(collection *mongo.Collection, data interface{}) func(context *gin.Context) {
	return func(c *gin.Context) {
		// 分页功能
		pageNo := utils.String2Int64(c.Query("page_no"))
		pageSize := utils.String2Int64(c.Query("page_size"))
		if total, err := collection.FindByPage(nil, pageNo, pageSize, &data); err != nil {
			noOk(c, err)
			return
		} else {
			ok(c, map[string]interface{}{
				"total": total,
				"list":  data,
			})
		}
	}
}

// 获取内容
func getData(collection *mongo.Collection, data interface{}) func(context *gin.Context) {
	return func(c *gin.Context) {
		if err := collection.FindOne(mongo.NewFilter().Id(c.Param("id")), data); err != nil {
			noOk(c, err)
			return
		}
		ok(c, data)
	}
}

// 更新内容
func updateData(collection *mongo.Collection, data interface{}) func(context *gin.Context) {
	return func(c *gin.Context) {
		c.Bind(data)
		// 解析为json然后删除id字段
		var tmp map[string]interface{}
		be, _ := json.Marshal(data)
		json.Unmarshal(be, &tmp)
		delete(tmp, "id")
		if err := collection.UpdateOne(
			mongo.NewFilter().Id(c.Param("id")),
			mongo.NewUpdate().Set(tmp),
		); err != nil {
			noOk(c, err)
			return
		}
		ok(c, nil)
	}
}

// 删除内容
func deleteData(collection *mongo.Collection) func(context *gin.Context) {
	return func(c *gin.Context) {
		if err := collection.DeleteOne(mongo.NewFilter().Id(c.Param("id"))); err != nil {
			noOk(c, err)
			return
		}
		ok(c, nil)
	}
}
