// Simple fetch wrapper for API calls

import { useAccountStore } from '@/stores/accountStore';

const API_BASE = '';

interface FetchOptions extends RequestInit {
  params?: Record<string, string>;
  authToken?: string;
}

export async function apiFetch<T>(
  endpoint: string,
  options: FetchOptions = {},
): Promise<T> {
  // Get token from Zustand store
  const { token } = useAccountStore.getState();
  
  const { params, authToken, ...fetchOptions } = options;

  const finalToken = authToken || token;

  let url = API_BASE + endpoint;
  if (params) {
    const searchParams = new URLSearchParams(params);
    url += `?${searchParams.toString()}`;
  }

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(finalToken && { 'X-Cube-Token': finalToken }),
    ...fetchOptions.headers,
  };

  const response = await fetch(url, {
    ...fetchOptions,
    headers,
  });

  if (response.status === 401) {
    if (!authToken) {
      useAccountStore.getState().clear();
      window.dispatchEvent(new Event('auth:unauthorized'));
    }
    throw new Error('Unauthorized: please configure your API token');
  }

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(
      `API error: ${response.status} ${response.statusText} - ${errorText}`,
    );
  }

  return response.json();
}
