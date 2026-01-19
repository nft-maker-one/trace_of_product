<template>
  <div class="personal-center">
    <div class="profile-container">
      <div class="card profile-card">
        <div class="card-header">
          <h3 class="card-title">
            <i class="fas fa-user-circle"></i>
            {{ t('personalCenter.userInfo') }}
          </h3>
          <button class="btn btn-outline" @click="showEditModal = true">
            <i class="fas fa-edit"></i>
            {{ t('personalCenter.editProfile') }}
          </button>
        </div>
        <div class="profile-content">
          <div class="avatar-section">
            <img :src="userAvatar" alt="User Avatar" class="profile-avatar">
            <h3>{{ userName }}</h3>
            <p class="user-role">农产品溯源系统用户</p>
            <!-- <div class="avatar-actions">
              <button class="btn btn-sm" @click="changeAvatar">
                <i class="fas fa-camera"></i>
                {{ t('personalCenter.changeAvatar') }}
              </button>
            </div> -->
          </div>

          <div class="user-details">
            <div class="detail-item">
              <div class="detail-label">
                <i class="fas fa-id-card"></i>
                {{ t('personalCenter.userId') }}:
              </div>
              <div class="detail-value">{{ userId }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-label">
                <i class="fas fa-sign-in-alt"></i>
                {{ t('personalCenter.loginTime') }}:
              </div>
              <div class="detail-value">{{ formatLoginTime }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-label">
                <i class="fas fa-user-tag"></i>
                {{ t('personalCenter.accountType') }}:
              </div>
              <div class="detail-value">
                <span class="user-type-badge">{{ userType }}</span>
              </div>
            </div>
            <div class="detail-item">
              <div class="detail-label">
                <i class="fas fa-history"></i>
                {{ t('personalCenter.lastActivity') }}:
              </div>
              <div class="detail-value">{{ lastActivity }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-label">
                <i class="fas fa-envelope"></i>
                {{ t('personalCenter.email') }}:
              </div>
              <div class="detail-value">{{ userEmail }}</div>
            </div>
          </div>
        </div>
      </div>

      <div class="card profile-card">
        <div class="card-header">
          <h3 class="card-title">
            <i class="fas fa-chart-line"></i>
            {{ t('personalCenter.systemStatus') }}
          </h3>
          <button class="btn btn-outline" @click="refreshSystemStatus">
            <i class="fas fa-sync-alt"></i>
            {{ t('common.refresh') }}
          </button>
        </div>
        <div class="status-content">
          <div class="status-grid">
            <div class="status-item">
              <div class="status-icon">
                <i class="fas fa-server"></i>
              </div>
              <div class="status-info">
                <div class="status-label">{{ t('personalCenter.onlineNodes') }}</div>
                <div class="status-value">{{ systemStatus.onlineNodes }}</div>
              </div>
            </div>
            <div class="status-item">
              <div class="status-icon">
                <i class="fas fa-database"></i>
              </div>
              <div class="status-info">
                <div class="status-label">{{ t('personalCenter.totalData') }}</div>
                <div class="status-value">{{ systemStatus.totalData }}</div>
              </div>
            </div>
            <div class="status-item">
              <div class="status-icon">
                <i class="fas fa-exchange-alt"></i>
              </div>
              <div class="status-info">
                <div class="status-label">{{ t('personalCenter.todayUploads') }}</div>
                <div class="status-value">{{ systemStatus.todayUploads }}</div>
              </div>
            </div>
            <div class="status-item">
              <div class="status-icon">
                <i class="fas fa-search"></i>
              </div>
              <div class="status-info">
                <div class="status-label">{{ t('personalCenter.todayQueries') }}</div>
                <div class="status-value">{{ systemStatus.todayQueries }}</div>
              </div>
            </div>
          </div>

          <div class="system-details">
            <div class="detail-item">
              <div class="detail-label">{{ t('personalCenter.systemStatusLabel') }}:</div>
              <div class="detail-value">
                <span class="status-badge" :class="systemStatus.statusClass">
                  {{ systemStatus.statusText }}
                </span>
              </div>
            </div>
            <div class="detail-item">
              <div class="detail-label">{{ t('personalCenter.apiResponse') }}:</div>
              <div class="detail-value">{{ systemStatus.apiResponse }} ms</div>
            </div>
            <div class="detail-item">
              <div class="detail-label">{{ t('personalCenter.systemVersion') }}:</div>
              <div class="detail-value">v{{ systemStatus.version }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-label">{{ t('personalCenter.lastUpdate') }}:</div>
              <div class="detail-value">{{ systemStatus.lastUpdate }}</div>
            </div>
          </div>

          <div class="progress-section">
            <div class="progress-item">
              <div class="progress-header">
                <span>{{ t('personalCenter.storageUsage') }}</span>
                <span>{{ systemStatus.storageUsage }}%</span>
              </div>
              <div class="progress-bar">
                <div class="progress-fill" :style="{ width: systemStatus.storageUsage + '%' }"></div>
              </div>
            </div>
            <div class="progress-item">
              <div class="progress-header">
                <span>{{ t('personalCenter.cpuUsage') }}</span>
                <span>{{ systemStatus.cpuUsage }}%</span>
              </div>
              <div class="progress-bar">
                <div class="progress-fill" :style="{ width: systemStatus.cpuUsage + '%' }"></div>
              </div>
            </div>
            <div class="progress-item">
              <div class="progress-header">
                <span>{{ t('personalCenter.memoryUsage') }}</span>
                <span>{{ systemStatus.memoryUsage }}%</span>
              </div>
              <div class="progress-bar">
                <div class="progress-fill" :style="{ width: systemStatus.memoryUsage + '%' }"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 活动日志 -->
    <div class="card activity-card">
      <div class="card-header">
        <h3 class="card-title">
          <i class="fas fa-history"></i>
          {{ t('personalCenter.activityLog') }}
        </h3>
        <button class="btn btn-outline" @click="clearLogs">
          <i class="fas fa-trash-alt"></i>
          {{ t('personalCenter.clearLog') }}
        </button>
      </div>
      <div class="activity-list">
        <div v-if="activityLogs.length === 0" class="empty-activity">
          <i class="fas fa-clipboard-list"></i>
          <p>{{ t('personalCenter.noActivity') }}</p>
        </div>
        <div v-else class="activity-items">
          <div v-for="log in activityLogs" :key="log.id" class="activity-item">
            <div class="activity-icon" :class="log.type">
              <i :class="log.icon"></i>
            </div>
            <div class="activity-content">
              <div class="activity-title">{{ log.title }}</div>
              <div class="activity-desc">{{ log.description }}</div>
              <div class="activity-time">{{ log.time }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- 模态框 -->
      <EditProfileModal
        :visible="showEditModal"
        @update:visible="showEditModal = $event"
      />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useI18n } from 'vue-i18n'
import EditProfileModal from './EditProfileModal.vue'

const authStore = useAuthStore()
const { t } = useI18n()

// 用户信息
const userName = computed(() => {
  return authStore.user?.nickname || authStore.user?.username || t('personalCenter.notLoggedIn')
})
const userId = computed(() => {
  return authStore.user?.id || 'USER' + Math.random().toString().slice(2, 8)
})
const userEmail = computed(() => {
  return authStore.user?.email
})
const userAvatar = computed(() => authStore.user?.avatarUrl)

const userType = ref(t('personalCenter.standardUser'))
const lastActivity = computed(() => {
  return new Date().toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
})

const formatLoginTime = computed(() => {
  if (authStore.user?.loginTime) {
    return new Date(authStore.user.loginTime).toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }
  return t('personalCenter.notLoggedIn')
})

// 系统状态
const systemStatus = reactive({
  onlineNodes: 5,
  totalData: 1250,
  todayUploads: 23,
  todayQueries: 45,
  apiResponse: 156,
  version: '1.0.0',
  lastUpdate: '2023-10-15 14:30:00',
  storageUsage: 65,
  cpuUsage: 42,
  memoryUsage: 78,
  statusClass: 'online',
  statusText: t('personalCenter.runningNormal')
})

// 活动日志
const activityLogs = ref([
  {
    id: 1,
    type: 'login',
    icon: 'fas fa-sign-in-alt',
    title: t('personalCenter.userLogin'),
    description: t('personalCenter.loginSuccess'),
    time: t('common.justNow')
  },
  {
    id: 2,
    type: 'upload',
    icon: 'fas fa-upload',
    title: t('personalCenter.dataUpload'),
    description: t('personalCenter.uploadProductData') + ' (ID: EG20231015001)',
    time: '10' + t('common.minutesAgo')
  },
  {
    id: 3,
    type: 'query',
    icon: 'fas fa-search',
    title: t('personalCenter.dataQuery'),
    description: t('personalCenter.queryProductInfo') + ' (ID: EG20231014045)',
    time: '30' + t('common.minutesAgo')
  },
  {
    id: 4,
    type: 'node',
    icon: 'fas fa-server',
    title: t('personalCenter.nodeView'),
    description: t('personalCenter.viewOnlineNodes'),
    time: '1' + t('common.hoursAgo')
  }
])

const showEditModal = ref(false)

// function changeAvatar() {
//   alert(t('personalCenter.changeAvatarDeveloping'))
// }

function refreshSystemStatus() {
  // 模拟刷新系统状态
  systemStatus.apiResponse = Math.floor(Math.random() * 100) + 50
  systemStatus.todayUploads += Math.floor(Math.random() * 5)
  systemStatus.todayQueries += Math.floor(Math.random() * 10)
  systemStatus.storageUsage = Math.min(100, systemStatus.storageUsage + Math.floor(Math.random() * 3))
  systemStatus.cpuUsage = Math.floor(Math.random() * 50) + 20
  systemStatus.memoryUsage = Math.floor(Math.random() * 30) + 50

  // 添加活动记录
  activityLogs.value.unshift({
    id: Date.now(),
    type: 'refresh',
    icon: 'fas fa-sync-alt',
    title: t('personalCenter.systemRefresh'),
    description: t('personalCenter.refreshStatus'),
    time: t('common.justNow')
  })
}

function clearLogs() {
  if (confirm(t('personalCenter.confirmClearLog'))) {
    activityLogs.value = []
  }
}

onMounted(() => {
  // 初始化时获取系统状态
  refreshSystemStatus()
})
</script>

<style scoped>
.personal-center {
  padding: 0;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4edf5 100%);
  min-height: 100vh;
}

.profile-container {
  display: flex;
  gap: 30px;
  flex-wrap: wrap;
  margin-bottom: 30px;
}

.profile-card {
  flex: 1;
  min-width: 300px;
  background: white;
  border-radius: 20px;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(46, 125, 50, 0.1);
  position: relative;
}

.profile-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--primary) 0%, var(--primary-light) 100%);
}

.profile-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 20px 40px rgba(46, 125, 50, 0.15);
}

.card-header {
  padding: 24px 30px 20px;
  border-bottom: 1px solid rgba(46, 125, 50, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(135deg, rgba(46, 125, 50, 0.02) 0%, rgba(76, 175, 80, 0.02) 100%);
}

.card-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--dark);
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 0;
}

.card-title i {
  color: var(--primary);
  font-size: 24px;
}

.btn {
  padding: 10px 24px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.3s ease;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  border: 2px solid transparent;
}

.btn-outline {
  background: transparent;
  color: var(--primary);
  border-color: var(--primary);
}

.btn-outline:hover {
  background: var(--primary);
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(46, 125, 50, 0.3);
}

.profile-content {
  padding: 30px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.avatar-section {
  text-align: center;
  margin-bottom: 40px;
  position: relative;
}

.profile-avatar {
  width: 140px;
  height: 140px;
  border-radius: 50%;
  object-fit: cover;
  margin-bottom: 20px;
  border: 8px solid white;
  box-shadow:
    0 10px 30px rgba(0, 0, 0, 0.1),
    inset 0 0 0 4px var(--primary-light);
  transition: all 0.3s ease;
  position: relative;
  z-index: 1;
}

.profile-avatar:hover {
  transform: scale(1.05);
  box-shadow:
    0 15px 40px rgba(46, 125, 50, 0.3),
    inset 0 0 0 4px var(--primary);
}

.avatar-section h3 {
  font-size: 28px;
  margin-bottom: 8px;
  color: var(--dark);
  font-weight: 700;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.user-role {
  color: var(--gray);
  margin-bottom: 20px;
  font-size: 15px;
  position: relative;
  display: inline-block;
  padding: 0 15px;
}

.user-role::before,
.user-role::after {
  content: '';
  position: absolute;
  top: 50%;
  width: 30px;
  height: 1px;
  background: linear-gradient(90deg, transparent, var(--primary-light));
}

.user-role::before {
  left: -30px;
}

.user-role::after {
  right: -30px;
  transform: rotate(180deg);
}

.avatar-actions {
  margin-top: 15px;
}

.btn-sm {
  padding: 10px 20px;
  font-size: 14px;
  border-radius: 10px;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
  color: white;
  border: none;
}

.btn-sm:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(46, 125, 50, 0.4);
}

.user-details {
  width: 100%;
  margin-top: 20px;
  background: #f8f9fa;
  border-radius: 16px;
  padding: 25px;
}

.detail-item {
  display: flex;
  margin-bottom: 18px;
  padding-bottom: 18px;
  border-bottom: 1px solid rgba(46, 125, 50, 0.1);
  align-items: center;
  transition: all 0.3s ease;
}

.detail-item:hover {
  transform: translateX(10px);
  border-bottom-color: var(--primary-light);
}

.detail-label {
  font-weight: 600;
  width: 140px;
  color: var(--gray);
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 15px;
}

.detail-label i {
  width: 20px;
  color: var(--primary);
  font-size: 18px;
}

.detail-value {
  color: var(--dark);
  flex: 1;
  font-size: 16px;
  font-weight: 500;
}

.user-type-badge {
  display: inline-block;
  padding: 6px 16px;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
  color: white;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.5px;
  box-shadow: 0 4px 12px rgba(46, 125, 50, 0.3);
  position: relative;
  overflow: hidden;
}

.user-type-badge::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  transition: left 0.5s ease;
}

.user-type-badge:hover::after {
  left: 100%;
}

/* 系统状态样式 */
.status-content {
  padding: 30px;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 25px;
  margin-bottom: 40px;
}

.status-item {
  display: flex;
  align-items: center;
  padding: 25px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 16px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border: 2px solid transparent;
  position: relative;
  overflow: hidden;
}

.status-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--primary) 0%, var(--primary-light) 100%);
  transform: scaleX(0);
  transform-origin: left;
  transition: transform 0.4s ease;
}

.status-item:hover {
  transform: translateY(-8px);
  border-color: var(--primary-light);
  box-shadow: 0 15px 35px rgba(46, 125, 50, 0.15);
}

.status-item:hover::before {
  transform: scaleX(1);
}

.status-icon {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
  flex-shrink: 0;
  transition: all 0.3s ease;
}

.status-item:hover .status-icon {
  transform: rotate(15deg) scale(1.1);
}

.status-icon i {
  font-size: 28px;
  color: white;
}

.status-info {
  flex: 1;
}

.status-label {
  font-size: 14px;
  color: var(--gray);
  margin-bottom: 8px;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  font-weight: 600;
}

.status-value {
  font-size: 32px;
  font-weight: 800;
  color: var(--primary-dark);
  line-height: 1;
}

.system-details {
  margin-bottom: 40px;
  background: #f8f9fa;
  border-radius: 16px;
  padding: 25px;
}

.system-details .detail-item {
  border-bottom: 1px dashed rgba(46, 125, 50, 0.2);
  margin-bottom: 15px;
  padding-bottom: 15px;
  transition: none;
}

.system-details .detail-item:hover {
  transform: none;
  border-bottom-color: var(--primary);
}

.system-details .detail-label {
  width: 120px;
  color: var(--dark);
  font-weight: 600;
}

.status-badge {
  display: inline-block;
  padding: 8px 18px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  letter-spacing: 0.5px;
  transition: all 0.3s ease;
}

.status-badge.online {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.15) 0%, rgba(76, 175, 80, 0.05) 100%);
  color: var(--success);
  border: 2px solid rgba(76, 175, 80, 0.3);
  position: relative;
  padding-left: 35px;
}

.status-badge.online::before {
  content: '';
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  width: 8px;
  height: 8px;
  background: var(--success);
  border-radius: 50%;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { box-shadow: 0 0 0 0 rgba(76, 175, 80, 0.7); }
  70% { box-shadow: 0 0 0 10px rgba(76, 175, 80, 0); }
  100% { box-shadow: 0 0 0 0 rgba(76, 175, 80, 0); }
}

/* 进度条样式 */
.progress-section {
  margin-top: 40px;
}

.progress-item {
  margin-bottom: 30px;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  font-size: 15px;
  color: var(--dark);
  font-weight: 600;
}

.progress-header span:last-child {
  color: var(--primary);
  font-weight: 700;
}

.progress-bar {
  height: 12px;
  background: linear-gradient(135deg, #e9ecef 0%, #dee2e6 100%);
  border-radius: 10px;
  overflow: hidden;
  position: relative;
}

.progress-bar::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: repeating-linear-gradient(
    45deg,
    transparent,
    transparent 5px,
    rgba(255, 255, 255, 0.1) 5px,
    rgba(255, 255, 255, 0.1) 10px
  );
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg,
    var(--primary) 0%,
    var(--primary-light) 50%,
    var(--primary) 100%);
  border-radius: 10px;
  transition: width 0.8s cubic-bezier(0.34, 1.56, 0.64, 1);
  position: relative;
  z-index: 1;
  background-size: 200% 100%;
  animation: shimmer 2s infinite linear;
}

@keyframes shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}

/* 活动日志样式 */
.activity-card {
  margin-top: 30px;
  border-radius: 20px;
  overflow: hidden;
  background: white;
  border: 1px solid rgba(46, 125, 50, 0.1);
}

.activity-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--primary) 0%, var(--primary-light) 100%);
}

.activity-list {
  min-height: 200px;
  padding: 0;
}

.empty-activity {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: var(--gray);
}

.empty-activity i {
  font-size: 64px;
  margin-bottom: 20px;
  color: #e0e0e0;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.empty-activity p {
  font-size: 18px;
  color: #b0bec5;
}

.activity-items {
  max-height: 400px;
  overflow-y: auto;
  padding: 20px;
}

.activity-items::-webkit-scrollbar {
  width: 8px;
}

.activity-items::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.activity-items::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
  border-radius: 4px;
}

.activity-item {
  display: flex;
  padding: 20px;
  border-bottom: 1px solid rgba(46, 125, 50, 0.1);
  transition: all 0.3s ease;
  border-radius: 12px;
  margin-bottom: 10px;
  background: rgba(248, 249, 250, 0.5);
}

.activity-item:hover {
  background: rgba(46, 125, 50, 0.05);
  transform: translateX(10px);
  box-shadow: 0 5px 20px rgba(46, 125, 50, 0.1);
}

.activity-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.activity-icon {
  width: 50px;
  height: 50px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
  flex-shrink: 0;
  transition: all 0.3s ease;
  position: relative;
}

.activity-item:hover .activity-icon {
  transform: scale(1.1) rotate(5deg);
}

.activity-icon::after {
  content: '';
  position: absolute;
  top: -3px;
  right: -3px;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: inherit;
  opacity: 0.5;
  filter: blur(3px);
}

.activity-icon.login {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.2) 0%, rgba(76, 175, 80, 0.1) 100%);
  color: var(--success);
  border: 2px solid rgba(76, 175, 80, 0.3);
}

.activity-icon.upload {
  background: linear-gradient(135deg, rgba(33, 150, 243, 0.2) 0%, rgba(33, 150, 243, 0.1) 100%);
  color: var(--secondary);
  border: 2px solid rgba(33, 150, 243, 0.3);
}

.activity-icon.query {
  background: linear-gradient(135deg, rgba(255, 152, 0, 0.2) 0%, rgba(255, 152, 0, 0.1) 100%);
  color: var(--warning);
  border: 2px solid rgba(255, 152, 0, 0.3);
}

.activity-icon.node {
  background: linear-gradient(135deg, rgba(156, 39, 176, 0.2) 0%, rgba(156, 39, 176, 0.1) 100%);
  color: #9c27b0;
  border: 2px solid rgba(156, 39, 176, 0.3);
}

.activity-icon.refresh {
  background: linear-gradient(135deg, rgba(96, 125, 139, 0.2) 0%, rgba(96, 125, 139, 0.1) 100%);
  color: var(--gray);
  border: 2px solid rgba(96, 125, 139, 0.3);
}

.activity-content {
  flex: 1;
}

.activity-title {
  font-weight: 700;
  margin-bottom: 8px;
  color: var(--dark);
  font-size: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.activity-desc {
  font-size: 14px;
  color: var(--gray);
  margin-bottom: 8px;
  line-height: 1.5;
}

.activity-time {
  font-size: 13px;
  color: var(--primary);
  font-weight: 600;
  background: rgba(46, 125, 50, 0.1);
  padding: 4px 12px;
  border-radius: 20px;
  display: inline-block;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .profile-container {
    flex-direction: column;
  }

  .profile-card {
    min-width: 100%;
  }
}

@media (max-width: 768px) {
  .status-grid {
    grid-template-columns: 1fr;
    gap: 15px;
  }

  .card-header {
    flex-direction: column;
    gap: 15px;
    text-align: center;
  }

  .profile-content {
    padding: 20px;
  }

  .avatar-section h3 {
    font-size: 24px;
  }

  .profile-avatar {
    width: 120px;
    height: 120px;
  }

  .user-details {
    padding: 20px;
  }

  .detail-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .detail-label {
    width: 100%;
  }

  .status-item {
    padding: 20px;
  }

  .status-value {
    font-size: 28px;
  }
}

@media (max-width: 480px) {
  .card-title {
    font-size: 18px;
  }

  .btn {
    width: 100%;
    justify-content: center;
  }

  .status-icon {
    width: 50px;
    height: 50px;
    margin-right: 15px;
  }

  .activity-item {
    padding: 15px;
  }

  .activity-icon {
    width: 40px;
    height: 40px;
    margin-right: 15px;
  }
}

/* 加载动画 */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.profile-card, .activity-card {
  animation: fadeInUp 0.6s ease-out;
}

.profile-card:nth-child(2) {
  animation-delay: 0.2s;
}

.activity-card {
  animation-delay: 0.4s;
}
</style>
