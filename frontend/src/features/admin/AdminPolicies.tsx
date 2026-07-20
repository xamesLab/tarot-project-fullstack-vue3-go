import { useEffect, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { api } from '@/api';
import type { AdminRolePolicyDto, UpdateAdminRolePolicyPayload } from '@/api/types';
import { Button } from '@/shared/ui/Button';
import { Toggle } from '@/shared/ui/Controls';
import { Field } from '@/shared/ui/Field';
import { AdminState } from './AdminDashboard';
import { adminKeys, getAdminError } from './shared';

export function AdminPolicies() {
  const { t } = useTranslation('admin');
  const query = useQuery({ queryKey: adminKeys.policies, queryFn: api.admin.getRolePolicies });
  if (query.isPending) return <AdminState>{t('states.loading')}</AdminState>;
  if (query.isError) return <AdminState error>{getAdminError(query.error, t('states.error'))}</AdminState>;
  return <div className="admin-policy-grid">{query.data.policies.map((policy) => <PolicyCard key={policy.role} policy={policy} />)}</div>;
}

function PolicyCard({ policy }: { policy: AdminRolePolicyDto }) {
  const { t } = useTranslation(['admin', 'common']);
  const queryClient = useQueryClient();
  const [form, setForm] = useState<UpdateAdminRolePolicyPayload>({ three_card_limit: policy.three_card_limit, ai_enabled: policy.ai_enabled, ai_daily_limit: policy.ai_daily_limit, history_retention_days: policy.history_retention_days });
  const [validation, setValidation] = useState('');
  useEffect(() => setForm({ three_card_limit: policy.three_card_limit, ai_enabled: policy.ai_enabled, ai_daily_limit: policy.ai_daily_limit, history_retention_days: policy.history_retention_days }), [policy]);
  const mutation = useMutation({
    mutationFn: () => api.admin.updateRolePolicy(policy.role, form),
    onSuccess: (response) => queryClient.setQueryData(adminKeys.policies, (current: { success: true; policies: AdminRolePolicyDto[] } | undefined) => current ? { ...current, policies: current.policies.map((item) => item.role === policy.role ? response.policy : item) } : current),
  });
  const setNumber = (key: 'three_card_limit' | 'ai_daily_limit' | 'history_retention_days', value: string) => setForm((current) => ({ ...current, [key]: Number(value) }));
  const submit = () => {
    if ([form.three_card_limit, form.ai_daily_limit, form.history_retention_days].some((value) => !Number.isFinite(value) || value < 0)) { setValidation(t('admin:policies.validation')); return; }
    setValidation('');
    mutation.mutate();
  };

  return (
    <section className="admin-surface admin-policy-card">
      <header><div><span>{t('admin:policies.eyebrow')}</span><h2>{t(`common:roles.${policy.role}`)}</h2></div><span className="admin-role-badge">{policy.role}</span></header>
      <form onSubmit={(event) => { event.preventDefault(); submit(); }}>
        <Field label={t('admin:policies.threeCardLimit')} type="number" min="0" value={form.three_card_limit} onChange={(event) => setNumber('three_card_limit', event.target.value)} />
        <div className="admin-toggle-field"><div><span>{t('admin:policies.aiEnabled')}</span><small>{t('admin:policies.aiEnabledHint')}</small></div><Toggle checked={form.ai_enabled} onChange={(ai_enabled) => setForm((current) => ({ ...current, ai_enabled }))} label={t('admin:policies.aiEnabled')} /></div>
        <Field label={t('admin:policies.aiDailyLimit')} type="number" min="0" value={form.ai_daily_limit} onChange={(event) => setNumber('ai_daily_limit', event.target.value)} />
        <Field label={t('admin:policies.retention')} type="number" min="0" value={form.history_retention_days} onChange={(event) => setNumber('history_retention_days', event.target.value)} />
        <p className="admin-form-hint">{t('admin:policies.hint')}</p>
        {(validation || mutation.isError) && <p className="admin-feedback error" role="alert">{validation || getAdminError(mutation.error, t('admin:states.error'))}</p>}
        {mutation.isSuccess && !mutation.isError && <p className="admin-feedback success" role="status">{t('admin:policies.saved')}</p>}
        <Button type="submit" disabled={mutation.isPending}>{mutation.isPending ? t('admin:actions.saving') : t('admin:actions.save')}</Button>
      </form>
    </section>
  );
}
