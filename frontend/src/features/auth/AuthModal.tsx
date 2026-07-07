import { zodResolver } from '@hookform/resolvers/zod';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { z } from 'zod';
import { ApiError } from '@/api/client';
import type { LoginPayload, RegisterPayload } from '@/api/types';
import { Button } from '@/shared/ui/Button';
import { Field, SelectField } from '@/shared/ui/Field';
import { Modal } from '@/shared/ui/Modal';
import { Ornament } from '@/shared/ui/Ornament';
import { useAuthStore } from './store';

export type AuthModalMode = 'login' | 'register';

const fallbackTimezones = ['UTC', 'Europe/Moscow', 'Europe/Amsterdam', 'Europe/London', 'America/New_York', 'America/Los_Angeles', 'Asia/Tokyo'];
const timezoneApi = Intl as typeof Intl & { supportedValuesOf?: (key: 'timeZone') => string[] };
const timezoneOptions = Array.from(new Set(['UTC', ...(timezoneApi.supportedValuesOf?.('timeZone') ?? fallbackTimezones)])).map((timezone) => ({ value: timezone, label: timezone.split('_').join(' ') }));

export function AuthModal({ open, initialMode, onClose }: { open: boolean; initialMode: AuthModalMode; onClose: () => void }) {
  const { t } = useTranslation();
  const login = useAuthStore((state) => state.login);
  const registerUser = useAuthStore((state) => state.register);
  const [mode, setMode] = useState<AuthModalMode>(initialMode);
  const [success, setSuccess] = useState(false);
  const [apiError, setApiError] = useState('');
  const firstFieldRef = useRef<HTMLInputElement | null>(null);
  const closeTimerRef = useRef<number>();
  const defaultTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC';

  const loginSchema = useMemo(() => z.object({
    email: z.string().trim().min(1, t('auth.errors.required')).email(t('auth.errors.email')),
    password: z.string().min(1, t('auth.errors.required')),
  }), [t]);
  const registerSchema = useMemo(() => z.object({
    email: z.string().trim().min(1, t('auth.errors.required')).email(t('auth.errors.email')),
    username: z.string().trim().refine((value) => !value || /^[a-zA-Z0-9._-]{3,20}$/.test(value), t('auth.errors.username')),
    password: z.string().refine((value) => value.length >= 8 && /[A-Za-z]/.test(value) && /\d/.test(value), t('auth.errors.password')),
    timezone: z.string().min(1, t('auth.errors.timezone')),
  }), [t]);

  const loginForm = useForm<LoginPayload>({ resolver: zodResolver(loginSchema), defaultValues: { email: '', password: '' } });
  const registerForm = useForm<RegisterPayload>({ resolver: zodResolver(registerSchema), defaultValues: { email: '', username: '', password: '', timezone: defaultTimezone } });

  useEffect(() => {
    if (open) {
      setMode(initialMode);
      setSuccess(false);
      setApiError('');
      return;
    }
    loginForm.reset();
    registerForm.reset({ email: '', username: '', password: '', timezone: defaultTimezone });
    window.clearTimeout(closeTimerRef.current);
  }, [defaultTimezone, initialMode, loginForm, open, registerForm]);

  useEffect(() => () => window.clearTimeout(closeTimerRef.current), []);

  const errorMessage = (error: unknown) => {
    if (!(error instanceof ApiError)) return t('auth.errors.network');
    if (error.code === 'registration_disabled') return t('auth.errors.registrationDisabled');
    if (error.code === 'invalid_credentials') return t('auth.errors.invalidCredentials');
    if (error.code === 'account_disabled') return t('auth.errors.accountDisabled');
    const key = `${error.details?.field ?? ''}:${error.details?.reason ?? ''}`;
    if (key === 'email:already_registered') return t('auth.errors.emailTaken');
    if (key === 'username:unavailable') return t('auth.errors.usernameTaken');
    return error.status >= 500 ? t('auth.errors.server') : error.message;
  };

  const submitLogin = loginForm.handleSubmit(async (payload) => {
    setApiError('');
    try {
      await login(payload);
      onClose();
    } catch (error) {
      setApiError(errorMessage(error));
    }
  });

  const submitRegistration = registerForm.handleSubmit(async (payload) => {
    setApiError('');
    try {
      await registerUser({ ...payload, username: payload.username?.trim() || undefined });
      setSuccess(true);
      closeTimerRef.current = window.setTimeout(onClose, 1200);
    } catch (error) {
      setApiError(errorMessage(error));
    }
  });

  const switchMode = (nextMode: AuthModalMode) => {
    setMode(nextMode);
    setApiError('');
    window.setTimeout(() => firstFieldRef.current?.focus(), 0);
  };

  const formId = mode === 'login' ? 'login-form' : 'registration-form';
  const isSubmitting = mode === 'login' ? loginForm.formState.isSubmitting : registerForm.formState.isSubmitting;
  const title = success ? t('auth.register.welcome') : t(`auth.${mode}.title`);

  return <Modal open={open} onClose={onClose} title={title} initialFocusRef={success ? undefined : firstFieldRef} footer={success ? <Button onClick={onClose}>{t('auth.continue')}</Button> : <><Button variant="ghost" onClick={onClose}>{t('auth.later')}</Button><Button type="submit" form={formId} disabled={isSubmitting}>{isSubmitting ? t(`auth.${mode}.submitting`) : t(`auth.${mode}.submit`)}</Button></>}>
    {success ? <div className="welcome-state"><Ornament/><p>{t('auth.register.welcomeText')}</p></div> : <><p className="registration-lead">{t(`auth.${mode}.lead`)}</p>{mode === 'login' ? <form id="login-form" className="registration-form" onSubmit={submitLogin} noValidate><Field label={t('auth.fields.email')} type="email" autoComplete="email" {...loginForm.register('email')} ref={(node) => { loginForm.register('email').ref(node); firstFieldRef.current = node; }} error={loginForm.formState.errors.email?.message}/><Field label={t('auth.fields.password')} type="password" autoComplete="current-password" {...loginForm.register('password')} error={loginForm.formState.errors.password?.message}/></form> : <form id="registration-form" className="registration-form" onSubmit={submitRegistration} noValidate><Field label={t('auth.fields.email')} type="email" autoComplete="email" {...registerForm.register('email')} ref={(node) => { registerForm.register('email').ref(node); firstFieldRef.current = node; }} error={registerForm.formState.errors.email?.message}/><Field label={t('auth.fields.username')} autoComplete="username" {...registerForm.register('username')} error={registerForm.formState.errors.username?.message}/><Field label={t('auth.fields.password')} type="password" autoComplete="new-password" {...registerForm.register('password')} error={registerForm.formState.errors.password?.message}/><SelectField label={t('auth.fields.timezone')} options={timezoneOptions} {...registerForm.register('timezone')} error={registerForm.formState.errors.timezone?.message}/></form>}{apiError && <p className="auth-api-error" role="alert">{apiError}</p>}<button type="button" className="auth-mode-switch" onClick={() => switchMode(mode === 'login' ? 'register' : 'login')}>{t(`auth.${mode}.switch`)}</button></>}
  </Modal>;
}
