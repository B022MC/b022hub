<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import Select, { type SelectOption } from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import type { OpsStatusMatrixTimeRange } from '@/api/admin/ops'
import type { GroupPlatform } from '@/types'

interface Props {
  timeRange: OpsStatusMatrixTimeRange
  platform: GroupPlatform | ''
  groupId: number | null
  query: string
  loading?: boolean
  groupsLoading?: boolean
  platformOptions: SelectOption[]
  groupOptions: SelectOption[]
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  groupsLoading: false
})

const emit = defineEmits<{
  (e: 'update:timeRange', value: OpsStatusMatrixTimeRange): void
  (e: 'update:platform', value: GroupPlatform | ''): void
  (e: 'update:groupId', value: number | null): void
  (e: 'update:query', value: string): void
  (e: 'refresh'): void
}>()

const { t } = useI18n()

const timeRangeOptions = computed<Array<{ value: OpsStatusMatrixTimeRange; label: string }>>(() => [
  { value: '90m', label: t('admin.ops.timeRange.90m') },
  { value: '24h', label: t('admin.ops.timeRange.24h') }
])
</script>

<template>
  <div class="status-matrix-filters-card">
    <div class="status-matrix-filters-topbar">
      <div class="status-matrix-range-toggle" role="tablist" :aria-label="t('admin.dashboard.timeRange')">
        <button
          v-for="option in timeRangeOptions"
          :key="option.value"
          type="button"
          class="status-matrix-range-button"
          :class="{ 'status-matrix-range-button-active': props.timeRange === option.value }"
          @click="emit('update:timeRange', option.value)"
        >
          {{ option.label }}
        </button>
      </div>

      <button
        type="button"
        class="btn btn-secondary"
        :disabled="props.loading"
        :title="t('common.refresh')"
        @click="emit('refresh')"
      >
        <Icon name="refresh" size="md" :class="props.loading ? 'animate-spin' : ''" />
      </button>
    </div>

    <div class="status-matrix-filters-grid">
      <label class="status-matrix-search-field">
        <Icon
          name="search"
          size="md"
          class="status-matrix-search-icon"
        />
        <input
          :value="props.query"
          type="text"
          :placeholder="t('admin.ops.statusMatrix.filters.searchPlaceholder')"
          class="input status-matrix-search-input"
          @input="emit('update:query', ($event.target as HTMLInputElement).value)"
        />
      </label>

      <div class="status-matrix-select">
        <Select
          :model-value="props.platform"
          :options="props.platformOptions"
          :placeholder="t('admin.ops.statusMatrix.filters.platform')"
          @update:model-value="emit('update:platform', ($event ?? '') as GroupPlatform | '')"
        />
      </div>

      <div class="status-matrix-select">
        <Select
          :model-value="props.groupId"
          :options="props.groupOptions"
          :disabled="props.groupsLoading"
          :placeholder="t('admin.ops.statusMatrix.filters.group')"
          @update:model-value="emit('update:groupId', ($event ?? null) as number | null)"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.status-matrix-filters-card {
  @apply rounded-3xl border p-4 shadow-sm;
  border-color: rgba(125, 94, 71, 0.12);
  background: rgba(255, 251, 245, 0.82);
  box-shadow: 0 18px 48px rgba(89, 62, 41, 0.08);
  backdrop-filter: blur(18px);
}

html.dark .status-matrix-filters-card {
  border-color: rgba(227, 224, 211, 0.08);
  background: rgba(20, 18, 16, 0.78);
  box-shadow: 0 22px 56px rgba(0, 0, 0, 0.28);
}

.status-matrix-filters-topbar {
  @apply mb-4 flex flex-wrap items-center justify-between gap-3;
}

.status-matrix-range-toggle {
  @apply inline-flex rounded-2xl p-1;
  background: rgba(125, 94, 71, 0.08);
}

html.dark .status-matrix-range-toggle {
  background: rgba(227, 224, 211, 0.08);
}

.status-matrix-range-button {
  @apply rounded-xl px-4 py-2 text-sm font-medium text-gray-600 transition-colors dark:text-dark-300;
}

.status-matrix-range-button-active {
  @apply text-white;
  background: linear-gradient(135deg, #1b85ff 0%, #11b3ff 100%);
  box-shadow: 0 10px 24px rgba(27, 133, 255, 0.28);
}

.status-matrix-filters-grid {
  @apply grid gap-3;
  grid-template-columns: minmax(0, 1.4fr) repeat(2, minmax(0, 220px));
}

.status-matrix-search-field {
  @apply relative block;
}

.status-matrix-search-icon {
  @apply pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-dark-500;
}

.status-matrix-search-input {
  @apply pl-10;
}

.status-matrix-select {
  @apply min-w-0;
}

@media (max-width: 1024px) {
  .status-matrix-filters-grid {
    grid-template-columns: 1fr;
  }
}
</style>
