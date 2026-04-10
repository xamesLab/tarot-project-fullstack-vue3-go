<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

import BaseButton from '../components/base/BaseButton.vue'
import BaseInput from '../components/base/BaseInput.vue'
import BaseSelect from '../components/base/BaseSelect.vue'
import AuthCard from '../components/form/AuthCard.vue'
import FormField from '../components/form/FormField.vue'
import PageContainer from '../components/layout/PageContainer.vue'
import { useAuthStore } from '../stores/auth'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const FALLBACK_TIMEZONES = [
  'UTC',
  'Europe/Moscow',
  'Europe/Amsterdam',
  'Europe/London',
  'America/New_York',
  'America/Los_Angeles',
  'Asia/Tokyo',
]

const form = reactive({
  email: '',
  username: '',
  password: '',
  timezone: '',
})

const fieldErrors = reactive({
  email: '',
  username: '',
  password: '',
  timezone: '',
})

const apiError = ref('')

const timezoneOptions = computed(() => {
  const values = typeof Intl.supportedValuesOf === 'function'
    ? Intl.supportedValuesOf('timeZone')
    : FALLBACK_TIMEZONES

  return values.map((timezone) => ({
    value: timezone,
    label: timezone.replaceAll('_', ' '),
  }))
})

function validateEmail(value) {
  return /\S+@\S+\.\S+/.test(value)
}

function validatePassword(value) {
  return value.length >= 8 && /[A-Za-z]/.test(value) && /\d/.test(value)
}

function validateUsername(value) {
  return /^[a-zA-Z0-9._-]{3,20}$/.test(value)
}

function validateForm() {
  fieldErrors.email = ''
  fieldErrors.username = ''
  fieldErrors.password = ''
  fieldErrors.timezone = ''

  if (!form.email) {
    fieldErrors.email = t('auth.validation.required')
  } else if (!validateEmail(form.email)) {
    fieldErrors.email = t('auth.validation.email')
  }

  if (form.username && !validateUsername(form.username)) {
    fieldErrors.username = t('auth.validation.username')
  }

  if (!form.password) {
    fieldErrors.password = t('auth.validation.required')
  } else if (!validatePassword(form.password)) {
    fieldErrors.password = t('auth.validation.password')
  }

  if (!form.timezone) {
    fieldErrors.timezone = t('auth.validation.required')
  }

  return !fieldErrors.email && !fieldErrors.username && !fieldErrors.password && !fieldErrors.timezone
}

async function submit() {
  apiError.value = ''

  if (!validateForm()) {
    return
  }

  const result = await authStore.register({
    email: form.email.trim(),
    username: form.username.trim(),
    password: form.password,
    timezone: form.timezone,
  })

  if (!result.ok) {
    apiError.value = result.error.message
    return
  }

  await router.push(typeof route.query.redirect === 'string' ? route.query.redirect : '/')
}

onMounted(() => {
  if (!form.timezone) {
    form.timezone = Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC'
  }
})
</script>

<template>
  <PageContainer tag="section">
    <div class="auth-page">
      <AuthCard :eyebrow="t('auth.registerEyebrow')" :title="t('auth.registerTitle')" :copy="t('auth.registerCopy')">
        <form class="stack-md" @submit.prevent="submit">
          <FormField :label="t('auth.email')" :error="fieldErrors.email" for-id="register-email" required>
            <BaseInput
              id="register-email"
              v-model="form.email"
              type="email"
              autocomplete="email"
              name="email"
              :placeholder="t('auth.emailPlaceholder')"
              :error="Boolean(fieldErrors.email)"
            />
          </FormField>

          <FormField
            :label="t('auth.username')"
            :hint="t('auth.usernameOptionalHint')"
            :error="fieldErrors.username"
            for-id="register-username"
          >
            <BaseInput
              id="register-username"
              v-model="form.username"
              autocomplete="off"
              name="register_username"
              :placeholder="t('auth.usernamePlaceholder')"
              :error="Boolean(fieldErrors.username)"
            />
          </FormField>

          <FormField :label="t('auth.password')" :error="fieldErrors.password" for-id="register-password" required>
            <BaseInput
              id="register-password"
              v-model="form.password"
              type="password"
              autocomplete="new-password"
              name="new-password"
              :placeholder="t('auth.passwordPlaceholder')"
              :error="Boolean(fieldErrors.password)"
            />
          </FormField>

          <FormField
            :label="t('auth.timezone')"
            :hint="t('auth.timezoneHint')"
            :error="fieldErrors.timezone"
            for-id="register-timezone"
            required
          >
            <BaseSelect
              id="register-timezone"
              v-model="form.timezone"
              :options="timezoneOptions"
              :placeholder="t('auth.timezonePlaceholder')"
              :error="Boolean(fieldErrors.timezone)"
            />
          </FormField>

          <p v-if="apiError" class="auth-page__error">{{ apiError }}</p>

          <BaseButton type="submit" full-width :loading="authStore.status === 'loading'">
            {{ t('auth.registerAction') }}
          </BaseButton>
        </form>

        <template #footer>
          <p>
            {{ t('auth.haveAccount') }}
            <RouterLink class="text-link" to="/login">{{ t('auth.loginAction') }}</RouterLink>
          </p>
        </template>
      </AuthCard>
    </div>
  </PageContainer>
</template>

<style scoped>
.auth-page {
  display: grid;
  min-height: calc(100vh - var(--ds-header-height) - 180px);
  place-items: center;
}

.auth-page__error {
  color: var(--ds-color-danger);
  font-family: var(--ds-font-ui);
  font-size: var(--ds-text-body-sm);
}

@media (max-width: 720px) {
  .auth-page {
    min-height: auto;
    place-items: start center;
    padding-block: var(--ds-space-4);
  }
}
</style>
