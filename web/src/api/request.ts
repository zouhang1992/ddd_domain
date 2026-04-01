import axios from 'axios';

// 使用相对路径配合 Vite 代理
const API_BASE_URL = import.meta.env.DEV ? '' : 'http://localhost:8080';

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

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
});

apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

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
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
    }
    return Promise.reject(error);
  }
);

export default apiClient;
