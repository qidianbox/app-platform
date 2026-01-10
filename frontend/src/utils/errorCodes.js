/**
 * APP中台管理系统 - 统一错误码定义
 * 与后端response包保持一致
 */

// 错误码定义
export const ErrorCodes = {
  // 成功
  SUCCESS: 0,
  
  // 通用错误 (1000-1999)
  UNKNOWN_ERROR: 1000,
  INVALID_PARAMS: 1001,
  INTERNAL_ERROR: 1002,
  DATABASE_ERROR: 1003,
  NETWORK_ERROR: 1004,
  TIMEOUT_ERROR: 1005,
  
  // 认证错误 (2000-2999)
  UNAUTHORIZED: 2000,
  TOKEN_EXPIRED: 2001,
  TOKEN_INVALID: 2002,
  LOGIN_FAILED: 2003,
  PERMISSION_DENIED: 2004,
  
  // 业务错误 (3000-3999)
  NOT_FOUND: 3000,
  ALREADY_EXISTS: 3001,
  OPERATION_FAILED: 3002,
  RATE_LIMITED: 3003,
  VALIDATION_FAILED: 3004,
  
  // APP相关错误 (4000-4999)
  APP_NOT_FOUND: 4000,
  APP_DISABLED: 4001,
  APP_NAME_EXISTS: 4002,
  APP_PACKAGE_EXISTS: 4003,
  
  // 用户相关错误 (5000-5999)
  USER_NOT_FOUND: 5000,
  USER_DISABLED: 5001,
  USER_EXISTS: 5002,
  
  // 模块相关错误 (6000-6999)
  MODULE_NOT_FOUND: 6000,
  MODULE_DISABLED: 6001,
  MODULE_CONFIG_ERROR: 6002
}

// 错误消息映射
export const ErrorMessages = {
  [ErrorCodes.SUCCESS]: '操作成功',
  
  // 通用错误
  [ErrorCodes.UNKNOWN_ERROR]: '未知错误，请稍后重试',
  [ErrorCodes.INVALID_PARAMS]: '请求参数无效',
  [ErrorCodes.INTERNAL_ERROR]: '服务器内部错误',
  [ErrorCodes.DATABASE_ERROR]: '数据库操作失败',
  [ErrorCodes.NETWORK_ERROR]: '网络连接失败，请检查网络',
  [ErrorCodes.TIMEOUT_ERROR]: '请求超时，请稍后重试',
  
  // 认证错误
  [ErrorCodes.UNAUTHORIZED]: '请先登录',
  [ErrorCodes.TOKEN_EXPIRED]: '登录已过期，请重新登录',
  [ErrorCodes.TOKEN_INVALID]: '登录凭证无效',
  [ErrorCodes.LOGIN_FAILED]: '用户名或密码错误',
  [ErrorCodes.PERMISSION_DENIED]: '没有操作权限',
  
  // 业务错误
  [ErrorCodes.NOT_FOUND]: '资源不存在',
  [ErrorCodes.ALREADY_EXISTS]: '资源已存在',
  [ErrorCodes.OPERATION_FAILED]: '操作失败',
  [ErrorCodes.RATE_LIMITED]: '请求过于频繁，请稍后重试',
  [ErrorCodes.VALIDATION_FAILED]: '数据验证失败',
  
  // APP相关错误
  [ErrorCodes.APP_NOT_FOUND]: 'APP不存在',
  [ErrorCodes.APP_DISABLED]: 'APP已禁用',
  [ErrorCodes.APP_NAME_EXISTS]: 'APP名称已存在',
  [ErrorCodes.APP_PACKAGE_EXISTS]: '包名已存在',
  
  // 用户相关错误
  [ErrorCodes.USER_NOT_FOUND]: '用户不存在',
  [ErrorCodes.USER_DISABLED]: '用户已禁用',
  [ErrorCodes.USER_EXISTS]: '用户已存在',
  
  // 模块相关错误
  [ErrorCodes.MODULE_NOT_FOUND]: '模块不存在',
  [ErrorCodes.MODULE_DISABLED]: '模块已禁用',
  [ErrorCodes.MODULE_CONFIG_ERROR]: '模块配置错误'
}

// 获取错误消息
export function getErrorMessage(code, defaultMessage = '操作失败') {
  return ErrorMessages[code] || defaultMessage
}

// 判断是否需要重新登录
export function isAuthError(code) {
  return code >= 2000 && code < 3000
}

// 判断是否可以重试
export function isRetryableError(code) {
  const retryableCodes = [
    ErrorCodes.NETWORK_ERROR,
    ErrorCodes.TIMEOUT_ERROR,
    ErrorCodes.INTERNAL_ERROR,
    ErrorCodes.DATABASE_ERROR,
    ErrorCodes.RATE_LIMITED
  ]
  return retryableCodes.includes(code)
}

// 判断是否是限流错误
export function isRateLimitError(code) {
  return code === ErrorCodes.RATE_LIMITED
}

export default {
  ErrorCodes,
  ErrorMessages,
  getErrorMessage,
  isAuthError,
  isRetryableError,
  isRateLimitError
}
