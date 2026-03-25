import apiClient from './request';
import type { LoginRequest, LoginResponse } from '../types/api';

export const authApi = {
  login: async (data: LoginRequest) => {
    const response = await apiClient.post<LoginResponse>('/login', data);
    return response.data;
  },

  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },

  setToken: (token: string) => {
    localStorage.setItem('token', token);
  },

  getToken: () => {
    return localStorage.getItem('token');
  },

  isAuthenticated: () => {
    return !!localStorage.getItem('token');
  },
};
