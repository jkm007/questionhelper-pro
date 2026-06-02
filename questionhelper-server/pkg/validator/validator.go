package validator

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"questionhelper-server/pkg/response"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("phone", validatePhone)
	validate.RegisterValidation("idcard", validateIDCard)
	validate.RegisterValidation("password", validatePassword)
}

// ValidatePhone 手机号验证
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}

// ValidateIDCard 身份证号验证
func validateIDCard(fl validator.FieldLevel) bool {
	idCard := fl.Field().String()
	matched, _ := regexp.MatchString(`^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`, idCard)
	return matched
}

// validatePassword 密码强度验证：必须同时包含字母和数字
func validatePassword(fl validator.FieldLevel) bool {
	pwd := fl.Field().String()
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(pwd)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(pwd)
	return hasLetter && hasNumber
}

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// BindAndValidate 绑定并验证请求参数
func BindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBind(obj); err != nil {
		response.Error(c, 400, formatError(err))
		return false
	}

	if err := ValidateStruct(obj); err != nil {
		response.Error(c, 400, formatError(err))
		return false
	}

	return true
}

// formatError 格式化错误信息
func formatError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := e.Field()
			tag := e.Tag()
			param := e.Param()

			switch tag {
			case "required":
				return field + "不能为空"
			case "min":
				return field + "最小长度为" + param
			case "max":
				return field + "最大长度为" + param
			case "phone":
				return field + "格式不正确"
			case "idcard":
				return field + "格式不正确"
			case "password":
				return field + "必须包含字母和数字"
			default:
				return field + "验证失败: " + tag
			}
		}
	}

	// 处理 JSON 解析错误
	if strings.Contains(err.Error(), "EOF") {
		return "请求体不能为空"
	}
	if strings.Contains(err.Error(), "cannot unmarshal") {
		return "参数类型错误"
	}

	return "参数错误"
}

// GetFieldName 获取字段的中文名称
func GetFieldName(s interface{}, fieldName string) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, found := t.FieldByName(fieldName)
	if !found {
		return fieldName
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag != "" {
		return strings.Split(jsonTag, ",")[0]
	}

	return fieldName
}
