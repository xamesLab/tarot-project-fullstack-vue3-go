<template>
  <div class="arcana-bg" :class="{ 'light-mode': isLight }">
    <div class="bg-blooms" aria-hidden="true"></div>

    <svg
      class="bg-svg"
      viewBox="0 0 1440 900"
      preserveAspectRatio="xMidYMid slice"
      xmlns="http://www.w3.org/2000/svg"
      aria-hidden="true"
      focusable="false"
    >
      <defs>
        <radialGradient id="d-top" cx="50%" cy="0%" r="75%" gradientUnits="objectBoundingBox">
          <stop offset="0%" stop-color="#3d2a80" stop-opacity="0.75" />
          <stop offset="100%" stop-color="#0d0d12" stop-opacity="0" />
        </radialGradient>
        <radialGradient id="d-br" cx="92%" cy="98%" r="65%" gradientUnits="objectBoundingBox">
          <stop offset="0%" stop-color="#6b2e08" stop-opacity="0.60" />
          <stop offset="100%" stop-color="#0d0d12" stop-opacity="0" />
        </radialGradient>
        <radialGradient id="d-bl" cx="0%" cy="65%" r="50%" gradientUnits="objectBoundingBox">
          <stop offset="0%" stop-color="#1e0e55" stop-opacity="0.55" />
          <stop offset="100%" stop-color="#0d0d12" stop-opacity="0" />
        </radialGradient>
        <radialGradient id="d-vig" cx="50%" cy="50%" r="72%" gradientUnits="objectBoundingBox">
          <stop offset="50%" stop-color="#0d0d12" stop-opacity="0" />
          <stop offset="100%" stop-color="#07070d" stop-opacity="0.72" />
        </radialGradient>

        <radialGradient id="l-top" cx="50%" cy="0%" r="75%" gradientUnits="objectBoundingBox">
          <stop offset="0%" stop-color="#c4b0f0" stop-opacity="0.55" />
          <stop offset="100%" stop-color="#f2ede4" stop-opacity="0" />
        </radialGradient>
        <radialGradient id="l-br" cx="92%" cy="98%" r="65%" gradientUnits="objectBoundingBox">
          <stop offset="0%" stop-color="#d4a84a" stop-opacity="0.45" />
          <stop offset="100%" stop-color="#f2ede4" stop-opacity="0" />
        </radialGradient>
        <radialGradient id="l-bl" cx="0%" cy="65%" r="50%" gradientUnits="objectBoundingBox">
          <stop offset="0%" stop-color="#a898d8" stop-opacity="0.38" />
          <stop offset="100%" stop-color="#f2ede4" stop-opacity="0" />
        </radialGradient>
        <radialGradient id="l-vig" cx="50%" cy="50%" r="72%" gradientUnits="objectBoundingBox">
          <stop offset="55%" stop-color="#f2ede4" stop-opacity="0" />
          <stop offset="100%" stop-color="#ddd5c4" stop-opacity="0.60" />
        </radialGradient>

        <radialGradient id="glow-gold" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stop-color="#d4af72" stop-opacity="1" />
          <stop offset="100%" stop-color="#d4af72" stop-opacity="0" />
        </radialGradient>
        <radialGradient id="glow-violet" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stop-color="#b08fff" stop-opacity="1" />
          <stop offset="100%" stop-color="#b08fff" stop-opacity="0" />
        </radialGradient>
        <radialGradient id="glow-gold-l" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stop-color="#7a5210" stop-opacity="0.8" />
          <stop offset="100%" stop-color="#7a5210" stop-opacity="0" />
        </radialGradient>
        <radialGradient id="glow-violet-l" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stop-color="#5530b8" stop-opacity="0.7" />
          <stop offset="100%" stop-color="#5530b8" stop-opacity="0" />
        </radialGradient>
      </defs>

      <rect class="bloom bloom-d" fill="url(#d-top)" width="1440" height="900" />
      <rect class="bloom bloom-d" fill="url(#d-br)" width="1440" height="900" />
      <rect class="bloom bloom-d" fill="url(#d-bl)" width="1440" height="900" />
      <rect class="bloom bloom-l" fill="url(#l-top)" width="1440" height="900" />
      <rect class="bloom bloom-l" fill="url(#l-br)" width="1440" height="900" />
      <rect class="bloom bloom-l" fill="url(#l-bl)" width="1440" height="900" />

      <g class="mandala" transform="translate(720, 450)" fill="none">
        <circle r="185" class="m-violet" stroke-width="0.6" />
        <circle r="148" class="m-gold" stroke-width="0.8" />
        <circle r="110" class="m-violet" stroke-width="0.6" />
        <circle r="72" class="m-gold" stroke-width="0.9" />
        <circle r="36" class="m-violet" stroke-width="0.8" />
        <circle r="12" class="m-gold" stroke-width="1" />

        <line x1="0" y1="-185" x2="0" y2="185" class="m-gold-line" stroke-width="0.4" />
        <line x1="-185" y1="0" x2="185" y2="0" class="m-gold-line" stroke-width="0.4" />
        <line x1="-131" y1="-131" x2="131" y2="131" class="m-violet-line" stroke-width="0.4" />
        <line x1="131" y1="-131" x2="-131" y2="131" class="m-violet-line" stroke-width="0.4" />

        <polygon points="0,-90 77.9,45 -77.9,45" class="m-gold" stroke-width="1" />
        <polygon points="0,90 77.9,-45 -77.9,-45" class="m-violet" stroke-width="1" />

        <g class="m-gold-line" stroke-width="0.6">
          <line x1="0" y1="-148" x2="0" y2="-138" />
          <line x1="0" y1="138" x2="0" y2="148" />
          <line x1="-148" y1="0" x2="-138" y2="0" />
          <line x1="138" y1="0" x2="148" y2="0" />
          <line x1="-104.6" y1="-104.6" x2="-97.6" y2="-97.6" />
          <line x1="97.6" y1="97.6" x2="104.6" y2="104.6" />
          <line x1="104.6" y1="-104.6" x2="97.6" y2="-97.6" />
          <line x1="-97.6" y1="97.6" x2="-104.6" y2="104.6" />
        </g>

        <circle r="4" class="dot-gold" />
        <circle r="2" class="dot-cream" />
      </g>

      <g class="constellation gold">
        <line x1="55" y1="110" x2="175" y2="72" />
        <line x1="175" y1="72" x2="295" y2="145" />
        <line x1="295" y1="145" x2="415" y2="92" />
        <line x1="415" y1="92" x2="478" y2="195" />
        <line x1="175" y1="72" x2="238" y2="192" />
        <line x1="238" y1="192" x2="295" y2="145" />

        <ellipse cx="175" cy="72" rx="10" ry="10" class="glow gold-glow" />
        <ellipse cx="415" cy="92" rx="8" ry="8" class="glow gold-glow-small" />

        <circle cx="55" cy="110" r="2.2" class="star gold s-a1" />
        <circle cx="175" cy="72" r="3" class="star gold s-a2" />
        <circle cx="295" cy="145" r="2.4" class="star gold s-a3" />
        <circle cx="415" cy="92" r="2.8" class="star gold s-a4" />
        <circle cx="478" cy="195" r="2" class="star gold s-a5" />
        <circle cx="238" cy="192" r="1.8" class="star gold s-a6" />
      </g>

      <g class="constellation violet">
        <line x1="72" y1="710" x2="195" y2="748" />
        <line x1="195" y1="748" x2="348" y2="688" />
        <line x1="348" y1="688" x2="458" y2="755" />
        <line x1="195" y1="748" x2="258" y2="828" />
        <line x1="348" y1="688" x2="420" y2="620" />

        <ellipse cx="195" cy="748" rx="10" ry="10" class="glow violet-glow" />

        <circle cx="72" cy="710" r="2.2" class="star violet s-b1" />
        <circle cx="195" cy="748" r="3" class="star violet s-b2" />
        <circle cx="348" cy="688" r="2" class="star violet s-b3" />
        <circle cx="458" cy="755" r="2.4" class="star violet s-b4" />
        <circle cx="258" cy="828" r="1.8" class="star violet s-b5" />
        <circle cx="420" cy="620" r="2" class="star violet s-b6" />
      </g>

      <g class="scattered-stars">
        <circle cx="380" cy="315" r="1.4" class="star gold s-m1" />
        <circle cx="128" cy="395" r="1.2" class="star gold s-m2" />
        <circle cx="28" cy="275" r="1.2" class="star gold s-m3" />
        <circle cx="488" cy="575" r="1.4" class="star gold s-m4" />
        <circle cx="62" cy="832" r="1.2" class="star gold s-m5" />
        <circle cx="408" cy="862" r="1.4" class="star gold s-m6" />
        <circle cx="168" cy="558" r="1.3" class="star violet s-m7" />
        <circle cx="318" cy="498" r="1.4" class="star violet s-m8" />
        <circle cx="455" cy="338" r="1.3" class="star violet s-m9" />
        <circle cx="92" cy="520" r="1.2" class="star violet s-m10" />
        <circle cx="1380" cy="200" r="1.5" class="star gold s-m11" />
        <circle cx="1200" cy="120" r="1.5" class="star violet s-m12" />
        <circle cx="1300" cy="680" r="1.5" class="star gold s-m13" />
        <circle cx="1100" cy="780" r="1.5" class="star violet s-m14" />
        <circle cx="1050" cy="400" r="1.4" class="star gold s-m15" />
        <circle cx="1380" cy="550" r="1.4" class="star violet s-m16" />
      </g>

      <g class="glyphs" font-family="Georgia, serif" text-anchor="middle">
        <text x="720" y="52" font-size="40" class="glyph-crescent">&#9790;</text>
        <text x="28" y="42" font-size="16" class="glyph-star">&#10022;</text>
        <text x="1412" y="42" font-size="16" class="glyph-star">&#10022;</text>
        <text x="28" y="876" font-size="16" class="glyph-star">&#10022;</text>
        <text x="1412" y="876" font-size="16" class="glyph-star">&#10022;</text>
        <line x1="80" y1="450" x2="680" y2="450" class="divider-line" />
        <line x1="760" y1="450" x2="1360" y2="450" class="divider-line" />
        <text x="720" y="455" font-size="14" class="glyph-center">&#8853;</text>
      </g>

      <rect x="18" y="18" width="1404" height="864" rx="4" fill="none" class="frame-border" />

      <rect class="vignette vignette-d" fill="url(#d-vig)" width="1440" height="900" />
      <rect class="vignette vignette-l" fill="url(#l-vig)" width="1440" height="900" />
    </svg>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  isLight: boolean;
}>();
</script>

<style scoped>
.arcana-bg {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  overflow: hidden;
}

.bg-blooms {
  position: absolute;
  inset: 0;
  background-color: #0d0d12;
  background-image:
    radial-gradient(ellipse 80% 60% at 50% -5%, rgba(61, 42, 128, 0.75) 0%, transparent 100%),
    radial-gradient(ellipse 65% 55% at 95% 105%, rgba(107, 46, 8, 0.6) 0%, transparent 100%),
    radial-gradient(ellipse 50% 45% at -5% 65%, rgba(30, 14, 85, 0.55) 0%, transparent 100%),
    radial-gradient(ellipse 40% 35% at 50% 50%, rgba(42, 18, 96, 0.3) 0%, transparent 100%);
}

.arcana-bg.light-mode .bg-blooms {
  background-color: #f2ede4;
  background-image:
    radial-gradient(ellipse 80% 60% at 50% -5%, rgba(196, 176, 240, 0.55) 0%, transparent 100%),
    radial-gradient(ellipse 65% 55% at 95% 105%, rgba(212, 168, 74, 0.45) 0%, transparent 100%),
    radial-gradient(ellipse 50% 45% at -5% 65%, rgba(168, 152, 216, 0.38) 0%, transparent 100%);
}

.bg-svg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

.arcana-bg .bloom-d,
.arcana-bg .vignette-d {
  display: block;
}

.arcana-bg .bloom-l,
.arcana-bg .vignette-l {
  display: none;
}

.arcana-bg.light-mode .bloom-d,
.arcana-bg.light-mode .vignette-d {
  display: none;
}

.arcana-bg.light-mode .bloom-l,
.arcana-bg.light-mode .vignette-l {
  display: block;
}

.arcana-bg .m-gold {
  stroke: #d4af72;
  opacity: 0.26;
}

.arcana-bg .m-violet {
  stroke: #b08fff;
  opacity: 0.22;
}

.arcana-bg .m-gold-line {
  stroke: #d4af72;
  opacity: 0.16;
}

.arcana-bg .m-violet-line {
  stroke: #b08fff;
  opacity: 0.16;
}

.arcana-bg .dot-gold {
  fill: #d4af72;
  opacity: 0.7;
}

.arcana-bg .dot-cream {
  fill: #f0ddb8;
  opacity: 0.9;
}

.arcana-bg.light-mode .m-gold {
  stroke: #7a5210;
  opacity: 0.22;
}

.arcana-bg.light-mode .m-violet {
  stroke: #5530b8;
  opacity: 0.18;
}

.arcana-bg.light-mode .m-gold-line {
  stroke: #7a5210;
  opacity: 0.14;
}

.arcana-bg.light-mode .m-violet-line {
  stroke: #5530b8;
  opacity: 0.14;
}

.arcana-bg.light-mode .dot-gold {
  fill: #7a5210;
  opacity: 0.55;
}

.arcana-bg.light-mode .dot-cream {
  fill: #5530b8;
  opacity: 0.7;
}

.arcana-bg .constellation {
  fill: none;
  stroke-width: 0.7;
}

.arcana-bg .constellation.gold {
  stroke: #d4af72;
  opacity: 0.22;
}

.arcana-bg .constellation.violet {
  stroke: #b08fff;
  opacity: 0.2;
}

.arcana-bg .star.gold {
  fill: #d4af72;
}

.arcana-bg .star.violet {
  fill: #b08fff;
}

.arcana-bg.light-mode .constellation.gold {
  stroke: #7a5210;
  opacity: 0.2;
}

.arcana-bg.light-mode .constellation.violet {
  stroke: #5530b8;
  opacity: 0.18;
}

.arcana-bg.light-mode .star.gold {
  fill: #7a5210;
}

.arcana-bg.light-mode .star.violet {
  fill: #5530b8;
}

.s-a1 { opacity: 0.75; }
.s-a2 { opacity: 0.85; }
.s-a3 { opacity: 0.75; }
.s-a4 { opacity: 0.8; }
.s-a5 { opacity: 0.65; }
.s-a6 { opacity: 0.65; }
.s-b1 { opacity: 0.7; }
.s-b2 { opacity: 0.82; }
.s-b3 { opacity: 0.68; }
.s-b4 { opacity: 0.72; }
.s-b5 { opacity: 0.62; }
.s-b6 { opacity: 0.65; }
.s-m1, .s-m2, .s-m3, .s-m4, .s-m5, .s-m6, .s-m7, .s-m8, .s-m9, .s-m10,
.s-m11, .s-m12, .s-m13, .s-m14, .s-m15, .s-m16 { opacity: 0.48; }

.arcana-bg.light-mode .s-a1 { opacity: 0.6; }
.arcana-bg.light-mode .s-a2 { opacity: 0.7; }
.arcana-bg.light-mode .s-a3 { opacity: 0.6; }
.arcana-bg.light-mode .s-a4 { opacity: 0.65; }
.arcana-bg.light-mode .s-a5 { opacity: 0.5; }
.arcana-bg.light-mode .s-a6 { opacity: 0.52; }
.arcana-bg.light-mode .s-b1 { opacity: 0.55; }
.arcana-bg.light-mode .s-b2 { opacity: 0.65; }
.arcana-bg.light-mode .s-b3 { opacity: 0.55; }
.arcana-bg.light-mode .s-b4 { opacity: 0.58; }
.arcana-bg.light-mode .s-b5 { opacity: 0.5; }
.arcana-bg.light-mode .s-b6 { opacity: 0.52; }
.arcana-bg.light-mode .s-m1,
.arcana-bg.light-mode .s-m2,
.arcana-bg.light-mode .s-m3,
.arcana-bg.light-mode .s-m4,
.arcana-bg.light-mode .s-m5,
.arcana-bg.light-mode .s-m6,
.arcana-bg.light-mode .s-m7,
.arcana-bg.light-mode .s-m8,
.arcana-bg.light-mode .s-m9,
.arcana-bg.light-mode .s-m10,
.arcana-bg.light-mode .s-m11,
.arcana-bg.light-mode .s-m12,
.arcana-bg.light-mode .s-m13,
.arcana-bg.light-mode .s-m14,
.arcana-bg.light-mode .s-m15,
.arcana-bg.light-mode .s-m16 { opacity: 0.4; }

.arcana-bg .gold-glow {
  fill: url(#glow-gold);
  opacity: 0.35;
}

.arcana-bg .gold-glow-small {
  fill: url(#glow-gold);
  opacity: 0.28;
}

.arcana-bg .violet-glow {
  fill: url(#glow-violet);
  opacity: 0.32;
}

.arcana-bg.light-mode .gold-glow {
  fill: url(#glow-gold-l);
  opacity: 0.28;
}

.arcana-bg.light-mode .gold-glow-small {
  fill: url(#glow-gold-l);
  opacity: 0.22;
}

.arcana-bg.light-mode .violet-glow {
  fill: url(#glow-violet-l);
  opacity: 0.26;
}

.arcana-bg .glyph-crescent {
  fill: #b08fff;
  opacity: 0.35;
}

.arcana-bg .glyph-star {
  fill: #d4af72;
  opacity: 0.38;
}

.arcana-bg .glyph-center {
  fill: #d4af72;
  opacity: 0.4;
}

.arcana-bg .divider-line {
  stroke: #d4af72;
  stroke-width: 0.6;
  opacity: 0.22;
}

.arcana-bg.light-mode .glyph-crescent {
  fill: #5530b8;
  opacity: 0.28;
}

.arcana-bg.light-mode .glyph-star {
  fill: #7a5210;
  opacity: 0.32;
}

.arcana-bg.light-mode .glyph-center {
  fill: #7a5210;
  opacity: 0.35;
}

.arcana-bg.light-mode .divider-line {
  stroke: #7a5210;
  opacity: 0.2;
}

.arcana-bg .frame-border {
  stroke: #b08fff;
  stroke-width: 0.5;
  opacity: 0.1;
}

.arcana-bg.light-mode .frame-border {
  stroke: #5530b8;
  opacity: 0.08;
}

.arcana-bg,
.arcana-bg .bg-blooms,
.arcana-bg .bg-svg {
  transition: background-color 0.5s ease, opacity 0.4s ease;
}

.arcana-bg .mandala,
.arcana-bg .constellation,
.arcana-bg .glyph-crescent,
.arcana-bg .glyph-star,
.arcana-bg .glyph-center,
.arcana-bg .divider-line,
.arcana-bg .frame-border {
  transition: opacity 0.4s ease, stroke 0.4s ease, fill 0.4s ease;
}
</style>
