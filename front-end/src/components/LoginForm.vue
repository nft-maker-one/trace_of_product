<template>
  <div class="login-container">
    <!-- 语言选择器 -->
    <div class="language-switcher">
      <LanguageSelector mode="compact" />
    </div>

    <div class="left-panel">
      <div class="logo">
        <i class="fas fa-shield-alt"></i>
        <span>{{ t('login.logo') }}</span>
      </div>
      <h1>{{ t('login.welcomeBack') }}</h1>
      <p>{{ t('login.description') }}</p>

      <ul class="features">
        <li><i class="fas fa-lock"></i> <span>{{ t('login.features.security') }}</span></li>
        <li><i class="fas fa-bolt"></i> <span>{{ t('login.features.fast') }}</span></li>
        <li><i class="fas fa-user-check"></i> <span>{{ t('login.features.personalized') }}</span></li>
      </ul>
    </div>

    <div class="login-form">
      <div class="form-header">
        <h2>{{ t('login.signIn') }}</h2>
        <p>{{ t('login.enterCredentials') }}</p>
      </div>

      <div v-if="errorMessage" class="error-message">
        <i class="fas fa-exclamation-circle"></i>
        <span>{{ errorMessage }}</span>
      </div>

      <div v-if="successMessage" class="success-message">
        <i class="fas fa-check-circle"></i>
        <span>{{ successMessage }}</span>
      </div>

      <div class="input-group">
        <i class="fas fa-user"></i>
        <input
          v-model="username"
          type="text"
          :placeholder="t('login.usernamePlaceholder')"
          autocomplete="username"
          @keyup.enter="handleLogin"
        >
      </div>

      <div class="input-group">
        <i class="fas fa-lock"></i>
        <input
          v-model="password"
          type="password"
          :placeholder="t('login.passwordPlaceholder')"
          autocomplete="current-password"
          @keyup.enter="handleLogin"
        >
      </div>

      <div class="form-options">
        <label class="remember-me">
          <input v-model="rememberMe" type="checkbox">
          <span>{{ t('login.rememberMe') }}</span>
        </label>
        <a href="#" class="forgot-password" @click.prevent="handleForgotPassword">
          {{ t('login.forgotPassword') }}
        </a>
      </div>

      <button class="btn btn-primary" @click="handleLogin" :disabled="loading">
        <i class="fas fa-sign-in-alt"></i>
        <span v-if="!loading">{{ t('login.signIn') }}</span>
        <span v-else class="loading-text">{{ t('login.signingIn') }}</span>
        <div v-if="loading" class="loading-spinner"></div>
      </button>

      <div class="divider">
        <span>{{ t('login.or') }}</span>
      </div>

      <button class="btn btn-secondary" @click="handleRegister">
        <i class="fas fa-user-plus"></i>
        <span>{{ t('login.createAccount') }}</span>
      </button>

      <div class="register-link">
        {{ t('login.noAccount') }}
        <a href="#" @click.prevent="handleRegister">{{ t('login.signUpNow') }}</a>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import LanguageSelector from './LanguageSelector.vue'

const router = useRouter()
const authStore = useAuthStore()
const { t } = useI18n()

const username = ref('')
const password = ref('')
const rememberMe = ref(false)
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

onMounted(() => {
  if (localStorage.getItem('rememberUser') === 'true') {
    const savedUsername = localStorage.getItem('username')
    if (savedUsername) {
      username.value = savedUsername
      rememberMe.value = true
    }
  }
})

async function handleLogin() {
  // 重置消息
  errorMessage.value = ''
  successMessage.value = ''

  // 验证输入
  if (!username.value || !password.value) {
    errorMessage.value = t('login.errors.emptyFields')
    return
  }

  loading.value = true

  try {
    const result = await authStore.login(username.value, password.value)

    if (result.success) {
      // 记住用户偏好
      if (rememberMe.value) {
        localStorage.setItem('rememberUser', 'true')
        localStorage.setItem('username', username.value)
      } else {
        localStorage.removeItem('rememberUser')
        localStorage.removeItem('username')
      }

      // 显示成功消息
      successMessage.value = t('login.success.loginSuccessful')

      // 延迟跳转
      setTimeout(() => {
        router.push('/dashboard')
      }, 1500)
    } else {
      errorMessage.value = result.message || t('login.errors.loginFailed')
    }
  } catch (error) {
    errorMessage.value = t('login.errors.networkError')
    console.error('Login error:', error)
  } finally {
    loading.value = false
  }
}

function handleRegister() {
  // 路由跳转到注册页面
  router.push('/register')
}

function handleForgotPassword() {
  errorMessage.value = t('login.errors.contactAdmin')
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
  box-shadow: var(--shadow);
  background: white;
  animation: fadeIn 0.8s ease-out;
  margin: 0 auto;
  position: relative;
}

/* 语言选择器 */
.language-switcher {
  position: absolute;
  top: 20px;
  right: 20px;
  z-index: 10;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 左侧面板 */
.left-panel {
  flex: 1;
  background: linear-gradient(135deg, var(--primary) 0%, var(--secondary) 100%);
  color: white;
  padding: 50px 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.left-panel::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -20%;
  width: 200px;
  height: 200px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
}

.left-panel::after {
  content: '';
  position: absolute;
  bottom: -30%;
  left: -10%;
  width: 150px;
  height: 150px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 50%;
}

/* Logo样式 */
.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 30px;
  font-size: 24px;
  font-weight: 700;
  z-index: 1;
}

.logo i {
  font-size: 28px;
  background: rgba(255, 255, 255, 0.2);
  padding: 10px;
  border-radius: 10px;
}

/* 标题和描述 */
.left-panel h1 {
  font-size: 32px;
  margin-bottom: 15px;
  line-height: 1.2;
  z-index: 1;
}

.left-panel p {
  font-size: 16px;
  opacity: 0.9;
  line-height: 1.6;
  margin-bottom: 30px;
  z-index: 1;
}

/* 特性列表 */
.features {
  list-style: none;
  margin-top: 30px;
  z-index: 1;
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

/* 右侧登录表单 */
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

.input-group input:focus + i {
  color: var(--primary);
}

.input-group input::placeholder {
  color: var(--gray);
  opacity: 0.7;
}

/* 表单选项 */
.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 25px;
  font-size: 14px;
}

.remember-me {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--gray);
  cursor: pointer;
  user-select: none;
}

.remember-me input[type="checkbox"] {
  width: 16px;
  height: 16px;
  cursor: pointer;
  accent-color: var(--primary);
}

.forgot-password {
  color: var(--primary);
  text-decoration: none;
  font-weight: 500;
  transition: var(--transition);
}

.forgot-password:hover {
  color: var(--primary-dark);
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

/* 注册链接 */
.register-link {
  text-align: center;
  margin-top: 25px;
  color: var(--gray);
  font-size: 15px;
}

.register-link a {
  color: var(--primary);
  text-decoration: none;
  font-weight: 600;
  margin-left: 5px;
  transition: var(--transition);
}

.register-link a:hover {
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
  }

  .left-panel,
  .login-form {
    padding: 30px 20px;
  }

  .form-options {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }

  .forgot-password {
    align-self: flex-end;
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
}

/* 页面容器样式（如果需要包裹） */
:global(.login-page) {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

/* 全局样式补充 */
:global(*) {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

:global(body) {
  font-family: 'Segoe UI', system-ui, -apple-system, BlinkMacSystemFont, 'Roboto', sans-serif;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  color: var(--dark);
  line-height: 1.6;
}

/* 输入框placeholder颜色修复 */
:global(::placeholder) {
  color: #a0aec0;
  opacity: 0.8;
}

/* 复选框样式增强 */
:global(input[type="checkbox"]) {
  width: 16px;
  height: 16px;
  cursor: pointer;
  accent-color: var(--primary);
}
</style>
