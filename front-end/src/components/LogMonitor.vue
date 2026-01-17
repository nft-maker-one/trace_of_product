<template>
  <div class="log-monitor">
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">
          <i class="fas fa-clipboard-list"></i>
          {{ t('logMonitor.title') }}
        </h3>
        <div class="header-actions">
          <button class="btn btn-primary" @click="fetchLogs">
            <i class="fas fa-sync-alt" :class="{ 'spin': loading }"></i>
            {{ loading ? t('common.loading') : t('logMonitor.refresh') }}
          </button>
          <button class="btn btn-outline" @click="clearLogs">
            <i class="fas fa-trash-alt"></i>
            {{ t('logMonitor.clear') }}
          </button>
        </div>
      </div>

      <!-- 日志统计 -->
      <div class="node-stats">
        <div class="stat-item">
          <div class="stat-icon total">
            <i class="fas fa-file-alt"></i>
          </div>
          <div class="stat-info">
            <div class="stat-label">{{ t('logMonitor.debug') }}</div>
            <div class="stat-value">{{ logStats.debug }}</div>
          </div>
        </div>
        <div class="stat-item">
          <div class="stat-icon info">
            <i class="fas fa-info-circle"></i>
          </div>
          <div class="stat-info">
            <div class="stat-label">{{ t('logMonitor.info') }}</div>
            <div class="stat-value">{{ logStats.info }}</div>
          </div>
        </div>
        <div class="stat-item">
          <div class="stat-icon warning">
            <i class="fas fa-exclamation-triangle"></i>
          </div>
          <div class="stat-info">
            <div class="stat-label">{{ t('logMonitor.warning') }}</div>
            <div class="stat-value">{{ logStats.warning }}</div>
          </div>
        </div>
        <div class="stat-item">
          <div class="stat-icon error">
            <i class="fas fa-times-circle"></i>
          </div>
          <div class="stat-info">
            <div class="stat-label">{{ t('logMonitor.error') }}</div>
            <div class="stat-value">{{ logStats.error }}</div>
          </div>
        </div>
      </div>

      <!-- 日志列表 -->
      <div class="node-list-container">
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>{{ t('logMonitor.loading') }}</p>
        </div>

        <div v-else-if="filteredLogs.length === 0" class="empty-state">
          <i class="fas fa-clipboard-list"></i>
          <h3>{{ t('logMonitor.noLogs') }}</h3>
          <p>{{ t('logMonitor.noLogsDescription') }}</p>
          <button class="btn btn-primary" @click="fetchLogs">
            <i class="fas fa-sync-alt"></i>
            {{ t('logMonitor.refresh') }}
          </button>
        </div>

        <div v-else class="log-table-wrapper">
          <div class="table-controls">
            <div class="search-box">
              <i class="fas fa-search"></i>
              <input
                v-model="searchQuery"
                type="text"
                :placeholder="t('logMonitor.searchLogs')"
                @input="filterLogs"
              >
            </div>
            <div class="filter-controls">
              <select v-model="levelFilter" @change="filterLogs" class="filter-select">
                <option value="">{{ t('logMonitor.allLevels') }}</option>
                <option value="info">{{ t('logMonitor.info') }}</option>
                <option value="warning">{{ t('logMonitor.warning') }}</option>
                <option value="error">{{ t('logMonitor.error') }}</option>
                <option value="debug">{{ t('logMonitor.debug') }}</option>
              </select>
              <select v-model="timeFilter" @change="filterLogs" class="filter-select">
                <option value="all">{{ t('logMonitor.allTime') }}</option>
                <option value="today">{{ t('logMonitor.today') }}</option>
                <option value="week">{{ t('logMonitor.lastWeek') }}</option>
                <option value="month">{{ t('logMonitor.lastMonth') }}</option>
              </select>
              <select v-model="limit" @change="fetchLogs" class="filter-select">
                <option value="50">{{ t('logMonitor.count50') }}</option>
                <option value="100">{{ t('logMonitor.count100') }}</option>
                <option value="200">{{ t('logMonitor.count200') }}</option>
                <option value="500">{{ t('logMonitor.count500') }}</option>
              </select>
            </div>
          </div>

          <div class="log-entries-container">
            <div class="log-entries">
              <div
                v-for="log in paginatedLogs"
                :key="log.timestamp + log.message"
                :class="['log-entry', log.level]"
              >
                <div class="log-entry-header">
                  <div class="log-timestamp">
                    <i class="fas fa-clock"></i>
                    {{ formatTime(log.timestamp) }}
                  </div>
                  <span :class="['log-level-badge', log.level]">
                    {{ getLevelText(log.level) }}
                  </span>
                </div>
                <div class="log-message">
                  <i class="fas fa-chevron-right"></i>
                  {{ log.message }}
                </div>
                <div v-if="log.fields && Object.keys(log.fields).length > 0" class="log-fields">
                  <div class="fields-title">
                    <i class="fas fa-list"></i>
                    {{ t('logMonitor.details') }}
                  </div>
                  <div class="fields-grid">
                    <div v-for="(value, key) in log.fields" :key="key" class="field-item">
                      <span class="field-key">{{ key }}:</span>
                      <span class="field-value">{{ value }}</span>
                    </div>
                  </div>
                </div>
                <div class="log-actions">
                  <button class="action-btn" @click="copyLog(log)" title="复制日志">
                    <i class="fas fa-copy"></i>
                  </button>
                  <button class="action-btn" @click="analyzeLog(log)" title="分析日志">
                    <i class="fas fa-chart-bar"></i>
                  </button>
                  <button v-if="log.level === 'error'" class="action-btn danger" @click="handleError(log)" title="处理错误">
                    <i class="fas fa-wrench"></i>
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- 分页控件 -->
          <div v-if="filteredLogs.length > 0" class="pagination">
            <div class="pagination-info">
              {{ t('logMonitor.paginationInfo', { start: pagination.start + 1, end: Math.min(pagination.end, filteredLogs.length), total: filteredLogs.length }) }}
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

    <!-- 日志详情对话框 -->
    <div v-if="showLogDetail" class="modal-overlay">
      <div class="modal-dialog">
        <div class="modal-header">
          <h3>
            <i class="fas fa-file-alt"></i>
            {{ t('logMonitor.logDetails') }}
          </h3>
          <button class="close-btn" @click="showLogDetail = false">
            <i class="fas fa-times"></i>
          </button>
        </div>
        <div class="modal-content">
          <div v-if="selectedLog" class="log-detail">
            <div class="detail-section">
              <h4 class="section-title">
                <i class="fas fa-info-circle"></i>
                {{ t('logMonitor.basicInfo') }}
              </h4>
              <div class="detail-grid">
                <div class="detail-item">
                  <div class="detail-label">{{ t('logMonitor.timestamp') }}:</div>
                  <div class="detail-value">{{ formatTime(selectedLog.timestamp) }}</div>
                </div>
                <div class="detail-item">
                  <div class="detail-label">{{ t('logMonitor.level') }}:</div>
                  <div class="detail-value">
                    <span :class="['level-badge', selectedLog.level]">
                      {{ getLevelText(selectedLog.level) }}
                    </span>
                  </div>
                </div>
                <div class="detail-item">
                  <div class="detail-label">{{ t('logMonitor.source') }}:</div>
                  <div class="detail-value">{{ selectedLog.source || 'System' }}</div>
                </div>
                <div class="detail-item">
                  <div class="detail-label">{{ t('logMonitor.module') }}:</div>
                  <div class="detail-value">{{ selectedLog.module || 'Core' }}</div>
                </div>
              </div>
            </div>

            <div class="detail-section">
              <h4 class="section-title">
                <i class="fas fa-comment"></i>
                {{ t('logMonitor.message') }}
              </h4>
              <div class="message-content">
                <pre>{{ selectedLog.message }}</pre>
              </div>
            </div>

            <div v-if="selectedLog.fields && Object.keys(selectedLog.fields).length > 0" class="detail-section">
              <h4 class="section-title">
                <i class="fas fa-list-ul"></i>
                {{ t('logMonitor.additionalInfo') }}
              </h4>
              <div class="fields-detail">
                <div v-for="(value, key) in selectedLog.fields" :key="key" class="field-row">
                  <div class="field-label">{{ key }}:</div>
                  <div class="field-content">
                    <pre>{{ typeof value === 'object' ? JSON.stringify(value, null, 2) : value }}</pre>
                  </div>
                </div>
              </div>
            </div>

            <div class="detail-section">
              <h4 class="section-title">
                <i class="fas fa-chart-line"></i>
                {{ t('logMonitor.analysis') }}
              </h4>
              <div class="analysis-stats">
                <div class="analysis-item">
                  <div class="analysis-label">{{ t('logMonitor.similarErrors') }}:</div>
                  <div class="analysis-value">{{ getSimilarCount(selectedLog) }}</div>
                </div>
                <div class="analysis-item">
                  <div class="analysis-label">{{ t('logMonitor.firstOccurrence') }}:</div>
                  <div class="analysis-value">{{ getFirstOccurrence(selectedLog) }}</div>
                </div>
                <div class="analysis-item">
                  <div class="analysis-label">{{ t('logMonitor.frequency') }}:</div>
                  <div class="analysis-value">
                    <span class="frequency-badge" :class="getFrequencyClass(selectedLog)">
                      {{ getFrequencyText(selectedLog) }}
                    </span>
                  </div>
                </div>
                <div class="analysis-item">
                  <div class="analysis-label">{{ t('logMonitor.impact') }}:</div>
                  <div class="analysis-value">
                    <span class="impact-badge" :class="getImpactClass(selectedLog)">
                      {{ getImpactText(selectedLog) }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="showLogDetail = false">
            {{ t('common.close') }}
          </button>
          <button class="btn btn-primary" @click="exportLog(selectedLog)">
            <i class="fas fa-download"></i>
            {{ t('logMonitor.export') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const logs = ref([])
const filteredLogs = ref([])
const loading = ref(false)
const showLogDetail = ref(false)
const selectedLog = ref(null)
const searchQuery = ref('')
const levelFilter = ref('')
const timeFilter = ref('all')
const limit = ref(100)

// 分页
const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  totalPages: 1,
  start: 0,
  end: 10
})

// 日志统计
const logStats = reactive({
  info: 0,
  warning: 0,
  error: 0,
  debug: 0
})

// 计算属性
const paginatedLogs = computed(() => {
  return filteredLogs.value.slice(pagination.start, pagination.end)
})

const levelMap = computed(() => ({
  info: t('logMonitor.info'),
  warning: t('logMonitor.warning'),
  error: t('logMonitor.error'),
  debug: t('logMonitor.debug')
}))

// 方法
function getLevelText(level) {
  return levelMap.value[level] || level
}

function formatTime(timestamp) {
  if (!timestamp) return 'Unknown'
  try {
    const date = new Date(timestamp)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch {
    return timestamp
  }
}

async function fetchLogs() {
  loading.value = true
  // try {
  //   const response = await api.getLogs(limit.value)
  //   logs.value = response.map(log => ({
  //     ...log,
  //     // 添加模拟数据用于演示
  //     source: log.source || 'server',
  //     module: log.module || 'api',
  //     timestamp: log.timestamp || new Date().toISOString()
  //   }))

  //   updateLogStats()
  //   filterLogs()
  // } catch (error) {
  //   console.error('Failed to fetch logs:', error)
  //   logs.value = []
  //   filteredLogs.value = []
  //   alert('加载日志失败: ' + error.message)
  // } finally {
  //   loading.value = false
  // }
}

function updateLogStats() {
  logStats.total = logs.value.length
  logStats.info = logs.value.filter(log => log.level === 'info').length
  logStats.warning = logs.value.filter(log => log.level === 'warning').length
  logStats.error = logs.value.filter(log => log.level === 'error').length
  logStats.debug = logs.value.filter(log => log.level === 'debug').length
}

function filterLogs() {
  let filtered = [...logs.value]

  // 搜索过滤
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(log =>
      log.message.toLowerCase().includes(query) ||
      (log.source && log.source.toLowerCase().includes(query)) ||
      (log.module && log.module.toLowerCase().includes(query))
    )
  }

  // 级别过滤
  if (levelFilter.value) {
    filtered = filtered.filter(log => log.level === levelFilter.value)
  }

  // 时间过滤
  if (timeFilter.value !== 'all') {
    const now = new Date()
    let cutoffDate = new Date()

    switch (timeFilter.value) {
      case 'today':
        cutoffDate.setDate(now.getDate() - 1)
        break
      case 'week':
        cutoffDate.setDate(now.getDate() - 7)
        break
      case 'month':
        cutoffDate.setMonth(now.getMonth() - 1)
        break
    }

    filtered = filtered.filter(log => new Date(log.timestamp) >= cutoffDate)
  }

  // 按时间倒序排序
  filtered.sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp))

  // 更新分页
  filteredLogs.value = filtered
  updatePagination()
}

function updatePagination() {
  pagination.totalPages = Math.ceil(filteredLogs.value.length / pagination.pageSize)
  pagination.currentPage = Math.min(pagination.currentPage, pagination.totalPages)
  pagination.start = (pagination.currentPage - 1) * pagination.pageSize
  pagination.end = pagination.start + pagination.pageSize
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

function getSimilarCount(log) {
  // 模拟相似错误数量
  return logs.value.filter(l =>
    l.level === log.level &&
    l.message.includes(log.message.substring(0, 50))
  ).length
}

function getFirstOccurrence(log) {
  // 模拟首次出现时间
  const firstTime = new Date(log.timestamp)
  firstTime.setHours(firstTime.getHours() - Math.floor(Math.random() * 24))
  return formatTime(firstTime.toISOString())
}

function getFrequencyClass(log) {
  const count = getSimilarCount(log)
  if (count > 10) return 'high'
  if (count > 5) return 'medium'
  return 'low'
}

function getFrequencyText(log) {
  const cls = getFrequencyClass(log)
  const textMap = {
    'high': '高',
    'medium': '中',
    'low': '低'
  }
  return textMap[cls] || '未知'
}

function getImpactClass(log) {
  if (log.level === 'error') return 'high'
  if (log.level === 'warning') return 'medium'
  return 'low'
}

function getImpactText(log) {
  if (log.level === 'error') return '高'
  if (log.level === 'warning') return '中'
  return '低'
}

function copyLog(log) {
  const logText = `[${formatTime(log.timestamp)}] [${log.level.toUpperCase()}] ${log.message}`
  navigator.clipboard.writeText(logText)
    .then(() => {
      alert('日志已复制到剪贴板')
    })
    .catch(err => {
      console.error('复制失败:', err)
    })
}

function analyzeLog(log) {
  selectedLog.value = log
  showLogDetail.value = true
}

function handleError(log) {
  if (confirm(`确定要处理此错误吗？\n${log.message.substring(0, 100)}...`)) {
    alert(`正在处理错误...\n实际实现中会调用错误处理API`)
  }
}

function clearLogs() {
  if (confirm(t('logMonitor.confirmClear'))) {
    filteredLogs.value = []
    updateLogStats()
    updatePagination()
    alert('日志已清空')
  }
}

function exportLog(log) {
  const logData = {
    ...log,
    exportedAt: new Date().toISOString()
  }

  const blob = new Blob([JSON.stringify(logData, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `log_${log.timestamp.replace(/[:.]/g, '-')}.json`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)

  alert('日志已导出')
}

onMounted(() => {
  fetchLogs()
})
</script>

<style scoped>
/* 继承第一个文件的基础样式，保持一致性 */
.log-monitor {
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

/* 卡片样式 - 与第一个文件一致 */
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

/* 卡片头部 - 与第一个文件一致 */
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
  padding: 12px 24px;
  border-radius: 25px;
  font-size: 14px;
  font-weight: 600;
  text-decoration: none;
  border: none;
  cursor: pointer;
  transition: all 0.3s ease;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  position: relative;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, transparent 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.btn:hover::before {
  opacity: 1;
}

.btn:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(46, 125, 50, 0.3);
}

.btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(46, 125, 50, 0.2);
}

.btn-primary {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(46, 125, 50, 0.3);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(46, 125, 50, 0.4);
}

.btn-primary:focus {
  box-shadow: 0 0 0 3px rgba(46, 125, 50, 0.4);
}

.btn-outline {
  background: transparent;
  color: var(--primary);
  border: 2px solid var(--primary);
  box-shadow: 0 4px 15px rgba(46, 125, 50, 0.1);
}

.btn-outline:hover {
  background: var(--primary);
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(46, 125, 50, 0.2);
}

.btn-outline:focus {
  box-shadow: 0 0 0 3px rgba(46, 125, 50, 0.2);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
  box-shadow: none !important;
  filter: grayscale(50%);
}

/* 统计卡片 - 调整图标颜色 */
.stat-icon.info {
  background: linear-gradient(135deg, var(--info) 0%, #64b5f6 100%);
  color: white;
}

.stat-icon.warning {
  background: linear-gradient(135deg, var(--warning) 0%, #ffb74d 100%);
  color: white;
}

/* 日志统计 */
.node-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 25px;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 15px;
  border: 1px solid var(--light-gray);
  transition: var(--transition);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

.stat-item:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
  border-color: var(--primary-light);
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  transition: var(--transition);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-icon.total {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
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
  font-size: 14px;
  color: var(--gray);
  margin-bottom: 8px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-value {
  font-size: 32px;
  font-weight: 800;
  color: var(--dark);
  line-height: 1;
}

/* 表格控制 */
.table-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  padding: 25px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 15px;
  margin-bottom: 25px;
  border: 1px solid var(--light-gray);
  flex-wrap: wrap;
}

.search-box {
  position: relative;
  flex: 1;
  min-width: 300px;
}

.search-box i {
  position: absolute;
  left: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--gray);
  font-size: 16px;
}

.search-box input {
  width: 100%;
  padding: 12px 15px 12px 45px;
  border: 2px solid var(--light-gray);
  border-radius: 25px;
  font-size: 14px;
  transition: var(--transition);
  background: white;
}

.search-box input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(46, 125, 50, 0.1);
}

.filter-controls {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.filter-select {
  padding: 10px 15px;
  border: 2px solid var(--light-gray);
  border-radius: 10px;
  font-size: 14px;
  background: white;
  cursor: pointer;
  transition: var(--transition);
  min-width: 120px;
}

.filter-select:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(46, 125, 50, 0.1);
}

/* 日志条目容器 */
.log-entries-container {
  max-height: 600px;
  overflow-y: auto;
  border: 1px solid var(--light-gray);
  border-radius: 15px;
  background: white;
  scrollbar-width: thin;
}

.log-entries {
  padding: 5px;
}

/* 日志条目 */
.log-entry {
  padding: 25px;
  margin: 10px;
  border-radius: 12px;
  border: 1px solid var(--light-gray);
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  transition: var(--transition);
  position: relative;
  overflow: hidden;
  animation: slideInRight 0.5s ease;
}

@keyframes slideInRight {
  from {
    opacity: 0;
    transform: translateX(20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.log-entry::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
  transition: var(--transition);
}

.log-entry:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
  border-color: var(--primary-light);
}

.log-entry.info {
  border-left: 4px solid var(--info);
}

.log-entry.info::before {
  background: linear-gradient(180deg, var(--info) 0%, rgba(33, 150, 243, 0.3) 100%);
}

.log-entry.warning {
  border-left: 4px solid var(--warning);
}

.log-entry.warning::before {
  background: linear-gradient(180deg, var(--warning) 0%, rgba(255, 152, 0, 0.3) 100%);
}

.log-entry.error {
  border-left: 4px solid var(--error);
}

.log-entry.error::before {
  background: linear-gradient(180deg, var(--error) 0%, rgba(244, 67, 54, 0.3) 100%);
}

.log-entry.debug {
  border-left: 4px solid #9c27b0;
}

.log-entry.debug::before {
  background: linear-gradient(180deg, #9c27b0 0%, rgba(156, 39, 176, 0.3) 100%);
}

/* 日志头部 */
.log-entry-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  flex-wrap: wrap;
  gap: 15px;
}

.log-timestamp {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--gray);
  font-size: 14px;
  font-family: 'Courier New', monospace;
  background: rgba(255, 255, 255, 0.8);
  padding: 8px 15px;
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.log-timestamp i {
  color: var(--primary);
}

.log-level-badge {
  padding: 8px 20px;
  border-radius: 25px;
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  transition: var(--transition);
  position: relative;
  overflow: hidden;
}

.log-level-badge::before {
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

.log-level-badge:hover::before {
  opacity: 1;
}

.log-level-badge.info {
  background: linear-gradient(135deg, rgba(33, 150, 243, 0.15) 0%, rgba(33, 150, 243, 0.1) 100%);
  color: var(--info);
  border: 1px solid rgba(33, 150, 243, 0.2);
  box-shadow: 0 4px 12px rgba(33, 150, 243, 0.1);
}

.log-level-badge.warning {
  background: linear-gradient(135deg, rgba(255, 152, 0, 0.15) 0%, rgba(255, 152, 0, 0.1) 100%);
  color: var(--warning);
  border: 1px solid rgba(255, 152, 0, 0.2);
  box-shadow: 0 4px 12px rgba(255, 152, 0, 0.1);
}

.log-level-badge.error {
  background: linear-gradient(135deg, rgba(244, 67, 54, 0.15) 0%, rgba(244, 67, 54, 0.1) 100%);
  color: var(--error);
  border: 1px solid rgba(244, 67, 54, 0.2);
  box-shadow: 0 4px 12px rgba(244, 67, 54, 0.1);
}

.log-level-badge.debug {
  background: linear-gradient(135deg, rgba(156, 39, 176, 0.15) 0%, rgba(156, 39, 176, 0.1) 100%);
  color: #9c27b0;
  border: 1px solid rgba(156, 39, 176, 0.2);
  box-shadow: 0 4px 12px rgba(156, 39, 176, 0.1);
}

/* 日志消息 */
.log-message {
  color: var(--dark);
  margin-bottom: 15px;
  line-height: 1.6;
  font-size: 15px;
  padding-left: 8px;
  position: relative;
  background: white;
  padding: 15px;
  border-radius: 10px;
  border: 1px solid var(--border-color);
  font-family: 'Courier New', monospace;
}

.log-message::before {
  content: '▶';
  position: absolute;
  left: 0;
  color: var(--primary);
  font-size: 12px;
}

/* 日志字段 */
.log-fields {
  background: linear-gradient(135deg, #ffffff 0%, #f5f7fa 100%);
  padding: 20px;
  border-radius: 10px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  border: 1px solid var(--border-color);
  margin-bottom: 15px;
  transition: var(--transition);
}

.fields-title {
  font-weight: 600;
  color: var(--primary-dark);
  margin-bottom: 15px;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  padding-bottom: 10px;
  border-bottom: 2px solid var(--light-gray);
}

.fields-title i {
  color: var(--primary);
}

.fields-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 15px;
}

.field-item {
  display: flex;
  flex-direction: column;
  gap: 5px;
  padding: 12px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 8px;
  border: 1px solid transparent;
  transition: var(--transition);
}

.field-item:hover {
  border-color: var(--primary);
  transform: translateY(-2px);
}

.field-key {
  font-weight: 600;
  color: var(--gray);
  text-transform: uppercase;
  font-size: 12px;
  letter-spacing: 0.5px;
}

.field-value {
  color: var(--dark);
  word-break: break-word;
  line-height: 1.4;
  font-size: 14px;
}

/* 日志操作 */
.log-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid var(--light-gray);
}

.action-btn {
  padding: 8px 12px;
  border: none;
  border-radius: 8px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  color: var(--gray);
  cursor: pointer;
  transition: var(--transition);
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.action-btn:hover {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(46, 125, 50, 0.2);
}

.action-btn.danger {
  background: linear-gradient(135deg, var(--error) 0%, #ef5350 100%);
  color: white;
}

.action-btn.danger:hover {
  background: linear-gradient(135deg, #d32f2f 0%, var(--error) 100%);
  box-shadow: 0 4px 15px rgba(244, 67, 54, 0.3);
}

/* 分页控件 */
.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 25px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 15px;
  margin-top: 25px;
  border: 1px solid var(--light-gray);
  flex-wrap: wrap;
  gap: 20px;
}

.pagination-info {
  color: var(--gray);
  font-size: 14px;
  font-weight: 500;
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 10px;
}

.pagination-btn {
  padding: 10px 15px;
  border: 2px solid var(--light-gray);
  background: white;
  color: var(--gray);
  border-radius: 8px;
  cursor: pointer;
  transition: var(--transition);
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pagination-btn:hover:not(:disabled) {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(46, 125, 50, 0.2);
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.page-numbers {
  display: flex;
  gap: 5px;
  margin: 0 15px;
}

.page-btn {
  padding: 8px 12px;
  border: 2px solid var(--light-gray);
  background: white;
  color: var(--gray);
  border-radius: 6px;
  cursor: pointer;
  transition: var(--transition);
  font-size: 14px;
  min-width: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.page-btn:hover {
  background: var(--primary-light);
  color: white;
  border-color: var(--primary-light);
}

.page-btn.active {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
  color: white;
  border-color: var(--primary);
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(46, 125, 50, 0.2);
}

/* 加载和空状态 */
.loading-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--gray);
}

.loading-state .spinner {
  width: 50px;
  height: 50px;
  border: 4px solid var(--light-gray);
  border-top: 4px solid var(--primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 20px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--gray);
}

.empty-state i {
  font-size: 64px;
  color: var(--light-gray);
  margin-bottom: 20px;
  display: block;
}

.empty-state h3 {
  font-size: 24px;
  margin-bottom: 10px;
  color: var(--dark);
}

.empty-state p {
  font-size: 16px;
  margin-bottom: 30px;
}

/* 模态对话框 - 与第一个文件一致 */
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

.close-btn:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(244, 67, 54, 0.3);
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

/* 日志详情样式 */
.log-detail {
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
.detail-section:nth-child(4) { animation-delay: 0.4s; }

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
}

/* 消息内容 */
.message-content {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  padding: 25px;
  border-radius: 12px;
  border: 1px solid var(--border-color);
  font-family: 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
  max-height: 300px;
  overflow-y: auto;
}

.message-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
  color: var(--dark);
}

/* 字段详情 */
.fields-detail {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.field-row {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 20px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 12px;
  border: 1px solid transparent;
  transition: var(--transition);
}

.field-row:hover {
  border-color: var(--primary);
  transform: translateY(-2px);
}

.field-label {
  font-weight: 600;
  color: var(--primary-dark);
  font-size: 15px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.field-content {
  color: var(--dark);
  font-family: 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
  max-height: 200px;
  overflow-y: auto;
  padding: 15px;
  background: white;
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.field-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* 分析统计 */
.analysis-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 25px;
}

.analysis-item {
  padding: 25px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 15px;
  border: 1px solid transparent;
  transition: var(--transition);
  text-align: center;
}

.analysis-item:hover {
  border-color: var(--primary);
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.08);
}

.analysis-label {
  font-size: 15px;
  color: var(--gray);
  margin-bottom: 15px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.analysis-value {
  font-size: 28px;
  font-weight: 800;
  color: var(--dark);
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 8px;
  line-height: 1;
}

.frequency-badge,
.impact-badge {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 10px 20px;
  border-radius: 25px;
  font-size: 14px;
  font-weight: 700;
  transition: var(--transition);
}

.frequency-badge.high,
.impact-badge.high {
  background: linear-gradient(135deg, rgba(244, 67, 54, 0.15) 0%, rgba(244, 67, 54, 0.1) 100%);
  color: var(--error);
  border: 1px solid rgba(244, 67, 54, 0.2);
}

.frequency-badge.medium,
.impact-badge.medium {
  background: linear-gradient(135deg, rgba(255, 152, 0, 0.15) 0%, rgba(255, 152, 0, 0.1) 100%);
  color: var(--warning);
  border: 1px solid rgba(255, 152, 0, 0.2);
}

.frequency-badge.low,
.impact-badge.low {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.15) 0%, rgba(76, 175, 80, 0.1) 100%);
  color: var(--success);
  border: 1px solid rgba(76, 175, 80, 0.2);
}

/* 响应式设计 - 与第一个文件保持一致 */
@media (max-width: 1400px) {
  .log-monitor {
    max-width: 1200px;
  }
}

@media (max-width: 1200px) {
  .log-monitor {
    padding: 15px;
  }

  .card {
    padding: 30px;
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

  .node-stats {
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 15px;
  }

  .stat-item {
    padding: 20px;
  }

  .stat-icon {
    width: 50px;
    height: 50px;
    font-size: 20px;
  }

  .stat-value {
    font-size: 28px;
  }

  .log-entry-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .pagination {
    flex-direction: column;
    align-items: center;
    gap: 15px;
  }

  .pagination-controls {
    order: -1;
  }

  .detail-grid,
  .analysis-stats {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .log-monitor {
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

  .btn {
    padding: 12px 20px;
    font-size: 15px;
  }

  .table-controls {
    padding: 20px;
    gap: 15px;
  }

  .search-box {
    min-width: 100%;
  }

  .filter-controls {
    flex-direction: column;
    width: 100%;
  }

  .filter-select {
    min-width: 100%;
    padding: 12px 15px;
  }

  .node-stats {
    grid-template-columns: 1fr;
    gap: 15px;
  }

  .stat-item {
    padding: 20px;
    flex-direction: column;
    text-align: center;
    gap: 15px;
  }

  .stat-icon {
    width: 50px;
    height: 50px;
    font-size: 20px;
  }

  .stat-value {
    font-size: 24px;
  }

  .log-entry {
    padding: 20px;
  }

  .log-message {
    font-size: 14px;
  }

  .fields-grid {
    grid-template-columns: 1fr;
  }

  .modal-dialog {
    margin: 10px;
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

  .pagination {
    padding: 20px;
  }

  .page-numbers {
    margin: 0 10px;
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

  .node-stats {
    gap: 10px;
  }

  .stat-item {
    padding: 15px;
  }

  .stat-icon {
    width: 40px;
    height: 40px;
    font-size: 16px;
  }

  .stat-value {
    font-size: 20px;
  }

  .log-entry {
    padding: 15px;
  }

  .log-timestamp {
    font-size: 12px;
    padding: 6px 10px;
  }

  .log-level-badge {
    padding: 6px 15px;
    font-size: 11px;
  }

  .log-message {
    font-size: 13px;
    padding: 12px;
  }

  .log-fields {
    padding: 15px;
  }

  .field-key {
    font-size: 11px;
  }

  .field-value {
    font-size: 13px;
  }

  .action-btn {
    padding: 6px 10px;
    font-size: 12px;
  }

  .pagination {
    padding: 15px;
    gap: 10px;
  }

  .pagination-info {
    font-size: 12px;
  }

  .page-btn {
    padding: 6px 10px;
    font-size: 12px;
    min-width: 35px;
  }

  .loading-state,
  .empty-state {
    padding: 40px 20px;
  }

  .loading-state .spinner {
    width: 40px;
    height: 40px;
  }

  .empty-state i {
    font-size: 48px;
  }

  .empty-state h3 {
    font-size: 18px;
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

/* 按钮样式完善 */
.btn:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(46, 125, 50, 0.3);
}

.btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(46, 125, 50, 0.2);
}

.btn-primary:focus {
  box-shadow: 0 0 0 3px rgba(46, 125, 50, 0.4);
}

.btn-outline:focus {
  box-shadow: 0 0 0 3px rgba(46, 125, 50, 0.2);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
  box-shadow: none !important;
  filter: grayscale(50%);
}

/* 旋转动画 */
.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
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
