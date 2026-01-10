package validator

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// AppCreateRequest APP创建请求
type AppCreateRequest struct {
	Name        string   `json:"name"`
	AppName     string   `json:"app_name"`
	PackageName string   `json:"package_name"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	Modules     []string `json:"modules"`
}

// AppUpdateRequest APP更新请求
type AppUpdateRequest struct {
	Name        string `json:"name"`
	PackageName string `json:"package_name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Status      *int   `json:"status"`
}

// ValidateAppCreate 验证APP创建请求
func ValidateAppCreate(req *AppCreateRequest) error {
	// 获取实际的名称
	name := req.Name
	if name == "" {
		name = req.AppName
	}

	// 验证名称
	if err := validateAppName(name); err != nil {
		return err
	}

	// 验证包名（如果提供）
	if req.PackageName != "" {
		if err := validatePackageName(req.PackageName); err != nil {
			return err
		}
	}

	// 验证描述长度
	if utf8.RuneCountInString(req.Description) > 500 {
		return errors.New("描述不能超过500个字符")
	}

	// 验证图标URL
	if req.Icon != "" {
		if err := validateURL(req.Icon); err != nil {
			return errors.New("图标URL格式不正确")
		}
	}

	// 验证模块列表
	if len(req.Modules) > 20 {
		return errors.New("启用的模块数量不能超过20个")
	}

	return nil
}

// ValidateAppUpdate 验证APP更新请求
func ValidateAppUpdate(req *AppUpdateRequest) error {
	// 验证名称（如果提供）
	if req.Name != "" {
		if err := validateAppName(req.Name); err != nil {
			return err
		}
	}

	// 验证包名（如果提供）
	if req.PackageName != "" {
		if err := validatePackageName(req.PackageName); err != nil {
			return err
		}
	}

	// 验证描述长度
	if utf8.RuneCountInString(req.Description) > 500 {
		return errors.New("描述不能超过500个字符")
	}

	// 验证状态值
	if req.Status != nil && (*req.Status != 0 && *req.Status != 1) {
		return errors.New("状态值只能是0或1")
	}

	return nil
}

// validateAppName 验证APP名称
func validateAppName(name string) error {
	if name == "" {
		return errors.New("应用名称不能为空")
	}

	nameLen := utf8.RuneCountInString(name)
	if nameLen < 2 {
		return errors.New("应用名称至少需要2个字符")
	}
	if nameLen > 50 {
		return errors.New("应用名称不能超过50个字符")
	}

	// 检查是否包含特殊字符
	if containsSpecialChars(name) {
		return errors.New("应用名称不能包含特殊字符")
	}

	return nil
}

// validatePackageName 验证包名
func validatePackageName(packageName string) error {
	// Android包名格式: com.example.app
	// iOS Bundle ID格式: com.example.app
	pattern := `^[a-zA-Z][a-zA-Z0-9]*(\.[a-zA-Z][a-zA-Z0-9]*)+$`
	matched, _ := regexp.MatchString(pattern, packageName)
	if !matched {
		return errors.New("包名格式不正确，应为类似 com.example.app 的格式")
	}

	if len(packageName) > 100 {
		return errors.New("包名长度不能超过100个字符")
	}

	return nil
}

// validateURL 验证URL格式
func validateURL(url string) error {
	if url == "" {
		return nil
	}

	// 简单的URL验证
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("URL必须以http://或https://开头")
	}

	if len(url) > 500 {
		return errors.New("URL长度不能超过500个字符")
	}

	return nil
}

// containsSpecialChars 检查是否包含特殊字符
func containsSpecialChars(s string) bool {
	// 允许中文、英文、数字、空格、下划线、连字符
	pattern := `^[\p{Han}a-zA-Z0-9\s_-]+$`
	matched, _ := regexp.MatchString(pattern, s)
	return !matched
}

// SanitizeString 清理字符串，防止XSS
func SanitizeString(s string) string {
	// 移除HTML标签
	s = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(s, "")
	// 转义特殊字符
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

// ValidateID 验证ID参数
func ValidateID(id string) (uint, error) {
	if id == "" {
		return 0, errors.New("ID不能为空")
	}

	// 检查是否为纯数字
	pattern := `^[1-9]\d*$`
	matched, _ := regexp.MatchString(pattern, id)
	if !matched {
		return 0, errors.New("ID格式不正确")
	}

	// 转换为uint
	var result uint
	for _, c := range id {
		result = result*10 + uint(c-'0')
	}

	return result, nil
}

// ValidatePagination 验证分页参数
func ValidatePagination(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return page, size
}

// ParsePagination 解析分页参数字符串
func ParsePagination(pageStr, sizeStr string) (int, int) {
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	return page, size
}

// ValidateURL 验证URL格式（导出版本）
func ValidateURL(url string) error {
	return validateURL(url)
}
