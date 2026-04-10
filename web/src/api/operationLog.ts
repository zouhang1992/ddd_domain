import apiClient from './request';

export interface OperationLog {
  id: string;
  timestamp: string;
  eventName: string;
  domainType: string;
  aggregateId?: string;
  operatorId?: string;
  action: string;
  details?: Record<string, any>;
  createdAt: string;
}

export interface OperationLogsQueryResult {
  items: OperationLog[];
  total: number;
  page: number;
  limit: number;
}

export interface OperationLogQueryParams {
  domainType?: string;
  eventName?: string;
  aggregateId?: string;
  operatorId?: string;
  startTime?: string;
  endTime?: string;
  offset?: number;
  limit?: number;
}

export const operationLogApi = {
  list: async (params?: OperationLogQueryParams) => {
    const response = await apiClient.get<OperationLogsQueryResult>('operation-logs', { params });
    return response.data;
  },

  get: async (id: string) => {
    const response = await apiClient.get<OperationLog>(`operation-logs/${id}`);
    return response.data;
  },
};
