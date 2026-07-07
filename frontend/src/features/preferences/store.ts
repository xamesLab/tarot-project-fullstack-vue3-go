import { create } from 'zustand';
import type { Language, Theme } from '@/api/types';

interface PreferencesState {
  theme: Theme;
  language: Language;
  toggleTheme: () => void;
  setTheme: (theme: Theme) => void;
  setLanguage: (language: Language) => void;
}

const storedTheme = localStorage.getItem('nebula-tarot:theme');
const storedLanguage = localStorage.getItem('nebula-tarot:lang');

export const usePreferencesStore = create<PreferencesState>((set) => ({
  theme: storedTheme === 'light' ? 'light' : 'dark',
  language: storedLanguage === 'en' ? 'en' : 'ru',
  toggleTheme: () => set((state) => {
    const theme = state.theme === 'dark' ? 'light' : 'dark';
    localStorage.setItem('nebula-tarot:theme', theme);
    return { theme };
  }),
  setTheme: (theme) => {
    localStorage.setItem('nebula-tarot:theme', theme);
    set({ theme });
  },
  setLanguage: (language) => {
    localStorage.setItem('nebula-tarot:lang', language);
    set({ language });
  },
}));
