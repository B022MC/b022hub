<template>
  <BaseDialog
    :show="show"
    :title="t('admin.accounts.dataImportTitle')"
    width="normal"
    close-on-click-outside
    @close="handleClose"
  >
    <form id="import-data-form" class="space-y-4" @submit.prevent="handleImport">
      <div class="text-sm text-gray-600 dark:text-dark-300">
        {{ t('admin.accounts.dataImportHint') }}
      </div>
      <div
        class="rounded-lg border border-amber-200 bg-amber-50 p-3 text-xs text-amber-600 dark:border-amber-800 dark:bg-amber-900/20 dark:text-amber-400"
      >
        {{ t('admin.accounts.dataImportWarning') }}
      </div>

      <div class="space-y-2">
        <GroupSelector v-model="selectedGroupIds" :groups="groups" />
        <div class="text-xs text-gray-500 dark:text-dark-400">
          {{ t('admin.accounts.dataImportGroupHint') }}
        </div>
      </div>

      <div>
        <label class="input-label">{{ t('admin.accounts.dataImportFile') }}</label>
        <div
          class="flex items-center justify-between gap-3 rounded-lg border border-dashed border-gray-300 bg-gray-50 px-4 py-3 dark:border-dark-600 dark:bg-dark-800"
        >
          <div class="min-w-0" :title="selectedFileNames || undefined">
            <div class="truncate text-sm text-gray-700 dark:text-dark-200">
              {{ fileSummary || t('admin.accounts.dataImportSelectFile') }}
            </div>
            <div class="truncate text-xs text-gray-500 dark:text-dark-400">
              {{ fileDetail }}
            </div>
          </div>
          <button type="button" class="btn btn-secondary shrink-0" @click="openFilePicker">
            {{ t('common.chooseFile') }}
          </button>
        </div>
        <input
          ref="fileInput"
          type="file"
          class="hidden"
          accept="application/json,.json"
          multiple
          @change="handleFileChange"
        />
      </div>

      <div
        v-if="result"
        class="space-y-2 rounded-xl border border-gray-200 p-4 dark:border-dark-700"
      >
        <div class="text-sm font-medium text-gray-900 dark:text-white">
          {{ t('admin.accounts.dataImportResult') }}
        </div>
        <div class="text-sm text-gray-700 dark:text-dark-300">
          {{ t('admin.accounts.dataImportResultSummary', result) }}
        </div>

        <div v-if="errorItems.length" class="mt-2">
          <div class="text-sm font-medium text-red-600 dark:text-red-400">
            {{ t('admin.accounts.dataImportErrors') }}
          </div>
          <div
            class="mt-2 max-h-48 overflow-auto rounded-lg bg-gray-50 p-3 font-mono text-xs dark:bg-dark-800"
          >
            <div v-for="(item, idx) in errorItems" :key="idx" class="whitespace-pre-wrap">
              {{ item.kind }} {{ item.name || item.proxy_key || '-' }} — {{ item.message }}
            </div>
          </div>
        </div>
      </div>
    </form>

    <template #footer>
      <div class="flex justify-end gap-3">
        <button class="btn btn-secondary" type="button" :disabled="importing" @click="handleClose">
          {{ t('common.cancel') }}
        </button>
        <button
          class="btn btn-primary"
          type="submit"
          form="import-data-form"
          :disabled="importing"
        >
          {{ importing ? t('admin.accounts.dataImporting') : t('admin.accounts.dataImportButton') }}
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import GroupSelector from '@/components/common/GroupSelector.vue'
import { adminAPI } from '@/api/admin'
import { useAppStore } from '@/stores/app'
import type { AdminDataImportResult, AdminGroup } from '@/types'

interface Props {
  show: boolean
  groups?: AdminGroup[]
}

interface Emits {
  (e: 'close'): void
  (e: 'imported'): void
}

interface BatchImportError {
  kind: 'proxy' | 'account' | 'file'
  name?: string
  proxy_key?: string
  message: string
}

interface BatchImportResult {
  proxy_created: number
  proxy_reused: number
  proxy_failed: number
  account_created: number
  account_failed: number
  errors: BatchImportError[]
  file_total: number
  file_successful: number
  file_failed: number
}

const props = withDefaults(defineProps<Props>(), {
  groups: () => []
})
const emit = defineEmits<Emits>()

const { t } = useI18n()
const appStore = useAppStore()

const importing = ref(false)
const files = ref<File[]>([])
const result = ref<BatchImportResult | null>(null)
const selectedGroupIds = ref<number[]>([])

const fileInput = ref<HTMLInputElement | null>(null)
const groups = computed(() => props.groups)
const selectedFileNames = computed(() => files.value.map((item) => item.name).join(', '))
const fileSummary = computed(() => {
  if (files.value.length === 1) {
    return files.value[0].name
  }
  if (files.value.length > 1) {
    return t('admin.accounts.dataImportSelectedFiles', { count: files.value.length })
  }
  return ''
})
const fileDetail = computed(() =>
  files.value.length > 1 ? selectedFileNames.value : t('admin.accounts.dataImportFileHint')
)

const errorItems = computed(() => result.value?.errors || [])

watch(
  () => props.show,
  (open) => {
    if (open) {
      files.value = []
      result.value = null
      selectedGroupIds.value = []
      if (fileInput.value) {
        fileInput.value.value = ''
      }
    }
  }
)

const openFilePicker = () => {
  fileInput.value?.click()
}

const handleFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  files.value = Array.from(target.files || [])
}

const handleClose = () => {
  if (importing.value) return
  emit('close')
}

const readFileAsText = async (sourceFile: File): Promise<string> => {
  if (typeof sourceFile.text === 'function') {
    return sourceFile.text()
  }

  if (typeof sourceFile.arrayBuffer === 'function') {
    const buffer = await sourceFile.arrayBuffer()
    return new TextDecoder().decode(buffer)
  }

  return await new Promise<string>((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(String(reader.result ?? ''))
    reader.onerror = () => reject(reader.error || new Error('Failed to read file'))
    reader.readAsText(sourceFile)
  })
}

const createEmptyResult = (): BatchImportResult => ({
  proxy_created: 0,
  proxy_reused: 0,
  proxy_failed: 0,
  account_created: 0,
  account_failed: 0,
  errors: [],
  file_total: 0,
  file_successful: 0,
  file_failed: 0
})

const mergeImportResult = (aggregate: BatchImportResult, current: AdminDataImportResult) => {
  aggregate.proxy_created += current.proxy_created
  aggregate.proxy_reused += current.proxy_reused
  aggregate.proxy_failed += current.proxy_failed
  aggregate.account_created += current.account_created
  aggregate.account_failed += current.account_failed
  if (current.errors?.length) {
    aggregate.errors.push(...current.errors)
  }
}

const buildMessageParams = (payload: BatchImportResult): Record<string, number> => ({
  account_created: payload.account_created,
  account_failed: payload.account_failed,
  proxy_created: payload.proxy_created,
  proxy_reused: payload.proxy_reused,
  proxy_failed: payload.proxy_failed,
  file_total: payload.file_total,
  file_successful: payload.file_successful,
  file_failed: payload.file_failed
})

const handleImport = async () => {
  if (files.value.length === 0) {
    appStore.showError(t('admin.accounts.dataImportSelectFile'))
    return
  }

  importing.value = true
  try {
    const aggregate = createEmptyResult()
    const groupIds = selectedGroupIds.value.length > 0 ? [...selectedGroupIds.value] : undefined

    for (const sourceFile of files.value) {
      aggregate.file_total += 1

      try {
        const text = await readFileAsText(sourceFile)
        const dataPayload = JSON.parse(text)
        const res = await adminAPI.accounts.importData({
          data: dataPayload,
          group_ids: groupIds,
          skip_default_group_bind: true
        })

        mergeImportResult(aggregate, res)
        aggregate.file_successful += 1
      } catch (error: any) {
        aggregate.file_failed += 1
        aggregate.errors.push({
          kind: 'file',
          name: sourceFile.name,
          message:
            error instanceof SyntaxError
              ? t('admin.accounts.dataImportParseFailed')
              : (error?.message ?? t('admin.accounts.dataImportFailed'))
        })
      }
    }

    result.value = aggregate

    const msgParams = buildMessageParams(aggregate)
    const singleFileFailure =
      files.value.length === 1 && aggregate.file_successful === 0 && aggregate.errors.length === 1
    const hasErrors =
      aggregate.file_failed > 0 ||
      aggregate.account_failed > 0 ||
      aggregate.proxy_failed > 0 ||
      aggregate.errors.length > 0

    if (aggregate.file_successful > 0) {
      emit('imported')
    }

    if (singleFileFailure) {
      appStore.showError(aggregate.errors[0].message)
    } else if (hasErrors) {
      appStore.showError(t('admin.accounts.dataImportCompletedWithErrors', msgParams))
    } else {
      appStore.showSuccess(t('admin.accounts.dataImportSuccess', msgParams))
    }
  } catch (error: any) {
    appStore.showError(error?.message || t('admin.accounts.dataImportFailed'))
  } finally {
    importing.value = false
  }
}
</script>
