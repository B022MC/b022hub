<script setup lang="ts">
import { OPENAI_LOGO_PATH, type B022HubVendorId } from './b022hub'

const props = withDefaults(defineProps<{
  vendor: B022HubVendorId
  width?: number
  height?: number
}>(), {
  width: 400,
  height: 400
})

const geminiGradientId = `b022hub-glyph-gemini-${Math.random().toString(36).slice(2, 10)}`
</script>

<template>
  <div
    class="vendor-glyph"
    :class="`vendor-glyph--${vendor}`"
    :style="{ width: `${props.width}px`, height: `${props.height}px` }"
  >
    <svg
      v-if="vendor === 'claude'"
      viewBox="0 0 440 440"
      class="vendor-glyph__svg"
      xmlns="http://www.w3.org/2000/svg"
    >
      <g class="vendor-glyph__ghost vendor-glyph__ghost--claude">
        <ellipse
          v-for="angle in [0, 30, 60, 90, 120, 150, 180, 210, 240, 270, 300, 330]"
          :key="`claude-ghost-${angle}`"
          class="vendor-glyph__claude-ghost-petal"
          cx="220"
          cy="108"
          rx="20"
          ry="106"
          :transform="`rotate(${angle} 220 220)`"
        />
      </g>

      <g class="vendor-glyph__main vendor-glyph__main--claude">
        <ellipse
          v-for="angle in [0, 30, 60, 90, 120, 150, 180, 210, 240, 270, 300, 330]"
          :key="`claude-main-${angle}`"
          class="vendor-glyph__claude-petal"
          cx="220"
          cy="116"
          rx="18"
          ry="96"
          :transform="`rotate(${angle} 220 220)`"
        />
        <circle class="vendor-glyph__claude-core" cx="220" cy="220" r="54" />
        <circle class="vendor-glyph__claude-cutout" cx="220" cy="220" r="16" />
      </g>
    </svg>

    <svg
      v-else-if="vendor === 'openai'"
      viewBox="0 0 24 24"
      class="vendor-glyph__svg"
      xmlns="http://www.w3.org/2000/svg"
    >
      <g class="vendor-glyph__ghost vendor-glyph__ghost--openai">
        <path class="vendor-glyph__openai-ghost-path" :d="OPENAI_LOGO_PATH" />
      </g>

      <g class="vendor-glyph__main vendor-glyph__main--openai">
        <path class="vendor-glyph__openai-path" :d="OPENAI_LOGO_PATH" />
      </g>
    </svg>

    <svg
      v-else
      viewBox="0 0 440 440"
      class="vendor-glyph__svg"
      xmlns="http://www.w3.org/2000/svg"
    >
      <defs>
        <linearGradient :id="geminiGradientId" x1="18%" y1="18%" x2="82%" y2="82%">
          <stop offset="0%" stop-color="#89d6ff" />
          <stop offset="48%" stop-color="#f3b3ff" />
          <stop offset="100%" stop-color="#ffd792" />
        </linearGradient>
      </defs>

      <g class="vendor-glyph__ghost vendor-glyph__ghost--gemini">
        <path
          class="vendor-glyph__gemini-ghost-star"
          d="M220 44 L251 189 L396 220 L251 251 L220 396 L189 251 L44 220 L189 189 Z"
        />
        <path
          class="vendor-glyph__gemini-ghost-star vendor-glyph__gemini-ghost-star--small"
          d="M335 70 L348 115 L392 128 L348 141 L335 186 L322 141 L278 128 L322 115 Z"
        />
        <path
          class="vendor-glyph__gemini-ghost-star vendor-glyph__gemini-ghost-star--tiny"
          d="M110 98 L119 128 L149 137 L119 146 L110 176 L101 146 L71 137 L101 128 Z"
        />
      </g>

      <g class="vendor-glyph__main vendor-glyph__main--gemini">
        <path
          class="vendor-glyph__gemini-star"
          d="M220 64 L246 194 L376 220 L246 246 L220 376 L194 246 L64 220 L194 194 Z"
          :fill="`url(#${geminiGradientId})`"
        />
        <path
          class="vendor-glyph__gemini-star vendor-glyph__gemini-star--small"
          d="M334 84 L345 119 L380 130 L345 141 L334 176 L323 141 L288 130 L323 119 Z"
          fill="#fff0c8"
        />
        <path
          class="vendor-glyph__gemini-star vendor-glyph__gemini-star--tiny"
          d="M112 120 L120 146 L146 154 L120 162 L112 188 L104 162 L78 154 L104 146 Z"
          fill="#a5deff"
        />
      </g>
    </svg>
  </div>
</template>

<style scoped>
.vendor-glyph {
  display: flex;
  align-items: center;
  justify-content: center;
}

.vendor-glyph__svg {
  height: 100%;
  width: 100%;
  overflow: visible;
}

.vendor-glyph__ghost {
  filter: blur(0.2px);
  opacity: 0.28;
}

.vendor-glyph__main {
  transform-origin: center;
}

.vendor-glyph__ghost--claude {
  animation: vendor-glyph-breathe 7.5s ease-in-out infinite;
}

.vendor-glyph__main--claude {
  animation: vendor-glyph-spin-slow 20s linear infinite;
}

.vendor-glyph__claude-ghost-petal {
  fill: none;
  stroke: rgba(204, 120, 92, 0.24);
  stroke-width: 7;
}

.vendor-glyph__claude-petal {
  fill: rgba(204, 120, 92, 0.88);
}

.vendor-glyph__claude-core {
  fill: rgba(204, 120, 92, 0.96);
}

.vendor-glyph__claude-cutout {
  fill: rgba(14, 13, 12, 0.94);
}

.vendor-glyph__ghost--openai {
  animation: vendor-glyph-breathe 8s ease-in-out infinite;
}

.vendor-glyph__main--openai {
  animation: vendor-glyph-float 8.5s ease-in-out infinite;
}

.vendor-glyph__openai-ghost-path {
  fill: rgba(248, 243, 231, 0.18);
}

.vendor-glyph__openai-path {
  fill: #f5f1e7;
}

.vendor-glyph__ghost--gemini {
  animation: vendor-glyph-breathe 8.5s ease-in-out infinite;
}

.vendor-glyph__main--gemini {
  animation: vendor-glyph-float 7.2s ease-in-out infinite;
}

.vendor-glyph__gemini-ghost-star {
  fill: none;
  stroke: rgba(149, 214, 255, 0.2);
  stroke-width: 9;
}

.vendor-glyph__gemini-ghost-star--small {
  stroke: rgba(243, 182, 255, 0.18);
}

.vendor-glyph__gemini-ghost-star--tiny {
  stroke: rgba(255, 215, 155, 0.18);
}

.vendor-glyph__gemini-star {
  filter: drop-shadow(0 0 32px rgba(137, 214, 255, 0.24));
}

.vendor-glyph__gemini-star--small {
  opacity: 0.92;
}

.vendor-glyph__gemini-star--tiny {
  opacity: 0.88;
}

@keyframes vendor-glyph-float {
  0%,
  100% {
    transform: translateY(0) scale(1);
  }

  50% {
    transform: translateY(-8px) scale(1.015);
  }
}

@keyframes vendor-glyph-breathe {
  0%,
  100% {
    opacity: 0.2;
    transform: scale(0.98);
  }

  50% {
    opacity: 0.34;
    transform: scale(1.02);
  }
}

@keyframes vendor-glyph-spin-slow {
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
}
</style>
