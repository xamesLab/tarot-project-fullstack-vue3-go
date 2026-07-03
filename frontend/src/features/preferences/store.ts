import { create } from 'zustand';
import type { Language } from '@/api/types';

type Theme = 'dark' | 'light';

interface PreferencesState {
  theme: Theme;
  language: Language;
  toggleTheme: () => void;
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
  setLanguage: (language) => {
    localStorage.setItem('nebula-tarot:lang', language);
    set({ language });
  },
}));
