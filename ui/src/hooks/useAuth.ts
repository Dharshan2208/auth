import { useState, useEffect, useCallback } from 'react';
import type { AuthState, User } from '../types';
import * as api from '../api/client';

const ACCESS_KEY = 'auth_access_token';
const REFRESH_KEY = 'auth_refresh_token';
const USER_KEY = 'auth_user';

function loadPersistedAuth(): AuthState {
  const accessToken = localStorage.getItem(ACCESS_KEY);
  const refreshToken = localStorage.getItem(REFRESH_KEY);
  const userJson = localStorage.getItem(USER_KEY);

  if (accessToken && refreshToken && userJson) {
    try {
      const user = JSON.parse(userJson) as User;
      return { status: 'authenticated', user, accessToken, refreshToken };
    } catch {
      localStorage.removeItem(ACCESS_KEY);
      localStorage.removeItem(REFRESH_KEY);
      localStorage.removeItem(USER_KEY);
    }
  }
  return { status: 'unauthenticated' };
}

export function useAuth() {
  const [auth, setAuth] = useState<AuthState>({ status: 'loading' });

  useEffect(() => {
    setAuth(loadPersistedAuth());
  }, []);

  const persist = useCallback(
    (user: User, accessToken: string, refreshToken: string) => {
      localStorage.setItem(ACCESS_KEY, accessToken);
      localStorage.setItem(REFRESH_KEY, refreshToken);
      localStorage.setItem(USER_KEY, JSON.stringify(user));
      setAuth({ status: 'authenticated', user, accessToken, refreshToken });
    },
    [],
  );

  const clear = useCallback(() => {
    localStorage.removeItem(ACCESS_KEY);
    localStorage.removeItem(REFRESH_KEY);
    localStorage.removeItem(USER_KEY);
    setAuth({ status: 'unauthenticated' });
  }, []);

  const loginAction = useCallback(
    async (identifier: string, password: string) => {
      const tokens = await api.login({ identifier, password });
      const userData = await api.profile(tokens.access_token);
      const user: User = {
        id: userData.user_id,
        username: userData.username,
        email: userData.email,
        role: userData.role,
      };
      persist(user, tokens.access_token, tokens.refresh_token);
      return user;
    },
    [persist],
  );

  const signupAction = useCallback(
    async (username: string, email: string, password: string) => {
      await api.signup({ username, email, password });
    },
    [],
  );

  const logoutAction = useCallback(async () => {
    const state = auth;
    if (state.status === 'authenticated') {
      try {
        await api.logout({ refresh_token: state.refreshToken });
      } catch {
        // Best-effort logout
      }
    }
    clear();
  }, [auth, clear]);

  const refreshAction = useCallback(async () => {
    if (auth.status !== 'authenticated') return;
    try {
      const tokens = await api.refresh({ refresh_token: auth.refreshToken });
      const userData = await api.profile(tokens.access_token);
      const user: User = {
        id: userData.user_id,
        username: userData.username,
        email: userData.email,
        role: userData.role,
      };
      persist(user, tokens.access_token, tokens.refresh_token);
      return tokens;
    } catch {
      clear();
      throw new Error('Session expired. Please log in again.');
    }
  }, [auth, persist, clear]);

  return {
    auth,
    login: loginAction,
    signup: signupAction,
    logout: logoutAction,
    refreshTokens: refreshAction,
    accessToken: auth.status === 'authenticated' ? auth.accessToken : null,
    refreshToken: auth.status === 'authenticated' ? auth.refreshToken : null,
  };
}
