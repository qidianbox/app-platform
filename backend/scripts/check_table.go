package main

import (
    "fmt"
    "app-platform-backend/internal/config"
    "app-platform-backend/internal/pkg/database"
)

func main() {
    cfg, _ := config.LoadConfig("config/config.yaml")
    database.InitDB(&cfg.Database)
    db := database.GetDB()
    
    var columns []struct {
        Field string
        Type  string
    }
    
    db.Raw("SHOW COLUMNS FROM data_collections").Scan(&columns)
    
    fmt.Println("data_collections columns:")
    for _, col := range columns {
        fmt.Printf("  %s: %s\n", col.Field, col.Type)
    }
}
