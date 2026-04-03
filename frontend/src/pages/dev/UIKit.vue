<script setup>
import { ref } from 'vue'

import BaseSwitchToggle from '../../components/base/BaseSwitchToggle.vue'

const activeLanguage = ref('EN')
const themeActive = ref(false)
const finalLanguage = ref('EN')
const finalThemeDark = ref(true)

const defaultTab = ref('overview')
const pillTab = ref('daily')
const boxedTab = ref('seeker')

const animationSpeed = ref(75)
const cardSize = ref(50)
const backgroundOpacity = ref(30)
const toggleDarkMode = ref(true)
const toggleSoundEffects = ref(false)
const toggleAutoShuffle = ref(true)
const toggleShowHints = ref(false)

const defaultTabs = [
  {
    id: 'overview',
    title: 'Overview',
    body: 'The Fool represents new beginnings, innocence, and spontaneity. This card signifies the start of a journey with an open heart and mind, stepping forward without fear into the unknown mysteries of life. It encourages taking risks and embracing new adventures.',
  },
  {
    id: 'interpretation',
    title: 'Interpretation',
    body: "When The Fool appears in a reading, it suggests that now is the time to take a leap of faith. Trust your instincts and don't let fear hold you back. This card reminds you that every expert was once a beginner, and every journey starts with a single step.",
  },
  {
    id: 'reversed',
    title: 'Reversed Meaning',
    body: "In the reversed position, The Fool can indicate recklessness, poor judgment, or naivety. It may suggest that you're taking unnecessary risks without proper consideration. This is a warning to think before you act and consider the consequences of your decisions.",
  },
  {
    id: 'symbolism',
    title: 'Symbolism',
    body: 'The Fool carries a white rose symbolizing purity and innocence. The small dog at his feet represents loyalty and protection. The mountain in the background signifies the challenges ahead, while his colorful clothing reflects his optimistic and carefree nature.',
  },
]

const pillTabs = [
  {
    id: 'daily',
    label: 'Card of Day',
    title: 'Card of the Day',
    body: 'Receive one card each day at 3:00 AM based on your timezone. This deterministic reading provides daily guidance and insight into the energies surrounding you for the next 24 hours.',
  },
  {
    id: 'spread3',
    label: '3-Card Spread',
    title: 'Three Card Spread',
    body: 'A classic spread representing Past, Present, and Future. This reading helps you understand how past events influence your current situation and what potential outcomes lie ahead based on your current path.',
  },
  {
    id: 'spread5',
    label: '5-Card Spread',
    title: 'Five Card Spread',
    body: 'A more comprehensive reading that explores multiple aspects of your situation. Each position offers deeper insight into different factors affecting your question or circumstance.',
  },
]

const boxedTabs = [
  {
    id: 'seeker',
    label: 'Искатель',
    title: 'Искатель (Seeker)',
    body: "The entry level for new users. Seekers can generate up to 3 three-card spreads per day and view reading history for 1 day. This role provides a solid introduction to the platform's core features.",
  },
  {
    id: 'adept',
    label: 'Адепт',
    title: 'Адепт (Adept)',
    body: 'Intermediate users with expanded access. Adepts can generate 5 spreads daily, access 7 days of reading history, and choose from multiple deck options. Perfect for regular practitioners.',
  },
  {
    id: 'master',
    label: 'Мастер',
    title: 'Мастер (Master)',
    body: 'Advanced users with premium features. Masters enjoy 10 daily spreads, 15 days of history, and full access to all decks and customization options. Ideal for serious tarot enthusiasts.',
  },
  {
    id: 'keeper',
    label: 'Смотритель',
    title: 'Смотритель (Keeper)',
    body: 'Administrative access with unlimited spreads, 30 days of history, and full platform management capabilities. Keepers can manage users, upload decks, and configure system settings.',
  },
]

function toggleTheme() {
  themeActive.value = !themeActive.value
}

function toggleLanguage() {
  activeLanguage.value = activeLanguage.value === 'EN' ? 'RU' : 'EN'
}

function toggleFinalLanguage() {
  finalLanguage.value = finalLanguage.value === 'EN' ? 'RU' : 'EN'
}

function toggleFinalTheme() {
  finalThemeDark.value = !finalThemeDark.value
}

function currentDefaultTab() {
  return defaultTabs.find((tab) => tab.id === defaultTab.value) ?? defaultTabs[0]
}

function currentPillTab() {
  return pillTabs.find((tab) => tab.id === pillTab.value) ?? pillTabs[0]
}

function currentBoxedTab() {
  return boxedTabs.find((tab) => tab.id === boxedTab.value) ?? boxedTabs[0]
}
</script>

<template>
  <div class="ui-kit-dark-page">
    <div class="starfield"></div>

    <header class="main-header">
      <div class="header-container">
        <div class="logo-section">
          <img class="logo" src="/branding/nebula-tarot-logo.png" alt="Nebula Tarot" />
          <span class="brand-name">Nebula Tarot</span>
        </div>

        <div class="header-actions">
          <div class="header-toggle-control">
            <button
              type="button"
              class="header-binary-toggle"
              :class="{ 'is-right': activeLanguage === 'RU' }"
              :aria-pressed="activeLanguage === 'RU'"
              @click="toggleLanguage"
            >
              <span class="header-binary-toggle__thumb" aria-hidden="true"></span>
              <span class="header-binary-toggle__option">EN</span>
              <span class="header-binary-toggle__option">RU</span>
            </button>
          </div>

          <div class="header-toggle-control">
            <button
              type="button"
              class="header-binary-toggle"
              :class="{ 'is-right': themeActive }"
              :aria-pressed="themeActive"
              @click="toggleTheme"
            >
              <span class="header-binary-toggle__thumb" aria-hidden="true"></span>
              <span class="header-binary-toggle__option">Light</span>
              <span class="header-binary-toggle__option">Dark</span>
            </button>
          </div>

          <div class="user-profile" role="button" tabindex="0">
            <div class="user-avatar">AM</div>
            <div class="user-info">
              <span class="user-name">Alice Moon</span>
              <span class="user-role">Адепт</span>
            </div>
          </div>
        </div>
      </div>
    </header>

    <div class="container">
      <header class="header">
        <div class="header-ornament">
          <svg viewBox="0 0 200 100" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M100 10 L110 30 L100 50 L90 30 Z" stroke="#D4AF37" stroke-width="1" fill="none" />
            <line x1="85" y1="30" x2="75" y2="30" stroke="#D4AF37" stroke-width="1" />
            <line x1="115" y1="30" x2="125" y2="30" stroke="#D4AF37" stroke-width="1" />
            <path d="M70 30 Q60 35 60 45 Q60 55 70 60" stroke="#D4AF37" stroke-width="1" fill="none" />
            <path d="M130 30 Q140 35 140 45 Q140 55 130 60" stroke="#D4AF37" stroke-width="1" fill="none" />
            <line x1="70" y1="60" x2="130" y2="60" stroke="#D4AF37" stroke-width="1" />
            <circle cx="75" cy="30" r="2" fill="#D4AF37" />
            <circle cx="125" cy="30" r="2" fill="#D4AF37" />
            <circle cx="70" cy="60" r="2" fill="#D4AF37" />
            <circle cx="130" cy="60" r="2" fill="#D4AF37" />
            <text x="50" y="50" font-size="8" fill="#D4AF37">✦</text>
            <text x="145" y="50" font-size="8" fill="#D4AF37">✦</text>
          </svg>
        </div>

        <h1 class="title">Nebula Tarot</h1>
        <p class="subtitle">Mystical Design System • Dark Theme</p>
      </header>

      <section id="typography" class="section">
        <div class="section-header">
          <h2 class="section-title">Typography</h2>
        </div>

        <div class="type-grid">
          <div class="type-item">
            <h1 class="h1-example">The Mystic Path Awaits</h1>
            <p class="type-meta">Heading 1 • Cinzel 56px / 600</p>
          </div>

          <div class="type-item">
            <h2 class="h2-example">Card of the Day</h2>
            <p class="type-meta">Heading 2 • Cinzel 42px / 500</p>
          </div>

          <div class="type-item">
            <h3 class="h3-example">Three Card Spread</h3>
            <p class="type-meta">Heading 3 • Cinzel 32px / 400</p>
          </div>

          <div class="type-item">
            <p class="body-example">
              New beginnings, innocence, spontaneity, and a free spirit. The Fool represents the start
              of a journey with an open heart and mind, stepping forward without fear into the unknown
              mysteries of life.
            </p>
            <p class="type-meta">Body Text • Cormorant Garamond 18px / 400</p>
          </div>
        </div>
      </section>

      <section id="colors" class="section">
        <div class="section-header">
          <h2 class="section-title">Color Palette</h2>
        </div>

        <div class="color-grid">
          <div class="color-category">
            <h3 class="color-category-title">Primary</h3>
            <div class="color-swatch-list">
              <div class="color-swatch">
                <div class="swatch-box" style="background: #D4AF37;"></div>
                <div class="swatch-info">
                  <div class="swatch-name">Gold Primary</div>
                  <div class="swatch-value">#D4AF37</div>
                </div>
              </div>
              <div class="color-swatch">
                <div class="swatch-box" style="background: #F4E5B0;"></div>
                <div class="swatch-info">
                  <div class="swatch-name">Gold Light</div>
                  <div class="swatch-value">#F4E5B0</div>
                </div>
              </div>
              <div class="color-swatch">
                <div class="swatch-box" style="background: #9B7E1F;"></div>
                <div class="swatch-info">
                  <div class="swatch-name">Gold Dark</div>
                  <div class="swatch-value">#9B7E1F</div>
                </div>
              </div>
            </div>
          </div>

          <div class="color-category">
            <h3 class="color-category-title">Background</h3>
            <div class="color-swatch-list">
              <div class="color-swatch">
                <div class="swatch-box" style="background: #000000;"></div>
                <div class="swatch-info">
                  <div class="swatch-name">Void</div>
                  <div class="swatch-value">#000000</div>
                </div>
              </div>
              <div class="color-swatch">
                <div class="swatch-box" style="background: #0D0D0D;"></div>
                <div class="swatch-info">
                  <div class="swatch-name">Dark</div>
                  <div class="swatch-value">#0D0D0D</div>
                </div>
              </div>
              <div class="color-swatch">
                <div class="swatch-box" style="background: #1A1A1A;"></div>
                <div class="swatch-info">
                  <div class="swatch-name">Surface</div>
                  <div class="swatch-value">#1A1A1A</div>
                </div>
              </div>
            </div>
          </div>

          <div class="color-category">
            <h3 class="color-category-title">Text</h3>
            <div class="color-swatch-list">
              <div class="color-swatch">
                <div class="swatch-box" style="background: #FFFFFF;"></div>
                <div class="swatch-info">
                  <div class="swatch-name">Primary</div>
                  <div class="swatch-value">#FFFFFF</div>
                </div>
              </div>
              <div class="color-swatch">
                <div class="swatch-box" style="background: #C4C4C4;"></div>
                <div class="swatch-info">
                  <div class="swatch-name">Secondary</div>
                  <div class="swatch-value">#C4C4C4</div>
                </div>
              </div>
              <div class="color-swatch">
                <div class="swatch-box" style="background: #6B6B6B;"></div>
                <div class="swatch-info">
                  <div class="swatch-name">Muted</div>
                  <div class="swatch-value">#6B6B6B</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section id="buttons" class="section">
        <div class="section-header">
          <h2 class="section-title">Buttons</h2>
        </div>

        <div class="button-showcase">
          <div class="button-group">
            <span class="button-label">Primary</span>
            <button class="btn btn-primary" type="button">Generate Reading</button>
            <button class="btn btn-primary" type="button" disabled>Generate Reading</button>
          </div>

          <div class="button-group">
            <span class="button-label">Secondary</span>
            <button class="btn btn-secondary" type="button">View History</button>
            <button class="btn btn-secondary" type="button" disabled>View History</button>
          </div>

          <div class="button-group">
            <span class="button-label">Ghost</span>
            <button class="btn btn-ghost" type="button">Learn More</button>
            <button class="btn btn-ghost" type="button" disabled>Learn More</button>
          </div>

          <div class="button-group">
            <span class="button-label">Outlined</span>
            <button class="btn btn-outlined" type="button">Get Started</button>
            <button class="btn btn-outlined" type="button" disabled>Get Started</button>
          </div>

          <div class="button-group">
            <span class="button-label">Text Only</span>
            <button class="btn btn-text" type="button">Skip for now</button>
            <button class="btn btn-text" type="button" disabled>Skip for now</button>
          </div>

          <div class="button-group">
            <span class="button-label">Sizes</span>
            <button class="btn btn-primary btn-small" type="button">Small</button>
            <button class="btn btn-primary" type="button">Medium</button>
            <button class="btn btn-primary btn-large" type="button">Large</button>
          </div>

          <div class="button-group">
            <span class="button-label">Icon Only</span>
            <button class="btn btn-primary btn-icon" type="button">✦</button>
            <button class="btn btn-secondary btn-icon" type="button">☽</button>
            <button class="btn btn-ghost btn-icon" type="button">⋆</button>
          </div>

          <div class="button-group">
            <span class="button-label">With Icons</span>
            <button class="btn btn-primary btn-icon-text" type="button">
              <span>✦</span>
              <span>Card of the Day</span>
            </button>
            <button class="btn btn-secondary btn-icon-text" type="button">
              <span>☽</span>
              <span>Three Cards</span>
            </button>
          </div>

          <div class="button-group">
            <span class="button-label">Star Sides</span>
            <button class="btn btn-primary btn-star-sides" type="button">Generate Reading</button>
            <button class="btn btn-outlined btn-star-sides" type="button">Get Started</button>
          </div>

          <div class="button-group button-group--block">
            <span class="button-label button-label--block">Full Width</span>
            <button class="btn btn-primary btn-full" type="button">Continue to Dashboard</button>
          </div>
        </div>
      </section>

      <section id="components" class="section">
        <div class="section-header">
          <h2 class="section-title">Cards</h2>
        </div>

        <div class="cards-showcase">
          <div class="card card--lg">
            <div class="card-header">
              <h3 class="card-title">Card of the Day</h3>
              <span class="card-badge">Today</span>
            </div>
            <div class="card-body">
              <h4 class="card-name">The Fool</h4>
              <div class="tarot-card-back">
                <div class="card-back-ornament">
                  <span class="card-symbol">✦</span>
                </div>
              </div>
              <p class="card-interpretation">
                New beginnings, innocence, spontaneity. The Fool represents the start of a journey
                with an open heart.
              </p>
              <div class="card-tags">
                <span class="tag">Freedom</span>
                <span class="tag">Adventure</span>
                <span class="tag">New Start</span>
              </div>
            </div>
          </div>

          <div class="card card--md">
            <div class="card-header">
              <h3 class="card-title">Three Card Spread</h3>
              <span class="card-badge">3 of 5</span>
            </div>
            <div class="card-body">
              <h4 class="card-name">The Moon</h4>
              <div class="tarot-card-back">
                <div class="card-back-ornament">
                  <span class="card-symbol">☽</span>
                </div>
              </div>
              <p class="card-interpretation">
                Uncertainty, fears, and deep feelings call for caution and trust in intuition.
              </p>
              <div class="card-tags">
                <span class="tag">Illusion</span>
                <span class="tag">Intuition</span>
                <span class="tag">Mystery</span>
              </div>
            </div>
          </div>

          <div class="card card--sm">
            <div class="card-header">
              <h3 class="card-title">The Star</h3>
              <span class="card-badge">Major</span>
            </div>
            <div class="card-body">
              <h4 class="card-name">The Star</h4>
              <div class="tarot-card-back">
                <div class="card-back-ornament">
                  <span class="card-symbol">⋆</span>
                </div>
              </div>
              <p class="card-interpretation">
                After hardship, hope, faith, and quiet healing begin to return.
              </p>
              <div class="card-tags">
                <span class="tag">Hope</span>
                <span class="tag">Inspiration</span>
                <span class="tag">Healing</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="section">
        <div class="section-header">
          <h2 class="section-title">Ornamental Dividers</h2>
        </div>

        <div class="divider-showcase">
          <div class="ornamental-divider">
            <span class="divider-content">✦</span>
          </div>

          <div class="ornamental-divider">
            <div class="divider-complex">
              <span>◯</span>
              <div class="divider-diamond"></div>
              <span>☽</span>
              <div class="divider-diamond"></div>
              <span>◯</span>
              <div class="divider-diamond"></div>
              <span>☾</span>
              <div class="divider-diamond"></div>
              <span>◯</span>
            </div>
          </div>

          <div class="ornamental-divider">
            <span class="divider-content">⋆ · ✦ · ⋆</span>
          </div>

          <div class="ornamental-divider">
            <div class="divider-complex">
              <span>✦</span>
              <span>✧</span>
              <span>⋆</span>
              <span>✧</span>
              <span>✦</span>
            </div>
          </div>
        </div>
      </section>

      <section class="section">
        <div class="section-header">
          <h2 class="section-title">Form Elements</h2>
        </div>

        <div class="form-grid">
          <div class="form-group">
            <label class="form-label">Email Address</label>
            <input type="email" class="form-input" placeholder="your@email.com" />
          </div>

          <div class="form-group">
            <label class="form-label">Password</label>
            <input type="password" class="form-input" placeholder="Enter password" />
          </div>

          <div class="form-group">
            <label class="form-label">Select Deck</label>
            <div class="select-wrapper">
              <select class="form-select select-custom">
                <option>Tarot Classic</option>
                <option>Rider-Waite Modern</option>
                <option>Ethereal Dreams</option>
              </select>
              <span class="select-arrow">▼</span>
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">Email (with error)</label>
            <input type="email" class="form-input input-error" value="invalid-email" />
            <span class="form-error">Please enter a valid email address</span>
          </div>
        </div>
      </section>

      <section id="controls" class="section">
        <div class="section-header">
          <h2 class="section-title">Controls & Inputs</h2>
        </div>

        <div class="controls-grid">
          <div class="control-group">
            <h3 class="control-title">Checkboxes</h3>
            <div class="control-items">
              <label class="checkbox-wrapper">
                <input type="checkbox" class="checkbox-input" checked />
                <span class="checkbox-custom"></span>
                <span class="checkbox-label">Enable daily notifications</span>
              </label>

              <label class="checkbox-wrapper">
                <input type="checkbox" class="checkbox-input" />
                <span class="checkbox-custom"></span>
                <span class="checkbox-label">Show reversed cards</span>
              </label>

              <label class="checkbox-wrapper">
                <input type="checkbox" class="checkbox-input" />
                <span class="checkbox-custom"></span>
                <span class="checkbox-label">Save reading history</span>
              </label>

              <label class="checkbox-wrapper">
                <input type="checkbox" class="checkbox-input" checked />
                <span class="checkbox-custom"></span>
                <span class="checkbox-label">Remember my preferences</span>
              </label>
            </div>
          </div>

          <div class="control-group">
            <h3 class="control-title">Radio Buttons</h3>
            <div class="control-items">
              <label class="radio-wrapper">
                <input type="radio" name="spread-type" class="radio-input" checked />
                <span class="radio-custom"></span>
                <span class="radio-label">Card of the Day</span>
              </label>

              <label class="radio-wrapper">
                <input type="radio" name="spread-type" class="radio-input" />
                <span class="radio-custom"></span>
                <span class="radio-label">Three Card Spread</span>
              </label>

              <label class="radio-wrapper">
                <input type="radio" name="spread-type" class="radio-input" />
                <span class="radio-custom"></span>
                <span class="radio-label">Celtic Cross</span>
              </label>

              <label class="radio-wrapper">
                <input type="radio" name="spread-type" class="radio-input" />
                <span class="radio-custom"></span>
                <span class="radio-label">Custom Spread</span>
              </label>
            </div>
          </div>

          <div class="control-group">
            <h3 class="control-title">Toggle Switches</h3>
            <div class="control-items">
              <BaseSwitchToggle v-model="toggleDarkMode" label="Dark Mode" />
              <BaseSwitchToggle v-model="toggleSoundEffects" label="Sound Effects" />
              <BaseSwitchToggle v-model="toggleAutoShuffle" label="Auto-shuffle" />
              <BaseSwitchToggle v-model="toggleShowHints" label="Show hints" />
            </div>
          </div>

          <div class="control-group">
            <h3 class="control-title">Final Toggles</h3>
            <div class="control-items control-items--final-toggles">
              <div class="final-toggle-row">
                <span class="final-toggle-label">Language</span>
                <button
                  type="button"
                  class="header-binary-toggle header-binary-toggle--demo"
                  :class="{ 'is-right': finalLanguage === 'RU' }"
                  :aria-pressed="finalLanguage === 'RU'"
                  @click="toggleFinalLanguage"
                >
                  <span class="header-binary-toggle__thumb" aria-hidden="true"></span>
                  <span class="header-binary-toggle__option">EN</span>
                  <span class="header-binary-toggle__option">RU</span>
                </button>
              </div>

              <div class="final-toggle-row">
                <span class="final-toggle-label">Theme</span>
                <button
                  type="button"
                  class="header-binary-toggle header-binary-toggle--demo"
                  :class="{ 'is-right': finalThemeDark }"
                  :aria-pressed="finalThemeDark"
                  @click="toggleFinalTheme"
                >
                  <span class="header-binary-toggle__thumb" aria-hidden="true"></span>
                  <span class="header-binary-toggle__option">Light</span>
                  <span class="header-binary-toggle__option">Dark</span>
                </button>
              </div>
            </div>
          </div>

          <div class="control-group">
            <h3 class="control-title">Range Sliders</h3>
            <div class="control-items">
              <div class="range-wrapper">
                <div class="range-label-row">
                  <span class="range-label">Animation Speed</span>
                  <span class="range-value">{{ animationSpeed }}%</span>
                </div>
                <input v-model="animationSpeed" type="range" class="range-input" min="0" max="100" />
              </div>

              <div class="range-wrapper">
                <div class="range-label-row">
                  <span class="range-label">Card Size</span>
                  <span class="range-value">{{ cardSize }}%</span>
                </div>
                <input v-model="cardSize" type="range" class="range-input" min="0" max="100" />
              </div>

              <div class="range-wrapper">
                <div class="range-label-row">
                  <span class="range-label">Background Opacity</span>
                  <span class="range-value">{{ backgroundOpacity }}%</span>
                </div>
                <input
                  v-model="backgroundOpacity"
                  type="range"
                  class="range-input"
                  min="0"
                  max="100"
                />
              </div>
            </div>
          </div>
        </div>
      </section>

      <section id="tabs" class="section">
        <div class="section-header">
          <h2 class="section-title">Tabs & Navigation</h2>
        </div>

        <div class="tabs-sections">
          <div class="tabs-container">
            <h3 class="control-title control-title--with-gap">Default Tabs</h3>
            <div class="tabs-header">
              <button
                v-for="tab in defaultTabs"
                :key="tab.id"
                type="button"
                :class="['tab-button', { active: defaultTab === tab.id }]"
                @click="defaultTab = tab.id"
              >
                {{ tab.id === 'overview' ? 'Overview' : tab.id === 'interpretation' ? 'Interpretation' : tab.id === 'reversed' ? 'Reversed' : 'Symbolism' }}
              </button>
            </div>
            <div class="tabs-content">
              <div class="tab-content active">
                <div class="tab-panel">
                  <h4>{{ currentDefaultTab().title }}</h4>
                  <p>{{ currentDefaultTab().body }}</p>
                </div>
              </div>
            </div>
          </div>

          <div class="tabs-container tabs-pills">
            <h3 class="control-title control-title--with-gap">Pill Tabs</h3>
            <div class="tabs-header">
              <button
                v-for="tab in pillTabs"
                :key="tab.id"
                type="button"
                :class="['tab-button', { active: pillTab === tab.id }]"
                @click="pillTab = tab.id"
              >
                {{ tab.label }}
              </button>
            </div>
            <div class="tabs-content">
              <div class="tab-content active">
                <div class="tab-panel">
                  <h4>{{ currentPillTab().title }}</h4>
                  <p>{{ currentPillTab().body }}</p>
                </div>
              </div>
            </div>
          </div>

          <div class="tabs-container tabs-boxed">
            <h3 class="control-title control-title--with-gap">Boxed Tabs</h3>
            <div class="tabs-header">
              <button
                v-for="tab in boxedTabs"
                :key="tab.id"
                type="button"
                :class="['tab-button', { active: boxedTab === tab.id }]"
                @click="boxedTab = tab.id"
              >
                {{ tab.label }}
              </button>
            </div>
            <div class="tabs-content">
              <div class="tab-content active">
                <div class="tab-panel">
                  <h4>{{ currentBoxedTab().title }}</h4>
                  <p>{{ currentBoxedTab().body }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <footer class="footer">
        <div class="footer-divider">
          <span>✦</span>
        </div>
        <p class="footer-text">Nebula Tarot • Mystical Design System</p>
        <p class="footer-meta">Dark Theme • Version 1.0</p>
      </footer>
    </div>

  </div>
</template>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Cinzel:wght@400;500;600;700&family=Cormorant+Garamond:wght@300;400;500;600;700&family=Inter:wght@300;400;500;600&display=swap");

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.ui-kit-dark-page {
  --gold-primary: #d4af37;
  --gold-light: #f4e5b0;
  --gold-dark: #9b7e1f;
  --bg-void: #000000;
  --bg-dark: #0d0d0d;
  --bg-surface: #1a1a1a;
  --bg-elevated: #242424;
  --text-primary: #ffffff;
  --text-secondary: #c4c4c4;
  --text-muted: #6b6b6b;
  --border-gold: rgba(212, 175, 55, 0.3);
  --border-subtle: rgba(255, 255, 255, 0.1);
  --font-display: "Cinzel", serif;
  --font-body: "Cormorant Garamond", serif;
  --font-sans: "Inter", sans-serif;
  --container-max: 1400px;
  min-height: 100vh;
  background: var(--bg-void);
  color: var(--text-primary);
  font-family: var(--font-body);
  line-height: 1.6;
  overflow-x: hidden;
}

.starfield {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  opacity: 0.6;
}

.starfield::before,
.starfield::after {
  content: "";
  position: absolute;
  width: 2px;
  height: 2px;
  background: white;
  box-shadow:
    100px 200px white,
    300px 400px white,
    500px 100px white,
    700px 300px white,
    900px 500px white,
    200px 600px white,
    400px 150px white,
    600px 450px white,
    800px 250px white,
    1000px 350px white,
    150px 450px white,
    350px 250px white,
    550px 550px white,
    750px 150px white,
    950px 450px white,
    1100px 200px white,
    1300px 400px white,
    250px 100px white,
    450px 500px white,
    650px 300px white;
}

.starfield::after {
  width: 1px;
  height: 1px;
  box-shadow:
    120px 180px white,
    320px 420px white,
    520px 120px white,
    720px 320px white,
    920px 520px white,
    220px 620px white,
    420px 170px white,
    620px 470px white,
    820px 270px white,
    1020px 370px white;
}

.container {
  position: relative;
  z-index: 1;
  max-width: var(--container-max);
  margin: 0 auto;
  padding: 80px 60px;
}

.header {
  text-align: center;
  margin-bottom: 120px;
  position: relative;
}

.header-ornament {
  margin-bottom: 40px;
}

.header-ornament svg {
  width: 200px;
  height: auto;
  opacity: 0.9;
}

.title {
  font-family: var(--font-display);
  font-size: 72px;
  font-weight: 600;
  letter-spacing: 0.15em;
  text-transform: uppercase;
  margin-bottom: 20px;
  background: linear-gradient(135deg, var(--gold-light), var(--gold-primary), var(--gold-dark));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  font-family: var(--font-sans);
  font-size: 14px;
  font-weight: 300;
  letter-spacing: 0.3em;
  text-transform: uppercase;
  color: var(--text-muted);
}

.section {
  margin-bottom: 140px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 30px;
  margin-bottom: 60px;
  position: relative;
}

.section-header::before,
.section-header::after {
  content: "";
  flex: 1;
  height: 1px;
  background: linear-gradient(to right, transparent, var(--border-gold), transparent);
}

.section-title {
  font-family: var(--font-display);
  font-size: 36px;
  font-weight: 500;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--gold-primary);
  white-space: nowrap;
  position: relative;
}

.section-title::before,
.section-title::after {
  content: "✦";
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  font-size: 14px;
  color: var(--gold-primary);
}

.section-title::before {
  left: -25px;
}

.section-title::after {
  right: -25px;
}

.type-grid {
  display: grid;
  gap: 40px;
}

.type-item {
  padding: 40px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  position: relative;
}

.h1-example {
  font-family: var(--font-display);
  font-size: 56px;
  font-weight: 600;
  letter-spacing: 0.05em;
  line-height: 1.2;
  margin-bottom: 12px;
}

.h2-example {
  font-family: var(--font-display);
  font-size: 42px;
  font-weight: 500;
  letter-spacing: 0.05em;
  line-height: 1.2;
  margin-bottom: 12px;
}

.h3-example {
  font-family: var(--font-display);
  font-size: 32px;
  font-weight: 400;
  letter-spacing: 0.05em;
  line-height: 1.3;
  margin-bottom: 12px;
}

.body-example {
  font-family: var(--font-body);
  font-size: 20px;
  font-weight: 400;
  line-height: 1.7;
  color: var(--text-secondary);
  margin-bottom: 12px;
}

.type-meta {
  font-family: var(--font-sans);
  font-size: 11px;
  font-weight: 400;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: var(--text-muted);
}

.color-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 30px;
}

.color-category {
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  padding: 30px;
  position: relative;
}

.color-category-title {
  font-family: var(--font-display);
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.15em;
  text-transform: uppercase;
  color: var(--gold-primary);
  margin-bottom: 20px;
}

.color-swatch-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.color-swatch {
  display: flex;
  align-items: center;
  gap: 16px;
}

.swatch-box {
  width: 60px;
  height: 60px;
  border: 1px solid var(--border-gold);
  flex-shrink: 0;
}

.swatch-info {
  flex: 1;
}

.swatch-name {
  font-family: var(--font-sans);
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.swatch-value {
  font-family: "Courier New", monospace;
  font-size: 11px;
  color: var(--text-muted);
}

.button-showcase {
  display: grid;
  gap: 50px;
}

.button-group {
  display: flex;
  align-items: center;
  gap: 20px;
  flex-wrap: wrap;
  padding: 40px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
}

.button-group--block {
  display: block;
}

.button-label {
  font-family: var(--font-sans);
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--text-muted);
  min-width: 100px;
}

.button-label--block {
  display: block;
  margin-bottom: 20px;
}

.btn {
  padding: 16px 40px;
  font-family: var(--font-sans);
  font-size: 12px;
  font-weight: 500;
  letter-spacing: 0.15em;
  text-transform: uppercase;
  border: none;
  cursor: pointer;
  position: relative;
  background: transparent;
  color: var(--text-primary);
  transition: all 0.3s ease;
}

.btn-primary {
  background: linear-gradient(135deg, var(--gold-dark), var(--gold-primary));
  border: 1px solid var(--gold-primary);
  box-shadow: 0 0 20px rgba(212, 175, 55, 0.2);
}

.btn-primary:hover {
  box-shadow: 0 0 30px rgba(212, 175, 55, 0.4);
  transform: translateY(-2px);
}

.btn-secondary {
  border: 1px solid var(--border-gold);
  background: var(--bg-elevated);
}

.btn-secondary:hover {
  border-color: var(--gold-primary);
  box-shadow: 0 0 20px rgba(212, 175, 55, 0.15);
}

.btn-ghost {
  border: 1px solid var(--border-subtle);
}

.btn-ghost:hover {
  border-color: var(--border-gold);
  background: rgba(212, 175, 55, 0.05);
}

.btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
  box-shadow: none;
}

.cards-showcase {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 40px;
}

.card {
  --card-frame-max: 320px;
  --card-header-padding: 30px;
  --card-header-gap: 12px;
  --card-body-padding: 30px;
  --card-title-size: 20px;
  --card-name-size: 24px;
  --card-name-margin: 16px;
  --card-copy-size: 16px;
  --card-copy-line-height: 1.7;
  --card-copy-margin: 20px;
  --card-tag-gap: 10px;
  --card-badge-padding-y: 6px;
  --card-badge-padding-x: 16px;
  --card-badge-font-size: 10px;
  --card-badge-letter-spacing: 0.1em;
  --card-back-margin: 20px;
  --card-symbol-size: 48px;
  --card-tag-padding-y: 6px;
  --card-tag-padding-x: 14px;
  --card-tag-font-size: 10px;
  --card-tag-letter-spacing: 0.05em;
  background: transparent;
  border: 0;
  box-shadow: none;
  position: relative;
  width: min(100%, var(--card-frame-max));
  justify-self: center;
}

.card--lg {
  --card-frame-max: 320px;
}

.card--md {
  --card-frame-max: 280px;
  --card-header-padding: 24px;
  --card-header-gap: 10px;
  --card-body-padding: 24px;
  --card-title-size: 18px;
  --card-name-size: 21px;
  --card-name-margin: 14px;
  --card-copy-size: 15px;
  --card-copy-line-height: 1.65;
  --card-copy-margin: 18px;
  --card-tag-gap: 8px;
  --card-badge-padding-y: 5px;
  --card-badge-padding-x: 14px;
  --card-badge-font-size: 9px;
  --card-badge-letter-spacing: 0.09em;
  --card-back-margin: 18px;
  --card-symbol-size: 42px;
  --card-tag-padding-y: 5px;
  --card-tag-padding-x: 12px;
  --card-tag-font-size: 9px;
  --card-tag-letter-spacing: 0.04em;
}

.card--sm {
  --card-frame-max: 220px;
  --card-header-padding: 18px;
  --card-header-gap: 8px;
  --card-body-padding: 18px;
  --card-title-size: 15px;
  --card-name-size: 18px;
  --card-name-margin: 12px;
  --card-copy-size: 14px;
  --card-copy-line-height: 1.55;
  --card-copy-margin: 16px;
  --card-tag-gap: 6px;
  --card-badge-padding-y: 4px;
  --card-badge-padding-x: 10px;
  --card-badge-font-size: 8px;
  --card-badge-letter-spacing: 0.08em;
  --card-back-margin: 16px;
  --card-symbol-size: 34px;
  --card-tag-padding-y: 4px;
  --card-tag-padding-x: 10px;
  --card-tag-font-size: 8px;
  --card-tag-letter-spacing: 0.03em;
}

.card-header {
  padding: var(--card-header-padding);
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--card-header-gap);
}

.card-title {
  font-family: var(--font-display);
  font-size: var(--card-title-size);
  font-weight: 500;
  letter-spacing: 0.05em;
}

.card-badge {
  padding: var(--card-badge-padding-y) var(--card-badge-padding-x);
  font-family: var(--font-sans);
  font-size: var(--card-badge-font-size);
  font-weight: 500;
  letter-spacing: var(--card-badge-letter-spacing);
  line-height: 1;
  white-space: nowrap;
  text-transform: uppercase;
  background: rgba(212, 175, 55, 0.1);
  border: 1px solid var(--border-gold);
  color: var(--gold-primary);
}

.card-body {
  padding: var(--card-body-padding);
}

.tarot-card-back {
  width: 100%;
  aspect-ratio: 2 / 3;
  background: var(--bg-dark);
  border: 2px solid var(--gold-primary);
  margin-bottom: var(--card-back-margin);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  transition: transform 0.4s ease, box-shadow 0.4s ease, border-color 0.4s ease;
}

.tarot-card-back:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.6);
  border-color: var(--gold-light);
}

.card-back-ornament {
  width: 80%;
  height: 80%;
  border: 1px solid var(--gold-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.card-back-ornament::before,
.card-back-ornament::after {
  content: "";
  position: absolute;
  background: var(--gold-primary);
}

.card-back-ornament::before {
  width: 60%;
  height: 1px;
  top: 50%;
  left: 20%;
}

.card-back-ornament::after {
  width: 1px;
  height: 60%;
  left: 50%;
  top: 20%;
}

.card-symbol {
  font-size: var(--card-symbol-size);
  color: var(--gold-primary);
  position: relative;
  z-index: 1;
}

.card-name {
  font-family: var(--font-display);
  font-size: var(--card-name-size);
  font-weight: 400;
  letter-spacing: 0.05em;
  text-align: center;
  line-height: 1.2;
  margin-bottom: var(--card-name-margin);
}

.card-interpretation {
  font-family: var(--font-body);
  font-size: var(--card-copy-size);
  line-height: var(--card-copy-line-height);
  color: var(--text-secondary);
  margin-bottom: var(--card-copy-margin);
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--card-tag-gap);
  justify-content: center;
}

.tag {
  padding: var(--card-tag-padding-y) var(--card-tag-padding-x);
  font-family: var(--font-sans);
  font-size: var(--card-tag-font-size);
  font-weight: 400;
  letter-spacing: var(--card-tag-letter-spacing);
  line-height: 1;
  background: rgba(212, 175, 55, 0.1);
  border: 1px solid var(--border-gold);
  color: var(--text-secondary);
}

.divider-showcase {
  display: flex;
  flex-direction: column;
  gap: 60px;
}

.ornamental-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 30px;
  color: var(--gold-primary);
}

.ornamental-divider::before,
.ornamental-divider::after {
  content: "";
  flex: 1;
  height: 1px;
  background: linear-gradient(to right, transparent, var(--border-gold), transparent);
}

.divider-content {
  font-size: 24px;
  letter-spacing: 0.5em;
}

.divider-complex {
  display: flex;
  align-items: center;
  gap: 15px;
}

.divider-diamond {
  width: 8px;
  height: 8px;
  background: var(--gold-primary);
  transform: rotate(45deg);
}

.form-grid {
  display: grid;
  gap: 30px;
  max-width: 600px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.form-label {
  font-family: var(--font-sans);
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--text-muted);
}

.form-input,
.form-select {
  padding: 16px 20px;
  font-family: var(--font-body);
  font-size: 16px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-subtle);
  color: var(--text-primary);
  transition: all 0.3s ease;
}

.form-input:focus,
.form-select:focus {
  outline: none;
  border-color: var(--gold-primary);
  box-shadow: 0 0 20px rgba(212, 175, 55, 0.15);
}

.form-input::placeholder {
  color: var(--text-muted);
}

.form-error {
  font-family: var(--font-sans);
  font-size: 12px;
  color: #ff4444;
}

.input-error {
  border-color: #ff4444;
}

.footer {
  text-align: center;
  padding-top: 80px;
  margin-top: 120px;
  border-top: 1px solid var(--border-subtle);
}

.footer-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-bottom: 30px;
  color: var(--gold-primary);
}

.footer-divider::before,
.footer-divider::after {
  content: "";
  flex: 1;
  height: 1px;
  max-width: 200px;
  background: linear-gradient(to right, transparent, var(--border-gold), transparent);
}

.footer-text {
  font-family: var(--font-display);
  font-size: 16px;
  letter-spacing: 0.1em;
  color: var(--text-secondary);
  margin-bottom: 10px;
}

.footer-meta {
  font-family: var(--font-sans);
  font-size: 11px;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: var(--text-muted);
}

.main-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: rgba(0, 0, 0, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--border-gold);
}

.header-container {
  max-width: var(--container-max);
  margin: 0 auto;
  padding: 20px 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 20px;
  justify-content: flex-start;
  flex-shrink: 0;
}

.logo {
  width: 56px;
  height: 56px;
  object-fit: contain;
  filter: drop-shadow(0 0 12px rgba(212, 175, 55, 0.18));
}

.brand-name {
  font-family: var(--font-display);
  font-size: 20px;
  font-weight: 600;
  letter-spacing: 0.15em;
  text-transform: uppercase;
  color: var(--gold-primary);
}

.header-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  flex-wrap: nowrap;
  margin-left: auto;
}

.header-toggle-control {
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: visible;
}

.header-binary-toggle {
  position: relative;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  align-items: center;
  width: 148px;
  padding: 4px;
  border: 1px solid var(--border-subtle);
  background: var(--bg-surface);
  color: var(--text-secondary);
  cursor: pointer;
  transition:
    border-color 0.3s ease,
    box-shadow 0.3s ease,
    background 0.3s ease;
}

.header-binary-toggle:hover {
  border-color: var(--gold-primary);
  box-shadow: 0 0 16px rgba(212, 175, 55, 0.14);
}

.header-binary-toggle__thumb {
  position: absolute;
  top: 4px;
  left: 4px;
  width: calc(50% - 4px);
  height: calc(100% - 8px);
  background: linear-gradient(135deg, var(--gold-dark), var(--gold-primary));
  transition: transform 0.3s ease;
}

.header-binary-toggle.is-right .header-binary-toggle__thumb {
  transform: translateX(100%);
}

.header-binary-toggle__option {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 34px;
  font-family: var(--font-sans);
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  transition: color 0.3s ease;
}

.header-binary-toggle:not(.is-right) .header-binary-toggle__option:first-of-type,
.header-binary-toggle.is-right .header-binary-toggle__option:last-of-type {
  color: var(--bg-void);
}

.header-binary-toggle:not(.is-right) .header-binary-toggle__option:last-of-type,
.header-binary-toggle.is-right .header-binary-toggle__option:first-of-type {
  color: var(--text-secondary);
}

.user-profile {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 8px 16px;
  background: var(--bg-surface);
  border: 1px solid var(--border-gold);
  cursor: pointer;
  transition: all 0.3s ease;
}

.user-profile:hover {
  box-shadow: 0 0 20px rgba(212, 175, 55, 0.2);
  border-color: var(--gold-light);
}

.user-avatar {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, var(--gold-dark), var(--gold-primary));
  border: 1px solid var(--gold-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-display);
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.user-name {
  font-family: var(--font-sans);
  font-size: 12px;
  font-weight: 500;
  color: var(--text-primary);
}

.user-role {
  font-family: var(--font-sans);
  font-size: 10px;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: var(--gold-primary);
}

.btn-small {
  padding: 10px 24px;
  font-size: 10px;
}

.btn-large {
  padding: 20px 50px;
  font-size: 14px;
}

.btn-full {
  width: 100%;
  justify-content: center;
}

.btn-icon {
  width: 40px;
  height: 40px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-icon-text {
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.btn-outlined {
  background: transparent;
  border: 2px solid var(--gold-primary);
  color: var(--gold-primary);
}

.btn-outlined:hover {
  background: var(--gold-primary);
  color: var(--bg-void);
  box-shadow: 0 0 25px rgba(212, 175, 55, 0.3);
}

.btn-text {
  background: transparent;
  border: none;
  color: var(--gold-primary);
  padding: 12px 24px;
}

.btn-text:hover {
  background: rgba(212, 175, 55, 0.1);
}

.btn-star-sides {
  display: inline-flex;
  align-items: center;
  gap: 15px;
  padding: 16px 30px;
}

.btn-star-sides::before,
.btn-star-sides::after {
  content: "✦";
  font-size: 16px;
  color: var(--gold-primary);
}

.tabs-container {
  padding: 40px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
}

.tabs-header {
  display: flex;
  gap: 0;
  border-bottom: 1px solid var(--border-subtle);
  margin-bottom: 30px;
}

.tab-button {
  padding: 16px 32px;
  font-family: var(--font-sans);
  font-size: 12px;
  font-weight: 500;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  background: transparent;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  position: relative;
  transition: all 0.3s ease;
  border-bottom: 2px solid transparent;
}

.tab-button:hover {
  color: var(--text-primary);
  background: rgba(255, 255, 255, 0.03);
}

.tab-button.active {
  color: var(--gold-primary);
  border-bottom-color: var(--gold-primary);
}

.tab-button.active::before {
  content: "✦";
  position: absolute;
  top: 50%;
  left: 10px;
  transform: translateY(-50%);
  font-size: 10px;
}

.tab-content {
  display: none;
  animation: fadeIn 0.3s ease;
}

.tab-content.active {
  display: block;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.tab-panel {
  padding: 20px 0;
}

.tab-panel h4 {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 500;
  color: var(--gold-primary);
  margin-bottom: 12px;
}

.tab-panel p {
  font-family: var(--font-body);
  font-size: 16px;
  line-height: 1.7;
  color: var(--text-secondary);
  margin: 0;
}

.tabs-pills .tab-button {
  border: 1px solid var(--border-subtle);
  border-radius: 0;
  margin-right: 10px;
  border-bottom: 1px solid var(--border-subtle);
}

.tabs-pills .tab-button.active {
  background: rgba(212, 175, 55, 0.1);
  border-color: var(--gold-primary);
  border-bottom-color: var(--gold-primary);
}

.tabs-pills .tabs-header {
  border-bottom: none;
  margin-bottom: 20px;
}

.tabs-boxed .tab-button {
  border: 1px solid var(--border-subtle);
  margin-right: -1px;
}

.tabs-boxed .tab-button.active {
  background: var(--bg-elevated);
  border-color: var(--gold-primary);
  z-index: 1;
}

.tabs-boxed .tabs-header {
  border-bottom: none;
}

.controls-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 40px;
}

.control-group {
  padding: 40px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
}

.control-title {
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 500;
  letter-spacing: 0.05em;
  color: var(--gold-primary);
  margin-bottom: 24px;
}

.control-title--with-gap {
  margin-bottom: 24px;
}

.control-title--no-margin {
  margin-bottom: 0;
}

.control-items {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.checkbox-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  user-select: none;
}

.checkbox-input {
  display: none;
}

.checkbox-custom {
  width: 22px;
  height: 22px;
  border: 1px solid var(--border-gold);
  background: transparent;
  position: relative;
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.checkbox-input:checked + .checkbox-custom {
  background: var(--gold-primary);
  border-color: var(--gold-primary);
  box-shadow: 0 0 10px rgba(212, 175, 55, 0.3);
}

.checkbox-input:checked + .checkbox-custom::after {
  content: "✓";
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: var(--bg-void);
  font-size: 14px;
  font-weight: bold;
}

.checkbox-label {
  font-family: var(--font-sans);
  font-size: 13px;
  color: var(--text-secondary);
}

.checkbox-wrapper:hover .checkbox-custom {
  border-color: var(--gold-primary);
}

.radio-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  user-select: none;
}

.radio-input {
  display: none;
}

.radio-custom {
  width: 22px;
  height: 22px;
  border: 1px solid var(--border-gold);
  border-radius: 50%;
  background: transparent;
  position: relative;
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.radio-input:checked + .radio-custom {
  border-color: var(--gold-primary);
  box-shadow: 0 0 10px rgba(212, 175, 55, 0.3);
}

.radio-input:checked + .radio-custom::after {
  content: "";
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 10px;
  height: 10px;
  background: var(--gold-primary);
  border-radius: 50%;
}

.radio-label {
  font-family: var(--font-sans);
  font-size: 13px;
  color: var(--text-secondary);
}

.radio-wrapper:hover .radio-custom {
  border-color: var(--gold-primary);
}

.control-items--final-toggles {
  gap: 18px;
}

.final-toggle-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.final-toggle-label {
  font-family: var(--font-sans);
  font-size: 13px;
  color: var(--text-secondary);
}

.header-binary-toggle--demo {
  width: 168px;
  flex-shrink: 0;
}

.range-wrapper {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.range-label-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.range-label {
  font-family: var(--font-sans);
  font-size: 13px;
  color: var(--text-secondary);
}

.range-value {
  font-family: var(--font-sans);
  font-size: 13px;
  font-weight: 600;
  color: var(--gold-primary);
}

.range-input {
  -webkit-appearance: none;
  width: 100%;
  height: 4px;
  background: var(--bg-elevated);
  outline: none;
  border: 1px solid var(--border-subtle);
}

.range-input::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 20px;
  height: 20px;
  background: var(--gold-primary);
  cursor: pointer;
  border: 2px solid var(--gold-primary);
  box-shadow: 0 0 10px rgba(212, 175, 55, 0.5);
}

.range-input::-moz-range-thumb {
  width: 20px;
  height: 20px;
  background: var(--gold-primary);
  cursor: pointer;
  border: 2px solid var(--gold-primary);
  box-shadow: 0 0 10px rgba(212, 175, 55, 0.5);
}

.select-wrapper {
  position: relative;
}

.select-custom {
  appearance: none;
  width: 100%;
  padding: 14px 40px 14px 16px;
  font-family: var(--font-sans);
  font-size: 13px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-subtle);
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.3s ease;
}

.select-custom:focus {
  outline: none;
  border-color: var(--gold-primary);
  box-shadow: 0 0 15px rgba(212, 175, 55, 0.15);
}

.select-arrow {
  position: absolute;
  right: 14px;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  color: var(--gold-primary);
  font-size: 10px;
}

.tabs-sections {
  display: grid;
  gap: 60px;
}

@media (max-width: 768px) {
  .header-container {
    flex-direction: column;
    padding: 12px 16px;
    gap: 14px;
  }

  .header-actions {
    margin-left: 0;
    flex-wrap: wrap;
    justify-content: center;
    gap: 8px;
  }

  .final-toggle-row {
    flex-direction: column;
    align-items: stretch;
  }

  .header-binary-toggle--demo {
    width: 100%;
  }

  .user-profile {
    justify-content: center;
    gap: 8px;
    padding: 6px 10px;
  }

  .container {
    padding: 28px 16px;
  }

  .title {
    font-size: 34px;
  }

  .section-title {
    font-size: 20px;
  }

  .h1-example {
    font-size: 30px;
  }

  .h2-example {
    font-size: 24px;
  }

  .header {
    margin-bottom: 48px;
  }

  .section {
    margin-bottom: 48px;
  }

  .section-header,
  .tabs-sections,
  .divider-showcase {
    gap: 28px;
  }

  .type-grid,
  .color-grid,
  .controls-grid,
  .form-grid {
    gap: 20px;
  }

  .btn {
    padding: 12px 24px;
    font-size: 11px;
  }

  .btn-large {
    padding: 14px 28px;
    font-size: 12px;
  }

  .btn-small,
  .btn-text {
    padding: 8px 14px;
  }

  .btn-icon {
    width: 34px;
    height: 34px;
  }

  .btn-star-sides {
    gap: 10px;
    padding: 12px 20px;
  }

  .tabs-container,
  .control-group {
    padding: 20px;
  }

  .tabs-header {
    margin-bottom: 18px;
    overflow-x: auto;
    scrollbar-width: none;
  }

  .tabs-header::-webkit-scrollbar {
    display: none;
  }

  .tab-button {
    flex-shrink: 0;
    padding: 12px 18px;
    font-size: 10px;
  }

  .tab-panel {
    padding: 12px 0;
  }

  .tab-panel h4 {
    font-size: 16px;
    margin-bottom: 8px;
  }

  .tab-panel p {
    font-size: 14px;
    line-height: 1.6;
  }

  .control-title {
    margin-bottom: 16px;
    font-size: 14px;
  }

  .control-items {
    gap: 14px;
  }

  .form-input,
  .form-select,
  .select-custom {
    padding: 12px 14px;
    font-size: 14px;
  }

  .select-custom {
    padding-right: 36px;
  }

  .range-input::-webkit-slider-thumb {
    width: 16px;
    height: 16px;
  }

  .range-input::-moz-range-thumb {
    width: 16px;
    height: 16px;
  }

  .header-binary-toggle {
    width: 130px;
  }

  .header-binary-toggle__option {
    min-height: 30px;
    font-size: 10px;
  }

  .logo {
    width: 42px;
    height: 42px;
  }

  .brand-name {
    font-size: 16px;
  }

  .user-avatar {
    width: 26px;
    height: 26px;
    font-size: 12px;
  }

  .user-name {
    font-size: 11px;
  }

  .user-role {
    font-size: 9px;
  }

  .cards-showcase {
    gap: 20px;
  }

  .card {
    --card-frame-width: 100%;
    --card-body-padding: 18px;
    --card-title-size: 14px;
    --card-header-gap: 10px;
    --card-badge-padding-y: 4px;
    --card-badge-padding-x: 8px;
    --card-badge-font-size: 8px;
    --card-name-size: 18px;
    --card-name-margin: 10px;
    --card-copy-size: 14px;
    --card-copy-line-height: 1.55;
    --card-copy-margin: 12px;
    --card-tag-gap: 6px;
    --card-tag-padding-y: 4px;
    --card-tag-padding-x: 8px;
    --card-tag-font-size: 8px;
    --card-symbol-size: 32px;
    --card-back-margin: 12px;
  }

  .cards-showcase {
    grid-template-columns: 1fr;
  }

  .button-group {
    flex-direction: column;
    align-items: stretch;
  }

  .button-label {
    min-width: auto;
  }

  .footer {
    padding-top: 48px;
    margin-top: 56px;
  }
}
</style>
