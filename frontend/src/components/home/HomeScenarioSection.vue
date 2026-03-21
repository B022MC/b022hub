<script setup lang="ts">
import Icon from '@/components/icons/Icon.vue'
import type { HomeScenarioItem, HomeTone } from './types'

defineProps<{
  title: string
  subtitle: string
  items: HomeScenarioItem[]
}>()

const toneRibbonClass: Record<HomeTone, string> = {
  primary: 'bg-primary-500/12 text-primary-700 dark:bg-primary-500/15 dark:text-primary-200',
  sky: 'bg-sky-500/12 text-sky-700 dark:bg-sky-500/15 dark:text-sky-200',
  amber: 'bg-amber-500/12 text-amber-700 dark:bg-amber-500/15 dark:text-amber-200',
  rose: 'bg-rose-500/12 text-rose-700 dark:bg-rose-500/15 dark:text-rose-200',
  violet: 'bg-violet-500/12 text-violet-700 dark:bg-violet-500/15 dark:text-violet-200',
  slate: 'bg-slate-500/10 text-slate-700 dark:bg-slate-500/15 dark:text-slate-200'
}

const toneIconClass: Record<HomeTone, string> = {
  primary: 'bg-primary-600 text-white',
  sky: 'bg-sky-600 text-white',
  amber: 'bg-amber-500 text-white',
  rose: 'bg-rose-500 text-white',
  violet: 'bg-violet-600 text-white',
  slate: 'bg-slate-700 text-white'
}
</script>

<template>
  <section class="space-y-8">
    <div class="max-w-3xl space-y-4">
      <p class="text-xs font-semibold uppercase tracking-[0.24em] text-primary-700/80 dark:text-primary-300/80">
        Scenarios
      </p>
      <h2 class="font-display text-4xl font-semibold leading-tight text-gray-950 dark:text-white">
        {{ title }}
      </h2>
      <p class="text-base leading-8 text-gray-600 dark:text-dark-300">
        {{ subtitle }}
      </p>
    </div>

    <div class="grid gap-5 lg:grid-cols-3">
      <article
        v-for="item in items"
        :key="item.title"
        class="rounded-[1.8rem] border border-white/70 bg-white/80 p-6 shadow-sm shadow-slate-900/5 backdrop-blur dark:border-white/10 dark:bg-white/5 dark:shadow-none"
      >
        <div class="mb-5 flex items-center justify-between gap-4">
          <span
            :class="[
              'rounded-full px-3 py-1.5 text-xs font-semibold uppercase tracking-[0.16em]',
              toneRibbonClass[item.tone]
            ]"
          >
            {{ item.eyebrow }}
          </span>

          <div
            :class="[
              'flex h-11 w-11 items-center justify-center rounded-2xl shadow-lg',
              toneIconClass[item.tone]
            ]"
          >
            <Icon :name="item.icon" size="md" />
          </div>
        </div>

        <h3 class="mb-3 font-display text-2xl font-semibold text-gray-950 dark:text-white">
          {{ item.title }}
        </h3>
        <p class="mb-5 text-sm leading-7 text-gray-600 dark:text-dark-300">
          {{ item.description }}
        </p>

        <ul class="space-y-3">
          <li
            v-for="point in item.points"
            :key="point"
            class="flex items-start gap-3 text-sm leading-7 text-gray-700 dark:text-dark-200"
          >
            <Icon name="checkCircle" size="sm" class="mt-1 shrink-0 text-primary-500 dark:text-primary-300" />
            <span>{{ point }}</span>
          </li>
        </ul>
      </article>
    </div>
  </section>
</template>
