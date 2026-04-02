<template>
  <div class="ui-kit">
    <div class="ui-kit__header">
      <div class="celestial-border">
        <h1 class="ui-kit__title">
          <span class="moon-phases">☽◯☾</span>
          Nebula Tarot UI Kit
          <span class="stars">✦✧✦</span>
        </h1>
        <p class="ui-kit__subtitle">Mystical Modern Design System</p>
      </div>
    </div>

    <!-- Typography Section -->
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
            Body Text • New beginnings, innocence, spontaneity, and a free spirit. 
            The Fool represents the start of a journey with an open heart and mind.
          </p>
          <p class="type-meta">16px • Inter • 400</p>
        </div>
        
        <div class="type-example">
          <p class="body-small">
            Small Text • Keywords: freedom, adventure, new start, innocence
          </p>
          <p class="type-meta">14px • Inter • 400</p>
        </div>
        
        <div class="type-example">
          <p class="caption">Caption • Generated at 3:00 AM • Next reset in 21h</p>
          <p class="type-meta">12px • Inter • 400</p>
        </div>
      </div>
    </section>

    <!-- Color Palette Section -->
    <section class="ui-kit__section">
      <div class="section-divider">
        <span class="divider-icon">☽</span>
        <h2 class="section-title">Color Palette</h2>
        <span class="divider-icon">☾</span>
      </div>
      
      <div class="color-grid">
        <div class="color-category">
          <h4 class="category-name">Primary</h4>
          <div class="color-swatches">
            <div class="color-swatch">
              <div class="swatch primary-purple"></div>
              <span class="swatch-name">Mystic Purple</span>
              <span class="swatch-value">#8B5CF6</span>
            </div>
            <div class="color-swatch">
              <div class="swatch primary-gold"></div>
              <span class="swatch-name">Celestial Gold</span>
              <span class="swatch-value">#F59E0B</span>
            </div>
            <div class="color-swatch">
              <div class="swatch primary-rose"></div>
              <span class="swatch-name">Rose Gold</span>
              <span class="swatch-value">#F3B8C6</span>
            </div>
          </div>
        </div>

        <div class="color-category">
          <h4 class="category-name">Background</h4>
          <div class="color-swatches">
            <div class="color-swatch">
              <div class="swatch bg-primary"></div>
              <span class="swatch-name">Deep Void</span>
              <span class="swatch-value">#0A0A14</span>
            </div>
            <div class="color-swatch">
              <div class="swatch bg-surface"></div>
              <span class="swatch-name">Dark Surface</span>
              <span class="swatch-value">#1A1A2E</span>
            </div>
            <div class="color-swatch">
              <div class="swatch bg-elevated"></div>
              <span class="swatch-name">Elevated</span>
              <span class="swatch-value">#252540</span>
            </div>
          </div>
        </div>

        <div class="color-category">
          <h4 class="category-name">Text</h4>
          <div class="color-swatches">
            <div class="color-swatch">
              <div class="swatch text-primary"></div>
              <span class="swatch-name">Primary Text</span>
              <span class="swatch-value">#F8F9FA</span>
            </div>
            <div class="color-swatch">
              <div class="swatch text-secondary"></div>
              <span class="swatch-name">Secondary</span>
              <span class="swatch-value">#B8B9C1</span>
            </div>
            <div class="color-swatch">
              <div class="swatch text-muted"></div>
              <span class="swatch-name">Muted</span>
              <span class="swatch-value">#6B6B7E</span>
            </div>
          </div>
        </div>

        <div class="color-category">
          <h4 class="category-name">Semantic</h4>
          <div class="color-swatches">
            <div class="color-swatch">
              <div class="swatch semantic-success"></div>
              <span class="swatch-name">Success</span>
              <span class="swatch-value">#10B981</span>
            </div>
            <div class="color-swatch">
              <div class="swatch semantic-warning"></div>
              <span class="swatch-name">Warning</span>
              <span class="swatch-value">#F59E0B</span>
            </div>
            <div class="color-swatch">
              <div class="swatch semantic-error"></div>
              <span class="swatch-name">Error</span>
              <span class="swatch-value">#EF4444</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Buttons Section -->
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
          <h4 class="showcase-label">With Icons</h4>
          <button class="btn btn-primary">
            <span class="btn-icon">✦</span>
            Get Card of the Day
          </button>
          <button class="btn btn-secondary">
            <span class="btn-icon">☽</span>
            Night Mode
          </button>
        </div>
      </div>
    </section>

    <!-- Inputs Section -->
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

    <!-- Cards Section -->
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
              New beginnings, innocence, spontaneity, and a free spirit. 
              The Fool represents the start of a journey.
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

    <!-- Dividers Section -->
    <section class="ui-kit__section">
      <div class="section-divider">
        <span class="divider-icon">✦</span>
        <h2 class="section-title">Dividers & Ornaments</h2>
        <span class="divider-icon">✦</span>
      </div>
      
      <div class="dividers-showcase">
        <div class="ornamental-divider divider-1">
          <span class="divider-symbol">✦</span>
        </div>

        <div class="ornamental-divider divider-2">
          <span class="divider-left">◯</span>
          <span class="divider-center">☽◯☾</span>
          <span class="divider-right">◯</span>
        </div>

        <div class="ornamental-divider divider-3">
          <span class="moon-phases-full">☽ ◐ ○ ◑ ☾</span>
        </div>

        <div class="ornamental-divider divider-4">
          <span class="celestial-symbols">✦ ✧ ⋆ ✧ ✦</span>
        </div>

        <div class="ornamental-divider divider-5">
          <span class="arrow-left">→</span>
          <span class="center-burst">✦</span>
          <span class="arrow-right">←</span>
        </div>
      </div>
    </section>

    <!-- Modals Section -->
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

    <!-- Badges & Tags Section -->
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

        <div class="badge-group">
          <h4 class="showcase-label">Card Tags</h4>
          <span class="tag">Freedom</span>
          <span class="tag">Adventure</span>
          <span class="tag">New Beginning</span>
          <span class="tag">Intuition</span>
          <span class="tag">Transformation</span>
        </div>
      </div>
    </section>

    <!-- Icons & Symbols Section -->
    <section class="ui-kit__section">
      <div class="section-divider">
        <span class="divider-icon">☆</span>
        <h2 class="section-title">Mystical Icons</h2>
        <span class="divider-icon">★</span>
      </div>
      
      <div class="icons-grid">
        <div class="icon-item">
          <span class="icon-display">☽</span>
          <span class="icon-label">Crescent Moon</span>
        </div>
        <div class="icon-item">
          <span class="icon-display">◯</span>
          <span class="icon-label">Full Moon</span>
        </div>
        <div class="icon-item">
          <span class="icon-display">☾</span>
          <span class="icon-label">Waning Moon</span>
        </div>
        <div class="icon-item">
          <span class="icon-display">✦</span>
          <span class="icon-label">Star</span>
        </div>
        <div class="icon-item">
          <span class="icon-display">✧</span>
          <span class="icon-label">Small Star</span>
        </div>
        <div class="icon-item">
          <span class="icon-display">⋆</span>
          <span class="icon-label">Asterisk</span>
        </div>
        <div class="icon-item">
          <span class="icon-display">◐</span>
          <span class="icon-label">Half Moon</span>
        </div>
        <div class="icon-item">
          <span class="icon-display">◑</span>
          <span class="icon-label">Reverse Half</span>
        </div>
        <div class="icon-item">
          <span class="icon-display">☆</span>
          <span class="icon-label">Outlined Star</span>
        </div>
        <div class="icon-item">
          <span class="icon-display">★</span>
          <span class="icon-label">Filled Star</span>
        </div>
        <div class="icon-item">
          <span class="icon-display">◈</span>
          <span class="icon-label">Diamond</span>
        </div>
        <div class="icon-item">
          <span class="icon-display">✺</span>
          <span class="icon-label">Starburst</span>
        </div>
      </div>
    </section>

    <!-- Loading States Section -->
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

    <!-- Footer -->
    <footer class="ui-kit__footer">
      <div class="ornamental-divider divider-2">
        <span class="divider-left">◯</span>
        <span class="divider-center">☽◯☾</span>
        <span class="divider-right">◯</span>
      </div>
      <p class="footer-text">Nebula Tarot • Mystical Modern Design System</p>
      <p class="footer-meta">Dark Theme • v1.0</p>
    </footer>
  </div>
</template>

<script setup>
// UI Kit is static, no logic needed for now
</script>

<style lang="scss" scoped>
// CSS Variables (Dark Theme)
:root {
  // Colors - Primary
  --color-primary: #8B5CF6;
  --color-primary-hover: #7C3AED;
  --color-primary-active: #6D28D9;
  --color-gold: #F59E0B;
  --color-rose: #F3B8C6;
  
  // Background
  --color-bg-primary: #0A0A14;
  --color-bg-surface: #1A1A2E;
  --color-bg-elevated: #252540;
  --color-bg-overlay: rgba(10, 10, 20, 0.95);
  
  // Text
  --color-text-primary: #F8F9FA;
  --color-text-secondary: #B8B9C1;
  --color-text-muted: #6B6B7E;
  
  // Semantic
  --color-success: #10B981;
  --color-warning: #F59E0B;
  --color-error: #EF4444;
  
  // Borders
  --color-border: rgba(139, 92, 246, 0.2);
  --color-border-muted: rgba(255, 255, 255, 0.1);
  
  // Typography
  --font-heading: 'Cinzel', serif;
  --font-body: 'Inter', sans-serif;
  
  // Shadows
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.5);
  --shadow-md: 0 4px 6px rgba(0, 0, 0, 0.4);
  --shadow-lg: 0 10px 15px rgba(0, 0, 0, 0.6);
  --shadow-glow: 0 0 20px rgba(139, 92, 246, 0.3);
  
  // Transitions
  --transition-fast: 0.15s ease;
  --transition-base: 0.3s ease;
}

// Base Styles
.ui-kit {
  min-height: 100vh;
  background: var(--color-bg-primary);
  color: var(--color-text-primary);
  padding: 4rem 2rem;
  font-family: var(--font-body);
  
  // Starry background
  background-image: 
    radial-gradient(2px 2px at 20% 30%, white, transparent),
    radial-gradient(2px 2px at 60% 70%, white, transparent),
    radial-gradient(1px 1px at 50% 50%, white, transparent),
    radial-gradient(1px 1px at 80% 10%, white, transparent),
    radial-gradient(2px 2px at 90% 60%, white, transparent),
    radial-gradient(1px 1px at 33% 80%, white, transparent);
  background-size: 200% 200%;
  background-position: 0% 0%;
}

// Header
.ui-kit__header {
  text-align: center;
  margin-bottom: 6rem;
}

.celestial-border {
  border: 2px solid var(--color-border);
  border-radius: 1rem;
  padding: 3rem 2rem;
  background: var(--color-bg-surface);
  box-shadow: var(--shadow-glow);
  position: relative;
  
  &::before,
  &::after {
    content: '✦';
    position: absolute;
    font-size: 1.5rem;
    color: var(--color-gold);
  }
  
  &::before {
    top: -1rem;
    left: 50%;
    transform: translateX(-50%);
  }
  
  &::after {
    bottom: -1rem;
    left: 50%;
    transform: translateX(-50%);
  }
}

.ui-kit__title {
  font-family: var(--font-heading);
  font-size: 3rem;
  font-weight: 500;
  margin: 0 0 1rem;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  background: linear-gradient(135deg, var(--color-primary), var(--color-gold));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  
  .moon-phases,
  .stars {
    display: inline-block;
    margin: 0 1rem;
    color: var(--color-rose);
    -webkit-text-fill-color: var(--color-rose);
  }
}

.ui-kit__subtitle {
  font-size: 1.125rem;
  color: var(--color-text-secondary);
  letter-spacing: 0.1em;
  text-transform: uppercase;
  font-weight: 300;
}

// Sections
.ui-kit__section {
  max-width: 1200px;
  margin: 0 auto 6rem;
}

.section-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1.5rem;
  margin-bottom: 3rem;
  
  &::before,
  &::after {
    content: '';
    flex: 1;
    height: 1px;
    background: linear-gradient(
      to right,
      transparent,
      var(--color-border),
      transparent
    );
  }
}

.section-title {
  font-family: var(--font-heading);
  font-size: 2rem;
  font-weight: 400;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  margin: 0;
}

.divider-icon {
  font-size: 1.5rem;
  color: var(--color-gold);
}

// Typography Examples
.typography-grid {
  display: grid;
  gap: 2rem;
}

.type-example {
  padding: 1.5rem;
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border-muted);
  border-radius: 0.5rem;
}

.heading-1 {
  font-family: var(--font-heading);
  font-size: 3rem;
  font-weight: 500;
  margin: 0;
  line-height: 1.2;
}

.heading-2 {
  font-family: var(--font-heading);
  font-size: 2.25rem;
  font-weight: 500;
  margin: 0;
  line-height: 1.2;
}

.heading-3 {
  font-family: var(--font-heading);
  font-size: 1.75rem;
  font-weight: 400;
  margin: 0;
  line-height: 1.3;
}

.heading-4 {
  font-family: var(--font-heading);
  font-size: 1.25rem;
  font-weight: 400;
  margin: 0;
  line-height: 1.4;
}

.body-text {
  font-size: 1rem;
  line-height: 1.6;
  margin: 0;
  color: var(--color-text-secondary);
}

.body-small {
  font-size: 0.875rem;
  line-height: 1.6;
  margin: 0;
  color: var(--color-text-secondary);
}

.caption {
  font-size: 0.75rem;
  line-height: 1.5;
  margin: 0;
  color: var(--color-text-muted);
}

.type-meta {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  margin-top: 0.5rem;
  font-family: 'Courier New', monospace;
}

// Color Grid
.color-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
}

.color-category {
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border-muted);
  border-radius: 0.5rem;
  padding: 1.5rem;
}

.category-name {
  font-family: var(--font-heading);
  font-size: 1rem;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  margin: 0 0 1rem;
  color: var(--color-gold);
}

.color-swatches {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.color-swatch {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.swatch {
  width: 3rem;
  height: 3rem;
  border-radius: 0.5rem;
  border: 2px solid var(--color-border-muted);
  flex-shrink: 0;
}

.swatch-name {
  font-size: 0.875rem;
  flex: 1;
}

.swatch-value {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  font-family: 'Courier New', monospace;
}

// Color Swatches
.primary-purple { background: #8B5CF6; }
.primary-gold { background: #F59E0B; }
.primary-rose { background: #F3B8C6; }
.bg-primary { background: #0A0A14; }
.bg-surface { background: #1A1A2E; }
.bg-elevated { background: #252540; }
.text-primary { background: #F8F9FA; }
.text-secondary { background: #B8B9C1; }
.text-muted { background: #6B6B7E; }
.semantic-success { background: #10B981; }
.semantic-warning { background: #F59E0B; }
.semantic-error { background: #EF4444; }

// Buttons
.button-showcase {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.button-row {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
  padding: 1.5rem;
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border-muted);
  border-radius: 0.5rem;
}

.showcase-label {
  font-family: var(--font-heading);
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--color-text-muted);
  min-width: 120px;
  margin: 0;
}

.btn {
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  border: none;
  cursor: pointer;
  transition: all var(--transition-base);
  font-family: var(--font-body);
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

.btn-primary {
  background: var(--color-primary);
  color: white;
  box-shadow: var(--shadow-md);
  
  &:hover:not(:disabled) {
    background: var(--color-primary-hover);
    box-shadow: var(--shadow-glow);
    transform: translateY(-2px);
  }
  
  &:active:not(:disabled) {
    transform: translateY(0);
  }
}

.btn-secondary {
  background: var(--color-bg-elevated);
  color: var(--color-text-primary);
  border: 1px solid var(--color-border);
  
  &:hover:not(:disabled) {
    background: var(--color-bg-surface);
    border-color: var(--color-primary);
  }
}

.btn-ghost {
  background: transparent;
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border-muted);
  
  &:hover:not(:disabled) {
    color: var(--color-text-primary);
    border-color: var(--color-border);
  }
}

.btn-small {
  padding: 0.5rem 1rem;
  font-size: 0.75rem;
}

.btn-large {
  padding: 1rem 2rem;
  font-size: 1rem;
}

.btn-full {
  width: 100%;
  justify-content: center;
}

.btn-icon {
  font-size: 1.125rem;
}

// Forms
.form-showcase {
  display: grid;
  gap: 2rem;
  max-width: 600px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.form-input,
.form-select {
  padding: 0.75rem 1rem;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border-muted);
  border-radius: 0.5rem;
  color: var(--color-text-primary);
  font-size: 1rem;
  font-family: var(--font-body);
  transition: all var(--transition-base);
  
  &:focus {
    outline: none;
    border-color: var(--color-primary);
    box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.1);
  }
  
  &::placeholder {
    color: var(--color-text-muted);
  }
}

.input-error {
  border-color: var(--color-error);
  
  &:focus {
    border-color: var(--color-error);
    box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
  }
}

.form-error {
  font-size: 0.75rem;
  color: var(--color-error);
}

.form-hint {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.checkbox-label,
.radio-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  
  &:hover {
    color: var(--color-text-primary);
  }
}

.form-checkbox,
.form-radio {
  width: 1.25rem;
  height: 1.25rem;
  cursor: pointer;
  accent-color: var(--color-primary);
}

// Cards
.cards-showcase {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
}

.card {
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border-muted);
  border-radius: 1rem;
  overflow: hidden;
  transition: all var(--transition-base);
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: var(--shadow-lg);
  }
}

.card-elevated {
  box-shadow: var(--shadow-md);
}

.card-outline {
  background: transparent;
  border: 2px solid var(--color-border);
}

.card-header {
  padding: 1.5rem;
  border-bottom: 1px solid var(--color-border-muted);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  font-family: var(--font-heading);
  font-size: 1.25rem;
  font-weight: 400;
  margin: 0;
}

.card-title-small {
  font-family: var(--font-heading);
  font-size: 1rem;
  font-weight: 400;
  margin: 0;
}

.card-icon {
  font-size: 1.5rem;
  color: var(--color-gold);
  margin-right: 0.5rem;
}

.card-badge {
  padding: 0.25rem 0.75rem;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: 999px;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text-secondary);
}

.badge-purple {
  background: rgba(139, 92, 246, 0.1);
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.card-body {
  padding: 1.5rem;
}

.tarot-card-placeholder {
  width: 100%;
  aspect-ratio: 2/3;
  background: var(--color-bg-elevated);
  border: 2px solid var(--color-border);
  border-radius: 0.5rem;
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-back {
  font-size: 2rem;
  color: var(--color-gold);
}

.moon-symbol {
  display: block;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.card-name {
  font-family: var(--font-heading);
  font-size: 1.5rem;
  margin: 0 0 0.5rem;
  text-align: center;
}

.card-interpretation {
  font-size: 0.875rem;
  line-height: 1.6;
  color: var(--color-text-secondary);
  margin: 0 0 1rem;
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag {
  padding: 0.25rem 0.75rem;
  background: rgba(139, 92, 246, 0.1);
  border: 1px solid var(--color-border);
  border-radius: 999px;
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.card-footer {
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--color-border-muted);
  text-align: center;
}

.card-meta {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.spread-preview {
  display: flex;
  justify-content: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.mini-card {
  width: 60px;
  height: 90px;
  background: var(--color-bg-elevated);
  border: 2px solid var(--color-border);
  border-radius: 0.25rem;
}

.card-text {
  font-size: 0.875rem;
  line-height: 1.6;
  color: var(--color-text-secondary);
  margin: 0;
}

// Dividers
.dividers-showcase {
  display: flex;
  flex-direction: column;
  gap: 3rem;
}

.ornamental-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  font-size: 1.5rem;
  color: var(--color-gold);
  padding: 1rem 0;
  
  &::before,
  &::after {
    content: '';
    flex: 1;
    height: 1px;
    background: linear-gradient(
      to right,
      transparent,
      var(--color-border),
      transparent
    );
  }
}

.divider-symbol {
  font-size: 1.5rem;
}

.divider-left,
.divider-right {
  font-size: 1rem;
}

.divider-center {
  font-size: 1.5rem;
  color: var(--color-rose);
}

.moon-phases-full {
  letter-spacing: 0.5rem;
}

.celestial-symbols {
  letter-spacing: 1rem;
}

// Modal
.modal-showcase {
  display: flex;
  justify-content: center;
  padding: 2rem;
  background: rgba(10, 10, 20, 0.5);
  border-radius: 1rem;
}

.modal-preview {
  max-width: 500px;
  width: 100%;
}

.modal {
  background: var(--color-bg-surface);
  border: 2px solid var(--color-border);
  border-radius: 1rem;
  box-shadow: var(--shadow-lg);
}

.modal-header {
  padding: 1.5rem;
  border-bottom: 1px solid var(--color-border-muted);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-title {
  font-family: var(--font-heading);
  font-size: 1.5rem;
  margin: 0;
}

.modal-close {
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0;
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color var(--transition-base);
  
  &:hover {
    color: var(--color-text-primary);
  }
}

.modal-body {
  padding: 1.5rem;
  
  p {
    margin: 0 0 0.5rem;
    color: var(--color-text-secondary);
  }
}

.modal-meta {
  font-size: 0.875rem;
  color: var(--color-text-muted);
}

.modal-footer {
  padding: 1.5rem;
  border-top: 1px solid var(--color-border-muted);
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

// Badges
.badges-showcase {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.badge-group {
  padding: 1.5rem;
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border-muted);
  border-radius: 0.5rem;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 1rem;
}

.badge {
  padding: 0.5rem 1rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.badge-success {
  background: rgba(16, 185, 129, 0.1);
  color: var(--color-success);
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.badge-warning {
  background: rgba(245, 158, 11, 0.1);
  color: var(--color-warning);
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.badge-error {
  background: rgba(239, 68, 68, 0.1);
  color: var(--color-error);
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.role-badge {
  padding: 0.5rem 1rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  border: 1px solid;
}

.role-seeker {
  background: rgba(107, 114, 128, 0.1);
  color: #9CA3AF;
  border-color: rgba(107, 114, 128, 0.3);
}

.role-adept {
  background: rgba(59, 130, 246, 0.1);
  color: #60A5FA;
  border-color: rgba(59, 130, 246, 0.3);
}

.role-master {
  background: rgba(139, 92, 246, 0.1);
  color: var(--color-primary);
  border-color: rgba(139, 92, 246, 0.3);
}

.role-keeper {
  background: rgba(245, 158, 11, 0.1);
  color: var(--color-gold);
  border-color: rgba(245, 158, 11, 0.3);
}

// Icons
.icons-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 2rem;
}

.icon-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 1.5rem;
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border-muted);
  border-radius: 0.5rem;
}

.icon-display {
  font-size: 2.5rem;
  color: var(--color-gold);
}

.icon-label {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  text-align: center;
}

// Loading States
.loading-showcase {
  display: flex;
  justify-content: space-around;
  align-items: center;
  gap: 2rem;
  flex-wrap: wrap;
  padding: 3rem;
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border-muted);
  border-radius: 1rem;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.spinner {
  width: 3rem;
  height: 3rem;
  border: 3px solid var(--color-border-muted);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.spinner-moon {
  font-size: 3rem;
  color: var(--color-gold);
  animation: rotate 2s linear infinite;
}

@keyframes rotate {
  to { transform: rotate(360deg); }
}

.loading-dots {
  display: flex;
  gap: 0.5rem;
}

.dot {
  width: 0.75rem;
  height: 0.75rem;
  background: var(--color-primary);
  border-radius: 50%;
  animation: pulse 1.5s ease-in-out infinite;
  
  &:nth-child(2) {
    animation-delay: 0.2s;
  }
  
  &:nth-child(3) {
    animation-delay: 0.4s;
  }
}

@keyframes pulse {
  0%, 100% { opacity: 0.3; transform: scale(0.8); }
  50% { opacity: 1; transform: scale(1.2); }
}

.loading-text {
  font-size: 0.875rem;
  color: var(--color-text-muted);
  margin: 0;
}

// Footer
.ui-kit__footer {
  text-align: center;
  padding-top: 4rem;
  margin-top: 6rem;
  border-top: 1px solid var(--color-border-muted);
}

.footer-text {
  font-family: var(--font-heading);
  font-size: 1.125rem;
  margin: 1rem 0 0.5rem;
  letter-spacing: 0.05em;
}

.footer-meta {
  font-size: 0.875rem;
  color: var(--color-text-muted);
  margin: 0;
}

// Responsive
@media (max-width: 768px) {
  .ui-kit {
    padding: 2rem 1rem;
  }
  
  .ui-kit__title {
    font-size: 2rem;
  }
  
  .section-title {
    font-size: 1.5rem;
  }
  
  .color-grid {
    grid-template-columns: 1fr;
  }
  
  .cards-showcase {
    grid-template-columns: 1fr;
  }
  
  .button-row {
    flex-direction: column;
    align-items: stretch;
  }
  
  .showcase-label {
    min-width: auto;
  }
  
  .icons-grid {
    grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  }
}
</style>
