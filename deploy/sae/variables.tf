# Terraform变量定义 - SAE部署

# ==================== 阿里云认证 ====================

variable "access_key" {
  description = "阿里云AccessKey ID"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "阿里云AccessKey Secret"
  type        = string
  sensitive   = true
}

variable "region" {
  description = "部署区域"
  type        = string
  default     = "cn-hangzhou"
}

# ==================== 项目配置 ====================

variable "project_name" {
  description = "项目名称（用于资源命名）"
  type        = string
  default     = "app-platform"
}

variable "environment" {
  description = "环境标识"
  type        = string
  default     = "production"
}

# ==================== SAE配置 ====================

variable "backend_image_url" {
  description = "后端镜像地址"
  type        = string
}

variable "sae_cpu" {
  description = "SAE实例CPU（毫核，1000=1核）"
  type        = number
  default     = 1000  # 1核
}

variable "sae_memory" {
  description = "SAE实例内存（MB）"
  type        = number
  default     = 2048  # 2GB
}

variable "sae_replicas" {
  description = "SAE初始实例数"
  type        = number
  default     = 2
}

variable "sae_min_replicas" {
  description = "SAE最小实例数"
  type        = number
  default     = 1
}

variable "sae_max_replicas" {
  description = "SAE最大实例数"
  type        = number
  default     = 20
}

# ==================== RDS配置 ====================

variable "rds_instance_type" {
  description = "RDS实例规格"
  type        = string
  default     = "mysql.n2.small.1"  # 1核2G
}

variable "rds_storage" {
  description = "RDS存储空间(GB)"
  type        = number
  default     = 20
}

variable "db_name" {
  description = "数据库名称"
  type        = string
  default     = "app_platform"
}

variable "db_user" {
  description = "数据库用户名"
  type        = string
  default     = "app_platform"
}

variable "db_password" {
  description = "数据库密码"
  type        = string
  sensitive   = true
}

# ==================== 应用配置 ====================

variable "jwt_secret" {
  description = "JWT密钥"
  type        = string
  sensitive   = true
}
