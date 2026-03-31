export interface BaseEntity {
  id: string;
  createdAt: string;
  updatedAt: string;
}

export interface Location {
  id: string;
  shortName: string;
  detail: string;
  createdAt: string;
  updatedAt: string;
}

export interface Room {
  id: string;
  locationId: string;
  roomNumber: string;
  status: string;
  tags: string[];
  note?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Landlord {
  id: string;
  name: string;
  phone: string;
  note: string;
  createdAt: string;
  updatedAt: string;
}

export interface Lease {
  id: string;
  roomId: string;
  landlordId: string;
  tenantName: string;
  tenantPhone: string;
  startDate: string;
  endDate: string;
  rentAmount: number;
  depositAmount: number;
  status: string;
  note: string;
  lastChargeAt: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface Bill {
  id: string;
  leaseId: string;
  type: string;
  status: string;
  amount: number;
  rentAmount: number;
  waterAmount: number;
  electricAmount: number;
  otherAmount: number;
  refundDepositAmount: number;
  billStart: string;
  billEnd: string;
  dueDate: string;
  paidAt: string | null;
  note: string;
  createdAt: string;
  updatedAt: string;
}

export interface PrintJob {
  id: string;
  type: string;
  status: string;
  content: string;
  createdAt: string;
  updatedAt: string;
}

export interface Deposit {
  id: string;
  leaseId: string;
  amount: number;
  status: string;
  refundedAt?: string;
  deductedAt?: string;
  note: string;
  createdAt: string;
  updatedAt: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: any;
}
