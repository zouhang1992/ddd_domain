import React, { useState, useEffect } from 'react';
import { Layout, Menu, Button, Dropdown, Avatar, Typography } from 'antd';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import {
  HomeOutlined,
  EnvironmentOutlined,
  ApartmentOutlined,
  UserOutlined,
  FileTextOutlined,
  DollarOutlined,
  WalletOutlined,
  BarChartOutlined,
  PrinterOutlined,
  LogoutOutlined,
  MenuOutlined,
  HistoryOutlined,
  DownOutlined,
} from '@ant-design/icons';
import { useAuth } from '../context/AuthContext';

const { Header, Sider, Content } = Layout;
const { Text } = Typography;

const menuItems = [
  { key: 'dashboard', icon: <HomeOutlined />, label: '仪表盘', path: '/' },
  { key: 'locations', icon: <EnvironmentOutlined />, label: '位置管理', path: '/locations' },
  { key: 'rooms', icon: <ApartmentOutlined />, label: '房间管理', path: '/rooms' },
  { key: 'landlords', icon: <UserOutlined />, label: '房东管理', path: '/landlords' },
  { key: 'leases', icon: <FileTextOutlined />, label: '租约管理', path: '/leases' },
  { key: 'bills', icon: <DollarOutlined />, label: '账单管理', path: '/bills' },
  { key: 'deposits', icon: <WalletOutlined />, label: '押金管理', path: '/deposits' },
  { key: 'income', icon: <BarChartOutlined />, label: '收入查询', path: '/income' },
  { key: 'print', icon: <PrinterOutlined />, label: '打印服务', path: '/print' },
  { key: 'operation-logs', icon: <HistoryOutlined />, label: '操作日志', path: '/operation-logs' },
];

const AppLayout: React.FC = () => {
  const [collapsed, setCollapsed] = useState(false);
  const [mobileCollapsed, setMobileCollapsed] = useState(false);
  const [isMobile, setIsMobile] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const { logout, user } = useAuth();

  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 768);
    };
    checkMobile();
    window.addEventListener('resize', checkMobile);
    return () => window.removeEventListener('resize', checkMobile);
  }, []);

  const handleMenuClick = ({ key }: { key: string }) => {
    const item = menuItems.find((item) => item.key === key);
    if (item) {
      navigate(item.path);
      if (isMobile) {
        setMobileCollapsed(false);
      }
    }
  };

  const handleLogout = () => {
    logout();
  };

  const selectedKey = menuItems.find((item) => {
    if (item.path === '/') {
      return location.pathname === '/';
    }
    return location.pathname.startsWith(item.path);
  })?.key || 'dashboard';

  const toggleMobileMenu = () => {
    setMobileCollapsed(!mobileCollapsed);
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      {isMobile ? (
        <>
          <Sider
            trigger={null}
            collapsible
            collapsed={mobileCollapsed}
            onCollapse={setMobileCollapsed}
            theme="dark"
            width={240}
            collapsedWidth={0}
            style={{
              position: 'fixed',
              left: 0,
              top: 0,
              bottom: 0,
              zIndex: 1000,
            }}
          >
            <div style={{ height: 64, background: '#001529', color: 'white', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 16, fontWeight: 'bold' }}>
              房屋管理系统
            </div>
            <Menu
              theme="dark"
              selectedKeys={[selectedKey]}
              mode="inline"
              items={menuItems}
              onClick={handleMenuClick}
            />
          </Sider>
          {!mobileCollapsed && (
            <div
              style={{
                position: 'fixed',
                top: 0,
                left: 0,
                right: 0,
                bottom: 0,
                background: 'rgba(0,0,0,0.45)',
                zIndex: 999,
              }}
              onClick={toggleMobileMenu}
            />
          )}
          <Layout style={{ marginLeft: mobileCollapsed ? 0 : 240 }}>
            <Header style={{ background: '#fff', padding: '0 16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <Button
                type="text"
                icon={<MenuOutlined />}
                onClick={toggleMobileMenu}
              />
              <span style={{ fontSize: 16, fontWeight: 'bold' }}>房屋管理系统</span>
              <Dropdown
                menu={{
                  items: [
                    { key: 'user-info', label: <span><UserOutlined /> {user?.name || user?.email || '用户'}</span>, disabled: true },
                    { type: 'divider' },
                    { key: 'logout', label: <span><LogoutOutlined /> 退出登录</span>, onClick: handleLogout },
                  ],
                }}
              >
                <Button type="text" style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <Avatar size="small" icon={<UserOutlined />} />
                  <Text>{user?.name || user?.email || '用户'}</Text>
                  <DownOutlined />
                </Button>
              </Dropdown>
            </Header>
            <Content style={{ margin: '16px', padding: 16, background: '#fff', minHeight: 280 }}>
              <Outlet />
            </Content>
          </Layout>
        </>
      ) : (
        <>
          <Sider collapsible collapsed={collapsed} onCollapse={setCollapsed}>
            <div style={{ height: 64, background: '#001529', color: 'white', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 18, fontWeight: 'bold' }}>
              {collapsed ? 'HS' : '房屋管理系统'}
            </div>
            <Menu
              theme="dark"
              selectedKeys={[selectedKey]}
              mode="inline"
              items={menuItems}
              onClick={handleMenuClick}
            />
          </Sider>
          <Layout>
            <Header style={{ background: '#fff', padding: '0 24px', display: 'flex', justifyContent: 'flex-end', alignItems: 'center' }}>
              <Dropdown
                menu={{
                  items: [
                    { key: 'user-info', label: <span><UserOutlined /> {user?.name || user?.email || '用户'}</span>, disabled: true },
                    { type: 'divider' },
                    { key: 'logout', label: <span><LogoutOutlined /> 退出登录</span>, onClick: handleLogout },
                  ],
                }}
              >
                <Button type="text" style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <Avatar size="small" icon={<UserOutlined />} />
                  <Text strong>{user?.name || user?.email || '用户'}</Text>
                  <DownOutlined />
                </Button>
              </Dropdown>
            </Header>
            <Content style={{ margin: '24px 16px', padding: 24, background: '#fff', minHeight: 280 }}>
              <Outlet />
            </Content>
          </Layout>
        </>
      )}
    </Layout>
  );
};

export default AppLayout;
