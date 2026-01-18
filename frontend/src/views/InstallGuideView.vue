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
            {{ t('installGuide.title') }}
          </h1>
          <p class="mx-auto max-w-2xl text-base text-gray-600 dark:text-dark-300 sm:text-lg">
            {{ t('installGuide.subtitle') }}
          </p>
        </div>

        <!-- Tool Cards -->
        <div class="space-y-6 sm:space-y-8">
          <!-- Claude Code -->
          <div id="claude-code" class="card overflow-hidden scroll-mt-24">
            <div class="card-header flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <div class="flex items-center gap-3 sm:gap-4">
                <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-[#d97757] to-[#c45a3a] shadow-lg shadow-primary-500/30 sm:h-12 sm:w-12">
                  <img src="/llmLogo/Claude.png" alt="Claude" class="h-6 w-6 rounded object-contain sm:h-7 sm:w-7" />
                </div>
                <div>
                  <h2 class="text-lg font-semibold text-gray-900 dark:text-white sm:text-xl">Claude Code</h2>
                  <p class="text-xs text-gray-500 dark:text-dark-400 sm:text-sm">{{ t('installGuide.claudeCode.description') }}</p>
                </div>
              </div>
              <a
                href="https://docs.anthropic.com/en/docs/claude-code"
                target="_blank"
                rel="noopener noreferrer"
                class="btn btn-secondary btn-sm"
              >
                <Icon name="externalLink" size="sm" />
                {{ t('installGuide.officialDocs') }}
              </a>
            </div>
            <div class="card-body">
              <!-- OS Tabs -->
              <div class="mb-4 flex flex-wrap gap-2 sm:mb-6">
                <button
                  v-for="os in osList"
                  :key="os.id"
                  @click="claudeCodeOS = os.id"
                  class="flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs font-medium transition-all sm:gap-2 sm:px-4 sm:py-2 sm:text-sm"
                  :class="claudeCodeOS === os.id
                    ? 'bg-primary-500 text-white shadow-md shadow-primary-500/25'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-dark-700 dark:text-dark-300 dark:hover:bg-dark-600'"
                >
                  <component :is="os.icon" class="h-3.5 w-3.5 sm:h-4 sm:w-4" />
                  {{ os.name }}
                </button>
              </div>

              <!-- Installation Steps -->
              <div class="space-y-4 sm:space-y-6">
                <!-- Prerequisites -->
                <div>
                  <h3 class="mb-2 flex items-center gap-2 text-sm font-medium text-gray-900 dark:text-white sm:mb-3 sm:text-base">
                    <span class="flex h-5 w-5 items-center justify-center rounded-full bg-primary-100 text-xs font-semibold text-primary-600 dark:bg-primary-900/30 dark:text-primary-400 sm:h-6 sm:w-6">1</span>
                    {{ t('installGuide.prerequisites') }}
                  </h3>
                  <div class="rounded-xl border border-gray-200 bg-gray-50 p-3 dark:border-dark-700 dark:bg-dark-800/50 sm:p-4">
                    <ul class="space-y-1.5 text-xs text-gray-600 dark:text-dark-300 sm:space-y-2 sm:text-sm">
                      <li class="flex items-start gap-2">
                        <Icon name="check" size="sm" class="mt-0.5 flex-shrink-0 text-emerald-500" />
                        <span>Node.js 18+ ({{ t('installGuide.recommended') }} 20 LTS)</span>
                      </li>
                      <li class="flex items-start gap-2">
                        <Icon name="check" size="sm" class="mt-0.5 flex-shrink-0 text-emerald-500" />
                        <span>npm {{ t('installGuide.or') }} yarn {{ t('installGuide.or') }} pnpm</span>
                      </li>
                    </ul>
                  </div>
                </div>

                <!-- Install Command -->
                <div>
                  <h3 class="mb-2 flex items-center gap-2 text-sm font-medium text-gray-900 dark:text-white sm:mb-3 sm:text-base">
                    <span class="flex h-5 w-5 items-center justify-center rounded-full bg-primary-100 text-xs font-semibold text-primary-600 dark:bg-primary-900/30 dark:text-primary-400 sm:h-6 sm:w-6">2</span>
                    {{ t('installGuide.installCommand') }}
                  </h3>
                  <CodeBlock :code="claudeCodeInstallCommands[claudeCodeOS]" language="bash" />
                </div>

                <!-- Config -->
                <div>
                  <h3 class="mb-2 flex items-center gap-2 text-sm font-medium text-gray-900 dark:text-white sm:mb-3 sm:text-base">
                    <span class="flex h-5 w-5 items-center justify-center rounded-full bg-primary-100 text-xs font-semibold text-primary-600 dark:bg-primary-900/30 dark:text-primary-400 sm:h-6 sm:w-6">3</span>
                    {{ t('installGuide.configApiKey') }}
                  </h3>
                  <CodeBlock :code="claudeCodeConfigCommands[claudeCodeOS]" language="bash" />
                </div>

                <!-- Verify -->
                <div>
                  <h3 class="mb-2 flex items-center gap-2 text-sm font-medium text-gray-900 dark:text-white sm:mb-3 sm:text-base">
                    <span class="flex h-5 w-5 items-center justify-center rounded-full bg-primary-100 text-xs font-semibold text-primary-600 dark:bg-primary-900/30 dark:text-primary-400 sm:h-6 sm:w-6">4</span>
                    {{ t('installGuide.verifyInstall') }}
                  </h3>
                  <CodeBlock code="claude --version" language="bash" />
                </div>
              </div>
            </div>
          </div>

          <!-- OpenAI Codex CLI -->
          <div id="codex" class="card overflow-hidden scroll-mt-24">
            <div class="card-header flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <div class="flex items-center gap-3 sm:gap-4">
                <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-[#10a37f] to-[#0d8a6a] shadow-lg shadow-emerald-500/30 sm:h-12 sm:w-12">
                  <img src="/llmLogo/ChatGPT.png" alt="Codex" class="h-6 w-6 rounded object-contain sm:h-7 sm:w-7" />
                </div>
                <div>
                  <h2 class="text-lg font-semibold text-gray-900 dark:text-white sm:text-xl">OpenAI Codex CLI</h2>
                  <p class="text-xs text-gray-500 dark:text-dark-400 sm:text-sm">{{ t('installGuide.codex.description') }}</p>
                </div>
              </div>
              <a
                href="https://github.com/openai/codex"
                target="_blank"
                rel="noopener noreferrer"
                class="btn btn-secondary btn-sm"
              >
                <Icon name="externalLink" size="sm" />
                GitHub
              </a>
            </div>
            <div class="card-body">
              <!-- OS Tabs -->
              <div class="mb-4 flex flex-wrap gap-2 sm:mb-6">
                <button
                  v-for="os in osList"
                  :key="os.id"
                  @click="codexOS = os.id"
                  class="flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs font-medium transition-all sm:gap-2 sm:px-4 sm:py-2 sm:text-sm"
                  :class="codexOS === os.id
                    ? 'bg-emerald-500 text-white shadow-md shadow-emerald-500/25'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-dark-700 dark:text-dark-300 dark:hover:bg-dark-600'"
                >
                  <component :is="os.icon" class="h-3.5 w-3.5 sm:h-4 sm:w-4" />
                  {{ os.name }}
                </button>
              </div>

              <!-- Installation Steps -->
              <div class="space-y-4 sm:space-y-6">
                <!-- Prerequisites -->
                <div>
                  <h3 class="mb-2 flex items-center gap-2 text-sm font-medium text-gray-900 dark:text-white sm:mb-3 sm:text-base">
                    <span class="flex h-5 w-5 items-center justify-center rounded-full bg-emerald-100 text-xs font-semibold text-emerald-600 dark:bg-emerald-900/30 dark:text-emerald-400 sm:h-6 sm:w-6">1</span>
                    {{ t('installGuide.prerequisites') }}
                  </h3>
                  <div class="rounded-xl border border-gray-200 bg-gray-50 p-3 dark:border-dark-700 dark:bg-dark-800/50 sm:p-4">
                    <ul class="space-y-1.5 text-xs text-gray-600 dark:text-dark-300 sm:space-y-2 sm:text-sm">
                      <li class="flex items-start gap-2">
                        <Icon name="check" size="sm" class="mt-0.5 flex-shrink-0 text-emerald-500" />
                        <span>Node.js 22+</span>
                      </li>
                      <li class="flex items-start gap-2">
                        <Icon name="check" size="sm" class="mt-0.5 flex-shrink-0 text-emerald-500" />
                        <span>npm {{ t('installGuide.or') }} yarn</span>
                      </li>
                      <li v-if="codexOS === 'linux'" class="flex items-start gap-2">
                        <Icon name="check" size="sm" class="mt-0.5 flex-shrink-0 text-emerald-500" />
                        <span>{{ t('installGuide.codex.linuxSandbox') }}</span>
                      </li>
                    </ul>
                  </div>
                </div>

                <!-- Install Command -->
                <div>
                  <h3 class="mb-2 flex items-center gap-2 text-sm font-medium text-gray-900 dark:text-white sm:mb-3 sm:text-base">
                    <span class="flex h-5 w-5 items-center justify-center rounded-full bg-emerald-100 text-xs font-semibold text-emerald-600 dark:bg-emerald-900/30 dark:text-emerald-400 sm:h-6 sm:w-6">2</span>
                    {{ t('installGuide.installCommand') }}
                  </h3>
                  <CodeBlock :code="codexInstallCommands[codexOS]" language="bash" />
                </div>

                <!-- Config -->
                <div>
                  <h3 class="mb-2 flex items-center gap-2 text-sm font-medium text-gray-900 dark:text-white sm:mb-3 sm:text-base">
                    <span class="flex h-5 w-5 items-center justify-center rounded-full bg-emerald-100 text-xs font-semibold text-emerald-600 dark:bg-emerald-900/30 dark:text-emerald-400 sm:h-6 sm:w-6">3</span>
                    {{ t('installGuide.configApiKey') }}
                  </h3>
                  <CodeBlock :code="codexConfigCommands[codexOS]" language="bash" />
                </div>

                <!-- Verify -->
                <div>
                  <h3 class="mb-2 flex items-center gap-2 text-sm font-medium text-gray-900 dark:text-white sm:mb-3 sm:text-base">
                    <span class="flex h-5 w-5 items-center justify-center rounded-full bg-emerald-100 text-xs font-semibold text-emerald-600 dark:bg-emerald-900/30 dark:text-emerald-400 sm:h-6 sm:w-6">4</span>
                    {{ t('installGuide.verifyInstall') }}
                  </h3>
                  <CodeBlock code="codex --version" language="bash" />
                </div>
              </div>
            </div>
          </div>

          <!-- Gemini CLI -->
          <div id="gemini" class="card overflow-hidden scroll-mt-24">
            <div class="card-header flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <div class="flex items-center gap-3 sm:gap-4">
                <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-[#4285f4] to-[#1a73e8] shadow-lg shadow-blue-500/30 sm:h-12 sm:w-12">
                  <img src="/llmLogo/Gemini.jpg" alt="Gemini" class="h-6 w-6 rounded object-contain sm:h-7 sm:w-7" />
                </div>
                <div>
                  <h2 class="text-lg font-semibold text-gray-900 dark:text-white sm:text-xl">Gemini CLI</h2>
                  <p class="text-xs text-gray-500 dark:text-dark-400 sm:text-sm">{{ t('installGuide.gemini.description') }}</p>
                </div>
              </div>
              <a
                href="https://github.com/google-gemini/gemini-cli"
                target="_blank"
                rel="noopener noreferrer"
                class="btn btn-secondary btn-sm"
              >
                <Icon name="externalLink" size="sm" />
                GitHub
              </a>
            </div>
            <div class="card-body">
              <!-- OS Tabs -->
              <div class="mb-4 flex flex-wrap gap-2 sm:mb-6">
                <button
                  v-for="os in osList"
                  :key="os.id"
                  @click="geminiOS = os.id"
                  class="flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs font-medium transition-all sm:gap-2 sm:px-4 sm:py-2 sm:text-sm"
                  :class="geminiOS === os.id
                    ? 'bg-blue-500 text-white shadow-md shadow-blue-500/25'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-dark-700 dark:text-dark-300 dark:hover:bg-dark-600'"
                >
                  <component :is="os.icon" class="h-3.5 w-3.5 sm:h-4 sm:w-4" />
                  {{ os.name }}
                </button>
              </div>

              <!-- Installation Steps -->
              <div class="space-y-4 sm:space-y-6">
                <!-- Prerequisites -->
                <div>
                  <h3 class="mb-2 flex items-center gap-2 text-sm font-medium text-gray-900 dark:text-white sm:mb-3 sm:text-base">
                    <span class="flex h-5 w-5 items-center justify-center rounded-full bg-blue-100 text-xs font-semibold text-blue-600 dark:bg-blue-900/30 dark:text-blue-400 sm:h-6 sm:w-6">1</span>
                    {{ t('installGuide.prerequisites') }}
                  </h3>
                  <div class="rounded-xl border border-gray-200 bg-gray-50 p-3 dark:border-dark-700 dark:bg-dark-800/50 sm:p-4">
                    <ul class="space-y-1.5 text-xs text-gray-600 dark:text-dark-300 sm:space-y-2 sm:text-sm">
                      <li class="flex items-start gap-2">
                        <Icon name="check" size="sm" class="mt-0.5 flex-shrink-0 text-emerald-500" />
                        <span>Node.js 20+</span>
                      </li>
                      <li class="flex items-start gap-2">
                        <Icon name="check" size="sm" class="mt-0.5 flex-shrink-0 text-emerald-500" />
                        <span>npm {{ t('installGuide.or') }} yarn {{ t('installGuide.or') }} pnpm</span>
                      </li>
                    </ul>
                  </div>
                </div>

                <!-- Install Command -->
                <div>
                  <h3 class="mb-2 flex items-center gap-2 text-sm font-medium text-gray-900 dark:text-white sm:mb-3 sm:text-base">
                    <span class="flex h-5 w-5 items-center justify-center rounded-full bg-blue-100 text-xs font-semibold text-blue-600 dark:bg-blue-900/30 dark:text-blue-400 sm:h-6 sm:w-6">2</span>
                    {{ t('installGuide.installCommand') }}
                  </h3>
                  <CodeBlock :code="geminiInstallCommands[geminiOS]" language="bash" />
                </div>

                <!-- Config -->
                <div>
                  <h3 class="mb-2 flex items-center gap-2 text-sm font-medium text-gray-900 dark:text-white sm:mb-3 sm:text-base">
                    <span class="flex h-5 w-5 items-center justify-center rounded-full bg-blue-100 text-xs font-semibold text-blue-600 dark:bg-blue-900/30 dark:text-blue-400 sm:h-6 sm:w-6">3</span>
                    {{ t('installGuide.configApiKey') }}
                  </h3>
                  <CodeBlock :code="geminiConfigCommands[geminiOS]" language="bash" />
                </div>

                <!-- Verify -->
                <div>
                  <h3 class="mb-2 flex items-center gap-2 text-sm font-medium text-gray-900 dark:text-white sm:mb-3 sm:text-base">
                    <span class="flex h-5 w-5 items-center justify-center rounded-full bg-blue-100 text-xs font-semibold text-blue-600 dark:bg-blue-900/30 dark:text-blue-400 sm:h-6 sm:w-6">4</span>
                    {{ t('installGuide.verifyInstall') }}
                  </h3>
                  <CodeBlock code="gemini --version" language="bash" />
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Tips Section -->
        <div class="mt-8 sm:mt-12">
          <div class="rounded-2xl border border-primary-200/50 bg-gradient-to-br from-primary-50 to-primary-100/50 p-4 dark:border-primary-800/50 dark:from-primary-900/20 dark:to-primary-800/10 sm:p-6">
            <div class="flex items-start gap-3 sm:gap-4">
              <div class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-primary-500 text-white sm:h-10 sm:w-10">
                <Icon name="infoCircle" size="md" />
              </div>
              <div>
                <h3 class="mb-2 text-base font-semibold text-gray-900 dark:text-white sm:text-lg">
                  {{ t('installGuide.tips.title') }}
                </h3>
                <ul class="space-y-1.5 text-xs text-gray-600 dark:text-dark-300 sm:space-y-2 sm:text-sm">
                  <li class="flex items-start gap-2">
                    <span class="mt-1.5 h-1 w-1 flex-shrink-0 rounded-full bg-primary-500 sm:mt-2"></span>
                    {{ t('installGuide.tips.tip1') }}
                  </li>
                  <li class="flex items-start gap-2">
                    <span class="mt-1.5 h-1 w-1 flex-shrink-0 rounded-full bg-primary-500 sm:mt-2"></span>
                    {{ t('installGuide.tips.tip2') }}
                  </li>
                  <li class="flex items-start gap-2">
                    <span class="mt-1.5 h-1 w-1 flex-shrink-0 rounded-full bg-primary-500 sm:mt-2"></span>
                    {{ t('installGuide.tips.tip3') }}
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
import { ref, computed, h, onMounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/stores'
import { useTheme } from '@/composables/useTheme'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import CodeBlock from '@/components/common/CodeBlock.vue'

const { t } = useI18n()
const route = useRoute()
const appStore = useAppStore()
const { isDark, themeMode, toggleTheme } = useTheme()

// Handle hash scroll on mount
onMounted(() => {
  nextTick(() => {
    const hash = route.hash
    if (hash) {
      const element = document.querySelector(hash)
      if (element) {
        element.scrollIntoView({ behavior: 'smooth' })
      }
    }
  })
})

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

// OS Icons
const WindowsIcon = {
  render() {
    return h('svg', { viewBox: '0 0 24 24', fill: 'currentColor' }, [
      h('path', { d: 'M3 5.548l7.024-.96v6.784H3V5.548zm0 12.904l7.024.96v-6.784H3v5.824zm7.984 1.088L21 21V12.628h-10.016v6.912zM11.984 4.912L21 3v8.372h-10.016V4.912z' })
    ])
  }
}

const MacOSIcon = {
  render() {
    return h('svg', { viewBox: '0 0 24 24', fill: 'currentColor' }, [
      h('path', { d: 'M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z' })
    ])
  }
}

const LinuxIcon = {
  render() {
    return h('svg', { viewBox: '0 0 24 24', fill: 'currentColor' }, [
      h('path', { d: 'M12.504 0c-.155 0-.311.001-.465.003-.653.014-1.297.067-1.934.152-1.277.17-2.515.46-3.677.941-.574.239-1.128.517-1.654.838-.524.322-1.019.686-1.479 1.09-.46.404-.884.849-1.267 1.329-.384.48-.727.995-1.026 1.539-.298.543-.553 1.117-.761 1.712-.208.595-.369 1.212-.48 1.845-.111.633-.172 1.284-.18 1.943-.008.659.035 1.327.133 1.992.098.665.248 1.325.449 1.97.201.645.452 1.273.751 1.878.299.604.646 1.182 1.038 1.725.393.543.83 1.051 1.305 1.518.476.466.989.89 1.535 1.269.545.379 1.121.711 1.719.992.598.281 1.217.51 1.852.684.636.174 1.288.292 1.946.353.657.061 1.32.063 1.976.01.657-.053 1.305-.163 1.937-.327.632-.164 1.247-.382 1.84-.651.593-.27 1.161-.59 1.699-.96.538-.37 1.044-.787 1.511-1.247.467-.46.894-.962 1.276-1.499.381-.537.715-1.107 1-1.705.283-.597.516-1.219.698-1.86.181-.64.309-1.298.384-1.964.074-.667.093-1.34.056-2.01-.038-.668-.131-1.33-.28-1.977-.149-.647-.353-1.277-.612-1.881-.26-.604-.573-1.18-.939-1.721-.365-.541-.782-1.046-1.245-1.507-.464-.461-.973-.877-1.519-1.242-.546-.365-1.128-.678-1.736-.937C16.26.623 15.612.424 14.945.283 14.278.142 13.592.054 12.902.02 12.77.013 12.637.007 12.504 0zm-.177 2.04c.096 0 .192.002.288.005.576.021 1.14.076 1.69.168.55.092 1.083.22 1.596.381.513.162 1.006.357 1.476.583.47.225.917.482 1.338.766.421.284.815.595 1.179.931.365.336.699.697 1 1.079.301.382.57.786.802 1.207.232.421.428.86.586 1.313.158.453.278.921.358 1.397.08.476.12.96.12 1.446 0 .487-.04.971-.12 1.447-.08.476-.2.944-.358 1.397-.158.453-.354.892-.586 1.313-.232.421-.5.825-.802 1.207-.301.382-.635.743-1 1.079-.364.336-.758.647-1.179.931-.421.284-.869.541-1.338.766-.47.226-.963.42-1.476.583-.513.161-1.046.289-1.596.381-.55.092-1.114.147-1.69.168-.192.007-.384.01-.576.01-.192 0-.384-.003-.576-.01-.576-.021-1.14-.076-1.69-.168-.55-.092-1.083-.22-1.596-.381-.513-.163-1.006-.357-1.476-.583-.47-.225-.917-.482-1.338-.766-.421-.284-.815-.595-1.179-.931-.365-.336-.699-.697-1-1.079-.301-.382-.57-.786-.802-1.207-.232-.421-.428-.86-.586-1.313-.158-.453-.278-.921-.358-1.397-.08-.476-.12-.96-.12-1.447 0-.486.04-.97.12-1.446.08-.476.2-.944.358-1.397.158-.453.354-.892.586-1.313.232-.421.5-.825.802-1.207.301-.382.635-.743 1-1.079.364-.336.758-.647 1.179-.931.421-.284.869-.541 1.338-.766.47-.226.963-.421 1.476-.583.513-.161 1.046-.289 1.596-.381.55-.092 1.114-.147 1.69-.168.096-.003.192-.005.288-.005h.288z' })
    ])
  }
}

// OS List
const osList = [
  { id: 'macos', name: 'macOS', icon: MacOSIcon },
  { id: 'windows', name: 'Windows', icon: WindowsIcon },
  { id: 'linux', name: 'Linux', icon: LinuxIcon }
]

// Selected OS for each tool
const claudeCodeOS = ref('macos')
const codexOS = ref('macos')
const geminiOS = ref('macos')

// Claude Code Commands
const claudeCodeInstallCommands: Record<string, string> = {
  macos: `# 使用 npm 全局安装
npm install -g @anthropic-ai/claude-code

# 或使用 Homebrew
brew install claude-code`,
  windows: `# 使用 npm 全局安装 (需要管理员权限)
npm install -g @anthropic-ai/claude-code

# 或使用 winget
winget install Anthropic.ClaudeCode`,
  linux: `# 使用 npm 全局安装
npm install -g @anthropic-ai/claude-code

# Ubuntu/Debian 用户也可以使用 snap
sudo snap install claude-code`
}

const claudeCodeConfigCommands: Record<string, string> = {
  macos: `# 设置 API Key (推荐使用环境变量)
export ANTHROPIC_API_KEY="your-api-key"

# 或添加到 ~/.zshrc 或 ~/.bashrc 永久生效
echo 'export ANTHROPIC_API_KEY="your-api-key"' >> ~/.zshrc
source ~/.zshrc`,
  windows: `# 设置环境变量 (PowerShell)
$env:ANTHROPIC_API_KEY="your-api-key"

# 永久设置 (需要管理员权限)
[Environment]::SetEnvironmentVariable("ANTHROPIC_API_KEY", "your-api-key", "User")`,
  linux: `# 设置 API Key
export ANTHROPIC_API_KEY="your-api-key"

# 添加到 ~/.bashrc 永久生效
echo 'export ANTHROPIC_API_KEY="your-api-key"' >> ~/.bashrc
source ~/.bashrc`
}

// Codex Commands
const codexInstallCommands: Record<string, string> = {
  macos: `# 使用 npm 全局安装
npm install -g @openai/codex`,
  windows: `# 使用 npm 全局安装 (PowerShell 管理员模式)
npm install -g @openai/codex`,
  linux: `# 安装沙箱依赖 (Ubuntu/Debian)
sudo apt-get update
sudo apt-get install -y bubblewrap

# 使用 npm 全局安装
npm install -g @openai/codex`
}

const codexConfigCommands: Record<string, string> = {
  macos: `# 设置 OpenAI API Key
export OPENAI_API_KEY="your-api-key"

# 添加到 ~/.zshrc 永久生效
echo 'export OPENAI_API_KEY="your-api-key"' >> ~/.zshrc
source ~/.zshrc`,
  windows: `# 设置环境变量 (PowerShell)
$env:OPENAI_API_KEY="your-api-key"

# 永久设置
[Environment]::SetEnvironmentVariable("OPENAI_API_KEY", "your-api-key", "User")`,
  linux: `# 设置 OpenAI API Key
export OPENAI_API_KEY="your-api-key"

# 添加到 ~/.bashrc 永久生效
echo 'export OPENAI_API_KEY="your-api-key"' >> ~/.bashrc
source ~/.bashrc`
}

// Gemini Commands
const geminiInstallCommands: Record<string, string> = {
  macos: `# 使用 npm 全局安装
npm install -g @google-gemini/cli

# 或使用 npx 直接运行 (无需安装)
npx @google-gemini/cli`,
  windows: `# 使用 npm 全局安装
npm install -g @google-gemini/cli

# 或使用 npx 直接运行
npx @google-gemini/cli`,
  linux: `# 使用 npm 全局安装
npm install -g @google-gemini/cli

# 或使用 npx 直接运行
npx @google-gemini/cli`
}

const geminiConfigCommands: Record<string, string> = {
  macos: `# 设置 Gemini API Key
export GEMINI_API_KEY="your-api-key"

# 添加到 ~/.zshrc 永久生效
echo 'export GEMINI_API_KEY="your-api-key"' >> ~/.zshrc
source ~/.zshrc`,
  windows: `# 设置环境变量 (PowerShell)
$env:GEMINI_API_KEY="your-api-key"

# 永久设置
[Environment]::SetEnvironmentVariable("GEMINI_API_KEY", "your-api-key", "User")`,
  linux: `# 设置 Gemini API Key
export GEMINI_API_KEY="your-api-key"

# 添加到 ~/.bashrc 永久生效
echo 'export GEMINI_API_KEY="your-api-key"' >> ~/.bashrc
source ~/.bashrc`
}
</script>
