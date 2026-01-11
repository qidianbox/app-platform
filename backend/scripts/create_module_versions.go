package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 直接使用配置文件中的数据库连接信息
	dsn := "app_platform:App@Platform123@tcp(rm-bp13s51058fu3r061.mysql.rds.aliyuncs.com:3306)/app_platform?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 创建模块版本表
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS module_versions (
		id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		app_id BIGINT UNSIGNED NOT NULL COMMENT 'APP ID',
		module_code VARCHAR(100) NOT NULL COMMENT '模块代码',
		version VARCHAR(50) NOT NULL COMMENT '版本号',
		config_snapshot TEXT COMMENT '配置快照',
		status VARCHAR(20) DEFAULT 'draft' COMMENT '状态: draft/published/deprecated',
		environment VARCHAR(20) DEFAULT 'dev' COMMENT '环境: dev/test/prod',
		changelog TEXT COMMENT '变更说明',
		created_by VARCHAR(100) COMMENT '创建人',
		published_at DATETIME COMMENT '发布时间',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_app_module (app_id, module_code),
		INDEX idx_status (status),
		INDEX idx_environment (environment)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='模块版本表';
	`

	if err := db.Exec(createTableSQL).Error; err != nil {
		fmt.Printf("创建表失败: %v\n", err)
	} else {
		fmt.Println("模块版本表创建成功")
	}
}
