import { createI18n } from 'vue-i18n'

const LOCALE_KEY = 'sub2api_locale'
const DEFAULT_LOCALE = 'en'
const SUPPORTED_LOCALES = ['en', 'zh'] as const

type LocaleCode = (typeof SUPPORTED_LOCALES)[number]

const localeLoaders: Record<LocaleCode, () => Promise<{ default: Record<string, unknown> }>> = {
  en: () => import('./locales/en'),
  zh: () => import('./locales/zh')
}

const normalizeLocale = (locale: string): LocaleCode => {
  return SUPPORTED_LOCALES.includes(locale as LocaleCode)
    ? (locale as LocaleCode)
    : DEFAULT_LOCALE
}

function getDefaultLocale(): string {
  // Check localStorage first
  const saved = localStorage.getItem(LOCALE_KEY)
  if (saved && SUPPORTED_LOCALES.includes(saved as LocaleCode)) {
    return saved as LocaleCode
  }

  // Check browser language
  const browserLang = navigator.language.toLowerCase()
  if (browserLang.startsWith('zh')) {
    return 'zh'
  }

  return DEFAULT_LOCALE
}

export const i18n = createI18n({
  legacy: false,
  locale: normalizeLocale(getDefaultLocale()),
  fallbackLocale: DEFAULT_LOCALE,
  messages: {},
  // 禁用 HTML 消息警告 - 引导步骤使用富文本内容（driver.js 支持 HTML）
  // 这些内容是内部定义的，不存在 XSS 风险
  warnHtmlMessage: false
})

export async function loadLocaleMessages(locale: string): Promise<void> {
  const normalizedLocale = normalizeLocale(locale)

  if (i18n.global.availableLocales.includes(normalizedLocale)) {
    return
  }

  const messages = (await localeLoaders[normalizedLocale]()).default
  i18n.global.setLocaleMessage(normalizedLocale, messages)
}

export function warmupLocaleMessages(locale: string): void {
  const normalizedLocale = normalizeLocale(locale)
  const warmup = () => {
    void loadLocaleMessages(normalizedLocale)
  }

  if (typeof window.requestIdleCallback === 'function') {
    window.requestIdleCallback(() => {
      warmup()
    }, { timeout: 3000 })
    return
  }

  window.setTimeout(warmup, 1500)
}

export async function setLocale(locale: string): Promise<void> {
  const normalizedLocale = normalizeLocale(locale)
  await loadLocaleMessages(normalizedLocale)

  i18n.global.locale.value = normalizedLocale
  localStorage.setItem(LOCALE_KEY, normalizedLocale)
  document.documentElement.setAttribute('lang', normalizedLocale)
}

export function getLocale(): string {
  return normalizeLocale(String(i18n.global.locale.value))
}

export const availableLocales = [
  { code: 'en', name: 'English', flag: '🇺🇸' },
  { code: 'zh', name: '中文', flag: '🇨🇳' }
]

export default i18n
