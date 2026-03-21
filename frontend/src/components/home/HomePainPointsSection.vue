<script setup lang="ts">
import Icon from '@/components/icons/Icon.vue'
import type { HomeInfoCard, HomeTone } from './types'

defineProps<{
  title: string
  description: string
  cards: HomeInfoCard[]
}>()

const toneCardClass: Record<HomeTone, string> = {
  primary:
    'border-primary-200/60 bg-primary-50/75 dark:border-primary-500/15 dark:bg-primary-500/10',
  sky: 'border-sky-200/70 bg-sky-50/70 dark:border-sky-500/15 dark:bg-sky-500/10',
  amber:
    'border-amber-200/70 bg-amber-50/70 dark:border-amber-500/15 dark:bg-amber-500/10',
  rose: 'border-rose-200/70 bg-rose-50/70 dark:border-rose-500/15 dark:bg-rose-500/10',
  violet:
    'border-violet-200/70 bg-violet-50/70 dark:border-violet-500/15 dark:bg-violet-500/10',
  slate:
    'border-slate-200/70 bg-slate-50/75 dark:border-white/10 dark:bg-white/5'
}

const toneIconClass: Record<HomeTone, string> = {
  primary: 'bg-primary-600 text-white shadow-primary-500/30',
  sky: 'bg-sky-600 text-white shadow-sky-500/30',
  amber: 'bg-amber-500 text-white shadow-amber-500/30',
  rose: 'bg-rose-500 text-white shadow-rose-500/30',
  violet: 'bg-violet-600 text-white shadow-violet-500/30',
  slate: 'bg-slate-700 text-white shadow-slate-500/20'
}
</script>

<template>
  <section class="grid gap-10 lg:grid-cols-[minmax(0,0.88fr)_minmax(0,1.12fr)]">
    <div class="space-y-5">
      <p class="text-xs font-semibold uppercase tracking-[0.24em] text-primary-700/80 dark:text-primary-300/80">
        Pain Points
      </p>
      <h2 class="font-display text-4xl font-semibold leading-tight text-gray-950 dark:text-white">
        {{ title }}
      </h2>
      <p class="max-w-xl text-base leading-8 text-gray-600 dark:text-dark-300">
        {{ description }}
      </p>
    </div>

    <div class="grid gap-4 sm:grid-cols-2">
      <article
        v-for="card in cards"
        :key="card.title"
        :class="[
          'group rounded-[1.75rem] border p-6 shadow-sm shadow-slate-900/5 transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:shadow-slate-900/10 dark:shadow-none',
          toneCardClass[card.tone]
        ]"
      >
        <div
          :class="[
            'mb-5 flex h-12 w-12 items-center justify-center rounded-2xl shadow-lg transition-transform duration-300 group-hover:scale-110',
            toneIconClass[card.tone]
          ]"
        >
          <Icon :name="card.icon" size="lg" />
        </div>

        <h3 class="mb-3 text-xl font-semibold text-gray-900 dark:text-white">
          {{ card.title }}
        </h3>
        <p class="text-sm leading-7 text-gray-600 dark:text-dark-300">
          {{ card.description }}
        </p>
      </article>
    </div>
  </section>
</template>
