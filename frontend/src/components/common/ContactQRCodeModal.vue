<template>
  <BaseDialog
    :show="show"
    :title="t('common.contactSupport')"
    width="narrow"
    :close-on-escape="true"
    :close-on-click-outside="true"
    @close="emit('close')"
  >
    <div class="flex flex-col items-center space-y-4">
      <!-- Tabs if both QR codes exist -->
      <div v-if="hasWechat && hasGroup" class="flex w-full rounded-lg bg-gray-100 p-1 dark:bg-dark-700">
        <button
          type="button"
          @click="activeTab = 'wechat'"
          :class="[
            'flex-1 rounded-md px-4 py-2 text-sm font-medium transition-colors',
            activeTab === 'wechat'
              ? 'bg-white text-gray-900 shadow dark:bg-dark-600 dark:text-white'
              : 'text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white'
          ]"
        >
          {{ t('contact.qrcode.wechat') }}
        </button>
        <button
          type="button"
          @click="activeTab = 'group'"
          :class="[
            'flex-1 rounded-md px-4 py-2 text-sm font-medium transition-colors',
            activeTab === 'group'
              ? 'bg-white text-gray-900 shadow dark:bg-dark-600 dark:text-white'
              : 'text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white'
          ]"
        >
          {{ t('contact.qrcode.group') }}
        </button>
      </div>

      <!-- QR Code Display -->
      <div class="flex flex-col items-center">
        <!-- WeChat QR Code -->
        <template v-if="(hasWechat && activeTab === 'wechat') || (hasWechat && !hasGroup)">
          <img
            :src="wechatQRCode"
            alt="WeChat QR Code"
            class="h-64 w-64 rounded-lg border border-gray-200 object-contain dark:border-dark-600"
          />
          <p class="mt-3 text-sm text-gray-600 dark:text-gray-400">
            {{ t('contact.qrcode.scanWechat') }}
          </p>
        </template>

        <!-- Group QR Code -->
        <template v-else-if="(hasGroup && activeTab === 'group') || (hasGroup && !hasWechat)">
          <img
            :src="groupQRCode"
            alt="Group QR Code"
            class="h-64 w-64 rounded-lg border border-gray-200 object-contain dark:border-dark-600"
          />
          <p class="mt-3 text-sm text-gray-600 dark:text-gray-400">
            {{ t('contact.qrcode.scanGroup') }}
          </p>
        </template>
      </div>

      <!-- Contact Info Text -->
      <div v-if="contactInfo" class="w-full border-t border-gray-100 pt-4 dark:border-dark-700">
        <p class="text-center text-sm text-gray-600 dark:text-gray-400">
          {{ contactInfo }}
        </p>
      </div>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from './BaseDialog.vue'

interface Props {
  show: boolean
  wechatQRCode?: string
  groupQRCode?: string
  contactInfo?: string
}

interface Emits {
  (e: 'close'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const { t } = useI18n()

const hasWechat = computed(() => !!props.wechatQRCode)
const hasGroup = computed(() => !!props.groupQRCode)

// Default to wechat tab, or group if wechat not available
const activeTab = ref<'wechat' | 'group'>('wechat')

// Reset tab when modal opens
watch(() => props.show, (isOpen) => {
  if (isOpen) {
    activeTab.value = hasWechat.value ? 'wechat' : 'group'
  }
})
</script>
