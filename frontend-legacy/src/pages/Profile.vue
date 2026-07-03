<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import * as catalogService from '../services/catalog'
import { normalizeApiError } from '../services/api'
import BaseButton from '../components/base/BaseButton.vue'
import BaseInput from '../components/base/BaseInput.vue'
import BaseSegmentedToggle from '../components/base/BaseSegmentedToggle.vue'
import BaseSelect from '../components/base/BaseSelect.vue'
import BaseSurface from '../components/base/BaseSurface.vue'
import BaseTabs from '../components/base/BaseTabs.vue'
import FormField from '../components/form/FormField.vue'
import PageContainer from '../components/layout/PageContainer.vue'
import PageSection from '../components/layout/PageSection.vue'
import SectionHeader from '../components/layout/SectionHeader.vue'
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

const { t, locale } = useI18n()
const authStore = useAuthStore()
const activeProfileTab = ref('account')

const profileForm = reactive({
  username: '',
  timezone: '',
})

const settingsForm = reactive({
  language: 'en',
  theme: 'light',
})

const deckSettingsForm = reactive({
  selectedDeckID: '',
  selectedBackID: '',
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

const deckSettingsState = reactive({
  error: '',
  success: '',
})

const catalogState = reactive({
  status: 'idle',
  error: '',
})

const availableDecks = ref([])
const availableCardBacks = ref([])
const defaultDeckID = ref('default')
const defaultBackID = ref('default')

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

const profileTabs = computed(() => [
  { value: 'account', label: t('profile.accountTitle') },
  { value: 'preferences', label: t('profile.preferencesTitle') },
  { value: 'deck-settings', label: t('profile.deckSettingsTitle') },
])

const canCustomizeDecks = computed(() => ['adept', 'master', 'keeper'].includes(authStore.user?.role || ''))

watch(
  () => authStore.user,
  (user) => {
    profileForm.username = user?.username || ''
    profileForm.timezone = user?.timezone || 'UTC'
    settingsForm.language = user?.language || 'en'
    settingsForm.theme = user?.theme || 'light'
    deckSettingsForm.selectedDeckID = user?.selected_deck_id || ''
    deckSettingsForm.selectedBackID = user?.selected_back_id || ''
  },
  { immediate: true },
)

function localizeText(value) {
  if (!value) {
    return ''
  }

  return value[locale.value] || value.en || value.ru || ''
}

const defaultDeckName = computed(() => {
  const deck = availableDecks.value.find((item) => item.id === defaultDeckID.value)
  return localizeText(deck?.name) || t('profile.systemDefaultFallback')
})

const defaultBackName = computed(() => {
  const cardBack = availableCardBacks.value.find((item) => item.id === defaultBackID.value)
  return localizeText(cardBack?.name) || t('profile.systemDefaultFallback')
})

const deckOptions = computed(() => [
  { value: '', label: t('profile.useSystemDefaultDeck', { name: defaultDeckName.value }) },
  ...availableDecks.value.map((deck) => ({
    value: deck.id,
    label: localizeText(deck.name),
  })),
])

const cardBackOptions = computed(() => [
  { value: '', label: t('profile.useSystemDefaultBack', { name: defaultBackName.value }) },
  ...availableCardBacks.value.map((cardBack) => ({
    value: cardBack.id,
    label: localizeText(cardBack.name),
  })),
])

function validateProfile() {
  profileErrors.username = ''
  profileErrors.timezone = ''

  if (!profileForm.username) {
    profileErrors.username = t('auth.validation.required')
  } else if (!/^[a-zA-Z0-9._-]{3,20}$/.test(profileForm.username)) {
    profileErrors.username = t('auth.validation.username')
  }

  if (!profileForm.timezone) {
    profileErrors.timezone = t('auth.validation.required')
  }

  return !profileErrors.username && !profileErrors.timezone
}

async function loadCatalogOptions() {
  catalogState.status = 'loading'
  catalogState.error = ''

  try {
    const [deckResponse, cardBackResponse] = await Promise.all([
      catalogService.getAvailableDecks(),
      catalogService.getAvailableCardBacks(),
    ])

    availableDecks.value = deckResponse.decks || []
    availableCardBacks.value = cardBackResponse.card_backs || []
    defaultDeckID.value = deckResponse.default_deck_id || 'default'
    defaultBackID.value = cardBackResponse.default_back_id || 'default'
    catalogState.status = 'ready'
  } catch (error) {
    const normalized = normalizeApiError(error)
    catalogState.status = 'error'
    catalogState.error = normalized.message
  }
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

async function submitDeckSettings() {
  deckSettingsState.error = ''
  deckSettingsState.success = ''

  if (!canCustomizeDecks.value) {
    deckSettingsState.error = t('profile.deckSettingsRestricted')
    return
  }

  const result = await authStore.updateSettings({
    selected_deck_id: deckSettingsForm.selectedDeckID,
    selected_back_id: deckSettingsForm.selectedBackID,
  })

  if (!result.ok) {
    deckSettingsState.error = result.error.message
    return
  }

  deckSettingsState.success = t('profile.deckSettingsSaved')
}

onMounted(() => {
  loadCatalogOptions()
})
</script>

<template>
  <PageContainer>
    <div class="page-stack">
      <PageSection :eyebrow="t('profile.eyebrow')" :title="t('profile.title')" :copy="t('profile.copy')" />

      <BaseSurface padding="lg" tag="section">
        <BaseTabs v-model="activeProfileTab" :items="profileTabs">
          <template #panel="{ activeItem }">
            <div v-if="activeItem?.value === 'account'" class="profile-panel">
              <SectionHeader
                class="profile-head"
                size="sm"
                title-tag="h3"
                :title="t('profile.accountTitle')"
                :copy="t('profile.accountCopy')"
              />

              <form class="stack-md" @submit.prevent="submitProfile">
                <FormField :label="t('profile.email')" :hint="t('profile.emailReadonly')">
                  <BaseInput :model-value="authStore.user?.email || ''" disabled />
                </FormField>

                <FormField :label="t('profile.role')" :hint="t('profile.roleReadonly')">
                  <BaseInput :model-value="authStore.user?.role || ''" disabled />
                </FormField>

                <FormField :label="t('profile.username')" :error="profileErrors.username" for-id="profile-username" required>
                  <BaseInput
                    id="profile-username"
                    v-model="profileForm.username"
                    autocomplete="username"
                    :placeholder="t('auth.usernamePlaceholder')"
                    :error="Boolean(profileErrors.username)"
                  />
                </FormField>

                <FormField :label="t('profile.timezone')" :error="profileErrors.timezone" for-id="profile-timezone" required>
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

            <div v-else-if="activeItem?.value === 'preferences'" class="profile-panel">
              <SectionHeader
                class="profile-head"
                size="sm"
                title-tag="h3"
                :title="t('profile.preferencesTitle')"
                :copy="t('profile.preferencesCopy')"
              />

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

            <div v-else class="profile-panel">
              <SectionHeader
                class="profile-head"
                size="sm"
                title-tag="h3"
                :title="t('profile.deckSettingsTitle')"
                :copy="t('profile.deckSettingsCopy')"
              />

              <div v-if="catalogState.error" class="form-message form-message--error">
                {{ catalogState.error }}
              </div>

              <form class="stack-md" @submit.prevent="submitDeckSettings">
                <FormField
                  :label="t('profile.selectedDeck')"
                  :hint="t('profile.selectedDeckHint', { name: defaultDeckName })"
                  for-id="profile-selected-deck"
                >
                  <BaseSelect
                    id="profile-selected-deck"
                    v-model="deckSettingsForm.selectedDeckID"
                    :options="deckOptions"
                    :disabled="!canCustomizeDecks || catalogState.status === 'loading'"
                  />
                </FormField>

                <FormField
                  :label="t('profile.selectedBack')"
                  :hint="t('profile.selectedBackHint', { name: defaultBackName })"
                  for-id="profile-selected-back"
                >
                  <BaseSelect
                    id="profile-selected-back"
                    v-model="deckSettingsForm.selectedBackID"
                    :options="cardBackOptions"
                    :disabled="!canCustomizeDecks || catalogState.status === 'loading'"
                  />
                </FormField>

                <p v-if="!canCustomizeDecks" class="text-body-muted">
                  {{ t('profile.deckSettingsRestricted') }}
                </p>

                <p v-if="deckSettingsState.error" class="form-message form-message--error">{{ deckSettingsState.error }}</p>
                <p v-else-if="deckSettingsState.success" class="form-message form-message--success">{{ deckSettingsState.success }}</p>

                <BaseButton
                  type="submit"
                  variant="secondary"
                  :disabled="!canCustomizeDecks || catalogState.status === 'loading'"
                  :loading="authStore.status === 'loading'"
                >
                  {{ t('profile.saveDeckSettings') }}
                </BaseButton>
              </form>
            </div>
          </template>
        </BaseTabs>
      </BaseSurface>
    </div>
  </PageContainer>
</template>

<style scoped>
.profile-head {
  margin-bottom: 24px;
}

.profile-panel {
  max-width: 760px;
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
