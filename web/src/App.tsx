import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { AuthProvider, useAuth } from './context/AuthContext';
import AppLayout from './components/Layout';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Locations from './pages/Locations';
import Rooms from './pages/Rooms';
import Landlords from './pages/Landlords';
import Leases from './pages/Leases';
import Bills from './pages/Bills';
import Print from './pages/Print';
import Income from './pages/Income';
import OperationLogs from './pages/OperationLogs';
import 'dayjs/locale/zh-cn';

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />;
};

const App: React.FC = () => {
  return (
    <ConfigProvider locale={zhCN}>
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<Login />} />
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
