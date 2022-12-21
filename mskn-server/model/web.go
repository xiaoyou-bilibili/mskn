package model

import (
	"mskn-server/core/proto/mn"
	"time"
)

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Code 代码相关
type Code struct {
	Id      string `json:"id" bson:"_id,omitempty"`
	Name    string `json:"name"`    // 代码名字
	Content string `json:"content"` // 代码内容
}

// CodeList 代码列表
type CodeList struct {
	Id   string `json:"id" bson:"_id"`
	Name string `json:"name"`
}

// Task 新增任务
type Task struct {
	Id       string `json:"id" bson:"_id,omitempty"`
	TaskName string `json:"task_name"` // 任务名称
	CodeName string `json:"code_name"` // 代码名称
	Args     string `json:"args"`      // 任务参数
}

// ExecuteTask 执行任务
type ExecuteTask struct {
	TaskId    string   `json:"task_id"`   // 任务ID
	NodeId    []string `json:"node_id"`   // 节点id
	ForcePush bool     `json:"fore_push"` // 强制推送
}

// GetDataTask 获取数据任务
type GetDataTask struct {
	DataName  string `json:"data_name"` // 数据名字
	NodeId    string `json:"node_id"`   // 节点id
	Args      string `json:"args"`      // 任务参数
	ForcePush bool   `json:"fore_push"` // 强制推送
}

// PushDataReq 数据推送
type PushDataReq struct {
	Topic  string `json:"topic"`   // 数据topic
	NodeId string `json:"node_id"` // 节点id
	Data   string `json:"data"`    // 数据内容
}

// PushTaskReq 执行任务
type PushTaskReq struct {
	NodeId   string `json:"node_id"`   // 节点id
	FuncName string `json:"func_name"` // 执行函数名称
	Args     string `json:"args"`      // 执行参数
	PushCode bool   `json:"push_code"` // 是否推送代码
}

// Node 集群与节点
type Node struct {
	Id          string         `json:"id" bson:"_id,omitempty"`
	Name        string         `json:"name"`         // 节点名字
	Addr        string         `json:"addr"`         // 节点地址
	Secret      string         `json:"secret"`       // 秘钥
	ConnectType int32          `json:"connect_type"` // 节点通讯方式 1 tcp
	EncryptType mn.EncryptType `json:"encrypt_type"` // 加密方式 0 无 1 base64
}

// DataPush 数据推送信息
type DataPush struct {
	Id     string    `json:"id" bson:"_id,omitempty"`
	NodeId string    `json:"node_id"` // 节点id
	Addr   string    `json:"addr"`    // 地址
	Topic  string    `json:"topic"`   // 数据topic
	Data   string    `json:"data"`    // 节点数据
	Update time.Time `json:"update"`  // 推送时间
}

// Record 执行记录
type Record struct {
	Id       string    `json:"id" bson:"_id,omitempty"` // 记录id
	TaskId   string    `json:"task_id"`                 // 任务id
	FuncName string    `json:"func_name"`               // 函数名称
	Args     string    `json:"args"`                    // 执行参数
	Status   int32     `json:"status" `                 // 状态 0 开始 1 进行中 2 完成
	Progress int32     `json:"progress"`                // 执行进度 0-100
	Create   time.Time `json:"create" `                 // 开始时间
	Finish   time.Time `json:"finish" `                 // 完成时间
}
