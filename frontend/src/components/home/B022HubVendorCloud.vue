<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import B022HubLineLogo from './B022HubLineLogo.vue'
import B022HubVendorMark from './B022HubVendorMark.vue'
import type { B022HubVendorId } from './b022hub'

interface VendorCard {
  id: B022HubVendorId
  offset: number
  eyebrow: string
  description: string
}

const props = defineProps<{
  progress: number
}>()

const { t } = useI18n()

const normalizedProgress = computed(() => Math.min(1, Math.max(0, props.progress)))

const vendorCards = computed<VendorCard[]>(() => [
  {
    id: 'claude',
    offset: 0.14,
    eyebrow: t('home.progress.orbit.vendors.claude.eyebrow'),
    description: t('home.progress.orbit.vendors.claude.description')
  },
  {
    id: 'openai',
    offset: 0.28,
    eyebrow: t('home.progress.orbit.vendors.openai.eyebrow'),
    description: t('home.progress.orbit.vendors.openai.description')
  },
  {
    id: 'gemini',
    offset: 0.42,
    eyebrow: t('home.progress.orbit.vendors.gemini.eyebrow'),
    description: t('home.progress.orbit.vendors.gemini.description')
  }
])

function resolveLocal(offset: number) {
  return Math.min(1, Math.max(0, (normalizedProgress.value - offset) / 0.32))
}

const hubStyle = computed(() => {
  const opacity = 0.16 + normalizedProgress.value * 0.84
  const scale = 0.86 + normalizedProgress.value * 0.14
  const y = (1 - normalizedProgress.value) * 18

  return {
    opacity,
    transform: `translate3d(-50%, calc(-50% + ${y}px), 0) scale(${scale})`
  }
})

function buildCardStyle(card: VendorCard) {
  const local = resolveLocal(card.offset)
  const opacity = local
  const scale = 0.72 + local * 0.28
  const blur = (1 - local) * 22

  const deltaX = card.id === 'claude'
    ? (1 - local) * -56
    : card.id === 'gemini'
      ? (1 - local) * 56
      : 0

  const deltaY = card.id === 'openai'
    ? (1 - local) * -56
    : (1 - local) * 48

  const rotate = card.id === 'claude'
    ? (1 - local) * -9
    : card.id === 'gemini'
      ? (1 - local) * 9
      : 0

  return {
    opacity,
    filter: `blur(${blur}px)`,
    transform: `translate3d(${deltaX}px, ${deltaY}px, 0) scale(${scale}) rotate(${rotate}deg)`
  }
}
</script>

<template>
  <div class="vendor-cloud">
    <div class="vendor-cloud__grid"></div>
    <div class="vendor-cloud__glow vendor-cloud__glow--warm"></div>
    <div class="vendor-cloud__glow vendor-cloud__glow--cool"></div>

    <span class="vendor-cloud__ring vendor-cloud__ring--outer"></span>
    <span class="vendor-cloud__ring vendor-cloud__ring--inner"></span>

    <div class="vendor-cloud__hub" :style="hubStyle">
      <div class="vendor-cloud__hub-shell">
        <p class="vendor-cloud__hub-eyebrow">
          {{ t('home.progress.orbit.hubEyebrow') }}
        </p>
        <B022HubLineLogo
          class="vendor-cloud__hub-logo"
          tone="progress"
          :animated="false"
          :width="216"
          :height="110"
        />
        <p class="vendor-cloud__hub-description">
          {{ t('home.progress.orbit.hubDescription') }}
        </p>
      </div>
    </div>

    <article
      v-for="card in vendorCards"
      :key="card.id"
      class="vendor-cloud__card"
      :class="`vendor-cloud__card--${card.id}`"
      :style="buildCardStyle(card)"
    >
      <p class="vendor-cloud__eyebrow">
        {{ card.eyebrow }}
      </p>
      <B022HubVendorMark
        class="vendor-cloud__brand"
        :vendor="card.id"
        :progress="resolveLocal(card.offset)"
      />
      <p class="vendor-cloud__description">
        {{ card.description }}
      </p>
    </article>
  </div>
</template>

<style scoped>
.vendor-cloud {
  position: relative;
  margin: 0 auto;
  min-height: 420px;
  max-width: 1040px;
  overflow: hidden;
  border: 1px solid rgba(227, 224, 211, 0.12);
  border-radius: 2.5rem;
  background:
    radial-gradient(circle at 26% 18%, rgba(204, 120, 92, 0.18), transparent 26%),
    radial-gradient(circle at 78% 18%, rgba(122, 185, 255, 0.16), transparent 28%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.05), rgba(255, 255, 255, 0.02));
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.08),
    0 36px 90px rgba(0, 0, 0, 0.24);
}

.vendor-cloud__grid,
.vendor-cloud__glow,
.vendor-cloud__ring {
  pointer-events: none;
  position: absolute;
}

.vendor-cloud__grid {
  inset: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.035) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.035) 1px, transparent 1px);
  background-size: 34px 34px;
  mask-image: radial-gradient(circle at center, rgba(0, 0, 0, 0.85), transparent 100%);
}

.vendor-cloud__glow {
  border-radius: 999px;
  filter: blur(34px);
  opacity: 0.5;
}

.vendor-cloud__glow--warm {
  left: -5%;
  top: 56%;
  height: 220px;
  width: 220px;
  background: rgba(204, 120, 92, 0.22);
}

.vendor-cloud__glow--cool {
  right: -5%;
  top: 10%;
  height: 240px;
  width: 240px;
  background: rgba(123, 211, 255, 0.18);
}

.vendor-cloud__ring {
  left: 50%;
  top: 50%;
  border: 1px solid rgba(212, 162, 127, 0.16);
  border-radius: 999px;
  transform: translate(-50%, -50%);
}

.vendor-cloud__ring--outer {
  height: 300px;
  width: 300px;
  animation: vendor-cloud-pulse 7s ease-in-out infinite;
}

.vendor-cloud__ring--inner {
  height: 214px;
  width: 214px;
  animation: vendor-cloud-pulse 7s ease-in-out infinite 1.3s;
}

.vendor-cloud__hub {
  position: absolute;
  left: 50%;
  top: 50%;
  will-change: transform, opacity;
}

.vendor-cloud__hub-shell {
  display: flex;
  width: min(280px, 72vw);
  flex-direction: column;
  align-items: center;
  gap: 0.7rem;
  border: 1px solid rgba(227, 224, 211, 0.12);
  border-radius: 2rem;
  background: rgba(14, 13, 12, 0.84);
  padding: 1.3rem 1.25rem;
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.08),
    0 24px 70px rgba(0, 0, 0, 0.34);
  text-align: center;
  backdrop-filter: blur(18px);
}

.vendor-cloud__hub-eyebrow {
  font-size: 0.7rem;
  font-weight: 600;
  letter-spacing: 0.22em;
  text-transform: uppercase;
  color: #9d998f;
}

.vendor-cloud__hub-logo {
  width: min(216px, 100%);
}

.vendor-cloud__hub-description {
  max-width: 18rem;
  font-size: 0.9rem;
  line-height: 1.65;
  color: #b8b0a2;
}

.vendor-cloud__card {
  position: absolute;
  width: min(280px, 72vw);
  overflow: hidden;
  border: 1px solid rgba(227, 224, 211, 0.14);
  border-radius: 1.75rem;
  background: rgba(18, 17, 15, 0.78);
  padding: 1.1rem 1.15rem 1.15rem;
  text-align: left;
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.08),
    0 24px 64px rgba(0, 0, 0, 0.26);
  backdrop-filter: blur(16px);
  will-change: transform, opacity, filter;
}

.vendor-cloud__card::before {
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  height: 1px;
  content: '';
}

.vendor-cloud__card--claude {
  left: 3%;
  bottom: 7%;
}

.vendor-cloud__card--claude::before {
  background: linear-gradient(90deg, rgba(228, 161, 131, 0.92), transparent 72%);
}

.vendor-cloud__card--openai {
  left: 50%;
  top: 6%;
  transform-origin: center bottom;
  margin-left: -140px;
}

.vendor-cloud__card--openai::before {
  background: linear-gradient(90deg, rgba(255, 246, 230, 0.82), transparent 72%);
}

.vendor-cloud__card--gemini {
  right: 3%;
  bottom: 8%;
}

.vendor-cloud__card--gemini::before {
  background: linear-gradient(90deg, rgba(149, 214, 255, 0.88), rgba(243, 182, 255, 0.7), transparent 78%);
}

.vendor-cloud__eyebrow {
  font-size: 0.7rem;
  font-weight: 600;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: #9d998f;
}

.vendor-cloud__brand {
  margin-top: 0.65rem;
}

.vendor-cloud__description {
  margin-top: 0.65rem;
  font-size: 0.88rem;
  line-height: 1.75;
  color: #c9c3b4;
}

@media (max-width: 1023px) {
  .vendor-cloud {
    min-height: 520px;
  }

  .vendor-cloud__ring--outer {
    height: 260px;
    width: 260px;
  }

  .vendor-cloud__ring--inner {
    height: 190px;
    width: 190px;
  }

  .vendor-cloud__card--claude {
    left: 2%;
  }

  .vendor-cloud__card--openai {
    margin-left: -130px;
  }

  .vendor-cloud__card--gemini {
    right: 2%;
  }
}

@media (max-width: 767px) {
  .vendor-cloud {
    display: flex;
    min-height: auto;
    flex-direction: column;
    align-items: center;
    gap: 0.9rem;
    padding: 1.1rem 0.85rem 1.25rem;
  }

  .vendor-cloud__ring--outer {
    height: 220px;
    width: 220px;
  }

  .vendor-cloud__ring--inner {
    height: 164px;
    width: 164px;
  }

  .vendor-cloud__hub {
    position: relative;
    left: auto;
    top: auto;
    margin-top: 0.35rem;
    transform: none !important;
  }

  .vendor-cloud__hub-shell {
    width: min(280px, 100%);
    padding: 1rem 1rem 1.1rem;
  }

  .vendor-cloud__card {
    position: relative;
    left: auto;
    top: auto;
    right: auto;
    bottom: auto;
    width: min(280px, 100%);
    margin-left: 0;
  }

  .vendor-cloud__card--claude,
  .vendor-cloud__card--openai,
  .vendor-cloud__card--gemini {
    left: auto;
    top: auto;
    right: auto;
    bottom: auto;
    margin-left: 0;
  }
}

@keyframes vendor-cloud-pulse {
  0%,
  100% {
    opacity: 0.22;
    transform: translate(-50%, -50%) scale(1);
  }

  50% {
    opacity: 0.06;
    transform: translate(-50%, -50%) scale(1.08);
  }
}
</style>
