import apiClient from './request';

export interface UserInfo {
  sub: string;
  email: string;
  name: string;
  roles: string[];
  permissions: string[];
}

export const authApi = {
  // 启动 OIDC 登录流程 - 直接跳转到后端
  login: () => {
    window.location.href = '/oauth2/login';
  },

  // 登出
  logout: async () => {
    try {
      await apiClient.post('/oauth2/logout');
    } finally {
      // 无论登出 API 是否成功，都刷新页面
      window.location.href = '/';
    }
  },

  // 获取当前用户信息
  getUserInfo: async (): Promise<UserInfo> => {
    const response = await apiClient.get<UserInfo>('/oauth2/userinfo');
    return response.data;
  },

  // 检查是否已认证（通过尝试获取用户信息）
  isAuthenticated: async (): Promise<boolean> => {
    try {
      await authApi.getUserInfo();
      return true;
    } catch {
      return false;
    }
  },
};
