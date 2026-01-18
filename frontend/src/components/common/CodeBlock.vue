<template>
  <div class="group relative">
    <!-- Code Block -->
    <div
      class="overflow-x-auto rounded-xl border border-gray-200 bg-gray-900 dark:border-dark-700"
    >
      <!-- Header with Language and Copy Button -->
      <div class="flex items-center justify-between border-b border-gray-700 px-4 py-2">
        <span class="text-xs font-medium text-gray-400">{{ language }}</span>
        <button
          @click="copyCode"
          class="flex items-center gap-1.5 rounded-lg px-2 py-1 text-xs text-gray-400 transition-colors hover:bg-gray-800 hover:text-gray-200"
          :title="t('common.copiedToClipboard')"
        >
          <Icon v-if="copied" name="check" size="sm" class="text-emerald-400" />
          <Icon v-else name="copy" size="sm" />
          <span>{{ copied ? t('codeBlock.copied') : t('codeBlock.copy') }}</span>
        </button>
      </div>
      <!-- Code Content -->
      <pre
        class="overflow-x-auto p-4 text-sm leading-relaxed text-gray-100"
      ><code>{{ code }}</code></pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()

const props = withDefaults(defineProps<{
  code: string
  language?: string
}>(), {
  language: 'bash'
})

const copied = ref(false)

async function copyCode() {
  try {
    await navigator.clipboard.writeText(props.code)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}
</script>
