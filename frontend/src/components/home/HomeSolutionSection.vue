<script setup lang="ts">
import Icon from '@/components/icons/Icon.vue'
import type { HomeInfoCard, HomeStepItem, HomeTone } from './types'

defineProps<{
  title: string
  subtitle: string
  stepTitle: string
  features: HomeInfoCard[]
  steps: HomeStepItem[]
}>()

const toneAccentClass: Record<HomeTone, string> = {
  primary: 'from-primary-500 to-cyan-400',
  sky: 'from-sky-500 to-cyan-300',
  amber: 'from-amber-500 to-orange-300',
  rose: 'from-rose-500 to-pink-300',
  violet: 'from-violet-500 to-fuchsia-300',
  slate: 'from-slate-700 to-slate-400'
}
</script>

<template>
  <section class="space-y-8">
    <div class="max-w-3xl space-y-4">
      <p class="text-xs font-semibold uppercase tracking-[0.24em] text-primary-700/80 dark:text-primary-300/80">
        Solutions
      </p>
      <h2 class="font-display text-4xl font-semibold leading-tight text-gray-950 dark:text-white">
        {{ title }}
      </h2>
      <p class="text-base leading-8 text-gray-600 dark:text-dark-300">
        {{ subtitle }}
      </p>
    </div>

    <div class="grid gap-6 xl:grid-cols-[minmax(0,1.05fr)_minmax(0,0.95fr)]">
      <div class="grid gap-5 md:grid-cols-3">
        <article
          v-for="feature in features"
          :key="feature.title"
          class="group rounded-[1.8rem] border border-white/70 bg-white/80 p-6 shadow-sm shadow-slate-900/5 backdrop-blur transition-all duration-300 hover:-translate-y-1 hover:shadow-[0_24px_50px_rgba(15,23,42,0.1)] dark:border-white/10 dark:bg-white/5 dark:shadow-none"
        >
          <div
            :class="[
              'mb-5 flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br text-white shadow-lg transition-transform duration-300 group-hover:scale-110',
              toneAccentClass[feature.tone]
            ]"
          >
            <Icon :name="feature.icon" size="lg" />
          </div>
          <h3 class="mb-3 text-xl font-semibold text-gray-900 dark:text-white">
            {{ feature.title }}
          </h3>
          <p class="text-sm leading-7 text-gray-600 dark:text-dark-300">
            {{ feature.description }}
          </p>
        </article>
      </div>

      <aside class="rounded-[2rem] border border-slate-200/70 bg-slate-950 px-6 py-7 text-white shadow-[0_24px_60px_rgba(15,23,42,0.18)] dark:border-white/10 sm:px-8">
        <div class="mb-6 space-y-2">
          <p class="text-xs font-semibold uppercase tracking-[0.24em] text-primary-300/80">
            Workflow
          </p>
          <h3 class="font-display text-2xl font-semibold text-white">
            {{ stepTitle }}
          </h3>
        </div>

        <div class="space-y-5">
          <div
            v-for="step in steps"
            :key="step.index"
            class="flex gap-4 rounded-[1.5rem] border border-white/10 bg-white/[0.04] p-4"
          >
            <div
              class="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-white/10 text-sm font-semibold text-white"
            >
              {{ step.index }}
            </div>

            <div class="space-y-2">
              <div class="flex items-center gap-2">
                <Icon :name="step.icon" size="sm" class="text-primary-300" />
                <p class="text-base font-semibold text-white">
                  {{ step.title }}
                </p>
              </div>
              <p class="text-sm leading-7 text-slate-300">
                {{ step.description }}
              </p>
            </div>
          </div>
        </div>
      </aside>
    </div>
  </section>
</template>
