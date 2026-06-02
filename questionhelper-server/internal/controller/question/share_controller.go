package question

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/question"
	"questionhelper-server/pkg/response"
)

type ShareController struct{}

func NewShareController() *ShareController {
	return &ShareController{}
}

// CreateShare 创建分享
// @Summary      创建分享
// @Description  管理员为指定题目创建分享链接
// @Tags         题目分享
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateShareRequest  true  "分享信息"
// @Success      200  {object}  response.Response{data=dto.ShareInfo}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/questions/{id}/share [post]
// @Security     BearerAuth
func (ctrl *ShareController) CreateShare(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	info, err := question.CreateShare(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// GetShare 获取分享
// @Summary      获取分享信息
// @Description  通过分享码获取分享的题目信息
// @Tags         题目分享
// @Accept       json
// @Produce      json
// @Param        code  path      string  true  "分享码"
// @Success      200  {object}  response.Response{data=dto.ShareInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /shares/{code} [get]
// @Security     BearerAuth
func (ctrl *ShareController) GetShare(c *gin.Context) {
	code := c.Param("code")

	info, err := question.GetShare(code)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// RevokeShare 撤销分享
// @Summary      撤销分享
// @Description  管理员撤销指定分享链接
// @Tags         题目分享
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "分享ID"
// @Success      200  {object}  response.Response  "撤销成功"
// @Failure      400  {object}  response.Response  "无效的分享ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/shares/{id} [delete]
// @Security     BearerAuth
func (ctrl *ShareController) RevokeShare(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分享ID")
		return
	}

	if err := question.RevokeShare(uint(id), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "撤销成功", nil)
}

// ListMyShares 我的分享列表
// @Summary      获取我的分享列表
// @Description  获取当前用户的分享列表，支持分页
// @Tags         题目分享
// @Accept       json
// @Produce      json
// @Param        page        query     int  false  "页码"
// @Param        page_size   query     int  false  "每页数量"
// @Success      200  {object}  response.Response{data=dto.PageResponse{list=[]dto.ShareInfo}}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/shares [get]
// @Security     BearerAuth
func (ctrl *ShareController) ListMyShares(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	list, total, err := question.ListMyShares(userID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, page, pageSize)
}
