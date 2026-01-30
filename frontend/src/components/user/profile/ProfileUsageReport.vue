<template>
  <div class="card">
    <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
      <h2 class="text-lg font-medium text-gray-900 dark:text-white">
        {{ t('profile.usageReport.title') }}
      </h2>
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
        {{ t('profile.usageReport.description') }}
      </p>
    </div>
    <div class="px-6 py-6">
      <!-- Loading state -->
      <div v-if="loading" class="flex items-center justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
      </div>

      <!-- Feature disabled globally -->
      <div v-else-if="config && !config.global_enabled" class="flex items-center gap-4 py-4">
        <div class="flex-shrink-0 rounded-full bg-gray-100 p-3 dark:bg-dark-700">
          <svg class="h-6 w-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
          </svg>
        </div>
        <div>
          <p class="font-medium text-gray-700 dark:text-gray-300">
            {{ t('profile.usageReport.featureDisabled') }}
          </p>
          <p class="text-sm text-gray-500 dark:text-gray-400">
            {{ t('profile.usageReport.featureDisabledHint') }}
          </p>
        </div>
      </div>

      <!-- Email not bound -->
      <div v-else-if="!emailBound" class="flex items-center gap-4 py-4">
        <div class="flex-shrink-0 rounded-full bg-amber-100 p-3 dark:bg-amber-900/30">
          <svg class="h-6 w-6 text-amber-600 dark:text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21.75 6.75v10.5a2.25 2.25 0 01-2.25 2.25h-15a2.25 2.25 0 01-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25m19.5 0v.243a2.25 2.25 0 01-1.07 1.916l-7.5 4.615a2.25 2.25 0 01-2.36 0L3.32 8.91a2.25 2.25 0 01-1.07-1.916V6.75" />
          </svg>
        </div>
        <div>
          <p class="font-medium text-amber-700 dark:text-amber-300">
            {{ t('profile.usageReport.emailRequired') }}
          </p>
          <p class="text-sm text-gray-500 dark:text-gray-400">
            {{ t('profile.usageReport.emailRequiredHint') }}
          </p>
        </div>
      </div>

      <!-- Main config form -->
      <div v-else class="space-y-6">
        <!-- Enable toggle -->
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <div class="flex-shrink-0 rounded-full p-3" :class="config?.enabled ? 'bg-green-100 dark:bg-green-900/30' : 'bg-gray-100 dark:bg-dark-700'">
              <svg class="h-6 w-6" :class="config?.enabled ? 'text-green-600 dark:text-green-400' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21.75 6.75v10.5a2.25 2.25 0 01-2.25 2.25h-15a2.25 2.25 0 01-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25m19.5 0v.243a2.25 2.25 0 01-1.07 1.916l-7.5 4.615a2.25 2.25 0 01-2.36 0L3.32 8.91a2.25 2.25 0 01-1.07-1.916V6.75" />
              </svg>
            </div>
            <div>
              <p class="font-medium text-gray-900 dark:text-white">
                {{ t('profile.usageReport.enableReport') }}
              </p>
              <p class="text-sm text-gray-500 dark:text-gray-400">
                {{ t('profile.usageReport.enableReportHint') }}
              </p>
            </div>
          </div>
          <label class="relative inline-flex items-center cursor-pointer">
            <input
              type="checkbox"
              v-model="formEnabled"
              :disabled="saving"
              class="sr-only peer"
              @change="handleToggle"
            />
            <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary-300 dark:peer-focus:ring-primary-800 rounded-full peer dark:bg-dark-600 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-dark-600 peer-checked:bg-primary-600"></div>
          </label>
        </div>

        <!-- Schedule settings (only show when enabled) -->
        <div v-if="formEnabled" class="space-y-4 pl-16">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {{ t('profile.usageReport.sendTime') }}
            </label>
            <input
              type="time"
              v-model="formSchedule"
              :disabled="saving"
              class="input w-40"
              @change="handleUpdate"
            />
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
              {{ t('profile.usageReport.sendTimeHint') }}
            </p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {{ t('profile.usageReport.timezone') }}
            </label>
            <select
              v-model="formTimezone"
              :disabled="saving"
              class="input w-64"
              @change="handleUpdate"
            >
              <option value="Asia/Shanghai">Asia/Shanghai (UTC+8)</option>
              <option value="Asia/Tokyo">Asia/Tokyo (UTC+9)</option>
              <option value="America/New_York">America/New_York (UTC-5)</option>
              <option value="America/Los_Angeles">America/Los_Angeles (UTC-8)</option>
              <option value="Europe/London">Europe/London (UTC+0)</option>
              <option value="Europe/Paris">Europe/Paris (UTC+1)</option>
            </select>
          </div>

          <!-- Test report button -->
          <div class="pt-4">
            <button
              type="button"
              class="btn btn-outline-primary"
              :disabled="sendingTest"
              @click="handleSendTest"
            >
              <span v-if="sendingTest" class="flex items-center gap-2">
                <svg class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                {{ t('profile.usageReport.sendingTest') }}
              </span>
              <span v-else>{{ t('profile.usageReport.sendTest') }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { usageReportAPI, type UsageReportConfig } from '@/api/usageReport'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const loading = ref(true)
const saving = ref(false)
const sendingTest = ref(false)
const config = ref<UsageReportConfig | null>(null)

const formEnabled = ref(false)
const formSchedule = ref('09:00')
const formTimezone = ref('Asia/Shanghai')

// Check if user has a valid email bound
const emailBound = computed(() => {
  const email = authStore.user?.email
  if (!email) return false
  if (email.endsWith('.invalid')) return false
  return true
})

const loadConfig = async () => {
  loading.value = true
  try {
    config.value = await usageReportAPI.getConfig()
    formEnabled.value = config.value.enabled
    formSchedule.value = config.value.schedule
    formTimezone.value = config.value.timezone
  } catch (error) {
    console.error('Failed to load usage report config:', error)
    appStore.showError(t('profile.usageReport.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleToggle = async () => {
  await handleUpdate()
}

const handleUpdate = async () => {
  saving.value = true
  try {
    const updated = await usageReportAPI.updateConfig({
      enabled: formEnabled.value,
      schedule: formSchedule.value,
      timezone: formTimezone.value
    })
    config.value = updated
    appStore.showSuccess(t('profile.usageReport.updateSuccess'))
  } catch (error: any) {
    // Revert to previous values
    if (config.value) {
      formEnabled.value = config.value.enabled
      formSchedule.value = config.value.schedule
      formTimezone.value = config.value.timezone
    }
    console.error('Failed to update usage report config:', error)
    const message = error.response?.data?.message || t('profile.usageReport.updateFailed')
    appStore.showError(message)
  } finally {
    saving.value = false
  }
}

const handleSendTest = async () => {
  sendingTest.value = true
  try {
    await usageReportAPI.sendTestReport()
    appStore.showSuccess(t('profile.usageReport.testSent'))
  } catch (error: any) {
    console.error('Failed to send test report:', error)
    const message = error.response?.data?.message || t('profile.usageReport.testFailed')
    appStore.showError(message)
  } finally {
    sendingTest.value = false
  }
}

onMounted(() => {
  loadConfig()
})
</script>
