<template>
  <!-- Opening Animation Overlay -->
  <div
    v-if="showAnimation"
    class="fixed inset-0 z-50 overflow-hidden pointer-events-none"
    :class="isDark ? 'bg-dark-950/50' : 'bg-gray-50/50'"
  >
    <!-- Car (clickable for horn sound) -->
    <img
      ref="carRef"
      src="/car.png"
      alt="car"
      class="car-animation absolute h-20 w-auto cursor-pointer pointer-events-auto"
      :style="carStyle"
      @click="playHorn"
    />
    <!-- Falling Logos -->
    <img
      v-for="(logo, index) in fallingLogos"
      :key="index"
      :src="logo.src"
      :alt="logo.name"
      class="absolute h-10 w-10 object-contain rounded-lg shadow-lg"
      :style="logo.style"
    />
  </div>

  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="homeContent" class="min-h-screen">
    <!-- iframe mode -->
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting, XSS risk is acceptable -->
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Default Home Page -->
  <div
    v-else
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
        <!-- Logo (hidden during animation, clickable for horn sound) -->
        <div class="flex items-center">
          <div
            ref="siteLogoRef"
            class="h-10 w-10 overflow-hidden rounded-xl shadow-md transition-all duration-500 cursor-pointer"
            :class="showSiteLogo ? 'opacity-100 scale-100' : 'opacity-0 scale-75'"
            @click="playHorn"
          >
            <img :src="currentLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </div>
        </div>

        <!-- Nav Actions -->
        <div class="flex items-center gap-3">
          <!-- Language Switcher -->
          <LocaleSwitcher />

          <!-- Doc Link -->
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
          </a>

          <!-- Theme Toggle -->
          <button
            @click="toggleTheme"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>

          <!-- Login / Dashboard Button -->
          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="inline-flex items-center gap-1.5 rounded-full bg-gray-900 py-1 pl-1 pr-2.5 transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
          >
            <span
              class="flex h-5 w-5 items-center justify-center rounded-full bg-gradient-to-br from-primary-400 to-primary-600 text-[10px] font-semibold text-white"
            >
              {{ userInitial }}
            </span>
            <span class="text-xs font-medium text-white">{{ t('home.dashboard') }}</span>
            <svg
              class="h-3 w-3 text-gray-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M4.5 19.5l15-15m0 0H8.25m11.25 0v11.25"
              />
            </svg>
          </router-link>
          <router-link
            v-else
            to="/login"
            class="inline-flex items-center rounded-full bg-gray-900 px-3 py-1 text-xs font-medium text-white transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
          >
            {{ t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <!-- Main Content -->
    <main class="relative z-10 flex-1 px-6 py-16">
      <div class="mx-auto max-w-6xl">
        <!-- Hero Section - Left/Right Layout -->
        <div class="mb-12 flex flex-col items-center justify-between gap-12 lg:flex-row lg:gap-16">
          <!-- Left: Text Content (Hidden during animation) -->
          <div
            class="flex-1 text-center lg:text-left transition-all duration-700"
            :class="showHeroContent ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
          >
            <h1
              class="mb-4 text-4xl font-bold text-gray-900 dark:text-white md:text-5xl lg:text-6xl"
            >
              {{ siteName }}
            </h1>
            <p class="mb-8 text-lg text-gray-600 dark:text-dark-300 md:text-xl">
              {{ siteSubtitle }}
            </p>

            <!-- CTA Button -->
            <div>
              <router-link
                :to="isAuthenticated ? dashboardPath : '/login'"
                class="btn btn-primary px-8 py-3 text-base shadow-lg shadow-primary-500/30"
              >
                {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
                <Icon name="arrowRight" size="md" class="ml-2" :stroke-width="2" />
              </router-link>
            </div>
          </div>

          <!-- Right: Terminal Animation -->
          <div class="flex flex-1 justify-center lg:justify-end">
            <div class="terminal-container">
              <div class="terminal-window">
                <!-- Window header -->
                <div class="terminal-header">
                  <div class="terminal-buttons">
                    <span class="btn-close"></span>
                    <span class="btn-minimize"></span>
                    <span class="btn-maximize"></span>
                  </div>
                  <span class="terminal-title">terminal</span>
                </div>
                <!-- Terminal content -->
                <div class="terminal-body">
                  <div class="code-line line-1">
                    <span class="code-prompt">$</span>
                    <span class="code-cmd">curl</span>
                    <span class="code-flag">-X POST</span>
                    <span class="code-url">/v1/messages</span>
                  </div>
                  <div class="code-line line-2">
                    <span class="code-comment"># Routing to upstream...</span>
                  </div>
                  <div class="code-line line-3">
                    <span class="code-success">200 OK</span>
                    <span class="code-response">{ "content": "Hello!" }</span>
                  </div>
                  <div class="code-line line-4">
                    <span class="code-prompt">$</span>
                    <span class="cursor"></span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Feature Tags - Centered -->
        <div class="mb-12 flex flex-wrap items-center justify-center gap-4 md:gap-6">
          <div
            class="inline-flex items-center gap-2.5 rounded-full border border-gray-200/50 bg-white/80 px-5 py-2.5 shadow-sm backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-800/80"
          >
            <Icon name="swap" size="sm" class="text-primary-500" />
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{
              t('home.tags.subscriptionToApi')
            }}</span>
          </div>
          <div
            class="inline-flex items-center gap-2.5 rounded-full border border-gray-200/50 bg-white/80 px-5 py-2.5 shadow-sm backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-800/80"
          >
            <Icon name="shield" size="sm" class="text-primary-500" />
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{
              t('home.tags.stickySession')
            }}</span>
          </div>
          <div
            class="inline-flex items-center gap-2.5 rounded-full border border-gray-200/50 bg-white/80 px-5 py-2.5 shadow-sm backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-800/80"
          >
            <Icon name="chart" size="sm" class="text-primary-500" />
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{
              t('home.tags.realtimeBilling')
            }}</span>
          </div>
        </div>

        <!-- Features Grid -->
        <div class="mb-12 grid gap-6 md:grid-cols-3">
          <!-- Feature 1: Unified Gateway -->
          <div
            class="group rounded-2xl border border-gray-200/50 bg-white/60 p-6 backdrop-blur-sm transition-all duration-300 hover:shadow-xl hover:shadow-primary-500/10 dark:border-dark-700/50 dark:bg-dark-800/60"
          >
            <div
              class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-blue-500 to-blue-600 shadow-lg shadow-blue-500/30 transition-transform group-hover:scale-110"
            >
              <Icon name="server" size="lg" class="text-white" />
            </div>
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('home.features.unifiedGateway') }}
            </h3>
            <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-400">
              {{ t('home.features.unifiedGatewayDesc') }}
            </p>
          </div>

          <!-- Feature 2: Account Pool -->
          <div
            class="group rounded-2xl border border-gray-200/50 bg-white/60 p-6 backdrop-blur-sm transition-all duration-300 hover:shadow-xl hover:shadow-primary-500/10 dark:border-dark-700/50 dark:bg-dark-800/60"
          >
            <div
              class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-primary-500 to-primary-600 shadow-lg shadow-primary-500/30 transition-transform group-hover:scale-110"
            >
              <svg
                class="h-6 w-6 text-white"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M18 18.72a9.094 9.094 0 003.741-.479 3 3 0 00-4.682-2.72m.94 3.198l.001.031c0 .225-.012.447-.037.666A11.944 11.944 0 0112 21c-2.17 0-4.207-.576-5.963-1.584A6.062 6.062 0 016 18.719m12 0a5.971 5.971 0 00-.941-3.197m0 0A5.995 5.995 0 0012 12.75a5.995 5.995 0 00-5.058 2.772m0 0a3 3 0 00-4.681 2.72 8.986 8.986 0 003.74.477m.94-3.197a5.971 5.971 0 00-.94 3.197M15 6.75a3 3 0 11-6 0 3 3 0 016 0zm6 3a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0zm-13.5 0a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z"
                />
              </svg>
            </div>
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('home.features.multiAccount') }}
            </h3>
            <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-400">
              {{ t('home.features.multiAccountDesc') }}
            </p>
          </div>

          <!-- Feature 3: Billing & Quota -->
          <div
            class="group rounded-2xl border border-gray-200/50 bg-white/60 p-6 backdrop-blur-sm transition-all duration-300 hover:shadow-xl hover:shadow-primary-500/10 dark:border-dark-700/50 dark:bg-dark-800/60"
          >
            <div
              class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-purple-500 to-purple-600 shadow-lg shadow-purple-500/30 transition-transform group-hover:scale-110"
            >
              <svg
                class="h-6 w-6 text-white"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z"
                />
              </svg>
            </div>
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('home.features.balanceQuota') }}
            </h3>
            <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-400">
              {{ t('home.features.balanceQuotaDesc') }}
            </p>
          </div>
        </div>

        <!-- Supported Providers -->
        <div class="mb-8 text-center">
          <h2 class="mb-3 text-2xl font-bold text-gray-900 dark:text-white">
            {{ t('home.providers.title') }}
          </h2>
          <p class="text-sm text-gray-600 dark:text-dark-400">
            {{ t('home.providers.description') }}
          </p>
        </div>

        <div class="mb-16 flex flex-wrap items-center justify-center gap-4">
          <!-- Claude - Supported -->
          <div
            ref="providerClaudeRef"
            class="flex items-center gap-2 rounded-xl border border-primary-200 bg-white/60 px-5 py-3 ring-1 ring-primary-500/20 backdrop-blur-sm dark:border-primary-800 dark:bg-dark-800/60"
          >
            <img src="/llmLogo/Claude.png" alt="Claude" class="h-8 w-8 rounded-lg object-contain" />
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.claude') }}</span>
            <span
              class="rounded bg-primary-100 px-1.5 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >{{ t('home.providers.supported') }}</span
            >
          </div>
          <!-- GPT - Supported -->
          <div
            ref="providerGPTRef"
            class="flex items-center gap-2 rounded-xl border border-primary-200 bg-white/60 px-5 py-3 ring-1 ring-primary-500/20 backdrop-blur-sm dark:border-primary-800 dark:bg-dark-800/60"
          >
            <img src="/llmLogo/ChatGPT.png" alt="ChatGPT" class="h-8 w-8 rounded-lg object-contain" />
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">GPT</span>
            <span
              class="rounded bg-primary-100 px-1.5 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >{{ t('home.providers.supported') }}</span
            >
          </div>
          <!-- Gemini - Supported -->
          <div
            ref="providerGeminiRef"
            class="flex items-center gap-2 rounded-xl border border-primary-200 bg-white/60 px-5 py-3 ring-1 ring-primary-500/20 backdrop-blur-sm dark:border-primary-800 dark:bg-dark-800/60"
          >
            <img src="/llmLogo/Gemini.jpg" alt="Gemini" class="h-8 w-8 rounded-lg object-contain" />
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.gemini') }}</span>
            <span
              class="rounded bg-primary-100 px-1.5 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >{{ t('home.providers.supported') }}</span
            >
          </div>
          <!-- Antigravity - Supported -->
          <div
            ref="providerAntigravityRef"
            class="flex items-center gap-2 rounded-xl border border-primary-200 bg-white/60 px-5 py-3 ring-1 ring-primary-500/20 backdrop-blur-sm dark:border-primary-800 dark:bg-dark-800/60"
          >
            <img src="/llmLogo/Antigravity.jpg" alt="Antigravity" class="h-8 w-8 rounded-lg object-contain" />
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.antigravity') }}</span>
            <span
              class="rounded bg-primary-100 px-1.5 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >{{ t('home.providers.supported') }}</span
            >
          </div>
          <!-- More - Coming Soon -->
          <div
            class="flex items-center gap-2 rounded-xl border border-gray-200/50 bg-white/40 px-5 py-3 opacity-60 backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-800/40"
          >
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-gray-500 to-gray-600"
            >
              <span class="text-xs font-bold text-white">+</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.more') }}</span>
            <span
              class="rounded bg-gray-100 px-1.5 py-0.5 text-[10px] font-medium text-gray-500 dark:bg-dark-700 dark:text-dark-400"
              >{{ t('home.providers.soon') }}</span
            >
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
        <a
          v-if="docUrl"
          :href="docUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-dark-400 dark:hover:text-white"
        >
          {{ t('home.docs') }}
        </a>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, reactive, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

// Site settings - directly from appStore (already initialized from injected config)
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Code80')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const siteLogoDark = computed(() => appStore.cachedPublicSettings?.site_logo_dark || appStore.siteLogoDark || '')
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'AI API Gateway Platform')
const docUrl = computed(() => appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')

// Check if homeContent is a URL (for iframe display)
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

// Theme
const isDark = ref(document.documentElement.classList.contains('dark'))

// Watch for theme changes via MutationObserver (for external theme changes)
let themeObserver: MutationObserver | null = null

function setupThemeObserver() {
  themeObserver = new MutationObserver((mutations) => {
    for (const mutation of mutations) {
      if (mutation.attributeName === 'class') {
        isDark.value = document.documentElement.classList.contains('dark')
      }
    }
  })
  themeObserver.observe(document.documentElement, { attributes: true })
}

onUnmounted(() => {
  if (themeObserver) {
    themeObserver.disconnect()
    themeObserver = null
  }
})

// Current logo based on theme
const currentLogo = computed(() => {
  if (isDark.value && siteLogoDark.value) {
    return siteLogoDark.value
  }
  return siteLogo.value
})

// Auth state
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const userInitial = computed(() => {
  const user = authStore.user
  if (!user || !user.email) return ''
  return user.email.charAt(0).toUpperCase()
})

// Current year for footer
const currentYear = computed(() => new Date().getFullYear())

// Toggle theme
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

// Initialize theme
function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

// ==================== Opening Animation ====================
const showAnimation = ref(true)
const showHeroContent = ref(false)
const showSiteLogo = ref(false)
const carRef = ref<HTMLImageElement | null>(null)
const siteLogoRef = ref<HTMLElement | null>(null)

// Provider refs for target positions
const providerClaudeRef = ref<HTMLElement | null>(null)
const providerGPTRef = ref<HTMLElement | null>(null)
const providerGeminiRef = ref<HTMLElement | null>(null)
const providerAntigravityRef = ref<HTMLElement | null>(null)

// LLM Logo list - mapped to providers
const llmLogos = [
  { name: 'Claude', src: '/llmLogo/Claude.png', providerRef: () => providerClaudeRef.value },
  { name: 'ChatGPT', src: '/llmLogo/ChatGPT.png', providerRef: () => providerGPTRef.value },
  { name: 'Gemini', src: '/llmLogo/Gemini.jpg', providerRef: () => providerGeminiRef.value },
  { name: 'Antigravity', src: '/llmLogo/Antigravity.jpg', providerRef: () => providerAntigravityRef.value }
]

// Car position state
const carPosition = reactive({
  x: -150,
  y: 0,
  scale: 1,
  opacity: 1
})

const carStyle = computed(() => ({
  transform: `translateX(${carPosition.x}px) scale(${carPosition.scale})`,
  opacity: carPosition.opacity,
  top: `${carPosition.y}px`,
  left: '0'
}))

// Falling logos state
interface FallingLogo {
  name: string
  src: string
  style: {
    transform: string
    opacity: number
    left: string
    top: string
    transition: string
  }
}

const fallingLogos = ref<FallingLogo[]>([])

// Audio for car animation
const busDrivingSound = ref<HTMLAudioElement | null>(null)
const busHornSound = ref<HTMLAudioElement | null>(null)

// Initialize audio elements
function initAudio() {
  busDrivingSound.value = new Audio('/audio/bus-driving.MP3')
  busDrivingSound.value.volume = 0.5
  busHornSound.value = new Audio('/audio/bus-horn.MP3')
  busHornSound.value.volume = 0.6
}

// Play horn sound when car is clicked
function playHorn() {
  // Initialize audio if not already done
  if (!busHornSound.value) {
    busHornSound.value = new Audio('/audio/bus-horn.MP3')
    busHornSound.value.volume = 0.6
  }
  busHornSound.value.currentTime = 0
  busHornSound.value.play().catch(() => {
    // Ignore autoplay errors
  })
}

// Animation controller
function startAnimation() {
  const screenWidth = window.innerWidth

  // Initialize and play driving sound
  initAudio()
  if (busDrivingSound.value) {
    busDrivingSound.value.play().catch(() => {
      // Ignore autoplay errors (browser may block without user interaction)
    })
  }

  // Get site logo position for car path and target
  let headerY = 16 // Default fallback
  let logoTargetX = 24
  let logoTargetY = 16

  if (siteLogoRef.value) {
    const rect = siteLogoRef.value.getBoundingClientRect()
    // Position car at same level as logo top, with minimum of 8px from top
    headerY = Math.max(8, rect.top)
    logoTargetX = rect.left
    logoTargetY = rect.top
  }

  // Car path: Same horizontal level as site logo in header
  const startX = -150
  const midX = screenWidth * 0.55 // Stop point before going to logo

  carPosition.x = startX
  carPosition.y = headerY

  // Phase 1: Car drives from left to right along header
  const phase1Duration = 4500
  const phase1Start = Date.now()

  // Track when to drop each logo
  const dropPositions = [0.2, 0.35, 0.5, 0.65, 0.8]
  let droppedCount = 0

  function animatePhase1() {
    const elapsed = Date.now() - phase1Start
    const progress = Math.min(elapsed / phase1Duration, 1)
    const eased = 1 - Math.pow(1 - progress, 2)

    carPosition.x = startX + (midX - startX) * eased
    carPosition.y = headerY + Math.sin(progress * Math.PI * 6) * 2

    // Drop logos at specific positions
    while (droppedCount < dropPositions.length && progress >= dropPositions[droppedCount]) {
      dropSingleLogo(droppedCount)
      droppedCount++
    }

    if (progress < 1) {
      requestAnimationFrame(animatePhase1)
    } else {
      // Phase 2: Move car to site logo position
      setTimeout(() => {
        // Get actual logo position
        if (siteLogoRef.value) {
          const rect = siteLogoRef.value.getBoundingClientRect()
          logoTargetX = rect.left
          logoTargetY = rect.top
        }
        moveCarToLogo(logoTargetX, logoTargetY)
      }, 400)
    }
  }

  animatePhase1()

  // Drop a single logo - falls downward to provider
  function dropSingleLogo(index: number) {
    const logo = llmLogos[index]
    if (!logo) return

    const logoStartX = carPosition.x + 50
    const logoStartY = carPosition.y + 50

    // Get target provider position
    const providerEl = logo.providerRef()
    let targetLogoX = logoStartX
    let targetLogoY = window.innerHeight - 150

    if (providerEl) {
      const rect = providerEl.getBoundingClientRect()
      targetLogoX = rect.left + 20
      targetLogoY = rect.top + 20
    }

    const fallingLogo: FallingLogo = {
      name: logo.name,
      src: logo.src,
      style: {
        transform: `translate(0px, 0px) scale(1) rotate(0deg)`,
        opacity: 1,
        left: `${logoStartX}px`,
        top: `${logoStartY}px`,
        transition: 'none'
      }
    }

    fallingLogos.value.push(fallingLogo)
    const logoIndex = fallingLogos.value.length - 1  // Capture index immediately after push

    // Animate falling downward to provider
    nextTick(() => {
      const deltaX = targetLogoX - logoStartX
      const deltaY = targetLogoY - logoStartY
      const flyDuration = 1200
      const rotation = (Math.random() - 0.5) * 360

      setTimeout(() => {
        if (fallingLogos.value[logoIndex]) {
          fallingLogos.value[logoIndex].style = {
            ...fallingLogos.value[logoIndex].style,
            transform: `translate(${deltaX}px, ${deltaY}px) scale(0.6) rotate(${rotation}deg)`,
            opacity: 0,
            transition: `transform ${flyDuration}ms cubic-bezier(0.25, 0.46, 0.45, 0.94), opacity ${flyDuration * 0.6}ms ease-out ${flyDuration * 0.4}ms`
          }
        }
      }, 50)
    })
  }

  // Phase 2: Move car to site logo and show logo
  function moveCarToLogo(targetX: number, targetY: number) {
    const phase2Duration = 1200
    const phase2Start = Date.now()
    const carStartX = carPosition.x
    const carStartY = carPosition.y

    function animatePhase2() {
      const elapsed = Date.now() - phase2Start
      const progress = Math.min(elapsed / phase2Duration, 1)
      const eased = progress < 0.5 ? 4 * progress * progress * progress : 1 - Math.pow(-2 * progress + 2, 3) / 2

      carPosition.x = carStartX + (targetX - carStartX) * eased
      carPosition.y = carStartY + (targetY - carStartY) * eased
      carPosition.scale = 1 - eased * 0.6
      carPosition.opacity = 1 - eased

      if (progress < 1) {
        requestAnimationFrame(animatePhase2)
      } else {
        // Car disappeared, show site logo
        showSiteLogo.value = true

        // Play horn sound when animation ends
        playHorn()

        // End animation and show hero content
        setTimeout(() => {
          showAnimation.value = false
          showHeroContent.value = true
        }, 400)
      }
    }

    animatePhase2()
  }
}

onMounted(() => {
  initTheme()
  setupThemeObserver()

  // Check auth state
  authStore.checkAuth()

  // Ensure public settings are loaded (will use cache if already loaded from injected config)
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }

  // Start opening animation after a short delay
  setTimeout(() => {
    startAnimation()
  }, 300)
})
</script>

<style scoped>
/* Car Animation */
.car-animation {
  will-change: transform, opacity;
  filter: drop-shadow(0 8px 16px rgba(0, 0, 0, 0.3));
}

/* Terminal Container */
.terminal-container {
  position: relative;
  display: inline-block;
}

/* Terminal Window */
.terminal-window {
  width: 420px;
  background: linear-gradient(145deg, #1e293b 0%, #0f172a 100%);
  border-radius: 14px;
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.4),
    0 0 0 1px rgba(255, 255, 255, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  overflow: hidden;
  transform: perspective(1000px) rotateX(2deg) rotateY(-2deg);
  transition: transform 0.3s ease;
}

.terminal-window:hover {
  transform: perspective(1000px) rotateX(0deg) rotateY(0deg) translateY(-4px);
}

/* Terminal Header */
.terminal-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: rgba(30, 41, 59, 0.8);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.terminal-buttons {
  display: flex;
  gap: 8px;
}

.terminal-buttons span {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.btn-close {
  background: #ef4444;
}
.btn-minimize {
  background: #eab308;
}
.btn-maximize {
  background: #22c55e;
}

.terminal-title {
  flex: 1;
  text-align: center;
  font-size: 12px;
  font-family: ui-monospace, monospace;
  color: #64748b;
  margin-right: 52px;
}

/* Terminal Body */
.terminal-body {
  padding: 20px 24px;
  font-family: ui-monospace, 'Fira Code', monospace;
  font-size: 14px;
  line-height: 2;
}

.code-line {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  opacity: 0;
  animation: line-appear 0.5s ease forwards;
}

.line-1 {
  animation-delay: 0.3s;
}
.line-2 {
  animation-delay: 1s;
}
.line-3 {
  animation-delay: 1.8s;
}
.line-4 {
  animation-delay: 2.5s;
}

@keyframes line-appear {
  from {
    opacity: 0;
    transform: translateY(5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.code-prompt {
  color: #22c55e;
  font-weight: bold;
}
.code-cmd {
  color: #38bdf8;
}
.code-flag {
  color: #a78bfa;
}
.code-url {
  color: #d97757;
}
.code-comment {
  color: #64748b;
  font-style: italic;
}
.code-success {
  color: #22c55e;
  background: rgba(34, 197, 94, 0.15);
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 600;
}
.code-response {
  color: #fbbf24;
}

/* Blinking Cursor */
.cursor {
  display: inline-block;
  width: 8px;
  height: 16px;
  background: #22c55e;
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  0%,
  50% {
    opacity: 1;
  }
  51%,
  100% {
    opacity: 0;
  }
}

/* Dark mode adjustments */
:deep(.dark) .terminal-window {
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.6),
    0 0 0 1px rgba(217, 119, 87, 0.2),
    0 0 40px rgba(217, 119, 87, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}
</style>
