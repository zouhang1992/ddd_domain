import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';
import type { ReactNode } from 'react';
import { authApi, type UserInfo } from '../api/auth';

interface AuthContextType {
  isAuthenticated: boolean;
  user: UserInfo | null;
  login: () => void;
  logout: () => void;
  loading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState<UserInfo | null>(null);
  const [loading, setLoading] = useState(true);
  const [loginInitiated, setLoginInitiated] = useState(false);

  // 检查认证状态
  const checkAuth = useCallback(async () => {
    // 防止重定向循环
    if (loginInitiated) {
      console.log('Login already initiated, skipping check');
      setLoading(false);
      return;
    }

    try {
      const userInfo = await authApi.getUserInfo();
      setUser(userInfo);
      setIsAuthenticated(true);
      setLoginInitiated(false);

      // 登录成功后，检查是否有 returnUrl 需要跳转
      const returnUrl = sessionStorage.getItem('returnUrl');
      if (returnUrl && returnUrl !== '/') {
        sessionStorage.removeItem('returnUrl');
        // 使用更平滑的跳转
        console.log('Redirecting to:', returnUrl);
        window.location.href = returnUrl;
      }
    } catch (error) {
      console.debug('Not authenticated:', error);
      setUser(null);
      setIsAuthenticated(false);
    } finally {
      setLoading(false);
    }
  }, [loginInitiated]);

  // 初始化时检查认证状态
  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  const login = useCallback(() => {
    // 防止重定向循环
    if (loginInitiated) {
      console.log('Login already in progress, skipping');
      return;
    }

    console.log('Initiating login flow...');
    setLoginInitiated(true);
    authApi.login();
  }, [loginInitiated]);

  const logout = useCallback(async () => {
    console.log('Initiating logout flow...');
    setLoginInitiated(false);
    // 清除本地状态
    sessionStorage.removeItem('returnUrl');
    // 直接调用 authApi.logout()，它会跳转到后端，后端再跳转到 Keycloak
    authApi.logout();
  }, []);

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, login, logout, loading }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
