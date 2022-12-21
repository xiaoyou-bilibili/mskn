import {useEffect, useRef, useState} from "react";
import {
    addNode,
    deleteNode,
    executeTask, getDataTask,
    getNodeList,
    getTaskList,
    Node, pushDataTask,
    Task
} from "../../api/api";
import {FormInstance} from "antd/es/form";
import {Button, Form, Input, message, Modal, Space, Table, Radio, Tag, Select, Checkbox} from "antd";
import {DeleteOutlined, FileAddOutlined, PlayCircleOutlined, CloudDownloadOutlined, CloudUploadOutlined} from "@ant-design/icons";
const {Option} = Select;
const { TextArea } = Input;

const getConnectType = (tp:number) => {
    switch (tp){
        case 1:
            return <Tag color="#108ee9">TCP</Tag>
        default:
            return <Tag color="#f50">未知</Tag>
    }
}

const getEncryptType = (tp:number) => {
    switch (tp){
        case 0:
            return <Tag color="green">不加密</Tag>
        case 1:
            return <Tag color="blue">Base64</Tag>
        default:
            return <Tag color="warning">未知</Tag>
    }
}

export function NodeSetting () {
    const [addNodeShow,setAddNodeShow] = useState(false)
    const [total,setTotal] = useState(0)
    const [current,setCurrent] = useState(1)
    const [size,setSize] = useState(10)
    const [nodeList, setNodeList] = useState<Node[]>([])
    const [executeTaskShow, setExecuteTaskShow] = useState(false)
    const [taskList, setTaskList] = useState<Task[]>([])
    // 获取数据视图
    const [getDataShow, setGetDataShow] = useState(false)
    const getDataForm = useRef<FormInstance>(null)
    // 数据推送
    const [pushDataShow, setPushDataShow] = useState(false)
    const pushDataForm = useRef<FormInstance>(null)

    // 待执行任务id
    const [nodeId, setNodeId] = useState('')
    // 执行表单
    const executeForm = useRef<FormInstance>(null)

    // form表单
    const reqForm = useRef<FormInstance>(null)
    const addNodeReq = () => {
        if (reqForm.current) {
            addNode(reqForm.current.getFieldsValue()).
            then(res=>{
                message.success('添加成功')
                getNodeListReq()
            }).finally(() => {setAddNodeShow(false)})
        }
    }
    const getNodeListReq = () => {
        getNodeList().then((res) => {
            setNodeList(res.list)
        })
    }

    const deleteNodeReq = (id:string) => {
        deleteNode(id).then(()=>{
            message.success('删除成功')
            getNodeListReq()
        })
    }

    const getTaskListReq = () => {
        getTaskList().then((res) => {
            setTaskList(res.list)
        })
    }

    // 执行任务
    const executeTaskReq = () => {
        if (executeForm.current) {
            let data = executeForm.current.getFieldsValue()
            data.node_id = [nodeId]
            executeTask(data).then(res=>{
                message.success('执行成功')
            }).finally(() => {setExecuteTaskShow(false)})
        }
    }

    // 获取数据
    const getDataReq = () => {
        if (getDataForm.current) {
            let data = getDataForm.current.getFieldsValue()
            data.node_id = nodeId
            getDataTask(data).then(res => {
                message.info(JSON.stringify(res))
            }).finally(() => {
                setGetDataShow(false)
            })
        }
    }

    // 数据推送
    const pushDataReq = () => {
        if (pushDataForm.current) {
            let data = pushDataForm.current.getFieldsValue()
            data.node_id = nodeId
            pushDataTask(data).then(res => {
              message.success('推送成功')
            }).finally(() => {
                setPushDataShow(false)
            })
        }
    }

    useEffect(() => {
        getNodeListReq()
        getTaskListReq()
    }, [])

    const columns = [
        {
            title: 'id',
            dataIndex: 'id',
        },
        {
            title: '名字',
            dataIndex: 'name',
        },
        {
            title: '地址',
            dataIndex: 'addr',
        },
        {
            title: '秘钥',
            dataIndex: 'secret',
        },
        {
            title: '通讯方式',
            dataIndex: 'connect_type',
            render: (value:number) => getConnectType(value)
        },
        {
            title: '加密方式',
            dataIndex: 'encrypt_type',
            render: (value:number) => getEncryptType(value)
        },
        {
            title: '操作',
            key: 'action',
            render: (_: any, record: Node) => (<Space>
                    <Button type="primary" onClick={() => {
                        setExecuteTaskShow(true)
                        setNodeId(record.id || "")
                    }} icon={<PlayCircleOutlined />}>执行任务</Button>
                    <Button type="default" onClick={() => {
                        setGetDataShow(true)
                        setNodeId(record.id || "")
                    }} icon={<CloudDownloadOutlined />}>获取数据</Button>
                    <Button type="default" onClick={() => {
                        setPushDataShow(true)
                        setNodeId(record.id || "")
                    }} icon={<CloudUploadOutlined />}>设置数据</Button>
                    <Button danger type="primary" onClick={()=>{deleteNodeReq(record.id||"")}} icon={<DeleteOutlined />}>删除</Button>
                </Space>
            ),
        }
    ]

    return <div>
        <Form style={{padding: "5px", marginBottom: "10px"}} layout="inline">
            <Form.Item><Button type="primary" onClick={()=>{setAddNodeShow(true)}} icon={<FileAddOutlined/>}>新增节点</Button></Form.Item>
        </Form>
        <Modal title="新增节点" open={addNodeShow} onCancel={()=>{setAddNodeShow(false)}} onOk={addNodeReq}><Form ref={reqForm}>
            <Form.Item label="节点名称" name={"name"}><Input placeholder={"请输入节点名称"} /></Form.Item>
            <Form.Item label="节点地址" name={"addr"}><Input placeholder={"请输入节点地址"} /></Form.Item>
            <Form.Item label="节点秘钥" name={"secret"}><Input placeholder={"请输入节点秘钥"} /></Form.Item>
            <Form.Item label="通讯方式" name={"connect_type"}>
                <Radio.Group><Radio value={1}>TCP</Radio></Radio.Group>
            </Form.Item>
            <Form.Item label="加密方式" name={"encrypt_type"}>
                <Radio.Group>
                    <Radio value={0}>无</Radio>
                    <Radio value={1}>base64</Radio>
                </Radio.Group>
            </Form.Item>
        </Form></Modal>
        <Modal title="执行任务" open={executeTaskShow} onCancel={()=>{setExecuteTaskShow(false)}} onOk={executeTaskReq}>
            <Form ref={executeForm}>
                <Form.Item label="执行任务" name={"task_id"}>
                    <Select style={{width: "150px"}}>
                        {taskList.map(item => <Option key={item.id} value={item.id}>{item.task_name}</Option>)}
                    </Select>
                </Form.Item>
                <Form.Item name={"fore_push"} valuePropName="checked">
                    <Checkbox>强制代码同步</Checkbox>
                </Form.Item>
            </Form>
        </Modal>
        <Modal title="获取数据" open={getDataShow} onCancel={()=>{setGetDataShow(false)}} onOk={getDataReq}>
            <Form ref={getDataForm}>
                <Form.Item label="数据名字" name={"data_name"}><Input /></Form.Item>
                <Form.Item label="数据参数" name={"args"}><TextArea rows={4} /></Form.Item>
                <Form.Item name={"fore_push"} valuePropName="checked">
                    <Checkbox>强制代码同步</Checkbox>
                </Form.Item>
            </Form>
        </Modal>
        <Modal title="推送数据" open={pushDataShow} onCancel={()=>{setPushDataShow(false)}} onOk={pushDataReq}>
            <Form ref={pushDataForm}>
                <Form.Item label="topic" name={"topic"}><Input /></Form.Item>
                <Form.Item label="数据" name={"data"}><TextArea rows={4} /></Form.Item>
            </Form>
        </Modal>
        <Table
            rowKey={"id"}
            onChange={(pagination, )=>{
                setCurrent(pagination.current || 1)
                setSize(pagination.pageSize || 10)
            }}
            pagination={{total: total, current: current, pageSize: size, showQuickJumper: true}}
            dataSource={nodeList}
            columns={columns}
        />;
    </div>
}