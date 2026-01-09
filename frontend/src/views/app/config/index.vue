<template>
  <div class="app-detail">
    <!-- 顶部导航栏 -->
    <div class="top-header">
      <div class="header-left">
        <el-button class="back-btn" text @click="$router.push('/apps')">
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        <div class="app-icon">
          <span>{{ appInfo.name?.charAt(0) || '应' }}</span>
        </div>
        <span class="app-name">{{ appInfo.name || '加载中...' }}</span>
      </div>
      <div class="header-tabs">
        <div 
          class="tab-item" 
          :class="{ active: activeTab === 'workspace' }"
          @click="activeTab = 'workspace'"
        >
          <span class="tab-icon">🚀</span>
          <span>工作台</span>
        </div>
        <div 
          class="tab-item" 
          :class="{ active: activeTab === 'config' }"
          @click="activeTab = 'config'"
        >
          <span class="tab-icon">⚙️</span>
          <span>配置中心</span>
        </div>
      </div>
      <div class="header-right">
        <el-dropdown>
          <el-button text class="user-btn">
            <el-avatar :size="32">{{ adminName?.charAt(0) || 'A' }}</el-avatar>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="$router.push('/apps')">返回APP列表</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- 工作台内容 -->
    <div v-if="activeTab === 'workspace'" class="workspace-container">
      <div class="workspace-sidebar">
        <div class="sidebar-title">功能菜单</div>
        <div 
          v-for="func in workspaceFunctions" 
          :key="func.id"
          class="sidebar-item"
          :class="{ active: activeFunction === func.id }"
          @click="activeFunction = func.id"
        >
          <span class="func-icon">{{ func.icon }}</span>
          <span class="func-name">{{ func.name }}</span>
        </div>
        <div v-if="workspaceFunctions.length === 0" class="empty-sidebar">
          <p>暂无功能</p>
          <p class="hint">请在配置中心添加模块</p>
        </div>
      </div>
      <div class="workspace-main">
        <div v-if="activeFunction" class="function-content">
          <h2>{{ currentFunction?.name }}</h2>
          <p class="function-desc">{{ currentFunction?.description }}</p>
          <div class="function-body">
            <!-- 根据不同功能显示不同内容 -->
            <component :is="getFunctionComponent(activeFunction)" :app-id="appId" />
          </div>
        </div>
        <div v-else class="empty-workspace">
          <el-empty description="请从左侧选择功能">
            <template #image>
              <div class="empty-icon">🚀</div>
            </template>
          </el-empty>
        </div>
      </div>
    </div>

    <!-- 配置中心内容 -->
    <div v-if="activeTab === 'config'" class="config-container">
      <div class="config-header">
        <h2>模块配置</h2>
        <p>管理该APP已启用的模块配置</p>
      </div>
      <div class="modules-grid">
        <div 
          v-for="module in appModules" 
          :key="module.id"
          class="module-card"
          @click="openModuleConfig(module)"
        >
          <div class="module-icon">{{ module.icon || '📦' }}</div>
          <div class="module-info">
            <h3>{{ module.name }}</h3>
            <p>{{ module.description || '暂无描述' }}</p>
          </div>
          <div class="module-status">
            <el-tag :type="module.enabled ? 'success' : 'info'" size="small">
              {{ module.enabled ? '已启用' : '未启用' }}
            </el-tag>
          </div>
          <el-icon class="arrow-icon"><ArrowRight /></el-icon>
        </div>
        <div v-if="appModules.length === 0" class="empty-modules">
          <el-empty description="该APP暂未配置任何模块">
            <el-button type="primary" @click="$router.push('/apps')">
              返回管理模块
            </el-button>
          </el-empty>
        </div>
      </div>
    </div>

    <!-- 模块配置弹窗 -->
    <el-dialog 
      v-model="configDialogVisible" 
      :title="`${currentModule?.name || ''} 配置`"
      width="600px"
    >
      <div v-if="currentModule" class="module-config-form">
        <el-form :model="moduleConfigForm" label-width="120px">
          <el-form-item label="启用状态">
            <el-switch v-model="moduleConfigForm.enabled" />
          </el-form-item>
          <el-form-item label="API端点">
            <el-input v-model="moduleConfigForm.apiEndpoint" placeholder="请输入API端点" />
          </el-form-item>
          <el-form-item label="超时时间(ms)">
            <el-input-number v-model="moduleConfigForm.timeout" :min="1000" :max="60000" />
          </el-form-item>
          <el-form-item label="重试次数">
            <el-input-number v-model="moduleConfigForm.retryCount" :min="0" :max="10" />
          </el-form-item>
          <el-form-item label="备注">
            <el-input v-model="moduleConfigForm.remark" type="textarea" rows="3" />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="configDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveModuleConfig">保存配置</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowLeft, ArrowRight } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const route = useRoute()
const appId = computed(() => route.params.id)

const activeTab = ref('workspace')
const activeFunction = ref('')
const configDialogVisible = ref(false)
const currentModule = ref(null)
const adminName = ref(localStorage.getItem('adminName') || 'Admin')

const appInfo = ref({
  name: '',
  app_id: '',
  description: ''
})

const appModules = ref([])

const workspaceFunctions = ref([])

const moduleConfigForm = ref({
  enabled: true,
  apiEndpoint: '',
  timeout: 5000,
  retryCount: 3,
  remark: ''
})

const currentFunction = computed(() => {
  return workspaceFunctions.value.find(f => f.id === activeFunction.value)
})

// 获取APP信息
const fetchAppInfo = async () => {
  try {
    const res = await request.get(`/api/v1/apps/${appId.value}`)
    if (res.code === 0 && res.data) {
      appInfo.value = res.data
    }
  } catch (error) {
    console.error('获取APP信息失败:', error)
  }
}

// 获取APP模块列表
const fetchAppModules = async () => {
  try {
    const res = await request.get(`/api/v1/apps/${appId.value}/modules`)
    if (res.code === 0 && res.data) {
      appModules.value = res.data
      // 根据模块生成工作台功能
      generateWorkspaceFunctions()
    }
  } catch (error) {
    console.error('获取APP模块失败:', error)
    // 使用模拟数据
    appModules.value = [
      { id: 1, name: '用户管理', icon: '👥', description: '管理APP用户', enabled: true },
      { id: 2, name: '消息中心', icon: '💬', description: '消息推送管理', enabled: true },
      { id: 3, name: '数据统计', icon: '📊', description: '数据分析统计', enabled: true },
      { id: 4, name: '版本管理', icon: '📦', description: 'APP版本控制', enabled: true },
      { id: 5, name: '配置管理', icon: '⚙️', description: '远程配置管理', enabled: false }
    ]
    generateWorkspaceFunctions()
  }
}

// 根据模块生成工作台功能
const generateWorkspaceFunctions = () => {
  const functions = []
  appModules.value.forEach(module => {
    if (module.enabled) {
      // 根据模块类型添加对应的工作台功能
      if (module.name.includes('用户')) {
        functions.push({ id: 'user-list', name: '用户列表', icon: '👥', description: '查看和管理用户', module: module.id })
        functions.push({ id: 'user-stats', name: '用户统计', icon: '📈', description: '用户数据统计', module: module.id })
      }
      if (module.name.includes('消息')) {
        functions.push({ id: 'send-message', name: '发送消息', icon: '✉️', description: '发送站内消息', module: module.id })
        functions.push({ id: 'message-list', name: '消息记录', icon: '📋', description: '查看消息历史', module: module.id })
      }
      if (module.name.includes('统计') || module.name.includes('数据')) {
        functions.push({ id: 'data-overview', name: '数据概览', icon: '📊', description: '数据统计概览', module: module.id })
        functions.push({ id: 'event-tracking', name: '事件追踪', icon: '🎯', description: '用户行为追踪', module: module.id })
      }
      if (module.name.includes('版本')) {
        functions.push({ id: 'version-list', name: '版本列表', icon: '📦', description: '管理APP版本', module: module.id })
        functions.push({ id: 'release', name: '发布版本', icon: '🚀', description: '发布新版本', module: module.id })
      }
    }
  })
  workspaceFunctions.value = functions
  if (functions.length > 0) {
    activeFunction.value = functions[0].id
  }
}

// 打开模块配置
const openModuleConfig = (module) => {
  currentModule.value = module
  moduleConfigForm.value = {
    enabled: module.enabled,
    apiEndpoint: module.apiEndpoint || '',
    timeout: module.timeout || 5000,
    retryCount: module.retryCount || 3,
    remark: module.remark || ''
  }
  configDialogVisible.value = true
}

// 保存模块配置
const saveModuleConfig = async () => {
  try {
    // 这里调用后端API保存配置
    ElMessage.success('配置保存成功')
    configDialogVisible.value = false
    // 更新本地数据
    const index = appModules.value.findIndex(m => m.id === currentModule.value.id)
    if (index !== -1) {
      appModules.value[index] = {
        ...appModules.value[index],
        ...moduleConfigForm.value
      }
    }
    generateWorkspaceFunctions()
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

// 获取功能组件
const getFunctionComponent = (funcId) => {
  // 返回对应的功能组件，这里可以根据需要扩展
  return 'div'
}

onMounted(() => {
  fetchAppInfo()
  fetchAppModules()
})
</script>

<style lang="scss" scoped>
.app-detail {
  min-height: 100vh;
  background: #f5f7fa;
}

.top-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
  
  .back-btn {
    color: white;
    &:hover {
      background: rgba(255,255,255,0.1);
    }
  }
  
  .app-icon {
    width: 48px;
    height: 48px;
    background: white;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    font-weight: bold;
    color: #667eea;
  }
  
  .app-name {
    font-size: 18px;
    font-weight: 600;
  }
}

.header-tabs {
  display: flex;
  gap: 8px;
  
  .tab-item {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 10px 20px;
    border-radius: 20px;
    cursor: pointer;
    transition: all 0.3s;
    font-size: 15px;
    
    .tab-icon {
      font-size: 16px;
    }
    
    &:hover {
      background: rgba(255,255,255,0.2);
    }
    
    &.active {
      background: white;
      color: #667eea;
      font-weight: 500;
    }
  }
}

.header-right {
  .user-btn {
    color: white;
  }
}

// 工作台样式
.workspace-container {
  display: flex;
  height: calc(100vh - 72px);
}

.workspace-sidebar {
  width: 240px;
  background: white;
  border-right: 1px solid #e4e7ed;
  padding: 16px 0;
  overflow-y: auto;
  
  .sidebar-title {
    padding: 8px 20px;
    font-size: 12px;
    color: #909399;
    text-transform: uppercase;
    letter-spacing: 1px;
  }
  
  .sidebar-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 20px;
    cursor: pointer;
    transition: all 0.2s;
    
    .func-icon {
      font-size: 18px;
    }
    
    .func-name {
      font-size: 14px;
      color: #303133;
    }
    
    &:hover {
      background: #f5f7fa;
    }
    
    &.active {
      background: #ecf5ff;
      border-right: 3px solid #409eff;
      
      .func-name {
        color: #409eff;
        font-weight: 500;
      }
    }
  }
  
  .empty-sidebar {
    padding: 40px 20px;
    text-align: center;
    color: #909399;
    
    .hint {
      font-size: 12px;
      margin-top: 8px;
    }
  }
}

.workspace-main {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  
  .function-content {
    background: white;
    border-radius: 8px;
    padding: 24px;
    
    h2 {
      margin: 0 0 8px;
      font-size: 20px;
      color: #303133;
    }
    
    .function-desc {
      color: #909399;
      margin-bottom: 24px;
    }
  }
  
  .empty-workspace {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    
    .empty-icon {
      font-size: 64px;
    }
  }
}

// 配置中心样式
.config-container {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  
  .config-header {
    margin-bottom: 24px;
    
    h2 {
      margin: 0 0 8px;
      font-size: 24px;
      color: #303133;
    }
    
    p {
      color: #909399;
      margin: 0;
    }
  }
}

.modules-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.module-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: white;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid #e4e7ed;
  
  &:hover {
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
    transform: translateY(-2px);
  }
  
  .module-icon {
    font-size: 32px;
    width: 56px;
    height: 56px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f5f7fa;
    border-radius: 12px;
  }
  
  .module-info {
    flex: 1;
    
    h3 {
      margin: 0 0 4px;
      font-size: 16px;
      color: #303133;
    }
    
    p {
      margin: 0;
      font-size: 13px;
      color: #909399;
    }
  }
  
  .arrow-icon {
    color: #c0c4cc;
  }
}

.empty-modules {
  grid-column: 1 / -1;
  padding: 60px;
}

.module-config-form {
  padding: 10px 0;
}

// 移动端适配
@media (max-width: 768px) {
  .top-header {
    flex-wrap: wrap;
    gap: 12px;
    padding: 12px 16px;
  }
  
  .header-tabs {
    order: 3;
    width: 100%;
    justify-content: center;
    
    .tab-item {
      padding: 8px 16px;
      font-size: 14px;
    }
  }
  
  .workspace-container {
    flex-direction: column;
    height: auto;
  }
  
  .workspace-sidebar {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid #e4e7ed;
    display: flex;
    overflow-x: auto;
    padding: 8px;
    
    .sidebar-title {
      display: none;
    }
    
    .sidebar-item {
      flex-shrink: 0;
      padding: 8px 12px;
      border-radius: 20px;
      
      &.active {
        background: #409eff;
        border-right: none;
        
        .func-name {
          color: white;
        }
      }
    }
  }
  
  .workspace-main {
    padding: 16px;
  }
  
  .config-container {
    padding: 16px;
  }
  
  .modules-grid {
    grid-template-columns: 1fr;
  }
}
</style>
