import { createI18n } from 'vue-i18n'
import zhCN from './locales/zh-CN'
import enUS from './locales/en-US'

// 从本地存储或默认获取语言设置
const getBrowserLanguage = () => {
  const storedLang = localStorage.getItem('preferred-language')
  if (storedLang) return storedLang

  const browserLang = navigator.language.toLowerCase()
  if (browserLang.startsWith('zh')) return 'zh-CN'
  if (browserLang.startsWith('en')) return 'en-US'

  return 'zh-CN' // 默认中文
}

const i18n = createI18n({
  legacy: false, // 使用Composition API模式
  locale: getBrowserLanguage(),
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS
  },
  globalInjection: true // 全局注入 $t
})

// 将i18n实例挂载到window对象，供store使用
if (typeof window !== 'undefined') {
  window.$i18n = i18n.global
}

export default i18n
