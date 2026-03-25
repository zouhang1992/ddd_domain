import apiClient from './request';
import type { Location } from '../types/api';

export const locationApi = {
  list: async () => {
    const response = await apiClient.get<Location[]>('/locations');
    return response.data || [];
  },

  get: async (id: string) => {
    const response = await apiClient.get<Location>(`/locations/${id}`);
    return response.data;
  },

  create: async (data: { shortName: string; detail: string }) => {
    const response = await apiClient.post<Location>('/locations', data);
    return response.data;
  },

  update: async (id: string, data: { shortName: string; detail: string }) => {
    const response = await apiClient.put<Location>(`/locations/${id}`, data);
    return response.data;
  },

  delete: async (id: string) => {
    await apiClient.delete(`/locations/${id}`);
  },
};
