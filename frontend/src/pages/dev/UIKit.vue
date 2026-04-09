<script setup>
import { reactive, ref } from 'vue'

import BaseBadge from '../../components/base/BaseBadge.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseInput from '../../components/base/BaseInput.vue'
import BaseSegmentedToggle from '../../components/base/BaseSegmentedToggle.vue'
import BaseSelect from '../../components/base/BaseSelect.vue'
import BaseSurface from '../../components/base/BaseSurface.vue'
import BaseSwitchToggle from '../../components/base/BaseSwitchToggle.vue'
import BaseTabs from '../../components/base/BaseTabs.vue'
import CardOfDay from '../../components/domain/CardOfDay.vue'
import ReadingHistory from '../../components/domain/ReadingHistory.vue'
import TarotCard from '../../components/domain/TarotCard.vue'
import ThreeCardSpread from '../../components/domain/ThreeCardSpread.vue'
import UserChip from '../../components/domain/UserChip.vue'
import EmptyStateCard from '../../components/feedback/EmptyStateCard.vue'
import LoadingCard from '../../components/feedback/LoadingCard.vue'
import AuthCard from '../../components/form/AuthCard.vue'
import FormField from '../../components/form/FormField.vue'
import PageContainer from '../../components/layout/PageContainer.vue'
import SectionHeader from '../../components/layout/SectionHeader.vue'

const activeTab = ref('overview')
const language = ref('en')
const theme = ref('dark')
const notifications = ref(true)
const cinematicAnimations = ref(false)

const formState = reactive({
  email: 'alice@nebula.tarot',
  deck: 'classic',
})

const tabItems = [
  {
    value: 'overview',
    label: 'Overview',
    content: 'The design system is built around the same production components used across the application shell.',
  },
  {
    value: 'compositions',
    label: 'Compositions',
    content: 'Screens are expected to compose primitives instead of inventing bespoke page-level styling patterns.',
  },
  {
    value: 'motion',
    label: 'Motion',
    content: 'Animations stay soft and meaningful: hover, focus, tab change, modal reveal, and card interactions.',
  },
]

const deckOptions = [
  { value: 'classic', label: 'Classic Waite' },
  { value: 'gold', label: 'Golden Etching' },
  { value: 'moon', label: 'Moon Archive' },
]

const languageOptions = [
  { value: 'en', label: 'EN' },
  { value: 'ru', label: 'RU' },
]

const themeOptions = [
  { value: 'light', label: 'Light' },
  { value: 'dark', label: 'Dark' },
]

const sampleUser = {
  username: 'Alice Moon',
  email: 'alice@nebula.tarot',
  role: 'adept',
}

function createCard({
  key,
  position = 0,
  name,
  nameRu,
  interpretation,
  interpretationRu,
  tags,
  tagsRu,
  viewed = true,
}) {
  return {
    key,
    position,
    name,
    name_ru: nameRu,
    interpretation,
    interpretation_ru: interpretationRu,
    tags,
    tags_ru: tagsRu,
    viewed,
    image_url: '',
    back_image_url: '',
  }
}

const showcaseCard = createCard({
  key: 'the-star',
  name: 'The Star',
  nameRu: 'Звезда',
  interpretation: 'Hope, calm clarity, and the feeling that your path is finally opening up.',
  interpretationRu: 'Надежда, ясность и ощущение, что путь наконец начинает раскрываться.',
  tags: 'Hope, Healing, Direction',
  tagsRu: 'Надежда, Исцеление, Направление',
})

const cardOfDayReading = {
  id: 'daily-reference',
  accounting_date: '2026-04-09',
  next_available_at: '2026-04-10T03:00:00Z',
  cards: [
    createCard({
      key: 'the-hermit',
      name: 'The Hermit',
      nameRu: 'Отшельник',
      interpretation: 'Step back, observe, and trust the quiet wisdom that appears away from the noise.',
      interpretationRu: 'Сделайте шаг назад и доверьтесь тихой мудрости, которая проявляется вдали от шума.',
      tags: 'Reflection, Solitude, Wisdom',
      tagsRu: 'Рефлексия, Одиночество, Мудрость',
      viewed: true,
    }),
  ],
}

const spreadReading = {
  id: 'spread-reference',
  remaining_today: {
    count: 2,
    limit: 3,
  },
  cards: [
    createCard({
      key: 'six-of-swords',
      position: 0,
      name: 'Six of Swords',
      nameRu: 'Шестерка Мечей',
      interpretation: 'A transition away from turbulence toward a steadier emotional climate.',
      interpretationRu: 'Переход от турбулентности к более спокойному эмоциональному состоянию.',
      tags: 'Transition, Distance, Recovery',
      tagsRu: 'Переход, Дистанция, Восстановление',
      viewed: true,
    }),
    createCard({
      key: 'strength',
      position: 1,
      name: 'Strength',
      nameRu: 'Сила',
      interpretation: 'Lead with patience and inner composure instead of force.',
      interpretationRu: 'Действуйте через терпение и внутреннюю собранность, а не через давление.',
      tags: 'Patience, Courage, Composure',
      tagsRu: 'Терпение, Мужество, Самообладание',
      viewed: true,
    }),
    createCard({
      key: 'ace-of-pentacles',
      position: 2,
      name: 'Ace of Pentacles',
      nameRu: 'Туз Пентаклей',
      interpretation: 'A grounded opportunity is approaching; build carefully and it will last.',
      interpretationRu: 'Приближается практичная возможность; стройте аккуратно, и результат сохранится.',
      tags: 'Opportunity, Grounding, Growth',
      tagsRu: 'Возможность, Опора, Рост',
      viewed: true,
    }),
  ],
}

const historyReadings = [
  {
    id: 'history-1',
    type: 'card_of_day',
    generated_at: '2026-04-08T05:40:00Z',
    accounting_date: '2026-04-08',
    cards: [
      createCard({
        key: 'temperance',
        name: 'Temperance',
        nameRu: 'Умеренность',
        interpretation: 'Blend opposites into something stable and sustainable.',
        interpretationRu: 'Соединяйте противоположности в нечто устойчивое и жизнеспособное.',
        tags: 'Balance, Harmony, Patience',
        tagsRu: 'Баланс, Гармония, Терпение',
      }),
    ],
  },
  {
    id: 'history-2',
    type: 'three_card',
    generated_at: '2026-04-07T11:10:00Z',
    accounting_date: '2026-04-07',
    cards: spreadReading.cards,
  },
]
</script>

<template>
  <div class="ui-kit-page">
    <div class="app-background" aria-hidden="true"></div>

    <PageContainer wide tag="main">
      <div class="page-stack">
        <BaseSurface class="ui-kit-hero" tone="outlined" padding="lg" tag="section">
          <SectionHeader
            size="lg"
            title-tag="h1"
            eyebrow="Nebula Tarot"
            title="UI Kit Reference"
            copy="This page now renders the same primitives and composite blocks used in production, so design changes stay synchronized instead of drifting into a parallel demo."
          >
            <template #meta>
              <BaseBadge label="Source of truth" />
              <BaseBadge variant="subtle" label="v0.1.0" />
            </template>
          </SectionHeader>

          <div class="ui-kit-hero__tokens">
            <BaseBadge label="Dark-first surfaces" />
            <BaseBadge variant="subtle" label="Gold accent" />
            <BaseBadge variant="subtle" label="Display + serif + UI fonts" />
            <BaseBadge variant="subtle" label="Rigid geometry" />
          </div>
        </BaseSurface>

        <div class="ui-kit-grid ui-kit-grid--two">
          <BaseSurface padding="lg" tag="section">
            <SectionHeader
              size="sm"
              eyebrow="Foundations"
              title="Color and semantic badges"
              copy="Swatches and badges pull directly from the active token layer."
            />

            <div class="ui-kit-swatch-grid">
              <div class="ui-kit-swatch">
                <span class="ui-kit-swatch__chip ui-kit-swatch__chip--page"></span>
                <div>
                  <p class="ui-kit-swatch__label">Page background</p>
                  <p class="text-body-muted"><code>--ds-color-bg-page</code></p>
                </div>
              </div>
              <div class="ui-kit-swatch">
                <span class="ui-kit-swatch__chip ui-kit-swatch__chip--surface"></span>
                <div>
                  <p class="ui-kit-swatch__label">Surface</p>
                  <p class="text-body-muted"><code>--ds-color-bg-surface</code></p>
                </div>
              </div>
              <div class="ui-kit-swatch">
                <span class="ui-kit-swatch__chip ui-kit-swatch__chip--accent"></span>
                <div>
                  <p class="ui-kit-swatch__label">Accent</p>
                  <p class="text-body-muted"><code>--ds-color-accent-primary</code></p>
                </div>
              </div>
            </div>

            <div class="cluster">
              <BaseBadge label="Accent" />
              <BaseBadge variant="subtle" label="Subtle" />
              <BaseBadge variant="danger" label="Danger" />
              <BaseBadge variant="success" label="Success" />
            </div>
          </BaseSurface>

          <BaseSurface padding="lg" tag="section">
            <SectionHeader
              size="sm"
              eyebrow="Navigation"
              title="Header-level controls"
              copy="Language, theme and user identity all use the same component set as the production header."
            />

            <div class="stack-md">
              <div class="ui-kit-inline">
                <BaseSegmentedToggle v-model="language" :options="languageOptions" aria-label="Language" />
                <BaseSegmentedToggle v-model="theme" :options="themeOptions" aria-label="Theme" />
              </div>

              <UserChip :user="sampleUser" />
            </div>
          </BaseSurface>
        </div>

        <div class="ui-kit-grid ui-kit-grid--two">
          <BaseSurface padding="lg" tag="section">
            <SectionHeader
              size="sm"
              eyebrow="Controls"
              title="Buttons, fields and toggles"
              copy="Buttons, inputs and switches should be consumed directly instead of being restyled per-page."
            />

            <div class="stack-lg">
              <div class="cluster">
                <BaseButton>Primary</BaseButton>
                <BaseButton variant="secondary">Secondary</BaseButton>
                <BaseButton variant="ghost">Ghost</BaseButton>
                <BaseButton variant="outlined">Outlined</BaseButton>
                <BaseButton loading>Loading</BaseButton>
              </div>

              <div class="stack-md">
                <FormField label="Email address" hint="FormField handles helper and error copy.">
                  <BaseInput v-model="formState.email" type="email" placeholder="seeker@nebula.tarot" />
                </FormField>

                <FormField label="Preferred deck">
                  <BaseSelect v-model="formState.deck" :options="deckOptions" placeholder="Select deck" />
                </FormField>

                <div class="ui-kit-inline">
                  <BaseSwitchToggle v-model="notifications" label="Daily reminders" />
                  <BaseSwitchToggle v-model="cinematicAnimations" label="Cinematic reveals" />
                </div>
              </div>
            </div>
          </BaseSurface>

          <BaseSurface padding="lg" tag="section">
            <SectionHeader
              size="sm"
              eyebrow="Tabs"
              title="Single navigation pattern"
              copy="The same tabs primitive is used for dashboard, profile and other screen-level filtering."
            />

            <BaseTabs v-model="activeTab" :items="tabItems">
              <template #panel="{ activeItem }">
                <div class="stack-md">
                  <SectionHeader
                    size="sm"
                    :eyebrow="activeItem?.label || ''"
                    :title="activeItem?.label || ''"
                    :copy="activeItem?.content || ''"
                  />
                  <div class="cluster">
                    <BaseBadge label="Shared primitive" />
                    <BaseBadge variant="subtle" label="Keyboard accessible" />
                  </div>
                </div>
              </template>
            </BaseTabs>
          </BaseSurface>
        </div>

        <div class="ui-kit-grid ui-kit-grid--two">
          <AuthCard eyebrow="Authentication" title="Auth shell" copy="Login and registration screens should stay inside this composite instead of inventing local hero wrappers.">
            <form class="stack-md" @submit.prevent>
              <FormField label="Email">
                <BaseInput model-value="alice@nebula.tarot" type="email" />
              </FormField>

              <FormField label="Password" hint="Minimum 8 characters with letters and numbers.">
                <BaseInput model-value="secure-pass-123" type="password" />
              </FormField>

              <BaseButton full-width>Continue</BaseButton>
            </form>

            <template #footer>
              <p>Footer slots stay inside the same system card.</p>
            </template>
          </AuthCard>

          <div class="page-stack">
            <LoadingCard title="Loading state" copy="Use the shared feedback card when data is still arriving." />
            <EmptyStateCard
              title="Empty state"
              description="Use the shared empty state instead of inventing one-off placeholders for each page."
            >
              <template #actions>
                <BaseButton variant="secondary">Create reading</BaseButton>
              </template>
            </EmptyStateCard>
          </div>
        </div>

        <BaseSurface padding="lg" tag="section">
          <SectionHeader
            size="sm"
            eyebrow="Domain components"
            title="Tarot-specific building blocks"
            copy="These are the components screens should compose directly for the reading experience."
          />

          <div class="ui-kit-domain-grid">
            <TarotCard :card="showcaseCard" size="sm" />
            <CardOfDay :reading="cardOfDayReading" />
          </div>
        </BaseSurface>

        <BaseSurface padding="lg" tag="section">
          <SectionHeader
            size="sm"
            eyebrow="Readings"
            title="Spread and history patterns"
            copy="Production reading screens now reuse the same headings and badge primitives shown here."
          />

          <div class="page-stack">
            <ThreeCardSpread :reading="spreadReading" />
            <ReadingHistory
              :readings="historyReadings"
              empty-title="No readings yet"
              empty-copy="Your saved readings will appear here."
            />
          </div>
        </BaseSurface>
      </div>
    </PageContainer>
  </div>
</template>

<style scoped>
.ui-kit-page {
  position: relative;
  min-height: 100vh;
  padding: var(--ds-space-10) 0 var(--ds-space-12);
  background:
    radial-gradient(circle at top, rgba(212, 175, 55, 0.08), transparent 28%),
    linear-gradient(180deg, var(--ds-color-bg-page), var(--ds-color-bg-page-alt));
  color: var(--ds-color-text-primary);
}

.ui-kit-page > :not(.app-background) {
  position: relative;
  z-index: 1;
}

.ui-kit-hero {
  overflow: hidden;
}

.ui-kit-hero__tokens {
  display: flex;
  flex-wrap: wrap;
  gap: var(--ds-space-3);
  margin-top: var(--ds-space-8);
}

.ui-kit-grid {
  display: grid;
  gap: var(--ds-space-8);
}

.ui-kit-grid--two {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.ui-kit-inline {
  display: flex;
  flex-wrap: wrap;
  gap: var(--ds-space-3);
  align-items: center;
}

.ui-kit-swatch-grid {
  display: grid;
  gap: var(--ds-space-4);
  margin-bottom: var(--ds-space-6);
}

.ui-kit-swatch {
  display: flex;
  align-items: center;
  gap: var(--ds-space-4);
}

.ui-kit-swatch__chip {
  width: 44px;
  height: 44px;
  border: var(--ds-border-width) solid var(--ds-color-border-accent);
  flex-shrink: 0;
}

.ui-kit-swatch__chip--page {
  background: var(--ds-color-bg-page);
}

.ui-kit-swatch__chip--surface {
  background: var(--ds-color-bg-surface);
}

.ui-kit-swatch__chip--accent {
  background: var(--ds-color-accent-primary);
}

.ui-kit-swatch__label {
  font-family: var(--ds-font-ui);
  font-size: var(--ds-text-body-sm);
  font-weight: 600;
  color: var(--ds-color-text-primary);
}

.ui-kit-page code {
  font-family: var(--ds-font-ui);
  font-size: var(--ds-text-caption);
  color: var(--ds-color-accent-primary);
}

.ui-kit-domain-grid {
  display: grid;
  gap: var(--ds-space-8);
  grid-template-columns: minmax(240px, 320px) minmax(0, 1fr);
  align-items: start;
}

@media (max-width: 1080px) {
  .ui-kit-grid--two,
  .ui-kit-domain-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 720px) {
  .ui-kit-page {
    padding: var(--ds-space-6) 0 var(--ds-space-8);
  }

  .ui-kit-grid,
  .ui-kit-domain-grid {
    gap: var(--ds-space-6);
  }

  .ui-kit-hero__tokens {
    margin-top: var(--ds-space-6);
  }
}
</style>
