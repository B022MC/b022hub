<script setup lang="ts">
import Icon from '@/components/icons/Icon.vue'
import type { HomeComparisonRow } from './types'

defineProps<{
  title: string
  featureHeader: string
  officialHeader: string
  platformHeader: string
  rows: HomeComparisonRow[]
}>()
</script>

<template>
  <section class="space-y-8">
    <div class="max-w-3xl space-y-4">
      <p class="text-xs font-semibold uppercase tracking-[0.24em] text-primary-700/80 dark:text-primary-300/80">
        Comparison
      </p>
      <h2 class="font-display text-4xl font-semibold leading-tight text-gray-950 dark:text-white">
        {{ title }}
      </h2>
    </div>

    <div class="overflow-hidden rounded-[2rem] border border-white/70 bg-white/85 shadow-sm shadow-slate-900/5 backdrop-blur dark:border-white/10 dark:bg-white/5 dark:shadow-none">
      <div class="grid grid-cols-[1.05fr_1fr_1fr] gap-px bg-slate-200/70 dark:bg-white/10">
        <div class="bg-slate-100/80 px-5 py-4 text-sm font-semibold text-gray-600 dark:bg-white/5 dark:text-dark-200">
          {{ featureHeader }}
        </div>
        <div class="bg-slate-100/80 px-5 py-4 text-sm font-semibold text-gray-600 dark:bg-white/5 dark:text-dark-200">
          {{ officialHeader }}
        </div>
        <div class="bg-primary-50 px-5 py-4 text-sm font-semibold text-primary-700 dark:bg-primary-500/10 dark:text-primary-200">
          {{ platformHeader }}
        </div>

        <template v-for="row in rows" :key="row.feature">
          <div class="bg-white/90 px-5 py-5 text-sm font-semibold text-gray-900 dark:bg-white/[0.03] dark:text-white">
            {{ row.feature }}
          </div>
          <div class="bg-white/90 px-5 py-5 text-sm leading-7 text-gray-600 dark:bg-white/[0.03] dark:text-dark-300">
            {{ row.official }}
          </div>
          <div class="bg-primary-50/75 px-5 py-5 text-sm leading-7 text-primary-900 dark:bg-primary-500/10 dark:text-primary-100">
            <div class="flex items-start gap-3">
              <Icon name="checkCircle" size="sm" class="mt-1 shrink-0 text-primary-600 dark:text-primary-300" />
              <span>{{ row.platform }}</span>
            </div>
          </div>
        </template>
      </div>
    </div>
  </section>
</template>
