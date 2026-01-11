<template>
  <div class="module-versions">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <el-button text @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回模块配置
        </el-button>
        <h2 class="page-title">{{ moduleName }} - 版本管理</h2>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="showCreateDialog = true">
          <el-icon><Plus /></el-icon>
          创建新版本
        </el-button>
        <el-button @click="loadVersions">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 筛选区域 -->
    <div class="filter-section">
      <el-select v-model="filterStatus" placeholder="状态筛选" clearable style="width: 150px;">
        <el-option label="全部状态" value="" />
        <el-option label="草稿" value="draft" />
        <el-option label="已发布" value="published" />
        <el-option label="已废弃" value="deprecated" />
      </el-select>
      <el-select v-model="filterEnvironment" placeholder="环境筛选" clearable style="width: 150px; margin-left: 10px;">
        <el-option label="全部环境" value="" />
        <el-option label="开发环境" value="dev" />
        <el-option label="测试环境" value="test" />
        <el-option label="生产环境" value="prod" />
      </el-select>
    </div>

    <!-- 版本列表 -->
    <el-table :data="versions" style="width: 100%" v-loading="loading">
      <el-table-column prop="version" label="版本号" width="120">
        <template #default="{ row }">
          <span class="version-tag">v{{ row.version }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="environment" label="环境" width="100">
        <template #default="{ row }">
          <el-tag :type="getEnvType(row.environment)" size="small" effect="plain">
            {{ getEnvText(row.environment) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="changelog" label="变更说明" min-width="200">
        <template #default="{ row }">
          <span class="changelog-text">{{ row.changelog || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_by" label="创建人" width="120" />
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column prop="published_at" label="发布时间" width="180">
        <template #default="{ row }">
          {{ row.published_at ? formatDate(row.published_at) : '-' }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button 
            v-if="row.status === 'draft'" 
            type="primary" 
            size="small" 
            @click="publishVersion(row)"
          >
            发布
          </el-button>
          <el-button 
            v-if="row.status === 'published'" 
            type="warning" 
            size="small" 
            @click="rollbackVersion(row)"
          >
            回滚
          </el-button>
          <el-dropdown @command="(cmd) => handleCommand(cmd, row)" style="margin-left: 8px;">
            <el-button size="small">
              更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="view">查看配置</el-dropdown-item>
                <el-dropdown-item command="compare">版本对比</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadVersions"
        @current-change="loadVersions"
      />
    </div>

    <!-- 创建版本对话框 -->
    <el-dialog v-model="showCreateDialog" title="创建新版本" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="版本号">
          <el-input v-model="createForm.version" placeholder="留空自动生成，如 1.0.0" />
        </el-form-item>
        <el-form-item label="目标环境">
          <el-select v-model="createForm.environment" style="width: 100%;">
            <el-option label="开发环境" value="dev" />
            <el-option label="测试环境" value="test" />
            <el-option label="生产环境" value="prod" />
          </el-select>
        </el-form-item>
        <el-form-item label="变更说明">
          <el-input 
            v-model="createForm.changelog" 
            type="textarea" 
            :rows="4"
            placeholder="描述本次版本的变更内容" 
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="createVersion" :loading="creating">创建</el-button>
      </template>
    </el-dialog>

    <!-- 查看配置对话框 -->
    <el-dialog v-model="showConfigDialog" title="配置详情" width="700px">
      <div class="config-viewer">
        <pre>{{ configContent }}</pre>
      </div>
    </el-dialog>

    <!-- 版本对比对话框 -->
    <el-dialog v-model="showCompareDialog" title="版本对比" width="900px">
      <div class="compare-section">
        <div class="compare-header">
          <el-select v-model="compareVersion1" placeholder="选择版本1" style="width: 200px;">
            <el-option 
              v-for="v in versions" 
              :key="v.id" 
              :label="`v${v.version}`" 
              :value="v.id" 
            />
          </el-select>
          <span class="compare-arrow">→</span>
          <el-select v-model="compareVersion2" placeholder="选择版本2" style="width: 200px;">
            <el-option 
              v-for="v in versions" 
              :key="v.id" 
              :label="`v${v.version}`" 
              :value="v.id" 
            />
          </el-select>
          <el-button type="primary" @click="compareVersions" :loading="comparing" style="margin-left: 20px;">
            对比
          </el-button>
        </div>
        <div class="compare-result" v-if="compareResult">
          <div class="compare-info">
            <div class="version-info">
              <strong>版本1:</strong> v{{ compareResult.version1?.version }} ({{ compareResult.version1?.status }})
            </div>
            <div class="version-info">
              <strong>版本2:</strong> v{{ compareResult.version2?.version }} ({{ compareResult.version2?.status }})
            </div>
          </div>
          <el-table :data="compareResult.changes" style="width: 100%; margin-top: 20px;">
            <el-table-column prop="field" label="字段" width="200" />
            <el-table-column prop="type" label="变更类型" width="120">
              <template #default="{ row }">
                <el-tag :type="getChangeType(row.type)" size="small">
                  {{ getChangeText(row.type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="old_value" label="旧值" min-width="200">
              <template #default="{ row }">
                {{ row.old_value !== undefined ? JSON.stringify(row.old_value) : '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="new_value" label="新值" min-width="200">
              <template #default="{ row }">
                {{ row.new_value !== undefined ? JSON.stringify(row.new_value) : '-' }}
              </template>
            </el-table-column>
          </el-table>
          <div v-if="compareResult.changes?.length === 0" class="no-changes">
            两个版本配置相同，无差异
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Plus, Refresh, ArrowDown } from '@element-plus/icons-vue'
import request from '@/utils/request'

const route = useRoute()
const router = useRouter()

// 路由参数
const appId = ref(route.params.appId || route.params.id)
const moduleCode = ref(route.params.moduleCode || route.params.module_code)
const moduleName = ref(route.query.name || moduleCode.value)

// 列表数据
const versions = ref([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 筛选
const filterStatus = ref('')
const filterEnvironment = ref('')

// 创建版本
const showCreateDialog = ref(false)
const creating = ref(false)
const createForm = ref({
  version: '',
  environment: 'dev',
  changelog: ''
})

// 查看配置
const showConfigDialog = ref(false)
const configContent = ref('')

// 版本对比
const showCompareDialog = ref(false)
const comparing = ref(false)
const compareVersion1 = ref(null)
const compareVersion2 = ref(null)
const compareResult = ref(null)

// 加载版本列表
const loadVersions = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      size: pageSize.value
    }
    if (filterStatus.value) params.status = filterStatus.value
    if (filterEnvironment.value) params.environment = filterEnvironment.value

    const res = await request.get(`/apps/${appId.value}/modules/${moduleCode.value}/versions`, { params })
    // 响应拦截器已经解包，res直接是data
    if (res && res.list) {
      versions.value = res.list || []
      total.value = res.total || 0
    } else if (Array.isArray(res)) {
      versions.value = res
      total.value = res.length
    } else {
      versions.value = []
      total.value = 0
    }
  } catch (error) {
    console.error('加载版本列表失败:', error)
    ElMessage.error('加载版本列表失败')
  } finally {
    loading.value = false
  }
}

// 创建版本
const createVersion = async () => {
  creating.value = true
  try {
    await request.post(`/apps/${appId.value}/modules/${moduleCode.value}/versions`, createForm.value)
    ElMessage.success('版本创建成功')
    showCreateDialog.value = false
    createForm.value = { version: '', environment: 'dev', changelog: '' }
    loadVersions()
  } catch (error) {
    ElMessage.error('创建版本失败')
  } finally {
    creating.value = false
  }
}

// 发布版本
const publishVersion = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要发布版本 v${row.version} 吗？`, '发布确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.post(`/apps/${appId.value}/modules/${moduleCode.value}/versions/${row.id}/publish`)
    ElMessage.success('版本发布成功')
    loadVersions()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('发布失败')
    }
  }
}

// 回滚版本
const rollbackVersion = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要回滚到版本 v${row.version} 吗？这将用该版本的配置覆盖当前配置。`, '回滚确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.post(`/apps/${appId.value}/modules/${moduleCode.value}/versions/${row.id}/rollback`)
    ElMessage.success('版本回滚成功')
    loadVersions()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('回滚失败')
    }
  }
}

// 版本对比
const compareVersions = async () => {
  if (!compareVersion1.value || !compareVersion2.value) {
    ElMessage.warning('请选择两个版本进行对比')
    return
  }
  comparing.value = true
  try {
    const res = await request.get(`/apps/${appId.value}/modules/${moduleCode.value}/versions/compare`, {
      params: {
        version1: compareVersion1.value,
        version2: compareVersion2.value
      }
    })
    compareResult.value = res
  } catch (error) {
    ElMessage.error('版本对比失败')
  } finally {
    comparing.value = false
  }
}

// 处理下拉菜单命令
const handleCommand = (command, row) => {
  if (command === 'view') {
    try {
      configContent.value = JSON.stringify(JSON.parse(row.config_snapshot || '{}'), null, 2)
    } catch {
      configContent.value = row.config_snapshot || '{}'
    }
    showConfigDialog.value = true
  } else if (command === 'compare') {
    compareVersion1.value = row.id
    compareVersion2.value = null
    compareResult.value = null
    showCompareDialog.value = true
  }
}

// 返回
const goBack = () => {
  router.push(`/apps/${appId.value}/config?tab=config`)
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

// 状态相关
const getStatusType = (status) => {
  const types = { draft: 'info', published: 'success', deprecated: 'warning' }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = { draft: '草稿', published: '已发布', deprecated: '已废弃' }
  return texts[status] || status
}

// 环境相关
const getEnvType = (env) => {
  const types = { dev: '', test: 'warning', prod: 'danger' }
  return types[env] || ''
}

const getEnvText = (env) => {
  const texts = { dev: '开发', test: '测试', prod: '生产' }
  return texts[env] || env
}

// 变更类型相关
const getChangeType = (type) => {
  const types = { added: 'success', modified: 'warning', removed: 'danger' }
  return types[type] || 'info'
}

const getChangeText = (type) => {
  const texts = { added: '新增', modified: '修改', removed: '删除' }
  return texts[type] || type
}

// 监听筛选变化
watch([filterStatus, filterEnvironment], () => {
  currentPage.value = 1
  loadVersions()
})

onMounted(() => {
  loadVersions()
})
</script>

<style scoped>
.module-versions {
  padding: 20px;
  background: #fff;
  min-height: calc(100vh - 60px);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #eee;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.page-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.filter-section {
  margin-bottom: 20px;
}

.version-tag {
  font-family: 'Monaco', 'Menlo', monospace;
  font-weight: 600;
  color: #409eff;
}

.changelog-text {
  color: #666;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.config-viewer {
  background: #f5f7fa;
  border-radius: 4px;
  padding: 15px;
  max-height: 500px;
  overflow: auto;
}

.config-viewer pre {
  margin: 0;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 13px;
  white-space: pre-wrap;
  word-break: break-all;
}

.compare-section {
  padding: 10px 0;
}

.compare-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
}

.compare-arrow {
  font-size: 20px;
  color: #909399;
}

.compare-info {
  display: flex;
  gap: 30px;
  padding: 10px 15px;
  background: #f5f7fa;
  border-radius: 4px;
}

.version-info {
  color: #666;
}

.no-changes {
  text-align: center;
  color: #909399;
  padding: 40px 0;
}
</style>
