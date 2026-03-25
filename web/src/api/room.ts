import apiClient from './request';
import type { Room } from '../types/api';

export const roomApi = {
  list: async (locationId?: string) => {
    const params = locationId ? { location_id: locationId } : {};
    const response = await apiClient.get<Room[]>('/rooms', { params });
    return response.data || [];
  },

  get: async (id: string) => {
    const response = await apiClient.get<Room>(`/rooms/${id}`);
    return response.data;
  },

  create: async (data: { locationId: string; roomNumber: string; tags: string[] }) => {
    const response = await apiClient.post<Room>('/rooms', {
      location_id: data.locationId,
      room_number: data.roomNumber,
      tags: data.tags,
    });
    return response.data;
  },

  update: async (id: string, data: { locationId: string; roomNumber: string; tags: string[] }) => {
    const response = await apiClient.put<Room>(`/rooms/${id}`, {
      location_id: data.locationId,
      room_number: data.roomNumber,
      tags: data.tags,
    });
    return response.data;
  },

  delete: async (id: string) => {
    await apiClient.delete(`/rooms/${id}`);
  },
};
