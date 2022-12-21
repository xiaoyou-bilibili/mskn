import * as monaco from 'monaco-editor';
import {useEffect, useRef, useState} from "react";
import {Button, Form, Input, Modal, Select, message, Checkbox} from "antd";
import {
    SaveOutlined,
    FileAddOutlined,
    DeleteOutlined,
    PlayCircleOutlined
} from '@ant-design/icons';
import {
    addCode,
    Code,
    deleteCode,
    getCodeContent,
    getCodeList,
    getNodeList,
    updateCode,
    Node,
    taskPush
} from "../api/api";
import {FormInstance} from "antd/es/form";

let editor: monaco.editor.IStandaloneCodeEditor;
const {Option} = Select;
const { TextArea } = Input;

export default function CodeView() {
    const editorElem = useRef(null)
    const [codeList, setCodeList] = useState<Code[]>([])
    const [addCodeShow,setAddCodeShow] = useState(false)
    const [codeName, setCodename] = useState('')
    const [currentCode, setCurrentCode] = useState<Code>({id: "", name: "", content: ""})
    // 节点列表
    const [nodeList, setNodeList] = useState<Node[]>([])
    // 当前选中节点
    const [currentNode,setCurrentNode] = useState('')
    // 执行任务
    const [executeCodeShow,setExecuteCodeShow] = useState(false)
    const executeCodeForm = useRef<FormInstance>(null)
    const [executeWithSave, setExecuteWithSave] = useState(false)

    const addCodeReq = () => {
        addCode({name: codeName, content: "123"}).then(res => {
            message.success('添加成功')
            getCodeListReq()
        }).finally(() => {setAddCodeShow(false)})
    }

    const getCodeListReq = () => {
        getCodeList().then(data => {
            setCodeList(data.list)
        })
    }

    const getCodeContentReq = (id: string) => {
        getCodeContent(id).then(setCurrentCode)
    }

    const saveCode = () => {
        setCurrentCode(code => {
            code.content = editor.getValue()
            return code
        })
        updateCode(currentCode).then(res=>{
            message.success('保存成功')
        })
    }

    const deleteCodeReq = () => {
        deleteCode(currentCode.id || "").then(res => {
            message.success('删除成功')
            getCodeListReq()
            setCurrentCode({id: "", name: "", content: ""})
        })
    }


    const handleKeyDown = (event: React.KeyboardEvent) => {
        if (event.key === "s" && (event.ctrlKey || event.metaKey)) {
            currentCode.id === ''? message.error('请选择代码'):saveCode()
            // 阻止默认事件
            event.preventDefault();
        }
    };

    const getNodeListReq = () => {
        getNodeList().then((res) => {
            setNodeList(res.list)
        })
    }

    const executeCodeReq = () => {
        if (executeWithSave) {
            // 先保存代码
            saveCode()
        }
        // 执行任务
        if (executeCodeForm.current){
            let data = executeCodeForm.current.getFieldsValue()
            data.node_id = currentNode
            data.push_code = executeWithSave
            taskPush(data).then(() => {
                message.success('执行成功')
            }).finally(()=>setExecuteCodeShow(false))
        }
    }

    useEffect(() => {
        editor = monaco.editor.create(editorElem.current!, {
            value: '', // 编辑器初始显示文字
            language: 'python', // 使用JavaScript语言
            automaticLayout: true, // 自适应布局
            theme: 'vs-dark', // 官方自带三种主题vs, hc-black, or vs-dark
            foldingStrategy: 'indentation', // 代码折叠
            renderLineHighlight: 'all', // 行亮
            selectOnLineNumbers: true, // 显示行号
            minimap: {enabled: false}, // 是否开启小地图（侧边栏的那个全览图）
            readOnly: false, // 只读
            fontSize: 18, // 字体大小
            scrollBeyondLastLine: false, // 取消代码后面一大段空白
            overviewRulerBorder: false, // 不要滚动条的边框
        })
        getCodeListReq()
        getNodeListReq()
    }, [])


    useEffect(() => {
        editor.setValue(currentCode.content)
    }, [currentCode])

    return <div>
        <Form style={{padding: "5px", marginBottom: "10px"}} layout="inline">
            <Form.Item><Button type="primary" onClick={()=>{setAddCodeShow(true)}} icon={<FileAddOutlined/>}>新增代码</Button></Form.Item>
            <Form.Item label="代码列表" name="plugin">
                <Select style={{width: "150px"}} onSelect={getCodeContentReq}>
                    {codeList.map(item => <Option key={item.id} value={item.id}>{item.name}</Option>)}
                </Select>
            </Form.Item>
            <Form.Item><Button type="primary" onClick={saveCode}  icon={<SaveOutlined />}>保存</Button></Form.Item>
            <Form.Item><Button danger type="primary" onClick={deleteCodeReq}  icon={<DeleteOutlined />}>删除</Button></Form.Item>
            <Form.Item label="测试节点" name="node">
                <Select style={{width: "150px"}} onSelect={setCurrentNode}>
                    {nodeList.map(item => <Option key={item.id} value={item.id}>{`${item.name}(${item.addr})`}</Option>)}
                </Select>
            </Form.Item>
            <Form.Item><Checkbox onChange={(e)=>setExecuteWithSave(e.target.checked)}>强制更新</Checkbox></Form.Item>
            <Form.Item><Button type="primary" onClick={()=>setExecuteCodeShow(true)}  icon={<PlayCircleOutlined />}>执行</Button></Form.Item>
        </Form>
        <Modal title="新增代码" open={addCodeShow} onCancel={()=>{setAddCodeShow(false)}} onOk={addCodeReq}>
            <Input addonBefore={"代码名称"} placeholder={"请输入代码名称"} value={codeName} onChange={(e)=>{setCodename(e.target.value)}} />
        </Modal>
        <Modal title="执行" open={executeCodeShow} onCancel={()=>{setExecuteCodeShow(false)}} onOk={executeCodeReq}>
            <Form ref={executeCodeForm}>
                <Form.Item label="代码名称" name={"func_name"}><Input placeholder={"请输入代码名称"} /></Form.Item>
                <Form.Item label="任务参数" name={"args"}><TextArea rows={4} /></Form.Item>
            </Form>
        </Modal>
        <div onKeyDown={handleKeyDown} style={{height: 800}} id="codeEditBox"  ref={editorElem} />
    </div>
}