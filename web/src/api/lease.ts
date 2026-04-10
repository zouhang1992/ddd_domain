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
    const response = await apiClient.get<LeasesQueryResult>('leases', { params });
    return response.data;
  },

  get: async (id: string) => {
    const response = await apiClient.get<Lease>(`leases/${id}`);
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
    const response = await apiClient.post<Lease>('leases', data);
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
    const response = await apiClient.put<Lease>(`leases/${id}`, data);
    return response.data;
  },

  delete: async (id: string) => {
    await apiClient.delete(`leases/${id}`);
  },

  renew: async (id: string, data: {
    newStartDate: string;
    newEndDate: string;
    newRentAmount: number;
    note: string;
  }) => {
    const response = await apiClient.post<Lease>(`leases/${id}/renew`, data);
    return response.data;
  },

  checkout: async (id: string) => {
    const response = await apiClient.post<Lease>(`leases/${id}/checkout`);
    return response.data;
  },

  checkoutWithBills: async (id: string, data: {
    refundRentAmount: number;
    refundDepositAmount: number;
    waterAmount: number;
    electricAmount: number;
    otherAmount: number;
    note: string;
  }) => {
    const response = await apiClient.post(`leases/${id}/checkout-with-bills`, data);
    return response.data;
  },

  printContract: async (id: string) => {
    const response = await apiClient.get(`leases/${id}/contract`, {
      responseType: 'blob',
    });
    const url = window.URL.createObjectURL(new Blob([response.data], { type: 'text/html' }));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', `contract_${id}.html`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  },

  activate: async (id: string) => {
    const response = await apiClient.put<Lease>(`leases/${id}/activate`);
    return response.data;
  },
};
