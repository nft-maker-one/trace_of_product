<template>
  <div class="login-container">
    <!-- 语言选择器 -->
    <div class="language-switcher">
      <LanguageSelector mode="compact" />
    </div>

    <div class="left-panel">
      <div class="logo">
        <i class="fas fa-shield-alt"></i>
        <span>{{ t('register.logo') }}</span>
      </div>
      <h1>{{ t('register.createAccount') }}</h1>
      <p>{{ t('register.description') }}</p>

      <ul class="features">
        <li><i class="fas fa-shield"></i> <span>{{ t('register.features.security') }}</span></li>
        <li><i class="fas fa-rocket"></i> <span>{{ t('register.features.fastSetup') }}</span></li>
        <li><i class="fas fa-chart-line"></i> <span>{{ t('register.features.advancedFeatures') }}</span></li>
        <li><i class="fas fa-headset"></i> <span>{{ t('register.features.support') }}</span></li>
      </ul>
    </div>

    <div class="login-form">
      <div class="form-header">
        <h2>{{ t('register.signUp') }}</h2>
        <p>{{ t('register.enterDetails') }}</p>
      </div>

      <div v-if="errorMessage" class="error-message">
        <i class="fas fa-exclamation-circle"></i>
        <span>{{ errorMessage }}</span>
      </div>

      <div v-if="successMessage" class="success-message">
        <i class="fas fa-check-circle"></i>
        <span>{{ successMessage }}</span>
      </div>

      <!-- 用户名 -->
      <div class="input-group">
        <i class="fas fa-user"></i>
        <input
          v-model="formData.username"
          type="text"
          :placeholder="t('register.usernamePlaceholder')"
          autocomplete="username"
        >
        <div v-if="errors.username" class="field-error">{{ errors.username }}</div>
      </div>

      <!-- 邮箱 -->
      <div class="input-group">
        <i class="fas fa-envelope"></i>
        <input
          v-model="formData.email"
          type="email"
          :placeholder="t('register.emailPlaceholder')"
          autocomplete="email"
        >
        <div v-if="errors.email" class="field-error">{{ errors.email }}</div>
      </div>

      <!-- 密码 -->
      <div class="input-group">
        <i class="fas fa-lock"></i>
        <input
          v-model="formData.password"
          :type="showPassword ? 'text' : 'password'"
          :placeholder="t('register.passwordPlaceholder')"
          autocomplete="new-password"
        >
        <i
          class="fas toggle-password"
          :class="showPassword ? 'fa-eye-slash' : 'fa-eye'"
          @click="showPassword = !showPassword"
        ></i>
        <div v-if="errors.password" class="field-error">{{ errors.password }}</div>
      </div>

      <!-- 确认密码 -->
      <div class="input-group">
        <i class="fas fa-lock"></i>
        <input
          v-model="formData.confirmPassword"
          :type="showConfirmPassword ? 'text' : 'password'"
          :placeholder="t('register.confirmPasswordPlaceholder')"
          autocomplete="new-password"
        >
        <i
          class="fas toggle-password"
          :class="showConfirmPassword ? 'fa-eye-slash' : 'fa-eye'"
          @click="showConfirmPassword = !showConfirmPassword"
        ></i>
        <div v-if="errors.confirmPassword" class="field-error">{{ errors.confirmPassword }}</div>
      </div>

      <!-- 可选：邀请码 -->
      <div v-if="showInviteCode" class="input-group">
        <i class="fas fa-ticket-alt"></i>
        <input
          v-model="formData.inviteCode"
          type="text"
          :placeholder="t('register.inviteCodePlaceholder')"
        >
      </div>

      <!-- 服务条款 -->
      <div class="terms-agreement">
        <label class="terms-checkbox">
          <input v-model="formData.agreeTerms" type="checkbox">
          <span>
            {{ t('register.agreeTo') }}
            <a href="#" @click.prevent="handleTermsClick">{{ t('register.termsOfService') }}</a>
            {{ t('register.and') }}
            <a href="#" @click.prevent="handlePrivacyClick">{{ t('register.privacyPolicy') }}</a>
          </span>
        </label>
        <div v-if="errors.agreeTerms" class="field-error">{{ errors.agreeTerms }}</div>
      </div>

      <button class="btn btn-primary" @click="handleRegister" :disabled="loading">
        <i class="fas fa-user-plus"></i>
        <span v-if="!loading">{{ t('register.createAccount') }}</span>
        <span v-else class="loading-text">{{ t('register.creatingAccount') }}</span>
        <div v-if="loading" class="loading-spinner"></div>
      </button>

      <div class="divider">
        <span>{{ t('register.alreadyHaveAccount') }}</span>
      </div>

      <button class="btn btn-secondary" @click="handleSwitchToLogin">
        <i class="fas fa-sign-in-alt"></i>
        <span>{{ t('register.signInInstead') }}</span>
      </button>

      <div class="login-link">
        {{ t('register.haveAccount') }}
        <a href="#" @click.prevent="handleSwitchToLogin">{{ t('register.signInNow') }}</a>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import LanguageSelector from './LanguageSelector.vue'

const authStore = useAuthStore()
const { t } = useI18n()

const emit = defineEmits(['switch-to-login'])

// 表单数据
const formData = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  inviteCode: '',
  agreeTerms: false
})

// 状态
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const showInviteCode = ref(false) // 可根据需要配置是否显示邀请码
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

// 错误信息
const errors = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  agreeTerms: ''
})

// 表单验证
function validateForm() {
  let isValid = true
  const errorMessages = {
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    agreeTerms: ''
  }

  // 用户名验证
  if (!formData.username.trim()) {
    errorMessages.username = t('register.errors.usernameRequired')
    isValid = false
  } else if (formData.username.length < 3) {
    errorMessages.username = t('register.errors.usernameTooShort')
    isValid = false
  } else if (!/^[a-zA-Z0-9_]+$/.test(formData.username)) {
    errorMessages.username = t('register.errors.usernameInvalid')
    isValid = false
  }

  // 邮箱验证
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!formData.email.trim()) {
    errorMessages.email = t('register.errors.emailRequired')
    isValid = false
  } else if (!emailRegex.test(formData.email)) {
    errorMessages.email = t('register.errors.emailInvalid')
    isValid = false
  }

  // 密码验证
  if (!formData.password) {
    errorMessages.password = t('register.errors.passwordRequired')
    isValid = false
  } else if (formData.password.length < 8) {
    errorMessages.password = t('register.errors.passwordTooShort')
    isValid = false
  } else if (!/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/.test(formData.password)) {
    errorMessages.password = t('register.errors.passwordWeak')
    isValid = false
  }

  // 确认密码验证
  if (!formData.confirmPassword) {
    errorMessages.confirmPassword = t('register.errors.confirmPasswordRequired')
    isValid = false
  } else if (formData.password !== formData.confirmPassword) {
    errorMessages.confirmPassword = t('register.errors.passwordsDontMatch')
    isValid = false
  }

  // 服务条款验证
  if (!formData.agreeTerms) {
    errorMessages.agreeTerms = t('register.errors.agreeTermsRequired')
    isValid = false
  }

  // 更新错误信息
  Object.keys(errors).forEach(key => {
    errors[key] = errorMessages[key]
  })

  return isValid
}

async function handleRegister() {
  // 重置消息
  errorMessage.value = ''
  successMessage.value = ''

  // 验证表单
  if (!validateForm()) {
    errorMessage.value = t('register.errors.formInvalid')
    return
  }

  loading.value = true

  try {
    const result = await authStore.register({
      username: formData.username,
      email: formData.email,
      password: formData.password,
      inviteCode: formData.inviteCode || null
    })

    if (result.success) {
      successMessage.value = t('register.success.accountCreated')

      // 延迟跳转或切换回登录
      setTimeout(() => {
        emit('switch-to-login')
        successMessage.value = t('register.success.pleaseLogin')
      }, 2000)
    } else {
      errorMessage.value = result.message || t('register.errors.registrationFailed')
    }
  } catch (error) {
    errorMessage.value = t('register.errors.networkError')
    console.error('Registration error:', error)
  } finally {
    loading.value = false
  }
}

function handleSwitchToLogin() {
  emit('switch-to-login')
}

function handleTermsClick() {
  // 打开服务条款页面或模态框
  window.open('/terms', '_blank')
}

function handlePrivacyClick() {
  // 打开隐私政策页面或模态框
  window.open('/privacy', '_blank')
}
</script>

<style scoped>
/* 样式变量 */
:root {
  --primary: #4361ee;
  --primary-dark: #3a56d4;
  --secondary: #7209b7;
  --success: #4cc9f0;
  --light: #f8f9fa;
  --dark: #212529;
  --gray: #6c757d;
  --light-gray: #e9ecef;
  --shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
  --shadow-hover: 0 15px 40px rgba(67, 97, 238, 0.15);
  --transition: all 0.3s ease;
}

/* 登录容器 */
.login-container {
  display: flex;
  width: 100%;
  max-width: 1000px;
  min-height: auto;
  border-radius: 20px;
  overflow: hidden;
  background: white;
  box-shadow: var(--shadow);
  margin: 40px auto;
  position: relative;
}

/* 语言选择器 */
.language-switcher {
  position: absolute;
  top: 20px;
  right: 20px;
  z-index: 10;
}

/* 左侧面板 */
.left-panel {
  flex: 1;
  padding: 50px 40px;
  background: linear-gradient(135deg, var(--primary) 0%, var(--secondary) 100%);
  color: white;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 30px;
}

.logo i {
  font-size: 32px;
}

.left-panel h1 {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 15px;
  line-height: 1.3;
}

.left-panel p {
  font-size: 16px;
  opacity: 0.95;
  margin-bottom: 30px;
  line-height: 1.6;
}

.features {
  list-style: none;
  padding: 0;
}

.features li {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 15px;
  font-size: 15px;
}

.features i {
  background: rgba(255, 255, 255, 0.2);
  padding: 8px;
  border-radius: 50%;
  width: 34px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 右侧注册表单 */
.login-form {
  flex: 1;
  padding: 50px 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  background: white;
}

/* 表单头部 */
.form-header {
  margin-bottom: 40px;
  text-align: center;
}

.form-header h2 {
  font-size: 28px;
  color: var(--dark);
  margin-bottom: 8px;
}

.form-header p {
  color: var(--gray);
  font-size: 15px;
}

/* 输入框组 */
.input-group {
  position: relative;
  margin-bottom: 25px;
}

.input-group i {
  position: absolute;
  left: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--gray);
  font-size: 18px;
  transition: var(--transition);
}

.input-group input {
  width: 100%;
  padding: 16px 16px 16px 50px;
  border: 2px solid var(--light-gray);
  border-radius: 12px;
  font-size: 16px;
  transition: var(--transition);
  background: white;
  color: var(--dark);
}

.input-group input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(67, 97, 238, 0.1);
}

.input-group input:focus ~ i {
  color: var(--primary);
}

.input-group input::placeholder {
  color: var(--gray);
  opacity: 0.7;
}

/* 密码切换按钮 */
.toggle-password {
  position: absolute;
  right: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--gray);
  cursor: pointer;
  z-index: 2;
  font-size: 16px;
  transition: var(--transition);
}

.toggle-password:hover {
  color: var(--primary);
}

/* 字段级错误提示 */
.field-error {
  color: #c33;
  font-size: 12px;
  margin-top: 5px;
  margin-left: 50px;
  animation: slideIn 0.3s ease;
}

/* 服务条款 */
.terms-agreement {
  margin-bottom: 25px;
  font-size: 14px;
}

.terms-checkbox {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  color: var(--gray);
  cursor: pointer;
  user-select: none;
  line-height: 1.4;
}

.terms-checkbox input[type="checkbox"] {
  margin-top: 3px;
  flex-shrink: 0;
  width: 16px;
  height: 16px;
  cursor: pointer;
  accent-color: var(--primary);
}

.terms-checkbox a {
  color: var(--primary);
  text-decoration: none;
  font-weight: 500;
}

.terms-checkbox a:hover {
  text-decoration: underline;
}

/* 按钮样式 */
.btn {
  padding: 16px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  width: 100%;
}

.btn-primary {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
  color: white;
  margin-bottom: 20px;
  box-shadow: 0 4px 15px rgba(67, 97, 238, 0.2);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: var(--shadow-hover);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: white;
  color: var(--primary);
  border: 2px solid var(--light-gray);
  margin-bottom: 20px;
}

.btn-secondary:hover {
  background: var(--light-gray);
  border-color: var(--gray);
}

/* 分隔线 */
.divider {
  display: flex;
  align-items: center;
  margin: 25px 0;
  color: var(--gray);
  font-size: 14px;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--light-gray);
}

.divider span {
  padding: 0 15px;
}

/* 登录链接 */
.login-link {
  text-align: center;
  margin-top: 25px;
  color: var(--gray);
  font-size: 15px;
}

.login-link a {
  color: var(--primary);
  text-decoration: none;
  font-weight: 600;
  margin-left: 5px;
  transition: var(--transition);
}

.login-link a:hover {
  text-decoration: underline;
}

/* 错误消息 */
.error-message {
  background: #fee;
  color: #c33;
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  border-left: 4px solid #c33;
  animation: slideIn 0.3s ease;
}

/* 成功消息 */
.success-message {
  background: #efffee;
  color: #2a8;
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  border-left: 4px solid #2a8;
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(-10px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

/* 加载状态 */
.loading-text {
  margin-left: 8px;
}

.loading-spinner {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s ease-in-out infinite;
  margin-left: 8px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* 响应式设计 */
@media (max-width: 992px) {
  .login-container {
    max-width: 800px;
  }

  .left-panel,
  .login-form {
    padding: 40px 30px;
  }
}

@media (max-width: 768px) {
  .login-container {
    flex-direction: column;
    max-width: 450px;
    min-height: auto;
  }

  .left-panel {
    padding: 40px 30px;
  }

  .login-form {
    padding: 40px 30px;
  }

  .left-panel h1 {
    font-size: 28px;
  }

  .left-panel p {
    font-size: 15px;
  }

  .form-header h2 {
    font-size: 24px;
  }
}

@media (max-width: 480px) {
  .login-container {
    border-radius: 16px;
    margin: 20px;
  }

  .left-panel,
  .login-form {
    padding: 30px 20px;
  }

  .logo {
    font-size: 20px;
  }

  .logo i {
    font-size: 24px;
    padding: 8px;
  }

  .left-panel h1 {
    font-size: 24px;
  }

  .features li {
    font-size: 14px;
  }

  .terms-checkbox {
    font-size: 13px;
    line-height: 1.3;
  }

  .field-error {
    margin-left: 0;
    padding-left: 50px;
  }
}
</style>
