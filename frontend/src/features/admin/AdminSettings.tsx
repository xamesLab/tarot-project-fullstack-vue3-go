import { useEffect, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { api } from '@/api';
import type { UpdateAdminSettingsPayload } from '@/api/types';
import { usePreferencesStore } from '@/features/preferences/store';
import { Button } from '@/shared/ui/Button';
import { Toggle } from '@/shared/ui/Controls';
import { SelectField, TextareaField } from '@/shared/ui/Field';
import { AdminState } from './AdminDashboard';
import { adminKeys, getAdminError } from './shared';

const emptySettings: UpdateAdminSettingsPayload = {
  maintenance_mode: false,
  maintenance_message: { en: '', ru: '' },
  ai_enabled: true,
  default_role: 'seeker',
  default_deck_id: '',
  default_back_id: '',
  registration_enabled: true,
};

export function AdminSettings() {
  const { t } = useTranslation(['admin', 'common']);
  const language = usePreferencesStore((state) => state.language);
  const queryClient = useQueryClient();
  const settingsQuery = useQuery({ queryKey: adminKeys.settings, queryFn: api.admin.getSettings });
  const decksQuery = useQuery({ queryKey: adminKeys.decks('active', 1, 100), queryFn: () => api.admin.getDecks({ scope: 'active', page: 1, limit: 100 }) });
  const backsQuery = useQuery({ queryKey: adminKeys.cardBacks('active', 1, 100), queryFn: () => api.admin.getCardBacks({ scope: 'active', page: 1, limit: 100 }) });
  const [form, setForm] = useState(emptySettings);
  useEffect(() => {
    const settings = settingsQuery.data?.settings;
    if (!settings) return;
    setForm({
      maintenance_mode: settings.maintenance_mode,
      maintenance_message: settings.maintenance_message,
      ai_enabled: settings.ai_enabled,
      default_role: settings.default_role,
      default_deck_id: settings.default_deck_id,
      default_back_id: settings.default_back_id,
      registration_enabled: settings.registration_enabled,
    });
  }, [settingsQuery.data]);
  const mutation = useMutation({ mutationFn: () => api.admin.updateSettings(form), onSuccess: (response) => queryClient.setQueryData(adminKeys.settings, response) });

  if (settingsQuery.isPending) return <AdminState>{t('admin:states.loading')}</AdminState>;
  if (settingsQuery.isError) return <AdminState error>{getAdminError(settingsQuery.error, t('admin:states.error'))}</AdminState>;

  const localize = (value: { ru: string; en: string }) => value[language] || value.en || value.ru;
  const deckOptions = (decksQuery.data?.decks || []).filter((deck) => deck.status === 'done').map((deck) => ({ value: deck.id, label: localize(deck.name) }));
  const backOptions = (backsQuery.data?.card_backs || []).filter((back) => back.status === 'done').map((back) => ({ value: back.id, label: localize(back.name) }));
  const roleOptions = (['seeker', 'adept', 'master'] as const).map((role) => ({ value: role, label: t(`common:roles.${role}`) }));
  const setBoolean = (key: 'maintenance_mode' | 'registration_enabled' | 'ai_enabled', value: boolean) => setForm((current) => ({ ...current, [key]: value }));

  return (
    <form className="admin-settings-form" onSubmit={(event) => { event.preventDefault(); mutation.mutate(); }}>
      <div className="admin-settings-grid">
        <section className="admin-surface admin-settings-card">
          <div className="admin-section-head"><div><span>{t('admin:settings.systemEyebrow')}</span><h2>{t('admin:settings.systemTitle')}</h2><p>{t('admin:settings.systemCopy')}</p></div></div>
          <div className="admin-settings-toggles">
            <SettingToggle title={t('admin:settings.maintenance')} copy={t('admin:settings.maintenanceHint')} checked={form.maintenance_mode} onChange={(value) => setBoolean('maintenance_mode', value)} />
            <SettingToggle title={t('admin:settings.registration')} copy={t('admin:settings.registrationHint')} checked={form.registration_enabled} onChange={(value) => setBoolean('registration_enabled', value)} />
            <SettingToggle title={t('admin:settings.ai')} copy={t('admin:settings.aiHint')} checked={form.ai_enabled} onChange={(value) => setBoolean('ai_enabled', value)} />
          </div>
          <div className="admin-settings-selects">
            <SelectField label={t('admin:settings.defaultRole')} options={roleOptions} value={form.default_role} onChange={(event) => setForm((current) => ({ ...current, default_role: event.target.value as UpdateAdminSettingsPayload['default_role'] }))} />
            <SelectField label={t('admin:settings.defaultDeck')} options={deckOptions.length ? deckOptions : [{ value: form.default_deck_id, label: t('admin:states.loading') }]} value={form.default_deck_id} disabled={decksQuery.isPending || !deckOptions.length} onChange={(event) => setForm((current) => ({ ...current, default_deck_id: event.target.value }))} />
            <SelectField label={t('admin:settings.defaultBack')} options={backOptions.length ? backOptions : [{ value: form.default_back_id, label: t('admin:states.loading') }]} value={form.default_back_id} disabled={backsQuery.isPending || !backOptions.length} onChange={(event) => setForm((current) => ({ ...current, default_back_id: event.target.value }))} />
          </div>
        </section>
        <section className="admin-surface admin-settings-card">
          <div className="admin-section-head"><div><span>{t('admin:settings.messageEyebrow')}</span><h2>{t('admin:settings.messageTitle')}</h2><p>{t('admin:settings.messageCopy')}</p></div></div>
          <TextareaField id="admin-maintenance-ru" label={t('admin:settings.messageRu')} value={form.maintenance_message.ru} onChange={(event) => setForm((current) => ({ ...current, maintenance_message: { ...current.maintenance_message, ru: event.target.value } }))} />
          <TextareaField id="admin-maintenance-en" label={t('admin:settings.messageEn')} value={form.maintenance_message.en} onChange={(event) => setForm((current) => ({ ...current, maintenance_message: { ...current.maintenance_message, en: event.target.value } }))} />
          <p className="admin-form-hint">{t('admin:settings.messageHint')}</p>
        </section>
      </div>
      <div className="admin-form-actions">
        <div>{mutation.isError && <p className="admin-feedback error" role="alert">{getAdminError(mutation.error, t('admin:states.error'))}</p>}{mutation.isSuccess && <p className="admin-feedback success" role="status">{t('admin:settings.saved')}</p>}</div>
        <Button type="submit" disabled={mutation.isPending}>{mutation.isPending ? t('admin:actions.saving') : t('admin:settings.save')}</Button>
      </div>
    </form>
  );
}

function SettingToggle({ title, copy, checked, onChange }: { title: string; copy: string; checked: boolean; onChange: (value: boolean) => void }) {
  return <div className="admin-toggle-field"><div><span>{title}</span><small>{copy}</small></div><Toggle checked={checked} onChange={onChange} label={title} /></div>;
}
