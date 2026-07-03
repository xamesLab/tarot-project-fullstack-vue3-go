<script setup>
import { onMounted, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import AdminLayout from '../../components/admin/AdminLayout.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseInput from '../../components/base/BaseInput.vue'
import BaseSurface from '../../components/base/BaseSurface.vue'
import BaseSwitchToggle from '../../components/base/BaseSwitchToggle.vue'
import FormField from '../../components/form/FormField.vue'
import { useAdminStore } from '../../stores/admin'

const { t } = useI18n()
const adminStore = useAdminStore()

const drafts = reactive({})

watch(
  () => adminStore.policies,
  (policies) => {
    policies.forEach((policy) => {
      drafts[policy.role] = {
        three_card_limit: String(policy.three_card_limit ?? 0),
        ai_enabled: policy.ai_enabled !== false,
        ai_daily_limit: String(policy.ai_daily_limit ?? 0),
        history_retention_days: String(policy.history_retention_days ?? 0),
      }
    })
  },
  { immediate: true },
)

function isSaving(role) {
  return adminStore.policyMutationStatusByRole[role] === 'loading'
}

async function savePolicy(role) {
  const draft = drafts[role]
  if (!draft) {
    return
  }

  await adminStore.updateRolePolicy(role, {
    three_card_limit: Number(draft.three_card_limit),
    ai_enabled: Boolean(draft.ai_enabled),
    ai_daily_limit: Number(draft.ai_daily_limit),
    history_retention_days: Number(draft.history_retention_days),
  })
}

onMounted(() => {
  if (adminStore.policiesStatus === 'idle') {
    adminStore.loadPolicies()
  }
})
</script>

<template>
  <AdminLayout :title="t('admin.policiesTitle')" :copy="t('admin.policiesCopy')">
    <div class="page-stack">
      <div v-if="adminStore.policiesError" class="admin-message admin-message--error">
        {{ adminStore.policiesError }}
      </div>

      <div class="policy-grid">
        <BaseSurface
          v-for="policy in adminStore.policies"
          :key="policy.role"
          padding="lg"
          tag="section"
          class="policy-card"
        >
          <div class="policy-card__head">
            <div>
              <p class="text-eyebrow">{{ t('admin.policyCard.eyebrow') }}</p>
              <h3 class="policy-card__title">{{ t(`roles.${policy.role}`) }}</h3>
            </div>
            <span class="policy-card__badge">{{ policy.role }}</span>
          </div>

          <form v-if="drafts[policy.role]" class="policy-form" @submit.prevent="savePolicy(policy.role)">
            <FormField :label="t('admin.policyCard.threeCardLimit')" :for-id="`policy-limit-${policy.role}`" required>
              <BaseInput
                :id="`policy-limit-${policy.role}`"
                v-model="drafts[policy.role].three_card_limit"
                type="number"
              />
            </FormField>

            <FormField :label="t('admin.policyCard.aiEnabled')">
              <BaseSwitchToggle
                v-model="drafts[policy.role].ai_enabled"
                :label="t('admin.policyCard.aiEnabledSwitch')"
              />
            </FormField>

            <FormField :label="t('admin.policyCard.aiDailyLimit')" :for-id="`policy-ai-limit-${policy.role}`" required>
              <BaseInput
                :id="`policy-ai-limit-${policy.role}`"
                v-model="drafts[policy.role].ai_daily_limit"
                type="number"
              />
            </FormField>

            <FormField :label="t('admin.policyCard.historyRetention')" :for-id="`policy-retention-${policy.role}`" required>
              <BaseInput
                :id="`policy-retention-${policy.role}`"
                v-model="drafts[policy.role].history_retention_days"
                type="number"
              />
            </FormField>

            <p class="text-body-muted">{{ t('admin.policyCard.copy') }}</p>

            <p v-if="adminStore.policyMutationErrorByRole[policy.role]" class="admin-message admin-message--inline">
              {{ adminStore.policyMutationErrorByRole[policy.role] }}
            </p>

            <BaseButton type="submit" :loading="isSaving(policy.role)">
              {{ t('admin.policyCard.save') }}
            </BaseButton>
          </form>
        </BaseSurface>
      </div>
    </div>
  </AdminLayout>
</template>

<style scoped>
.policy-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20px;
}

.policy-card,
.policy-form {
  display: grid;
  gap: 18px;
}

.policy-card__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.policy-card__title {
  font-size: var(--ds-text-heading-2);
  color: var(--ds-color-accent-primary);
}

.policy-card__badge {
  padding: 6px 12px;
  border: 1px solid var(--ds-color-border-accent);
  color: var(--ds-color-accent-primary);
  font-family: var(--ds-font-ui);
  font-size: var(--ds-text-caption);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.admin-message {
  padding: 14px 16px;
  border: 1px solid var(--ds-color-danger);
  color: var(--ds-color-danger);
  background: rgba(178, 58, 44, 0.12);
}

.admin-message--inline {
  padding: 0;
  border: 0;
  background: transparent;
}

@media (max-width: 860px) {
  .policy-grid {
    grid-template-columns: 1fr;
  }
}
</style>
