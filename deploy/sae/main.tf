# Terraform配置 - 阿里云SAE部署
# 版本：1.0
# 作者：Manus AI
# 说明：使用Serverless应用引擎(SAE)部署Go后端应用

terraform {
  required_version = ">= 1.0"
  required_providers {
    alicloud = {
      source  = "aliyun/alicloud"
      version = "~> 1.200"
    }
  }
}

# 阿里云Provider配置
provider "alicloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

# 数据源：获取可用区
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

# ==================== VPC网络 ====================

# 创建VPC
resource "alicloud_vpc" "main" {
  vpc_name   = "${var.project_name}-vpc"
  cidr_block = "10.0.0.0/16"
  tags = {
    Project = var.project_name
    Env     = var.environment
  }
}

# 创建交换机 - 应用层
resource "alicloud_vswitch" "app" {
  vpc_id       = alicloud_vpc.main.id
  cidr_block   = "10.0.1.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = "${var.project_name}-app-vswitch"
  tags = {
    Project = var.project_name
    Layer   = "app"
  }
}

# 创建交换机 - 数据层
resource "alicloud_vswitch" "data" {
  vpc_id       = alicloud_vpc.main.id
  cidr_block   = "10.0.2.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = "${var.project_name}-data-vswitch"
  tags = {
    Project = var.project_name
    Layer   = "data"
  }
}

# ==================== 安全组 ====================

# SAE安全组
resource "alicloud_security_group" "sae" {
  name        = "${var.project_name}-sae-sg"
  vpc_id      = alicloud_vpc.main.id
  description = "Security group for SAE applications"
  tags = {
    Project = var.project_name
  }
}

# 安全组规则 - 允许内部通信
resource "alicloud_security_group_rule" "internal" {
  type              = "ingress"
  ip_protocol       = "all"
  port_range        = "-1/-1"
  security_group_id = alicloud_security_group.sae.id
  cidr_ip           = "10.0.0.0/16"
}

# 安全组规则 - 允许HTTP
resource "alicloud_security_group_rule" "http" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "80/80"
  security_group_id = alicloud_security_group.sae.id
  cidr_ip           = "0.0.0.0/0"
}

# 安全组规则 - 允许HTTPS
resource "alicloud_security_group_rule" "https" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "443/443"
  security_group_id = alicloud_security_group.sae.id
  cidr_ip           = "0.0.0.0/0"
}

# ==================== RDS MySQL ====================

# 创建RDS实例
resource "alicloud_db_instance" "main" {
  engine               = "MySQL"
  engine_version       = "8.0"
  instance_type        = var.rds_instance_type
  instance_storage     = var.rds_storage
  instance_name        = "${var.project_name}-rds"
  vswitch_id           = alicloud_vswitch.data.id
  security_ips         = ["10.0.0.0/16"]
  instance_charge_type = "Postpaid"
  
  tags = {
    Project = var.project_name
    Env     = var.environment
  }
}

# 创建数据库
resource "alicloud_db_database" "main" {
  instance_id   = alicloud_db_instance.main.id
  name          = var.db_name
  character_set = "utf8mb4"
}

# 创建数据库账号
resource "alicloud_db_account" "main" {
  db_instance_id   = alicloud_db_instance.main.id
  account_name     = var.db_user
  account_password = var.db_password
  account_type     = "Super"
}

# 授权数据库访问
resource "alicloud_db_account_privilege" "main" {
  instance_id  = alicloud_db_instance.main.id
  account_name = alicloud_db_account.main.account_name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.main.name]
}

# ==================== SAE命名空间 ====================

resource "alicloud_sae_namespace" "main" {
  namespace_id          = "${var.region}:${var.project_name}"
  namespace_name        = var.project_name
  namespace_description = "APP中台管理系统命名空间"
}

# ==================== SAE应用 ====================

resource "alicloud_sae_application" "backend" {
  app_name          = "${var.project_name}-backend"
  app_description   = "APP中台管理系统后端服务"
  namespace_id      = alicloud_sae_namespace.main.id
  
  # 镜像配置
  image_url         = var.backend_image_url
  package_type      = "Image"
  
  # 实例规格
  cpu               = var.sae_cpu
  memory            = var.sae_memory
  replicas          = var.sae_replicas
  
  # 网络配置
  vpc_id            = alicloud_vpc.main.id
  vswitch_id        = alicloud_vswitch.app.id
  security_group_id = alicloud_security_group.sae.id
  
  # 环境变量
  envs = jsonencode([
    {
      name  = "DB_HOST"
      value = alicloud_db_instance.main.connection_string
    },
    {
      name  = "DB_PORT"
      value = tostring(alicloud_db_instance.main.port)
    },
    {
      name  = "DB_USER"
      value = var.db_user
    },
    {
      name  = "DB_PASSWORD"
      value = var.db_password
    },
    {
      name  = "DB_NAME"
      value = var.db_name
    },
    {
      name  = "JWT_SECRET"
      value = var.jwt_secret
    },
    {
      name  = "GIN_MODE"
      value = "release"
    }
  ])
  
  # 健康检查
  liveness = jsonencode({
    exec = {
      command = ["curl", "-f", "http://localhost:8080/api/v1/health"]
    }
    initialDelaySeconds = 30
    periodSeconds       = 10
    timeoutSeconds      = 5
  })
  
  readiness = jsonencode({
    exec = {
      command = ["curl", "-f", "http://localhost:8080/api/v1/health"]
    }
    initialDelaySeconds = 10
    periodSeconds       = 5
    timeoutSeconds      = 3
  })
  
  # 自动扩缩容配置
  auto_config = true
  
  # 日志配置
  sls_configs = jsonencode([
    {
      type       = "stdout"
      logDir     = ""
      project    = ""
      logstore   = "${var.project_name}-stdout"
      createIfNotExist = true
    }
  ])
  
  # 时区设置
  timezone = "Asia/Shanghai"
  
  tags = {
    Project = var.project_name
    Env     = var.environment
  }
}

# ==================== SAE自动扩缩容 ====================

resource "alicloud_sae_application_scaling_rule" "backend" {
  app_id            = alicloud_sae_application.backend.id
  scaling_rule_name = "${var.project_name}-scaling"
  scaling_rule_type = "mix"
  
  scaling_rule_metric {
    max_replicas = var.sae_max_replicas
    min_replicas = var.sae_min_replicas
    
    metrics {
      metric_type                       = "CPU"
      metric_target_average_utilization = 70
    }
    
    metrics {
      metric_type                       = "MEMORY"
      metric_target_average_utilization = 80
    }
    
    scale_up_rules {
      step                         = 2
      stabilization_window_seconds = 60
    }
    
    scale_down_rules {
      step                         = 1
      stabilization_window_seconds = 300
    }
  }
  
  scaling_rule_enable = true
}

# ==================== SAE公网SLB ====================

resource "alicloud_sae_load_balancer_internet" "backend" {
  app_id = alicloud_sae_application.backend.id
  
  internet {
    protocol    = "HTTP"
    port        = 80
    target_port = 8080
  }
}

# ==================== OSS存储桶 ====================

# 创建OSS存储桶（前端静态资源）
resource "alicloud_oss_bucket" "frontend" {
  bucket = "${var.project_name}-frontend-${var.region}"
  acl    = "public-read"
  
  website {
    index_document = "index.html"
    error_document = "index.html"
  }
  
  cors_rule {
    allowed_origins = ["*"]
    allowed_methods = ["GET", "HEAD"]
    allowed_headers = ["*"]
    max_age_seconds = 3600
  }
  
  tags = {
    Project = var.project_name
    Env     = var.environment
  }
}

# ==================== ACR容器镜像仓库 ====================

# 创建命名空间
resource "alicloud_cr_namespace" "main" {
  name               = var.project_name
  auto_create        = false
  default_visibility = "PRIVATE"
}

# 创建镜像仓库
resource "alicloud_cr_repo" "backend" {
  namespace = alicloud_cr_namespace.main.name
  name      = "backend"
  summary   = "Backend service image"
  repo_type = "PRIVATE"
}

# ==================== 输出 ====================

output "vpc_id" {
  description = "VPC ID"
  value       = alicloud_vpc.main.id
}

output "sae_namespace_id" {
  description = "SAE命名空间ID"
  value       = alicloud_sae_namespace.main.id
}

output "sae_app_id" {
  description = "SAE应用ID"
  value       = alicloud_sae_application.backend.id
}

output "sae_app_name" {
  description = "SAE应用名称"
  value       = alicloud_sae_application.backend.app_name
}

output "backend_slb_address" {
  description = "后端SLB公网地址"
  value       = alicloud_sae_load_balancer_internet.backend.internet_ip
}

output "rds_connection_string" {
  description = "RDS连接地址"
  value       = alicloud_db_instance.main.connection_string
}

output "rds_port" {
  description = "RDS端口"
  value       = alicloud_db_instance.main.port
}

output "acr_namespace" {
  description = "ACR命名空间"
  value       = alicloud_cr_namespace.main.name
}

output "acr_repo_domain" {
  description = "ACR仓库域名"
  value       = "registry.${var.region}.aliyuncs.com/${alicloud_cr_namespace.main.name}/backend"
}

output "oss_bucket" {
  description = "OSS存储桶名称"
  value       = alicloud_oss_bucket.frontend.bucket
}

output "oss_endpoint" {
  description = "OSS访问端点"
  value       = alicloud_oss_bucket.frontend.extranet_endpoint
}

output "frontend_url" {
  description = "前端访问地址"
  value       = "http://${alicloud_oss_bucket.frontend.bucket}.${alicloud_oss_bucket.frontend.extranet_endpoint}"
}
