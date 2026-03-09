<template>
  <main class="page stack">
    <AppHeader />

    <section class="panel stack">
      <h2>{{ t("profileTitle") }}</h2>
      <p class="muted">{{ t("profileHint") }}</p>

      <BaseLoader v-if="loading" :text="t('profileLoading')" />

      <form v-else class="stack" @submit.prevent="save">
        <div class="field">
          <label for="profile-name">{{ t("profileName") }}</label>
          <input id="profile-name" v-model.trim="displayName" maxlength="64" />
        </div>

        <div class="field">
          <label for="profile-timezone">{{ t("profileTimezone") }}</label>
          <input id="profile-timezone" v-model.trim="timezone" placeholder="UTC+00:00 or Europe/Moscow" required />
        </div>

        <div class="field">
          <label for="profile-deck">{{ t("profileDeck") }}</label>
          <select id="profile-deck" v-model="preferredDeckId">
            <option value="">{{ t("profileDeckDefault") }}</option>
            <option v-for="deck in decks" :key="deck.id" :value="deck.id">{{ deck.name }}</option>
          </select>
        </div>

        <p v-if="error" class="error">{{ error }}</p>
        <p v-if="success" class="muted">{{ success }}</p>

        <button class="btn" type="submit" :disabled="saving">
          {{ saving ? t("profileSaving") : t("profileSave") }}
        </button>
      </form>
    </section>
  </main>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";

import AppHeader from "../components/AppHeader.vue";
import BaseLoader from "../components/BaseLoader.vue";
import { useI18n } from "../composables/useI18n";
import * as profileApi from "../services/profile.api";
import { useAuthStore } from "../stores/auth.store";

const authStore = useAuthStore();
const { t } = useI18n();

const loading = ref(false);
const saving = ref(false);
const error = ref("");
const success = ref("");

const decks = ref<profileApi.AvailableDeck[]>([]);
const displayName = ref("");
const timezone = ref("");
const preferredDeckId = ref("");

async function load() {
  if (!authStore.token) return;
  loading.value = true;
  error.value = "";
  success.value = "";
  try {
    const [profile, availableDecks] = await Promise.all([
      profileApi.getProfile(authStore.token),
      profileApi.getAvailableDecks(authStore.token)
    ]);

    authStore.setUser(profile);
    decks.value = availableDecks;
    displayName.value = profile.displayName ?? "";
    timezone.value = profile.timezone;

    const deckIds = new Set(availableDecks.map((item) => item.id));
    preferredDeckId.value = profile.preferredDeckId && deckIds.has(profile.preferredDeckId) ? profile.preferredDeckId : "";
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Failed to load profile";
  } finally {
    loading.value = false;
  }
}

async function save() {
  if (!authStore.token) return;
  saving.value = true;
  error.value = "";
  success.value = "";
  try {
    const updated = await profileApi.saveProfile(authStore.token, {
      displayName: displayName.value,
      timezone: timezone.value,
      preferredDeckId: preferredDeckId.value
    });
    authStore.setUser(updated);
    success.value = t("profileSaved");
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Failed to save profile";
  } finally {
    saving.value = false;
  }
}

onMounted(load);
</script>
