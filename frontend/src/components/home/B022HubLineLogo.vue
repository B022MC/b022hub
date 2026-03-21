<script setup lang="ts">
import { computed } from 'vue'

type B022HubLogoTone = 'home' | 'claude' | 'codex' | 'gemini' | 'progress'

interface StrokeOffset {
  x: number
  y: number
  opacity: number
  width: number
}

const props = withDefaults(defineProps<{
  tone?: B022HubLogoTone
  width?: number
  height?: number
  animated?: boolean
}>(), {
  tone: 'home',
  width: 420,
  height: 260,
  animated: true
})

const gradientId = `b022hub-line-${Math.random().toString(36).slice(2, 10)}`

const pathDefs = [
  'M26 20 C20 48 20 86 24 126',
  'M27 81 C38 70 59 68 71 80 C82 92 79 113 66 124 C54 134 35 131 24 117',
  'M110 62 C91 63 81 79 82 97 C83 116 96 131 114 132 C132 133 145 117 145 96 C145 76 132 60 110 62 Z',
  'M177 74 C188 61 209 58 221 68 C232 78 228 94 214 105 C203 114 190 120 181 132',
  'M181 132 H226',
  'M246 74 C257 61 279 58 290 68 C301 78 297 95 283 106 C272 115 259 121 250 132',
  'M250 132 H296'
] as const

const ghostOffsets: StrokeOffset[] = [
  { x: -4, y: 2.5, opacity: 0.16, width: 3.3 },
  { x: 0, y: 0, opacity: 0.22, width: 2.8 },
  { x: 4, y: -2.5, opacity: 0.16, width: 2.4 }
]

const drawOffsets: StrokeOffset[] = [
  { x: -3.5, y: 2, opacity: 0.48, width: 2.6 },
  { x: 0, y: 0, opacity: 1, width: 2.9 },
  { x: 3.5, y: -2, opacity: 0.62, width: 2.45 }
]

const toneStops = computed(() => ({
  home: ['#cc785c', '#f0c8b2', '#d38b6a'],
  claude: ['#cc785c', '#f1c2a7', '#f4e2d5'],
  codex: ['#f8f3e7', '#d4a27f', '#fff6ea'],
  gemini: ['#84ccff', '#d0b9ff', '#f8a7d0'],
  progress: ['#d89e79', '#f8f3e7', '#8dcfff']
}[props.tone]))

const ghostStroke = computed(() => ({
  home: '#5a2f25',
  claude: '#693328',
  codex: '#5c5448',
  gemini: '#324657',
  progress: '#4f3b32'
}[props.tone]))

function buildTransform(offset: StrokeOffset) {
  return `translate(${offset.x} ${offset.y})`
}

function buildDrawStyle(pathIndex: number, layerIndex: number, offset: StrokeOffset) {
  const delay = layerIndex * 0.16 + pathIndex * 0.1

  if (!props.animated) {
    return {
      opacity: offset.opacity,
      strokeWidth: offset.width,
      strokeDashoffset: 0,
      animation: 'none'
    }
  }

  return {
    opacity: offset.opacity,
    strokeWidth: offset.width,
    animationDelay: `${delay}s`
  }
}
</script>

<template>
  <div
    class="line-logo"
    :class="{ 'line-logo--static': !animated }"
    :style="{ width: `${width}px`, height: `${height}px` }"
  >
    <svg viewBox="0 0 320 150" class="line-logo__svg" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <linearGradient :id="gradientId" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" :stop-color="toneStops[0]" />
          <stop offset="50%" :stop-color="toneStops[1]" />
          <stop offset="100%" :stop-color="toneStops[2]" />
        </linearGradient>
      </defs>

      <g class="line-logo__ghost">
        <g
          v-for="offset in ghostOffsets"
          :key="`ghost-${offset.x}-${offset.y}`"
          :transform="buildTransform(offset)"
        >
          <path
            v-for="(d, pathIndex) in pathDefs"
            :key="`ghost-${pathIndex}`"
            :d="d"
            class="line-logo__ghost-path"
            :stroke="ghostStroke"
            :stroke-width="offset.width"
          />
        </g>
      </g>

      <g class="line-logo__main">
        <g
          v-for="(offset, layerIndex) in drawOffsets"
          :key="`draw-${offset.x}-${offset.y}`"
          :transform="buildTransform(offset)"
        >
          <path
            v-for="(d, pathIndex) in pathDefs"
            :key="`draw-${layerIndex}-${pathIndex}`"
            :d="d"
            pathLength="1"
            class="line-logo__draw-path"
            :stroke="`url(#${gradientId})`"
            :style="buildDrawStyle(pathIndex, layerIndex, offset)"
          />
        </g>
      </g>
    </svg>
  </div>
</template>

<style scoped>
.line-logo {
  display: flex;
  align-items: center;
  justify-content: center;
}

.line-logo__svg {
  height: 100%;
  width: 100%;
  overflow: visible;
}

.line-logo__ghost-path,
.line-logo__draw-path {
  fill: none;
  stroke-linecap: round;
  stroke-linejoin: round;
  vector-effect: non-scaling-stroke;
}

.line-logo__ghost {
  filter: blur(0.25px);
}

.line-logo__ghost-path {
  opacity: 0.22;
}

.line-logo__main {
  filter: drop-shadow(0 16px 36px rgba(0, 0, 0, 0.22));
}

.line-logo__draw-path {
  stroke-dasharray: 1;
  stroke-dashoffset: 1;
  animation: line-logo-draw 1.7s cubic-bezier(0.22, 1, 0.36, 1) forwards;
}

.line-logo--static .line-logo__draw-path {
  animation: none;
  stroke-dashoffset: 0;
}

@keyframes line-logo-draw {
  0% {
    stroke-dashoffset: 1;
  }

  70%,
  100% {
    stroke-dashoffset: 0;
  }
}
</style>
