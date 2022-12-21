import { BrowserRouter as Router, Routes, Route, Link } from "react-router-dom";
import './App.css'

import {
    FormOutlined,
    LineChartOutlined,
    AuditOutlined,
    AppstoreAddOutlined,
    SettingOutlined,
    ProfileOutlined
} from '@ant-design/icons';
import { Layout, Menu } from 'antd';
import React, {useCallback, useEffect, useState} from 'react';
import CodeView from "./views/code";
import IndexView from "./views/index";
import TaskView from "./views/task";
import SettingView from "./views/setting";
import DataView from "./views/data";
const { Sider, Content } = Layout;

function App() {
    const items =  [
        {
            key: 'index',
            icon: <LineChartOutlined />,
            label: <Link to="/">数据总览</Link>,
        },
        {
            key: 'code',
            icon: <FormOutlined />,
            label:  <Link to="/code">代码编辑</Link>,
        },
        {
            key: 'task',
            icon: <AuditOutlined />,
            label:  <Link to="/task">任务管理</Link>,
        },
        {
            key: 'data',
            icon: <ProfileOutlined />,
            label:  <Link to="/data">数据管理</Link>,
        },
        {
            key: 'setting',
            icon: <SettingOutlined />,
            label:  <Link to="/setting">服务设置</Link>,
        }
    ]

    return (
    <Router>
      <Layout>
          <Sider collapsible theme="light">
              <Menu
                  theme="light"
                  mode="inline"
                  defaultSelectedKeys={['1']}
                  items={items}
              />
          </Sider>
          <Layout className="site-layout">
              <Content
                  className="site-layout-background"
                  style={{margin: '12px 12px', minHeight: 700, background: "#ffffff"}}
              >
                  <Routes>
                      <Route path="/" element={<IndexView />} />
                      <Route path="/code" element={<CodeView />} />
                      <Route path="/task" element={<TaskView />} />
                      <Route path="/setting" element={<SettingView />} />
                      <Route path="/data" element={<DataView />} />
                  </Routes>
              </Content>
          </Layout>
      </Layout>
    </Router>
  )
}

export default App
