<script setup>
const colorTokens = [
  { name: 'Mystic Purple', value: '#8B5CF6', className: 'primary-purple' },
  { name: 'Celestial Gold', value: '#F59E0B', className: 'primary-gold' },
  { name: 'Rose Gold', value: '#F3B8C6', className: 'primary-rose' },
  { name: 'Deep Void', value: '#0A0A14', className: 'bg-primary' },
  { name: 'Dark Surface', value: '#1A1A2E', className: 'bg-surface' },
  { name: 'Elevated', value: '#252540', className: 'bg-elevated' },
  { name: 'Primary Text', value: '#F8F9FA', className: 'text-primary' },
  { name: 'Secondary Text', value: '#B8B9C1', className: 'text-secondary' },
  { name: 'Muted Text', value: '#6B6B7E', className: 'text-muted' },
  { name: 'Success', value: '#10B981', className: 'semantic-success' },
  { name: 'Warning', value: '#F59E0B', className: 'semantic-warning' },
  { name: 'Error', value: '#EF4444', className: 'semantic-error' },
]

const icons = [
  ['☽', 'Crescent Moon'],
  ['◯', 'Full Moon'],
  ['☾', 'Waning Moon'],
  ['✦', 'Star'],
  ['✧', 'Small Star'],
  ['⋆', 'Spark'],
  ['◐', 'Half Moon'],
  ['◑', 'Reverse Half'],
  ['☆', 'Outlined Star'],
  ['★', 'Filled Star'],
  ['◈', 'Diamond'],
  ['✺', 'Burst'],
]
</script>

<template>
  <div class="ui-kit">
    <div class="ui-kit__header">
      <div class="celestial-border">
        <h1 class="ui-kit__title">
          <span class="moon-phases">☽◯☾</span>
          Tarot Daily UI Kit
          <span class="stars">✦✧✦</span>
        </h1>
        <p class="ui-kit__subtitle">Mystical Modern Design System</p>
      </div>
      <div class="ui-kit__pill-row">
        <span class="kit-pill">Dev only</span>
        <span class="kit-pill">No auth</span>
        <span class="kit-pill">Autonomous</span>
        <span class="kit-pill">Frontend source of truth</span>
      </div>
    </div>

    <section class="ui-kit__section">
      <div class="section-divider">
        <span class="divider-icon">✦</span>
        <h2 class="section-title">Typography</h2>
        <span class="divider-icon">✦</span>
      </div>
      <div class="typography-grid">
        <div class="type-example">
          <h1 class="heading-1">Heading 1 • The Mystic Path</h1>
          <p class="type-meta">48px • Cinzel • 500</p>
        </div>
        <div class="type-example">
          <h2 class="heading-2">Heading 2 • Card of the Day</h2>
          <p class="type-meta">36px • Cinzel • 500</p>
        </div>
        <div class="type-example">
          <h3 class="heading-3">Heading 3 • Three Card Spread</h3>
          <p class="type-meta">28px • Cinzel • 400</p>
        </div>
        <div class="type-example">
          <h4 class="heading-4">Heading 4 • The Fool's Journey</h4>
          <p class="type-meta">20px • Cinzel • 400</p>
        </div>
        <div class="type-example">
          <p class="body-text">
            Body Text • New beginnings, innocence, spontaneity, and a free spirit. The Fool represents
            the start of a journey with an open heart and mind.
          </p>
          <p class="type-meta">16px • Inter • 400</p>
        </div>
        <div class="type-example">
          <p class="body-small">Small Text • Keywords: freedom, adventure, new start, innocence</p>
          <p class="type-meta">14px • Inter • 400</p>
        </div>
        <div class="type-example">
          <p class="caption">Caption • Generated at 3:00 AM • Next reset in 21h</p>
          <p class="type-meta">12px • Inter • 400</p>
        </div>
      </div>
    </section>

    <section class="ui-kit__section">
      <div class="section-divider">
        <span class="divider-icon">☽</span>
        <h2 class="section-title">Color Palette</h2>
        <span class="divider-icon">☾</span>
      </div>
      <div class="color-grid">
        <div v-for="token in colorTokens" :key="token.name" class="color-swatch">
          <div class="swatch" :class="token.className"></div>
          <span class="swatch-name">{{ token.name }}</span>
          <span class="swatch-value">{{ token.value }}</span>
        </div>
      </div>
    </section>

    <section class="ui-kit__section">
      <div class="section-divider">
        <span class="divider-icon">✧</span>
        <h2 class="section-title">Buttons</h2>
        <span class="divider-icon">✧</span>
      </div>
      <div class="button-showcase">
        <div class="button-row">
          <h4 class="showcase-label">Primary</h4>
          <button class="btn btn-primary">Generate Reading</button>
          <button class="btn btn-primary" disabled>Generate Reading</button>
        </div>
        <div class="button-row">
          <h4 class="showcase-label">Secondary</h4>
          <button class="btn btn-secondary">View History</button>
          <button class="btn btn-secondary" disabled>View History</button>
        </div>
        <div class="button-row">
          <h4 class="showcase-label">Ghost</h4>
          <button class="btn btn-ghost">Learn More</button>
          <button class="btn btn-ghost" disabled>Learn More</button>
        </div>
        <div class="button-row">
          <h4 class="showcase-label">Sizes</h4>
          <button class="btn btn-primary btn-small">Small</button>
          <button class="btn btn-primary">Medium</button>
          <button class="btn btn-primary btn-large">Large</button>
        </div>
        <div class="button-row">
          <h4 class="showcase-label">States</h4>
          <button class="btn btn-primary is-hover">Hover</button>
          <button class="btn btn-primary is-focus">Focus</button>
          <button class="btn btn-secondary is-loading">Loading</button>
        </div>
      </div>
    </section>

    <section class="ui-kit__section">
      <div class="section-divider">
        <span class="divider-icon">◯</span>
        <h2 class="section-title">Form Elements</h2>
        <span class="divider-icon">◯</span>
      </div>
      <div class="form-showcase">
        <div class="form-group">
          <label class="form-label">Email Address</label>
          <input type="email" class="form-input" placeholder="your@email.com" />
        </div>
        <div class="form-group">
          <label class="form-label">Email Address (with error)</label>
          <input type="email" class="form-input input-error" value="invalid-email" />
          <span class="form-error">Please enter a valid email address</span>
        </div>
        <div class="form-group">
          <label class="form-label">Password</label>
          <input type="password" class="form-input" placeholder="••••••••" />
          <span class="form-hint">Minimum 8 characters with letter and number</span>
        </div>
        <div class="form-group">
          <label class="form-label">Select Deck</label>
          <select class="form-select">
            <option>Tarot Classic</option>
            <option>Rider-Waite Modern</option>
            <option>Ethereal Dreams</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">Timezone</label>
          <select class="form-select">
            <option>Europe/Amsterdam</option>
            <option>America/New_York</option>
            <option>Asia/Tokyo</option>
          </select>
        </div>
        <div class="form-group">
          <label class="checkbox-label">
            <input type="checkbox" class="form-checkbox" />
            <span class="checkbox-text">I agree to the terms and conditions</span>
          </label>
        </div>
        <div class="form-group">
          <label class="radio-label">
            <input type="radio" name="theme" class="form-radio" checked />
            <span class="radio-text">Dark Theme</span>
          </label>
          <label class="radio-label">
            <input type="radio" name="theme" class="form-radio" />
            <span class="radio-text">Light Theme</span>
          </label>
        </div>
      </div>
    </section>

    <section class="ui-kit__section">
      <div class="section-divider">
        <span class="divider-icon">☆</span>
        <h2 class="section-title">Cards & Containers</h2>
        <span class="divider-icon">☆</span>
      </div>
      <div class="cards-showcase">
        <div class="card">
          <div class="card-header">
            <h3 class="card-title">Card of the Day</h3>
            <span class="card-badge">Today</span>
          </div>
          <div class="card-body">
            <div class="tarot-card-placeholder">
              <div class="card-back">
                <span class="moon-symbol">☽✦☾</span>
              </div>
            </div>
            <h4 class="card-name">The Fool</h4>
            <p class="card-interpretation">
              New beginnings, innocence, spontaneity, and a free spirit. The Fool represents the start
              of a journey.
            </p>
            <div class="card-tags">
              <span class="tag">Freedom</span>
              <span class="tag">Adventure</span>
              <span class="tag">New Start</span>
            </div>
          </div>
          <div class="card-footer">
            <span class="card-meta">Next card in 14h 32m</span>
          </div>
        </div>

        <div class="card card-elevated">
          <div class="card-header">
            <h3 class="card-title">Three Card Spread</h3>
            <span class="card-badge badge-purple">3 of 5 today</span>
          </div>
          <div class="card-body">
            <div class="spread-preview">
              <div class="mini-card"></div>
              <div class="mini-card"></div>
              <div class="mini-card"></div>
            </div>
            <button class="btn btn-primary btn-full">Generate New Spread</button>
          </div>
        </div>

        <div class="card card-outline">
          <div class="card-header">
            <span class="card-icon">☽</span>
            <h4 class="card-title-small">Reading History</h4>
          </div>
          <div class="card-body">
            <p class="card-text">Your history is kept for 7 days based on your current role.</p>
          </div>
        </div>
      </div>
    </section>

    <section class="ui-kit__section">
      <div class="section-divider">
        <span class="divider-icon">✦</span>
        <h2 class="section-title">Dividers & Ornaments</h2>
        <span class="divider-icon">✦</span>
      </div>
      <div class="dividers-showcase">
        <div class="ornamental-divider divider-1"><span class="divider-symbol">✦</span></div>
        <div class="ornamental-divider divider-2">
          <span class="divider-left">◯</span>
          <span class="divider-center">☽◯☾</span>
          <span class="divider-right">◯</span>
        </div>
        <div class="ornamental-divider divider-3"><span class="moon-phases-full">☽ ◐ ○ ◑ ☾</span></div>
      </div>
    </section>

    <section class="ui-kit__section">
      <div class="section-divider">
        <span class="divider-icon">◐</span>
        <h2 class="section-title">Modals & Overlays</h2>
        <span class="divider-icon">◑</span>
      </div>
      <div class="modal-showcase">
        <div class="modal-preview">
          <div class="modal">
            <div class="modal-header">
              <h3 class="modal-title">Daily Limit Reached</h3>
              <button class="modal-close">✕</button>
            </div>
            <div class="modal-body">
              <p>You've reached your daily limit of 5 spreads.</p>
              <p class="modal-meta">Your limit resets tomorrow at 3:00 AM</p>
            </div>
            <div class="modal-footer">
              <button class="btn btn-secondary">View History</button>
              <button class="btn btn-primary">Understood</button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="ui-kit__section">
      <div class="section-divider">
        <span class="divider-icon">⋆</span>
        <h2 class="section-title">Badges & Tags</h2>
        <span class="divider-icon">⋆</span>
      </div>
      <div class="badges-showcase">
        <div class="badge-group">
          <h4 class="showcase-label">Status Badges</h4>
          <span class="badge badge-success">Active</span>
          <span class="badge badge-warning">Processing</span>
          <span class="badge badge-error">Disabled</span>
          <span class="badge badge-purple">Today</span>
        </div>
        <div class="badge-group">
          <h4 class="showcase-label">Role Badges</h4>
          <span class="role-badge role-seeker">Искатель</span>
          <span class="role-badge role-adept">Адепт</span>
          <span class="role-badge role-master">Мастер</span>
          <span class="role-badge role-keeper">Смотритель</span>
        </div>
      </div>
    </section>

    <section class="ui-kit__section">
      <div class="section-divider">
        <span class="divider-icon">★</span>
        <h2 class="section-title">Mystical Icons</h2>
        <span class="divider-icon">☆</span>
      </div>
      <div class="icons-grid">
        <div v-for="[glyph, label] in icons" :key="label" class="icon-item">
          <span class="icon-display">{{ glyph }}</span>
          <span class="icon-label">{{ label }}</span>
        </div>
      </div>
    </section>

    <section class="ui-kit__section">
      <div class="section-divider">
        <span class="divider-icon">◐</span>
        <h2 class="section-title">Loading States</h2>
        <span class="divider-icon">◑</span>
      </div>
      <div class="loading-showcase">
        <div class="loading-state">
          <div class="spinner"></div>
          <p class="loading-text">Consulting the cosmos...</p>
        </div>
        <div class="loading-state">
          <div class="spinner-moon">☽</div>
          <p class="loading-text">Shuffling the deck...</p>
        </div>
        <div class="loading-state">
          <div class="loading-dots">
            <span class="dot"></span>
            <span class="dot"></span>
            <span class="dot"></span>
          </div>
          <p class="loading-text">Drawing your cards...</p>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.ui-kit {
  --kit-bg: #0a0a14;
  --kit-surface: #171724;
  --kit-surface-strong: #1d1d2e;
  --kit-elevated: #252540;
  --kit-border: rgba(245, 158, 11, 0.18);
  --kit-text: #f8f9fa;
  --kit-text-secondary: #b8b9c1;
  --kit-text-muted: #6b6b7e;
  --kit-purple: #8b5cf6;
  --kit-gold: #f59e0b;
  --kit-rose: #f3b8c6;
  --kit-success: #10b981;
  --kit-warning: #f59e0b;
  --kit-error: #ef4444;
  color: var(--kit-text);
  display: grid;
  gap: 28px;
}

.ui-kit__header,
.ui-kit__section {
  position: relative;
  overflow: hidden;
  padding: 34px;
  border: 1px solid rgba(245, 158, 11, 0.1);
  border-radius: 10px;
  background:
    radial-gradient(circle at top, rgba(139, 92, 246, 0.1), transparent 24%),
    linear-gradient(180deg, rgba(14, 14, 24, 0.98) 0%, rgba(10, 10, 20, 0.99) 100%);
  box-shadow:
    0 18px 48px rgba(5, 5, 10, 0.24),
    inset 0 1px 0 rgba(255, 255, 255, 0.02);
}

.celestial-border {
  padding: 12px 0 20px;
  text-align: center;
}

.ui-kit__title,
.section-title,
.heading-1,
.heading-2,
.heading-3,
.heading-4,
.card-title,
.modal-title {
  margin: 0;
  font-family: Georgia, "Times New Roman", serif;
  font-weight: 500;
  letter-spacing: 0.02em;
}

.ui-kit__title {
  font-size: clamp(2.2rem, 5vw, 4rem);
  color: var(--kit-text);
}

.moon-phases,
.stars,
.divider-icon,
.moon-symbol,
.card-icon {
  color: var(--kit-gold);
}

.ui-kit__subtitle,
.type-meta,
.card-meta,
.modal-meta,
.icon-label,
.loading-text,
.swatch-value,
.swatch-name,
.card-text,
.showcase-label,
.category-name {
  color: var(--kit-text-secondary);
}

.ui-kit__pill-row,
.button-row,
.badge-group,
.card-tags,
.modal-footer {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.kit-pill,
.tag,
.badge,
.role-badge,
.card-badge {
  display: inline-flex;
  align-items: center;
  border-radius: 3px;
  padding: 7px 11px;
  font-size: 0.85rem;
  line-height: 1;
}

.kit-pill,
.tag,
.badge-purple,
.role-seeker {
  border: 1px solid rgba(139, 92, 246, 0.26);
  background: rgba(139, 92, 246, 0.12);
  color: #d8c5ff;
}

.section-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 14px;
  margin-bottom: 20px;
}

.section-title {
  font-size: clamp(1.4rem, 3vw, 2rem);
  text-align: center;
}

.typography-grid,
.color-grid,
.cards-showcase,
.icons-grid,
.loading-showcase {
  display: grid;
  gap: 18px;
}

.typography-grid,
.color-grid,
.icons-grid {
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
}

.cards-showcase {
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
}

.type-example,
.color-swatch,
.card,
.icon-item,
.loading-state,
.form-group,
.modal,
.ornamental-divider {
  border: 1px solid rgba(245, 158, 11, 0.08);
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.02);
  padding: 18px;
}

.heading-1 {
  font-size: 3rem;
}

.heading-2 {
  font-size: 2.2rem;
}

.heading-3 {
  font-size: 1.75rem;
}

.heading-4 {
  font-size: 1.25rem;
}

.body-text,
.card-interpretation,
.modal-body p,
.checkbox-text,
.radio-text {
  margin: 0;
  color: var(--kit-text-secondary);
  line-height: 1.6;
}

.body-small,
.caption {
  margin: 0;
  color: var(--kit-text-muted);
}

.swatch {
  width: 100%;
  height: 92px;
  border-radius: 18px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  margin-bottom: 12px;
}

.primary-purple { background: #8b5cf6; }
.primary-gold { background: #f59e0b; }
.primary-rose { background: #f3b8c6; }
.bg-primary { background: #0a0a14; }
.bg-surface { background: #1a1a2e; }
.bg-elevated { background: #252540; }
.text-primary { background: #f8f9fa; }
.text-secondary { background: #b8b9c1; }
.text-muted { background: #6b6b7e; }
.semantic-success { background: #10b981; }
.semantic-warning { background: #f59e0b; }
.semantic-error { background: #ef4444; }

.button-showcase,
.form-showcase,
.badges-showcase,
.dividers-showcase {
  display: grid;
  gap: 16px;
}

.btn {
  position: relative;
  border: 1px solid transparent;
  border-radius: 4px;
  padding: 12px 18px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 160ms ease, box-shadow 160ms ease, opacity 160ms ease;
}

.btn:hover,
.btn.is-hover {
  transform: translateY(-1px);
}

.btn:focus-visible,
.btn.is-focus {
  outline: 2px solid rgba(245, 158, 11, 0.7);
  outline-offset: 2px;
  box-shadow: 0 0 0 4px rgba(245, 158, 11, 0.16);
}

.btn-primary {
  background: linear-gradient(135deg, var(--kit-gold), #f3c26d);
  color: #24150a;
}

.btn-secondary {
  background: rgba(139, 92, 246, 0.14);
  border-color: rgba(139, 92, 246, 0.28);
  color: #e9ddff;
}

.btn-ghost {
  background: transparent;
  border-color: rgba(255, 255, 255, 0.12);
  color: var(--kit-text);
}

.btn-small { padding: 9px 14px; font-size: 0.86rem; }
.btn-large { padding: 15px 24px; font-size: 1rem; }
.btn-full { width: 100%; justify-content: center; }

.btn[disabled],
.is-loading {
  opacity: 0.55;
  cursor: not-allowed;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.form-label {
  color: var(--kit-text);
  font-weight: 600;
}

.form-input,
.form-select {
  width: 100%;
  padding: 13px 15px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 4px;
  background: rgba(10, 10, 20, 0.88);
  color: var(--kit-text);
}

.form-input:focus,
.form-select:focus {
  outline: none;
  border-color: rgba(245, 158, 11, 0.58);
  box-shadow: 0 0 0 4px rgba(245, 158, 11, 0.12);
}

.input-error {
  border-color: rgba(239, 68, 68, 0.72);
}

.form-error {
  color: #ff9f9f;
  font-size: 0.9rem;
}

.form-hint {
  color: var(--kit-text-muted);
  font-size: 0.9rem;
}

.checkbox-label,
.radio-label {
  display: flex;
  align-items: center;
  gap: 10px;
}

.card {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-height: 100%;
}

.card-elevated {
  background: linear-gradient(180deg, rgba(37, 37, 64, 0.86), rgba(19, 19, 34, 0.92));
}

.card-outline {
  border-style: dashed;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.card-title,
.card-name,
.card-title-small {
  color: var(--kit-text);
}

.card-title-small,
.card-name {
  margin: 0;
  font-family: Georgia, "Times New Roman", serif;
}

.card-body {
  display: grid;
  gap: 14px;
}

.tarot-card-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-back,
.mini-card {
  border-radius: 6px;
  background:
    radial-gradient(circle at top, rgba(243, 184, 198, 0.3), transparent 32%),
    linear-gradient(180deg, #1d1030 0%, #12091f 100%);
  border: 1px solid rgba(243, 184, 198, 0.2);
}

.card-back {
  width: 150px;
  height: 240px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.moon-symbol {
  font-size: 2rem;
}

.spread-preview {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

.mini-card {
  height: 160px;
}

.card-footer {
  margin-top: auto;
}

.ornamental-divider {
  text-align: center;
  color: var(--kit-gold);
  font-size: 1.2rem;
}

.modal-showcase {
  display: grid;
}

.modal-preview {
  padding: 20px;
  border-radius: 6px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.02), rgba(0, 0, 0, 0.18));
}

.modal {
  max-width: 520px;
  margin: 0 auto;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.modal-close {
  border: 0;
  background: transparent;
  color: var(--kit-text-secondary);
  font-size: 1.1rem;
  cursor: pointer;
}

.badge-success { background: rgba(16, 185, 129, 0.14); color: #93f1ca; }
.badge-warning { background: rgba(245, 158, 11, 0.14); color: #ffd58a; }
.badge-error { background: rgba(239, 68, 68, 0.14); color: #ffafaf; }
.role-adept { background: rgba(243, 184, 198, 0.16); color: #ffd7e1; }
.role-master { background: rgba(245, 158, 11, 0.16); color: #ffd48b; }
.role-keeper { background: rgba(16, 185, 129, 0.16); color: #9ef2ce; }

.icons-grid {
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
}

.icon-item {
  display: grid;
  place-items: center;
  gap: 10px;
  text-align: center;
}

.icon-display {
  font-size: 2rem;
  color: var(--kit-gold);
}

.loading-showcase {
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
}

.loading-state {
  display: grid;
  place-items: center;
  gap: 12px;
  min-height: 160px;
}

.spinner {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  border: 3px solid rgba(255, 255, 255, 0.12);
  border-top-color: var(--kit-gold);
  animation: spin 1s linear infinite;
}

.spinner-moon {
  font-size: 2rem;
  color: var(--kit-gold);
  animation: float 1.8s ease-in-out infinite;
}

.loading-dots {
  display: flex;
  gap: 8px;
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: var(--kit-purple);
  animation: pulse 1.1s ease-in-out infinite;
}

.dot:nth-child(2) { animation-delay: 0.12s; }
.dot:nth-child(3) { animation-delay: 0.24s; }

@keyframes spin {
  to { transform: rotate(360deg); }
}

@keyframes float {
  50% { transform: translateY(-6px); }
}

@keyframes pulse {
  50% { transform: translateY(-4px); opacity: 0.45; }
}

@media (max-width: 720px) {
  .ui-kit__header,
  .ui-kit__section {
    padding: 22px 18px;
  }

  .ui-kit__title {
    font-size: 2rem;
  }

  .cards-showcase,
  .typography-grid,
  .color-grid,
  .icons-grid,
  .loading-showcase {
    grid-template-columns: 1fr;
  }

  .spread-preview {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}
</style>
