/* eslint-disable react-refresh/only-export-components */
import { lazy, Suspense, type ReactNode } from 'react';
import { createBrowserRouter, Navigate } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';
import { HomePage } from '@/pages/home/HomePage';

const UiKitPage = import.meta.env.DEV ? lazy(() => import('@/pages/ui-kit/UiKitPage')) : null;
const ProfilePage = lazy(() => import('@/pages/profile/ProfilePage'));
const AdminPage = lazy(() => import('@/pages/admin/AdminPage'));
const SpreadsPage = lazy(() => import('@/pages/spreads/SpreadsPage'));

export function ProtectedRoute({ children }: { children: ReactNode }) {
  const hydrated = useAuthStore((state) => state.hydrated);
  const status = useAuthStore((state) => state.status);
  const token = useAuthStore((state) => state.token);
  const user = useAuthStore((state) => state.user);

  if (!hydrated || (status === 'loading' && token && !user)) {
    return <div className="route-loading">Loading…</div>;
  }

  if (!token || !user || status !== 'authenticated') {
    return <Navigate to="/" replace />;
  }

  return children;
}

export function KeeperRoute({ children }: { children: ReactNode }) {
  const hydrated = useAuthStore((state) => state.hydrated);
  const status = useAuthStore((state) => state.status);
  const token = useAuthStore((state) => state.token);
  const user = useAuthStore((state) => state.user);

  if (!hydrated || (status === 'loading' && token && !user)) return <div className="route-loading">Loading…</div>;
  if (!token || !user || status !== 'authenticated') return <Navigate to="/" replace />;
  if (user.role !== 'keeper') return <Navigate to="/profile" replace state={{ accessDenied: 'admin' }} />;
  return children;
}

const adminElement = (
  <KeeperRoute>
    <Suspense fallback={<div className="route-loading">Loading…</div>}><AdminPage /></Suspense>
  </KeeperRoute>
);

export const routes = [
  { path: '/', element: <HomePage /> },
  {
    path: '/profile',
    element: (
      <ProtectedRoute>
        <Suspense fallback={<div className="route-loading">Loading…</div>}>
          <ProfilePage />
        </Suspense>
      </ProtectedRoute>
    ),
  },
  {
    path: '/spreads',
    element: (
      <ProtectedRoute>
        <Suspense fallback={<div className="route-loading">Loading…</div>}>
          <SpreadsPage />
        </Suspense>
      </ProtectedRoute>
    ),
  },
  { path: '/admin', element: adminElement },
  { path: '/admin/users', element: <KeeperRoute><Navigate to="/admin?tab=users" replace /></KeeperRoute> },
  { path: '/admin/role-policies', element: <KeeperRoute><Navigate to="/admin?tab=policies" replace /></KeeperRoute> },
  { path: '/admin/settings', element: <KeeperRoute><Navigate to="/admin?tab=settings" replace /></KeeperRoute> },
  { path: '/admin/decks', element: <KeeperRoute><Navigate to="/admin?tab=assets" replace /></KeeperRoute> },
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

export const router = createBrowserRouter(routes);
