import { create } from 'zustand';
import type { UserDto } from '@/api/types';

interface AuthState {
  user: UserDto | null;
  setUser: (user: UserDto) => void;
}

export const useAuthStore = create<AuthState>((set) => ({ user: null, setUser: (user) => set({ user }) }));
