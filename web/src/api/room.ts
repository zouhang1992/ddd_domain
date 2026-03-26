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
    const queryParams: any = {};
    if (params?.locationId) queryParams.location_id = params.locationId;
    if (params?.roomNumber) queryParams.room_number = params.roomNumber;
    if (params?.tags) queryParams.tags = params.tags;
    if (params?.startDate) queryParams.start_date = params.startDate;
    if (params?.endDate) queryParams.end_date = params.endDate;
    if (params?.offset !== undefined) queryParams.offset = params.offset;
    if (params?.limit !== undefined) queryParams.limit = params.limit;

    const response = await apiClient.get<RoomsQueryResult>('/rooms', { params: queryParams });
    return response.data;
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
