-- APP中台管理系统数据库备份
-- 日期: 2026-01-10

-- 注意：此备份仅包含表结构，实际数据存储在Manus平台数据库中

-- 应用表
CREATE TABLE IF NOT EXISTS apps (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  app_id VARCHAR(50) NOT NULL UNIQUE,
  app_secret VARCHAR(100) NOT NULL,
  package_name VARCHAR(100),
  description TEXT,
  status TINYINT DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 应用模块关联表
CREATE TABLE IF NOT EXISTS app_modules (
  id INT AUTO_INCREMENT PRIMARY KEY,
  app_id INT NOT NULL,
  module_code VARCHAR(50) NOT NULL,
  source_module VARCHAR(50),
  enabled TINYINT DEFAULT 1,
  config JSON,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 模块模板表
CREATE TABLE IF NOT EXISTS module_templates (
  id INT AUTO_INCREMENT PRIMARY KEY,
  code VARCHAR(50) NOT NULL UNIQUE,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  category VARCHAR(50),
  icon VARCHAR(100),
  config_schema JSON,
  default_config JSON,
  status TINYINT DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 模块配置表
CREATE TABLE IF NOT EXISTS app_module_configs (
  id INT AUTO_INCREMENT PRIMARY KEY,
  app_id INT NOT NULL,
  module_code VARCHAR(50) NOT NULL,
  config_key VARCHAR(100) NOT NULL,
  config_value TEXT,
  environment VARCHAR(20) DEFAULT 'production',
  version INT DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_app_module_key_env (app_id, module_code, config_key, environment)
);

-- 版本管理表
CREATE TABLE IF NOT EXISTS app_versions (
  id INT AUTO_INCREMENT PRIMARY KEY,
  app_id INT NOT NULL,
  version VARCHAR(20) NOT NULL,
  version_code INT NOT NULL,
  description TEXT,
  download_url VARCHAR(500),
  file_size BIGINT,
  md5 VARCHAR(32),
  status TINYINT DEFAULT 0,
  force_update TINYINT DEFAULT 0,
  gray_scale INT DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  published_at TIMESTAMP NULL
);

-- 日志表
CREATE TABLE IF NOT EXISTS app_logs (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  app_id INT NOT NULL,
  level VARCHAR(10) NOT NULL,
  module VARCHAR(50),
  message TEXT,
  context JSON,
  user_id VARCHAR(50),
  device_id VARCHAR(100),
  ip VARCHAR(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_app_level_time (app_id, level, created_at)
);

-- 消息表
CREATE TABLE IF NOT EXISTS messages (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  app_id INT NOT NULL,
  user_id VARCHAR(50),
  type VARCHAR(20) NOT NULL,
  title VARCHAR(200),
  content TEXT,
  extra JSON,
  status TINYINT DEFAULT 0,
  read_at TIMESTAMP NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 推送任务表
CREATE TABLE IF NOT EXISTS push_tasks (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  app_id INT NOT NULL,
  title VARCHAR(200) NOT NULL,
  content TEXT,
  type VARCHAR(20) DEFAULT 'all',
  target_users JSON,
  status TINYINT DEFAULT 0,
  sent_count INT DEFAULT 0,
  success_count INT DEFAULT 0,
  fail_count INT DEFAULT 0,
  scheduled_at TIMESTAMP NULL,
  sent_at TIMESTAMP NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 事件埋点表
CREATE TABLE IF NOT EXISTS events (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  app_id INT NOT NULL,
  event_name VARCHAR(100) NOT NULL,
  event_type VARCHAR(50),
  user_id VARCHAR(50),
  device_id VARCHAR(100),
  properties JSON,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_app_event_time (app_id, event_name, created_at)
);

-- 文件存储表
CREATE TABLE IF NOT EXISTS files (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  app_id INT NOT NULL,
  filename VARCHAR(255) NOT NULL,
  original_name VARCHAR(255),
  file_path VARCHAR(500) NOT NULL,
  file_size BIGINT,
  mime_type VARCHAR(100),
  storage_type VARCHAR(20) DEFAULT 'local',
  uploader_id VARCHAR(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 监控指标表
CREATE TABLE IF NOT EXISTS monitor_metrics (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  app_id INT NOT NULL,
  metric_name VARCHAR(100) NOT NULL,
  metric_value DECIMAL(20,4),
  tags JSON,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_app_metric_time (app_id, metric_name, created_at)
);

-- 告警规则表
CREATE TABLE IF NOT EXISTS alert_rules (
  id INT AUTO_INCREMENT PRIMARY KEY,
  app_id INT NOT NULL,
  name VARCHAR(100) NOT NULL,
  metric_name VARCHAR(100) NOT NULL,
  condition_type VARCHAR(20) NOT NULL,
  threshold DECIMAL(20,4) NOT NULL,
  duration INT DEFAULT 60,
  notify_channels JSON,
  status TINYINT DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 告警记录表
CREATE TABLE IF NOT EXISTS alert_records (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  app_id INT NOT NULL,
  rule_id INT NOT NULL,
  metric_value DECIMAL(20,4),
  status TINYINT DEFAULT 0,
  resolved_at TIMESTAMP NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
