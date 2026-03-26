import apiClient from './request';
import type { Lease } from '../types/api';

export interface LeasesQueryResult {
  items: Lease[];
  total: number;
  page: number;
  limit: number;
}

export interface LeaseQueryParams {
  tenantName?: string;
  tenantPhone?: string;
  status?: string;
  roomId?: string;
  startDate?: string;
  endDate?: string;
  offset?: number;
  limit?: number;
}

export const leaseApi = {
  list: async (params?: LeaseQueryParams) => {
    const queryParams: any = {};
    if (params?.tenantName) queryParams.tenant_name = params.tenantName;
    if (params?.tenantPhone) queryParams.tenant_phone = params.tenantPhone;
    if (params?.status) queryParams.status = params.status;
    if (params?.roomId) queryParams.room_id = params.roomId;
    if (params?.startDate) queryParams.start_date = params.startDate;
    if (params?.endDate) queryParams.end_date = params.endDate;
    if (params?.offset !== undefined) queryParams.offset = params.offset;
    if (params?.limit !== undefined) queryParams.limit = params.limit;

    const response = await apiClient.get<LeasesQueryResult>('/leases', { params: queryParams });
    return response.data;
  },

  get: async (id: string) => {
    const response = await apiClient.get<Lease>(`/leases/${id}`);
    return response.data;
  },

  create: async (data: {
    roomId: string;
    landlordId: string;
    tenantName: string;
    tenantPhone: string;
    startDate: string;
    endDate: string;
    rentAmount: number;
    note: string;
    depositAmount: number;
    depositNote: string;
  }) => {
    const response = await apiClient.post<Lease>('/leases', {
      room_id: data.roomId,
      landlord_id: data.landlordId,
      tenant_name: data.tenantName,
      tenant_phone: data.tenantPhone,
      start_date: data.startDate,
      end_date: data.endDate,
      rent_amount: data.rentAmount,
      note: data.note,
      deposit_amount: data.depositAmount,
      deposit_note: data.depositNote,
    });
    return response.data;
  },

  update: async (id: string, data: {
    tenantName: string;
    tenantPhone: string;
    startDate: string;
    endDate: string;
    rentAmount: number;
    note: string;
  }) => {
    const response = await apiClient.put<Lease>(`/leases/${id}`, {
      tenant_name: data.tenantName,
      tenant_phone: data.tenantPhone,
      start_date: data.startDate,
      end_date: data.endDate,
      rent_amount: data.rentAmount,
      note: data.note,
    });
    return response.data;
  },

  delete: async (id: string) => {
    await apiClient.delete(`/leases/${id}`);
  },

  renew: async (id: string, data: {
    newStartDate: string;
    newEndDate: string;
    newRentAmount: number;
    note: string;
  }) => {
    const response = await apiClient.post<Lease>(`/leases/${id}/renew`, {
      new_start_date: data.newStartDate,
      new_end_date: data.newEndDate,
      new_rent_amount: data.newRentAmount,
      note: data.note,
    });
    return response.data;
  },

  checkout: async (id: string) => {
    const response = await apiClient.post<Lease>(`/leases/${id}/checkout`);
    return response.data;
  },

  printContract: async (id: string) => {
    const response = await apiClient.get(`/leases/${id}/contract`, {
      responseType: 'blob',
    });
    // 创建下载链接
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', `租约合同_${id}.pdf`);
    document.body.appendChild(link);
    link.click();
    link.remove();
  },

  activate: async (id: string) => {
    const response = await apiClient.put<Lease>(`/leases/${id}/activate`);
    return response.data;
  },
};
