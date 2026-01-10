package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type IndexDef struct {
	Table   string
	Name    string
	Columns string
}

var indexes = []IndexDef{
	// apps表
	{"apps", "idx_apps_status_created", "status, created_at"},
	// app_modules表
	{"app_modules", "idx_app_modules_app_status", "app_id, status"},
	{"app_modules", "idx_app_modules_module_code", "module_code, app_id"},
	// monitor_metrics表
	{"monitor_metrics", "idx_monitor_metrics_app_metric", "app_id, metric_name"},
	{"monitor_metrics", "idx_monitor_metrics_app_time", "app_id, created_at"},
	// push_records表
	{"push_records", "idx_push_records_app_status", "app_id, status"},
	// events表
	{"events", "idx_events_app_code_time", "app_id, event_code, created_at"},
	// logs表
	{"logs", "idx_logs_app_level_time", "app_id, level, created_at"},
	// files表
	{"files", "idx_files_app_created", "app_id, created_at"},
	// versions表
	{"versions", "idx_versions_app_status_code", "app_id, status, version_code"},
	// configs表
	{"configs", "idx_configs_app_published", "app_id, is_published"},
	// messages表
	{"messages", "idx_messages_app_status", "app_id, status"},
	// audit_logs表
	{"audit_logs", "idx_audit_logs_app_action", "app_id, action"},
	{"audit_logs", "idx_audit_logs_user_time", "user_id, created_at"},
	{"audit_logs", "idx_audit_logs_created_at", "created_at"},
}

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	// 解析URL格式的连接字符串
	dsn, err := parseDBURL(databaseURL)
	if err != nil {
		log.Fatalf("解析连接字符串失败: %v", err)
	}

	// 注册TLS配置
	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: "gateway03.us-east-1.prod.aws.tidbcloud.com",
	})

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("连接测试失败: %v", err)
	}

	fmt.Println("=== 数据库索引迁移 ===")
	fmt.Println()
	success, skip, fail := 0, 0, 0

	for _, idx := range indexes {
		// 检查表是否存在
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", idx.Table).Scan(&count)
		if err != nil || count == 0 {
			fmt.Printf("[跳过] %s.%s: 表不存在\n", idx.Table, idx.Name)
			skip++
			continue
		}

		// 检查索引是否存在
		err = db.QueryRow("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?", idx.Table, idx.Name).Scan(&count)
		if err == nil && count > 0 {
			fmt.Printf("[跳过] %s.%s: 已存在\n", idx.Table, idx.Name)
			skip++
			continue
		}

		// 创建索引
		sqlStmt := fmt.Sprintf("CREATE INDEX %s ON %s (%s)", idx.Name, idx.Table, idx.Columns)
		_, err = db.Exec(sqlStmt)
		if err != nil {
			fmt.Printf("[失败] %s.%s: %v\n", idx.Table, idx.Name, err)
			fail++
		} else {
			fmt.Printf("[成功] %s.%s\n", idx.Table, idx.Name)
			success++
		}
	}

	fmt.Printf("\n总计: %d成功, %d跳过, %d失败\n", success, skip, fail)
}

// parseDBURL 将mysql://格式转换为Go MySQL DSN格式
func parseDBURL(dbURL string) (string, error) {
	u, err := url.Parse(dbURL)
	if err != nil {
		return "", err
	}

	password, _ := u.User.Password()
	host := u.Host
	dbName := strings.TrimPrefix(u.Path, "/")

	// 构建DSN: user:password@tcp(host)/dbname?params
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=tidb&timeout=30s",
		u.User.Username(),
		password,
		host,
		dbName,
	)

	return dsn, nil
}
