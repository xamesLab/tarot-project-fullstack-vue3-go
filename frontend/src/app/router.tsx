/* eslint-disable react-refresh/only-export-components */
import { lazy, Suspense } from 'react';
import { createBrowserRouter, Navigate } from 'react-router-dom';
import { HomePage } from '@/pages/home/HomePage';

const UiKitPage = import.meta.env.DEV ? lazy(() => import('@/pages/ui-kit/UiKitPage')) : null;

export const routes = [
  { path: '/', element: <HomePage /> },
  ...(UiKitPage
    ? [
        {
          path: '/ui-kit',
          element: (
            <Suspense fallback={<div className="route-loading">Loading UI Kit…</div>}>
              <UiKitPage />
            </Suspense>
          ),
        },
      ]
    : []),
  { path: '*', element: <Navigate to="/" replace /> },
];

export const router = createBrowserRouter(routes, { future: { v7_relativeSplatPath: true } });
