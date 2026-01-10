-- =====================================================
-- APP中台管理系统 - TiDB数据库索引优化脚本
-- 创建时间: 2026-01-11
-- 说明: 适用于TiDB的索引创建语句（不支持IF NOT EXISTS和DESC）
-- =====================================================

-- 1. apps表索引优化
CREATE INDEX idx_apps_status_created ON apps(status, created_at);

-- 2. app_modules表索引优化
CREATE INDEX idx_app_modules_app_status ON app_modules(app_id, status);
CREATE INDEX idx_app_modules_module_code ON app_modules(module_code, app_id);

-- 3. monitor_metrics表索引优化
CREATE INDEX idx_monitor_metrics_app_metric ON monitor_metrics(app_id, metric_name);
CREATE INDEX idx_monitor_metrics_app_time ON monitor_metrics(app_id, created_at);

-- 4. push_records表索引优化
CREATE INDEX idx_push_records_app_status ON push_records(app_id, status);
CREATE INDEX idx_push_records_status_scheduled ON push_records(status, scheduled_at);

-- 5. events表索引优化
CREATE INDEX idx_events_app_code_time ON events(app_id, event_code, created_at);

-- 6. logs表索引优化
CREATE INDEX idx_logs_app_level_time ON logs(app_id, level, created_at);

-- 7. files表索引优化
CREATE INDEX idx_files_app_created ON files(app_id, created_at);

-- 8. versions表索引优化
CREATE INDEX idx_versions_app_status_code ON versions(app_id, status, version_code);

-- 9. configs表索引优化
CREATE INDEX idx_configs_app_published ON configs(app_id, is_published);

-- 10. messages表索引优化
CREATE INDEX idx_messages_app_status ON messages(app_id, status);
CREATE INDEX idx_messages_user_status ON messages(user_id, status);

-- 11. audit_logs表索引优化
CREATE INDEX idx_audit_logs_app_action ON audit_logs(app_id, action);
CREATE INDEX idx_audit_logs_user_time ON audit_logs(user_id, created_at);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
