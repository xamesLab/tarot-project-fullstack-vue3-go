import { create } from 'zustand';
import { api } from '@/api';
import { ApiError, TOKEN_STORAGE_KEY } from '@/api/client';
import type { AuthResponse, LoginPayload, RegisterPayload, UserDto } from '@/api/types';
import { usePreferencesStore } from '@/features/preferences/store';

export type AuthStatus = 'idle' | 'loading' | 'authenticated' | 'guest' | 'error';

interface AuthState {
  user: UserDto | null;
  token: string;
  status: AuthStatus;
  hydrated: boolean;
  login: (payload: LoginPayload) => Promise<UserDto>;
  register: (payload: RegisterPayload) => Promise<UserDto>;
  hydrateSession: () => Promise<UserDto | null>;
  logout: () => void;
}

const storedToken = localStorage.getItem(TOKEN_STORAGE_KEY) || '';
let hydrationPromise: Promise<UserDto | null> | null = null;

function syncPreferences(user: UserDto) {
  const preferences = usePreferencesStore.getState();
  preferences.setLanguage(user.language);
  preferences.setTheme(user.theme);
}

function persistToken(token: string) {
  if (token) localStorage.setItem(TOKEN_STORAGE_KEY, token);
  else localStorage.removeItem(TOKEN_STORAGE_KEY);
}

export const useAuthStore = create<AuthState>((set, get) => {
  const authenticate = async (request: () => Promise<AuthResponse>): Promise<UserDto> => {
    set({ status: 'loading' });
    try {
      const response = await request();
      persistToken(response.token);
      syncPreferences(response.user);
      set({ user: response.user, token: response.token, status: 'authenticated', hydrated: true });
      return response.user;
    } catch (error) {
      set({ status: 'error' });
      throw error;
    }
  };

  return {
    user: null,
    token: storedToken,
    status: storedToken ? 'idle' : 'guest',
    hydrated: false,

    login: (payload) => authenticate(() => api.auth.login(payload)),

    register: (payload) => authenticate(() => api.auth.register(payload)),

    async hydrateSession() {
      if (get().hydrated) return get().user;
      if (hydrationPromise) return hydrationPromise;
      if (!get().token) {
        set({ hydrated: true, status: 'guest' });
        return null;
      }

      set({ status: 'loading' });
      hydrationPromise = api.auth.getCurrentUser()
        .then((response) => {
          syncPreferences(response.user);
          set({ user: response.user, status: 'authenticated', hydrated: true });
          return response.user;
        })
        .catch((error: unknown) => {
          if (error instanceof ApiError && (error.status === 401 || error.status === 403)) {
            persistToken('');
            set({ user: null, token: '', status: 'guest', hydrated: true });
          } else {
            set({ user: null, status: 'error', hydrated: true });
          }
          return null;
        })
        .finally(() => { hydrationPromise = null; });
      return hydrationPromise;
    },

    logout() {
      persistToken('');
      set({ user: null, token: '', status: 'guest', hydrated: true });
    },
  };
});
