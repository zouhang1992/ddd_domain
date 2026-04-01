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

  // 检查认证状态
  const checkAuth = useCallback(async () => {
    try {
      const userInfo = await authApi.getUserInfo();
      setUser(userInfo);
      setIsAuthenticated(true);
    } catch {
      setUser(null);
      setIsAuthenticated(false);
    } finally {
      setLoading(false);
    }
  }, []);

  // 初始化时检查认证状态
  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  const login = useCallback(() => {
    authApi.login();
  }, []);

  const logout = useCallback(async () => {
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
