<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import { formatDateTime, formatRelativeTime } from '@/utils/format'

import OpsStatusMatrixFilters from './components/OpsStatusMatrixFilters.vue'
import OpsStatusMatrixTable from './components/OpsStatusMatrixTable.vue'
import { useOpsStatusMatrix } from './composables/useOpsStatusMatrix'

const { t } = useI18n()

const {
  timeRange,
  platform,
  groupId,
  query,
  rows,
  loading,
  groupsLoading,
  errorMessage,
  lastUpdated,
  platformOptions,
  groupOptions,
  refresh
} = useOpsStatusMatrix()

const timeRangeLabel = computed(() => t(`admin.ops.timeRange.${timeRange.value}`))
const lastUpdatedLabel = computed(() => (
  lastUpdated.value ? formatRelativeTime(lastUpdated.value) : '-'
))
</script>

<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="space-y-4">
          <section class="status-matrix-hero">
            <div class="status-matrix-hero-copy">
              <p class="status-matrix-hero-eyebrow">{{ t('admin.ops.statusMatrix.eyebrow') }}</p>
              <h1 class="status-matrix-hero-title">{{ t('admin.ops.statusMatrix.title') }}</h1>
              <p class="status-matrix-hero-description">{{ t('admin.ops.statusMatrix.description') }}</p>
            </div>

            <div class="status-matrix-hero-metrics">
              <div class="status-matrix-hero-metric">
                <span class="status-matrix-hero-label">{{ t('admin.ops.statusMatrix.summary.rows') }}</span>
                <strong class="status-matrix-hero-value">{{ rows.length }}</strong>
              </div>
              <div class="status-matrix-hero-metric">
                <span class="status-matrix-hero-label">{{ t('admin.ops.statusMatrix.summary.window') }}</span>
                <strong class="status-matrix-hero-value">{{ timeRangeLabel }}</strong>
              </div>
              <div class="status-matrix-hero-metric">
                <span class="status-matrix-hero-label">{{ t('admin.ops.statusMatrix.summary.lastUpdated') }}</span>
                <strong
                  class="status-matrix-hero-value"
                  :title="lastUpdated ? formatDateTime(lastUpdated) : undefined"
                >
                  {{ lastUpdatedLabel }}
                </strong>
              </div>
            </div>
          </section>

          <OpsStatusMatrixFilters
            :time-range="timeRange"
            :platform="platform"
            :group-id="groupId"
            :query="query"
            :loading="loading"
            :groups-loading="groupsLoading"
            :platform-options="platformOptions"
            :group-options="groupOptions"
            @update:time-range="timeRange = $event"
            @update:platform="platform = $event"
            @update:group-id="groupId = $event"
            @update:query="query = $event"
            @refresh="refresh"
          />

          <div
            v-if="errorMessage"
            class="status-matrix-error-banner"
          >
            {{ errorMessage }}
          </div>
        </div>
      </template>

      <template #table>
        <OpsStatusMatrixTable
          :rows="rows"
          :loading="loading"
          :time-range="timeRange"
        />
      </template>
    </TablePageLayout>
  </AppLayout>
</template>

<style scoped>
.status-matrix-hero {
  @apply relative overflow-hidden rounded-[28px] border p-6;
  border-color: rgba(56, 189, 248, 0.16);
  background:
    radial-gradient(circle at top right, rgba(56, 189, 248, 0.22), transparent 34%),
    radial-gradient(circle at bottom left, rgba(14, 165, 233, 0.16), transparent 30%),
    linear-gradient(135deg, rgba(11, 18, 32, 0.98) 0%, rgba(15, 35, 62, 0.96) 58%, rgba(12, 74, 110, 0.82) 100%);
  box-shadow: 0 24px 56px rgba(6, 15, 28, 0.28);
}

.status-matrix-hero-copy {
  @apply relative z-[1] max-w-3xl;
}

.status-matrix-hero-eyebrow {
  @apply mb-2 text-xs font-semibold uppercase tracking-[0.28em];
  color: rgba(165, 243, 252, 0.8);
}

.status-matrix-hero-title {
  @apply text-3xl font-black tracking-tight text-white;
}

.status-matrix-hero-description {
  @apply mt-3 max-w-2xl text-sm leading-6;
  color: rgba(226, 232, 240, 0.8);
}

.status-matrix-hero-metrics {
  @apply relative z-[1] mt-5 grid gap-3;
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.status-matrix-hero-metric {
  @apply rounded-2xl border px-4 py-3;
  border-color: rgba(255, 255, 255, 0.1);
  background: rgba(8, 15, 29, 0.34);
  backdrop-filter: blur(14px);
}

.status-matrix-hero-label {
  @apply block text-[11px] font-medium uppercase tracking-[0.18em];
  color: rgba(203, 213, 225, 0.72);
}

.status-matrix-hero-value {
  @apply mt-2 block text-lg font-semibold text-white;
}

.status-matrix-error-banner {
  @apply rounded-2xl border px-4 py-3 text-sm font-medium text-rose-700 dark:text-rose-300;
  border-color: rgba(251, 113, 133, 0.24);
  background: rgba(254, 226, 226, 0.72);
}

html.dark .status-matrix-error-banner {
  background: rgba(127, 29, 29, 0.24);
}

@media (max-width: 1024px) {
  .status-matrix-hero-metrics {
    grid-template-columns: 1fr;
  }
}
</style>
