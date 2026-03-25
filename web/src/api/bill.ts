import apiClient from './request';
import type { Bill } from '../types/api';

export const billApi = {
  list: async (leaseId?: string, roomId?: string, month?: string) => {
    const params: any = {};
    if (leaseId) params.lease_id = leaseId;
    if (roomId) params.room_id = roomId;
    if (month) params.month = month;
    const response = await apiClient.get<Bill[]>('/bills', { params });
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
