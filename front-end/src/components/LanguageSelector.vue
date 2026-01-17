<template>
  <div class="language-selector">
    <!-- 简洁模式：只显示当前语言和切换按钮 -->
    <div v-if="mode === 'simple'" class="simple-mode">
      <button
        class="lang-toggle-btn"
        @click="toggleLanguage"
        :title="`Switch to ${currentLanguageInfo.name === '中文' ? 'English' : '中文'}`"
      >
        <span class="lang-flag">{{ currentLanguageInfo.flag }}</span>
        <span class="lang-code">{{ currentLanguageInfo.code.split('-')[1] }}</span>
        <i class="fas fa-exchange-alt toggle-icon"></i>
      </button>
    </div>

    <!-- 完整模式：下拉选择器 -->
    <div v-else-if="mode === 'dropdown'" class="dropdown-mode">
      <select
        v-model="selectedLanguage"
        @change="handleLanguageChange"
        class="lang-select"
        :class="size"
      >
        <option
          v-for="lang in supportedLanguages"
          :key="lang.code"
          :value="lang.code"
        >
          {{ lang.flag }} {{ lang.name }}
        </option>
      </select>
    </div>

    <!-- 按钮组模式：两个按钮 -->
    <div v-else-if="mode === 'buttons'" class="buttons-mode">
      <div class="lang-buttons">
        <button
          v-for="lang in supportedLanguages"
          :key="lang.code"
          class="lang-btn"
          :class="{ active: lang.code === currentLanguage }"
          @click="setLanguage(lang.code)"
          :title="lang.name"
        >
          <span class="lang-flag">{{ lang.flag }}</span>
          <span class="lang-text">{{ lang.name }}</span>
        </button>
      </div>
    </div>

    <!-- 紧凑模式：只有国旗 -->
    <div v-else-if="mode === 'compact'" class="compact-mode">
      <button
        class="compact-btn"
        @click="toggleLanguage"
        :title="currentLanguageInfo.name"
      >
        {{ currentLanguageInfo.flag }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18nStore } from '../stores/i18n'
import { useI18n } from 'vue-i18n'

// Props
const props = defineProps({
  mode: {
    type: String,
    default: 'simple',
    validator: (value) => ['simple', 'dropdown', 'buttons', 'compact'].includes(value)
  },
  size: {
    type: String,
    default: 'medium',
    validator: (value) => ['small', 'medium', 'large'].includes(value)
  }
})

// Store and i18n
const i18nStore = useI18nStore()
const { locale } = useI18n()

// 响应式数据
const selectedLanguage = ref(i18nStore.currentLanguage)

// 计算属性
const currentLanguage = computed(() => i18nStore.currentLanguage)
const currentLanguageInfo = computed(() => i18nStore.currentLanguageInfo)
const supportedLanguages = computed(() => i18nStore.supportedLanguages)

// 方法
function setLanguage(languageCode) {
  console.log('Setting language to:', languageCode)
  locale.value = languageCode
  localStorage.setItem('preferred-language', languageCode)
  i18nStore.setLanguage(languageCode)
}

function toggleLanguage() {
  console.log('Toggling language')
  const newLang = currentLanguage.value === 'zh-CN' ? 'en-US' : 'zh-CN'
  setLanguage(newLang)
}

function handleLanguageChange() {
  setLanguage(selectedLanguage.value)
}
</script>

<style scoped>
/* 基础样式变量 */
.language-selector {
  display: inline-block;
}

/* 简洁模式 */
.simple-mode .lang-toggle-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.9);
  border: 2px solid rgba(46, 125, 50, 0.2);
  border-radius: 25px;
  color: var(--dark, #212529);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.simple-mode .lang-toggle-btn:hover {
  border-color: var(--primary, #4361ee);
  background: white;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(67, 97, 238, 0.15);
}

.simple-mode .lang-flag {
  font-size: 16px;
}

.simple-mode .lang-code {
  font-weight: 600;
  color: var(--primary, #4361ee);
}

.simple-mode .toggle-icon {
  font-size: 12px;
  color: var(--gray, #6c757d);
  transition: transform 0.3s ease;
}

.simple-mode .lang-toggle-btn:hover .toggle-icon {
  transform: rotate(180deg);
}

/* 下拉模式 */
.dropdown-mode .lang-select {
  padding: 8px 12px;
  border: 2px solid rgba(46, 125, 50, 0.2);
  border-radius: 8px;
  background: white;
  color: var(--dark, #212529);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  outline: none;
  min-width: 120px;
}

.dropdown-mode .lang-select:focus {
  border-color: var(--primary, #4361ee);
  box-shadow: 0 0 0 3px rgba(67, 97, 238, 0.1);
}

.dropdown-mode .lang-select:hover {
  border-color: var(--primary, #4361ee);
}

/* 按钮组模式 */
.buttons-mode .lang-buttons {
  display: flex;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.buttons-mode .lang-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.9);
  border: none;
  color: var(--gray, #6c757d);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.buttons-mode .lang-btn:hover {
  background: white;
  color: var(--dark, #212529);
  transform: translateY(-1px);
}

.buttons-mode .lang-btn.active {
  background: var(--primary, #4361ee);
  color: white;
  font-weight: 600;
}

.buttons-mode .lang-flag {
  font-size: 16px;
}

.buttons-mode .lang-text {
  white-space: nowrap;
}

/* 紧凑模式 */
.compact-mode .compact-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border: 2px solid rgba(46, 125, 50, 0.2);
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.9);
  color: var(--dark, #212529);
  font-size: 18px;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.compact-mode .compact-btn:hover {
  border-color: var(--primary, #4361ee);
  background: white;
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(67, 97, 238, 0.15);
}

/* 尺寸变体 */
.lang-select.small {
  padding: 6px 10px;
  font-size: 13px;
  min-width: 100px;
}

.lang-select.large {
  padding: 10px 14px;
  font-size: 15px;
  min-width: 140px;
}

.lang-btn.small {
  padding: 6px 12px;
  font-size: 13px;
}

.lang-btn.large {
  padding: 10px 20px;
  font-size: 15px;
}

.compact-btn.small {
  width: 32px;
  height: 32px;
  font-size: 16px;
}

.compact-btn.large {
  width: 48px;
  height: 48px;
  font-size: 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .simple-mode .lang-toggle-btn {
    padding: 6px 10px;
    font-size: 13px;
  }

  .dropdown-mode .lang-select {
    min-width: 100px;
    font-size: 13px;
  }

  .buttons-mode .lang-btn {
    padding: 6px 12px;
    font-size: 13px;
  }

  .compact-mode .compact-btn {
    width: 36px;
    height: 36px;
    font-size: 16px;
  }
}

/* 深色主题支持 */
@media (prefers-color-scheme: dark) {
  .simple-mode .lang-toggle-btn,
  .buttons-mode .lang-btn,
  .compact-mode .compact-btn {
    background: rgba(33, 37, 41, 0.9);
    border-color: rgba(255, 255, 255, 0.2);
    color: #ffffff;
  }

  .simple-mode .lang-toggle-btn:hover,
  .buttons-mode .lang-btn:hover,
  .compact-mode .compact-btn:hover {
    background: rgba(33, 37, 41, 1);
    border-color: var(--primary, #4361ee);
  }

  .dropdown-mode .lang-select {
    background: rgba(33, 37, 41, 0.9);
    border-color: rgba(255, 255, 255, 0.2);
    color: #ffffff;
  }
}

/* 无障碍支持 */
@media (prefers-reduced-motion: reduce) {
  .simple-mode .lang-toggle-btn,
  .buttons-mode .lang-btn,
  .compact-mode .compact-btn,
  .dropdown-mode .lang-select {
    transition: none;
  }

  .simple-mode .toggle-icon {
    transition: none;
  }
}

/* 高对比度支持 */
@media (prefers-contrast: high) {
  .simple-mode .lang-toggle-btn,
  .buttons-mode .lang-btn,
  .compact-mode .compact-btn {
    border-width: 3px;
  }

  .dropdown-mode .lang-select {
    border-width: 3px;
  }
}
</style>
