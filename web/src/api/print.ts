import apiClient from './request';

export interface PrintJobQueryParams {
  status?: string;
  type?: string;
  startDate?: string;
  endDate?: string;
  offset?: number;
  limit?: number;
}

export interface PrintJob {
  id: string;
  type: string;
  status: string;
  reference_id: string;
  tenant_name: string;
  tenant_phone: string;
  room_id: string;
  room_number: string;
  address: string;
  landlord_name: string;
  amount: number;
  error_msg?: string;
  created_at: string;
  updated_at: string;
  completed_at?: string;
  type_text: string;
  status_text: string;
  amount_yuan: string;
}

export interface PrintJobsQueryResult {
  items: PrintJob[];
  total: number;
  page: number;
  limit: number;
}

export const printApi = {
  printBill: async (billId: string) => {
    const response = await apiClient.post<{ jobId: string }>('/print/bill', { billId });
    return response.data;
  },

  printLease: async (leaseId: string) => {
    const response = await apiClient.post<{ jobId: string }>('/print/lease', { leaseId });
    return response.data;
  },

  printInvoice: async (billId: string) => {
    const response = await apiClient.post<{ jobId: string }>('/print/invoice', { billId });
    return response.data;
  },

  getPrintContent: async (billId: string) => {
    const response = await apiClient.get(`/print/content/${billId}`, {
      responseType: 'blob',
    });
    return response.data;
  },

  listPrintJobs: async (params?: PrintJobQueryParams) => {
    const response = await apiClient.get<PrintJobsQueryResult>('/print/jobs', { params });
    return response.data;
  },

  // Directly download receipt (simpler flow)
  downloadReceipt: async (billId: string) => {
    const response = await apiClient.get(`/bills/${billId}/receipt`, {
      responseType: 'blob',
    });
    const url = window.URL.createObjectURL(new Blob([response.data], { type: 'text/html' }));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', `receipt_${billId}.html`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  },

  // Download lease contract
  downloadContract: async (leaseId: string) => {
    const response = await apiClient.get(`/leases/${leaseId}/contract`, {
      responseType: 'blob',
    });
    const url = window.URL.createObjectURL(new Blob([response.data], { type: 'text/html' }));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', `contract_${leaseId}.html`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  },
};
