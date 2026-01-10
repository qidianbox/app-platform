/**
 * 请求重试工具
 * 提供自动重试、指数退避等功能
 */

/**
 * 延迟函数
 * @param {number} ms 延迟毫秒数
 * @returns {Promise}
 */
export const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms))

/**
 * 判断错误是否可重试
 * @param {Error} error 错误对象
 * @returns {boolean}
 */
export const isRetryableError = (error) => {
  // 网络错误
  if (!error.response) {
    return true
  }

  const status = error.response.status

  // 服务器错误（5xx）可重试
  if (status >= 500 && status < 600) {
    return true
  }

  // 请求超时可重试
  if (error.code === 'ECONNABORTED') {
    return true
  }

  // 429 Too Many Requests 可重试（等待后）
  if (status === 429) {
    return true
  }

  // 其他错误不重试
  return false
}

/**
 * 计算退避时间（指数退避）
 * @param {number} attempt 当前尝试次数（从0开始）
 * @param {number} baseDelay 基础延迟（毫秒）
 * @param {number} maxDelay 最大延迟（毫秒）
 * @returns {number}
 */
export const calculateBackoff = (attempt, baseDelay = 1000, maxDelay = 30000) => {
  // 指数退避 + 随机抖动
  const exponentialDelay = baseDelay * Math.pow(2, attempt)
  const jitter = Math.random() * 1000
  return Math.min(exponentialDelay + jitter, maxDelay)
}

/**
 * 带重试的请求封装
 * @param {Function} requestFn 请求函数
 * @param {Object} options 配置选项
 * @returns {Promise}
 */
export const withRetry = async (requestFn, options = {}) => {
  const {
    maxRetries = 3,
    baseDelay = 1000,
    maxDelay = 30000,
    shouldRetry = isRetryableError,
    onRetry = null
  } = options

  let lastError

  for (let attempt = 0; attempt <= maxRetries; attempt++) {
    try {
      return await requestFn()
    } catch (error) {
      lastError = error

      // 检查是否应该重试
      if (attempt >= maxRetries || !shouldRetry(error)) {
        throw error
      }

      // 计算退避时间
      const backoffTime = calculateBackoff(attempt, baseDelay, maxDelay)

      // 回调通知
      if (onRetry) {
        onRetry({
          attempt: attempt + 1,
          maxRetries,
          error,
          nextRetryIn: backoffTime
        })
      }

      // 等待后重试
      await delay(backoffTime)
    }
  }

  throw lastError
}

/**
 * 创建可重试的请求函数
 * @param {Function} requestFn 原始请求函数
 * @param {Object} defaultOptions 默认配置
 * @returns {Function}
 */
export const createRetryableRequest = (requestFn, defaultOptions = {}) => {
  return async (...args) => {
    return withRetry(() => requestFn(...args), defaultOptions)
  }
}

/**
 * 批量请求（带并发控制和重试）
 * @param {Array<Function>} requests 请求函数数组
 * @param {Object} options 配置选项
 * @returns {Promise<Array>}
 */
export const batchRequest = async (requests, options = {}) => {
  const {
    concurrency = 3,
    maxRetries = 2,
    onProgress = null
  } = options

  const results = []
  const errors = []
  let completed = 0

  // 分批执行
  for (let i = 0; i < requests.length; i += concurrency) {
    const batch = requests.slice(i, i + concurrency)
    const batchPromises = batch.map(async (request, index) => {
      try {
        const result = await withRetry(request, { maxRetries })
        results[i + index] = { success: true, data: result }
      } catch (error) {
        results[i + index] = { success: false, error }
        errors.push({ index: i + index, error })
      } finally {
        completed++
        if (onProgress) {
          onProgress({
            completed,
            total: requests.length,
            progress: completed / requests.length
          })
        }
      }
    })

    await Promise.all(batchPromises)
  }

  return {
    results,
    errors,
    hasErrors: errors.length > 0
  }
}

/**
 * 请求队列（串行执行，带重试）
 */
export class RequestQueue {
  constructor(options = {}) {
    this.queue = []
    this.processing = false
    this.options = {
      maxRetries: 2,
      delayBetween: 100,
      ...options
    }
  }

  /**
   * 添加请求到队列
   * @param {Function} requestFn 请求函数
   * @returns {Promise}
   */
  add(requestFn) {
    return new Promise((resolve, reject) => {
      this.queue.push({ requestFn, resolve, reject })
      this.process()
    })
  }

  /**
   * 处理队列
   */
  async process() {
    if (this.processing || this.queue.length === 0) {
      return
    }

    this.processing = true

    while (this.queue.length > 0) {
      const { requestFn, resolve, reject } = this.queue.shift()

      try {
        const result = await withRetry(requestFn, {
          maxRetries: this.options.maxRetries
        })
        resolve(result)
      } catch (error) {
        reject(error)
      }

      if (this.queue.length > 0) {
        await delay(this.options.delayBetween)
      }
    }

    this.processing = false
  }

  /**
   * 清空队列
   */
  clear() {
    this.queue.forEach(({ reject }) => {
      reject(new Error('Queue cleared'))
    })
    this.queue = []
  }

  /**
   * 获取队列长度
   */
  get length() {
    return this.queue.length
  }
}

export default {
  delay,
  isRetryableError,
  calculateBackoff,
  withRetry,
  createRetryableRequest,
  batchRequest,
  RequestQueue
}
