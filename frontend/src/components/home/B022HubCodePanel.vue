<script setup lang="ts">
defineProps<{
  path: string
  languageLabel: string
  code: string
}>()

defineEmits<{
  copy: [value: string]
}>()
</script>

<template>
  <div
    class="b022hub-theme-panel-strong overflow-hidden rounded-[1.35rem]"
  >
    <div
      class="b022hub-theme-subpanel flex items-center justify-between px-4 py-3"
    >
      <div>
        <p class="b022hub-theme-text-muted text-xs font-medium uppercase tracking-[0.18em]">
          {{ languageLabel }}
        </p>
        <p class="b022hub-theme-code mt-1 text-sm">
          {{ path }}
        </p>
      </div>

      <button
        class="b022hub-theme-outline-button inline-flex items-center rounded-full px-3 py-1.5 text-xs font-medium transition"
        @click="$emit('copy', code)"
      >
        {{ $t('common.copy') }}
      </button>
    </div>

    <pre class="b022hub-theme-code overflow-x-auto px-4 py-4 text-sm leading-7"><code>{{ code }}</code></pre>
  </div>
</template>

<style scoped>
pre {
  font-family: 'JetBrains Mono', 'SFMono-Regular', Menlo, Monaco, Consolas, monospace;
}
</style>
