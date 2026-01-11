package baas

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DataCollection 数据模型定义
type DataCollection struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	AppID       uint            `json:"app_id" gorm:"index"`
	Name        string          `json:"name" gorm:"size:50;not null"`
	DisplayName string          `json:"display_name" gorm:"size:100"`
	Description string          `json:"description" gorm:"type:text"`
	Schema      json.RawMessage `json:"schema" gorm:"type:json"`
	Indexes     json.RawMessage `json:"indexes" gorm:"type:json"`
	Permissions json.RawMessage `json:"permissions" gorm:"type:json"`
	Hooks       json.RawMessage `json:"hooks" gorm:"type:json"`
	Status      int             `json:"status" gorm:"default:1"`
	CreatedBy   uint            `json:"created_by"`
	Fields      json.RawMessage `json:"fields" gorm:"type:json"`
	ReadPerm    string          `json:"read_perm" gorm:"size:50;default:public"`
	CreatePerm  string          `json:"create_perm" gorm:"size:50;default:authenticated"`
	UpdatePerm  string          `json:"update_perm" gorm:"size:50;default:creator"`
	DeletePerm  string          `json:"delete_perm" gorm:"size:50;default:admin"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (DataCollection) TableName() string {
	return "data_collections"
}

// DataDocument 数据文档定义
type DataDocument struct {
	ID           uint            `json:"id" gorm:"primaryKey"`
	CollectionID uint            `json:"collection_id" gorm:"index"`
	AppID        uint            `json:"app_id" gorm:"index"`
	Data         json.RawMessage `json:"data" gorm:"type:json"`
	CreatedBy    uint            `json:"created_by"`
	UpdatedBy    uint            `json:"updated_by"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

func (DataDocument) TableName() string {
	return "data_documents"
}

// Handler BaaS API处理器
type Handler struct {
	db *gorm.DB
}

// NewHandler 创建新的Handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	baas := r.Group("/baas")
	{
		apps := baas.Group("/apps/:appId")
		{
			// 数据模型管理
			apps.GET("/collections", h.ListCollections)
			apps.POST("/collections", h.CreateCollection)
			apps.GET("/collections/:collectionId", h.GetCollection)
			apps.PUT("/collections/:collectionId", h.UpdateCollection)
			apps.DELETE("/collections/:collectionId", h.DeleteCollection)

			// 数据文档管理
			apps.GET("/data/:collectionName", h.ListDocuments)
			apps.POST("/data/:collectionName", h.CreateDocument)
			apps.GET("/data/:collectionName/:docId", h.GetDocument)
			apps.PUT("/data/:collectionName/:docId", h.UpdateDocument)
			apps.DELETE("/data/:collectionName/:docId", h.DeleteDocument)
		}
	}
}

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// ListCollections 获取数据模型列表
func (h *Handler) ListCollections(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Param("appId"), 10, 64)
	if err != nil {
		fail(c, 400, "无效的应用ID")
		return
	}

	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	var collections []DataCollection
	var total int64

	query := h.db.Model(&DataCollection{}).Where("app_id = ?", appID)
	if search != "" {
		query = query.Where("name LIKE ? OR display_name LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	query.Offset((page - 1) * size).Limit(size).Order("id DESC").Find(&collections)

	success(c, gin.H{
		"list":  collections,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// CreateCollection 创建数据模型
func (h *Handler) CreateCollection(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Param("appId"), 10, 64)
	if err != nil {
		fail(c, 400, "无效的应用ID")
		return
	}

	var req struct {
		Name        string          `json:"name" binding:"required"`
		DisplayName string          `json:"display_name"`
		Description string          `json:"description"`
		Fields      json.RawMessage `json:"fields"`
		ReadPerm    string          `json:"read_perm"`
		CreatePerm  string          `json:"create_perm"`
		UpdatePerm  string          `json:"update_perm"`
		DeletePerm  string          `json:"delete_perm"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 检查名称是否已存在
	var existing DataCollection
	if err := h.db.Where("app_id = ? AND name = ?", appID, req.Name).First(&existing).Error; err == nil {
		fail(c, 400, "数据模型名称已存在")
		return
	}

	// 设置默认权限
	if req.ReadPerm == "" {
		req.ReadPerm = "public"
	}
	if req.CreatePerm == "" {
		req.CreatePerm = "authenticated"
	}
	if req.UpdatePerm == "" {
		req.UpdatePerm = "creator"
	}
	if req.DeletePerm == "" {
		req.DeletePerm = "creator"
	}

	// 创建空的schema
	emptySchema := json.RawMessage(`{}`)

	collection := DataCollection{
		AppID:       uint(appID),
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Schema:      emptySchema,
		Fields:      req.Fields,
		ReadPerm:    req.ReadPerm,
		CreatePerm:  req.CreatePerm,
		UpdatePerm:  req.UpdatePerm,
		DeletePerm:  req.DeletePerm,
		Status:      1,
	}

	if err := h.db.Create(&collection).Error; err != nil {
		fail(c, 500, "创建数据模型失败: "+err.Error())
		return
	}

	success(c, collection)
}

// GetCollection 获取单个数据模型
func (h *Handler) GetCollection(c *gin.Context) {
	appID, _ := strconv.ParseUint(c.Param("appId"), 10, 64)
	collectionID, _ := strconv.ParseUint(c.Param("collectionId"), 10, 64)

	var collection DataCollection
	if err := h.db.Where("id = ? AND app_id = ?", collectionID, appID).First(&collection).Error; err != nil {
		fail(c, 404, "数据模型不存在")
		return
	}

	success(c, collection)
}

// UpdateCollection 更新数据模型
func (h *Handler) UpdateCollection(c *gin.Context) {
	appID, _ := strconv.ParseUint(c.Param("appId"), 10, 64)
	collectionID, _ := strconv.ParseUint(c.Param("collectionId"), 10, 64)

	var collection DataCollection
	if err := h.db.Where("id = ? AND app_id = ?", collectionID, appID).First(&collection).Error; err != nil {
		fail(c, 404, "数据模型不存在")
		return
	}

	var req struct {
		Name        string          `json:"name"`
		DisplayName string          `json:"display_name"`
		Description string          `json:"description"`
		Fields      json.RawMessage `json:"fields"`
		ReadPerm    string          `json:"read_perm"`
		CreatePerm  string          `json:"create_perm"`
		UpdatePerm  string          `json:"update_perm"`
		DeletePerm  string          `json:"delete_perm"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 检查名称是否与其他模型冲突
	if req.Name != "" && req.Name != collection.Name {
		var existing DataCollection
		if err := h.db.Where("app_id = ? AND name = ? AND id != ?", appID, req.Name, collectionID).First(&existing).Error; err == nil {
			fail(c, 400, "数据模型名称已存在")
			return
		}
		collection.Name = req.Name
	}

	if req.DisplayName != "" {
		collection.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		collection.Description = req.Description
	}
	if req.Fields != nil {
		collection.Fields = req.Fields
	}
	if req.ReadPerm != "" {
		collection.ReadPerm = req.ReadPerm
	}
	if req.CreatePerm != "" {
		collection.CreatePerm = req.CreatePerm
	}
	if req.UpdatePerm != "" {
		collection.UpdatePerm = req.UpdatePerm
	}
	if req.DeletePerm != "" {
		collection.DeletePerm = req.DeletePerm
	}

	if err := h.db.Save(&collection).Error; err != nil {
		fail(c, 500, "更新数据模型失败: "+err.Error())
		return
	}

	success(c, collection)
}

// DeleteCollection 删除数据模型
func (h *Handler) DeleteCollection(c *gin.Context) {
	appID, _ := strconv.ParseUint(c.Param("appId"), 10, 64)
	collectionID, _ := strconv.ParseUint(c.Param("collectionId"), 10, 64)

	var collection DataCollection
	if err := h.db.Where("id = ? AND app_id = ?", collectionID, appID).First(&collection).Error; err != nil {
		fail(c, 404, "数据模型不存在")
		return
	}

	// 删除关联的文档
	h.db.Where("collection_id = ?", collectionID).Delete(&DataDocument{})

	// 删除数据模型
	if err := h.db.Delete(&collection).Error; err != nil {
		fail(c, 500, "删除数据模型失败: "+err.Error())
		return
	}

	success(c, nil)
}

// ListDocuments 获取文档列表
func (h *Handler) ListDocuments(c *gin.Context) {
	appID, _ := strconv.ParseUint(c.Param("appId"), 10, 64)
	collectionName := c.Param("collectionName")

	// 查找数据模型
	var collection DataCollection
	if err := h.db.Where("app_id = ? AND name = ?", appID, collectionName).First(&collection).Error; err != nil {
		fail(c, 404, "数据模型不存在")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	var documents []DataDocument
	var total int64

	query := h.db.Model(&DataDocument{}).Where("collection_id = ?", collection.ID)
	query.Count(&total)
	query.Offset((page - 1) * size).Limit(size).Order("id DESC").Find(&documents)

	success(c, gin.H{
		"list":  documents,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// FieldDefinition 字段定义
type FieldDefinition struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Unique      bool   `json:"unique"`
}

// Schema 数据模型结构
type Schema struct {
	Fields []FieldDefinition `json:"fields"`
}

// validateDocument 验证文档数据
func (h *Handler) validateDocument(collection *DataCollection, data map[string]interface{}) error {
	if collection.Fields == nil {
		return nil
	}

	var schema Schema
	if err := json.Unmarshal(collection.Fields, &schema); err != nil {
		return nil // 如果解析失败，跳过验证
	}

	for _, field := range schema.Fields {
		value, exists := data[field.Name]

		// 必填验证
		if field.Required && (!exists || value == nil || value == "") {
			return &ValidationError{Field: field.DisplayName, Message: "不能为空"}
		}

		if !exists || value == nil {
			continue
		}

		// 类型验证
		switch field.Type {
		case "string":
			if _, ok := value.(string); !ok {
				return &ValidationError{Field: field.DisplayName, Message: "应为字符串类型"}
			}
		case "number":
			switch value.(type) {
			case float64, int, int64, float32:
				// OK
			default:
				return &ValidationError{Field: field.DisplayName, Message: "应为数字类型"}
			}
		case "boolean":
			if _, ok := value.(bool); !ok {
				return &ValidationError{Field: field.DisplayName, Message: "应为布尔类型"}
			}
		case "array":
			if _, ok := value.([]interface{}); !ok {
				return &ValidationError{Field: field.DisplayName, Message: "应为数组类型"}
			}
		case "object":
			if _, ok := value.(map[string]interface{}); !ok {
				return &ValidationError{Field: field.DisplayName, Message: "应为对象类型"}
			}
		}

		// 唯一性验证
		if field.Unique && value != nil && value != "" {
			var count int64
			// 使用JSON查询检查唯一性
			h.db.Model(&DataDocument{}).Where("collection_id = ? AND JSON_EXTRACT(data, ?) = ?", 
				collection.ID, "$."+field.Name, value).Count(&count)
			if count > 0 {
				return &ValidationError{Field: field.DisplayName, Message: "已存在相同的值"}
			}
		}
	}

	return nil
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// validateDocumentForUpdate 更新时验证文档数据（排除当前文档的唯一性检查）
func (h *Handler) validateDocumentForUpdate(collection *DataCollection, data map[string]interface{}, excludeDocID uint) error {
	if collection.Fields == nil {
		return nil
	}

	var schema Schema
	if err := json.Unmarshal(collection.Fields, &schema); err != nil {
		return nil
	}

	for _, field := range schema.Fields {
		value, exists := data[field.Name]

		// 必填验证
		if field.Required && (!exists || value == nil || value == "") {
			return &ValidationError{Field: field.DisplayName, Message: "不能为空"}
		}

		if !exists || value == nil {
			continue
		}

		// 类型验证
		switch field.Type {
		case "string":
			if _, ok := value.(string); !ok {
				return &ValidationError{Field: field.DisplayName, Message: "应为字符串类型"}
			}
		case "number":
			switch value.(type) {
			case float64, int, int64, float32:
				// OK
			default:
				return &ValidationError{Field: field.DisplayName, Message: "应为数字类型"}
			}
		case "boolean":
			if _, ok := value.(bool); !ok {
				return &ValidationError{Field: field.DisplayName, Message: "应为布尔类型"}
			}
		case "array":
			if _, ok := value.([]interface{}); !ok {
				return &ValidationError{Field: field.DisplayName, Message: "应为数组类型"}
			}
		case "object":
			if _, ok := value.(map[string]interface{}); !ok {
				return &ValidationError{Field: field.DisplayName, Message: "应为对象类型"}
			}
		}

		// 唯一性验证（排除当前文档）
		if field.Unique && value != nil && value != "" {
			var count int64
			h.db.Model(&DataDocument{}).Where("collection_id = ? AND id != ? AND JSON_EXTRACT(data, ?) = ?", 
				collection.ID, excludeDocID, "$."+field.Name, value).Count(&count)
			if count > 0 {
				return &ValidationError{Field: field.DisplayName, Message: "已存在相同的值"}
			}
		}
	}

	return nil
}

// CreateDocument 创建文档
func (h *Handler) CreateDocument(c *gin.Context) {
	appID, _ := strconv.ParseUint(c.Param("appId"), 10, 64)
	collectionName := c.Param("collectionName")

	// 查找数据模型
	var collection DataCollection
	if err := h.db.Where("app_id = ? AND name = ?", appID, collectionName).First(&collection).Error; err != nil {
		fail(c, 404, "数据模型不存在")
		return
	}

	var req struct {
		Data json.RawMessage `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 解析数据进行验证
	var dataMap map[string]interface{}
	if err := json.Unmarshal(req.Data, &dataMap); err != nil {
		fail(c, 400, "数据格式错误")
		return
	}

	// 验证数据
	if err := h.validateDocument(&collection, dataMap); err != nil {
		fail(c, 400, "数据验证失败: "+err.Error())
		return
	}

	// 获取当前用户ID (从 JWT中)
	userID := uint(0)
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	document := DataDocument{
		CollectionID: collection.ID,
		AppID:        uint(appID),
		Data:         req.Data,
		CreatedBy:    userID,
		UpdatedBy:    userID,
	}

	if err := h.db.Create(&document).Error; err != nil {
		fail(c, 500, "创建文档失败: "+err.Error())
		return
	}

	success(c, document)
}

// GetDocument 获取单个文档
func (h *Handler) GetDocument(c *gin.Context) {
	appID, _ := strconv.ParseUint(c.Param("appId"), 10, 64)
	collectionName := c.Param("collectionName")
	docID, _ := strconv.ParseUint(c.Param("docId"), 10, 64)

	// 查找数据模型
	var collection DataCollection
	if err := h.db.Where("app_id = ? AND name = ?", appID, collectionName).First(&collection).Error; err != nil {
		fail(c, 404, "数据模型不存在")
		return
	}

	var document DataDocument
	if err := h.db.Where("id = ? AND collection_id = ?", docID, collection.ID).First(&document).Error; err != nil {
		fail(c, 404, "文档不存在")
		return
	}

	success(c, document)
}

// UpdateDocument 更新文档
func (h *Handler) UpdateDocument(c *gin.Context) {
	appID, _ := strconv.ParseUint(c.Param("appId"), 10, 64)
	collectionName := c.Param("collectionName")
	docID, _ := strconv.ParseUint(c.Param("docId"), 10, 64)

	// 查找数据模型
	var collection DataCollection
	if err := h.db.Where("app_id = ? AND name = ?", appID, collectionName).First(&collection).Error; err != nil {
		fail(c, 404, "数据模型不存在")
		return
	}

	var document DataDocument
	if err := h.db.Where("id = ? AND collection_id = ?", docID, collection.ID).First(&document).Error; err != nil {
		fail(c, 404, "文档不存在")
		return
	}

	var req struct {
		Data json.RawMessage `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 解析数据进行验证
	var dataMap map[string]interface{}
	if err := json.Unmarshal(req.Data, &dataMap); err != nil {
		fail(c, 400, "数据格式错误")
		return
	}

	// 验证数据（更新时跳过唯一性检查）
	if err := h.validateDocumentForUpdate(&collection, dataMap, document.ID); err != nil {
		fail(c, 400, "数据验证失败: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID := uint(0)
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	document.Data = req.Data
	document.UpdatedBy = userID

	if err := h.db.Save(&document).Error; err != nil {
		fail(c, 500, "更新文档失败: "+err.Error())
		return
	}

	success(c, document)
}

// DeleteDocument 删除文档
func (h *Handler) DeleteDocument(c *gin.Context) {
	appID, _ := strconv.ParseUint(c.Param("appId"), 10, 64)
	collectionName := c.Param("collectionName")
	docID, _ := strconv.ParseUint(c.Param("docId"), 10, 64)

	// 查找数据模型
	var collection DataCollection
	if err := h.db.Where("app_id = ? AND name = ?", appID, collectionName).First(&collection).Error; err != nil {
		fail(c, 404, "数据模型不存在")
		return
	}

	var document DataDocument
	if err := h.db.Where("id = ? AND collection_id = ?", docID, collection.ID).First(&document).Error; err != nil {
		fail(c, 404, "文档不存在")
		return
	}

	if err := h.db.Delete(&document).Error; err != nil {
		fail(c, 500, "删除文档失败: "+err.Error())
		return
	}

	success(c, nil)
}

// MigrateDB 数据库迁移
func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(&DataCollection{}, &DataDocument{})
}
