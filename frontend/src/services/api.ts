// Simple fetch wrapper for API calls

const API_BASE = '';
const TOKEN_KEY = 'cubeToken';
const PROJECT_ID_KEY = 'cubeProjectId';

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token);
}

export function removeToken(): void {
  localStorage.removeItem(TOKEN_KEY);
}

export function getProjectId(): number | null {
  const id = localStorage.getItem(PROJECT_ID_KEY);
  return id ? parseInt(id, 10) : null;
}

export function setProjectId(id: number): void {
  localStorage.setItem(PROJECT_ID_KEY, id.toString());
}

export function removeProjectId(): void {
  localStorage.removeItem(PROJECT_ID_KEY);
}

export function isConfigured(): boolean {
  return getToken() !== null;
}

interface FetchOptions extends RequestInit {
  params?: Record<string, string>;
  authToken?: string;
}

export async function apiFetch<T>(
  endpoint: string,
  options: FetchOptions = {},
): Promise<T> {
  const storedToken = getToken();
  const { params, authToken, ...fetchOptions } = options;

  const token = authToken || storedToken;

  let url = API_BASE + endpoint;
  if (params) {
    const searchParams = new URLSearchParams(params);
    url += `?${searchParams.toString()}`;
  }

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(token && { 'X-Cube-Token': token }),
    ...fetchOptions.headers,
  };

  const response = await fetch(url, {
    ...fetchOptions,
    headers,
  });

  if (response.status === 401) {
    if (!authToken) {
      removeToken();
      removeProjectId();
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
