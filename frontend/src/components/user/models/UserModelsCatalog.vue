<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import EmptyState from '@/components/common/EmptyState.vue'
import GroupBadge from '@/components/common/GroupBadge.vue'
import SearchInput from '@/components/common/SearchInput.vue'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import { userGroupsAPI, type AvailableModelsGroup } from '@/api/groups'
import { useAppStore } from '@/stores/app'

interface ModelCatalogItem {
  id: string
  groups: AvailableModelsGroup[]
}

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(false)
const searchQuery = ref('')
const selectedGroupId = ref<'all' | number>('all')
const groups = ref<AvailableModelsGroup[]>([])
const copiedModelId = ref<string | null>(null)

let copyResetTimer: number | null = null

const groupOptions = computed(() => [
  { value: 'all', label: t('models.allGroups') },
  ...groups.value.map((entry) => ({
    value: entry.group.id,
    label: `${entry.group.name} · ${String(entry.group.platform || '').toUpperCase()}`
  }))
])

const selectedGroup = computed(() =>
  selectedGroupId.value === 'all'
    ? null
    : groups.value.find((entry) => entry.group.id === selectedGroupId.value) ?? null
)

const selectedGroupModelCount = computed(() => {
  if (selectedGroupId.value === 'all') {
    return modelCatalog.value.length
  }
  return selectedGroup.value?.models.length ?? 0
})

const modelCatalog = computed<ModelCatalogItem[]>(() => {
  const catalog = new Map<string, AvailableModelsGroup[]>()

  for (const entry of groups.value) {
    if (selectedGroupId.value !== 'all' && entry.group.id !== selectedGroupId.value) {
      continue
    }

    for (const modelId of entry.models) {
      const current = catalog.get(modelId) ?? []
      current.push(entry)
      catalog.set(modelId, current)
    }
  }

  return [...catalog.entries()]
    .map(([id, modelGroups]) => ({ id, groups: modelGroups }))
    .sort((a, b) => a.id.localeCompare(b.id))
})

const filteredCatalog = computed(() => {
  const keyword = searchQuery.value.trim().toLowerCase()
  if (!keyword) {
    return modelCatalog.value
  }

  return modelCatalog.value.filter((item) => {
    if (item.id.toLowerCase().includes(keyword)) {
      return true
    }

    return item.groups.some((entry) => {
      return (
        entry.group.name.toLowerCase().includes(keyword) ||
        String(entry.group.platform || '').toLowerCase().includes(keyword)
      )
    })
  })
})

const totalModels = computed(() => {
  const ids = new Set<string>()
  for (const entry of groups.value) {
    for (const modelId of entry.models) {
      ids.add(modelId)
    }
  }
  return ids.size
})

const totalGroups = computed(() => groups.value.length)

async function loadModels() {
  loading.value = true
  try {
    const response = await userGroupsAPI.getModels()
    groups.value = response.groups ?? []

    if (
      selectedGroupId.value !== 'all' &&
      !groups.value.some((entry) => entry.group.id === selectedGroupId.value)
    ) {
      selectedGroupId.value = 'all'
    }
  } catch (error: any) {
    console.error('[UserModelsCatalog] Failed to load available models', error)
    groups.value = []
    appStore.showError(error?.response?.data?.detail || t('models.loadFailed'))
  } finally {
    loading.value = false
  }
}

async function copyModel(modelId: string) {
  try {
    await navigator.clipboard.writeText(modelId)
    copiedModelId.value = modelId
    appStore.showSuccess(t('models.copySuccess'))

    if (copyResetTimer !== null) {
      window.clearTimeout(copyResetTimer)
    }
    copyResetTimer = window.setTimeout(() => {
      copiedModelId.value = null
      copyResetTimer = null
    }, 1800)
  } catch {
    appStore.showError(t('common.copyFailed'))
  }
}

onMounted(loadModels)

onBeforeUnmount(() => {
  if (copyResetTimer !== null) {
    window.clearTimeout(copyResetTimer)
  }
})
</script>

<template>
  <TablePageLayout>
    <template #actions>
      <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
        <div class="card p-4">
          <div class="flex items-center gap-3">
            <div class="rounded-lg bg-blue-100 p-2 dark:bg-blue-900/30">
              <Icon name="users" size="md" class="text-blue-600 dark:text-blue-400" />
            </div>
            <div>
              <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
                {{ t('models.totalGroups') }}
              </p>
              <p class="text-xl font-bold text-gray-900 dark:text-white">
                {{ totalGroups }}
              </p>
              <p class="text-xs text-gray-500 dark:text-gray-400">
                {{ t('models.groupsHint') }}
              </p>
            </div>
          </div>
        </div>

        <div class="card p-4">
          <div class="flex items-center gap-3">
            <div class="rounded-lg bg-emerald-100 p-2 dark:bg-emerald-900/30">
              <Icon name="cube" size="md" class="text-emerald-600 dark:text-emerald-400" />
            </div>
            <div>
              <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
                {{ t('models.totalModels') }}
              </p>
              <p class="text-xl font-bold text-gray-900 dark:text-white">
                {{ totalModels }}
              </p>
              <p class="text-xs text-gray-500 dark:text-gray-400">
                {{ t('models.uniqueModelsHint') }}
              </p>
            </div>
          </div>
        </div>

        <div class="card p-4">
          <div class="flex items-center gap-3">
            <div class="rounded-lg bg-amber-100 p-2 dark:bg-amber-900/30">
              <Icon name="search" size="md" class="text-amber-600 dark:text-amber-400" />
            </div>
            <div>
              <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
                {{ t('models.selectedModels') }}
              </p>
              <p class="text-xl font-bold text-gray-900 dark:text-white">
                {{ selectedGroupModelCount }}
              </p>
              <p class="text-xs text-gray-500 dark:text-gray-400">
                {{
                  selectedGroup
                    ? t('models.selectedGroupHint', { group: selectedGroup.group.name })
                    : t('models.allGroupsHint')
                }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </template>

    <template #filters>
      <div class="card p-5">
        <div class="flex flex-col gap-4 lg:flex-row lg:items-end">
          <div class="w-full lg:max-w-sm">
            <label class="input-label">{{ t('models.searchLabel') }}</label>
            <SearchInput
              v-model="searchQuery"
              :placeholder="t('models.searchPlaceholder')"
            />
          </div>

          <div class="w-full lg:w-72">
            <label class="input-label">{{ t('models.groupFilterLabel') }}</label>
            <Select
              :model-value="selectedGroupId"
              :options="groupOptions"
              @update:model-value="selectedGroupId = ($event as 'all' | number)"
            />
          </div>

          <div class="ml-auto flex items-center gap-3">
            <button
              class="btn btn-secondary"
              :disabled="loading"
              @click="loadModels"
            >
              <Icon name="refresh" size="md" class="mr-2" :class="loading ? 'animate-spin' : ''" />
              {{ t('common.refresh') }}
            </button>
          </div>
        </div>
      </div>
    </template>

    <template #table>
      <div class="models-scroll-area">
        <div class="flex min-h-full flex-col gap-5 p-5">
          <div class="flex flex-col gap-2 md:flex-row md:items-start md:justify-between">
            <div>
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ t('models.catalogTitle') }}
              </h2>
              <p class="text-sm text-gray-500 dark:text-gray-400">
                {{ t('models.catalogDescription') }}
              </p>
            </div>
            <div class="flex flex-wrap items-center gap-2">
              <GroupBadge
                v-if="selectedGroup"
                :name="selectedGroup.group.name"
                :platform="selectedGroup.group.platform"
                :subscription-type="selectedGroup.group.subscription_type"
                :rate-multiplier="selectedGroup.group.rate_multiplier"
                :show-rate="false"
              />
              <span class="rounded-full bg-gray-100 px-3 py-1 text-xs font-medium text-gray-600 dark:bg-dark-700 dark:text-gray-300">
                {{ t('models.resultsCount', { count: filteredCatalog.length }) }}
              </span>
            </div>
          </div>

          <div v-if="loading" class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
            <div
              v-for="index in 6"
              :key="index"
              class="model-card animate-pulse"
            >
              <div class="h-6 w-1/2 rounded bg-gray-200 dark:bg-dark-700"></div>
              <div class="mt-3 h-4 w-2/3 rounded bg-gray-200 dark:bg-dark-700"></div>
              <div class="mt-4 flex flex-wrap gap-2">
                <div class="h-6 w-24 rounded-full bg-gray-200 dark:bg-dark-700"></div>
                <div class="h-6 w-20 rounded-full bg-gray-200 dark:bg-dark-700"></div>
              </div>
            </div>
          </div>

          <div
            v-else-if="filteredCatalog.length > 0"
            class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3"
          >
            <article
              v-for="item in filteredCatalog"
              :key="item.id"
              class="model-card"
            >
              <div class="flex items-start justify-between gap-3">
                <div class="min-w-0">
                  <code class="model-code">{{ item.id }}</code>
                  <p class="mt-2 text-xs text-gray-500 dark:text-gray-400">
                    {{ t('models.groupCount', { count: item.groups.length }) }}
                  </p>
                </div>

                <button
                  class="rounded-xl border border-gray-200 bg-white/70 p-2 text-gray-500 transition hover:border-primary-300 hover:text-primary-600 dark:border-dark-700 dark:bg-dark-900/40 dark:text-gray-400 dark:hover:border-primary-500/40 dark:hover:text-primary-300"
                  :title="copiedModelId === item.id ? t('common.copied') : t('models.copyModel')"
                  @click="copyModel(item.id)"
                >
                  <Icon
                    :name="copiedModelId === item.id ? 'check' : 'clipboard'"
                    size="sm"
                    :stroke-width="copiedModelId === item.id ? 2 : 1.75"
                  />
                </button>
              </div>

              <div class="mt-4 flex flex-wrap gap-2">
                <GroupBadge
                  v-for="entry in item.groups"
                  :key="`${item.id}-${entry.group.id}`"
                  :name="entry.group.name"
                  :platform="entry.group.platform"
                  :subscription-type="entry.group.subscription_type"
                  :rate-multiplier="entry.group.rate_multiplier"
                  :show-rate="false"
                />
              </div>
            </article>
          </div>

          <EmptyState
            v-else
            :title="t('models.emptyTitle')"
            :description="t('models.emptyDescription')"
          />
        </div>
      </div>
    </template>
  </TablePageLayout>
</template>

<style scoped>
.models-scroll-area {
  @apply h-full min-h-0 overflow-y-auto overflow-x-hidden;
  background: rgba(255, 251, 245, 0.28);
  border-radius: inherit;
  scrollbar-gutter: stable;
}

html.dark .models-scroll-area {
  background: rgba(18, 16, 14, 0.92);
}

.model-card {
  @apply rounded-2xl border border-gray-100 bg-gray-50 p-[1.125rem] shadow-sm transition-all duration-200;
  @apply dark:border-dark-700 dark:bg-dark-800/50;
}

.model-card:hover {
  @apply border-gray-200 bg-white dark:border-dark-600 dark:bg-dark-800;
  transform: translateY(-1px);
}

.model-code {
  @apply inline-block max-w-full break-all rounded-xl bg-gray-100 px-3 py-2 font-mono text-sm font-semibold leading-relaxed text-primary-600;
  @apply dark:bg-dark-900/70 dark:text-primary-400;
}
</style>
