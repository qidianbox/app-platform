import axios from 'axios'
import { ElMessage, ElNotification } from 'element-plus'

// 根据环境自动选择API地址
const getBaseURL = () => {
  // 如果是通过代理访问（开发模式），使用相对路径
  if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
    return '/api/v1'
  }
  // 生产环境或通过公网访问时，使用绝对路径
  // 将前端域名中的5173或5174替换为8080
  const apiHost = window.location.origin.replace(/517[34]/, '8080')
  return `${apiHost}/api/v1`
}

// 系统日志收集器
const systemLogger = {
  logs: [],
  errors: [],
  requests: [],
  maxLogs: 200,
  
  // 记录通用日志
  log(level, module, message, data = null) {
    const log = {
      id: Date.now() + Math.random().toString(36).substr(2, 9),
      timestamp: new Date().toISOString(),
      level,
      module,
      message,
      data,
      url: window.location.href
    }
    this.logs.push(log)
    if (this.logs.length > this.maxLogs) {
      this.logs.shift()
    }
    
    // 控制台输出
    const consoleMethod = level === 'error' ? 'error' : level === 'warn' ? 'warn' : 'log'
    console[consoleMethod](`[${level.toUpperCase()}] [${module}] ${message}`, data || '')
    
    return log
  },
  
  info(module, message, data) {
    return this.log('info', module, message, data)
  },
  
  warn(module, message, data) {
    return this.log('warn', module, message, data)
  },
  
  error(module, message, data) {
    const log = this.log('error', module, message, data)
    this.errors.push(log)
    if (this.errors.length > this.maxLogs) {
      this.errors.shift()
    }
    return log
  },
  
  // 记录API请求
  logRequest(config) {
    const log = {
      id: Date.now() + Math.random().toString(36).substr(2, 9),
      timestamp: new Date().toISOString(),
      type: 'request',
      method: config.method?.toUpperCase(),
      url: config.url,
      fullUrl: config.baseURL + config.url,
      params: config.params,
      data: config.data,
      headers: {
        Authorization: config.headers?.Authorization ? '[REDACTED]' : undefined,
        'Content-Type': config.headers?.['Content-Type']
      }
    }
    this.requests.push(log)
    if (this.requests.length > this.maxLogs) {
      this.requests.shift()
    }
    this.info('API', `Request: ${log.method} ${log.url}`, { params: log.params })
    return log
  },
  
  // 记录API响应
  logResponse(response, requestLog) {
    const log = {
      id: Date.now() + Math.random().toString(36).substr(2, 9),
      timestamp: new Date().toISOString(),
      type: 'response',
      method: response.config?.method?.toUpperCase(),
      url: response.config?.url,
      status: response.status,
      statusText: response.statusText,
      duration: requestLog ? Date.now() - new Date(requestLog.timestamp).getTime() : null,
      dataSize: JSON.stringify(response.data || {}).length
    }
    this.info('API', `Response: ${log.method} ${log.url} - ${log.status} (${log.duration}ms)`, {
      status: log.status,
      dataSize: log.dataSize
    })
    return log
  },
  
  // 记录API错误
  logApiError(error) {
    const log = {
      id: Date.now() + Math.random().toString(36).substr(2, 9),
      timestamp: new Date().toISOString(),
      type: 'api_error',
      method: error.config?.method?.toUpperCase(),
      url: error.config?.url,
      status: error.response?.status,
      statusText: error.response?.statusText,
      message: error.response?.data?.message || error.message,
      errorCode: error.response?.data?.code,
      responseData: error.response?.data,
      stack: error.stack
    }
    this.error('API', `Error: ${log.method} ${log.url} - ${log.status}: ${log.message}`, log)
    return log
  },
  
  // 获取最近的错误日志
  getRecentErrors(count = 20) {
    return this.errors.slice(-count)
  },
  
  // 获取最近的请求日志
  getRecentRequests(count = 20) {
    return this.requests.slice(-count)
  },
  
  // 获取所有日志
  getAllLogs() {
    return {
      logs: this.logs,
      errors: this.errors,
      requests: this.requests
    }
  },
  
  // 导出日志为JSON
  exportLogs() {
    const data = {
      exportTime: new Date().toISOString(),
      userAgent: navigator.userAgent,
      url: window.location.href,
      ...this.getAllLogs()
    }
    return JSON.stringify(data, null, 2)
  },
  
  // 清空日志
  clear() {
    this.logs = []
    this.errors = []
    this.requests = []
    this.info('System', 'Logs cleared')
  }
}

// 将日志收集器暴露到全局，方便调试
window.__systemLogger = systemLogger
window.__exportLogs = () => {
  const data = systemLogger.exportLogs()
  const blob = new Blob([data], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `system-logs-${new Date().toISOString().replace(/[:.]/g, '-')}.json`
  a.click()
  URL.revokeObjectURL(url)
  console.log('Logs exported successfully')
}

// 全局错误捕获
window.addEventListener('error', (event) => {
  systemLogger.error('Global', `Uncaught error: ${event.message}`, {
    filename: event.filename,
    lineno: event.lineno,
    colno: event.colno,
    error: event.error?.stack
  })
})

window.addEventListener('unhandledrejection', (event) => {
  systemLogger.error('Global', `Unhandled promise rejection: ${event.reason}`, {
    reason: event.reason?.stack || event.reason
  })
})

// 记录页面加载
systemLogger.info('System', 'Application initialized', {
  baseURL: getBaseURL(),
  userAgent: navigator.userAgent,
  timestamp: new Date().toISOString()
})

const request = axios.create({
  baseURL: getBaseURL(),
  timeout: 30000
})

systemLogger.info('API', `API client created with base URL: ${getBaseURL()}`)

// 请求拦截器
request.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  config._requestLog = systemLogger.logRequest(config)
  config._startTime = Date.now()
  return config
}, error => {
  systemLogger.error('API', 'Request interceptor error', error)
  return Promise.reject(error)
})

// 响应拦截器
request.interceptors.response.use(
  response => {
    systemLogger.logResponse(response, response.config._requestLog)
    return response.data
  },
  error => {
    const errorLog = systemLogger.logApiError(error)
    
    // 更详细的错误提示
    let errorMessage = '请求失败'
    let showNotification = false
    
    if (error.response) {
      const status = error.response.status
      const data = error.response.data
      
      if (status === 401) {
        errorMessage = '登录已过期，请重新登录'
        localStorage.removeItem('token')
        setTimeout(() => {
          window.location.href = '/login'
        }, 1500)
      } else if (status === 403) {
        errorMessage = '没有权限执行此操作'
      } else if (status === 404) {
        errorMessage = data?.message || 'API接口不存在'
      } else if (status === 500) {
        errorMessage = data?.message || '服务器内部错误'
        showNotification = true
      } else {
        errorMessage = data?.message || `请求失败 (${status})`
      }
    } else if (error.code === 'ECONNABORTED') {
      errorMessage = '请求超时，请检查网络连接'
      showNotification = true
    } else if (!navigator.onLine) {
      errorMessage = '网络连接已断开'
      showNotification = true
    }
    
    // 对于严重错误，显示通知而不是简单的消息
    if (showNotification) {
      ElNotification({
        title: '系统错误',
        message: `${errorMessage}\n\n错误ID: ${errorLog.id}\n\n您可以在控制台输入 __exportLogs() 导出日志`,
        type: 'error',
        duration: 10000
      })
    } else {
      ElMessage.error(errorMessage)
    }
    
    return Promise.reject(error)
  }
)

// 导出日志收集器供其他模块使用
export const logger = systemLogger

export default request
