<template>
  <div class="node-management">
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">
          <i class="fas fa-network-wired"></i>
          {{ t('nodeManagement.title') }}
        </h3>
        <div class="header-actions">
          <button class="btn btn-outline" @click="showAddNodeDialog">
            <i class="fas fa-plus"></i>
            {{ t('nodeManagement.addNode') }}
          </button>
          <button class="btn btn-secondary" @click="refreshNodes">
            <i class="fas fa-sync-alt" :class="{ 'spin': refreshing }"></i>
            {{ refreshing ? t('common.loading') : t('nodeManagement.refresh') }}
          </button>
        </div>
      </div>

      <!-- 节点统计 -->
      <div class="node-stats">
        <div class="stat-item">
          <div class="stat-icon online">
            <i class="fas fa-server"></i>
          </div>
          <div class="stat-info">
            <div class="stat-label">{{ t('nodeManagement.onlineNodes') }}</div>
            <div class="stat-value">{{ nodeStats.online }}</div>
          </div>
        </div>
        <div class="stat-item">
          <div class="stat-icon total">
            <i class="fas fa-hdd"></i>
          </div>
          <div class="stat-info">
            <div class="stat-label">{{ t('nodeManagement.totalNodes') }}</div>
            <div class="stat-value">{{ nodeStats.total }}</div>
          </div>
        </div>
        <div class="stat-item">
          <div class="stat-icon verify">
            <i class="fas fa-check-circle"></i>
          </div>
          <div class="stat-info">
            <div class="stat-label">{{ t('nodeManagement.verifiedNodes') }}</div>
            <div class="stat-value">{{ nodeStats.verified }}</div>
          </div>
        </div>
        <div class="stat-item">
          <div class="stat-icon error">
            <i class="fas fa-exclamation-triangle"></i>
          </div>
          <div class="stat-info">
            <div class="stat-label">{{ t('nodeManagement.offlineNodes') }}</div>
            <div class="stat-value">{{ nodeStats.error }}</div>
          </div>
        </div>
      </div>

      <!-- 节点列表 -->
      <div class="node-list-container">
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>{{ t('nodeManagement.loadNodesFailed') }}</p>
        </div>

        <div v-else-if="nodes.length === 0" class="empty-state">
          <i class="fas fa-server"></i>
          <h3>{{ t('common.noData') }}</h3>
          <p>{{ t('nodeManagement.description') }}</p>
          <button class="btn btn-primary" @click="showAddNodeDialog">
            <i class="fas fa-plus"></i>
            {{ t('nodeManagement.addNode') }}
          </button>
        </div>

        <div v-else class="node-table-wrapper">
          <div class="table-controls">
            <div class="search-box">
              <i class="fas fa-search"></i>
              <input
                v-model="searchQuery"
                type="text"
                :placeholder="t('nodeManagement.searchNodes')"
                @input="filterNodes"
              >
            </div>
            <div class="filter-controls">
              <select v-model="statusFilter" @change="filterNodes" class="filter-select">
                <option value="all">{{ t('common.all') }}</option>
                <option value="online">{{ t('nodeManagement.online') }}</option>
                <option value="offline">{{ t('nodeManagement.offline') }}</option>
                <option value="pending">{{ t('nodeManagement.connecting') }}</option>
              </select>
              <select v-model="sortBy" @change="sortNodes" class="filter-select">
                <option value="id">{{ t('nodeManagement.nodeId') }}</option>
                <option value="verify_time">{{ t('nodeManagement.verifiedNodes') }}</option>
                <option value="create_time">{{ t('nodeManagement.lastSeen') }}</option>
              </select>
            </div>
          </div>

          <table class="node-table">
            <thead>
              <tr>
                <th>{{ t('nodeManagement.nodeId') }}</th>
                <th>{{ t('nodeManagement.nodeAddress') }}</th>
                <th>{{ t('nodeManagement.lastSeen') }}</th>
                <th>{{ t('nodeManagement.verifiedNodes') }}</th>
                <th>{{ t('nodeManagement.nodeStatus') }}</th>
                <th>{{ t('nodeManagement.operations') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="node in filteredNodes" :key="node.id">
                <td>
                  <div class="node-id">
                    <i class="fas fa-circle" :class="getNodeStatusClass(node)"></i>
                    {{ node.id }}
                  </div>
                </td>
                <td class="node-address">
                  <div class="address-content">
                    <code>{{ node.addr }}</code>
                    <button class="copy-btn" @click="copyAddress(node.addr)">
                      <i class="fas fa-copy"></i>
                    </button>
                  </div>
                </td>
                <td>{{ formatTime(node.create_time) }}</td>
                <td>
                  <div class="verify-count">
                    <span class="count-badge">{{ node.verify_time || 0 }}</span>
                  </div>
                </td>
                <td>
                  <span class="status-badge" :class="getNodeStatusClass(node)">
                    {{ getNodeStatusText(node) }}
                  </span>
                </td>
                <td>
                  <div class="response-time">
                    {{ getResponseTime(node) }}
                    <span class="time-unit">ms</span>
                  </div>
                </td>
                <td>
                  <div class="action-buttons">
                    <button class="action-btn" @click="pingNode(node)" title="测试连接">
                      <i class="fas fa-wifi"></i>
                    </button>
                    <button class="action-btn" @click="showNodeDetail(node)" title="查看详情">
                      <i class="fas fa-eye"></i>
                    </button>
                    <button class="action-btn" @click="editNode(node)" title="编辑节点">
                      <i class="fas fa-edit"></i>
                    </button>
                    <button class="action-btn danger" @click="deleteNode(node)" title="删除节点">
                      <i class="fas fa-trash"></i>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>

          <!-- 分页控件 -->
          <div v-if="filteredNodes.length > 0" class="pagination">
            <div class="pagination-info">
              {{ t('nodeManagement.paginationInfo', { start: pagination.start + 1, end: Math.min(pagination.end, filteredNodes.length), total: filteredNodes.length }) }}
            </div>
            <div class="pagination-controls">
              <button
                class="pagination-btn"
                :disabled="pagination.currentPage === 1"
                @click="prevPage"
              >
                <i class="fas fa-chevron-left"></i>
              </button>
              <div class="page-numbers">
                <button
                  v-for="page in pagination.totalPages"
                  :key="page"
                  class="page-btn"
                  :class="{ active: page === pagination.currentPage }"
                  @click="goToPage(page)"
                >
                  {{ page }}
                </button>
              </div>
              <button
                class="pagination-btn"
                :disabled="pagination.currentPage === pagination.totalPages"
                @click="nextPage"
              >
                <i class="fas fa-chevron-right"></i>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 节点详情对话框 -->
    <div v-if="showDetailDialog" class="modal-overlay">
      <div class="modal-dialog">
        <div class="modal-header">
          <h3>{{ t('nodeManagement.nodeDetails') }}</h3>
          <button class="close-btn" @click="showDetailDialog = false">
            <i class="fas fa-times"></i>
          </button>
        </div>
        <div class="modal-content">
          <div v-if="selectedNode" class="node-detail">
            <div class="detail-section">
              <h4 class="section-title">
                <i class="fas fa-info-circle"></i>
                {{ t('nodeManagement.nodeInfo') }}
              </h4>
              <div class="detail-grid">
                <div class="detail-item">
                  <div class="detail-label">{{ t('nodeManagement.nodeId') }}:</div>
                  <div class="detail-value">{{ selectedNode.id }}</div>
                </div>
                <div class="detail-item">
                  <div class="detail-label">{{ t('nodeManagement.nodeAddress') }}:</div>
                  <div class="detail-value">
                    <code>{{ selectedNode.addr }}</code>
                  </div>
                </div>
                <div class="detail-item">
                  <div class="detail-label">{{ t('nodeManagement.lastSeen') }}:</div>
                  <div class="detail-value">{{ formatTime(selectedNode.create_time) }}</div>
                </div>
                <div class="detail-item">
                  <div class="detail-label">{{ t('nodeManagement.connectionStatus') }}:</div>
                  <div class="detail-value">{{ selectedNode.last_active || t('common.noData') }}</div>
                </div>
              </div>
            </div>

            <div class="detail-section">
              <h4 class="section-title">
                <i class="fas fa-chart-bar"></i>
                {{ t('nodeManagement.networkInfo') }}
              </h4>
              <div class="performance-stats">
                <div class="performance-item">
                  <div class="performance-label">{{ t('nodeManagement.verifiedNodes') }}:</div>
                  <div class="performance-value">{{ selectedNode.verify_time || 0 }}</div>
                </div>
                <div class="performance-item">
                  <div class="performance-label">{{ t('nodeManagement.operations') }}:</div>
                  <div class="performance-value">
                    {{ getResponseTime(selectedNode) }}
                    <span class="unit">ms</span>
                  </div>
                </div>
                <div class="performance-item">
                  <div class="performance-label">{{ t('nodeManagement.nodeStatus') }}:</div>
                  <div class="performance-value">
                    <span class="availability-badge" :class="getAvailabilityClass(selectedNode)">
                      {{ getAvailabilityText(selectedNode) }}
                    </span>
                  </div>
                </div>
                <div class="performance-item">
                  <div class="performance-label">{{ t('nodeManagement.version') }}:</div>
                  <div class="performance-value">{{ selectedNode.version || '1.0.0' }}</div>
                </div>
              </div>
            </div>

            <div class="detail-section">
              <h4 class="section-title">
                <i class="fas fa-cogs"></i>
                {{ t('nodeManagement.nodeInfo') }}
              </h4>
              <div class="config-info">
                <div class="config-item">
                  <div class="config-label">{{ t('nodeManagement.blockHeight') }}:</div>
                  <div class="config-value">{{ selectedNode.block_height || 'N/A' }}</div>
                </div>
                <div class="config-item">
                  <div class="config-label">{{ t('nodeManagement.connectionStatus') }}:</div>
                  <div class="config-value">
                    <span class="sync-status" :class="getSyncClass(selectedNode)">
                      {{ getSyncText(selectedNode) }}
                    </span>
                  </div>
                </div>
                <div class="config-item">
                  <div class="config-label">{{ t('nodeManagement.peers') }}:</div>
                  <div class="config-value">{{ selectedNode.connections || 0 }}</div>
                </div>
                <div class="config-item">
                  <div class="config-label">{{ t('nodeManagement.diskUsage') }}:</div>
                  <div class="config-value">
                    <div class="storage-bar">
                      <div class="storage-fill" :style="{ width: getStorageUsage(selectedNode) + '%' }"></div>
                    </div>
                    <span class="storage-percent">{{ getStorageUsage(selectedNode) }}%</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="showDetailDialog = false">
            {{ t('common.cancel') }}
          </button>
          <button class="btn btn-primary" @click="testNodeConnection(selectedNode)">
            <i class="fas fa-wifi"></i>
            {{ t('nodeManagement.testConnection') }}
          </button>
        </div>
      </div>
    </div>

    <!-- 添加节点对话框 -->
    <div v-if="showAddDialog" class="modal-overlay">
      <div class="modal-dialog">
        <div class="modal-header">
          <h3>{{ t('nodeManagement.addNewNode') }}</h3>
          <button class="close-btn" @click="closeAddDialog">
            <i class="fas fa-times"></i>
          </button>
        </div>
        <div class="modal-content">
          <div class="add-node-form">
            <div class="form-group">
              <label class="form-label">{{ t('nodeManagement.nodeId') }} *</label>
              <input
                v-model="newNode.id"
                type="text"
                class="form-control"
                :placeholder="t('nodeManagement.enterNodeId')"
              >
            </div>
            <div class="form-group">
              <label class="form-label">{{ t('nodeManagement.nodeAddress') }} *</label>
              <input
                v-model="newNode.addr"
                type="text"
                class="form-control"
                :placeholder="t('nodeManagement.nodeAddressFormat')"
              >
            </div>
            <div class="form-group">
              <label class="form-label">节点类型</label>
              <select v-model="newNode.type" class="form-control">
                <option value="full">全节点</option>
                <option value="light">轻节点</option>
                <option value="validator">验证节点</option>
              </select>
            </div>
            <div class="form-group">
              <label class="form-label">{{ t('nodeManagement.nodeDescriptionLabel') }}</label>
              <textarea
                v-model="newNode.description"
                class="form-control"
                rows="3"
                :placeholder="t('nodeManagement.nodeDescription')"
              ></textarea>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="closeAddDialog">
            {{ t('nodeManagement.cancel') }}
          </button>
          <button class="btn btn-primary" @click="addNewNode" :disabled="!isValidNewNode">
            <i class="fas fa-plus"></i>
            {{ t('nodeManagement.confirmAdd') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '../services/api'

const { t } = useI18n()

const nodes = ref([])
const filteredNodes = ref([])
const loading = ref(false)
const refreshing = ref(false)
const showDetailDialog = ref(false)
const showAddDialog = ref(false)
const selectedNode = ref(null)
const searchQuery = ref('')
const statusFilter = ref('all')
const sortBy = ref('id')

// 分页
const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  totalPages: 1,
  start: 0,
  end: 10
})

// 新节点表单
const newNode = reactive({
  id: '',
  addr: '',
  type: 'full',
  description: ''
})

// 节点统计
const nodeStats = reactive({
  total: 0,
  online: 0,
  offline: 0,
  verified: 0,
  error: 0
})

// 计算属性
const isValidNewNode = computed(() => {
  return newNode.id.trim() !== '' && newNode.addr.trim() !== ''
})

// 方法
async function loadNodes() {
  loading.value = true
  try {
    const data = await api.getNodes()
    if (Array.isArray(data)) {
      nodes.value = data.map(node => ({
        ...node,
        // 模拟一些额外数据用于演示
        status: Math.random() > 0.3 ? 'online' : 'offline',
        response_time: Math.floor(Math.random() * 200) + 50,
        last_active: new Date(Date.now() - Math.random() * 86400000).toISOString(),
        version: '1.0.' + Math.floor(Math.random() * 10),
        block_height: Math.floor(Math.random() * 10000) + 1000,
        connections: Math.floor(Math.random() * 50) + 5,
        storage_usage: Math.floor(Math.random() * 80) + 10
      }))

      updateNodeStats()
      filterNodes()
    }
  } catch (error) {
    console.error('加载节点失败:', error)
    alert('加载节点列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

function updateNodeStats() {
  nodeStats.total = nodes.value.length
  nodeStats.online = nodes.value.filter(n => n.status === 'online').length
  nodeStats.offline = nodes.value.filter(n => n.status === 'offline').length
  nodeStats.verified = nodes.value.filter(n => (n.verify_time || 0) > 0).length
  nodeStats.error = nodes.value.filter(n => n.status === 'error').length
}

function filterNodes() {
  let filtered = [...nodes.value]

  // 搜索过滤
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(node =>
      node.id.toLowerCase().includes(query) ||
      node.addr.toLowerCase().includes(query)
    )
  }

  // 状态过滤
  if (statusFilter.value !== 'all') {
    filtered = filtered.filter(node => node.status === statusFilter.value)
  }

  // 排序
  sortNodes(filtered)

  // 更新分页
  filteredNodes.value = filtered
  updatePagination()
}

function sortNodes(list = filteredNodes.value) {
  if (sortBy.value === 'id') {
    list.sort((a, b) => a.id.localeCompare(b.id))
  } else if (sortBy.value === 'verify_time') {
    list.sort((a, b) => (b.verify_time || 0) - (a.verify_time || 0))
  } else if (sortBy.value === 'create_time') {
    list.sort((a, b) => new Date(b.create_time) - new Date(a.create_time))
  }
}

function updatePagination() {
  pagination.totalPages = Math.ceil(filteredNodes.value.length / pagination.pageSize)
  pagination.currentPage = Math.min(pagination.currentPage, pagination.totalPages)
  pagination.start = (pagination.currentPage - 1) * pagination.pageSize
  pagination.end = pagination.start + pagination.pageSize

  // 更新当前页数据
  filteredNodes.value = filteredNodes.value
}

function goToPage(page) {
  if (page >= 1 && page <= pagination.totalPages) {
    pagination.currentPage = page
    pagination.start = (page - 1) * pagination.pageSize
    pagination.end = page * pagination.pageSize
  }
}

function prevPage() {
  if (pagination.currentPage > 1) {
    goToPage(pagination.currentPage - 1)
  }
}

function nextPage() {
  if (pagination.currentPage < pagination.totalPages) {
    goToPage(pagination.currentPage + 1)
  }
}

function getNodeStatusClass(node) {
  return node.status || 'unknown'
}

function getNodeStatusText(node) {
  const status = node.status || 'unknown'
  const statusMap = {
    'online': '在线',
    'offline': '离线',
    'pending': '待验证',
    'error': '异常',
    'unknown': '未知'
  }
  return statusMap[status] || '未知'
}

function getResponseTime(node) {
  return node.response_time || Math.floor(Math.random() * 200) + 50
}

function getAvailabilityClass(node) {
  const status = node.status || 'unknown'
  if (status === 'online') return 'high'
  if (status === 'offline') return 'low'
  return 'medium'
}

function getAvailabilityText(node) {
  const status = node.status || 'unknown'
  if (status === 'online') return '高'
  if (status === 'offline') return '低'
  return '中'
}

function getSyncClass(node) {
  const height = node.block_height || 0
  const avgHeight = Math.max(...nodes.value.map(n => n.block_height || 0))
  const diff = avgHeight - height
  if (diff === 0) return 'synced'
  if (diff <= 10) return 'syncing'
  return 'outdated'
}

function getSyncText(node) {
  const cls = getSyncClass(node)
  const textMap = {
    'synced': '已同步',
    'syncing': '同步中',
    'outdated': '落后'
  }
  return textMap[cls] || '未知'
}

function getStorageUsage(node) {
  return node.storage_usage || Math.floor(Math.random() * 80) + 10
}

function formatTime(timeString) {
  if (!timeString) return '未知'
  try {
    return new Date(timeString).toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return timeString
  }
}

function copyAddress(address) {
  navigator.clipboard.writeText(address)
    .then(() => {
      alert('地址已复制到剪贴板')
    })
    .catch(err => {
      console.error('复制失败:', err)
    })
}

function showNodeDetail(node) {
  selectedNode.value = node
  showDetailDialog.value = true
}

async function pingNode(node) {
  try {
    // 显示正在测试的提示
    node.status = 'pending'

    const result = await api.pingNode(node.addr)

    if (result.online) {
      // 更新节点状态为在线
      node.status = 'online'
      node.response_time = result.response_time
      alert(`✅ 节点 ${node.id} 连接成功！\n响应时间: ${result.response_time}ms`)
    } else {
      // 更新节点状态为离线
      node.status = 'offline'
      node.response_time = -1
      alert(`❌ 节点 ${node.id} 连接失败\n${result.msg}`)
    }

    // 更新统计信息
    updateNodeStats()
  } catch (error) {
    console.error('测试节点连接失败:', error)
    node.status = 'offline'
    node.response_time = -1
    alert(`❌ 测试节点 ${node.id} 连接失败\n${error.message || '网络错误'}`)
    updateNodeStats()
  }
}

async function testNodeConnection(node) {
  await pingNode(node)
}

function editNode(node) {
  selectedNode.value = node
  alert(`编辑节点 ${node.id}\n实际实现中会打开编辑表单`)
}

function deleteNode(node) {
  if (confirm(`确定要删除节点 ${node.id} 吗？`)) {
    // 实际实现应该调用API删除节点
    nodes.value = nodes.value.filter(n => n.id !== node.id)
    updateNodeStats()
    filterNodes()
    alert(`节点 ${node.id} 已删除`)
  }
}

function showAddNodeDialog() {
  showAddDialog.value = true
}

function closeAddDialog() {
  showAddDialog.value = false
  // 重置表单
  Object.assign(newNode, {
    id: '',
    addr: '',
    type: 'full',
    description: ''
  })
}

function addNewNode() {
  if (!isValidNewNode.value) {
    alert('请填写完整的节点信息')
    return
  }

  // 检查ID是否重复
  if (nodes.value.some(node => node.id === newNode.id)) {
    alert('节点ID已存在，请使用其他ID')
    return
  }

  // 创建新节点对象
  const newNodeObj = {
    id: newNode.id,
    addr: newNode.addr,
    create_time: new Date().toISOString(),
    verify_time: 0,
    status: 'pending',
    type: newNode.type,
    description: newNode.description,
    response_time: Math.floor(Math.random() * 200) + 50,
    last_active: new Date().toISOString(),
    version: '1.0.0',
    block_height: Math.floor(Math.random() * 10000) + 1000,
    connections: 0,
    storage_usage: 0
  }

  // 添加到列表
  nodes.value.push(newNodeObj)
  updateNodeStats()
  filterNodes()

  // 关闭对话框
  closeAddDialog()

  alert(`节点 ${newNodeObj.id} 添加成功`)
}

async function refreshNodes() {
  refreshing.value = true
  try {
    await loadNodes()
  } finally {
    refreshing.value = false
  }
}

onMounted(() => {
  loadNodes()
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

.node-management {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
  min-height: 100vh;
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

/* 卡片样式 */
.card {
  background: white;
  border-radius: 20px;
  padding: 35px;
  margin-bottom: 35px;
  box-shadow: var(--card-shadow);
  transition: var(--transition);
  border: 1px solid var(--light-gray);
  position: relative;
  overflow: hidden;
}

.card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--primary) 0%, var(--primary-light) 100%);
}

.card:hover {
  box-shadow: var(--card-shadow-hover);
  transform: translateY(-2px);
}

/* 卡片头部 */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 40px;
  padding-bottom: 25px;
  border-bottom: 2px solid var(--light-gray);
  flex-wrap: wrap;
  gap: 20px;
}

.card-title {
  font-size: 26px;
  font-weight: 700;
  color: var(--primary-dark);
  display: flex;
  align-items: center;
  gap: 15px;
  margin: 0;
  position: relative;
  padding-left: 15px;
}

.card-title::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 6px;
  height: 24px;
  background: linear-gradient(180deg, var(--primary) 0%, var(--primary-light) 100%);
  border-radius: 3px;
}

.card-title i {
  color: var(--primary);
  font-size: 28px;
  background: linear-gradient(135deg, rgba(46, 125, 50, 0.1) 0%, rgba(46, 125, 50, 0.05) 100%);
  padding: 12px;
  border-radius: 12px;
}

.header-actions {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

/* 按钮样式 */
.btn {
  padding: 14px 28px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  position: relative;
  overflow: hidden;
}

.btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.6s;
}

.btn:hover::before {
  left: 100%;
}

.btn-primary {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
  color: white;
  box-shadow: 0 6px 20px rgba(46, 125, 50, 0.3);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-3px);
  box-shadow: 0 10px 25px rgba(46, 125, 50, 0.4);
}

.btn-secondary {
  background: linear-gradient(135deg, var(--secondary) 0%, var(--secondary-dark) 100%);
  color: white;
  box-shadow: 0 6px 20px rgba(2, 136, 209, 0.3);
}

.btn-secondary:hover:not(:disabled) {
  transform: translateY(-3px);
  box-shadow: 0 10px 25px rgba(2, 136, 209, 0.4);
}

.btn-outline {
  background: transparent;
  border: 2px solid var(--primary);
  color: var(--primary);
  box-shadow: 0 4px 15px rgba(46, 125, 50, 0.1);
}

.btn-outline:hover {
  background: var(--primary);
  color: white;
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(46, 125, 50, 0.2);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
  box-shadow: none !important;
}

/* 节点统计卡片 */
.node-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 25px;
  margin-bottom: 40px;
  animation: slideIn 0.6s ease;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.stat-item {
  display: flex;
  align-items: center;
  padding: 25px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 16px;
  transition: var(--transition);
  border: 2px solid transparent;
  position: relative;
  overflow: hidden;
}

.stat-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--primary) 0%, var(--primary-light) 100%);
}

.stat-item:hover {
  transform: translateY(-5px);
  border-color: var(--primary);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  width: 70px;
  height: 70px;
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 25px;
  font-size: 32px;
  box-shadow: 0 6px 15px rgba(0, 0, 0, 0.1);
  transition: var(--transition);
}

.stat-item:hover .stat-icon {
  transform: scale(1.1);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.15);
}

.stat-icon.online {
  background: linear-gradient(135deg, var(--success) 0%, #66bb6a 100%);
  color: white;
}

.stat-icon.total {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
  color: white;
}

.stat-icon.verify {
  background: linear-gradient(135deg, var(--warning) 0%, #ffb74d 100%);
  color: white;
}

.stat-icon.error {
  background: linear-gradient(135deg, var(--error) 0%, #ef5350 100%);
  color: white;
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 15px;
  color: var(--gray);
  margin-bottom: 8px;
  font-weight: 500;
}

.stat-value {
  font-size: 36px;
  font-weight: 800;
  color: var(--dark);
  line-height: 1;
  font-family: 'Segoe UI', sans-serif;
}

/* 表格控件 */
.table-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  flex-wrap: wrap;
  gap: 20px;
  padding: 20px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 15px;
  border: 1px solid var(--border-color);
}

.search-box {
  display: flex;
  align-items: center;
  background: white;
  border: 2px solid var(--light-gray);
  border-radius: 12px;
  padding: 12px 20px;
  flex: 1;
  min-width: 300px;
  transition: var(--transition);
}

.search-box:focus-within {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(46, 125, 50, 0.1);
  transform: translateY(-2px);
}

.search-box i {
  color: var(--gray);
  margin-right: 15px;
  font-size: 16px;
}

.search-box input {
  border: none;
  outline: none;
  flex: 1;
  font-size: 15px;
  background: transparent;
  color: var(--dark);
}

.search-box input::placeholder {
  color: var(--gray);
  opacity: 0.7;
}

.filter-controls {
  display: flex;
  gap: 15px;
}

.filter-select {
  padding: 12px 20px;
  border: 2px solid var(--light-gray);
  border-radius: 12px;
  font-size: 15px;
  background: white;
  color: var(--dark);
  cursor: pointer;
  transition: var(--transition);
  min-width: 160px;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%2378909c' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 20px center;
  background-size: 16px;
  padding-right: 50px;
}

.filter-select:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(46, 125, 50, 0.1);
  transform: translateY(-2px);
}

/* 节点表格 */
.node-table-wrapper {
  overflow-x: auto;
  border-radius: 15px;
  border: 1px solid var(--border-color);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

.node-table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  min-width: 1100px;
  background: white;
}

.node-table thead {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
}

.node-table th {
  padding: 22px 20px;
  text-align: left;
  font-weight: 600;
  font-size: 15px;
  color: white;
  border: none;
  white-space: nowrap;
  position: relative;
  overflow: hidden;
}

.node-table th::after {
  content: '';
  position: absolute;
  right: 0;
  top: 25%;
  height: 50%;
  width: 1px;
  background: rgba(255, 255, 255, 0.2);
}

.node-table th:last-child::after {
  display: none;
}

.node-table tbody tr {
  transition: var(--transition);
  border-bottom: 1px solid var(--light-gray);
  position: relative;
}

.node-table tbody tr::after {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  background: transparent;
  transition: var(--transition);
}

.node-table tbody tr:hover {
  background: linear-gradient(135deg, rgba(46, 125, 50, 0.03) 0%, rgba(46, 125, 50, 0.01) 100%);
}

.node-table tbody tr:hover::after {
  background: linear-gradient(180deg, var(--primary) 0%, var(--primary-light) 100%);
}

.node-table tbody tr:last-child {
  border-bottom: none;
}

.node-table td {
  padding: 22px 20px;
  color: var(--dark);
  border-bottom: 1px solid var(--light-gray);
  transition: var(--transition);
}

.node-table tbody tr:hover td {
  transform: translateX(5px);
}

.node-table tbody tr:last-child td {
  border-bottom: none;
}

/* 节点ID样式 */
.node-id {
  display: flex;
  align-items: center;
  gap: 15px;
  font-weight: 600;
  color: var(--dark);
  font-family: 'Courier New', monospace;
}

.node-id i {
  font-size: 10px;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.node-id i.online {
  color: var(--success);
  text-shadow: 0 0 8px rgba(76, 175, 80, 0.4);
}

.node-id i.offline {
  color: var(--error);
}

.node-id i.pending {
  color: var(--warning);
}

/* 节点地址 */
.node-address {
  max-width: 350px;
}

.address-content {
  display: flex;
  align-items: center;
  gap: 15px;
  background: linear-gradient(135deg, #f5f5f5 0%, #eeeeee 100%);
  padding: 10px 15px;
  border-radius: 10px;
  border: 1px solid var(--border-color);
}

.address-content code {
  flex: 1;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  word-break: break-all;
  color: var(--dark);
  background: transparent;
  padding: 0;
}

.copy-btn {
  background: transparent;
  border: none;
  color: var(--primary);
  cursor: pointer;
  padding: 6px;
  font-size: 15px;
  border-radius: 8px;
  transition: var(--transition);
  display: flex;
  align-items: center;
  justify-content: center;
}

.copy-btn:hover {
  background: var(--primary);
  color: white;
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(46, 125, 50, 0.2);
}

/* 验证数徽章 */
.verify-count {
  display: flex;
  align-items: center;
}

.count-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 6px 14px;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
  color: white;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 700;
  min-width: 40px;
  height: 40px;
  box-shadow: 0 4px 10px rgba(46, 125, 50, 0.2);
  transition: var(--transition);
}

.count-badge:hover {
  transform: scale(1.05);
  box-shadow: 0 6px 15px rgba(46, 125, 50, 0.3);
}

/* 状态徽章 */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 8px 18px;
  border-radius: 25px;
  font-size: 13px;
  font-weight: 700;
  transition: var(--transition);
  position: relative;
  overflow: hidden;
}

.status-badge::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, transparent 100%);
  opacity: 0;
  transition: var(--transition);
}

.status-badge:hover::before {
  opacity: 1;
}

.status-badge.online {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.15) 0%, rgba(76, 175, 80, 0.1) 100%);
  color: var(--success);
  border: 1px solid rgba(76, 175, 80, 0.2);
  box-shadow: 0 4px 12px rgba(76, 175, 80, 0.1);
}

.status-badge.offline {
  background: linear-gradient(135deg, rgba(244, 67, 54, 0.15) 0%, rgba(244, 67, 54, 0.1) 100%);
  color: var(--error);
  border: 1px solid rgba(244, 67, 54, 0.2);
  box-shadow: 0 4px 12px rgba(244, 67, 54, 0.1);
}

.status-badge.pending {
  background: linear-gradient(135deg, rgba(255, 152, 0, 0.15) 0%, rgba(255, 152, 0, 0.1) 100%);
  color: var(--warning);
  border: 1px solid rgba(255, 152, 0, 0.2);
  box-shadow: 0 4px 12px rgba(255, 152, 0, 0.1);
}

/* 响应时间 */
.response-time {
  font-family: 'Courier New', monospace;
  font-weight: 700;
  color: var(--dark);
  font-size: 16px;
  display: flex;
  align-items: baseline;
  gap: 5px;
}

.time-unit {
  font-size: 13px;
  color: var(--gray);
  font-weight: normal;
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  gap: 10px;
}

.action-btn {
  width: 40px;
  height: 40px;
  border: none;
  border-radius: 10px;
  background: linear-gradient(135deg, var(--light-gray) 0%, #e0e0e0 100%);
  color: var(--dark);
  cursor: pointer;
  transition: var(--transition);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  position: relative;
  overflow: hidden;
}

.action-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.2) 0%, transparent 100%);
  opacity: 0;
  transition: var(--transition);
}

.action-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 15px rgba(0, 0, 0, 0.1);
}

.action-btn:hover::before {
  opacity: 1;
}

.action-btn:hover {
  background: var(--primary);
  color: white;
}

.action-btn.danger:hover {
  background: var(--error);
}

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 30px;
  min-height: 400px;
  animation: fadeIn 0.5s ease;
}

.loading-state .spinner {
  width: 50px;
  height: 50px;
  border: 3px solid var(--light-gray);
  border-top-color: var(--primary);
  border-radius: 50%;
  margin-bottom: 25px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-state p {
  color: var(--gray);
  font-size: 16px;
  font-weight: 500;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 80px 30px;
  min-height: 400px;
  animation: fadeIn 0.6s ease;
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
  font-size: 24px;
  margin-bottom: 15px;
  color: var(--dark);
  font-weight: 700;
}

.empty-state p {
  color: var(--gray);
  margin-bottom: 30px;
  font-size: 16px;
  max-width: 400px;
  margin-left: auto;
  margin-right: auto;
  line-height: 1.6;
}

/* 分页控件 */
.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 30px;
  padding-top: 30px;
  border-top: 2px solid var(--light-gray);
  flex-wrap: wrap;
  gap: 20px;
}

.pagination-info {
  color: var(--gray);
  font-size: 15px;
  font-weight: 500;
  background: var(--light);
  padding: 10px 20px;
  border-radius: 10px;
  border: 1px solid var(--border-color);
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 15px;
}

.pagination-btn {
  width: 44px;
  height: 44px;
  border: 2px solid var(--light-gray);
  border-radius: 12px;
  background: white;
  color: var(--dark);
  cursor: pointer;
  transition: var(--transition);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}

.pagination-btn:not(:disabled):hover {
  border-color: var(--primary);
  color: var(--primary);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(46, 125, 50, 0.1);
}

.pagination-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  background: var(--light-gray);
}

.page-numbers {
  display: flex;
  gap: 8px;
}

.page-btn {
  min-width: 44px;
  height: 44px;
  border: 2px solid var(--light-gray);
  border-radius: 12px;
  background: white;
  color: var(--dark);
  cursor: pointer;
  transition: var(--transition);
  font-size: 15px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
}

.page-btn:hover {
  border-color: var(--primary);
  color: var(--primary);
  transform: translateY(-2px);
}

.page-btn.active {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
  border-color: var(--primary);
  color: white;
  box-shadow: 0 4px 15px rgba(46, 125, 50, 0.2);
}

/* 模态对话框 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
  padding: 20px;
  animation: fadeIn 0.3s ease;
  backdrop-filter: blur(5px);
}

.modal-dialog {
  background: white;
  border-radius: 20px;
  max-width: 850px;
  width: 100%;
  max-height: 85vh;
  overflow-y: auto;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.15);
  animation: slideUp 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 25px 35px;
  border-bottom: 2px solid var(--light-gray);
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 20px 20px 0 0;
}

.modal-header h3 {
  font-size: 22px;
  color: var(--primary-dark);
  font-weight: 700;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.modal-header h3 i {
  color: var(--primary);
}

.close-btn {
  background: transparent;
  border: none;
  font-size: 22px;
  color: var(--gray);
  cursor: pointer;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: var(--transition);
}

.close-btn:hover {
  background: rgba(244, 67, 54, 0.1);
  color: var(--error);
  transform: rotate(90deg);
}

.modal-content {
  padding: 35px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 20px;
  padding: 25px 35px;
  border-top: 2px solid var(--light-gray);
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 0 0 20px 20px;
}

/* 节点详情样式 */
.node-detail {
  display: flex;
  flex-direction: column;
  gap: 35px;
}

.detail-section {
  padding-bottom: 30px;
  border-bottom: 1px solid var(--light-gray);
  animation: fadeInUp 0.5s ease;
  animation-fill-mode: both;
}

.detail-section:nth-child(1) { animation-delay: 0.1s; }
.detail-section:nth-child(2) { animation-delay: 0.2s; }
.detail-section:nth-child(3) { animation-delay: 0.3s; }

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

.detail-section:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.section-title {
  font-size: 19px;
  margin-bottom: 25px;
  color: var(--primary-dark);
  display: flex;
  align-items: center;
  gap: 12px;
  font-weight: 700;
  padding-bottom: 12px;
  border-bottom: 2px solid var(--light-gray);
}

.section-title i {
  color: var(--primary);
  background: linear-gradient(135deg, rgba(46, 125, 50, 0.1) 0%, rgba(46, 125, 50, 0.05) 100%);
  padding: 10px;
  border-radius: 10px;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 20px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 12px;
  border: 1px solid transparent;
  transition: var(--transition);
}

.detail-item:hover {
  border-color: var(--primary);
  transform: translateY(-2px);
  box-shadow: 0 6px 15px rgba(0, 0, 0, 0.05);
}

.detail-label {
  font-weight: 600;
  color: var(--gray);
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.detail-value {
  color: var(--dark);
  word-break: break-all;
  font-size: 16px;
  font-weight: 500;
  font-family: 'Courier New', monospace;
}

/* 性能统计 */
.performance-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 25px;
}

.performance-item {
  padding: 25px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 15px;
  border: 1px solid transparent;
  transition: var(--transition);
  text-align: center;
}

.performance-item:hover {
  border-color: var(--primary);
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.08);
}

.performance-label {
  font-size: 15px;
  color: var(--gray);
  margin-bottom: 15px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.performance-value {
  font-size: 32px;
  font-weight: 800;
  color: var(--dark);
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 8px;
  line-height: 1;
}

.performance-value .unit {
  font-size: 16px;
  color: var(--gray);
  font-weight: normal;
}

.availability-badge {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 10px 20px;
  border-radius: 25px;
  font-size: 14px;
  font-weight: 700;
  transition: var(--transition);
}

.availability-badge.high {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.15) 0%, rgba(76, 175, 80, 0.1) 100%);
  color: var(--success);
  border: 1px solid rgba(76, 175, 80, 0.2);
}

.availability-badge.medium {
  background: linear-gradient(135deg, rgba(255, 152, 0, 0.15) 0%, rgba(255, 152, 0, 0.1) 100%);
  color: var(--warning);
  border: 1px solid rgba(255, 152, 0, 0.2);
}

.availability-badge.low {
  background: linear-gradient(135deg, rgba(244, 67, 54, 0.15) 0%, rgba(244, 67, 54, 0.1) 100%);
  color: var(--error);
  border: 1px solid rgba(244, 67, 54, 0.2);
}

/* 配置信息 */
.config-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.config-item {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 20px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 12px;
  border: 1px solid transparent;
  transition: var(--transition);
}

.config-item:hover {
  border-color: var(--primary);
  transform: translateY(-2px);
}

.config-label {
  font-weight: 600;
  color: var(--gray);
  font-size: 15px;
}

.config-value {
  color: var(--dark);
  display: flex;
  align-items: center;
  gap: 15px;
  flex-wrap: wrap;
}

.sync-status {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 18px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 700;
  transition: var(--transition);
}

.sync-status.synced {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.15) 0%, rgba(76, 175, 80, 0.1) 100%);
  color: var(--success);
  border: 1px solid rgba(76, 175, 80, 0.2);
}

.sync-status.syncing {
  background: linear-gradient(135deg, rgba(255, 152, 0, 0.15) 0%, rgba(255, 152, 0, 0.1) 100%);
  color: var(--warning);
  border: 1px solid rgba(255, 152, 0, 0.2);
}

.sync-status.outdated {
  background: linear-gradient(135deg, rgba(244, 67, 54, 0.15) 0%, rgba(244, 67, 54, 0.1) 100%);
  color: var(--error);
  border: 1px solid rgba(244, 67, 54, 0.2);
}

.storage-bar {
  flex: 1;
  height: 10px;
  background: var(--light-gray);
  border-radius: 5px;
  overflow: hidden;
  position: relative;
}

.storage-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--primary) 0%, var(--primary-light) 100%);
  border-radius: 5px;
  transition: width 0.8s cubic-bezier(0.34, 1.56, 0.64, 1);
  position: relative;
  overflow: hidden;
}

.storage-fill::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.4), transparent);
  animation: shimmer 2s infinite;
}

@keyframes shimmer {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}

.storage-percent {
  min-width: 50px;
  text-align: right;
  font-weight: 700;
  color: var(--primary-dark);
  font-size: 16px;
}

/* 添加节点表单 */
.add-node-form {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

.add-node-form .form-group {
  margin-bottom: 0;
  animation: fadeInRight 0.5s ease;
  animation-fill-mode: both;
}

.add-node-form .form-group:nth-child(1) { animation-delay: 0.1s; }
.add-node-form .form-group:nth-child(2) { animation-delay: 0.2s; }
.add-node-form .form-group:nth-child(3) { animation-delay: 0.3s; }
.add-node-form .form-group:nth-child(4) { animation-delay: 0.4s; }

@keyframes fadeInRight {
  from {
    opacity: 0;
    transform: translateX(20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.add-node-form .form-label {
  font-weight: 700;
  margin-bottom: 12px;
  display: block;
  color: var(--dark);
  font-size: 16px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.add-node-form .form-control {
  width: 100%;
  padding: 16px 20px;
  border: 2px solid var(--light-gray);
  border-radius: 12px;
  font-size: 16px;
  transition: var(--transition);
  background: white;
  color: var(--dark);
}

.add-node-form .form-control:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 4px rgba(46, 125, 50, 0.1);
  transform: translateY(-2px);
}

.add-node-form textarea {
  resize: vertical;
  min-height: 100px;
  font-family: inherit;
}

/* 旋转动画 */
.spin {
  animation: spin 1s linear infinite;
}

/* 响应式设计 */
@media (max-width: 1400px) {
  .node-management {
    max-width: 1200px;
  }
}

@media (max-width: 1200px) {
  .node-management {
    padding: 15px;
  }

  .card {
    padding: 30px;
  }

  .node-stats {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 992px) {
  .card-header {
    flex-direction: column;
    gap: 20px;
    align-items: flex-start;
  }

  .header-actions {
    width: 100%;
    justify-content: flex-start;
  }

  .table-controls {
    flex-direction: column;
    align-items: stretch;
  }

  .search-box {
    min-width: 100%;
  }

  .filter-controls {
    width: 100%;
  }

  .filter-select {
    flex: 1;
  }

  .detail-grid,
  .performance-stats,
  .config-info {
    grid-template-columns: 1fr;
  }

  .modal-dialog {
    margin: 10px;
  }
}

@media (max-width: 768px) {
  .node-management {
    padding: 10px;
  }

  .card {
    padding: 25px;
    border-radius: 16px;
  }

  .card-title {
    font-size: 22px;
  }

  .card-title i {
    font-size: 24px;
    padding: 10px;
  }

  .node-stats {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .stat-icon {
    width: 60px;
    height: 60px;
    margin-right: 20px;
  }

  .stat-value {
    font-size: 30px;
  }

  .btn {
    padding: 12px 20px;
    font-size: 15px;
  }

  .node-table th,
  .node-table td {
    padding: 18px 15px;
    font-size: 14px;
  }

  .action-buttons {
    gap: 8px;
  }

  .action-btn {
    width: 36px;
    height: 36px;
    font-size: 15px;
  }

  .pagination {
    flex-direction: column;
    gap: 15px;
  }

  .modal-header {
    padding: 20px 25px;
  }

  .modal-content {
    padding: 25px;
  }

  .modal-footer {
    padding: 20px 25px;
  }

  .add-node-form .form-control {
    padding: 14px 18px;
    font-size: 15px;
  }
}

@media (max-width: 480px) {
  .card {
    padding: 20px;
    border-radius: 14px;
  }

  .card-title {
    font-size: 20px;
  }

  .stat-item {
    padding: 20px;
  }

  .stat-icon {
    width: 50px;
    height: 50px;
    margin-right: 15px;
    font-size: 28px;
  }

  .stat-value {
    font-size: 28px;
  }

  .table-controls {
    padding: 15px;
  }

  .search-box {
    padding: 10px 15px;
  }

  .filter-select {
    min-width: 100%;
    padding: 10px 15px;
  }

  .node-table th,
  .node-table td {
    padding: 15px 12px;
  }

  .empty-state,
  .loading-state {
    padding: 60px 20px;
  }

  .empty-state i {
    font-size: 56px;
  }

  .empty-state h3 {
    font-size: 22px;
  }
}

/* 打印样式 */
@media print {
  .node-management {
    padding: 0;
  }

  .card {
    box-shadow: none;
    border: 1px solid #ddd;
    break-inside: avoid;
  }

  .btn,
  .search-box,
  .filter-controls,
  .action-buttons,
  .pagination {
    display: none !important;
  }

  .node-table {
    font-size: 12px;
  }
}

/* 滚动条美化 */
::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}

::-webkit-scrollbar-track {
  background: var(--light-gray);
  border-radius: 5px;
}

::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
  border-radius: 5px;
  border: 2px solid var(--light-gray);
}

::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, var(--primary-dark) 0%, var(--primary) 100%);
}

/* 焦点样式 */
:focus-visible {
  outline: 2px solid var(--primary);
  outline-offset: 2px;
  border-radius: 4px;
}

/* 工具类 */
.text-truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.d-inline-flex {
  display: inline-flex;
}

.align-middle {
  vertical-align: middle;
}
</style>
