import apiClient from './request';

export interface Deposit {
  id: string;
  leaseId: string;
  amount: number;
  status: string;
  refundedAt?: string;
  deductedAt?: string;
  note: string;
  createdAt: string;
  updatedAt: string;
}

export interface DepositsQueryResult {
  items: Deposit[];
  total: number;
  page: number;
  limit: number;
}

export interface DepositQueryParams {
  leaseId?: string;
  status?: string;
  offset?: number;
  limit?: number;
  page?: number;
}

export const depositApi = {
  list: async (params?: DepositQueryParams) => {
    const response = await apiClient.get<DepositsQueryResult>('deposits', { params });
    return response.data;
  },

  get: async (id: string) => {
    const response = await apiClient.get<Deposit>(`deposits/${id}`);
    return response.data;
  },

  markReturning: async (id: string) => {
    const response = await apiClient.post<Deposit>(`deposits/${id}/mark-returning`);
    return response.data;
  },

  markReturned: async (id: string) => {
    const response = await apiClient.post<Deposit>(`deposits/${id}/mark-returned`);
    return response.data;
  },
};
