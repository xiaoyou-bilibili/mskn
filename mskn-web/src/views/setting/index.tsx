import {Tabs} from "antd";
import {NodeSetting} from "./node";

export default function SettingView() {
    return  <Tabs
        style={{margin: 10}}
        defaultActiveKey="1"
        items={[
            {
                label: `节点设置`,
                key: '1',
                children: <NodeSetting></NodeSetting>,
            }
        ]}
    />
}