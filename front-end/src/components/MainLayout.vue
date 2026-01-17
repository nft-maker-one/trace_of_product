<!-- eslint-disable vue/multi-word-component-names -->
<template>
  <div class="dashboard">
    <!-- 移动端菜单按钮 -->
    <button class="menu-toggle" @click="toggleMobileMenu">
      <i class="fas fa-bars"></i>
    </button>

    <!-- 侧边栏 -->
    <div class="sidebar" :class="{ 'active': mobileMenuOpen }">
      <div class="logo">
        <i class="fas fa-seedling"></i>
        <h1>{{ t('layout.systemName') }}</h1>
      </div>

      <ul class="nav-menu">
        <li v-for="item in menuItems" :key="item.id" class="nav-item">
          <router-link
            :to="item.path"
            class="nav-link"
            :class="{ 'active': isActive(item.id) }"
            @click="closeMobileMenu"
          >
            <i :class="item.icon"></i>
            <span>{{ item.name }}</span>
          </router-link>
        </li>
      </ul>

      <div class="sidebar-footer">
        <p class="system-version">v1.0.0</p>
        <p class="system-status">
          <i class="fas fa-circle" :class="{ 'online': systemOnline }"></i>
          {{ systemOnline ? t('layout.systemOnline') : t('layout.systemOffline') }}
        </p>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content" @click="closeMobileMenuOnClick">
      <div class="header">
        <h2>{{ pageTitle }}</h2>
        <div class="user-profile">
          <img :src="userAvatar" alt="User Avatar" class="user-avatar">
          <div class="user-info">
            <h4>{{ userName }}</h4>
            <p>{{ $t('layout.loginTime') }}: {{ formatLoginTime }}</p>
          </div>
          <LanguageSelector mode="simple" size="small" />
          <button class="logout-btn" @click="handleLogout">
            <i class="fas fa-sign-out-alt"></i>
            <span>{{ t('layout.logout') }}</span>
          </button>
        </div>
      </div>

      <slot name="content"></slot>

      <!-- 页脚 -->
      <div class="footer">
        <p>{{ t('layout.copyright') }}</p>
        <p class="footer-links">
          <a href="#" @click.prevent="showAbout">{{ t('layout.aboutUs') }}</a> |
          <a href="#" @click.prevent="showHelp">{{ t('layout.helpDoc') }}</a> |
          <a href="#" @click.prevent="showContact">{{ t('layout.contactUs') }}</a>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useI18nStore } from '../stores/i18n'
import LanguageSelector from './LanguageSelector.vue'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const authStore = useAuthStore()
const i18nStore = useI18nStore()
const { t } = useI18n()

const mobileMenuOpen = ref(false)
const systemOnline = ref(true)
const loginTime = ref(new Date())

const menuItems = computed(() => [
  { id: 'personal-center', name: t('layout.personalCenter'), path: '/dashboard#personal-center', icon: 'fas fa-user-circle' },
  { id: 'data-upload', name: t('layout.dataUpload'), path: '/dashboard#data-upload', icon: 'fas fa-upload' },
  { id: 'data-query', name: t('layout.dataQuery'), path: '/dashboard#data-query', icon: 'fas fa-search' },
  { id: 'node-management', name: t('layout.nodeManagement'), path: '/dashboard#node-management', icon: 'fas fa-server' },
  { id: 'log-monitor', name: t('layout.logMonitor'), path: '/dashboard#log-monitor', icon: 'fas fa-list-alt' }
])

const pageTitle = computed(() => {
  const tab = route.hash.replace('#', '') || 'personal-center'
  const titles = {
    'data-query': t('layout.dataQuery'),
    'node-management': t('layout.nodeManagement'),
    'log-monitor': t('layout.logMonitor')
  }
  return titles[tab] || t('layout.personalCenter')
})

const userName = computed(() => {
  return authStore.user?.username || localStorage.getItem('username') || t('common.noData')
})

const userAvatar = computed(() => {
  return '/profile.jpg'
})

const formatLoginTime = computed(() => {
  return loginTime.value.toLocaleString(i18nStore.currentLanguage === 'zh-CN' ? 'zh-CN' : 'en-US', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
})

function isActive(menuId) {
  const currentTab = route.hash.replace('#', '') || 'personal-center'
  return currentTab === menuId
}

function toggleMobileMenu() {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

function closeMobileMenu() {
  mobileMenuOpen.value = false
}

function closeMobileMenuOnClick() {
  if (window.innerWidth <= 768) {
    mobileMenuOpen.value = false
  }
}

function handleLogout() {
  if (confirm(t('common.confirm') + t('layout.logout') + '？')) {
    authStore.clearAuth()
  }
}

function showAbout() {
  alert(t('layout.aboutContent'))
}

function showHelp() {
  alert(t('layout.helpContent'))
}

function showContact() {
  alert(t('layout.contactContent'))
}

// 监听窗口大小变化
onMounted(() => {
  const handleResize = () => {
    if (window.innerWidth > 768) {
      mobileMenuOpen.value = false
    }
  }

  window.addEventListener('resize', handleResize)

  // 检查系统状态
  checkSystemStatus()

  // 定期检查系统状态
  const intervalId = setInterval(checkSystemStatus, 30000)

  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
    clearInterval(intervalId)
  })
})

function checkSystemStatus() {
  // 这里可以调用API检查系统状态
  systemOnline.value = true // 暂时硬编码为在线
}
</script>

<style scoped>
/* 布局样式 */
.dashboard {
  display: flex;
  min-height: 100vh;
  background-color: var(--light);
}

/* 侧边栏样式 */
.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  width: 260px;
  height: 100vh;
  background: linear-gradient(180deg, var(--primary-dark) 0%, var(--primary) 100%);
  color: white;
  padding: 25px 0;
  box-shadow: 3px 0 15px rgba(0, 0, 0, 0.1);
  z-index: 100;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.logo {
  display: flex;
  align-items: center;
  padding: 0 25px 30px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  margin-bottom: 25px;
}

.logo i {
  font-size: 28px;
  margin-right: 12px;
  color: #a5d6a7;
}

.logo h1 {
  font-size: 22px;
  font-weight: 700;
  letter-spacing: 0.5px;
}

.logo span {
  color: #c8e6c9;
  font-weight: 300;
}

.nav-menu {
  list-style: none;
  padding: 0 15px;
  flex: 1;
}

.nav-item {
  margin-bottom: 8px;
}

.nav-link {
  display: flex;
  align-items: center;
  padding: 14px 20px;
  color: rgba(255, 255, 255, 0.85);
  text-decoration: none;
  border-radius: 10px;
  transition: var(--transition);
  font-weight: 500;
}

.nav-link i {
  font-size: 18px;
  margin-right: 15px;
  width: 24px;
  text-align: center;
}

.nav-link:hover {
  background-color: rgba(255, 255, 255, 0.1);
  color: white;
}

.nav-link.active {
  background-color: rgba(255, 255, 255, 0.15);
  color: white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.sidebar-footer {
  padding: 20px 25px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
}

.system-version {
  margin-bottom: 5px;
}

.system-status .online {
  color: #4caf50;
}

/* 主要内容区域 */
.main-content {
  margin-left: 260px;
  padding: 20px;
  min-height: 100vh;
  width: calc(100% - 260px);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 0;
  margin-bottom: 30px;
  border-bottom: 1px solid var(--light-gray);
}

.header h2 {
  font-size: 28px;
  color: var(--dark);
  font-weight: 600;
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 15px;
}

.user-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid var(--primary-light);
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.1);
}

.user-info h4 {
  font-size: 16px;
  margin-bottom: 5px;
}

.user-info p {
  font-size: 13px;
  color: var(--gray);
}

.logout-btn {
  background: var(--light-gray);
  border: none;
  border-radius: 8px;
  padding: 8px 16px;
  color: var(--dark);
  font-weight: 500;
  cursor: pointer;
  transition: var(--transition);
  display: flex;
  align-items: center;
  gap: 8px;
}

.logout-btn:hover {
  background: #e0e0e0;
}

/* 页脚 */
.footer {
  margin-top: 50px;
  padding: 20px 0;
  border-top: 1px solid var(--light-gray);
  text-align: center;
  color: var(--gray);
  font-size: 14px;
}

.footer-links {
  margin-top: 10px;
}

.footer-links a {
  color: var(--primary);
  text-decoration: none;
  margin: 0 10px;
  transition: var(--transition);
}

.footer-links a:hover {
  text-decoration: underline;
  color: var(--primary-dark);
}

/* 移动端菜单按钮 */
.menu-toggle {
  display: none;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 10px 15px;
  font-size: 20px;
  cursor: pointer;
  position: fixed;
  top: 20px;
  left: 20px;
  z-index: 101;
}

/* 响应式设计 */
@media (max-width: 992px) {
  .sidebar {
    width: 80px;
    overflow: hidden;
  }

  .sidebar .logo h1,
  .sidebar .nav-link span,
  .sidebar-footer {
    display: none;
  }

  .sidebar .logo {
    justify-content: center;
    padding: 0 15px 30px;
  }

  .sidebar .logo i {
    margin-right: 0;
    font-size: 32px;
  }

  .main-content {
    margin-left: 80px;
    width: calc(100% - 80px);
  }

  .nav-link {
    justify-content: center;
    padding: 15px;
  }

  .nav-link i {
    margin-right: 0;
    font-size: 20px;
  }
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
    padding: 15px;
    width: 100%;
  }

  .sidebar {
    transform: translateX(-100%);
    transition: transform 0.3s ease;
    width: 280px;
  }

  .sidebar.active {
    transform: translateX(0);
  }

  .sidebar .logo h1,
  .sidebar .nav-link span,
  .sidebar-footer {
    display: block;
  }

  .sidebar .logo {
    justify-content: flex-start;
  }

  .sidebar .logo i {
    margin-right: 12px;
  }

  .nav-link {
    justify-content: flex-start;
  }

  .nav-link i {
    margin-right: 15px;
  }

  .header {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }

  .user-profile {
    width: 100%;
    justify-content: space-between;
  }

  .menu-toggle {
    display: block;
  }
}
</style>
