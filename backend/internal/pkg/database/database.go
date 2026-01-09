package database

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"app-platform-backend/internal/config"

	"github.com/go-sql-driver/mysql"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB(cfg *config.DatabaseConfig) error {
	var dsn string
	
	// 检查是否有DATABASE_URL环境变量（Manus平台数据库）
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		// 解析DATABASE_URL格式: mysql://user:pass@host:port/dbname?ssl=...
		dsn = parseDatabaseURL(databaseURL)
	} else {
		// 使用配置文件中的数据库配置
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	}

	var err error
	db, err = gorm.Open(gormMysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	return nil
}

// parseDatabaseURL 解析DATABASE_URL格式并返回DSN
func parseDatabaseURL(url string) string {
	// 注册TLS配置
	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: "",
	})
	
	// 移除mysql://前缀
	url = strings.TrimPrefix(url, "mysql://")
	
	// 分离用户信息和主机信息
	// 格式: user:pass@host:port/dbname?ssl=...
	atIndex := strings.Index(url, "@")
	if atIndex == -1 {
		return url
	}
	
	userPass := url[:atIndex]
	rest := url[atIndex+1:]
	
	// 分离主机和数据库名
	slashIndex := strings.Index(rest, "/")
	if slashIndex == -1 {
		return url
	}
	
	hostPort := rest[:slashIndex]
	dbAndParams := rest[slashIndex+1:]
	
	// 分离数据库名和参数
	questionIndex := strings.Index(dbAndParams, "?")
	var dbName string
	if questionIndex == -1 {
		dbName = dbAndParams
	} else {
		dbName = dbAndParams[:questionIndex]
	}
	
	// 构建DSN
	dsn := fmt.Sprintf("%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=tidb",
		userPass, hostPort, dbName)
	
	return dsn
}

func GetDB() *gorm.DB {
	return db
}

func Close() {
	if db != nil {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}
