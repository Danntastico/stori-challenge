import { useState, useEffect, useCallback } from 'react';

interface UseFetchState<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
}

interface UseFetchOptions<T> {
  immediate?: boolean; // Whether to fetch immediately on mount
  onSuccess?: (data: T) => void;
  onError?: (error: Error) => void;
}

/**
 * Custom hook for data fetching with loading and error states
 * @param fetchFn - Async function that fetches data
 * @param options - Configuration options
 * @returns Object containing data, loading state, error, and refetch function
 */
export function useFetch<T>(
  fetchFn: () => Promise<T>,
  options: UseFetchOptions<T> = {}
): UseFetchState<T> & { refetch: () => Promise<void> } {
  const { immediate = true, onSuccess, onError } = options;

  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState<boolean>(immediate);
  const [error, setError] = useState<string | null>(null);

  const executeFetch = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await fetchFn();
      setData(result);
      onSuccess?.(result);
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : 'An unexpected error occurred';
      setError(errorMessage);
      console.error('Fetch error:', err);
      onError?.(err as Error);
    } finally {
      setLoading(false);
    }
  }, [fetchFn, onSuccess, onError]);

  useEffect(() => {
    if (immediate) {
      executeFetch();
    }
  }, [immediate, executeFetch]);

  return {
    data,
    loading,
    error,
    refetch: executeFetch,
  };
}

/**
 * Hook specifically for API data fetching with built-in retry logic
 * @param fetchFn - Async function that fetches data
 * @param dependencies - Dependencies array for refetching
 */
export function useApiData<T>(
  fetchFn: () => Promise<T>,
  dependencies: unknown[] = []
): UseFetchState<T> & { refetch: () => Promise<void> } {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const executeFetch = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await fetchFn();
      setData(result);
    } catch (err) {
      console.error('API fetch error:', err);
      setError('Failed to load data. Please try again.');
    } finally {
      setLoading(false);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, dependencies);

  useEffect(() => {
    executeFetch();
  }, [executeFetch]);

  return {
    data,
    loading,
    error,
    refetch: executeFetch,
  };
}

