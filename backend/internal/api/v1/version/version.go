package version

import (
	"app-platform-backend/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(database *gorm.DB) {
	db = database
}

// List 版本列表
func List(c *gin.Context) {
	appID := c.Query("app_id")
	platform := c.Query("platform")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	query := db.Model(&model.Version{}).Where("app_id = ?", appID)

	if platform != "" {
		// 暂时忽略平台筛选，因为模型中没有platform字段
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var versions []model.Version
	offset := (page - 1) * size
	query.Offset(offset).Limit(size).Order("created_at DESC").Find(&versions)

	// 转换为前端需要的格式
	var result []map[string]interface{}
	for _, v := range versions {
		result = append(result, map[string]interface{}{
			"id":           v.ID,
			"app_id":       v.AppID,
			"version":      v.VersionName,
			"version_code": v.VersionCode,
			"platform":     "android", // 默认平台
			"description":  v.Description,
			"download_url": v.DownloadURL,
			"force_update": v.IsForceUpdate == 1,
			"status":       v.Status,
			"published_at": v.PublishedAt,
			"created_at":   v.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  result,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// Create 创建版本
func Create(c *gin.Context) {
	var req struct {
		AppID       uint   `json:"app_id" binding:"required"`
		Version     string `json:"version" binding:"required"`
		Platform    string `json:"platform"`
		Description string `json:"description"`
		DownloadURL string `json:"download_url"`
		ForceUpdate bool   `json:"force_update"`
		GrayRelease bool   `json:"gray_release"`
		GrayPercent int    `json:"gray_percent"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 获取最大版本号
	var maxVersionCode int
	db.Model(&model.Version{}).Where("app_id = ?", req.AppID).Select("COALESCE(MAX(version_code), 0)").Scan(&maxVersionCode)

	forceUpdate := 0
	if req.ForceUpdate {
		forceUpdate = 1
	}

	version := model.Version{
		AppID:         req.AppID,
		VersionName:   req.Version,
		VersionCode:   maxVersionCode + 1,
		Description:   req.Description,
		DownloadURL:   req.DownloadURL,
		IsForceUpdate: forceUpdate,
		Status:        "draft",
	}

	if err := db.Create(&version).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create version"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Version created successfully",
		"data":    version,
	})
}

// Update 更新版本
func Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid version ID"})
		return
	}

	var req struct {
		Version     string `json:"version"`
		Description string `json:"description"`
		DownloadURL string `json:"download_url"`
		ForceUpdate bool   `json:"force_update"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Version != "" {
		updates["version_name"] = req.Version
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.DownloadURL != "" {
		updates["download_url"] = req.DownloadURL
	}
	updates["is_force_update"] = 0
	if req.ForceUpdate {
		updates["is_force_update"] = 1
	}

	if err := db.Model(&model.Version{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update version"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Version updated successfully",
	})
}

// Publish 发布版本
func Publish(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid version ID"})
		return
	}

	now := time.Now()
	if err := db.Model(&model.Version{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       "published",
		"published_at": now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to publish version"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Version published successfully",
	})
}

// Offline 下线版本
func Offline(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid version ID"})
		return
	}

	if err := db.Model(&model.Version{}).Where("id = ?", id).Update("status", "offline").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to offline version"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Version offline successfully",
	})
}

// Delete 删除版本
func Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid version ID"})
		return
	}

	if err := db.Delete(&model.Version{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to delete version"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Version deleted successfully",
	})
}

// CheckUpdate 检查更新
func CheckUpdate(c *gin.Context) {
	appID := c.Query("app_id")
	currentVersion := c.Query("version")

	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	// 获取最新发布的版本
	var latestVersion model.Version
	err := db.Where("app_id = ? AND status = ?", appID, "published").
		Order("version_code DESC").
		First(&latestVersion).Error

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"has_update": false,
			},
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to check update"})
		return
	}

	hasUpdate := latestVersion.VersionName != currentVersion

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"has_update":   hasUpdate,
			"version":      latestVersion.VersionName,
			"version_code": latestVersion.VersionCode,
			"description":  latestVersion.Description,
			"download_url": latestVersion.DownloadURL,
			"force_update": latestVersion.IsForceUpdate == 1,
		},
	})
}

// Stats 版本统计
func Stats(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	var total, published, draft, offline int64
	db.Model(&model.Version{}).Where("app_id = ?", appID).Count(&total)
	db.Model(&model.Version{}).Where("app_id = ? AND status = ?", appID, "published").Count(&published)
	db.Model(&model.Version{}).Where("app_id = ? AND status = ?", appID, "draft").Count(&draft)
	db.Model(&model.Version{}).Where("app_id = ? AND status = ?", appID, "offline").Count(&offline)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"total":     total,
			"published": published,
			"draft":     draft,
			"offline":   offline,
		},
	})
}
