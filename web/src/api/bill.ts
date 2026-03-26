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
    const queryParams: any = {};
    if (params?.type) queryParams.type = params.type;
    if (params?.status) queryParams.status = params.status;
    if (params?.leaseId) queryParams.lease_id = params.leaseId;
    if (params?.roomId) queryParams.room_id = params.roomId;
    if (params?.month) queryParams.month = params.month;
    if (params?.minAmount !== undefined) queryParams.min_amount = params.minAmount;
    if (params?.maxAmount !== undefined) queryParams.max_amount = params.maxAmount;
    if (params?.startDate) queryParams.start_date = params.startDate;
    if (params?.endDate) queryParams.end_date = params.endDate;
    if (params?.offset !== undefined) queryParams.offset = params.offset;
    if (params?.limit !== undefined) queryParams.limit = params.limit;

    const response = await apiClient.get<BillsQueryResult>('/bills', { params: queryParams });
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
    const response = await apiClient.post<Bill>('/bills', {
      lease_id: data.leaseId,
      type: data.type,
      amount: data.amount,
      rent_amount: data.rentAmount,
      water_amount: data.waterAmount,
      electric_amount: data.electricAmount,
      other_amount: data.otherAmount,
      paid_at: data.paidAt,
      note: data.note,
    });
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
    const response = await apiClient.put<Bill>(`/bills/${id}`, {
      amount: data.amount,
      rent_amount: data.rentAmount,
      water_amount: data.waterAmount,
      electric_amount: data.electricAmount,
      other_amount: data.otherAmount,
      paid_at: data.paidAt,
      note: data.note,
    });
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
    const response = await apiClient.post(`/bills/${id}/confirm-arrival`, paidAt ? { paid_at: paidAt } : {});
    return response.data;
  },
};
