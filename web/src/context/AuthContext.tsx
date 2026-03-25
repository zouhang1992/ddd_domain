import React, { createContext, useContext, useState, useCallback } from 'react';
import type { ReactNode } from 'react';
import { authApi } from '../api/auth';

interface AuthContextType {
  isAuthenticated: boolean;
  login: (token: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(() => authApi.isAuthenticated());

  const login = useCallback((token: string) => {
    authApi.setToken(token);
    setIsAuthenticated(true);
  }, []);

  const logout = useCallback(() => {
    authApi.logout();
    setIsAuthenticated(false);
  }, []);

  return (
    <AuthContext.Provider value={{ isAuthenticated, login, logout }}>
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
