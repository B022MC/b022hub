<script setup lang="ts">
import { computed } from 'vue'
import type { B022HubVendorId } from './b022hub'

const props = withDefaults(defineProps<{
  vendor: B022HubVendorId
  compact?: boolean
  progress?: number
}>(), {
  compact: false,
  progress: 1
})

const normalizedProgress = computed(() => Math.min(1, Math.max(0, props.progress)))

const label = computed(() => ({
  claude: 'Claude',
  openai: 'OpenAI',
  gemini: 'Gemini'
}[props.vendor]))

const componentStyle = computed<Record<string, string>>(() => {
  const progress = normalizedProgress.value

  return {
    '--vendor-symbol-opacity': `${0.5 + progress * 0.5}`,
    '--vendor-wordmark-opacity': `${0.34 + progress * 0.66}`,
    '--vendor-symbol-translate': `${(1 - progress) * 10}px`,
    '--vendor-wordmark-translate': `${(1 - progress) * 8}px`,
    '--vendor-scale': `${0.82 + progress * 0.18}`,
    '--vendor-glow-opacity': `${0.12 + progress * 0.4}`
  }
})

const geminiGradientId = `b022hub-gemini-${Math.random().toString(36).slice(2, 10)}`
</script>

<template>
  <div
    class="vendor-mark"
    :class="[`vendor-mark--${vendor}`, { 'vendor-mark--compact': compact }]"
    :style="componentStyle"
  >
    <span class="vendor-mark__symbol">
      <svg
        v-if="vendor === 'claude'"
        viewBox="0 0 80 80"
        class="vendor-mark__icon"
        xmlns="http://www.w3.org/2000/svg"
      >
        <g class="vendor-mark__claude-flower">
          <ellipse class="vendor-mark__claude-petal" cx="40" cy="14" rx="7" ry="13" />
          <ellipse class="vendor-mark__claude-petal" cx="40" cy="14" rx="7" ry="13" transform="rotate(45 40 40)" />
          <ellipse class="vendor-mark__claude-petal" cx="40" cy="14" rx="7" ry="13" transform="rotate(90 40 40)" />
          <ellipse class="vendor-mark__claude-petal" cx="40" cy="14" rx="7" ry="13" transform="rotate(135 40 40)" />
          <ellipse class="vendor-mark__claude-petal" cx="40" cy="14" rx="7" ry="13" transform="rotate(180 40 40)" />
          <ellipse class="vendor-mark__claude-petal" cx="40" cy="14" rx="7" ry="13" transform="rotate(225 40 40)" />
          <ellipse class="vendor-mark__claude-petal" cx="40" cy="14" rx="7" ry="13" transform="rotate(270 40 40)" />
          <ellipse class="vendor-mark__claude-petal" cx="40" cy="14" rx="7" ry="13" transform="rotate(315 40 40)" />
        </g>
        <circle class="vendor-mark__claude-core" cx="40" cy="40" r="12.5" />
        <circle class="vendor-mark__claude-cutout" cx="40" cy="40" r="4.5" />
      </svg>

      <svg
        v-else-if="vendor === 'openai'"
        viewBox="0 0 80 80"
        class="vendor-mark__icon"
        xmlns="http://www.w3.org/2000/svg"
      >
        <g class="vendor-mark__openai-knot" transform="translate(40 40)">
          <rect class="vendor-mark__openai-loop" x="-8.5" y="-33" width="17" height="30" rx="8.5" />
          <rect class="vendor-mark__openai-loop" x="-8.5" y="-33" width="17" height="30" rx="8.5" transform="rotate(60)" />
          <rect class="vendor-mark__openai-loop" x="-8.5" y="-33" width="17" height="30" rx="8.5" transform="rotate(120)" />
          <rect class="vendor-mark__openai-loop" x="-8.5" y="-33" width="17" height="30" rx="8.5" transform="rotate(180)" />
          <rect class="vendor-mark__openai-loop" x="-8.5" y="-33" width="17" height="30" rx="8.5" transform="rotate(240)" />
          <rect class="vendor-mark__openai-loop" x="-8.5" y="-33" width="17" height="30" rx="8.5" transform="rotate(300)" />
        </g>
        <circle class="vendor-mark__openai-cutout" cx="40" cy="40" r="12.5" />
      </svg>

      <svg
        v-else
        viewBox="0 0 80 80"
        class="vendor-mark__icon"
        xmlns="http://www.w3.org/2000/svg"
      >
        <defs>
          <linearGradient :id="geminiGradientId" x1="10%" y1="10%" x2="90%" y2="90%">
            <stop offset="0%" stop-color="#8fd6ff" />
            <stop offset="52%" stop-color="#f2b7ff" />
            <stop offset="100%" stop-color="#ffd79b" />
          </linearGradient>
        </defs>
        <path
          class="vendor-mark__gemini-star"
          d="M40 6 L48.5 31.5 L74 40 L48.5 48.5 L40 74 L31.5 48.5 L6 40 L31.5 31.5 Z"
          :fill="`url(#${geminiGradientId})`"
        />
        <path
          class="vendor-mark__gemini-star vendor-mark__gemini-star--small"
          d="M59 12 L62.8 23.2 L74 27 L62.8 30.8 L59 42 L55.2 30.8 L44 27 L55.2 23.2 Z"
          fill="#fff2cb"
        />
      </svg>
    </span>

    <span class="vendor-mark__wordmark">
      {{ label }}
    </span>
  </div>
</template>

<style scoped>
.vendor-mark {
  display: inline-flex;
  align-items: center;
  gap: 0.9rem;
}

.vendor-mark__symbol {
  position: relative;
  display: grid;
  height: 3.2rem;
  width: 3.2rem;
  place-items: center;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 1.2rem;
  background: rgba(255, 255, 255, 0.04);
  opacity: var(--vendor-symbol-opacity);
  transform: translateY(var(--vendor-symbol-translate)) scale(var(--vendor-scale));
  transition: transform 0.45s ease, opacity 0.45s ease;
}

.vendor-mark__symbol::before {
  position: absolute;
  inset: 18%;
  border-radius: 999px;
  content: '';
  filter: blur(20px);
  opacity: var(--vendor-glow-opacity);
}

.vendor-mark__icon {
  position: relative;
  z-index: 1;
  height: 72%;
  width: 72%;
  overflow: visible;
}

.vendor-mark__wordmark {
  opacity: var(--vendor-wordmark-opacity);
  transform: translateY(var(--vendor-wordmark-translate));
  transition: transform 0.45s ease, opacity 0.45s ease;
  font-size: clamp(1.7rem, 3vw, 2.35rem);
  font-weight: 600;
  letter-spacing: -0.06em;
  line-height: 0.92;
  color: #f8f3e7;
}

.vendor-mark--compact {
  gap: 0.72rem;
}

.vendor-mark--compact .vendor-mark__symbol {
  height: 2.7rem;
  width: 2.7rem;
  border-radius: 1rem;
}

.vendor-mark--compact .vendor-mark__wordmark {
  font-size: 1.18rem;
}

.vendor-mark--claude .vendor-mark__symbol::before {
  background: rgba(226, 146, 112, 0.4);
}

.vendor-mark--claude .vendor-mark__wordmark {
  font-family: 'Cormorant Garamond', 'Noto Serif SC', serif;
  color: #e4a183;
}

.vendor-mark__claude-petal,
.vendor-mark__claude-core {
  fill: #e4a183;
}

.vendor-mark__claude-flower {
  opacity: 0.9;
}

.vendor-mark__claude-cutout {
  fill: rgba(18, 17, 15, 0.92);
}

.vendor-mark--openai .vendor-mark__symbol::before {
  background: rgba(255, 248, 235, 0.28);
}

.vendor-mark--openai .vendor-mark__wordmark {
  letter-spacing: -0.05em;
  color: #f5f1e7;
}

.vendor-mark__openai-loop {
  fill: #f5f1e7;
}

.vendor-mark__openai-cutout {
  fill: rgba(18, 17, 15, 0.95);
}

.vendor-mark--gemini .vendor-mark__symbol::before {
  background: linear-gradient(135deg, rgba(143, 214, 255, 0.42), rgba(242, 183, 255, 0.34), rgba(255, 215, 155, 0.28));
}

.vendor-mark--gemini .vendor-mark__wordmark {
  background-image: linear-gradient(120deg, #95d6ff 0%, #f3b6ff 48%, #ffd48e 100%);
  color: transparent;
  -webkit-background-clip: text;
  background-clip: text;
}

.vendor-mark__gemini-star {
  filter: drop-shadow(0 0 10px rgba(143, 214, 255, 0.24));
}

.vendor-mark__gemini-star--small {
  opacity: 0.86;
}
</style>
