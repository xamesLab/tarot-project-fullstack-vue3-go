<script setup>
import { computed, onMounted, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import AdminLayout from '../../components/admin/AdminLayout.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseInput from '../../components/base/BaseInput.vue'
import BaseSelect from '../../components/base/BaseSelect.vue'
import BaseSurface from '../../components/base/BaseSurface.vue'
import BaseSwitchToggle from '../../components/base/BaseSwitchToggle.vue'
import FormField from '../../components/form/FormField.vue'
import { useAdminStore } from '../../stores/admin'

const { t, locale } = useI18n()
const adminStore = useAdminStore()

const form = reactive({
  maintenance_mode: false,
  registration_enabled: true,
  ai_enabled: true,
  default_role: 'seeker',
  default_deck_id: 'default',
  default_back_id: 'default',
  maintenance_message_en: '',
  maintenance_message_ru: '',
})

function localizeText(value) {
  if (!value) {
    return ''
  }

  return value[locale.value] || value.en || value.ru || ''
}

const roleOptions = computed(() => [
  { value: 'seeker', label: t('roles.seeker') },
  { value: 'adept', label: t('roles.adept') },
  { value: 'master', label: t('roles.master') },
  { value: 'keeper', label: t('roles.keeper') },
])

const deckOptions = computed(() =>
  adminStore.deckList
    .filter((deck) => deck.status === 'done')
    .map((deck) => ({
      value: deck.id,
      label: localizeText(deck.name),
    })),
)

const cardBackOptions = computed(() =>
  adminStore.cardBackList
    .filter((cardBack) => cardBack.status === 'done')
    .map((cardBack) => ({
      value: cardBack.id,
      label: localizeText(cardBack.name),
    })),
)

watch(
  () => adminStore.settings,
  (settings) => {
    if (!settings) {
      return
    }

    form.maintenance_mode = Boolean(settings.maintenance_mode)
    form.registration_enabled = Boolean(settings.registration_enabled)
    form.ai_enabled = settings.ai_enabled !== false
    form.default_role = settings.default_role || 'seeker'
    form.default_deck_id = settings.default_deck_id || 'default'
    form.default_back_id = settings.default_back_id || 'default'
    form.maintenance_message_en = settings.maintenance_message?.en || ''
    form.maintenance_message_ru = settings.maintenance_message?.ru || ''
  },
  { immediate: true },
)

async function submitSettings() {
  await adminStore.updateSettings({
    maintenance_mode: form.maintenance_mode,
    registration_enabled: form.registration_enabled,
    ai_enabled: form.ai_enabled,
    default_role: form.default_role,
    default_deck_id: form.default_deck_id,
    default_back_id: form.default_back_id,
    maintenance_message: {
      en: form.maintenance_message_en.trim(),
      ru: form.maintenance_message_ru.trim(),
    },
  })
}

onMounted(async () => {
  const tasks = []
  if (adminStore.settingsStatus === 'idle') {
    tasks.push(adminStore.loadSettings())
  }
  if (adminStore.decksStatus === 'idle') {
    tasks.push(adminStore.loadDecks({ page: 1, limit: 100 }))
  }
  if (adminStore.cardBacksStatus === 'idle') {
    tasks.push(adminStore.loadCardBacks({ page: 1, limit: 100 }))
  }
  await Promise.all(tasks)
})
</script>

<template>
  <AdminLayout :title="t('admin.settingsTitle')" :copy="t('admin.settingsCopy')">
    <BaseSurface padding="lg" tag="section">
      <form class="settings-form" @submit.prevent="submitSettings">
        <div class="settings-grid">
          <BaseSurface tone="outlined" tag="section" class="settings-card">
            <div>
              <p class="text-eyebrow">{{ t('admin.settingsCard.systemEyebrow') }}</p>
              <h3 class="settings-card__title">{{ t('admin.settingsCard.systemTitle') }}</h3>
            </div>

            <FormField :label="t('admin.settingsCard.maintenanceMode')">
              <BaseSwitchToggle v-model="form.maintenance_mode" :label="t('admin.settingsCard.maintenanceSwitch')" />
            </FormField>

            <FormField :label="t('admin.settingsCard.registrationEnabled')">
              <BaseSwitchToggle v-model="form.registration_enabled" :label="t('admin.settingsCard.registrationSwitch')" />
            </FormField>

            <FormField :label="t('admin.settingsCard.aiEnabled')">
              <BaseSwitchToggle v-model="form.ai_enabled" :label="t('admin.settingsCard.aiSwitch')" />
            </FormField>

            <FormField :label="t('admin.settingsCard.defaultRole')" for-id="admin-default-role" required>
              <BaseSelect id="admin-default-role" v-model="form.default_role" :options="roleOptions" />
            </FormField>

            <FormField :label="t('admin.settingsCard.defaultDeck')" for-id="admin-default-deck" required>
              <BaseSelect id="admin-default-deck" v-model="form.default_deck_id" :options="deckOptions" />
            </FormField>

            <FormField :label="t('admin.settingsCard.defaultBack')" for-id="admin-default-back" required>
              <BaseSelect id="admin-default-back" v-model="form.default_back_id" :options="cardBackOptions" />
            </FormField>
          </BaseSurface>

          <BaseSurface tone="outlined" tag="section" class="settings-card">
            <div>
              <p class="text-eyebrow">{{ t('admin.settingsCard.messageEyebrow') }}</p>
              <h3 class="settings-card__title">{{ t('admin.settingsCard.messageTitle') }}</h3>
            </div>

            <FormField :label="t('admin.settingsCard.messageEn')" for-id="admin-message-en">
              <BaseInput id="admin-message-en" v-model="form.maintenance_message_en" :placeholder="t('admin.settingsCard.messagePlaceholder')" />
            </FormField>

            <FormField :label="t('admin.settingsCard.messageRu')" for-id="admin-message-ru">
              <BaseInput id="admin-message-ru" v-model="form.maintenance_message_ru" :placeholder="t('admin.settingsCard.messagePlaceholder')" />
            </FormField>

            <p class="text-body-muted">{{ t('admin.settingsCard.messageHint') }}</p>
          </BaseSurface>
        </div>

        <div v-if="adminStore.settingsError" class="admin-message admin-message--error">
          {{ adminStore.settingsError }}
        </div>

        <div v-if="adminStore.settingsMutationError" class="admin-message admin-message--error">
          {{ adminStore.settingsMutationError }}
        </div>

        <div class="settings-actions">
          <BaseButton type="submit" :loading="adminStore.settingsMutationStatus === 'loading'">
            {{ t('admin.settingsCard.save') }}
          </BaseButton>
          <span v-if="adminStore.settingsMutationStatus === 'ready'" class="text-body-muted">
            {{ t('admin.settingsCard.saved') }}
          </span>
        </div>
      </form>
    </BaseSurface>
  </AdminLayout>
</template>

<style scoped>
.settings-form,
.settings-card {
  display: grid;
  gap: 18px;
}

.settings-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20px;
}

.settings-card__title {
  font-size: var(--ds-text-heading-2);
  color: var(--ds-color-accent-primary);
}

.settings-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.admin-message {
  padding: 14px 16px;
  border: 1px solid var(--ds-color-danger);
  color: var(--ds-color-danger);
  background: rgba(178, 58, 44, 0.12);
}

@media (max-width: 860px) {
  .settings-grid {
    grid-template-columns: 1fr;
  }
}
</style>
