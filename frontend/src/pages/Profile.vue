<script setup>
import { computed, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseButton from '../components/base/BaseButton.vue'
import BaseInput from '../components/base/BaseInput.vue'
import BaseSegmentedToggle from '../components/base/BaseSegmentedToggle.vue'
import BaseSelect from '../components/base/BaseSelect.vue'
import BaseSurface from '../components/base/BaseSurface.vue'
import FormField from '../components/form/FormField.vue'
import PageSection from '../components/layout/PageSection.vue'
import { useAuthStore } from '../stores/auth'

const FALLBACK_TIMEZONES = [
  'UTC',
  'Europe/Moscow',
  'Europe/Amsterdam',
  'Europe/London',
  'America/New_York',
  'America/Los_Angeles',
  'Asia/Tokyo',
]

const { t } = useI18n()
const authStore = useAuthStore()

const profileForm = reactive({
  username: '',
  timezone: '',
})

const settingsForm = reactive({
  language: 'en',
  theme: 'light',
})

const profileErrors = reactive({
  username: '',
  timezone: '',
})

const settingsState = reactive({
  error: '',
  success: '',
})

const profileState = reactive({
  error: '',
  success: '',
})

const timezoneOptions = computed(() => {
  const values = typeof Intl.supportedValuesOf === 'function'
    ? Intl.supportedValuesOf('timeZone')
    : FALLBACK_TIMEZONES

  return values.map((timezone) => ({
    value: timezone,
    label: timezone.replaceAll('_', ' '),
  }))
})

const languageOptions = computed(() => [
  { label: 'EN', value: 'en' },
  { label: 'RU', value: 'ru' },
])

const themeOptions = computed(() => [
  { label: t('controls.light'), value: 'light' },
  { label: t('controls.dark'), value: 'dark' },
])

watch(
  () => authStore.user,
  (user) => {
    profileForm.username = user?.username || ''
    profileForm.timezone = user?.timezone || 'UTC'
    settingsForm.language = user?.language || 'en'
    settingsForm.theme = user?.theme || 'light'
  },
  { immediate: true },
)

function validateProfile() {
  profileErrors.username = ''
  profileErrors.timezone = ''

  if (!profileForm.username) {
    profileErrors.username = t('auth.validation.required')
  } else if (!/^[a-zA-Z0-9]{3,20}$/.test(profileForm.username)) {
    profileErrors.username = t('auth.validation.username')
  }

  if (!profileForm.timezone) {
    profileErrors.timezone = t('auth.validation.required')
  }

  return !profileErrors.username && !profileErrors.timezone
}

async function submitProfile() {
  profileState.error = ''
  profileState.success = ''

  if (!validateProfile()) {
    return
  }

  const result = await authStore.updateProfile({
    username: profileForm.username.trim(),
    timezone: profileForm.timezone,
  })

  if (!result.ok) {
    profileState.error = result.error.message
    return
  }

  profileState.success = t('profile.profileSaved')
}

async function submitSettings() {
  settingsState.error = ''
  settingsState.success = ''

  const result = await authStore.updateSettings({
    language: settingsForm.language,
    theme: settingsForm.theme,
  })

  if (!result.ok) {
    settingsState.error = result.error.message
    return
  }

  settingsState.success = t('profile.settingsSaved')
}
</script>

<template>
  <div class="page-stack">
    <PageSection :eyebrow="t('profile.eyebrow')" :title="t('profile.title')" :copy="t('profile.copy')" />

    <div class="page-grid page-grid--two">
      <BaseSurface padding="lg" tag="section">
        <div class="stack-md">
          <div class="profile-head">
            <div>
              <h3 class="section-head__title section-head__title--small">{{ t('profile.accountTitle') }}</h3>
              <p class="text-body-muted">{{ t('profile.accountCopy') }}</p>
            </div>
          </div>

          <form class="stack-md" @submit.prevent="submitProfile">
            <FormField :label="t('profile.email')" :hint="t('profile.emailReadonly')">
              <BaseInput :model-value="authStore.user?.email || ''" disabled />
            </FormField>

            <FormField :label="t('profile.role')" :hint="t('profile.roleReadonly')">
              <BaseInput :model-value="authStore.user?.role || ''" disabled />
            </FormField>

            <FormField :label="t('profile.username')" :error="profileErrors.username" for-id="profile-username">
              <BaseInput
                id="profile-username"
                v-model="profileForm.username"
                autocomplete="username"
                :placeholder="t('auth.usernamePlaceholder')"
                :error="Boolean(profileErrors.username)"
              />
            </FormField>

            <FormField :label="t('profile.timezone')" :error="profileErrors.timezone" for-id="profile-timezone">
              <BaseSelect
                id="profile-timezone"
                v-model="profileForm.timezone"
                :options="timezoneOptions"
                :placeholder="t('auth.timezonePlaceholder')"
                :error="Boolean(profileErrors.timezone)"
              />
            </FormField>

            <p v-if="profileState.error" class="form-message form-message--error">{{ profileState.error }}</p>
            <p v-else-if="profileState.success" class="form-message form-message--success">{{ profileState.success }}</p>

            <BaseButton type="submit" :loading="authStore.status === 'loading'">
              {{ t('profile.saveProfile') }}
            </BaseButton>
          </form>
        </div>
      </BaseSurface>

      <BaseSurface padding="lg" tag="section">
        <div class="stack-md">
          <div class="profile-head">
            <div>
              <h3 class="section-head__title section-head__title--small">{{ t('profile.preferencesTitle') }}</h3>
              <p class="text-body-muted">{{ t('profile.preferencesCopy') }}</p>
            </div>
          </div>

          <form class="stack-md" @submit.prevent="submitSettings">
            <FormField :label="t('profile.language')">
              <BaseSegmentedToggle
                v-model="settingsForm.language"
                :options="languageOptions"
                :aria-label="t('controls.language')"
              />
            </FormField>

            <FormField :label="t('profile.theme')">
              <BaseSegmentedToggle
                v-model="settingsForm.theme"
                :options="themeOptions"
                :aria-label="t('controls.theme')"
              />
            </FormField>

            <p v-if="settingsState.error" class="form-message form-message--error">{{ settingsState.error }}</p>
            <p v-else-if="settingsState.success" class="form-message form-message--success">{{ settingsState.success }}</p>

            <BaseButton type="submit" variant="secondary" :loading="authStore.status === 'loading'">
              {{ t('profile.saveSettings') }}
            </BaseButton>
          </form>
        </div>
      </BaseSurface>
    </div>
  </div>
</template>

<style scoped>
.section-head__title--small {
  font-size: var(--ds-text-heading-2);
  color: var(--ds-color-accent-primary);
}

.profile-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.form-message {
  font-family: var(--ds-font-ui);
  font-size: var(--ds-text-body-sm);
}

.form-message--error {
  color: var(--ds-color-danger);
}

.form-message--success {
  color: var(--ds-color-accent-primary);
}
</style>
