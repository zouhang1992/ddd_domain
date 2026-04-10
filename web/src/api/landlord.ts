import apiClient from './request';
import type { Landlord } from '../types/api';

export interface LandlordsQueryResult {
  items: Landlord[];
  total: number;
  page: number;
  limit: number;
}

export interface LandlordQueryParams {
  name?: string;
  phone?: string;
  offset?: number;
  limit?: number;
}

export const landlordApi = {
  list: async (params?: LandlordQueryParams) => {
    const response = await apiClient.get<LandlordsQueryResult>('landlords', { params });
    return response.data;
  },

  get: async (id: string) => {
    const response = await apiClient.get<Landlord>(`landlords/${id}`);
    return response.data;
  },

  create: async (data: { name: string; phone: string; note: string }) => {
    const response = await apiClient.post<Landlord>('landlords', data);
    return response.data;
  },

  update: async (id: string, data: { name: string; phone: string; note: string }) => {
    const response = await apiClient.put<Landlord>(`landlords/${id}`, data);
    return response.data;
  },

  delete: async (id: string) => {
    await apiClient.delete(`landlords/${id}`);
  },
};
