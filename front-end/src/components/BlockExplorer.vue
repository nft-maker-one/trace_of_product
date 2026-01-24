<template>
  <div class="block-explorer">
    <div class="decor-lines decor-lines-top"></div>
    <div class="decor-lines decor-lines-bottom"></div>

    <!-- 顶部信息与语言切换 -->
    <div class="page-hero">
      <div class="hero-text">
        <div class="eyebrow">
          <i class="fas fa-cubes"></i>
          {{ t('blockExplorer.title') }}
        </div>
        <h2 class="hero-heading">{{ t('blockExplorer.blockList') }}</h2>
        <p class="hero-subtext">{{ t('blockExplorer.queryBlocks') }}</p>
      </div>
      <div class="hero-actions">
        <LanguageSelector mode="simple" />
        <button class="btn ghost" @click="loadLatestBlocks" :disabled="loading">
          <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }"></i>
          {{ t('blockExplorer.loadLatest') }}
        </button>
      </div>
    </div>

    <!-- 区块查询 -->
    <div class="card query-card">
      <div class="card-header">
        <h3 class="card-title">
          <i class="fas fa-search"></i>
          {{ t('blockExplorer.queryBlocks') }}
        </h3>
      </div>
      <div class="query-form">
        <div class="form-group">
          <label class="form-label">
            <i class="fas fa-arrow-up"></i>
            {{ t('blockExplorer.startHeight') }}
          </label>
          <input 
            v-model.number="startHeight" 
            type="number" 
            class="form-control"
            min="0"
            :placeholder="t('blockExplorer.enterStartHeight')"
          >
        </div>
        <div class="form-group">
          <label class="form-label">
            <i class="fas fa-arrow-down"></i>
            {{ t('blockExplorer.endHeight') }}
          </label>
          <input 
            v-model.number="endHeight" 
            type="number" 
            class="form-control"
            min="0"
            :placeholder="t('blockExplorer.enterEndHeight')"
          >
        </div>
        <button class="btn btn-primary query-btn" @click="queryBlocks" :disabled="loading">
          <i class="fas fa-search"></i>
          {{ loading ? t('blockExplorer.loading') : t('blockExplorer.query') }}
        </button>
        <button class="btn btn-secondary" @click="loadLatestBlocks" :disabled="loading">
          <i class="fas fa-history"></i>
          {{ t('blockExplorer.loadLatest') }}
        </button>
      </div>
    </div>

    <!-- 区块列表 -->
    <div class="card blocks-card">
      <div class="card-header">
        <h3 class="card-title">
          <i class="fas fa-list"></i>
          {{ t('blockExplorer.blockList') }}
          <span class="block-count">({{ blocks.length }} {{ t('blockExplorer.blocks') }})</span>
        </h3>
      </div>

      <div v-if="loading" class="loading-container">
        <div class="spinner"></div>
        <p>{{ t('blockExplorer.loadingBlocks') }}</p>
      </div>

      <div v-else-if="blocks.length === 0" class="empty-state">
        <i class="fas fa-inbox"></i>
        <p>{{ t('blockExplorer.noBlocks') }}</p>
      </div>

      <div v-else class="blocks-list">
        <div 
          v-for="block in blocks" 
          :key="block.height" 
          class="block-item"
          :class="{ 'genesis': block.height === 0, 'expanded': expandedBlock === block.height }"
          @click="toggleBlock(block.height)"
        >
          <div class="block-header">
            <div class="block-height">
              <i class="fas fa-cube"></i>
              <span class="height-label">{{ t('blockExplorer.height') }}</span>
              <span class="height-value">{{ block.height }}</span>
              <span v-if="block.height === 0" class="genesis-badge">Genesis</span>
            </div>
            <div class="block-meta">
              <span class="block-time">
                <i class="fas fa-clock"></i>
                {{ formatTimestamp(block.timestamp) }}
              </span>
              <span class="block-leader">
                <i class="fas fa-user"></i>
                {{ block.leader }}
              </span>
            </div>
            <i class="fas fa-chevron-down expand-icon"></i>
          </div>

          <div class="block-details" v-show="expandedBlock === block.height">
            <div class="detail-grid">
              <div class="detail-item">
                <span class="detail-label">{{ t('blockExplorer.dataHash') }}</span>
                <div class="hash-value">
                  <code>{{ formatHash(block.dataHash) }}</code>
                  <button class="copy-btn" @click.stop="copyToClipboard(block.dataHash)">
                    <i class="fas fa-copy"></i>
                  </button>
                </div>
              </div>
              <div class="detail-item">
                <span class="detail-label">{{ t('blockExplorer.prevHash') }}</span>
                <div class="hash-value">
                  <code>{{ formatHash(block.prevHash) }}</code>
                  <button class="copy-btn" @click.stop="copyToClipboard(block.prevHash)">
                    <i class="fas fa-copy"></i>
                  </button>
                </div>
              </div>
              <div class="detail-item">
                <span class="detail-label">{{ t('blockExplorer.nonce') }}</span>
                <span class="detail-value">{{ block.nonce }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">{{ t('blockExplorer.version') }}</span>
                <span class="detail-value">{{ block.version }}</span>
              </div>
              <div class="detail-item full-width">
                <span class="detail-label">{{ t('blockExplorer.scores') }}</span>
                <div class="scores-list">
                  <span v-if="Object.keys(block.scores || {}).length === 0" class="no-scores">
                    {{ t('blockExplorer.noScores') }}
                  </span>
                  <span v-else v-for="(score, nodeId) in block.scores" :key="nodeId" class="score-item">
                    {{ nodeId }}: {{ score }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '../services/api'
import LanguageSelector from './LanguageSelector.vue'

const { t } = useI18n()

const loading = ref(false)
const blocks = ref([])
const startHeight = ref(0)
const endHeight = ref(10)
const expandedBlock = ref(null)

// 方法
async function queryBlocks() {
  if (startHeight.value < 0 || endHeight.value < 0) {
    alert(t('blockExplorer.invalidHeight'))
    return
  }
  if (startHeight.value > endHeight.value) {
    alert(t('blockExplorer.startGreaterThanEnd'))
    return
  }

  loading.value = true
  try {
    const response = await api.getBlocksByRange(startHeight.value, endHeight.value)
    if (response && response.success) {
      blocks.value = response.blocks || []
      // 按高度降序排列
      blocks.value.sort((a, b) => b.height - a.height)
    } else {
      console.error('获取区块失败:', response)
      blocks.value = []
    }
  } catch (error) {
    console.error('查询区块失败:', error)
    alert(t('blockExplorer.queryFailed'))
  } finally {
    loading.value = false
  }
}

async function loadLatestBlocks() {
  loading.value = true
  try {
    // 先获取所有区块（从0到一个较大的数）
    const response = await api.getBlocksByRange(0, 100)
    if (response && response.success) {
      blocks.value = response.blocks || []
      blocks.value.sort((a, b) => b.height - a.height)
      
      // 更新查询范围
      if (blocks.value.length > 0) {
        endHeight.value = Math.max(...blocks.value.map(b => b.height))
        startHeight.value = Math.max(0, endHeight.value - 10)
      }
    }
  } catch (error) {
    console.error('加载最新区块失败:', error)
  } finally {
    loading.value = false
  }
}

function toggleBlock(height) {
  expandedBlock.value = expandedBlock.value === height ? null : height
}

function formatTimestamp(timestamp) {
  if (!timestamp || timestamp === 0) return 'Genesis'
  return new Date(timestamp * 1000).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

function formatHash(hash) {
  if (!hash) return '-'
  if (hash.length > 20) {
    return hash.substring(0, 10) + '...' + hash.substring(hash.length - 10)
  }
  return hash
}

function copyToClipboard(text) {
  navigator.clipboard.writeText(text).then(() => {
    alert(t('blockExplorer.copied'))
  }).catch(err => {
    console.error('复制失败:', err)
  })
}

onMounted(() => {
  loadLatestBlocks()
})
</script>

<style scoped>
/* 变量 */
:root {
  --primary: #2e7d32;
  --primary-light: #4caf50;
  --primary-dark: #1b5e20;
  --secondary: #0288d1;
  --dark: #333;
  --light-gray: #e0e0e0;
  --border-color: #ddd;
}

.block-explorer {
  position: relative;
  padding: 30px 20px 60px;
  max-width: 1200px;
  margin: 0 auto;
  background: radial-gradient(circle at 10% 20%, rgba(46, 125, 50, 0.08), transparent 35%),
              radial-gradient(circle at 90% 10%, rgba(2, 136, 209, 0.08), transparent 30%),
              linear-gradient(180deg, #f7faf9 0%, #eef5f0 100%);
  border-radius: 24px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.decor-lines {
  position: absolute;
  left: -20%;
  right: -20%;
  height: 120px;
  background: linear-gradient(90deg, rgba(46, 125, 50, 0.08), rgba(255, 255, 255, 0), rgba(2, 136, 209, 0.08));
  filter: blur(6px);
  pointer-events: none;
}

.decor-lines-top { top: 0; }
.decor-lines-bottom { bottom: 0; }

.page-hero {
  position: relative;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  padding: 10px 5px 25px;
}

.hero-text {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.eyebrow {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  border-radius: 999px;
  background: rgba(46, 125, 50, 0.1);
  color: var(--primary);
  font-weight: 700;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  width: fit-content;
}

.hero-heading {
  margin: 0;
  font-size: 28px;
  font-weight: 800;
  color: #1f2d3d;
}

.hero-subtext {
  margin: 0;
  color: #50616e;
  font-size: 14px;
}

.hero-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.btn.ghost {
  background: transparent;
  border: 1px solid rgba(46, 125, 50, 0.3);
  color: var(--primary);
  padding: 10px 16px;
  border-radius: 10px;
  font-weight: 600;
  transition: all 0.2s;
}

.btn.ghost:hover:not(:disabled) {
  background: rgba(46, 125, 50, 0.08);
  box-shadow: 0 8px 20px rgba(46, 125, 50, 0.15);
}
}

/* 卡片样式 */
.card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.08);
  border: 1px solid var(--light-gray);
}


.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid var(--light-gray);
}

.card-title i {
  color: var(--primary);
}

.block-count {
  font-size: 14px;
  color: #666;
  font-weight: normal;
}
  background: #9c27b0;
}

.stat-content {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--dark);
}

.stat-label {
  font-size: 12px;
  color: #666;
}

/* 查询表单 */
.query-form {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  align-items: flex-end;
}

.query-form .form-group {
  flex: 1;
  min-width: 150px;
}

.form-label {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
  font-weight: 600;
  color: var(--dark);
  font-size: 14px;
}

.form-label i {
  color: var(--primary);
}

.form-control {
  width: 100%;
  padding: 10px 14px;
  border: 2px solid var(--light-gray);
  border-radius: 8px;
  font-size: 14px;
  transition: all 0.2s;
}

.form-control:focus {
  outline: none;
  border-color: var(--primary);
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: all 0.2s;
}

.btn-primary {
  background: var(--primary);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: var(--primary-dark);
}

.btn-secondary {
  background: #f5f5f5;
  color: var(--dark);
  border: 1px solid var(--light-gray);
}

.btn-secondary:hover:not(:disabled) {
  background: #e0e0e0;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 加载状态 */
.loading-container {
  text-align: center;
  padding: 40px;
  color: #666;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid var(--light-gray);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 15px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 40px;
  color: #999;
}

.empty-state i {
  font-size: 48px;
  margin-bottom: 15px;
}

/* 区块列表 */
.blocks-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.block-item {
  border: 2px solid var(--light-gray);
  border-radius: 10px;
  overflow: hidden;
  transition: all 0.2s;
  cursor: pointer;
}

.block-item:hover {
  border-color: var(--primary-light);
  box-shadow: 0 2px 8px rgba(46, 125, 50, 0.1);
}

.block-item.genesis {
  border-color: #ff9800;
  background: linear-gradient(135deg, #fff8e1 0%, #ffecb3 100%);
}

.block-item.expanded {
  border-color: var(--primary);
}

.block-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 15px;
  background: #fafafa;
}

.block-item.genesis .block-header {
  background: transparent;
}

.block-height {
  display: flex;
  align-items: center;
  gap: 8px;
}

.block-height i {
  color: var(--primary);
  font-size: 18px;
}

.height-label {
  font-size: 12px;
  color: #666;
}

.height-value {
  font-size: 20px;
  font-weight: 700;
  color: var(--dark);
}

.genesis-badge {
  background: #ff9800;
  color: white;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 10px;
  font-weight: 600;
}

.block-meta {
  display: flex;
  gap: 20px;
  font-size: 13px;
  color: #666;
}

.block-meta span {
  display: flex;
  align-items: center;
  gap: 5px;
}

.expand-icon {
  color: #999;
  transition: transform 0.2s;
}

.block-item.expanded .expand-icon {
  transform: rotate(180deg);
}

/* 区块详情 */
.block-details {
  padding: 15px;
  background: white;
  border-top: 1px solid var(--light-gray);
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.detail-item.full-width {
  grid-column: span 2;
}

.detail-label {
  font-size: 12px;
  color: #666;
  font-weight: 600;
}

.detail-value {
  font-size: 14px;
  color: var(--dark);
}

.hash-value {
  display: flex;
  align-items: center;
  gap: 8px;
}

.hash-value code {
  font-family: 'Courier New', monospace;
  font-size: 12px;
  background: #f5f5f5;
  padding: 5px 10px;
  border-radius: 5px;
  word-break: break-all;
}

.copy-btn {
  background: none;
  border: none;
  color: var(--primary);
  cursor: pointer;
  padding: 5px;
  border-radius: 5px;
  transition: background 0.2s;
}

.copy-btn:hover {
  background: #e8f5e9;
}

.scores-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.score-item {
  background: #e3f2fd;
  color: #1976d2;
  padding: 4px 10px;
  border-radius: 15px;
  font-size: 12px;
}

.no-scores {
  color: #999;
  font-style: italic;
}

/* 响应式 */
@media (max-width: 768px) {
  .page-hero {
    flex-direction: column;
    align-items: flex-start;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .query-form {
    flex-direction: column;
  }

  .query-form .form-group {
    width: 100%;
  }

  .detail-grid {
    grid-template-columns: 1fr;
  }

  .detail-item.full-width {
    grid-column: span 1;
  }

  .block-meta {
    flex-direction: column;
    gap: 5px;
  }
}
</style>
