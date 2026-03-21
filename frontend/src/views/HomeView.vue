<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import B022HubLanding from '@/components/home/B022HubLanding.vue'
import { useAuthStore, useAppStore } from '@/stores'

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'b022hub')
const rawSiteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle?.trim() || '')
const siteSubtitle = computed(() => {
  if (
    !rawSiteSubtitle.value ||
    rawSiteSubtitle.value === 'AI API Gateway Platform' ||
    rawSiteSubtitle.value === 'Subscription to API Conversion Platform'
  ) {
    return t('home.heroSubtitle')
  }

  return rawSiteSubtitle.value
})
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const docUrl = computed(() => appStore.docUrl || '')
const apiBaseUrl = computed(
  () => appStore.cachedPublicSettings?.api_base_url || appStore.apiBaseUrl || '',
)
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const registrationEnabled = computed(() => appStore.cachedPublicSettings?.registration_enabled ?? true)
const isHomeContentUrl = computed(() => /^https?:\/\//.test(homeContent.value.trim()))
const isDark = computed(() => appStore.isDark)

const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => (isAdmin.value ? '/admin/dashboard' : '/dashboard'))

const primaryPath = computed(() => {
  if (isAuthenticated.value) {
    return dashboardPath.value
  }

  return registrationEnabled.value ? '/register' : '/login'
})

const primaryLabel = computed(() => {
  if (isAuthenticated.value) {
    return t('home.goToDashboard')
  }

  return registrationEnabled.value ? t('home.cta.button') : t('home.login')
})

function toggleTheme() {
  appStore.toggleTheme()
}

onMounted(() => {
  authStore.checkAuth()

  if (!appStore.publicSettingsLoaded) {
    void appStore.fetchPublicSettings()
  }
})
</script>

<template>
  <div v-if="homeContent" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <div v-else v-html="homeContent"></div>
  </div>

  <B022HubLanding
    v-else
    :site-name="siteName"
    :site-logo="siteLogo"
    :site-subtitle="siteSubtitle"
    :doc-url="docUrl"
    :api-base-url="apiBaseUrl"
    :is-dark="isDark"
    :is-authenticated="isAuthenticated"
    :dashboard-path="dashboardPath"
    :primary-path="primaryPath"
    :primary-label="primaryLabel"
    @toggle-theme="toggleTheme"
  />
</template>
