import { ref, computed, onMounted, onUnmounted } from 'vue'

/**
 * 主题模式类型
 * - light: 始终日间模式
 * - dark: 始终夜间模式
 * - auto: 根据时间自动切换 (6:00-18:00 日间, 18:00-6:00 夜间)
 */
export type ThemeMode = 'light' | 'dark' | 'auto'

const THEME_KEY = 'theme'
const DAY_START_HOUR = 6
const DAY_END_HOUR = 18

// 全局状态，确保所有组件共享同一个状态
const isDark = ref(false)
const themeMode = ref<ThemeMode>('auto')
let checkInterval: ReturnType<typeof setInterval> | null = null
let initialized = false

/**
 * 判断当前时间是否为夜间
 */
function isNightTime(): boolean {
  const hour = new Date().getHours()
  return hour < DAY_START_HOUR || hour >= DAY_END_HOUR
}

/**
 * 应用主题到 DOM
 */
function applyTheme(dark: boolean): void {
  isDark.value = dark
  if (dark) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

/**
 * 根据当前模式更新主题
 */
function updateTheme(): void {
  switch (themeMode.value) {
    case 'light':
      applyTheme(false)
      break
    case 'dark':
      applyTheme(true)
      break
    case 'auto':
      applyTheme(isNightTime())
      break
  }
}

/**
 * 启动定时检查（auto 模式下使用）
 */
function startAutoCheck(): void {
  if (checkInterval) return
  // 每分钟检查一次
  checkInterval = setInterval(() => {
    if (themeMode.value === 'auto') {
      updateTheme()
    }
  }, 60 * 1000)
}

/**
 * 停止定时检查
 */
function stopAutoCheck(): void {
  if (checkInterval) {
    clearInterval(checkInterval)
    checkInterval = null
  }
}

/**
 * 初始化主题
 */
function initTheme(): void {
  if (initialized) return
  initialized = true

  const savedTheme = localStorage.getItem(THEME_KEY) as ThemeMode | null

  if (savedTheme && ['light', 'dark', 'auto'].includes(savedTheme)) {
    themeMode.value = savedTheme
  } else {
    // 默认使用 auto 模式
    themeMode.value = 'auto'
  }

  updateTheme()

  // 如果是 auto 模式，启动定时检查
  if (themeMode.value === 'auto') {
    startAutoCheck()
  }
}

/**
 * 主题管理 Composable
 */
export function useTheme() {
  /**
   * 设置主题模式
   */
  function setThemeMode(mode: ThemeMode): void {
    themeMode.value = mode
    localStorage.setItem(THEME_KEY, mode)
    updateTheme()

    if (mode === 'auto') {
      startAutoCheck()
    } else {
      stopAutoCheck()
    }
  }

  /**
   * 切换主题（在 light/dark/auto 之间循环）
   */
  function toggleTheme(): void {
    const modes: ThemeMode[] = ['light', 'dark', 'auto']
    const currentIndex = modes.indexOf(themeMode.value)
    const nextIndex = (currentIndex + 1) % modes.length
    setThemeMode(modes[nextIndex])
  }

  /**
   * 简单切换（仅在 light/dark 之间切换，会退出 auto 模式）
   */
  function toggleDarkMode(): void {
    setThemeMode(isDark.value ? 'light' : 'dark')
  }

  /**
   * 获取当前主题图标名称
   */
  const themeIcon = computed(() => {
    switch (themeMode.value) {
      case 'light':
        return 'sun'
      case 'dark':
        return 'moon'
      case 'auto':
        return 'clock'
      default:
        return 'sun'
    }
  })

  /**
   * 获取当前主题模式显示名称
   */
  const themeModeLabel = computed(() => {
    switch (themeMode.value) {
      case 'light':
        return 'nav.lightMode'
      case 'dark':
        return 'nav.darkMode'
      case 'auto':
        return 'nav.autoMode'
      default:
        return 'nav.lightMode'
    }
  })

  onMounted(() => {
    initTheme()
  })

  onUnmounted(() => {
    // 注意：不要在这里停止 autoCheck，因为其他组件可能还在使用
  })

  return {
    isDark,
    themeMode,
    themeIcon,
    themeModeLabel,
    setThemeMode,
    toggleTheme,
    toggleDarkMode,
    initTheme
  }
}

/**
 * 用于在 Vue 应用挂载前初始化主题，避免闪烁
 */
export function initThemeBeforeMount(): void {
  initTheme()
}
