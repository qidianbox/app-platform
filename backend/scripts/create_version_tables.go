package main

import (
    "log"
    "os"
    "app-platform-backend/internal/config"
    "app-platform-backend/internal/pkg/database"
)

func main() {
    // 加载配置
    cfg, err := config.LoadConfig("configs/config.yaml")
    if err != nil {
        log.Printf("Warning: Could not load config: %v", err)
    }

    // 初始化数据库
    if err := database.InitDB(&cfg.Database); err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    db := database.GetDB()
    sqlDB, _ := db.DB()

    // 创建功能版本表
    _, err = sqlDB.Exec(`
        CREATE TABLE IF NOT EXISTS feature_versions (
            id BIGINT PRIMARY KEY AUTO_INCREMENT,
            app_id BIGINT NOT NULL COMMENT '应用ID',
            collection_id BIGINT NOT NULL COMMENT '数据模型ID',
            version VARCHAR(20) NOT NULL COMMENT '版本号',
            version_num INT NOT NULL COMMENT '版本序号',
            schema_snapshot JSON NOT NULL COMMENT '字段结构快照',
            status VARCHAR(20) DEFAULT 'draft' COMMENT '版本状态',
            changelog TEXT COMMENT '变更说明',
            created_by VARCHAR(100) COMMENT '创建人',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            published_at TIMESTAMP NULL COMMENT '发布时间',
            INDEX idx_collection_id (collection_id),
            INDEX idx_app_id (app_id),
            INDEX idx_status (status)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
    `)
    if err != nil {
        log.Printf("Create feature_versions table result: %v", err)
    } else {
        log.Println("Table feature_versions created successfully")
    }

    // 创建模块版本表
    _, err = sqlDB.Exec(`
        CREATE TABLE IF NOT EXISTS module_versions (
            id BIGINT PRIMARY KEY AUTO_INCREMENT,
            app_id BIGINT NOT NULL COMMENT '应用ID',
            module_code VARCHAR(50) NOT NULL COMMENT '模块代码',
            version VARCHAR(20) NOT NULL COMMENT '版本号',
            version_num INT NOT NULL COMMENT '版本序号',
            config_snapshot JSON NOT NULL COMMENT '配置快照',
            status VARCHAR(20) DEFAULT 'draft' COMMENT '版本状态',
            environment VARCHAR(20) DEFAULT 'dev' COMMENT '环境',
            changelog TEXT COMMENT '变更说明',
            created_by VARCHAR(100),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            published_at TIMESTAMP NULL,
            INDEX idx_app_module (app_id, module_code),
            INDEX idx_environment (environment),
            INDEX idx_status (status)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
    `)
    if err != nil {
        log.Printf("Create module_versions table result: %v", err)
    } else {
        log.Println("Table module_versions created successfully")
    }
    
    os.Exit(0)
}
