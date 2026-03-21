<script setup lang="ts">
import { computed } from 'vue'
import Icon from '@/components/icons/Icon.vue'
import type { HomeProviderItem, HomeTagItem, HomeTone } from './types'

const props = defineProps<{
  siteName: string
  siteSubtitle: string
  heroEyebrow: string
  heroDescription: string
  primaryCtaPath: string
  primaryCtaLabel: string
  secondaryCtaHref: string
  secondaryCtaLabel: string
  tags: HomeTagItem[]
  providers: HomeProviderItem[]
  showcaseTitle: string
  showcaseRequestLabel: string
  showcaseRoutingLabel: string
  showcasePolicyLabel: string
  showcaseResponseLabel: string
  liveStatusLabel: string
}>()

const providerPreview = computed(() => props.providers.slice(0, 4))

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

const toneDotClass: Record<HomeTone, string> = {
  primary: 'bg-primary-500',
  sky: 'bg-sky-500',
  amber: 'bg-amber-500',
  rose: 'bg-rose-500',
  violet: 'bg-violet-500',
  slate: 'bg-slate-400'
}
</script>

<template>
  <section class="grid gap-12 lg:grid-cols-[minmax(0,1.02fr)_minmax(0,0.98fr)] lg:items-center">
    <div class="space-y-8">
      <div
        class="inline-flex items-center gap-2 rounded-full border border-white/70 bg-white/85 px-4 py-2 text-xs font-semibold uppercase tracking-[0.22em] text-gray-600 shadow-sm shadow-primary-500/10 backdrop-blur dark:border-white/10 dark:bg-white/5 dark:text-dark-200"
      >
        <span class="h-2 w-2 rounded-full bg-primary-500 shadow-[0_0_18px_rgba(204,120,92,0.55)]"></span>
        {{ heroEyebrow }}
      </div>

      <div class="space-y-5">
        <div class="space-y-2">
          <p class="font-display text-[clamp(2.8rem,8vw,5.6rem)] font-bold leading-none text-gray-950 dark:text-white">
            {{ siteName }}
          </p>
          <h1
            class="max-w-3xl font-display text-[clamp(2rem,5vw,3.35rem)] font-semibold leading-[1.04] text-gray-800 dark:text-dark-50"
          >
            {{ siteSubtitle }}
          </h1>
        </div>

        <p class="max-w-2xl text-lg leading-8 text-gray-600 dark:text-dark-300">
          {{ heroDescription }}
        </p>
      </div>

      <div class="flex flex-wrap items-center gap-3">
        <router-link :to="primaryCtaPath" class="btn btn-primary rounded-full px-7 py-3 text-base">
          {{ primaryCtaLabel }}
          <Icon name="arrowRight" size="md" :stroke-width="2" />
        </router-link>

        <a
          :href="secondaryCtaHref"
          target="_blank"
          rel="noopener noreferrer"
          class="btn btn-secondary rounded-full px-6 py-3 text-base"
        >
          <Icon name="book" size="md" />
          {{ secondaryCtaLabel }}
        </a>
      </div>

      <div class="flex flex-wrap gap-3">
        <div
          v-for="tag in tags"
          :key="tag.label"
          class="inline-flex items-center gap-2 rounded-full border border-white/70 bg-white/75 px-4 py-2 text-sm font-medium text-gray-700 shadow-sm shadow-primary-500/5 backdrop-blur dark:border-white/10 dark:bg-white/5 dark:text-dark-100"
        >
          <Icon :name="tag.icon" size="sm" class="text-primary-600 dark:text-primary-300" />
          <span>{{ tag.label }}</span>
        </div>
      </div>
    </div>

    <div class="home-showcase rounded-[2rem] border border-white/70 bg-white/80 p-5 shadow-[0_30px_90px_rgba(15,23,42,0.14)] backdrop-blur-xl dark:border-white/10 dark:bg-[#071118]/88 dark:shadow-[0_30px_120px_rgba(0,0,0,0.45)] sm:p-8">
      <div class="home-showcase-noise"></div>

      <div class="relative z-10 space-y-5">
        <div class="flex flex-wrap items-start justify-between gap-4">
          <div class="space-y-2">
            <p class="text-xs font-semibold uppercase tracking-[0.22em] text-primary-700/80 dark:text-primary-300/80">
              {{ showcaseTitle }}
            </p>
            <div class="space-y-1">
              <h2 class="font-display text-2xl font-semibold text-gray-950 dark:text-white">
                {{ siteName }} Runtime
              </h2>
              <p class="text-sm text-gray-500 dark:text-dark-400">
                /v1/messages · /v1/chat/completions
              </p>
            </div>
          </div>

          <div class="inline-flex items-center gap-2 rounded-full border border-emerald-500/20 bg-emerald-500/12 px-3 py-1.5 text-sm font-medium text-emerald-700 dark:border-emerald-400/20 dark:bg-emerald-500/15 dark:text-emerald-200">
            <span class="home-live-dot h-2 w-2 rounded-full bg-emerald-400"></span>
            {{ liveStatusLabel }}
          </div>
        </div>

        <div class="grid gap-4">
          <div class="rounded-[1.75rem] border border-white/70 bg-white/85 p-4 shadow-sm shadow-slate-900/5 dark:border-white/10 dark:bg-white/5">
            <p class="mb-3 text-xs font-semibold uppercase tracking-[0.18em] text-gray-500 dark:text-dark-400">
              {{ showcaseRequestLabel }}
            </p>
            <div class="rounded-2xl bg-[#05131a] p-4 font-mono text-sm text-slate-100">
              <div class="flex flex-wrap items-center gap-2">
                <span class="text-primary-300">$</span>
                <span class="text-white">curl</span>
                <span class="text-sky-300">-X POST</span>
                <span class="text-emerald-300">/v1/messages</span>
              </div>
              <div class="mt-3 text-slate-400">
                model=auto · account_pool=smart · retry_policy=sticky
              </div>
            </div>
          </div>

          <div class="grid gap-4 lg:grid-cols-[1.18fr_0.82fr]">
            <div class="rounded-[1.75rem] border border-white/70 bg-white/85 p-4 shadow-sm shadow-slate-900/5 dark:border-white/10 dark:bg-white/5">
              <p class="mb-3 text-xs font-semibold uppercase tracking-[0.18em] text-gray-500 dark:text-dark-400">
                {{ showcaseRoutingLabel }}
              </p>

              <div class="space-y-3">
                <div
                  v-for="provider in providerPreview"
                  :key="provider.name"
                  class="flex items-center justify-between gap-4 rounded-2xl border border-slate-200/70 bg-slate-900/[0.03] px-4 py-3 dark:border-white/10 dark:bg-white/[0.04]"
                >
                  <div class="flex items-center gap-3">
                    <div
                      :class="[
                        'flex h-11 w-11 items-center justify-center rounded-2xl border text-sm font-semibold shadow-sm',
                        toneBadgeClass[provider.tone]
                      ]"
                    >
                      {{ provider.shortLabel }}
                    </div>
                    <div>
                      <p class="text-sm font-semibold text-gray-900 dark:text-white">
                        {{ provider.name }}
                      </p>
                      <p class="text-xs text-gray-500 dark:text-dark-400">
                        {{ provider.detail }}
                      </p>
                    </div>
                  </div>

                  <div class="flex items-center gap-2 text-xs font-medium text-gray-500 dark:text-dark-300">
                    <span :class="['h-2.5 w-2.5 rounded-full', toneDotClass[provider.tone]]"></span>
                    <span>{{ provider.status === 'supported' ? 'LIVE' : 'NEXT' }}</span>
                  </div>
                </div>
              </div>
            </div>

            <div class="rounded-[1.75rem] border border-white/70 bg-white/85 p-4 shadow-sm shadow-slate-900/5 dark:border-white/10 dark:bg-white/5">
              <p class="mb-3 text-xs font-semibold uppercase tracking-[0.18em] text-gray-500 dark:text-dark-400">
                {{ showcasePolicyLabel }}
              </p>

              <div class="space-y-3">
                <div
                  v-for="tag in tags"
                  :key="tag.label"
                  class="rounded-2xl border border-slate-200/70 bg-slate-900/[0.03] px-4 py-3 dark:border-white/10 dark:bg-white/[0.04]"
                >
                  <div class="mb-1 flex items-center gap-2">
                    <Icon :name="tag.icon" size="sm" class="text-primary-600 dark:text-primary-300" />
                    <p class="text-sm font-semibold text-gray-900 dark:text-white">
                      {{ tag.label }}
                    </p>
                  </div>
                  <div class="h-2 overflow-hidden rounded-full bg-slate-200 dark:bg-white/10">
                    <div class="home-meter h-full rounded-full bg-gradient-to-r from-primary-500 via-cyan-400 to-emerald-400"></div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="rounded-[1.75rem] border border-white/70 bg-white/85 p-4 shadow-sm shadow-slate-900/5 dark:border-white/10 dark:bg-white/5">
            <p class="mb-3 text-xs font-semibold uppercase tracking-[0.18em] text-gray-500 dark:text-dark-400">
              {{ showcaseResponseLabel }}
            </p>
            <div class="rounded-2xl bg-[#081018] p-4 font-mono text-sm text-slate-100">
              <div class="flex items-center gap-3 text-emerald-300">
                <span class="rounded-full bg-emerald-500/15 px-2 py-0.5 text-[11px] font-semibold uppercase tracking-[0.12em] text-emerald-200">
                  200 OK
                </span>
                <span class="text-slate-400">latency 418ms</span>
              </div>
              <div class="mt-4 space-y-2 text-slate-300">
                <p>{</p>
                <p class="pl-4">&quot;model&quot;: &quot;smart-route&quot;,</p>
                <p class="pl-4">&quot;provider&quot;: &quot;{{ providerPreview[0]?.name || siteName }}&quot;,</p>
                <p class="pl-4 text-primary-200">&quot;content&quot;: &quot;Route stable. Usage accounted.&quot;</p>
                <p>}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.home-showcase {
  position: relative;
  isolation: isolate;
}

.home-showcase::before {
  position: absolute;
  inset: 14px;
  border-radius: 1.65rem;
  background:
    radial-gradient(circle at top, rgba(204, 120, 92, 0.16), transparent 40%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.42), rgba(255, 255, 255, 0.06));
  content: '';
  pointer-events: none;
}

.home-showcase-noise {
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.08) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.08) 1px, transparent 1px);
  background-position: center;
  background-size: 26px 26px;
  mask-image: linear-gradient(180deg, rgba(0, 0, 0, 0.8), transparent 92%);
  opacity: 0.26;
  pointer-events: none;
}

.home-live-dot {
  animation: pulse-live 1.8s ease-in-out infinite;
}

.home-meter {
  width: 72%;
  animation: flow-meter 4.8s ease-in-out infinite;
}

@keyframes pulse-live {
  0%,
  100% {
    opacity: 1;
    transform: scale(1);
  }

  50% {
    opacity: 0.45;
    transform: scale(0.82);
  }
}

@keyframes flow-meter {
  0%,
  100% {
    width: 64%;
  }

  50% {
    width: 88%;
  }
}
</style>
