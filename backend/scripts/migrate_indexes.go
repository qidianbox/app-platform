package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// IndexDefinition 索引定义
type IndexDefinition struct {
	Table   string
	Name    string
	Columns string
	Comment string
}

// 索引列表
var indexes = []IndexDefinition{
	// apps表
	{Table: "apps", Name: "idx_apps_status_created", Columns: "status, created_at DESC", Comment: "APP列表查询优化"},

	// app_modules表
	{Table: "app_modules", Name: "idx_app_modules_app_status", Columns: "app_id, status", Comment: "获取APP模块列表"},
	{Table: "app_modules", Name: "idx_app_modules_module_code", Columns: "module_code, app_id", Comment: "按模块代码查询"},

	// users表
	{Table: "users", Name: "idx_users_app_status", Columns: "app_id, status", Comment: "用户列表查询"},

	// monitor_metrics表
	{Table: "monitor_metrics", Name: "idx_monitor_metrics_app_metric", Columns: "app_id, metric_name", Comment: "按APP和指标查询"},
	{Table: "monitor_metrics", Name: "idx_monitor_metrics_app_time", Columns: "app_id, created_at DESC", Comment: "时间范围查询"},

	// monitor_alerts表
	{Table: "monitor_alerts", Name: "idx_monitor_alerts_app_status", Columns: "app_id, status", Comment: "告警状态查询"},
	{Table: "monitor_alerts", Name: "idx_monitor_alerts_app_metric_active", Columns: "app_id, metric_name, is_active", Comment: "活跃告警规则"},

	// push_records表
	{Table: "push_records", Name: "idx_push_records_app_status", Columns: "app_id, status", Comment: "推送记录查询"},
	{Table: "push_records", Name: "idx_push_records_status_scheduled", Columns: "status, scheduled_at", Comment: "定时推送查询"},

	// events表
	{Table: "events", Name: "idx_events_app_code_time", Columns: "app_id, event_code, created_at DESC", Comment: "事件统计查询"},

	// logs表
	{Table: "logs", Name: "idx_logs_app_level_time", Columns: "app_id, level, created_at DESC", Comment: "日志查询"},

	// files表
	{Table: "files", Name: "idx_files_app_created", Columns: "app_id, created_at DESC", Comment: "文件列表查询"},

	// versions表
	{Table: "versions", Name: "idx_versions_app_status_code", Columns: "app_id, status, version_code DESC", Comment: "版本检查更新"},

	// configs表
	{Table: "configs", Name: "idx_configs_app_published", Columns: "app_id, is_published", Comment: "已发布配置查询"},

	// messages表
	{Table: "messages", Name: "idx_messages_app_status", Columns: "app_id, status", Comment: "消息列表查询"},
	{Table: "messages", Name: "idx_messages_user_status", Columns: "user_id, status", Comment: "用户消息查询"},

	// audit_logs表
	{Table: "audit_logs", Name: "idx_audit_logs_app_action", Columns: "app_id, action", Comment: "审计日志查询"},
	{Table: "audit_logs", Name: "idx_audit_logs_user_time", Columns: "user_id, created_at DESC", Comment: "用户操作记录"},
	{Table: "audit_logs", Name: "idx_audit_logs_created_at", Columns: "created_at DESC", Comment: "时间范围查询"},
}

func main() {
	// 从环境变量获取数据库连接信息
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/app_platform?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}

	fmt.Println("=== 数据库索引迁移工具 ===")
	fmt.Println()

	// 获取命令行参数
	action := "check"
	if len(os.Args) > 1 {
		action = os.Args[1]
	}

	switch action {
	case "check":
		checkIndexes(db)
	case "create":
		createIndexes(db)
	case "drop":
		dropIndexes(db)
	default:
		fmt.Println("用法: migrate_indexes [check|create|drop]")
		fmt.Println("  check  - 检查索引状态")
		fmt.Println("  create - 创建缺失的索引")
		fmt.Println("  drop   - 删除所有自定义索引")
	}
}

// checkIndexes 检查索引状态
func checkIndexes(db *sql.DB) {
	fmt.Println("检查索引状态...")
	fmt.Println()

	existingCount := 0
	missingCount := 0

	for _, idx := range indexes {
		exists, err := indexExists(db, idx.Table, idx.Name)
		if err != nil {
			fmt.Printf("  [错误] %s.%s: %v\n", idx.Table, idx.Name, err)
			continue
		}

		if exists {
			fmt.Printf("  [存在] %s.%s\n", idx.Table, idx.Name)
			existingCount++
		} else {
			fmt.Printf("  [缺失] %s.%s (%s)\n", idx.Table, idx.Name, idx.Comment)
			missingCount++
		}
	}

	fmt.Println()
	fmt.Printf("总计: %d 个索引存在, %d 个索引缺失\n", existingCount, missingCount)
}

// createIndexes 创建索引
func createIndexes(db *sql.DB) {
	fmt.Println("创建索引...")
	fmt.Println()

	successCount := 0
	skipCount := 0
	failCount := 0

	for _, idx := range indexes {
		// 检查表是否存在
		if !tableExists(db, idx.Table) {
			fmt.Printf("  [跳过] %s.%s: 表不存在\n", idx.Table, idx.Name)
			skipCount++
			continue
		}

		// 检查索引是否已存在
		exists, _ := indexExists(db, idx.Table, idx.Name)
		if exists {
			fmt.Printf("  [跳过] %s.%s: 索引已存在\n", idx.Table, idx.Name)
			skipCount++
			continue
		}

		// 创建索引
		sql := fmt.Sprintf("CREATE INDEX %s ON %s (%s)", idx.Name, idx.Table, idx.Columns)
		_, err := db.Exec(sql)
		if err != nil {
			fmt.Printf("  [失败] %s.%s: %v\n", idx.Table, idx.Name, err)
			failCount++
		} else {
			fmt.Printf("  [成功] %s.%s\n", idx.Table, idx.Name)
			successCount++
		}
	}

	fmt.Println()
	fmt.Printf("总计: %d 个成功, %d 个跳过, %d 个失败\n", successCount, skipCount, failCount)
}

// dropIndexes 删除索引
func dropIndexes(db *sql.DB) {
	fmt.Println("删除索引...")
	fmt.Println()

	successCount := 0
	skipCount := 0
	failCount := 0

	for _, idx := range indexes {
		// 检查索引是否存在
		exists, _ := indexExists(db, idx.Table, idx.Name)
		if !exists {
			fmt.Printf("  [跳过] %s.%s: 索引不存在\n", idx.Table, idx.Name)
			skipCount++
			continue
		}

		// 删除索引
		sql := fmt.Sprintf("DROP INDEX %s ON %s", idx.Name, idx.Table)
		_, err := db.Exec(sql)
		if err != nil {
			fmt.Printf("  [失败] %s.%s: %v\n", idx.Table, idx.Name, err)
			failCount++
		} else {
			fmt.Printf("  [成功] %s.%s\n", idx.Table, idx.Name)
			successCount++
		}
	}

	fmt.Println()
	fmt.Printf("总计: %d 个成功, %d 个跳过, %d 个失败\n", successCount, skipCount, failCount)
}

// indexExists 检查索引是否存在
func indexExists(db *sql.DB, table, indexName string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM information_schema.statistics 
		WHERE table_schema = DATABASE() 
		AND table_name = ? 
		AND index_name = ?
	`
	var count int
	err := db.QueryRow(query, table, indexName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// tableExists 检查表是否存在
func tableExists(db *sql.DB, table string) bool {
	query := `
		SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_schema = DATABASE() 
		AND table_name = ?
	`
	var count int
	err := db.QueryRow(query, table).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

// 辅助函数：移除DESC关键字用于索引创建
func normalizeColumns(columns string) string {
	// MySQL不支持在CREATE INDEX中使用DESC
	columns = strings.ReplaceAll(columns, " DESC", "")
	columns = strings.ReplaceAll(columns, " ASC", "")
	return columns
}
