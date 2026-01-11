# 版本管理功能设计文档

## 1. 功能概述

为低代码平台的生成功能（数据模型）和开发模块（配置）增加版本管理能力，实现：
- 版本追踪：记录每次变更
- 版本对比：查看不同版本差异
- 版本发布：控制版本上线
- 版本回滚：快速恢复到历史版本

## 2. 数据表设计

### 2.1 功能版本表 (feature_versions)

记录数据模型（Collection）的版本历史。

```sql
CREATE TABLE feature_versions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    app_id BIGINT NOT NULL COMMENT '应用ID',
    collection_id BIGINT NOT NULL COMMENT '数据模型ID',
    version VARCHAR(20) NOT NULL COMMENT '版本号 (如 1.0.0)',
    version_num INT NOT NULL COMMENT '版本序号 (自增)',
    schema_snapshot JSON NOT NULL COMMENT '字段结构快照',
    status ENUM('draft', 'published', 'deprecated') DEFAULT 'draft' COMMENT '版本状态',
    changelog TEXT COMMENT '变更说明',
    created_by VARCHAR(100) COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP NULL COMMENT '发布时间',
    INDEX idx_collection_id (collection_id),
    INDEX idx_app_id (app_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 2.2 模块版本表 (module_versions)

记录模块配置的版本历史。

```sql
CREATE TABLE module_versions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    app_id BIGINT NOT NULL COMMENT '应用ID',
    module_code VARCHAR(50) NOT NULL COMMENT '模块代码',
    version VARCHAR(20) NOT NULL COMMENT '版本号',
    version_num INT NOT NULL COMMENT '版本序号',
    config_snapshot JSON NOT NULL COMMENT '配置快照',
    status ENUM('draft', 'published', 'deprecated') DEFAULT 'draft',
    environment ENUM('dev', 'test', 'prod') DEFAULT 'dev' COMMENT '环境',
    changelog TEXT COMMENT '变更说明',
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP NULL,
    INDEX idx_app_module (app_id, module_code),
    INDEX idx_environment (environment),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 3. API设计

### 3.1 功能版本API

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /baas/apps/:appId/collections/:id/versions | 获取版本列表 |
| GET | /baas/apps/:appId/collections/:id/versions/:versionId | 获取版本详情 |
| POST | /baas/apps/:appId/collections/:id/versions | 创建新版本 |
| PUT | /baas/apps/:appId/collections/:id/versions/:versionId/publish | 发布版本 |
| PUT | /baas/apps/:appId/collections/:id/versions/:versionId/rollback | 回滚到此版本 |
| GET | /baas/apps/:appId/collections/:id/versions/compare | 对比两个版本 |

### 3.2 模块版本API

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /apps/:appId/modules/:moduleCode/versions | 获取版本列表 |
| GET | /apps/:appId/modules/:moduleCode/versions/:versionId | 获取版本详情 |
| POST | /apps/:appId/modules/:moduleCode/versions | 创建新版本 |
| PUT | /apps/:appId/modules/:moduleCode/versions/:versionId/publish | 发布版本 |
| PUT | /apps/:appId/modules/:moduleCode/versions/:versionId/rollback | 回滚到此版本 |
| GET | /apps/:appId/modules/:moduleCode/versions/compare | 对比两个版本 |

## 4. 版本号规则

采用语义化版本号 (Semantic Versioning)：
- 主版本号.次版本号.修订号 (如 1.2.3)
- 主版本号：不兼容的API修改
- 次版本号：向下兼容的功能性新增
- 修订号：向下兼容的问题修正

自动版本号递增规则：
- 默认递增修订号
- 用户可手动指定版本号

## 5. 状态流转

```
draft (草稿) -> published (已发布) -> deprecated (已废弃)
                    |
                    v
              (可回滚到此版本)
```

## 6. 前端页面设计

### 6.1 功能版本管理入口
- 在数据模型管理页面，每个已生成的功能增加"版本管理"按钮
- 点击进入版本管理页面

### 6.2 版本列表页面
- 显示所有版本（版本号、状态、创建时间、变更说明）
- 支持筛选（按状态、时间范围）
- 操作按钮：查看详情、发布、回滚、对比

### 6.3 版本详情页面
- 显示版本的完整字段结构
- 显示变更说明
- 显示创建人和时间

### 6.4 版本对比页面
- 左右对比两个版本的差异
- 高亮显示新增、删除、修改的字段
