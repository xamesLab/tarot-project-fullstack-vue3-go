import { useEffect, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { api } from '@/api';
import type { AdminUserDto, UserRole } from '@/api/types';
import { useAuthStore } from '@/features/auth/store';
import { usePreferencesStore } from '@/features/preferences/store';
import { Button } from '@/shared/ui/Button';
import { Field, SelectField } from '@/shared/ui/Field';
import { AdminState } from './AdminDashboard';
import { adminKeys, getAdminError } from './shared';

interface UserFilters {
  search: string;
  role: UserRole | 'all';
  status: 'active' | 'disabled' | 'all';
  page: number;
  limit: number;
}

const defaultFilters: UserFilters = { search: '', role: 'all', status: 'all', page: 1, limit: 20 };

export function AdminUsers() {
  const { t } = useTranslation(['admin', 'common']);
  const [draft, setDraft] = useState(defaultFilters);
  const [filters, setFilters] = useState(defaultFilters);
  const query = useQuery({ queryKey: adminKeys.users(filters), queryFn: () => api.admin.getUsers(filters) });
  const roleOptions = ['all', 'seeker', 'adept', 'master', 'keeper'].map((role) => ({ value: role, label: role === 'all' ? t('admin:users.allRoles') : t(`common:roles.${role}`) }));
  const statusOptions = ['all', 'active', 'disabled'].map((status) => ({ value: status, label: status === 'all' ? t('admin:users.allStatuses') : t(`admin:statuses.${status}`) }));

  return (
    <div className="admin-stack">
      <section className="admin-surface admin-filter-panel">
        <form onSubmit={(event) => { event.preventDefault(); setFilters({ ...draft, page: 1 }); }}>
          <Field label={t('admin:users.search')} placeholder={t('admin:users.searchPlaceholder')} value={draft.search} onChange={(event) => setDraft((current) => ({ ...current, search: event.target.value }))} />
          <SelectField label={t('admin:users.role')} options={roleOptions} value={draft.role} onChange={(event) => setDraft((current) => ({ ...current, role: event.target.value as UserFilters['role'] }))} />
          <SelectField label={t('admin:users.status')} options={statusOptions} value={draft.status} onChange={(event) => setDraft((current) => ({ ...current, status: event.target.value as UserFilters['status'] }))} />
          <div className="admin-filter-actions"><Button type="submit" size="sm" disabled={query.isFetching}>{t('admin:actions.apply')}</Button><Button type="button" variant="ghost" size="sm" onClick={() => { setDraft(defaultFilters); setFilters(defaultFilters); }}>{t('admin:actions.reset')}</Button></div>
        </form>
      </section>

      <section className="admin-surface admin-list-section">
        <div className="admin-section-head"><div><span>{t('admin:users.eyebrow')}</span><h2>{t('admin:users.title')}</h2><p>{t('admin:users.copy', { total: query.data?.pagination.total ?? 0 })}</p></div></div>
        {query.isPending ? <AdminState>{t('admin:states.loading')}</AdminState> : query.isError ? <AdminState error>{getAdminError(query.error, t('admin:states.error'))}</AdminState> : query.data.users.length ? (
          <div className="admin-user-list">{query.data.users.map((user) => <AdminUserCard key={user.id} user={user} />)}</div>
        ) : <div className="admin-empty">{t('admin:users.empty')}</div>}
        {query.data && <div className="admin-pagination"><Button variant="ghost" size="sm" disabled={filters.page <= 1 || query.isFetching} onClick={() => setFilters((current) => ({ ...current, page: current.page - 1 }))}>{t('admin:actions.previous')}</Button><span>{t('admin:pagination', { page: query.data.pagination.page, total: query.data.pagination.total_pages })}</span><Button variant="ghost" size="sm" disabled={filters.page >= query.data.pagination.total_pages || query.isFetching} onClick={() => setFilters((current) => ({ ...current, page: current.page + 1 }))}>{t('admin:actions.next')}</Button></div>}
      </section>
    </div>
  );
}

function AdminUserCard({ user }: { user: AdminUserDto }) {
  const { t } = useTranslation(['admin', 'common']);
  const language = usePreferencesStore((state) => state.language);
  const currentUserId = useAuthStore((state) => state.user?.id);
  const queryClient = useQueryClient();
  const [role, setRole] = useState(user.role);
  const [status, setStatus] = useState(user.status);
  const [feedback, setFeedback] = useState('');
  useEffect(() => { setRole(user.role); setStatus(user.status); }, [user.role, user.status]);

  const refresh = async () => {
    await Promise.all([queryClient.invalidateQueries({ queryKey: ['admin', 'users'] }), queryClient.invalidateQueries({ queryKey: adminKeys.dashboard })]);
  };
  const roleMutation = useMutation({ mutationFn: () => api.admin.updateUserRole(user.id, role), onSuccess: async () => { setFeedback(t('admin:users.roleSaved')); await refresh(); }, onError: (error) => { setRole(user.role); setFeedback(getAdminError(error, t('admin:states.error'))); } });
  const statusMutation = useMutation({ mutationFn: () => api.admin.updateUserStatus(user.id, status), onSuccess: async () => { setFeedback(t('admin:users.statusSaved')); await refresh(); }, onError: (error) => { setStatus(user.status); setFeedback(getAdminError(error, t('admin:states.error'))); } });
  const busy = roleMutation.isPending || statusMutation.isPending;
  const formatDate = (value?: string) => value ? new Intl.DateTimeFormat(language === 'ru' ? 'ru-RU' : 'en-US', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value)) : '—';
  const roleOptions = (['seeker', 'adept', 'master', 'keeper'] as const).map((value) => ({ value, label: t(`common:roles.${value}`) }));
  const statusOptions = (['active', 'disabled'] as const).map((value) => ({ value, label: t(`admin:statuses.${value}`) }));

  return (
    <article className="admin-user-card">
      <div className="admin-user-identity"><div><h3>{user.username || t('admin:users.noName')}</h3><p>{user.email}</p></div><span className={`admin-status admin-status-${user.status}`}>{t(`common:roles.${user.role}`)}</span></div>
      <dl className="admin-user-meta"><div><dt>{t('admin:users.status')}</dt><dd>{t(`admin:statuses.${user.status}`)}</dd></div><div><dt>{t('admin:users.timezone')}</dt><dd>{user.timezone || 'UTC'}</dd></div><div><dt>{t('admin:users.created')}</dt><dd>{formatDate(user.created_at)}</dd></div><div><dt>{t('admin:users.lastLogin')}</dt><dd>{formatDate(user.last_login_at)}</dd></div></dl>
      <div className="admin-user-controls">
        <SelectField label={t('admin:users.role')} options={roleOptions} value={role} disabled={busy} onChange={(event) => setRole(event.target.value as UserRole)} />
        <Button variant="outline" size="sm" disabled={busy || role === user.role} onClick={() => { setFeedback(''); roleMutation.mutate(); }}>{roleMutation.isPending ? t('admin:actions.saving') : t('admin:users.saveRole')}</Button>
        <SelectField label={t('admin:users.status')} options={statusOptions} value={status} disabled={busy} onChange={(event) => setStatus(event.target.value as AdminUserDto['status'])} />
        <Button variant="outline" size="sm" disabled={busy || status === user.status} onClick={() => { setFeedback(''); statusMutation.mutate(); }}>{statusMutation.isPending ? t('admin:actions.saving') : t('admin:users.saveStatus')}</Button>
      </div>
      {feedback && <p className={roleMutation.isError || statusMutation.isError ? 'admin-feedback error' : 'admin-feedback success'} role={roleMutation.isError || statusMutation.isError ? 'alert' : 'status'}>{feedback}</p>}
      {currentUserId === user.id && <p className="admin-self-hint">{t('admin:users.selfHint')}</p>}
    </article>
  );
}
