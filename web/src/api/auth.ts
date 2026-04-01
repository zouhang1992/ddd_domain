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

  // 登出 - 直接跳转到后端，后端会处理重定向
  logout: async (): Promise<void> => {
    try {
      // 直接跳转到后端登出地址，后端会处理重定向到 Keycloak
      window.location.href = '/oauth2/logout';
    } catch (error) {
      console.warn('Logout failed:', error);
      // 出错时还是清除本地状态并刷新
      localStorage.removeItem('token');
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
