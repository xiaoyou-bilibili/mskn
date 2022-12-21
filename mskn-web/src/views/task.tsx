import {Button, Checkbox, Form, Input, message, Modal, Select, Space, Table} from "antd";
import {FileAddOutlined, PlayCircleOutlined, DeleteOutlined} from "@ant-design/icons";
import {useEffect, useRef, useState} from "react";
import {addTask, Code, Node, deleteTask, getTaskList, Task, getNodeList, executeTask} from "../api/api";
const { TextArea } = Input;
import type { FormInstance } from 'antd/es/form';
const {Option} = Select;

export default function TaskView () {
    const [addTaskShow,setAddTaskShow] = useState(false)
    const [executeTaskShow, setExecuteTaskShow] = useState(false)
    const [taskList, setTaskList] = useState<Task[]>([])
    // 待执行任务id
    const [executeTaskId, setExecuteTaskId] = useState('')
    // 节点列表
    const [nodeList, setNodeList] = useState<Node[]>([])
    // form表单
    const reqForm = useRef<FormInstance>(null)
    // 执行表单
    const executeForm = useRef<FormInstance>(null)
    const addTaskReq = () => {
        if (reqForm.current) {
            addTask(reqForm.current.getFieldsValue()).
            then(res=>{
                message.success('添加成功')
                getTaskListReq()
            }).finally(() => {setAddTaskShow(false)})
        }
    }
    const getTaskListReq = () => {
        getTaskList().then((res)=>{
            setTaskList(res.list)
        })
    }

    const deleteTaskReq = (id:string) => {
        deleteTask(id).then(()=>{
            message.success('删除成功')
            getTaskListReq()
        })
    }

    const getNodeListReq = () => {
        getNodeList().then((res) => {
            setNodeList(res.list)
        })
    }

    // 执行任务
    const executeTaskReq = () => {
        if (executeForm.current) {
            let data = executeForm.current.getFieldsValue()
            data.task_id = executeTaskId
            executeTask(data).then(res=>{
                message.success('执行成功')
            }).finally(() => {setExecuteTaskShow(false)})
        }
    }

    useEffect(() => {
        getTaskListReq()
        getNodeListReq()
    }, [])

    const columns = [
        {
            title: 'id',
            dataIndex: 'id',
            key: 'id',
        },
        {
            title: '任务名字',
            dataIndex: 'task_name',
            key: 'task_name',
        },
        {
            title: '代码名称',
            dataIndex: 'code_name',
            key: 'code_name',
        },
        {
            title: '参数',
            dataIndex: 'args',
            key: 'args',
        },
        {
            title: '操作',
            key: 'action',
            render: (_: any, record: Task) => (<Space>
                    <Button type="primary" onClick={()=>{
                        setExecuteTaskShow(true)
                        setExecuteTaskId(record.id || "")
                    }} icon={<PlayCircleOutlined />}>执行</Button>
                    <Button danger type="primary" onClick={()=>{deleteTaskReq(record.id||"")}} icon={<DeleteOutlined />}>删除</Button>
                </Space>
            ),
        }
    ]

    return <div>
        <Form style={{padding: "5px", marginBottom: "10px"}} layout="inline">
            <Form.Item><Button type="primary" onClick={()=>{setAddTaskShow(true)}} icon={<FileAddOutlined/>}>新增任务</Button></Form.Item>
        </Form>
        <Modal title="新增任务" open={addTaskShow} onCancel={()=>{setAddTaskShow(false)}} onOk={addTaskReq}>
            <Form ref={reqForm}>
                <Form.Item label="任务名称" name={"task_name"}><Input placeholder={"请输入任务名称"} /></Form.Item>
                <Form.Item label="代码名称" name={"code_name"}><Input placeholder={"请输入代码名称"} /></Form.Item>
                <Form.Item label="任务参数" name={"args"}><TextArea rows={4} /></Form.Item>
            </Form>
        </Modal>
        <Modal title="执行任务" open={executeTaskShow} onCancel={()=>{setExecuteTaskShow(false)}} onOk={executeTaskReq}>
            <Form ref={executeForm}>
                <Form.Item label="执行节点" name={"node_id"}>
                    <Select mode="multiple" style={{width: "150px"}}>
                        {nodeList.map(item => <Option key={item.id} value={item.id}>{`${item.name}(${item.addr})`}</Option>)}
                    </Select>
                </Form.Item>
                <Form.Item name={"fore_push"} valuePropName="checked">
                    <Checkbox>强制代码同步</Checkbox>
                </Form.Item>
            </Form>
        </Modal>
        <Table
            rowKey={"id"}
            dataSource={taskList}
            columns={columns}
        />;
    </div>
}