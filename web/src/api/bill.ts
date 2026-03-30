import apiClient from './request';
import type { Bill } from '../types/api';

export interface BillsQueryResult {
  items: Bill[];
  total: number;
  page: number;
  limit: number;
}

export interface BillQueryParams {
  type?: string;
  status?: string;
  leaseId?: string;
  roomId?: string;
  month?: string;
  minAmount?: number;
  maxAmount?: number;
  startDate?: string;
  endDate?: string;
  offset?: number;
  limit?: number;
}

export const billApi = {
  list: async (params?: BillQueryParams) => {
    const response = await apiClient.get<BillsQueryResult>('/bills', { params });
    return response.data;
  },

  get: async (id: string) => {
    const response = await apiClient.get<Bill>(`/bills/${id}`);
    return response.data;
  },

  create: async (data: {
    leaseId: string;
    type: string;
    amount: number;
    rentAmount: number;
    waterAmount: number;
    electricAmount: number;
    otherAmount: number;
    paidAt: string | null;
    note: string;
  }) => {
    const response = await apiClient.post<Bill>('/bills', data);
    return response.data;
  },

  update: async (id: string, data: {
    amount: number;
    rentAmount: number;
    waterAmount: number;
    electricAmount: number;
    otherAmount: number;
    paidAt: string | null;
    note: string;
  }) => {
    const response = await apiClient.put<Bill>(`/bills/${id}`, data);
    return response.data;
  },

  delete: async (id: string) => {
    await apiClient.delete(`/bills/${id}`);
  },

  printReceipt: async (id: string) => {
    const response = await apiClient.get(`/bills/${id}/receipt`, {
      responseType: 'blob',
    });
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', `receipt_${id}.rtf`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  },

  confirmArrival: async (id: string, paidAt?: string) => {
    const response = await apiClient.post(`/bills/${id}/confirm-arrival`, paidAt ? { paidAt } : {});
    return response.data;
  },
};
