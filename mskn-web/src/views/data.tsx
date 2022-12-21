import {Button, Form, Space, Switch, Table, Typography} from "antd";
import {Code, Data, getDataList, Task} from "../api/api";
import {DeleteOutlined, SyncOutlined} from "@ant-design/icons";
import {useEffect, useRef, useState} from "react";
const { Paragraph, Text } = Typography;

export default function DataView () {
    const [total,setTotal] = useState(0)
    const [current,setCurrent] = useState(1)
    const [size,setSize] = useState(10)
    const [dataList, setDataList] = useState<Data[]>([])
    const [loading, seLoading] = useState(false)
    const [autoRefresh, setAutoRefresh] = useState(false)
    const timer = useRef<NodeJS.Timer>()

    const getDataListReq = () => {
        seLoading(true)
        getDataList(current, size).then((res) =>{
            setTotal(res.total)
            setDataList(res.list)
        }).finally(()=>{
            setTimeout(()=>{seLoading(false)}, 50)
        })
    }


    useEffect(() => {
        getDataListReq()
    }, [current, size])

    useEffect(() => {
        timer.current = setInterval(()=>{
            if (autoRefresh) {
                getDataListReq()
            }
        }, 2000)
        return () => {
            clearInterval(timer.current)
        }
    }, [autoRefresh])

    const columns = [
        {
            title: 'id',
            dataIndex: 'id',
            key: 'id',
        },
        {
            title: '节点id',
            dataIndex: 'node_id',
        },
        {
            title: '客户端',
            dataIndex: 'addr',
        },
        {
            title: 'topic',
            dataIndex: 'topic',
        },
        {
            title: 'data',
            dataIndex: 'data',
            width: 800,
            render: (data: string) => (<Paragraph  copyable={true} ellipsis={{ rows: 2, expandable: false}}>{data}</Paragraph>)
        },
        {
            title: '更新时间',
            dataIndex: 'update',
        },
    ]


    return <div style={{margin: 10}}>
        <Form layout={"inline"}>
            <Form.Item label="自动刷新">
                <Switch onChange={setAutoRefresh} />
            </Form.Item>
            <Form.Item>
                <Button type="primary" onClick={getDataListReq} icon={<SyncOutlined />}>刷新</Button>
            </Form.Item>
        </Form>
        <Table
            loading={loading}
            style={{marginTop: 10}}
            onChange={(pagination, )=>{
                setCurrent(pagination.current || 1)
                setSize(pagination.pageSize || 10)
            }}
            pagination={{total: total, current: current, pageSize: size, showQuickJumper: true}}
            rowKey={"id"}
            dataSource={dataList}
            columns={columns}
        />
    </div>;
}