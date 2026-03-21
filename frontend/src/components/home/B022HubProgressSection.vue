<script setup lang="ts">
import { computed } from 'vue'
import Icon from '@/components/icons/Icon.vue'
import B022HubVendorCloud from './B022HubVendorCloud.vue'
import type { B022HubFeatureCard } from './b022hub'

const props = defineProps<{
  badge: string
  title: string
  description: string
  cards: B022HubFeatureCard[]
  docHref: string
  docLabel: string
  primaryPath: string
  primaryLabel: string
  progress: number
}>()

const normalizedProgress = computed(() => Math.min(1, Math.max(0, props.progress)))

function buildRevealStyle(offset: number) {
  const local = Math.min(1, Math.max(0, (normalizedProgress.value - offset) / (1 - offset || 1)))
  const opacity = 0.18 + local * 0.82
  const y = (1 - local) * 26
  const scale = 0.94 + local * 0.06

  return {
    opacity,
    transform: `translate3d(0, ${y}px, 0) scale(${scale})`
  }
}
</script>

<template>
  <section class="min-h-screen snap-start px-5 py-16 sm:px-8 md:px-14 lg:px-20">
    <div class="mx-auto flex min-h-[calc(100vh-8rem)] max-w-7xl items-center">
      <div class="w-full space-y-8 text-center">
        <div class="space-y-4" :style="buildRevealStyle(0.04)">
        <div
          class="b022hub-theme-chip inline-flex items-center gap-2 rounded-full px-4 py-2 text-xs font-medium uppercase tracking-[0.18em]"
        >
            <Icon name="sparkles" size="sm" />
            <span>{{ badge }}</span>
          </div>

          <h2 class="b022hub-theme-text-strong font-display text-4xl font-semibold md:text-5xl">
            {{ title }}
          </h2>
          <p class="b022hub-theme-text-body mx-auto max-w-2xl text-base leading-8 md:text-lg">
            {{ description }}
          </p>
        </div>

        <div :style="buildRevealStyle(0.14)">
          <B022HubVendorCloud :progress="normalizedProgress" />
        </div>

        <div class="grid gap-4 md:grid-cols-3" :style="buildRevealStyle(0.22)">
          <article
            v-for="card in cards"
            :key="card.title"
            class="rounded-[1.5rem] border p-5 text-left backdrop-blur-xl transition"
            :class="card.status === 'completed'
              ? 'b022hub-theme-panel-completed'
              : 'b022hub-theme-panel border-dashed'"
          >
            <div
              class="b022hub-theme-icon-soft mb-4 flex h-11 w-11 items-center justify-center rounded-2xl"
            >
              <Icon :name="card.icon" size="md" />
            </div>

            <h3 class="b022hub-theme-text-strong text-lg font-semibold">
              {{ card.title }}
            </h3>
            <p class="b022hub-theme-text-body mt-2 text-sm leading-7">
              {{ card.description }}
            </p>
            <p
              class="mt-4 inline-flex rounded-full border px-2.5 py-1 text-xs font-medium"
              :class="card.status === 'completed'
                ? 'b022hub-theme-chip'
                : 'b022hub-theme-badge'"
            >
              {{ card.status === 'completed' ? $t('home.progress.completed') : $t('home.progress.building') }}
            </p>
          </article>
        </div>

        <div class="flex flex-wrap items-center justify-center gap-4" :style="buildRevealStyle(0.32)">
          <a
            v-if="docHref"
            :href="docHref"
            target="_blank"
            rel="noopener noreferrer"
            class="b022hub-theme-outline-button inline-flex items-center gap-2 rounded-full px-6 py-3 text-sm font-medium transition"
          >
            <Icon name="book" size="sm" />
            <span>{{ docLabel }}</span>
          </a>

          <router-link
            :to="primaryPath"
            class="inline-flex items-center gap-2 rounded-full bg-[#cc785c] px-6 py-3 text-sm font-medium text-white shadow-[0_18px_45px_rgba(204,120,92,0.28)] transition hover:bg-[#d58e74]"
          >
            <Icon name="arrowRight" size="sm" />
            <span>{{ primaryLabel }}</span>
          </router-link>
        </div>
      </div>
    </div>
  </section>
</template>
