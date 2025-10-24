/**
 * Format a number as currency (USD)
 * @param {number} amount - The amount to format
 * @param {boolean} showDecimals - Whether to show decimal places
 * @returns {string} Formatted currency string
 */
export const formatCurrency = (amount, showDecimals = true) => {
  if (amount === null || amount === undefined || isNaN(amount)) {
    return '$0.00';
  }

  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: showDecimals ? 2 : 0,
    maximumFractionDigits: showDecimals ? 2 : 0,
  }).format(amount);
};

/**
 * Format a date string
 * @param {string} dateString - ISO date string (YYYY-MM-DD)
 * @param {string} format - Format type ('short', 'long', 'month-year')
 * @returns {string} Formatted date string
 */
export const formatDate = (dateString, format = 'short') => {
  if (!dateString) return '';

  const date = new Date(dateString);
  
  // Check for invalid date
  if (isNaN(date.getTime())) {
    return dateString;
  }

  switch (format) {
    case 'long':
      return new Intl.DateTimeFormat('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      }).format(date);
    
    case 'month-year':
      return new Intl.DateTimeFormat('en-US', {
        year: 'numeric',
        month: 'short',
      }).format(date);
    
    case 'month':
      return new Intl.DateTimeFormat('en-US', {
        month: 'short',
      }).format(date);
    
    case 'short':
    default:
      return new Intl.DateTimeFormat('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
      }).format(date);
  }
};

/**
 * Format a percentage
 * @param {number} value - The percentage value
 * @param {number} decimals - Number of decimal places
 * @returns {string} Formatted percentage string
 */
export const formatPercentage = (value, decimals = 1) => {
  if (value === null || value === undefined || isNaN(value)) {
    return '0%';
  }

  return `${value.toFixed(decimals)}%`;
};

/**
 * Format a number with commas
 * @param {number} value - The number to format
 * @returns {string} Formatted number string
 */
export const formatNumber = (value) => {
  if (value === null || value === undefined || isNaN(value)) {
    return '0';
  }

  return new Intl.NumberFormat('en-US').format(value);
};

/**
 * Parse month-year string (YYYY-MM) to readable format
 * @param {string} monthStr - Month string in YYYY-MM format
 * @returns {string} Formatted month string
 */
export const parseMonthYear = (monthStr) => {
  if (!monthStr) return '';
  
  const [year, month] = monthStr.split('-');
  const date = new Date(year, parseInt(month) - 1);
  
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
  }).format(date);
};

/**
 * Capitalize first letter of string
 * @param {string} str - String to capitalize
 * @returns {string} Capitalized string
 */
export const capitalize = (str) => {
  if (!str) return '';
  return str.charAt(0).toUpperCase() + str.slice(1);
};

/**
 * Get color for transaction type
 * @param {string} type - Transaction type ('income' or 'expense')
 * @returns {string} Tailwind color class
 */
export const getTypeColor = (type) => {
  return type === 'income' ? 'text-green-600' : 'text-red-600';
};

/**
 * Get background color for transaction type
 * @param {string} type - Transaction type ('income' or 'expense')
 * @returns {string} Tailwind background color class
 */
export const getTypeBgColor = (type) => {
  return type === 'income' ? 'bg-green-50' : 'bg-red-50';
};

