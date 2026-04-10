<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import type { Column } from '@/components/common/types'
import DataTable from '@/components/common/DataTable.vue'
import PlatformIcon from '@/components/common/PlatformIcon.vue'
import type { OpsStatusMatrixRow, OpsStatusMatrixTimeRange } from '@/api/admin/ops'
import { formatDateTime, formatRelativeTime } from '@/utils/format'

import OpsStatusTimeline from './OpsStatusTimeline.vue'

interface Props {
  rows: OpsStatusMatrixRow[]
  loading?: boolean
  timeRange: OpsStatusMatrixTimeRange
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const { t } = useI18n()

const platformLabelMap: Record<string, string> = {
  openai: 'OpenAI',
  anthropic: 'Anthropic',
  gemini: 'Gemini',
  antigravity: 'Antigravity',
  sora: 'Sora'
}

const columns = computed<Column[]>(() => [
  { key: 'platform', label: t('admin.ops.statusMatrix.table.platform'), class: '!whitespace-normal min-w-[120px]' },
  { key: 'group_name', label: t('admin.ops.statusMatrix.table.group'), class: '!whitespace-normal min-w-[160px]' },
  { key: 'model', label: t('admin.ops.statusMatrix.table.model'), class: '!whitespace-normal min-w-[180px]' },
  { key: 'cache_hit_rate', label: t('admin.ops.statusMatrix.table.cacheHitRate'), class: 'min-w-[120px]' },
  { key: 'last_latency_ms', label: t('admin.ops.statusMatrix.table.lastLatency'), class: 'min-w-[120px]' },
  { key: 'last_checked_at', label: t('admin.ops.statusMatrix.table.lastChecked'), class: 'min-w-[140px]' },
  { key: 'availability', label: t('admin.ops.statusMatrix.table.availability'), class: 'min-w-[120px]' },
  {
    key: 'timeline',
    label: t('admin.ops.statusMatrix.table.timeline'),
    class: `!whitespace-normal ${props.timeRange === '24h' ? 'min-w-[340px]' : 'min-w-[280px]'}`
  }
])

const resolveRowKey = (row: OpsStatusMatrixRow) => `${row.platform}:${row.group_id ?? 'nogroup'}:${row.model}`

const formatLatency = (value?: number | null) => (
  typeof value === 'number' ? `${Math.round(value)}ms` : '-'
)

const formatAvailability = (value?: number | null) => {
  if (typeof value !== 'number') {
    return '-'
  }
  return `${(value * 100).toFixed(2).replace(/\.00$/, '').replace(/(\.\d)0$/, '$1')}%`
}

const formatCacheHitRate = (value?: number | null) => {
  if (typeof value !== 'number') {
    return '-'
  }
  return `${(value * 100).toFixed(2)}%`
}

const getAvailabilityClass = (value?: number | null) => {
  if (typeof value !== 'number') {
    return 'status-matrix-availability-neutral'
  }
  if (value >= 0.95) {
    return 'status-matrix-availability-good'
  }
  if (value >= 0.8) {
    return 'status-matrix-availability-warn'
  }
  return 'status-matrix-availability-bad'
}

const getPlatformLabel = (value: string) => platformLabelMap[value] ?? (value || '-')

const getExcludedTooltip = (row: OpsStatusMatrixRow) => t('admin.ops.statusMatrix.tooltips.excludedErrors', {
  count: row.excluded_error_count
})
</script>

<template>
  <DataTable
    :columns="columns"
    :data="props.rows"
    :loading="props.loading"
    :row-key="resolveRowKey"
    :sticky-first-column="false"
    :sticky-actions-column="false"
    :expandable-actions="false"
  >
    <template #cell-platform="{ row }">
      <div class="flex items-center gap-2">
        <PlatformIcon :platform="row.platform" size="md" />
        <span class="font-medium text-gray-900 dark:text-white">{{ getPlatformLabel(row.platform) }}</span>
      </div>
    </template>

    <template #cell-group_name="{ row }">
      <span class="block max-w-[220px] truncate text-gray-700 dark:text-gray-200" :title="row.group_name || '-'">
        {{ row.group_name || '-' }}
      </span>
    </template>

    <template #cell-model="{ row }">
      <code class="rounded-lg bg-gray-100 px-2 py-1 text-xs text-gray-700 dark:bg-dark-700 dark:text-gray-100">
        {{ row.model }}
      </code>
    </template>

    <template #cell-cache_hit_rate="{ value }">
      <span data-testid="cache-hit-rate-value" class="font-medium text-sky-600 dark:text-sky-400">
        {{ formatCacheHitRate(value) }}
      </span>
    </template>

    <template #cell-last_latency_ms="{ value }">
      <span class="font-medium text-gray-700 dark:text-gray-200">{{ formatLatency(value) }}</span>
    </template>

    <template #cell-last_checked_at="{ value }">
      <span class="text-gray-700 dark:text-gray-200" :title="value ? formatDateTime(value) : undefined">
        {{ value ? formatRelativeTime(value) : '-' }}
      </span>
    </template>

    <template #cell-availability="{ row, value }">
      <div class="flex items-center gap-2">
        <span
          data-testid="availability-value"
          class="status-matrix-availability"
          :class="getAvailabilityClass(value)"
        >
          {{ formatAvailability(value) }}
        </span>

        <span
          v-if="row.excluded_error_count > 0"
          data-testid="availability-tooltip"
          class="status-matrix-excluded-indicator"
          :title="getExcludedTooltip(row)"
        >
          {{ t('admin.ops.statusMatrix.table.excludedShort') }}
        </span>
      </div>
    </template>

    <template #cell-timeline="{ row }">
      <OpsStatusTimeline :buckets="row.buckets" />
    </template>
  </DataTable>
</template>

<style scoped>
.status-matrix-availability {
  @apply inline-flex items-center rounded-full px-2.5 py-1 text-xs font-semibold;
}

.status-matrix-availability-good {
  @apply bg-emerald-100 text-emerald-700 dark:text-emerald-300;
}

html.dark .status-matrix-availability-good {
  background: rgba(16, 185, 129, 0.15);
}

.status-matrix-availability-warn {
  @apply bg-amber-100 text-amber-700 dark:text-amber-300;
}

html.dark .status-matrix-availability-warn {
  background: rgba(245, 158, 11, 0.15);
}

.status-matrix-availability-bad {
  @apply bg-rose-100 text-rose-700 dark:text-rose-300;
}

html.dark .status-matrix-availability-bad {
  background: rgba(244, 63, 94, 0.15);
}

.status-matrix-availability-neutral {
  @apply bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-dark-300;
}

.status-matrix-excluded-indicator {
  @apply inline-flex cursor-help items-center rounded-full border px-2 py-0.5 text-[11px] font-medium;
  border-color: rgba(245, 158, 11, 0.28);
  color: rgb(180 83 9);
  background: rgba(251, 191, 36, 0.16);
}

html.dark .status-matrix-excluded-indicator {
  border-color: rgba(245, 158, 11, 0.24);
  color: rgb(252 211 77);
  background: rgba(245, 158, 11, 0.12);
}
</style>
