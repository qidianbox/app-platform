# 阿里云SAE部署指南

## 概述

本目录包含将APP中台管理系统部署到阿里云SAE（Serverless应用引擎）的所有配置文件和脚本。

## 架构说明

```
                     ┌─────────────────┐
                     │    用户访问      │
                     └────────┬────────┘
                              │
              ┌───────────────┼───────────────┐
              │               │               │
        ┌─────▼─────┐   ┌─────▼─────┐   ┌─────▼─────┐
        │  OSS+CDN  │   │    SLB    │   │  OSS存储   │
        │  (前端)   │   │  (负载)   │   │  (文件)   │
        └───────────┘   └─────┬─────┘   └───────────┘
                              │
                        ┌─────▼─────┐
                        │    SAE    │
                        │ (Go后端)  │
                        │ 自动扩缩容 │
                        └─────┬─────┘
                              │
                        ┌─────▼─────┐
                        │  RDS MySQL │
                        │  (数据库)  │
                        └───────────┘
```

## 文件说明

| 文件 | 说明 |
|------|------|
| `main.tf` | Terraform主配置，定义所有云资源 |
| `variables.tf` | Terraform变量定义 |
| `terraform.tfvars.example` | 变量示例文件 |
| `deploy.sh` | 自动化部署脚本 |

## 前置条件

### 1. 安装必要工具

```bash
# Terraform
wget https://releases.hashicorp.com/terraform/1.6.0/terraform_1.6.0_linux_amd64.zip
unzip terraform_1.6.0_linux_amd64.zip
sudo mv terraform /usr/local/bin/

# Docker
curl -fsSL https://get.docker.com | sh

# 阿里云CLI
curl -O https://aliyuncli.alicdn.com/aliyun-cli-linux-latest-amd64.tgz
tar xzvf aliyun-cli-linux-latest-amd64.tgz
sudo mv aliyun /usr/local/bin/
```

### 2. 创建RAM用户并获取AccessKey

登录阿里云控制台，创建RAM用户并授予以下权限：

- `AliyunSAEFullAccess` - SAE完全访问
- `AliyunVPCFullAccess` - VPC完全访问
- `AliyunRDSFullAccess` - RDS完全访问
- `AliyunOSSFullAccess` - OSS完全访问
- `AliyunCRFullAccess` - 容器镜像服务完全访问
- `AliyunSLBFullAccess` - 负载均衡完全访问

### 3. 配置环境变量

```bash
export ALICLOUD_ACCESS_KEY="your-access-key-id"
export ALICLOUD_SECRET_KEY="your-access-key-secret"
export REGION="cn-hangzhou"
```

## 快速部署

### 方式一：一键部署（推荐）

```bash
# 1. 复制并编辑配置文件
cp terraform.tfvars.example terraform.tfvars
vim terraform.tfvars  # 填写实际配置

# 2. 执行部署
./deploy.sh deploy
```

### 方式二：分步部署

```bash
# 1. 部署基础设施
./deploy.sh infra

# 2. 构建并推送镜像
./deploy.sh image

# 3. 更新SAE应用
./deploy.sh update

# 4. 部署前端
./deploy.sh frontend
```

## 配置说明

### terraform.tfvars 配置项

```hcl
# 阿里云认证
access_key = "LTAI5t..."           # AccessKey ID
secret_key = "..."                  # AccessKey Secret
region     = "cn-hangzhou"          # 部署区域

# SAE配置
sae_cpu          = 1000             # CPU（毫核，1000=1核）
sae_memory       = 2048             # 内存（MB）
sae_replicas     = 2                # 初始实例数
sae_min_replicas = 1                # 最小实例数
sae_max_replicas = 20               # 最大实例数

# RDS配置
rds_instance_type = "mysql.n2.small.1"  # 实例规格
rds_storage       = 20                   # 存储空间(GB)
db_password       = "YourPassword123!"   # 数据库密码

# 应用配置
jwt_secret = "your-32-char-secret"  # JWT密钥
```

### 区域选择

| 区域 | Region ID | 推荐场景 |
|------|-----------|----------|
| 华东1（杭州） | cn-hangzhou | 默认推荐 |
| 华东2（上海） | cn-shanghai | 金融业务 |
| 华北2（北京） | cn-beijing | 北方用户 |
| 华南1（深圳） | cn-shenzhen | 南方用户 |

## 成本估算

| 资源 | 规格 | 月费用（约） |
|------|------|-------------|
| SAE | 1核2G × 2实例 | ¥200 |
| RDS MySQL | 1核2G | ¥200 |
| OSS | 10GB | ¥5 |
| SLB | 按量 | ¥50 |
| **合计** | | **¥455/月** |

*注：实际费用根据使用量浮动，SAE按实际运行时间计费*

## 自动扩缩容

SAE已配置自动扩缩容策略：

- **CPU阈值**：70%时触发扩容
- **内存阈值**：80%时触发扩容
- **扩容步长**：每次增加2个实例
- **缩容步长**：每次减少1个实例
- **扩容冷却**：60秒
- **缩容冷却**：300秒

## 常用命令

```bash
# 查看部署状态
terraform show

# 更新应用（代码变更后）
./deploy.sh update

# 仅更新前端
./deploy.sh frontend

# 销毁所有资源
./deploy.sh cleanup

# 查看SAE日志
aliyun sae DescribeApplicationSlbStatus --AppId <app-id>
```

## 故障排查

### 1. SAE应用启动失败

```bash
# 查看应用日志
aliyun sae DescribeApplicationStatus --AppId <app-id>

# 常见原因：
# - 镜像拉取失败：检查ACR权限
# - 健康检查失败：检查/api/v1/health接口
# - 环境变量错误：检查数据库连接配置
```

### 2. 数据库连接失败

```bash
# 检查安全组
# 确保SAE和RDS在同一VPC
# 检查RDS白名单是否包含SAE网段
```

### 3. 前端访问404

```bash
# 检查OSS静态网站配置
# 确保index.html和error.html都指向index.html（SPA应用）
```

## 后续优化建议

1. **添加CDN加速**：为OSS配置CDN，提升前端访问速度
2. **配置自定义域名**：绑定已备案域名，配置HTTPS
3. **启用日志服务**：接入SLS，实现日志分析和告警
4. **配置监控告警**：设置CPU/内存/请求量告警

## 迁移到多云

当业务发展需要多云部署时，可参考 `../docs/deployment-optimization.md` 中的K8s多云方案进行迁移。
