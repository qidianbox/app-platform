package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		fmt.Println("DATABASE_URL not set")
		os.Exit(1)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// 检查并创建 push_records 表
	createPushRecords := `
	CREATE TABLE IF NOT EXISTS push_records (
		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		app_id BIGINT UNSIGNED NOT NULL,
		title VARCHAR(255),
		content TEXT,
		target_type VARCHAR(50) DEFAULT 'all',
		target_ids TEXT,
		status VARCHAR(50) DEFAULT 'pending',
		sent_count INT DEFAULT 0,
		success_count INT DEFAULT 0,
		failed_count INT DEFAULT 0,
		scheduled_at DATETIME,
		sent_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		INDEX idx_push_records_app_id (app_id),
		INDEX idx_push_records_status (status),
		INDEX idx_push_records_app_status (app_id, status),
		INDEX idx_push_records_created_at (created_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	// 检查并创建 versions 表
	createVersions := `
	CREATE TABLE IF NOT EXISTS versions (
		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		app_id BIGINT UNSIGNED NOT NULL,
		version_name VARCHAR(50) NOT NULL,
		version_code INT NOT NULL,
		description TEXT,
		download_url VARCHAR(500),
		is_force_update TINYINT DEFAULT 0,
		status VARCHAR(50) DEFAULT 'draft',
		published_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		INDEX idx_versions_app_id (app_id),
		INDEX idx_versions_version_code (version_code),
		INDEX idx_versions_app_status_code (app_id, status, version_code),
		INDEX idx_versions_created_at (created_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	// 执行创建表
	tables := []struct {
		name string
		sql  string
	}{
		{"push_records", createPushRecords},
		{"versions", createVersions},
	}

	for _, t := range tables {
		// 检查表是否存在
		var tableName string
		err := db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", t.name).Scan(&tableName)
		if err == sql.ErrNoRows {
			// 表不存在，创建
			_, err = db.Exec(t.sql)
			if err != nil {
				fmt.Printf("[失败] 创建表 %s: %v\n", t.name, err)
			} else {
				fmt.Printf("[成功] 创建表 %s\n", t.name)
			}
		} else if err != nil {
			fmt.Printf("[错误] 检查表 %s: %v\n", t.name, err)
		} else {
			fmt.Printf("[跳过] 表 %s 已存在\n", t.name)
		}
	}

	// 添加缺失的索引
	indexes := []struct {
		table string
		name  string
		sql   string
	}{
		{"push_records", "idx_push_records_app_status", "CREATE INDEX idx_push_records_app_status ON push_records(app_id, status)"},
		{"versions", "idx_versions_app_status_code", "CREATE INDEX idx_versions_app_status_code ON versions(app_id, status, version_code)"},
	}

	for _, idx := range indexes {
		// 检查索引是否存在
		var indexName string
		err := db.QueryRow(`
			SELECT index_name FROM information_schema.statistics 
			WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?
		`, idx.table, idx.name).Scan(&indexName)
		
		if err == sql.ErrNoRows {
			// 索引不存在，创建
			_, err = db.Exec(idx.sql)
			if err != nil {
				fmt.Printf("[失败] 创建索引 %s.%s: %v\n", idx.table, idx.name, err)
			} else {
				fmt.Printf("[成功] 创建索引 %s.%s\n", idx.table, idx.name)
			}
		} else if err != nil {
			fmt.Printf("[错误] 检查索引 %s.%s: %v\n", idx.table, idx.name, err)
		} else {
			fmt.Printf("[跳过] 索引 %s.%s 已存在\n", idx.table, idx.name)
		}
	}

	fmt.Println("\n数据表迁移完成!")
}
