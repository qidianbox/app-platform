package push

import (
	"app-platform-backend/internal/model"
	"app-platform-backend/internal/response"
	"app-platform-backend/internal/validator"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(database *gorm.DB) {
	db = database
}

// List 推送列表
func List(c *gin.Context) {
	appID := c.Query("app_id")
	status := c.Query("status")
	page, size := validator.ParsePagination(c.DefaultQuery("page", "1"), c.DefaultQuery("size", "20"))

	if appID == "" {
		response.ParamError(c, "app_id 不能为空")
		return
	}

	// 验证分页参数
	page, size = validator.ValidatePagination(page, size)

	query := db.Model(&model.PushRecord{}).Where("app_id = ?", appID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.DBError(c, err)
		return
	}

	var records []model.PushRecord
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&records).Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.PageSuccess(c, records, total, page, size)
}

// Create 创建推送任务
func Create(c *gin.Context) {
	var req struct {
		AppID       uint     `json:"app_id" binding:"required"`
		Title       string   `json:"title" binding:"required"`
		Content     string   `json:"content" binding:"required"`
		TargetType  string   `json:"target_type"`
		TargetIDs   []string `json:"target_ids"`
		ScheduledAt string   `json:"scheduled_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "参数错误: "+err.Error())
		return
	}

	// 验证标题长度
	if len(req.Title) < 1 || len(req.Title) > 100 {
		response.ParamError(c, "标题长度应在1-100个字符之间")
		return
	}

	// 验证内容长度
	if len(req.Content) < 1 || len(req.Content) > 1000 {
		response.ParamError(c, "内容长度应在1-1000个字符之间")
		return
	}

	if req.TargetType == "" {
		req.TargetType = "all"
	}

	// 验证目标类型
	validTargetTypes := map[string]bool{"all": true, "user": true, "tag": true, "segment": true}
	if !validTargetTypes[req.TargetType] {
		response.ParamError(c, "无效的目标类型，请使用: all, user, tag, segment")
		return
	}

	record := model.PushRecord{
		AppID:      req.AppID,
		Title:      req.Title,
		Content:    req.Content,
		TargetType: req.TargetType,
		TargetIDs:  strings.Join(req.TargetIDs, ","),
		Status:     "pending",
	}

	if req.ScheduledAt != "" {
		scheduledTime, err := time.Parse("2006-01-02 15:04:05", req.ScheduledAt)
		if err != nil {
			response.ParamError(c, "计划发送时间格式错误，请使用: 2006-01-02 15:04:05")
			return
		}
		if scheduledTime.Before(time.Now()) {
			response.ParamError(c, "计划发送时间不能早于当前时间")
			return
		}
		record.ScheduledAt = &scheduledTime
	}

	if err := db.Create(&record).Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, record, "推送任务创建成功")
}

// Detail 推送详情
func Detail(c *gin.Context) {
	id := c.Param("id")

	// 验证ID
	if _, err := validator.ValidateID(id); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	var record model.PushRecord
	if err := db.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "推送记录不存在")
			return
		}
		response.DBError(c, err)
		return
	}

	response.Success(c, record)
}

// Send 立即发送推送
func Send(c *gin.Context) {
	id := c.Param("id")

	// 验证ID
	if _, err := validator.ValidateID(id); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	var record model.PushRecord
	if err := db.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "推送记录不存在")
			return
		}
		response.DBError(c, err)
		return
	}

	if record.Status != "pending" {
		response.ParamError(c, "只有待发送状态的推送可以发送")
		return
	}

	// 模拟发送推送
	now := time.Now()
	sentCount := 100
	successCount := 95
	failedCount := 5

	if err := db.Model(&record).Updates(map[string]interface{}{
		"status":        "sent",
		"sent_at":       now,
		"sent_count":    sentCount,
		"success_count": successCount,
		"failed_count":  failedCount,
	}).Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, gin.H{
		"sent_count":    sentCount,
		"success_count": successCount,
		"failed_count":  failedCount,
	}, "推送发送成功")
}

// Cancel 取消推送任务
func Cancel(c *gin.Context) {
	id := c.Param("id")

	// 验证ID
	if _, err := validator.ValidateID(id); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	var record model.PushRecord
	if err := db.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "推送记录不存在")
			return
		}
		response.DBError(c, err)
		return
	}

	if record.Status != "pending" {
		response.ParamError(c, "只有待发送状态的推送可以取消")
		return
	}

	if err := db.Model(&record).Update("status", "cancelled").Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, nil, "推送已取消")
}

// Delete 删除推送记录
func Delete(c *gin.Context) {
	id := c.Param("id")

	// 验证ID
	if _, err := validator.ValidateID(id); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	var record model.PushRecord
	if err := db.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "推送记录不存在")
			return
		}
		response.DBError(c, err)
		return
	}

	if err := db.Delete(&record).Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, nil, "推送记录删除成功")
}

// Stats 推送统计
func Stats(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		response.ParamError(c, "app_id 不能为空")
		return
	}

	var total, pending, sent, cancelled int64
	var totalSent, totalSuccess, totalFailed int64

	db.Model(&model.PushRecord{}).Where("app_id = ?", appID).Count(&total)
	db.Model(&model.PushRecord{}).Where("app_id = ? AND status = ?", appID, "pending").Count(&pending)
	db.Model(&model.PushRecord{}).Where("app_id = ? AND status = ?", appID, "sent").Count(&sent)
	db.Model(&model.PushRecord{}).Where("app_id = ? AND status = ?", appID, "cancelled").Count(&cancelled)

	db.Model(&model.PushRecord{}).Where("app_id = ?", appID).
		Select("COALESCE(SUM(sent_count), 0)").Scan(&totalSent)
	db.Model(&model.PushRecord{}).Where("app_id = ?", appID).
		Select("COALESCE(SUM(success_count), 0)").Scan(&totalSuccess)
	db.Model(&model.PushRecord{}).Where("app_id = ?", appID).
		Select("COALESCE(SUM(failed_count), 0)").Scan(&totalFailed)

	successRate := float64(0)
	if totalSent > 0 {
		successRate = float64(totalSuccess) / float64(totalSent) * 100
	}

	response.Success(c, gin.H{
		"total":         total,
		"pending":       pending,
		"sent":          sent,
		"cancelled":     cancelled,
		"total_sent":    totalSent,
		"total_success": totalSuccess,
		"total_failed":  totalFailed,
		"success_rate":  successRate,
	})
}

// Tasks 推送任务列表（兼容旧接口）
func Tasks(c *gin.Context) {
	List(c)
}

// Templates 推送模板（兼容旧接口）
func Templates(c *gin.Context) {
	response.Success(c, []gin.H{
		{"id": 1, "name": "系统通知", "title_template": "【系统通知】{{title}}", "content_template": "{{content}}"},
		{"id": 2, "name": "活动推送", "title_template": "【活动】{{title}}", "content_template": "{{content}}，点击查看详情"},
		{"id": 3, "name": "订单通知", "title_template": "订单{{order_id}}状态更新", "content_template": "您的订单{{order_id}}{{status}}"},
	})
}
