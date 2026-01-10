-- =====================================================
-- APP中台管理系统 - 数据库索引优化脚本
-- 创建时间: 2026-01-11
-- 说明: 基于压测分析结果，为高频查询字段添加复合索引
-- =====================================================

-- 1. apps表索引优化
-- 用于APP列表查询（按状态筛选、按创建时间排序）
CREATE INDEX IF NOT EXISTS idx_apps_status_created ON apps(status, created_at DESC);

-- 2. app_modules表索引优化
-- 用于获取APP启用的模块列表
CREATE INDEX IF NOT EXISTS idx_app_modules_app_status ON app_modules(app_id, status);
-- 用于按模块代码查询
CREATE INDEX IF NOT EXISTS idx_app_modules_module_code ON app_modules(module_code, app_id);

-- 3. users表索引优化
-- 用于用户搜索（按名称或邮箱）
CREATE INDEX IF NOT EXISTS idx_users_name ON users(name(50));
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email(100));
-- 用于活跃用户统计
CREATE INDEX IF NOT EXISTS idx_users_last_signed ON users(lastSignedIn);
-- 用于按角色筛选
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- 4. monitor_metrics表索引优化（高频写入表，需要谨慎添加索引）
-- 用于按APP和指标名称查询
CREATE INDEX IF NOT EXISTS idx_monitor_metrics_app_metric ON monitor_metrics(app_id, metric_name);
-- 用于时间范围查询
CREATE INDEX IF NOT EXISTS idx_monitor_metrics_app_time ON monitor_metrics(app_id, created_at DESC);
-- 复合索引用于统计查询
CREATE INDEX IF NOT EXISTS idx_monitor_metrics_app_metric_time ON monitor_metrics(app_id, metric_name, created_at DESC);

-- 5. monitor_alerts表索引优化
-- 用于按APP和状态查询告警
CREATE INDEX IF NOT EXISTS idx_monitor_alerts_app_status ON monitor_alerts(app_id, status);
-- 用于检查活跃告警规则
CREATE INDEX IF NOT EXISTS idx_monitor_alerts_app_metric_active ON monitor_alerts(app_id, metric_name, is_active);

-- 6. push_records表索引优化
-- 用于按APP和状态查询推送记录
CREATE INDEX IF NOT EXISTS idx_push_records_app_status ON push_records(app_id, status);
-- 用于定时任务查询待发送推送
CREATE INDEX IF NOT EXISTS idx_push_records_status_scheduled ON push_records(status, scheduled_at);

-- 7. events表索引优化（高频写入表）
-- 用于事件统计查询
CREATE INDEX IF NOT EXISTS idx_events_app_code_time ON events(app_id, event_code, created_at DESC);
-- 用于用户行为分析
CREATE INDEX IF NOT EXISTS idx_events_app_user_time ON events(app_id, user_id, created_at DESC);

-- 8. logs表索引优化（高频写入表）
-- 用于日志查询（按级别和模块筛选）
CREATE INDEX IF NOT EXISTS idx_logs_app_level_time ON logs(app_id, level, created_at DESC);
-- 用于按模块查询
CREATE INDEX IF NOT EXISTS idx_logs_app_module_time ON logs(app_id, module, created_at DESC);

-- 9. files表索引优化
-- 用于按APP和MIME类型查询
CREATE INDEX IF NOT EXISTS idx_files_app_mime ON files(app_id, mime_type(50));
-- 用于按上传时间排序
CREATE INDEX IF NOT EXISTS idx_files_app_created ON files(app_id, created_at DESC);

-- 10. versions表索引优化
-- 用于检查更新（获取最新发布版本）
CREATE INDEX IF NOT EXISTS idx_versions_app_status_code ON versions(app_id, status, version_code DESC);

-- 11. configs表索引优化
-- 用于按APP和配置键查询
CREATE INDEX IF NOT EXISTS idx_configs_app_key ON configs(app_id, config_key(100));
-- 用于获取已发布配置
CREATE INDEX IF NOT EXISTS idx_configs_app_published ON configs(app_id, is_published);

-- 12. messages表索引优化
-- 用于按APP和状态查询消息
CREATE INDEX IF NOT EXISTS idx_messages_app_status ON messages(app_id, status);
-- 用于按用户查询消息
CREATE INDEX IF NOT EXISTS idx_messages_user_status ON messages(user_id, status);

-- 13. audit_logs表索引优化
-- 用于按APP和操作类型查询
CREATE INDEX IF NOT EXISTS idx_audit_logs_app_action ON audit_logs(app_id, action);
-- 用于按用户查询操作记录
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_time ON audit_logs(user_id, created_at DESC);
-- 用于按资源类型查询
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
-- 用于时间范围查询
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- =====================================================
-- 索引使用建议:
-- 1. 对于高频写入表（monitor_metrics, events, logs），索引会影响写入性能
--    建议在业务低峰期执行索引创建
-- 2. 定期使用 ANALYZE TABLE 更新索引统计信息
-- 3. 使用 EXPLAIN 验证查询是否使用了索引
-- 4. 监控索引使用情况，删除未使用的索引
-- =====================================================

-- 查看索引使用情况的示例查询:
-- SELECT * FROM sys.schema_index_statistics WHERE table_schema = 'your_database';

-- 查看表的索引信息:
-- SHOW INDEX FROM table_name;

-- 分析表以更新统计信息:
-- ANALYZE TABLE apps, app_modules, users, monitor_metrics, monitor_alerts;
