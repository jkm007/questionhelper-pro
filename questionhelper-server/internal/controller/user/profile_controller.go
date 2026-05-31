package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/user"
	"questionhelper-server/pkg/response"
)

type ProfileController struct{}

func NewProfileController() *ProfileController {
	return &ProfileController{}
}

// GetPrivacy 获取隐私设置
func (ctrl *ProfileController) GetPrivacy(c *gin.Context) {
	userID := c.GetUint("user_id")

	info, err := user.GetPrivacy(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// UpdatePrivacy 更新隐私设置
func (ctrl *ProfileController) UpdatePrivacy(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.UpdatePrivacyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.UpdatePrivacy(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// BindPhone 绑定手机号
func (ctrl *ProfileController) BindPhone(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.BindPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.BindPhone(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "绑定成功", nil)
}

// BindEmail 绑定邮箱
func (ctrl *ProfileController) BindEmail(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.BindEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.BindEmail(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "绑定成功", nil)
}

// SubmitRealName 提交实名认证
func (ctrl *ProfileController) SubmitRealName(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.RealNameSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.SubmitRealName(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "实名认证提交成功", nil)
}

// GetRealName 获取实名认证信息
func (ctrl *ProfileController) GetRealName(c *gin.Context) {
	userID := c.GetUint("user_id")

	info, err := user.GetRealName(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// ReviewRealName 审核实名认证(管理员)
func (ctrl *ProfileController) ReviewRealName(c *gin.Context) {
	reviewerID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的认证ID")
		return
	}

	var req dto.ReviewRealNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.ReviewRealName(uint(id), reviewerID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "审核完成", nil)
}

// GetOAuthBindings 获取第三方账号绑定
func (ctrl *ProfileController) GetOAuthBindings(c *gin.Context) {
	userID := c.GetUint("user_id")

	list, err := user.GetOAuthBindings(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// UnbindOAuth 解绑第三方账号
func (ctrl *ProfileController) UnbindOAuth(c *gin.Context) {
	userID := c.GetUint("user_id")

	provider := c.Param("provider")
	if provider == "" {
		response.Error(c, http.StatusBadRequest, "请指定第三方账号类型")
		return
	}

	if err := user.UnbindOAuth(userID, provider); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "解绑成功", nil)
}

// ListRealNames 实名认证列表(管理员)
func (ctrl *ProfileController) ListRealNames(c *gin.Context) {
	// TODO: 实现实名认证列表查询
	response.Error(c, 501, "接口未实现")
}
