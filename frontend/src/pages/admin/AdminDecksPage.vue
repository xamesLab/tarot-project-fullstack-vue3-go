<template>
  <section class="panel stack">
    <div class="heading-row">
      <h2>{{ t("adminDecksTitle") }}</h2>
      <button class="btn btn-secondary" type="button" @click="reload" :disabled="loadingDecks || savingDeck">
        {{ loadingDecks ? t("adminRefreshing") : t("adminRefresh") }}
      </button>
    </div>
    <p class="muted">{{ t("adminDecksHint") }}</p>

    <p v-if="error" class="error">{{ error }}</p>
    <p v-if="success" class="muted">{{ success }}</p>

    <section class="panel stack">
      <h3>{{ t("adminDeckCreateTitle") }}</h3>
      <form class="stack" @submit.prevent="createDeck">
        <div class="field">
          <label for="deck-name">{{ t("adminDeckCreateName") }}</label>
          <input id="deck-name" v-model.trim="deckName" minlength="2" required />
        </div>
        <button class="btn" type="submit" :disabled="loadingDecks || savingDeck">
          {{ savingDeck ? t("adminSavingTheme") : t("adminDeckCreateAction") }}
        </button>
      </form>
    </section>

    <section class="panel stack">
      <h3>{{ t("adminDeckListTitle") }}</h3>

      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>{{ t("adminDeckTableName") }}</th>
              <th>{{ t("adminDeckTableStatus") }}</th>
              <th>{{ t("adminDeckTableCreatedBy") }}</th>
              <th>{{ t("adminDeckTableCards") }}</th>
              <th>{{ t("adminDeckTableAction") }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loadingDecks">
              <td colspan="5" class="muted">{{ t("commonLoading") }}</td>
            </tr>
            <tr v-else-if="decks.length === 0">
              <td colspan="5" class="muted">{{ t("adminDeckListEmpty") }}</td>
            </tr>
            <tr v-for="deck in decks" :key="deck.id">
              <td>{{ deck.name }}</td>
              <td>{{ deck.status }}</td>
              <td>{{ deck.createdByName }}</td>
              <td>{{ deck.cardsUploadedCount }} / {{ deck.requiredCardsCount }}</td>
              <td>
                <button class="btn btn-secondary" type="button" @click="selectDeck(deck.id)">
                  {{ t("adminDeckSelectAction") }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <section v-if="selectedDeck" class="panel stack">
      <h3>{{ t("adminDeckDetailsTitle") }}: {{ selectedDeck.name }}</h3>
      <p class="muted">
        {{ t("adminDeckDetailsStatus") }}: {{ selectedDeck.status }} |
        {{ selectedDeck.cardsUploadedCount }} / {{ selectedDeck.requiredCardsCount }}
      </p>

      <div class="details-grid">
        <section class="panel stack">
          <h4>{{ t("adminDeckBackTitle") }}</h4>
          <p class="muted" v-if="!selectedDeck.backImageUrl">{{ t("adminDeckBackMissing") }}</p>
          <a v-else :href="selectedDeck.backImageUrl" target="_blank" rel="noreferrer">{{ selectedDeck.backImageUrl }}</a>
          <div class="field">
            <label for="back-file">{{ t("adminDeckBackUpload") }}</label>
            <input id="back-file" type="file" accept="image/png,image/jpeg,image/webp" @change="onBackFileChange" />
          </div>
          <button class="btn" type="button" :disabled="!backFile || uploadingBack" @click="uploadBack">
            {{ uploadingBack ? t("adminDeckUploading") : t("adminDeckUploadAction") }}
          </button>
        </section>

        <section class="panel stack">
          <h4>{{ t("adminDeckCardsTitle") }}</h4>
          <p class="muted">{{ t("adminDeckCardsMissing") }}: {{ missingCardKeys.length }}</p>
          <div class="field">
            <label for="card-key">{{ t("adminDeckCardKey") }}</label>
            <select id="card-key" v-model="selectedCardKey">
              <option v-for="key in missingCardKeys" :key="key" :value="key">{{ key }}</option>
            </select>
          </div>
          <div class="field">
            <label for="card-file">{{ t("adminDeckCardUpload") }}</label>
            <input id="card-file" type="file" accept="image/png,image/jpeg,image/webp" @change="onCardFileChange" />
          </div>
          <button class="btn" type="button" :disabled="!cardFile || !selectedCardKey || uploadingCard" @click="uploadCard">
            {{ uploadingCard ? t("adminDeckUploading") : t("adminDeckUploadAction") }}
          </button>
        </section>
      </div>

      <section class="panel stack">
        <h4>{{ t("adminDeckUploadedCardsTitle") }}</h4>
        <p v-if="uploadedCards.length === 0" class="muted">{{ t("adminDeckUploadedCardsEmpty") }}</p>
        <ul v-else class="cards-list">
          <li v-for="item in uploadedCards" :key="item.key">
            <strong>{{ item.key }}</strong>:
            <a :href="item.url" target="_blank" rel="noreferrer">{{ item.url }}</a>
          </li>
        </ul>
      </section>
    </section>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

import { useI18n } from "../../composables/useI18n";
import * as adminApi from "../../services/admin.api";
import { useAuthStore } from "../../stores/auth.store";

const authStore = useAuthStore();
const { t } = useI18n();

const loadingDecks = ref(false);
const savingDeck = ref(false);
const uploadingBack = ref(false);
const uploadingCard = ref(false);
const error = ref("");
const success = ref("");

const decks = ref<adminApi.AdminDeck[]>([]);
const selectedDeck = ref<adminApi.AdminDeck | null>(null);

const deckName = ref("");
const backFile = ref<File | null>(null);
const cardFile = ref<File | null>(null);
const selectedCardKey = ref("");

const allCardKeys = buildCardKeys();

const uploadedCards = computed(() => {
  if (!selectedDeck.value?.cards) return [] as Array<{ key: string; url: string }>;
  return Object.entries(selectedDeck.value.cards)
    .map(([key, url]) => ({ key, url }))
    .sort((a, b) => Number(a.key) - Number(b.key));
});

const missingCardKeys = computed(() => {
  const used = new Set(Object.keys(selectedDeck.value?.cards ?? {}));
  return allCardKeys.filter((key) => !used.has(key));
});

async function loadDecks() {
  if (!authStore.token) return;
  loadingDecks.value = true;
  error.value = "";
  try {
    decks.value = await adminApi.getAdminDecks(authStore.token);
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Failed to load decks";
  } finally {
    loadingDecks.value = false;
  }
}

async function createDeck() {
  if (!authStore.token || !deckName.value) return;
  savingDeck.value = true;
  error.value = "";
  success.value = "";
  try {
    const created = await adminApi.createAdminDeck(authStore.token, { name: deckName.value });
    deckName.value = "";
    success.value = `${t("adminDeckCreated")}: ${created.name}`;
    await loadDecks();
    await selectDeck(created.id);
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Failed to create deck";
  } finally {
    savingDeck.value = false;
  }
}

async function selectDeck(deckId: string) {
  if (!authStore.token) return;
  try {
    const details = await adminApi.getAdminDeckById(authStore.token, deckId);
    selectedDeck.value = details;
    if (!selectedCardKey.value || !missingCardKeys.value.includes(selectedCardKey.value)) {
      selectedCardKey.value = missingCardKeys.value[0] ?? "";
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Failed to load deck details";
  }
}

async function uploadBack() {
  if (!authStore.token || !selectedDeck.value || !backFile.value) return;
  uploadingBack.value = true;
  error.value = "";
  success.value = "";
  try {
    const updated = await adminApi.uploadAdminDeckBack(authStore.token, selectedDeck.value.id, backFile.value);
    selectedDeck.value = updated;
    backFile.value = null;
    success.value = t("adminDeckBackUpdated");
    await loadDecks();
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Failed to upload back image";
  } finally {
    uploadingBack.value = false;
  }
}

async function uploadCard() {
  if (!authStore.token || !selectedDeck.value || !cardFile.value || !selectedCardKey.value) return;
  uploadingCard.value = true;
  error.value = "";
  success.value = "";
  try {
    const updated = await adminApi.uploadAdminDeckCard(
      authStore.token,
      selectedDeck.value.id,
      selectedCardKey.value,
      cardFile.value
    );
    selectedDeck.value = updated;
    cardFile.value = null;
    selectedCardKey.value = missingCardKeys.value[0] ?? "";
    success.value = t("adminDeckCardUpdated");
    await loadDecks();
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Failed to upload card image";
  } finally {
    uploadingCard.value = false;
  }
}

function onBackFileChange(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0] ?? null;
  backFile.value = file;
}

function onCardFileChange(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0] ?? null;
  cardFile.value = file;
}

async function reload() {
  await loadDecks();
  if (selectedDeck.value) {
    await selectDeck(selectedDeck.value.id);
  }
}

onMounted(loadDecks);

function buildCardKeys() {
  const keys: string[] = [];
  for (let i = 22; i <= 77; i += 1) {
    keys.push(String(i));
  }
  for (let i = 100; i <= 121; i += 1) {
    keys.push(String(i));
  }
  return keys;
}
</script>

<style scoped>
.heading-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.heading-row h2 {
  margin: 0;
}

.table-wrap {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  border: 1px solid var(--border);
  text-align: left;
  padding: 10px;
}

thead th {
  background: var(--surface-soft);
}

.details-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 10px;
}

.cards-list {
  margin: 0;
  padding-left: 18px;
  display: grid;
  gap: 6px;
}
</style>
