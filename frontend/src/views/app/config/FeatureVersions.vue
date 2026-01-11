<template>
  <div class="feature-versions">
    <!-- 返回按钮 -->
    <div class="header-bar">
      <el-button @click="goBack" :icon="ArrowLeft">返回数据模型管理</el-button>
      <span class="collection-name">{{ collectionName }} - 版本管理</span>
    </div>

    <!-- 操作栏 -->
    <div class="action-bar">
      <el-button type="primary" @click="createVersion" :icon="Plus">创建新版本</el-button>
      <el-button @click="loadVersions" :icon="Refresh">刷新</el-button>
    </div>

    <!-- 版本列表 -->
    <el-table :data="versions" v-loading="loading" stripe>
      <el-table-column prop="version" label="版本号" width="120">
        <template #default="{ row }">
          <el-tag>v{{ row.version }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="changelog" label="变更说明" min-width="200">
        <template #default="{ row }">
          {{ row.changelog || '无说明' }}
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
      <el-table-column label="操作" width="280" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="viewVersion(row)" :icon="View">查看</el-button>
          <el-button 
            v-if="row.status === 'draft'" 
            size="small" 
            type="success" 
            @click="publishVersion(row)"
            :icon="Upload"
          >发布</el-button>
          <el-button 
            v-if="row.status !== 'draft'" 
            size="small" 
            type="warning" 
            @click="rollbackVersion(row)"
            :icon="RefreshLeft"
          >回滚</el-button>
          <el-button size="small" @click="selectForCompare(row)" :icon="Switch">对比</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="size"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="loadVersions"
        @current-change="loadVersions"
      />
    </div>

    <!-- 创建版本对话框 -->
    <el-dialog v-model="createDialogVisible" title="创建新版本" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="版本号">
          <el-input v-model="createForm.version" placeholder="留空自动生成，如 1.0.0" />
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
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCreateVersion" :loading="submitting">创建</el-button>
      </template>
    </el-dialog>

    <!-- 版本详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="版本详情" width="700px">
      <div v-if="currentVersion" class="version-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="版本号">v{{ currentVersion.version }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(currentVersion.status)">{{ getStatusText(currentVersion.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建人">{{ currentVersion.created_by || '-' }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDate(currentVersion.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="发布时间" :span="2">
            {{ currentVersion.published_at ? formatDate(currentVersion.published_at) : '未发布' }}
          </el-descriptions-item>
          <el-descriptions-item label="变更说明" :span="2">
            {{ currentVersion.changelog || '无说明' }}
          </el-descriptions-item>
        </el-descriptions>

        <div class="schema-section">
          <h4>字段结构快照</h4>
          <el-table :data="parseSchema(currentVersion.schema_snapshot)" size="small">
            <el-table-column prop="name" label="字段名" width="150" />
            <el-table-column prop="display_name" label="显示名" width="150" />
            <el-table-column prop="type" label="类型" width="100" />
            <el-table-column prop="required" label="必填" width="80">
              <template #default="{ row }">
                <el-tag v-if="row.required" type="danger" size="small">是</el-tag>
                <el-tag v-else type="info" size="small">否</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="unique" label="唯一" width="80">
              <template #default="{ row }">
                <el-tag v-if="row.unique" type="warning" size="small">是</el-tag>
                <el-tag v-else type="info" size="small">否</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </el-dialog>

    <!-- 版本对比对话框 -->
    <el-dialog v-model="compareDialogVisible" title="版本对比" width="900px">
      <div v-if="compareResult" class="compare-container">
        <div class="compare-header">
          <div class="version-col">
            <span>版本 v{{ compareResult.version1?.version }}</span>
          </div>
          <div class="version-col">
            <span>版本 v{{ compareResult.version2?.version }}</span>
          </div>
        </div>

        <div v-if="compareResult.diff" class="diff-summary">
          <el-tag type="success" v-if="compareResult.diff.added?.length">
            新增 {{ compareResult.diff.added.length }} 个字段
          </el-tag>
          <el-tag type="danger" v-if="compareResult.diff.removed?.length">
            删除 {{ compareResult.diff.removed.length }} 个字段
          </el-tag>
          <el-tag type="warning" v-if="compareResult.diff.modified?.length">
            修改 {{ compareResult.diff.modified.length }} 个字段
          </el-tag>
          <el-tag type="info" v-if="!compareResult.diff.added?.length && !compareResult.diff.removed?.length && !compareResult.diff.modified?.length">
            无差异
          </el-tag>
        </div>

        <div class="compare-content">
          <div class="compare-col">
            <el-table :data="compareResult.version1?.fields || []" size="small">
              <el-table-column prop="name" label="字段名" />
              <el-table-column prop="type" label="类型" width="80" />
            </el-table>
          </div>
          <div class="compare-col">
            <el-table :data="compareResult.version2?.fields || []" size="small">
              <el-table-column prop="name" label="字段名" />
              <el-table-column prop="type" label="类型" width="80" />
            </el-table>
          </div>
        </div>
      </div>

      <div v-else class="compare-select">
        <p>请选择两个版本进行对比：</p>
        <el-form :inline="true">
          <el-form-item label="版本1">
            <el-select v-model="compareV1" placeholder="选择版本">
              <el-option 
                v-for="v in versions" 
                :key="v.id" 
                :label="'v' + v.version" 
                :value="v.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="版本2">
            <el-select v-model="compareV2" placeholder="选择版本">
              <el-option 
                v-for="v in versions" 
                :key="v.id" 
                :label="'v' + v.version" 
                :value="v.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="doCompare" :disabled="!compareV1 || !compareV2">对比</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Plus, Refresh, View, Upload, RefreshLeft, Switch } from '@element-plus/icons-vue'
import request from '@/utils/request'

const route = useRoute()
const router = useRouter()

const appId = ref(route.params.appId || route.query.appId)
const collectionId = ref(route.params.collectionId || route.query.collectionId)
const collectionName = ref(route.query.name || '数据模型')

const loading = ref(false)
const submitting = ref(false)
const versions = ref([])
const total = ref(0)
const page = ref(1)
const size = ref(20)

const createDialogVisible = ref(false)
const createForm = ref({
  version: '',
  changelog: ''
})

const detailDialogVisible = ref(false)
const currentVersion = ref(null)

const compareDialogVisible = ref(false)
const compareV1 = ref(null)
const compareV2 = ref(null)
const compareResult = ref(null)

// 加载版本列表
const loadVersions = async () => {
  loading.value = true
  try {
    const res = await request.get(`/baas/apps/${appId.value}/collections/${collectionId.value}/versions`, {
      params: { page: page.value, size: size.value }
    })
    // 响应拦截器已经解包，res直接是data内容
    if (res && res.list) {
      versions.value = res.list || []
      total.value = res.total || 0
    } else if (res && res.code === 0) {
      // 兼容未解包的响应格式
      versions.value = res.data?.list || []
      total.value = res.data?.total || 0
    }
  } catch (error) {
    console.error('加载版本列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 创建版本
const createVersion = () => {
  createForm.value = { version: '', changelog: '' }
  createDialogVisible.value = true
}

const submitCreateVersion = async () => {
  submitting.value = true
  try {
    const res = await request.post(`/baas/apps/${appId.value}/collections/${collectionId.value}/versions`, createForm.value)
    if (res.code === 0) {
      ElMessage.success('版本创建成功')
      createDialogVisible.value = false
      loadVersions()
    } else {
      ElMessage.error(res.message || '创建失败')
    }
  } catch (error) {
    ElMessage.error('创建失败')
  } finally {
    submitting.value = false
  }
}

// 查看版本详情
const viewVersion = (row) => {
  currentVersion.value = row
  detailDialogVisible.value = true
}

// 发布版本
const publishVersion = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要发布版本 v${row.version} 吗？发布后其他已发布版本将被标记为废弃。`, '确认发布')
    const res = await request.put(`/baas/apps/${appId.value}/collections/${collectionId.value}/versions/${row.id}/publish`)
    if (res.code === 0) {
      ElMessage.success('发布成功')
      loadVersions()
    } else {
      ElMessage.error(res.message || '发布失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('发布失败')
    }
  }
}

// 回滚版本
const rollbackVersion = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要回滚到版本 v${row.version} 吗？这将恢复数据模型的字段结构到该版本。`, '确认回滚')
    const res = await request.put(`/baas/apps/${appId.value}/collections/${collectionId.value}/versions/${row.id}/rollback`)
    if (res.code === 0) {
      ElMessage.success('回滚成功')
      loadVersions()
    } else {
      ElMessage.error(res.message || '回滚失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('回滚失败')
    }
  }
}

// 选择对比
const selectForCompare = (row) => {
  compareV1.value = null
  compareV2.value = row.id
  compareResult.value = null
  compareDialogVisible.value = true
}

// 执行对比
const doCompare = async () => {
  if (!compareV1.value || !compareV2.value) return
  try {
    const res = await request.get(`/baas/apps/${appId.value}/collections/${collectionId.value}/versions/compare`, {
      params: { v1: compareV1.value, v2: compareV2.value }
    })
    if (res.code === 0) {
      compareResult.value = res.data
    } else {
      ElMessage.error(res.message || '对比失败')
    }
  } catch (error) {
    ElMessage.error('对比失败')
  }
}

// 返回
const goBack = () => {
  router.back()
}

// 工具函数
const getStatusType = (status) => {
  const map = {
    draft: 'info',
    published: 'success',
    deprecated: 'warning'
  }
  return map[status] || 'info'
}

const getStatusText = (status) => {
  const map = {
    draft: '草稿',
    published: '已发布',
    deprecated: '已废弃'
  }
  return map[status] || status
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const parseSchema = (schema) => {
  if (!schema) return []
  try {
    return typeof schema === 'string' ? JSON.parse(schema) : schema
  } catch {
    return []
  }
}

onMounted(() => {
  loadVersions()
})
</script>

<style scoped>
.feature-versions {
  padding: 20px;
}

.header-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.collection-name {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.action-bar {
  margin-bottom: 16px;
}

.pagination-wrapper {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.version-detail {
  .schema-section {
    margin-top: 20px;
    
    h4 {
      margin-bottom: 12px;
      color: #303133;
    }
  }
}

.compare-container {
  .compare-header {
    display: flex;
    gap: 20px;
    margin-bottom: 16px;
    
    .version-col {
      flex: 1;
      text-align: center;
      font-weight: 600;
      padding: 8px;
      background: #f5f7fa;
      border-radius: 4px;
    }
  }
  
  .diff-summary {
    display: flex;
    gap: 8px;
    margin-bottom: 16px;
  }
  
  .compare-content {
    display: flex;
    gap: 20px;
    
    .compare-col {
      flex: 1;
    }
  }
}

.compare-select {
  text-align: center;
  padding: 20px;
  
  p {
    margin-bottom: 20px;
    color: #606266;
  }
}
</style>
