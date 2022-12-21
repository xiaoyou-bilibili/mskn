import {Card, Col, Progress, Row, Statistic, Table, Tag} from "antd";
import { DashboardOutlined, DatabaseOutlined, BarChartOutlined, AreaChartOutlined } from '@ant-design/icons';
import {useEffect, useState} from "react";
import {getCurrentStatus, getRecordList, Record, Status} from "../api/api";

const getStatus = (status:number) => {
    switch (status){
        case 0:
            return <Tag color="">待开始</Tag>
        case 1:
            return <Tag color="green">进行中</Tag>
        case 2:
            return <Tag color="blue">已完成</Tag>
        default:
            return <Tag color="warning">未知</Tag>
    }
}

const columns = [
    {
        title: 'id',
        dataIndex: 'id',
    },
    {
        title: '任务id',
        dataIndex: 'task_id',
    },
    {
        title: '函数名称',
        dataIndex: 'func_name',
    },
    {
        title: '参数',
        dataIndex: 'args',
    },
    {
        title: '状态',
        dataIndex: 'status',
        render: (status:number) => getStatus(status)
    },
    {
        title: '进度',
        dataIndex: 'progress',
        render: (progress: number) => <Progress percent={progress} size="small" />
    },
    {
        title: '开始时间',
        dataIndex: 'create'
    },
]

export default function IndexView() {
    const [total,setTotal] = useState(0)
    const [current,setCurrent] = useState(1)
    const [size,setSize] = useState(10)
    const [recordList, setRecordList] = useState<Record[]>([])
    const [status, setStatus] = useState<Status>()
    const getRecordListReq = () => {
        getRecordList(current, size).then((res) => {
            setRecordList(res.list)
            setTotal(res.total)
        })
    }

    useEffect(()=>{
        getRecordListReq()
    }, [current, size])

    useEffect(() => {
        getCurrentStatus().then(setStatus)
    }, [])

    return <div style={{margin: 10}}>
        <Row gutter={16}>
            <Col span={6}>
                <Card>
                    <Statistic
                        title="节点数"
                        value={status?.node}
                        prefix={<DatabaseOutlined />}
                    />
                </Card>
            </Col>
            <Col span={6}>
                <Card>
                    <Statistic
                        title="任务数"
                        value={status?.task}
                        prefix={<DashboardOutlined />}
                    />
                </Card>
            </Col>
            <Col span={6}>
                <Card>
                    <Statistic
                        title="执行数"
                        value={status?.record}
                        prefix={<AreaChartOutlined />}
                    />
                </Card>
            </Col>
            <Col span={6}>
                <Card>
                    <Statistic
                        title="数据量"
                        value={status?.data}
                        prefix={<BarChartOutlined />}
                    />
                </Card>
            </Col>
        </Row>
        <Table
            style={{marginTop: 20}}
            onChange={(pagination, )=>{
                setCurrent(pagination.current || 1)
                setSize(pagination.pageSize || 10)
            }}
            pagination={{total: total, current: current, pageSize: size, showQuickJumper: true}}
            rowKey={"id"}
            dataSource={recordList}
            columns={columns}
        />;
    </div>
}