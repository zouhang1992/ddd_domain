import apiClient from './request';
import type { Location } from '../types/api';

export interface LocationsQueryResult {
  items: Location[];
  total: number;
  page: number;
  limit: number;
}

export interface LocationQueryParams {
  shortName?: string;
  detail?: string;
  offset?: number;
  limit?: number;
}

export const locationApi = {
  list: async (params?: LocationQueryParams) => {
    const response = await apiClient.get<LocationsQueryResult>('/api/locations', { params });
    return response.data;
  },

  get: async (id: string) => {
    const response = await apiClient.get<Location>(`/api/locations/${id}`);
    return response.data;
  },

  create: async (data: { shortName: string; detail: string }) => {
    const response = await apiClient.post<Location>('/api/locations', data);
    return response.data;
  },

  update: async (id: string, data: { shortName: string; detail: string }) => {
    const response = await apiClient.put<Location>(`/api/locations/${id}`, data);
    return response.data;
  },

  delete: async (id: string) => {
    await apiClient.delete(`/api/locations/${id}`);
  },
};
