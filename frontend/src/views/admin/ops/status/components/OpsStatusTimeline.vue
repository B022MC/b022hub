<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import type { OpsStatusMatrixBucket, OpsStatusMatrixBucketStatus } from '@/api/admin/ops'
import { formatDateTime } from '@/utils/format'

interface Props {
  buckets: OpsStatusMatrixBucket[]
}

const props = defineProps<Props>()
const { t } = useI18n()

const normalizedBuckets = computed(() => props.buckets ?? [])

const cellClassMap: Record<OpsStatusMatrixBucketStatus, string> = {
  ok: 'status-timeline-cell-ok',
  warn: 'status-timeline-cell-warn',
  down: 'status-timeline-cell-down',
  nodata: 'status-timeline-cell-nodata'
}

const getCellClass = (status: OpsStatusMatrixBucketStatus) => cellClassMap[status] ?? cellClassMap.nodata

const getStatusLabel = (status: OpsStatusMatrixBucketStatus) => t(`admin.ops.statusMatrix.status.${status}`)

const getCellTitle = (bucket: OpsStatusMatrixBucket) => {
  const summary = t('admin.ops.statusMatrix.tooltips.timeline', {
    status: getStatusLabel(bucket.status),
    success: bucket.success_count,
    error: bucket.error_count,
    excluded: bucket.excluded_error_count
  })
  return `${formatDateTime(bucket.bucket_start)} - ${formatDateTime(bucket.bucket_end)}\n${summary}`
}
</script>

<template>
  <div
    class="status-timeline"
    :style="{ '--timeline-columns': String(normalizedBuckets.length || 1) }"
    :data-bucket-count="normalizedBuckets.length"
  >
    <div
      v-for="bucket in normalizedBuckets"
      :key="`${bucket.bucket_start}-${bucket.bucket_end}`"
      data-testid="status-cell"
      class="status-timeline-cell"
      :class="getCellClass(bucket.status)"
      :data-status="bucket.status"
      :title="getCellTitle(bucket)"
    ></div>
  </div>
</template>

<style scoped>
.status-timeline {
  display: grid;
  grid-template-columns: repeat(var(--timeline-columns), minmax(0, 1fr));
  gap: 0.25rem;
  min-width: 18rem;
}

.status-timeline-cell {
  @apply h-5 rounded-md border;
  border-color: transparent;
}

.status-timeline-cell-ok {
  background: linear-gradient(180deg, rgba(34, 197, 94, 0.92) 0%, rgba(22, 163, 74, 0.92) 100%);
}

.status-timeline-cell-warn {
  background: linear-gradient(180deg, rgba(245, 158, 11, 0.92) 0%, rgba(217, 119, 6, 0.92) 100%);
}

.status-timeline-cell-down {
  background: linear-gradient(180deg, rgba(239, 68, 68, 0.92) 0%, rgba(220, 38, 38, 0.92) 100%);
}

.status-timeline-cell-nodata {
  border-color: rgba(148, 163, 184, 0.28);
  background: rgba(148, 163, 184, 0.18);
}

html.dark .status-timeline-cell-nodata {
  border-color: rgba(148, 163, 184, 0.2);
  background: rgba(148, 163, 184, 0.12);
}
</style>
