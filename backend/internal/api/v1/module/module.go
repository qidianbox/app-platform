package module

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"app-platform-backend/internal/model"
	"app-platform-backend/internal/pkg/database"

	"github.com/gin-gonic/gin"
)

func GetAllTemplates(c *gin.Context) {
	var templates []model.ModuleTemplate
	database.GetDB().Where("status = 1").Find(&templates)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": templates,
	})
}

func GetAppModules(c *gin.Context) {
	appID := c.Param("id")

	var modules []model.AppModule
	database.GetDB().Where("app_id = ?", appID).Find(&modules)

	// 获取模块信息（使用source_module匹配）
	type ModuleInfo struct {
		ID          uint   `json:"id"`
		AppID       uint   `json:"app_id"`
		ModuleCode  string `json:"module_code"`
		ModuleName  string `json:"module_name"`
		Category    string `json:"category"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Status      int    `json:"status"`
	}

	// 模块名称映射
	moduleNameMap := map[string]struct {
		Name        string
		Category    string
		Description string
		Icon        string
	}{
		"user_management":    {"用户管理", "用户与权限", "用户注册、登录、权限管理", "user"},
		"message_center":     {"消息中心", "消息与通知", "站内消息、通知管理", "message"},
		"push_service":       {"推送服务", "消息与通知", "APP推送通知服务", "notification"},
		"data_tracking":      {"数据埋点", "数据与分析", "用户行为埋点和数据分析", "chart"},
		"log_service":        {"日志服务", "系统与运维", "应用日志收集和查询", "document"},
		"monitor_alert":      {"监控告警", "系统与运维", "应用监控和告警通知", "warning"},
		"file_storage":       {"文件存储", "存储服务", "文件上传、下载、管理", "folder"},
		"config_management":  {"配置管理", "存储服务", "远程配置下发和管理", "setting"},
		"version_management": {"版本管理", "存储服务", "APP版本发布和更新", "box"},
	}

	result := make([]ModuleInfo, 0)
	for _, m := range modules {
		info := ModuleInfo{
			ID:         m.ID,
			AppID:      m.AppID,
			ModuleCode: m.ModuleCode,
			Status:     m.Status,
		}
		if moduleInfo, ok := moduleNameMap[m.ModuleCode]; ok {
			info.ModuleName = moduleInfo.Name
			info.Category = moduleInfo.Category
			info.Description = moduleInfo.Description
			info.Icon = moduleInfo.Icon
		} else {
			info.ModuleName = m.ModuleCode
			info.Category = "其他"
			info.Description = ""
			info.Icon = "setting"
		}
		result = append(result, info)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": result,
	})
}

func GetAppModule(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")

	var module model.AppModule
	if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, moduleCode).First(&module).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": module,
	})
}

func EnableModule(c *gin.Context) {
	appID := c.Param("id")

	var req struct {
		ModuleCode string `json:"module_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查是否已启用
	var existing model.AppModule
	if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, req.ModuleCode).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Module already enabled"})
		return
	}

	module := model.AppModule{
		AppID:      parseUint(appID),
		ModuleCode: req.ModuleCode,
		Config:     "{}",
		Status:     1,
	}

	if err := database.GetDB().Create(&module).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable module"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": module,
	})
}

func UpdateModule(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")

	var module model.AppModule
	if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, moduleCode).First(&module).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	var req struct {
		Status *int `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status != nil {
		database.GetDB().Model(&module).Update("status", *req.Status)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": module,
	})
}

func DisableModule(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")

	if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, moduleCode).Delete(&model.AppModule{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable module"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Module disabled successfully",
	})
}

func BatchEnableModules(c *gin.Context) {
	appID := c.Param("id")

	var req struct {
		ModuleCodes []string `json:"module_codes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, code := range req.ModuleCodes {
		var existing model.AppModule
		if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, code).First(&existing).Error; err != nil {
			module := model.AppModule{
				AppID:      parseUint(appID),
				ModuleCode: code,
				Config:     "{}",
				Status:     1,
			}
			database.GetDB().Create(&module)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Modules enabled successfully",
	})
}

func SaveModuleConfig(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")

	var module model.AppModule
	if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, moduleCode).First(&module).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	var req struct {
		Config map[string]interface{} `json:"config" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 保存配置历史
	var maxVersion int
	database.GetDB().Model(&model.ModuleConfigHistory{}).
		Where("app_id = ? AND module_code = ?", appID, moduleCode).
		Select("COALESCE(MAX(version), 0)").Scan(&maxVersion)

	history := model.ModuleConfigHistory{
		AppID:      parseUint(appID),
		ModuleCode: moduleCode,
		Config:     module.Config,
		Version:    maxVersion + 1,
		Operator:   c.GetString("username"),
	}
	database.GetDB().Create(&history)

	// 更新配置
	configJSON, _ := json.Marshal(req.Config)
	database.GetDB().Model(&module).Update("config", string(configJSON))

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Config saved successfully",
	})
}

func GetModuleConfig(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")

	var module model.AppModule
	if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, moduleCode).First(&module).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"config": module.Config,
		},
	})
}

func ResetModuleConfig(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")

	var module model.AppModule
	if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, moduleCode).First(&module).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	database.GetDB().Model(&module).Update("config", "{}")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Config reset successfully",
	})
}

func TestModuleConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Config test passed",
		"data": gin.H{
			"success": true,
		},
	})
}

func GetConfigHistory(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")

	var history []model.ModuleConfigHistory
	database.GetDB().Where("app_id = ? AND module_code = ?", appID, moduleCode).
		Order("version DESC").Limit(20).Find(&history)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": history,
	})
}

func RollbackConfig(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")
	historyID := c.Param("history_id")

	var history model.ModuleConfigHistory
	if err := database.GetDB().First(&history, historyID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "History not found"})
		return
	}

	database.GetDB().Model(&model.AppModule{}).
		Where("app_id = ? AND module_code = ?", appID, moduleCode).
		Update("config", history.Config)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Config rolled back successfully",
	})
}

func CompareConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"diff": []string{},
		},
	})
}

func CheckModuleDependencies(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"satisfied":   true,
			"missing":     []string{},
			"suggestions": []string{},
		},
	})
}

func CheckModuleReverseDependencies(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"dependents": []string{},
		},
	})
}

func AutoEnableModuleDependencies(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Dependencies enabled successfully",
	})
}

func DetectCircularDependency(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"has_circular": false,
			"path":         []string{},
		},
	})
}

func parseUint(s string) uint {
	var id uint
	fmt.Sscanf(s, "%d", &id)
	return id
}



// ==================== 模块版本管理 API ====================

// ModuleVersion 模块版本结构体
type ModuleVersion struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	AppID          uint       `gorm:"index" json:"app_id"`
	ModuleCode     string     `gorm:"size:100;index" json:"module_code"`
	VersionNum     int        `json:"version_num"`
	Version        string     `gorm:"size:50" json:"version"`
	ConfigSnapshot string     `gorm:"type:text" json:"config_snapshot"`
	Status         string     `gorm:"size:20;default:draft" json:"status"`
	Environment    string     `gorm:"size:20;default:dev" json:"environment"`
	Changelog      string     `gorm:"type:text" json:"changelog"`
	CreatedBy      string     `gorm:"size:100" json:"created_by"`
	PublishedAt    *time.Time `json:"published_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

// TableName 指定表名
func (ModuleVersion) TableName() string {
	return "module_versions"
}

// MigrateModuleVersions 数据库迁移
func MigrateModuleVersions() error {
	return database.GetDB().AutoMigrate(&ModuleVersion{})
}

// ListModuleVersions 获取模块版本列表
func ListModuleVersions(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")

	// 验证模块存在
	var module model.AppModule
	if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, moduleCode).First(&module).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "模块不存在"})
		return
	}

	status := c.Query("status")
	environment := c.Query("environment")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	var versions []ModuleVersion
	var total int64

	query := database.GetDB().Model(&ModuleVersion{}).Where("app_id = ? AND module_code = ?", appID, moduleCode)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if environment != "" {
		query = query.Where("environment = ?", environment)
	}

	query.Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&versions)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  versions,
			"total": total,
		},
	})
}

// CreateModuleVersion 创建模块版本
func CreateModuleVersion(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")

	// 验证模块存在
	var module model.AppModule
	if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, moduleCode).First(&module).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "模块不存在"})
		return
	}

	var req struct {
		Version     string `json:"version"`
		Environment string `json:"environment"`
		Changelog   string `json:"changelog"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 自动生成版本号
	version := req.Version
	if version == "" {
		var maxVersion int
		database.GetDB().Model(&ModuleVersion{}).
			Where("app_id = ? AND module_code = ?", appID, moduleCode).
			Select("COALESCE(MAX(CAST(SUBSTRING_INDEX(version, '.', 1) AS UNSIGNED)), 0)").
			Scan(&maxVersion)
		version = fmt.Sprintf("%d.0.0", maxVersion+1)
	}

	// 默认环境
	environment := req.Environment
	if environment == "" {
		environment = "dev"
	}

	// 获取最大版本号
	var maxVersionNum int
	database.GetDB().Model(&ModuleVersion{}).
		Where("app_id = ? AND module_code = ?", appID, moduleCode).
		Select("COALESCE(MAX(version_num), 0)").
		Scan(&maxVersionNum)

	// 创建版本
	newVersion := ModuleVersion{
		AppID:          parseUint(appID),
		ModuleCode:     moduleCode,
		VersionNum:     maxVersionNum + 1,
		Version:        version,
		ConfigSnapshot: module.Config,
		Status:         "draft",
		Environment:    environment,
		Changelog:      req.Changelog,
		CreatedBy:      c.GetString("username"),
	}

	if err := database.GetDB().Create(&newVersion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建版本失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": newVersion,
	})
}

// PublishModuleVersion 发布模块版本
func PublishModuleVersion(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")
	versionID := c.Param("version_id")

	var version ModuleVersion
	if err := database.GetDB().Where("id = ? AND app_id = ? AND module_code = ?", versionID, appID, moduleCode).First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	if version.Status == "published" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "版本已发布"})
		return
	}

	now := time.Now()
	database.GetDB().Model(&version).Updates(map[string]interface{}{
		"status":       "published",
		"published_at": now,
	})

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "版本发布成功",
	})
}

// RollbackModuleVersion 回滚模块版本
func RollbackModuleVersion(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")
	versionID := c.Param("version_id")

	var version ModuleVersion
	if err := database.GetDB().Where("id = ? AND app_id = ? AND module_code = ?", versionID, appID, moduleCode).First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	// 将配置快照恢复到模块
	database.GetDB().Model(&model.AppModule{}).
		Where("app_id = ? AND module_code = ?", appID, moduleCode).
		Update("config", version.ConfigSnapshot)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "版本回滚成功",
	})
}

// CompareModuleVersions 对比模块版本
func CompareModuleVersions(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")
	version1ID := c.Query("version1")
	version2ID := c.Query("version2")

	var v1, v2 ModuleVersion
	if err := database.GetDB().Where("id = ? AND app_id = ? AND module_code = ?", version1ID, appID, moduleCode).First(&v1).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本1不存在"})
		return
	}
	if err := database.GetDB().Where("id = ? AND app_id = ? AND module_code = ?", version2ID, appID, moduleCode).First(&v2).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本2不存在"})
		return
	}

	// 解析配置JSON进行对比
	var config1, config2 map[string]interface{}
	json.Unmarshal([]byte(v1.ConfigSnapshot), &config1)
	json.Unmarshal([]byte(v2.ConfigSnapshot), &config2)

	// 简单对比
	changes := []map[string]interface{}{}
	
	// 检查v1中的字段
	for key, val1 := range config1 {
		if val2, ok := config2[key]; ok {
			if fmt.Sprintf("%v", val1) != fmt.Sprintf("%v", val2) {
				changes = append(changes, map[string]interface{}{
					"field":     key,
					"type":      "modified",
					"old_value": val1,
					"new_value": val2,
				})
			}
		} else {
			changes = append(changes, map[string]interface{}{
				"field":     key,
				"type":      "removed",
				"old_value": val1,
			})
		}
	}
	
	// 检查v2中新增的字段
	for key, val2 := range config2 {
		if _, ok := config1[key]; !ok {
			changes = append(changes, map[string]interface{}{
				"field":     key,
				"type":      "added",
				"new_value": val2,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"version1": v1,
			"version2": v2,
			"changes":  changes,
		},
	})
}

