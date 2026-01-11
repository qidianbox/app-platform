<template>
  <div class="app-detail">
    <!-- 顶部导航栏 -->
    <header class="top-header" role="banner">
      <!-- 移动端汉堡菜单 -->
      <MobileMenu 
        v-model="mobileMenuOpen" 
        logo-text="拓" 
        app-name="拓客APP中台"
        @close="handleMobileMenuClose"
      >
        <!-- 移动端菜单内容 -->
        <div class="mobile-nav-tabs">
          <div 
            class="mobile-nav-item" 
            :class="{ active: activeTab === 'workspace' }"
            @click="switchMobileTab('workspace')"
          >
            <el-icon><Monitor /></el-icon>
            <span>工作台</span>
          </div>
          <div 
            class="mobile-nav-item" 
            :class="{ active: activeTab === 'config' }"
            @click="switchMobileTab('config')"
          >
            <el-icon><Setting /></el-icon>
            <span>配置中心</span>
          </div>
        </div>
        
        <!-- 工作台菜单 -->
        <div v-if="activeTab === 'workspace'" class="mobile-sidebar-menu">
          <div 
            v-for="item in workspaceMenuItems" 
            :key="item.key"
            class="mobile-menu-item"
            :class="{ active: workspaceMenu === item.key }"
            @click="switchWorkspaceMenu(item.key)"
          >
            <el-icon><component :is="item.icon" /></el-icon>
            <span>{{ item.label }}</span>
          </div>
        </div>
        
        <!-- 配置中心菜单 -->
        <div v-if="activeTab === 'config'" class="mobile-sidebar-menu">
          <div 
            class="mobile-menu-item"
            :class="{ active: currentPage === 'overview' }"
            @click="switchMobilePage('overview')"
          >
            <el-icon><House /></el-icon>
            <span>概览</span>
          </div>
          <div 
            class="mobile-menu-item"
            :class="{ active: currentPage === 'basic' }"
            @click="switchMobilePage('basic')"
          >
            <el-icon><Setting /></el-icon>
            <span>基础配置</span>
          </div>
          <template v-for="group in moduleGroups" :key="group.key">
            <div v-if="hasModulesInGroup(group.key)" class="mobile-menu-group">
              <div class="mobile-group-title">
                <el-icon><component :is="group.icon" /></el-icon>
                <span>{{ group.name }}</span>
              </div>
              <div 
                v-for="module in getModulesInGroup(group.key)" 
                :key="module.source_module"
                class="mobile-menu-item sub-item"
                :class="{ active: currentPage === module.source_module }"
                @click="switchMobilePage(module.source_module)"
              >
                <span>{{ module.name }}</span>
              </div>
            </div>
          </template>
        </div>
        
        <template #footer>
          <div 
            class="mobile-menu-item back-item" 
            @click="goBackToList"
          >
            <el-icon><ArrowLeft /></el-icon>
            <span>返回APP列表</span>
          </div>
        </template>
      </MobileMenu>
      
      <!-- 左侧：Logo + APP信息 + 工作台/配置中心 Tab -->
      <div class="header-left">
        <div class="header-logo">
          <div class="app-icon">
            <span>拓</span>
          </div>
          <span class="app-name">拓客APP中台</span>
        </div>
        
        <!-- 工作台 | 配置中心 Tab -->
        <div class="header-nav">
          <div 
            class="nav-item" 
            :class="{ active: activeTab === 'workspace' }"
            @click="activeTab = 'workspace'"
          >
            工作台
          </div>
          <div 
            class="nav-item" 
            :class="{ active: activeTab === 'config' }"
            @click="activeTab = 'config'"
          >
            配置中心
          </div>
        </div>
      </div>
      
      <!-- 右侧空白区域 -->
      <div class="header-right"></div>
    </header>

    <div class="main-container">
      <!-- 左侧边栏 - 仅在配置中心模式显示 -->
      <aside 
        class="sidebar" 
        v-show="activeTab === 'config'"
        role="navigation"
        aria-label="配置中心导航"
      >
        <nav class="sidebar-menu" role="menu">
          <!-- 概览 -->
          <div 
            class="sidebar-item"
            :class="{ active: currentPage === 'overview' }"
            @click="switchPage('overview')"
            role="menuitem"
            tabindex="0"
            aria-label="概览"
            @keydown.enter="switchPage('overview')"
          >
            <el-icon aria-hidden="true"><House /></el-icon>
            <span>概览</span>
          </div>
          
          <!-- 基础配置 -->
          <div 
            class="sidebar-item"
            :class="{ active: currentPage === 'basic' }"
            @click="switchPage('basic')"
            role="menuitem"
            tabindex="0"
            aria-label="基础配置"
            @keydown.enter="switchPage('basic')"
          >
            <el-icon aria-hidden="true"><Setting /></el-icon>
            <span>基础配置</span>
          </div>
          
          <!-- 菜单管理 -->
          <div 
            class="sidebar-item"
            :class="{ active: currentPage === 'menu_management' }"
            @click="switchPage('menu_management')"
            role="menuitem"
            tabindex="0"
            aria-label="菜单管理"
            @keydown.enter="switchPage('menu_management')"
          >
            <el-icon aria-hidden="true"><Menu /></el-icon>
            <span>菜单管理</span>
          </div>

        <!-- 模块分组 -->
        <template v-for="group in moduleGroups" :key="group.key">
          <div 
            v-if="hasModulesInGroup(group.key)"
            class="sidebar-group"
          >
            <div 
              class="group-header"
              @click="toggleGroup(group.key)"
            >
              <div class="group-title">
                <el-icon><component :is="group.icon" /></el-icon>
                <span>{{ group.name }}</span>
              </div>
              <el-icon class="expand-icon" :class="{ expanded: expandedGroups.includes(group.key) }">
                <ArrowRight />
              </el-icon>
            </div>
            <div v-show="expandedGroups.includes(group.key)" class="group-items">
              <div 
                v-for="module in getModulesInGroup(group.key)" 
                :key="module.source_module"
                class="sidebar-item sub-item"
                :class="{ active: currentPage === module.source_module }"
                @click="switchPage(module.source_module)"
              >
                <span>{{ module.name }}</span>
              </div>
            </div>
          </div>
        </template>
        </nav>
        <div class="sidebar-footer">
          <div 
            class="sidebar-item back-item" 
            @click="$router.push('/apps')"
            role="button"
            tabindex="0"
            aria-label="返回APP列表"
            @keydown.enter="$router.push('/apps')"
          >
            <el-icon aria-hidden="true"><ArrowLeft /></el-icon>
            <span>返回APP列表</span>
          </div>
        </div>
      </aside>

      <!-- 右侧内容区 -->
      <div class="content-area">
        <!-- 工作台模式 -->
        <template v-if="activeTab === 'workspace'">
          <!-- 只有当appId有效时才渲染Workspace组件，避免空的appId触发API请求 -->
          <Workspace v-if="appId" :app-id="appId" :app-info="appInfo" :initial-menu="workspaceMenu" />
          <div v-else class="loading-placeholder">
            <el-skeleton :rows="5" animated />
          </div>
        </template>

        <!-- 配置中心模式 -->
        <template v-else-if="activeTab === 'config'">
        <!-- 概览页面 -->
        <div v-if="currentPage === 'overview'" class="page-content">
          <h2 class="page-title">APP概览</h2>
          
          <!-- 统计卡片 -->
          <div class="stats-cards">
            <div class="stat-card">
              <div class="stat-icon users"><el-icon><User /></el-icon></div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.userCount }}</div>
                <div class="stat-label">用户数</div>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon modules"><el-icon><Grid /></el-icon></div>
              <div class="stat-info">
                <div class="stat-value">{{ appModules.length }}</div>
                <div class="stat-label">启用模块</div>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon requests"><el-icon><DataLine /></el-icon></div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.todayRequests }}</div>
                <div class="stat-label">今日请求</div>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon errors"><el-icon><Warning /></el-icon></div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.todayErrors }}</div>
                <div class="stat-label">今日异常</div>
              </div>
            </div>
          </div>

          <!-- APP信息 -->
          <div class="info-section">
            <h3>APP信息</h3>
            <div class="info-grid">
              <div class="info-item">
                <label>APP名称</label>
                <span>{{ appInfo.name }}</span>
              </div>
              <div class="info-item">
                <label>APP标识</label>
                <span class="copyable" @click="copyText(appInfo.app_id)">
                  {{ appInfo.app_id }}
                  <el-icon><CopyDocument /></el-icon>
                </span>
              </div>
              <div class="info-item">
                <label>AppSecret</label>
                <span class="copyable" @click="copyText(appInfo.app_secret)">
                  {{ maskSecret(appInfo.app_secret) }}
                  <el-icon><CopyDocument /></el-icon>
                </span>
              </div>
              <div class="info-item">
                <label>包名</label>
                <span>{{ appInfo.package_name || '-' }}</span>
              </div>
              <div class="info-item">
                <label>状态</label>
                <el-tag :type="appInfo.status === 1 ? 'success' : 'info'" size="small">
                  {{ appInfo.status === 1 ? '正常' : '禁用' }}
                </el-tag>
              </div>
              <div class="info-item">
                <label>创建时间</label>
                <span>{{ formatDate(appInfo.created_at) }}</span>
              </div>
            </div>
          </div>

          <!-- 已启用模块 -->
          <div class="info-section">
            <h3>已启用模块</h3>
            <div class="module-tags">
              <el-tag 
                v-for="module in appModules" 
                :key="module.id"
                type="primary"
                effect="plain"
              >
                {{ module.module_name || moduleNameMap[module.module_code] || module.name }}
              </el-tag>
              <el-empty v-if="appModules.length === 0" description="暂无启用模块" />
            </div>
          </div>
        </div>

        <!-- BaaS数据模型管理页面 -->
        <div v-else-if="currentPage === 'baas_data'" class="page-content">
          <div class="page-header">
            <div>
              <h2 class="page-title">数据模型管理</h2>
              <p class="page-desc">定义数据结构，自动生成CRUD API</p>
            </div>
            <el-button type="primary" @click="showCreateCollectionDialog = true">
              <el-icon><Plus /></el-icon>
              新建数据模型
            </el-button>
          </div>
          
          <!-- 搜索栏 -->
          <div class="search-bar">
            <el-input 
              v-model="collectionSearch" 
              placeholder="搜索数据模型" 
              style="width: 300px;"
              clearable
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
          
          <!-- 数据模型列表 -->
          <div class="collection-list">
            <div v-if="filteredCollections.length === 0" class="empty-state">
              <el-empty description="暂无数据模型">
                <el-button type="primary" @click="showCreateCollectionDialog = true">创建第一个数据模型</el-button>
              </el-empty>
            </div>
            <div v-else class="collection-grid">
              <div v-for="collection in filteredCollections" :key="collection.id" class="collection-card">
                <div class="card-header">
                  <div class="card-icon">
                    <el-icon size="24"><Grid /></el-icon>
                  </div>
                  <div class="card-info">
                    <h3>{{ collection.display_name || collection.name }}</h3>
                    <span class="card-name">{{ collection.name }}</span>
                  </div>
                  <el-tag :type="collection.is_generated ? 'success' : 'warning'" size="small">
                    {{ collection.is_generated ? '已生成' : '未生成' }}
                  </el-tag>
                </div>
                <p class="card-desc">{{ collection.description || '暂无描述' }}</p>
                <div class="card-stats">
                  <span>字段数: {{ (collection.fields || []).length }}</span>
                  <span>文档数: {{ collection.document_count || 0 }}</span>
                </div>
                <div class="card-actions">
                  <el-button 
                    v-if="!collection.is_generated" 
                    size="small" 
                    type="success" 
                    @click="generateFeature(collection)"
                  >
                    <el-icon><MagicStick /></el-icon>生成功能
                  </el-button>
                  <el-button 
                    v-else 
                    size="small" 
                    type="primary" 
                    @click="goToWorkspace(collection)"
                  >
                    <el-icon><View /></el-icon>查看功能
                  </el-button>
                  <el-button size="small" @click="editCollection(collection)">
                    <el-icon><Edit /></el-icon>编辑
                  </el-button>
                  <el-button size="small" @click="viewApiDoc(collection)">
                    <el-icon><Document /></el-icon>API文档
                  </el-button>
                  <el-button size="small" type="danger" @click="deleteCollection(collection)">
                    <el-icon><Delete /></el-icon>删除
                  </el-button>
                </div>
                <div class="card-api">
                  <span class="api-label">API端点:</span>
                  <code>/api/v1/baas/apps/{{ appId }}/data/{{ collection.name }}</code>
                  <el-button link size="small" @click="copyApiEndpoint(collection)">
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- API生成器页面 -->
        <div v-else-if="currentPage === 'api_generator'" class="page-content">
          <div class="page-header">
            <div>
              <h2 class="page-title">API生成器</h2>
              <p class="page-desc">基于数据模型自动生成RESTful API</p>
            </div>
          </div>
          
          <!-- 步骤导航 -->
          <div class="generator-steps">
            <el-steps :active="apiGeneratorStep" finish-status="success" align-center>
              <el-step title="选择数据模型" />
              <el-step title="配置API选项" />
              <el-step title="预览与生成" />
            </el-steps>
          </div>
          
          <!-- 步骤1: 选择数据模型 -->
          <div v-if="apiGeneratorStep === 0" class="generator-content">
            <h3>选择要生成API的数据模型</h3>
            <div class="model-select-grid">
              <div 
                v-for="collection in collections" 
                :key="collection.id" 
                class="model-select-card"
                :class="{ selected: selectedApiModels.includes(collection.id) }"
                @click="toggleApiModel(collection.id)"
              >
                <el-checkbox :model-value="selectedApiModels.includes(collection.id)" />
                <div class="model-info">
                  <h4>{{ collection.display_name || collection.name }}</h4>
                  <span class="model-name">{{ collection.name }}</span>
                  <p>{{ collection.description || '暂无描述' }}</p>
                  <div class="model-stats">
                    <span>字段数: {{ (collection.fields || []).length }}</span>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="collections.length === 0" class="empty-state">
              <el-empty description="暂无数据模型，请先在BaaS数据服务中创建">
                <el-button type="primary" @click="switchPage('baas_data')">+创建数据模型</el-button>
              </el-empty>
            </div>
            <div class="step-actions">
              <el-button type="primary" :disabled="selectedApiModels.length === 0" @click="apiGeneratorStep = 1">
                下一步
              </el-button>
            </div>
          </div>
          
          <!-- 步骤2: 配置API选项 -->
          <div v-else-if="apiGeneratorStep === 1" class="generator-content">
            <h3>配置API生成选项</h3>
            <el-form :model="apiGeneratorConfig" label-width="160px" class="config-form">
              <div class="form-section">
                <h4>接口配置</h4>
                <el-form-item label="API前缀">
                  <el-input v-model="apiGeneratorConfig.prefix" placeholder="/api/v1" />
                  <span class="form-hint">所有生成的API都会以此前缀开头</span>
                </el-form-item>
                <el-form-item label="生成接口">
                  <el-checkbox-group v-model="apiGeneratorConfig.endpoints">
                    <el-checkbox label="list">列表查询 (GET /list)</el-checkbox>
                    <el-checkbox label="detail">详情查询 (GET /:id)</el-checkbox>
                    <el-checkbox label="create">创建 (POST /)</el-checkbox>
                    <el-checkbox label="update">更新 (PUT /:id)</el-checkbox>
                    <el-checkbox label="delete">删除 (DELETE /:id)</el-checkbox>
                    <el-checkbox label="batch_delete">批量删除 (POST /batch-delete)</el-checkbox>
                  </el-checkbox-group>
                </el-form-item>
              </div>
              <div class="form-section">
                <h4>高级选项</h4>
                <el-form-item label="启用分页">
                  <el-switch v-model="apiGeneratorConfig.pagination" />
                </el-form-item>
                <el-form-item label="启用排序">
                  <el-switch v-model="apiGeneratorConfig.sorting" />
                </el-form-item>
                <el-form-item label="启用筛选">
                  <el-switch v-model="apiGeneratorConfig.filtering" />
                </el-form-item>
                <el-form-item label="启用认证">
                  <el-switch v-model="apiGeneratorConfig.authentication" />
                  <span class="form-hint">启用后需要携带Token访问</span>
                </el-form-item>
              </div>
            </el-form>
            <div class="step-actions">
              <el-button @click="apiGeneratorStep = 0">上一步</el-button>
              <el-button type="primary" @click="generateApiPreview">下一步</el-button>
            </div>
          </div>
          
          <!-- 步骤3: 预览与生成 -->
          <div v-else-if="apiGeneratorStep === 2" class="generator-content">
            <h3>API预览</h3>
            <div class="preview-section">
              <div class="preview-header">
                <span>生成的API列表</span>
                <el-button size="small" @click="copyAllApiCode">
                  <el-icon><CopyDocument /></el-icon>复制代码
                </el-button>
              </div>
              <div class="api-preview-list">
                <div v-for="api in generatedApis" :key="api.path" class="api-preview-item">
                  <el-tag :type="getMethodTagType(api.method)" size="small">{{ api.method }}</el-tag>
                  <code>{{ api.path }}</code>
                  <span class="api-desc">{{ api.description }}</span>
                </div>
              </div>
            </div>
            <div class="preview-section">
              <div class="preview-header">
                <span>后端代码预览 (Go)</span>
                <el-button size="small" @click="copyApiCode('go')">
                  <el-icon><CopyDocument /></el-icon>复制
                </el-button>
              </div>
              <pre class="code-preview"><code>{{ generatedGoCode }}</code></pre>
            </div>
            <div class="step-actions">
              <el-button @click="apiGeneratorStep = 1">上一步</el-button>
              <el-button type="primary" @click="downloadApiCode">
                <el-icon><Download /></el-icon>下载代码
              </el-button>
              <el-button type="success" @click="deployApi">
                <el-icon><Upload /></el-icon>部署API
              </el-button>
            </div>
          </div>
        </div>

        <!-- 页面生成器页面 -->
        <div v-else-if="currentPage === 'page_generator'" class="page-content">
          <div class="page-header">
            <div>
              <h2 class="page-title">页面生成器</h2>
              <p class="page-desc">基于数据模型自动生成前端页面</p>
            </div>
          </div>
          
          <!-- 步骤导航 -->
          <div class="generator-steps">
            <el-steps :active="pageGeneratorStep" finish-status="success" align-center>
              <el-step title="选择数据模型" />
              <el-step title="配置页面选项" />
              <el-step title="预览与生成" />
            </el-steps>
          </div>
          
          <!-- 步骤1: 选择数据模型 -->
          <div v-if="pageGeneratorStep === 0" class="generator-content">
            <h3>选择要生成页面的数据模型</h3>
            <div class="model-select-grid">
              <div 
                v-for="collection in collections" 
                :key="collection.id" 
                class="model-select-card"
                :class="{ selected: selectedPageModel === collection.id }"
                @click="selectedPageModel = collection.id"
              >
                <el-radio :model-value="selectedPageModel === collection.id" />
                <div class="model-info">
                  <h4>{{ collection.display_name || collection.name }}</h4>
                  <span class="model-name">{{ collection.name }}</span>
                  <p>{{ collection.description || '暂无描述' }}</p>
                </div>
              </div>
            </div>
            <div v-if="collections.length === 0" class="empty-state">
              <el-empty description="暂无数据模型，请先在BaaS数据服务中创建">
                <el-button type="primary" @click="switchPage('baas_data')">+创建数据模型</el-button>
              </el-empty>
            </div>
            <div class="step-actions">
              <el-button type="primary" :disabled="!selectedPageModel" @click="pageGeneratorStep = 1">
                下一步
              </el-button>
            </div>
          </div>
          
          <!-- 步骤2: 配置页面选项 -->
          <div v-else-if="pageGeneratorStep === 1" class="generator-content">
            <h3>配置页面生成选项</h3>
            <el-form :model="pageGeneratorConfig" label-width="160px" class="config-form">
              <div class="form-section">
                <h4>页面类型</h4>
                <el-form-item label="生成页面">
                  <el-checkbox-group v-model="pageGeneratorConfig.pages">
                    <el-checkbox label="list">列表页</el-checkbox>
                    <el-checkbox label="form">表单页（新增/编辑）</el-checkbox>
                    <el-checkbox label="detail">详情页</el-checkbox>
                  </el-checkbox-group>
                </el-form-item>
              </div>
              <div class="form-section">
                <h4>列表页配置</h4>
                <el-form-item label="显示字段">
                  <el-checkbox-group v-model="pageGeneratorConfig.listFields">
                    <el-checkbox 
                      v-for="field in selectedModelFields" 
                      :key="field.name" 
                      :label="field.name"
                    >
                      {{ field.display_name || field.name }}
                    </el-checkbox>
                  </el-checkbox-group>
                </el-form-item>
                <el-form-item label="启用搜索">
                  <el-switch v-model="pageGeneratorConfig.enableSearch" />
                </el-form-item>
                <el-form-item label="启用分页">
                  <el-switch v-model="pageGeneratorConfig.enablePagination" />
                </el-form-item>
              </div>
              <div class="form-section">
                <h4>样式配置</h4>
                <el-form-item label="UI框架">
                  <el-select v-model="pageGeneratorConfig.uiFramework">
                    <el-option label="Element Plus" value="element-plus" />
                    <el-option label="Ant Design Vue" value="ant-design" />
                    <el-option label="Naive UI" value="naive-ui" />
                  </el-select>
                </el-form-item>
                <el-form-item label="布局样式">
                  <el-select v-model="pageGeneratorConfig.layout">
                    <el-option label="卡片布局" value="card" />
                    <el-option label="表格布局" value="table" />
                  </el-select>
                </el-form-item>
              </div>
            </el-form>
            <div class="step-actions">
              <el-button @click="pageGeneratorStep = 0">上一步</el-button>
              <el-button type="primary" @click="generatePagePreview">下一步</el-button>
            </div>
          </div>
          
          <!-- 步骤3: 预览与生成 -->
          <div v-else-if="pageGeneratorStep === 2" class="generator-content">
            <h3>页面预览</h3>
            <div class="preview-tabs">
              <el-tabs v-model="pagePreviewTab">
                <el-tab-pane label="列表页" name="list" v-if="pageGeneratorConfig.pages.includes('list')">
                  <div class="page-preview-frame">
                    <div class="preview-toolbar">
                      <el-input placeholder="搜索..." style="width: 200px;" v-if="pageGeneratorConfig.enableSearch" />
                      <el-button type="primary">新增</el-button>
                    </div>
                    <el-table :data="previewTableData" style="width: 100%">
                      <el-table-column 
                        v-for="field in pageGeneratorConfig.listFields" 
                        :key="field" 
                        :prop="field" 
                        :label="getFieldLabel(field)"
                      />
                      <el-table-column label="操作" width="180">
                        <template #default>
                          <el-button size="small">编辑</el-button>
                          <el-button size="small" type="danger">删除</el-button>
                        </template>
                      </el-table-column>
                    </el-table>
                  </div>
                </el-tab-pane>
                <el-tab-pane label="表单页" name="form" v-if="pageGeneratorConfig.pages.includes('form')">
                  <div class="page-preview-frame">
                    <el-form label-width="120px">
                      <el-form-item 
                        v-for="field in selectedModelFields" 
                        :key="field.name" 
                        :label="field.display_name || field.name"
                      >
                        <el-input v-if="field.type === 'string'" :placeholder="`请输入${field.display_name || field.name}`" />
                        <el-input-number v-else-if="field.type === 'number'" />
                        <el-switch v-else-if="field.type === 'boolean'" />
                        <el-input v-else type="textarea" />
                      </el-form-item>
                      <el-form-item>
                        <el-button type="primary">保存</el-button>
                        <el-button>取消</el-button>
                      </el-form-item>
                    </el-form>
                  </div>
                </el-tab-pane>
              </el-tabs>
            </div>
            <div class="preview-section">
              <div class="preview-header">
                <span>Vue代码预览</span>
                <el-button size="small" @click="copyPageCode">
                  <el-icon><CopyDocument /></el-icon>复制
                </el-button>
              </div>
              <pre class="code-preview"><code>{{ generatedVueCode }}</code></pre>
            </div>
            <div class="step-actions">
              <el-button @click="pageGeneratorStep = 1">上一步</el-button>
              <el-button type="primary" @click="downloadPageCode">
                <el-icon><Download /></el-icon>下载代码
              </el-button>
            </div>
          </div>
        </div>

        <!-- 代码生成器页面 -->
        <div v-else-if="currentPage === 'code_generator'" class="page-content">
          <div class="page-header">
            <div>
              <h2 class="page-title">代码生成器</h2>
              <p class="page-desc">一键生成完整的前后端代码</p>
            </div>
          </div>
          
          <div class="generator-content">
            <div class="code-gen-options">
              <h3>选择要生成的内容</h3>
              
              <div class="option-section">
                <h4>数据模型</h4>
                <div class="model-checkbox-list">
                  <el-checkbox 
                    v-for="collection in collections" 
                    :key="collection.id"
                    :label="collection.id"
                    v-model="codeGeneratorConfig.selectedModels"
                  >
                    {{ collection.display_name || collection.name }}
                  </el-checkbox>
                </div>
                <el-checkbox v-model="codeGeneratorConfig.selectAll" @change="toggleSelectAllModels">
                  全选
                </el-checkbox>
              </div>
              
              <div class="option-section">
                <h4>生成内容</h4>
                <el-checkbox-group v-model="codeGeneratorConfig.generateItems">
                  <el-checkbox label="backend_api">后端API代码 (Go)</el-checkbox>
                  <el-checkbox label="frontend_pages">前端页面 (Vue)</el-checkbox>
                  <el-checkbox label="api_docs">API文档 (Swagger)</el-checkbox>
                  <el-checkbox label="database_sql">数据库脚本 (SQL)</el-checkbox>
                </el-checkbox-group>
              </div>
              
              <div class="option-section">
                <h4>后端配置</h4>
                <el-form label-width="140px">
                  <el-form-item label="后端框架">
                    <el-select v-model="codeGeneratorConfig.backendFramework">
                      <el-option label="Gin (Go)" value="gin" />
                      <el-option label="Echo (Go)" value="echo" />
                      <el-option label="Fiber (Go)" value="fiber" />
                    </el-select>
                  </el-form-item>
                  <el-form-item label="数据库">
                    <el-select v-model="codeGeneratorConfig.database">
                      <el-option label="MySQL" value="mysql" />
                      <el-option label="PostgreSQL" value="postgres" />
                      <el-option label="SQLite" value="sqlite" />
                    </el-select>
                  </el-form-item>
                </el-form>
              </div>
              
              <div class="option-section">
                <h4>前端配置</h4>
                <el-form label-width="140px">
                  <el-form-item label="前端框架">
                    <el-select v-model="codeGeneratorConfig.frontendFramework">
                      <el-option label="Vue 3" value="vue3" />
                      <el-option label="React" value="react" />
                    </el-select>
                  </el-form-item>
                  <el-form-item label="UI组件库">
                    <el-select v-model="codeGeneratorConfig.uiLibrary">
                      <el-option label="Element Plus" value="element-plus" />
                      <el-option label="Ant Design" value="antd" />
                      <el-option label="Naive UI" value="naive-ui" />
                    </el-select>
                  </el-form-item>
                </el-form>
              </div>
            </div>
            
            <div class="code-gen-actions">
              <el-button 
                type="primary" 
                size="large"
                :disabled="codeGeneratorConfig.selectedModels.length === 0 || codeGeneratorConfig.generateItems.length === 0"
                @click="generateFullCode"
              >
                <el-icon><Download /></el-icon>
                生成并下载代码
              </el-button>
              <p class="hint">将生成ZIP压缩包，包含所有选中的代码文件</p>
            </div>
          </div>
        </div>

        <!-- 基础配置页面 -->
        <div v-else-if="currentPage === 'basic'" class="page-content">
          <h2 class="page-title">基础配置</h2>
          <p class="page-desc">配置APP的基本信息和通用设置</p>
          
          <el-form :model="basicConfig" label-width="140px" class="config-form">
            <div class="form-section">
              <h4>基本信息</h4>
              <el-form-item label="APP名称">
                <el-input v-model="basicConfig.name" placeholder="请输入APP名称" />
              </el-form-item>
              <el-form-item label="APP描述">
                <el-input v-model="basicConfig.description" type="textarea" :rows="3" placeholder="请输入APP描述" />
              </el-form-item>
              <el-form-item label="包名">
                <el-input v-model="basicConfig.package_name" placeholder="如：com.example.app" />
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>安全设置</h4>
              <el-form-item label="启用签名验证">
                <el-switch v-model="basicConfig.enableSignature" />
                <span class="form-hint">启用后所有API请求需要携带签名</span>
              </el-form-item>
              <el-form-item label="IP白名单">
                <el-input v-model="basicConfig.ipWhitelist" type="textarea" :rows="2" placeholder="每行一个IP，留空表示不限制" />
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveBasicConfig">保存配置</el-button>
              <el-button @click="resetBasicConfig">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 菜单管理页面 -->
        <div v-else-if="currentPage === 'menu_management'" class="page-content">
          <div class="page-header">
            <div>
              <h2 class="page-title">菜单管理</h2>
              <p class="page-desc">管理工作台侧边栏显示的功能菜单，只有开启的菜单才会在工作台中显示</p>
            </div>
          </div>
          
          <!-- 已生成功能列表 -->
          <el-card class="menu-list-card">
            <template #header>
              <div class="card-header">
                <span>已生成的功能</span>
                <span class="hint-text">点击开关控制是否在工作台显示</span>
              </div>
            </template>
            
            <el-table :data="generatedCollections" style="width: 100%" v-loading="menuLoading">
              <el-table-column prop="display_name" label="功能名称" min-width="150">
                <template #default="{ row }">
                  <div class="menu-name">
                    <el-icon><Document /></el-icon>
                    <span>{{ row.display_name || row.name }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="name" label="模型标识" width="150" />
              <el-table-column prop="description" label="描述" min-width="200">
                <template #default="{ row }">
                  <span class="text-gray">{{ row.description || '-' }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="created_at" label="创建时间" width="180">
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>
              <el-table-column label="工作台显示" width="120" align="center">
                <template #default="{ row }">
                  <el-switch
                    v-model="row.is_visible"
                    @change="toggleMenuVisibility(row)"
                    :loading="row.loading"
                  />
                </template>
              </el-table-column>
            </el-table>
            
            <el-empty v-if="generatedCollections.length === 0 && !menuLoading" description="暂无已生成的功能">
              <el-button type="primary" @click="switchPage('baas_data')">去创建数据模型</el-button>
            </el-empty>
          </el-card>
          
          <div class="menu-tips">
            <el-alert
              title="使用说明"
              type="info"
              :closable="false"
            >
              <template #default>
                <ul class="tips-list">
                  <li>1. 在「数据模型管理」中创建数据模型并点击「生成功能」</li>
                  <li>2. 生成后的功能会出现在此列表中</li>
                  <li>3. 开启「工作台显示」开关，该功能将显示在工作台侧边栏</li>
                  <li>4. 关闭开关后，功能仍然存在，只是不在工作台显示</li>
                </ul>
              </template>
            </el-alert>
          </div>
        </div>

        <!-- 用户管理配置 -->
        <div v-else-if="currentPage === 'user_management'" class="page-content">
          <h2 class="page-title">用户管理配置</h2>
          <p class="page-desc">配置用户注册、登录和管理相关设置</p>
          
          <el-form :model="userConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>🔐 登录配置</h4>
              <el-form-item label="密码最小长度">
                <el-input-number v-model="userConfig.passwordMinLength" :min="6" :max="32" />
                <span class="form-hint">建议8位以上</span>
              </el-form-item>
              <el-form-item label="密码复杂度要求">
                <el-checkbox-group v-model="userConfig.passwordRequirements">
                  <el-checkbox label="number">必须包含数字</el-checkbox>
                  <el-checkbox label="letter">必须包含字母</el-checkbox>
                  <el-checkbox label="special">必须包含特殊字符</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item label="登录失败锁定">
                <el-switch v-model="userConfig.enableLoginLock" />
                <span class="form-hint">防止暴力破解</span>
              </el-form-item>
              <el-form-item v-if="userConfig.enableLoginLock" label="失败次数限制">
                <el-input-number v-model="userConfig.maxLoginAttempts" :min="3" :max="10" />
                <span class="form-hint">次</span>
              </el-form-item>
              <el-form-item v-if="userConfig.enableLoginLock" label="锁定时长">
                <el-input-number v-model="userConfig.lockDuration" :min="5" :max="1440" />
                <span class="form-hint">分钟</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>👤 用户信息管理</h4>
              <el-form-item label="必填字段">
                <el-checkbox-group v-model="userConfig.requiredFields">
                  <el-checkbox label="nickname">昵称</el-checkbox>
                  <el-checkbox label="avatar">头像</el-checkbox>
                  <el-checkbox label="gender">性别</el-checkbox>
                  <el-checkbox label="birthday">生日</el-checkbox>
                  <el-checkbox label="region">地区</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item label="允许修改用户名">
                <el-switch v-model="userConfig.allowChangeUsername" />
                <span class="form-hint">关闭后用户名不可修改</span>
              </el-form-item>
              <el-form-item label="昵称敏感词过滤">
                <el-switch v-model="userConfig.enableNicknameFilter" />
              </el-form-item>
              <el-form-item label="头像审核">
                <el-switch v-model="userConfig.enableAvatarReview" />
                <span class="form-hint">自动检测头像是否违规</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>🪪 实名认证配置</h4>
              <el-form-item label="启用实名认证">
                <el-switch v-model="userConfig.enableRealName" />
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>🗑️ 账号注销配置</h4>
              <el-form-item label="允许账号注销">
                <el-switch v-model="userConfig.allowAccountDeletion" />
              </el-form-item>
              <el-form-item v-if="userConfig.allowAccountDeletion" label="注销冷静期">
                <el-input-number v-model="userConfig.deletionCooldown" :min="0" :max="30" />
                <span class="form-hint">天，0表示立即注销</span>
              </el-form-item>
              <el-form-item v-if="userConfig.allowAccountDeletion" label="注销前置条件">
                <el-checkbox-group v-model="userConfig.deletionRequirements">
                  <el-checkbox label="clearData">清空个人数据</el-checkbox>
                  <el-checkbox label="unbindThirdParty">解除第三方账号</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item v-if="userConfig.allowAccountDeletion" label="注销确认方式">
                <el-radio-group v-model="userConfig.deletionConfirmMethod">
                  <el-radio label="sms">短信验证码</el-radio>
                  <el-radio label="password">密码验证</el-radio>
                </el-radio-group>
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('user_management')">保存配置</el-button>
              <el-button @click="testConfig('user_management')">测试配置</el-button>
              <el-button @click="resetConfig('user_management')">重置</el-button>
              <el-button @click="showConfigHistory('user_management')">查看历史</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 消息中心配置 -->
        <div v-else-if="currentPage === 'message_center'" class="page-content">
          <h2 class="page-title">消息中心配置</h2>
          <p class="page-desc">配置站内消息和通知相关设置</p>
          
          <el-form :model="messageConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>📬 基础配置</h4>
              <el-form-item label="启用消息服务">
                <el-switch v-model="messageConfig.enabled" />
              </el-form-item>
              <el-form-item label="消息保留天数">
                <el-input-number v-model="messageConfig.retentionDays" :min="7" :max="365" />
                <span class="form-hint">天</span>
              </el-form-item>
              <el-form-item label="单用户消息上限">
                <el-input-number v-model="messageConfig.maxMessagesPerUser" :min="100" :max="10000" />
                <span class="form-hint">条</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>📝 消息类型</h4>
              <el-form-item label="支持的消息类型">
                <el-checkbox-group v-model="messageConfig.supportedTypes">
                  <el-checkbox label="system">系统通知</el-checkbox>
                  <el-checkbox label="activity">活动消息</el-checkbox>
                  <el-checkbox label="transaction">交易消息</el-checkbox>
                  <el-checkbox label="social">社交消息</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('message_center')">保存配置</el-button>
              <el-button @click="resetConfig('message_center')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 推送服务配置 -->
        <div v-else-if="currentPage === 'push_service'" class="page-content">
          <h2 class="page-title">推送服务配置</h2>
          <p class="page-desc">配置APP推送通知服务</p>
          
          <el-form :model="pushConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>🔔 基础配置</h4>
              <el-form-item label="启用推送服务">
                <el-switch v-model="pushConfig.enabled" />
              </el-form-item>
              <el-form-item label="推送服务商">
                <el-select v-model="pushConfig.provider" placeholder="请选择">
                  <el-option label="极光推送" value="jpush" />
                  <el-option label="个推" value="getui" />
                  <el-option label="友盟推送" value="umeng" />
                  <el-option label="Firebase" value="firebase" />
                </el-select>
              </el-form-item>
              <el-form-item label="AppKey">
                <el-input v-model="pushConfig.appKey" placeholder="请输入AppKey" />
              </el-form-item>
              <el-form-item label="MasterSecret">
                <el-input v-model="pushConfig.masterSecret" type="password" placeholder="请输入MasterSecret" show-password />
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>⏰ 推送策略</h4>
              <el-form-item label="静默时段">
                <el-switch v-model="pushConfig.enableQuietHours" />
                <span class="form-hint">在指定时段不发送推送</span>
              </el-form-item>
              <el-form-item v-if="pushConfig.enableQuietHours" label="静默时间">
                <el-time-picker v-model="pushConfig.quietStart" placeholder="开始时间" format="HH:mm" />
                <span style="margin: 0 8px;">至</span>
                <el-time-picker v-model="pushConfig.quietEnd" placeholder="结束时间" format="HH:mm" />
              </el-form-item>
              <el-form-item label="每日推送上限">
                <el-input-number v-model="pushConfig.dailyLimit" :min="1" :max="100" />
                <span class="form-hint">条/用户</span>
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('push_service')">保存配置</el-button>
              <el-button @click="testConfig('push_service')">测试推送</el-button>
              <el-button @click="resetConfig('push_service')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 支付中心配置 -->
        <div v-else-if="currentPage === 'payment'" class="page-content">
          <h2 class="page-title">支付中心配置</h2>
          <p class="page-desc">配置支付渠道和安全设置</p>
          
          <el-form :model="paymentConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>🔐 安全验证配置</h4>
              <el-form-item label="启用安全验证">
                <el-switch v-model="paymentConfig.enableSecurityVerify" />
              </el-form-item>
              <el-form-item v-if="paymentConfig.enableSecurityVerify" label="验证方式">
                <el-checkbox-group v-model="paymentConfig.verifyMethods">
                  <el-checkbox label="password">支付密码</el-checkbox>
                  <el-checkbox label="fingerprint">指纹识别</el-checkbox>
                  <el-checkbox label="face">面容识别</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item label="验证触发金额">
                <el-input-number v-model="paymentConfig.verifyThreshold" :min="0" :max="100000" />
                <span class="form-hint">元，0表示所有支付都需要验证</span>
              </el-form-item>
              <el-form-item label="密码错误锁定">
                <el-input-number v-model="paymentConfig.maxPasswordAttempts" :min="3" :max="10" />
                <span class="form-hint">次</span>
              </el-form-item>
              <el-form-item label="锁定时长">
                <el-input-number v-model="paymentConfig.lockDuration" :min="5" :max="1440" />
                <span class="form-hint">分钟</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>💰 限额控制配置</h4>
              <el-form-item label="启用限额控制">
                <el-switch v-model="paymentConfig.enableLimitControl" />
              </el-form-item>
              <el-form-item v-if="paymentConfig.enableLimitControl" label="单笔支付限额">
                <el-input-number v-model="paymentConfig.singleLimit" :min="0" :max="1000000" />
                <span class="form-hint">元，0表示不限制</span>
              </el-form-item>
              <el-form-item v-if="paymentConfig.enableLimitControl" label="每日支付限额">
                <el-input-number v-model="paymentConfig.dailyLimit" :min="0" :max="10000000" />
                <span class="form-hint">元，0表示不限制</span>
              </el-form-item>
              <el-form-item v-if="paymentConfig.enableLimitControl" label="每月支付限额">
                <el-input-number v-model="paymentConfig.monthlyLimit" :min="0" :max="100000000" />
                <span class="form-hint">元，0表示不限制</span>
              </el-form-item>
              <el-form-item v-if="paymentConfig.enableLimitControl" label="每日支付次数">
                <el-input-number v-model="paymentConfig.dailyCount" :min="0" :max="1000" />
                <span class="form-hint">次，0表示不限制</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>🔗 回调配置</h4>
              <el-form-item label="支付成功回调">
                <el-input v-model="paymentConfig.successCallback" placeholder="请输入支付成功回调地址" />
              </el-form-item>
              <el-form-item label="支付失败回调">
                <el-input v-model="paymentConfig.failCallback" placeholder="请输入支付失败回调地址" />
              </el-form-item>
              <el-form-item label="退款回调">
                <el-input v-model="paymentConfig.refundCallback" placeholder="请输入退款回调地址" />
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>⚙️ 其他配置</h4>
              <el-form-item label="支付超时时间">
                <el-input-number v-model="paymentConfig.timeout" :min="5" :max="60" />
                <span class="form-hint">分钟</span>
              </el-form-item>
              <el-form-item label="启用自动退款">
                <el-switch v-model="paymentConfig.enableAutoRefund" />
                <span class="form-hint">订单超时自动退款</span>
              </el-form-item>
              <el-form-item label="启用支付日志">
                <el-switch v-model="paymentConfig.enablePaymentLog" />
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('payment')">保存配置</el-button>
              <el-button @click="testConfig('payment')">测试配置</el-button>
              <el-button @click="resetConfig('payment')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 短信服务配置 -->
        <div v-else-if="currentPage === 'sms_service'" class="page-content">
          <h2 class="page-title">短信服务配置</h2>
          <p class="page-desc">配置短信发送服务和验证码设置</p>
          
          <el-form :model="smsConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>📱 基础配置</h4>
              <el-form-item label="启用短信服务">
                <el-switch v-model="smsConfig.enabled" />
              </el-form-item>
              <el-form-item label="短信服务提供商">
                <el-select v-model="smsConfig.provider" placeholder="请选择">
                  <el-option label="阿里云短信" value="aliyun" />
                  <el-option label="腾讯云短信" value="tencent" />
                  <el-option label="华为云短信" value="huawei" />
                </el-select>
              </el-form-item>
              <el-form-item label="AccessKey">
                <el-input v-model="smsConfig.accessKey" placeholder="请输入AccessKey" />
              </el-form-item>
              <el-form-item label="SecretKey">
                <el-input v-model="smsConfig.secretKey" type="password" placeholder="请输入SecretKey" show-password />
              </el-form-item>
              <el-form-item label="短信签名">
                <el-input v-model="smsConfig.signName" placeholder="例如：我的应用" />
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>🔢 验证码短信配置</h4>
              <el-form-item label="验证码长度">
                <el-input-number v-model="smsConfig.codeLength" :min="4" :max="8" />
                <span class="form-hint">位</span>
              </el-form-item>
              <el-form-item label="验证码有效期">
                <el-input-number v-model="smsConfig.codeExpiry" :min="1" :max="30" />
                <span class="form-hint">分钟</span>
              </el-form-item>
              <el-form-item label="验证码模板ID">
                <el-input v-model="smsConfig.codeTemplateId" placeholder="例如：SMS_123456789" />
              </el-form-item>
              <el-form-item label="发送间隔">
                <el-input-number v-model="smsConfig.sendInterval" :min="30" :max="300" />
                <span class="form-hint">秒</span>
              </el-form-item>
              <el-form-item label="每日发送限制">
                <el-input-number v-model="smsConfig.dailyLimit" :min="1" :max="50" />
                <span class="form-hint">条/手机号</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>📢 通知短信配置</h4>
              <el-form-item label="启用通知短信">
                <el-switch v-model="smsConfig.enableNotification" />
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>⚙️ 高级配置</h4>
              <el-form-item label="失败重试次数">
                <el-input-number v-model="smsConfig.retryCount" :min="0" :max="5" />
                <span class="form-hint">次</span>
              </el-form-item>
              <el-form-item label="请求超时时间">
                <el-input-number v-model="smsConfig.timeout" :min="5" :max="60" />
                <span class="form-hint">秒</span>
              </el-form-item>
              <el-form-item label="状态回调URL">
                <el-input v-model="smsConfig.callbackUrl" placeholder="请输入状态回调地址" />
              </el-form-item>
              <el-form-item label="余额告警阈值">
                <el-input-number v-model="smsConfig.balanceAlert" :min="100" :max="10000" />
                <span class="form-hint">条</span>
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('sms_service')">保存配置</el-button>
              <el-button @click="testConfig('sms_service')">测试发送</el-button>
              <el-button @click="resetConfig('sms_service')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 数据埋点配置 -->
        <div v-else-if="currentPage === 'data_tracking'" class="page-content">
          <h2 class="page-title">数据埋点配置</h2>
          <p class="page-desc">配置用户行为埋点和数据分析</p>
          
          <el-form :model="trackingConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>📊 基础配置</h4>
              <el-form-item label="启用数据埋点">
                <el-switch v-model="trackingConfig.enabled" />
              </el-form-item>
              <el-form-item label="数据上报方式">
                <el-radio-group v-model="trackingConfig.reportMethod">
                  <el-radio label="realtime">实时上报</el-radio>
                  <el-radio label="batch">批量上报</el-radio>
                </el-radio-group>
              </el-form-item>
              <el-form-item v-if="trackingConfig.reportMethod === 'batch'" label="批量上报间隔">
                <el-input-number v-model="trackingConfig.batchInterval" :min="10" :max="300" />
                <span class="form-hint">秒</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>🎯 事件配置</h4>
              <el-form-item label="自动采集事件">
                <el-checkbox-group v-model="trackingConfig.autoEvents">
                  <el-checkbox label="pageView">页面浏览</el-checkbox>
                  <el-checkbox label="click">点击事件</el-checkbox>
                  <el-checkbox label="scroll">滚动事件</el-checkbox>
                  <el-checkbox label="error">错误事件</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('data_tracking')">保存配置</el-button>
              <el-button @click="resetConfig('data_tracking')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 日志服务配置 -->
        <div v-else-if="currentPage === 'log_service'" class="page-content">
          <h2 class="page-title">日志服务配置</h2>
          <p class="page-desc">配置日志收集、存储、分析等功能</p>
          
          <el-form :model="logConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>📝 基础配置</h4>
              <el-form-item label="启用日志服务">
                <el-switch v-model="logConfig.enabled" />
              </el-form-item>
              <el-form-item label="日志级别">
                <el-select v-model="logConfig.level" placeholder="请选择">
                  <el-option label="DEBUG" value="debug" />
                  <el-option label="INFO" value="info" />
                  <el-option label="WARN" value="warn" />
                  <el-option label="ERROR" value="error" />
                </el-select>
              </el-form-item>
              <el-form-item label="日志存储方式">
                <el-checkbox-group v-model="logConfig.storageTypes">
                  <el-checkbox label="local">本地存储</el-checkbox>
                  <el-checkbox label="cloud">云端存储</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item label="日志保留时间">
                <el-input-number v-model="logConfig.retentionDays" :min="7" :max="365" />
                <span class="form-hint">天</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>📤 上报配置</h4>
              <el-form-item label="实时上报">
                <el-switch v-model="logConfig.realtimeReport" />
              </el-form-item>
              <el-form-item label="批量上报数量">
                <el-input-number v-model="logConfig.batchSize" :min="10" :max="1000" />
                <span class="form-hint">条</span>
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('log_service')">保存配置</el-button>
              <el-button @click="resetConfig('log_service')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 监控告警配置 -->
        <div v-else-if="currentPage === 'monitor_alert'" class="page-content">
          <h2 class="page-title">监控告警配置</h2>
          <p class="page-desc">配置应用监控和告警通知</p>
          
          <el-form :model="monitorConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>📡 监控配置</h4>
              <el-form-item label="启用监控服务">
                <el-switch v-model="monitorConfig.enabled" />
              </el-form-item>
              <el-form-item label="监控指标">
                <el-checkbox-group v-model="monitorConfig.metrics">
                  <el-checkbox label="cpu">CPU使用率</el-checkbox>
                  <el-checkbox label="memory">内存使用率</el-checkbox>
                  <el-checkbox label="api">API响应时间</el-checkbox>
                  <el-checkbox label="error">错误率</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item label="采集间隔">
                <el-input-number v-model="monitorConfig.interval" :min="10" :max="300" />
                <span class="form-hint">秒</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>🚨 告警配置</h4>
              <el-form-item label="启用告警">
                <el-switch v-model="monitorConfig.alertEnabled" />
              </el-form-item>
              <el-form-item v-if="monitorConfig.alertEnabled" label="告警方式">
                <el-checkbox-group v-model="monitorConfig.alertMethods">
                  <el-checkbox label="email">邮件</el-checkbox>
                  <el-checkbox label="sms">短信</el-checkbox>
                  <el-checkbox label="webhook">Webhook</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item v-if="monitorConfig.alertEnabled" label="告警接收人">
                <el-input v-model="monitorConfig.alertReceivers" type="textarea" :rows="2" placeholder="多个接收人用逗号分隔" />
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('monitor_alert')">保存配置</el-button>
              <el-button @click="testConfig('monitor_alert')">测试告警</el-button>
              <el-button @click="resetConfig('monitor_alert')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 文件存储配置 -->
        <div v-else-if="currentPage === 'file_storage'" class="page-content">
          <h2 class="page-title">文件存储配置</h2>
          <p class="page-desc">配置文件上传、下载和存储服务</p>
          
          <el-form :model="storageConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>☁️ 存储配置</h4>
              <el-form-item label="启用文件存储">
                <el-switch v-model="storageConfig.enabled" />
              </el-form-item>
              <el-form-item label="存储服务商">
                <el-select v-model="storageConfig.provider" placeholder="请选择">
                  <el-option label="阿里云OSS" value="aliyun" />
                  <el-option label="腾讯云COS" value="tencent" />
                  <el-option label="七牛云" value="qiniu" />
                  <el-option label="AWS S3" value="aws" />
                </el-select>
              </el-form-item>
              <el-form-item label="Bucket名称">
                <el-input v-model="storageConfig.bucket" placeholder="请输入Bucket名称" />
              </el-form-item>
              <el-form-item label="AccessKey">
                <el-input v-model="storageConfig.accessKey" placeholder="请输入AccessKey" />
              </el-form-item>
              <el-form-item label="SecretKey">
                <el-input v-model="storageConfig.secretKey" type="password" placeholder="请输入SecretKey" show-password />
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>📁 上传限制</h4>
              <el-form-item label="最大文件大小">
                <el-input-number v-model="storageConfig.maxFileSize" :min="1" :max="1024" />
                <span class="form-hint">MB</span>
              </el-form-item>
              <el-form-item label="允许的文件类型">
                <el-input v-model="storageConfig.allowedTypes" placeholder="例如：jpg,png,pdf" />
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('file_storage')">保存配置</el-button>
              <el-button @click="testConfig('file_storage')">测试连接</el-button>
              <el-button @click="resetConfig('file_storage')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 配置管理配置 -->
        <div v-else-if="currentPage === 'config_management'" class="page-content">
          <h2 class="page-title">配置管理</h2>
          <p class="page-desc">管理远程配置下发和动态配置</p>
          
          <el-form :model="configMgmtConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>⚙️ 基础配置</h4>
              <el-form-item label="启用配置管理">
                <el-switch v-model="configMgmtConfig.enabled" />
              </el-form-item>
              <el-form-item label="配置刷新间隔">
                <el-input-number v-model="configMgmtConfig.refreshInterval" :min="60" :max="3600" />
                <span class="form-hint">秒</span>
              </el-form-item>
              <el-form-item label="启用配置缓存">
                <el-switch v-model="configMgmtConfig.enableCache" />
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('config_management')">保存配置</el-button>
              <el-button @click="resetConfig('config_management')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 版本管理配置 -->
        <div v-else-if="currentPage === 'version_management'" class="page-content">
          <h2 class="page-title">版本管理配置</h2>
          <p class="page-desc">配置APP版本发布和更新策略</p>
          
          <el-form :model="versionConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>📦 更新配置</h4>
              <el-form-item label="启用版本管理">
                <el-switch v-model="versionConfig.enabled" />
              </el-form-item>
              <el-form-item label="强制更新">
                <el-switch v-model="versionConfig.forceUpdate" />
                <span class="form-hint">开启后用户必须更新到最新版本</span>
              </el-form-item>
              <el-form-item label="更新提示方式">
                <el-radio-group v-model="versionConfig.promptType">
                  <el-radio label="dialog">弹窗提示</el-radio>
                  <el-radio label="toast">轻提示</el-radio>
                  <el-radio label="silent">静默更新</el-radio>
                </el-radio-group>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>📥 下载配置</h4>
              <el-form-item label="Android下载地址">
                <el-input v-model="versionConfig.androidUrl" placeholder="请输入Android安装包下载地址" />
              </el-form-item>
              <el-form-item label="iOS下载地址">
                <el-input v-model="versionConfig.iosUrl" placeholder="请输入iOS App Store地址" />
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('version_management')">保存配置</el-button>
              <el-button @click="resetConfig('version_management')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 默认页面 -->
        <div v-else class="page-content">
          <el-empty description="请从左侧选择配置项" />
        </div>
        </template>
      </div>
    </div>

    <!-- 配置历史记录对话框 -->
    <el-dialog 
      v-model="historyDialogVisible" 
      title="配置历史记录" 
      width="800px"
      :close-on-click-modal="false"
    >
      <el-table 
        :data="configHistory" 
        v-loading="loadingHistory"
        style="width: 100%"
      >
        <el-table-column prop="version" label="版本" width="80" />
        <el-table-column prop="operator" label="操作人" width="120" />
        <el-table-column label="配置内容" min-width="200">
          <template #default="{ row }">
            <pre style="margin: 0; font-size: 12px; max-height: 100px; overflow: auto;">{{ formatConfig(row.config) }}</pre>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button 
              type="primary" 
              size="small" 
              @click="rollbackToHistory(row.id)"
            >
              回滚
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <template #footer>
        <el-button @click="historyDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
    
    <!-- 创建数据模型对话框 -->
    <el-dialog v-model="showCreateCollectionDialog" title="新建数据模型" width="700px">
      <el-form :model="collectionForm" label-width="120px">
        <el-form-item label="模型名称" required>
          <el-input v-model="collectionForm.name" placeholder="英文名称，如 products" />
          <div class="form-hint">用于API路径，只能包含小写字母、数字和下划线</div>
        </el-form-item>
        <el-form-item label="显示名称">
          <el-input v-model="collectionForm.display_name" placeholder="中文名称，如 商品" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="collectionForm.description" type="textarea" :rows="2" placeholder="数据模型的用途说明" />
        </el-form-item>
        
        <el-divider content-position="left">字段定义</el-divider>
        
        <div class="field-list">
          <div v-for="(field, index) in collectionForm.fields" :key="index" class="field-item">
            <el-row :gutter="10">
              <el-col :span="6">
                <el-input v-model="field.name" placeholder="字段名" size="small" />
              </el-col>
              <el-col :span="6">
                <el-input v-model="field.display_name" placeholder="显示名" size="small" />
              </el-col>
              <el-col :span="5">
                <el-select v-model="field.type" placeholder="类型" size="small">
                  <el-option label="字符串" value="string" />
                  <el-option label="数字" value="number" />
                  <el-option label="布尔" value="boolean" />
                  <el-option label="数组" value="array" />
                  <el-option label="对象" value="object" />
                </el-select>
              </el-col>
              <el-col :span="4">
                <el-checkbox v-model="field.required" size="small">必填</el-checkbox>
              </el-col>
              <el-col :span="3">
                <el-button type="danger" size="small" @click="removeField(index)" link>
                  <el-icon><Delete /></el-icon>
                </el-button>
              </el-col>
            </el-row>
          </div>
          <el-button type="primary" link @click="addField">
            <el-icon><Plus /></el-icon>添加字段
          </el-button>
        </div>
        
        <el-divider content-position="left">权限设置</el-divider>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="读取权限">
              <el-select v-model="collectionForm.permissions.read">
                <el-option label="公开" value="public" />
                <el-option label="登录用户" value="authenticated" />
                <el-option label="仅创建者" value="owner" />
                <el-option label="仅管理员" value="admin" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="创建权限">
              <el-select v-model="collectionForm.permissions.create">
                <el-option label="公开" value="public" />
                <el-option label="登录用户" value="authenticated" />
                <el-option label="仅管理员" value="admin" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="更新权限">
              <el-select v-model="collectionForm.permissions.update">
                <el-option label="登录用户" value="authenticated" />
                <el-option label="仅创建者" value="owner" />
                <el-option label="仅管理员" value="admin" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="删除权限">
              <el-select v-model="collectionForm.permissions.delete">
                <el-option label="仅创建者" value="owner" />
                <el-option label="仅管理员" value="admin" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="showCreateCollectionDialog = false">取消</el-button>
        <el-button type="primary" @click="createCollection">创建</el-button>
      </template>
    </el-dialog>
    
    <!-- 编辑数据模型对话框 -->
    <el-dialog v-model="showEditCollectionDialog" title="编辑数据模型" width="700px">
      <el-form :model="collectionForm" label-width="120px">
        <el-form-item label="模型名称">
          <el-input v-model="collectionForm.name" disabled />
          <div class="form-hint">模型名称创建后不可修改</div>
        </el-form-item>
        <el-form-item label="显示名称">
          <el-input v-model="collectionForm.display_name" placeholder="中文名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="collectionForm.description" type="textarea" :rows="2" />
        </el-form-item>
        
        <el-divider content-position="left">字段定义</el-divider>
        
        <div class="field-list">
          <div v-for="(field, index) in collectionForm.fields" :key="index" class="field-item">
            <el-row :gutter="10">
              <el-col :span="6">
                <el-input v-model="field.name" placeholder="字段名" size="small" />
              </el-col>
              <el-col :span="6">
                <el-input v-model="field.display_name" placeholder="显示名" size="small" />
              </el-col>
              <el-col :span="5">
                <el-select v-model="field.type" placeholder="类型" size="small">
                  <el-option label="字符串" value="string" />
                  <el-option label="数字" value="number" />
                  <el-option label="布尔" value="boolean" />
                  <el-option label="数组" value="array" />
                  <el-option label="对象" value="object" />
                </el-select>
              </el-col>
              <el-col :span="4">
                <el-checkbox v-model="field.required" size="small">必填</el-checkbox>
              </el-col>
              <el-col :span="3">
                <el-button type="danger" size="small" @click="removeField(index)" link>
                  <el-icon><Delete /></el-icon>
                </el-button>
              </el-col>
            </el-row>
          </div>
          <el-button type="primary" link @click="addField">
            <el-icon><Plus /></el-icon>添加字段
          </el-button>
        </div>
        
        <el-divider content-position="left">权限设置</el-divider>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="读取权限">
              <el-select v-model="collectionForm.permissions.read">
                <el-option label="公开" value="public" />
                <el-option label="登录用户" value="authenticated" />
                <el-option label="仅创建者" value="owner" />
                <el-option label="仅管理员" value="admin" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="创建权限">
              <el-select v-model="collectionForm.permissions.create">
                <el-option label="公开" value="public" />
                <el-option label="登录用户" value="authenticated" />
                <el-option label="仅管理员" value="admin" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="更新权限">
              <el-select v-model="collectionForm.permissions.update">
                <el-option label="登录用户" value="authenticated" />
                <el-option label="仅创建者" value="owner" />
                <el-option label="仅管理员" value="admin" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="删除权限">
              <el-select v-model="collectionForm.permissions.delete">
                <el-option label="仅创建者" value="owner" />
                <el-option label="仅管理员" value="admin" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="showEditCollectionDialog = false">取消</el-button>
        <el-button type="primary" @click="updateCollection">保存</el-button>
      </template>
    </el-dialog>
    
    <!-- API文档对话框 -->
    <el-dialog v-model="showApiDocDialog" title="API文档" width="900px">
      <div v-if="currentCollection" class="api-doc">
        <h3>{{ currentCollection.display_name || currentCollection.name }} API</h3>
        <p class="api-desc">{{ currentCollection.description || '暂无描述' }}</p>
        
        <el-divider content-position="left">基础信息</el-divider>
        <div class="api-info">
          <p><strong>API基础路径:</strong> <code>/api/v1/baas/apps/{{ appId }}/data/{{ currentCollection.name }}</code></p>
          <p><strong>认证方式:</strong> Bearer Token</p>
          <p><strong>请求头:</strong> <code>Authorization: Bearer &lt;your_token&gt;</code></p>
        </div>
        
        <el-divider content-position="left">权限配置</el-divider>
        <div class="api-info">
          <el-table :data="[
            { action: '读取', permission: getPermissionLabel(currentCollection.permissions?.read) },
            { action: '创建', permission: getPermissionLabel(currentCollection.permissions?.create) },
            { action: '更新', permission: getPermissionLabel(currentCollection.permissions?.update) },
            { action: '删除', permission: getPermissionLabel(currentCollection.permissions?.delete) }
          ]" border size="small" style="margin-bottom: 16px">
            <el-table-column prop="action" label="操作" width="100" />
            <el-table-column prop="permission" label="权限要求" />
          </el-table>
        </div>
        
        <el-divider content-position="left">接口列表</el-divider>
        
        <el-collapse>
          <el-collapse-item title="GET - 获取列表" name="list">
            <div class="api-detail">
              <p><strong>请求方式:</strong> GET</p>
              <p><strong>URL:</strong> <code>/api/v1/baas/apps/{{ appId }}/data/{{ currentCollection.name }}</code></p>
              <p><strong>查询参数:</strong></p>
              <ul>
                <li><code>page</code> - 页码，默认1</li>
                <li><code>page_size</code> - 每页数量，默认20</li>
                <li><code>search</code> - 搜索关键词</li>
              </ul>
              <p><strong>响应示例:</strong></p>
              <pre>{
  "code": 0,
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}</pre>
            </div>
          </el-collapse-item>
          
          <el-collapse-item title="POST - 创建文档" name="create">
            <div class="api-detail">
              <p><strong>请求方式:</strong> POST</p>
              <p><strong>URL:</strong> <code>/api/v1/baas/apps/{{ appId }}/data/{{ currentCollection.name }}</code></p>
              <p><strong>请求体:</strong></p>
              <pre>{
  "data": {
    // 根据字段定义填写
  }
}</pre>
            </div>
          </el-collapse-item>
          
          <el-collapse-item title="GET - 获取单个文档" name="get">
            <div class="api-detail">
              <p><strong>请求方式:</strong> GET</p>
              <p><strong>URL:</strong> <code>/api/v1/baas/apps/{{ appId }}/data/{{ currentCollection.name }}/:id</code></p>
            </div>
          </el-collapse-item>
          
          <el-collapse-item title="PUT - 更新文档" name="update">
            <div class="api-detail">
              <p><strong>请求方式:</strong> PUT</p>
              <p><strong>URL:</strong> <code>/api/v1/baas/apps/{{ appId }}/data/{{ currentCollection.name }}/:id</code></p>
              <p><strong>请求体:</strong></p>
              <pre>{
  "data": {
    // 要更新的字段
  }
}</pre>
            </div>
          </el-collapse-item>
          
          <el-collapse-item title="DELETE - 删除文档" name="delete">
            <div class="api-detail">
              <p><strong>请求方式:</strong> DELETE</p>
              <p><strong>URL:</strong> <code>/api/v1/baas/apps/{{ appId }}/data/{{ currentCollection.name }}/:id</code></p>
            </div>
          </el-collapse-item>
        </el-collapse>
        
        <el-divider content-position="left">字段定义</el-divider>
        <el-table :data="currentCollection.fields || []" border size="small">
          <el-table-column prop="name" label="字段名" width="150" />
          <el-table-column prop="display_name" label="显示名" width="150" />
          <el-table-column prop="type" label="类型" width="100" />
          <el-table-column prop="required" label="必填" width="80">
            <template #default="{ row }">
              <el-tag :type="row.required ? 'danger' : 'info'" size="small">
                {{ row.required ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="default_value" label="默认值" />
        </el-table>
      </div>
      <template #footer>
        <el-button @click="showApiDocDialog = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  ArrowLeft, ArrowRight, House, Setting, User, UserFilled, 
  CreditCard, ChatDotRound, DataLine, Document, Monitor, 
  FolderOpened, Tools, Box, Grid, Warning, CopyDocument,
  Bell, DataAnalysis, Promotion, Lock, Plus, Edit, Delete, Search,
  Download, Upload, MagicStick, Cpu, Files, View, Menu
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'
import Workspace from './Workspace.vue'
import MobileMenu from '@/components/MobileMenu.vue'

const route = useRoute()
const router = useRouter()
const appId = computed(() => route.params.id ? String(route.params.id) : '')

// 根据URL参数初始化activeTab
const initialTab = route.query.tab === 'workspace' ? 'workspace' : 'config'
const activeTab = ref(initialTab) // 默认显示配置中心
const mobileMenuOpen = ref(false) // 移动端菜单状态
const currentPage = ref('overview')
const workspaceMenu = ref('overview') // 工作台子菜单

// 工作台菜单配置
const workspaceMenuItems = [
  { key: 'overview', label: '数据概览', icon: House },
  { key: 'users', label: '用户管理', icon: User },
  { key: 'messages', label: '消息推送', icon: Bell },
  { key: 'storage', label: '存储服务', icon: FolderOpened },
  { key: 'events', label: '数据埋点', icon: DataAnalysis },
  { key: 'monitor', label: '监控告警', icon: Monitor },
  { key: 'logs', label: '日志查询', icon: Document },
  { key: 'versions', label: '版本管理', icon: Promotion },
  { key: 'audit', label: '审计日志', icon: Lock }
]
const expandedGroups = ref(['generator', 'baas', 'user', 'message', 'data', 'system', 'storage'])
const adminName = ref(localStorage.getItem('adminName') || 'Admin')

// 模块分组定义
const moduleGroups = [
  { key: 'generator', name: '功能生成器', icon: 'MagicStick', modules: ['api_generator', 'page_generator', 'code_generator'] },
  { key: 'baas', name: 'BaaS数据服务', icon: 'Grid', modules: ['baas_data'] },
  { key: 'user', name: '用户与权限', icon: 'UserFilled', modules: ['user_management'] },
  { key: 'payment', name: '交易与支付', icon: 'CreditCard', modules: ['payment'] },
  { key: 'message', name: '消息与通知', icon: 'ChatDotRound', modules: ['message_center', 'push_service', 'sms_service'] },
  { key: 'data', name: '数据与分析', icon: 'DataLine', modules: ['data_tracking'] },
  { key: 'system', name: '系统与运维', icon: 'Monitor', modules: ['log_service', 'monitor_alert'] },
  { key: 'storage', name: '存储服务', icon: 'FolderOpened', modules: ['file_storage', 'config_management', 'version_management'] }
]

// 模块名称映射
const moduleNameMap = {
  api_generator: 'API生成器',
  page_generator: '页面生成器',
  code_generator: '代码生成器',
  baas_data: '数据模型管理',
  user_management: '用户管理',
  message_center: '消息中心',
  push_service: '推送服务',
  data_tracking: '数据埋点',
  log_service: '日志服务',
  monitor_alert: '监控告警',
  file_storage: '文件存储',
  config_management: '配置管理',
  version_management: '版本管理',
  payment: '支付中心',
  sms_service: '短信服务'
}

const appInfo = ref({
  name: '',
  app_id: '',
  app_secret: '',
  package_name: '',
  description: '',
  status: 1,
  created_at: ''
})

const appModules = ref([])

// BaaS数据模型相关变量
const collections = ref([])
const collectionSearch = ref('')

// 菜单管理相关变量
const menuLoading = ref(false)
const generatedCollections = computed(() => {
  return collections.value.filter(c => c.is_generated)
})
const showCreateCollectionDialog = ref(false)
const showEditCollectionDialog = ref(false)
const showApiDocDialog = ref(false)
const currentCollection = ref(null)
const collectionForm = ref({
  name: '',
  display_name: '',
  description: '',
  fields: [],
  permissions: {
    read: 'public',
    create: 'authenticated',
    update: 'owner',
    delete: 'admin'
  }
})

// 计算属性：过滤后的数据模型列表
const filteredCollections = computed(() => {
  if (!collectionSearch.value) return collections.value
  const search = collectionSearch.value.toLowerCase()
  return collections.value.filter(c => 
    c.name.toLowerCase().includes(search) || 
    (c.display_name && c.display_name.toLowerCase().includes(search)) ||
    (c.description && c.description.toLowerCase().includes(search))
  )
})

// ==================== 功能生成器相关变量 ====================
// API生成器
const apiGeneratorStep = ref(0)
const selectedApiModels = ref([])
const apiGeneratorConfig = ref({
  prefix: '/api/v1',
  endpoints: ['list', 'detail', 'create', 'update', 'delete'],
  pagination: true,
  sorting: true,
  filtering: true,
  authentication: true
})
const generatedApis = ref([])
const generatedGoCode = ref('')

// 页面生成器
const pageGeneratorStep = ref(0)
const selectedPageModel = ref(null)
const pageGeneratorConfig = ref({
  pages: ['list', 'form'],
  listFields: [],
  enableSearch: true,
  enablePagination: true,
  uiFramework: 'element-plus',
  layout: 'table'
})
const pagePreviewTab = ref('list')
const previewTableData = ref([
  { id: 1, name: '示例数据1' },
  { id: 2, name: '示例数据2' }
])
const generatedVueCode = ref('')

// 代码生成器
const codeGeneratorConfig = ref({
  selectedModels: [],
  selectAll: false,
  generateItems: ['backend_api', 'frontend_pages'],
  backendFramework: 'gin',
  database: 'mysql',
  frontendFramework: 'vue3',
  uiLibrary: 'element-plus'
})

// 计算属性：选中模型的字段
const selectedModelFields = computed(() => {
  if (!selectedPageModel.value) return []
  const model = collections.value.find(c => c.id === selectedPageModel.value)
  return model?.fields || []
})

const stats = ref({
  userCount: 0,
  todayRequests: 0,
  todayErrors: 0
})

// 各模块配置表单
const basicConfig = ref({
  name: '',
  description: '',
  package_name: '',
  enableSignature: false,
  ipWhitelist: ''
})

const userConfig = ref({
  passwordMinLength: 8,
  passwordRequirements: ['number', 'letter'],
  enableLoginLock: true,
  maxLoginAttempts: 5,
  lockDuration: 30,
  requiredFields: ['nickname'],
  allowChangeUsername: false,
  enableNicknameFilter: true,
  enableAvatarReview: false,
  enableRealName: false,
  allowAccountDeletion: true,
  deletionCooldown: 7,
  deletionRequirements: ['clearData'],
  deletionConfirmMethod: 'sms'
})

const messageConfig = ref({
  enabled: true,
  retentionDays: 30,
  maxMessagesPerUser: 1000,
  supportedTypes: ['system', 'activity']
})

const pushConfig = ref({
  enabled: true,
  provider: 'jpush',
  appKey: '',
  masterSecret: '',
  enableQuietHours: false,
  quietStart: null,
  quietEnd: null,
  dailyLimit: 10
})

const paymentConfig = ref({
  enableSecurityVerify: true,
  verifyMethods: ['password'],
  verifyThreshold: 500,
  maxPasswordAttempts: 5,
  lockDuration: 30,
  enableLimitControl: true,
  singleLimit: 50000,
  dailyLimit: 100000,
  monthlyLimit: 500000,
  dailyCount: 100,
  successCallback: '',
  failCallback: '',
  refundCallback: '',
  timeout: 30,
  enableAutoRefund: false,
  enablePaymentLog: true
})

const smsConfig = ref({
  enabled: true,
  provider: 'aliyun',
  accessKey: '',
  secretKey: '',
  signName: '',
  codeLength: 6,
  codeExpiry: 5,
  codeTemplateId: '',
  sendInterval: 60,
  dailyLimit: 10,
  enableNotification: true,
  retryCount: 3,
  timeout: 10,
  callbackUrl: '',
  balanceAlert: 1000
})

const trackingConfig = ref({
  enabled: true,
  reportMethod: 'batch',
  batchInterval: 60,
  autoEvents: ['pageView', 'click']
})

const logConfig = ref({
  enabled: true,
  level: 'info',
  storageTypes: ['local'],
  retentionDays: 30,
  realtimeReport: false,
  batchSize: 100
})

const monitorConfig = ref({
  enabled: true,
  metrics: ['api', 'error'],
  interval: 60,
  alertEnabled: true,
  alertMethods: ['email'],
  alertReceivers: ''
})

const storageConfig = ref({
  enabled: true,
  provider: 'aliyun',
  bucket: '',
  accessKey: '',
  secretKey: '',
  maxFileSize: 100,
  allowedTypes: 'jpg,png,gif,pdf,doc,docx'
})

const configMgmtConfig = ref({
  enabled: true,
  refreshInterval: 300,
  enableCache: true
})

const versionConfig = ref({
  enabled: true,
  forceUpdate: false,
  promptType: 'dialog',
  androidUrl: '',
  iosUrl: ''
})

// 切换页面
const switchPage = (page) => {
  currentPage.value = page
  // 切换到模块配置页面时加载配置
  if (page !== 'overview' && page !== 'basic') {
    loadModuleConfig(page)
  }
}

// 移动端菜单关闭处理
const handleMobileMenuClose = () => {
  mobileMenuOpen.value = false
}

// 移动端切换Tab
const switchMobileTab = (tab) => {
  activeTab.value = tab
  // 不再自动关闭菜单，让用户选择子菜单
}

// 移动端切换工作台子菜单
const switchWorkspaceMenu = (menu) => {
  // 确保先切换到工作台Tab
  activeTab.value = 'workspace'
  // 延迟设置菜单，确保Workspace组件已渲染并接收到appId
  setTimeout(() => {
    workspaceMenu.value = menu
  }, 100)
  mobileMenuOpen.value = false
}

// 移动端切换页面
const switchMobilePage = (page) => {
  currentPage.value = page
  mobileMenuOpen.value = false
  if (page !== 'overview' && page !== 'basic') {
    loadModuleConfig(page)
  }
}

// 返回APP列表
const goBackToList = () => {
  mobileMenuOpen.value = false
  router.push('/apps')
}

// 切换分组展开/收起
const toggleGroup = (groupKey) => {
  const index = expandedGroups.value.indexOf(groupKey)
  if (index > -1) {
    expandedGroups.value.splice(index, 1)
  } else {
    expandedGroups.value.push(groupKey)
  }
}

// 检查分组是否有模块
const hasModulesInGroup = (groupKey) => {
  const group = moduleGroups.find(g => g.key === groupKey)
  if (!group) return false
  // 功能生成器和BaaS分组始终显示
  if (groupKey === 'generator' || groupKey === 'baas') return true
  // 使用module_code匹配（后端返回的字段）
  return appModules.value.some(m => group.modules.includes(m.module_code))
}

// 获取分组内的模块
const getModulesInGroup = (groupKey) => {
  const group = moduleGroups.find(g => g.key === groupKey)
  if (!group) return []
  
  // 功能生成器分组返回固定模块
  if (groupKey === 'generator') {
    return [
      { source_module: 'api_generator', name: 'API生成器' },
      { source_module: 'page_generator', name: '页面生成器' },
      { source_module: 'code_generator', name: '代码生成器' }
    ]
  }
  
  // BaaS分组返回固定模块
  if (groupKey === 'baas') {
    return [{
      source_module: 'baas_data',
      name: '数据模型管理'
    }]
  }
  
  // 去重：使用Map确保每个module_code只出现一次
  const uniqueModules = new Map()
  appModules.value
    .filter(m => group.modules.includes(m.module_code))
    .forEach(m => {
      if (!uniqueModules.has(m.module_code)) {
        uniqueModules.set(m.module_code, {
          ...m,
          source_module: m.module_code, // 兼容侧边栏点击
          name: moduleNameMap[m.module_code] || m.module_name || m.name
        })
      }
    })
  
  return Array.from(uniqueModules.values())
}

// ==================== BaaS数据模型相关函数 ====================

// 获取数据模型列表
const fetchCollections = async () => {
  if (!appId.value) return
  try {
    const res = await request.get(`/baas/apps/${appId.value}/collections`)
    collections.value = res.list || res || []
  } catch (error) {
    console.error('获取数据模型列表失败:', error)
    collections.value = []
  }
}

// 创建数据模型
const createCollection = async () => {
  if (!collectionForm.value.name) {
    ElMessage.warning('请输入数据模型名称')
    return
  }
  try {
    await request.post(`/baas/apps/${appId.value}/collections`, {
      name: collectionForm.value.name,
      display_name: collectionForm.value.display_name,
      description: collectionForm.value.description,
      fields: collectionForm.value.fields,
      permissions: collectionForm.value.permissions
    })
    ElMessage.success('创建成功')
    showCreateCollectionDialog.value = false
    resetCollectionForm()
    fetchCollections()
  } catch (error) {
    ElMessage.error('创建失败: ' + (error.message || '未知错误'))
  }
}

// 编辑数据模型
const editCollection = (collection) => {
  currentCollection.value = collection
  collectionForm.value = {
    name: collection.name,
    display_name: collection.display_name || '',
    description: collection.description || '',
    fields: collection.fields || [],
    permissions: collection.permissions || {
      read: 'public',
      create: 'authenticated',
      update: 'owner',
      delete: 'admin'
    }
  }
  showEditCollectionDialog.value = true
}

// 更新数据模型
const updateCollection = async () => {
  if (!currentCollection.value) return
  try {
    await request.put(`/baas/apps/${appId.value}/collections/${currentCollection.value.id}`, {
      display_name: collectionForm.value.display_name,
      description: collectionForm.value.description,
      fields: collectionForm.value.fields,
      permissions: collectionForm.value.permissions
    })
    ElMessage.success('更新成功')
    showEditCollectionDialog.value = false
    resetCollectionForm()
    fetchCollections()
  } catch (error) {
    ElMessage.error('更新失败: ' + (error.message || '未知错误'))
  }
}

// 删除数据模型
const deleteCollection = async (collection) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除数据模型 "${collection.display_name || collection.name}" 吗？\n此操作将删除该模型下的所有数据，且不可恢复。`,
      '删除确认',
      { confirmButtonText: '删除', cancelButtonText: '取消', type: 'warning' }
    )
    await request.delete(`/baas/apps/${appId.value}/collections/${collection.id}`)
    ElMessage.success('删除成功')
    fetchCollections()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败: ' + (error.message || '未知错误'))
    }
  }
}

// 获取权限标签
const getPermissionLabel = (permission) => {
  const labels = {
    'public': '公开访问',
    'auth': '登录用户',
    'owner': '仅创建者',
    'admin': '仅管理员'
  }
  return labels[permission] || permission || '未设置'
}

// 查看API文档
const viewApiDoc = (collection) => {
  currentCollection.value = collection
  showApiDocDialog.value = true
}

// 复制API端点
const copyApiEndpoint = (collection) => {
  const endpoint = `/api/v1/baas/apps/${appId.value}/data/${collection.name}`
  navigator.clipboard.writeText(endpoint)
  ElMessage.success('API端点已复制到剪贴板')
}

// 生成功能
const generateFeature = async (collection) => {
  try {
    await ElMessageBox.confirm(
      `确定要为数据模型 "${collection.display_name || collection.name}" 生成功能吗？\n\n生成后将在工作台中创建对应的数据管理页面。`,
      '生成功能',
      {
        confirmButtonText: '确定生成',
        cancelButtonText: '取消',
        type: 'info'
      }
    )
    
    // 调用后端API更新is_generated状态
    await request.post(`/baas/apps/${appId.value}/collections/${collection.id}/generate`)
    
    ElMessage.success('功能生成成功！可以在工作台中使用了')
    // 刷新列表
    fetchCollections()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('生成功能失败:', error)
      ElMessage.error('生成功能失败')
    }
  }
}

// 跳转到工作台查看功能
const goToWorkspace = (collection) => {
  activeTab.value = 'workspace'
  // 通过事件或状态通知Workspace组件切换到对应的数据模型
  setTimeout(() => {
    // 触发自定义事件通知Workspace组件
    window.dispatchEvent(new CustomEvent('switch-to-collection', { 
      detail: { collectionId: collection.id, collectionName: collection.name } 
    }))
  }, 100)
}

// 切换菜单显示状态
const toggleMenuVisibility = async (collection) => {
  collection.loading = true
  try {
    await request.put(`/baas/apps/${appId.value}/collections/${collection.id}/visibility`, {
      is_visible: collection.is_visible
    })
    ElMessage.success(collection.is_visible ? '已开启工作台显示' : '已关闭工作台显示')
  } catch (error) {
    console.error('更新显示状态失败:', error)
    // 回滚状态
    collection.is_visible = !collection.is_visible
    ElMessage.error('更新失败，请重试')
  } finally {
    collection.loading = false
  }
}

// 重置表单
const resetCollectionForm = () => {
  collectionForm.value = {
    name: '',
    display_name: '',
    description: '',
    fields: [],
    permissions: {
      read: 'public',
      create: 'authenticated',
      update: 'owner',
      delete: 'admin'
    }
  }
  currentCollection.value = null
}

// 添加字段
const addField = () => {
  collectionForm.value.fields.push({
    name: '',
    display_name: '',
    type: 'string',
    required: false,
    default_value: ''
  })
}

// 删除字段
const removeField = (index) => {
  collectionForm.value.fields.splice(index, 1)
}

// ==================== 结束BaaS相关函数 ====================

// ==================== 功能生成器相关函数 ====================

// 切换选中的API模型
const toggleApiModel = (modelId) => {
  const index = selectedApiModels.value.indexOf(modelId)
  if (index === -1) {
    selectedApiModels.value.push(modelId)
  } else {
    selectedApiModels.value.splice(index, 1)
  }
}

// 生成API预览
const generateApiPreview = () => {
  const apis = []
  const selectedModels = collections.value.filter(c => selectedApiModels.value.includes(c.id))
  
  selectedModels.forEach(model => {
    const basePath = `${apiGeneratorConfig.value.prefix}/${model.name}`
    
    if (apiGeneratorConfig.value.endpoints.includes('list')) {
      apis.push({ method: 'GET', path: basePath, description: `获取${model.display_name || model.name}列表` })
    }
    if (apiGeneratorConfig.value.endpoints.includes('detail')) {
      apis.push({ method: 'GET', path: `${basePath}/:id`, description: `获取${model.display_name || model.name}详情` })
    }
    if (apiGeneratorConfig.value.endpoints.includes('create')) {
      apis.push({ method: 'POST', path: basePath, description: `创建${model.display_name || model.name}` })
    }
    if (apiGeneratorConfig.value.endpoints.includes('update')) {
      apis.push({ method: 'PUT', path: `${basePath}/:id`, description: `更新${model.display_name || model.name}` })
    }
    if (apiGeneratorConfig.value.endpoints.includes('delete')) {
      apis.push({ method: 'DELETE', path: `${basePath}/:id`, description: `删除${model.display_name || model.name}` })
    }
    if (apiGeneratorConfig.value.endpoints.includes('batch_delete')) {
      apis.push({ method: 'POST', path: `${basePath}/batch-delete`, description: `批量删除${model.display_name || model.name}` })
    }
  })
  
  generatedApis.value = apis
  
  // 生成Go代码
  generatedGoCode.value = generateGoCode(selectedModels)
  
  apiGeneratorStep.value = 2
}

// 生成Go代码
const generateGoCode = (models) => {
  let code = `package api

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

`
  
  models.forEach(model => {
    const modelName = model.name.charAt(0).toUpperCase() + model.name.slice(1)
    code += `// ${modelName} API Handlers

func Get${modelName}List(c *gin.Context) {
    // TODO: Implement list query
    c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
}

func Get${modelName}Detail(c *gin.Context) {
    id := c.Param("id")
    // TODO: Implement detail query
    c.JSON(http.StatusOK, gin.H{"id": id})
}

func Create${modelName}(c *gin.Context) {
    // TODO: Implement create
    c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func Update${modelName}(c *gin.Context) {
    id := c.Param("id")
    // TODO: Implement update
    c.JSON(http.StatusOK, gin.H{"id": id, "message": "updated"})
}

func Delete${modelName}(c *gin.Context) {
    id := c.Param("id")
    // TODO: Implement delete
    c.JSON(http.StatusOK, gin.H{"id": id, "message": "deleted"})
}

`
  })
  
  // 添加路由注册
code += `// RegisterRoutes registers all API routes
func RegisterRoutes(r *gin.Engine) {
`
  models.forEach(model => {
    const modelName = model.name.charAt(0).toUpperCase() + model.name.slice(1)
    code += `    // ${model.display_name || model.name} routes
    r.GET("${apiGeneratorConfig.value.prefix}/${model.name}", Get${modelName}List)
    r.GET("${apiGeneratorConfig.value.prefix}/${model.name}/:id", Get${modelName}Detail)
    r.POST("${apiGeneratorConfig.value.prefix}/${model.name}", Create${modelName})
    r.PUT("${apiGeneratorConfig.value.prefix}/${model.name}/:id", Update${modelName})
    r.DELETE("${apiGeneratorConfig.value.prefix}/${model.name}/:id", Delete${modelName})
`
  })
  code += `}
`
  
  return code
}

// 获取HTTP方法的标签类型
const getMethodTagType = (method) => {
  const types = {
    'GET': 'success',
    'POST': 'primary',
    'PUT': 'warning',
    'DELETE': 'danger'
  }
  return types[method] || 'info'
}

// 复制所有API代码
const copyAllApiCode = () => {
  navigator.clipboard.writeText(generatedGoCode.value)
  ElMessage.success('API代码已复制到剪贴板')
}

// 复制指定语言的代码
const copyApiCode = (lang) => {
  if (lang === 'go') {
    navigator.clipboard.writeText(generatedGoCode.value)
    ElMessage.success('Go代码已复制到剪贴板')
  }
}

// 下载API代码
const downloadApiCode = () => {
  const blob = new Blob([generatedGoCode.value], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'api_handlers.go'
  a.click()
  URL.revokeObjectURL(url)
  ElMessage.success('代码文件已下载')
}

// 部署API
const deployApi = () => {
  ElMessage.info('部署功能开发中...')
}

// 生成页面预览
const generatePagePreview = () => {
  const model = collections.value.find(c => c.id === selectedPageModel.value)
  if (!model) return
  
  // 如果没有选择字段，默认选择所有字段
  if (pageGeneratorConfig.value.listFields.length === 0) {
    pageGeneratorConfig.value.listFields = (model.fields || []).map(f => f.name)
  }
  
  // 生成Vue代码
  generatedVueCode.value = generateVueCode(model)
  
  pageGeneratorStep.value = 2
}

// 生成Vue代码
const generateVueCode = (model) => {
  const modelName = model.name.charAt(0).toUpperCase() + model.name.slice(1)
  const fields = model.fields || []
  
  let code = `<template>
  <div class="${model.name}-list">
    <div class="page-header">
      <h2>${model.display_name || model.name}管理</h2>
      <el-button type="primary" @click="showAddDialog = true">新增</el-button>
    </div>
`
  
  if (pageGeneratorConfig.value.enableSearch) {
    code += `    
    <div class="search-bar">
      <el-input v-model="searchKeyword" placeholder="搜索..." style="width: 300px;" clearable />
      <el-button type="primary" @click="handleSearch">搜索</el-button>
    </div>
`
  }
  
  code += `
    <el-table :data="tableData" style="width: 100%">
`
  
  pageGeneratorConfig.value.listFields.forEach(fieldName => {
    const field = fields.find(f => f.name === fieldName)
    if (field) {
      code += `      <el-table-column prop="${field.name}" label="${field.display_name || field.name}" />
`
    }
  })
  
  code += `      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
`
  
  if (pageGeneratorConfig.value.enablePagination) {
    code += `
    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="fetchData"
    />
`
  }
  
  code += `  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const tableData = ref([])
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const showAddDialog = ref(false)

const fetchData = async () => {
  // TODO: 调用API获取数据
}

const handleSearch = () => {
  currentPage.value = 1
  fetchData()
}

const handleEdit = (row) => {
  // TODO: 编辑逻辑
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定要删除吗？', '提示')
  // TODO: 调用删除API
  ElMessage.success('删除成功')
  fetchData()
}

onMounted(() => {
  fetchData()
})
<\/script>
`
  
  return code
}

// 获取字段标签
const getFieldLabel = (fieldName) => {
  const model = collections.value.find(c => c.id === selectedPageModel.value)
  if (!model) return fieldName
  const field = (model.fields || []).find(f => f.name === fieldName)
  return field?.display_name || fieldName
}

// 复制页面代码
const copyPageCode = () => {
  navigator.clipboard.writeText(generatedVueCode.value)
  ElMessage.success('Vue代码已复制到剪贴板')
}

// 下载页面代码
const downloadPageCode = () => {
  const model = collections.value.find(c => c.id === selectedPageModel.value)
  const fileName = model ? `${model.name}List.vue` : 'page.vue'
  
  const blob = new Blob([generatedVueCode.value], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = fileName
  a.click()
  URL.revokeObjectURL(url)
  ElMessage.success('页面代码已下载')
}

// 切换全选模型
const toggleSelectAllModels = (val) => {
  if (val) {
    codeGeneratorConfig.value.selectedModels = collections.value.map(c => c.id)
  } else {
    codeGeneratorConfig.value.selectedModels = []
  }
}

// 生成完整代码
const generateFullCode = async () => {
  const selectedModels = collections.value.filter(c => 
    codeGeneratorConfig.value.selectedModels.includes(c.id)
  )
  
  if (selectedModels.length === 0) {
    ElMessage.warning('请选择至少一个数据模型')
    return
  }
  
  // 创建ZIP文件内容
  let zipContent = ''
  
  // 后端API代码
  if (codeGeneratorConfig.value.generateItems.includes('backend_api')) {
    const goCode = generateGoCode(selectedModels)
    zipContent += '=== backend/api/handlers.go ===\n' + goCode + '\n\n'
  }
  
  // 前端页面代码
  if (codeGeneratorConfig.value.generateItems.includes('frontend_pages')) {
    selectedModels.forEach(model => {
      const vueCode = generateVueCode(model)
      zipContent += `=== frontend/pages/${model.name}List.vue ===\n` + vueCode + '\n\n'
    })
  }
  
  // 数据库SQL
  if (codeGeneratorConfig.value.generateItems.includes('database_sql')) {
    let sqlCode = '-- Database Schema\n\n'
    selectedModels.forEach(model => {
      sqlCode += `CREATE TABLE IF NOT EXISTS ${model.name} (\n`
      sqlCode += '  id BIGINT PRIMARY KEY AUTO_INCREMENT,\n'
      ;(model.fields || []).forEach(field => {
        let sqlType = 'VARCHAR(255)'
        if (field.type === 'number') sqlType = 'DECIMAL(10,2)'
        if (field.type === 'boolean') sqlType = 'TINYINT(1)'
        if (field.type === 'object' || field.type === 'array') sqlType = 'JSON'
        sqlCode += `  ${field.name} ${sqlType}${field.required ? ' NOT NULL' : ''},\n`
      })
      sqlCode += '  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,\n'
      sqlCode += '  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP\n'
      sqlCode += ');\n\n'
    })
    zipContent += '=== database/schema.sql ===\n' + sqlCode + '\n\n'
  }
  
  // 下载文件
  const blob = new Blob([zipContent], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `generated_code_${Date.now()}.txt`
  a.click()
  URL.revokeObjectURL(url)
  ElMessage.success('代码文件已生成并下载')
}

// ==================== 结束功能生成器相关函数 ====================

// 获取APP信息
const fetchAppInfo = async () => {
  if (!appId.value || appId.value === '') return
  try {
    const res = await request.get(`/apps/${appId.value}`)
    // request.js已解包，res直接是数据对象
    if (res) {
      appInfo.value = res
      basicConfig.value.name = res.name
      basicConfig.value.description = res.description || ''
      basicConfig.value.package_name = res.package_name || ''
    }
  } catch (error) {
    console.error('获取APP信息失败:', error)
  }
}

// 获取APP模块列表
const fetchAppModules = async () => {
  if (!appId.value || appId.value === '') return
  try {
    const res = await request.get(`/apps/${appId.value}/modules`)
    // request.js已解包，res直接是数据数组
    if (res) {
      appModules.value = res
    }
  } catch (error) {
    console.error('获取APP模块失败:', error)
  }
}

// 复制文本
const copyText = (text) => {
  if (!text) return
  navigator.clipboard.writeText(text)
  ElMessage.success('已复制到剪贴板')
}

// 遮盖密钥
const maskSecret = (secret) => {
  if (!secret) return '-'
  if (secret.length <= 8) return '********'
  return secret.substring(0, 4) + '****' + secret.substring(secret.length - 4)
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

// 保存基础配置
const saveBasicConfig = async () => {
  try {
    await request.put(`/apps/${appId.value}`, basicConfig.value)
    ElMessage({
      message: '基础配置保存成功',
      type: 'success',
      duration: 3000,
      showClose: true
    })
    fetchAppInfo()
  } catch (error) {
    ElMessage({
      message: '配置保存失败，请稍后重试',
      type: 'error',
      duration: 5000,
      showClose: true
    })
  }
}

// 重置基础配置
const resetBasicConfig = () => {
  basicConfig.value.name = appInfo.value.name
  basicConfig.value.description = appInfo.value.description || ''
  basicConfig.value.package_name = appInfo.value.package_name || ''
}

// 获取模块配置数据
const getModuleConfigData = (moduleKey) => {
  const configMap = {
    'user_management': userConfig.value,
    'message_center': messageConfig.value,
    'push_service': pushConfig.value,
    'payment': paymentConfig.value,
    'sms_service': smsConfig.value,
    'data_tracking': trackingConfig.value,
    'log_service': logConfig.value,
    'monitor_alert': monitorConfig.value,
    'file_storage': fileConfig.value,
    'config_management': configMgmtConfig.value,
    'version_management': versionConfig.value
  }
  return configMap[moduleKey] || {}
}

// 加载模块配置
const loadModuleConfig = async (moduleKey) => {
  if (!appId.value || appId.value === '') return
  try {
    const res = await request.get(`/apps/${appId.value}/modules/${moduleKey}/config`)
    // request.js已解包，res直接是数据对象
    if (res && res.config) {
      const config = typeof res.config === 'string' ? JSON.parse(res.config) : res.config
      // 更新对应的配置对象
      const configMap = {
        'user_management': userConfig,
        'message_center': messageConfig,
        'push_service': pushConfig,
        'payment': paymentConfig,
        'sms_service': smsConfig,
        'data_tracking': trackingConfig,
        'log_service': logConfig,
        'monitor_alert': monitorConfig,
        'file_storage': fileConfig,
        'config_management': configMgmtConfig,
        'version_management': versionConfig
      }
      if (configMap[moduleKey]) {
        Object.assign(configMap[moduleKey].value, config)
      }
    }
  } catch (error) {
    console.error('加载配置失败:', error)
  }
}

// 保存模块配置
const saveModuleConfig = async (moduleKey) => {
  try {
    const configData = getModuleConfigData(moduleKey)
    await request.put(`/apps/${appId.value}/modules/${moduleKey}/config`, {
      config: configData
    })
    ElMessage({
      message: '配置保存成功',
      type: 'success',
      duration: 3000,
      showClose: true
    })
  } catch (error) {
    console.error('保存配置失败:', error)
    ElMessage({
      message: '配置保存失败，请稍后重试',
      type: 'error',
      duration: 5000,
      showClose: true
    })
  }
}

// 测试配置
const testConfig = (moduleKey) => {
  ElMessage.info('测试功能开发中...')
}

// 重置配置
const resetConfig = async (moduleKey) => {
  try {
    await ElMessageBox.confirm('确定要重置此模块的配置吗？重置后将恢复为默认配置。', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await request.post(`/apps/${appId.value}/modules/${moduleKey}/config/reset`)
    ElMessage.success('配置已重置')
    // 重新加载配置
    await loadModuleConfig(moduleKey)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重置配置失败:', error)
      ElMessage.error('重置失败')
    }
  }
}

// 配置历史记录
const historyDialogVisible = ref(false)
const currentHistoryModule = ref('')
const configHistory = ref([])
const loadingHistory = ref(false)

// 显示配置历史
const showConfigHistory = async (moduleKey) => {
  currentHistoryModule.value = moduleKey
  historyDialogVisible.value = true
  await loadConfigHistory(moduleKey)
}

// 加载配置历史
const loadConfigHistory = async (moduleKey) => {
  if (!appId.value || appId.value === '') return
  loadingHistory.value = true
  try {
    const res = await request.get(`/apps/${appId.value}/modules/${moduleKey}/config/history`)
    // request.js已解包，res直接是数据数组
    if (res) {
      configHistory.value = res
    }
  } catch (error) {
    console.error('加载历史记录失败:', error)
    ElMessage.error('加载失败')
  } finally {
    loadingHistory.value = false
  }
}

// 回滚配置
const rollbackToHistory = async (historyId) => {
  try {
    await ElMessageBox.confirm('确定要回滚到该版本的配置吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await request.post(`/apps/${appId.value}/modules/${currentHistoryModule.value}/config/rollback/${historyId}`)
    ElMessage.success('配置已回滚')
    historyDialogVisible.value = false
    // 重新加载配置
    await loadModuleConfig(currentHistoryModule.value)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('回滚配置失败:', error)
      ElMessage.error('回滚失败')
    }
  }
}

// 格式化配置显示
const formatConfig = (config) => {
  try {
    const configObj = typeof config === 'string' ? JSON.parse(config) : config
    return JSON.stringify(configObj, null, 2)
  } catch {
    return config
  }
}

onMounted(() => {
  // 只有当appId有效时才加载数据
  if (appId.value && appId.value !== '') {
    fetchAppInfo()
    fetchAppModules()
    fetchCollections()
  }
})

// 监听appId变化，当appId从空变为有效时加载数据
watch(appId, (newVal, oldVal) => {
  if (newVal && newVal !== '' && (!oldVal || oldVal === '')) {
    fetchAppInfo()
    fetchAppModules()
    fetchCollections()
  }
})

// 监听currentPage变化，当切换到baas_data时加载数据
watch(() => currentPage.value, (newVal) => {
  if (newVal === 'baas_data') {
    fetchCollections()
  }
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
  padding: 0 24px;
  height: 60px;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 24px;
}

.header-right {
  flex: 1;
}

.header-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  
  .back-btn {
    color: white;
    &:hover {
      background: rgba(255, 255, 255, 0.1);
    }
  }
  
  .app-icon {
    width: 36px;
    height: 36px;
    background: rgba(255, 255, 255, 0.2);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    font-size: 16px;
  }
  
  .app-name {
    font-size: 18px;
    font-weight: 600;
  }
}

.header-nav {
  display: flex;
  align-items: center;
  gap: 4px;
  
  .nav-item {
    padding: 8px 20px;
    font-size: 14px;
    font-weight: 500;
    color: rgba(255, 255, 255, 0.7);
    cursor: pointer;
    border-radius: 6px;
    transition: all 0.2s;
    
    &:hover {
      color: white;
      background: rgba(255, 255, 255, 0.1);
    }
    
    &.active {
      color: white;
      background: #409eff;
      font-weight: 600;
    }
  }
}

.main-container {
  display: flex;
  height: calc(100vh - 60px);
}

.sidebar {
  width: 240px;
  background: white;
  border-right: 1px solid #e4e7ed;
  overflow-y: auto;
  padding: 16px 0;
  display: flex;
  flex-direction: column;
}

.sidebar-menu {
  flex: 1;
}

.sidebar-footer {
  border-top: 1px solid #e4e7ed;
  padding-top: 8px;
  margin-top: 8px;
  
  .back-item {
    color: #909399;
    
    &:hover {
      color: #409eff;
      background: #f5f7fa;
    }
  }
}

.sidebar-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  cursor: pointer;
  color: #606266;
  transition: all 0.2s;
  
  &:hover {
    background: #f5f7fa;
    color: #409eff;
  }
  
  &.active {
    background: #ecf5ff;
    color: #409eff;
    border-right: 3px solid #409eff;
  }
  
  &.sub-item {
    padding-left: 48px;
    font-size: 14px;
  }
}

.sidebar-group {
  .group-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 20px;
    cursor: pointer;
    color: #303133;
    font-weight: 500;
    
    &:hover {
      background: #f5f7fa;
    }
    
    .group-title {
      display: flex;
      align-items: center;
      gap: 10px;
    }
    
    .expand-icon {
      transition: transform 0.2s;
      
      &.expanded {
        transform: rotate(90deg);
      }
    }
  }
  
  .group-items {
    background: #fafafa;
  }
}

.content-area {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.page-content {
  max-width: 900px;
}

.page-title {
  font-size: 22px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.page-desc {
  color: #909399;
  margin-bottom: 24px;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  
  .stat-icon {
    width: 48px;
    height: 48px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
    
    &.users { background: #e6f7ff; color: #1890ff; }
    &.modules { background: #f6ffed; color: #52c41a; }
    &.requests { background: #fff7e6; color: #fa8c16; }
    &.errors { background: #fff1f0; color: #f5222d; }
  }
  
  .stat-info {
    .stat-value {
      font-size: 28px;
      font-weight: 600;
      color: #303133;
    }
    
    .stat-label {
      font-size: 14px;
      color: #909399;
    }
  }
}

.info-section {
  background: white;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  
  h3 {
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 1px solid #ebeef5;
  }
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  
  label {
    font-size: 13px;
    color: #909399;
  }
  
  span {
    font-size: 14px;
    color: #303133;
    
    &.copyable {
      cursor: pointer;
      display: flex;
      align-items: center;
      gap: 4px;
      
      &:hover {
        color: #409eff;
      }
    }
  }
}

.module-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.config-form {
  background: white;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.form-section {
  margin-bottom: 32px;
  
  h4 {
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 20px;
    padding-bottom: 12px;
    border-bottom: 1px solid #ebeef5;
  }
}

.form-hint {
  margin-left: 12px;
  font-size: 12px;
  color: #909399;
}

:deep(.el-form-item) {
  margin-bottom: 20px;
}

:deep(.el-input-number) {
  width: 150px;
}

/* 移动端菜单样式 */
.mobile-nav-tabs {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 0 12px;
  margin-bottom: 16px;
  border-bottom: 1px solid #e4e7ed;
  padding-bottom: 16px;
}

.mobile-nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 8px;
  cursor: pointer;
  color: #606266;
  font-size: 15px;
  transition: all 0.2s;

  &:hover {
    background: #f5f7fa;
    color: #409eff;
  }

  &.active {
    background: #ecf5ff;
    color: #409eff;
    font-weight: 600;
  }

  .el-icon {
    font-size: 18px;
  }
}

.mobile-sidebar-menu {
  padding: 0 12px;
}

.mobile-menu-group {
  margin-bottom: 8px;
}

.mobile-group-title {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  color: #909399;
  font-size: 13px;
  font-weight: 500;
  text-transform: uppercase;
}

.mobile-menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 8px;
  cursor: pointer;
  color: #606266;
  font-size: 15px;
  transition: all 0.2s;

  &:hover {
    background: #f5f7fa;
    color: #409eff;
  }

  &.active {
    background: #ecf5ff;
    color: #409eff;
    font-weight: 600;
  }

  &.sub-item {
    padding-left: 48px;
    font-size: 14px;
  }

  &.back-item {
    color: #909399;
    
    &:hover {
      color: #409eff;
      background: #f5f7fa;
    }
  }

  .el-icon {
    font-size: 18px;
  }
}

/* 移动端响应式样式 */
@media (max-width: 768px) {
  .header-logo,
  .header-nav {
    display: none;
  }

  .sidebar {
    display: none !important;
  }

  .content-area {
    width: 100% !important;
    padding: 0 !important;
  }

  .page-content {
    padding: 16px !important;
  }

  .page-title {
    font-size: 20px !important;
  }

  .page-desc {
    font-size: 13px !important;
  }

  .stats-cards {
    grid-template-columns: repeat(2, 1fr) !important;
    gap: 12px !important;
  }

  .info-grid {
    grid-template-columns: 1fr !important;
  }

  .config-form {
    padding: 12px !important;
  }

  /* 表单样式优化 */
  :deep(.el-form) {
    --el-form-label-font-size: 13px;
  }

  :deep(.el-form-item) {
    flex-direction: column !important;
    align-items: flex-start !important;
    margin-bottom: 16px !important;
  }

  :deep(.el-form-item__label) {
    width: 100% !important;
    text-align: left !important;
    margin-bottom: 8px !important;
    padding-right: 0 !important;
    line-height: 1.4 !important;
    white-space: normal !important;
  }

  :deep(.el-form-item__content) {
    width: 100% !important;
    margin-left: 0 !important;
    flex-wrap: wrap !important;
  }

  /* 输入框全宽 */
  :deep(.el-input),
  :deep(.el-select),
  :deep(.el-input-number) {
    width: 100% !important;
  }

  :deep(.el-input-number) {
    max-width: 150px !important;
  }

  /* 复选框组换行 */
  :deep(.el-checkbox-group) {
    display: flex !important;
    flex-direction: column !important;
    gap: 8px !important;
  }

  /* 提示文字换行 */
  .form-hint {
    display: block !important;
    margin-top: 4px !important;
    margin-left: 0 !important;
  }

  /* 表单分区标题 */
  .form-section h4 {
    font-size: 15px !important;
  }

  /* 按钮组 */
  :deep(.el-form-item:last-child .el-form-item__content) {
    flex-wrap: wrap !important;
    gap: 8px !important;
  }

  :deep(.el-button) {
    margin-left: 0 !important;
  }

  /* 时间选择器 */
  :deep(.el-time-picker) {
    width: 100% !important;
    margin-bottom: 8px !important;
  }
}

@media (max-width: 480px) {
  .stats-cards {
    grid-template-columns: 1fr !important;
  }

  .page-title {
    font-size: 18px !important;
  }

  .page-desc {
    font-size: 12px !important;
  }

  /* 更小屏幕的表单优化 */
  .config-form {
    padding: 8px !important;
  }

  :deep(.el-form-item__label) {
    font-size: 13px !important;
  }
}

// BaaS数据模型管理样式
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.search-bar {
  margin-bottom: 20px;
}

.collection-list {
  .empty-state {
    padding: 60px 0;
    text-align: center;
  }
}

.collection-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
  gap: 20px;
}

.collection-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: box-shadow 0.3s;
  
  &:hover {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  }
  
  .card-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;
    
    .card-icon {
      width: 48px;
      height: 48px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
    }
    
    .card-info {
      flex: 1;
      
      h3 {
        margin: 0;
        font-size: 16px;
        font-weight: 600;
        color: #303133;
      }
      
      .card-name {
        font-size: 12px;
        color: #909399;
      }
    }
  }
  
  .card-desc {
    font-size: 14px;
    color: #606266;
    margin: 0 0 12px;
    line-height: 1.5;
  }
  
  .card-stats {
    display: flex;
    gap: 20px;
    margin-bottom: 16px;
    
    span {
      font-size: 13px;
      color: #909399;
    }
  }
  
  .card-actions {
    display: flex;
    gap: 8px;
    margin-bottom: 12px;
  }
  
  .card-api {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px;
    background: #f5f7fa;
    border-radius: 6px;
    
    .api-label {
      font-size: 12px;
      color: #909399;
    }
    
    code {
      flex: 1;
      font-size: 12px;
      color: #606266;
      word-break: break-all;
    }
  }
}

.field-list {
  .field-item {
    margin-bottom: 12px;
    padding: 12px;
    background: #f5f7fa;
    border-radius: 6px;
  }
}

.form-hint {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.api-doc {
  h3 {
    margin: 0 0 8px;
    font-size: 18px;
  }
  
  .api-desc {
    color: #606266;
    margin-bottom: 16px;
  }
  
  .api-info {
    p {
      margin: 8px 0;
    }
    
    code {
      background: #f5f7fa;
      padding: 4px 8px;
      border-radius: 4px;
      font-size: 13px;
    }
  }
  
  .api-detail {
    padding: 16px;
    
    p {
      margin: 8px 0;
    }
    
    ul {
      margin: 8px 0;
      padding-left: 20px;
    }
    
    pre {
      background: #f5f7fa;
      padding: 12px;
      border-radius: 6px;
      overflow-x: auto;
      font-size: 13px;
    }
    
    code {
      background: #f5f7fa;
      padding: 2px 6px;
      border-radius: 4px;
    }
  }
}

/* 功能生成器样式 */
.generator-steps {
  margin: 20px 0 30px;
  padding: 0 40px;
}

.generator-content {
  padding: 20px;
  
  h3 {
    margin-bottom: 20px;
    font-size: 16px;
    color: #303133;
  }
}

.model-select-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.model-select-card {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  
  &:hover {
    border-color: #409eff;
    box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
  }
  
  &.selected {
    border-color: #409eff;
    background: #ecf5ff;
  }
  
  .model-info {
    flex: 1;
    
    h4 {
      margin: 0 0 4px;
      font-size: 14px;
      color: #303133;
    }
    
    .model-name {
      font-size: 12px;
      color: #909399;
      font-family: monospace;
    }
    
    p {
      margin: 8px 0 0;
      font-size: 12px;
      color: #606266;
    }
    
    .model-stats {
      margin-top: 8px;
      font-size: 12px;
      color: #909399;
    }
  }
}

.step-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #ebeef5;
}

.preview-section {
  margin-bottom: 24px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
  
  .preview-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    background: #f5f7fa;
    border-bottom: 1px solid #e4e7ed;
    
    span {
      font-weight: 500;
      color: #303133;
    }
  }
}

.api-preview-list {
  padding: 16px;
  
  .api-preview-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 8px 0;
    border-bottom: 1px dashed #ebeef5;
    
    &:last-child {
      border-bottom: none;
    }
    
    code {
      font-family: monospace;
      color: #409eff;
    }
    
    .api-desc {
      color: #909399;
      font-size: 13px;
    }
  }
}

.code-preview {
  margin: 0;
  padding: 16px;
  background: #1e1e1e;
  color: #d4d4d4;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 13px;
  line-height: 1.5;
  overflow-x: auto;
  max-height: 400px;
  
  code {
    background: transparent;
    padding: 0;
    color: inherit;
  }
}

.preview-tabs {
  margin-bottom: 24px;
}

.page-preview-frame {
  padding: 20px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: #fff;
  
  .preview-toolbar {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;
  }
}

.code-gen-options {
  h3 {
    margin-bottom: 24px;
  }
  
  .option-section {
    margin-bottom: 24px;
    padding: 20px;
    background: #f5f7fa;
    border-radius: 8px;
    
    h4 {
      margin: 0 0 16px;
      font-size: 14px;
      color: #303133;
    }
  }
  
  .model-checkbox-list {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    margin-bottom: 12px;
  }
}

.code-gen-actions {
  text-align: center;
  padding: 30px 0;
  
  .hint {
    margin-top: 12px;
    color: #909399;
    font-size: 13px;
  }
}

.empty-state {
  padding: 40px 0;
  text-align: center;
}

/* 菜单管理页面样式 */
.menu-list-card {
  margin-bottom: 20px;
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    
    .hint-text {
      font-size: 12px;
      color: #909399;
    }
  }
  
  .menu-name {
    display: flex;
    align-items: center;
    gap: 8px;
    
    .el-icon {
      color: #409eff;
    }
  }
  
  .text-gray {
    color: #909399;
  }
}

.menu-tips {
  margin-top: 20px;
  
  .tips-list {
    margin: 0;
    padding-left: 20px;
    
    li {
      line-height: 2;
      color: #606266;
    }
  }
}
</style>
