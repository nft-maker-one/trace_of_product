<template>
  <div class="explorer-page">
    <!-- 动态背景 -->
    <div class="animated-bg">
      <div class="particle" v-for="i in 20" :key="i"></div>
    </div>

    <!-- 顶部导航栏 -->
    <nav class="navbar">
      <div class="nav-brand" @click="goHome">
        <i class="fas fa-seedling"></i>
        <span>{{ t('explorer.systemName') }}</span>
      </div>
      <div class="nav-links">
        <LanguageSelector mode="compact" />
        <router-link to="/login" class="nav-link login-btn">
          <i class="fas fa-sign-in-alt"></i>
          {{ t('explorer.login') }}
        </router-link>
      </div>
    </nav>

    <!-- 主标题区域 -->
    <header class="hero-section">
      <div class="hero-content">
        <div class="title-icon-wrapper">
          <i class="fas fa-cubes"></i>
        </div>
        <h1 class="hero-title">{{ t('explorer.title') }}</h1>
        <p class="hero-subtitle">{{ t('explorer.subtitle') }}</p>
      </div>
    </header>

    <!-- 搜索区域 - 强调 -->
    <section class="search-section">
      <div class="container">
        <div class="search-box">
          <div class="search-header">
            <i class="fas fa-search"></i>
            <h2>{{ t('explorer.queryBlocks') }}</h2>
          </div>
          <div class="search-form">
            <div class="search-inputs">
              <div class="input-wrapper">
                <label>{{ t('explorer.startHeight') }}</label>
                <div class="input-with-icon">
                  <i class="fas fa-arrow-up"></i>
                  <input 
                    v-model.number="startHeight" 
                    type="number" 
                    min="0" 
                    :placeholder="t('explorer.enterHeight')"
                  >
                </div>
              </div>
              <div class="input-divider">
                <i class="fas fa-arrows-alt-h"></i>
              </div>
              <div class="input-wrapper">
                <label>{{ t('explorer.endHeight') }}</label>
                <div class="input-with-icon">
                  <i class="fas fa-arrow-down"></i>
                  <input 
                    v-model.number="endHeight" 
                    type="number" 
                    min="0" 
                    :placeholder="t('explorer.enterHeight')"
                  >
                </div>
              </div>
            </div>
            <div class="search-actions">
              <button class="btn btn-primary btn-large" @click="queryBlocks" :disabled="loading">
                <i class="fas fa-search"></i>
                {{ loading ? t('explorer.loading') : t('explorer.search') }}
              </button>
              <button class="btn btn-secondary" @click="loadLatest" :disabled="loading">
                <i class="fas fa-sync-alt"></i>
                {{ t('explorer.loadLatest') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 统计信息 -->
    <section class="stats-section">
      <div class="container">
        <div class="stats-bar">
          <div class="stat-item" v-for="(stat, index) in statsData" :key="index">
            <i :class="stat.icon"></i>
            <span class="stat-value">{{ stat.value }}</span>
            <span class="stat-label">{{ stat.label }}</span>
          </div>
        </div>
      </div>
    </section>

    <!-- 区块列表 -->
    <section class="blocks-section">
      <div class="container">
        <div class="section-header">
          <h2>
            <i class="fas fa-cubes"></i>
            {{ t('explorer.blockList') }}
            <span class="count">({{ blocks.length }} {{ t('explorer.blocksUnit') }})</span>
          </h2>
        </div>

        <div v-if="loading" class="loading-state">
          <div class="loader">
            <div class="loader-ring"></div>
          </div>
          <p>{{ t('explorer.loadingBlocks') }}</p>
        </div>

        <div v-else-if="blocks.length === 0" class="empty-state">
          <i class="fas fa-inbox"></i>
          <p>{{ t('explorer.noBlocks') }}</p>
        </div>

        <div v-else class="blocks-list">
          <div 
            v-for="(block, index) in blocks" 
            :key="block.height" 
            class="block-card"
            :class="{ 'genesis': block.height === 0, 'expanded': expandedBlock === block.height }"
              :style="{ animationDelay: (index * 0.03) + 's' }"
            @click="toggleBlock(block.height)"
          >
            <div class="block-main">
              <div class="block-height">
                <span class="height-label">#</span>
                <span class="height-value">{{ block.height }}</span>
                <span v-if="block.height === 0" class="genesis-tag">Genesis</span>
              </div>
              <div class="block-info">
                <div class="info-item">
                  <i class="fas fa-clock"></i>
                  <span>{{ formatTime(block.timestamp) }}</span>
                </div>
                <div class="info-item">
                  <i class="fas fa-crown"></i>
                  <span>{{ block.leader }}</span>
                </div>
                <div class="info-item hash">
                  <i class="fas fa-fingerprint"></i>
                  <span>{{ formatHash(block.dataHash) }}</span>
                </div>
              </div>
              <div class="expand-icon">
                <i class="fas" :class="expandedBlock === block.height ? 'fa-chevron-up' : 'fa-chevron-down'"></i>
              </div>
            </div>

            <div class="block-details" v-show="expandedBlock === block.height">
              <div class="detail-grid">
                <div class="detail-item">
                  <span class="label">{{ t('explorer.dataHash') }}</span>
                  <code class="value">{{ block.dataHash }}</code>
                </div>
                <div class="detail-item">
                  <span class="label">{{ t('explorer.prevHash') }}</span>
                  <code class="value">{{ block.prevHash }}</code>
                </div>
                <div class="detail-item">
                  <span class="label">Nonce</span>
                  <span class="value">{{ block.nonce }}</span>
                </div>
                <div class="detail-item">
                  <span class="label">{{ t('explorer.scores') }}</span>
                  <div class="scores">
                    <span v-if="Object.keys(block.scores || {}).length === 0" class="no-data">-</span>
                    <span v-else v-for="(score, nodeId) in block.scores" :key="nodeId" class="score-tag">
                      {{ nodeId }}: {{ score }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 页脚 -->
    <footer class="footer">
      <div class="container">
        <div class="footer-brand">
          <i class="fas fa-seedling"></i>
          <span>{{ t('explorer.systemName') }}</span>
        </div>
        <p class="footer-text">{{ t('explorer.footerText') }}</p>
        <p class="copyright">© 2026 All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import api from '../services/api'
import LanguageSelector from '../components/LanguageSelector.vue'

const router = useRouter()
const { t } = useI18n()

const loading = ref(false)
const blocks = ref([])
const startHeight = ref(0)
const endHeight = ref(20)
const expandedBlock = ref(null)

// 统计数据
const statsData = computed(() => [
  { icon: 'fas fa-layer-group', value: blocks.value.length, label: t('explorer.totalBlocks') },
  { icon: 'fas fa-arrow-up', value: latestHeight.value, label: t('explorer.latestHeight') },
  { icon: 'fas fa-server', value: '4', label: t('explorer.activeNodes') }
])

const latestHeight = computed(() => {
  if (blocks.value.length === 0) return 0
  return Math.max(...blocks.value.map(b => b.height))
})

function goHome() {
  router.push('/explorer')
}

async function queryBlocks() {
  if (startHeight.value < 0 || endHeight.value < 0) {
    alert(t('explorer.invalidHeight'))
    return
  }
  if (startHeight.value > endHeight.value) {
    alert(t('explorer.startGreaterThanEnd'))
    return
  }

  loading.value = true
  try {
    const response = await api.getBlocksByRange(startHeight.value, endHeight.value)
    if (response && response.success) {
      blocks.value = response.blocks || []
      blocks.value.sort((a, b) => b.height - a.height)
    }
  } catch (error) {
    console.error('查询区块失败:', error)
  } finally {
    loading.value = false
  }
}

async function loadLatest() {
  loading.value = true
  try {
    const response = await api.getBlocksByRange(0, 100)
    if (response && response.success) {
      blocks.value = response.blocks || []
      blocks.value.sort((a, b) => b.height - a.height)
      if (blocks.value.length > 0) {
        endHeight.value = Math.max(...blocks.value.map(b => b.height))
        startHeight.value = 0
      }
    }
  } catch (error) {
    console.error('加载区块失败:', error)
  } finally {
    loading.value = false
  }
}

function toggleBlock(height) {
  expandedBlock.value = expandedBlock.value === height ? null : height
}

function formatTime(timestamp) {
  if (!timestamp || timestamp === 0) return 'Genesis'
  return new Date(timestamp * 1000).toLocaleString('zh-CN')
}

function formatHash(hash) {
  if (!hash) return '-'
  if (hash.length > 20) {
    return hash.substring(0, 10) + '...' + hash.substring(hash.length - 10)
  }
  return hash
}

onMounted(() => {
  loadLatest()
})
</script>

<style scoped>
/* 页面容器 */
.explorer-page {
  min-height: 100vh;
  background: linear-gradient(180deg, #0a0a1a 0%, #1a1a2e 50%, #0f1f3d 100%);
  color: white;
  position: relative;
  overflow-x: hidden;
}

/* 动态背景粒子 */
.animated-bg {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 0;
  pointer-events: none;
}

.particle {
  position: absolute;
  width: 4px;
  height: 4px;
  background: #4caf50;
  border-radius: 50%;
  opacity: 0.3;
  animation: floatParticle 15s infinite;
}

.particle:nth-child(odd) { background: #00e676; }
.particle:nth-child(1) { left: 5%; animation-delay: 0s; }
.particle:nth-child(2) { left: 10%; animation-delay: 1s; }
.particle:nth-child(3) { left: 15%; animation-delay: 2s; }
.particle:nth-child(4) { left: 20%; animation-delay: 0.5s; }
.particle:nth-child(5) { left: 25%; animation-delay: 1.5s; }
.particle:nth-child(6) { left: 35%; animation-delay: 3s; }
.particle:nth-child(7) { left: 45%; animation-delay: 0.8s; }
.particle:nth-child(8) { left: 55%; animation-delay: 2.5s; }
.particle:nth-child(9) { left: 65%; animation-delay: 1.2s; }
.particle:nth-child(10) { left: 75%; animation-delay: 3.5s; }
.particle:nth-child(11) { left: 85%; animation-delay: 0.3s; }
.particle:nth-child(12) { left: 90%; animation-delay: 2.8s; }
.particle:nth-child(13) { left: 30%; animation-delay: 4s; }
.particle:nth-child(14) { left: 40%; animation-delay: 1.8s; }
.particle:nth-child(15) { left: 50%; animation-delay: 0.6s; }
.particle:nth-child(16) { left: 60%; animation-delay: 3.2s; }
.particle:nth-child(17) { left: 70%; animation-delay: 2.2s; }
.particle:nth-child(18) { left: 80%; animation-delay: 4.5s; }
.particle:nth-child(19) { left: 95%; animation-delay: 1.1s; }
.particle:nth-child(20) { left: 3%; animation-delay: 2.7s; }

@keyframes floatParticle {
  0% { transform: translateY(100vh) scale(0); opacity: 0; }
  10% { opacity: 0.5; }
  90% { opacity: 0.3; }
  100% { transform: translateY(-100vh) scale(1); opacity: 0; }
}

/* 导航栏 */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 40px;
  background: rgba(10, 10, 26, 0.9);
  backdrop-filter: blur(20px);
  z-index: 1000;
  border-bottom: 1px solid rgba(46, 125, 50, 0.2);
}

.nav-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 20px;
  font-weight: 700;
  color: #4caf50;
  cursor: pointer;
  transition: all 0.3s;
}

.nav-brand:hover {
  color: #00e676;
  transform: scale(1.02);
}

.nav-brand i {
  font-size: 26px;
}

.nav-links {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 12px;
}

.login-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: #2e7d32;
  color: white;
  text-decoration: none;
  border-radius: 25px;
  font-weight: 600;
  transition: all 0.3s;
}

.login-btn:hover {
  background: #4caf50;
  transform: translateY(-2px);
  box-shadow: 0 5px 20px rgba(46, 125, 50, 0.4);
}

/* Hero区域 */
.hero-section {
  position: relative;
  z-index: 1;
  padding: 120px 40px 40px;
  text-align: center;
}

.hero-content {
  max-width: 800px;
  margin: 0 auto;
}

.title-icon-wrapper {
  width: 100px;
  height: 100px;
  margin: 0 auto 25px;
  background: linear-gradient(135deg, #2e7d32 0%, #00e676 100%);
  border-radius: 25px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 45px;
  animation: pulse 2s infinite, float 3s ease-in-out infinite;
  box-shadow: 0 10px 40px rgba(0, 230, 118, 0.3);
}

@keyframes pulse {
  0%, 100% { box-shadow: 0 10px 40px rgba(0, 230, 118, 0.3); }
  50% { box-shadow: 0 15px 60px rgba(0, 230, 118, 0.5); }
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.hero-title {
  font-size: 42px;
  font-weight: 800;
  margin-bottom: 15px;
  background: linear-gradient(90deg, #fff 0%, #00e676 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-subtitle {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.6);
  max-width: 600px;
  margin: 0 auto;
}

/* 容器 */
.container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 20px;
}

/* 搜索区域 - 强调 */
.search-section {
  position: relative;
  z-index: 1;
  padding: 30px 0;
}

.search-box {
  background: linear-gradient(135deg, rgba(46, 125, 50, 0.15) 0%, rgba(0, 230, 118, 0.1) 100%);
  border: 2px solid rgba(46, 125, 50, 0.3);
  border-radius: 20px;
  padding: 30px;
  backdrop-filter: blur(10px);
  transition: all 0.3s;
}

.search-box:hover {
  border-color: #4caf50;
  box-shadow: 0 10px 40px rgba(46, 125, 50, 0.2);
}

.search-header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 25px;
}

.search-header i {
  font-size: 28px;
  color: #00e676;
}

.search-header h2 {
  font-size: 24px;
  font-weight: 700;
  margin: 0;
}

.search-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.search-inputs {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
  flex-wrap: wrap;
}

.input-wrapper {
  flex: 1;
  min-width: 180px;
  max-width: 250px;
}

.input-wrapper label {
  display: block;
  margin-bottom: 8px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.7);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.input-with-icon {
  position: relative;
}

.input-with-icon i {
  position: absolute;
  left: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: #4caf50;
}

.input-with-icon input {
  width: 100%;
  padding: 14px 14px 14px 45px;
  background: rgba(255, 255, 255, 0.05);
  border: 2px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  color: white;
  font-size: 16px;
  transition: all 0.3s;
}

.input-with-icon input:focus {
  outline: none;
  border-color: #00e676;
  background: rgba(255, 255, 255, 0.1);
  box-shadow: 0 0 20px rgba(0, 230, 118, 0.2);
}

.input-divider {
  color: rgba(255, 255, 255, 0.3);
  padding-top: 25px;
}

.search-actions {
  display: flex;
  justify-content: center;
  gap: 15px;
  flex-wrap: wrap;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 12px 25px;
  border: none;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-large {
  padding: 16px 40px;
  font-size: 17px;
}

.btn-primary {
  background: linear-gradient(135deg, #2e7d32 0%, #00e676 100%);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-3px);
  box-shadow: 0 10px 30px rgba(0, 230, 118, 0.4);
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-secondary:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.2);
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 统计条 */
.stats-section {
  position: relative;
  z-index: 1;
  padding: 20px 0;
}

.stats-bar {
  display: flex;
  justify-content: center;
  gap: 40px;
  flex-wrap: wrap;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 20px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 30px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.stat-item i {
  color: #00e676;
  font-size: 18px;
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
}

.stat-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
}

/* 区块列表 */
.blocks-section {
  position: relative;
  z-index: 1;
  padding: 40px 0;
}

.section-header {
  margin-bottom: 25px;
}

.section-header h2 {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 22px;
  font-weight: 700;
  margin: 0;
}

.section-header i {
  color: #4caf50;
}

.count {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.4);
  font-weight: normal;
}

/* 加载状态 */
.loading-state {
  text-align: center;
  padding: 60px;
}

.loader {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.loader-ring {
  width: 50px;
  height: 50px;
  border: 4px solid rgba(255, 255, 255, 0.1);
  border-top-color: #00e676;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 60px;
  color: rgba(255, 255, 255, 0.4);
}

.empty-state i {
  font-size: 50px;
  margin-bottom: 15px;
}

/* 区块列表 */
.blocks-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.block-card {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  animation: slideIn 0.4s ease both;
  transition: all 0.3s;
}

.block-card:hover {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(46, 125, 50, 0.3);
  transform: translateX(5px);
}

.block-card.genesis {
  border-color: rgba(255, 152, 0, 0.4);
  background: linear-gradient(135deg, rgba(255, 152, 0, 0.1) 0%, rgba(255, 152, 0, 0.02) 100%);
}

.block-card.expanded {
  border-color: #4caf50;
}

@keyframes slideIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.block-main {
  display: flex;
  align-items: center;
  padding: 18px 20px;
  gap: 20px;
}

.block-height {
  display: flex;
  align-items: center;
  gap: 5px;
  min-width: 100px;
}

.height-label {
  font-size: 22px;
  font-weight: 800;
  color: #4caf50;
}

.height-value {
  font-size: 22px;
  font-weight: 800;
}

.genesis-tag {
  background: #ff9800;
  color: #1a1a2e;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 10px;
  font-weight: 700;
  margin-left: 8px;
}

.block-info {
  flex: 1;
  display: flex;
  gap: 25px;
  flex-wrap: wrap;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
}

.info-item i {
  color: #4caf50;
  font-size: 12px;
}

.info-item.hash span {
  font-family: 'Courier New', monospace;
  color: #00e676;
}

.expand-icon {
  color: rgba(255, 255, 255, 0.3);
  transition: color 0.3s;
}

.block-card:hover .expand-icon {
  color: #00e676;
}

/* 区块详情 */
.block-details {
  padding: 20px;
  background: rgba(0, 0, 0, 0.3);
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.detail-grid {
  display: grid;
  gap: 15px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.detail-item .label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.detail-item .value {
  font-size: 13px;
  word-break: break-all;
  font-family: 'Courier New', monospace;
  color: rgba(255, 255, 255, 0.8);
}

.scores {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.score-tag {
  background: rgba(46, 125, 50, 0.3);
  padding: 3px 10px;
  border-radius: 12px;
  font-size: 12px;
}

.no-data {
  color: rgba(255, 255, 255, 0.3);
}

/* 页脚 */
.footer {
  position: relative;
  z-index: 1;
  padding: 50px 20px;
  text-align: center;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  margin-top: 40px;
}

.footer-brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 700;
  color: #4caf50;
  margin-bottom: 15px;
}

.footer-text {
  color: rgba(255, 255, 255, 0.5);
  margin-bottom: 10px;
  font-size: 14px;
}

.copyright {
  color: rgba(255, 255, 255, 0.3);
  font-size: 12px;
}

/* 响应式 */
@media (max-width: 768px) {
  .navbar {
    padding: 12px 20px;
  }

  .nav-brand span {
    display: none;
  }

  .hero-section {
    padding: 100px 20px 30px;
  }

  .hero-title {
    font-size: 28px;
  }

  .title-icon-wrapper {
    width: 70px;
    height: 70px;
    font-size: 30px;
  }

  .search-inputs {
    flex-direction: column;
    align-items: stretch;
  }

  .input-wrapper {
    max-width: none;
  }

  .input-divider {
    display: none;
  }

  .stats-bar {
    gap: 15px;
  }

  .block-main {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .block-info {
    flex-direction: column;
    gap: 8px;
  }
}
</style>
