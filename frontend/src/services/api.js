import axios from 'axios';

// API base URL - will be configured via environment variable
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';

// Create axios instance with default config
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000, // 10 seconds
});

// Request interceptor (for future auth tokens, etc.)
apiClient.interceptors.request.use(
  (config) => {
    // Could add auth token here in the future
    // config.headers.Authorization = `Bearer ${token}`;
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor (for global error handling)
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    // Handle specific error cases
    if (error.response) {
      // Server responded with error status
      console.error('API Error:', error.response.status, error.response.data);
    } else if (error.request) {
      // Request made but no response
      console.error('Network Error:', error.message);
    } else {
      // Something else happened
      console.error('Error:', error.message);
    }
    return Promise.reject(error);
  }
);

// API Service Methods

/**
 * Health Check
 * GET /health
 */
export const checkHealth = async () => {
  const response = await apiClient.get('/health');
  return response.data;
};

/**
 * Get All Transactions
 * GET /transactions
 * @param {Object} params - Query parameters
 * @param {string} params.startDate - ISO date string (YYYY-MM-DD)
 * @param {string} params.endDate - ISO date string (YYYY-MM-DD)
 * @param {string} params.type - Transaction type ('income' or 'expense')
 * @param {string} params.category - Category name
 */
export const getTransactions = async (params = {}) => {
  const response = await apiClient.get('/transactions', { params });
  return response.data;
};

/**
 * Get Category Summary
 * GET /summary/categories
 */
export const getCategorySummary = async () => {
  const response = await apiClient.get('/summary/categories');
  return response.data;
};

/**
 * Get Timeline Data
 * GET /summary/timeline
 */
export const getTimeline = async () => {
  const response = await apiClient.get('/summary/timeline');
  return response.data;
};

/**
 * Get AI Financial Advice
 * POST /advice
 * @param {Object} data - Request body
 * @param {string} data.context - Context for advice ('general', 'savings', etc.)
 * @param {string} data.category - Optional category for specific advice
 */
export const getFinancialAdvice = async (data = {}) => {
  const response = await apiClient.post('/advice', data);
  return response.data;
};

// Export default API client for custom requests
export default apiClient;

