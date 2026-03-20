<template>
  <main class="ui-kit-root" :class="{ 'light-mode': !isDark }">
    <div class="ui-kit-shell">
      <header class="kit-header">
        <div class="arcana-logo">
          <span class="logo-symbol">&#9790;</span>
          <span class="logo-wordmark">ARCANA</span>
        </div>
        <button class="theme-toggle" type="button" aria-label="Toggle theme" @click="toggleTheme">
          {{ isDark ? "\u2600" : "\u263D" }}
        </button>
      </header>

      <nav class="section-nav" aria-label="Section navigation">
        <button
          v-for="section in sections"
          :key="section.id"
          type="button"
          class="section-tab"
          :class="{ active: activeSection === section.id }"
          @click="scrollToSection(section.id)"
        >
          {{ section.label }}
        </button>
      </nav>

      <div class="sections">
        <section id="brand" :ref="(el) => setSectionRef('brand', el)" class="kit-section">
          <h2 class="section-title">Section 1 - Brand Identity</h2>
          <div class="brand-grid">
            <article class="brand-card">
              <p class="brand-card-label">Dark variant</p>
              <div class="logo-preview logo-preview-dark">
                <span class="logo-preview-symbol">&#9790;</span>
                <span class="logo-preview-wordmark">ARCANA</span>
              </div>
            </article>
            <article class="brand-card">
              <p class="brand-card-label">Light variant</p>
              <div class="logo-preview logo-preview-light">
                <span class="logo-preview-symbol">&#9790;</span>
                <span class="logo-preview-wordmark logo-preview-wordmark-light">ARCANA</span>
              </div>
            </article>
          </div>
          <div class="brand-voice">
            <h3 class="sub-title">Brand Voice</h3>
            <div class="pill-row">
              <span v-for="tone in brandVoice" :key="tone" class="voice-pill">{{ tone }}</span>
            </div>
          </div>
        </section>

        <section id="palette" :ref="(el) => setSectionRef('palette', el)" class="kit-section">
          <h2 class="section-title">Section 2 - Color Palette</h2>
          <div class="palette-block">
            <h3 class="sub-title">Dark Scheme</h3>
            <div class="swatch-row">
              <div v-for="swatch in darkScheme" :key="swatch.token" class="swatch-item">
                <div
                  class="swatch-color"
                  :class="{ 'swatch-deep': swatch.deepTone }"
                  :style="{ backgroundColor: `var(${swatch.variable})` }"
                />
                <p class="swatch-token">{{ swatch.token }}</p>
                <p class="swatch-hex">{{ swatch.hex }}</p>
              </div>
            </div>
          </div>
          <div class="palette-block">
            <h3 class="sub-title">Light Scheme</h3>
            <div class="swatch-row">
              <div v-for="swatch in lightScheme" :key="swatch.token" class="swatch-item">
                <div class="swatch-color" :style="{ backgroundColor: `var(${swatch.variable})` }" />
                <p class="swatch-token">{{ swatch.token }}</p>
                <p class="swatch-hex">{{ swatch.hex }}</p>
              </div>
            </div>
          </div>
        </section>

        <section id="typography" :ref="(el) => setSectionRef('typography', el)" class="kit-section">
          <h2 class="section-title">Section 3 - Typography</h2>
          <div class="type-table">
            <div v-for="row in typographyRows" :key="row.role" class="type-row">
              <p class="type-label">{{ row.role }}</p>
              <p class="type-specimen" :class="row.className">{{ row.text }}</p>
            </div>
          </div>
          <blockquote class="font-callout">
            <p><strong>Display / Headings / Pullquotes</strong> -> Cormorant Garamond</p>
            <p><strong>Body / UI / Captions</strong> -> DM Sans</p>
          </blockquote>
        </section>

        <section id="iconography" :ref="(el) => setSectionRef('iconography', el)" class="kit-section">
          <h2 class="section-title">Section 4 - Iconography &amp; Symbols</h2>
          <div class="icon-grid">
            <article v-for="item in iconography" :key="item.name" class="icon-card">
              <p class="icon-glyph">{{ item.glyph }}</p>
              <p class="icon-name">{{ item.name }}</p>
              <p class="icon-usage">{{ item.usage }}</p>
            </article>
          </div>
        </section>
        <section id="buttons" :ref="(el) => setSectionRef('buttons', el)" class="kit-section">
          <h2 class="section-title">Section 5 - Components: Buttons</h2>
          <div class="button-grid">
            <article v-for="variant in buttonVariants" :key="variant.name" class="button-card">
              <button type="button" class="arc-btn" :class="variant.className">{{ variant.label }}</button>
              <p class="variant-name">{{ variant.name }}</p>
            </article>
          </div>
          <div class="button-states">
            <p class="sub-title">Primary states: Default · Hover · Disabled</p>
            <div class="states-row">
              <button type="button" class="arc-btn btn-primary">&#10022; Reveal Cards</button>
              <button type="button" class="arc-btn btn-primary hover-demo">&#10022; Reveal Cards</button>
              <button type="button" class="arc-btn btn-primary" disabled>&#10022; Reveal Cards</button>
            </div>
          </div>
        </section>

        <section id="forms" :ref="(el) => setSectionRef('forms', el)" class="kit-section">
          <h2 class="section-title">Section 6 - Components: Form Controls</h2>
          <div class="form-stack">
            <div class="form-field">
              <label class="field-label" for="cards-input">Text Input</label>
              <input id="cards-input" class="arc-input" type="text" placeholder="Ask the cards..." />
            </div>
            <div class="form-field">
              <p class="field-label">Language Toggle</p>
              <div class="language-toggle" role="tablist" aria-label="Language">
                <button type="button" class="lang-segment" :class="{ active: activeLanguage === 'EN' }" @click="activeLanguage = 'EN'">EN</button>
                <button type="button" class="lang-segment" :class="{ active: activeLanguage === 'RU' }" @click="activeLanguage = 'RU'">RU</button>
              </div>
            </div>
            <div class="form-field">
              <p class="field-label">Theme Toggle</p>
              <button class="theme-toggle control-toggle" type="button" @click="toggleControlTheme">
                {{ isControlDark ? "\u2600" : "\u263D" }}
              </button>
            </div>
            <div class="form-field">
              <p class="field-label">Status Badge</p>
              <div class="status-badge"><span class="status-dot" /><span class="status-text">&#9679; PREMIUM</span></div>
            </div>
          </div>
        </section>

        <section id="cards" :ref="(el) => setSectionRef('cards', el)" class="kit-section">
          <h2 class="section-title">Section 7 - Components: Cards</h2>
          <div class="cards-grid">
            <article class="card-surface">
              <h3>Surface Card</h3>
              <p>The veil is thin tonight and intuition aligns with focus.</p>
              <p>Use this surface as a neutral container for rich content.</p>
            </article>
            <article class="card-surface">
              <h3>Reading Card</h3>
              <p class="reading-quote">"What the Moon conceals, the Emperor reveals."</p>
              <div class="tag-row">
                <span class="tag-pill">Past</span>
                <span class="tag-pill">Present</span>
                <span class="tag-pill">Future</span>
              </div>
            </article>
            <article class="card-surface tarot-showcase">
              <h3>Tarot Card States</h3>
              <div class="tarot-grid">
                <div class="tarot-card tarot-back">
                  <svg class="tarot-pattern" viewBox="0 0 140 240" aria-hidden="true">
                    <rect x="8" y="8" width="124" height="224" fill="none" stroke="rgba(176,143,255,0.35)" stroke-width="1" />
                    <rect x="16" y="16" width="108" height="208" fill="none" stroke="rgba(212,175,114,0.28)" stroke-width="0.5" />
                    <circle cx="70" cy="120" r="36" fill="none" stroke="rgba(212,175,114,0.45)" stroke-width="1" />
                    <circle cx="70" cy="120" r="24" fill="none" stroke="rgba(176,143,255,0.35)" stroke-width="1" />
                    <circle cx="70" cy="120" r="12" fill="none" stroke="rgba(212,175,114,0.45)" stroke-width="1" />
                    <path d="M70 84 L101 138 L39 138 Z" fill="none" stroke="rgba(212,175,114,0.45)" stroke-width="1" />
                    <path d="M70 156 L101 102 L39 102 Z" fill="none" stroke="rgba(176,143,255,0.35)" stroke-width="1" />
                    <line x1="70" y1="76" x2="70" y2="164" stroke="rgba(212,175,114,0.35)" stroke-width="1" />
                    <line x1="26" y1="120" x2="114" y2="120" stroke="rgba(176,143,255,0.35)" stroke-width="1" />
                    <circle cx="70" cy="120" r="2.5" fill="rgba(212,175,114,0.7)" />
                    <text x="18" y="28" fill="rgba(212,175,114,0.55)" font-size="10">&#10022;</text>
                    <text x="118" y="28" fill="rgba(212,175,114,0.55)" font-size="10">&#10022;</text>
                    <text x="18" y="222" fill="rgba(212,175,114,0.55)" font-size="10">&#10022;</text>
                    <text x="118" y="222" fill="rgba(212,175,114,0.55)" font-size="10">&#10022;</text>
                  </svg>
                </div>
                <div class="tarot-card tarot-front">
                  <p class="roman">XIX</p>
                  <p class="major-glyph">&#9737;</p>
                  <p class="major-name">The Sun</p>
                </div>
              </div>
            </article>
          </div>
        </section>

        <section id="navigation" :ref="(el) => setSectionRef('navigation', el)" class="kit-section">
          <h2 class="section-title">Section 8 - Components: Navigation</h2>
          <div class="top-nav-demo">
            <button v-for="item in navItems" :key="item.name" type="button" class="top-nav-item" :class="{ active: item.active }">
              <span class="item-glyph">{{ item.glyph }}</span>
              <span>{{ item.name }}</span>
            </button>
          </div>
        </section>

        <section id="dividers" :ref="(el) => setSectionRef('dividers', el)" class="kit-section">
          <h2 class="section-title">Section 9 - Dividers &amp; Ornaments</h2>
          <div class="divider-stack">
            <div class="divider-item"><p class="divider-label">Simple line</p><div class="divider-simple" /></div>
            <div class="divider-item"><p class="divider-label">Gold gradient line</p><div class="divider-gradient" /></div>
            <div class="divider-item">
              <p class="divider-label">Ornament divider</p>
              <div class="divider-ornament"><span class="line" /><span class="ornament-text">&#10022; &#9790; &#10022;</span><span class="line" /></div>
            </div>
            <div class="divider-item"><p class="divider-label">Symbol divider</p><div class="divider-symbol">&#8853;</div></div>
          </div>
        </section>

        <section id="motion" :ref="(el) => setSectionRef('motion', el)" class="kit-section">
          <h2 class="section-title">Section 10 - Motion &amp; Animation Principles</h2>
          <div class="motion-table-wrap">
            <table class="motion-table">
              <thead><tr><th>Name</th><th>CSS spec</th><th>Live demo</th></tr></thead>
              <tbody>
                <tr>
                  <td>Card flip</td>
                  <td><code>rotateY 0deg->90deg->0deg, 0.6s ease</code></td>
                  <td><div class="flip-demo"><div class="flip-inner">&#9790;</div></div></td>
                </tr>
                <tr>
                  <td>Hover lift</td>
                  <td><code>translateY 0->-6px, scale 1->1.02, 0.3s ease</code></td>
                  <td><div class="lift-demo">Hover me</div></td>
                </tr>
                <tr>
                  <td>Text reveal</td>
                  <td><code>opacity 0->1, translateY 12px->0, 0.5s ease</code></td>
                  <td>
                    <div class="text-reveal-demo">
                      <button type="button" class="arc-btn btn-secondary reveal-btn" @click="toggleRevealText">Toggle reveal</button>
                      <p class="reveal-text" :class="{ visible: isRevealVisible }">Clarity arrives where attention settles.</p>
                    </div>
                  </td>
                </tr>
                <tr>
                  <td>Status pulse</td>
                  <td><code>opacity 1->0.4->1, 2s infinite</code></td>
                  <td><div class="status-badge"><span class="status-dot" /><span class="status-text">&#9679; LIVE</span></div></td>
                </tr>
                <tr>
                  <td>Button press</td>
                  <td><code>scale 1->0.97, 0.15s</code></td>
                  <td>
                    <button
                      type="button"
                      class="arc-btn btn-primary press-demo"
                      :class="{ pressed: isPressing }"
                      @mousedown="startPress"
                      @mouseup="endPress"
                      @mouseleave="endPress"
                      @touchstart="startPress"
                      @touchend="endPress"
                    >
                      Press
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        <section id="rules" :ref="(el) => setSectionRef('rules', el)" class="kit-section">
          <h2 class="section-title">Section 11 - Do's and Don'ts</h2>
          <div class="rules-grid">
            <article class="rules-column do-column">
              <h3 class="rules-title">Do</h3>
              <div class="rule-list"><p v-for="rule in doRules" :key="rule" class="rule-card">{{ rule }}</p></div>
            </article>
            <article class="rules-column dont-column">
              <h3 class="rules-title">Don't</h3>
              <div class="rule-list"><p v-for="rule in dontRules" :key="rule" class="rule-card">{{ rule }}</p></div>
            </article>
          </div>
        </section>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, ref } from "vue";

interface SectionLink {
  id: string;
  label: string;
}

interface ColorSwatch {
  token: string;
  variable: string;
  hex: string;
  deepTone?: boolean;
}

interface TypographyRow {
  role: string;
  className: string;
  text: string;
}

interface IconItem {
  glyph: string;
  name: string;
  usage: string;
}

interface ButtonVariant {
  name: string;
  className: string;
  label: string;
}

interface NavItem {
  glyph: string;
  name: string;
  active?: boolean;
}
const sections: SectionLink[] = [
  { id: "brand", label: "Brand" },
  { id: "palette", label: "Palette" },
  { id: "typography", label: "Typography" },
  { id: "iconography", label: "Iconography" },
  { id: "buttons", label: "Buttons" },
  { id: "forms", label: "Forms" },
  { id: "cards", label: "Cards" },
  { id: "navigation", label: "Navigation" },
  { id: "dividers", label: "Dividers" },
  { id: "motion", label: "Motion" },
  { id: "rules", label: "Do / Don't" }
];

const brandVoice = ["Mystical", "Precise", "Trustworthy", "Intimate"];

const darkScheme: ColorSwatch[] = [
  { token: "--bg", variable: "--dark-bg", hex: "#0d0d12", deepTone: true },
  { token: "--surface", variable: "--dark-surface", hex: "#16161f", deepTone: true },
  { token: "--surface2", variable: "--dark-surface2", hex: "#1e1e2e", deepTone: true },
  { token: "--surface3", variable: "--dark-surface3", hex: "#2a2a3e", deepTone: true },
  { token: "--accent", variable: "--dark-accent", hex: "#d4af72" },
  { token: "--accent2", variable: "--dark-accent2", hex: "#b08fff" },
  { token: "--accent3", variable: "--dark-accent3", hex: "#f0ddb8" },
  { token: "--text-primary", variable: "--dark-text-primary", hex: "#f5f2ed" },
  { token: "--text-secondary", variable: "--dark-text-secondary", hex: "#c8c2d4" },
  { token: "--text-muted", variable: "--dark-text-muted", hex: "#7a748a" }
];

const lightScheme: ColorSwatch[] = [
  { token: "--bg", variable: "--light-bg", hex: "#f2ede4" },
  { token: "--surface", variable: "--light-surface", hex: "#faf6ee" },
  { token: "--surface2", variable: "--light-surface2", hex: "#ede7da" },
  { token: "--surface3", variable: "--light-surface3", hex: "#e0d8c8" },
  { token: "--accent", variable: "--light-accent", hex: "#7a5210" },
  { token: "--accent2", variable: "--light-accent2", hex: "#5530b8" },
  { token: "--accent3", variable: "--light-accent3", hex: "#4a3010" },
  { token: "--text-primary", variable: "--light-text-primary", hex: "#14101c" },
  { token: "--text-secondary", variable: "--light-text-secondary", hex: "#3a3248" },
  { token: "--text-muted", variable: "--light-text-muted", hex: "#7a6e8a" }
];

const typographyRows: TypographyRow[] = [
  { role: "Display", className: "spec-display", text: "The World" },
  { role: "Heading 1", className: "spec-h1", text: "Three Card Reading" },
  { role: "Heading 2", className: "spec-h2", text: "The Narrative of Your Cards" },
  { role: "Subheading", className: "spec-sub", text: "Past · Present · Future" },
  {
    role: "Body",
    className: "spec-body",
    text: "A cycle closes as wisdom crystallizes - the turbulence of uncertainty gives way to clarity."
  },
  { role: "Caption", className: "spec-caption", text: "\u2726 READING ACTIVE" },
  { role: "Pullquote", className: "spec-quote", text: "\"What the Moon conceals, the Emperor reveals.\"" }
];

const iconography: IconItem[] = [
  { glyph: "\u263D", name: "Moon", usage: "Logo, nav markers, lunar themes" },
  { glyph: "\u2726", name: "Star", usage: "CTA buttons, section accents" },
  { glyph: "\u2295", name: "Circle Cross", usage: "Section dividers, card labels" },
  { glyph: "\u25B3", name: "Triangle", usage: "Fire element, upward motion" },
  { glyph: "\u25C7", name: "Diamond", usage: "Air element, clarity" },
  { glyph: "\u2643", name: "Jupiter", usage: "Planet markers, expansion" },
  { glyph: "\u2727", name: "Sparkle", usage: "Hover states, decorative accents" },
  { glyph: "\u263F", name: "Mercury", usage: "Communication, swift themes" }
];

const buttonVariants: ButtonVariant[] = [
  { name: "Primary", className: "btn-primary", label: "\u2726 Reveal Cards" },
  { name: "Secondary", className: "btn-secondary", label: "\u263D New Reading" },
  { name: "Ghost", className: "btn-ghost", label: "Save to Journal" },
  { name: "Violet Accent", className: "btn-violet", label: "\u2295 Celtic Cross" }
];

const navItems: NavItem[] = [
  { glyph: "\u2726", name: "Today" },
  { glyph: "\u263D", name: "Three Card", active: true },
  { glyph: "\u25C7", name: "Celtic Cross" },
  { glyph: "\u2643", name: "Year Ahead" },
  { glyph: "\u25B3", name: "Relationship" },
  { glyph: "\u2295", name: "Journal" }
];

const doRules = [
  "Use Cormorant Garamond for all display and heading text",
  "Use --accent gold sparingly as the primary highlight color",
  "Keep symbol glyphs small and precise (9-16px) as accents only",
  "Maintain 1px solid --border2 on all interactive containers"
];

const dontRules = [
  "Use system fonts (Arial, Inter, Roboto) anywhere in the UI",
  "Use more than 2 accent colors simultaneously in one component",
  "Use filled backgrounds on symbol/glyph elements - always transparent",
  "Use border-radius larger than --radius-card except on pill buttons"
];

const isDark = ref(true);
const activeSection = ref("brand");
const activeLanguage = ref<"EN" | "RU">("EN");
const isControlDark = ref(true);
const isRevealVisible = ref(false);
const isPressing = ref(false);

const sectionElements = new Map<string, HTMLElement>();
let sectionObserver: IntersectionObserver | null = null;

function toggleTheme() {
  isDark.value = !isDark.value;
}

function toggleControlTheme() {
  isControlDark.value = !isControlDark.value;
}

function toggleRevealText() {
  isRevealVisible.value = !isRevealVisible.value;
}

function startPress() {
  isPressing.value = true;
}

function endPress() {
  isPressing.value = false;
}

function setSectionRef(id: string, el: unknown) {
  if (el instanceof HTMLElement) {
    sectionElements.set(id, el);
    return;
  }

  sectionElements.delete(id);
}

function scrollToSection(id: string) {
  const element = sectionElements.get(id);
  if (!element) return;

  const y = element.getBoundingClientRect().top + window.scrollY - 160;
  window.scrollTo({ top: y, behavior: "smooth" });
  activeSection.value = id;
}

onMounted(async () => {
  await nextTick();

  sectionObserver = new IntersectionObserver(
    (entries) => {
      const visible = entries.filter((entry) => entry.isIntersecting).sort((a, b) => b.intersectionRatio - a.intersectionRatio);
      const selected = visible[0]?.target.getAttribute("id");
      if (selected) activeSection.value = selected;
    },
    {
      threshold: [0.2, 0.35, 0.5, 0.75],
      rootMargin: "-18% 0px -55% 0px"
    }
  );

  sectionElements.forEach((element) => {
    sectionObserver?.observe(element);
  });
});

onUnmounted(() => {
  sectionObserver?.disconnect();
});
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Cormorant+Garamond:ital,wght@0,300;0,400;0,500;0,600;1,300;1,400;1,500&family=DM+Sans:wght@300;400;500&display=swap");

.ui-kit-root {
  --dark-bg: #0d0d12;
  --dark-surface: #16161f;
  --dark-surface2: #1e1e2e;
  --dark-surface3: #2a2a3e;
  --dark-border: rgba(200, 180, 255, 0.22);
  --dark-border2: rgba(200, 180, 255, 0.38);
  --dark-border-gold: rgba(201, 169, 110, 0.35);
  --dark-accent: #d4af72;
  --dark-accent2: #b08fff;
  --dark-accent3: #f0ddb8;
  --dark-text-primary: #f5f2ed;
  --dark-text-secondary: #c8c2d4;
  --dark-text-muted: #7a748a;
  --dark-card-bg: #13131e;

  --light-bg: #f2ede4;
  --light-surface: #faf6ee;
  --light-surface2: #ede7da;
  --light-surface3: #e0d8c8;
  --light-border: rgba(80, 50, 140, 0.2);
  --light-border2: rgba(80, 50, 140, 0.35);
  --light-border-gold: rgba(120, 80, 20, 0.3);
  --light-accent: #7a5210;
  --light-accent2: #5530b8;
  --light-accent3: #4a3010;
  --light-text-primary: #14101c;
  --light-text-secondary: #3a3248;
  --light-text-muted: #7a6e8a;
  --light-card-bg: #e8e2d4;

  --bg: var(--dark-bg);
  --surface: var(--dark-surface);
  --surface2: var(--dark-surface2);
  --surface3: var(--dark-surface3);
  --border: var(--dark-border);
  --border2: var(--dark-border2);
  --border-gold: var(--dark-border-gold);
  --accent: var(--dark-accent);
  --accent2: var(--dark-accent2);
  --accent3: var(--dark-accent3);
  --text-primary: var(--dark-text-primary);
  --text-secondary: var(--dark-text-secondary);
  --text-muted: var(--dark-text-muted);
  --card-bg: var(--dark-card-bg);

  --font-display: "Cormorant Garamond", serif;
  --font-body: "DM Sans", sans-serif;

  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-card: 16px;
  --radius-pill: 40px;

  --ink-strong: #0d0b08;
  --white: #ffffff;
  --do-color: #4a9e6f;
  --dont-color: #c0392b;
  --swatch-contrast: rgba(255, 255, 255, 0.1);
  --shadow-card: 0 20px 34px rgba(0, 0, 0, 0.24);
  --shadow-lift: 0 16px 26px rgba(0, 0, 0, 0.2);
  --focus-ring: rgba(212, 175, 114, 0.15);

  min-height: 100vh;
  background: radial-gradient(circle at 18% -6%, var(--border-gold), transparent 35%), radial-gradient(circle at 92% 6%, var(--border), transparent 25%), var(--bg);
  color: var(--text-primary);
  font-family: var(--font-body);
  transition: background 0.4s, color 0.4s, border-color 0.3s;
}
.ui-kit-root.light-mode {
  --bg: var(--light-bg);
  --surface: var(--light-surface);
  --surface2: var(--light-surface2);
  --surface3: var(--light-surface3);
  --border: var(--light-border);
  --border2: var(--light-border2);
  --border-gold: var(--light-border-gold);
  --accent: var(--light-accent);
  --accent2: var(--light-accent2);
  --accent3: var(--light-accent3);
  --text-primary: var(--light-text-primary);
  --text-secondary: var(--light-text-secondary);
  --text-muted: var(--light-text-muted);
  --card-bg: var(--light-card-bg);
  --swatch-contrast: rgba(0, 0, 0, 0.1);
  --shadow-card: 0 18px 30px rgba(58, 50, 72, 0.16);
  --shadow-lift: 0 14px 24px rgba(58, 50, 72, 0.14);
}

.ui-kit-root *,
.ui-kit-root *::before,
.ui-kit-root *::after {
  box-sizing: border-box;
}

.ui-kit-shell {
  max-width: 1100px;
  margin: 0 auto;
  padding: clamp(12px, 2.4vw, 16px) clamp(14px, 4.2vw, 48px) clamp(40px, 6vw, 64px);
  display: grid;
  gap: clamp(30px, 4vw, 48px);
}

.kit-header,
.section-nav,
.kit-section,
.arc-btn,
.theme-toggle,
.section-tab,
.arc-input,
.lang-segment,
.top-nav-item {
  transition: background 0.4s, color 0.4s, border-color 0.3s;
}

.kit-header {
  position: sticky;
  top: 8px;
  z-index: 60;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  padding: 12px 16px;
  background: var(--surface);
  border: 1px solid var(--border2);
  border-radius: var(--radius-pill);
  box-shadow: var(--shadow-lift);
}

.arcana-logo {
  display: inline-flex;
  align-items: center;
  gap: 9px;
}

.logo-symbol,
.logo-preview-symbol {
  font-family: var(--font-display);
  font-size: clamp(18px, 2.5vw, 22px);
  color: var(--accent2);
}

.logo-wordmark,
.logo-preview-wordmark {
  font-family: var(--font-display);
  font-size: clamp(21px, 3vw, 28px);
  font-weight: 500;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--accent);
}

.logo-preview-wordmark-light {
  color: var(--light-accent);
}

.theme-toggle {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  border: 1px solid var(--border2);
  background: transparent;
  color: var(--text-secondary);
  font-size: 15px;
  line-height: 1;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.theme-toggle:hover {
  border-color: var(--accent);
  color: var(--accent);
}

.section-nav {
  position: sticky;
  top: 74px;
  z-index: 55;
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding: 10px;
  background: var(--surface);
  border: 1px solid var(--border2);
  border-radius: var(--radius-pill);
}

.section-tab {
  border: 1px solid var(--border2);
  border-radius: var(--radius-pill);
  background: transparent;
  color: var(--text-secondary);
  padding: 7px clamp(12px, 2vw, 16px);
  font-size: clamp(10px, 1.6vw, 11px);
  font-weight: 500;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  white-space: nowrap;
  cursor: pointer;
}

.section-tab.active,
.section-tab:hover {
  color: var(--accent);
  border-color: var(--accent);
}

.section-tab.active {
  background: var(--surface3);
}

.sections {
  display: grid;
  gap: clamp(28px, 3.8vw, 48px);
}

.kit-section {
  display: grid;
  gap: clamp(16px, 2.5vw, 22px);
  background: var(--surface);
  border: 1px solid var(--border2);
  border-radius: var(--radius-card);
  padding: clamp(16px, 2.4vw, 26px);
  box-shadow: var(--shadow-card);
}

.section-title {
  margin: 0;
  font-family: var(--font-display);
  font-size: clamp(24px, 3.1vw, 31px);
  font-weight: 400;
  letter-spacing: 0.05em;
}

.sub-title,
.field-label,
.variant-name,
.divider-label {
  margin: 0;
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--text-muted);
}

.brand-grid,
.rules-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.brand-card,
.button-card,
.rule-card,
.icon-card,
.card-surface {
  background: var(--surface2);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
}

.brand-card {
  display: grid;
  gap: 10px;
  padding: 16px;
}

.logo-preview {
  padding: 24px 16px;
  border-radius: var(--radius-md);
  border: 1px solid var(--border2);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.logo-preview-dark {
  background: var(--dark-bg);
}

.logo-preview-light {
  background: var(--light-bg);
}

.pill-row,
.tag-row,
.states-row {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.voice-pill,
.tag-pill {
  border: 1px solid var(--border2);
  border-radius: var(--radius-pill);
  padding: 6px 14px;
  font-size: 10px;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--text-secondary);
}

.palette-block,
.button-states,
.form-stack,
.form-field,
.divider-stack,
.divider-item {
  display: grid;
  gap: 10px;
}

.swatch-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(90px, 1fr));
  gap: 14px;
}

.swatch-item {
  display: grid;
  justify-items: center;
  gap: 6px;
}

.swatch-color {
  width: 80px;
  height: 64px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
}

.swatch-deep {
  border-color: var(--swatch-contrast);
}

.swatch-token {
  margin: 0;
  font-size: 10px;
  color: var(--text-secondary);
}

.swatch-hex {
  margin: 0;
  font-size: 10px;
  color: var(--text-muted);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Courier New", monospace;
}

.type-table {
  border-top: 1px solid var(--border);
}

.type-row {
  display: flex;
  gap: 18px;
  align-items: flex-start;
  padding: 14px 0;
  border-bottom: 1px solid var(--border);
}

.type-label {
  margin: 0;
  min-width: 90px;
  font-size: 10px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--text-muted);
  padding-top: 6px;
}

.type-specimen {
  margin: 0;
  flex: 1;
}
.spec-display { font-family: var(--font-display); font-size: clamp(34px, 7vw, 52px); font-weight: 300; letter-spacing: 0.06em; }
.spec-h1 { font-family: var(--font-display); font-size: clamp(28px, 5.4vw, 36px); font-weight: 400; letter-spacing: 0.05em; }
.spec-h2 { font-family: var(--font-display); font-size: clamp(23px, 4.8vw, 28px); font-weight: 400; letter-spacing: 0.04em; }
.spec-sub { font-family: var(--font-display); font-size: clamp(18px, 3.6vw, 20px); font-weight: 300; font-style: italic; color: var(--text-secondary); }
.spec-body { font-size: clamp(14px, 2.8vw, 15px); font-weight: 300; color: var(--text-secondary); line-height: 1.65; }
.spec-caption { font-size: clamp(10px, 2.3vw, 11px); font-weight: 500; letter-spacing: 0.15em; text-transform: uppercase; color: var(--text-muted); }
.spec-quote { font-family: var(--font-display); font-size: clamp(17px, 3.5vw, 19px); font-weight: 300; font-style: italic; color: var(--accent3); }

.font-callout {
  margin: 0;
  border-left: 2px solid var(--accent);
  padding: 16px 24px;
  background: var(--surface2);
  display: grid;
  gap: 8px;
}

.font-callout p,
.icon-usage,
.card-surface p,
.rule-card,
.motion-table td {
  margin: 0;
  color: var(--text-secondary);
}

.icon-grid,
.button-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.icon-card,
.button-card {
  padding: 14px;
  text-align: center;
  display: grid;
  gap: 8px;
  justify-items: center;
}

.icon-glyph { margin: 0; font-size: clamp(26px, 4.6vw, 32px); font-family: var(--font-display); color: var(--accent); }
.icon-name { margin: 0; font-size: 10px; text-transform: uppercase; letter-spacing: 0.12em; color: var(--text-muted); }

.arc-btn {
  border: 1px solid transparent;
  border-radius: var(--radius-pill);
  padding: 13px 36px;
  min-height: 44px;
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  cursor: pointer;
}

.arc-btn:disabled { opacity: 0.4; cursor: not-allowed; }
.btn-primary { background: var(--accent); color: var(--ink-strong); }
.btn-primary:hover:not(:disabled), .btn-primary.hover-demo { opacity: 0.88; letter-spacing: 0.26em; }
.btn-secondary { background: transparent; border: 1.5px solid var(--accent); color: var(--accent); }
.btn-secondary:hover { background: var(--focus-ring); }
.btn-ghost { background: transparent; border: 1px solid var(--border2); color: var(--text-secondary); }
.btn-ghost:hover { background: var(--surface2); }
.btn-violet { background: var(--accent2); color: var(--white); }
.btn-violet:hover { opacity: 0.88; }

.arc-input {
  width: min(320px, 100%);
  height: 44px;
  border-radius: var(--radius-md);
  border: 1px solid var(--border2);
  padding: 0 16px;
  background: var(--surface2);
  color: var(--text-primary);
}

.arc-input::placeholder { color: var(--text-muted); }
.arc-input:focus { border-color: var(--accent); outline: none; box-shadow: 0 0 0 3px var(--focus-ring); }

.language-toggle {
  display: inline-flex;
  border: 1px solid var(--border2);
  border-radius: 20px;
  overflow: hidden;
}

.lang-segment {
  border: none;
  background: transparent;
  color: var(--text-muted);
  padding: 6px 14px;
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.06em;
  cursor: pointer;
}

.lang-segment.active { background: var(--surface3); color: var(--accent); }
.control-toggle { width: 34px; height: 34px; }

.status-badge { display: inline-flex; align-items: center; gap: 6px; }
.status-dot { width: 5px; height: 5px; border-radius: 50%; background: var(--accent); animation: pulse 2s infinite; }
.status-text { font-size: 10px; font-weight: 500; letter-spacing: 0.08em; text-transform: uppercase; color: var(--accent); }

.cards-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
}

.card-surface {
  padding: 28px;
  display: grid;
  gap: 10px;
}

.card-surface h3,
.rules-title {
  margin: 0;
  font-family: var(--font-display);
  font-size: 20px;
  font-weight: 400;
}

.reading-quote {
  border-left: 3px solid var(--accent);
  padding-left: 12px;
  font-family: var(--font-display);
  font-style: italic;
  color: var(--accent3);
}

.tarot-grid { display: flex; gap: 14px; flex-wrap: wrap; justify-content: center; }
.tarot-card { width: 140px; height: 240px; border-radius: 12px; border: 1px solid var(--border2); background: var(--card-bg); overflow: hidden; }
.tarot-pattern { width: 100%; height: 100%; }
.tarot-front { display: grid; justify-items: center; align-content: space-between; padding: 12px 0 18px; }
.roman { margin: 0; font-size: 11px; letter-spacing: 0.12em; color: var(--text-secondary); }
.major-glyph { margin: 0; font-family: var(--font-display); font-size: 42px; color: var(--accent); }
.major-name { margin: 0; font-family: var(--font-display); font-size: 13px; font-weight: 500; color: var(--accent3); }

.top-nav-demo {
  height: 48px;
  display: flex;
  align-items: center;
  gap: 2px;
  overflow-x: auto;
  border: 1px solid var(--border2);
  background: var(--surface);
}

.top-nav-item {
  height: 100%;
  border: none;
  border-bottom: 2px solid transparent;
  background: transparent;
  color: var(--text-muted);
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  padding: 0 20px;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  white-space: nowrap;
  cursor: pointer;
}

.top-nav-item:hover { color: var(--text-secondary); }
.top-nav-item.active { color: var(--accent); border-bottom-color: var(--accent); }
.item-glyph { font-size: 9px; vertical-align: super; }

.divider-simple { height: 1px; background: var(--border2); }
.divider-gradient { height: 1px; background: linear-gradient(to right, transparent, var(--border-gold), transparent); }
.divider-ornament { display: flex; align-items: center; gap: 14px; }
.divider-ornament .line { flex: 1; height: 1px; background: var(--border-gold); }
.ornament-text { font-size: 10px; letter-spacing: 0.3em; color: var(--accent); }
.divider-symbol { text-align: center; font-size: 16px; color: var(--accent2); }

.motion-table-wrap { overflow-x: auto; }
.motion-table { width: 100%; border-collapse: collapse; min-width: 620px; }
.motion-table th, .motion-table td { border: 1px solid var(--border); padding: 12px; text-align: left; vertical-align: middle; }
.motion-table th { font-size: 11px; letter-spacing: 0.12em; text-transform: uppercase; color: var(--text-muted); background: var(--surface2); }
.motion-table code { color: var(--text-primary); font-size: 12px; }

.flip-demo { width: 60px; height: 90px; perspective: 800px; }
.flip-inner { width: 100%; height: 100%; border-radius: var(--radius-sm); border: 1px solid var(--border2); background: var(--surface2); display: grid; place-items: center; font-size: 24px; color: var(--accent); }
.flip-demo:hover .flip-inner { animation: cardFlip 0.6s ease; }

.lift-demo { display: inline-flex; align-items: center; justify-content: center; min-width: 88px; min-height: 58px; border-radius: var(--radius-md); border: 1px solid var(--border2); background: var(--surface2); color: var(--text-secondary); transition: transform 0.3s ease; }
.lift-demo:hover { transform: translateY(-6px) scale(1.02); }

.text-reveal-demo { display: grid; gap: 8px; }
.reveal-btn { padding: 9px 14px; }
.reveal-text { margin: 0; opacity: 0; transform: translateY(12px); transition: opacity 0.5s ease, transform 0.5s ease; }
.reveal-text.visible { opacity: 1; transform: translateY(0); }
.press-demo { transition: transform 0.15s; }
.press-demo.pressed { transform: scale(0.97); }

.rules-column { border-top: 3px solid; border-radius: var(--radius-md); padding-top: 12px; display: grid; gap: 12px; }
.do-column { border-top-color: var(--do-color); }
.dont-column { border-top-color: var(--dont-color); }
.rule-list { display: grid; gap: 10px; }
.rule-card { padding: 20px; line-height: 1.5; }

@keyframes pulse {
  0%,
  100% { opacity: 1; }
  50% { opacity: 0.4; }
}

@keyframes cardFlip {
  0% { transform: rotateY(0deg); }
  50% { transform: rotateY(90deg); }
  100% { transform: rotateY(0deg); }
}

@media (max-width: 1100px) {
  .ui-kit-shell { padding-left: 24px; padding-right: 24px; }
  .cards-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}

@media (max-width: 900px) {
  .button-grid,
  .icon-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}

@media (max-width: 760px) {
  .ui-kit-shell { gap: 28px; padding: 12px 14px 44px; }
  .kit-header { padding: 10px 12px; }
  .section-nav { top: 68px; padding: 8px; gap: 6px; }
  .section-tab { padding: 6px 12px; letter-spacing: 0.1em; }
  .kit-section { padding: 16px; }
  .section-title { font-size: 25px; }
  .brand-grid,
  .rules-grid,
  .cards-grid,
  .button-grid { grid-template-columns: 1fr; }
  .icon-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .type-row { flex-direction: column; gap: 6px; }
  .type-label { min-width: auto; padding-top: 0; }
  .font-callout { padding: 14px 16px; }
  .button-card,
  .icon-card,
  .card-surface,
  .rule-card { padding: 16px; }
  .arc-btn { width: 100%; padding: 11px 16px; letter-spacing: 0.14em; }
  .states-row { flex-direction: column; }
  .top-nav-demo { height: 44px; }
  .top-nav-item { padding: 0 14px; }
  .motion-table { min-width: 560px; }
  .motion-table th,
  .motion-table td { padding: 10px; }
  .tarot-card { width: 126px; height: 216px; }
  .major-glyph { font-size: 36px; }
}

@media (max-width: 480px) {
  .ui-kit-shell { gap: 24px; padding: 10px 10px 36px; }
  .kit-header { top: 6px; }
  .section-nav { top: 60px; }
  .section-tab { padding: 5px 10px; font-size: 9.5px; letter-spacing: 0.08em; }
  .kit-section { padding: 14px; gap: 14px; }
  .section-title { font-size: 22px; }
  .swatch-row { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .swatch-color { width: 100%; height: 58px; max-width: 90px; }
  .icon-grid { grid-template-columns: 1fr; }
  .type-row { padding: 10px 0; }
  .type-label { font-size: 9px; }
  .tarot-grid { gap: 10px; }
  .tarot-card { width: 112px; height: 192px; }
  .roman { font-size: 10px; }
  .major-glyph { font-size: 31px; }
  .major-name { font-size: 12px; }
  .motion-table { min-width: 500px; }
  .divider-ornament { gap: 8px; }
  .ornament-text { letter-spacing: 0.2em; }
  .rules-title { font-size: 18px; }
}
</style>

