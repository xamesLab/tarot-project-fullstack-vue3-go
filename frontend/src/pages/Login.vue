<script setup>
import { reactive, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

import BaseButton from '../components/base/BaseButton.vue'
import BaseInput from '../components/base/BaseInput.vue'
import AuthCard from '../components/form/AuthCard.vue'
import FormField from '../components/form/FormField.vue'
import PageContainer from '../components/layout/PageContainer.vue'
import { useAuthStore } from '../stores/auth'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const form = reactive({
  email: '',
  password: '',
})

const fieldErrors = reactive({
  email: '',
  password: '',
})

const apiError = ref('')

function validateEmail(value) {
  return /\S+@\S+\.\S+/.test(value)
}

function validateForm() {
  fieldErrors.email = ''
  fieldErrors.password = ''

  if (!form.email) {
    fieldErrors.email = t('auth.validation.required')
  } else if (!validateEmail(form.email)) {
    fieldErrors.email = t('auth.validation.email')
  }

  if (!form.password) {
    fieldErrors.password = t('auth.validation.required')
  }

  return !fieldErrors.email && !fieldErrors.password
}

async function submit() {
  apiError.value = ''

  if (!validateForm()) {
    return
  }

  const result = await authStore.login({
    email: form.email.trim(),
    password: form.password,
  })

  if (!result.ok) {
    apiError.value = result.error.message
    return
  }

  await router.push(typeof route.query.redirect === 'string' ? route.query.redirect : '/')
}
</script>

<template>
  <PageContainer tag="section">
    <div class="auth-page">
      <AuthCard :eyebrow="t('auth.loginEyebrow')" :title="t('auth.loginTitle')" :copy="t('auth.loginCopy')">
        <form class="stack-md" @submit.prevent="submit">
          <FormField :label="t('auth.email')" :error="fieldErrors.email" for-id="login-email" required>
            <BaseInput
              id="login-email"
              v-model="form.email"
              type="email"
              autocomplete="email"
              name="email"
              :placeholder="t('auth.emailPlaceholder')"
              :error="Boolean(fieldErrors.email)"
            />
          </FormField>

          <FormField :label="t('auth.password')" :error="fieldErrors.password" for-id="login-password" required>
            <BaseInput
              id="login-password"
              v-model="form.password"
              type="password"
              autocomplete="current-password"
              name="password"
              :placeholder="t('auth.passwordPlaceholder')"
              :error="Boolean(fieldErrors.password)"
            />
          </FormField>

          <p v-if="apiError" class="auth-page__error">{{ apiError }}</p>

          <BaseButton type="submit" full-width :loading="authStore.status === 'loading'">
            {{ t('auth.loginAction') }}
          </BaseButton>
        </form>

        <template #footer>
          <p>
            {{ t('auth.noAccount') }}
            <RouterLink class="text-link" to="/register">{{ t('auth.registerAction') }}</RouterLink>
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
</style>
