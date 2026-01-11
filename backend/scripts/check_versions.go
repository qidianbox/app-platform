package main

import (
    "fmt"
    "log"
    "os"
    "gopkg.in/yaml.v3"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type Config struct {
    Database struct {
        DSN string `yaml:"dsn"`
    } `yaml:"database"`
}

type FeatureVersion struct {
    ID           uint   `gorm:"primaryKey"`
    CollectionID uint   `gorm:"column:collection_id"`
    Version      string `gorm:"column:version"`
    Status       string `gorm:"column:status"`
    Changelog    string `gorm:"column:changelog"`
}

func main() {
    data, err := os.ReadFile("config/config.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        log.Fatal(err)
    }
    
    db, err := gorm.Open(mysql.Open(config.Database.DSN), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }
    
    var versions []FeatureVersion
    result := db.Table("feature_versions").Find(&versions)
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    
    fmt.Printf("Found %d versions:\n", len(versions))
    for _, v := range versions {
        fmt.Printf("ID: %d, CollectionID: %d, Version: %s, Status: %s, Changelog: %s\n", 
            v.ID, v.CollectionID, v.Version, v.Status, v.Changelog)
    }
}
