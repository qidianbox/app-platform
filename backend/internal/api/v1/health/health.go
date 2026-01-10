package health

import (
	"net/http"
	"runtime"
	"time"

	"app-platform-backend/internal/pkg/database"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// HealthStatus 健康状态
type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp int64                  `json:"timestamp"`
	Uptime    float64                `json:"uptime"`
	Version   string                 `json:"version"`
	Checks    map[string]CheckResult `json:"checks"`
	System    SystemInfo             `json:"system"`
}

// CheckResult 检查结果
type CheckResult struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Latency int64  `json:"latency_ms,omitempty"`
}

// SystemInfo 系统信息
type SystemInfo struct {
	GoVersion    string `json:"go_version"`
	NumGoroutine int    `json:"num_goroutine"`
	NumCPU       int    `json:"num_cpu"`
	MemAlloc     uint64 `json:"mem_alloc_mb"`
	MemSys       uint64 `json:"mem_sys_mb"`
}

// Check 健康检查
func Check(c *gin.Context) {
	checks := make(map[string]CheckResult)

	// 数据库检查
	checks["database"] = checkDatabase()

	// 计算整体状态
	overallStatus := "healthy"
	for _, check := range checks {
		if check.Status != "healthy" {
			overallStatus = "unhealthy"
			break
		}
	}

	// 获取系统信息
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	status := HealthStatus{
		Status:    overallStatus,
		Timestamp: time.Now().Unix(),
		Uptime:    time.Since(startTime).Seconds(),
		Version:   "1.0.0",
		Checks:    checks,
		System: SystemInfo{
			GoVersion:    runtime.Version(),
			NumGoroutine: runtime.NumGoroutine(),
			NumCPU:       runtime.NumCPU(),
			MemAlloc:     memStats.Alloc / 1024 / 1024,
			MemSys:       memStats.Sys / 1024 / 1024,
		},
	}

	httpStatus := http.StatusOK
	if overallStatus != "healthy" {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, gin.H{
		"code": 0,
		"data": status,
	})
}

// checkDatabase 检查数据库连接
func checkDatabase() CheckResult {
	start := time.Now()

	db := database.GetDB()
	if db == nil {
		return CheckResult{
			Status:  "unhealthy",
			Message: "Database not initialized",
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return CheckResult{
			Status:  "unhealthy",
			Message: "Failed to get database connection: " + err.Error(),
		}
	}

	if err := sqlDB.Ping(); err != nil {
		return CheckResult{
			Status:  "unhealthy",
			Message: "Database ping failed: " + err.Error(),
			Latency: time.Since(start).Milliseconds(),
		}
	}

	return CheckResult{
		Status:  "healthy",
		Latency: time.Since(start).Milliseconds(),
	}
}

// Liveness 存活探针
func Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

// Readiness 就绪探针
func Readiness(c *gin.Context) {
	// 检查数据库连接
	dbCheck := checkDatabase()
	if dbCheck.Status != "healthy" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not ready",
			"message": dbCheck.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// Metrics 简单的指标端点
func Metrics(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"uptime_seconds":     time.Since(startTime).Seconds(),
			"goroutines":         runtime.NumGoroutine(),
			"memory_alloc_bytes": memStats.Alloc,
			"memory_sys_bytes":   memStats.Sys,
			"gc_runs":            memStats.NumGC,
			"gc_pause_total_ns":  memStats.PauseTotalNs,
		},
	})
}
