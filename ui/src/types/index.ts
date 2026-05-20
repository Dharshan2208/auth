// ─── Request Types ──────────────────────────────────────────

export interface SignupRequest {
  username: string;
  email: string;
  password: string;
}

export interface LoginRequest {
  identifier: string;
  password: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface LogoutRequest {
  refresh_token: string;
}

export interface ChangePasswordRequest {
  old_password: string;
  confirm_old_password: string;
  new_password: string;
}

// ─── Response Types ─────────────────────────────────────────

export interface TokenPairResponse {
  access_token: string;
  refresh_token: string;
}

export interface ProfileResponse {
  user_id: number;
  username: string;
  email: string;
  role: string;
}

export interface MessageResponse {
  message: string;
}

export interface ErrorResponse {
  error: string;
}

export interface HealthResponse {
  status: string;
}

// ─── Derived Types ──────────────────────────────────────────

export interface User {
  id: number;
  username: string;
  email: string;
  role: string;
}

export type AuthState =
  | { status: 'loading' }
  | { status: 'unauthenticated' }
  | { status: 'authenticated'; user: User; accessToken: string; refreshToken: string };

export type Page = 'home' | 'dashboard';
