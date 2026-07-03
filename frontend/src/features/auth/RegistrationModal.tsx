import { zodResolver } from '@hookform/resolvers/zod';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { z } from 'zod';
import { api } from '@/api';
import type { RegisterPayload } from '@/api/types';
import { Button } from '@/shared/ui/Button';
import { Field } from '@/shared/ui/Field';
import { Modal } from '@/shared/ui/Modal';
import { Ornament } from '@/shared/ui/Ornament';
import { useAuthStore } from './store';

export function RegistrationModal({ open, onClose }: { open: boolean; onClose: () => void }) {
  const { t } = useTranslation();
  const user = useAuthStore((state) => state.user);
  const setUser = useAuthStore((state) => state.setUser);
  const firstFieldRef = useRef<HTMLInputElement | null>(null);
  const [success, setSuccess] = useState(false);
  const schema = useMemo(() => z.object({ name: z.string().trim().min(1, t('registration.errors.name')), email: z.string().email(t('registration.errors.email')), password: z.string().min(8, t('registration.errors.password')), birthDate: z.string().optional() }), [t]);
  const { register, handleSubmit, reset, formState: { errors, isSubmitting } } = useForm<RegisterPayload>({ resolver: zodResolver(schema), defaultValues: { name: '', email: '', password: '', birthDate: '' } });
  const nameRegistration = register('name');

  useEffect(() => { if (!open) { reset(); setSuccess(false); } }, [open, reset]);

  const submit = handleSubmit(async (payload) => {
    const created = await api.auth.register({ ...payload, birthDate: payload.birthDate || undefined });
    setUser(created);
    setSuccess(true);
    window.setTimeout(onClose, 1200);
  });

  const welcome = Boolean(user) || success;
  return <Modal open={open} onClose={onClose} title={welcome ? t('registration.welcome') : t('registration.title')} initialFocusRef={welcome ? undefined : firstFieldRef} footer={welcome ? <Button onClick={onClose}>{t('registration.close')}</Button> : <><Button variant="ghost" onClick={onClose}>{t('registration.later')}</Button><Button type="submit" form="registration-form" disabled={isSubmitting}>{isSubmitting ? t('registration.submitting') : t('registration.create')}</Button></>}>
    {welcome ? <div className="welcome-state"><Ornament/><p>{t('registration.welcomeText')}</p></div> : <><p className="registration-lead">{t('registration.lead')}</p><form id="registration-form" className="registration-form" onSubmit={submit} noValidate><Field label={t('registration.name')} autoComplete="name" {...nameRegistration} ref={(node) => { nameRegistration.ref(node); firstFieldRef.current = node; }} error={errors.name?.message}/><Field label={t('registration.email')} type="email" autoComplete="email" {...register('email')} error={errors.email?.message}/><Field label={t('registration.password')} type="password" autoComplete="new-password" {...register('password')} error={errors.password?.message}/><Field label={t('registration.birthDate')} type="date" {...register('birthDate')} error={errors.birthDate?.message}/></form><button type="button" className="sign-in-stub">{t('registration.existing')}</button></>}
  </Modal>;
}
