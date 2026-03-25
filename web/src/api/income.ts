import apiClient from './request';

export interface IncomeReport {
  year: number;
  month: number;
  rentTotal: number;
  waterTotal: number;
  electricTotal: number;
  otherTotal: number;
  depositIncome: number;
  depositExpense: number;
  total: number;
  totalFormatted: string;
  rentFormatted: string;
  waterFormatted: string;
  electricFormatted: string;
  otherFormatted: string;
  depositIncomeFormatted: string;
  depositExpenseFormatted: string;
}

export const incomeApi = {
  getReport: async (month?: string) => {
    const params: any = {};
    if (month) {
      params.month = month;
    }
    const response = await apiClient.get<IncomeReport>('/income', { params });
    return response.data;
  },
};
