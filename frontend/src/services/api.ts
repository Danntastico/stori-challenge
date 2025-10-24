import axios, { AxiosInstance } from 'axios';

// API Configuration
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';

// Types
export interface HealthResponse {
  status: string;
  timestamp: string;
}

export interface Transaction {
  date: string;
  amount: number;
  category: string;
  description: string;
  type: 'income' | 'expense';
}

export interface Period {
  start: string;
  end: string;
  months?: number;
}

export interface TransactionsResponse {
  transactions: Transaction[];
  count: number;
  period: Period;
}

export interface CategoryDetail {
  total: number;
  count: number;
  percentage: number;
}

export interface FinancialSummary {
  total_income: number;
  total_expenses: number;
  net_savings: number;
  savings_rate: number;
}

export interface CategorySummaryResponse {
  income: Record<string, CategoryDetail>;
  expenses: Record<string, CategoryDetail>;
  summary: FinancialSummary;
  period: Period;
}

export interface TimelinePoint {
  period: string;
  income: number;
  expenses: number;
  net: number;
}

export interface TimelineResponse {
  timeline: TimelinePoint[];
  aggregation: string;
}

export interface AIAdviceRequest {
  context: string;
  category?: string;
}

export interface AIAdviceResponse {
  advice: string;
  insights: string[];
  recommendations: string[];
  timestamp: string;
}

// Create axios instance with base configuration
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for debugging
apiClient.interceptors.request.use(
  (config) => {
    console.log(`[API] ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    console.error('[API] Request error:', error);
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => {
    console.log(`[API] Response:`, response.status, response.data);
    return response;
  },
  (error) => {
    console.error('[API] Response error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

// API Methods

/**
 * Check API health
 */
export const checkHealth = async (): Promise<HealthResponse> => {
  const response = await apiClient.get<HealthResponse>('/health');
  return response.data;
};

/**
 * Get all transactions with optional filters
 */
export const getTransactions = async (params?: {
  startDate?: string;
  endDate?: string;
  type?: 'income' | 'expense';
  category?: string;
}): Promise<TransactionsResponse> => {
  const response = await apiClient.get<TransactionsResponse>('/transactions', { params });
  return response.data;
};

/**
 * Get category summary
 */
export const getCategorySummary = async (): Promise<CategorySummaryResponse> => {
  const response = await apiClient.get<CategorySummaryResponse>('/summary/categories');
  return response.data;
};

/**
 * Get timeline data
 */
export const getTimeline = async (): Promise<TimelineResponse> => {
  const response = await apiClient.get<TimelineResponse>('/summary/timeline');
  return response.data;
};

/**
 * Get AI financial advice
 */
export const getAIAdvice = async (request: AIAdviceRequest): Promise<AIAdviceResponse> => {
  const response = await apiClient.post<AIAdviceResponse>('/advice', request);
  return response.data;
};

export default apiClient;

