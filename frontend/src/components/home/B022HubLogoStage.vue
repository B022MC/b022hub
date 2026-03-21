<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import B022HubLineLogo from './B022HubLineLogo.vue'
import type { B022HubVendorId } from './b022hub'
import B022HubVendorGlyph from './B022HubVendorGlyph.vue'

type B022HubLogoTone = 'home' | 'claude' | 'codex' | 'gemini' | 'progress'

const props = defineProps<{
  siteName: string
  siteLogo: string
  activeSectionId: string
  vendorProgress: number
}>()

const { t } = useI18n()

const normalizedVendorProgress = computed(() => Math.min(1, Math.max(0, props.vendorProgress)))
const glyphVendors: B022HubVendorId[] = ['claude', 'openai', 'gemini']

const logoTone = computed<B022HubLogoTone>(() => {
  switch (props.activeSectionId) {
    case 'claude':
      return 'claude'
    case 'codex':
      return 'codex'
    case 'gemini':
      return 'gemini'
    case 'progress':
      return 'progress'
    default:
      return 'home'
  }
})

const activeLabel = computed(() => ({
  home: props.siteName,
  claude: t('home.logo.layers.claude'),
  codex: t('home.logo.layers.codex'),
  gemini: t('home.logo.layers.gemini'),
  progress: t('home.logo.layers.progress')
}[props.activeSectionId] || props.siteName))

const activeGlyphVendor = computed<B022HubVendorId | null>(() => {
  switch (props.activeSectionId) {
    case 'claude':
      return 'claude'
    case 'codex':
      return 'openai'
    case 'gemini':
      return 'gemini'
    default:
      return null
  }
})

const stageStyle = computed(() => ({
  home: {
    transform: 'scale(1.1) translateY(-12vh)',
    opacity: 0.25
  },
  claude: {
    transform: 'translateX(-25vw) scale(1)',
    opacity: 1
  },
  codex: {
    transform: 'translateX(25vw) scale(1)',
    opacity: 1
  },
  gemini: {
    transform: 'translateX(-25vw) scale(1)',
    opacity: 1
  },
  progress: {
    transform: 'translateX(0) scale(1)',
    opacity: 0.15
  }
}[props.activeSectionId] || {
  transform: 'scale(1.1) translateY(-12vh)',
  opacity: 0.25
}))

const lineLogoStyle = computed(() => {
  if (props.activeSectionId === 'home') {
    return {
      opacity: 1,
      transform: 'translate3d(0, 0, 0) scale(1)'
    }
  }

  if (props.activeSectionId === 'progress') {
    const opacity = 0.18 + normalizedVendorProgress.value * 0.18
    const scale = 0.96 + normalizedVendorProgress.value * 0.06

    return {
      opacity,
      transform: `translate3d(0, 0, 0) scale(${scale})`
    }
  }

  return {
    opacity: 0.08,
    transform: 'translate3d(0, 0, 0) scale(0.88)'
  }
})

function buildGlyphStyle(vendor: B022HubVendorId) {
  const isActive = activeGlyphVendor.value === vendor
  const baseX = vendor === 'claude'
    ? -18
    : vendor === 'openai'
      ? 22
      : -12
  const baseY = vendor === 'openai'
    ? 10
    : vendor === 'gemini'
      ? -6
      : 0
  const scale = isActive ? 1 : 0.76
  const blur = isActive ? 0 : 18

  return {
    opacity: isActive ? 1 : 0,
    filter: `blur(${blur}px)`,
    transform: `translate3d(${isActive ? 0 : baseX}px, ${isActive ? 0 : baseY}px, 0) scale(${scale})`
  }
}
</script>

<template>
  <div class="pointer-events-none fixed inset-0 z-10 hidden items-center justify-center overflow-hidden xl:flex">
    <div
      class="b022hub-stage"
      :style="{
        ...stageStyle,
        transition: 'transform 0.8s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.6s ease-out'
      }"
    >
      <div class="b022hub-stage__logo-wrap">
        <div class="b022hub-stage__symbol-stack">
          <B022HubLineLogo
            class="b022hub-stage__logo"
            :tone="logoTone"
            :animated="true"
            :width="430"
            :height="270"
            :style="{
              ...lineLogoStyle,
              transition: 'transform 0.8s cubic-bezier(0.22, 1, 0.36, 1), opacity 0.45s ease-out'
            }"
          />

          <B022HubVendorGlyph
            v-for="vendor in glyphVendors"
            :key="vendor"
            class="b022hub-stage__glyph"
            :class="`b022hub-stage__glyph--${vendor}`"
            :vendor="vendor"
            :width="392"
            :height="392"
            :style="{
              ...buildGlyphStyle(vendor),
              transition: 'transform 0.8s cubic-bezier(0.22, 1, 0.36, 1), opacity 0.45s ease-out, filter 0.45s ease-out'
            }"
          />
        </div>

        <div
          class="b022hub-stage__glow"
          :class="`b022hub-stage__glow--${logoTone}`"
        ></div>
      </div>

      <div class="b022hub-stage__label">
        <p class="b022hub-stage__eyebrow">
          {{ t('home.logo.activeLayer') }}
        </p>
        <p class="b022hub-stage__title">
          {{ activeLabel }}
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.b022hub-stage {
  position: relative;
  display: flex;
  height: 420px;
  width: 560px;
  flex-direction: column;
  align-items: center;
  gap: 1.05rem;
  will-change: transform, opacity;
}

.b022hub-stage__logo-wrap {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.b022hub-stage__symbol-stack {
  position: relative;
  display: grid;
  height: 392px;
  width: 392px;
  place-items: center;
}

.b022hub-stage__logo {
  position: relative;
  z-index: 1;
  will-change: transform, opacity;
}

.b022hub-stage__glyph {
  position: absolute;
  inset: 0;
  z-index: 2;
  will-change: transform, opacity, filter;
}

.b022hub-stage__glyph--claude {
  transform-origin: center;
}

.b022hub-stage__glyph--openai {
  transform-origin: center;
}

.b022hub-stage__glyph--gemini {
  transform-origin: center;
}

.b022hub-stage__glow {
  position: absolute;
  left: 50%;
  top: 50%;
  height: 300px;
  width: 300px;
  border-radius: 999px;
  filter: blur(46px);
  opacity: 0.34;
  transform: translate(-50%, -50%);
}

.b022hub-stage__glow--home,
.b022hub-stage__glow--claude {
  background: radial-gradient(circle, rgba(204, 120, 92, 0.34), transparent 68%);
}

.b022hub-stage__glow--codex {
  background: radial-gradient(circle, rgba(248, 243, 231, 0.2), transparent 70%);
}

.b022hub-stage__glow--gemini {
  background: radial-gradient(circle, rgba(145, 214, 255, 0.28), transparent 68%);
}

.b022hub-stage__glow--progress {
  background: radial-gradient(circle, rgba(143, 207, 255, 0.18), transparent 72%);
}

.b022hub-stage__label {
  position: absolute;
  left: 50%;
  top: calc(50% + 9.4rem);
  transform: translateX(-50%);
  text-align: center;
}

.b022hub-stage__eyebrow {
  font-size: 0.72rem;
  letter-spacing: 0.24em;
  text-transform: uppercase;
  color: #9d998f;
}

.b022hub-stage__title {
  margin-top: 0.55rem;
  font-family: 'Cormorant Garamond', 'Noto Serif SC', serif;
  font-size: 2.15rem;
  line-height: 0.95;
  color: #f8f3e7;
}
</style>
