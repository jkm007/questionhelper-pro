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
// @Summary      获取隐私设置
// @Description  获取当前用户的隐私设置
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=dto.PrivacyInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /user/privacy [get]
// @Security     BearerAuth
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
// @Summary      更新隐私设置
// @Description  更新当前用户的隐私设置
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Param        req  body      dto.UpdatePrivacyRequest  true  "隐私设置"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /user/privacy [put]
// @Security     BearerAuth
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
// @Summary      绑定手机号
// @Description  绑定或更换手机号
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Param        req  body      dto.BindPhoneRequest  true  "手机号信息"
// @Success      200  {object}  response.Response  "绑定成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /user/bind-phone [post]
// @Security     BearerAuth
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
// @Summary      绑定邮箱
// @Description  绑定或更换邮箱
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Param        req  body      dto.BindEmailRequest  true  "邮箱信息"
// @Success      200  {object}  response.Response  "绑定成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /user/bind-email [post]
// @Security     BearerAuth
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
// @Summary      提交实名认证
// @Description  提交实名认证信息
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Param        req  body      dto.RealNameSubmitRequest  true  "实名认证信息"
// @Success      200  {object}  response.Response  "实名认证提交成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /user/realname [post]
// @Security     BearerAuth
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
// @Summary      获取实名认证信息
// @Description  获取当前用户的实名认证信息
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=dto.RealNameInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /user/realname [get]
// @Security     BearerAuth
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
// @Summary      审核实名认证
// @Description  审核用户的实名认证申请（管理员）
// @Tags         实名认证
// @Accept       json
// @Produce      json
// @Param        id   path      uint                          true  "认证ID"
// @Param        req  body      dto.ReviewRealNameRequest     true  "审核信息"
// @Success      200  {object}  response.Response  "审核完成"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/realnames/{id}/review [put]
// @Security     BearerAuth
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
// @Summary      获取第三方账号绑定
// @Description  获取当前用户绑定的第三方账号列表
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.OAuthInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /user/oauth [get]
// @Security     BearerAuth
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
// @Summary      解绑第三方账号
// @Description  解绑指定的第三方账号
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Param        provider  path      string  true  "第三方账号类型"
// @Success      200       {object}  response.Response  "解绑成功"
// @Failure      400       {object}  response.Response  "参数错误"
// @Failure      500       {object}  response.Response  "服务器内部错误"
// @Router       /user/oauth/{provider} [delete]
// @Security     BearerAuth
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
// @Summary      实名认证列表
// @Description  分页查询实名认证申请列表（管理员）
// @Tags         实名认证
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"       default(1)
// @Param        page_size  query     int     false  "每页数量"   default(10)
// @Param        status     query     int     false  "审核状态"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/realnames [get]
// @Security     BearerAuth
func (ctrl *ProfileController) ListRealNames(c *gin.Context) {
	// TODO: 实现实名认证列表查询
	response.Error(c, 501, "接口未实现")
}
