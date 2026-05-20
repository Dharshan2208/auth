import type {
  SignupRequest,
  LoginRequest,
  RefreshTokenRequest,
  LogoutRequest,
  ChangePasswordRequest,
  TokenPairResponse,
  ProfileResponse,
  MessageResponse,
  HealthResponse,
} from '../types';

const BASE = '/api/v1';

class ApiError extends Error {
  status: number;
  body: string;

  constructor(status: number, body: string) {
    super(`API error ${status}: ${body}`);
    this.status = status;
    this.body = body;
    this.name = 'ApiError';
  }
}

async function request<T>(
  method: string,
  path: string,
  body?: unknown,
  token?: string,
): Promise<T> {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const res = await fetch(`${BASE}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  const text = await res.text();

  if (!res.ok) {
    throw new ApiError(res.status, text);
  }

  return JSON.parse(text) as T;
}

function get<T>(path: string, token?: string): Promise<T> {
  return request<T>('GET', path, undefined, token);
}

function post<T>(path: string, body?: unknown, token?: string): Promise<T> {
  return request<T>('POST', path, body, token);
}

// ─── Public Endpoints ───────────────────────────────────────

export function healthCheck(): Promise<HealthResponse> {
  return get<HealthResponse>('/health');
}

export function signup(data: SignupRequest): Promise<MessageResponse> {
  return post<MessageResponse>('/signup', data);
}

export function login(data: LoginRequest): Promise<TokenPairResponse> {
  return post<TokenPairResponse>('/login', data);
}

export function refresh(data: RefreshTokenRequest): Promise<TokenPairResponse> {
  return post<TokenPairResponse>('/refresh', data);
}

export function logout(data: LogoutRequest): Promise<MessageResponse> {
  return post<MessageResponse>('/logout', data);
}

// ─── Authenticated Endpoints ────────────────────────────────

export function profile(token: string): Promise<ProfileResponse> {
  return get<ProfileResponse>('/profile', token);
}

export function adminCheck(token: string): Promise<MessageResponse> {
  return get<MessageResponse>('/admin', token);
}

export function changePassword(
  data: ChangePasswordRequest,
  token: string,
): Promise<MessageResponse> {
  return post<MessageResponse>('/password/change', data, token);
}

export { ApiError };
