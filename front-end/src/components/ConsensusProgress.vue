<template>
  <div class="consensus-progress" v-if="visible">
    <div class="consensus-overlay"></div>
    <div class="consensus-modal">
      <div class="consensus-header">
        <h3>
          <i class="fas fa-network-wired"></i>
          {{ t('consensus.title') }}
        </h3>
        <div class="header-badge" :class="statusClass">
          {{ statusText }}
        </div>
      </div>

      <div class="consensus-content">
        <!-- 上传信息 -->
        <div class="upload-info">
          <div class="info-row">
            <span class="label">{{ t('consensus.productId') }}:</span>
            <span class="value">{{ uploadData?.eggplant_id }}</span>
          </div>
          <div class="info-row">
            <span class="label">{{ t('consensus.uploadNode') }}:</span>
            <span class="value highlight">{{ formatNodeId(uploadNode?.id) }}</span>
          </div>
        </div>

        <!-- 实时进度条 -->
        <div class="progress-section">
          <div class="progress-bar-container">
            <div class="progress-bar" :style="{ width: progressPercent + '%' }"></div>
          </div>
          <div class="progress-text">{{ progressPercent }}%</div>
        </div>

        <!-- 网络节点状态图 -->
        <div class="network-graph">
          <div class="graph-title">
            <i class="fas fa-project-diagram"></i>
            {{ t('consensus.networkStatus') }}
          </div>
          
          <div class="nodes-container">
            <!-- 上传节点 (Leader) -->
            <div class="upload-node-wrapper">
              <div class="node-card leader" :class="{ active: currentPhase >= 1 }">
                <div class="node-icon">
                  <i class="fas fa-crown"></i>
                </div>
                <div class="node-info">
                  <span class="node-role">{{ t('consensus.leader') }}</span>
                  <span class="node-id">{{ formatNodeId(uploadNode?.id) }}</span>
                </div>
                <div class="node-status-indicator" :class="{ active: currentPhase >= 1 }">
                  <i class="fas fa-broadcast-tower"></i>
                </div>
              </div>
              
              <!-- 广播动画 -->
              <div class="broadcast-lines" v-if="currentPhase >= 1">
                <div 
                  v-for="(node, index) in otherNodes" 
                  :key="'line-' + node.id"
                  class="broadcast-line"
                  :class="{ active: currentPhase >= 2 }"
                  :style="{ '--delay': index * 0.2 + 's' }"
                ></div>
              </div>
            </div>

            <!-- 其他节点 -->
            <div class="peer-nodes-grid">
              <div 
                v-for="(node, index) in otherNodes" 
                :key="node.id"
                class="node-card peer"
                :class="{ 
                  receiving: receivingNodes.includes(node.id),
                  preparing: preparingNodes.includes(node.id),
                  committing: committingNodes.includes(node.id),
                  confirmed: confirmedNodes.includes(node.id)
                }"
                :style="{ '--delay': index * 0.15 + 's' }"
              >
                <div class="node-icon">
                  <i class="fas fa-server"></i>
                </div>
                <div class="node-info">
                  <span class="node-role">Node {{ index + 2 }}</span>
                  <span class="node-id">{{ formatNodeId(node.id) }}</span>
                </div>
                <div class="node-status-badge">
                  <template v-if="confirmedNodes.includes(node.id)">
                    <i class="fas fa-check-circle"></i>
                    <span>{{ t('consensus.confirmed') }}</span>
                  </template>
                  <template v-else-if="committingNodes.includes(node.id)">
                    <i class="fas fa-handshake"></i>
                    <span>{{ t('consensus.committing') }}</span>
                  </template>
                  <template v-else-if="preparingNodes.includes(node.id)">
                    <i class="fas fa-sync-alt spinning"></i>
                    <span>{{ t('consensus.preparing') }}</span>
                  </template>
                  <template v-else-if="receivingNodes.includes(node.id)">
                    <i class="fas fa-download"></i>
                    <span>{{ t('consensus.receiving') }}</span>
                  </template>
                  <template v-else>
                    <i class="fas fa-clock"></i>
                    <span>{{ t('consensus.waiting') }}</span>
                  </template>
                </div>
                
                <!-- 时间戳 -->
                <div class="node-timestamp" v-if="nodeTimestamps[node.id]">
                  {{ nodeTimestamps[node.id] }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- PBFT 阶段指示器 -->
        <div class="phase-indicators">
          <div class="phase" :class="{ active: currentPhase >= 1, completed: currentPhase > 1 }">
            <div class="phase-dot">
              <i v-if="currentPhase > 1" class="fas fa-check"></i>
              <span v-else>1</span>
            </div>
            <span class="phase-label">{{ t('consensus.prePrepare') }}</span>
          </div>
          <div class="phase-line" :class="{ active: currentPhase >= 2 }"></div>
          <div class="phase" :class="{ active: currentPhase >= 2, completed: currentPhase > 2 }">
            <div class="phase-dot">
              <i v-if="currentPhase > 2" class="fas fa-check"></i>
              <span v-else>2</span>
            </div>
            <span class="phase-label">{{ t('consensus.prepare') }}</span>
          </div>
          <div class="phase-line" :class="{ active: currentPhase >= 3 }"></div>
          <div class="phase" :class="{ active: currentPhase >= 3, completed: currentPhase > 3 }">
            <div class="phase-dot">
              <i v-if="currentPhase > 3" class="fas fa-check"></i>
              <span v-else>3</span>
            </div>
            <span class="phase-label">{{ t('consensus.commit') }}</span>
          </div>
          <div class="phase-line" :class="{ active: currentPhase >= 4 }"></div>
          <div class="phase" :class="{ active: currentPhase >= 4, completed: currentPhase >= 4 }">
            <div class="phase-dot">
              <i v-if="currentPhase >= 4" class="fas fa-check"></i>
              <span v-else>4</span>
            </div>
            <span class="phase-label">{{ t('consensus.reply') }}</span>
          </div>
        </div>

        <!-- 日志区域 -->
        <div class="consensus-logs">
          <div class="logs-header">
            <i class="fas fa-terminal"></i>
            {{ t('consensus.logs') }}
          </div>
          <div class="logs-content" ref="logsContainer">
            <div 
              v-for="(log, index) in logs" 
              :key="index" 
              class="log-entry"
              :class="log.type"
            >
              <span class="log-time">{{ log.time }}</span>
              <span class="log-message">{{ log.message }}</span>
            </div>
          </div>
        </div>

        <!-- 完成提示 -->
        <div class="completion-banner" v-if="isCompleted">
          <div class="completion-icon">
            <i class="fas fa-shield-alt"></i>
          </div>
          <div class="completion-text">
            <h4>{{ t('consensus.consensusSuccess') }}</h4>
            <p>{{ t('consensus.allNodesConfirmed', { count: confirmedNodes.length + 1 }) }}</p>
          </div>
        </div>
      </div>

      <div class="consensus-footer">
        <button 
          class="btn btn-primary" 
          @click="handleClose"
          :disabled="!isCompleted && !hasError"
        >
          {{ isCompleted ? t('consensus.done') : t('consensus.processing') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  uploadData: {
    type: Object,
    default: null
  },
  uploadNode: {
    type: Object,
    default: null
  },
  allNodes: {
    type: Array,
    default: () => []
  },
  uploadStatus: {
    type: String,
    default: 'pending' // pending, uploading, success, error
  }
})

const emit = defineEmits(['close', 'completed'])

const { t } = useI18n()

// 状态
const currentPhase = ref(0)
const receivingNodes = ref([])
const preparingNodes = ref([])
const committingNodes = ref([])
const confirmedNodes = ref([])
const nodeTimestamps = ref({})
const logs = ref([])
const logsContainer = ref(null)
const hasError = ref(false)

// 计算属性
const otherNodes = computed(() => {
  if (!props.uploadNode || !props.allNodes.length) return props.allNodes.slice(0, 3)
  return props.allNodes.filter(n => n.id !== props.uploadNode.id).slice(0, 3)
})

const isCompleted = computed(() => {
  return currentPhase.value >= 4 && confirmedNodes.value.length >= otherNodes.value.length
})

const progressPercent = computed(() => {
  if (currentPhase.value === 0) return 0
  if (currentPhase.value === 1) return 15
  if (currentPhase.value === 2) {
    const base = 15
    const nodeProgress = (preparingNodes.value.length / Math.max(otherNodes.value.length, 1)) * 25
    return Math.round(base + nodeProgress)
  }
  if (currentPhase.value === 3) {
    const base = 40
    const nodeProgress = (committingNodes.value.length / Math.max(otherNodes.value.length, 1)) * 30
    return Math.round(base + nodeProgress)
  }
  if (currentPhase.value >= 4) {
    const base = 70
    const nodeProgress = (confirmedNodes.value.length / Math.max(otherNodes.value.length, 1)) * 30
    return Math.round(base + nodeProgress)
  }
  return 0
})

const statusClass = computed(() => {
  if (hasError.value) return 'error'
  if (isCompleted.value) return 'success'
  return 'processing'
})

const statusText = computed(() => {
  if (hasError.value) return t('consensus.failed')
  if (isCompleted.value) return t('consensus.completed')
  return t('consensus.inProgress')
})

// 工具函数
function formatNodeId(id) {
  if (!id) return 'N/A'
  const idStr = String(id)
  if (idStr.length <= 8) return idStr
  return `${idStr.substring(0, 4)}...${idStr.substring(idStr.length - 4)}`
}

function getCurrentTime() {
  return new Date().toLocaleTimeString('zh-CN', { hour12: false })
}

function addLog(message, type = 'info') {
  logs.value.push({
    time: getCurrentTime(),
    message,
    type
  })
  
  // 自动滚动到底部
  nextTick(() => {
    if (logsContainer.value) {
      logsContainer.value.scrollTop = logsContainer.value.scrollHeight
    }
  })
}

// 模拟共识过程
function startConsensusAnimation() {
  // 重置状态
  currentPhase.value = 0
  receivingNodes.value = []
  preparingNodes.value = []
  committingNodes.value = []
  confirmedNodes.value = []
  nodeTimestamps.value = {}
  logs.value = []
  hasError.value = false

  const nodes = otherNodes.value

  // 阶段1: Pre-Prepare - Leader广播
  setTimeout(() => {
    currentPhase.value = 1
    addLog(t('consensus.logPrePrepare', { node: formatNodeId(props.uploadNode?.id) }), 'info')
  }, 500)

  // 节点开始接收
  setTimeout(() => {
    nodes.forEach((node, index) => {
      setTimeout(() => {
        receivingNodes.value.push(node.id)
        addLog(t('consensus.logReceiving', { node: formatNodeId(node.id) }), 'info')
      }, index * 300)
    })
  }, 1200)

  // 阶段2: Prepare - 节点验证并广播准备消息
  setTimeout(() => {
    currentPhase.value = 2
    addLog(t('consensus.logPreparePhase'), 'info')
    
    nodes.forEach((node, index) => {
      setTimeout(() => {
        preparingNodes.value.push(node.id)
        nodeTimestamps.value[node.id] = getCurrentTime()
        addLog(t('consensus.logNodePrepare', { node: formatNodeId(node.id) }), 'success')
      }, index * 400 + 200)
    })
  }, 2500)

  // 阶段3: Commit - 节点广播提交消息
  setTimeout(() => {
    currentPhase.value = 3
    addLog(t('consensus.logCommitPhase'), 'info')
    
    nodes.forEach((node, index) => {
      setTimeout(() => {
        committingNodes.value.push(node.id)
        nodeTimestamps.value[node.id] = getCurrentTime()
        addLog(t('consensus.logNodeCommit', { node: formatNodeId(node.id) }), 'success')
      }, index * 350 + 200)
    })
  }, 4500)

  // 阶段4: Reply - 达成共识
  setTimeout(() => {
    currentPhase.value = 4
    addLog(t('consensus.logReplyPhase'), 'info')
    
    nodes.forEach((node, index) => {
      setTimeout(() => {
        confirmedNodes.value.push(node.id)
        nodeTimestamps.value[node.id] = getCurrentTime()
        addLog(t('consensus.logNodeConfirm', { node: formatNodeId(node.id) }), 'success')
        
        // 最后一个节点确认后
        if (index === nodes.length - 1) {
          setTimeout(() => {
            addLog(t('consensus.logComplete'), 'success')
            emit('completed')
          }, 300)
        }
      }, index * 300 + 200)
    })
  }, 6500)
}

function handleClose() {
  emit('close')
}

// 监听显示状态
watch(() => props.visible, (newVal) => {
  if (newVal) {
    startConsensusAnimation()
  }
})

// 监听上传状态
watch(() => props.uploadStatus, (newVal) => {
  if (newVal === 'error') {
    hasError.value = true
    addLog(t('consensus.logError'), 'error')
  }
})
</script>

<style scoped>
.consensus-progress {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.consensus-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(8px);
}

.consensus-modal {
  position: relative;
  background: #1a1f2e;
  border-radius: 20px;
  width: 95%;
  max-width: 900px;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 25px 80px rgba(0, 0, 0, 0.5);
  animation: modalSlideIn 0.4s ease-out;
  color: #e0e0e0;
}

@keyframes modalSlideIn {
  from {
    opacity: 0;
    transform: translateY(-30px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.consensus-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
  border-bottom: 1px solid #3d4a5c;
}

.consensus-header h3 {
  margin: 0;
  font-size: 1.25rem;
  display: flex;
  align-items: center;
  gap: 10px;
  color: #fff;
}

.header-badge {
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
}

.header-badge.processing {
  background: linear-gradient(135deg, #3182ce, #2c5282);
  color: white;
  animation: pulse 2s infinite;
}

.header-badge.success {
  background: linear-gradient(135deg, #38a169, #276749);
  color: white;
}

.header-badge.error {
  background: linear-gradient(135deg, #e53e3e, #c53030);
  color: white;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

.consensus-content {
  padding: 24px;
  max-height: calc(90vh - 140px);
  overflow-y: auto;
}

/* 上传信息 */
.upload-info {
  display: flex;
  gap: 24px;
  padding: 16px 20px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  margin-bottom: 20px;
}

.info-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.info-row .label {
  color: #a0aec0;
  font-size: 0.9rem;
}

.info-row .value {
  color: #fff;
  font-weight: 600;
}

.info-row .value.highlight {
  color: #68d391;
}

/* 进度条 */
.progress-section {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.progress-bar-container {
  flex: 1;
  height: 8px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: linear-gradient(90deg, #38a169, #68d391);
  border-radius: 4px;
  transition: width 0.5s ease;
}

.progress-text {
  font-size: 1rem;
  font-weight: 600;
  color: #68d391;
  min-width: 50px;
  text-align: right;
}

/* 网络图 */
.network-graph {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 24px;
}

.graph-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 1rem;
  color: #a0aec0;
  margin-bottom: 20px;
}

.nodes-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.upload-node-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
}

.node-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  border: 2px solid transparent;
  transition: all 0.4s ease;
  position: relative;
}

.node-card.leader {
  background: linear-gradient(135deg, rgba(251, 191, 36, 0.1), rgba(245, 158, 11, 0.05));
  border-color: rgba(251, 191, 36, 0.3);
  max-width: 300px;
}

.node-card.leader.active {
  border-color: #fbbf24;
  box-shadow: 0 0 30px rgba(251, 191, 36, 0.2);
}

.node-card.peer {
  flex: 1;
  min-width: 200px;
}

.node-card.receiving {
  border-color: #3182ce;
  background: rgba(49, 130, 206, 0.1);
}

.node-card.preparing {
  border-color: #805ad5;
  background: rgba(128, 90, 213, 0.1);
}

.node-card.committing {
  border-color: #dd6b20;
  background: rgba(221, 107, 32, 0.1);
}

.node-card.confirmed {
  border-color: #38a169;
  background: rgba(56, 161, 105, 0.1);
  box-shadow: 0 0 20px rgba(56, 161, 105, 0.2);
}

.node-icon {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.1rem;
  background: rgba(255, 255, 255, 0.1);
  color: #a0aec0;
}

.node-card.leader .node-icon {
  background: linear-gradient(135deg, #fbbf24, #f59e0b);
  color: #1a1f2e;
}

.node-card.confirmed .node-icon {
  background: linear-gradient(135deg, #38a169, #2f855a);
  color: white;
}

.node-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.node-role {
  font-size: 0.8rem;
  color: #a0aec0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.node-id {
  font-size: 0.95rem;
  font-weight: 600;
  color: #fff;
  font-family: 'Monaco', 'Consolas', monospace;
}

.node-status-indicator {
  margin-left: auto;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  color: #a0aec0;
}

.node-status-indicator.active {
  background: linear-gradient(135deg, #38a169, #2f855a);
  color: white;
  animation: broadcast 1.5s ease-in-out infinite;
}

@keyframes broadcast {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.1); box-shadow: 0 0 20px rgba(56, 161, 105, 0.5); }
}

.node-status-badge {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 0.8rem;
  background: rgba(255, 255, 255, 0.1);
  color: #a0aec0;
}

.node-card.receiving .node-status-badge {
  background: rgba(49, 130, 206, 0.2);
  color: #63b3ed;
}

.node-card.preparing .node-status-badge {
  background: rgba(128, 90, 213, 0.2);
  color: #b794f4;
}

.node-card.committing .node-status-badge {
  background: rgba(221, 107, 32, 0.2);
  color: #fbd38d;
}

.node-card.confirmed .node-status-badge {
  background: rgba(56, 161, 105, 0.2);
  color: #68d391;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.node-timestamp {
  position: absolute;
  right: 12px;
  bottom: -20px;
  font-size: 0.7rem;
  color: #718096;
  font-family: 'Monaco', 'Consolas', monospace;
}

.peer-nodes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.broadcast-lines {
  display: flex;
  justify-content: center;
  gap: 40px;
  margin-top: 12px;
  height: 30px;
}

.broadcast-line {
  width: 2px;
  height: 0;
  background: linear-gradient(180deg, #fbbf24, transparent);
  transition: height 0.5s ease;
  transition-delay: var(--delay);
}

.broadcast-line.active {
  height: 30px;
}

/* 阶段指示器 */
.phase-indicators {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0;
  padding: 20px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 12px;
  margin-bottom: 20px;
}

.phase {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  opacity: 0.4;
  transition: all 0.4s ease;
}

.phase.active {
  opacity: 1;
}

.phase-dot {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  color: #a0aec0;
  font-weight: 600;
  transition: all 0.4s ease;
}

.phase.active .phase-dot {
  background: linear-gradient(135deg, #3182ce, #2c5282);
  color: white;
}

.phase.completed .phase-dot {
  background: linear-gradient(135deg, #38a169, #2f855a);
  color: white;
}

.phase-label {
  font-size: 0.75rem;
  color: #a0aec0;
  text-align: center;
  max-width: 80px;
}

.phase-line {
  width: 60px;
  height: 2px;
  background: rgba(255, 255, 255, 0.1);
  margin: 0 8px;
  margin-bottom: 24px;
  transition: all 0.4s ease;
}

.phase-line.active {
  background: linear-gradient(90deg, #38a169, #3182ce);
}

/* 日志区域 */
.consensus-logs {
  background: #0d1117;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 20px;
}

.logs-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.05);
  color: #a0aec0;
  font-size: 0.9rem;
}

.logs-content {
  padding: 12px 16px;
  max-height: 150px;
  overflow-y: auto;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 0.8rem;
}

.log-entry {
  display: flex;
  gap: 12px;
  padding: 4px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.log-entry:last-child {
  border-bottom: none;
}

.log-time {
  color: #718096;
  min-width: 70px;
}

.log-message {
  color: #e2e8f0;
}

.log-entry.success .log-message {
  color: #68d391;
}

.log-entry.error .log-message {
  color: #fc8181;
}

.log-entry.warning .log-message {
  color: #fbd38d;
}

/* 完成提示 */
.completion-banner {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: linear-gradient(135deg, rgba(56, 161, 105, 0.2), rgba(47, 133, 90, 0.1));
  border-radius: 12px;
  border: 1px solid rgba(56, 161, 105, 0.3);
  animation: fadeInUp 0.5s ease;
}

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

.completion-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #38a169, #2f855a);
  color: white;
  font-size: 1.5rem;
}

.completion-text h4 {
  margin: 0 0 4px 0;
  color: #68d391;
  font-size: 1.1rem;
}

.completion-text p {
  margin: 0;
  color: #a0aec0;
  font-size: 0.9rem;
}

/* 底部 */
.consensus-footer {
  padding: 16px 24px;
  background: rgba(255, 255, 255, 0.03);
  border-top: 1px solid #3d4a5c;
  display: flex;
  justify-content: flex-end;
}

.btn {
  padding: 12px 24px;
  border-radius: 10px;
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-primary {
  background: linear-gradient(135deg, #38a169, #2f855a);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(56, 161, 105, 0.4);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 响应式 */
@media (max-width: 768px) {
  .consensus-modal {
    width: 98%;
    max-height: 95vh;
  }

  .upload-info {
    flex-direction: column;
    gap: 12px;
  }

  .peer-nodes-grid {
    grid-template-columns: 1fr;
  }

  .phase-indicators {
    flex-wrap: wrap;
    gap: 8px;
  }

  .phase-line {
    width: 30px;
  }
}
</style>
