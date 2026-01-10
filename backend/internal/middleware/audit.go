package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var auditDB *gorm.DB

// AuditLog 审计日志模型
type AuditLog struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	AppID         uint      `json:"app_id" gorm:"index"`
	UserID        string    `json:"user_id" gorm:"index"`
	UserName      string    `json:"user_name"`
	Action        string    `json:"action" gorm:"index"`
	Resource      string    `json:"resource" gorm:"index"`
	ResourceID    string    `json:"resource_id"`
	Description   string    `json:"description"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	RequestPath   string    `json:"request_path"`
	RequestMethod string    `json:"request_method"`
	RequestBody   string    `json:"request_body" gorm:"type:text"`
	StatusCode    int       `json:"status_code"`
	Duration      int64     `json:"duration"`
	Extra         string    `json:"extra" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at" gorm:"index"`
}

// AuditConfig 审计配置
type AuditConfig struct {
	EnableRequestBody bool     // 是否记录请求体
	MaxBodySize       int      // 最大请求体大小
	SensitiveFields   []string // 敏感字段（需要脱敏）
	SkipPaths         []string // 跳过的路径
}

var defaultAuditConfig = AuditConfig{
	EnableRequestBody: true,
	MaxBodySize:       4096,
	SensitiveFields:   []string{"password", "token", "secret", "app_secret"},
	SkipPaths: []string{
		"/api/v1/health",
		"/api/v1/ws",
		"/api/v1/monitor/metrics",
		"/api/v1/logs/report",
		"/api/v1/events/report",
		"/api/v1/audit",
	},
}

// InitAuditDB 初始化审计数据库连接
func InitAuditDB(db *gorm.DB) {
	auditDB = db
}

// AuditMiddleware 审计日志中间件
func AuditMiddleware() gin.HandlerFunc {
	return AuditMiddlewareWithConfig(defaultAuditConfig)
}

// AuditMiddlewareWithConfig 带配置的审计日志中间件
func AuditMiddlewareWithConfig(config AuditConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要审计的路径
		path := c.Request.URL.Path
		if shouldSkipAuditPath(path, config.SkipPaths) {
			c.Next()
			return
		}

		startTime := time.Now()

		// 读取请求体（用于记录）
		var bodyBytes []byte
		var sanitizedBody string
		if c.Request.Body != nil && config.EnableRequestBody {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			
			// 限制请求体大小
			if len(bodyBytes) > config.MaxBodySize {
				bodyBytes = bodyBytes[:config.MaxBodySize]
			}
			
			// 脱敏处理
			sanitizedBody = sanitizeRequestBody(string(bodyBytes), config.SensitiveFields)
		}

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(startTime).Milliseconds()

		// 异步记录审计日志
		go recordAuditLogEnhanced(c, duration, sanitizedBody)
	}
}

// shouldSkipAuditPath 判断是否跳过审计
func shouldSkipAuditPath(path string, skipPaths []string) bool {
	for _, skip := range skipPaths {
		if path == skip || strings.HasPrefix(path, skip+"/") {
			return true
		}
	}
	return false
}

// sanitizeRequestBody 脱敏请求体中的敏感字段
func sanitizeRequestBody(body string, sensitiveFields []string) string {
	if body == "" {
		return ""
	}

	// 尝试解析JSON
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		// 非JSON格式，直接返回
		return body
	}

	// 递归脱敏
	sanitizeMap(data, sensitiveFields)

	// 重新序列化
	result, err := json.Marshal(data)
	if err != nil {
		return body
	}
	return string(result)
}

// sanitizeMap 递归脱敏map中的敏感字段
func sanitizeMap(data map[string]interface{}, sensitiveFields []string) {
	for key, value := range data {
		// 检查是否是敏感字段
		for _, field := range sensitiveFields {
			if strings.EqualFold(key, field) {
				data[key] = "***REDACTED***"
				break
			}
		}

		// 递归处理嵌套对象
		if nested, ok := value.(map[string]interface{}); ok {
			sanitizeMap(nested, sensitiveFields)
		}
	}
}

// recordAuditLogEnhanced 增强的审计日志记录
func recordAuditLogEnhanced(c *gin.Context, duration int64, requestBody string) {
	if auditDB == nil {
		return
	}

	// 从上下文获取用户信息
	userID := getStringFromContext(c, "user_id")
	userName := getStringFromContext(c, "user_name")
	
	// 尝试从JWT claims获取
	if userID == "" {
		if claims, exists := c.Get("claims"); exists {
			if claimsMap, ok := claims.(map[string]interface{}); ok {
				if uid, ok := claimsMap["user_id"]; ok {
					userID = toString(uid)
				}
				if uname, ok := claimsMap["username"]; ok {
					userName = toString(uname)
				}
			}
		}
	}

	// 解析操作类型和资源
	action, resource := parseActionAndResourceEnhanced(c.Request.Method, c.Request.URL.Path)

	// 获取资源ID
	resourceID := extractResourceID(c)

	// 获取AppID
	appID := extractAppID(c)

	// 生成详细描述
	description := generateDetailedDescription(action, resource, c.Request.Method, c.Request.URL.Path, resourceID)

	// 构建额外信息
	extra := buildExtraInfo(c)

	log := &AuditLog{
		AppID:         appID,
		UserID:        userID,
		UserName:      userName,
		Action:        action,
		Resource:      resource,
		ResourceID:    resourceID,
		Description:   description,
		IPAddress:     c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		RequestPath:   c.Request.URL.Path,
		RequestMethod: c.Request.Method,
		RequestBody:   requestBody,
		StatusCode:    c.Writer.Status(),
		Duration:      duration,
		Extra:         extra,
		CreatedAt:     time.Now(),
	}

	if err := auditDB.Create(log).Error; err != nil {
		// 记录失败时打印日志，但不影响业务
		println("Failed to create audit log:", err.Error())
	}
}

// getStringFromContext 从上下文获取字符串值
func getStringFromContext(c *gin.Context, key string) string {
	if value, exists := c.Get(key); exists {
		return toString(value)
	}
	return ""
}

// toString 将任意类型转换为字符串
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return intToString(int64(val))
	case int64:
		return intToString(val)
	case uint:
		return uintToString(uint64(val))
	case uint64:
		return uintToString(val)
	case float64:
		return intToString(int64(val))
	default:
		if b, err := json.Marshal(val); err == nil {
			return string(b)
		}
		return ""
	}
}

// intToString 将int64转换为字符串
func intToString(n int64) string {
	if n == 0 {
		return "0"
	}
	var result []byte
	negative := n < 0
	if negative {
		n = -n
	}
	for n > 0 {
		result = append([]byte{byte('0' + n%10)}, result...)
		n /= 10
	}
	if negative {
		result = append([]byte{'-'}, result...)
	}
	return string(result)
}

// uintToString 将uint64转换为字符串
func uintToString(n uint64) string {
	if n == 0 {
		return "0"
	}
	var result []byte
	for n > 0 {
		result = append([]byte{byte('0' + n%10)}, result...)
		n /= 10
	}
	return string(result)
}

// extractResourceID 提取资源ID
func extractResourceID(c *gin.Context) string {
	// 优先从路径参数获取
	if id := c.Param("id"); id != "" {
		return id
	}
	// 其次从查询参数获取
	if id := c.Query("id"); id != "" {
		return id
	}
	// 尝试从路径中提取数字ID
	parts := strings.Split(c.Request.URL.Path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if isNumeric(parts[i]) {
			return parts[i]
		}
	}
	return ""
}

// extractAppID 提取AppID
func extractAppID(c *gin.Context) uint {
	// 从路径参数获取
	if appID := c.Param("app_id"); appID != "" {
		if id := parseUint(appID); id > 0 {
			return id
		}
	}
	// 从查询参数获取
	if appID := c.Query("app_id"); appID != "" {
		if id := parseUint(appID); id > 0 {
			return id
		}
	}
	// 从上下文获取
	if appID, exists := c.Get("app_id"); exists {
		if id, ok := appID.(uint); ok {
			return id
		}
	}
	return 0
}

// isNumeric 检查字符串是否为数字
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// parseUint 解析无符号整数
func parseUint(s string) uint {
	var result uint
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + uint(c-'0')
		} else {
			return 0
		}
	}
	return result
}

// parseActionAndResourceEnhanced 增强的操作类型和资源解析
func parseActionAndResourceEnhanced(method, path string) (action, resource string) {
	// 根据HTTP方法判断基本操作类型
	switch method {
	case "GET":
		action = "view"
	case "POST":
		action = "create"
	case "PUT", "PATCH":
		action = "update"
	case "DELETE":
		action = "delete"
	default:
		action = "unknown"
	}

	// 解析资源类型
	resourceMap := map[string]string{
		"/api/v1/users":    "user",
		"/api/v1/apps":     "app",
		"/api/v1/configs":  "config",
		"/api/v1/messages": "message",
		"/api/v1/push":     "push",
		"/api/v1/files":    "file",
		"/api/v1/versions": "version",
		"/api/v1/logs":     "log",
		"/api/v1/events":   "event",
		"/api/v1/monitor":  "monitor",
		"/api/v1/admin":    "admin",
		"/api/v1/modules":  "module",
		"/api/v1/audit":    "audit",
		"/api/v1/alerts":   "alert",
	}

	resource = "unknown"
	for prefix, res := range resourceMap {
		if strings.HasPrefix(path, prefix) {
			resource = res
			break
		}
	}

	// 特殊操作识别（覆盖基本操作类型）
	specialActions := map[string]string{
		"/login":     "login",
		"/logout":    "logout",
		"/export":    "export",
		"/import":    "import",
		"/publish":   "publish",
		"/unpublish": "unpublish",
		"/send":      "send",
		"/enable":    "enable",
		"/disable":   "disable",
		"/activate":  "activate",
		"/deactivate": "deactivate",
		"/upload":    "upload",
		"/download":  "download",
		"/batch":     "batch",
		"/sync":      "sync",
		"/reset":     "reset",
		"/verify":    "verify",
		"/approve":   "approve",
		"/reject":    "reject",
	}

	for suffix, act := range specialActions {
		if strings.Contains(path, suffix) {
			action = act
			break
		}
	}

	return action, resource
}

// generateDetailedDescription 生成详细的操作描述
func generateDetailedDescription(action, resource, method, path, resourceID string) string {
	actionNames := map[string]string{
		"view":       "查看",
		"create":     "创建",
		"update":     "更新",
		"delete":     "删除",
		"login":      "登录",
		"logout":     "登出",
		"export":     "导出",
		"import":     "导入",
		"publish":    "发布",
		"unpublish":  "取消发布",
		"send":       "发送",
		"enable":     "启用",
		"disable":    "禁用",
		"activate":   "激活",
		"deactivate": "停用",
		"upload":     "上传",
		"download":   "下载",
		"batch":      "批量操作",
		"sync":       "同步",
		"reset":      "重置",
		"verify":     "验证",
		"approve":    "审批通过",
		"reject":     "审批拒绝",
	}

	resourceNames := map[string]string{
		"user":    "用户",
		"app":     "应用",
		"config":  "配置",
		"message": "消息",
		"push":    "推送",
		"file":    "文件",
		"version": "版本",
		"log":     "日志",
		"event":   "事件",
		"monitor": "监控",
		"admin":   "管理员",
		"module":  "模块",
		"audit":   "审计日志",
		"alert":   "告警",
	}

	actionName := actionNames[action]
	if actionName == "" {
		actionName = action
	}
	resourceName := resourceNames[resource]
	if resourceName == "" {
		resourceName = resource
	}

	desc := actionName + resourceName
	if resourceID != "" {
		desc += " (ID: " + resourceID + ")"
	}

	return desc
}

// buildExtraInfo 构建额外信息
func buildExtraInfo(c *gin.Context) string {
	extra := map[string]interface{}{
		"query_params": c.Request.URL.RawQuery,
		"referer":      c.Request.Referer(),
		"content_type": c.ContentType(),
	}

	// 添加响应大小
	extra["response_size"] = c.Writer.Size()

	result, _ := json.Marshal(extra)
	return string(result)
}

// 保留旧函数以保持兼容性
func shouldSkipAudit(path string) bool {
	return shouldSkipAuditPath(path, defaultAuditConfig.SkipPaths)
}

func recordAuditLog(c *gin.Context, duration int64) {
	recordAuditLogEnhanced(c, duration, "")
}

func parseActionAndResource(method, path string) (action, resource string) {
	return parseActionAndResourceEnhanced(method, path)
}

func generateDescription(action, resource, method, path string) string {
	return generateDetailedDescription(action, resource, method, path, "")
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
