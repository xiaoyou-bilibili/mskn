
// 代码推送
struct CodePush{
   1: string addr // 服务地址
   2: string name // 代码名字
   3: string code // 代码内容
}

// 代码获取
struct CodeGet {
    1: string addr // 服务地址
    2: string name // 代码名字
}

// 数据推送
struct DataPush {
    1: string addr // 服务地址
    2: string topic // 推送topic信息
    3: binary data // 推送数据
}

// 执行任务
struct Task {
    1: string addr // 服务地址
    2: string name // 任务名称
    3: binary param // 任务参数
}

// 获取数据
struct DataGet {
    1: string addr // 服务地址
    2: string name // 数据名称
    3: binary param // 参数
}

// 数据返回
struct DataBack {
    1: string addr // 服务地址
    2: string name // 数据名称
    3: binary data // 返回的数据
}
