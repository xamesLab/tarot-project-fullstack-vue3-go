import { useEffect, useMemo, useState } from 'react';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { z } from 'zod';
import { ApiError } from '@/api/client';
import type { Language, Theme, UserDto } from '@/api/types';
import { Button } from '@/shared/ui/Button/Button';
import { SegmentedControl } from '@/shared/ui/Controls/Controls';
import { Field, SelectField } from '@/shared/ui/Field/Field';
import { useProfileCatalog, useUpdateProfile, useUpdateSettings } from './useProfileSettings';

const FALLBACK_TIMEZONES = [
  'Europe/Moscow',
  'Europe/Kaliningrad',
  'Europe/Samara',
  'Asia/Yekaterinburg',
  'Asia/Omsk',
  'Asia/Novosibirsk',
  'Asia/Irkutsk',
  'Asia/Yakutsk',
  'Asia/Vladivostok',
  'Asia/Magadan',
  'Asia/Kamchatka',
  'UTC',
] as const;

type AccountFormValues = {
  username: string;
  timezone: string;
};

function getTimezones(currentTimezone: string) {
  const intl = Intl as typeof Intl & { supportedValuesOf?: (key: 'timeZone') => string[] };
  let timezones = [...FALLBACK_TIMEZONES] as string[];
  try {
    timezones = intl.supportedValuesOf?.('timeZone') ?? timezones;
  } catch {
    // Older browsers use the compact fallback list.
  }
  if (currentTimezone && !timezones.includes(currentTimezone)) timezones.unshift(currentTimezone);
  return timezones;
}

export function ProfileSettings({ user }: { user: UserDto }) {
  const { t } = useTranslation(['common', 'profile']);
  const accountMutation = useUpdateProfile();
  const appearanceMutation = useUpdateSettings();
  const deckMutation = useUpdateSettings();
  const catalog = useProfileCatalog();
  const canCustomizeDeck = ['adept', 'master', 'keeper'].includes(user.role);
  const [accountSaved, setAccountSaved] = useState(false);
  const [appearanceSaved, setAppearanceSaved] = useState(false);
  const [deckSaved, setDeckSaved] = useState(false);
  const [language, setLanguage] = useState<Language>(user.language);
  const [theme, setTheme] = useState<Theme>(user.theme);
  const [selectedDeckId, setSelectedDeckId] = useState(user.selected_deck_id ?? '');
  const [selectedBackId, setSelectedBackId] = useState(user.selected_back_id ?? '');

  const accountSchema = useMemo(() => z.object({
    username: z.string()
      .trim()
      .min(3, t('profile:settings.validation.usernameLength'))
      .max(20, t('profile:settings.validation.usernameLength'))
      .regex(/^[a-zA-Z0-9._-]+$/, t('profile:settings.validation.usernameFormat')),
    timezone: z.string().min(1, t('profile:settings.validation.timezone')),
  }), [t]);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isDirty },
  } = useForm<AccountFormValues>({
    resolver: zodResolver(accountSchema),
    defaultValues: { username: user.username, timezone: user.timezone },
  });

  useEffect(() => {
    reset({ username: user.username, timezone: user.timezone });
  }, [reset, user.timezone, user.username]);

  useEffect(() => {
    setLanguage(user.language);
    setTheme(user.theme);
  }, [user.language, user.theme]);

  useEffect(() => {
    setSelectedDeckId(user.selected_deck_id ?? '');
    setSelectedBackId(user.selected_back_id ?? '');
  }, [user.selected_back_id, user.selected_deck_id]);

  const timezoneOptions = useMemo(
    () => getTimezones(user.timezone).map((timezone) => ({ value: timezone, label: timezone })),
    [user.timezone],
  );
  const decks = catalog.data?.decks.decks ?? [];
  const cardBacks = catalog.data?.cardBacks.card_backs ?? [];
  const defaultDeck = decks.find((deck) => deck.id === catalog.data?.decks.default_deck_id);
  const defaultBack = cardBacks.find((back) => back.id === catalog.data?.cardBacks.default_back_id);
  const deckOptions = [
    {
      value: '',
      label: t('profile:settings.deck.defaultOption', {
        name: defaultDeck?.name[user.language] ?? t('profile:settings.deck.systemDefault'),
      }),
    },
    ...decks.map((deck) => ({ value: deck.id, label: deck.name[user.language] || deck.name.en })),
  ];
  const cardBackOptions = [
    {
      value: '',
      label: t('profile:settings.deck.defaultOption', {
        name: defaultBack?.name[user.language] ?? t('profile:settings.deck.systemDefault'),
      }),
    },
    ...cardBacks.map((back) => ({ value: back.id, label: back.name[user.language] || back.name.en })),
  ];

  const mutationError = (error: unknown) => {
    if (error instanceof ApiError && error.details?.field === 'username') {
      return t('profile:settings.errors.usernameUnavailable');
    }
    return t('profile:settings.errors.save');
  };

  const saveAccount = handleSubmit(async (values) => {
    setAccountSaved(false);
    try {
      await accountMutation.mutateAsync(values);
      setAccountSaved(true);
    } catch {
      // The mutation error is rendered below the form.
    }
  });

  const saveAppearance = async () => {
    setAppearanceSaved(false);
    try {
      await appearanceMutation.mutateAsync({ language, theme });
      setAppearanceSaved(true);
    } catch {
      // The mutation error is rendered below the form.
    }
  };

  const saveDeck = async () => {
    if (!canCustomizeDeck) return;
    setDeckSaved(false);
    try {
      await deckMutation.mutateAsync({
        selected_deck_id: selectedDeckId,
        selected_back_id: selectedBackId,
      });
      setDeckSaved(true);
    } catch {
      // The mutation error is rendered below the form.
    }
  };

  return (
    <div className="profile-settings-stack">
      <section className="profile-settings-section" aria-labelledby="profile-account-settings-title">
        <header className="profile-settings-heading">
          <span aria-hidden="true">01</span>
          <div>
            <h2 id="profile-account-settings-title">{t('profile:settings.account.title')}</h2>
            <p>{t('profile:settings.account.description')}</p>
          </div>
        </header>
        <form className="profile-settings-form" onSubmit={saveAccount} aria-busy={accountMutation.isPending}>
          <div className="profile-settings-fields">
            <Field
              label={t('profile:settings.account.username')}
              autoComplete="username"
              error={errors.username?.message}
              {...register('username')}
            />
            <SelectField
              label={t('profile:settings.account.timezone')}
              options={timezoneOptions}
              error={errors.timezone?.message}
              {...register('timezone')}
            />
            <Field label={t('profile:settings.account.email')} value={user.email} disabled readOnly />
            <Field label={t('profile:settings.account.role')} value={t(`common:roles.${user.role}`)} disabled readOnly />
          </div>
          <div className="profile-settings-actions">
            <ProfileFeedback
              saved={accountSaved && !isDirty}
              error={accountMutation.error ? mutationError(accountMutation.error) : ''}
              savedText={t('profile:settings.account.saved')}
            />
            <Button type="submit" disabled={accountMutation.isPending}>
              {accountMutation.isPending ? t('profile:settings.saving') : t('profile:settings.account.save')}
            </Button>
          </div>
        </form>
      </section>

      <section className="profile-settings-section" aria-labelledby="profile-appearance-settings-title">
        <header className="profile-settings-heading">
          <span aria-hidden="true">02</span>
          <div>
            <h2 id="profile-appearance-settings-title">{t('profile:settings.appearance.title')}</h2>
            <p>{t('profile:settings.appearance.description')}</p>
          </div>
        </header>
        <div className="profile-choice-grid">
          <div className="profile-setting-choice">
            <span className="field-label">{t('profile:settings.appearance.language')}</span>
            <SegmentedControl
              label={t('profile:settings.appearance.language')}
              value={language}
              onChange={(value) => { setLanguage(value); setAppearanceSaved(false); }}
              options={[
                { value: 'ru', label: t('profile:settings.appearance.russian') },
                { value: 'en', label: t('profile:settings.appearance.english') },
              ]}
            />
          </div>
          <div className="profile-setting-choice">
            <span className="field-label">{t('profile:settings.appearance.theme')}</span>
            <SegmentedControl
              label={t('profile:settings.appearance.theme')}
              value={theme}
              onChange={(value) => { setTheme(value); setAppearanceSaved(false); }}
              options={[
                { value: 'dark', label: t('profile:settings.appearance.dark') },
                { value: 'light', label: t('profile:settings.appearance.light') },
              ]}
            />
          </div>
        </div>
        <div className="profile-settings-actions">
          <ProfileFeedback
            saved={appearanceSaved}
            error={appearanceMutation.error ? mutationError(appearanceMutation.error) : ''}
            savedText={t('profile:settings.appearance.saved')}
          />
          <Button type="button" variant="outline" disabled={appearanceMutation.isPending} onClick={saveAppearance}>
            {appearanceMutation.isPending ? t('profile:settings.saving') : t('profile:settings.appearance.save')}
          </Button>
        </div>
      </section>

      <section className="profile-settings-section" aria-labelledby="profile-deck-settings-title">
        <header className="profile-settings-heading">
          <span aria-hidden="true">03</span>
          <div>
            <h2 id="profile-deck-settings-title">{t('profile:settings.deck.title')}</h2>
            <p>{t('profile:settings.deck.description')}</p>
          </div>
        </header>
        {!canCustomizeDeck && <p className="profile-settings-access-note">{t('profile:settings.deck.restricted')}</p>}
        {catalog.isError && <p className="profile-settings-feedback error" role="alert">{t('profile:settings.deck.catalogError')}</p>}
        <div className="profile-settings-fields profile-settings-fields-deck">
          <SelectField
            label={t('profile:settings.deck.deck')}
            options={deckOptions}
            value={selectedDeckId}
            onChange={(event) => { setSelectedDeckId(event.target.value); setDeckSaved(false); }}
            disabled={!canCustomizeDeck || catalog.isPending || catalog.isError}
          />
          <SelectField
            label={t('profile:settings.deck.cardBack')}
            options={cardBackOptions}
            value={selectedBackId}
            onChange={(event) => { setSelectedBackId(event.target.value); setDeckSaved(false); }}
            disabled={!canCustomizeDeck || catalog.isPending || catalog.isError}
          />
        </div>
        <div className="profile-settings-actions">
          <ProfileFeedback
            saved={deckSaved}
            error={deckMutation.error ? mutationError(deckMutation.error) : ''}
            savedText={t('profile:settings.deck.saved')}
          />
          <Button
            type="button"
            variant="outline"
            disabled={!canCustomizeDeck || catalog.isPending || catalog.isError || deckMutation.isPending}
            onClick={saveDeck}
          >
            {deckMutation.isPending ? t('profile:settings.saving') : t('profile:settings.deck.save')}
          </Button>
        </div>
      </section>
    </div>
  );
}

function ProfileFeedback({ saved, error, savedText }: { saved: boolean; error: string; savedText: string }) {
  if (error) return <p className="profile-settings-feedback error" role="alert">{error}</p>;
  if (saved) return <p className="profile-settings-feedback success" role="status">{savedText}</p>;
  return null;
}
