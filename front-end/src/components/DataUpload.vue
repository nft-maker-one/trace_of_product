<template>
  <div class="data-upload">
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">
          <i class="fas fa-upload"></i>
          {{ t('dataUpload.title') }}
        </h3>
        <div style="font-size: 14px; color: var(--gray);">
          <i class="fas fa-info-circle"></i>
          {{ t('dataUpload.description') }}
        </div>
      </div>

      <!-- eggplant ID 单独一行 -->
      <div class="form-row">
        <div class="form-col">
          <div class="form-group">
            <label class="form-label" for="upload-eggplant-id">
              <i class="fas fa-fingerprint"></i>
              {{ t('dataUpload.eggplantId') }} <span>*</span>
            </label>
            <input
              type="text"
              id="upload-eggplant-id"
              v-model="formData.eggplant_id"
              class="form-control"
              :placeholder="t('dataUpload.enterProductId')"
            >
          </div>
        </div>
        <div class="form-col"></div>
      </div>

      <!-- 产品高度 + 产品哈希 -->
      <div class="form-row">
        <div class="form-col">
          <div class="form-group">
            <label class="form-label" for="upload-product-height">
              <i class="fas fa-seedling"></i>
              {{ t('dataUpload.productHeight') }}
            </label>
            <input
              type="number"
              id="upload-product-height"
              v-model="formData.product_height"
              class="form-control"
              :placeholder="t('dataUpload.productHeight')"
            >
          </div>
        </div>
        <div class="form-col">
          <div class="form-group">
            <label class="form-label" for="upload-product-hash">
              <i class="fas fa-hashtag"></i>
              {{ t('dataUpload.productHash') }} <span>*</span>
            </label>
            <div class="input-with-button">
              <input
                type="text"
                id="upload-product-hash"
                v-model="formData.product_hash"
                class="form-control hash-readonly"
                :placeholder="t('dataUpload.productHash')"
                readonly
              >
              <button type="button" class="btn btn-generate" @click="generateHashForField('product_hash')" :title="t('dataUpload.generateHash')">
                <i class="fas fa-random"></i>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 运输高度 + 运输哈希 -->
      <div class="form-row">
        <div class="form-col">
          <div class="form-group">
            <label class="form-label" for="upload-transport-height">
              <i class="fas fa-truck"></i>
              {{ t('dataUpload.transportHeight') }}
            </label>
            <input
              type="number"
              id="upload-transport-height"
              v-model="formData.transport_height"
              class="form-control"
              :placeholder="t('dataUpload.transportHeight')"
            >
          </div>
        </div>
        <div class="form-col">
          <div class="form-group">
            <label class="form-label" for="upload-transport-hash">
              <i class="fas fa-hashtag"></i>
              {{ t('dataUpload.transportHash') }} <span>*</span>
            </label>
            <div class="input-with-button">
              <input
                type="text"
                id="upload-transport-hash"
                v-model="formData.transport_hash"
                class="form-control hash-readonly"
                :placeholder="t('dataUpload.transportHash')"
                readonly
              >
              <button type="button" class="btn btn-generate" @click="generateHashForField('transport_hash')" :title="t('dataUpload.generateHash')">
                <i class="fas fa-random"></i>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 加工高度 + 加工哈希 -->
      <div class="form-row">
        <div class="form-col">
          <div class="form-group">
            <label class="form-label" for="upload-process-height">
              <i class="fas fa-industry"></i>
              {{ t('dataUpload.processHeight') }}
            </label>
            <input
              type="number"
              id="upload-process-height"
              v-model="formData.process_height"
              class="form-control"
              :placeholder="t('dataUpload.processHeight')"
            >
          </div>
        </div>
        <div class="form-col">
          <div class="form-group">
            <label class="form-label" for="upload-process-hash">
              <i class="fas fa-hashtag"></i>
              {{ t('dataUpload.processHash') }} <span>*</span>
            </label>
            <div class="input-with-button">
              <input
                type="text"
                id="upload-process-hash"
                v-model="formData.process_hash"
                class="form-control hash-readonly"
                :placeholder="t('dataUpload.processHash')"
                readonly
              >
              <button type="button" class="btn btn-generate" @click="generateHashForField('process_hash')" :title="t('dataUpload.generateHash')">
                <i class="fas fa-random"></i>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 存储高度 + 存储哈希 -->
      <div class="form-row">
        <div class="form-col">
          <div class="form-group">
            <label class="form-label" for="upload-storage-height">
              <i class="fas fa-warehouse"></i>
              {{ t('dataUpload.storageHeight') }}
            </label>
            <input
              type="number"
              id="upload-storage-height"
              v-model="formData.storage_height"
              class="form-control"
              :placeholder="t('dataUpload.storageHeight')"
            >
          </div>
        </div>
        <div class="form-col">
          <div class="form-group">
            <label class="form-label" for="upload-storage-hash">
              <i class="fas fa-hashtag"></i>
              {{ t('dataUpload.storageHash') }} <span>*</span>
            </label>
            <div class="input-with-button">
              <input
                type="text"
                id="upload-storage-hash"
                v-model="formData.storage_hash"
                class="form-control hash-readonly"
                :placeholder="t('dataUpload.storageHash')"
                readonly
              >
              <button type="button" class="btn btn-generate" @click="generateHashForField('storage_hash')" :title="t('dataUpload.generateHash')">
                <i class="fas fa-random"></i>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 销售高度 + 销售哈希 -->
      <div class="form-row">
        <div class="form-col">
          <div class="form-group">
            <label class="form-label" for="upload-sell-height">
              <i class="fas fa-shopping-cart"></i>
              {{ t('dataUpload.sellHeight') }}
            </label>
            <input
              type="number"
              id="upload-sell-height"
              v-model="formData.sell_height"
              class="form-control"
              :placeholder="t('dataUpload.sellHeight')"
            >
          </div>
        </div>
        <div class="form-col">
          <div class="form-group">
            <label class="form-label" for="upload-sell-hash">
              <i class="fas fa-hashtag"></i>
              {{ t('dataUpload.sellHash') }} <span>*</span>
            </label>
            <div class="input-with-button">
              <input
                type="text"
                id="upload-sell-hash"
                v-model="formData.sell_hash"
                class="form-control hash-readonly"
                :placeholder="t('dataUpload.sellHash')"
                readonly
              >
              <button type="button" class="btn btn-generate" @click="generateHashForField('sell_hash')" :title="t('dataUpload.generateHash')">
                <i class="fas fa-random"></i>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 节点选择 + 上传按钮 -->
      <div class="form-row">
        <div class="form-col">
          <div class="select-wrapper">
            <label class="form-label" for="select-node">
              <i class="fas fa-server"></i>
              {{ t('dataUpload.selectUploadNode') }} *
            </label>
            <select id="select-node" v-model="selectedNode">
              <option value="">{{ t('dataUpload.selectNode') }}</option>
              <option v-for="node in nodes" :key="node.id" :value="node.id">
                {{ node.id }}
              </option>
            </select>
          </div>
        </div>
        <div class="form-col">
          <button class="btn btn-primary" @click="handleUpload" :disabled="uploading" style="margin-top: 28px; width: 100%;">
            <i class="fas fa-cloud-upload-alt"></i>
            {{ uploading ? t('dataUpload.uploading') : t('dataUpload.uploadData') }}
            <div v-if="uploading" class="spinner"></div>
          </button>
        </div>
      </div>

      <div v-if="uploadResult.message" :class="['alert', uploadResult.type]">
        <i :class="uploadResult.icon"></i>
        <span>{{ uploadResult.message }}</span>
      </div>
    </div>

    <div class="card">
      <div class="card-header">
        <h3 class="card-title">
          <i class="fas fa-server"></i>
          {{ t('dataUpload.onlineNodeStatus') }}
        </h3>
        <button class="btn btn-outline" @click="toggleNodeList">
          <i :class="showNodeList ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
          {{ showNodeList ? t('dataUpload.hideNodeList') : t('dataUpload.viewBlockchainNodes') }}
        </button>
      </div>

      <div class="node-list-container">
        <div v-if="!showNodeList" class="empty-state">
          <i class="fas fa-server"></i>
          <h3>{{ t('dataUpload.clickToViewNodes') }}</h3>
          <p>{{ t('dataUpload.nodeInfoHere') }}</p>
        </div>

        <div v-else-if="loadingNodes" class="empty-state">
          <div class="spinner"></div>
          <h3>{{ t('dataUpload.loadingNodes') }}</h3>
        </div>

        <div v-else-if="nodes.length === 0" class="empty-state">
          <i class="fas fa-exclamation-circle"></i>
          <h3>{{ t('dataUpload.noNodeInfo') }}</h3>
          <p>{{ t('dataUpload.cannotGetNodeData') }}</p>
        </div>

        <table v-else class="node-table">
          <thead>
            <tr>
              <th>{{ t('dataUpload.nodeId') }}</th>
              <th>{{ t('dataUpload.address') }}</th>
              <th>{{ t('dataUpload.registerTime') }}</th>
              <th>{{ t('dataUpload.verifyMessageCount') }}</th>
              <th>{{ t('dataUpload.status') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="node in nodes" :key="node.id">
              <td>{{ node.id }}</td>
              <td>{{ node.addr }}</td>
              <td>{{ formatTime(node.create_time) }}</td>
              <td>{{ node.verify_time }}</td>
              <td>
                <span class="node-status status-online">{{ t('dataUpload.online') }}</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '../services/api'

const { t } = useI18n()

const showNodeList = ref(false)
const loadingNodes = ref(false)
const uploading = ref(false)
const selectedNode = ref('')
const nodes = ref([])

// 简单的哈希函数：将字符串转换为64位十六进制数（模拟SHA-256）
async function sha256(message) {
  // 使用 Web Crypto API 进行真正的 SHA-256 哈希
  const msgBuffer = new TextEncoder().encode(message)
  const hashBuffer = await crypto.subtle.digest('SHA-256', msgBuffer)
  const hashArray = Array.from(new Uint8Array(hashBuffer))
  const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
  return hashHex // 返回64位十六进制字符串
}

// 根据高度值生成哈希（加入随机因素，相同输入产生不同输出）
async function generateHashFromHeight(height) {
  if (!height && height !== 0) return ''
  
  const heightNum = parseInt(height) || 0
  // 在高度值基础上加入时间戳和随机数，确保每次生成不同的哈希
  const randomSalt = Date.now().toString() + Math.random().toString(36).substring(2)
  const dataToHash = heightNum.toString() + '_' + randomSalt
  
  // 使用 SHA-256 对组合数据进行哈希，生成64位十六进制数
  const hash = await sha256(dataToHash)
  return hash
}

// 生成64位16进制随机哈希值（保留用于其他用途）
function generateRandomHash() {
  const hexChars = '0123456789abcdef'
  let hash = ''
  for (let i = 0; i < 64; i++) {
    hash += hexChars[Math.floor(Math.random() * 16)]
  }
  return hash
}

// 为指定字段生成哈希值（现在从对应的高度字段生成）
async function generateHashForField(fieldName) {
  // 获取对应的高度字段名
  const heightFieldName = fieldName.replace('_hash', '_height')
  const heightValue = formData[heightFieldName]
  
  if (heightValue || heightValue === 0) {
    formData[fieldName] = await generateHashFromHeight(heightValue)
  } else {
    // 如果没有高度值，生成随机哈希
    formData[fieldName] = generateRandomHash()
  }
}

const formData = reactive({
  eggplant_id: '',
  product_height: '',
  product_hash: '',
  transport_height: '',
  transport_hash: '',
  process_height: '',
  process_hash: '',
  storage_height: '',
  storage_hash: '',
  sell_height: '',
  sell_hash: ''
})

// 监听高度值变化，自动生成对应的哈希
watch(() => formData.product_height, async (newVal) => {
  if (newVal || newVal === 0) {
    formData.product_hash = await generateHashFromHeight(newVal)
  }
})

watch(() => formData.transport_height, async (newVal) => {
  if (newVal || newVal === 0) {
    formData.transport_hash = await generateHashFromHeight(newVal)
  }
})

watch(() => formData.process_height, async (newVal) => {
  if (newVal || newVal === 0) {
    formData.process_hash = await generateHashFromHeight(newVal)
  }
})

watch(() => formData.storage_height, async (newVal) => {
  if (newVal || newVal === 0) {
    formData.storage_hash = await generateHashFromHeight(newVal)
  }
})

watch(() => formData.sell_height, async (newVal) => {
  if (newVal || newVal === 0) {
    formData.sell_hash = await generateHashFromHeight(newVal)
  }
})

const uploadResult = reactive({
  message: '',
  type: '',
  icon: ''
})

// formFields 已不再使用，表单字段现在直接在模板中定义

async function toggleNodeList() {
  showNodeList.value = !showNodeList.value

  if (showNodeList.value && nodes.value.length === 0) {
    await fetchNodes()
  }
}

async function fetchNodes() {
  loadingNodes.value = true
  try {
    const data = await api.getNodes()
    if (Array.isArray(data)) {
      nodes.value = data
      // 更新节点选择器
      if (nodes.value.length > 0 && !selectedNode.value) {
        selectedNode.value = nodes.value[0].id
      }
    }
  } catch (error) {
    console.error('获取节点失败:', error)
  } finally {
    loadingNodes.value = false
  }
}

async function handleUpload() {
  // 验证必填字段
  if (!formData.eggplant_id || !formData.product_hash || !formData.transport_hash ||
      !formData.process_hash || !formData.storage_hash || !formData.sell_hash || !selectedNode.value) {
    showUploadResult(t('dataUpload.fillAllRequired'), 'error', 'fas fa-exclamation-circle')
    return
  }

  uploading.value = true

  try {
    const selectedNodeObj = nodes.value.find(node => node.id === selectedNode.value)
    if (!selectedNodeObj) {
      throw new Error('选择的节点不存在')
    }

    const uploadData = {
      ...formData,
      eggplant_id: parseInt(formData.eggplant_id),
      product_height: parseInt(formData.product_height) || 0,
      transport_height: parseInt(formData.transport_height) || 0,
      process_height: parseInt(formData.process_height) || 0,
      storage_height: parseInt(formData.storage_height) || 0,
      sell_height: parseInt(formData.sell_height) || 0,
      node_ip: selectedNodeObj.addr
    }
    const response = await api.uploadData(uploadData)

    if (response.status === 'ok' && response.msg === '上传成功') {
      showUploadResult(t('dataUpload.uploadSuccess'), 'success', 'fas fa-check-circle')
      // 清空表单
      Object.keys(formData).forEach(key => {
        formData[key] = ''
      })
    } else {
      showUploadResult(response.msg || t('dataUpload.uploadFailed'), 'warning', 'fas fa-exclamation-triangle')
    }
  } catch (error) {
    console.error('上传失败:', error)
    showUploadResult(t('dataUpload.uploadFailedWithError', { error: error.message }), 'error', 'fas fa-times-circle')
  } finally {
    uploading.value = false
  }
}

function showUploadResult(message, type, icon) {
  uploadResult.message = message
  uploadResult.type = type
  uploadResult.icon = icon

  // 3秒后清除消息
  setTimeout(() => {
    uploadResult.message = ''
  }, 3000)
}

function formatTime(timestamp) {
  if (!timestamp) return t('dataUpload.unknown')
  // create_time 是纳秒时间戳，需要转换为毫秒
  const ms = typeof timestamp === 'number' ? Math.floor(timestamp / 1000000) : timestamp
  return new Date(ms).toLocaleString()
}

onMounted(() => {
  // 初始化时加载节点
  fetchNodes()
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

.data-upload {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
}

/* 卡片样式优化 */
.card {
  background: white;
  border-radius: 16px;
  padding: 30px;
  margin-bottom: 30px;
  box-shadow: var(--card-shadow);
  transition: var(--transition);
  border: 1px solid var(--light-gray);
}

.card:hover {
  box-shadow: var(--card-shadow-hover);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid var(--light-gray);
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

/* 表单样式优化 */
.form-row {
  display: flex;
  gap: 20px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.form-col {
  flex: 1;
  min-width: 280px;
}

.form-group {
  margin-bottom: 12px;
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
  width: 20px;
  color: var(--primary);
}

.form-label span {
  color: var(--error);
}

.form-control {
  width: 100%;
  padding: 10px 14px;
  border: 2px solid var(--light-gray);
  border-radius: 8px;
  font-size: 14px;
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

/* 输入框带按钮样式 */
.input-with-button {
  display: flex;
  gap: 8px;
  align-items: stretch;
}

.input-with-button .form-control {
  flex: 1;
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
}

.btn-generate {
  padding: 14px 16px;
  background: linear-gradient(135deg, var(--secondary) 0%, var(--secondary-dark) 100%);
  color: white;
  border: none;
  border-radius: 0 10px 10px 0;
  cursor: pointer;
  transition: var(--transition);
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 50px;
}

.btn-generate:hover {
  background: linear-gradient(135deg, var(--secondary-dark) 0%, var(--secondary) 100%);
  transform: scale(1.02);
}

.btn-generate i {
  font-size: 16px;
}

/* 只读哈希输入框样式 */
.hash-readonly {
  background-color: #f5f5f5 !important;
  color: #666 !important;
  cursor: not-allowed;
  font-family: 'Courier New', monospace;
  font-size: 13px;
}

.hash-readonly:focus {
  border-color: var(--light-gray) !important;
  box-shadow: none !important;
}

/* 选择框样式优化 */
.select-wrapper {
  margin-bottom: 20px;
}

.select-wrapper select {
  width: 100%;
  padding: 14px 16px;
  border: 2px solid var(--light-gray);
  border-radius: 10px;
  font-size: 15px;
  background-color: white;
  color: var(--dark);
  cursor: pointer;
  transition: var(--transition);
  appearance: none;
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: right 16px center;
  background-size: 16px;
  padding-right: 45px;
}

.select-wrapper select:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(46, 125, 50, 0.1);
}

/* 按钮样式优化 */
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
}

.btn-outline:hover {
  background: var(--primary);
  color: white;
}

/* 警告框样式优化 */
.alert {
  padding: 18px 20px;
  border-radius: 12px;
  margin-top: 25px;
  display: flex;
  align-items: center;
  gap: 15px;
  animation: slideIn 0.3s ease;
  border-left: 4px solid;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.alert.error {
  background-color: rgba(244, 67, 54, 0.1);
  border-left-color: var(--error);
  color: var(--error);
}

.alert.success {
  background-color: rgba(76, 175, 80, 0.1);
  border-left-color: var(--success);
  color: var(--success);
}

.alert.warning {
  background-color: rgba(255, 152, 0, 0.1);
  border-left-color: var(--warning);
  color: var(--warning);
}

.alert i {
  font-size: 20px;
}

/* 节点列表容器 */
.node-list-container {
  margin-top: 20px;
}

/* 空状态样式 */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--gray);
}

.empty-state i {
  font-size: 64px;
  margin-bottom: 20px;
  color: #cfd8dc;
}

.empty-state h3 {
  font-size: 20px;
  margin-bottom: 10px;
  color: var(--dark);
}

.empty-state p {
  font-size: 15px;
  margin-bottom: 20px;
}

/* 节点表格样式 */
.node-table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  margin-top: 20px;
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.node-table thead {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
}

.node-table th {
  padding: 18px 16px;
  text-align: left;
  font-weight: 600;
  font-size: 15px;
  color: white;
  border: none;
}

.node-table tbody tr {
  transition: var(--transition);
}

.node-table tbody tr:hover {
  background-color: rgba(46, 125, 50, 0.03);
}

.node-table td {
  padding: 16px;
  color: var(--dark);
  border-bottom: 1px solid var(--light-gray);
}

.node-table tbody tr:last-child td {
  border-bottom: none;
}

/* 节点状态标签 */
.node-status {
  display: inline-flex;
  align-items: center;
  padding: 6px 15px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
  gap: 6px;
}

.node-status::before {
  content: '';
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}

.status-online {
  background-color: rgba(76, 175, 80, 0.1);
  color: var(--success);
}

.status-online::before {
  background-color: var(--success);
}

.status-offline {
  background-color: rgba(244, 67, 54, 0.1);
  color: var(--error);
}

.status-offline::before {
  background-color: var(--error);
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
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* 表单容器内间距 */
.data-upload .card:first-child {
  margin-bottom: 40px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .data-upload {
    padding: 15px;
  }

  .card {
    padding: 25px;
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
}

@media (max-width: 768px) {
  .data-upload {
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

  .node-table {
    font-size: 14px;
  }

  .node-table th,
  .node-table td {
    padding: 12px 10px;
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

  .alert {
    padding: 14px 16px;
    font-size: 14px;
  }
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

/* 焦点样式 */
:focus-visible {
  outline: 2px solid var(--primary);
  outline-offset: 2px;
}

/* 禁用状态样式 */
:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 工具类 */
.text-muted {
  color: var(--gray);
}

.text-success {
  color: var(--success);
}

.text-error {
  color: var(--error);
}

.text-warning {
  color: var(--warning);
}

.mt-3 {
  margin-top: 15px;
}

.mb-3 {
  margin-bottom: 15px;
}

.d-flex {
  display: flex;
}

.align-items-center {
  align-items: center;
}

.justify-content-between {
  justify-content: space-between;
}

.gap-2 {
  gap: 8px;
}

.gap-3 {
  gap: 12px;
}
</style>
