import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { userGroupsAPI, userOpsAPI } from '@/api'
import { adminAPI } from '@/api/admin'
import type {
  OpsStatusMatrixResponse,
  OpsStatusMatrixRow,
  OpsStatusMatrixSort,
  OpsStatusMatrixTimeRange,
} from '@/api/admin/ops'
import { useAppStore, useAuthStore } from '@/stores'
import type { Group, GroupPlatform } from '@/types'

const SEARCH_DEBOUNCE_MS = 250

const PLATFORM_OPTIONS: Array<{ value: GroupPlatform | ''; label: string }> = [
  { value: '', label: 'All' },
  { value: 'openai', label: 'OpenAI' },
  { value: 'anthropic', label: 'Anthropic' },
  { value: 'gemini', label: 'Gemini' },
  { value: 'antigravity', label: 'Antigravity' },
  { value: 'sora', label: 'Sora' }
]

export function useOpsStatusMatrix() {
  const { t } = useI18n()
  const appStore = useAppStore()
  const authStore = useAuthStore()

  const timeRange = ref<OpsStatusMatrixTimeRange>('90m')
  const platform = ref<GroupPlatform | ''>('')
  const groupId = ref<number | null>(null)
  const query = ref('')
  const sort = ref<OpsStatusMatrixSort>('availability_asc')

  const groups = ref<Array<Pick<Group, 'id' | 'name' | 'platform'>>>([])
  const response = ref<OpsStatusMatrixResponse | null>(null)
  const loading = ref(false)
  const groupsLoading = ref(false)
  const errorMessage = ref('')
  const lastUpdated = ref<Date | null>(null)

  let requestController: AbortController | null = null
  let requestSeq = 0
  let searchTimer: number | undefined

  const rows = computed<OpsStatusMatrixRow[]>(() => response.value?.rows ?? [])
  const platformOptions = computed(() => (
    PLATFORM_OPTIONS.map((option) => ({
      value: option.value,
      label: option.value ? option.label : t('common.all')
    }))
  ))

  const filteredGroups = computed(() => (
    platform.value
      ? groups.value.filter((item) => item.platform === platform.value)
      : groups.value
  ))

  const groupOptions = computed(() => [
    { value: null, label: t('common.all') },
    ...filteredGroups.value.map((item) => ({
      value: item.id,
      label: item.name
    }))
  ])

  watch(platform, (nextPlatform) => {
    if (!groupId.value) {
      return
    }
    const selectedGroup = groups.value.find((item) => item.id === groupId.value)
    if (selectedGroup && nextPlatform && selectedGroup.platform !== nextPlatform) {
      groupId.value = null
    }
  })

  const buildParams = () => ({
    time_range: timeRange.value,
    platform: platform.value || undefined,
    group_id: groupId.value ?? undefined,
    q: query.value.trim() || undefined,
    sort: sort.value
  })

  const loadGroups = async () => {
    groupsLoading.value = true
    try {
      groups.value = authStore.isAdmin
        ? await adminAPI.groups.getAll()
        : await userGroupsAPI.getAvailable()
    } catch (error) {
      console.error('[OpsStatusMatrix] Failed to load groups', error)
      groups.value = []
    } finally {
      groupsLoading.value = false
    }
  }

  const loadStatusMatrix = async () => {
    requestController?.abort()

    const currentSeq = ++requestSeq
    const controller = new AbortController()
    requestController = controller
    loading.value = true
    errorMessage.value = ''

    try {
      const data = authStore.isAdmin
        ? await adminAPI.ops.getStatusMatrix(buildParams(), { signal: controller.signal })
        : await userOpsAPI.getStatusMatrix(buildParams(), { signal: controller.signal })
      if (currentSeq !== requestSeq || controller.signal.aborted) {
        return
      }
      response.value = data
      lastUpdated.value = new Date()
    } catch (error) {
      if (controller.signal.aborted) {
        return
      }
      console.error('[OpsStatusMatrix] Failed to load status matrix', error)
      errorMessage.value = t('admin.ops.statusMatrix.errors.loadFailed')
      appStore.showError(errorMessage.value)
    } finally {
      if (currentSeq === requestSeq && !controller.signal.aborted) {
        loading.value = false
      }
    }
  }

  const refresh = async () => {
    await loadStatusMatrix()
  }

  watch([timeRange, platform, groupId, sort], () => {
    void loadStatusMatrix()
  }, { immediate: true })

  watch(query, () => {
    if (searchTimer) {
      window.clearTimeout(searchTimer)
    }
    searchTimer = window.setTimeout(() => {
      void loadStatusMatrix()
    }, SEARCH_DEBOUNCE_MS)
  })

  onMounted(() => {
    void loadGroups()
  })

  onUnmounted(() => {
    requestController?.abort()
    if (searchTimer) {
      window.clearTimeout(searchTimer)
    }
  })

  return {
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
  }
}
