<template>
  <div class="data-query">
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">
          <i class="fas fa-search"></i>
          {{ t('dataQuery.title') }}
        </h3>
        <div class="query-help">
          <i class="fas fa-info-circle"></i>
          <span>{{ t('dataQuery.description') }}</span>
        </div>
      </div>

      <div class="form-row">
        <div class="form-col">
          <div class="form-group">
            <label class="form-label" for="query-id">
              <i class="fas fa-fingerprint"></i>
              {{ t('dataQuery.eggplantId') }} *
            </label>
            <input
              v-model="queryId"
              type="text"
              id="query-id"
              class="form-control"
              :placeholder="t('dataQuery.enterEggplantId')"
              @keyup.enter="handleQuery"
            >
          </div>
        </div>

        <div class="form-col">
          <div class="select-wrapper">
            <label class="form-label" for="select-node1">
              <i class="fas fa-server"></i>
              {{ t('dataQuery.selectNode') }} *
            </label>
            <select id="select-node1" v-model="selectedNode">
              <option value="">{{ t('dataQuery.chooseNode') }}</option>
              <option v-for="node in nodes" :key="node.id" :value="node.id">
                {{ node.id }}
              </option>
            </select>
          </div>
        </div>

        <div class="form-col">
          <button class="btn btn-primary" @click="handleQuery" :disabled="querying" style="margin-top: 28px; width: 100%;">
            <i class="fas fa-search"></i>
            {{ querying ? t('dataQuery.querying') : t('dataQuery.queryProductInfo') }}
            <div v-if="querying" class="spinner"></div>
          </button>
        </div>
      </div>

      <!-- 最近查询节点 -->
      <div v-if="recentQueryNodes.length > 0" class="recent-nodes">
        <h4 class="history-title">
          <i class="fas fa-server"></i>
          {{ t('dataQuery.recentNodes') || '最近查询节点' }}
        </h4>
        <div class="recent-nodes-list">
          <div
            v-for="(rnode, index) in recentQueryNodes"
            :key="rnode.id"
            class="recent-node-item"
            :class="{ active: selectedNode === rnode.id }"
            @click="selectedNode = rnode.id"
          >
            <span class="node-rank">#{{ index + 1 }}</span>
            <span class="node-id">{{ rnode.id }}</span>
            <span class="node-time">{{ rnode.time }}</span>
          </div>
        </div>
      </div>

      <!-- 查询历史 -->
      <div v-if="queryHistory.length > 0" class="query-history">
        <h4 class="history-title">
          <i class="fas fa-history"></i>
          {{ t('dataQuery.queryHistory') }}
        </h4>
        <div class="history-tags">
          <span
            v-for="history in queryHistory.slice(0, 5)"
            :key="history.id"
            class="history-tag"
            @click="selectHistory(history.id)"
          >
            {{ history.id }}
            <i class="fas fa-times" @click.stop="removeHistory(history.id)"></i>
          </span>
        </div>
      </div>
    </div>

    <!-- 查询结果 -->
    <div v-if="queryResult" class="card result-card">
      <div class="card-header">
        <h3 class="card-title">
          <i class="fas fa-file-alt"></i>
          {{ t('dataQuery.queryResults') }}
          <span class="result-id">ID: {{ queryId }}</span>
        </h3>
        <div class="result-actions">
          <button class="btn btn-outline" @click="exportResult">
            <i class="fas fa-download"></i>
            {{ t('dataQuery.export') }}
          </button>
          <button class="btn btn-outline" @click="printResult">
            <i class="fas fa-print"></i>
            {{ t('dataQuery.print') }}
          </button>
        </div>
      </div>

      <div class="result-content">
        <!-- 基本信息 -->
        <div class="result-section">
          <h4 class="section-title">
            <i class="fas fa-info-circle"></i>
            {{ t('dataQuery.basicInfo') }}
          </h4>
          <div class="result-grid">
            <div class="result-item">
              <div class="result-label">{{ t('dataQuery.productId') }}:</div>
              <div class="result-value">{{ queryResult.product_id }}</div>
            </div>
            <div class="result-item">
              <div class="result-label">{{ t('dataQuery.queryTime') }}:</div>
              <div class="result-value">{{ queryTime }}</div>
            </div>
            <div class="result-item">
              <div class="result-label">{{ t('dataQuery.queryNode') }}:</div>
              <div class="result-value">{{ selectedNode }}</div>
            </div>
            <div class="result-item">
              <div class="result-label">{{ t('dataQuery.dataStatus') }}:</div>
              <div class="result-value">
                <span class="status-badge" :class="getDataStatusClass">
                  {{ getDataStatusText }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- 详细数据 -->
        <div class="result-section">
          <h4 class="section-title">
            <i class="fas fa-list-alt"></i>
            {{ t('dataQuery.detailedData') }}
          </h4>

          <div class="data-timeline">
            <div class="timeline-item" v-for="stage in productStages" :key="stage.key">
              <div class="timeline-header">
                <div class="stage-icon" :class="stage.status">
                  <i :class="stage.icon"></i>
                </div>
                <div class="stage-info">
                  <h5>{{ t(`dataQuery.productStages.${stage.key}`) }}</h5>
                  <div class="stage-time">{{ t('dataQuery.blockHeightLabel') }}: {{ getStageHeight(stage.key) }}</div>
                </div>
                <div class="stage-status">
                  <span class="status-indicator" :class="getStageStatus(stage.key)"></span>
                  {{ getStageStatusText(stage.key) }}
                </div>
              </div>

              <div class="timeline-content">
                <div class="hash-display">
                  <div class="hash-label">{{ t('dataQuery.hashValue') }}:</div>
                  <div class="hash-value">
                    <code>{{ getStageHash(stage.key) }}</code>
                    <button class="copy-btn" @click="copyHash(getStageHash(stage.key))">
                      <i class="fas fa-copy"></i>
                    </button>
                  </div>
                </div>

                <div v-if="stage.details" class="stage-details">
                  <div class="detail-row" v-for="detail in stage.details" :key="detail.key">
                    <div class="detail-label">{{ t(`dataQuery.stageDetails.${detail.key}`) }}:</div>
                    <div class="detail-value">{{ detail.value }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 区块链验证 -->
        <div class="result-section">
          <h4 class="section-title">
            <i class="fas fa-link"></i>
            {{ t('dataQuery.blockchainVerification') }}
          </h4>
          <div class="verification-info">
            <div class="verification-item">
              <div class="verification-label">{{ t('dataQuery.blockHeightLabel') }}:</div>
              <div class="verification-value">{{ queryResult.block_height || 'N/A' }}</div>
            </div>
            <div class="verification-item">
              <div class="verification-label">{{ t('dataQuery.txHash') }}:</div>
              <div class="verification-value hash-value">
                <code>{{ queryResult.tx_hash || 'N/A' }}</code>
                <button v-if="queryResult.tx_hash" class="copy-btn" @click="copyHash(queryResult.tx_hash)">
                  <i class="fas fa-copy"></i>
                </button>
              </div>
            </div>
            <div class="verification-item">
              <div class="verification-label">{{ t('dataQuery.verificationStatus') }}:</div>
              <div class="verification-value">
                <span class="status-badge" :class="getVerificationClass">
                  {{ getVerificationText }}
                </span>
              </div>
            </div>
            <div class="verification-item">
              <div class="verification-label">{{ t('dataQuery.timestampLabel') }}:</div>
              <div class="verification-value">{{ queryResult.timestamp || 'N/A' }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 无结果提示 -->
    <div v-else-if="showNoResult" class="card empty-card">
      <div class="empty-state">
        <i class="fas fa-search"></i>
        <h3>{{ t('dataQuery.noResults') }}</h3>
        <p>{{ t('dataQuery.enterCorrectId') }}</p>
        <button class="btn btn-outline" @click="showHelp">
          <i class="fas fa-question-circle"></i>
          {{ t('dataQuery.viewHelp') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '../services/api'

const { t } = useI18n()

const queryId = ref('')
const selectedNode = ref('')
const nodes = ref([])
const querying = ref(false)
const queryResult = ref(null)
const queryTime = ref('')
const showNoResult = ref(false)

// 查询历史
const queryHistory = ref([])

// 最近三次查询使用的节点
const recentQueryNodes = ref([])

// 从 localStorage 加载最近查询节点
function loadRecentNodes() {
  try {
    const saved = localStorage.getItem('recentQueryNodes')
    if (saved) {
      recentQueryNodes.value = JSON.parse(saved)
    }
  } catch (e) {
    console.error('加载最近查询节点失败:', e)
  }
}

// 保存最近查询节点到 localStorage
function saveRecentNodes() {
  try {
    localStorage.setItem('recentQueryNodes', JSON.stringify(recentQueryNodes.value))
  } catch (e) {
    console.error('保存最近查询节点失败:', e)
  }
}

// 添加节点到最近查询记录
function addToRecentNodes(nodeId, nodeAddr) {
  // 移除已存在的相同节点
  recentQueryNodes.value = recentQueryNodes.value.filter(n => n.id !== nodeId)
  
  // 添加到最前面
  recentQueryNodes.value.unshift({
    id: nodeId,
    addr: nodeAddr,
    time: new Date().toLocaleString('zh-CN', { hour12: false })
  })
  
  // 只保留最近3个
  if (recentQueryNodes.value.length > 3) {
    recentQueryNodes.value = recentQueryNodes.value.slice(0, 3)
  }
  
  saveRecentNodes()
}

// 产品阶段定义
const productStages = [
  {
    key: 'product',
    icon: 'fas fa-box',
    status: 'completed',
    details: [
      { key: 'height', label: t('dataQuery.stageDetails.height') },
      { key: 'time', label: t('dataQuery.stageDetails.time') },
      { key: 'operator', label: t('dataQuery.stageDetails.operator') }
    ]
  },
  {
    key: 'transport',
    icon: 'fas fa-truck',
    status: 'completed',
    details: [
      { key: 'height', label: t('dataQuery.stageDetails.height') },
      { key: 'time', label: t('dataQuery.stageDetails.time') },
      { key: 'vehicle', label: t('dataQuery.stageDetails.vehicle') },
      { key: 'driver', label: t('dataQuery.stageDetails.driver') }
    ]
  },
  {
    key: 'process',
    icon: 'fas fa-industry',
    status: 'completed',
    details: [
      { key: 'height', label: t('dataQuery.stageDetails.height') },
      { key: 'time', label: t('dataQuery.stageDetails.time') },
      { key: 'factory', label: t('dataQuery.stageDetails.factory') },
      { key: 'method', label: t('dataQuery.stageDetails.method') }
    ]
  },
  {
    key: 'storage',
    icon: 'fas fa-warehouse',
    status: 'completed',
    details: [
      { key: 'height', label: t('dataQuery.stageDetails.height') },
      { key: 'time', label: t('dataQuery.stageDetails.time') },
      { key: 'location', label: t('dataQuery.stageDetails.location') },
      { key: 'temperature', label: t('dataQuery.stageDetails.temperature') }
    ]
  },
  {
    key: 'sell',
    icon: 'fas fa-shopping-cart',
    status: 'completed',
    details: [
      { key: 'height', label: t('dataQuery.stageDetails.height') },
      { key: 'time', label: t('dataQuery.stageDetails.time') },
      { key: 'seller', label: t('dataQuery.stageDetails.seller') },
      { key: 'buyer', label: t('dataQuery.stageDetails.buyer') }
    ]
  }
]

// 计算属性
const getDataStatusClass = computed(() => {
  if (!queryResult.value) return 'unknown'
  return queryResult.value.is_complete ? 'success' : 'warning'
})

const getDataStatusText = computed(() => {
  if (!queryResult.value) return t('dataUpload.unknown')
  return queryResult.value.is_complete ? t('dataQuery.dataComplete') : t('dataQuery.dataIncomplete')
})

const getVerificationClass = computed(() => {
  if (!queryResult.value) return 'unknown'
  return queryResult.value.verified ? 'success' : 'error'
})

const getVerificationText = computed(() => {
  if (!queryResult.value) return t('dataQuery.notVerified')
  return queryResult.value.verified ? t('dataQuery.verified') : t('dataQuery.notVerified')
})

// 方法
async function handleQuery() {
  if (!queryId.value.trim()) {
    alert(t('dataQuery.enterEggplantIdAlert'))
    return
  }

  if (!selectedNode.value) {
    alert(t('dataQuery.selectNodeAlert'))
    return
  }

  querying.value = true
  showNoResult.value = false
  queryResult.value = null

  try {
    // 查找节点地址
    const node = nodes.value.find(n => n.id === selectedNode.value)
    if (!node) {
      throw new Error('选择的节点不存在')
    }

    // 调用API查询数据
    const response = await api.queryData(queryId.value, node.addr)

    if (response && response.msg) {
      // 解析返回的数据
      const data = typeof response.msg === 'string' ? JSON.parse(response.msg) : response.msg

      // 设置查询结果
      queryResult.value = {
        product_id: queryId.value,
        product_height: data.product_height || 0,
        product_hash: data.product_hash || '',
        transport_height: data.transport_height || 0,
        transport_hash: data.transport_hash || '',
        process_height: data.process_height || 0,
        process_hash: data.process_hash || '',
        storage_height: data.storage_height || 0,
        storage_hash: data.storage_hash || '',
        sell_height: data.sell_height || 0,
        sell_hash: data.sell_hash || '',
        block_height: data.block_height,
        tx_hash: data.tx_hash,
        timestamp: data.timestamp,
        verified: true,
        is_complete: true
      }

      // 记录查询时间
      queryTime.value = new Date().toLocaleString('zh-CN')

      // 添加到查询历史
      addToHistory(queryId.value)
      
      // 记录使用的查询节点
      addToRecentNodes(node.id, node.addr)
    } else {
      showNoResult.value = true
    }
  } catch (error) {
    console.error('查询失败:', error)
    alert(t('dataQuery.queryFailedAlert', { error: error.message }))
    showNoResult.value = true
  } finally {
    querying.value = false
  }
}

function getStageHeight(stageKey) {
  if (!queryResult.value) return 'N/A'
  const height = queryResult.value[`${stageKey}_height`]
  return height || height === 0 ? height : 'N/A'
}

function getStageHash(stageKey) {
  if (!queryResult.value) return 'N/A'
  const hash = queryResult.value[`${stageKey}_hash`]
  return hash || 'N/A'
}

function getStageStatus(stageKey) {
  const hash = getStageHash(stageKey)
  return hash !== 'N/A' ? 'verified' : 'missing'
}

function getStageStatusText(stageKey) {
  const hash = getStageHash(stageKey)
  return hash !== 'N/A' ? t('dataQuery.status.completed') : t('dataQuery.status.missing')
}

function addToHistory(id) {
  // 检查是否已存在
  const exists = queryHistory.value.some(item => item.id === id)
  if (!exists) {
    queryHistory.value.unshift({
      id,
      time: new Date().toLocaleString('zh-CN', { hour12: false })
    })

    // 限制历史记录数量
    if (queryHistory.value.length > 10) {
      queryHistory.value = queryHistory.value.slice(0, 10)
    }
  }
}

function selectHistory(id) {
  queryId.value = id
}

function removeHistory(id) {
  queryHistory.value = queryHistory.value.filter(item => item.id !== id)
}

function copyHash(hash) {
  if (hash && hash !== 'N/A') {
    navigator.clipboard.writeText(hash)
      .then(() => {
        alert(t('dataQuery.hashCopied'))
      })
      .catch(err => {
        console.error('复制失败:', err)
      })
  }
}

function exportResult() {
  if (!queryResult.value) return

  const data = {
    查询ID: queryId.value,
    查询时间: queryTime.value,
    查询节点: selectedNode.value,
    ...queryResult.value
  }

  const dataStr = JSON.stringify(data, null, 2)
  const dataBlob = new Blob([dataStr], { type: 'application/json' })
  const url = URL.createObjectURL(dataBlob)

  const link = document.createElement('a')
  link.href = url
  link.download = `农产品追溯_${queryId.value}_${new Date().getTime()}.json`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

function printResult() {
  window.print()
}

function showHelp() {
  alert(t('dataQuery.queryHelp'))
}

// 加载节点数据
async function loadNodes() {
  try {
    const data = await api.getNodes()
    if (Array.isArray(data)) {
      nodes.value = data
      if (nodes.value.length > 0 && !selectedNode.value) {
        selectedNode.value = nodes.value[0].id
      }
    }
  } catch (error) {
    console.error('加载节点失败:', error)
  }
}

onMounted(() => {
  loadNodes()
  loadRecentNodes()
})
</script>

<style scoped>
/* 样式变量 */
:root {
  --primary: #2e7d32;
  --primary-light: #4caf50;
  --primary-dark: #1b5e20;
  --secondary: #0288d1;
  --secondary-dark: #0277bd;
  --success: #4caf50;
  --warning: #ff9800;
  --error: #f44336;
  --info: #2196f3;
  --light: #f5f7fa;
  --dark: #263238;
  --gray: #78909c;
  --light-gray: #eceff1;
  --border-color: #e0e0e0;
  --card-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  --card-shadow-hover: 0 8px 30px rgba(0, 0, 0, 0.12);
  --transition: all 0.3s ease;
}

.data-query {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
  min-height: 100vh;
}

/* 卡片样式 */
.card {
  background: white;
  border-radius: 16px;
  padding: 30px;
  margin-bottom: 30px;
  box-shadow: var(--card-shadow);
  transition: var(--transition);
  border: 1px solid var(--light-gray);
  animation: fadeIn 0.5s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.card:hover {
  box-shadow: var(--card-shadow-hover);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--light-gray);
  flex-wrap: wrap;
  gap: 15px;
}

.card-title {
  font-size: 22px;
  font-weight: 600;
  color: var(--primary-dark);
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 0;
}

.card-title i {
  color: var(--primary);
  font-size: 24px;
}

/* 查询帮助样式 */
.query-help {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 15px;
  color: var(--gray);
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  padding: 10px 16px;
  border-radius: 10px;
  border-left: 3px solid var(--info);
}

.query-help i {
  color: var(--info);
  font-size: 16px;
}

/* 表单样式 */
.form-row {
  display: flex;
  gap: 30px;
  margin-bottom: 30px;
  flex-wrap: wrap;
  align-items: flex-end;
}

.form-col {
  flex: 1;
  min-width: 250px;
}

.form-group {
  margin-bottom: 0;
}

.form-label {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
  font-weight: 600;
  color: var(--dark);
  font-size: 15px;
}

.form-label i {
  width: 20px;
  color: var(--primary);
}

.form-control {
  width: 100%;
  padding: 14px 16px;
  border: 2px solid var(--light-gray);
  border-radius: 10px;
  font-size: 15px;
  transition: var(--transition);
  background-color: white;
  color: var(--dark);
}

.form-control:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(46, 125, 50, 0.1);
}

.form-control::placeholder {
  color: var(--gray);
  opacity: 0.7;
}

/* 选择框样式 */
.select-wrapper {
  position: relative;
}

.select-wrapper select {
  width: 100%;
  padding: 14px 16px;
  padding-right: 45px;
  border: 2px solid var(--light-gray);
  border-radius: 10px;
  font-size: 15px;
  background-color: white;
  color: var(--dark);
  cursor: pointer;
  transition: var(--transition);
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%2378909c' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 16px center;
  background-size: 16px;
}

.select-wrapper select:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(46, 125, 50, 0.1);
}

/* 按钮样式 */
.btn {
  padding: 14px 28px;
  border: none;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  white-space: nowrap;
}

.btn-primary {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(46, 125, 50, 0.2);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(46, 125, 50, 0.3);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-outline {
  background: transparent;
  border: 2px solid var(--primary);
  color: var(--primary);
  padding: 10px 20px;
  font-size: 14px;
}

.btn-outline:hover {
  background: var(--primary);
  color: white;
}

/* 最近查询节点 */
.recent-nodes {
  margin-top: 20px;
  padding: 15px;
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
  border-radius: 12px;
  border: 1px solid #90caf9;
}

.recent-nodes-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.recent-node-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 15px;
  background: white;
  border: 2px solid #90caf9;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.recent-node-item:hover {
  background: #e3f2fd;
  border-color: #2196f3;
}

.recent-node-item.active {
  background: #2196f3;
  border-color: #1976d2;
  color: white;
}

.recent-node-item.active .node-time {
  color: rgba(255, 255, 255, 0.8);
}

.node-rank {
  font-weight: 700;
  color: #1976d2;
  font-size: 12px;
}

.recent-node-item.active .node-rank {
  color: white;
}

.node-id {
  font-weight: 600;
  font-size: 14px;
}

.node-time {
  font-size: 12px;
  color: #666;
}

/* 查询历史 */
.query-history {
  margin-top: 20px;
  padding: 15px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 12px;
  border: 1px solid var(--border-color);
  animation: slideIn 0.4s ease;
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

.history-title {
  font-size: 17px;
  margin-bottom: 15px;
  color: var(--dark);
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 600;
}

.history-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.history-tag {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 15px;
  background: white;
  border: 1px solid var(--border-color);
  border-radius: 25px;
  font-size: 14px;
  color: var(--dark);
  cursor: pointer;
  transition: var(--transition);
  font-family: 'Courier New', monospace;
  font-weight: 500;
}

.history-tag:hover {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(46, 125, 50, 0.2);
}

.history-tag i {
  font-size: 12px;
  opacity: 0.7;
  transition: var(--transition);
}

.history-tag:hover i {
  opacity: 1;
  color: white;
}

/* 查询结果卡片 */
.result-card {
  margin-top: 30px;
  animation: scaleIn 0.5s ease;
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.result-id {
  font-size: 14px;
  font-weight: normal;
  color: var(--gray);
  margin-left: 12px;
  background: var(--light-gray);
  padding: 4px 12px;
  border-radius: 15px;
  font-family: 'Courier New', monospace;
}

.result-actions {
  display: flex;
  gap: 12px;
}

.result-content {
  padding: 10px 0;
}

/* 结果区块 */
.result-section {
  margin-bottom: 35px;
  padding-bottom: 25px;
  border-bottom: 1px solid var(--light-gray);
  animation: fadeInUp 0.5s ease;
  animation-fill-mode: both;
}

.result-section:nth-child(1) { animation-delay: 0.1s; }
.result-section:nth-child(2) { animation-delay: 0.2s; }
.result-section:nth-child(3) { animation-delay: 0.3s; }

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

.result-section:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.section-title {
  font-size: 19px;
  margin-bottom: 25px;
  color: var(--primary-dark);
  display: flex;
  align-items: center;
  gap: 12px;
  font-weight: 600;
  padding-bottom: 10px;
  border-bottom: 2px solid var(--light-gray);
}

.section-title i {
  color: var(--primary);
}

/* 基本信息网格 */
.result-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.result-item {
  display: flex;
  align-items: center;
  padding: 18px 20px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 12px;
  transition: var(--transition);
  border: 1px solid transparent;
}

.result-item:hover {
  transform: translateY(-2px);
  border-color: var(--primary);
  box-shadow: 0 4px 15px rgba(46, 125, 50, 0.1);
}

.result-label {
  font-weight: 600;
  min-width: 90px;
  color: var(--gray);
  margin-right: 15px;
  font-size: 14px;
}

.result-value {
  color: var(--dark);
  flex: 1;
  font-size: 15px;
  font-weight: 500;
}

/* 时间线样式 */
.data-timeline {
  max-width: 900px;
  margin: 0 auto;
  position: relative;
  padding-left: 30px;
}

.data-timeline::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: linear-gradient(to bottom, var(--primary), var(--primary-light));
  border-radius: 3px;
}

.timeline-item {
  margin-bottom: 25px;
  padding: 25px;
  background: white;
  border-radius: 15px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.06);
  border: 1px solid var(--border-color);
  position: relative;
  transition: var(--transition);
  animation: slideInRight 0.5s ease;
}

.timeline-item:hover {
  transform: translateX(5px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
}

@keyframes slideInRight {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.timeline-item::before {
  content: '';
  position: absolute;
  left: -30px;
  top: 40px;
  width: 20px;
  height: 20px;
  background: var(--primary);
  border-radius: 50%;
  border: 4px solid white;
  box-shadow: 0 0 0 2px var(--primary);
}

.timeline-header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 15px;
}

.stage-icon {
  width: 50px;
  height: 50px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 15px;
  font-size: 22px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
}

.stage-icon.completed {
  background: linear-gradient(135deg, var(--success) 0%, #66bb6a 100%);
  color: white;
}

.stage-info {
  flex: 1;
  min-width: 200px;
}

.stage-info h5 {
  font-size: 18px;
  margin-bottom: 8px;
  color: var(--dark);
  font-weight: 600;
}

.stage-time {
  font-size: 14px;
  color: var(--gray);
  display: flex;
  align-items: center;
  gap: 8px;
}

.stage-status {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  color: var(--gray);
  background: var(--light);
  padding: 8px 15px;
  border-radius: 20px;
}

.status-indicator {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.status-indicator.verified {
  background-color: var(--success);
  box-shadow: 0 0 0 3px rgba(76, 175, 80, 0.2);
}

.status-indicator.missing {
  background-color: var(--error);
  box-shadow: 0 0 0 3px rgba(244, 67, 54, 0.2);
}

/* 哈希显示 */
.hash-display {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
  padding: 15px;
  background: linear-gradient(135deg, #f5f5f5 0%, #eeeeee 100%);
  border-radius: 10px;
  border: 1px solid var(--border-color);
}

.hash-label {
  font-weight: 600;
  min-width: 70px;
  color: var(--gray);
  margin-right: 15px;
  font-size: 14px;
}

.hash-value {
  flex: 1;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  color: var(--dark);
  word-break: break-all;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid var(--border-color);
}

.copy-btn {
  background: transparent;
  border: none;
  color: var(--primary);
  cursor: pointer;
  padding: 6px;
  font-size: 15px;
  border-radius: 6px;
  transition: var(--transition);
  margin-left: 10px;
}

.copy-btn:hover {
  background: var(--primary);
  color: white;
  transform: scale(1.1);
}

/* 阶段详情 */
.stage-details {
  margin-top: 15px;
  padding: 15px;
  background: linear-gradient(135deg, #f9f9f9 0%, #f5f5f5 100%);
  border-radius: 10px;
  border: 1px solid var(--border-color);
}

.detail-row {
  display: flex;
  margin-bottom: 12px;
  font-size: 14px;
  padding: 8px 0;
  border-bottom: 1px dashed var(--border-color);
}

.detail-row:last-child {
  margin-bottom: 0;
  border-bottom: none;
}

.detail-label {
  font-weight: 600;
  min-width: 100px;
  color: var(--gray);
}

.detail-value {
  color: var(--dark);
  flex: 1;
  font-weight: 500;
}

/* 验证信息 */
.verification-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 20px;
}

.verification-item {
  display: flex;
  align-items: center;
  padding: 18px 20px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 12px;
  transition: var(--transition);
  border: 1px solid transparent;
}

.verification-item:hover {
  transform: translateY(-2px);
  border-color: var(--info);
  box-shadow: 0 4px 15px rgba(33, 150, 243, 0.1);
}

.verification-label {
  font-weight: 600;
  min-width: 110px;
  color: var(--gray);
  margin-right: 15px;
  font-size: 14px;
}

.verification-value {
  color: var(--dark);
  flex: 1;
  font-size: 15px;
  font-weight: 500;
}

/* 状态徽章 */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 18px;
  border-radius: 25px;
  font-size: 13px;
  font-weight: 600;
  transition: var(--transition);
}

.status-badge.success {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.15) 0%, rgba(76, 175, 80, 0.1) 100%);
  color: var(--success);
  border: 1px solid rgba(76, 175, 80, 0.2);
}

.status-badge.warning {
  background: linear-gradient(135deg, rgba(255, 152, 0, 0.15) 0%, rgba(255, 152, 0, 0.1) 100%);
  color: var(--warning);
  border: 1px solid rgba(255, 152, 0, 0.2);
}

.status-badge.error {
  background: linear-gradient(135deg, rgba(244, 67, 54, 0.15) 0%, rgba(244, 67, 54, 0.1) 100%);
  color: var(--error);
  border: 1px solid rgba(244, 67, 54, 0.2);
}

.status-badge.unknown {
  background: linear-gradient(135deg, rgba(96, 125, 139, 0.15) 0%, rgba(96, 125, 139, 0.1) 100%);
  color: var(--gray);
  border: 1px solid rgba(96, 125, 139, 0.2);
}

.status-badge::before {
  content: '';
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}

.status-badge.success::before {
  background-color: var(--success);
  box-shadow: 0 0 0 3px rgba(76, 175, 80, 0.2);
}

.status-badge.warning::before {
  background-color: var(--warning);
  box-shadow: 0 0 0 3px rgba(255, 152, 0, 0.2);
}

.status-badge.error::before {
  background-color: var(--error);
  box-shadow: 0 0 0 3px rgba(244, 67, 54, 0.2);
}

.status-badge.unknown::before {
  background-color: var(--gray);
  box-shadow: 0 0 0 3px rgba(96, 125, 139, 0.2);
}

/* 空状态 */
.empty-card {
  margin-top: 30px;
}

.empty-state {
  text-align: center;
  padding: 60px 30px;
  animation: fadeIn 0.8s ease;
}

.empty-state i {
  font-size: 72px;
  color: #cfd8dc;
  margin-bottom: 25px;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.empty-state h3 {
  font-size: 22px;
  margin-bottom: 15px;
  color: var(--dark);
  font-weight: 600;
}

.empty-state p {
  color: var(--gray);
  margin-bottom: 25px;
  font-size: 16px;
  max-width: 400px;
  margin-left: auto;
  margin-right: auto;
}

/* 加载动画 */
.spinner {
  display: inline-block;
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s ease-in-out infinite;
  margin-left: 10px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 滚动条美化 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: var(--light-gray);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: var(--primary);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--primary-dark);
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .data-query {
    padding: 15px;
  }

  .card {
    padding: 25px;
  }

  .result-grid,
  .verification-info {
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  }
}

@media (max-width: 992px) {
  .form-row {
    gap: 20px;
  }

  .card-header {
    flex-direction: column;
    gap: 15px;
    align-items: flex-start;
  }

  .result-actions {
    width: 100%;
    justify-content: flex-start;
  }

  .data-timeline {
    padding-left: 20px;
  }

  .timeline-item::before {
    left: -20px;
  }
}

@media (max-width: 768px) {
  .data-query {
    padding: 10px;
  }

  .card {
    padding: 20px;
    border-radius: 12px;
  }

  .card-title {
    font-size: 20px;
  }

  .form-col {
    min-width: 100%;
  }

  .btn {
    padding: 12px 20px;
    font-size: 15px;
  }

  .history-tag {
    font-size: 13px;
    padding: 6px 12px;
  }

  .timeline-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }

  .stage-icon {
    margin-bottom: 10px;
    margin-right: 0;
  }

  .stage-status {
    margin-top: 10px;
  }

  .result-grid,
  .verification-info {
    grid-template-columns: 1fr;
    gap: 15px;
  }

  .empty-state {
    padding: 40px 20px;
  }

  .empty-state i {
    font-size: 56px;
  }

  .empty-state h3 {
    font-size: 20px;
  }
}

@media (max-width: 480px) {
  .card {
    padding: 15px;
  }

  .card-title {
    font-size: 18px;
  }

  .card-title i {
    font-size: 20px;
  }

  .form-control,
  .select-wrapper select {
    padding: 12px 14px;
    font-size: 14px;
  }

  .btn-outline {
    padding: 8px 16px;
    font-size: 13px;
  }

  .section-title {
    font-size: 17px;
  }

  .result-item,
  .verification-item {
    padding: 15px;
  }

  .timeline-item {
    padding: 20px;
  }
}

/* 打印样式 */
@media print {
  .data-query {
    padding: 0;
  }

  .card {
    box-shadow: none;
    border: 1px solid #ddd;
    break-inside: avoid;
  }

  .btn,
  .query-help,
  .history-tag,
  .copy-btn {
    display: none !important;
  }

  .result-content {
    padding: 0;
  }
}
</style>
