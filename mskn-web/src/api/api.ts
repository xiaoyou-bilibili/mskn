import request from "./request"
import {fdatasync} from "fs";

export const base = ''


export interface Code {
    id?: string
    name: string
    content: string
}

export interface Task {
    id?: string
    task_name: string
    code_name: string
    args: string
}

export interface Node {
    id?: string
    name: string
    addr: string
    secret: string
    connect_type: number
    encrypt_type: number
}

export interface Data {
    id?: string
    node_id: string
    addr: string
    topic: string
    data: object
}

export interface Record {
    id?: string
    task_id: string
    func_name: string
    args: string
    status: number
    progress: number
    create: string
    finish: string
}

export interface Page<T> {
    total: number
    list: T
}

export interface Status {
    node: number
    task: number
    data: number
    record: number
}


// 代码相关
export const addCode = (data:Code) => request(`${base}/api/code`, data, 'post')
export const getCodeList = () => request<Page<Code[]>>(`${base}/api/code`, {}, 'get')
export const getCodeContent = (id: string) => request<Code>(`${base}/api/code/${id}`, {}, 'get')
export const updateCode = (code: Code) => request(`${base}/api/code/${code.id}`, code, 'put')
export const deleteCode = (id: string) => request(`${base}/api/code/${id}`, {}, 'delete')

// 任务相关
export const addTask = (data:Task) => request(`${base}/api/task`, data, 'post')
export const getTaskList = () => request<Page<Task[]>>(`${base}/api/task`, {}, 'get')
export const deleteTask = (id:string) => request<Task[]>(`${base}/api/task/${id}`, {}, 'delete')
export const executeTask = (data:any) => request(`${base}/api/task/execute`, data, 'post')
export const getDataTask = (data:any) => request(`${base}/api/task/data`, data, 'post')
export const pushDataTask = (data:any) => request(`${base}/api/task/data/push`, data, 'post')
export const taskPush = (data:any) => request(`${base}/api/task/push`, data, 'post')


// 节点
export const addNode = (data:Node) => request(`${base}/api/node`, data, 'post')
export const getNodeList = () => request<Page<Node[]>>(`${base}/api/node`, {}, 'get')
export const updateNode = (data: Node) => request(`${base}/api/node/${data.id}`, data, 'put')
export const deleteNode = (id: string) => request(`${base}/api/node/${id}`, {}, 'delete')

// 数据
export const getDataList = (page_no:number, page_size: number) => request<Page<Data[]>>(`${base}/api/data?page_no=${page_no}&page_size=${page_size}`, {}, 'get')

// 获取记录
export const getRecordList = (page_no:number, page_size: number) => request<Page<Record[]>>(`${base}/api/record?page_no=${page_no}&page_size=${page_size}`, {}, 'get')

// 获取首页统计
export const getCurrentStatus = () => request<Status>(`${base}/api/status`, {}, 'get')
