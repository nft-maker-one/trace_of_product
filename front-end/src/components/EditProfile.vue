<template>
  <div class="profile-edit-container">
    <!-- 语言选择器 -->
    <div class="language-switcher">
      <LanguageSelector mode="compact" />
    </div>

    <div class="profile-container">
      <!-- 左侧信息面板 -->
      <div class="left-panel">
        <div class="profile-header">
          <div class="avatar-container" @click="triggerAvatarUpload">
            <div class="avatar-wrapper">
              <img
                v-if="formData.avatarUrl"
                :src="formData.avatarUrl"
                alt="Avatar"
                class="avatar-image"
              >
              <i v-else class="fas fa-user-circle avatar-placeholder"></i>
              <div class="avatar-overlay">
                <i class="fas fa-camera"></i>
                <span>{{ t('profile.changePhoto') }}</span>
              </div>
            </div>
          </div>

          <h2>{{ displayName }}</h2>
          <p class="user-email">{{ userData.email }}</p>
          <p class="member-since">{{ t('profile.memberSince') }} {{ formatDate(userData.createdAt) }}</p>
        </div>

        <div class="profile-stats">
          <div class="stat-item">
            <i class="fas fa-sign-in-alt"></i>
            <div class="stat-content">
              <span class="stat-label">{{ t('profile.lastLogin') }}</span>
              <span class="stat-value">{{ formatDate(userData.last_login_at) }}</span>
            </div>
          </div>

          <div class="stat-item">
            <i class="fas fa-shield-alt"></i>
            <div class="stat-content">
              <span class="stat-label">{{ t('profile.accountStatus') }}</span>
              <span class="stat-value status-active">{{ t('profile.active') }}</span>
            </div>
          </div>

          <div class="stat-item">
            <i class="fas fa-id-card"></i>
            <div class="stat-content">
              <span class="stat-label">{{ t('profile.userId') }}</span>
              <span class="stat-value">#{{ userData.id }}</span>
            </div>
          </div>
        </div>

        <div class="privacy-tips">
          <h3><i class="fas fa-lock"></i> {{ t('profile.privacyTips') }}</h3>
          <ul>
            <li><i class="fas fa-check-circle"></i> {{ t('profile.tipEmail') }}</li>
            <li><i class="fas fa-check-circle"></i> {{ t('profile.tipPassword') }}</li>
            <li><i class="fas fa-check-circle"></i> {{ t('profile.tipProfile') }}</li>
          </ul>
        </div>
      </div>

      <!-- 右侧编辑表单 -->
      <div class="edit-form">
        <div class="form-header">
          <h2>{{ t('profile.editProfile') }}</h2>
          <p>{{ t('profile.editDescription') }}</p>
        </div>

        <div v-if="errorMessage" class="error-message">
          <i class="fas fa-exclamation-circle"></i>
          <span>{{ errorMessage }}</span>
        </div>

        <div v-if="successMessage" class="success-message">
          <i class="fas fa-check-circle"></i>
          <span>{{ successMessage }}</span>
        </div>

        <!-- 基本信息部分 -->
        <div class="form-section">
          <h3><i class="fas fa-user"></i> {{ t('profile.basicInfo') }}</h3>

          <div class="input-group">
            <i class="fas fa-user"></i>
            <input
              v-model="formData.nickname"
              type="text"
              :placeholder="t('profile.usernamePlaceholder')"
              :class="{ 'has-error': errors.nickname }"
            >
            <div v-if="errors.nickname" class="field-error">{{ errors.nickname }}</div>
          </div>

          <div class="input-group">
            <i class="fas fa-envelope"></i>
            <input
              v-model="formData.email"
              type="email"
              :placeholder="t('profile.emailPlaceholder')"
              :class="{ 'has-error': errors.email }"
              disabled
            >
            <span class="field-note">{{ t('profile.emailNote') }}</span>
          </div>

        </div>

        <!-- 密码修改部分 -->
        <div class="form-section">
          <h3><i class="fas fa-lock"></i> {{ t('profile.security') }}</h3>

          <div class="input-group">
            <i class="fas fa-lock"></i>
            <input
              v-model="formData.current_password"
              :type="showCurrentPassword ? 'text' : 'password'"
              :placeholder="t('profile.currentPasswordPlaceholder')"
            >
            <i
              class="fas toggle-password"
              :class="showCurrentPassword ? 'fa-eye-slash' : 'fa-eye'"
              @click="showCurrentPassword = !showCurrentPassword"
            ></i>
          </div>

          <div class="input-group">
            <i class="fas fa-key"></i>
            <input
              v-model="formData.new_password"
              :type="showNewPassword ? 'text' : 'password'"
              :placeholder="t('profile.newPasswordPlaceholder')"
              :class="{ 'has-error': errors.new_password }"
            >
            <i
              class="fas toggle-password"
              :class="showNewPassword ? 'fa-eye-slash' : 'fa-eye'"
              @click="showNewPassword = !showNewPassword"
            ></i>
            <div v-if="errors.new_password" class="field-error">{{ errors.new_password }}</div>
          </div>

          <div class="input-group">
            <i class="fas fa-key"></i>
            <input
              v-model="formData.confirm_password"
              :type="showConfirmPassword ? 'text' : 'password'"
              :placeholder="t('profile.confirmPasswordPlaceholder')"
              :class="{ 'has-error': errors.confirm_password }"
            >
            <i
              class="fas toggle-password"
              :class="showConfirmPassword ? 'fa-eye-slash' : 'fa-eye'"
              @click="showConfirmPassword = !showConfirmPassword"
            ></i>
            <div v-if="errors.confirm_password" class="field-error">{{ errors.confirm_password }}</div>
          </div>

          <div class="password-strength" v-if="formData.new_password">
            <div class="strength-bar" :class="passwordStrengthClass"></div>
            <span class="strength-text">{{ passwordStrengthText }}</span>
          </div>
        </div>

        <!-- 隐私设置部分 -->
        <div class="form-section">
          <h3><i class="fas fa-user-secret"></i> {{ t('profile.privacy') }}</h3>

          <div class="privacy-options">
            <label class="privacy-option">
              <input v-model="formData.profile_public" type="checkbox">
              <div class="option-content">
                <i class="fas fa-globe-americas"></i>
                <div>
                  <span class="option-title">{{ t('profile.publicProfile') }}</span>
                  <span class="option-description">{{ t('profile.publicProfileDesc') }}</span>
                </div>
              </div>
            </label>

            <label class="privacy-option">
              <input v-model="formData.email_notifications" type="checkbox">
              <div class="option-content">
                <i class="fas fa-envelope"></i>
                <div>
                  <span class="option-title">{{ t('profile.emailNotifications') }}</span>
                  <span class="option-description">{{ t('profile.emailNotificationsDesc') }}</span>
                </div>
              </div>
            </label>
          </div>
        </div>

        <!-- 按钮组 -->
        <div class="form-actions">
          <button class="btn btn-secondary" @click="handleCancel" :disabled="loading">
            <i class="fas fa-times"></i>
            <span>{{ t('profile.cancel') }}</span>
          </button>

          <button class="btn btn-primary" @click="handleSave" :disabled="loading">
            <i class="fas fa-save"></i>
            <span v-if="!loading">{{ t('profile.saveChanges') }}</span>
            <span v-else class="loading-text">{{ t('profile.saving') }}</span>
            <div v-if="loading" class="loading-spinner"></div>
          </button>
        </div>
      </div>
    </div>

    <!-- 头像上传模态框 -->
    <div v-if="showAvatarModal" class="modal-overlay" @click="closeAvatarModal">
      <div class="modal-content" @click.stop>
        <h3>{{ t('profile.changePhoto') }}</h3>
        <div class="avatar-options">
          <div class="avatar-option" v-for="option in avatarOptions" :key="option.id"
               @click="selectAvatar(option.url)">
            <img :src="option.url" :alt="option.name">
            <span>{{ option.name }}</span>
          </div>
        </div>
        <div class="modal-actions">
          <button class="btn btn-secondary" @click="closeAvatarModal">
            {{ t('profile.cancel') }}
          </button>
          <button class="btn btn-primary" @click="confirmAvatarChange">
            {{ t('profile.confirm') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import LanguageSelector from './LanguageSelector.vue'
import { useAuthStore } from '../stores/auth'
import api from '@/services/api'

const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()

// 用户数据
const userData = ref(authStore.user || {})

// 会流数据
const formData = reactive({
  nickname: userData.value?.nickname || '',
  email: userData.value?.email || '',
  avatarUrl: userData.value?.avatarUrl || '/public/avatars/profile.jpg',
  current_password: '',
  new_password: '',
  confirm_password: '',
  profile_public: userData.value?.profile_public ?? true,
  email_notifications: userData.value?.email_notifications ?? true,
})

// 状态
const showCurrentPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const showAvatarModal = ref(false)

// 错误信息
const errors = reactive({
  nickname: '',
  email: '',
  phone: '',
  new_password: '',
  confirm_password: ''
})

// 计算属性
const displayName = computed(() => {
  return formData.nickname || userData.value?.nickname || userData.value?.username || '用户'
})

const passwordStrengthClass = computed(() => {
  const password = formData.new_password
  if (!password) return ''

  const strength = calculatePasswordStrength(password)
  return `strength-${strength}`
})

const passwordStrengthText = computed(() => {
  const password = formData.new_password
  if (!password) return ''

  const strength = calculatePasswordStrength(password)
  const texts = [
    t('profile.veryWeak'),
    t('profile.weak'),
    t('profile.fair'),
    t('profile.good'),
    t('profile.strong')
  ]
  return texts[strength]
})

// 头像选项
const avatarOptions = ref([
  { id: 1, name: t('profile.avatarDefault'), url: '/public/avatars/profile.jpg' },
  { id: 2, name: t('profile.avatarMale'), url: '/public/avatars/boy.png' },
  { id: 3, name: t('profile.avatarFemale'), url: '/public/avatars/girl.png' },
  { id: 4, name: t('profile.avatarCustom'), url: '/public/avatars/image.jpg' }
])

function calculatePasswordStrength(password) {
  let score = 0

  // 长度
  if (password.length >= 8) score++
  if (password.length >= 12) score++

  // 字符类型
  if (/[a-z]/.test(password)) score++
  if (/[A-Z]/.test(password)) score++
  if (/[0-9]/.test(password)) score++
  if (/[^A-Za-z0-9]/.test(password)) score++

  return Math.min(Math.floor(score / 2), 4) // 0-4 分
}

function validateForm() {
  let isValid = true
  const errorMsgs = {
    nickname: '',
    email: '',
    phone: '',
    new_password: '',
    confirm_password: ''
  }

  // 昵称验证
  if (!formData.nickname.trim()) {
    errorMsgs.nickname = t('profile.errors.usernameRequired')
    isValid = false
  } else if (formData.nickname.length < 3) {
    errorMsgs.nickname = t('profile.errors.usernameTooShort')
    isValid = false
  }

  // 邮箱验证（只读，不验证）

  // 密码验证（如果提供了当前密码）
  if (formData.current_password) {
    if (!formData.new_password) {
      errorMsgs.new_password = t('profile.errors.newPasswordRequired')
      isValid = false
    } else if (formData.new_password.length < 8) {
      errorMsgs.new_password = t('profile.errors.passwordTooShort')
      isValid = false
    } else if (formData.new_password !== formData.confirm_password) {
      errorMsgs.confirm_password = t('profile.errors.passwordsDontMatch')
      isValid = false
    }
  }

  // 更新错误信息
  Object.keys(errors).forEach(key => {
    errors[key] = errorMsgs[key]
  })

  return isValid
}

function handleSave() {
  console.log('保存按钮被点击')
  // 重置消息
  errorMessage.value = ''
  successMessage.value = ''
  console.log('准备保存用户资料:', formData)

  // 验证表单
  if (!validateForm()) {
    errorMessage.value = t('profile.errors.formInvalid')
    return
  }

  loading.value = true

  try {
    // 构建更新数据
    console.log('准备更新用户资料:', formData)
    const updateData = {
      nickname: formData.nickname,
      profile_public: formData.profile_public,
      email_notifications: formData.email_notifications,
    }
    console.log('准备更新用户资料:', updateData)

    // 如果提供了密码
    if (formData.current_password && formData.new_password) {
      updateData.current_password = formData.current_password
      updateData.new_password = formData.new_password
    }
    console.log('更新数据:', updateData)

    // 调用 API 更新用户信息
    api.updateProfile(updateData)

    // 更新本地数据
    if (userData.value) {
      Object.assign(userData.value, updateData)
    }

    // 同时更新authStore中的user信息
    authStore.setUser({
      ...authStore.user,
      ...updateData
    })

    successMessage.value = t('profile.success.profileUpdated')

    // 清空密码字段
    formData.current_password = ''
    formData.new_password = ''
    formData.confirm_password = ''

    // 3秒后清除成功消息
    setTimeout(() => {
      successMessage.value = ''
    }, 3000)

  } catch (error) {
    errorMessage.value = error.message || t('profile.errors.updateFailed')
  } finally {
    loading.value = false
  }
}

function handleCancel() {
  // 重置表单数据
  if (userData.value) {
    Object.assign(formData, userData.value)
  }
  // 清空密码字段
  formData.current_password = ''
  formData.new_password = ''
  formData.confirm_password = ''

  // 清空错误信息
  Object.keys(errors).forEach(key => {
    errors[key] = ''
  })

  errorMessage.value = ''
  successMessage.value = ''

  router.back()
}

function triggerAvatarUpload() {
  showAvatarModal.value = true
}

function closeAvatarModal() {
  showAvatarModal.value = false
}

function selectAvatar(url) {
  formData.avatar = url
}

function confirmAvatarChange() {
  // 这里可以调用 API 更新头像
  if (userData.value) {
    userData.value.avatar = formData.avatar
  }
  closeAvatarModal()
}

// 格式化时间戳函数
function formatDate(timestamp) {
  if (!timestamp) return '';

  try {
    const date = new Date(timestamp);

    // 格式化为本地时间：2026-01-17 14:18
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    }).replace(/\//g, '-');

  } catch (error) {
    console.error('时间格式错误:', error);
    return timestamp;
  }
}
</script>

<style scoped>
:root {
  --primary: #4361ee;
  --primary-dark: #3a56d4;
  --secondary: #7209b7;
  --danger: #e63946;
  --success: #4cc9f0;
  --warning: #f8961e;
  --light: #f8f9fa;
  --dark: #212529;
  --gray: #6c757d;
  --light-gray: #e9ecef;
  --border-radius: 12px;
  --shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
  --shadow-hover: 0 15px 40px rgba(67, 97, 238, 0.15);
  --transition: all 0.3s ease;
}

.profile-edit-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  animation: fadeIn 0.8s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 语言选择器 */
.language-switcher {
  position: absolute;
  top: 20px;
  right: 20px;
  z-index: 10;
}

/* 主容器 */
.profile-container {
  display: flex;
  width: 100%;
  max-width: 1200px;
  min-height: 700px;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: var(--shadow);
  background: white;
}

/* 左侧面板 */
.left-panel {
  flex: 1;
  background: linear-gradient(135deg, var(--primary) 0%, var(--secondary) 100%);
  color: white;
  padding: 40px;
  display: flex;
  flex-direction: column;
  position: relative;
  overflow: hidden;
}

.left-panel::before,
.left-panel::after {
  content: '';
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
}

.left-panel::before {
  top: -30%;
  right: -20%;
  width: 200px;
  height: 200px;
}

.left-panel::after {
  bottom: -20%;
  left: -10%;
  width: 150px;
  height: 150px;
}

/* 用户头像 */
.profile-header {
  text-align: center;
  margin-bottom: 40px;
  z-index: 1;
}

.avatar-container {
  display: inline-block;
  cursor: pointer;
  margin-bottom: 20px;
}

.avatar-wrapper {
  position: relative;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  overflow: hidden;
  border: 4px solid rgba(255, 255, 255, 0.3);
  transition: var(--transition);
}

.avatar-wrapper:hover {
  transform: scale(1.05);
  border-color: rgba(255, 255, 255, 0.5);
}

.avatar-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-size: 120px;
  color: rgba(255, 255, 255, 0.8);
}

.avatar-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 8px;
  transform: translateY(100%);
  transition: var(--transition);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.avatar-wrapper:hover .avatar-overlay {
  transform: translateY(0);
}

.avatar-overlay i {
  font-size: 16px;
}

.avatar-overlay span {
  font-size: 12px;
}

.profile-header h2 {
  font-size: 28px;
  margin-bottom: 5px;
}

.user-email {
  font-size: 14px;
  opacity: 0.9;
  margin-bottom: 10px;
}

.member-since {
  font-size: 13px;
  opacity: 0.8;
}

/* 统计信息 */
.profile-stats {
  background: rgba(255, 255, 255, 0.1);
  border-radius: var(--border-radius);
  padding: 20px;
  margin-bottom: 30px;
  z-index: 1;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
  padding-bottom: 15px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.stat-item:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.stat-item i {
  font-size: 20px;
  width: 40px;
  height: 40px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-content {
  flex: 1;
}

.stat-label {
  display: block;
  font-size: 12px;
  opacity: 0.8;
  margin-bottom: 4px;
}

.stat-value {
  display: block;
  font-size: 15px;
  font-weight: 500;
}

.status-active {
  color: #4ade80;
}

/* 隐私提示 */
.privacy-tips {
  background: rgba(255, 255, 255, 0.1);
  border-radius: var(--border-radius);
  padding: 20px;
  z-index: 1;
}

.privacy-tips h3 {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  margin-bottom: 15px;
}

.privacy-tips ul {
  list-style: none;
  padding-left: 5px;
}

.privacy-tips li {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
  font-size: 14px;
  opacity: 0.9;
}

.privacy-tips li i {
  color: #4ade80;
}

/* 右侧编辑表单 */
.edit-form {
  flex: 2;
  padding: 40px;
  display: flex;
  flex-direction: column;
  overflow: auto;
}

.form-header {
  margin-bottom: 30px;
  border-bottom: 2px solid var(--light-gray);
  padding-bottom: 20px;
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

/* 表单部分 */
.form-section {
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--light-gray);
}

.form-section:last-child {
  border-bottom: none;
}

.form-section h3 {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  color: var(--dark);
  margin-bottom: 20px;
}

.form-section h3 i {
  color: var(--primary);
}

/* 输入框组 */
.input-group {
  position: relative;
  margin-bottom: 20px;
}

.input-group i:first-child {
  position: absolute;
  left: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--gray);
  font-size: 18px;
  transition: var(--transition);
}

.input-group input,
.input-group .select-input {
  width: 100%;
  padding: 16px 16px 16px 50px;
  border: 2px solid var(--light-gray);
  border-radius: var(--border-radius);
  font-size: 16px;
  transition: var(--transition);
  background: white;
  color: var(--dark);
}

.select-input {
  appearance: none;
  cursor: pointer;
}

.input-group input:focus,
.input-group .select-input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(67, 97, 238, 0.1);
}

.input-group input:focus + i,
.input-group .select-input:focus + i {
  color: var(--primary);
}

.input-group input.has-error {
  border-color: var(--danger);
}

.toggle-password {
  position: absolute;
  right: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--gray);
  cursor: pointer;
  z-index: 2;
}

.toggle-password:hover {
  color: var(--primary);
}

.field-error {
  color: var(--danger);
  font-size: 12px;
  margin-top: 5px;
  margin-left: 50px;
  animation: slideIn 0.3s ease;
}

.field-note {
  display: block;
  font-size: 12px;
  color: var(--gray);
  margin-top: 5px;
  margin-left: 50px;
}

/* 行内输入框 */
.input-row {
  display: flex;
  gap: 15px;
}

.half-width {
  flex: 1;
}

/* 密码强度指示器 */
.password-strength {
  margin-top: 10px;
  margin-left: 50px;
}

.strength-bar {
  height: 4px;
  border-radius: 2px;
  margin-bottom: 5px;
  transition: var(--transition);
}

.strength-0 { width: 20%; background: var(--danger); }
.strength-1 { width: 40%; background: #ff6b6b; }
.strength-2 { width: 60%; background: var(--warning); }
.strength-3 { width: 80%; background: #51cf66; }
.strength-4 { width: 100%; background: var(--success); }

.strength-text {
  font-size: 12px;
  color: var(--gray);
}

/* 隐私选项 */
.privacy-options {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.privacy-option {
  display: flex;
  align-items: center;
  padding: 15px;
  border: 2px solid var(--light-gray);
  border-radius: var(--border-radius);
  cursor: pointer;
  transition: var(--transition);
}

.privacy-option:hover {
  border-color: var(--primary);
  background: rgba(67, 97, 238, 0.02);
}

.privacy-option input[type="checkbox"] {
  margin-right: 15px;
  width: 18px;
  height: 18px;
  accent-color: var(--primary);
}

.option-content {
  display: flex;
  align-items: center;
  gap: 15px;
  flex: 1;
}

.option-content i {
  font-size: 20px;
  color: var(--primary);
  width: 40px;
  height: 40px;
  background: rgba(67, 97, 238, 0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.option-title {
  display: block;
  font-weight: 600;
  color: var(--dark);
  margin-bottom: 4px;
}

.option-description {
  display: block;
  font-size: 13px;
  color: var(--gray);
}

/* 按钮组 */
.form-actions {
  display: flex;
  justify-content: space-between;
  gap: 15px;
  margin: 30px 0;
}

.btn {
  padding: 14px 24px;
  border: none;
  border-radius: var(--border-radius);
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  flex: 1;
}

.btn-primary {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(67, 97, 238, 0.2);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: var(--shadow-hover);
}

.btn-secondary {
  background: white;
  color: var(--primary);
  border: 2px solid var(--light-gray);
}

.btn-secondary:hover:not(:disabled) {
  background: var(--light-gray);
  border-color: var(--gray);
}

.btn-danger {
  background: linear-gradient(135deg, var(--danger) 0%, #d62839 100%);
  color: white;
}

.btn-danger:hover:not(:disabled) {
  opacity: 0.9;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 错误和成功消息 */
.error-message,
.success-message {
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  animation: slideIn 0.3s ease;
}

.error-message {
  background: #fee;
  color: #c33;
  border-left: 4px solid #c33;
}

.success-message {
  background: #efffee;
  color: #2a8;
  border-left: 4px solid #2a8;
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
  to { transform: rotate(360deg); }
}

/* 模态框 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease;
}

.modal-content {
  background: white;
  border-radius: var(--border-radius);
  overflow-y: auto;
  padding: 30px;
  max-width: 500px;
  width: 90%;
  max-height: 80vh;
}

.modal-content h3 {
  font-size: 22px;
  color: var(--dark);
  margin-bottom: 20px;
  text-align: center;
}

.avatar-options {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
  margin-bottom: 20px;
}

.avatar-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 15px;
  border: 2px solid var(--light-gray);
  border-radius: var(--border-radius);
  cursor: pointer;
  transition: var(--transition);
}

.avatar-option:hover {
  border-color: var(--primary);
  transform: translateY(-2px);
}

.avatar-option img {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
}

.avatar-option span {
  font-size: 14px;
  color: var(--dark);
}

.modal-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

/* 响应式设计 */
@media (max-width: 992px) {
  .profile-container {
    flex-direction: column;
    max-width: 600px;
  }

  .left-panel,
  .edit-form {
    padding: 30px;
  }

  .edit-form {
    max-height: none;
  }

  .form-actions {
    flex-direction: column;
  }
}

@media (max-width: 576px) {
  .profile-container {
    border-radius: 16px;
  }

  .left-panel,
  .edit-form {
    padding: 20px;
  }

  .input-row {
    flex-direction: column;
    gap: 0;
  }

  .avatar-options {
    grid-template-columns: 1fr;
  }

  .modal-content {
    width: 95%;
    padding: 20px;
  }
}
</style>
