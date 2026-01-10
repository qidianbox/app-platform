package file

import (
	"app-platform-backend/internal/model"
	"app-platform-backend/internal/response"
	"app-platform-backend/internal/validator"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB
var uploadDir = "/tmp/uploads"

// 允许的文件类型
var allowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
	"application/pdf": true,
	"application/zip": true,
	"text/plain":      true,
	"text/csv":        true,
	"application/json": true,
	"application/vnd.ms-excel": true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
}

// 最大文件大小 (50MB)
const maxFileSize = 50 * 1024 * 1024

func InitDB(database *gorm.DB) {
	db = database
	// 确保上传目录存在
	os.MkdirAll(uploadDir, 0755)
}

// Upload 上传文件
func Upload(c *gin.Context) {
	appIDStr := c.PostForm("app_id")
	if appIDStr == "" {
		response.ParamError(c, "app_id 不能为空")
		return
	}

	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		response.ParamError(c, "无效的 app_id")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.ParamError(c, "请选择要上传的文件")
		return
	}
	defer file.Close()

	// 验证文件大小
	if header.Size > maxFileSize {
		response.ParamError(c, fmt.Sprintf("文件大小不能超过 %dMB", maxFileSize/1024/1024))
		return
	}

	// 获取MIME类型
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// 验证文件类型（可选，根据需求开启）
	// if !allowedMimeTypes[mimeType] {
	// 	response.ParamError(c, "不支持的文件类型")
	// 	return
	// }

	// 生成唯一文件名
	ext := filepath.Ext(header.Filename)
	hash := md5.New()
	io.Copy(hash, file)
	file.Seek(0, 0)

	timestamp := time.Now().UnixNano()
	newFilename := fmt.Sprintf("%x_%d%s", hash.Sum(nil), timestamp, ext)

	// 按APP和日期组织目录
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(uploadDir, fmt.Sprintf("%d", appID), dateDir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		response.ServerError(c, "创建目录失败")
		return
	}

	filePath := filepath.Join(fullDir, newFilename)

	// 保存文件
	out, err := os.Create(filePath)
	if err != nil {
		response.ServerError(c, "保存文件失败")
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		response.ServerError(c, "写入文件失败")
		return
	}

	// 保存到数据库
	fileRecord := model.File{
		AppID:    uint(appID),
		Filename: header.Filename,
		FilePath: filePath,
		FileSize: header.Size,
		MimeType: mimeType,
	}

	if err := db.Create(&fileRecord).Error; err != nil {
		// 删除已上传的文件
		os.Remove(filePath)
		response.DBError(c, err)
		return
	}

	// 生成访问URL
	fileURL := fmt.Sprintf("/api/v1/files/download/%d", fileRecord.ID)

	response.SuccessWithMessage(c, gin.H{
		"id":        fileRecord.ID,
		"filename":  header.Filename,
		"size":      header.Size,
		"mime_type": mimeType,
		"url":       fileURL,
	}, "文件上传成功")
}

// List 文件列表
func List(c *gin.Context) {
	appID := c.Query("app_id")
	mimeType := c.Query("mime_type")
	page, size := validator.ParsePagination(c.DefaultQuery("page", "1"), c.DefaultQuery("size", "20"))

	if appID == "" {
		response.ParamError(c, "app_id 不能为空")
		return
	}

	// 验证分页参数
	page, size = validator.ValidatePagination(page, size)

	query := db.Model(&model.File{}).Where("app_id = ?", appID)

	if mimeType != "" {
		query = query.Where("mime_type LIKE ?", mimeType+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.DBError(c, err)
		return
	}

	var files []model.File
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&files).Error; err != nil {
		response.DBError(c, err)
		return
	}

	// 添加URL
	var result []gin.H
	for _, f := range files {
		result = append(result, gin.H{
			"id":         f.ID,
			"filename":   f.Filename,
			"file_size":  f.FileSize,
			"mime_type":  f.MimeType,
			"url":        fmt.Sprintf("/api/v1/files/download/%d", f.ID),
			"created_at": f.CreatedAt,
		})
	}

	response.PageSuccess(c, result, total, page, size)
}

// Detail 文件详情
func Detail(c *gin.Context) {
	id := c.Param("id")

	// 验证ID
	if _, err := validator.ValidateID(id); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	var file model.File
	if err := db.First(&file, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "文件不存在")
			return
		}
		response.DBError(c, err)
		return
	}

	response.Success(c, gin.H{
		"id":         file.ID,
		"filename":   file.Filename,
		"file_size":  file.FileSize,
		"mime_type":  file.MimeType,
		"url":        fmt.Sprintf("/api/v1/files/download/%d", file.ID),
		"created_at": file.CreatedAt,
	})
}

// Download 下载文件
func Download(c *gin.Context) {
	id := c.Param("id")

	// 验证ID
	if _, err := validator.ValidateID(id); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	var file model.File
	if err := db.First(&file, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "文件不存在")
			return
		}
		response.DBError(c, err)
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		response.NotFound(c, "文件已被删除")
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
	c.Header("Content-Type", file.MimeType)
	c.File(file.FilePath)
}

// Delete 删除文件
func Delete(c *gin.Context) {
	id := c.Param("id")

	// 验证ID
	if _, err := validator.ValidateID(id); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	var file model.File
	if err := db.First(&file, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "文件不存在")
			return
		}
		response.DBError(c, err)
		return
	}

	// 删除物理文件
	os.Remove(file.FilePath)

	// 删除数据库记录
	if err := db.Delete(&file).Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, nil, "文件删除成功")
}

// BatchDelete 批量删除文件
func BatchDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "参数错误: "+err.Error())
		return
	}

	if len(req.IDs) == 0 {
		response.ParamError(c, "请选择要删除的文件")
		return
	}

	if len(req.IDs) > 100 {
		response.ParamError(c, "单次最多删除100个文件")
		return
	}

	var files []model.File
	if err := db.Find(&files, req.IDs).Error; err != nil {
		response.DBError(c, err)
		return
	}

	// 删除物理文件
	for _, file := range files {
		os.Remove(file.FilePath)
	}

	// 删除数据库记录
	result := db.Delete(&model.File{}, req.IDs)
	if result.Error != nil {
		response.DBError(c, result.Error)
		return
	}

	response.SuccessWithMessage(c, gin.H{
		"affected": result.RowsAffected,
	}, "文件批量删除成功")
}

// Stats 文件统计
func Stats(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		response.ParamError(c, "app_id 不能为空")
		return
	}

	var total int64
	var totalSize int64
	var todayCount int64

	db.Model(&model.File{}).Where("app_id = ?", appID).Count(&total)
	db.Model(&model.File{}).Where("app_id = ?", appID).Select("COALESCE(SUM(file_size), 0)").Scan(&totalSize)

	today := time.Now().Format("2006-01-02")
	db.Model(&model.File{}).Where("app_id = ? AND DATE(created_at) = ?", appID, today).Count(&todayCount)

	// 按类型统计
	var typeStats []struct {
		MimeType string `json:"mime_type"`
		Count    int64  `json:"count"`
		Size     int64  `json:"size"`
	}
	db.Model(&model.File{}).
		Select("SUBSTRING_INDEX(mime_type, '/', 1) as mime_type, COUNT(*) as count, SUM(file_size) as size").
		Where("app_id = ?", appID).
		Group("SUBSTRING_INDEX(mime_type, '/', 1)").
		Scan(&typeStats)

	response.Success(c, gin.H{
		"total":       total,
		"total_size":  totalSize,
		"today_count": todayCount,
		"type_stats":  typeStats,
	})
}
