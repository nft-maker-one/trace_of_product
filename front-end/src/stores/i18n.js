import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useI18nStore = defineStore('i18n', () => {
  // 状态
  const currentLanguage = ref(localStorage.getItem('preferred-language') || 'zh-CN')
  const supportedLanguages = ref([
    { code: 'zh-CN', name: '中文', flag: '🇨🇳' },
    { code: 'en-US', name: 'English', flag: '🇺🇸' }
  ])

  // 计算属性
  const currentLanguageInfo = computed(() => {
    return supportedLanguages.value.find(lang => lang.code === currentLanguage.value) || supportedLanguages.value[0]
  })

  const isChinese = computed(() => currentLanguage.value === 'zh-CN')
  const isEnglish = computed(() => currentLanguage.value === 'en-US')

  // 方法
  function setLanguage(languageCode) {
    if (supportedLanguages.value.some(lang => lang.code === languageCode)) {
      currentLanguage.value = languageCode
      localStorage.setItem('preferred-language', languageCode)
    }
  }

  function toggleLanguage() {
    const newLang = currentLanguage.value === 'zh-CN' ? 'en-US' : 'zh-CN'
    setLanguage(newLang)
  }

  function getLanguageName(code) {
    const lang = supportedLanguages.value.find(lang => lang.code === code)
    return lang ? lang.name : code
  }

  function getLanguageFlag(code) {
    const lang = supportedLanguages.value.find(lang => lang.code === code)
    return lang ? lang.flag : ''
  }

  return {
    // 状态
    currentLanguage,
    supportedLanguages,

    // 计算属性
    currentLanguageInfo,
    isChinese,
    isEnglish,

    // 方法
    setLanguage,
    toggleLanguage,
    getLanguageName,
    getLanguageFlag
  }
})
