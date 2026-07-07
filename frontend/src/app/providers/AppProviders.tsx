import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { type ReactNode, useEffect, useState } from 'react';
import { I18nextProvider } from 'react-i18next';
import { i18n } from '@/features/i18n/i18n';
import { usePreferencesStore } from '@/features/preferences/store';
import { useAuthStore } from '@/features/auth/store';
import { AUTH_UNAUTHORIZED_EVENT } from '@/api/client';

export function AppProviders({ children }: { children: ReactNode }) {
  const [queryClient] = useState(
    () => new QueryClient({ defaultOptions: { queries: { staleTime: 30 * 60 * 1000, retry: 1 } } }),
  );
  const theme = usePreferencesStore((state) => state.theme);
  const language = usePreferencesStore((state) => state.language);

  useEffect(() => {
    void useAuthStore.getState().hydrateSession();
    const handleUnauthorized = () => useAuthStore.getState().logout();
    window.addEventListener(AUTH_UNAUTHORIZED_EVENT, handleUnauthorized);
    return () => window.removeEventListener(AUTH_UNAUTHORIZED_EVENT, handleUnauthorized);
  }, []);

  useEffect(() => {
    document.documentElement.classList.toggle('light', theme === 'light');
    document.documentElement.dataset.theme = theme;
  }, [theme]);

  useEffect(() => {
    document.documentElement.lang = language;
    void i18n.changeLanguage(language);
  }, [language]);

  return (
    <I18nextProvider i18n={i18n}>
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    </I18nextProvider>
  );
}
