<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import B022HubCodePanel from './B022HubCodePanel.vue'
import type { B022HubConfigFile, B022HubIconName, B022HubInstallOption } from './b022hub'

const props = defineProps<{
  badge: string
  badgeIcon: B022HubIconName
  title: string
  description: string
  installOptions: B022HubInstallOption[]
  selectedInstallKey: string
  configFiles: B022HubConfigFile[]
  progress: number
  align?: 'left' | 'right'
}>()

const emit = defineEmits<{
  copy: [value: string]
  'update:selectedInstallKey': [value: string]
}>()

const { t } = useI18n()
const normalizedProgress = computed(() => Math.min(1, Math.max(0, props.progress)))
const direction = computed(() => (props.align === 'right' ? 1 : -1))
const sideClass = computed(() => (props.align === 'right' ? 'xl:col-start-2' : 'xl:col-start-1'))

function buildRevealStyle(offset: number) {
  const local = Math.min(1, Math.max(0, (normalizedProgress.value - offset) / (1 - offset || 1)))
  const opacity = 0.18 + local * 0.82
  const x = (1 - local) * 34 * direction.value
  const y = (1 - local) * 24

  return {
    opacity,
    transform: `translate3d(${x}px, ${y}px, 0)`
  }
}
</script>

<template>
  <section class="min-h-screen snap-start px-5 py-16 sm:px-8 md:px-14 lg:px-20">
    <div class="mx-auto grid min-h-[calc(100vh-8rem)] max-w-7xl items-center xl:grid-cols-2">
      <div :class="['max-w-2xl space-y-5 xl:max-w-xl', sideClass]" :style="buildRevealStyle(0.04)">
        <div
          class="b022hub-theme-chip inline-flex items-center gap-2 rounded-full px-4 py-2 text-xs font-medium uppercase tracking-[0.18em]"
        >
          <Icon :name="badgeIcon" size="sm" />
          <span>{{ badge }}</span>
        </div>

        <div class="space-y-4">
          <h2 class="b022hub-theme-text-strong font-display text-4xl font-semibold leading-tight md:text-5xl">
            {{ title }}
          </h2>
          <p class="b022hub-theme-text-body text-base leading-8 md:text-lg">
            {{ description }}
          </p>
        </div>

        <div
          class="b022hub-theme-panel rounded-[1.5rem] p-4"
          :style="buildRevealStyle(0.16)"
        >
          <div class="mb-4 flex flex-wrap gap-2">
            <button
              v-for="option in installOptions"
              :key="option.key"
              class="rounded-full border px-3 py-1.5 text-sm transition"
              :class="option.key === selectedInstallKey
                ? 'b022hub-theme-option-active'
                : 'b022hub-theme-option'"
              @click="emit('update:selectedInstallKey', option.key)"
            >
              <span>{{ option.label }}</span>
              <span class="ml-2 text-xs opacity-75">{{ option.hint }}</span>
            </button>
          </div>

          <div
            class="b022hub-theme-panel-strong rounded-[1.15rem] px-4 py-4"
          >
            <div class="b022hub-theme-text-muted mb-2 text-xs font-medium uppercase tracking-[0.18em]">
              {{ t('home.installCommand') }}
            </div>
            <pre class="b022hub-theme-code overflow-x-auto text-sm leading-7"><code>{{ installOptions.find((option) => option.key === selectedInstallKey)?.command }}</code></pre>
          </div>
        </div>

        <div class="space-y-4" :style="buildRevealStyle(0.24)">
          <B022HubCodePanel
            v-for="file in configFiles"
            :key="file.path"
            :path="file.path"
            :language-label="file.languageLabel"
            :code="file.content"
            @copy="emit('copy', $event)"
          />
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
pre {
  font-family: 'JetBrains Mono', 'SFMono-Regular', Menlo, Monaco, Consolas, monospace;
}
</style>
