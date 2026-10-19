import { RouterProvider } from 'react-router-dom';
import { AppProviders } from '@/app/providers/AppProviders';
import { router as defaultRouter } from '@/app/router';

export function App({ router = defaultRouter }: { router?: typeof defaultRouter }) {
  return (
    <AppProviders>
      <RouterProvider router={router} />
    </AppProviders>
  );
}
