<script setup lang="ts">
import Icon from '@/components/icons/Icon.vue'
import type { HomeTagItem } from './types'

defineProps<{
  title: string
  description: string
  primaryCtaPath: string
  primaryCtaLabel: string
  secondaryCtaHref: string
  secondaryCtaLabel: string
  tags: HomeTagItem[]
}>()
</script>

<template>
  <section class="overflow-hidden rounded-[2.25rem] border border-white/70 bg-[linear-gradient(135deg,#05262a_0%,#0d3b42_38%,#111827_100%)] px-6 py-10 text-white shadow-[0_30px_120px_rgba(15,23,42,0.22)] dark:border-white/10 sm:px-10">
    <div class="grid gap-10 lg:grid-cols-[minmax(0,1fr)_auto] lg:items-center">
      <div class="space-y-5">
        <p class="text-xs font-semibold uppercase tracking-[0.24em] text-primary-200/85">
          Ready
        </p>
        <h2 class="max-w-3xl font-display text-4xl font-semibold leading-tight text-white">
          {{ title }}
        </h2>
        <p class="max-w-2xl text-base leading-8 text-slate-200">
          {{ description }}
        </p>

        <div class="flex flex-wrap gap-3">
          <div
            v-for="tag in tags"
            :key="tag.label"
            class="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/[0.06] px-4 py-2 text-sm font-medium text-slate-100"
          >
            <Icon :name="tag.icon" size="sm" class="text-primary-300" />
            <span>{{ tag.label }}</span>
          </div>
        </div>
      </div>

      <div class="flex flex-wrap items-center gap-3 lg:justify-end">
        <router-link :to="primaryCtaPath" class="btn rounded-full bg-white px-6 py-3 text-base font-semibold text-slate-950 hover:bg-slate-100">
          {{ primaryCtaLabel }}
          <Icon name="arrowRight" size="md" />
        </router-link>

        <a
          :href="secondaryCtaHref"
          target="_blank"
          rel="noopener noreferrer"
          class="btn rounded-full border border-white/15 bg-white/[0.06] px-6 py-3 text-base font-medium text-white hover:bg-white/[0.12]"
        >
          <Icon name="book" size="md" />
          {{ secondaryCtaLabel }}
        </a>
      </div>
    </div>
  </section>
</template>
