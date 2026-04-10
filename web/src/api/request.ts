import axios from 'axios';

// API基础配置
const API_BASE_URL = '/api';

// 创建两个API客户端：一个用于业务API，一个用于oauth2

// 业务API客户端（使用/api前缀）
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  withCredentials: true, // 重要：携带cookie
});

// OAuth2客户端（不使用/api前缀）
export const oauthClient = axios.create({
  baseURL: '',
  timeout: 10000,
  withCredentials: true, // 重要：携带cookie
});

// 将蛇形命名转换为驼峰命名
const snakeToCamel = (str: string): string => {
  return str.replace(/_([a-z])/g, (_match, letter) => letter.toUpperCase());
};

// 递归转换对象的所有字段
const convertKeysToCamelCase = (obj: any): any => {
  if (typeof obj !== 'object' || obj === null) {
    return obj;
  }

  if (Array.isArray(obj)) {
    return obj.map(convertKeysToCamelCase);
  }

  const convertedObj: any = {};
  for (let [key, value] of Object.entries(obj)) {
    const camelKey = snakeToCamel(key);
    convertedObj[camelKey] = convertKeysToCamelCase(value);
  }
  return convertedObj;
};

// 将驼峰命名转换为蛇形命名
const camelToSnake = (str: string): string => {
  return str.replace(/([A-Z])/g, '_$1').toLowerCase().replace(/^_/, '');
};

// 递归转换对象的所有字段
const convertKeysToSnakeCase = (obj: any): any => {
  if (typeof obj !== 'object' || obj === null) {
    return obj;
  }

  if (Array.isArray(obj)) {
    return obj.map(convertKeysToSnakeCase);
  }

  const convertedObj: any = {};
  for (let [key, value] of Object.entries(obj)) {
    const snakeKey = camelToSnake(key);
    convertedObj[snakeKey] = convertKeysToSnakeCase(value);
  }
  return convertedObj;
};

// 业务API客户端拦截器
apiClient.interceptors.request.use(
  (config) => {
    // 请求时将驼峰命名转换为蛇形命名
    if (config.data) {
      config.data = convertKeysToSnakeCase(config.data);
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

apiClient.interceptors.response.use(
  (response) => {
    // 响应时将蛇形命名转换为驼峰命名，但 blob 类型除外
    if (response.data && response.config.responseType !== 'blob') {
      response.data = convertKeysToCamelCase(response.data);
    }
    return response;
  },
  (error) => {
    // 401 由 AuthContext 处理，不要在这里重定向
    return Promise.reject(error);
  }
);

export default apiClient;
