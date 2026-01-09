// 9个大模块定义 - 对应后端的9个source_module
export const moduleList = [
  { 
    id: 'user_management', 
    name: '用户管理', 
    icon: '👥', 
    description: '用户注册、登录、权限管理',
    category: 'user'
  },
  { 
    id: 'message_center', 
    name: '消息中心', 
    icon: '💬', 
    description: '站内消息、通知管理',
    category: 'message'
  },
  { 
    id: 'push_service', 
    name: '推送服务', 
    icon: '🔔', 
    description: 'APP推送通知服务',
    category: 'message'
  },
  { 
    id: 'event_tracking', 
    name: '数据埋点', 
    icon: '📊', 
    description: '用户行为埋点和数据分析',
    category: 'data'
  },
  { 
    id: 'log_service', 
    name: '日志服务', 
    icon: '📝', 
    description: '应用日志收集和查询',
    category: 'system'
  },
  { 
    id: 'monitor_service', 
    name: '监控告警', 
    icon: '📡', 
    description: '应用监控和告警通知',
    category: 'system'
  },
  { 
    id: 'file_storage', 
    name: '文件存储', 
    icon: '📁', 
    description: '文件上传、下载、管理',
    category: 'storage'
  },
  { 
    id: 'config_management', 
    name: '配置管理', 
    icon: '⚙️', 
    description: '远程配置下发和管理',
    category: 'storage'
  },
  { 
    id: 'version_management', 
    name: '版本管理', 
    icon: '📦', 
    description: 'APP版本发布和更新',
    category: 'other'
  }
]

// 模块分组（用于UI展示）
export const moduleCategories = [
  { id: 'user', name: '用户与权限', icon: '👤', description: '用户管理和权限控制' },
  { id: 'message', name: '消息与推送', icon: '📬', description: '消息中心和推送服务' },
  { id: 'data', name: '数据与分析', icon: '📊', description: '埋点和数据分析' },
  { id: 'system', name: '系统服务', icon: '⚙️', description: '日志、监控等系统服务' },
  { id: 'storage', name: '存储服务', icon: '📁', description: '文件存储和配置管理' },
  { id: 'other', name: '其他', icon: '📦', description: '其他功能模块' }
]

// 获取分组后的大模块列表
export const getGroupedModules = () => {
  return moduleCategories.map(cat => ({
    ...cat,
    modules: moduleList.filter(m => m.category === cat.id)
  }))
}

// 根据模块ID获取模块信息
export const getModuleById = (id) => {
  return moduleList.find(m => m.id === id)
}

// 获取所有模块ID列表
export const getAllModuleIds = () => {
  return moduleList.map(m => m.id)
}
