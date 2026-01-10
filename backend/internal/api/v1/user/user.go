package user

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB(database *gorm.DB) {
	db = database
	log.Println("[UserAPI] Database connection initialized")
}

// ManusUser Manus平台用户表结构
type ManusUser struct {
	ID           int        `gorm:"column:id;primaryKey" json:"id"`
	OpenID       string     `gorm:"column:openId" json:"open_id"`
	Name         *string    `gorm:"column:name" json:"name"`
	Email        *string    `gorm:"column:email" json:"email"`
	LoginMethod  *string    `gorm:"column:loginMethod" json:"login_method"`
	Role         string     `gorm:"column:role" json:"role"`
	CreatedAt    time.Time  `gorm:"column:createdAt" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updatedAt" json:"updated_at"`
	LastSignedIn *time.Time `gorm:"column:lastSignedIn" json:"last_signed_in"`
}

// TableName 指定表名
func (ManusUser) TableName() string {
	return "users"
}

// UserResponse 用户响应结构（适配前端）
type UserResponse struct {
	ID          int        `json:"id"`
	Nickname    string     `json:"nickname"`
	Phone       string     `json:"phone"`
	Email       string     `json:"email"`
	Status      int        `json:"status"`
	Role        string     `json:"role"`
	CreatedAt   time.Time  `json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at"`
}

// ListRequest 用户列表请求参数
type ListRequest struct {
	AppID  uint   `form:"app_id"`
	Page   int    `form:"page" binding:"min=1"`
	Size   int    `form:"size" binding:"min=1,max=100"`
	Status *int   `form:"status"`
	Search string `form:"search"`
}

// List 获取用户列表
func List(c *gin.Context) {
	var req ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("[UserAPI] List - Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 20
	}

	log.Printf("[UserAPI] List - Request: page=%d, size=%d, search=%s", req.Page, req.Size, req.Search)

	// 查询Manus平台的users表
	query := db.Model(&ManusUser{})

	// 搜索（按名称或邮箱）
	if req.Search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?",
			"%"+req.Search+"%", "%"+req.Search+"%")
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Printf("[UserAPI] List - Count error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to count users"})
		return
	}

	// 分页查询
	var users []ManusUser
	offset := (req.Page - 1) * req.Size
	if err := query.Offset(offset).Limit(req.Size).Order("createdAt DESC").Find(&users).Error; err != nil {
		log.Printf("[UserAPI] List - Query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query users"})
		return
	}

	// 转换为前端期望的格式
	var responseList []UserResponse
	for _, u := range users {
		name := ""
		if u.Name != nil {
			name = *u.Name
		}
		email := ""
		if u.Email != nil {
			email = *u.Email
		}
		responseList = append(responseList, UserResponse{
			ID:          u.ID,
			Nickname:    name,
			Phone:       "", // Manus平台没有phone字段
			Email:       email,
			Status:      1, // 默认正常状态
			Role:        u.Role,
			CreatedAt:   u.CreatedAt,
			LastLoginAt: u.LastSignedIn,
		})
	}

	log.Printf("[UserAPI] List - Found %d users, total %d", len(users), total)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  responseList,
			"total": total,
			"page":  req.Page,
			"size":  req.Size,
		},
	})
}

// Detail 获取用户详情
func Detail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("[UserAPI] Detail - Invalid user ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid user ID"})
		return
	}

	log.Printf("[UserAPI] Detail - Getting user %d", id)

	var user ManusUser
	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("[UserAPI] Detail - User %d not found", id)
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "User not found"})
			return
		}
		log.Printf("[UserAPI] Detail - Query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query user"})
		return
	}

	name := ""
	if user.Name != nil {
		name = *user.Name
	}
	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	response := UserResponse{
		ID:          user.ID,
		Nickname:    name,
		Phone:       "",
		Email:       email,
		Status:      1,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		LastLoginAt: user.LastSignedIn,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": response,
	})
}

// UpdateStatusRequest 更新用户状态请求参数
type UpdateStatusRequest struct {
	Status int `json:"status" binding:"oneof=0 1"`
}

// UpdateStatus 更新用户状态（Manus平台不支持，返回成功但不实际修改）
func UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("[UserAPI] UpdateStatus - Invalid user ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid user ID"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[UserAPI] UpdateStatus - Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	log.Printf("[UserAPI] UpdateStatus - User %d status to %d (note: Manus platform doesn't support status field)", id, req.Status)

	// Manus平台的users表没有status字段，这里只返回成功
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "User status updated successfully",
	})
}

// Stats 用户统计
func Stats(c *gin.Context) {
	log.Printf("[UserAPI] Stats - Getting user statistics")

	// 总用户数
	var total int64
	if err := db.Model(&ManusUser{}).Count(&total).Error; err != nil {
		log.Printf("[UserAPI] Stats - Count error: %v", err)
		total = 0
	}

	// 活跃用户数（最近7天登录）
	var active int64
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	if err := db.Model(&ManusUser{}).Where("lastSignedIn > ?", sevenDaysAgo).Count(&active).Error; err != nil {
		log.Printf("[UserAPI] Stats - Active count error: %v", err)
		active = 0
	}

	// 今日新增
	var todayNew int64
	today := time.Now().Format("2006-01-02")
	if err := db.Model(&ManusUser{}).Where("DATE(createdAt) = ?", today).Count(&todayNew).Error; err != nil {
		log.Printf("[UserAPI] Stats - Today new count error: %v", err)
		todayNew = 0
	}

	// 管理员数量
	var adminCount int64
	if err := db.Model(&ManusUser{}).Where("role = ?", "admin").Count(&adminCount).Error; err != nil {
		log.Printf("[UserAPI] Stats - Admin count error: %v", err)
		adminCount = 0
	}

	log.Printf("[UserAPI] Stats - total=%d, active=%d, todayNew=%d, admin=%d", total, active, todayNew, adminCount)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"total":     total,
			"active":    active,
			"today_new": todayNew,
			"disabled":  0,
			"normal":    total,
			"admin":     adminCount,
		},
	})
}
