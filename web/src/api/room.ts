import apiClient from './request';
import type { Room } from '../types/api';

export interface RoomsQueryResult {
  items: Room[];
  total: number;
  page: number;
  limit: number;
}

export interface RoomQueryParams {
  locationId?: string;
  roomNumber?: string;
  tags?: string[];
  startDate?: string;
  endDate?: string;
  offset?: number;
  limit?: number;
}

export const roomApi = {
  list: async (params?: RoomQueryParams) => {
    const response = await apiClient.get<RoomsQueryResult>('rooms', { params });
    return response.data;
  },

  get: async (id: string) => {
    const response = await apiClient.get<Room>(`rooms/${id}`);
    return response.data;
  },

  create: async (data: { locationId: string; roomNumber: string; tags: string[] }) => {
    const response = await apiClient.post<Room>('rooms', data);
    return response.data;
  },

  update: async (id: string, data: { locationId: string; roomNumber: string; tags: string[] }) => {
    const response = await apiClient.put<Room>(`rooms/${id}`, data);
    return response.data;
  },

  delete: async (id: string) => {
    await apiClient.delete(`rooms/${id}`);
  },
};
