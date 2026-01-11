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

    // 添加is_visible字段
    _, err = sqlDB.Exec("ALTER TABLE data_collections ADD COLUMN is_visible TINYINT(1) DEFAULT 0 AFTER is_generated")
    if err != nil {
        log.Printf("Result: %v", err)
    } else {
        log.Println("Column is_visible added successfully")
    }
    
    os.Exit(0)
}
