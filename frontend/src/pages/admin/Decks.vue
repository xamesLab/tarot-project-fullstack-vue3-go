<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import AdminLayout from '../../components/admin/AdminLayout.vue'
import CardBackCropModal from '../../components/admin/CardBackCropModal.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseInput from '../../components/base/BaseInput.vue'
import BaseSurface from '../../components/base/BaseSurface.vue'
import BaseTabs from '../../components/base/BaseTabs.vue'
import FormField from '../../components/form/FormField.vue'
import { useAdminStore } from '../../stores/admin'
import { readDeckMetaFromZip } from '../../utils/zip-meta'

const { t, locale } = useI18n()
const adminStore = useAdminStore()
const activeAssetTab = ref('catalog')

const deckUploadForm = reactive({
  nameEn: '',
  nameRu: '',
  descriptionEn: '',
  descriptionRu: '',
  author: '',
  year: '',
  tags: '',
})

const cardBackUploadForm = reactive({
  nameEn: '',
  nameRu: '',
  descriptionEn: '',
  descriptionRu: '',
  tags: '',
})

const deckUploadState = reactive({
  error: '',
  success: '',
})

const cardBackUploadState = reactive({
  error: '',
  success: '',
})

const zipFile = ref(null)
const cardBackFile = ref(null)
const pendingCardBackName = ref('')
const cropModalOpen = ref(false)
const cropSourceUrl = ref('')
const zipInputKey = ref(0)
const backInputKey = ref(0)

const isDeckFormReady = computed(() =>
  Boolean((deckUploadForm.nameEn.trim() || deckUploadForm.nameRu.trim()) && zipFile.value),
)

const isCardBackFormReady = computed(() =>
  Boolean((cardBackUploadForm.nameEn.trim() || cardBackUploadForm.nameRu.trim()) && cardBackFile.value),
)

const deckSummary = computed(() => [
  {
    key: 'done',
    label: t('admin.deckStatuses.done'),
    value: adminStore.dashboardMetrics?.decks?.by_status?.done ?? 0,
  },
  {
    key: 'processing',
    label: t('admin.deckStatuses.processing'),
    value: adminStore.dashboardMetrics?.decks?.by_status?.processing ?? 0,
  },
  {
    key: 'disabled',
    label: t('admin.deckStatuses.disabled'),
    value: adminStore.dashboardMetrics?.decks?.by_status?.disabled ?? 0,
  },
  {
    key: 'archived',
    label: t('admin.deckStatuses.archived'),
    value: adminStore.dashboardMetrics?.decks?.by_status?.archived ?? 0,
  },
])

const assetTabs = computed(() => [
  { value: 'catalog', label: t('admin.assetTabs.catalog') },
  { value: 'archive', label: t('admin.assetTabs.archive') },
])

function localizeText(value) {
  if (!value) {
    return '—'
  }

  return value[locale.value] || value.en || value.ru || '—'
}

function localizeDate(value) {
  if (!value) {
    return '—'
  }

  return new Intl.DateTimeFormat(locale.value, {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value))
}

function loadDecks(page = 1) {
  return adminStore.loadDecks({ page })
}

function loadCardBacks(page = 1) {
  return adminStore.loadCardBacks({ page })
}

function loadArchivedDecks(page = 1) {
  return adminStore.loadArchivedDecks({ page })
}

function loadArchivedCardBacks(page = 1) {
  return adminStore.loadArchivedCardBacks({ page })
}

function applyDeckMeta(meta) {
  deckUploadForm.nameEn = meta.nameEn
  deckUploadForm.nameRu = meta.nameRu
  deckUploadForm.descriptionEn = meta.descriptionEn
  deckUploadForm.descriptionRu = meta.descriptionRu
  deckUploadForm.author = meta.author
  deckUploadForm.year = meta.year
  deckUploadForm.tags = meta.tags
}

async function handleZipChange(event) {
  const [file] = event.target.files || []
  zipFile.value = file || null
  if (!file) {
    return
  }

  try {
    const meta = await readDeckMetaFromZip(file)
    if (meta) {
      applyDeckMeta(meta)
    }
  } catch (error) {
    console.warn('Failed to read 0_meta.json from deck archive', error)
  }
}

function handleCardBackChange(event) {
  const [file] = event.target.files || []
  if (!file) {
    return
  }

  if (cropSourceUrl.value) {
    URL.revokeObjectURL(cropSourceUrl.value)
  }

  pendingCardBackName.value = file.name
  cropSourceUrl.value = URL.createObjectURL(file)
  cropModalOpen.value = true
}

function closeCropModal() {
  cropModalOpen.value = false
  if (cropSourceUrl.value) {
    URL.revokeObjectURL(cropSourceUrl.value)
    cropSourceUrl.value = ''
  }
  pendingCardBackName.value = ''
  backInputKey.value += 1
}

function applyCroppedCardBack(file) {
  cardBackFile.value = file
  cropModalOpen.value = false
  if (cropSourceUrl.value) {
    URL.revokeObjectURL(cropSourceUrl.value)
    cropSourceUrl.value = ''
  }
  pendingCardBackName.value = file.name
}

function resetDeckForm() {
  deckUploadForm.nameEn = ''
  deckUploadForm.nameRu = ''
  deckUploadForm.descriptionEn = ''
  deckUploadForm.descriptionRu = ''
  deckUploadForm.author = ''
  deckUploadForm.year = ''
  deckUploadForm.tags = ''
  zipFile.value = null
  zipInputKey.value += 1
}

function resetCardBackForm() {
  cardBackUploadForm.nameEn = ''
  cardBackUploadForm.nameRu = ''
  cardBackUploadForm.descriptionEn = ''
  cardBackUploadForm.descriptionRu = ''
  cardBackUploadForm.tags = ''
  cardBackFile.value = null
  pendingCardBackName.value = ''
  backInputKey.value += 1
}

async function submitDeck() {
  deckUploadState.error = ''
  deckUploadState.success = ''

  if (!deckUploadForm.nameEn.trim() && !deckUploadForm.nameRu.trim()) {
    deckUploadState.error = t('admin.deckUpload.validation.name')
    return
  }
  if (!zipFile.value) {
    deckUploadState.error = t('admin.deckUpload.validation.zip')
    return
  }

  const formData = new FormData()
  formData.append('name_en', deckUploadForm.nameEn.trim())
  formData.append('name_ru', deckUploadForm.nameRu.trim())
  formData.append('description_en', deckUploadForm.descriptionEn.trim())
  formData.append('description_ru', deckUploadForm.descriptionRu.trim())
  formData.append('author', deckUploadForm.author.trim())
  formData.append('year', deckUploadForm.year.trim())
  formData.append('tags', deckUploadForm.tags.trim())
  formData.append('cards_zip', zipFile.value)

  const result = await adminStore.createDeck(formData)
  if (!result.ok) {
    deckUploadState.error = result.error.message
    return
  }

  deckUploadState.success = result.message || t('admin.deckUpload.success')
  resetDeckForm()
}

async function submitCardBack() {
  cardBackUploadState.error = ''
  cardBackUploadState.success = ''

  if (!cardBackUploadForm.nameEn.trim() && !cardBackUploadForm.nameRu.trim()) {
    cardBackUploadState.error = t('admin.cardBackUpload.validation.name')
    return
  }
  if (!cardBackFile.value) {
    cardBackUploadState.error = t('admin.cardBackUpload.validation.backImage')
    return
  }

  const formData = new FormData()
  formData.append('name_en', cardBackUploadForm.nameEn.trim())
  formData.append('name_ru', cardBackUploadForm.nameRu.trim())
  formData.append('description_en', cardBackUploadForm.descriptionEn.trim())
  formData.append('description_ru', cardBackUploadForm.descriptionRu.trim())
  formData.append('tags', cardBackUploadForm.tags.trim())
  formData.append('back_image', cardBackFile.value)

  const result = await adminStore.createCardBack(formData)
  if (!result.ok) {
    cardBackUploadState.error = result.error.message
    return
  }

  cardBackUploadState.success = result.message || t('admin.cardBackUpload.success')
  resetCardBackForm()
}

async function archiveDeck(deck) {
  if (!window.confirm(t('admin.archiveDeckConfirm', { name: localizeText(deck.name) }))) {
    return
  }

  await adminStore.archiveDeck(deck.id)
}

async function restoreDeck(deck) {
  await adminStore.restoreDeck(deck.id)
}

async function deleteDeck(deck) {
  if (!window.confirm(t('admin.deleteDeckConfirm', { name: localizeText(deck.name) }))) {
    return
  }

  await adminStore.deleteDeck(deck.id)
}

async function archiveCardBack(cardBack) {
  if (!window.confirm(t('admin.archiveCardBackConfirm', { name: localizeText(cardBack.name) }))) {
    return
  }

  await adminStore.archiveCardBack(cardBack.id)
}

async function restoreCardBack(cardBack) {
  await adminStore.restoreCardBack(cardBack.id)
}

async function deleteCardBack(cardBack) {
  if (!window.confirm(t('admin.deleteCardBackConfirm', { name: localizeText(cardBack.name) }))) {
    return
  }

  await adminStore.deleteCardBack(cardBack.id)
}

onMounted(async () => {
  const tasks = []

  if (adminStore.dashboardStatus === 'idle') {
    tasks.push(adminStore.loadDashboard())
  }

  if (adminStore.decksStatus === 'idle') {
    tasks.push(loadDecks())
  }

  if (adminStore.cardBacksStatus === 'idle') {
    tasks.push(loadCardBacks())
  }

  if (adminStore.archivedDecksStatus === 'idle') {
    tasks.push(loadArchivedDecks())
  }

  if (adminStore.archivedCardBacksStatus === 'idle') {
    tasks.push(loadArchivedCardBacks())
  }

  await Promise.all(tasks)
})

onBeforeUnmount(() => {
  if (cropSourceUrl.value) {
    URL.revokeObjectURL(cropSourceUrl.value)
  }
})
</script>

<template>
  <AdminLayout :title="t('admin.decksTitle')" :copy="t('admin.decksCopy')">
    <div class="page-stack">
      <section class="deck-summary">
        <BaseSurface v-for="item in deckSummary" :key="item.key" tag="article" class="summary-card">
          <p class="text-eyebrow">{{ item.label }}</p>
          <p class="summary-card__value">{{ item.value }}</p>
        </BaseSurface>
      </section>

      <BaseSurface padding="lg" tag="section">
        <div class="list-head">
          <div>
            <p class="text-eyebrow">{{ t('admin.deckUpload.title') }}</p>
            <p class="text-body-muted">{{ t('admin.deckUpload.copy') }}</p>
          </div>
          <div class="deck-table__placeholder">
            {{ t('admin.deckUpload.structureHint') }}
          </div>
        </div>

        <form class="deck-upload-grid" @submit.prevent="submitDeck">
          <FormField :label="t('admin.deckUpload.nameEn')" for-id="deck-name-en" required>
            <BaseInput id="deck-name-en" v-model="deckUploadForm.nameEn" :placeholder="t('admin.deckUpload.namePlaceholder')" />
          </FormField>

          <FormField :label="t('admin.deckUpload.nameRu')" for-id="deck-name-ru">
            <BaseInput id="deck-name-ru" v-model="deckUploadForm.nameRu" :placeholder="t('admin.deckUpload.namePlaceholder')" />
          </FormField>

          <FormField :label="t('admin.deckUpload.author')" for-id="deck-author">
            <BaseInput id="deck-author" v-model="deckUploadForm.author" :placeholder="t('admin.deckUpload.authorPlaceholder')" />
          </FormField>

          <FormField :label="t('admin.deckUpload.year')" for-id="deck-year">
            <BaseInput id="deck-year" v-model="deckUploadForm.year" type="number" placeholder="1909" />
          </FormField>

          <FormField :label="t('admin.deckUpload.descriptionEn')" for-id="deck-description-en">
            <BaseInput
              id="deck-description-en"
              v-model="deckUploadForm.descriptionEn"
              :placeholder="t('admin.deckUpload.descriptionPlaceholder')"
            />
          </FormField>

          <FormField :label="t('admin.deckUpload.descriptionRu')" for-id="deck-description-ru">
            <BaseInput
              id="deck-description-ru"
              v-model="deckUploadForm.descriptionRu"
              :placeholder="t('admin.deckUpload.descriptionPlaceholder')"
            />
          </FormField>

          <FormField :label="t('admin.deckUpload.tags')" for-id="deck-tags">
            <BaseInput id="deck-tags" v-model="deckUploadForm.tags" :placeholder="t('admin.deckUpload.tagsPlaceholder')" />
          </FormField>

          <FormField :label="t('admin.deckUpload.zip')" :hint="t('admin.deckUpload.zipHint')" required>
            <label class="file-field">
              <input :key="zipInputKey" type="file" accept=".zip,application/zip" @change="handleZipChange" />
              <span class="file-field__content">
                <span>{{ zipFile?.name || t('admin.deckUpload.chooseZip') }}</span>
                <span v-if="zipFile" class="file-field__check" aria-hidden="true">✓</span>
              </span>
            </label>
          </FormField>

          <div class="deck-upload-actions">
            <BaseButton type="submit" :loading="adminStore.deckCreateStatus === 'loading'" :disabled="!isDeckFormReady">
              {{ t('admin.deckUpload.submit') }}
            </BaseButton>
          </div>
        </form>

        <p v-if="deckUploadState.error" class="admin-message admin-message--error">{{ deckUploadState.error }}</p>
        <p v-else-if="deckUploadState.success" class="admin-message admin-message--success">{{ deckUploadState.success }}</p>
      </BaseSurface>

      <BaseSurface padding="lg" tag="section">
        <div class="list-head">
          <div>
            <p class="text-eyebrow">{{ t('admin.cardBackUpload.title') }}</p>
            <p class="text-body-muted">{{ t('admin.cardBackUpload.copy') }}</p>
          </div>
        </div>

        <form class="deck-upload-grid" @submit.prevent="submitCardBack">
          <FormField :label="t('admin.cardBackUpload.nameEn')" for-id="card-back-name-en" required>
            <BaseInput id="card-back-name-en" v-model="cardBackUploadForm.nameEn" :placeholder="t('admin.cardBackUpload.namePlaceholder')" />
          </FormField>

          <FormField :label="t('admin.cardBackUpload.nameRu')" for-id="card-back-name-ru">
            <BaseInput id="card-back-name-ru" v-model="cardBackUploadForm.nameRu" :placeholder="t('admin.cardBackUpload.namePlaceholder')" />
          </FormField>

          <FormField :label="t('admin.cardBackUpload.descriptionEn')" for-id="card-back-description-en">
            <BaseInput
              id="card-back-description-en"
              v-model="cardBackUploadForm.descriptionEn"
              :placeholder="t('admin.cardBackUpload.descriptionPlaceholder')"
            />
          </FormField>

          <FormField :label="t('admin.cardBackUpload.descriptionRu')" for-id="card-back-description-ru">
            <BaseInput
              id="card-back-description-ru"
              v-model="cardBackUploadForm.descriptionRu"
              :placeholder="t('admin.cardBackUpload.descriptionPlaceholder')"
            />
          </FormField>

          <FormField :label="t('admin.cardBackUpload.tags')" for-id="card-back-tags">
            <BaseInput id="card-back-tags" v-model="cardBackUploadForm.tags" :placeholder="t('admin.cardBackUpload.tagsPlaceholder')" />
          </FormField>

          <FormField :label="t('admin.cardBackUpload.backImage')" :hint="t('admin.cardBackUpload.backImageHint')" required>
            <label class="file-field">
              <input
                :key="backInputKey"
                type="file"
                accept=".png,.jpg,.jpeg,image/png,image/jpeg"
                @change="handleCardBackChange"
              />
              <span class="file-field__content">
                <span>{{ cardBackFile?.name || t('admin.cardBackUpload.chooseBackImage') }}</span>
                <span v-if="cardBackFile" class="file-field__check" aria-hidden="true">✓</span>
              </span>
            </label>
          </FormField>

          <div class="deck-upload-actions">
            <BaseButton type="submit" :loading="adminStore.cardBackCreateStatus === 'loading'" :disabled="!isCardBackFormReady">
              {{ t('admin.cardBackUpload.submit') }}
            </BaseButton>
          </div>
        </form>

        <p v-if="cardBackUploadState.error" class="admin-message admin-message--error">{{ cardBackUploadState.error }}</p>
        <p v-else-if="cardBackUploadState.success" class="admin-message admin-message--success">{{ cardBackUploadState.success }}</p>
      </BaseSurface>

      <BaseSurface padding="lg" tag="section">
        <BaseTabs v-model="activeAssetTab" :items="assetTabs">
          <template #panel="{ activeItem }">
            <div v-if="activeItem?.value === 'catalog'" class="page-stack">
              <section>
                <div class="list-head">
                  <div>
                    <p class="text-eyebrow">{{ t('admin.deckTable.title') }}</p>
                    <p class="text-body-muted">{{ t('admin.deckTable.copy') }}</p>
                  </div>
                </div>

                <div v-if="adminStore.decksError" class="admin-message admin-message--error">
                  {{ adminStore.decksError }}
                </div>

                <div v-if="!adminStore.deckList.length && adminStore.decksStatus === 'ready'" class="empty-state">
                  {{ t('admin.emptyDecks') }}
                </div>

                <div v-else class="deck-list">
                  <article v-for="deck in adminStore.deckList" :key="deck.id" class="deck-card">
                    <div class="deck-card__head">
                      <div>
                        <h3 class="deck-card__title">{{ localizeText(deck.name) }}</h3>
                        <p class="text-body-muted">{{ localizeText(deck.description) }}</p>
                      </div>
                      <span class="deck-card__badge">{{ t(`admin.deckStatuses.${deck.status}`) }}</span>
                    </div>

                    <dl class="deck-meta">
                      <div>
                        <dt>{{ t('admin.deckTable.author') }}</dt>
                        <dd>{{ deck.author || '—' }}</dd>
                      </div>
                      <div>
                        <dt>{{ t('admin.deckTable.year') }}</dt>
                        <dd>{{ deck.year || '—' }}</dd>
                      </div>
                      <div>
                        <dt>{{ t('admin.deckTable.cardCount') }}</dt>
                        <dd>{{ deck.card_count || 0 }}</dd>
                      </div>
                      <div>
                        <dt>{{ t('admin.deckTable.updatedAt') }}</dt>
                        <dd>{{ localizeDate(deck.updated_at) }}</dd>
                      </div>
                    </dl>

                    <div class="card-actions">
                      <BaseButton
                        variant="text"
                        :loading="adminStore.deckMutationStatusById[deck.id] === 'loading'"
                        @click="archiveDeck(deck)"
                      >
                        {{ t('admin.actions.archive') }}
                      </BaseButton>
                    </div>
                    <p v-if="adminStore.deckMutationErrorById[deck.id]" class="admin-inline-error">
                      {{ adminStore.deckMutationErrorById[deck.id] }}
                    </p>
                  </article>
                </div>

                <div class="pagination-bar">
                  <BaseButton
                    variant="text"
                    :disabled="adminStore.deckPagination.page <= 1 || adminStore.decksStatus === 'loading'"
                    @click="loadDecks(adminStore.deckPagination.page - 1)"
                  >
                    {{ t('admin.pagination.previous') }}
                  </BaseButton>
                  <span class="text-body-muted">
                    {{ t('admin.pagination.page', { page: adminStore.deckPagination.page, total: adminStore.deckPagination.total_pages }) }}
                  </span>
                  <BaseButton
                    variant="text"
                    :disabled="adminStore.deckPagination.page >= adminStore.deckPagination.total_pages || adminStore.decksStatus === 'loading'"
                    @click="loadDecks(adminStore.deckPagination.page + 1)"
                  >
                    {{ t('admin.pagination.next') }}
                  </BaseButton>
                </div>
              </section>

              <section>
                <div class="list-head">
                  <div>
                    <p class="text-eyebrow">{{ t('admin.cardBackTable.title') }}</p>
                    <p class="text-body-muted">{{ t('admin.cardBackTable.copy') }}</p>
                  </div>
                </div>

                <div v-if="adminStore.cardBacksError" class="admin-message admin-message--error">
                  {{ adminStore.cardBacksError }}
                </div>

                <div v-if="!adminStore.cardBackList.length && adminStore.cardBacksStatus === 'ready'" class="empty-state">
                  {{ t('admin.emptyCardBacks') }}
                </div>

                <div v-else class="deck-list">
                  <article v-for="cardBack in adminStore.cardBackList" :key="cardBack.id" class="deck-card">
                    <div class="deck-card__head">
                      <div>
                        <h3 class="deck-card__title">{{ localizeText(cardBack.name) }}</h3>
                        <p class="text-body-muted">{{ localizeText(cardBack.description) }}</p>
                      </div>
                      <span class="deck-card__badge">{{ t(`admin.deckStatuses.${cardBack.status}`) }}</span>
                    </div>

                    <dl class="deck-meta">
                      <div>
                        <dt>{{ t('admin.cardBackTable.imagePath') }}</dt>
                        <dd>{{ cardBack.image_path || '—' }}</dd>
                      </div>
                      <div>
                        <dt>{{ t('admin.cardBackTable.updatedAt') }}</dt>
                        <dd>{{ localizeDate(cardBack.updated_at) }}</dd>
                      </div>
                    </dl>

                    <div class="card-actions">
                      <BaseButton
                        variant="text"
                        :loading="adminStore.cardBackMutationStatusById[cardBack.id] === 'loading'"
                        @click="archiveCardBack(cardBack)"
                      >
                        {{ t('admin.actions.archive') }}
                      </BaseButton>
                    </div>
                    <p v-if="adminStore.cardBackMutationErrorById[cardBack.id]" class="admin-inline-error">
                      {{ adminStore.cardBackMutationErrorById[cardBack.id] }}
                    </p>
                  </article>
                </div>

                <div class="pagination-bar">
                  <BaseButton
                    variant="text"
                    :disabled="adminStore.cardBackPagination.page <= 1 || adminStore.cardBacksStatus === 'loading'"
                    @click="loadCardBacks(adminStore.cardBackPagination.page - 1)"
                  >
                    {{ t('admin.pagination.previous') }}
                  </BaseButton>
                  <span class="text-body-muted">
                    {{ t('admin.pagination.page', { page: adminStore.cardBackPagination.page, total: adminStore.cardBackPagination.total_pages }) }}
                  </span>
                  <BaseButton
                    variant="text"
                    :disabled="adminStore.cardBackPagination.page >= adminStore.cardBackPagination.total_pages || adminStore.cardBacksStatus === 'loading'"
                    @click="loadCardBacks(adminStore.cardBackPagination.page + 1)"
                  >
                    {{ t('admin.pagination.next') }}
                  </BaseButton>
                </div>
              </section>
            </div>

            <div v-else class="page-stack">
              <section>
                <div class="list-head">
                  <div>
                    <p class="text-eyebrow">{{ t('admin.deckTable.title') }}</p>
                    <p class="text-body-muted">{{ t('admin.deckTable.archiveCopy') }}</p>
                  </div>
                </div>

                <div v-if="adminStore.archivedDecksError" class="admin-message admin-message--error">
                  {{ adminStore.archivedDecksError }}
                </div>

                <div v-if="!adminStore.archivedDeckList.length && adminStore.archivedDecksStatus === 'ready'" class="empty-state">
                  {{ t('admin.emptyArchivedDecks') }}
                </div>

                <div v-else class="deck-list">
                  <article v-for="deck in adminStore.archivedDeckList" :key="deck.id" class="deck-card">
                    <div class="deck-card__head">
                      <div>
                        <h3 class="deck-card__title">{{ localizeText(deck.name) }}</h3>
                        <p class="text-body-muted">{{ localizeText(deck.description) }}</p>
                      </div>
                      <span class="deck-card__badge">{{ t(`admin.deckStatuses.${deck.status}`) }}</span>
                    </div>

                    <dl class="deck-meta">
                      <div>
                        <dt>{{ t('admin.deckTable.author') }}</dt>
                        <dd>{{ deck.author || '—' }}</dd>
                      </div>
                      <div>
                        <dt>{{ t('admin.deckTable.year') }}</dt>
                        <dd>{{ deck.year || '—' }}</dd>
                      </div>
                      <div>
                        <dt>{{ t('admin.deckTable.cardCount') }}</dt>
                        <dd>{{ deck.card_count || 0 }}</dd>
                      </div>
                      <div>
                        <dt>{{ t('admin.deckTable.updatedAt') }}</dt>
                        <dd>{{ localizeDate(deck.updated_at) }}</dd>
                      </div>
                    </dl>

                    <div class="card-actions">
                      <BaseButton
                        variant="text"
                        :loading="adminStore.deckMutationStatusById[deck.id] === 'loading'"
                        @click="restoreDeck(deck)"
                      >
                        {{ t('admin.actions.restore') }}
                      </BaseButton>
                      <BaseButton
                        variant="text"
                        :loading="adminStore.deckMutationStatusById[deck.id] === 'loading'"
                        @click="deleteDeck(deck)"
                      >
                        {{ t('admin.actions.deletePermanent') }}
                      </BaseButton>
                    </div>
                    <p v-if="adminStore.deckMutationErrorById[deck.id]" class="admin-inline-error">
                      {{ adminStore.deckMutationErrorById[deck.id] }}
                    </p>
                  </article>
                </div>

                <div class="pagination-bar">
                  <BaseButton
                    variant="text"
                    :disabled="adminStore.archivedDeckPagination.page <= 1 || adminStore.archivedDecksStatus === 'loading'"
                    @click="loadArchivedDecks(adminStore.archivedDeckPagination.page - 1)"
                  >
                    {{ t('admin.pagination.previous') }}
                  </BaseButton>
                  <span class="text-body-muted">
                    {{ t('admin.pagination.page', { page: adminStore.archivedDeckPagination.page, total: adminStore.archivedDeckPagination.total_pages }) }}
                  </span>
                  <BaseButton
                    variant="text"
                    :disabled="adminStore.archivedDeckPagination.page >= adminStore.archivedDeckPagination.total_pages || adminStore.archivedDecksStatus === 'loading'"
                    @click="loadArchivedDecks(adminStore.archivedDeckPagination.page + 1)"
                  >
                    {{ t('admin.pagination.next') }}
                  </BaseButton>
                </div>
              </section>

              <section>
                <div class="list-head">
                  <div>
                    <p class="text-eyebrow">{{ t('admin.cardBackTable.title') }}</p>
                    <p class="text-body-muted">{{ t('admin.cardBackTable.archiveCopy') }}</p>
                  </div>
                </div>

                <div v-if="adminStore.archivedCardBacksError" class="admin-message admin-message--error">
                  {{ adminStore.archivedCardBacksError }}
                </div>

                <div v-if="!adminStore.archivedCardBackList.length && adminStore.archivedCardBacksStatus === 'ready'" class="empty-state">
                  {{ t('admin.emptyArchivedCardBacks') }}
                </div>

                <div v-else class="deck-list">
                  <article v-for="cardBack in adminStore.archivedCardBackList" :key="cardBack.id" class="deck-card">
                    <div class="deck-card__head">
                      <div>
                        <h3 class="deck-card__title">{{ localizeText(cardBack.name) }}</h3>
                        <p class="text-body-muted">{{ localizeText(cardBack.description) }}</p>
                      </div>
                      <span class="deck-card__badge">{{ t(`admin.deckStatuses.${cardBack.status}`) }}</span>
                    </div>

                    <dl class="deck-meta">
                      <div>
                        <dt>{{ t('admin.cardBackTable.imagePath') }}</dt>
                        <dd>{{ cardBack.image_path || '—' }}</dd>
                      </div>
                      <div>
                        <dt>{{ t('admin.cardBackTable.updatedAt') }}</dt>
                        <dd>{{ localizeDate(cardBack.updated_at) }}</dd>
                      </div>
                    </dl>

                    <div class="card-actions">
                      <BaseButton
                        variant="text"
                        :loading="adminStore.cardBackMutationStatusById[cardBack.id] === 'loading'"
                        @click="restoreCardBack(cardBack)"
                      >
                        {{ t('admin.actions.restore') }}
                      </BaseButton>
                      <BaseButton
                        variant="text"
                        :loading="adminStore.cardBackMutationStatusById[cardBack.id] === 'loading'"
                        @click="deleteCardBack(cardBack)"
                      >
                        {{ t('admin.actions.deletePermanent') }}
                      </BaseButton>
                    </div>
                    <p v-if="adminStore.cardBackMutationErrorById[cardBack.id]" class="admin-inline-error">
                      {{ adminStore.cardBackMutationErrorById[cardBack.id] }}
                    </p>
                  </article>
                </div>

                <div class="pagination-bar">
                  <BaseButton
                    variant="text"
                    :disabled="adminStore.archivedCardBackPagination.page <= 1 || adminStore.archivedCardBacksStatus === 'loading'"
                    @click="loadArchivedCardBacks(adminStore.archivedCardBackPagination.page - 1)"
                  >
                    {{ t('admin.pagination.previous') }}
                  </BaseButton>
                  <span class="text-body-muted">
                    {{ t('admin.pagination.page', { page: adminStore.archivedCardBackPagination.page, total: adminStore.archivedCardBackPagination.total_pages }) }}
                  </span>
                  <BaseButton
                    variant="text"
                    :disabled="adminStore.archivedCardBackPagination.page >= adminStore.archivedCardBackPagination.total_pages || adminStore.archivedCardBacksStatus === 'loading'"
                    @click="loadArchivedCardBacks(adminStore.archivedCardBackPagination.page + 1)"
                  >
                    {{ t('admin.pagination.next') }}
                  </BaseButton>
                </div>
              </section>
            </div>
          </template>
        </BaseTabs>
      </BaseSurface>
    </div>

    <CardBackCropModal
      :is-open="cropModalOpen"
      :image-url="cropSourceUrl"
      :file-name="pendingCardBackName"
      :title="t('admin.cardBackUpload.cropTitle')"
      :copy="t('admin.cardBackUpload.cropCopy')"
      :zoom-label="t('admin.cardBackUpload.zoomLabel')"
      :cancel-label="t('admin.cardBackUpload.cancelCrop')"
      :confirm-label="t('admin.cardBackUpload.applyCrop')"
      @close="closeCropModal"
      @confirm="applyCroppedCardBack"
    />
  </AdminLayout>
</template>

<style scoped>
.deck-summary {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

.summary-card {
  display: grid;
  gap: 10px;
}

.summary-card__value {
  margin: 0;
  font-size: var(--ds-text-display-sm);
  color: var(--ds-color-accent-primary);
}

.list-head {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  align-items: flex-start;
  margin-bottom: 20px;
}

.deck-upload-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px;
}

.deck-upload-actions {
  grid-column: 1 / -1;
}

.file-field {
  display: grid;
  gap: 10px;
  padding: 16px 18px;
  border: 1px dashed var(--ds-color-border-strong);
  background: var(--ds-color-surface-secondary);
  color: var(--ds-color-text-primary);
  cursor: pointer;
}

.file-field input {
  display: none;
}

.file-field__content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.file-field__check {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border: 1px solid var(--ds-color-accent-primary);
  color: var(--ds-color-accent-primary);
  font-size: 14px;
  line-height: 1;
  flex-shrink: 0;
}

.admin-message {
  margin-top: 16px;
  padding: 14px 16px;
  border: 1px solid var(--ds-color-danger);
  color: var(--ds-color-danger);
  background: rgba(178, 58, 44, 0.12);
}

.admin-message--success {
  border-color: var(--ds-color-accent-primary);
  color: var(--ds-color-accent-primary);
  background: rgba(190, 155, 91, 0.12);
}

.deck-list {
  display: grid;
  gap: 16px;
}

.deck-card {
  display: grid;
  gap: 18px;
  padding: 18px;
  border: 1px solid var(--ds-color-border-subtle);
  background: var(--ds-color-surface-secondary);
}

.card-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.admin-inline-error {
  margin: 0;
  color: var(--ds-color-danger);
  font-size: var(--ds-text-body-sm);
}

.deck-card__head {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: flex-start;
}

.deck-card__title {
  margin: 0 0 8px;
  font-size: var(--ds-text-heading-3);
  color: var(--ds-color-accent-primary);
}

.deck-card__badge {
  padding: 6px 10px;
  border: 1px solid var(--ds-color-border-strong);
  color: var(--ds-color-text-secondary);
  font-size: var(--ds-text-label-sm);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.deck-meta {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  margin: 0;
}

.deck-meta dt {
  color: var(--ds-color-text-secondary);
  font-size: var(--ds-text-label-sm);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.deck-meta dd {
  margin: 6px 0 0;
  color: var(--ds-color-text-primary);
}

.pagination-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-top: 20px;
}

.empty-state {
  padding: 24px 0;
  color: var(--ds-color-text-secondary);
}

@media (max-width: 980px) {
  .deck-upload-grid,
  .deck-summary,
  .deck-meta {
    grid-template-columns: 1fr;
  }

  .deck-card__head,
  .list-head,
  .pagination-bar {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
