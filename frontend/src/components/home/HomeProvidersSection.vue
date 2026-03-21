<script setup lang="ts">
import type { HomeProviderItem, HomeTone } from './types'

defineProps<{
  title: string
  description: string
  supportedLabel: string
  soonLabel: string
  providers: HomeProviderItem[]
}>()

const toneBadgeClass: Record<HomeTone, string> = {
  primary:
    'border-primary-500/20 bg-primary-500/12 text-primary-700 dark:border-primary-400/20 dark:bg-primary-500/18 dark:text-primary-200',
  sky: 'border-sky-500/20 bg-sky-500/12 text-sky-700 dark:border-sky-400/20 dark:bg-sky-500/18 dark:text-sky-200',
  amber:
    'border-amber-500/20 bg-amber-500/12 text-amber-700 dark:border-amber-400/20 dark:bg-amber-500/18 dark:text-amber-200',
  rose: 'border-rose-500/20 bg-rose-500/12 text-rose-700 dark:border-rose-400/20 dark:bg-rose-500/18 dark:text-rose-200',
  violet:
    'border-violet-500/20 bg-violet-500/12 text-violet-700 dark:border-violet-400/20 dark:bg-violet-500/18 dark:text-violet-200',
  slate:
    'border-slate-400/20 bg-slate-500/10 text-slate-700 dark:border-slate-400/20 dark:bg-slate-500/15 dark:text-slate-200'
}
</script>

<template>
  <section class="space-y-8">
    <div class="max-w-3xl space-y-4">
      <p class="text-xs font-semibold uppercase tracking-[0.24em] text-primary-700/80 dark:text-primary-300/80">
        Providers
      </p>
      <h2 class="font-display text-4xl font-semibold leading-tight text-gray-950 dark:text-white">
        {{ title }}
      </h2>
      <p class="text-base leading-8 text-gray-600 dark:text-dark-300">
        {{ description }}
      </p>
    </div>

    <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-5">
      <article
        v-for="provider in providers"
        :key="provider.name"
        class="rounded-[1.75rem] border border-white/70 bg-white/80 p-5 shadow-sm shadow-slate-900/5 backdrop-blur dark:border-white/10 dark:bg-white/5 dark:shadow-none"
      >
        <div class="mb-4 flex items-center justify-between gap-3">
          <div
            :class="[
              'flex h-12 w-12 items-center justify-center rounded-2xl border text-sm font-semibold',
              toneBadgeClass[provider.tone]
            ]"
          >
            {{ provider.shortLabel }}
          </div>

          <span
            class="rounded-full px-2.5 py-1 text-[11px] font-semibold uppercase tracking-[0.14em]"
            :class="provider.status === 'supported'
              ? 'bg-primary-500/12 text-primary-700 dark:bg-primary-500/15 dark:text-primary-200'
              : 'bg-slate-200/80 text-slate-600 dark:bg-white/10 dark:text-dark-200'"
          >
            {{ provider.status === 'supported' ? supportedLabel : soonLabel }}
          </span>
        </div>

        <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
          {{ provider.name }}
        </h3>
        <p class="text-sm leading-7 text-gray-600 dark:text-dark-300">
          {{ provider.detail }}
        </p>
      </article>
    </div>
  </section>
</template>
