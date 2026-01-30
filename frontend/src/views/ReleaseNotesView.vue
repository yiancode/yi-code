<template>
  <div
    class="relative flex min-h-screen flex-col overflow-hidden bg-gradient-to-br from-gray-50 via-primary-50/30 to-gray-100 dark:from-dark-950 dark:via-dark-900 dark:to-dark-950"
  >
    <!-- Background Decorations -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <div
        class="absolute -right-40 -top-40 h-96 w-96 rounded-full bg-primary-400/20 blur-3xl"
      ></div>
      <div
        class="absolute -bottom-40 -left-40 h-96 w-96 rounded-full bg-primary-500/15 blur-3xl"
      ></div>
      <div
        class="absolute left-1/3 top-1/4 h-72 w-72 rounded-full bg-primary-300/10 blur-3xl"
      ></div>
      <div
        class="absolute bottom-1/4 right-1/4 h-64 w-64 rounded-full bg-primary-400/10 blur-3xl"
      ></div>
      <div
        class="absolute inset-0 bg-[linear-gradient(rgba(217,119,87,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(217,119,87,0.03)_1px,transparent_1px)] bg-[size:64px_64px]"
      ></div>
    </div>

    <!-- Header -->
    <header class="relative z-20 px-6 py-4">
      <nav class="mx-auto flex max-w-6xl items-center justify-between">
        <!-- Logo -->
        <router-link to="/home" class="flex items-center gap-3">
          <div class="h-10 w-10 overflow-hidden rounded-xl shadow-md">
            <img :src="currentLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <span class="text-lg font-semibold text-gray-900 dark:text-white">{{ siteName }}</span>
        </router-link>

        <!-- Nav Actions -->
        <div class="flex items-center gap-3">
          <!-- Language Switcher -->
          <LocaleSwitcher />

          <!-- Theme Toggle -->
          <button
            @click="toggleTheme"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t(`nav.${themeMode === 'light' ? 'darkMode' : themeMode === 'dark' ? 'autoMode' : 'lightMode'}`)"
          >
            <Icon v-if="themeMode === 'light'" name="sun" size="md" class="text-amber-500" />
            <Icon v-else-if="themeMode === 'dark'" name="moon" size="md" class="text-indigo-400" />
            <Icon v-else name="clock" size="md" class="text-emerald-500" />
          </button>

          <!-- Back to Home -->
          <router-link
            to="/home"
            class="inline-flex items-center rounded-full bg-gray-900 px-3 py-1 text-xs font-medium text-white transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
          >
            {{ t('common.back') }}
          </router-link>
        </div>
      </nav>
    </header>

    <!-- Main Content -->
    <main class="relative z-10 flex-1 px-4 py-8 sm:px-6 sm:py-16">
      <div class="mx-auto max-w-6xl">
        <!-- Page Title -->
        <div class="mb-8 text-center sm:mb-12">
          <h1 class="mb-3 text-3xl font-bold text-gray-900 dark:text-white sm:mb-4 sm:text-4xl lg:text-5xl">
            {{ t('releaseNotes.title') }}
          </h1>
          <p class="mx-auto max-w-2xl text-base text-gray-600 dark:text-dark-300 sm:text-lg">
            {{ t('releaseNotes.subtitle') }}
          </p>
        </div>

        <!-- Version List -->
        <div class="space-y-6 sm:space-y-8">
          <!-- Version Card -->
          <div
            v-for="version in versions"
            :key="version.version"
            class="card overflow-hidden"
          >
            <div class="card-header flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <div class="flex items-center gap-3 sm:gap-4">
                <div
                  class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br shadow-lg sm:h-12 sm:w-12"
                  :class="version.type === 'major' ? 'from-red-500 to-red-600 shadow-red-500/30' :
                           version.type === 'minor' ? 'from-blue-500 to-blue-600 shadow-blue-500/30' :
                           'from-emerald-500 to-emerald-600 shadow-emerald-500/30'"
                >
                  <Icon name="badge" size="md" class="text-white" />
                </div>
                <div>
                  <h2 class="text-lg font-semibold text-gray-900 dark:text-white sm:text-xl">
                    {{ version.version }}
                  </h2>
                  <p class="text-xs text-gray-500 dark:text-dark-400 sm:text-sm">
                    {{ formatDate(version.date) }}
                  </p>
                </div>
              </div>
              <span
                class="inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-medium"
                :class="version.type === 'major' ? 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400' :
                         version.type === 'minor' ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400' :
                         'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400'"
              >
                {{ t(`releaseNotes.types.${version.type}`) }}
              </span>
            </div>

            <div class="card-body">
              <!-- Features -->
              <div v-if="version.features.length > 0" class="mb-6">
                <h3 class="mb-3 flex items-center gap-2 text-base font-semibold text-gray-900 dark:text-white">
                  <Icon name="sparkles" size="sm" class="text-blue-500" />
                  {{ t('releaseNotes.sections.features') }}
                </h3>
                <ul class="space-y-2">
                  <li
                    v-for="(feature, index) in version.features"
                    :key="index"
                    class="flex items-start gap-3 text-sm text-gray-700 dark:text-dark-300"
                  >
                    <Icon name="check" size="sm" class="mt-0.5 flex-shrink-0 text-blue-500" />
                    <span>{{ feature }}</span>
                  </li>
                </ul>
              </div>

              <!-- Improvements -->
              <div v-if="version.improvements.length > 0" class="mb-6">
                <h3 class="mb-3 flex items-center gap-2 text-base font-semibold text-gray-900 dark:text-white">
                  <Icon name="arrowUp" size="sm" class="text-emerald-500" />
                  {{ t('releaseNotes.sections.improvements') }}
                </h3>
                <ul class="space-y-2">
                  <li
                    v-for="(improvement, index) in version.improvements"
                    :key="index"
                    class="flex items-start gap-3 text-sm text-gray-700 dark:text-dark-300"
                  >
                    <Icon name="check" size="sm" class="mt-0.5 flex-shrink-0 text-emerald-500" />
                    <span>{{ improvement }}</span>
                  </li>
                </ul>
              </div>

              <!-- Bug Fixes -->
              <div v-if="version.bugFixes.length > 0">
                <h3 class="mb-3 flex items-center gap-2 text-base font-semibold text-gray-900 dark:text-white">
                  <Icon name="exclamationCircle" size="sm" class="text-red-500" />
                  {{ t('releaseNotes.sections.bugFixes') }}
                </h3>
                <ul class="space-y-2">
                  <li
                    v-for="(fix, index) in version.bugFixes"
                    :key="index"
                    class="flex items-start gap-3 text-sm text-gray-700 dark:text-dark-300"
                  >
                    <Icon name="check" size="sm" class="mt-0.5 flex-shrink-0 text-red-500" />
                    <span>{{ fix }}</span>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-gray-200/50 px-6 py-8 dark:border-dark-800/50">
      <div
        class="mx-auto flex max-w-6xl flex-col items-center justify-center gap-4 text-center sm:flex-row sm:text-left"
      >
        <p class="text-sm text-gray-500 dark:text-dark-400">
          &copy; {{ currentYear }} {{ siteName }}. {{ t('home.footer.allRightsReserved') }}
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import { useTheme } from '@/composables/useTheme'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'

const { t, locale } = useI18n()
const appStore = useAppStore()
const { isDark, themeMode, toggleTheme } = useTheme()

// Site settings
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Code80')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const siteLogoDark = computed(() => appStore.cachedPublicSettings?.site_logo_dark || appStore.siteLogoDark || '')

const currentLogo = computed(() => {
  if (isDark.value && siteLogoDark.value) {
    return siteLogoDark.value
  }
  return siteLogo.value
})

// Current year for footer
const currentYear = computed(() => new Date().getFullYear())

// Version data structure
interface Version {
  version: string
  date: string
  type: 'major' | 'minor' | 'patch'
  features: string[]
  improvements: string[]
  bugFixes: string[]
}

// Release versions (sample data - replace with actual release notes)
const versions = computed<Version[]>(() => {
  if (locale.value === 'zh') {
    return [
      {
        version: 'v0.1.64',
        date: '2026-01-30',
        type: 'minor',
        features: [
          '增加用户每日使用报告邮件功能',
          '支持管理员配置报告发送时间和邮件模板',
          '新增使用报告测试发送功能'
        ],
        improvements: [
          '优化邮件发送队列性能',
          '改进使用量统计精度',
          '增强邮件模板可定制性'
        ],
        bugFixes: [
          '修复使用报告配置保存后状态重置的问题',
          '修复SSE流式响应中usage数据被覆盖的问题',
          '修复邮件发送失败时的错误处理'
        ]
      },
      {
        version: 'v0.1.63',
        date: '2026-01-25',
        type: 'minor',
        features: [
          '添加订阅购买功能（iframe集成）',
          '支持微信支付配置',
          '新增订单管理系统'
        ],
        improvements: [
          '优化前端构建性能',
          '改进主题切换体验',
          '增强移动端响应式设计'
        ],
        bugFixes: [
          '修复订单创建时的并发问题',
          '修复支付回调处理异常',
          '修复订阅状态同步延迟'
        ]
      },
      {
        version: 'v0.1.60',
        date: '2026-01-20',
        type: 'minor',
        features: [
          'Linux.DO OAuth登录集成',
          '添加运维监控面板',
          '新增错误日志分析功能'
        ],
        improvements: [
          '优化API网关性能',
          '改进账户调度算法',
          '增强安全策略配置'
        ],
        bugFixes: [
          '修复OAuth回调处理问题',
          '修复并发限制计数错误',
          '修复会话粘性失效问题'
        ]
      }
    ]
  } else {
    return [
      {
        version: 'v0.1.64',
        date: '2026-01-30',
        type: 'minor',
        features: [
          'Added daily usage report email feature',
          'Admin can configure report sending time and email template',
          'Added test send function for usage reports'
        ],
        improvements: [
          'Optimized email queue performance',
          'Improved usage statistics accuracy',
          'Enhanced email template customization'
        ],
        bugFixes: [
          'Fixed usage report configuration reset issue after saving',
          'Fixed usage data overwrite issue in SSE streaming response',
          'Fixed error handling when email sending fails'
        ]
      },
      {
        version: 'v0.1.63',
        date: '2026-01-25',
        type: 'minor',
        features: [
          'Added subscription purchase feature (iframe integration)',
          'WeChat Pay configuration support',
          'New order management system'
        ],
        improvements: [
          'Optimized frontend build performance',
          'Improved theme switching experience',
          'Enhanced mobile responsive design'
        ],
        bugFixes: [
          'Fixed concurrency issue in order creation',
          'Fixed payment callback handling exception',
          'Fixed subscription status sync delay'
        ]
      },
      {
        version: 'v0.1.60',
        date: '2026-01-20',
        type: 'minor',
        features: [
          'Linux.DO OAuth login integration',
          'Added ops monitoring dashboard',
          'New error log analysis feature'
        ],
        improvements: [
          'Optimized API gateway performance',
          'Improved account scheduling algorithm',
          'Enhanced security policy configuration'
        ],
        bugFixes: [
          'Fixed OAuth callback handling issue',
          'Fixed concurrent limit counting error',
          'Fixed sticky session failure issue'
        ]
      }
    ]
  }
})

// Format date based on locale
function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  if (locale.value === 'zh') {
    return date.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
  }
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })
}
</script>
