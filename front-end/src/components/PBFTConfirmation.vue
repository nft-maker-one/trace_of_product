<template>
  <div class="pbft-confirmation" v-if="visible">
    <div class="pbft-overlay" @click="handleClose"></div>
    <div class="pbft-modal">
      <div class="pbft-header">
        <h3>
          <i class="fas fa-network-wired"></i>
          {{ t('pbft.title') }}
        </h3>
        <button class="close-btn" @click="handleClose">
          <i class="fas fa-times"></i>
        </button>
      </div>

      <div class="pbft-content">
        <!-- 交易信息 -->
        <div class="transaction-info">
          <div class="info-item">
            <span class="label">{{ t('pbft.productId') }}:</span>
            <span class="value">{{ transactionData?.eggplant_id || 'N/A' }}</span>
          </div>
          <div class="info-item">
            <span class="label">{{ t('pbft.uploadNode') }}:</span>
            <span class="value node-highlight">
              {{ matchedNode?.id || t('pbft.unknown') }}
            </span>
          </div>
          <div class="info-item">
            <span class="label">{{ t('pbft.publicKey') }}:</span>
            <span class="value pubkey">{{ truncateKey(transactionData?.public_key) }}</span>
          </div>
        </div>

        <!-- PBFT 阶段动画 -->
        <div class="pbft-stages">
          <div class="stage-timeline">
            <!-- Pre-Prepare 阶段 -->
            <div class="stage" :class="{ active: currentStage >= 1, completed: currentStage > 1 }">
              <div class="stage-icon">
                <i class="fas fa-paper-plane"></i>
              </div>
              <div class="stage-content">
                <h4>{{ t('pbft.prePrepare') }}</h4>
                <p>{{ t('pbft.prePrepareDesc') }}</p>
                <div class="stage-animation" v-if="currentStage === 1">
                  <div class="pulse-ring"></div>
                </div>
              </div>
              <div class="stage-status">
                <i v-if="currentStage > 1" class="fas fa-check-circle"></i>
                <div v-else-if="currentStage === 1" class="spinner-small"></div>
              </div>
            </div>

            <!-- Prepare 阶段 -->
            <div class="stage" :class="{ active: currentStage >= 2, completed: currentStage > 2 }">
              <div class="stage-icon">
                <i class="fas fa-sync-alt"></i>
              </div>
              <div class="stage-content">
                <h4>{{ t('pbft.prepare') }}</h4>
                <p>{{ t('pbft.prepareDesc') }}</p>
                <div class="node-votes" v-if="currentStage >= 2">
                  <div 
                    v-for="(node, index) in participatingNodes" 
                    :key="node.id"
                    class="vote-badge"
                    :class="{ voted: prepareVotes.includes(node.id), 'vote-animating': animatingPrepareVotes.includes(node.id) }"
                    :style="{ animationDelay: `${index * 0.15}s` }"
                  >
                    <span class="node-mini-id">{{ formatNodeId(node.id) }}</span>
                    <i v-if="prepareVotes.includes(node.id)" class="fas fa-check"></i>
                  </div>
                </div>
              </div>
              <div class="stage-status">
                <span class="vote-count" v-if="currentStage >= 2">
                  {{ prepareVotes.length }}/{{ participatingNodes.length }}
                </span>
              </div>
            </div>

            <!-- Commit 阶段 -->
            <div class="stage" :class="{ active: currentStage >= 3, completed: currentStage > 3 }">
              <div class="stage-icon">
                <i class="fas fa-handshake"></i>
              </div>
              <div class="stage-content">
                <h4>{{ t('pbft.commit') }}</h4>
                <p>{{ t('pbft.commitDesc') }}</p>
                <div class="node-votes" v-if="currentStage >= 3">
                  <div 
                    v-for="(node, index) in participatingNodes" 
                    :key="node.id"
                    class="vote-badge commit"
                    :class="{ voted: commitVotes.includes(node.id), 'vote-animating': animatingCommitVotes.includes(node.id) }"
                    :style="{ animationDelay: `${index * 0.15}s` }"
                  >
                    <span class="node-mini-id">{{ formatNodeId(node.id) }}</span>
                    <i v-if="commitVotes.includes(node.id)" class="fas fa-check"></i>
                  </div>
                </div>
              </div>
              <div class="stage-status">
                <span class="vote-count" v-if="currentStage >= 3">
                  {{ commitVotes.length }}/{{ participatingNodes.length }}
                </span>
              </div>
            </div>

            <!-- Reply 阶段 -->
            <div class="stage" :class="{ active: currentStage >= 4, completed: currentStage >= 4 }">
              <div class="stage-icon">
                <i class="fas fa-check-double"></i>
              </div>
              <div class="stage-content">
                <h4>{{ t('pbft.reply') }}</h4>
                <p>{{ t('pbft.replyDesc') }}</p>
              </div>
              <div class="stage-status">
                <i v-if="currentStage >= 4" class="fas fa-check-circle success"></i>
              </div>
            </div>
          </div>
        </div>

        <!-- 节点网络可视化 -->
        <div class="network-visualization">
          <h4>
            <i class="fas fa-project-diagram"></i>
            {{ t('pbft.networkView') }}
          </h4>
          <div class="network-container">
            <div class="center-node" :class="{ active: currentStage >= 1 }">
              <div class="node-circle leader">
                <i class="fas fa-crown"></i>
              </div>
              <span class="node-label">{{ t('pbft.leader') }}</span>
              <span class="node-id">{{ formatNodeId(matchedNode?.id) }}</span>
            </div>
            
            <div class="peer-nodes">
              <div 
                v-for="(node, index) in otherNodes" 
                :key="node.id"
                class="peer-node"
                :class="{ 
                  active: currentStage >= 2,
                  confirmed: commitVotes.includes(node.id)
                }"
                :style="{ 
                  '--angle': `${(index * 360 / otherNodes.length)}deg`,
                  '--delay': `${index * 0.1}s`
                }"
              >
                <div class="connection-line" :class="{ active: currentStage >= 2 }"></div>
                <div class="node-circle">
                  <i class="fas fa-server"></i>
                </div>
                <span class="node-id">{{ formatNodeId(node.id) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 确认结果 -->
        <div class="confirmation-result" v-if="currentStage >= 4">
          <div class="result-icon success">
            <i class="fas fa-shield-alt"></i>
          </div>
          <h4>{{ t('pbft.consensusReached') }}</h4>
          <p>{{ t('pbft.consensusDesc', { nodes: participatingNodes.length }) }}</p>
        </div>
      </div>

      <div class="pbft-footer">
        <button class="btn btn-outline" @click="restartAnimation" v-if="currentStage >= 4">
          <i class="fas fa-redo"></i>
          {{ t('pbft.replay') }}
        </button>
        <button class="btn btn-primary" @click="handleClose">
          {{ t('pbft.close') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  transactionData: {
    type: Object,
    default: null
  },
  nodes: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['close'])

const { t } = useI18n()

// 状态
const currentStage = ref(0)
const prepareVotes = ref([])
const commitVotes = ref([])
const animatingPrepareVotes = ref([])
const animatingCommitVotes = ref([])

// 计算属性：匹配的上传节点（通过公钥比对）
const matchedNode = computed(() => {
  if (!props.transactionData?.public_key || !props.nodes.length) return null
  
  // 通过公钥匹配节点
  const pubKey = props.transactionData.public_key
  return props.nodes.find(node => node.pub_key === pubKey) || props.nodes[0]
})

// 参与共识的节点
const participatingNodes = computed(() => {
  return props.nodes.slice(0, 4) // PBFT 至少需要 3f+1 个节点
})

// 其他节点（非leader）
const otherNodes = computed(() => {
  if (!matchedNode.value) return props.nodes.slice(1, 4)
  return props.nodes.filter(n => n.id !== matchedNode.value.id).slice(0, 3)
})

// 工具函数
function truncateKey(key) {
  if (!key) return 'N/A'
  if (key.length <= 20) return key
  return `${key.substring(0, 10)}...${key.substring(key.length - 10)}`
}

function formatNodeId(id) {
  if (!id) return 'N/A'
  const idStr = String(id)
  if (idStr.length <= 6) return idStr
  return `${idStr.substring(0, 4)}...`
}

// 动画控制
let animationTimer = null

function startAnimation() {
  currentStage.value = 0
  prepareVotes.value = []
  commitVotes.value = []
  animatingPrepareVotes.value = []
  animatingCommitVotes.value = []

  // 阶段1: Pre-Prepare
  setTimeout(() => {
    currentStage.value = 1
  }, 500)

  // 阶段2: Prepare - 节点逐个投票
  setTimeout(() => {
    currentStage.value = 2
    animatePrepareVotes()
  }, 2000)

  // 阶段3: Commit - 节点逐个确认
  setTimeout(() => {
    currentStage.value = 3
    animateCommitVotes()
  }, 4500)

  // 阶段4: Reply - 完成
  setTimeout(() => {
    currentStage.value = 4
  }, 7000)
}

function animatePrepareVotes() {
  const nodes = participatingNodes.value
  nodes.forEach((node, index) => {
    setTimeout(() => {
      animatingPrepareVotes.value.push(node.id)
      setTimeout(() => {
        prepareVotes.value.push(node.id)
        animatingPrepareVotes.value = animatingPrepareVotes.value.filter(id => id !== node.id)
      }, 300)
    }, index * 400)
  })
}

function animateCommitVotes() {
  const nodes = participatingNodes.value
  nodes.forEach((node, index) => {
    setTimeout(() => {
      animatingCommitVotes.value.push(node.id)
      setTimeout(() => {
        commitVotes.value.push(node.id)
        animatingCommitVotes.value = animatingCommitVotes.value.filter(id => id !== node.id)
      }, 300)
    }, index * 400)
  })
}

function restartAnimation() {
  startAnimation()
}

function handleClose() {
  emit('close')
}

// 监听可见性变化
watch(() => props.visible, (newVal) => {
  if (newVal) {
    startAnimation()
  } else {
    // 重置状态
    currentStage.value = 0
    prepareVotes.value = []
    commitVotes.value = []
  }
})

onMounted(() => {
  if (props.visible) {
    startAnimation()
  }
})

onUnmounted(() => {
  if (animationTimer) {
    clearTimeout(animationTimer)
  }
})
</script>

<style scoped>
.pbft-confirmation {
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

.pbft-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
}

.pbft-modal {
  position: relative;
  background: white;
  border-radius: 20px;
  width: 90%;
  max-width: 800px;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 25px 80px rgba(0, 0, 0, 0.3);
  animation: modalSlideIn 0.4s ease-out;
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

.pbft-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  background: linear-gradient(135deg, #2e7d32 0%, #1b5e20 100%);
  color: white;
}

.pbft-header h3 {
  margin: 0;
  font-size: 1.25rem;
  display: flex;
  align-items: center;
  gap: 10px;
}

.close-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: rotate(90deg);
}

.pbft-content {
  padding: 24px;
  max-height: 60vh;
  overflow-y: auto;
}

/* 交易信息 */
.transaction-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 12px;
  margin-bottom: 24px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-item .label {
  font-size: 0.85rem;
  color: #78909c;
}

.info-item .value {
  font-size: 1rem;
  font-weight: 600;
  color: #263238;
}

.info-item .value.node-highlight {
  color: #2e7d32;
}

.info-item .value.pubkey {
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 0.85rem;
  word-break: break-all;
}

/* PBFT 阶段 */
.pbft-stages {
  margin-bottom: 24px;
}

.stage-timeline {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stage {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 16px;
  border-radius: 12px;
  background: #f9fafb;
  border: 2px solid transparent;
  opacity: 0.5;
  transition: all 0.4s ease;
}

.stage.active {
  opacity: 1;
  background: white;
  border-color: #e8f5e9;
  box-shadow: 0 4px 15px rgba(46, 125, 50, 0.1);
}

.stage.completed {
  opacity: 1;
  background: #e8f5e9;
  border-color: #4caf50;
}

.stage-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: linear-gradient(135deg, #e0e0e0, #bdbdbd);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 1.2rem;
  flex-shrink: 0;
  transition: all 0.4s ease;
}

.stage.active .stage-icon {
  background: linear-gradient(135deg, #4caf50, #2e7d32);
  animation: pulse 2s infinite;
}

.stage.completed .stage-icon {
  background: linear-gradient(135deg, #2e7d32, #1b5e20);
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

.stage-content {
  flex: 1;
}

.stage-content h4 {
  margin: 0 0 4px 0;
  font-size: 1rem;
  color: #263238;
}

.stage-content p {
  margin: 0;
  font-size: 0.85rem;
  color: #78909c;
}

.stage-status {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 60px;
}

.stage-status .fa-check-circle {
  color: #4caf50;
  font-size: 1.5rem;
}

.stage-status .fa-check-circle.success {
  color: #2e7d32;
  animation: bounceIn 0.5s ease;
}

@keyframes bounceIn {
  0% { transform: scale(0); }
  50% { transform: scale(1.2); }
  100% { transform: scale(1); }
}

.spinner-small {
  width: 24px;
  height: 24px;
  border: 3px solid #e0e0e0;
  border-top-color: #4caf50;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 节点投票 */
.node-votes {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.vote-badge {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #eceff1;
  border-radius: 20px;
  font-size: 0.8rem;
  transition: all 0.3s ease;
}

.vote-badge.vote-animating {
  animation: votePopIn 0.3s ease-out;
}

@keyframes votePopIn {
  0% { transform: scale(0.5); opacity: 0; }
  100% { transform: scale(1); opacity: 1; }
}

.vote-badge.voted {
  background: linear-gradient(135deg, #4caf50, #2e7d32);
  color: white;
}

.vote-badge.commit.voted {
  background: linear-gradient(135deg, #0288d1, #01579b);
}

.vote-badge .fa-check {
  font-size: 0.7rem;
}

.vote-count {
  font-size: 0.9rem;
  font-weight: 600;
  color: #2e7d32;
}

/* 网络可视化 */
.network-visualization {
  padding: 20px;
  background: #f5f7fa;
  border-radius: 16px;
  margin-bottom: 24px;
}

.network-visualization h4 {
  margin: 0 0 20px 0;
  font-size: 1rem;
  color: #263238;
  display: flex;
  align-items: center;
  gap: 8px;
}

.network-container {
  position: relative;
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.center-node {
  position: absolute;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  z-index: 2;
  opacity: 0.5;
  transition: all 0.5s ease;
}

.center-node.active {
  opacity: 1;
}

.node-circle {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: linear-gradient(135deg, #90a4ae, #607d8b);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 1.2rem;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
  transition: all 0.4s ease;
}

.node-circle.leader {
  width: 70px;
  height: 70px;
  background: linear-gradient(135deg, #ffd54f, #ff8f00);
  box-shadow: 0 6px 20px rgba(255, 143, 0, 0.4);
}

.center-node.active .node-circle.leader {
  animation: leaderPulse 2s infinite;
}

@keyframes leaderPulse {
  0%, 100% { 
    box-shadow: 0 6px 20px rgba(255, 143, 0, 0.4);
  }
  50% { 
    box-shadow: 0 6px 30px rgba(255, 143, 0, 0.6), 0 0 40px rgba(255, 143, 0, 0.3);
  }
}

.node-label {
  font-size: 0.75rem;
  font-weight: 600;
  color: #ff8f00;
  text-transform: uppercase;
}

.node-id {
  font-size: 0.7rem;
  color: #78909c;
  font-family: 'Monaco', 'Consolas', monospace;
}

.peer-nodes {
  position: absolute;
  width: 100%;
  height: 100%;
}

.peer-node {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: rotate(var(--angle)) translateX(100px) rotate(calc(-1 * var(--angle)));
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  opacity: 0.5;
  transition: all 0.5s ease;
  transition-delay: var(--delay);
}

.peer-node.active {
  opacity: 1;
}

.peer-node.confirmed .node-circle {
  background: linear-gradient(135deg, #4caf50, #2e7d32);
  box-shadow: 0 4px 15px rgba(46, 125, 50, 0.4);
}

.peer-node .node-circle {
  width: 50px;
  height: 50px;
  font-size: 1rem;
}

.connection-line {
  position: absolute;
  width: 60px;
  height: 2px;
  background: linear-gradient(90deg, transparent, #bdbdbd);
  transform: rotate(calc(180deg + var(--angle))) translateX(30px);
  transform-origin: right center;
  opacity: 0;
  transition: all 0.5s ease;
}

.peer-node.active .connection-line {
  opacity: 1;
  background: linear-gradient(90deg, transparent, #4caf50);
}

/* 确认结果 */
.confirmation-result {
  text-align: center;
  padding: 24px;
  background: linear-gradient(135deg, #e8f5e9, #c8e6c9);
  border-radius: 16px;
  animation: resultFadeIn 0.5s ease-out;
}

@keyframes resultFadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.result-icon {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  margin: 0 auto 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
}

.result-icon.success {
  background: linear-gradient(135deg, #4caf50, #2e7d32);
  color: white;
  box-shadow: 0 8px 25px rgba(46, 125, 50, 0.3);
}

.confirmation-result h4 {
  margin: 0 0 8px 0;
  font-size: 1.25rem;
  color: #1b5e20;
}

.confirmation-result p {
  margin: 0;
  color: #2e7d32;
}

/* 底部按钮 */
.pbft-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  background: #f9fafb;
  border-top: 1px solid #e0e0e0;
}

.btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border-radius: 10px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
}

.btn-primary {
  background: linear-gradient(135deg, #4caf50, #2e7d32);
  color: white;
}

.btn-primary:hover {
  background: linear-gradient(135deg, #2e7d32, #1b5e20);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(46, 125, 50, 0.3);
}

.btn-outline {
  background: white;
  color: #2e7d32;
  border: 2px solid #4caf50;
}

.btn-outline:hover {
  background: #e8f5e9;
}

/* 响应式 */
@media (max-width: 600px) {
  .pbft-modal {
    width: 95%;
    max-height: 95vh;
  }

  .transaction-info {
    grid-template-columns: 1fr;
  }

  .stage {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .network-container {
    height: 250px;
  }

  .peer-node {
    transform: rotate(var(--angle)) translateX(80px) rotate(calc(-1 * var(--angle)));
  }
}
</style>
