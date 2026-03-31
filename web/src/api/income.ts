import apiClient from './request';

export interface IncomeReport {
  year: number;
  month: number;
  rentIncome: number;
  waterIncome: number;
  electricIncome: number;
  otherIncome: number;
  depositIncome: number;
  rentExpense: number;
  depositExpense: number;
  totalIncome: number;
  totalExpense: number;
  netIncome: number;
  rentIncomeFormatted: string;
  waterIncomeFormatted: string;
  electricIncomeFormatted: string;
  otherIncomeFormatted: string;
  depositIncomeFormatted: string;
  rentExpenseFormatted: string;
  depositExpenseFormatted: string;
  totalIncomeFormatted: string;
  totalExpenseFormatted: string;
  netIncomeFormatted: string;
}

export const incomeApi = {
  getReport: async (month?: string, locationId?: string) => {
    const params: any = {};
    if (month) {
      params.month = month;
    }
    if (locationId) {
      params.location_id = locationId;
    }
    const response = await apiClient.get<IncomeReport>('/income', { params });
    return response.data;
  },
};
