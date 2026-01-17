<template>
  <div class="relative flex min-h-screen items-center justify-center overflow-hidden p-4">
    <!-- Background -->
    <div
      class="absolute inset-0 bg-gradient-to-br from-gray-50 via-primary-50/30 to-gray-100 dark:from-dark-950 dark:via-dark-900 dark:to-dark-950"
    ></div>

    <!-- Decorative Elements -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <!-- Gradient Orbs -->
      <div
        class="absolute -right-40 -top-40 h-80 w-80 rounded-full bg-primary-400/20 blur-3xl"
      ></div>
      <div
        class="absolute -bottom-40 -left-40 h-80 w-80 rounded-full bg-primary-500/15 blur-3xl"
      ></div>
      <div
        class="absolute left-1/2 top-1/2 h-96 w-96 -translate-x-1/2 -translate-y-1/2 rounded-full bg-primary-300/10 blur-3xl"
      ></div>

      <!-- Grid Pattern -->
      <div
        class="absolute inset-0 bg-[linear-gradient(rgba(217,119,87,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(217,119,87,0.03)_1px,transparent_1px)] bg-[size:64px_64px]"
      ></div>
    </div>

    <!-- Content Container -->
    <div class="relative z-10 w-full max-w-md">
      <!-- Logo/Brand -->
      <div class="mb-8 text-center">
        <!-- Custom Logo or Default Logo -->
        <div
          class="mb-4 inline-flex h-16 w-16 cursor-pointer items-center justify-center overflow-hidden rounded-2xl shadow-lg shadow-primary-500/30 transition-transform hover:scale-105"
          @click="handleLogoClick"
        >
          <img :src="currentLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
        </div>
        <h1 class="text-gradient mb-2 text-3xl font-bold">
          {{ siteName }}
        </h1>
        <p class="text-sm text-gray-500 dark:text-dark-400">
          {{ siteSubtitle }}
        </p>
      </div>

      <!-- Card Container -->
      <div class="card-glass rounded-2xl p-8 shadow-glass">
        <slot />
      </div>

      <!-- Footer Links -->
      <div class="mt-6 text-center text-sm">
        <slot name="footer" />
      </div>

      <!-- Copyright -->
      <div class="mt-8 text-center text-xs text-gray-400 dark:text-dark-500">
        &copy; {{ currentYear }} {{ siteName }}. All rights reserved.
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { getPublicSettings } from '@/api/auth'
import { sanitizeUrl } from '@/utils/url'

const router = useRouter()

const siteName = ref('Code80')
const siteLogo = ref('')
const siteLogoDark = ref('')
const siteSubtitle = ref('Subscription to API Conversion Platform')
const isDark = ref(document.documentElement.classList.contains('dark'))

// Watch for theme changes via MutationObserver
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

const currentYear = computed(() => new Date().getFullYear())

// Audio for horn sound
const busHornSound = ref<HTMLAudioElement | null>(null)

function playHorn() {
  if (!busHornSound.value) {
    busHornSound.value = new Audio('/audio/bus-horn.MP3')
    busHornSound.value.volume = 0.6
  }
  busHornSound.value.currentTime = 0
  busHornSound.value.play().catch(() => {
    // Ignore autoplay errors
  })
}

function handleLogoClick() {
  playHorn()
  router.push('/')
}

onMounted(async () => {
  setupThemeObserver()
  try {
    const settings = await getPublicSettings()
    siteName.value = settings.site_name || 'Code80'
    siteLogo.value = sanitizeUrl(settings.site_logo || '', { allowRelative: true })
    siteLogoDark.value = sanitizeUrl(settings.site_logo_dark || '', { allowRelative: true })
    siteSubtitle.value = settings.site_subtitle || 'Subscription to API Conversion Platform'
  } catch (error) {
    console.error('Failed to load public settings:', error)
  }
})
</script>

<style scoped>
.text-gradient {
  @apply bg-gradient-to-r from-primary-600 to-primary-500 bg-clip-text text-transparent;
}
</style>
