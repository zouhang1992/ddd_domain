import apiClient from './request';

export const printApi = {
  printBill: async (billId: string) => {
    const response = await apiClient.post<{ job_id: string }>('/print/bill', { bill_id: billId });
    return response.data;
  },

  printLease: async (leaseId: string) => {
    const response = await apiClient.post<{ job_id: string }>('/print/lease', { lease_id: leaseId });
    return response.data;
  },

  printInvoice: async (billId: string) => {
    const response = await apiClient.post<{ job_id: string }>('/print/invoice', { bill_id: billId });
    return response.data;
  },

  getPrintContent: async (billId: string) => {
    const response = await apiClient.get(`/print/content/${billId}`, {
      responseType: 'blob',
    });
    return response.data;
  },

  // Directly download receipt (simpler flow)
  downloadReceipt: async (billId: string) => {
    const response = await apiClient.get(`/bills/${billId}/receipt`, {
      responseType: 'blob',
    });
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', `receipt_${billId}.rtf`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  },

  // Download lease contract
  downloadContract: async (leaseId: string) => {
    const response = await apiClient.get(`/leases/${leaseId}/contract`, {
      responseType: 'blob',
    });
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', `contract_${leaseId}.rtf`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  },
};
