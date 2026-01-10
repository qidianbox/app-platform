<template>
  <div class="error-handler">
    <!-- 全局错误提示 -->
    <el-dialog
      v-model="showErrorDialog"
      title="操作失败"
      width="500px"
      :close-on-click-modal="false"
    >
      <div class="error-content">
        <el-icon class="error-icon" :size="48">
          <CircleClose />
        </el-icon>
        <div class="error-info">
          <h3>{{ errorTitle }}</h3>
          <p class="error-message">{{ errorMessage }}</p>
          <div v-if="errorDetails" class="error-details">
            <el-collapse>
              <el-collapse-item title="错误详情" name="details">
                <pre>{{ errorDetails }}</pre>
              </el-collapse-item>
            </el-collapse>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button v-if="canRetry" type="primary" @click="handleRetry">
            <el-icon><Refresh /></el-icon>
            重试
          </el-button>
          <el-button @click="handleClose">关闭</el-button>
          <el-button type="info" @click="exportLogs">
            <el-icon><Download /></el-icon>
            导出日志
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 网络状态提示 -->
    <el-alert
      v-if="!isOnline"
      title="网络连接已断开"
      type="error"
      :closable="false"
      show-icon
      class="network-alert"
    >
      <template #default>
        请检查您的网络连接，系统将在网络恢复后自动重试
      </template>
    </el-alert>

    <!-- 重试进度提示 -->
    <el-dialog
      v-model="showRetryDialog"
      title="正在重试"
      width="400px"
      :close-on-click-modal="false"
      :show-close="false"
    >
      <div class="retry-content">
        <el-progress
          :percentage="retryProgress"
          :stroke-width="10"
          :format="retryProgressFormat"
        />
        <p class="retry-message">
          正在进行第 {{ currentRetry }} 次重试，共 {{ maxRetries }} 次
        </p>
        <p class="retry-countdown" v-if="retryCountdown > 0">
          {{ retryCountdown }} 秒后重试...
        </p>
      </div>
      <template #footer>
        <el-button @click="cancelRetry">取消重试</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { CircleClose, Refresh, Download } from '@element-plus/icons-vue'
import { logger } from '@/utils/request'

// 状态
const showErrorDialog = ref(false)
const showRetryDialog = ref(false)
const isOnline = ref(navigator.onLine)
const errorTitle = ref('')
const errorMessage = ref('')
const errorDetails = ref('')
const canRetry = ref(false)
const retryCallback = ref(null)

// 重试相关
const currentRetry = ref(0)
const maxRetries = ref(3)
const retryCountdown = ref(0)
const retryTimer = ref(null)

// 计算属性
const retryProgress = computed(() => {
  return Math.round((currentRetry.value / maxRetries.value) * 100)
})

// 方法
const retryProgressFormat = (percentage) => {
  return `${currentRetry.value}/${maxRetries.value}`
}

const showError = (options) => {
  errorTitle.value = options.title || '操作失败'
  errorMessage.value = options.message || '发生未知错误'
  errorDetails.value = options.details || ''
  canRetry.value = options.canRetry || false
  retryCallback.value = options.onRetry || null
  showErrorDialog.value = true
}

const handleClose = () => {
  showErrorDialog.value = false
  errorTitle.value = ''
  errorMessage.value = ''
  errorDetails.value = ''
  canRetry.value = false
  retryCallback.value = null
}

const handleRetry = () => {
  if (retryCallback.value) {
    showErrorDialog.value = false
    retryCallback.value()
  }
}

const startRetry = (options) => {
  currentRetry.value = options.currentRetry || 1
  maxRetries.value = options.maxRetries || 3
  retryCountdown.value = options.delay || 3
  showRetryDialog.value = true
  
  retryTimer.value = setInterval(() => {
    retryCountdown.value--
    if (retryCountdown.value <= 0) {
      clearInterval(retryTimer.value)
      showRetryDialog.value = false
      if (options.onRetry) {
        options.onRetry()
      }
    }
  }, 1000)
}

const cancelRetry = () => {
  if (retryTimer.value) {
    clearInterval(retryTimer.value)
  }
  showRetryDialog.value = false
}

const exportLogs = () => {
  if (window.__exportLogs) {
    window.__exportLogs()
  }
}

// 网络状态监听
const handleOnline = () => {
  isOnline.value = true
  logger.info('Network', 'Network connection restored')
}

const handleOffline = () => {
  isOnline.value = false
  logger.warn('Network', 'Network connection lost')
}

onMounted(() => {
  window.addEventListener('online', handleOnline)
  window.addEventListener('offline', handleOffline)
  
  // 暴露方法到全局
  window.__showError = showError
  window.__startRetry = startRetry
})

onUnmounted(() => {
  window.removeEventListener('online', handleOnline)
  window.removeEventListener('offline', handleOffline)
  if (retryTimer.value) {
    clearInterval(retryTimer.value)
  }
})

// 暴露方法
defineExpose({
  showError,
  startRetry,
  cancelRetry
})
</script>

<style scoped>
.error-handler {
  position: relative;
}

.error-content {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.error-icon {
  color: var(--el-color-danger);
  flex-shrink: 0;
}

.error-info {
  flex: 1;
}

.error-info h3 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: var(--el-text-color-primary);
}

.error-message {
  margin: 0 0 12px 0;
  color: var(--el-text-color-regular);
  line-height: 1.6;
}

.error-details {
  margin-top: 12px;
}

.error-details pre {
  background: var(--el-fill-color-light);
  padding: 12px;
  border-radius: 4px;
  font-size: 12px;
  overflow-x: auto;
  max-height: 200px;
  margin: 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.network-alert {
  position: fixed;
  top: 60px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 9999;
  width: auto;
  max-width: 500px;
}

.retry-content {
  text-align: center;
  padding: 20px 0;
}

.retry-message {
  margin: 16px 0 8px 0;
  color: var(--el-text-color-regular);
}

.retry-countdown {
  color: var(--el-color-primary);
  font-weight: 500;
}
</style>
