import { createI18n } from 'vue-i18n'
import nexusmindOverlay from './nexusmind-i18n' // nexusmind

type LocaleCode = 'en' | 'zh' | 'pt-BR' // nexusmind

type LocaleMessages = Record<string, any>

const LOCALE_KEY = 'sub2api_locale'
const DEFAULT_LOCALE: LocaleCode = 'pt-BR' // nexusmind: pt-BR default for the Brazilian audience

const localeLoaders: Record<LocaleCode, () => Promise<{ default: LocaleMessages }>> = {
  en: () => import('./locales/en'),
  zh: () => import('./locales/zh'),
  'pt-BR': () => import('./locales/pt-BR') // nexusmind
}

function isLocaleCode(value: string): value is LocaleCode {
  return value === 'en' || value === 'zh' || value === 'pt-BR' // nexusmind
}

function getDefaultLocale(): LocaleCode {
  const saved = localStorage.getItem(LOCALE_KEY)
  if (saved && isLocaleCode(saved)) {
    return saved
  }

  const browserLang = navigator.language.toLowerCase()
  if (browserLang.startsWith('zh')) {
    return 'zh'
  }

  // nexusmind
  if (browserLang.startsWith('pt')) {
    return 'pt-BR'
  }

  return DEFAULT_LOCALE
}

export const i18n = createI18n({
  legacy: false,
  locale: getDefaultLocale(),
  fallbackLocale: DEFAULT_LOCALE,
  messages: {},
  // 禁用 HTML 消息警告 - 引导步骤使用富文本内容（driver.js 支持 HTML）
  // 这些内容是内部定义的，不存在 XSS 风险
  warnHtmlMessage: false
})

const loadedLocales = new Set<LocaleCode>()

// nexusmind: deep-merge the NexusMind overlay (Ciabra keys) into a locale so
// upstream en.ts / zh.ts never need editing.
function isPlainObject(value: unknown): value is LocaleMessages {
  return !!value && typeof value === 'object' && !Array.isArray(value)
}

function deepMerge(base: LocaleMessages, patch: LocaleMessages): LocaleMessages {
  const result: LocaleMessages = { ...base }
  for (const [key, value] of Object.entries(patch)) {
    result[key] =
      isPlainObject(value) && isPlainObject(result[key])
        ? deepMerge(result[key], value)
        : value
  }
  return result
}

function mergeNexusmindOverlay(locale: LocaleCode, messages: LocaleMessages): LocaleMessages {
  const extra = nexusmindOverlay[locale]
  return extra ? deepMerge(messages, extra) : messages
}

export async function loadLocaleMessages(locale: LocaleCode): Promise<void> {
  if (loadedLocales.has(locale)) {
    return
  }

  const loader = localeLoaders[locale]
  const module = await loader()
  i18n.global.setLocaleMessage(locale, mergeNexusmindOverlay(locale, module.default)) // nexusmind
  loadedLocales.add(locale)
}

export async function initI18n(): Promise<void> {
  const current = getLocale()
  await loadLocaleMessages(current)
  document.documentElement.setAttribute('lang', current)
}

export async function setLocale(locale: string): Promise<void> {
  if (!isLocaleCode(locale)) {
    return
  }

  await loadLocaleMessages(locale)
  i18n.global.locale.value = locale
  localStorage.setItem(LOCALE_KEY, locale)
  document.documentElement.setAttribute('lang', locale)

  // 同步更新浏览器页签标题，使其跟随语言切换
  const { resolveDocumentTitle } = await import('@/router/title')
  const { default: router } = await import('@/router')
  const { useAppStore } = await import('@/stores/app')
  const route = router.currentRoute.value
  const appStore = useAppStore()
  document.title = resolveDocumentTitle(route.meta.title, appStore.siteName, route.meta.titleKey as string)
}

export function getLocale(): LocaleCode {
  const current = i18n.global.locale.value
  return isLocaleCode(current) ? current : DEFAULT_LOCALE
}

export const availableLocales = [
  { code: 'en', name: 'English', flag: '🇺🇸' },
  { code: 'zh', name: '中文', flag: '🇨🇳' },
  { code: 'pt-BR', name: 'Português (Brasil)', flag: '🇧🇷' } // nexusmind
] as const

export default i18n
