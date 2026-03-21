<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, shallowRef, type ComponentPublicInstance } from 'vue'
import { useI18n } from 'vue-i18n'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { useClipboard } from '@/composables/useClipboard'
import B022HubLogoStage from './B022HubLogoStage.vue'
import B022HubProgressSection from './B022HubProgressSection.vue'
import B022HubToolSection from './B022HubToolSection.vue'
import type { B022HubFeatureCard, B022HubInstallOption } from './b022hub'

const props = defineProps<{
  siteName: string
  siteLogo: string
  siteSubtitle: string
  docUrl: string
  apiBaseUrl: string
  isDark: boolean
  isAuthenticated: boolean
  dashboardPath: string
  primaryPath: string
  primaryLabel: string
}>()

const emit = defineEmits<{
  toggleTheme: []
}>()

const { t } = useI18n()
const { copyToClipboard } = useClipboard()

const scrollContainer = ref<HTMLElement | null>(null)
const sectionRefs = ref<HTMLElement[]>([])
const activeSectionIndex = shallowRef(0)
const sectionProgresses = ref<number[]>([1, 0, 0, 0, 0])
const typewrittenBrand = shallowRef('')
const showCursor = shallowRef(true)
const origin = shallowRef('')
const claudeInstallKey = shallowRef('mac')
const codexInstallKey = shallowRef('node')
const geminiInstallKey = shallowRef('node')

let scrollRafId = 0
let cursorTimer: ReturnType<typeof setInterval> | null = null
let typeTimer: ReturnType<typeof setTimeout> | null = null
let deleteTimer: ReturnType<typeof setTimeout> | null = null
let pauseTimer: ReturnType<typeof setTimeout> | null = null
let startTimer: ReturnType<typeof setTimeout> | null = null

const fallbackMark = computed(() => 'b022')

const navItems = computed(() => [
  { id: 'home', label: t('home.nav.home') },
  { id: 'claude', label: t('home.nav.claude') },
  { id: 'codex', label: t('home.nav.codex') },
  { id: 'gemini', label: t('home.nav.gemini') },
  { id: 'progress', label: t('home.nav.more') }
])

const activeSectionId = computed(() => navItems.value[activeSectionIndex.value]?.id || 'home')
const runtimeBaseUrl = computed(() => props.apiBaseUrl.trim() || origin.value || '')
const docsHref = computed(() => props.docUrl.trim())

const claudeInstallOptions = computed<B022HubInstallOption[]>(() => [
  {
    key: 'mac',
    label: 'Mac / Linux',
    hint: 'Terminal',
    command: 'curl -fsSL https://claude.ai/install.sh | bash'
  }
])

const codexInstallOptions = computed<B022HubInstallOption[]>(() => [
  {
    key: 'node',
    label: 'Node.js',
    hint: 'npm',
    command: 'npm install -g @openai/codex'
  },
  {
    key: 'brew',
    label: 'Mac',
    hint: 'Homebrew',
    command: 'brew install --cask codex'
  }
])

const geminiInstallOptions = computed<B022HubInstallOption[]>(() => [
  {
    key: 'node',
    label: 'Node.js',
    hint: 'npm',
    command: 'npm install -g @google/gemini-cli'
  },
  {
    key: 'brew',
    label: 'Mac',
    hint: 'Homebrew',
    command: 'brew install gemini-cli'
  }
])

const claudeConfig = computed(() => `{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "your-api-key",
    "ANTHROPIC_BASE_URL": "${runtimeBaseUrl.value}"
  }
}`)

const codexConfig = computed(() => `model_provider = "b022hub"
model = "latest-model-name"
model_reasoning_effort = "high"
network_access = "enabled"
disable_response_storage = true

[model_providers.b022hub]
name = "b022hub"
base_url = "${runtimeBaseUrl.value}/v1"
wire_api = "responses"
requires_openai_auth = true`)

const codexAuth = computed(() => `{
  "OPENAI_API_KEY": "your-api-key"
}`)

const geminiEnv = computed(() => `GOOGLE_GEMINI_BASE_URL=${runtimeBaseUrl.value}
GEMINI_API_KEY=your-api-key
GEMINI_MODEL=latest-model-name`)

const geminiSettings = computed(() => `{
  "ide": {
    "enabled": true
  },
  "security": {
    "auth": {
      "selectedType": "gemini-api-key"
    }
  }
}`)

const progressCards = computed<B022HubFeatureCard[]>(() => [
  {
    icon: 'terminal',
    title: t('home.progress.cards.gateway.title'),
    description: t('home.progress.cards.gateway.description'),
    status: 'completed'
  },
  {
    icon: 'swap',
    title: t('home.progress.cards.routing.title'),
    description: t('home.progress.cards.routing.description'),
    status: 'completed'
  },
  {
    icon: 'cpu',
    title: t('home.progress.cards.workflow.title'),
    description: t('home.progress.cards.workflow.description'),
    status: 'completed'
  }
])

function registerSection(index: number, element: Element | ComponentPublicInstance | null) {
  if (element instanceof HTMLElement) {
    sectionRefs.value[index] = element
    return
  }

  if (element && '$el' in element && element.$el instanceof HTMLElement) {
    sectionRefs.value[index] = element.$el
  }
}

function scrollToSection(index: number) {
  sectionRefs.value[index]?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

function updateScrollState() {
  scrollRafId = 0

  if (!sectionRefs.value.length) {
    return
  }

  const viewportCenter = window.innerHeight / 2
  const progresses = sectionRefs.value.map((section) => {
    const rect = section.getBoundingClientRect()
    const center = rect.top + rect.height / 2
    const distance = Math.abs(center - viewportCenter)
    const maxDistance = viewportCenter + rect.height / 2
    return Math.max(0, Math.min(1, 1 - distance / maxDistance))
  })

  sectionProgresses.value = progresses

  const maxProgress = Math.max(...progresses)
  const nextIndex = progresses.findIndex((progress) => progress === maxProgress)

  if (nextIndex >= 0 && nextIndex !== activeSectionIndex.value) {
    activeSectionIndex.value = nextIndex

    if (nextIndex === 0) {
      startTypewriter()
    }
  }
}

function requestScrollUpdate() {
  if (scrollRafId) {
    return
  }

  scrollRafId = window.requestAnimationFrame(updateScrollState)
}

function clearTypewriterTimers() {
  if (cursorTimer) {
    clearInterval(cursorTimer)
    cursorTimer = null
  }
  if (typeTimer) {
    clearTimeout(typeTimer)
    typeTimer = null
  }
  if (deleteTimer) {
    clearTimeout(deleteTimer)
    deleteTimer = null
  }
  if (pauseTimer) {
    clearTimeout(pauseTimer)
    pauseTimer = null
  }
  if (startTimer) {
    clearTimeout(startTimer)
    startTimer = null
  }
}

function startTypewriter() {
  clearTypewriterTimers()
  typewrittenBrand.value = ''
  showCursor.value = true

  cursorTimer = setInterval(() => {
    showCursor.value = !showCursor.value
  }, 600)

  const brand = props.siteName

  const typeStep = (index: number) => {
    if (index <= brand.length) {
      typewrittenBrand.value = brand.slice(0, index)
      typeTimer = setTimeout(() => typeStep(index + 1), 140)
      return
    }

    pauseTimer = setTimeout(() => deleteStep(brand.length - 1), 1600)
  }

  const deleteStep = (index: number) => {
    if (activeSectionIndex.value !== 0) {
      return
    }

    if (index >= 0) {
      typewrittenBrand.value = brand.slice(0, index)
      deleteTimer = setTimeout(() => deleteStep(index - 1), 80)
      return
    }

    startTimer = setTimeout(() => typeStep(1), 700)
  }

  startTimer = setTimeout(() => typeStep(1), 250)
}

function heroStyle(offset: number, axisX = 0) {
  const progress = sectionProgresses.value[0] || 0
  const local = Math.min(1, Math.max(0, (progress - offset) / (1 - offset || 1)))
  const opacity = 0.15 + local * 0.85
  const y = (1 - local) * 32
  const x = (1 - local) * axisX

  return {
    opacity,
    transform: `translate3d(${x}px, ${y}px, 0)`
  }
}

async function copyConfig(value: string) {
  await copyToClipboard(value)
}

function onResize() {
  requestScrollUpdate()
}

onMounted(() => {
  origin.value = window.location.origin
  startTypewriter()
  scrollContainer.value?.addEventListener('scroll', requestScrollUpdate, { passive: true })
  window.addEventListener('resize', onResize, { passive: true })
  requestScrollUpdate()
})

onBeforeUnmount(() => {
  clearTypewriterTimers()
  if (scrollRafId) {
    cancelAnimationFrame(scrollRafId)
    scrollRafId = 0
  }
  scrollContainer.value?.removeEventListener('scroll', requestScrollUpdate)
  window.removeEventListener('resize', onResize)
})
</script>

<template>
  <div
    ref="scrollContainer"
    class="b022hub-shell b022hub-theme relative h-screen overflow-y-auto overflow-x-hidden snap-y snap-mandatory scroll-smooth"
    :class="{ 'b022hub-theme--light': !isDark }"
  >
    <div class="b022hub-grid"></div>
    <div class="b022hub-noise"></div>

    <B022HubLogoStage
      :site-name="siteName"
      :site-logo="siteLogo"
      :active-section-id="activeSectionId"
      :vendor-progress="sectionProgresses[4] || 0"
    />

    <nav class="fixed right-8 top-1/2 z-30 hidden -translate-y-1/2 flex-col gap-4 xl:flex">
      <button
        v-for="(item, index) in navItems"
        :key="item.id"
        class="group relative flex items-center justify-end"
        @click="scrollToSection(index)"
      >
        <span
          class="pointer-events-none absolute right-5 whitespace-nowrap rounded-full bg-[rgba(18,17,15,0.86)] px-3 py-1 text-xs text-[#b7b1a5] opacity-0 transition group-hover:opacity-100"
        >
          {{ item.label }}
        </span>
        <span
          class="h-3 w-3 rounded-full border-2 transition"
          :class="activeSectionIndex === index
            ? 'border-[#d4a27f] bg-[#cc785c] scale-125'
            : 'border-[rgba(227,224,211,0.32)] bg-transparent'"
        ></span>
      </button>
    </nav>

    <header class="b022hub-theme-header sticky top-0 z-40">
      <div class="mx-auto flex h-16 max-w-7xl items-center justify-between gap-4 px-4 sm:px-6 lg:px-8">
        <button class="flex items-center gap-3" @click="scrollToSection(0)">
          <div class="flex h-10 w-10 items-center justify-center rounded-2xl border border-[rgba(227,224,211,0.14)] bg-[rgba(255,255,255,0.04)]">
            <img v-if="siteLogo" :src="siteLogo" :alt="siteName" class="h-7 w-7 object-contain" />
            <span v-else class="b022hub-theme-text-strong text-[10px] font-semibold tracking-[-0.08em]">{{ fallbackMark }}</span>
          </div>
          <div class="hidden text-left sm:block">
            <p class="b022hub-theme-text-strong font-display text-xl">
              {{ siteName }}
            </p>
            <p class="b022hub-theme-text-muted text-[11px] uppercase tracking-[0.22em]">
              {{ siteSubtitle }}
            </p>
          </div>
        </button>

        <div class="hidden items-center gap-8 lg:flex">
          <button
            v-for="(item, index) in navItems"
            :key="item.id"
            class="relative text-sm transition"
            :class="activeSectionIndex === index
              ? 'b022hub-theme-text-accent'
              : 'b022hub-theme-text-body hover:text-[var(--b022-text-strong)]'"
            @click="scrollToSection(index)"
          >
            {{ item.label }}
            <span
              class="absolute -bottom-[22px] left-0 right-0 h-px bg-[#cc785c] transition"
              :class="activeSectionIndex === index ? 'scale-x-100' : 'scale-x-0'"
            ></span>
          </button>
        </div>

        <div class="flex items-center gap-2">
          <LocaleSwitcher />

          <button
            class="b022hub-theme-icon-button inline-flex h-10 w-10 items-center justify-center rounded-full transition"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="emit('toggleTheme')"
          >
            <Icon :name="isDark ? 'sun' : 'moon'" size="md" />
          </button>

          <router-link
            :to="isAuthenticated ? dashboardPath : primaryPath"
            class="inline-flex items-center gap-2 rounded-full bg-[#cc785c] px-4 py-2 text-sm font-medium text-white shadow-[0_18px_45px_rgba(204,120,92,0.28)] transition hover:bg-[#d58e74]"
          >
            <Icon :name="isAuthenticated ? 'grid' : 'login'" size="sm" />
            <span>{{ isAuthenticated ? t('home.dashboard') : primaryLabel }}</span>
          </router-link>
        </div>
      </div>
    </header>

    <main class="relative z-20">
      <section
        :ref="(element) => registerSection(0, element)"
        class="min-h-screen snap-start px-5 py-16 sm:px-8 md:px-14 lg:px-20"
      >
        <div class="mx-auto flex min-h-[calc(100vh-8rem)] max-w-7xl items-center">
          <div class="max-w-3xl space-y-7">
            <div
              class="b022hub-theme-chip inline-flex items-center gap-2 rounded-full px-4 py-2 text-xs font-medium uppercase tracking-[0.18em]"
              :style="heroStyle(0.05, -18)"
            >
              <Icon name="sparkles" size="sm" />
              <span>{{ t('home.heroEyebrow') }}</span>
            </div>

            <div class="space-y-4">
              <h1
                class="b022hub-theme-text-strong font-display text-[clamp(3rem,8vw,6rem)] leading-[0.95]"
                :style="heroStyle(0.1, -28)"
              >
                <span class="block">{{ t('home.heroPrefix') }}</span>
                <span class="b022hub-theme-text-accent">
                  {{ typewrittenBrand }}<span v-if="showCursor" class="type-cursor">_</span>
                </span>
              </h1>

              <p class="b022hub-theme-text-body max-w-2xl text-lg leading-8" :style="heroStyle(0.18, -16)">
                {{ t('home.heroDescription') }}
              </p>
            </div>

            <div class="flex flex-wrap items-center gap-4" :style="heroStyle(0.24, -10)">
              <router-link
                :to="primaryPath"
                class="inline-flex items-center gap-2 rounded-full bg-[#cc785c] px-6 py-3 text-sm font-medium text-white shadow-[0_18px_45px_rgba(204,120,92,0.28)] transition hover:bg-[#d58e74]"
              >
                <Icon name="arrowRight" size="sm" />
                <span>{{ primaryLabel }}</span>
              </router-link>

              <a
                v-if="docsHref"
                :href="docsHref"
                target="_blank"
                rel="noopener noreferrer"
                class="b022hub-theme-outline-button inline-flex items-center gap-2 rounded-full px-6 py-3 text-sm font-medium transition"
              >
                <Icon name="book" size="sm" />
                <span>{{ t('home.docsGuide') }}</span>
              </a>
            </div>

            <div class="grid gap-3 sm:grid-cols-3" :style="heroStyle(0.3, -8)">
              <div
                v-for="item in [
                  t('home.heroTags.gateway'),
                  t('home.heroTags.routing'),
                  t('home.heroTags.cli')
                ]"
                :key="item"
                class="b022hub-theme-panel b022hub-theme-text-body rounded-2xl px-4 py-4 text-sm"
              >
                {{ item }}
              </div>
            </div>
          </div>
        </div>
      </section>

      <div :ref="(element) => registerSection(1, element)">
        <B022HubToolSection
          align="right"
          badge-icon="cpu"
          :badge="t('home.sections.claude.badge')"
          :title="t('home.sections.claude.title')"
          :description="t('home.sections.claude.description')"
          :install-options="claudeInstallOptions"
          :selected-install-key="claudeInstallKey"
          :config-files="[
            {
              path: '~/.claude/settings.json',
              languageLabel: 'JSON',
              content: claudeConfig
            }
          ]"
          :progress="sectionProgresses[1] || 0"
          @copy="copyConfig"
          @update:selected-install-key="claudeInstallKey = $event"
        />
      </div>

      <div :ref="(element) => registerSection(2, element)">
        <B022HubToolSection
          align="left"
          badge-icon="terminal"
          :badge="t('home.sections.codex.badge')"
          :title="t('home.sections.codex.title')"
          :description="t('home.sections.codex.description')"
          :install-options="codexInstallOptions"
          :selected-install-key="codexInstallKey"
          :config-files="[
            {
              path: '~/.codex/config.toml',
              languageLabel: 'TOML',
              content: codexConfig
            },
            {
              path: '~/.codex/auth.json',
              languageLabel: 'JSON',
              content: codexAuth
            }
          ]"
          :progress="sectionProgresses[2] || 0"
          @copy="copyConfig"
          @update:selected-install-key="codexInstallKey = $event"
        />
      </div>

      <div :ref="(element) => registerSection(3, element)">
        <B022HubToolSection
          align="right"
          badge-icon="sparkles"
          :badge="t('home.sections.gemini.badge')"
          :title="t('home.sections.gemini.title')"
          :description="t('home.sections.gemini.description')"
          :install-options="geminiInstallOptions"
          :selected-install-key="geminiInstallKey"
          :config-files="[
            {
              path: '~/.gemini/.env',
              languageLabel: 'ENV',
              content: geminiEnv
            },
            {
              path: '~/.gemini/settings.json',
              languageLabel: 'JSON',
              content: geminiSettings
            }
          ]"
          :progress="sectionProgresses[3] || 0"
          @copy="copyConfig"
          @update:selected-install-key="geminiInstallKey = $event"
        />
      </div>

      <div :ref="(element) => registerSection(4, element)">
        <B022HubProgressSection
          :badge="t('home.progress.badge')"
          :title="t('home.progress.title')"
          :description="t('home.progress.description')"
          :cards="progressCards"
          :doc-href="docsHref"
          :doc-label="t('home.docsGuide')"
          :primary-path="primaryPath"
          :primary-label="primaryLabel"
          :progress="sectionProgresses[4] || 0"
        />
      </div>
    </main>

    <footer class="relative z-20 px-5 pb-10 pt-2 sm:px-8 md:px-14 lg:px-20">
      <div class="mx-auto max-w-7xl rounded-[1.75rem] border border-[rgba(227,224,211,0.12)] bg-[rgba(13,12,11,0.58)] px-5 py-4 text-center backdrop-blur-xl">
        <p class="b022hub-theme-text-muted text-sm tracking-[0.04em]">
          {{ t('home.footer.forkNotice', { siteName }) }}
        </p>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.b022hub-shell {
  background: var(--b022-bg);
}

.b022hub-grid,
.b022hub-noise {
  pointer-events: none;
  position: fixed;
  inset: 0;
}

.b022hub-grid {
  background-image:
    linear-gradient(var(--b022-grid-line) 1px, transparent 1px),
    linear-gradient(90deg, var(--b022-grid-line) 1px, transparent 1px);
  background-size: 42px 42px;
  mask-image: linear-gradient(180deg, rgba(0, 0, 0, 0.8), transparent 100%);
}

.b022hub-noise {
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.85' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E");
  opacity: var(--b022-noise-opacity);
  mix-blend-mode: soft-light;
}

.type-cursor {
  display: inline-block;
  min-width: 0.75ch;
}
</style>
