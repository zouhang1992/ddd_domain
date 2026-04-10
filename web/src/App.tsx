import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { AuthProvider, useAuth } from './context/AuthContext';
import AppLayout from './components/Layout';
import { AuthLoadingPage, RedirectingPage } from './components/LoadingPage';
import Dashboard from './pages/Dashboard';
import Locations from './pages/Locations';
import Rooms from './pages/Rooms';
import Landlords from './pages/Landlords';
import Leases from './pages/Leases';
import Bills from './pages/Bills';
import Deposits from './pages/Deposits';
import Print from './pages/Print';
import Income from './pages/Income';
import OperationLogs from './pages/OperationLogs';
import 'dayjs/locale/zh-cn';

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, loading, login } = useAuth();

  // 显示认证加载状态
  if (loading) {
    return <AuthLoadingPage />;
  }

  // 未认证则跳转到 OIDC 登录
  if (!isAuthenticated) {
    login();
    return <RedirectingPage to="登录页面" />;
  }

  return <>{children}</>;
};

const App: React.FC = () => {
  return (
    <ConfigProvider locale={zhCN}>
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/" element={
              <ProtectedRoute>
                <AppLayout />
              </ProtectedRoute>
            }>
              <Route index element={<Dashboard />} />
              <Route path="locations" element={<Locations />} />
              <Route path="rooms" element={<Rooms />} />
              <Route path="landlords" element={<Landlords />} />
              <Route path="leases" element={<Leases />} />
              <Route path="bills" element={<Bills />} />
              <Route path="deposits" element={<Deposits />} />
              <Route path="print" element={<Print />} />
              <Route path="income" element={<Income />} />
              <Route path="operation-logs" element={<OperationLogs />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </ConfigProvider>
  );
};

export default App;
